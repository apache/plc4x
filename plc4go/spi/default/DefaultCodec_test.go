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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
			m := &DefaultExpectation{
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
			m := &DefaultExpectation{
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
			m := &DefaultExpectation{
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
			m := &DefaultExpectation{
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
			m := &DefaultExpectation{
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
			m := &DefaultExpectation{
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
				expectations: []spi.Expectation{},
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
				expectations: []spi.Expectation{},
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
				expectations: []spi.Expectation{},
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "connect it",
			fields: fields{
				transportInstance: testTransportInstance{},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				running:                       tt.fields.running,
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "connect it",
			fields: fields{
				transportInstance: testTransportInstance{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "connect it (fails)",
			fields: fields{
				transportInstance: testTransportInstance{
					connectWithContext: func(_ context.Context) error {
						return errors.New("nope")
					},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "connect it already connected",
			fields: fields{
				transportInstance: testTransportInstance{
					connected: true,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				running:                       tt.fields.running,
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
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "disconnect it",
			fields: fields{
				transportInstance: testTransportInstance{},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				running:                       tt.fields.running,
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			tt.wantErr(t, m.Disconnect(), fmt.Sprintf("Disconnect()"))
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "expect it",
			args: args{
				ctx: context.Background(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				running:                       tt.fields.running,
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
				running:                       tt.fields.running,
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
				running:                       tt.fields.running,
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
					&DefaultExpectation{ // doesn't accept
						AcceptsMessage: func(_ spi.Message) bool {
							return false
						},
					},
					&DefaultExpectation{ // accepts but fails
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
					&DefaultExpectation{ // accepts but fails and fails to handle the error
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
					&DefaultExpectation{ // accepts
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&DefaultExpectation{ // accepts
						AcceptsMessage: func(_ spi.Message) bool {
							return true
						},
						HandleMessage: func(_ spi.Message) error {
							return nil
						},
					},
					&DefaultExpectation{ // accepts
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
				running:                       tt.fields.running,
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
				running:                       tt.fields.running,
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "send it",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "send it canceled",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			wantErr: assert.Error,
		},
		{
			name: "send it errors",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{
					send: func(_ spi.Message) error {
						return errors.New("nope")
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				running:                       tt.fields.running,
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
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
				running:                       tt.fields.running,
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
	type args struct {
		codec DefaultCodecRequirements
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "work it (nothing to do)",
			fields: fields{
				running: true,
			},
		},
		{
			name: "work hard (panics everywhere)",
			fields: fields{
				running: true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
		},
		{
			name: "work harder (nil message)",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{},
				running:                  true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
		},
		{
			name: "work harder (message)",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{
					receive: func() (spi.Message, error) {
						return testMessage{}, nil
					},
				},
				running: true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
		},
		{
			name: "work harder (message with default channel)",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{
					receive: func() (spi.Message, error) {
						return testMessage{}, nil
					},
				},
				defaultIncomingMessageChannel: make(chan spi.Message, 1),
				running:                       true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
				},
			},
		},
		{
			name: "work harder (message receive error)",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{
					receive: func() (spi.Message, error) {
						return nil, errors.New("nope")
					},
				},
				running: true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
		},
		{
			name: "work harder (message custom not handled)",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{
					receive: func() (spi.Message, error) {
						return testMessage{}, nil
					},
				},
				customMessageHandling: func(_ DefaultCodecRequirements, _ spi.Message) bool {
					return false
				},
				running: true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
		},
		{
			name: "work harder (message custom handled)",
			fields: fields{
				DefaultCodecRequirements: testMessageCodec{
					receive: func() (spi.Message, error) {
						return testMessage{}, nil
					},
				},
				customMessageHandling: func(_ DefaultCodecRequirements, _ spi.Message) bool {
					return true
				},
				running: true,
				expectations: []spi.Expectation{
					&DefaultExpectation{ // Expired
						Context: context.Background(),
						HandleError: func(err error) error {
							return nil
						},
					},
					&DefaultExpectation{ // Expired errors
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
					},
					&DefaultExpectation{ // Fine
						Context: context.Background(),
						HandleError: func(err error) error {
							return errors.New("yep")
						},
						Expiration: time.Time{}.Add(3 * time.Hour),
					},
					&DefaultExpectation{ // Context error
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &defaultCodec{
				DefaultCodecRequirements:      tt.fields.DefaultCodecRequirements,
				transportInstance:             tt.fields.transportInstance,
				defaultIncomingMessageChannel: tt.fields.defaultIncomingMessageChannel,
				expectations:                  tt.fields.expectations,
				running:                       tt.fields.running,
				customMessageHandling:         tt.fields.customMessageHandling,
			}
			go func() {
				// Stop after 200ms
				time.Sleep(200 * time.Millisecond)
				m.running = false
			}()
			m.Work(tt.args.codec)
		})
	}
}
