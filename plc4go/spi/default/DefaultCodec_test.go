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
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultExpectation_GetAcceptsMessage(t *testing.T) {
	type fields struct {
		Context        context.Context
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
				Context:        tt.fields.Context,
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
		Context        context.Context
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
				Context:        tt.fields.Context,
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
		Context        context.Context
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
				Context:        tt.fields.Context,
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
		Context        context.Context
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
				Context:        tt.fields.Context,
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
		Context        context.Context
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
				Context:        tt.fields.Context,
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
		Context        context.Context
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
			want: "Expectation(expires at 0001-01-01 00:00:00 +0000 UTC)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultExpectation{
				Context:        tt.fields.Context,
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
		name string
		args args
		want DefaultCodec
	}{
		{
			name: "create it",
			want: &defaultCodec{
				expectations:   []spi.Expectation{},
				receiveTimeout: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultCodec(tt.args.requirements, tt.args.transportInstance, tt.args.options...)
			assert.NotNil(t, got.(*defaultCodec).defaultIncomingMessageChannel)
			got.(*defaultCodec).defaultIncomingMessageChannel = nil // Not comparable
			assert.Equalf(t, tt.want, got, "NewDefaultCodec(%v, %v, %v)", tt.args.requirements, tt.args.transportInstance, tt.args.options)
		})
	}
}

func TestWithCustomMessageHandler(t *testing.T) {
	type args struct {
		customMessageHandler func(codec DefaultCodecRequirements, message spi.Message) bool
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
		name string
		args args
		want DefaultCodec
	}{
		{
			name: "build it",
			want: &defaultCodec{
				expectations:   []spi.Expectation{},
				receiveTimeout: 10 * time.Second,
			},
		},
		{
			name: "build it with custom handler",
			args: args{
				options: []options.WithOption{
					withCustomMessageHandler{},
				},
			},
			want: &defaultCodec{
				expectations:   []spi.Expectation{},
				receiveTimeout: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildDefaultCodec(tt.args.defaultCodecRequirements, tt.args.transportInstance, tt.args.options...)
			assert.NotNil(t, got.(*defaultCodec).defaultIncomingMessageChannel)
			got.(*defaultCodec).defaultIncomingMessageChannel = nil // Not comparable
			assert.Equalf(t, tt.want, got, "buildDefaultCodec(%v, %v, %v)", tt.args.defaultCodecRequirements, tt.args.transportInstance, tt.args.options)
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(t *testing.T, fields *fields)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "connect it",
			setup: func(t *testing.T, fields *fields) {
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
				tt.setup(t, &tt.fields)
			}
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			tt.wantErr(t, m.Connect(), fmt.Sprintf("Connect()"))
		})
	}
}

func Test_defaultCodec_ConnectWithContext(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
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
			setup: func(t *testing.T, fields *fields, args *args) {
				instance := NewMockTransportInstance(t)
				instance.EXPECT().IsConnected().Return(true)
				fields.transportInstance = instance
			},
			wantErr: assert.NoError,
		},
		{
			name: "connect it (fails)",
			setup: func(t *testing.T, fields *fields, args *args) {
				instance := NewMockTransportInstance(t)
				expect := instance.EXPECT()
				expect.IsConnected().Return(false)
				expect.ConnectWithContext(mock.Anything).Return(errors.New("nope"))
				fields.transportInstance = instance
			},
			wantErr: assert.Error,
		},
		{
			name: "connect it already connected",
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
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			tt.wantErr(t, m.ConnectWithContext(tt.args.ctx), fmt.Sprintf("ConnectWithContext(%v)", tt.args.ctx))
		})
	}
}

func Test_defaultCodec_Disconnect(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
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
				customMessageHandling:         tt.fields.customMessageHandling,
			}
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
	}
	type args struct {
		ctx            context.Context
		acceptsMessage spi.AcceptsMessage
		handleMessage  spi.HandleMessage
		handleError    spi.HandleError
		ttl            time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "expect it",
			setup: func(t *testing.T, fields *fields, args *args) {
				args.ctx = testutils.TestContext(t)
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
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			tt.wantErr(t, m.Expect(tt.args.ctx, tt.args.acceptsMessage, tt.args.handleMessage, tt.args.handleError, tt.args.ttl), fmt.Sprintf("Expect(%v, func(), func(), func(), %v)", tt.args.ctx, tt.args.ttl))
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
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
				customMessageHandling:         tt.fields.customMessageHandling,
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
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
				customMessageHandling:         tt.fields.customMessageHandling,
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
	}
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
						AcceptsMessage: func(_ spi.Message) bool {
							return false
						},
					},
					&defaultExpectation{ // accepts but fails
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
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&defaultExpectation{ // accepts
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				customMessageHandling:         tt.fields.customMessageHandling,
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
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
				customMessageHandling:         tt.fields.customMessageHandling,
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
	}
	type args struct {
		ctx            context.Context
		message        spi.Message
		acceptsMessage spi.AcceptsMessage
		handleMessage  spi.HandleMessage
		handleError    spi.HandleError
		ttl            time.Duration
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
				requirements.EXPECT().Send(mock.Anything).Return(nil)
				fields.DefaultCodecRequirements = requirements

				args.ctx = testutils.TestContext(t)
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
			},
			wantErr: assert.Error,
		},
		{
			name: "send it errors",
			setup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Send(mock.Anything).Return(errors.New("nope"))
				fields.DefaultCodecRequirements = requirements

				args.ctx = testutils.TestContext(t)
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
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			tt.wantErr(t, m.SendRequest(tt.args.ctx, tt.args.message, tt.args.acceptsMessage, tt.args.handleMessage, tt.args.handleError, tt.args.ttl), fmt.Sprintf("SendRequest(%v, %v, func(), func(), func(), %v)", tt.args.ctx, tt.args.message, tt.args.ttl))
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
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
	}
	type args struct {
		now time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "timeout it (no expectations)",
		},
		{
			name: "timeout some",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			args: args{now: time.Time{}.Add(2 * time.Hour)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			m.TimeoutExpectations(tt.args.now)
		})
	}
}

func Test_defaultCodec_Work(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
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
			},
		},
		{
			name: "work hard (panics everywhere)",
			fields: fields{
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(nil, errors.New("nope"))
				fields.DefaultCodecRequirements = requirements
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
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(nil, nil)
				fields.DefaultCodecRequirements = requirements
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
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
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
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
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
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(nil, errors.New("nope"))
				fields.DefaultCodecRequirements = requirements
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message custom not handled)",
			fields: fields{
				customMessageHandling: func(_ DefaultCodecRequirements, _ spi.Message) bool {
					return false
				},
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
			},
			manipulator: func(t *testing.T, codec *defaultCodec) {
				codec.running.Store(true)
				codec.activeWorker.Add(1)
			},
		},
		{
			name: "work harder (message custom handled)",
			fields: fields{
				customMessageHandling: func(_ DefaultCodecRequirements, _ spi.Message) bool {
					return true
				},
				expectations: []spi.Expectation{
					&defaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&defaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&defaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&defaultExpectation{ // Context error
						Context: func() context.Context {
							ctx, cancelFunc := context.WithCancel(context.Background())
							cancelFunc() // Cancel it instantly
							return ctx
						}(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
			setup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultCodecRequirements(t)
				requirements.EXPECT().Receive().Return(NewMockMessage(t), nil)
				fields.DefaultCodecRequirements = requirements
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
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			go func() {
				// Stop after 200ms
				time.Sleep(200 * time.Millisecond)
				m.running.Store(false)
			}()
			m.Work()
		})
	}
}

func Test_defaultCodec_String(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements      DefaultCodecRequirements
		transportInstance             transports.TransportInstance
		defaultIncomingMessageChannel chan spi.Message
		expectations                  []spi.Expectation
		running                       bool
		customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			fields: fields{
				transportInstance: test.NewTransportInstance(test.NewTransport()),
				defaultIncomingMessageChannel: func() chan spi.Message {
					messages := make(chan spi.Message, 1)
					messages <- NewMockMessage(t)
					return messages
				}(),
				expectations: []spi.Expectation{
					func() spi.Expectation {
						expectation := NewMockExpectation(t)
						expectation.EXPECT().String().Return("yoink1")
						return expectation
					}(),
					func() spi.Expectation {
						expectation := NewMockExpectation(t)
						expectation.EXPECT().String().Return("yoink2")
						return expectation
					}(),
				},
				customMessageHandling: nil,
			},
			want: `
╔═defaultCodec═══════════════════════════════════════════════════════════════════════════════════════════╗
║╔═transportInstance╗╔═expectations═══╗╔═defaultIncomingMessageChannel╗╔═customMessageHandling╗╔═running╗║
║║       test       ║║╔═value╗╔═value╗║║         1 element(s)         ║║       b0 false       ║║b0 false║║
║╚══════════════════╝║║yoink1║║yoink2║║╚══════════════════════════════╝╚══════════════════════╝╚════════╝║
║                    ║╚══════╝╚══════╝║                                                                  ║
║                    ╚════════════════╝                                                                  ║
║╔═receiveTimeout╗╔═traceDefaultMessageCodecWorker╗                                                      ║
║║      0s       ║║           b0 false            ║                                                      ║
║╚═══════════════╝╚═══════════════════════════════╝                                                      ║
╚════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			assert.Equalf(t, tt.want, m.String(), "String()")
		})
	}
}

func Test_defaultCodec_startWorker(t *testing.T) {
	type fields struct {
		DefaultCodecRequirements       DefaultCodecRequirements
		transportInstance              transports.TransportInstance
		expectations                   []spi.Expectation
		defaultIncomingMessageChannel  chan spi.Message
		customMessageHandling          func(codec DefaultCodecRequirements, message spi.Message) bool
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
			m.startWorker()
		})
	}
}
