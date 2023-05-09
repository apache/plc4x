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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

func TestDefaultConnectionMetadata_CanBrowse(t *testing.T) {
	type fields struct {
		ConnectionAttributes map[string]string
		ProvidesReading      bool
		ProvidesWriting      bool
		ProvidesSubscribing  bool
		ProvidesBrowsing     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "can't browse",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultConnectionMetadata{
				ConnectionAttributes: tt.fields.ConnectionAttributes,
				ProvidesReading:      tt.fields.ProvidesReading,
				ProvidesWriting:      tt.fields.ProvidesWriting,
				ProvidesSubscribing:  tt.fields.ProvidesSubscribing,
				ProvidesBrowsing:     tt.fields.ProvidesBrowsing,
			}
			assert.Equalf(t, tt.want, m.CanBrowse(), "CanBrowse()")
		})
	}
}

func TestDefaultConnectionMetadata_CanRead(t *testing.T) {
	type fields struct {
		ConnectionAttributes map[string]string
		ProvidesReading      bool
		ProvidesWriting      bool
		ProvidesSubscribing  bool
		ProvidesBrowsing     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "can't read",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultConnectionMetadata{
				ConnectionAttributes: tt.fields.ConnectionAttributes,
				ProvidesReading:      tt.fields.ProvidesReading,
				ProvidesWriting:      tt.fields.ProvidesWriting,
				ProvidesSubscribing:  tt.fields.ProvidesSubscribing,
				ProvidesBrowsing:     tt.fields.ProvidesBrowsing,
			}
			assert.Equalf(t, tt.want, m.CanRead(), "CanRead()")
		})
	}
}

func TestDefaultConnectionMetadata_CanSubscribe(t *testing.T) {
	type fields struct {
		ConnectionAttributes map[string]string
		ProvidesReading      bool
		ProvidesWriting      bool
		ProvidesSubscribing  bool
		ProvidesBrowsing     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "can't subscribe",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultConnectionMetadata{
				ConnectionAttributes: tt.fields.ConnectionAttributes,
				ProvidesReading:      tt.fields.ProvidesReading,
				ProvidesWriting:      tt.fields.ProvidesWriting,
				ProvidesSubscribing:  tt.fields.ProvidesSubscribing,
				ProvidesBrowsing:     tt.fields.ProvidesBrowsing,
			}
			assert.Equalf(t, tt.want, m.CanSubscribe(), "CanSubscribe()")
		})
	}
}

func TestDefaultConnectionMetadata_CanWrite(t *testing.T) {
	type fields struct {
		ConnectionAttributes map[string]string
		ProvidesReading      bool
		ProvidesWriting      bool
		ProvidesSubscribing  bool
		ProvidesBrowsing     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "can't write",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultConnectionMetadata{
				ConnectionAttributes: tt.fields.ConnectionAttributes,
				ProvidesReading:      tt.fields.ProvidesReading,
				ProvidesWriting:      tt.fields.ProvidesWriting,
				ProvidesSubscribing:  tt.fields.ProvidesSubscribing,
				ProvidesBrowsing:     tt.fields.ProvidesBrowsing,
			}
			assert.Equalf(t, tt.want, m.CanWrite(), "CanWrite()")
		})
	}
}

func TestDefaultConnectionMetadata_GetConnectionAttributes(t *testing.T) {
	type fields struct {
		ConnectionAttributes map[string]string
		ProvidesReading      bool
		ProvidesWriting      bool
		ProvidesSubscribing  bool
		ProvidesBrowsing     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultConnectionMetadata{
				ConnectionAttributes: tt.fields.ConnectionAttributes,
				ProvidesReading:      tt.fields.ProvidesReading,
				ProvidesWriting:      tt.fields.ProvidesWriting,
				ProvidesSubscribing:  tt.fields.ProvidesSubscribing,
				ProvidesBrowsing:     tt.fields.ProvidesBrowsing,
			}
			assert.Equalf(t, tt.want, m.GetConnectionAttributes(), "GetConnectionAttributes()")
		})
	}
}

func TestNewDefaultConnection(t *testing.T) {
	type args struct {
		requirements DefaultConnectionRequirements
		options      []options.WithOption
	}
	tests := []struct {
		name string
		args args
		want DefaultConnection
	}{
		{
			name: "just create it",
			want: &defaultConnection{
				defaultTtl: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultConnection(tt.args.requirements, tt.args.options...), "NewDefaultConnection(%v, %v)", tt.args.requirements, tt.args.options)
		})
	}
}

func TestNewDefaultPlcConnectionCloseResult(t *testing.T) {
	type args struct {
		connection plc4go.PlcConnection
		err        error
	}
	tests := []struct {
		name string
		args args
		want plc4go.PlcConnectionCloseResult
	}{
		{
			name: "create it",
			want: &plcConnectionCloseResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcConnectionCloseResult(tt.args.connection, tt.args.err), "NewDefaultPlcConnectionCloseResult(%v, %v)", tt.args.connection, tt.args.err)
		})
	}
}

func TestNewDefaultPlcConnectionCloseResultWithTraces(t *testing.T) {
	type args struct {
		connection plc4go.PlcConnection
		err        error
		traces     []spi.TraceEntry
	}
	tests := []struct {
		name string
		args args
		want plc4go.PlcConnectionCloseResult
	}{
		{
			name: "create it",
			want: &plcConnectionCloseResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcConnectionCloseResultWithTraces(tt.args.connection, tt.args.err, tt.args.traces), "NewDefaultPlcConnectionCloseResultWithTraces(%v, %v, %v)", tt.args.connection, tt.args.err, tt.args.traces)
		})
	}
}

func TestNewDefaultPlcConnectionConnectResult(t *testing.T) {
	type args struct {
		connection plc4go.PlcConnection
		err        error
	}
	tests := []struct {
		name string
		args args
		want DefaultPlcConnectionConnectResult
	}{
		{
			name: "create it",
			want: &plcConnectionConnectResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcConnectionConnectResult(tt.args.connection, tt.args.err), "NewDefaultPlcConnectionConnectResult(%v, %v)", tt.args.connection, tt.args.err)
		})
	}
}

func TestNewDefaultPlcConnectionPingResult(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want plc4go.PlcConnectionPingResult
	}{
		{
			name: "create it",
			want: &plcConnectionPingResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcConnectionPingResult(tt.args.err), "NewDefaultPlcConnectionPingResult(%v)", tt.args.err)
		})
	}
}

func TestWithDefaultTtl(t *testing.T) {
	type args struct {
		defaultTtl time.Duration
	}
	tests := []struct {
		name string
		args args
		want options.WithOption
	}{
		{
			name: "create it",
			want: withDefaultTtl{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WithDefaultTtl(tt.args.defaultTtl), "WithDefaultTtl(%v)", tt.args.defaultTtl)
		})
	}
}

func TestWithPlcTagHandler(t *testing.T) {
	type args struct {
		tagHandler spi.PlcTagHandler
	}
	tests := []struct {
		name string
		args args
		want options.WithOption
	}{
		{
			name: "create it",
			want: withPlcTagHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WithPlcTagHandler(tt.args.tagHandler), "WithPlcTagHandler(%v)", tt.args.tagHandler)
		})
	}
}

func TestWithPlcValueHandler(t *testing.T) {
	type args struct {
		plcValueHandler spi.PlcValueHandler
	}
	tests := []struct {
		name string
		args args
		want options.WithOption
	}{
		{
			name: "create it",
			want: withPlcValueHandler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, WithPlcValueHandler(tt.args.plcValueHandler), "WithPlcValueHandler(%v)", tt.args.plcValueHandler)
		})
	}
}

func Test_buildDefaultConnection(t *testing.T) {
	type args struct {
		requirements DefaultConnectionRequirements
		options      []options.WithOption
	}
	tests := []struct {
		name string
		args args
		want DefaultConnection
	}{
		{
			name: "build it",
			want: &defaultConnection{
				defaultTtl: 10 * time.Second,
			},
		},
		{
			name: "build it with ttl",
			args: args{
				options: []options.WithOption{
					withDefaultTtl{},
				},
			},
			want: &defaultConnection{},
		},
		{
			name: "build it with plc tag handler",
			args: args{
				options: []options.WithOption{
					withPlcTagHandler{},
				},
			},
			want: &defaultConnection{
				defaultTtl: 10 * time.Second,
			},
		},
		{
			name: "build it with plc value handler",
			args: args{
				options: []options.WithOption{
					withPlcValueHandler{},
				},
			},
			want: &defaultConnection{
				defaultTtl: 10 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, buildDefaultConnection(tt.args.requirements, tt.args.options...), "buildDefaultConnection(%v, %v)", tt.args.requirements, tt.args.options)
		})
	}
}

func Test_defaultConnection_BlockingClose(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name      string
		fields    fields
		mockSetup func(t *testing.T, fields *fields)
	}{
		{
			name: "close",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultConnectionRequirements(t)
				connection := NewMockPlcConnection(t)
				connection.EXPECT().Close().Return(nil)
				requirements.EXPECT().GetConnection().Return(connection)
				fields.DefaultConnectionRequirements = requirements
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			d.BlockingClose()
		})
	}
}

func Test_defaultConnection_BrowseRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   model.PlcBrowseRequestBuilder
	}{
		{
			name: "create it",
			fields: fields{
				DefaultConnectionRequirements: NewMockDefaultConnectionRequirements(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "not provided by actual connection", err)
				} else {
					t.Error("should fail")
				}
			}()
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.BrowseRequestBuilder(), "BrowseRequestBuilder()")
		})
	}
}

func Test_defaultConnection_Close(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name         string
		fields       fields
		mockSetup    func(t *testing.T, fields *fields)
		wantAsserter func(t *testing.T, results <-chan plc4go.PlcConnectionCloseResult) bool
	}{
		{
			name: "close it",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultConnectionRequirements(t)
				codec := NewMockMessageCodec(t)
				{
					expect := codec.EXPECT()
					expect.Disconnect().Return(nil)
					instance := NewMockTransportInstance(t)
					instance.EXPECT().Close().Return(nil)
					expect.GetTransportInstance().Return(instance)
				}
				{
					expect := requirements.EXPECT()
					expect.GetMessageCodec().Return(codec)
					expect.GetConnection().Return(nil)
				}
				fields.DefaultConnectionRequirements = requirements
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionCloseResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.Nil(t, result.GetErr())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Truef(t, tt.wantAsserter(t, d.Close()), "Close()")
		})
	}
}

func Test_defaultConnection_Connect(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name         string
		fields       fields
		mockSetup    func(t *testing.T, fields *fields)
		wantAsserter func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "connect it",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultConnectionRequirements(t)
				results := make(chan plc4go.PlcConnectionConnectResult, 1)
				results <- NewMockPlcConnectionConnectResult(t)
				expect := requirements.EXPECT()
				expect.ConnectWithContext(mock.Anything).Return(results)
				fields.DefaultConnectionRequirements = requirements
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				// Delegated call is tested below
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Truef(t, tt.wantAsserter(t, d.Connect()), "Connect()")
		})
	}
}

func Test_defaultConnection_ConnectWithContext(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		mockSetup    func(t *testing.T, fields *fields, args *args)
		wantAsserter func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "connect it",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				requirements := NewMockDefaultConnectionRequirements(t)
				codec := NewMockMessageCodec(t)
				{
					codec.EXPECT().ConnectWithContext(mock.Anything).Return(nil)
				}
				expect := requirements.EXPECT()
				expect.GetMessageCodec().Return(codec)
				expect.GetConnection().Return(NewMockPlcConnection(t))
				fields.DefaultConnectionRequirements = requirements
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.Nil(t, result.GetErr())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Truef(t, tt.wantAsserter(t, d.ConnectWithContext(tt.args.ctx)), "ConnectWithContext(%v)", tt.args.ctx)
		})
	}
}

func Test_defaultConnection_GetMetadata(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   model.PlcConnectionMetadata
	}{
		{
			name: "get it",
			want: DefaultConnectionMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.GetMetadata(), "GetMetadata()")
		})
	}
}

func Test_defaultConnection_GetPlcTagHandler(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.PlcTagHandler
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.GetPlcTagHandler(), "GetPlcTagHandler()")
		})
	}
}

func Test_defaultConnection_GetPlcValueHandler(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.PlcValueHandler
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.GetPlcValueHandler(), "GetPlcValueHandler()")
		})
	}
}

func Test_defaultConnection_GetTransportInstance(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	theInstance := NewMockTransportInstance(t)
	tests := []struct {
		name      string
		fields    fields
		mockSetup func(t *testing.T, fields *fields)
		want      transports.TransportInstance
	}{
		{
			name: "get it",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultConnectionRequirements(t)
				codec := NewMockMessageCodec(t)
				{
					codec.EXPECT().GetTransportInstance().Return(theInstance)
				}
				requirements.EXPECT().GetMessageCodec().Return(codec)
				fields.DefaultConnectionRequirements = requirements
			},
			want: theInstance,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.GetTransportInstance(), "GetTransportInstance()")
		})
	}
}

func Test_defaultConnection_GetTtl(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.GetTtl(), "GetTtl()")
		})
	}
}

func Test_defaultConnection_IsConnected(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
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
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.IsConnected(), "IsConnected()")
		})
	}
}

func Test_defaultConnection_Ping(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name         string
		fields       fields
		mockSetup    func(t *testing.T, fields *fields)
		wantAsserter func(t *testing.T, results <-chan plc4go.PlcConnectionPingResult) bool
	}{
		{
			name: "ping it",
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultConnectionRequirements(t)
				connection := NewMockPlcConnection(t)
				{
					connection.EXPECT().IsConnected().Return(false)
				}
				requirements.EXPECT().GetConnection().Return(connection)
				fields.DefaultConnectionRequirements = requirements
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionPingResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.NotNil(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "ping it connected",
			fields: fields{
				connected: true,
			},
			mockSetup: func(t *testing.T, fields *fields) {
				requirements := NewMockDefaultConnectionRequirements(t)
				connection := NewMockPlcConnection(t)
				{
					connection.EXPECT().IsConnected().Return(true)
				}
				requirements.EXPECT().GetConnection().Return(connection)
				fields.DefaultConnectionRequirements = requirements
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionPingResult) bool {
				timeout := time.NewTimer(2 * time.Second)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.Nil(t, result.GetErr())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Truef(t, tt.wantAsserter(t, d.Ping()), "Ping()")
		})
	}
}

func Test_defaultConnection_ReadRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   model.PlcReadRequestBuilder
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "not provided by actual connection", err)
				} else {
					t.Error("should fail")
				}
			}()
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.ReadRequestBuilder(), "ReadRequestBuilder()")
		})
	}
}

func Test_defaultConnection_SetConnected(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	type args struct {
		connected bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			d.SetConnected(tt.args.connected)
		})
	}
}

func Test_defaultConnection_SubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   model.PlcSubscriptionRequestBuilder
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "not provided by actual connection", err)
				} else {
					t.Error("should fail")
				}
			}()
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.SubscriptionRequestBuilder(), "SubscriptionRequestBuilder()")
		})
	}
}

func Test_defaultConnection_UnsubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   model.PlcUnsubscriptionRequestBuilder
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "not provided by actual connection", err)
				} else {
					t.Error("should fail")
				}
			}()
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.UnsubscriptionRequestBuilder(), "UnsubscriptionRequestBuilder()")
		})
	}
}

func Test_defaultConnection_WriteRequestBuilder(t *testing.T) {
	type fields struct {
		DefaultConnectionRequirements DefaultConnectionRequirements
		defaultTtl                    time.Duration
		connected                     bool
		tagHandler                    spi.PlcTagHandler
		valueHandler                  spi.PlcValueHandler
	}
	tests := []struct {
		name   string
		fields fields
		want   model.PlcWriteRequestBuilder
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Equal(t, "not provided by actual connection", err)
				} else {
					t.Error("should fail")
				}
			}()
			d := &defaultConnection{
				DefaultConnectionRequirements: tt.fields.DefaultConnectionRequirements,
				defaultTtl:                    tt.fields.defaultTtl,
				connected:                     tt.fields.connected,
				tagHandler:                    tt.fields.tagHandler,
				valueHandler:                  tt.fields.valueHandler,
			}
			assert.Equalf(t, tt.want, d.WriteRequestBuilder(), "WriteRequestBuilder()")
		})
	}
}

func Test_plcConnectionCloseResult_GetConnection(t *testing.T) {
	type fields struct {
		connection plc4go.PlcConnection
		err        error
		traces     []spi.TraceEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   plc4go.PlcConnection
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plcConnectionCloseResult{
				connection: tt.fields.connection,
				err:        tt.fields.err,
				traces:     tt.fields.traces,
			}
			assert.Equalf(t, tt.want, d.GetConnection(), "GetConnection()")
		})
	}
}

func Test_plcConnectionCloseResult_GetErr(t *testing.T) {
	type fields struct {
		connection plc4go.PlcConnection
		err        error
		traces     []spi.TraceEntry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get it",
			fields: fields{
				connection: NewMockPlcConnection(t),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plcConnectionCloseResult{
				connection: tt.fields.connection,
				err:        tt.fields.err,
				traces:     tt.fields.traces,
			}
			tt.wantErr(t, d.GetErr(), fmt.Sprintf("GetErr()"))
		})
	}
}

func Test_plcConnectionCloseResult_GetTraces(t *testing.T) {
	type fields struct {
		connection plc4go.PlcConnection
		err        error
		traces     []spi.TraceEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   []spi.TraceEntry
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plcConnectionCloseResult{
				connection: tt.fields.connection,
				err:        tt.fields.err,
				traces:     tt.fields.traces,
			}
			assert.Equalf(t, tt.want, d.GetTraces(), "GetTraces()")
		})
	}
}

func Test_plcConnectionConnectResult_GetConnection(t *testing.T) {
	type fields struct {
		connection plc4go.PlcConnection
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		want   plc4go.PlcConnection
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plcConnectionConnectResult{
				connection: tt.fields.connection,
				err:        tt.fields.err,
			}
			assert.Equalf(t, tt.want, d.GetConnection(), "GetConnection()")
		})
	}
}

func Test_plcConnectionConnectResult_GetErr(t *testing.T) {
	type fields struct {
		connection plc4go.PlcConnection
		err        error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get it",
			fields: fields{
				connection: NewMockPlcConnection(t),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plcConnectionConnectResult{
				connection: tt.fields.connection,
				err:        tt.fields.err,
			}
			tt.wantErr(t, d.GetErr(), fmt.Sprintf("GetErr()"))
		})
	}
}

func Test_plcConnectionPingResult_GetErr(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "get it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plcConnectionPingResult{
				err: tt.fields.err,
			}
			tt.wantErr(t, d.GetErr(), fmt.Sprintf("GetErr()"))
		})
	}
}
