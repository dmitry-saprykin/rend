// Copyright 2015 Netflix, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orcas

import (
	"github.com/netflix/rend/common"
	"github.com/netflix/rend/handlers"
	"github.com/netflix/rend/metrics"
)

type OrcaConst func(l1, l2 handlers.Handler, res common.Responder) Orca

type Orca interface {
	Set(req common.SetRequest) error
	Add(req common.SetRequest) error
	Replace(req common.SetRequest) error
	Delete(req common.DeleteRequest) error
	Touch(req common.TouchRequest) error
	Get(req common.GetRequest) error
	GetE(req common.GetRequest) error
	Gat(req common.GATRequest) error
	Noop(req common.NoopRequest) error
	Quit(req common.QuitRequest) error
	Version(req common.VersionRequest) error
	Unknown(req common.Request) error
	Error(req common.Request, reqType common.RequestType, err error)
}

var (
	MetricCmdGet         = metrics.AddCounter("cmd_get")
	MetricCmdGetL1       = metrics.AddCounter("cmd_get_l1")
	MetricCmdGetL2       = metrics.AddCounter("cmd_get_l2")
	MetricCmdGetHits     = metrics.AddCounter("cmd_get_hits")
	MetricCmdGetHitsL1   = metrics.AddCounter("cmd_get_hits_l1")
	MetricCmdGetHitsL2   = metrics.AddCounter("cmd_get_hits_l2")
	MetricCmdGetMisses   = metrics.AddCounter("cmd_get_misses")
	MetricCmdGetMissesL1 = metrics.AddCounter("cmd_get_misses_l1")
	MetricCmdGetMissesL2 = metrics.AddCounter("cmd_get_misses_l2")
	MetricCmdGetErrors   = metrics.AddCounter("cmd_get_errors")
	MetricCmdGetErrorsL1 = metrics.AddCounter("cmd_get_errors_l1")
	MetricCmdGetErrorsL2 = metrics.AddCounter("cmd_get_errors_l2")
	MetricCmdGetKeys     = metrics.AddCounter("cmd_get_keys")
	MetricCmdGetKeysL1   = metrics.AddCounter("cmd_get_keys_l1")
	MetricCmdGetKeysL2   = metrics.AddCounter("cmd_get_keys_l2")

	// Batch L1L2 get metrics
	MetricCmdGetSetL1       = metrics.AddCounter("cmd_get_set_l1")
	MetricCmdGetSetErrorsL1 = metrics.AddCounter("cmd_get_set_errors_l1")
	MetricCmdGetSetSucessL1 = metrics.AddCounter("cmd_get_set_success_l1")

	MetricCmdGetE         = metrics.AddCounter("cmd_gete")
	MetricCmdGetEL1       = metrics.AddCounter("cmd_gete_l1")
	MetricCmdGetEL2       = metrics.AddCounter("cmd_gete_l2")
	MetricCmdGetEHits     = metrics.AddCounter("cmd_gete_hits")
	MetricCmdGetEHitsL1   = metrics.AddCounter("cmd_gete_hits_l1")
	MetricCmdGetEHitsL2   = metrics.AddCounter("cmd_gete_hits_l2")
	MetricCmdGetEMisses   = metrics.AddCounter("cmd_gete_misses")
	MetricCmdGetEMissesL1 = metrics.AddCounter("cmd_gete_misses_l1")
	MetricCmdGetEMissesL2 = metrics.AddCounter("cmd_gete_misses_l2")
	MetricCmdGetEErrors   = metrics.AddCounter("cmd_gete_errors")
	MetricCmdGetEErrorsL1 = metrics.AddCounter("cmd_gete_errors_l1")
	MetricCmdGetEErrorsL2 = metrics.AddCounter("cmd_gete_errors_l2")
	MetricCmdGetEKeys     = metrics.AddCounter("cmd_gete_keys")
	MetricCmdGetEKeysL1   = metrics.AddCounter("cmd_gete_keys_l1")
	MetricCmdGetEKeysL2   = metrics.AddCounter("cmd_gete_keys_l2")

	MetricCmdSet          = metrics.AddCounter("cmd_set")
	MetricCmdSetL1        = metrics.AddCounter("cmd_set_l1")
	MetricCmdSetL2        = metrics.AddCounter("cmd_set_l2")
	MetricCmdSetSuccess   = metrics.AddCounter("cmd_set_success")
	MetricCmdSetSuccessL1 = metrics.AddCounter("cmd_set_success_l1")
	MetricCmdSetSuccessL2 = metrics.AddCounter("cmd_set_success_l2")
	MetricCmdSetErrors    = metrics.AddCounter("cmd_set_errors")
	MetricCmdSetErrorsL1  = metrics.AddCounter("cmd_set_errors_l1")
	MetricCmdSetErrorsL2  = metrics.AddCounter("cmd_set_errors_l2")

	// Batch L1L2 set metrics
	MetricCmdSetReplaceL1          = metrics.AddCounter("cmd_set_replace_l1")
	MetricCmdSetReplaceNotStoredL1 = metrics.AddCounter("cmd_set_replace_not_stored_l1")
	MetricCmdSetReplaceErrorsL1    = metrics.AddCounter("cmd_set_replace_errors_l1")
	MetricCmdSetReplaceStoredL1    = metrics.AddCounter("cmd_set_replace_stored_l1")

	MetricCmdAdd            = metrics.AddCounter("cmd_add")
	MetricCmdAddL1          = metrics.AddCounter("cmd_add_l1")
	MetricCmdAddL2          = metrics.AddCounter("cmd_add_l2")
	MetricCmdAddStored      = metrics.AddCounter("cmd_add_stored")
	MetricCmdAddStoredL1    = metrics.AddCounter("cmd_add_stored_l1")
	MetricCmdAddStoredL2    = metrics.AddCounter("cmd_add_stored_l2")
	MetricCmdAddNotStored   = metrics.AddCounter("cmd_add_not_stored")
	MetricCmdAddNotStoredL1 = metrics.AddCounter("cmd_add_not_stored_l1")
	MetricCmdAddNotStoredL2 = metrics.AddCounter("cmd_add_not_stored_l2")
	MetricCmdAddErrors      = metrics.AddCounter("cmd_add_errors")
	MetricCmdAddErrorsL1    = metrics.AddCounter("cmd_add_errors_l1")
	MetricCmdAddErrorsL2    = metrics.AddCounter("cmd_add_errors_l2")

	// Batch L1L2 add metrics
	MetricCmdAddReplaceL1          = metrics.AddCounter("cmd_add_replace_l1")
	MetricCmdAddReplaceNotStoredL1 = metrics.AddCounter("cmd_add_replace_not_stored_l1")
	MetricCmdAddReplaceErrorsL1    = metrics.AddCounter("cmd_add_replace_errors_l1")
	MetricCmdAddReplaceStoredL1    = metrics.AddCounter("cmd_add_replace_stored_l1")

	MetricCmdReplace            = metrics.AddCounter("cmd_replace")
	MetricCmdReplaceL1          = metrics.AddCounter("cmd_replace_l1")
	MetricCmdReplaceL2          = metrics.AddCounter("cmd_replace_l2")
	MetricCmdReplaceStored      = metrics.AddCounter("cmd_replace_stored")
	MetricCmdReplaceStoredL1    = metrics.AddCounter("cmd_replace_stored_l1")
	MetricCmdReplaceStoredL2    = metrics.AddCounter("cmd_replace_stored_l2")
	MetricCmdReplaceNotStored   = metrics.AddCounter("cmd_replace_not_stored")
	MetricCmdReplaceNotStoredL1 = metrics.AddCounter("cmd_replace_not_stored_l1")
	MetricCmdReplaceNotStoredL2 = metrics.AddCounter("cmd_replace_not_stored_l2")
	MetricCmdReplaceErrors      = metrics.AddCounter("cmd_replace_errors")
	MetricCmdReplaceErrorsL1    = metrics.AddCounter("cmd_replace_errors_l1")
	MetricCmdReplaceErrorsL2    = metrics.AddCounter("cmd_replace_errors_l2")

	// Batch L1L2 replace metrics
	MetricCmdReplaceReplaceL1          = metrics.AddCounter("cmd_replace_replace_l1")
	MetricCmdReplaceReplaceNotStoredL1 = metrics.AddCounter("cmd_replace_replace_not_stored_l1")
	MetricCmdReplaceReplaceErrorsL1    = metrics.AddCounter("cmd_replace_replace_errors_l1")
	MetricCmdReplaceReplaceStoredL1    = metrics.AddCounter("cmd_replace_replace_stored_l1")

	MetricCmdDelete         = metrics.AddCounter("cmd_delete")
	MetricCmdDeleteL1       = metrics.AddCounter("cmd_delete_l1")
	MetricCmdDeleteL2       = metrics.AddCounter("cmd_delete_l2")
	MetricCmdDeleteHits     = metrics.AddCounter("cmd_delete_hits")
	MetricCmdDeleteHitsL1   = metrics.AddCounter("cmd_delete_hits_l1")
	MetricCmdDeleteHitsL2   = metrics.AddCounter("cmd_delete_hits_l2")
	MetricCmdDeleteMisses   = metrics.AddCounter("cmd_delete_misses")
	MetricCmdDeleteMissesL1 = metrics.AddCounter("cmd_delete_misses_l1")
	MetricCmdDeleteMissesL2 = metrics.AddCounter("cmd_delete_misses_l2")
	MetricCmdDeleteErrors   = metrics.AddCounter("cmd_delete_errors")
	MetricCmdDeleteErrorsL1 = metrics.AddCounter("cmd_delete_errors_l1")
	MetricCmdDeleteErrorsL2 = metrics.AddCounter("cmd_delete_errors_l2")

	MetricCmdTouch         = metrics.AddCounter("cmd_touch")
	MetricCmdTouchL1       = metrics.AddCounter("cmd_touch_l1")
	MetricCmdTouchL2       = metrics.AddCounter("cmd_touch_l2")
	MetricCmdTouchHits     = metrics.AddCounter("cmd_touch_hits")
	MetricCmdTouchHitsL1   = metrics.AddCounter("cmd_touch_hits_l1")
	MetricCmdTouchHitsL2   = metrics.AddCounter("cmd_touch_hits_l2")
	MetricCmdTouchMisses   = metrics.AddCounter("cmd_touch_misses")
	MetricCmdTouchMissesL1 = metrics.AddCounter("cmd_touch_misses_l1")
	MetricCmdTouchMissesL2 = metrics.AddCounter("cmd_touch_misses_l2")
	MetricCmdTouchErrors   = metrics.AddCounter("cmd_touch_errors")
	MetricCmdTouchErrorsL1 = metrics.AddCounter("cmd_touch_errors_l1")
	MetricCmdTouchErrorsL2 = metrics.AddCounter("cmd_touch_errors_l2")

	// Batch L1L2 touch metrics
	MetricCmdTouchTouchL1       = metrics.AddCounter("cmd_touch_touch_l1")
	MetricCmdTouchTouchMissesL1 = metrics.AddCounter("cmd_touch_touch_misses_l1")
	MetricCmdTouchTouchErrorsL1 = metrics.AddCounter("cmd_touch_touch_errors_l1")
	MetricCmdTouchTouchHitsL1   = metrics.AddCounter("cmd_touch_touch_hits_l1")

	MetricCmdGat         = metrics.AddCounter("cmd_gat")
	MetricCmdGatL1       = metrics.AddCounter("cmd_gat_l1")
	MetricCmdGatL2       = metrics.AddCounter("cmd_gat_l2")
	MetricCmdGatHits     = metrics.AddCounter("cmd_gat_hits")
	MetricCmdGatHitsL1   = metrics.AddCounter("cmd_gat_hits_l1")
	MetricCmdGatHitsL2   = metrics.AddCounter("cmd_gat_hits_l2")
	MetricCmdGatMisses   = metrics.AddCounter("cmd_gat_misses")
	MetricCmdGatMissesL1 = metrics.AddCounter("cmd_gat_misses_l1")
	MetricCmdGatMissesL2 = metrics.AddCounter("cmd_gat_misses_l2")
	MetricCmdGatErrors   = metrics.AddCounter("cmd_gat_errors")
	MetricCmdGatErrorsL1 = metrics.AddCounter("cmd_gat_errors_l1")
	MetricCmdGatErrorsL2 = metrics.AddCounter("cmd_gat_errors_l2")

	// Secondary metrics under GAT that refer to other kinds of operations to
	// backing datastores as a part of the overall request
	MetricCmdGatAddL1          = metrics.AddCounter("cmd_gat_add_l1")
	MetricCmdGatAddErrorsL1    = metrics.AddCounter("cmd_gat_add_errors_l1")
	MetricCmdGatAddStoredL1    = metrics.AddCounter("cmd_gat_add_stored_l1")
	MetricCmdGatAddNotStoredL1 = metrics.AddCounter("cmd_gat_add_not_stored_l1")
	MetricCmdGatSetL2          = metrics.AddCounter("cmd_gat_set_l2")
	MetricCmdGatSetErrorsL2    = metrics.AddCounter("cmd_gat_set_errors_l2")
	MetricCmdGatSetSuccessL2   = metrics.AddCounter("cmd_gat_set_success_l2")

	// Batch L1L2 gat metrics
	MetricCmdGatTouchL1       = metrics.AddCounter("cmd_gat_touch_l1")
	MetricCmdGatTouchMissesL1 = metrics.AddCounter("cmd_gat_touch_misses_l1")
	MetricCmdGatTouchErrorsL1 = metrics.AddCounter("cmd_gat_touch_errors_l1")
	MetricCmdGatTouchHitsL1   = metrics.AddCounter("cmd_gat_touch_hits_l1")

	MetricCmdUnknown = metrics.AddCounter("cmd_unknown")
	MetricCmdNoop    = metrics.AddCounter("cmd_noop")
	MetricCmdQuit    = metrics.AddCounter("cmd_quit")
	MetricCmdVersion = metrics.AddCounter("cmd_version")
)
