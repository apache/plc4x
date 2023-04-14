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
	"sync"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Writer struct {
	alphaGenerator *AlphaGenerator
	messageCodec   *MessageCodec
	tm             spi.RequestTransactionManager
}

func NewWriter(tpduGenerator *AlphaGenerator, messageCodec *MessageCodec, tm spi.RequestTransactionManager) Writer {
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
		numTags := len(writeRequest.GetTagNames())
		if numTags > 20 { // letters g-z
			result <- &spiModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("Only 20 tags can be handled at once"),
			}
			return
		}

		messages := make(map[string]readWriteModel.CBusMessage)
		for _, tagName := range writeRequest.GetTagNames() {
			tag := writeRequest.GetTag(tagName)
			plcValue := writeRequest.GetValue(tagName)
			message, _, supportsWrite, _, err := TagToCBusMessage(tag, plcValue, m.alphaGenerator, m.messageCodec)
			if !supportsWrite {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding cbus message for tag %s. Tag is not meant to be written.", tagName),
				}
				return
			}
			if err != nil {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
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
		for tagName, messageToSend := range messages {
			if err := ctx.Err(); err != nil {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request: writeRequest,
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
					// Convert the response into an
					addResponseCode(tagName, apiModel.PlcResponseCode_OK)
					return transaction.EndRequest()
				}, func(err error) error {
					log.Debug().Msgf("Error waiting for tag %s", tagNameCopy)
					addResponseCode(tagNameCopy, apiModel.PlcResponseCode_REQUEST_TIMEOUT)
					// TODO: ok or not ok?
					return transaction.EndRequest()
				}, time.Second*1); err != nil {
					log.Debug().Err(err).Msgf("Error sending message for tag %s", tagNameCopy)
					addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INTERNAL_ERROR)
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
