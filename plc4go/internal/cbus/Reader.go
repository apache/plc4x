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
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult)
	go func() {
		numTags := len(readRequest.GetTagNames())
		if numTags > 20 { // letters g-z
			result <- &spiModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("Only 20 tags can be handled at once"),
			}
			return
		}
		messages := make(map[string]readWriteModel.CBusMessage)
		for _, tagName := range readRequest.GetTagNames() {
			tag := readRequest.GetTag(tagName)
			message, supportsRead, _, _, err := TagToCBusMessage(tag, nil, m.alphaGenerator, m.messageCodec.(*MessageCodec))
			if !supportsRead {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding cbus message for tag %s. Tag is not meant to be read.", tagName),
				}
				return
			}
			if err != nil {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding cbus message for tag %s", tagName),
				}
				return
			}
			messages[tagName] = message
		}
		responseMu := sync.Mutex{}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		addResponseCode := func(name string, responseCode apiModel.PlcResponseCode) {
			responseMu.Lock()
			defer responseMu.Unlock()
			responseCodes[name] = responseCode
		}
		valueMu := sync.Mutex{}
		plcValues := map[string]apiValues.PlcValue{}
		addPlcValue := func(name string, plcValue apiValues.PlcValue) {
			valueMu.Lock()
			defer valueMu.Unlock()
			plcValues[name] = plcValue
		}
		for tagName, messageToSend := range messages {
			if err := ctx.Err(); err != nil {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request: readRequest,
					Err:     err,
				}
				return
			}
			tagNameCopy := tagName
			// Start a new request-transaction (Is ended in the response-handler)
			transaction := m.tm.StartTransaction()
			transaction.Submit(func() {
				// Send the  over the wire
				log.Trace().Msg("Send ")
				if err := m.messageCodec.SendRequest(ctx, messageToSend, func(receivedMessage spi.Message) bool {
					cbusMessage, ok := receivedMessage.(readWriteModel.CBusMessageExactly)
					if !ok {
						return false
					}
					messageToClient, ok := cbusMessage.(readWriteModel.CBusMessageToClientExactly)
					if !ok {
						return false
					}
					// Check if this errored
					if _, ok = messageToClient.GetReply().(readWriteModel.ServerErrorReplyExactly); ok {
						// This means we must handle this below
						return true
					}

					confirmation, ok := messageToClient.GetReply().(readWriteModel.ReplyOrConfirmationConfirmationExactly)
					if !ok {
						return false
					}
					return confirmation.GetConfirmation().GetAlpha().GetCharacter() == messageToSend.(readWriteModel.CBusMessageToServer).GetRequest().(readWriteModel.RequestCommand).GetAlpha().GetCharacter()
				}, func(receivedMessage spi.Message) error {
					defer func(transaction *spi.RequestTransaction) {
						// This is just to make sure we don't forget to close the transaction here
						_ = transaction.EndRequest()
					}(transaction)
					// Convert the response into an
					log.Trace().Msg("convert response to ")
					cbusMessage := receivedMessage.(readWriteModel.CBusMessage)
					messageToClient := cbusMessage.(readWriteModel.CBusMessageToClient)
					if _, ok := messageToClient.GetReply().(readWriteModel.ServerErrorReplyExactly); ok {
						log.Trace().Msg("We got a server failure")
						addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INVALID_DATA)
						return transaction.EndRequest()
					}
					replyOrConfirmationConfirmation := messageToClient.GetReply().(readWriteModel.ReplyOrConfirmationConfirmationExactly)
					if !replyOrConfirmationConfirmation.GetConfirmation().GetIsSuccess() {
						var responseCode apiModel.PlcResponseCode
						switch replyOrConfirmationConfirmation.GetConfirmation().GetConfirmationType() {
						case readWriteModel.ConfirmationType_NOT_TRANSMITTED_TO_MANY_RE_TRANSMISSIONS:
							responseCode = apiModel.PlcResponseCode_REMOTE_ERROR
						case readWriteModel.ConfirmationType_NOT_TRANSMITTED_CORRUPTION:
							responseCode = apiModel.PlcResponseCode_INVALID_DATA
						case readWriteModel.ConfirmationType_NOT_TRANSMITTED_SYNC_LOSS:
							responseCode = apiModel.PlcResponseCode_REMOTE_BUSY
						case readWriteModel.ConfirmationType_NOT_TRANSMITTED_TOO_LONG:
							responseCode = apiModel.PlcResponseCode_INVALID_DATA
						default:
							return transaction.FailRequest(errors.Errorf("Every code should be mapped here: %v", replyOrConfirmationConfirmation.GetConfirmation().GetConfirmationType()))
						}
						log.Trace().Msgf("Was no success %s:%v", tagNameCopy, responseCode)
						addResponseCode(tagNameCopy, responseCode)
						return transaction.EndRequest()
					}

					alpha := replyOrConfirmationConfirmation.GetConfirmation().GetAlpha()
					// TODO: it could be double confirmed but this is not implemented yet
					embeddedReply, ok := replyOrConfirmationConfirmation.GetEmbeddedReply().(readWriteModel.ReplyOrConfirmationReplyExactly)
					if !ok {
						log.Trace().Msgf("Is a confirm only, no data. Alpha: %c", alpha.GetCharacter())
						addResponseCode(tagNameCopy, apiModel.PlcResponseCode_NOT_FOUND)
						return transaction.EndRequest()
					}

					log.Trace().Msg("Handling confirmed data")
					// TODO: check if we can use a plcValueSerializer
					switch reply := embeddedReply.GetReply().(readWriteModel.ReplyEncodedReply).GetEncodedReply().(type) {
					case readWriteModel.EncodedReplyCALReplyExactly:
						calData := reply.GetCalReply().GetCalData()
						addResponseCode(tagNameCopy, apiModel.PlcResponseCode_OK)
						switch calData := calData.(type) {
						case readWriteModel.CALDataStatusExactly:
							application := calData.GetApplication()
							// TODO: verify application... this should be the same
							_ = application
							blockStart := calData.GetBlockStart()
							// TODO: verify application... this should be the same
							_ = blockStart
							statusBytes := calData.GetStatusBytes()
							addResponseCode(tagNameCopy, apiModel.PlcResponseCode_OK)
							plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
							for i, statusByte := range statusBytes {
								plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
								plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
								plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
								plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
							}
							addPlcValue(tagNameCopy, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
								"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
								"blockStart":  spiValues.NewPlcBYTE(blockStart),
								"values":      spiValues.NewPlcList(plcListValues),
							}))
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
								addResponseCode(tagNameCopy, apiModel.PlcResponseCode_OK)
								plcListValues := make([]apiValues.PlcValue, len(statusBytes)*4)
								for i, statusByte := range statusBytes {
									plcListValues[i*4+0] = spiValues.NewPlcSTRING(statusByte.GetGav0().String())
									plcListValues[i*4+1] = spiValues.NewPlcSTRING(statusByte.GetGav1().String())
									plcListValues[i*4+2] = spiValues.NewPlcSTRING(statusByte.GetGav2().String())
									plcListValues[i*4+3] = spiValues.NewPlcSTRING(statusByte.GetGav3().String())
								}
								addPlcValue(tagNameCopy, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
									"application": spiValues.NewPlcSTRING(application.PLC4XEnumName()),
									"blockStart":  spiValues.NewPlcBYTE(blockStart),
									"values":      spiValues.NewPlcList(plcListValues),
								}))
							case readWriteModel.StatusCoding_LEVEL_BY_THIS_SERIAL_INTERFACE:
								fallthrough
							case readWriteModel.StatusCoding_LEVEL_BY_ELSEWHERE:
								levelInformation := calData.GetLevelInformation()
								addResponseCode(tagNameCopy, apiModel.PlcResponseCode_OK)
								plcListValues := make([]apiValues.PlcValue, len(levelInformation))
								for i, levelInformation := range levelInformation {
									switch levelInformation := levelInformation.(type) {
									case readWriteModel.LevelInformationAbsentExactly:
										plcListValues[i] = spiValues.NewPlcSTRING("is absent")
									case readWriteModel.LevelInformationCorruptedExactly:
										plcListValues[i] = spiValues.NewPlcSTRING("corrupted")
									case readWriteModel.LevelInformationNormalExactly:
										plcListValues[i] = spiValues.NewPlcUSINT(levelInformation.GetActualLevel())
									default:
										return transaction.FailRequest(errors.Errorf("Impossible case %v", levelInformation))
									}
								}
								addPlcValue(tagNameCopy, spiValues.NewPlcList(plcListValues))
							}
						case readWriteModel.CALDataIdentifyReplyExactly:
							switch identifyReplyCommand := calData.GetIdentifyReplyCommand().(type) {
							case readWriteModel.IdentifyReplyCommandCurrentSenseLevelsExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetCurrentSenseLevels()))
							case readWriteModel.IdentifyReplyCommandDelaysExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
									"ReStrikeDelay": spiValues.NewPlcUSINT(identifyReplyCommand.GetReStrikeDelay()),
									"TerminalLevel": spiValues.NewPlcRawByteArray(identifyReplyCommand.GetTerminalLevels()),
								}))
							case readWriteModel.IdentifyReplyCommandDSIStatusExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
									"ChannelStatus1":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus1().String()),
									"ChannelStatus2":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus2().String()),
									"ChannelStatus3":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus3().String()),
									"ChannelStatus4":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus4().String()),
									"ChannelStatus5":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus5().String()),
									"ChannelStatus6":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus6().String()),
									"ChannelStatus7":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus7().String()),
									"ChannelStatus8":          spiValues.NewPlcSTRING(identifyReplyCommand.GetChannelStatus8().String()),
									"UnitStatus":              spiValues.NewPlcSTRING(identifyReplyCommand.GetUnitStatus().String()),
									"DimmingUCRevisionNumber": spiValues.NewPlcUSINT(identifyReplyCommand.GetDimmingUCRevisionNumber()),
								}))
							case readWriteModel.IdentifyReplyCommandExtendedDiagnosticSummaryExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
									"LowApplication":         spiValues.NewPlcSTRING(identifyReplyCommand.GetLowApplication().String()),
									"HighApplication":        spiValues.NewPlcSTRING(identifyReplyCommand.GetHighApplication().String()),
									"Area":                   spiValues.NewPlcUSINT(identifyReplyCommand.GetArea()),
									"Crc":                    spiValues.NewPlcUINT(identifyReplyCommand.GetCrc()),
									"SerialNumber":           spiValues.NewPlcUDINT(identifyReplyCommand.GetSerialNumber()),
									"NetworkVoltage":         spiValues.NewPlcUSINT(identifyReplyCommand.GetNetworkVoltage()),
									"UnitInLearnMode":        spiValues.NewPlcBOOL(identifyReplyCommand.GetUnitInLearnMode()),
									"NetworkVoltageLow":      spiValues.NewPlcBOOL(identifyReplyCommand.GetNetworkVoltageLow()),
									"NetworkVoltageMarginal": spiValues.NewPlcBOOL(identifyReplyCommand.GetNetworkVoltageMarginal()),
									"EnableChecksumAlarm":    spiValues.NewPlcBOOL(identifyReplyCommand.GetEnableChecksumAlarm()),
									"OutputUnit":             spiValues.NewPlcBOOL(identifyReplyCommand.GetOutputUnit()),
									"InstallationMMIError":   spiValues.NewPlcBOOL(identifyReplyCommand.GetInstallationMMIError()),
									"EEWriteError":           spiValues.NewPlcBOOL(identifyReplyCommand.GetEEWriteError()),
									"EEChecksumError":        spiValues.NewPlcBOOL(identifyReplyCommand.GetEEChecksumError()),
									"EEDataError":            spiValues.NewPlcBOOL(identifyReplyCommand.GetEEDataError()),
									"MicroReset":             spiValues.NewPlcBOOL(identifyReplyCommand.GetMicroReset()),
									"CommsTxError":           spiValues.NewPlcBOOL(identifyReplyCommand.GetCommsTxError()),
									"InternalStackOverflow":  spiValues.NewPlcBOOL(identifyReplyCommand.GetInternalStackOverflow()),
									"MicroPowerReset":        spiValues.NewPlcBOOL(identifyReplyCommand.GetMicroPowerReset()),
								}))
							case readWriteModel.IdentifyReplyCommandSummaryExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
									"PartName":        spiValues.NewPlcSTRING(identifyReplyCommand.GetPartName()),
									"UnitServiceType": spiValues.NewPlcUSINT(identifyReplyCommand.GetUnitServiceType()),
									"Version":         spiValues.NewPlcSTRING(identifyReplyCommand.GetVersion()),
								}))
							case readWriteModel.IdentifyReplyCommandFirmwareVersionExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcSTRING(identifyReplyCommand.GetFirmwareVersion()))
							case readWriteModel.IdentifyReplyCommandGAVPhysicalAddressesExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetValues()))
							case readWriteModel.IdentifyReplyCommandGAVValuesCurrentExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetValues()))
							case readWriteModel.IdentifyReplyCommandGAVValuesStoredExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetValues()))
							case readWriteModel.IdentifyReplyCommandLogicalAssignmentExactly:
								var plcValues []apiValues.PlcValue
								for _, logicAssigment := range identifyReplyCommand.GetLogicAssigment() {
									plcValues = append(plcValues, spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
										"GreaterOfOrLogic": spiValues.NewPlcBOOL(logicAssigment.GetGreaterOfOrLogic()),
										"ReStrikeDelay":    spiValues.NewPlcBOOL(logicAssigment.GetReStrikeDelay()),
										"AssignedToGav16":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav16()),
										"AssignedToGav15":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav15()),
										"AssignedToGav14":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav14()),
										"AssignedToGav13":  spiValues.NewPlcBOOL(logicAssigment.GetAssignedToGav13()),
									}))
								}
								addPlcValue(tagNameCopy, spiValues.NewPlcList(plcValues))
							case readWriteModel.IdentifyReplyCommandManufacturerExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcSTRING(identifyReplyCommand.GetManufacturerName()))
							case readWriteModel.IdentifyReplyCommandMaximumLevelsExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetMaximumLevels()))
							case readWriteModel.IdentifyReplyCommandMinimumLevelsExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetMinimumLevels()))
							case readWriteModel.IdentifyReplyCommandNetworkTerminalLevelsExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetNetworkTerminalLevels()))
							case readWriteModel.IdentifyReplyCommandNetworkVoltageExactly:
								volts := identifyReplyCommand.GetVolts()
								voltsFloat, err := strconv.ParseFloat(volts, 0)
								if err != nil {
									addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INTERNAL_ERROR)
									return transaction.FailRequest(errors.Wrap(err, "Error parsing volts"))
								}
								voltsDecimalPlace := identifyReplyCommand.GetVoltsDecimalPlace()
								voltsDecimalPlaceFloat, err := strconv.ParseFloat(voltsDecimalPlace, 0)
								if err != nil {
									addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INTERNAL_ERROR)
									return transaction.FailRequest(errors.Wrap(err, "Error parsing volts decimal place"))
								}
								voltsFloat += voltsDecimalPlaceFloat / 10
								addPlcValue(tagNameCopy, spiValues.NewPlcLREAL(voltsFloat))
							case readWriteModel.IdentifyReplyCommandOutputUnitSummaryExactly:
								unitFlags := identifyReplyCommand.GetUnitFlags()
								structContent := map[string]apiValues.PlcValue{
									"UnitFlags": spiValues.NewPlcStruct(map[string]apiValues.PlcValue{
										"AssertingNetworkBurden": spiValues.NewPlcBOOL(unitFlags.GetAssertingNetworkBurden()),
										"RestrikeTimingActive":   spiValues.NewPlcBOOL(unitFlags.GetRestrikeTimingActive()),
										"RemoteOFFInputAsserted": spiValues.NewPlcBOOL(unitFlags.GetRemoteOFFInputAsserted()),
										"RemoteONInputAsserted":  spiValues.NewPlcBOOL(unitFlags.GetRemoteONInputAsserted()),
										"LocalToggleEnabled":     spiValues.NewPlcBOOL(unitFlags.GetLocalToggleEnabled()),
										"LocalToggleActiveState": spiValues.NewPlcBOOL(unitFlags.GetLocalToggleActiveState()),
										"ClockGenerationEnabled": spiValues.NewPlcBOOL(unitFlags.GetClockGenerationEnabled()),
										"UnitGeneratingClock":    spiValues.NewPlcBOOL(unitFlags.GetUnitGeneratingClock()),
									}),
									"TimeFromLastRecoverOfMainsInSeconds": spiValues.NewPlcUSINT(identifyReplyCommand.GetTimeFromLastRecoverOfMainsInSeconds()),
								}
								if gavStoreEnabledByte1 := identifyReplyCommand.GetGavStoreEnabledByte1(); gavStoreEnabledByte1 != nil {
									structContent["GavStoreEnabledByte1"] = spiValues.NewPlcUSINT(*gavStoreEnabledByte1)
								}
								if gavStoreEnabledByte2 := identifyReplyCommand.GetGavStoreEnabledByte2(); gavStoreEnabledByte2 != nil {
									structContent["GavStoreEnabledByte2"] = spiValues.NewPlcUSINT(*gavStoreEnabledByte2)
								}
								addPlcValue(tagNameCopy, spiValues.NewPlcStruct(structContent))
							case readWriteModel.IdentifyReplyCommandTerminalLevelsExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcRawByteArray(identifyReplyCommand.GetTerminalLevels()))
							case readWriteModel.IdentifyReplyCommandTypeExactly:
								addPlcValue(tagNameCopy, spiValues.NewPlcSTRING(identifyReplyCommand.GetUnitType()))
							default:
								addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INVALID_DATA)
								return transaction.FailRequest(errors.Errorf("Unmapped type %T", identifyReplyCommand))
							}
						default:
							wbpcb := spiValues.NewWriteBufferPlcValueBased()
							if err := calData.SerializeWithWriteBuffer(wbpcb); err != nil {
								log.Warn().Err(err).Msgf("Unmapped cal data type %T. Returning raw to string", calData)
								addPlcValue(tagNameCopy, spiValues.NewPlcSTRING(fmt.Sprintf("%s", calData)))
							} else {
								addPlcValue(tagNameCopy, wbpcb.GetPlcValue())
							}
						}
					default:
						return transaction.FailRequest(errors.Errorf("All types should be mapped here. Not mapped: %T", reply))
					}
					return transaction.EndRequest()
				}, func(err error) error {
					addResponseCode(tagNameCopy, apiModel.PlcResponseCode_REQUEST_TIMEOUT)
					return transaction.FailRequest(err)
				}, time.Second*1); err != nil {
					log.Debug().Err(err).Msgf("Error sending message for tag %s", tagNameCopy)
					addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INTERNAL_ERROR)
					_ = transaction.FailRequest(errors.Errorf("timeout after %ss", time.Second*1))
				}
			})
			if err := transaction.AwaitCompletion(); err != nil {
				log.Warn().Err(err).Msg("Error while awaiting completion")
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
