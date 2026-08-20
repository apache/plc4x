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

package umas

import (
	"context"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// requester is the one way a UMAS PDU leaves this driver. Every caller - the connect handshake, the
// data dictionary download, the reader, the writer, the ping - goes through it, which is what keeps
// the transaction identifiers, the one-request-at-a-time rule and the request timeout in one place.
type requester struct {
	codec         spi.MessageCodec
	tm            transactions.RequestTransactionManager
	session       *session
	configuration Configuration

	log zerolog.Logger
}

func newRequester(codec spi.MessageCodec, tm transactions.RequestTransactionManager, session *session, configuration Configuration, _options ...options.WithOption) *requester {
	return &requester{
		codec:         codec,
		tm:            tm,
		session:       session,
		configuration: configuration,
		log:           options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

// exchange wraps a UMAS PDU in a Modbus/TCP ADU, sends it and hands back the response's UMAS item.
// The response is matched to the request by transaction identifier - the only thing tying the two
// together, since a UMAS response repeats neither the function key nor anything else of the request.
func (r *requester) exchange(ctx context.Context, name string, item readWriteModel.UmasPDUItem) (readWriteModel.UmasPDUItem, error) {
	transactionIdentifier := r.session.nextTransactionIdentifier()
	request := readWriteModel.NewModbusTcpADU(
		transactionIdentifier, r.configuration.unitIdentifier, readWriteModel.NewUmasPDU(item))

	requestCtx, cancelRequest := context.WithTimeout(ctx, r.configuration.requestTimeout)
	defer cancelRequest()

	message, err := sendTransactedRequestAndWait(requestCtx, r.log, r.codec, r.tm, name, request,
		func(message spi.Message) bool {
			response, ok := message.(readWriteModel.ModbusTcpADU)
			return ok && response.GetTransactionIdentifier() == transactionIdentifier
		})
	if err != nil {
		return nil, errors.Wrapf(err, "error exchanging the %s request", name)
	}
	response, ok := message.(readWriteModel.ModbusTcpADU)
	if !ok {
		return nil, errors.Errorf("%s: expected a ModbusTcpADU, got %T", name, message)
	}
	return extractUmasItem(name, response)
}

// extractUmasItem digs the UMAS item out of a response ADU, turning a Modbus exception into an error.
// Ported from plc4j's UmasConnection.extractUmasResponse.
func extractUmasItem(name string, response readWriteModel.ModbusTcpADU) (readWriteModel.UmasPDUItem, error) {
	switch pdu := response.GetPdu().(type) {
	case readWriteModel.ModbusPDUError:
		return nil, errors.Errorf("%s received the Modbus exception %s", name, pdu.GetExceptionCode())
	case readWriteModel.UmasPDU:
		item := pdu.GetItem()
		if item == nil {
			return nil, errors.Errorf("%s received a UMAS PDU without an item", name)
		}
		return item, nil
	default:
		return nil, errors.Errorf("%s received an unexpected PDU of type %T", name, response.GetPdu())
	}
}

// expectItem narrows a response item to the type the caller asked for, reporting a UMAS error
// response as such: a PLC which refuses a request answers with function key 0xFD rather than with a
// Modbus exception.
func expectItem[T readWriteModel.UmasPDUItem](name string, item readWriteModel.UmasPDUItem) (T, error) {
	if typed, ok := item.(T); ok {
		return typed, nil
	}
	var zero T
	if errorResponse, isError := item.(readWriteModel.UmasPDUErrorResponse); isError {
		return zero, errors.Errorf("%s was refused by the PLC (UMAS error response, %d bytes of detail)",
			name, len(errorResponse.GetBlock()))
	}
	return zero, errors.Errorf("%s received an unexpected response of type %T", name, item)
}

// isUmasErrorResponse says whether the PLC refused the request. The reader and the writer report that
// as a remote error rather than as an internal one.
func isUmasErrorResponse(item readWriteModel.UmasPDUItem) bool {
	_, isError := item.(readWriteModel.UmasPDUErrorResponse)
	return isError
}

// isTimeout says whether an exchange failed because it ran out of time rather than for a reason the
// PLC or the transport gave. The request timeout surfaces as a context deadline (the ctx the exchange
// derives) or as the codec's own expectation timeout, depending on which of the two fires first.
func isTimeout(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, utils.TimeoutError{})
}

// sendTransactedRequestAndWait sends one request through the codec and blocks until the accepted
// response, an error, or ctx cancellation - every outcome resolves, no code path may leave the caller
// waiting forever. It is routed through the request transaction manager, which is what keeps UMAS's one-request-at-a-time promise: the manager is created with a
// concurrency of one, so a request waits for the request before it to finish. plc4j's
// UmasConnection.getMaxConcurrentRequests returns 1 for the same reason.
func sendTransactedRequestAndWait(
	ctx context.Context,
	log zerolog.Logger,
	codec spi.MessageCodec,
	tm transactions.RequestTransactionManager,
	name string,
	request spi.Message,
	accepts func(spi.Message) bool,
) (spi.Message, error) {
	messageChan := make(chan spi.Message, 1)
	errChan := make(chan error, 1)
	transaction := tm.StartTransaction(name)
	transaction.Submit(name+"Operation", func(transactionContext context.Context, transaction transactions.RequestTransaction) {
		ctx, cancel := context.WithCancel(ctx)
		context.AfterFunc(transactionContext, cancel)
		if err := codec.SendRequest(ctx, name, request, accepts, func(message spi.Message) error {
			messageChan <- message
			return transaction.EndRequest()
		}, func(err error) error {
			errChan <- err
			return transaction.EndRequest()
		}); err != nil {
			errChan <- errors.Wrapf(err, "error sending %s request", name)
			if failErr := transaction.FailRequest(err); failErr != nil {
				log.Debug().Err(failErr).Msg("error failing request")
			}
		}
	})
	select {
	case <-ctx.Done():
		return nil, errors.Wrapf(ctx.Err(), "%s aborted", name)
	case err := <-errChan:
		return nil, errors.Wrapf(err, "error waiting for %s response", name)
	case message := <-messageChan:
		return message, nil
	}
}
