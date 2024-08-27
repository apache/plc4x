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

package bacnetip

// BVLPDUTypes is a dictionary of message type values and structs
var BVLPDUTypes map[uint8]func() interface{ Decode(Arg) error }

func init() {
	BVLPDUTypes = map[uint8]func() interface{ Decode(Arg) error }{
		0x00: func() interface{ Decode(Arg) error } {
			v, _ := NewResult()
			return v
		},
		0x01: func() interface{ Decode(Arg) error } {
			v, _ := NewWriteBroadcastDistributionTable()
			return v
		},
		0x02: func() interface{ Decode(Arg) error } {
			v, _ := NewReadBroadcastDistributionTable()
			return v
		},
		0x03: func() interface{ Decode(Arg) error } {
			v, _ := NewReadBroadcastDistributionTableAck()
			return v
		},
		0x04: func() interface{ Decode(Arg) error } {
			v, _ := NewForwardedNPDU(nil)
			return v
		},
		0x05: func() interface{ Decode(Arg) error } {
			v, _ := NewRegisterForeignDevice()
			return v
		},
		0x06: func() interface{ Decode(Arg) error } {
			v, _ := NewReadForeignDeviceTable()
			return v
		},
		0x07: func() interface{ Decode(Arg) error } {
			v, _ := NewReadForeignDeviceTableAck()
			return v
		},
		0x08: func() interface{ Decode(Arg) error } {
			v, _ := NewDeleteForeignDeviceTableEntry()
			return v
		},
		0x09: func() interface{ Decode(Arg) error } {
			v, _ := NewDistributeBroadcastToNetwork(nil)
			return v
		},
		0x0A: func() interface{ Decode(Arg) error } {
			v, _ := NewOriginalUnicastNPDU(nil)
			return v
		},
		0x0B: func() interface{ Decode(Arg) error } {
			v, _ := NewOriginalBroadcastNPDU(nil)
			return v
		},
	}
}
