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
	result := make(chan apiModel.PlcReadRequestResult, 1)
	go m.readSync(ctx, readRequest, result)
	return result
}

func (m *Reader) readSync(ctx context.Context, readRequest apiModel.PlcReadRequest, result chan apiModel.PlcReadRequestResult) {
	defer func() {
		if err := recover(); err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v", err))
		}
	}()
	numTags := len(readRequest.GetTagNames())
	if numTags > 20 { // letters g-z
		result <- &spiModel.DefaultPlcReadRequestResult{
			Request: readRequest,
			Err:     errors.New("Only 20 tags can be handled at once"),
		}
		return
	}
	messages := make(map[string]readWriteModel.CBusMessage)
	for _, tagName := range readRequest.GetTagNames() {
		tag := readRequest.GetTag(tagName)
		message, supportsRead, _, _, err := TagToCBusMessage(tag, nil, m.alphaGenerator, m.messageCodec)
		switch {
		case err != nil:
			result <- &spiModel.DefaultPlcReadRequestResult{
				Request: readRequest,
				Err:     errors.Wrapf(err, "Error encoding cbus message for tag %s", tagName),
			}
			return
		case !supportsRead: // Note this should not be reachable
			panic("this should not be possible as we always should then get the error above")
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
		m.createMessageTransactionAndWait(ctx, messageToSend, addResponseCode, tagName, addPlcValue)
	}
	readResponse := spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues)
	result <- &spiModel.DefaultPlcReadRequestResult{
		Request:  readRequest,
		Response: readResponse,
	}
}

func (m *Reader) createMessageTransactionAndWait(ctx context.Context, messageToSend readWriteModel.CBusMessage, addResponseCode func(name string, responseCode apiModel.PlcResponseCode), tagName string, addPlcValue func(name string, plcValue apiValues.PlcValue)) {
	// Start a new request-transaction (Is ended in the response-handler)
	transaction := m.tm.StartTransaction()
	transaction.Submit(func(transaction spi.RequestTransaction) {
		m.sendMessageOverTheWire(ctx, transaction, messageToSend, addResponseCode, tagName, addPlcValue)
	})
	if err := transaction.AwaitCompletion(ctx); err != nil {
		log.Warn().Err(err).Msg("Error while awaiting completion")
	}
}

func (m *Reader) sendMessageOverTheWire(ctx context.Context, transaction spi.RequestTransaction, messageToSend readWriteModel.CBusMessage, addResponseCode func(name string, responseCode apiModel.PlcResponseCode), tagName string, addPlcValue func(name string, plcValue apiValues.PlcValue)) {
	// Send the  over the wire
	log.Trace().Msg("Send ")
	if err := m.messageCodec.SendRequest(ctx, messageToSend, func(cbusMessage spi.Message) bool {
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
		actualAlpha := confirmation.GetConfirmation().GetAlpha().GetCharacter()
		// TODO: assert that this is a CBusMessageToServer indeed (by changing param for example)
		expectedAlpha := messageToSend.(readWriteModel.CBusMessageToServer).GetRequest().(interface{ GetAlpha() readWriteModel.Alpha }).GetAlpha().GetCharacter()
		return actualAlpha == expectedAlpha
	}, func(receivedMessage spi.Message) error {
		defer func(transaction spi.RequestTransaction) {
			// This is just to make sure we don't forget to close the transaction here
			_ = transaction.EndRequest()
		}(transaction)
		// Convert the response into an
		log.Trace().Msg("convert response to ")
		messageToClient := receivedMessage.(readWriteModel.CBusMessageToClient)
		if _, ok := messageToClient.GetReply().(readWriteModel.ServerErrorReplyExactly); ok {
			log.Trace().Msg("We got a server failure")
			addResponseCode(tagName, apiModel.PlcResponseCode_INVALID_DATA)
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
			log.Trace().Msgf("Was no success %s:%v", tagName, responseCode)
			addResponseCode(tagName, responseCode)
			return transaction.EndRequest()
		}

		alpha := replyOrConfirmationConfirmation.GetConfirmation().GetAlpha()
		// TODO: it could be double confirmed but this is not implemented yet
		embeddedReply, ok := replyOrConfirmationConfirmation.GetEmbeddedReply().(readWriteModel.ReplyOrConfirmationReplyExactly)
		if !ok {
			log.Trace().Msgf("Is a confirm only, no data. Alpha: %c", alpha.GetCharacter())
			addResponseCode(tagName, apiModel.PlcResponseCode_NOT_FOUND)
			return transaction.EndRequest()
		}

		log.Trace().Msg("Handling confirmed data")
		// TODO: check if we can use a plcValueSerializer
		encodedReply := embeddedReply.GetReply().(readWriteModel.ReplyEncodedReply).GetEncodedReply()
		if err := MapEncodedReply(transaction, encodedReply, tagName, addResponseCode, addPlcValue); err != nil {
			return errors.Wrap(err, "error encoding reply")
		}
		return transaction.EndRequest()
	}, func(err error) error {
		addResponseCode(tagName, apiModel.PlcResponseCode_REQUEST_TIMEOUT)
		return transaction.FailRequest(err)
	}, time.Second*1); err != nil {
		log.Debug().Err(err).Msgf("Error sending message for tag %s", tagName)
		addResponseCode(tagName, apiModel.PlcResponseCode_INTERNAL_ERROR)
		if err := transaction.FailRequest(errors.Errorf("timeout after %s", time.Second*1)); err != nil {
			log.Debug().Err(err).Msg("Error failing request")
		}
	}
}
