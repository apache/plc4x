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

package canopen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
)

func TestParseTag(t *testing.T) {
	for address, want := range map[string]Tag{
		"SDO:1:1000/22:UNSIGNED32":      {Kind: TagKindSDO, Node: 1, Index: 1000, SubIndex: 22, DataType: readWriteModel.CANOpenDataType_UNSIGNED32},
		"SDO:1/2:0x1018/0x01:INTEGER8":  {Kind: TagKindSDO, Node: 1, AnswerNode: 2, Index: 0x1018, SubIndex: 1, DataType: readWriteModel.CANOpenDataType_INTEGER8},
		"SDO:127:0xFFFF/255:REAL64":     {Kind: TagKindSDO, Node: 127, Index: 0xFFFF, SubIndex: 255, DataType: readWriteModel.CANOpenDataType_REAL64},
		"TRANSMIT_PDO_3:5:UNSIGNED16":   {Kind: TagKindPDO, Node: 5, Service: readWriteModel.CANOpenService_TRANSMIT_PDO_3, DataType: readWriteModel.CANOpenDataType_UNSIGNED16},
		"RECEIVE_PDO_1:9:REAL32":        {Kind: TagKindPDO, Node: 9, Service: readWriteModel.CANOpenService_RECEIVE_PDO_1, DataType: readWriteModel.CANOpenDataType_REAL32},
		"NMT":                           {Kind: TagKindNMT},
		"NMT:12":                        {Kind: TagKindNMT, Node: 12},
		"HEARTBEAT":                     {Kind: TagKindHeartbeat},
		" HEARTBEAT:3 ":                 {Kind: TagKindHeartbeat, Node: 3},
		"SDO:3:0x1008/0:VISIBLE_STRING": {Kind: TagKindSDO, Node: 3, Index: 0x1008, DataType: readWriteModel.CANOpenDataType_VISIBLE_STRING},
	} {
		t.Run(address, func(t *testing.T) {
			parsed, err := NewTagHandler().ParseTag(address)
			require.NoError(t, err)
			assert.Equal(t, want, parsed)
			// GetAddressString has to parse back into the same tag: subscriptions depend on it.
			again, err := NewTagHandler().ParseTag(parsed.GetAddressString())
			require.NoError(t, err)
			assert.Equal(t, want, again)
		})
	}
}

func TestParseTagRejects(t *testing.T) {
	for _, address := range []string{
		"",
		"SDO:0:1000/1:UNSIGNED8",     // node 0 has no SDO server
		"SDO:128:1000/1:UNSIGNED8",   // node ids stop at 127
		"SDO:1:0x10000/1:UNSIGNED8",  // index is 16 bit
		"SDO:1:1000/256:UNSIGNED8",   // subindex is 8 bit
		"SDO:1:1000/1:NOPE",          // unknown type
		"SDO:1:1000/1:TIME_OF_DAY",   // no data item for it
		"SDO:1:1000/1:UNSIGNED8[2]",  // arrays are not supported
		"TRANSMIT_PDO_5:1:UNSIGNED8", // there are four PDOs
		"TRANSMIT_PDO_1:0:UNSIGNED8", // a PDO belongs to a node
		"NMT:128",
		"HEARTBEAT:x",
	} {
		_, err := NewTagHandler().ParseTag(address)
		assert.Error(t, err, address)
	}
}

func TestValueTypes(t *testing.T) {
	for address, want := range map[string]string{
		"SDO:1:1/0:UNSIGNED24":     "UDINT",
		"SDO:1:1/0:INTEGER56":      "LINT",
		"SDO:1:1/0:VISIBLE_STRING": "STRING",
		"SDO:1:1/0:RECORD":         "RAW_BYTE_ARRAY",
		"SDO:1:1/0:BOOLEAN":        "BOOL",
		"TRANSMIT_PDO_1:1:REAL32":  "REAL",
		"NMT":                      "USINT",
		"HEARTBEAT":                "Struct",
	} {
		tag, err := NewTagHandler().ParseTag(address)
		require.NoError(t, err)
		assert.Equal(t, want, tag.GetValueType().String(), address)
	}
}
