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

package _default

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// A fatal transport error stops the codec from inside its receive worker, which then returns through its restart
// defer. That defer used to start a new receive worker whenever running was true, and running is true again as
// soon as a Connect follows the fatal stop: the stopped run's worker then started a second receive worker next to
// the one of the new run. Both read the transport, and in a driver like cbus a later fatal error in one of them
// closed the channels the other one was sending to.
//
// The test holds the stopped run's worker between the fatal stop and its restart defer (on its "transport error
// requested worker shutdown" log line, the last thing it does before it returns), connects again, lets the worker
// return and counts the receive workers reading the transport.
func Test_defaultCodec_Connect_afterFatalStop_startsOneReceiveWorker(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var receives, reading, maxReading atomic.Int32
		requirements := NewMockDefaultCodecRequirements(t)
		requirements.EXPECT().Receive(mock.Anything).RunAndReturn(func(ctx context.Context) (spi.Message, error) {
			if receives.Add(1) == 1 {
				return nil, errors.New("read failed on purpose")
			}
			n := reading.Add(1)
			defer reading.Add(-1)
			for {
				if seen := maxReading.Load(); n <= seen || maxReading.CompareAndSwap(seen, n) {
					break
				}
			}
			<-ctx.Done()
			return nil, ctx.Err()
		})
		instance := NewMockTransportInstance(t)
		instance.EXPECT().IsConnected().Return(true)
		instance.EXPECT().Close().Return(nil)
		instance.EXPECT().ClassifyError(mock.Anything).Return(transports.TransportErrorFatal)

		fatalStopped := make(chan struct{})
		release := make(chan struct{})
		var holdOnce sync.Once
		logger := zerolog.New(io.Discard).Level(zerolog.TraceLevel).Hook(zerolog.HookFunc(func(_ *zerolog.Event, _ zerolog.Level, msg string) {
			if msg == "transport error requested worker shutdown" {
				holdOnce.Do(func() {
					close(fatalStopped)
					<-release
				})
			}
		}))
		codec := buildDefaultCodec(requirements, instance,
			options.WithCustomLogger(logger),
			options.WithTraceDefaultMessageCodecWorker(true),
			// Receive also without expectations.
			WithCustomMessageHandler(func(context.Context, DefaultCodecRequirements, spi.Message) bool { return true }),
		).(*defaultCodec)

		require.NoError(t, codec.Connect(t.Context()))
		<-fatalStopped
		require.False(t, codec.IsRunning(), "running after the fatal transport error")
		require.NoError(t, codec.Connect(t.Context()))
		close(release)

		// Fake time in the bubble: lets the workers pass the receive loop's cycle sleep and reach Receive.
		time.Sleep(time.Second)
		synctest.Wait()
		assert.EqualValues(t, 1, reading.Load(), "receive workers reading the transport")

		require.NoError(t, codec.Disconnect())
		assert.EqualValues(t, 1, maxReading.Load(), "most receive workers reading the transport at once")
	})
}

// An expire worker of a stopped run can take an expire notification on its way out: after a fatal transport error,
// a Connect and an Expect, its select may find both that notification and its cancelled context ready. The
// notification was meant for the next run's expire worker, which then slept on until its 30 s timer and expired
// the new expectation late. An ending expire worker hands one notification on.
func Test_defaultCodec_ExpireWork_ofStoppedRun_handsOnNotification(t *testing.T) {
	codec := buildDefaultCodec(NewMockDefaultCodecRequirements(t), NewMockTransportInstance(t),
		options.WithCustomLogger(zerolog.New(io.Discard)),
	).(*defaultCodec)
	// The next run is active, the stopped run's context is cancelled.
	codec.running.Store(true)
	stoppedRun, cancelStoppedRun := context.WithCancel(t.Context())
	cancelStoppedRun()
	// The notification of the next run's Expect, which the stopped run's worker took.
	codec.Expect(t.Context(), "next run", func(spi.Message) bool { return false }, func(spi.Message) error { return nil }, func(error) error { return nil })
	<-codec.notifyExpireWorker

	codec.ExpireWork(stoppedRun)

	select {
	case <-codec.notifyExpireWorker:
	default:
		t.Fatal("the stopped run's expire worker ended without handing on the notification")
	}
}
