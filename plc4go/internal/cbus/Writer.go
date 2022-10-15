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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Writer struct {
	alphaGenerator *AlphaGenerator
	messageCodec   spi.MessageCodec
	tm             *spi.RequestTransactionManager
}

func NewWriter(tpduGenerator *AlphaGenerator, messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager) Writer {
	return Writer{
		alphaGenerator: tpduGenerator,
		messageCodec:   messageCodec,
		tm:             tm,
	}
}

func (m Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	log.Trace().Msg("Writing")
	result := make(chan apiModel.PlcWriteRequestResult)
	go func() {
		numFields := len(writeRequest.GetFieldNames())
		if numFields > 20 { // letters g-z
			result <- &spiModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("Only 20 fields can be handled at once"),
			}
			return
		}

		messages := make(map[string]readWriteModel.CBusMessage)
		for _, fieldName := range writeRequest.GetFieldNames() {
			field := writeRequest.GetField(fieldName)
			plcValue := writeRequest.GetValue(fieldName)
			message, _, supportsWrite, _, err := FieldToCBusMessage(field, plcValue, m.alphaGenerator, m.messageCodec.(*MessageCodec))
			if !supportsWrite {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding cbus message for field %s. Field is not meant to be written.", fieldName),
				}
				return
			}
			if err != nil {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding cbus message for field %s", fieldName),
				}
				return
			}
			messages[fieldName] = message
		}
		responseMu := sync.Mutex{}
		responseCodes := map[string]apiModel.PlcResponseCode{}
		addResponseCode := func(name string, responseCode apiModel.PlcResponseCode) {
			responseMu.Lock()
			defer responseMu.Unlock()
			responseCodes[name] = responseCode
		}
		for fieldName, messageToSend := range messages {
			if err := ctx.Err(); err != nil {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request: writeRequest,
					Err:     err,
				}
				return
			}
			fieldNameCopy := fieldName
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
					// Convert the response into an
					addResponseCode(fieldName, apiModel.PlcResponseCode_OK)
					return transaction.EndRequest()
				}, func(err error) error {
					log.Debug().Msgf("Error waiting for field %s", fieldNameCopy)
					addResponseCode(fieldNameCopy, apiModel.PlcResponseCode_REQUEST_TIMEOUT)
					// TODO: ok or not ok?
					return transaction.EndRequest()
				}, time.Second*1); err != nil {
					log.Debug().Err(err).Msgf("Error sending message for field %s", fieldNameCopy)
					addResponseCode(fieldNameCopy, apiModel.PlcResponseCode_INTERNAL_ERROR)
					_ = transaction.EndRequest()
				}
			})
		}
		readResponse := spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes)
		result <- &spiModel.DefaultPlcWriteRequestResult{
			Request:  writeRequest,
			Response: readResponse,
		}
	}()
	return result
}
