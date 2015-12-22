/**
 * Copyright 2015 Netflix, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/**
 * Local request handlers that perform higher level logic.
 */
package local

import "bufio"
import "bytes"
import "encoding/binary"
import "fmt"
import "io"
import "io/ioutil"
import "math"

import "../binprot"
import "../common"
import "../stream"

// Chunk size, leaving room for the token
// Make sure the value subtracted from chunk size stays in sync
// with the size of the Metadata struct
const CHUNK_SIZE = 1024 - 16
const FULL_DATA_SIZE = 1024

func readResponseHeader(r *bufio.Reader) (binprot.ResponseHeader, error) {
	resHeader, err := binprot.ReadResponseHeader(r)
	if err != nil {
		return binprot.ResponseHeader{}, err
	}

	err = binprot.DecodeError(resHeader)
	if err != nil {
		return resHeader, err
	}

	return resHeader, nil
}

func HandleSet(cmd common.SetRequest, src *bufio.Reader, rw *bufio.ReadWriter) error {
	// For writing chunks, the specialized chunked reader is appropriate.
	// for unchunked, a limited reader will be needed since the text protocol
	// includes a /r/n at the end and there's no EOF to be had with a long-l`ived
	// connection.
	limChunkReader := stream.NewChunkLimitedReader(src, int64(CHUNK_SIZE), int64(cmd.Length))
	numChunks := int(math.Ceil(float64(cmd.Length) / float64(CHUNK_SIZE)))
	token := <-tokens

	metaKey := metaKey(cmd.Key)
	metaData := common.Metadata{
		Length:    cmd.Length,
		OrigFlags: cmd.Flags,
		NumChunks: uint32(numChunks),
		ChunkSize: CHUNK_SIZE,
		Token:     *token,
	}

	metaDataBuf := new(bytes.Buffer)
	binary.Write(metaDataBuf, binary.BigEndian, metaData)

	// Write metadata key
	// TODO: should there be a unique flags value for chunked data?
	localCmd := binprot.SetCmd(metaKey, cmd.Flags, cmd.Exptime, common.METADATA_SIZE)
	err := setLocal(rw.Writer, localCmd, nil, metaDataBuf)
	if err != nil {
		return err
	}

	// Read server's response
	resHeader, err := readResponseHeader(rw.Reader)
	if err != nil {
		// Discard request body
		if _, ioerr := src.Discard(int(cmd.Length)); ioerr != nil {
			return ioerr
		}

		// Discard response body
		if _, ioerr := rw.Discard(int(resHeader.TotalBodyLength)); ioerr != nil {
			return ioerr
		}

		return err
	}

	// Write all the data chunks
	// TODO: Clean up if a data chunk write fails
	// Failure can mean the write failing at the I/O level
	// or at the memcached level, e.g. response == ERROR
	chunkNum := 0
	for limChunkReader.More() {
		// Build this chunk's key
		key := chunkKey(cmd.Key, chunkNum)

		// Write the key
		localCmd = binprot.SetCmd(key, cmd.Flags, cmd.Exptime, FULL_DATA_SIZE)
		err = setLocal(rw.Writer, localCmd, token, limChunkReader)
		if err != nil {
			return err
		}

		// Read server's response
		resHeader, err = readResponseHeader(rw.Reader)
		if err != nil {
			// Discard request body
			for limChunkReader.More() {
				_, ioerr := io.Copy(ioutil.Discard, limChunkReader)

				if ioerr != nil {
					return ioerr
				}

				limChunkReader.NextChunk()
			}

			// Discard repsonse body
			if _, ioerr := rw.Discard(int(resHeader.TotalBodyLength)); ioerr != nil {
				return ioerr
			}

			return err
		}

		// Reset for next iteration
		limChunkReader.NextChunk()
		chunkNum++
	}

	return nil
}

func HandleGet(cmd common.GetRequest, rw *bufio.ReadWriter) (chan common.GetResponse, chan error) {
	// No buffering here so there's not multiple gets in memory
	dataOut := make(chan common.GetResponse)
	errorOut := make(chan error)
	go realHandleGet(cmd, dataOut, errorOut, rw)
	return dataOut, errorOut
}

func realHandleGet(cmd common.GetRequest, dataOut chan common.GetResponse, errorOut chan error, rw *bufio.ReadWriter) {
	// read index
	// make buf
	// for numChunks do
	//   read chunk, append to buffer
	// send response

	defer close(errorOut)
	defer close(dataOut)

outer:
	for idx, key := range cmd.Keys {
		_, metaData, err := getMetadata(rw, key)
		if err != nil {
			// TODO: Better error management
			if err == common.MISS || err == common.ERROR_KEY_NOT_FOUND {
				//fmt.Println("Get miss because of missing metadata. Key:", key)
				dataOut <- common.GetResponse{
					Miss:     true,
					Key:      key,
					Opaque:   cmd.Opaques[idx],
					Quiet:    cmd.Quiet[idx],
					Metadata: metaData,
					Data:     nil,
				}
				continue outer
			}

			errorOut <- err
			return
		}

		// Retrieve all the data from memcached
		dataBuf := make([]byte, metaData.Length)
		tokenBuf := make([]byte, 16)

		for i := 0; i < int(metaData.NumChunks); i++ {
			chunkKey := chunkKey(key, i)

			// indices for slicing, end exclusive
			start, end := chunkSliceIndices(int(metaData.ChunkSize), i, int(metaData.Length))

			// Get the data directly into our buf
			chunkBuf := dataBuf[start:end]
			getCmd := binprot.GetCmd(chunkKey)
			err = getLocalIntoBuf(rw, getCmd, tokenBuf, chunkBuf, int(metaData.ChunkSize))

			if err != nil {
				// TODO: Better error management
				if err == common.MISS || err == common.ERROR_KEY_NOT_FOUND {
					//fmt.Println("Get miss because of missing chunk. Cmd:", getCmd)
					dataOut <- common.GetResponse{
						Miss:     true,
						Key:      key,
						Opaque:   cmd.Opaques[idx],
						Quiet:    cmd.Quiet[idx],
						Metadata: metaData,
						Data:     nil,
					}
					continue outer
				}

				errorOut <- err
				return
			}

			if !bytes.Equal(metaData.Token[:], tokenBuf) {
				fmt.Println("Get miss because of invalid chunk token. Cmd:", getCmd)
				fmt.Printf("Expected: %v\n", metaData.Token)
				fmt.Printf("Got:      %v\n", tokenBuf)
				dataOut <- common.GetResponse{
					Miss:     true,
					Key:      key,
					Opaque:   cmd.Opaques[idx],
					Quiet:    cmd.Quiet[idx],
					Metadata: metaData,
					Data:     nil,
				}
				continue outer
			}
		}

		dataOut <- common.GetResponse{
			Miss:     false,
			Key:      key,
			Opaque:   cmd.Opaques[idx],
			Quiet:    cmd.Quiet[idx],
			Metadata: metaData,
			Data:     dataBuf,
		}
	}
}

func HandleGAT(cmd common.GATRequest, rw *bufio.ReadWriter) (common.GetResponse, error) {
	_, metaData, err := getAndTouchMetadata(rw, cmd.Key, cmd.Exptime)
	if err != nil {
		// TODO: Better error management
		if err == common.MISS || err == common.ERROR_KEY_NOT_FOUND {
			//fmt.Println("Get miss because of missing metadata. Key:", key)
			return common.GetResponse{
				Miss:     true,
				Key:      cmd.Key,
				Opaque:   cmd.Opaque,
				Quiet:    false,
				Metadata: metaData,
				Data:     nil,
			}, nil
		}

		return common.GetResponse{}, err
	}

	// Retrieve all the data from memcached while touching each segment
	dataBuf := make([]byte, metaData.Length)
	tokenBuf := make([]byte, 16)

	for i := 0; i < int(metaData.NumChunks); i++ {
		chunkKey := chunkKey(cmd.Key, i)

		// indices for slicing, end exclusive
		start, end := chunkSliceIndices(int(metaData.ChunkSize), i, int(metaData.Length))

		// Get the data directly into our buf
		chunkBuf := dataBuf[start:end]
		getCmd := binprot.GATCmd(chunkKey, cmd.Exptime)
		err = getLocalIntoBuf(rw, getCmd, tokenBuf, chunkBuf, int(metaData.ChunkSize))

		if err != nil {
			// TODO: Better error management
			if err == common.MISS || err == common.ERROR_KEY_NOT_FOUND {
				//fmt.Println("Get miss because of missing chunk. Cmd:", getCmd)
				return common.GetResponse{
					Miss:     true,
					Key:      cmd.Key,
					Opaque:   cmd.Opaque,
					Quiet:    false,
					Metadata: metaData,
					Data:     nil,
				}, nil
			}

			return common.GetResponse{}, err
		}

		if !bytes.Equal(metaData.Token[:], tokenBuf) {
			fmt.Println("Get miss because of invalid chunk token. Cmd:", getCmd)
			fmt.Printf("Expected: %v\n", metaData.Token)
			fmt.Printf("Got:      %v\n", tokenBuf)

			return common.GetResponse{
				Miss:     true,
				Key:      cmd.Key,
				Opaque:   cmd.Opaque,
				Quiet:    false,
				Metadata: metaData,
				Data:     nil,
			}, nil
		}
	}

	return common.GetResponse{
		Miss:     false,
		Key:      cmd.Key,
		Opaque:   cmd.Opaque,
		Quiet:    false,
		Metadata: metaData,
		Data:     dataBuf,
	}, nil
}

func HandleDelete(cmd common.DeleteRequest, rw *bufio.ReadWriter) error {
	// read metadata
	// delete metadata
	// for 0 to metadata.numChunks
	//  delete item

	metaKey, metaData, err := getMetadata(rw, cmd.Key)

	if err != nil {
		if err == common.MISS {
			fmt.Println("Delete miss because of missing metadata. Key:", cmd.Key)
		}
		return err
	}

	deleteCmd := binprot.DeleteCmd(metaKey)
	err = simpleCmdLocal(rw, deleteCmd)
	if err != nil {
		return err
	}

	for i := 0; i < int(metaData.NumChunks); i++ {
		chunkKey := chunkKey(cmd.Key, i)
		deleteCmd = binprot.DeleteCmd(chunkKey)
		err := simpleCmdLocal(rw, deleteCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func HandleTouch(cmd common.TouchRequest, rw *bufio.ReadWriter) error {
	// read metadata
	// for 0 to metadata.numChunks
	//  touch item
	// touch metadata

	metaKey, metaData, err := getMetadata(rw, cmd.Key)

	if err != nil {
		if err == common.MISS {
			fmt.Println("Touch miss because of missing metadata. Key:", cmd.Key)
			return err
		}

		return err
	}

	for i := 0; i < int(metaData.NumChunks); i++ {
		chunkKey := chunkKey(cmd.Key, i)
		touchCmd := binprot.TouchCmd(chunkKey, cmd.Exptime)
		err := simpleCmdLocal(rw, touchCmd)
		if err != nil {
			return err
		}
	}

	touchCmd := binprot.TouchCmd(metaKey, cmd.Exptime)
	err = simpleCmdLocal(rw, touchCmd)
	if err != nil {
		return err
	}

	return nil
}
