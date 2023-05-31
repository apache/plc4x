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
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

func TestAlphaGenerator_getAndIncrement(t *testing.T) {
	type fields struct {
		currentAlpha byte
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		{
			name: "get a alpha invalid instance",
		},
		{
			name: "get a alpha",
			fields: fields{
				currentAlpha: 'g',
			},
			want: 'g',
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AlphaGenerator{
				currentAlpha: tt.fields.currentAlpha,
			}
			assert.Equalf(t, tt.want, a.getAndIncrement(), "getAndIncrement()")
		})
	}
}

func TestAlphaGenerator_getAndIncrement_Turnaround(t *testing.T) {
	a := &AlphaGenerator{
		currentAlpha: 'y',
	}
	// Currently it is 'y' so the next call should return 'z'
	assert.Equal(t, a.getAndIncrement(), uint8('y'))
	// Currently it is 'z' so the next call should return 'g' as we roll over
	assert.Equal(t, a.getAndIncrement(), uint8('z'))
	// Currently it is 'g' so the next call should return 'h'
	assert.Equal(t, a.getAndIncrement(), uint8('g'))
	// the final 'h'
	assert.Equal(t, a.getAndIncrement(), uint8('h'))
}

func TestConnection_BrowseRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcBrowseRequestBuilder) bool
	}{
		{
			name: "return not nil",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
			wantAssert: func(t *testing.T, builder apiModel.PlcBrowseRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.True(t, tt.wantAssert(t, c.BrowseRequestBuilder()), "BrowseRequestBuilder()")
		})
	}
}

func TestConnection_ConnectWithContext(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		setup        func(t *testing.T, fields *fields)
		wantAsserter func(*testing.T, <-chan plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "just connect and fail",
			fields: fields{
				configuration: Configuration{
					Srchk:                 false,
					Exstat:                false,
					Pun:                   false,
					LocalSal:              false,
					Pcn:                   false,
					Idmon:                 false,
					Monitor:               false,
					Smart:                 false,
					XonXoff:               false,
					Connect:               false,
					MonitoredApplication1: 0,
					MonitoredApplication2: 0,
				},
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
				connectionId: "connectionId13",
				tracer:       nil,
			},
			args: args{ctx: context.Background()},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)
				fields.messageCodec = NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport(loggerOption)
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, loggerOption)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}(), loggerOption)
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				assert.NotNil(t, results)
				result := <-results
				assert.Nil(t, result.GetConnection())
				assert.NotNil(t, result.GetErr())
				return true
			},
		},
		// TODO: add error case for failing messageCodec connect
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.True(t, tt.wantAsserter(t, c.ConnectWithContext(tt.args.ctx)), "ConnectWithContext(%v)", tt.args.ctx)
		})
	}
}

func TestConnection_GetConnection(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name         string
		fields       fields
		wantAsserter func(*testing.T, plc4go.PlcConnection) bool
	}{
		{
			name: "not nil",
			wantAsserter: func(t *testing.T, connection plc4go.PlcConnection) bool {
				return assert.NotNil(t, connection)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Truef(t, tt.wantAsserter(t, c.GetConnection()), "GetConnection()")
		})
	}
}

func TestConnection_GetConnectionId(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.GetConnectionId(), "GetConnectionId()")
		})
	}
}

func TestConnection_GetMessageCodec(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.MessageCodec
	}{
		{
			name: "just get",
			fields: fields{
				messageCodec: &MessageCodec{},
			},
			want: &MessageCodec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.GetMessageCodec(), "GetMessageCodec()")
		})
	}
}

func TestConnection_GetMetadata(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcConnectionMetadata
	}{
		{
			name: "give metadata",
			want: _default.DefaultConnectionMetadata{
				ConnectionAttributes: nil,
				ProvidesReading:      true,
				ProvidesWriting:      true,
				ProvidesSubscribing:  true,
				ProvidesBrowsing:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.GetMetadata(), "GetMetadata()")
		})
	}
}

func TestConnection_GetTracer(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   *spi.Tracer
	}{
		{
			name: "just nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.GetTracer(), "GetTracer()")
		})
	}
}

func TestConnection_IsTraceEnabled(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "not enabled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.IsTraceEnabled(), "IsTraceEnabled()")
		})
	}
}

func TestConnection_ReadRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcReadRequestBuilder) bool
	}{
		{
			name: "return not nil",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
			wantAssert: func(t *testing.T, builder apiModel.PlcReadRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Truef(t, tt.wantAssert(t, c.ReadRequestBuilder()), "ReadRequestBuilder()")
		})
	}
}

func TestConnection_String(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "a string",
			want: "cbus.Connection",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func TestConnection_SubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcSubscriptionRequestBuilder) bool
	}{
		{
			name: "return not nil",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
			wantAssert: func(t *testing.T, builder apiModel.PlcSubscriptionRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Truef(t, tt.wantAssert(t, c.SubscriptionRequestBuilder()), "SubscriptionRequestBuilder()")
		})
	}
}

func TestConnection_UnsubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcUnsubscriptionRequestBuilder
	}{
		{
			name: "create one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.UnsubscriptionRequestBuilder(), "UnsubscriptionRequestBuilder()")
		})
	}
}

func TestConnection_WriteRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcWriteRequestBuilder) bool
	}{
		{
			name: "return not nil",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
			wantAssert: func(t *testing.T, builder apiModel.PlcWriteRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Truef(t, tt.wantAssert(t, c.WriteRequestBuilder()), "WriteRequestBuilder()")
		})
	}
}

func TestConnection_addSubscriber(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		subscriber *Subscriber
	}
	theOneSubscriber := NewSubscriber(nil)
	tests := []struct {
		name        string
		fields      fields
		args        args
		subElevator func(*testing.T, []*Subscriber) bool
	}{
		{
			name: "new subscriber",
			args: args{subscriber: NewSubscriber(nil)},
			subElevator: func(t *testing.T, subscribers []*Subscriber) bool {
				return len(subscribers) == 1
			},
		},
		{
			name:   "existing subscriber should not be added",
			fields: fields{subscribers: []*Subscriber{theOneSubscriber}},
			args:   args{subscriber: theOneSubscriber},
			subElevator: func(t *testing.T, subscribers []*Subscriber) bool {
				return len(subscribers) == 1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			c.addSubscriber(tt.args.subscriber)
			assert.Truef(t, tt.subElevator(t, c.subscribers), "addSubscriber(%v)", tt.args.subscriber)
		})
	}
}

func TestConnection_fireConnected(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ch chan<- plc4go.PlcConnectionConnectResult
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		chanValidator func(*testing.T, chan<- plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "instant connect",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
			chanValidator: func(_ *testing.T, _ chan<- plc4go.PlcConnectionConnectResult) bool {
				return true
			},
		},
		{
			name: "notified connect",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				driverContext: DriverContext{
					awaitSetupComplete: true,
				},
			},
			args: args{ch: make(chan<- plc4go.PlcConnectionConnectResult, 1)},
			chanValidator: func(t *testing.T, results chan<- plc4go.PlcConnectionConnectResult) bool {
				time.Sleep(time.Millisecond * 50)
				return len(results) == 1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			c.fireConnected(tt.args.ch)
			assert.True(t, tt.chanValidator(t, tt.args.ch))
		})
	}
}

func TestConnection_fireConnectionError(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		err error
		ch  chan<- plc4go.PlcConnectionConnectResult
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		chanValidator func(*testing.T, chan<- plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "instant connect",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}()),
			},
			chanValidator: func(_ *testing.T, _ chan<- plc4go.PlcConnectionConnectResult) bool {
				return true
			},
		},
		{
			name: "notified connect",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}()),
				driverContext: DriverContext{
					awaitSetupComplete: true,
				},
			},
			args: args{ch: make(chan<- plc4go.PlcConnectionConnectResult, 1)},
			chanValidator: func(t *testing.T, results chan<- plc4go.PlcConnectionConnectResult) bool {
				time.Sleep(time.Millisecond * 50)
				return len(results) == 1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			c.fireConnectionError(tt.args.err, tt.args.ch)
			assert.True(t, tt.chanValidator(t, tt.args.ch))
		})
	}
}

func TestConnection_sendCalDataWrite(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx            context.Context
		ch             chan plc4go.PlcConnectionConnectResult
		paramNo        readWriteModel.Parameter
		parameterValue readWriteModel.ParameterValue
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		want   bool
	}{
		{
			name: "send something",
			args: args{
				ctx:            context.Background(),
				ch:             make(chan plc4go.PlcConnectionConnectResult, 1),
				paramNo:        readWriteModel.Parameter_APPLICATION_ADDRESS_2,
				parameterValue: readWriteModel.NewParameterValueApplicationAddress2(readWriteModel.NewApplicationAddress2(1), nil, 0),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
				cbusOptions: func() *readWriteModel.CBusOptions {
					var cBusOptions readWriteModel.CBusOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
					return &cBusOptions
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)
				fields.messageCodec = NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport(loggerOption)
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, loggerOption)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}())
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.sendCalDataWrite(tt.args.ctx, tt.args.ch, tt.args.paramNo, tt.args.parameterValue, tt.args.requestContext, tt.args.cbusOptions), "sendCalDataWrite(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.paramNo, tt.args.parameterValue, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_sendReset(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx                      context.Context
		ch                       chan plc4go.PlcConnectionConnectResult
		cbusOptions              *readWriteModel.CBusOptions
		requestContext           *readWriteModel.RequestContext
		sendOutErrorNotification bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		{
			name: "send reset",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}()),
			},
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
				cbusOptions: func() *readWriteModel.CBusOptions {
					var cBusOptions readWriteModel.CBusOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
					return &cBusOptions
				}(),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
				sendOutErrorNotification: false,
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.wantOk, c.sendReset(tt.args.ctx, tt.args.ch, tt.args.cbusOptions, tt.args.requestContext, tt.args.sendOutErrorNotification), "sendReset(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.cbusOptions, tt.args.requestContext, tt.args.sendOutErrorNotification)
		})
	}
}

func TestConnection_setApplicationFilter(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx            context.Context
		ch             chan plc4go.PlcConnectionConnectResult
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		wantOk bool
	}{
		{
			name: "set application filter (failing)",
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
				cbusOptions: func() *readWriteModel.CBusOptions {
					var cBusOptions readWriteModel.CBusOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
					return &cBusOptions
				}(),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Setup connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)
				fields.messageCodec = NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport(loggerOption)
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, loggerOption)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}(), loggerOption)
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.wantOk, c.setApplicationFilter(tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions), "setApplicationFilter(%v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setInterface1PowerUpSettings(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx            context.Context
		ch             chan plc4go.PlcConnectionConnectResult
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args)
		wantOk bool
	}{
		{
			name: "set interface 1 PUN options (failing)",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}()),
			},
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
				cbusOptions: func() *readWriteModel.CBusOptions {
					var cBusOptions readWriteModel.CBusOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
					return &cBusOptions
				}(),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.wantOk, c.setInterface1PowerUpSettings(tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions), "setInterface1PowerUpSettings(%v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setInterfaceOptions1(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx            context.Context
		ch             chan plc4go.PlcConnectionConnectResult
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		want   bool
	}{
		{
			name: "set interface 1 options (failing)",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}()),
			},
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
				cbusOptions: func() *readWriteModel.CBusOptions {
					var cBusOptions readWriteModel.CBusOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
					return &cBusOptions
				}(),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.want, c.setInterfaceOptions1(tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions), "setInterfaceOptions1(%v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setInterfaceOptions3(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx            context.Context
		ch             chan plc4go.PlcConnectionConnectResult
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		wantOk bool
	}{
		{
			name: "set interface 3 options (failing)",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}()),
			},
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
				cbusOptions: func() *readWriteModel.CBusOptions {
					var cBusOptions readWriteModel.CBusOptions = readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
					return &cBusOptions
				}(),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			assert.Equalf(t, tt.wantOk, c.setInterfaceOptions3(tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions), "setInterfaceOptions3(%v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setupConnection(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	type args struct {
		ctx context.Context
		ch  chan plc4go.PlcConnectionConnectResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
	}{
		{
			name: "setup connection (failing)",
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)

				// Build the message codec
				fields.messageCodec = NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}(), loggerOption)
			},
		},
		{
			name: "setup connection (failing after reset)",
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)

				// Build the message codec
				transport := test.NewTransport(loggerOption)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, loggerOption)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				type MockState uint8
				const (
					RESET MockState = iota
					DONE
				)
				currentState := atomic.Value{}
				currentState.Store(RESET)
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					switch currentState.Load().(MockState) {
					case RESET:
						t.Log("Dispatching reset echo")
						transportInstance.FillReadBuffer([]byte("~~~\r"))
						currentState.Store(DONE)
					case DONE:
						t.Log("Done")
					}
				})
				codec := NewMessageCodec(transportInstance, loggerOption)
				err = codec.Connect()
				if err != nil {
					t.Error(err)
					t.FailNow()
				}

				fields.messageCodec = codec
			},
		},
		{
			name: "setup connection (failing after app filters)",
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)

				// Build the message codec
				transport := test.NewTransport(loggerOption)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, loggerOption)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				type MockState uint8
				const (
					RESET MockState = iota
					APPLICATION_FILTER_1
					APPLICATION_FILTER_2
					DONE
				)
				currentState := atomic.Value{}
				currentState.Store(RESET)
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					switch currentState.Load().(MockState) {
					case RESET:
						t.Log("Dispatching reset echo")
						transportInstance.FillReadBuffer([]byte("~~~\r"))
						currentState.Store(APPLICATION_FILTER_1)
					case APPLICATION_FILTER_1:
						t.Log("Dispatching app1 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A32100FF\r"))
						transportInstance.FillReadBuffer([]byte("322100AD\r\n"))
						currentState.Store(APPLICATION_FILTER_2)
					case APPLICATION_FILTER_2:
						t.Log("Dispatching app2 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A32200FF\r"))
						transportInstance.FillReadBuffer([]byte("322200AC\r\n"))
						currentState.Store(DONE)
					case DONE:
						t.Log("Done")
					}
				})
				codec := NewMessageCodec(transportInstance, loggerOption)
				err = codec.Connect()
				if err != nil {
					t.Error(err)
					t.FailNow()
				}

				fields.messageCodec = codec
			},
		},
		{
			name: "setup connection (failing after interface options 3",
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)

				// Build the message codec
				transport := test.NewTransport(loggerOption)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, loggerOption)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				type MockState uint8
				const (
					RESET MockState = iota
					APPLICATION_FILTER_1
					APPLICATION_FILTER_2
					INTERFACE_OPTIONS_3
					DONE
				)
				currentState := atomic.Value{}
				currentState.Store(RESET)
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					switch currentState.Load().(MockState) {
					case RESET:
						t.Log("Dispatching reset echo")
						transportInstance.FillReadBuffer([]byte("~~~\r"))
						currentState.Store(APPLICATION_FILTER_1)
					case APPLICATION_FILTER_1:
						t.Log("Dispatching app1 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A32100FF\r"))
						transportInstance.FillReadBuffer([]byte("322100AD\r\n"))
						currentState.Store(APPLICATION_FILTER_2)
					case APPLICATION_FILTER_2:
						t.Log("Dispatching app2 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A32200FF\r"))
						transportInstance.FillReadBuffer([]byte("322200AC\r\n"))
						currentState.Store(INTERFACE_OPTIONS_3)
					case INTERFACE_OPTIONS_3:
						t.Log("Dispatching interface 3 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A342000A\r"))
						transportInstance.FillReadBuffer([]byte("3242008C\r\n"))
						currentState.Store(DONE)
					case DONE:
						t.Log("Done")
					}
				})
				codec := NewMessageCodec(transportInstance, loggerOption)
				err = codec.Connect()
				if err != nil {
					t.Error(err)
					t.FailNow()
				}

				fields.messageCodec = codec
			},
		},
		{
			name: "setup connection (failing after interface options 1 pun)",
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)

				// Build the message codec
				transport := test.NewTransport(loggerOption)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, loggerOption)
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				type MockState uint8
				const (
					RESET MockState = iota
					APPLICATION_FILTER_1
					APPLICATION_FILTER_2
					INTERFACE_OPTIONS_3
					INTERFACE_OPTIONS_1_PUN
					DONE
				)
				currentState := atomic.Value{}
				currentState.Store(RESET)
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					switch currentState.Load().(MockState) {
					case RESET:
						t.Log("Dispatching reset echo")
						transportInstance.FillReadBuffer([]byte("~~~\r"))
						currentState.Store(APPLICATION_FILTER_1)
					case APPLICATION_FILTER_1:
						t.Log("Dispatching app1 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A32100FF\r"))
						transportInstance.FillReadBuffer([]byte("322100AD\r\n"))
						currentState.Store(APPLICATION_FILTER_2)
					case APPLICATION_FILTER_2:
						t.Log("Dispatching app2 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A32200FF\r"))
						transportInstance.FillReadBuffer([]byte("322200AC\r\n"))
						currentState.Store(INTERFACE_OPTIONS_3)
					case INTERFACE_OPTIONS_3:
						t.Log("Dispatching interface 3 echo and confirm")
						transportInstance.FillReadBuffer([]byte("@A342000A\r"))
						transportInstance.FillReadBuffer([]byte("3242008C\r\n"))
						currentState.Store(INTERFACE_OPTIONS_1_PUN)
					case INTERFACE_OPTIONS_1_PUN:
						t.Log("Dispatching interface 1 PUN echo and confirm???")
						transportInstance.FillReadBuffer([]byte("@A3410079\r"))
						transportInstance.FillReadBuffer([]byte("3241008D\r\n"))
						currentState.Store(DONE)
					case DONE:
						t.Log("Done")
					}
				})
				codec := NewMessageCodec(transportInstance, loggerOption)
				err = codec.Connect()
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
				fields.messageCodec = codec
			},
		},
		{
			name: "setup connection",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
				messageCodec: func() *MessageCodec {
					transport := test.NewTransport()
					transportUrl := url.URL{Scheme: "test"}
					transportInstance, err := transport.CreateTransportInstance(transportUrl, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
						return nil
					}
					type MockState uint8
					const (
						RESET MockState = iota
						APPLICATION_FILTER_1
						APPLICATION_FILTER_2
						INTERFACE_OPTIONS_3
						INTERFACE_OPTIONS_1_PUN
						INTERFACE_OPTIONS_1
						DONE
					)
					currentState := atomic.Value{}
					currentState.Store(RESET)
					transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
						switch currentState.Load().(MockState) {
						case RESET:
							t.Log("Dispatching reset echo")
							transportInstance.FillReadBuffer([]byte("~~~\r"))
							currentState.Store(APPLICATION_FILTER_1)
						case APPLICATION_FILTER_1:
							t.Log("Dispatching app1 echo and confirm")
							transportInstance.FillReadBuffer([]byte("@A32100FF\r"))
							transportInstance.FillReadBuffer([]byte("322100AD\r\n"))
							currentState.Store(APPLICATION_FILTER_2)
						case APPLICATION_FILTER_2:
							t.Log("Dispatching app2 echo and confirm")
							transportInstance.FillReadBuffer([]byte("@A32200FF\r"))
							transportInstance.FillReadBuffer([]byte("322200AC\r\n"))
							currentState.Store(INTERFACE_OPTIONS_3)
						case INTERFACE_OPTIONS_3:
							t.Log("Dispatching interface 3 echo and confirm")
							transportInstance.FillReadBuffer([]byte("@A342000A\r"))
							transportInstance.FillReadBuffer([]byte("3242008C\r\n"))
							currentState.Store(INTERFACE_OPTIONS_1_PUN)
						case INTERFACE_OPTIONS_1_PUN:
							t.Log("Dispatching interface 1 PUN echo and confirm???")
							transportInstance.FillReadBuffer([]byte("@A3410079\r"))
							transportInstance.FillReadBuffer([]byte("3241008D\r\n"))
							currentState.Store(INTERFACE_OPTIONS_1)
						case INTERFACE_OPTIONS_1:
							t.Log("Dispatching interface 1 echo and confirm???")
							transportInstance.FillReadBuffer([]byte("@A3300079\r"))
							transportInstance.FillReadBuffer([]byte("3230009E\r\n"))
							currentState.Store(DONE)
						case DONE:
							t.Log("Done")
						}
					})
					codec := NewMessageCodec(transportInstance)
					err = codec.Connect()
					if err != nil {
						t.Error(err)
						t.FailNow()
						return nil
					}
					return codec
				}(),
			},
			args: args{
				ctx: context.Background(),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)
				fields.log = logger

				testutils.SetToTestingLogger(t, readWriteModel.Plc4xModelLog)

				// Custom option for that
				loggerOption := options.WithCustomLogger(logger)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, loggerOption)

				// Build the message codec
				fields.messageCodec = NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil)
					if err != nil {
						t.Error(err)
						t.FailNow()
					}
					return ti
				}(), loggerOption)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			c.setupConnection(tt.args.ctx, tt.args.ch)
		})
	}
}

func TestConnection_startSubscriptionHandler(t *testing.T) {
	type fields struct {
		DefaultConnection _default.DefaultConnection
		messageCodec      *MessageCodec
		subscribers       []*Subscriber
		tm                transactions.RequestTransactionManager
		configuration     Configuration
		driverContext     DriverContext
		connectionId      string
		tracer            *spi.Tracer
		log               zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		setup  func(t *testing.T, fields *fields)
	}{
		{
			name: "just start",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
		},
		{
			name: "just start and feed (no subs)",
			fields: fields{
				DefaultConnection: func() _default.DefaultConnection {
					defaultConnection := _default.NewDefaultConnection(nil)
					defaultConnection.SetConnected(true)
					return defaultConnection
				}(),
				messageCodec: func() *MessageCodec {
					messageCodec := NewMessageCodec(nil)
					go func() {
						messageCodec.monitoredMMIs <- nil
						messageCodec.monitoredSALs <- nil
					}()
					return messageCodec
				}(),
			},
		},
		{
			name: "just start and feed",
			fields: fields{
				DefaultConnection: func() _default.DefaultConnection {
					defaultConnection := _default.NewDefaultConnection(nil)
					defaultConnection.SetConnected(true)
					return defaultConnection
				}(),
				messageCodec: func() *MessageCodec {
					messageCodec := NewMessageCodec(nil)
					go func() {
						messageCodec.monitoredMMIs <- readWriteModel.NewCALReplyShort(0, nil, nil, nil)
						messageCodec.monitoredSALs <- readWriteModel.NewMonitoredSAL(0, nil)
					}()
					return messageCodec
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				fields.subscribers = []*Subscriber{NewSubscriber(nil, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               tt.fields.log,
			}
			c.startSubscriptionHandler()
			time.Sleep(50 * time.Millisecond)
		})
	}
}

func TestNewConnection(t *testing.T) {
	type args struct {
		messageCodec  *MessageCodec
		configuration Configuration
		driverContext DriverContext
		tagHandler    spi.PlcTagHandler
		tm            transactions.RequestTransactionManager
		options       map[string][]string
	}
	tests := []struct {
		name       string
		args       args
		wantAssert func(*testing.T, *Connection) bool
	}{
		{
			name: "just create the connection",
			wantAssert: func(t *testing.T, connection *Connection) bool {
				return assert.NotNil(t, connection)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.wantAssert(t, NewConnection(tt.args.messageCodec, tt.args.configuration, tt.args.driverContext, tt.args.tagHandler, tt.args.tm, tt.args.options)), "NewConnection(%v, %v, %v, %v, %v, %v)", tt.args.messageCodec, tt.args.configuration, tt.args.driverContext, tt.args.tagHandler, tt.args.tm, tt.args.options)
		})
	}
}
