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
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
)

type Analyzer struct {
	Client         net.IP
	requestContext model.RequestContext
	cBusOptions    model.CBusOptions
	initialized    bool
}

func (a *Analyzer) PackageParse(packetInformation common.PacketInformation, payload []byte) (interface{}, error) {
	if !a.initialized {
		log.Warn().Msg("Not initialized... doing that now")
		a.requestContext = model.NewRequestContext(false, false, false)
		a.cBusOptions = model.NewCBusOptions(config.CBusConfigInstance.Connect, config.CBusConfigInstance.Smart, config.CBusConfigInstance.Idmon, config.CBusConfigInstance.Exstat, config.CBusConfigInstance.Monitor, config.CBusConfigInstance.Monall, config.CBusConfigInstance.Pun, config.CBusConfigInstance.Pcn, config.CBusConfigInstance.Srchk)
		a.initialized = true
	}
	log.Debug().Msgf("Parsing %s with requestContext\n%v\nBusOptions\n%s", packetInformation, a.requestContext, a.cBusOptions)
	// TODO: srcchk we need to pull that out of the config
	isResponse := packetInformation.DstIp.Equal(a.Client)
	log.Debug().Stringer("packetInformation", packetInformation).Msgf("isResponse: %t", isResponse)
	parse, err := model.CBusMessageParse(utils.NewReadBufferByteBased(payload), isResponse, a.requestContext, a.cBusOptions, uint16(len(payload)))
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing CBusCommand")
	}
	switch cBusMessage := parse.(type) {
	case model.CBusMessageToServerExactly:
		switch request := cBusMessage.GetRequest().(type) {
		case model.RequestDirectCommandAccessExactly:
			sendIdentifyRequestBefore := false
			log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
			switch calDataOrSetParameter := request.GetCalDataOrSetParameter().(type) {
			case model.CALDataOrSetParameterValueExactly:
				switch calDataOrSetParameter.GetCalData().(type) {
				case model.CALDataIdentifyExactly:
					sendIdentifyRequestBefore = true
				}
			}
			a.requestContext = model.NewRequestContext(true, false, sendIdentifyRequestBefore)
		case model.RequestCommandExactly:
			switch command := request.GetCbusCommand().(type) {
			case model.CBusCommandDeviceManagementExactly:
				log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
				a.requestContext = model.NewRequestContext(true, false, false)
			case model.CBusCommandPointToPointExactly:
				log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
				a.requestContext = model.NewRequestContext(true, false, false)
			case model.CBusCommandPointToMultiPointExactly:
				switch command.GetCommand().(type) {
				case model.CBusPointToMultiPointCommandStatusExactly:
					log.Debug().Msgf("No.[%d] SAL status request detected", packetInformation.PacketNumber)
					a.requestContext = model.NewRequestContext(false, true, false)
				}
			case model.CBusCommandPointToPointToMultiPointExactly:
				switch command.GetCommand().(type) {
				case model.CBusPointToPointToMultipointCommandStatusExactly:
					log.Debug().Msgf("No.[%d] SAL status request detected", packetInformation.PacketNumber)
					a.requestContext = model.NewRequestContext(false, true, false)
				}
			}
		case model.RequestObsoleteExactly:
			sendIdentifyRequestBefore := false
			log.Debug().Msgf("No.[%d] CAL request detected", packetInformation.PacketNumber)
			switch calDataOrSetParameter := request.GetCalDataOrSetParameter().(type) {
			case model.CALDataOrSetParameterValueExactly:
				switch calDataOrSetParameter.GetCalData().(type) {
				case model.CALDataIdentifyExactly:
					sendIdentifyRequestBefore = true
				}
			}
			a.requestContext = model.NewRequestContext(true, false, sendIdentifyRequestBefore)
		}
	case model.CBusMessageToClientExactly:
	}
	log.Debug().Msgf("Parsed c-bus command \n%v", parse)
	return parse, nil
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
				fmt.Printf("%v\n", request.GetCalDataOrSetParameter())
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
					case model.ReplyCALReplyExactly:
						// We print this a second time as the first print contains only the hex part
						fmt.Printf("%v\n", reply.GetCalReply())
					}
				}
			case model.ReplyOrConfirmationReplyExactly:
				switch reply := reply.GetReply().(type) {
				case model.ReplyCALReplyExactly:
					// We print this a second time as the first print contains only the hex part
					fmt.Printf("%v\n", reply.GetCalReply())
				}
			}
		}
	}
}
