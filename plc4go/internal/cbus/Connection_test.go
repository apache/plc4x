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
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
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
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcBrowseRequestBuilder) bool
	}{
		{
			name: "return not nil",
			wantAssert: func(t *testing.T, builder apiModel.PlcBrowseRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.True(t, tt.wantAssert(t, c.BrowseRequestBuilder()), "BrowseRequestBuilder()")
		})
	}
}

func TestConnection_Connect(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(*testing.T, *fields, *args)
		wantErr assert.ErrorAssertionFunc
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
				connectionId: "connectionId13",
				tracer:       nil,
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = t.Context()
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
		// TODO: add error case for failing messageCodec connect
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			err := c.Connect(tt.args.ctx)
			assert.True(t, tt.wantErr(t, err), "Connect(%v)", tt.args.ctx)
			// To shut down properly we always do that
			c.SetConnected(false)
			c.handlerWaitGroup.Wait()
		})
	}
}

func TestConnection_GetConnection(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantAsserter(t, c.GetConnection()), "GetConnection()")
		})
	}
}

func TestConnection_GetConnectionId(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Equalf(t, tt.want, c.GetConnectionId(), "GetConnectionId()")
		})
	}
}

func TestConnection_GetMessageCodec(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Equalf(t, tt.want, c.GetMessageCodec(), "GetMessageCodec()")
		})
	}
}

func TestConnection_GetMetadata(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcConnectionMetadata
	}{
		{
			name: "give metadata",
			want: &_default.DefaultConnectionMetadata{
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
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Equalf(t, tt.want, c.GetMetadata(), "GetMetadata()")
		})
	}
}

func TestConnection_GetTracer(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Equalf(t, tt.want, c.GetTracer(), "GetTracer()")
		})
	}
}

func TestConnection_IsTraceEnabled(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Equalf(t, tt.want, c.IsTraceEnabled(), "IsTraceEnabled()")
		})
	}
}

func TestConnection_ReadRequestBuilder(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcReadRequestBuilder) bool
	}{
		{
			name: "return not nil",
			wantAssert: func(t *testing.T, builder apiModel.PlcReadRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantAssert(t, c.ReadRequestBuilder()), "ReadRequestBuilder()")
		})
	}
}

func TestConnection_SubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcSubscriptionRequestBuilder) bool
	}{
		{
			name: "return not nil",
			wantAssert: func(t *testing.T, builder apiModel.PlcSubscriptionRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantAssert(t, c.SubscriptionRequestBuilder()), "SubscriptionRequestBuilder()")
		})
	}
}

func TestConnection_UnsubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcUnsubscriptionRequestBuilder
	}{
		{
			name: "create one",
			want: spiModel.NewDefaultPlcUnsubscriptionRequestBuilder(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Equalf(t, tt.want, c.UnsubscriptionRequestBuilder(), "UnsubscriptionRequestBuilder()")
		})
	}
}

func TestConnection_WriteRequestBuilder(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	tests := []struct {
		name       string
		fields     fields
		wantAssert func(*testing.T, apiModel.PlcWriteRequestBuilder) bool
	}{
		{
			name: "return not nil",
			wantAssert: func(t *testing.T, builder apiModel.PlcWriteRequestBuilder) bool {
				return assert.NotNil(t, builder)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantAssert(t, c.WriteRequestBuilder()), "WriteRequestBuilder()")
		})
	}
}

func TestConnection_addSubscriber(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
			name: "existing subscriber should not be added",
			args: args{subscriber: theOneSubscriber},
			subElevator: func(t *testing.T, subscribers []*Subscriber) bool {
				return len(subscribers) == 1
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			c.addSubscriber(tt.args.subscriber)
			assert.Truef(t, tt.subElevator(t, c.subscribers), "addSubscriber(%v)", tt.args.subscriber)
		})
	}
}

func TestConnection_sendCalDataWrite(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	type args struct {
		ctx            context.Context
		paramNo        readWriteModel.Parameter
		parameterValue readWriteModel.ParameterValue
		requestContext *readWriteModel.RequestContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "send something",
			args: args{
				paramNo:        readWriteModel.Parameter_APPLICATION_ADDRESS_2,
				parameterValue: readWriteModel.NewParameterValueApplicationAddress2(readWriteModel.NewApplicationAddress2(1), nil),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
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
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantErr(t, c.sendCalDataWrite(tt.args.ctx, tt.args.paramNo, tt.args.parameterValue, tt.args.requestContext)), "sendCalDataWrite(%v, %v, %v, %v, %v, %v)", tt.args.ctx, tt.args.paramNo, tt.args.parameterValue, tt.args.requestContext)
		})
	}
}

func TestConnection_sendReset(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
			name: "send reset",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
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
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantErr(t, c.sendReset(tt.args.ctx)), "sendReset(%v)", tt.args.ctx)
		})
	}
}

func TestConnection_setApplicationFilter(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	type args struct {
		ctx            context.Context
		requestContext *readWriteModel.RequestContext
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "set application filter (failing)",
			args: args{
				ctx: t.Context(),
				requestContext: func() *readWriteModel.RequestContext {
					var requestContext readWriteModel.RequestContext = readWriteModel.NewRequestContext(false)
					return &requestContext
				}(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Setup connection
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
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
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantErr(t, c.setApplicationFilter(tt.args.ctx, tt.args.requestContext)), "setApplicationFilter(%v, %v)", tt.args.ctx, tt.args.requestContext)
		})
	}
}

func TestConnection_setInterface1PowerUpSettings(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	type args struct {
		ctx            context.Context
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "set interface 1 PUN options (failing)",
			args: args{
				ctx: t.Context(),
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
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					t.Log("disconnecting codec")
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
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
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantErr(t, c.setInterface1PowerUpSettings(tt.args.ctx, tt.args.requestContext, tt.args.cbusOptions)), "setInterface1PowerUpSettings(%v, %v, %v)", tt.args.ctx, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setInterfaceOptions1(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	type args struct {
		ctx            context.Context
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "set interface 1 options (failing)",
			args: args{
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
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
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
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantErr(t, c.setInterfaceOptions1(tt.args.ctx, tt.args.requestContext, tt.args.cbusOptions)), "setInterfaceOptions1(%v, %v, %v)", tt.args.ctx, tt.args.requestContext, tt.args.cbusOptions)
		})
	}
}

func TestConnection_setInterfaceOptions3(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
	}
	type args struct {
		ctx            context.Context
		requestContext *readWriteModel.RequestContext
		cbusOptions    *readWriteModel.CBusOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "set interface 3 options (failing)",
			args: args{
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
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
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
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
			assert.Truef(t, tt.wantErr(t, c.setInterfaceOptions3(tt.args.ctx, tt.args.requestContext, tt.args.cbusOptions)), "setInterfaceOptions3(%v, %v, %v)", tt.args.ctx, tt.args.requestContext, tt.args.cbusOptions)
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
		connectionId       string
		tracer             tracer.Tracer
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
			name: "setup connection (failing)",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
				require.NoError(t, err)
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 20*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
		{
			name: "setup connection (failing after reset)",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
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
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})

				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
		{
			name: "setup connection (failing after app filters)",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
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
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})

				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
		{
			name: "setup connection (failing after interface options 3",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
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
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})

				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
		{
			name: "setup connection (failing after interface options 1 pun)",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
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
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})

				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 2*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.Error,
		},
		{
			name: "setup connection",
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				// Build the message codec
				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

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
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec

				args.ctx = testutils.TestContext(t)
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 20*time.Second)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				if tt.setup != nil {
					tt.setup(t, &tt.fields, &tt.args)
				}
				c := &Connection{
					messageCodec:  tt.fields.messageCodec,
					subscribers:   tt.fields.subscribers,
					tm:            tt.fields.tm,
					configuration: tt.fields.configuration,
					driverContext: driverContextForTesting(),
					connectionId:  tt.fields.connectionId,
					tracer:        tt.fields.tracer,
					log:           testutils.ProduceTestingLogger(t),
				}
				c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
				err := c.setupConnection(tt.args.ctx)
				tt.wantErr(t, err)
			})
		})
	}
}

func TestConnection_startSubscriptionHandler(t *testing.T) {
	type fields struct {
		messageCodec  *MessageCodec
		subscribers   []*Subscriber
		tm            transactions.RequestTransactionManager
		configuration Configuration
		connectionId  string
		tracer        tracer.Tracer
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
				dispatchWg := new(sync.WaitGroup)
				t.Cleanup(dispatchWg.Wait)
				dispatchWg.Go(func() {
					codec.monitoredMMIs <- readWriteModel.NewCALReplyShort(0, nil)
					codec.monitoredSALs <- readWriteModel.NewMonitoredSALShortFormBasicMode(
						0,
						0,
						nil,
						nil,
						nil,
						readWriteModel.ApplicationIdContainer_ACCESS_CONTROL_D5,
						nil,
					)
				})
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				connection.SetConnected(true)
			},
		},
		{
			name: "just start and feed",
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				fields.subscribers = []*Subscriber{NewSubscriber(nil, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))}
				codec := NewMessageCodec(nil, _options...)
				written := make(chan struct{})
				dispatchWg := new(sync.WaitGroup)
				t.Cleanup(dispatchWg.Wait)
				dispatchWg.Go(func() {
					codec.monitoredMMIs <- readWriteModel.NewCALReplyShort(0, nil)
					codec.monitoredSALs <- readWriteModel.NewMonitoredSALShortFormBasicMode(
						0,
						0,
						nil,
						nil,
						nil,
						readWriteModel.ApplicationIdContainer_ACCESS_CONTROL_D5,
						nil,
					)
					close(written)
				})
				t.Cleanup(func() {
					<-written
				})
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.messageCodec = codec
			},
			manipulator: func(t *testing.T, connection *Connection) {
				connection.SetConnected(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				messageCodec:  tt.fields.messageCodec,
				subscribers:   tt.fields.subscribers,
				tm:            tt.fields.tm,
				configuration: tt.fields.configuration,
				driverContext: driverContextForTesting(),
				connectionId:  tt.fields.connectionId,
				tracer:        tt.fields.tracer,
				log:           testutils.ProduceTestingLogger(t),
			}
			c.DefaultConnection = _default.NewDefaultConnection(c, testutils.EnrichOptionsWithOptionsForTesting(t)...)
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
		_options      []options.WithOption
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
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				args.messageCodec = codec

				args._options = _options
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
			connection := NewConnection(tt.args.messageCodec, tt.args.configuration, tt.args.driverContext, tt.args.tagHandler, tt.args.tm, tt.args.options, tt.args._options...)
			t.Cleanup(func() {
				t.Log("Disconnecting connection", connection.Close())
			})
			assert.True(t, tt.wantAssert(t, connection), "NewConnection(%v, %v, %v, %v, %v, %v)", tt.args.messageCodec, tt.args.configuration, tt.args.driverContext, tt.args.tagHandler, tt.args.tm, tt.args.options)
		})
	}
}
