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
var NPDUTypes map[uint8]func(args Args, kwArgs KWArgs, options ...Option) (Decoder, error)

func init() {
	NPDUTypes = map[uint8]func(args Args, kwArgs KWArgs, options ...Option) (Decoder, error){
		0x00: convertNPDUType(NewWhoIsRouterToNetwork),
		0x01: convertNPDUType(NewIAmRouterToNetwork),
		0x02: convertNPDUType(NewICouldBeRouterToNetwork),
		0x03: convertNPDUType(NewRejectMessageToNetwork),
		0x04: convertNPDUType(NewRouterBusyToNetwork),
		0x05: convertNPDUType(NewRouterAvailableToNetwork),
		0x06: convertNPDUType(NewInitializeRoutingTable),
		0x07: convertNPDUType(NewInitializeRoutingTableAck),
		0x08: convertNPDUType(NewEstablishConnectionToNetwork),
		0x09: convertNPDUType(NewDisconnectConnectionToNetwork),
		// 0x0A: NewChallengeRequest, // TODO: not present upstream
		// 0x0B: NewSecurityPayload, // TODO: not present upstream
		// 0x0C: NewSecurityResponse, // TODO: not present upstream
		// 0x0D: NewRequestKeyUpdate, // TODO: not present upstream
		// 0x0E: NewUpdateKeyUpdate, // TODO: not present upstream
		// 0x0F: NewUpdateKeyDistributionKey, // TODO: not present upstream
		// 0x10: NewRequestMasterKey, // TODO: not present upstream
		// 0x11: NewSetMasterKey, // TODO: not present upstream
		0x12: convertNPDUType(NewWhatIsNetworkNumber),
		0x13: convertNPDUType(NewNetworkNumberIs),
	}
}

func convertNPDUType[T Decoder](in func(args Args, kwArgs KWArgs, options ...Option) (T, error)) func(args Args, kwArgs KWArgs, options ...Option) (Decoder, error) {
	return func(args Args, kwArgs KWArgs, options ...Option) (Decoder, error) {
		return in(args, kwArgs, options...)
	}
}
