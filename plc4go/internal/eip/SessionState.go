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

package eip

import (
	"context"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
)

// SessionState is the per-connection EIP/CIP session context shared between the
// Connection (which fills it during the handshake) and Reader/Writer (which pick
// the messaging mode and addressing from it).
type SessionState struct {
	sessionHandle          uint32
	senderContext          []byte
	connectionId           uint32
	sequenceCount          atomic.Uint32
	useMessageRouter       bool
	useConnectionManager   bool
	connectionSerialNumber uint16
	routingAddress         []readWriteModel.PathSegment
	connectionPathSize     uint8
}

func NewSessionState(log zerolog.Logger, configuration Configuration) *SessionState {
	routingAddress, connectionPathSize := buildRoutingAddress(log, configuration)
	serialNumber := configuration.connectionSerialNumber
	if serialNumber == 0 {
		serialNumber = uint16(1 + rand.IntN(0xFFFF)) // [1, 0xFFFF], like plc4j
	}
	state := &SessionState{
		senderContext:          []byte(DefaultSenderContext),
		connectionSerialNumber: serialNumber,
		routingAddress:         routingAddress,
		connectionPathSize:     connectionPathSize,
	}
	state.sequenceCount.Store(1)
	return state
}

// nextSequenceCount hands out the connected-messaging sequence number and
// advances it, starting at 1 like plc4j.
func (s *SessionState) nextSequenceCount() uint16 {
	return uint16(s.sequenceCount.Add(1) - 1)
}

func buildRoutingAddress(log zerolog.Logger, configuration Configuration) ([]readWriteModel.PathSegment, uint8) {
	var routingAddress []readWriteModel.PathSegment
	if configuration.communicationPath != "" {
		parts := strings.Split(configuration.communicationPath, ",")
		if len(parts)%2 == 0 {
			for i := 0; i+1 < len(parts); i += 2 {
				switch parts[i] {
				case "1":
					slot, err := strconv.Atoi(parts[i+1])
					if err != nil {
						log.Error().Str("slot", parts[i+1]).Msg("Invalid slot in communication path")
						continue
					}
					routingAddress = append(routingAddress,
						readWriteModel.NewPortSegment(readWriteModel.NewPortSegmentNormal(1, uint8(slot))))
				case "2":
					address := parts[i+1]
					addressLength := uint8(len(address))
					if len(address)%2 != 0 {
						address += "\x00"
					}
					routingAddress = append(routingAddress,
						readWriteModel.NewPortSegment(readWriteModel.NewPortSegmentExtended(2, addressLength, address)))
				default:
					log.Error().Str("port", parts[i]).Msg("Only backplane or Ethernet module routing is supported")
				}
			}
		}
	} else {
		routingAddress = append(routingAddress,
			readWriteModel.NewPortSegment(readWriteModel.NewPortSegmentNormal(1, uint8(configuration.slot))))
	}
	routingAddress = append(routingAddress,
		readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 2)),
		readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1)),
	)
	totalBytes := 0
	for _, segment := range routingAddress {
		totalBytes += int(segment.GetLengthInBytes(context.Background()))
	}
	if totalBytes%2 != 0 {
		totalBytes++
	}
	return routingAddress, uint8(totalBytes / 2)
}
