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

package npdu

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

// NPDUTypes is a dictionary of message type values and structs
var NPDUTypes map[uint8]func() Decoder

func init() {
	NPDUTypes = map[uint8]func() Decoder{
		0x00: func() Decoder {
			v, _ := NewWhoIsRouterToNetwork(Nothing())
			return v
		},
		0x01: func() Decoder {
			v, _ := NewIAmRouterToNetwork(Nothing())
			return v
		},
		0x02: func() Decoder {
			v, _ := NewICouldBeRouterToNetwork(Nothing())
			return v
		},
		0x03: func() Decoder {
			v, _ := NewRejectMessageToNetwork(Nothing())
			return v
		},
		0x04: func() Decoder {
			v, _ := NewRouterBusyToNetwork(Nothing())
			return v
		},
		0x05: func() Decoder {
			v, _ := NewRouterAvailableToNetwork(Nothing())
			return v
		},
		0x06: func() Decoder {
			v, _ := NewInitializeRoutingTable(Nothing())
			return v
		},
		0x07: func() Decoder {
			v, _ := NewInitializeRoutingTableAck(Nothing())
			return v
		},
		0x08: func() Decoder {
			v, _ := NewEstablishConnectionToNetwork(Nothing())
			return v
		},
		0x09: func() Decoder {
			v, _ := NewDisconnectConnectionToNetwork(Nothing())
			return v
		},
		// 0x0A: NewChallengeRequest, // TODO: not present upstream
		// 0x0B: NewSecurityPayload, // TODO: not present upstream
		// 0x0C: NewSecurityResponse, // TODO: not present upstream
		// 0x0D: NewRequestKeyUpdate, // TODO: not present upstream
		// 0x0E: NewUpdateKeyUpdate, // TODO: not present upstream
		// 0x0F: NewUpdateKeyDistributionKey, // TODO: not present upstream
		// 0x10: NewRequestMasterKey, // TODO: not present upstream
		// 0x11: NewSetMasterKey, // TODO: not present upstream
		0x12: func() Decoder {
			v, _ := NewWhatIsNetworkNumber(Nothing())
			return v
		},
		0x13: func() Decoder {
			v, _ := NewNetworkNumberIs(Nothing())
			return v
		},
	}
}
