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

package cbusanalyzer

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/google/gopacket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"reflect"
)

type Analyzer struct {
	Client                          net.IP
	requestContext                  model.RequestContext
	cBusOptions                     model.CBusOptions
	initialized                     bool
	currentInboundPayloads          map[string][]byte
	currentPrefilterInboundPayloads map[string][]byte
	mappedPacketChan                chan gopacket.Packet
}

func (a *Analyzer) Init() {
	if a.initialized {
		return
	}
	a.requestContext = model.NewRequestContext(false, false, false)
	a.cBusOptions = model.NewCBusOptions(config.CBusConfigInstance.Connect, config.CBusConfigInstance.Smart, config.CBusConfigInstance.Idmon, config.CBusConfigInstance.Exstat, config.CBusConfigInstance.Monitor, config.CBusConfigInstance.Monall, config.CBusConfigInstance.Pun, config.CBusConfigInstance.Pcn, config.CBusConfigInstance.Srchk)
	a.currentInboundPayloads = make(map[string][]byte)
	a.currentPrefilterInboundPayloads = make(map[string][]byte)
	a.initialized = true
}

func (a *Analyzer) PackageParse(packetInformation common.PacketInformation, payload []byte) (interface{}, error) {
	if !a.initialized {
		log.Warn().Msg("Not initialized... doing that now")
		a.Init()
	}
	cBusOptions := a.cBusOptions
	log.Debug().Msgf("Parsing %s with requestContext\n%v\nBusOptions\n%s\npayload:%+q", packetInformation, a.requestContext, cBusOptions, payload)
	isResponse := a.isResponse(packetInformation)
	if isResponse {
		// Responses should have a checksum
		cBusOptions = model.NewCBusOptions(
			cBusOptions.GetConnect(),
			cBusOptions.GetSmart(),
			cBusOptions.GetIdmon(),
			cBusOptions.GetExstat(),
			cBusOptions.GetMonitor(),
			cBusOptions.GetMonall(),
			cBusOptions.GetPun(),
			cBusOptions.GetPcn(),
			true,
		)
	}
	currentPayload, err := a.getCurrentPayload(packetInformation, payload, isResponse, true, false)
	if err != nil {
		return nil, err
	}
	// TODO: apparently we only do crc on receive with our tests so we need to implement that
	parse, err := model.CBusMessageParse(utils.NewReadBufferByteBased(currentPayload), isResponse, a.requestContext, cBusOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing CBusCommand")
	}
	switch cBusMessage := parse.(type) {
	case model.CBusMessageToServerExactly:
		switch request := cBusMessage.GetRequest().(type) {
		case model.RequestDirectCommandAccessExactly:
			sendIdentifyRequestBefore := false
			log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
			switch request.GetCalData().(type) {
			case model.CALDataIdentifyExactly:
				sendIdentifyRequestBefore = true
			}
			a.requestContext = model.NewRequestContext(true, false, sendIdentifyRequestBefore)
		case model.RequestCommandExactly:
			switch command := request.GetCbusCommand().(type) {
			case model.CBusCommandDeviceManagementExactly:
				log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
				a.requestContext = model.NewRequestContext(true, false, false)
			case model.CBusCommandPointToPointExactly:
				sendIdentifyRequestBefore := false
				log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
				switch command.GetCommand().GetCalData().(type) {
				case model.CALDataIdentifyExactly:
					sendIdentifyRequestBefore = true
				}
				a.requestContext = model.NewRequestContext(true, false, sendIdentifyRequestBefore)
			case model.CBusCommandPointToMultiPointExactly:
				switch command.GetCommand().(type) {
				case model.CBusPointToMultiPointCommandStatusExactly:
					log.Debug().Msgf("No.[%d] SAL status request detected", packetInformation.PacketNumber)
					a.requestContext = model.NewRequestContext(false, true, false)
				}
			case model.CBusCommandPointToPointToMultiPointExactly:
				switch command.GetCommand().(type) {
				case model.CBusPointToPointToMultiPointCommandStatusExactly:
					log.Debug().Msgf("No.[%d] SAL status request detected", packetInformation.PacketNumber)
					a.requestContext = model.NewRequestContext(false, true, false)
				}
			}
		case model.RequestObsoleteExactly:
			sendIdentifyRequestBefore := false
			log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
			switch request.GetCalData().(type) {
			case model.CALDataIdentifyExactly:
				sendIdentifyRequestBefore = true
			}
			a.requestContext = model.NewRequestContext(true, false, sendIdentifyRequestBefore)
		}
	case model.CBusMessageToClientExactly:
		// We received a request so we need to reset our flags
		a.requestContext = model.NewRequestContext(false, false, false)
	}
	log.Debug().Msgf("Parsed c-bus command \n%v", parse)
	return parse, nil
}

func (a *Analyzer) isResponse(packetInformation common.PacketInformation) bool {
	isResponse := packetInformation.DstIp.Equal(a.Client)
	log.Debug().Stringer("packetInformation", packetInformation).Msgf("isResponse: %t", isResponse)
	return isResponse
}

func (a *Analyzer) getCurrentPayload(packetInformation common.PacketInformation, payload []byte, isResponse bool, warnAboutSplit bool, usePrefilterPayloads bool) ([]byte, error) {
	srcUip := packetInformation.SrcIp.String()
	payload = filterXOnXOff(payload)
	if len(payload) == 0 {
		return nil, common.ErrEmptyPackage
	}
	// Check if we have a termination in the middle
	currentInboundPayloads := a.currentInboundPayloads
	if usePrefilterPayloads {
		currentInboundPayloads = a.currentPrefilterInboundPayloads
	}
	currentPayload := currentInboundPayloads[srcUip]
	currentPayload = append(currentPayload, payload...)
	shouldClearInboundPayload := true
	isMergedMessage := false
	// Check if we have a merged message
mergeCheck:
	for i, b := range currentPayload {
		if i == 0 {
			// TODO: we ignore the first byte as this is typical for reset etc... so maybe this is good or bad we will see
			continue
		}
		switch b {
		case 0x0D:
			if i+1 < len(currentPayload) && currentPayload[i+1] == 0x0A {
				// If we know the next is a newline we jump to that index...
				i++
			}
			// ... other than that the logic is the same
			fallthrough
		case 0x0A:
			// We have a merged message if we are not at the end
			if i < len(currentPayload)-1 {
				event := log.Warn()
				if !warnAboutSplit {
					event = log.Debug()
				}
				event.Stringer("packetInformation", packetInformation).Msgf("we have a split at index %d (usePrefilterPayloads=%t)", i, usePrefilterPayloads)
				// In this case we need to put the tail into our "buffer"
				currentInboundPayloads[srcUip] = currentPayload[i+1:]
				// and use the beginning as current payload
				currentPayload = currentPayload[:i+1]
				shouldClearInboundPayload = false
				isMergedMessage = true
				break mergeCheck
			}
		}
	}
	if !isMergedMessage {
		// When we have a merge message we already set the current payload to the tail
		currentInboundPayloads[srcUip] = currentPayload
	} else {
		log.Debug().Stringer("packetInformation", packetInformation).Msgf("Remainder %+q", currentInboundPayloads[srcUip])
	}
	if lastElement := currentPayload[len(currentPayload)-1]; (!isResponse /*a request must end with cr*/ && lastElement != 0x0D /*cr*/) || (isResponse /*a response must end with lf*/ && lastElement != 0x0A /*lf*/) {
		return nil, common.ErrUnterminatedPackage
	} else {
		log.Debug().Msgf("Last element 0x%x", lastElement)
		if shouldClearInboundPayload {
			if currentSavedPayload := currentInboundPayloads[srcUip]; currentSavedPayload != nil {
				// We remove our current payload from the beginning of the cache
				for i, b := range currentPayload {
					if currentSavedPayload[i] != b {
						panic("programming error... at this point they should start with the identical bytes")
					}
				}
			}
			currentInboundPayloads[srcUip] = nil
		}
	}
	log.Debug().Stringer("packetInformation", packetInformation).Msgf("Returning payload %+q", currentPayload)
	return currentPayload, nil
}

func filterXOnXOff(payload []byte) []byte {
	n := 0
	for _, b := range payload {
		switch b {
		case 0x11: // Filter XON
		case 0x13: // Filter XOFF
		default:
			payload[n] = b
			n++
		}
	}
	return payload[:n]
}

func (a *Analyzer) SerializePackage(message interface{}) ([]byte, error) {
	if message, ok := message.(model.CBusMessage); !ok {
		log.Fatal().Msgf("Unsupported type %T supplied", message)
		panic("unreachable statement")
	} else {
		based := utils.NewWriteBufferByteBased()
		if err := message.Serialize(based); err != nil {
			return nil, errors.Wrap(err, "Error serializing")
		}
		return based.GetBytes(), nil
	}
}

// PrettyPrint as messages switch encoding (hex strings) we print the inner messages too
func (a *Analyzer) PrettyPrint(message interface{}) {
	if message, ok := message.(model.CBusMessage); !ok {
		log.Fatal().Msgf("Unsupported type %T supplied", message)
		panic("unreachable statement")
	} else {
		fmt.Printf("%v\n", message)
		switch message := message.(type) {
		case model.CBusMessageToServerExactly:
			switch request := message.GetRequest().(type) {
			case model.RequestDirectCommandAccessExactly:
				fmt.Printf("%v\n", request.GetCalData())
			case model.RequestObsoleteExactly:
				fmt.Printf("%v\n", request.GetCalData())
			case model.RequestCommandExactly:
				fmt.Printf("%v\n", request.GetCbusCommand())
			}
		case model.CBusMessageToClientExactly:
			switch reply := message.GetReply().(type) {
			case model.ReplyOrConfirmationConfirmationExactly:
				switch reply := reply.GetEmbeddedReply().(type) {
				// TODO: add recursion
				case model.ReplyOrConfirmationReplyExactly:
					switch reply := reply.GetReply().(type) {
					case model.ReplyEncodedReplyExactly:
						switch reply := reply.GetEncodedReply().(type) {
						case model.EncodedReplyExtendedFormatStatusReplyExactly:
							// We print this a second time as the first print contains only the hex part
							fmt.Printf("%v\n", reply.GetReply())
						case model.EncodedReplyStandardFormatStatusReplyExactly:
							// We print this a second time as the first print contains only the hex part
							fmt.Printf("%v\n", reply.GetReply())
						case model.EncodedReplyCALReplyExactly:
							// We print this a second time as the first print contains only the hex part
							fmt.Printf("%v\n", reply.GetCalReply())
						case model.MonitoredSALReplyExactly:
							// We print this a second time as the first print contains only the hex part
							fmt.Printf("%v\n", reply.GetMonitoredSAL())
						}
					}
				}
			case model.ReplyOrConfirmationReplyExactly:
				switch reply := reply.GetReply().(type) {
				case model.ReplyEncodedReplyExactly:
					switch reply := reply.GetEncodedReply().(type) {
					case model.EncodedReplyExtendedFormatStatusReplyExactly:
						// We print this a second time as the first print contains only the hex part
						fmt.Printf("%v\n", reply.GetReply())
					case model.EncodedReplyStandardFormatStatusReplyExactly:
						// We print this a second time as the first print contains only the hex part
						fmt.Printf("%v\n", reply.GetReply())
					case model.EncodedReplyCALReplyExactly:
						// We print this a second time as the first print contains only the hex part
						fmt.Printf("%v\n", reply.GetCalReply())
					}
				}
			}
		}
	}
}

// MapPackets reorders the packages as they were not split
func (a *Analyzer) MapPackets(in chan gopacket.Packet, packetInformationCreator func(packet gopacket.Packet) common.PacketInformation) chan gopacket.Packet {
	if a.mappedPacketChan == nil {
		a.mappedPacketChan = make(chan gopacket.Packet)
		go func() {
			defer close(a.mappedPacketChan)
		mappingLoop:
			for packet := range in {
				switch {
				case packet == nil:
					log.Debug().Msg("Done reading packages. (nil returned)")
					a.mappedPacketChan <- nil
					break mappingLoop
				case packet.ApplicationLayer() == nil:
					a.mappedPacketChan <- packet
				default:
					packetInformation := packetInformationCreator(packet)
					if payload, err := a.getCurrentPayload(packetInformation, packet.ApplicationLayer().Payload(), a.isResponse(packetInformation), false, true); err != nil {
						log.Debug().Err(err).Stringer("packetInformation", packetInformation).Msg("Filtering message")
						a.mappedPacketChan <- common.NewFilteredPackage(err, packet)
					} else {
						currentApplicationLayer := packet.ApplicationLayer()
						newPayload := gopacket.Payload(payload)
						if !reflect.DeepEqual(currentApplicationLayer.Payload(), newPayload) {
							log.Debug().Msgf("Replacing payload %q with %q", currentApplicationLayer.Payload(), payload)
							packet = &manipulatedPackage{Packet: packet, newApplicationLayer: newPayload}
						}
						a.mappedPacketChan <- packet
					}
				}
			}
		}()
	}
	return a.mappedPacketChan
}

// ByteOutput returns the string representation as usually this is ASCII over serial... so this output is much more useful in that context
func (a *Analyzer) ByteOutput(data []byte) string {
	return fmt.Sprintf("%+q\n", data)
}

type manipulatedPackage struct {
	gopacket.Packet
	newApplicationLayer gopacket.ApplicationLayer
}

func (p *manipulatedPackage) SetApplicationLayer(l gopacket.ApplicationLayer) {
	p.newApplicationLayer = l
}

func (p *manipulatedPackage) ApplicationLayer() gopacket.ApplicationLayer {
	return p.newApplicationLayer
}
