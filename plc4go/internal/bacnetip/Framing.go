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

import (
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// wrapAPDU encapsulates an APDU in the BACnet/IP NPDU + BVLC layers expected
// by MessageCodec.Send (which type-asserts to model.BVLC). Without this
// wrapping the codec panics on the cast and the request silently dies.
//
// expectingReply is set for confirmed requests; unconfirmed requests pass
// false. The NPDU is intentionally local-only (no DNET/SNET) because routed
// addressing happens at the Tag layer in Phase 6.
func wrapAPDU(apdu model.APDU, expectingReply bool) model.BVLC {
	control := model.NewNPDUControl(
		false, // messageTypeFieldPresent
		false, // destinationSpecified
		false, // sourceSpecified
		expectingReply,
		model.NPDUNetworkPriority_NORMAL_MESSAGE,
	)
	npdu := model.NewNPDU(
		1, // protocolVersionNumber
		control,
		nil, // destinationNetworkAddress
		nil, // destinationLength
		nil, // destinationAddress
		nil, // sourceNetworkAddress
		nil, // sourceLength
		nil, // sourceAddress
		nil, // hopCount
		nil, // nlm
		apdu,
	)
	return model.NewBVLCOriginalUnicastNPDU(npdu)
}
