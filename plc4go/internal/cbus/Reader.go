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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Reader struct {
	alphaGenerator *AlphaGenerator
	messageCodec   *MessageCodec
	tm             spi.RequestTransactionManager
}

func NewReader(tpduGenerator *AlphaGenerator, messageCodec *MessageCodec, tm spi.RequestTransactionManager) *Reader {
	return &Reader{
		alphaGenerator: tpduGenerator,
		messageCodec:   messageCodec,
		tm:             tm,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult)
	go m.readSync(ctx, readRequest, result)
	return result
}

func (m *Reader) readSync(ctx context.Context, readRequest apiModel.PlcReadRequest, result chan apiModel.PlcReadRequestResult) {
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
		message, supportsRead, _, _, err := TagToCBusMessage(tag, nil, m.alphaGenerator, m.messageCodec)
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
				defer func(transaction spi.RequestTransaction) {
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
				encodedReply := embeddedReply.GetReply().(readWriteModel.ReplyEncodedReply).GetEncodedReply()
				if err := MapEncodedReply(transaction, encodedReply, tagNameCopy, addResponseCode, addPlcValue); err != nil {
					return errors.Wrap(err, "error encoding reply")
				}
				return transaction.EndRequest()
			}, func(err error) error {
				addResponseCode(tagNameCopy, apiModel.PlcResponseCode_REQUEST_TIMEOUT)
				return transaction.FailRequest(err)
			}, time.Second*1); err != nil {
				log.Debug().Err(err).Msgf("Error sending message for tag %s", tagNameCopy)
				addResponseCode(tagNameCopy, apiModel.PlcResponseCode_INTERNAL_ERROR)
				if err := transaction.FailRequest(errors.Errorf("timeout after %ss", time.Second*1)); err != nil {
					log.Debug().Err(err).Msg("Error failing request")
				}
			}
		})
		if err := transaction.AwaitCompletion(ctx); err != nil {
			log.Warn().Err(err).Msg("Error while awaiting completion")
		}
	}
	readResponse := spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues)
	result <- &spiModel.DefaultPlcReadRequestResult{
		Request:  readRequest,
		Response: readResponse,
	}
}
