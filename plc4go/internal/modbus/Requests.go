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

package modbus

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transactions"
)

// maxConcurrentRequests is how many modbus PDUs a connection keeps on the wire at the same time.
//
// One, the same number all three of plc4j's modbus connections return from
// getMaxConcurrentRequests(). It is not a tuning knob this driver leaves open, because the devices
// that need it most are the ones that can say least about themselves: a modbus TCP endpoint is
// very often a gateway that serialises everything it is told onto a slow serial line, and floods
// it under a flood of requests by resetting the connection - which takes every request that was
// in flight with it, so the caller is left with nothing at all rather than with less.
//
// The cost is latency for requests that can't be merged. Reads buy it back by being merged into
// block requests (see ReadOptimizer.go); writes can't be merged - modbus has no function code that
// writes two unrelated addresses - so a multi-tag write is as many round trips as it has tags.
const maxConcurrentRequests = 1

// sendTransacted sends one request through the connection's request transaction manager, so that
// at most maxConcurrentRequests modbus PDUs are on the wire per connection.
//
// It mirrors spi.MessageCodec.SendRequest, which is what it ends up calling: it returns once the
// request has been sent - which is once every request queued ahead of it has been answered -
// handing back whatever the send failed with, and the answer reaches handleMessage or handleError
// exactly as it would without the manager in the middle.
//
// The transaction it takes out is ended on every path out of here, cancellation included. A leaked
// permit is worse than no throttle at all: the connection would simply stop sending, forever.
func sendTransacted(
	ctx context.Context,
	log zerolog.Logger,
	codec spi.MessageCodec,
	tm transactions.RequestTransactionManager,
	name string,
	request spi.Message,
	acceptsMessage spi.AcceptsMessage,
	handleMessage spi.HandleMessage,
	handleError spi.HandleError,
) error {
	transaction := tm.StartTransaction(name)
	if transaction.IsCompleted() {
		// A manager that has been closed hands out a transaction that is completed before it is
		// ever submitted, and submitting one of those does nothing at all - so waiting for it to
		// be sent would wait out the request timeout for a request nobody is going to send. It
		// never took a permit, so there is nothing to end.
		return errors.Errorf("can't send the %s request, the connection is shutting down", name)
	}
	sent := make(chan error, 1)
	transaction.Submit(name+"Operation", func(_ context.Context, transaction transactions.RequestTransaction) {
		// The permit has to be handed back exactly once, on every path through this operation.
		var release sync.Once
		finish := func(sendErr error) {
			release.Do(func() {
				var err error
				if sendErr != nil {
					err = transaction.FailRequest(sendErr)
				} else {
					err = transaction.EndRequest()
				}
				if err != nil {
					log.Debug().Err(err).Str("name", name).Msg("error ending the transaction")
				}
			})
		}
		// A panic in here would otherwise keep the permit forever and wedge the connection for
		// good, which is strictly worse than the flood this throttle exists to stop: the pool
		// worker recovers the panic, so nothing else would ever notice. finish is NOT deferred
		// unconditionally - the transaction has to stay open until the response handler fires,
		// or the permit is released before the answer arrives and the throttle does nothing.
		defer func() {
			if r := recover(); r != nil {
				panicErr := errors.Errorf("panic while sending the %s request: %v", name, r)
				finish(panicErr)
				select {
				case sent <- panicErr:
				default:
				}
				log.Error().Str("name", name).Interface("panic", r).
					Msg("a modbus request panicked; the connection's permit was handed back")
			}
		}()

		err := codec.SendRequest(ctx, name, request, acceptsMessage, func(message spi.Message) error {
			defer finish(nil)
			// The codec answers a handler that returns an error by calling the error handler as
			// well and leaving the expectation registered (spi/default.defaultCodec.HandleMessages),
			// so this one never returns one: the request has been answered, and a second outcome
			// for it would be a second result for the caller.
			if err := handleMessage(message); err != nil {
				log.Debug().Err(err).Str("name", name).Msg("error handling the response")
			}
			return nil
		}, func(err error) error {
			defer finish(nil)
			if err := handleError(err); err != nil {
				log.Debug().Err(err).Str("name", name).Msg("error handling the failure")
			}
			return nil
		})
		if err != nil {
			// Nothing reached the wire and the codec has dropped the expectation again, so no
			// handler is going to end this one.
			finish(err)
		}
		sent <- err
	})
	select {
	case err := <-sent:
		return err
	case <-ctx.Done():
		// The caller ran out of time while this was still waiting its turn behind another
		// request. Ending the transaction is still the operation's to do: either it never runs,
		// in which case it holds no permit, or it runs and the send fails on the expired context,
		// which ends it.
		return errors.Wrapf(ctx.Err(), "the %s request didn't get its turn on the wire", name)
	}
}
