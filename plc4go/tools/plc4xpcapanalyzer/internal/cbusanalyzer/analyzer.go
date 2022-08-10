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
	"github.com/apache/plc4x/plc4go/internal/cbus"
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/config"
	"github.com/apache/plc4x/plc4go/tools/plc4xpcapanalyzer/internal/common"
	"github.com/google/gopacket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
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

	lastParsePayload []byte
	lastMapPayload   []byte
}

func (a *Analyzer) Init() {
	if a.initialized {
		return
	}
	a.requestContext = model.NewRequestContext(false)
	a.cBusOptions = model.NewCBusOptions(config.CBusConfigInstance.Connect, config.CBusConfigInstance.Smart, config.CBusConfigInstance.Idmon, config.CBusConfigInstance.Exstat, config.CBusConfigInstance.Monitor, config.CBusConfigInstance.Monall, config.CBusConfigInstance.Pun, config.CBusConfigInstance.Pcn, config.CBusConfigInstance.Srchk)
	a.currentInboundPayloads = make(map[string][]byte)
	a.currentPrefilterInboundPayloads = make(map[string][]byte)
	a.initialized = true
}

func (a *Analyzer) PackageParse(packetInformation common.PacketInformation, payload []byte) (spi.Message, error) {
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
	mergeCallback := func(index int) {
		log.Warn().Stringer("packetInformation", packetInformation).Msgf("we have a split at index %d", index)
	}
	currentPayload, err := a.getCurrentPayload(packetInformation, payload, mergeCallback, a.currentInboundPayloads, &a.lastParsePayload)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(currentPayload, a.lastParsePayload) {
		return nil, common.ErrEcho
	}
	a.lastParsePayload = currentPayload
	parse, err := model.CBusMessageParse(utils.NewReadBufferByteBased(currentPayload), isResponse, a.requestContext, cBusOptions)
	if err != nil {
		if secondParse, err := model.CBusMessageParse(utils.NewReadBufferByteBased(currentPayload), isResponse, model.NewRequestContext(false), model.NewCBusOptions(false, false, false, false, false, false, false, false, false)); err != nil {
			log.Debug().Err(err).Msg("Second parse failed too")
			return nil, errors.Wrap(err, "Error parsing CBusCommand")
		} else {
			log.Warn().Stringer("packetInformation", packetInformation).Msgf("package got overridden by second parse... probably a MMI\n%s", secondParse)
			parse = secondParse
		}
	}
	a.requestContext = cbus.CreateRequestContextWithInfoCallback(parse, func(infoString string) {
		log.Debug().Msgf("No.[%d] %s", packetInformation.PacketNumber, infoString)
	})
	log.Debug().Msgf("Parsed c-bus command \n%v", parse)
	return parse, nil
}

func (a *Analyzer) isResponse(packetInformation common.PacketInformation) bool {
	isResponse := packetInformation.DstIp.Equal(a.Client)
	log.Debug().Stringer("packetInformation", packetInformation).Msgf("isResponse: %t", isResponse)
	return isResponse
}

func (a *Analyzer) getCurrentPayload(packetInformation common.PacketInformation, payload []byte, mergeCallback func(int), currentInboundPayloads map[string][]byte, lastPayload *[]byte) ([]byte, error) {
	srcUip := packetInformation.SrcIp.String()
	payload = filterXOnXOff(payload)
	if len(payload) == 0 {
		return nil, common.ErrEmptyPackage
	}
	currentPayload := currentInboundPayloads[srcUip]
	if currentPayload != nil {
		log.Debug().Func(func(e *zerolog.Event) {
			e.Msgf("Prepending current payload %+q to actual payload %+q: %+q", currentPayload, payload, append(currentPayload, payload...))
		})
		currentPayload = append(currentPayload, payload...)
	} else {
		currentPayload = payload
	}
	if len(currentPayload) == 1 && currentPayload[0] == '!' {
		// This is an errormessage from the server
		return currentPayload, nil
	}
	containsError := false
	// We ensure that there are no random ! in the string
	currentPayload, containsError = filterOneServerError(currentPayload)
	if containsError {
		// Save the current inbound payload for the next try
		currentInboundPayloads[srcUip] = currentPayload
		return []byte{'!'}, nil
	}
	// Check if we have a termination in the middle
	isMergedMessage, shouldClearInboundPayload := mergeCheck(&currentPayload, srcUip, mergeCallback, currentInboundPayloads, lastPayload)
	if !isMergedMessage {
		// When we have a merge message we already set the current payload to the tail
		currentInboundPayloads[srcUip] = currentPayload
	} else {
		log.Debug().Stringer("packetInformation", packetInformation).Msgf("Remainder %+q", currentInboundPayloads[srcUip])
	}
	if lastElement := currentPayload[len(currentPayload)-1]; (lastElement != '\r') && (lastElement != '\n') {
		return nil, common.ErrUnterminatedPackage
	} else {
		log.Debug().Stringer("packetInformation", packetInformation).Msgf("Last element 0x%x", lastElement)
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

func mergeCheck(currentPayload *[]byte, srcUip string, mergeCallback func(int), currentInboundPayloads map[string][]byte, lastPayload *[]byte) (isMergedMessage, shouldClearInboundPayload bool) {
	// Check if we have a merged message
	for i, b := range *currentPayload {
		if i == 0 {
			// we ignore the first byte as this is typical for reset etc... so maybe this is good or bad we will see
			continue
		}
		switch b {
		case 0x0D:
			if i+1 < len(*currentPayload) && (*currentPayload)[i+1] == 0x0A {
				// If we know the next is a newline we jump to that index...
				i++
			}
			// ... other than that the logic is the same
			fallthrough
		case 0x0A:
			// We have a merged message if we are not at the end
			if i < len(*currentPayload)-1 {
				headPayload := (*currentPayload)[:i+1]
				tailPayload := (*currentPayload)[i+1:]
				if reflect.DeepEqual(headPayload, *lastPayload) {
					// This means that we have a merge where the last payload is an echo. In that case we discard that here to not offset all numbers
					*currentPayload = tailPayload
					log.Debug().Msgf("We cut the echo message %s out of the response to keep numbering", headPayload)
					return mergeCheck(currentPayload, srcUip, mergeCallback, currentInboundPayloads, lastPayload)
				} else {
					if mergeCallback != nil {
						mergeCallback(i)
					}
					// In this case we need to put the tail into our "buffer"
					currentInboundPayloads[srcUip] = tailPayload
					// and use the beginning as current payload
					*currentPayload = headPayload
					return true, false
				}
			}
		}
	}
	return false, true
}

func filterXOnXOff(payload []byte) []byte {
	n := 0
	for i, b := range payload {
		switch b {
		case 0x11: // Filter XON
			fallthrough
		case 0x13: // Filter XOFF
			log.Trace().Msgf("Filtering %x at %d for %+q", b, i, payload)
		default:
			payload[n] = b
			n++
		}
	}
	return payload[:n]
}

func filterOneServerError(unfilteredPayload []byte) (filteredPayload []byte, containsError bool) {
	for i, b := range unfilteredPayload {
		if b == '!' {
			return append(unfilteredPayload[:i], unfilteredPayload[i+1:]...), true

		}
	}
	return unfilteredPayload, false
}

func (a *Analyzer) SerializePackage(message spi.Message) ([]byte, error) {
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
					mergeCallback := func(index int) {
						log.Warn().Stringer("packetInformation", packetInformation).Msgf("we have a split at index %d", index)
					}
					if payload, err := a.getCurrentPayload(packetInformation, packet.ApplicationLayer().Payload(), mergeCallback, a.currentPrefilterInboundPayloads, &a.lastMapPayload); err != nil {
						log.Debug().Err(err).Stringer("packetInformation", packetInformation).Msg("Filtering message")
						a.mappedPacketChan <- common.NewFilteredPackage(err, packet)
					} else {
						currentApplicationLayer := packet.ApplicationLayer()
						newPayload := gopacket.Payload(payload)
						if !reflect.DeepEqual(currentApplicationLayer.Payload(), payload) {
							log.Debug().Msgf("Replacing payload %q with %q", currentApplicationLayer.Payload(), payload)
							packet = &manipulatedPackage{Packet: packet, newApplicationLayer: newPayload}
						}
						a.lastMapPayload = payload
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
