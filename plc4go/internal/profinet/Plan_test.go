/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package profinet

import (
	"encoding/hex"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/profinet/gsdml"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/profinet/readwrite/model"
)

// capturedSuite is the protocol's parser/serializer suite, which holds a real controller's Connect
// request to a SIMOCODE pro V PN (case 39) and the device's response (case 40).
const capturedSuite = "../../../protocols/profinet/src/test/resources/protocols/profinet/ParserSerializerTestsuite.xml"

func capturedConnect(t *testing.T, name string) readWriteModel.DceRpc_Packet {
	t.Helper()
	text, err := os.ReadFile(capturedSuite)
	if err != nil {
		t.Skipf("%s is not checked out", capturedSuite)
	}
	m := regexp.MustCompile(`<name>` + regexp.QuoteMeta(name) + `[^<]*</name>\s*<raw>([0-9a-f]+)</raw>`).FindSubmatch(text)
	require.NotNil(t, m, "case %s", name)
	data, err := hex.DecodeString(string(m[1]))
	require.NoError(t, err)
	packet, err := readWriteModel.DceRpc_PacketParse(t.Context(), data)
	require.NoError(t, err)
	return packet
}

// The frame layout a real controller chose for this device has to come out of our planner, given
// the same submodules and lengths.
func TestPlanMatchesCapturedConnect(t *testing.T) {
	request := capturedConnect(t, "39\tConnect request").GetPayload().(readWriteModel.PnIoCm_Packet_Req)
	var reported []realSubmodule
	layouts := map[[2]uint16]gsdml.Layout{}
	var crs []readWriteModel.PnIoCm_Block_IoCrReq
	for _, block := range request.GetBlocks() {
		switch b := block.(type) {
		case readWriteModel.PnIoCm_Block_ExpectedSubmoduleReq:
			for _, api := range b.GetApis() {
				for _, submodule := range api.GetSubmodules() {
					key := [2]uint16{api.GetSlotNumber(), submodule.GetSlotNumber()}
					reported = append(reported, realSubmodule{Slot: key[0], Subslot: key[1], ModuleIdent: api.GetModuleIdentNumber(), SubmoduleIdent: submodule.GetSubmoduleIdentNumber()})
					layout := gsdml.Layout{IOPSLength: 1, IOCSLength: 1}
					switch s := submodule.(type) {
					case readWriteModel.PnIoCm_Submodule_InputData:
						layout.InputLength = int(s.GetInputSubmoduleDataLength())
					case readWriteModel.PnIoCm_Submodule_OutputData:
						layout.OutputLength = int(s.GetOutputSubmoduleDataLength())
					case readWriteModel.PnIoCm_Submodule_InputAndOutputData:
						layout.InputLength = int(s.GetInputSubmoduleDataLength())
						layout.OutputLength = int(s.GetOutputSubmoduleDataLength())
					}
					layouts[key] = layout
				}
			}
		case readWriteModel.PnIoCm_Block_IoCrReq:
			crs = append(crs, b)
		}
	}
	require.Len(t, crs, 2)
	require.NotEmpty(t, reported)

	p, err := planFrom(reported, func(real realSubmodule) (gsdml.Resolved, gsdml.Layout, error) {
		return gsdml.Resolved{Submodule: &gsdml.SubmoduleItem{}}, layouts[[2]uint16{real.Slot, real.Subslot}], nil
	})
	require.NoError(t, err)

	for _, cr := range crs {
		input := cr.GetIoCrType() == readWriteModel.PnIoCm_IoCrType_INPUT_CR
		ours := p.ioCrApi(input)
		theirs := cr.GetApis()[0]
		assert.Equal(t, describeObjects(theirs.GetIoDataObjects()), describeObjects(ours.GetIoDataObjects()), "data objects, input=%v", input)
		assert.Equal(t, describeStatuses(theirs.GetIoCss()), describeStatuses(ours.GetIoCss()), "consumer statuses, input=%v", input)
		length := p.OutputLength
		if input {
			length = p.InputLength
		}
		assert.Equal(t, int(cr.GetDataLength()), length, "C_SDU length, input=%v", input)
	}

	// The expected submodule blocks we would send describe the same modules.
	assert.Len(t, p.expectedSubmoduleBlocks(), countSlots(reported))
}

func describeObjects(objects []readWriteModel.PnIoCm_IoDataObject) [][3]uint16 {
	var out [][3]uint16
	for _, o := range objects {
		out = append(out, [3]uint16{o.GetSlotNumber(), o.GetSubSlotNumber(), o.GetIoDataObjectFrameOffset()})
	}
	return out
}

func describeStatuses(statuses []readWriteModel.PnIoCm_IoCs) [][3]uint16 {
	var out [][3]uint16
	for _, s := range statuses {
		out = append(out, [3]uint16{s.GetSlotNumber(), s.GetSubSlotNumber(), s.GetIoFrameOffset()})
	}
	return out
}

func countSlots(reported []realSubmodule) int {
	slots := map[uint16]bool{}
	for _, r := range reported {
		slots[r.Slot] = true
	}
	return len(slots)
}

// The device's captured Connect response parses into the frame ids our connect reads.
func TestCapturedConnectResponse(t *testing.T) {
	response := capturedConnect(t, "40\tConnect response").GetPayload().(readWriteModel.PnIoCm_Packet_Res)
	assert.Zero(t, response.GetErrorCode())
	frameIDs := map[readWriteModel.PnIoCm_IoCrType]uint16{}
	for _, block := range response.GetBlocks() {
		if cr, ok := block.(readWriteModel.PnIoCm_Block_IoCrRes); ok {
			frameIDs[cr.GetIoCrType()] = cr.GetFrameId()
		}
	}
	assert.NotZero(t, frameIDs[readWriteModel.PnIoCm_IoCrType_INPUT_CR])
	assert.NotZero(t, frameIDs[readWriteModel.PnIoCm_IoCrType_OUTPUT_CR])
}
