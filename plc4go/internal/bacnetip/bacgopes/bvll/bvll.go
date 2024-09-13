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

package bvll

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

var _debug = CreateDebugPrinter()

const (
	BVLCIResult                          = uint8(0x00)
	BVLCIWriteBroadcastDistributionTable = uint8(0x01)

	BVLCIReadBroadcastDistributionTable    = uint8(0x02)
	BVLCIReadBroadcastDistributionTableAck = uint8(0x03)
	BVLCIForwardedNPDU                     = uint8(0x04)
	BVLCIRegisterForeignDevice             = uint8(0x05)
	BVLCIReadForeignDeviceTable            = uint8(0x06)
	BVLCIReadForeignDeviceTableAck         = uint8(0x07)
	BVLCIDeleteForeignDeviceTableEntry     = uint8(0x08)
	BVLCIDistributeBroadcastToNetwork      = uint8(0x09)
	BVLCIOriginalUnicastNPDU               = uint8(0x0A)
	BVLCIOriginalBroadcastNPDU             = uint8(0x0B)
)

// BVLPDUTypes is a dictionary of message type values and structs
var BVLPDUTypes map[uint8]func() Decoder

func init() {
	BVLPDUTypes = map[uint8]func() Decoder{
		BVLCIResult: func() Decoder {
			v, _ := NewResult(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIWriteBroadcastDistributionTable: func() Decoder {
			v, _ := NewWriteBroadcastDistributionTable(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIReadBroadcastDistributionTable: func() Decoder {
			v, _ := NewReadBroadcastDistributionTable(Nothing())
			return v
		},
		BVLCIReadBroadcastDistributionTableAck: func() Decoder {
			v, _ := NewReadBroadcastDistributionTableAck(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIForwardedNPDU: func() Decoder {
			v, _ := NewForwardedNPDU(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIRegisterForeignDevice: func() Decoder {
			v, _ := NewRegisterForeignDevice(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIReadForeignDeviceTable: func() Decoder {
			v, _ := NewReadForeignDeviceTable(Nothing())
			return v
		},
		BVLCIReadForeignDeviceTableAck: func() Decoder {
			v, _ := NewReadForeignDeviceTableAck(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIDeleteForeignDeviceTableEntry: func() Decoder {
			v, _ := NewDeleteForeignDeviceTableEntry(nil, NoArgs, NoKWArgs())
			return v
		},
		BVLCIDistributeBroadcastToNetwork: func() Decoder {
			v, _ := NewDistributeBroadcastToNetwork(Nothing())
			return v
		},
		BVLCIOriginalUnicastNPDU: func() Decoder {
			v, _ := NewOriginalUnicastNPDU(Nothing())
			return v
		},
		BVLCIOriginalBroadcastNPDU: func() Decoder {
			v, _ := NewOriginalBroadcastNPDU(Nothing())
			return v
		},
	}
}
