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
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"uuid"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

func TestDefaultExpectation_GetAcceptsMessage(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		CancelFunc     context.CancelCauseFunc
		Expiration     time.Time
		AcceptsMessage spi.AcceptsMessage
		HandleMessage  spi.HandleMessage
		HandleError    spi.HandleError
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get it",
			fields: fields{
				AcceptsMessage: func(message spi.Message) bool {
					return true
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Ctx:            tt.fields.Ctx,
				CancelFunc:     tt.fields.CancelFunc,
				Expiration:     tt.fields.Expiration,
				AcceptsMessage: tt.fields.AcceptsMessage,
				HandleMessage:  tt.fields.HandleMessage,
				HandleError:    tt.fields.HandleError,
			}
			assert.NotNilf(t, m.GetAcceptsMessage(), "GetAcceptsMessage()")
		})
	}
}

func TestDefaultExpectation_GetContext(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		CancelFunc     context.CancelCauseFunc
		Expiration     time.Time
		AcceptsMessage spi.AcceptsMessage
		HandleMessage  spi.HandleMessage
		HandleError    spi.HandleError
	}
	tests := []struct {
		name   string
		fields fields
		want   context.Context
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Ctx:            tt.fields.Ctx,
				CancelFunc:     tt.fields.CancelFunc,
				Expiration:     tt.fields.Expiration,
				AcceptsMessage: tt.fields.AcceptsMessage,
				HandleMessage:  tt.fields.HandleMessage,
				HandleError:    tt.fields.HandleError,
			}
			assert.Equalf(t, tt.want, m.GetContext(), "GetContext()")
		})
	}
}

func TestDefaultExpectation_GetExpiration(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		CancelFunc     context.CancelCauseFunc
		Expiration     time.Time
		AcceptsMessage spi.AcceptsMessage
		HandleMessage  spi.HandleMessage
		HandleError    spi.HandleError
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Ctx:            tt.fields.Ctx,
				CancelFunc:     tt.fields.CancelFunc,
				Expiration:     tt.fields.Expiration,
				AcceptsMessage: tt.fields.AcceptsMessage,
				HandleMessage:  tt.fields.HandleMessage,
				HandleError:    tt.fields.HandleError,
			}
			assert.Equalf(t, tt.want, m.GetExpiration(), "GetExpiration()")
		})
	}
}

func TestDefaultExpectation_GetHandleError(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		CancelFunc     context.CancelCauseFunc
		Expiration     time.Time
		AcceptsMessage spi.AcceptsMessage
		HandleMessage  spi.HandleMessage
		HandleError    spi.HandleError
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get it",
			fields: fields{
				HandleError: func(err error) error {
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Ctx:            tt.fields.Ctx,
				CancelFunc:     tt.fields.CancelFunc,
				Expiration:     tt.fields.Expiration,
				AcceptsMessage: tt.fields.AcceptsMessage,
				HandleMessage:  tt.fields.HandleMessage,
				HandleError:    tt.fields.HandleError,
			}
			assert.NotNilf(t, m.GetHandleError(), "GetHandleError()")
		})
	}
}

func TestDefaultExpectation_GetHandleMessage(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		CancelFunc     context.CancelCauseFunc
		Expiration     time.Time
		AcceptsMessage spi.AcceptsMessage
		HandleMessage  spi.HandleMessage
		HandleError    spi.HandleError
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get it",
			fields: fields{
				HandleMessage: func(message spi.Message) error {
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Ctx:            tt.fields.Ctx,
				CancelFunc:     tt.fields.CancelFunc,
				Expiration:     tt.fields.Expiration,
				AcceptsMessage: tt.fields.AcceptsMessage,
				HandleMessage:  tt.fields.HandleMessage,
				HandleError:    tt.fields.HandleError,
			}
			assert.NotNilf(t, m.GetHandleMessage(), "GetHandleMessage()")
		})
	}
}

func TestDefaultExpectation_String(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		CancelFunc     context.CancelCauseFunc
		Expiration     time.Time
		AcceptsMessage spi.AcceptsMessage
		HandleMessage  spi.HandleMessage
		HandleError    spi.HandleError
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			want: "Expectation '' 00000000-0000-0000-0000-000000000000 (expires at 0001-01-01 00:00:00 +0000 UTC in -2562047h47m16.854775808s)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Ctx:            tt.fields.Ctx,
				CancelFunc:     tt.fields.CancelFunc,
				Expiration:     tt.fields.Expiration,
				AcceptsMessage: tt.fields.AcceptsMessage,
				HandleMessage:  tt.fields.HandleMessage,
				HandleError:    tt.fields.HandleError,
			}
			assert.Equalf(t, tt.want, m.String(), "String()")
		})
	}
}

func TestNewDefaultCodec(t *testing.T) {
	type args struct {
		requirements      DefaultCodecRequirements
		transportInstance transports.TransportInstance
		options           []options.WithOption
	}
	tests := []struct {
		name       string
		args       args
		wantAssert func(*testing.T, DefaultCodec) bool
	}{
		{
			name: "create it",
			wantAssert: func(t *testing.T, got DefaultCodec) bool {
				require.IsType(t, &defaultCodec{}, got)
				d := got.(*defaultCodec)
				assert.NotNil(t, d.defaultIncomingMessageChannel)
				assert.Equal(t, 1*time.Minute, d.receiveTimeout)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultCodec(tt.args.requirements, tt.args.transportInstance, tt.args.options...)
			assert.Truef(t, tt.wantAssert(t, got), "NewDefaultCodec(%v, %v, %v)", tt.args.requirements, tt.args.transportInstance, tt.args.options)
		})
	}
}

func TestWithCustomMessageHandler(t *testing.T) {
	type args struct {
		customMessageHandler CustomMessageHandler
	}
	tests := []struct {
		name string
		args args
		want options.WithOption
	}{
		{
			name: "create it",
			want: withCustomMessageHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WithCustomMessageHandler(tt.args.customMessageHandler), "WithCustomMessageHandler(func())")
		})
	}
}

func Test_buildDefaultCodec(t *testing.T) {
	type args struct {
		defaultCodecRequirements DefaultCodecRequirements
		transportInstance        transports.TransportInstance
		options                  []options.WithOption
	}
	tests := []struct {
		name       string
		args       args
		wantAssert func(*testing.T, DefaultCodec) bool
	}{
		{
			name: "build it",
			wantAssert: func(t *testing.T, got DefaultCodec) bool {
				require.IsType(t, &defaultCodec{}, got)
				d := got.(*defaultCodec)
				assert.NotNil(t, d.defaultIncomingMessageChannel)
				assert.Equal(t, 1*time.Minute, d.receiveTimeout)
				return true
			},
		},
		{
			name: "build it with custom handler",
			args: args{
				options: []options.WithOption{
					withCustomMessageHandler{
						customMessageHandler: func(_ context.Context, _ DefaultCodecRequirements, _ spi.Message) bool {
							return true
						},
					},
				},
			},
			wantAssert: func(t *testing.T, got DefaultCodec) bool {
				require.IsType(t, &defaultCodec{}, got)
				d := got.(*defaultCodec)
				assert.NotNil(t, d.defaultIncomingMessageChannel)
				assert.Equal(t, 1*time.Minute, d.receiveTimeout)
				assert.NotNil(t, d.customMessageHandling)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDefaultCodec(tt.args.defaultCodecRequirements, tt.args.transportInstance, tt.args.options...)
			assert.Truef(t, tt.wantAssert(t, got), "buildDefaultCodec(%v, %v, %v)", tt.args.defaultCodecRequirements, tt.args.transportInstance, tt.args.options)
		})
	}
}

func Test_defaultCodec_Connect(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "connect it",
			args: args{
				ctx: testutils.TestContext(t),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				instance := NewMockTransportInstance(t)
				instance.EXPECT().IsConnected().Return(true)
				fields.transportInstance = instance
			},
			wantErr: assert.NoError,
		},
		{
			name: "connect it (fails)",
			args: args{
				ctx: testutils.TestContext(t),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				instance := NewMockTransportInstance(t)
				expect := instance.EXPECT()
				expect.IsConnected().Return(false)
				expect.Connect(mock.Anything).Return(errors.New("nope"))
				fields.transportInstance = instance
			},
			wantErr: assert.Error,
		},
		{
			name: "connect it already connected",
			args: args{
				ctx: testutils.TestContext(t),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				instance := NewMockTransportInstance(t)
				instance.EXPECT().IsConnected().Return(true)
				fields.transportInstance = instance
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			// Connect starts the expire/receive workers, which select on m.ctx and
			// log through m.log on panic — so ctx must be non-nil and the workers
			// (tracked by activeWorker, not wg) must be joined before the test ends.
			m.ctx, m.ctxCancel = context.WithCancel(t.Context())
			t.Cleanup(func() {
				m.running.Store(false)
				m.ctxCancel()
				m.activeWorker.Wait()
				m.wg.Wait()
			})
			tt.wantErr(t, m.Connect(tt.args.ctx), fmt.Sprintf("Connect(%v)", tt.args.ctx))
		})
	}
}

// probeWorkersAlive verifies both codec workers are alive. The notify channels are
// buffered by one, so the probe sends twice: the first send at worst parks in the
// buffer, the second blocking send can only proceed once a live worker drained one.
func probeWorkersAlive(t *testing.T, codec *defaultCodec, iteration int) {
	t.Helper()
	for _, probe := range []struct {
		worker string
		notify chan struct{}
	}{
		{"expire", codec.notifyExpireWorker},
		{"receive", codec.notifyReceiveWorker},
	} {
		for range 2 {
			select {
			case probe.notify <- struct{}{}:
			case <-time.After(5 * time.Second):
				t.Fatalf("iteration %d: %s worker is dead", iteration, probe.worker)
			}
		}
	}
}

// Test_defaultCodec_Connect_workersAliveAfterConnect is a falsification test
// for the worker-startup race: Connect used to start the workers BEFORE
// running.Store(true), so a worker observing running==false in both its loop
// condition and its restart defer terminated permanently — leaving a
// successfully "connected" codec whose expectations never expire and whose
// messages are never received (observed as an infinite hang in the cbus
// Reader tests).
func Test_defaultCodec_Connect_workersAliveAfterConnect(t *testing.T) {
	for i := range 1000 {
		requirements := NewMockDefaultCodecRequirements(t)
		requirements.EXPECT().Receive(mock.Anything).Return(nil, nil).Maybe()
		instance := NewMockTransportInstance(t)
		instance.EXPECT().IsConnected().Return(true)
		instance.EXPECT().Close().Return(nil)
		// A receive cycle in flight during Disconnect classifies the synthetic
		// "not running" error before the worker notices the shutdown.
		instance.EXPECT().ClassifyError(mock.Anything).Return(transports.TransportErrorTransient).Maybe()
		codec := buildDefaultCodec(requirements, instance, options.WithCustomLogger(testutils.ProduceTestingLogger(t))).(*defaultCodec)
		require.NoError(t, codec.Connect(testutils.TestContext(t)))
		t.Cleanup(func() {
			if codec.IsRunning() {
				_ = codec.Disconnect()
			}
		})
		probeWorkersAlive(t, codec, i)
		require.NoError(t, codec.Disconnect())
	}
}

// Test_defaultCodec_Reconnect_workersAliveAfterReconnect covers codec reuse: the
// codec context is cancelled on Disconnect and used to stay cancelled forever, so a
// reconnected codec's workers exited immediately through their ctx.Done() paths.
func Test_defaultCodec_Reconnect_workersAliveAfterReconnect(t *testing.T) {
	for i := range 100 {
		requirements := NewMockDefaultCodecRequirements(t)
		requirements.EXPECT().Receive(mock.Anything).Return(nil, nil).Maybe()
		instance := NewMockTransportInstance(t)
		instance.EXPECT().IsConnected().Return(true)
		instance.EXPECT().Close().Return(nil)
		instance.EXPECT().ClassifyError(mock.Anything).Return(transports.TransportErrorTransient).Maybe()
		codec := buildDefaultCodec(requirements, instance, options.WithCustomLogger(testutils.ProduceTestingLogger(t))).(*defaultCodec)
		require.NoError(t, codec.Connect(testutils.TestContext(t)))
		t.Cleanup(func() {
			if codec.IsRunning() {
				_ = codec.Disconnect()
			}
		})
		require.NoError(t, codec.Disconnect())
		require.NoError(t, codec.Connect(testutils.TestContext(t)))
		probeWorkersAlive(t, codec, i)
		require.NoError(t, codec.Disconnect())
	}
}

// Test_defaultCodec_Expect_afterDisconnect_doesNotPanic: Disconnect used to close the
// notify channels, so a late Expect (e.g. a request racing a connection close) panicked
// with a send on a closed channel.
func Test_defaultCodec_Expect_afterDisconnect_doesNotPanic(t *testing.T) {
	requirements := NewMockDefaultCodecRequirements(t)
	requirements.EXPECT().Receive(mock.Anything).Return(nil, nil).Maybe()
	instance := NewMockTransportInstance(t)
	instance.EXPECT().IsConnected().Return(true)
	instance.EXPECT().Close().Return(nil)
	instance.EXPECT().ClassifyError(mock.Anything).Return(transports.TransportErrorTransient).Maybe()
	codec := buildDefaultCodec(requirements, instance, options.WithCustomLogger(testutils.ProduceTestingLogger(t))).(*defaultCodec)
	require.NoError(t, codec.Connect(testutils.TestContext(t)))
	require.NoError(t, codec.Disconnect())
	assert.NotPanics(t, func() {
		codec.Expect(testutils.TestContext(t), "late expect",
			func(message spi.Message) bool { return false },
			func(message spi.Message) error { return nil },
			func(err error) error { return nil })
	})
}

func Test_defaultCodec_Disconnect(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	tests := []struct {
		name        string
		fields      fields
		setup       func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, codec *defaultCodec)
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name:    "disconnect it (not running)",
			wantErr: assert.Error,
		},
		{
			name: "disconnect it",
			setup: func(t *testing.T, fields *fields) {
				instance := NewMockTransportInstance(t)
				instance.EXPECT().Close().Return(nil)
				fields.transportInstance = instance
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			c.ctx, c.ctxCancel = context.WithCancel(t.Context())
			if tt.manipulator != nil {
				tt.manipulator(t, c)
			}
			tt.wantErr(t, c.Disconnect(), fmt.Sprintf("Disconnect()"))
		})
	}
}

func Test_defaultCodec_Expect(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	type args struct {
		ctx             context.Context
		interactionInfo string
		acceptsMessage  spi.AcceptsMessage
		handleMessage   spi.HandleMessage
		handleError     spi.HandleError
		ttl             time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args)
	}{
		{
			name: "expect it",
			setup: func(t *testing.T, fields *fields, args *args) {
				args.ctx = testutils.TestContext(t)
				args.interactionInfo = t.Name()
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 20*time.Second)
				t.Cleanup(cancelFunc)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			m.Expect(tt.args.ctx, tt.args.interactionInfo, tt.args.acceptsMessage, tt.args.handleMessage, tt.args.handleError)
		})
	}
}

func Test_defaultCodec_GetDefaultIncomingMessageChannel(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   chan spi.Message
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, m.GetDefaultIncomingMessageChannel(), "GetDefaultIncomingMessageChannel()")
		})
	}
}

func Test_defaultCodec_GetTransportInstance(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   transports.TransportInstance
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, m.GetTransportInstance(), "GetTransportInstance()")
		})
	}
}

func Test_defaultCodec_HandleMessages(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args)
		want   bool
	}{
		{
			name: "handle it (nothing to handle)",
		},
		{
			name: "handle some",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // doesn't accept
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return false
						},
					},
					&defaultExpectation{ // accepts but fails
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return errors.New("oh noes")
						},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts but fails and fails to handle the error
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return errors.New("oh noes")
						},
						HandleError: func(err error) error {
							return errors.New("I failed completely")
						},
					},
					&defaultExpectation{ // accepts
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // not accept
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return false
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
				},
			},
			want: true,
		},
		{
			name: "handle some (ensure everyone get's it)",
			setup: func(t *testing.T, fields *fields, args *args) {
				accept1 := atomic.Bool{}
				accept2 := atomic.Bool{}
				accept3 := atomic.Bool{}
				accept4 := atomic.Bool{}
				accept5 := atomic.Bool{}
				accept6 := atomic.Bool{}
				t.Cleanup(func() {
					assert.True(t, accept1.Load(), "accept1 not called")
					assert.True(t, accept2.Load(), "accept2 not called")
					assert.True(t, accept3.Load(), "accept3 not called")
					assert.True(t, accept4.Load(), "accept4 not called")
					assert.True(t, accept5.Load(), "accept5 not called")
					assert.True(t, accept6.Load(), "accept6 not called")
				})
				fields.expectations = []spi.Expectation{
					&defaultExpectation{ // doesn't accept
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return false
						},
					},
					&defaultExpectation{ // accepts but fails // accept1
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							accept1.Store(true)
							return errors.New("oh noes")
						},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts but fails and fails to handle the error // accept2
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							accept2.Store(true)
							return errors.New("oh noes")
						},
						HandleError: func(err error) error {
							return errors.New("I failed completely")
						},
					},
					&defaultExpectation{ // accepts // accept3
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							accept3.Store(true)
							return nil
						},
					},
					&defaultExpectation{ // accepts // accept4
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							accept4.Store(true)
							return nil
						},
					},
					&defaultExpectation{ // not accept // accept5
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							accept5.Store(true)
							return nil
						},
					},
					&defaultExpectation{ // not accept
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return false
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts // accept6
						Uuid:       uuid.New(),
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							accept6.Store(true)
							return nil
						},
					},
				}
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, m.HandleMessages(tt.args.message), "HandleMessages(%v)", tt.args.message)
		})
	}
}

func Test_defaultCodec_IsRunning(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "is it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, m.IsRunning(), "IsRunning()")
		})
	}
}

func Test_defaultCodec_SendRequest(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	type args struct {
		ctx             context.Context
		interactionInfo string
		message         spi.Message
		acceptsMessage  spi.AcceptsMessage
		handleMessage   spi.HandleMessage
		handleError     spi.HandleError
		ttl             time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "send it",
			setup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Send(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.DefaultCodecRequirements = requirements

				args.ctx = testutils.TestContext(t)
				args.interactionInfo = t.Name()
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 20*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.NoError,
		},
		{
			name: "send it canceled",
			setup: func(t *testing.T, fields *fields, args *args) {
				fields.DefaultCodecRequirements = NewMockDefaultCodecRequirements(t)

				ctx, cancelFunc := context.WithCancel(testutils.TestContext(t))
				cancelFunc()
				args.ctx = ctx
				args.interactionInfo = t.Name()
			},
			wantErr: assert.Error,
		},
		{
			name: "send it errors",
			setup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Send(mock.Anything, mock.Anything, mock.Anything).Return(errors.New("nope"))
				fields.DefaultCodecRequirements = requirements

				args.ctx = testutils.TestContext(t)
				args.interactionInfo = t.Name()
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 20*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			tt.wantErr(t, m.SendRequest(tt.args.ctx, tt.args.interactionInfo, tt.args.message, tt.args.acceptsMessage, tt.args.handleMessage, tt.args.handleError), fmt.Sprintf("SendRequest(%v, %v, func(), func(), func(), %v)", tt.args.ctx, tt.args.message, tt.args.ttl))
		})
	}
}

func Test_defaultCodec_TimeoutExpectations(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args)
		want   time.Duration
	}{
		{
			name: "timeout it (no expectations)",
			want: 30 * time.Second,
		},
		{
			name: "timeout some",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			args: args{now: time.Time{}.Add(2 * time.Hour)},
			want: time.Until(time.Time{}.Add(3 * time.Hour)),
		},
		{
			name: "timeout some (ensure everyone is called)",
			args: args{now: time.Time{}.Add(2 * time.Hour)},
			setup: func(t *testing.T, fields *fields, args *args) {
				handle1 := atomic.Bool{}
				handle2 := atomic.Bool{}
				handle3 := atomic.Bool{}
				handle4 := atomic.Bool{}
				handle5 := atomic.Bool{}
				t.Cleanup(func() {
					time.Sleep(100 * time.Millisecond) // TODO: doing a sleep as handle error is called in a gofunc
					assert.True(t, handle1.Load(), "handle1 not called")
					assert.True(t, handle2.Load(), "handle2 not called")
					assert.False(t, handle3.Load(), "handle3 called")
					assert.True(t, handle4.Load(), "handle4 not called")
					assert.False(t, handle5.Load(), "handle5 called")
				})
				fields.expectations = []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							handle1.Store(true)
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							handle2.Store(true)
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							handle3.Store(true)
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							handle4.Store(true)
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							handle5.Store(true)
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				}
			},
			want: time.Until(time.Time{}.Add(3 * time.Hour)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				log:                           testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, m.TimeoutExpectations(tt.args.now), "TimeoutExpectations(%v)", tt.args.now)
			m.wg.Wait()
		})
	}
}

func Test_defaultCodec_ReceiveWork(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         CustomMessageHandler
		ctx                           context.Context
		ctxCancel                     context.CancelFunc
	}
	tests := []struct {
		name        string
		fields      fields
		setup       func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, codec *defaultCodec)
	}{
		{
			name: "work it (nothing to do)",
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
		},
		{
			name: "work hard (panics everywhere)",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(nil, errors.New("nope"))
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (nil message)",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(nil, nil)
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message)",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message with default channel)",
			fields: fields{
				defaultIncomingMessageChannel: make(chan spi.Message, 1),
				expectations: []spi.Expectation{
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message receive error)",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(nil, errors.New("nope"))
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message custom not handled)",
			fields: fields{
				customMessageHandling: func(_ context.Context, _ DefaultCodecRequirements, _ spi.Message) bool {
					return false
				},
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message custom handled)",
			fields: fields{
				customMessageHandling: func(_ context.Context, _ DefaultCodecRequirements, _ spi.Message) bool {
					return true
				},
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Ctx:        t.Context(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Ctx: func() context.Context {
							ctx, cancelFunc := context.WithCancel(t.Context())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						CancelFunc: func(_ error) {},
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive(mock.Anything).Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
				fields.ctx, fields.ctxCancel = context.WithTimeout(t.Context(), 2*time.Second)
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				notifyExpireWorker:            make(chan struct{}, 100),
				notifyReceiveWorker:           make(chan struct{}, 100),
				customMessageHandling:         tt.fields.customMessageHandling,
				ctx:                           tt.fields.ctx,
				ctxCancel:                     tt.fields.ctxCancel,
				log:                           testutils.ProduceTestingLogger(t),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			t.Cleanup(m.wg.Wait)
			m.wg.Go(func() {
				// Stop after 200ms
				timer := time.NewTimer(200 * time.Millisecond)
				select {
				case <-timer.C:
				case <-t.Context().Done():
				}
				m.running.Store(false)
			})
			m.ReceiveWork()
		})
	}
}

func Test_defaultCodec_startWorkers(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements       DefaultCodecRequirements
		transportInstance              transports.TransportInstance
		expectations                   []spi.Expectation
		defaultIncomingMessageChannel  chan spi.Message
		customMessageHandling          CustomMessageHandler
		receiveTimeout                 time.Duration
		traceDefaultMessageCodecWorker bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "start it not running",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:       tt.fields.DefaultCodecRequirements,
				transportInstance:              tt.fields.transportInstance,
				expectations:                   tt.fields.expectations,
				defaultIncomingMessageChannel:  tt.fields.defaultIncomingMessageChannel,
				customMessageHandling:          tt.fields.customMessageHandling,
				receiveTimeout:                 tt.fields.receiveTimeout,
				traceDefaultMessageCodecWorker: tt.fields.traceDefaultMessageCodecWorker,
				log:                            testutils.ProduceTestingLogger(t),
			}
			m.startWorkers()
		})
	}
}

func Test_defaultCodec_integration(t *testing.T) {
	mockDefaultCodecRequirements := NewMockDefaultCodecRequirements(t)
	{
		expect := mockDefaultCodecRequirements.EXPECT()
		message := NewMockMessage(t)
		{
			expect := message.EXPECT()
			expect.String().Return("message for " + t.Name())
			expect.SerializeWithWriteBuffer(mock.Anything, mock.Anything).Return(nil)
		}
		expect.Receive(mock.Anything).RunAndReturn(func(_ context.Context) (spi.Message, error) {
			// Simulate a bit read delay
			timer := time.NewTimer(100 * time.Millisecond)
			select {
			case <-timer.C:
			case <-t.Context().Done():
			}
			if err := t.Context().Err(); err != nil {
				return nil, err
			}
			return message, nil
		})
	}
	mockTransportInstance := NewMockTransportInstance(t)
	{
		expect := mockTransportInstance.EXPECT()
		expect.IsConnected().Return(true)
		expect.Close().Return(nil)
	}
	sut := NewDefaultCodec(mockDefaultCodecRequirements, mockTransportInstance,
		options.WithCustomLogger(testutils.ProduceTestingLogger(t)),
		options.WithTraceDefaultMessageCodecWorker(true),
	)
	t.Cleanup(func() {
		_ = sut.Disconnect()
	})
	// First expect
	var firstHandled bool
	sut.Expect(t.Context(), "first", func(message spi.Message) bool {
		t.Log("accepts message", message)
		return true
	}, func(message spi.Message) error {
		t.Log("handle message", message)
		firstHandled = true
		return nil
	}, func(err error) error {
		t.Log("error", err)
		return nil
	})
	// Second expect
	var secondHandled bool
	sut.Expect(t.Context(), "second", func(message spi.Message) bool {
		t.Log("accepts message", message)
		return true
	}, func(message spi.Message) error {
		t.Log("handle message", message)
		secondHandled = true
		return nil
	}, func(err error) error {
		t.Log("error", err)
		return nil
	})
	// Third expect
	var thridErrorCalled bool
	sut.Expect(t.Context(), "third", func(message spi.Message) bool {
		t.Log("does not accept message", message)
		return false
	}, func(message spi.Message) error {
		t.Error("should not be called")
		return nil
	}, func(err error) error {
		thridErrorCalled = true
		return nil
	})
	// Fourth expect
	var fourthHandled bool
	sut.Expect(t.Context(), "fourth", func(message spi.Message) bool {
		t.Log("accepts message", message)
		return true
	}, func(message spi.Message) error {
		t.Log("handle message", message)
		fourthHandled = true
		return nil
	}, func(err error) error {
		t.Log("error", err)
		return nil
	})

	err := sut.Connect(t.Context())
	assert.NoError(t, err)
	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
	case <-t.Context().Done():
	}
	assert.NoError(t, sut.Disconnect())
	assert.True(t, firstHandled)
	assert.True(t, secondHandled)
	assert.True(t, thridErrorCalled) // because of our disconnect
	assert.True(t, fourthHandled)
}
