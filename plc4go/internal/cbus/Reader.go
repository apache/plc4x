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

package cbus

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi"
	spiModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/internal/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Reader struct {
	alphaGenerator *AlphaGenerator
	messageCodec   spi.MessageCodec
	tm             *spi.RequestTransactionManager
}

func NewReader(tpduGenerator *AlphaGenerator, messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager) *Reader {
	return &Reader{
		alphaGenerator: tpduGenerator,
		messageCodec:   messageCodec,
		tm:             tm,
	}
}

func (m *Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {
		numFields := len(readRequest.GetFieldNames())
		if numFields > 20 {
			result <- &spiModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("Only 20 fields can be handled at once"),
			}
			return
		}
		messages := make(map[string]readWriteModel.CBusMessage)
		for _, fieldName := range readRequest.GetFieldNames() {
			field := readRequest.GetField(fieldName)
			message, err := m.fieldToCBusMessage(field)
			if err != nil {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding cbus message for field %s", fieldName),
				}
				return
			}
			messages[fieldName] = message
		}
		responseMu := sync.Mutex{}
		responseCodes := map[string]model.PlcResponseCode{}
		addResponseCode := func(name string, responseCode model.PlcResponseCode) {
			responseMu.Lock()
			defer responseMu.Unlock()
			responseCodes[name] = responseCode
		}
		valueMu := sync.Mutex{}
		plcValues := map[string]values.PlcValue{}
		addPlcValue := func(name string, plcValue values.PlcValue) {
			valueMu.Lock()
			defer valueMu.Unlock()
			plcValues[name] = plcValue
		}
		for fieldName, messageToSend := range messages {
			fieldNameCopy := fieldName
			// Start a new request-transaction (Is ended in the response-handler)
			requestWasOk := make(chan bool)
			transaction := m.tm.StartTransaction()
			transaction.Submit(func() {

				// Send the  over the wire
				log.Trace().Msg("Send ")
				if err := m.messageCodec.SendRequest(
					messageToSend,
					func(receivedMessage spi.Message) bool {
						cbusMessage, ok := receivedMessage.(readWriteModel.CBusMessageExactly)
						if !ok {
							return false
						}
						messageToClient, ok := cbusMessage.(readWriteModel.CBusMessageToClientExactly)
						if !ok {
							return false
						}
						confirmation, ok := messageToClient.GetReply().(readWriteModel.ReplyOrConfirmationConfirmationExactly)
						if !ok {
							return false
						}
						return confirmation.GetConfirmation().GetAlpha().GetCharacter() == messageToSend.(readWriteModel.CBusMessageToServer).GetRequest().(readWriteModel.RequestCommand).GetAlpha().GetCharacter()
					},
					func(receivedMessage spi.Message) error {
						// Convert the response into an
						log.Trace().Msg("convert response to ")
						cbusMessage := receivedMessage.(readWriteModel.CBusMessage)
						messageToClient := cbusMessage.(readWriteModel.CBusMessageToClient)
						confirmation := messageToClient.GetReply().(readWriteModel.ReplyOrConfirmationConfirmationExactly)
						if !confirmation.GetConfirmation().GetIsSuccess() {
							var responseCode model.PlcResponseCode
							switch confirmation.GetConfirmation().GetConfirmationType() {
							case readWriteModel.ConfirmationType_CONFIRMATION_SUCCESSFUL:
								responseCode = model.PlcResponseCode_OK
							case readWriteModel.ConfirmationType_NOT_TRANSMITTED_TO_MANY_RE_TRANSMISSIONS:
								responseCode = model.PlcResponseCode_REMOTE_ERROR
							case readWriteModel.ConfirmationType_NOT_TRANSMITTED_CORRUPTION:
								responseCode = model.PlcResponseCode_INVALID_DATA
							case readWriteModel.ConfirmationType_NOT_TRANSMITTED_SYNC_LOSS:
								responseCode = model.PlcResponseCode_REMOTE_BUSY
							case readWriteModel.ConfirmationType_NOT_TRANSMITTED_TOO_LONG:
								responseCode = model.PlcResponseCode_INVALID_DATA
							}
							addResponseCode(fieldNameCopy, responseCode)
							return nil
						}

						// TODO: it could be double confirmed but this is not implemented yet
						embeddedReply := confirmation.GetEmbeddedReply().(readWriteModel.ReplyOrConfirmationReplyExactly)

						switch reply := embeddedReply.GetReply().(readWriteModel.ReplyEncodedReply).GetEncodedReply().(type) {
						case readWriteModel.EncodedReplyCALReplyExactly:
							calData := reply.GetCalReply().GetCalData()
							addResponseCode(fieldNameCopy, model.PlcResponseCode_OK)
							switch calData := calData.(type) {
							case readWriteModel.CALDataStatusExactly:
								application := calData.GetApplication()
								// TODO: verify application... this should be the same
								_ = application
								blockStart := calData.GetBlockStart()
								// TODO: verify application... this should be the same
								_ = blockStart
								statusBytes := calData.GetStatusBytes()
								addResponseCode(fieldNameCopy, model.PlcResponseCode_OK)
								plcListValues := make([]values.PlcValue, len(statusBytes)*4)
								for i, statusByte := range statusBytes {
									plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
									plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
									plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
									plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
								}
								addPlcValue(fieldNameCopy, spiValues.NewPlcList(plcListValues))
							case readWriteModel.CALDataStatusExtendedExactly:
								coding := calData.GetCoding()
								// TODO: verify coding... this should be the same
								_ = coding
								application := calData.GetApplication()
								// TODO: verify application... this should be the same
								_ = application
								blockStart := calData.GetBlockStart()
								// TODO: verify application... this should be the same
								_ = blockStart
								switch coding {
								case readWriteModel.StatusCoding_BINARY_BY_THIS_SERIAL_INTERFACE:
									fallthrough
								case readWriteModel.StatusCoding_BINARY_BY_ELSEWHERE:
									statusBytes := calData.GetStatusBytes()
									addResponseCode(fieldNameCopy, model.PlcResponseCode_OK)
									plcListValues := make([]values.PlcValue, len(statusBytes)*4)
									for i, statusByte := range statusBytes {
										plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
										plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
										plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
										plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
									}
									addPlcValue(fieldNameCopy, spiValues.NewPlcList(plcListValues))
								case readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE:
									fallthrough
								case readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE:
									levelInformation := calData.GetLevelInformation()
									addResponseCode(fieldNameCopy, model.PlcResponseCode_OK)
									plcListValues := make([]values.PlcValue, len(levelInformation))
									for i, levelInformation := range levelInformation {
										switch levelInformation := levelInformation.(type) {
										case readWriteModel.LevelInformationAbsentExactly:
											plcListValues[i] = spiValues.NewPlcSTRING("is absent")
										case readWriteModel.LevelInformationCorruptedExactly:
											plcListValues[i] = spiValues.NewPlcSTRING("corrupted")
										case readWriteModel.LevelInformationNormalExactly:
											plcListValues[i] = spiValues.NewPlcUSINT(levelInformation.GetActualLevel())
										default:
											panic("Impossible case")
										}
									}
									addPlcValue(fieldNameCopy, spiValues.NewPlcList(plcListValues))
								}
							}
							// TODO: how should we serialize that???
							addPlcValue(fieldNameCopy, spiValues.NewPlcSTRING(fmt.Sprintf("%s", calData)))
						}
						requestWasOk <- true
						return transaction.EndRequest()
					},
					func(err error) error {
						log.Debug().Msgf("Error waiting for field %s", fieldNameCopy)
						addResponseCode(fieldNameCopy, model.PlcResponseCode_REQUEST_TIMEOUT)
						// TODO: ok or not ok?
						requestWasOk <- true
						return transaction.EndRequest()
					},
					time.Second*1); err != nil {
					log.Debug().Err(err).Msgf("Error sending message for field %s", fieldNameCopy)
					addResponseCode(fieldNameCopy, model.PlcResponseCode_INTERNAL_ERROR)
					_ = transaction.EndRequest()
					requestWasOk <- false
				}
			})
			if !<-requestWasOk {
				// TODO: if we found a error we can abort
				break
			}
		}
		readResponse := spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues)
		result <- &spiModel.DefaultPlcReadRequestResult{
			Request:  readRequest,
			Response: readResponse,
		}
	}()
	return result
}

var defaultRequestContext = readWriteModel.NewRequestContext(false)
var defaultOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)

func (m *Reader) fieldToCBusMessage(field model.PlcField) (readWriteModel.CBusMessage, error) {
	switch field := field.(type) {
	case *statusField:
		var statusRequest readWriteModel.StatusRequest
		switch field.statusRequestType {
		case StatusRequestTypeBinaryState:
			statusRequest = readWriteModel.NewStatusRequestBinaryState(field.application, 0x7A)
		case StatusRequestTypeLevel:
			statusRequest = readWriteModel.NewStatusRequestLevel(field.application, *field.startingGroupAddressLabel, 0x73)
		}
		command := readWriteModel.NewCBusPointToMultiPointCommandStatus(statusRequest, byte(field.application), defaultOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToMultiPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToMultiPoint(command, header, defaultOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(m.alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), defaultOptions)
		return readWriteModel.NewCBusMessageToServer(request, defaultRequestContext, defaultOptions), nil
	case *calRecallField:
		calData := readWriteModel.NewCALDataRecall(field.parameter, field.count, readWriteModel.CALCommandTypeContainer_CALCommandRecall, nil, defaultRequestContext)
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToPointCommandDirect(field.unitAddress, 0x0000, calData, defaultOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToPoint(command, header, defaultOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(m.alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), defaultOptions)
		return readWriteModel.NewCBusMessageToServer(request, defaultRequestContext, defaultOptions), nil
	case *calIdentifyField:
		calData := readWriteModel.NewCALDataIdentify(field.attribute, readWriteModel.CALCommandTypeContainer_CALCommandIdentify, nil, defaultRequestContext)
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToPointCommandDirect(field.unitAddress, 0x0000, calData, defaultOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToPoint(command, header, defaultOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(m.alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), defaultOptions)
		return readWriteModel.NewCBusMessageToServer(request, defaultRequestContext, defaultOptions), nil
	case *calGetstatusField:
		calData := readWriteModel.NewCALDataGetStatus(field.parameter, field.count, readWriteModel.CALCommandTypeContainer_CALCommandGetStatus, nil, defaultRequestContext)
		//TODO: we need support for bridged commands
		command := readWriteModel.NewCBusPointToPointCommandDirect(field.unitAddress, 0x0000, calData, defaultOptions)
		header := readWriteModel.NewCBusHeader(readWriteModel.PriorityClass_Class4, false, 0, readWriteModel.DestinationAddressType_PointToPoint)
		cbusCommand := readWriteModel.NewCBusCommandPointToPoint(command, header, defaultOptions)
		request := readWriteModel.NewRequestCommand(cbusCommand, nil, readWriteModel.NewAlpha(m.alphaGenerator.getAndIncrement()), readWriteModel.RequestType_REQUEST_COMMAND, nil, nil, readWriteModel.RequestType_EMPTY, readWriteModel.NewRequestTermination(), defaultOptions)
		return readWriteModel.NewCBusMessageToServer(request, defaultRequestContext, defaultOptions), nil
	default:
		return nil, errors.Errorf("Unmapped type %T", field)
	}
}
