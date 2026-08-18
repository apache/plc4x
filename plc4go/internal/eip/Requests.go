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

package eip

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// sendRequestAndWait sends one request through the codec and blocks until the
// accepted response, an error, or ctx cancellation. Every outcome resolves -
// no code path may leave the caller waiting forever.
func sendRequestAndWait(
	ctx context.Context,
	log zerolog.Logger,
	codec spi.MessageCodec,
	name string,
	request spi.Message,
	accepts func(spi.Message) bool,
) (spi.Message, error) {
	messageChan := make(chan spi.Message, 1)
	errChan := make(chan error, 1)
	if err := codec.SendRequest(ctx, name, request, accepts, func(message spi.Message) error {
		messageChan <- message
		return nil
	}, func(err error) error {
		errChan <- err
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "error sending %s request", name)
	}
	select {
	case <-ctx.Done():
		return nil, errors.Wrapf(ctx.Err(), "%s aborted", name)
	case err := <-errChan:
		return nil, errors.Wrapf(err, "error waiting for %s response", name)
	case message := <-messageChan:
		return message, nil
	}
}

// sendTransactedRequestAndWait is sendRequestAndWait routed through the
// request transaction manager, for the read/write paths.
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
