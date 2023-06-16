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
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/stretchr/testify/assert"
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Build the default connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
			}
			assert.True(t, tt.wantAsserter(t, c.ConnectWithContext(tt.args.ctx)), "ConnectWithContext(%v)", tt.args.ctx)
			// To shut down properly we always do that
			c.SetConnected(false)
			c.handlerWaitGroup.Wait()
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
		tracer            tracer.Tracer
	}
	tests := []struct {
		name         string
		fields       fields
		setup        func(t *testing.T, fields *fields)
		wantAsserter func(t *testing.T, connection plc4go.PlcConnection) bool
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
	}
	tests := []struct {
		name   string
		fields fields
		want   tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "a string",
			fields: fields{
				DefaultConnection: _default.NewDefaultConnection(nil),
			},
			want: `
╔═Connection══════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═defaultConnection═══════╗╔═alphaGenerator════════════════════╗╔═messageCodec╗                          ║
║║╔═defaultTtl╗╔═connected╗║║╔═AlphaGenerator/currentAlpha═════╗║║    <nil>    ║                          ║
║║║    10s    ║║ b0 false ║║║║            0x67 'g'             ║║╚═════════════╝                          ║
║║╚═══════════╝╚══════════╝║║╚═════════════════════════════════╝║                                         ║
║╚═════════════════════════╝╚═══════════════════════════════════╝                                         ║
║╔═configuration═════════════════════════════════════════════════════════════════════════════════════════╗║
║║╔═Configuration═══════════════════════════════════════════════════════════════════════════════════════╗║║
║║║╔═srchk══╗╔═exstat═╗╔═pun════╗╔═localSal╗╔═pcn════╗╔═idmon══╗╔═monitor╗╔═smart══╗╔═xonXoff╗╔═connect╗║║║
║║║║b0 false║║b0 false║║b0 false║║b0 false ║║b0 false║║b0 false║║b0 false║║b0 false║║b0 false║║b0 false║║║║
║║║╚════════╝╚════════╝╚════════╝╚═════════╝╚════════╝╚════════╝╚════════╝╚════════╝╚════════╝╚════════╝║║║
║║║╔═monitoredApplication1╗╔═monitoredApplication2╗                                                     ║║║
║║║║       0x00 '.'       ║║       0x00 '.'       ║                                                     ║║║
║║║╚══════════════════════╝╚══════════════════════╝                                                     ║║║
║║╚═════════════════════════════════════════════════════════════════════════════════════════════════════╝║║
║╚═══════════════════════════════════════════════════════════════════════════════════════════════════════╝║
║╔═driverContext═══════════════════════════════════╗                                                      ║
║║╔═DriverContext═════════════════════════════════╗║                                                      ║
║║║╔═awaitSetupComplete╗╔═awaitDisconnectComplete╗║║                                                      ║
║║║║     b0 false      ║║        b0 false        ║║║                                                      ║
║║║╚═══════════════════╝╚════════════════════════╝║║                                                      ║
║║╚═══════════════════════════════════════════════╝║                                                      ║
║╚═════════════════════════════════════════════════╝                                                      ║
╚═════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:], // TODO: configuration is not redered right now
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				DefaultConnection: tt.fields.DefaultConnection,
				alphaGenerator:    AlphaGenerator{currentAlpha: 'g'},
				messageCodec:      tt.fields.messageCodec,
				subscribers:       tt.fields.subscribers,
				tm:                tt.fields.tm,
				configuration:     tt.fields.configuration,
				driverContext:     tt.fields.driverContext,
				connectionId:      tt.fields.connectionId,
				tracer:            tt.fields.tracer,
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
	}
	type args struct {
		err error
		ch  chan<- plc4go.PlcConnectionConnectResult
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		setup         func(t *testing.T, fields *fields, args *args)
		chanValidator func(*testing.T, chan<- plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "instant connect",
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			chanValidator: func(_ *testing.T, _ chan<- plc4go.PlcConnectionConnectResult) bool {
				return true
			},
		},
		{
			name: "notified connect",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete: true,
				},
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
		setup  func(t *testing.T, fields *fields, args *args)
		wantOk bool
	}{
		{
			name: "send reset",
			args: args{
				ctx: testutils.TestContext(t),
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
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
				ctx: testutils.TestContext(t),
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
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Setup connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
			args: args{
				ctx: testutils.TestContext(t),
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
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Setup connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
			args: args{
				ctx: testutils.TestContext(t),
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
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Setup connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
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
		tracer            tracer.Tracer
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
			args: args{
				ctx: testutils.TestContext(t),
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
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Setup connection
				fields.DefaultConnection = _default.NewDefaultConnection(nil, _options...)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
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
				log:               testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.wantOk, c.setInterfaceOptions3(tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions), "setInterfaceOptions3(%v, %v, %v, %v)", tt.args.ctx, tt.args.ch, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setupConnection(t *testing.T) {
	type fields struct {
		_DefaultConnection _default.DefaultConnection
		messageCodec       *MessageCodec
		subscribers        []*Subscriber
		tm                 transactions.RequestTransactionManager
		configuration      Configuration
		driverContext      DriverContext
		connectionId       string
		tracer             tracer.Tracer
	}
	type args struct {
		ctx context.Context
		ch  chan plc4go.PlcConnectionConnectResult
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		setup       func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, connection *Connection)
		validator   func(t *testing.T, result plc4go.PlcConnectionConnectResult)
	}{
		{
			name: "setup connection (failing)",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
			},
			args: args{
				ctx: testutils.TestContext(t),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect())
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
			validator: func(t *testing.T, result plc4go.PlcConnectionConnectResult) {
				assert.NotNil(t, result)
				assert.Error(t, result.GetErr())
			},
		},
		{
			name: "setup connection (failing after reset)",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
			},
			args: args{
				ctx: testutils.TestContext(t),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Build the message codec
				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				ti, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
				type MockState uint8
				const (
					RESET MockState = iota
					DONE
				)
				currentState := atomic.Value{}
				currentState.Store(RESET)
				stateChangeMutex := sync.Mutex{}
				ti.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					stateChangeMutex.Lock()
					defer stateChangeMutex.Unlock()
					switch currentState.Load().(MockState) {
					case RESET:
						t.Log("Dispatching reset echo")
						transportInstance.FillReadBuffer([]byte("~~~\r"))
						currentState.Store(DONE)
					case DONE:
						t.Log("Done")
					}
				})
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect())
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})

				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
			validator: func(t *testing.T, result plc4go.PlcConnectionConnectResult) {
				assert.NotNil(t, result)
				assert.Error(t, result.GetErr())
			},
		},
		{
			name: "setup connection (failing after app filters)",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
			},
			args: args{
				ctx: testutils.TestContext(t),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				// Build the message codec
				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				ti, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
				type MockState uint8
				const (
					RESET MockState = iota
					APPLICATION_FILTER_1
					APPLICATION_FILTER_2
					DONE
				)
				currentState := atomic.Value{}
				currentState.Store(RESET)
				stateChangeMutex := sync.Mutex{}
				ti.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					stateChangeMutex.Lock()
					defer stateChangeMutex.Unlock()
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
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect())
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})

				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
			validator: func(t *testing.T, result plc4go.PlcConnectionConnectResult) {
				assert.NotNil(t, result)
				assert.Error(t, result.GetErr())
			},
		},
		{
			name: "setup connection (failing after interface options 3",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
			},
			args: args{
				ctx: testutils.TestContext(t),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Build the message codec
				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				ti, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
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
				stateChangeMutex := sync.Mutex{}
				ti.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					t.Logf("Reacting to\n%s", hex.Dump(data))
					stateChangeMutex.Lock()
					defer stateChangeMutex.Unlock()
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
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect())
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})

				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
			validator: func(t *testing.T, result plc4go.PlcConnectionConnectResult) {
				assert.NotNil(t, result)
				assert.Error(t, result.GetErr())
			},
		},
		{
			name: "setup connection (failing after interface options 1 pun)",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
			},
			args: args{
				ctx: testutils.TestContext(t),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Build the message codec
				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				ti, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
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
				stateChangeMutex := sync.Mutex{}
				ti.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					stateChangeMutex.Lock()
					defer stateChangeMutex.Unlock()
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
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect())
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})

				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
			validator: func(t *testing.T, result plc4go.PlcConnectionConnectResult) {
				assert.NotNil(t, result)
				assert.Error(t, result.GetErr())
			},
		},
		{
			name: "setup connection",
			fields: fields{
				driverContext: DriverContext{
					awaitSetupComplete:      true,
					awaitDisconnectComplete: true,
				},
			},
			args: args{
				ctx: testutils.TestContext(t),
				ch:  make(chan plc4go.PlcConnectionConnectResult, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Build the message codec
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)

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
				stateChangeMutex := sync.Mutex{}
				ti.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					stateChangeMutex.Lock()
					defer stateChangeMutex.Unlock()
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
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect())
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
			validator: func(t *testing.T, result plc4go.PlcConnectionConnectResult) {
				assert.NotNil(t, result)
				assert.NoError(t, result.GetErr())
				assert.NotNil(t, result.GetConnection())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: tt.fields.driverContext,
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, c)
			}
			c.setupConnection(tt.args.ctx, tt.args.ch)
			assert.NotNil(t, tt.args.ch, "We always need a result channel")
			chanTimeout := time.NewTimer(10 * time.Second)
			t.Cleanup(func() {
				utils.CleanupTimer(chanTimeout)
			})
			select {
			case <-chanTimeout.C:
				t.Fatal("setup connection doesn't fill chan in time")
			case result := <-tt.args.ch:
				if tt.validator != nil {
					tt.validator(t, result)
				}
			}
			// To shut down properly we always do that
			closeTimeout := time.NewTimer(10 * time.Second)
			t.Cleanup(func() {
				utils.CleanupTimer(closeTimeout)
			})
			select {
			case <-closeTimeout.C:
				t.Fatal("close didn't react in time")
			case <-c.Close():
				t.Log("connection closed")
			}
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
		tracer            tracer.Tracer
	}
	tests := []struct {
		name        string
		fields      fields
		setup       func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, connection *Connection)
	}{
		{
			name: "just start",
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				connection.DefaultConnection = _default.NewDefaultConnection(connection, _options...)
			},
		},
		{
			name: "just start and feed (no subs)",
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				codec := NewMessageCodec(nil, _options...)
				codec.monitoredMMIs = make(chan readWriteModel.CALReply, 1)
				codec.monitoredSALs = make(chan readWriteModel.MonitoredSAL, 1)
				dispatchWg := sync.WaitGroup{}
				dispatchWg.Add(1)
				t.Cleanup(dispatchWg.Wait)
				go func() {
					defer dispatchWg.Done()
					codec.monitoredMMIs <- readWriteModel.NewCALReplyShort(0, nil, nil, nil)
					codec.monitoredSALs <- readWriteModel.NewMonitoredSAL(0, nil)
				}()
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				defaultConnection := _default.NewDefaultConnection(connection, _options...)
				defaultConnection.SetConnected(true)
				connection.DefaultConnection = defaultConnection
			},
		},
		{
			name: "just start and feed",
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				fields.subscribers = []*Subscriber{NewSubscriber(nil, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))}
				codec := NewMessageCodec(nil, _options...)
				written := make(chan struct{})
				dispatchWg := sync.WaitGroup{}
				dispatchWg.Add(1)
				t.Cleanup(dispatchWg.Wait)
				go func() {
					defer dispatchWg.Done()
					codec.monitoredMMIs <- readWriteModel.NewCALReplyShort(0, nil, nil, nil)
					codec.monitoredSALs <- readWriteModel.NewMonitoredSAL(0, nil)
					close(written)
				}()
				t.Cleanup(func() {
					<-written
				})
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				defaultConnection := _default.NewDefaultConnection(connection, _options...)
				defaultConnection.SetConnected(true)
				connection.DefaultConnection = defaultConnection
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
				log:               testutils.ProduceTestingLogger(t),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, c)
			}
			c.startSubscriptionHandler()
			// To shut down properly we always do that
			c.SetConnected(false)
			c.handlerWaitGroup.Wait()
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
		setup      func(t *testing.T, args *args)
		wantAssert func(*testing.T, *Connection) bool
	}{
		{
			name: "just create the connection",
			setup: func(t *testing.T, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				codec := NewMessageCodec(test.NewTransportInstance(transport, _options...), _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				args.messageCodec = codec
			},
			wantAssert: func(t *testing.T, connection *Connection) bool {
				return assert.NotNil(t, connection)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
			connection := NewConnection(tt.args.messageCodec, tt.args.configuration, tt.args.driverContext, tt.args.tagHandler, tt.args.tm, tt.args.options)
			t.Cleanup(func() {
				timer := time.NewTimer(1 * time.Second)
				t.Cleanup(func() {
					utils.CleanupTimer(timer)
				})
				select {
				case <-connection.Close():
				case <-timer.C:
					t.Error("timeout")
				}
			})
			assert.True(t, tt.wantAssert(t, connection), "NewConnection(%v, %v, %v, %v, %v, %v)", tt.args.messageCodec, tt.args.configuration, tt.args.driverContext, tt.args.tagHandler, tt.args.tm, tt.args.options)
		})
	}
}
