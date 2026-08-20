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

package slmp

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

const (
	// commandBatchRead is Batch Read, the only read command this version sends (SH-080008 chapter 7).
	commandBatchRead = uint16(0x0401)
	// commandBatchWrite is Batch Write, the only write command this version sends.
	commandBatchWrite = uint16(0x1401)
	// subCommandWordUnits selects word units for both of them. Bit units (0x0001) are out of scope,
	// exactly as in plc4j's first version.
	subCommandWordUnits = uint16(0x0000)
	// endCodeNormalCompletion is the end code of a response that completed normally; every other
	// value is an error whose meaning is in SH-080008's error-code list.
	endCodeNormalCompletion = uint16(0x0000)
)

// formatEndCode spells an end code the way the SLMP manual spells it, so a log line can be looked up
// in the manual's error-code list without converting anything by hand.
func formatEndCode(endCode uint16) string {
	return fmt.Sprintf("0x%04X", endCode)
}

// acceptsAnyResponseFrame is the only correlation a 3E frame allows: it carries no transaction or
// sequence id whatsoever, so the first response frame to arrive is taken as the answer to the
// request that is in flight.
//
// That is sound only because the request transaction manager is built with a concurrency of one
// (maxConcurrentRequests), so exactly one request is ever in flight. plc4j relies on the same
// property with its single pendingResponse slot and getMaxConcurrentRequests() == 1.
//
// The caveat it inherits from the wire format: if a request times out, a late response for it may
// arrive while the *next* request is in flight and would then be taken as that next request's
// answer. Nothing in 3E can prevent this, so a request that timed out has to be treated as
// unreliable by the caller - and a timed-out write in particular may still have been applied by the
// device, so retrying it can apply it twice.
func acceptsAnyResponseFrame(message spi.Message) bool {
	_, ok := message.(readWriteModel.SlmpResponseFrame3E)
	return ok
}

// sendTransactedRequestAndWait sends one 3E request through the transaction manager and blocks until
// the response, an error, or ctx cancellation. Every outcome resolves - no code path may leave the
// caller waiting forever.
func sendTransactedRequestAndWait(
	ctx context.Context,
	log zerolog.Logger,
	codec spi.MessageCodec,
	tm transactions.RequestTransactionManager,
	name string,
	request readWriteModel.SlmpRequestFrame3E,
) (readWriteModel.SlmpResponseFrame3E, error) {
	messageChan := make(chan spi.Message, 1)
	errChan := make(chan error, 1)
	transaction := tm.StartTransaction(name)
	transaction.Submit(name+"Operation", func(transactionContext context.Context, transaction transactions.RequestTransaction) {
		ctx, cancel := context.WithCancel(ctx)
		context.AfterFunc(transactionContext, cancel)
		if err := codec.SendRequest(ctx, name, request, acceptsAnyResponseFrame, func(message spi.Message) error {
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
		response, ok := message.(readWriteModel.SlmpResponseFrame3E)
		if !ok {
			// Unreachable while acceptsAnyResponseFrame is the accepts predicate, but a wrong type
			// here has to be an error rather than a nil the caller dereferences.
			return nil, errors.Errorf("expected an SlmpResponseFrame3E, got %T", message)
		}
		return response, nil
	}
}

// newRequestFrame wraps one piece of request data in a 3E request frame.
func newRequestFrame(configuration Configuration, command uint16, requestData readWriteModel.SlmpRequestData) readWriteModel.SlmpRequestFrame3E {
	return readWriteModel.NewSlmpRequestFrame3E(configuration.monitoringTimer, command, subCommandWordUnits, requestData)
}
