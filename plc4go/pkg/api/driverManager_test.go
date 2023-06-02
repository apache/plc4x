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

package plc4go

import (
	"context"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog"
	"os"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewPlcDriverManager(t *testing.T) {
	tests := []struct {
		name string
		want PlcDriverManager
	}{
		{
			name: "create one",
			want: &plcDriverManger{
				drivers:    map[string]PlcDriver{},
				transports: map[string]transports.Transport{},
				log:        zerolog.Nop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPlcDriverManager(config.WithCustomLogger(zerolog.Nop())); !assert.Equal(t, got, tt.want) {
				t.Errorf("NewPlcDriverManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionDeviceName(t *testing.T) {
	type args struct {
		deviceName string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "with it",
			args: args{
				deviceName: "eth0",
			},
			want: withDiscoveryOption{options.WithDiscoveryOptionDeviceName("eth0")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionDeviceName(tt.args.deviceName); !assert.Equal(t, got, tt.want) {
				t.Errorf("WithDiscoveryOptionDeviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionLocalAddress(t *testing.T) {
	type args struct {
		localAddress string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "with it",
			args: args{
				localAddress: "1.1.1.1",
			},
			want: withDiscoveryOption{options.WithDiscoveryOptionLocalAddress("1.1.1.1")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionLocalAddress(tt.args.localAddress); !assert.Equal(t, got, tt.want) {
				t.Errorf("WithDiscoveryOptionLocalAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionProtocol(t *testing.T) {
	type args struct {
		protocolName string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "with it",
			args: args{
				protocolName: "test",
			},
			want: withDiscoveryOption{options.WithDiscoveryOptionProtocol("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionProtocol(tt.args.protocolName); !assert.Equal(t, got, tt.want) {
				t.Errorf("WithDiscoveryOptionProtocol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionProtocolSpecific(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "with it",
			args: args{
				key:   "some option",
				value: 13,
			},
			want: withDiscoveryOption{options.WithDiscoveryOptionProtocolSpecific("some option", 13)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionProtocolSpecific(tt.args.key, tt.args.value); !assert.Equal(t, got, tt.want) {
				t.Errorf("WithDiscoveryOptionProtocolSpecific() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionRemoteAddress(t *testing.T) {
	type args struct {
		remoteAddress string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "with it",
			args: args{
				remoteAddress: "127.0.0.1",
			},
			want: withDiscoveryOption{options.WithDiscoveryOptionRemoteAddress("127.0.0.1")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionRemoteAddress(tt.args.remoteAddress); !assert.Equal(t, got, tt.want) {
				t.Errorf("WithDiscoveryOptionRemoteAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDiscoveryOptionTransport(t *testing.T) {
	type args struct {
		transportName string
	}
	tests := []struct {
		name string
		args args
		want WithDiscoveryOption
	}{
		{
			name: "with it",
			args: args{
				transportName: "udp",
			},
			want: withDiscoveryOption{options.WithDiscoveryOptionTransport("udp")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDiscoveryOptionTransport(tt.args.transportName); !assert.Equal(t, got, tt.want) {
				t.Errorf("WithDiscoveryOptionTransport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertToInternalOptions(t *testing.T) {
	type args struct {
		withDiscoveryOptions []WithDiscoveryOption
	}
	tests := []struct {
		name string
		args args
		want []options.WithDiscoveryOption
	}{
		{
			name: "convert nothing",
			want: []options.WithDiscoveryOption{},
		},
		{
			name: "convert something",
			args: args{
				withDiscoveryOptions: []WithDiscoveryOption{
					WithDiscoveryOptionRemoteAddress("remote"),
					WithDiscoveryOptionTransport("udp"),
				},
			},
			want: []options.WithDiscoveryOption{
				options.WithDiscoveryOptionRemoteAddress("remote"),
				options.WithDiscoveryOptionTransport("udp"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertToInternalOptions(tt.args.withDiscoveryOptions...); !assert.Equal(t, got, tt.want) {
				t.Errorf("convertToInternalOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plcConnectionConnectResult_GetConnection(t *testing.T) {
	type fields struct {
		connection PlcConnection
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		want   PlcConnection
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
			if got := d.GetConnection(); !assert.Equal(t, got, tt.want) {
				t.Errorf("GetConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plcConnectionConnectResult_GetErr(t *testing.T) {
	type fields struct {
		connection PlcConnection
		err        error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
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
			if err := d.GetErr(); (err != nil) != tt.wantErr {
				t.Errorf("GetErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_plcDriverManger_Discover(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		callback         func(event model.PlcDiscoveryItem)
		discoveryOptions []WithDiscoveryOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "discover it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			if err := m.Discover(tt.args.callback, tt.args.discoveryOptions...); (err != nil) != tt.wantErr {
				t.Errorf("Discover() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_plcDriverManger_DiscoverWithContext(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		ctx              context.Context
		callback         func(event model.PlcDiscoveryItem)
		discoveryOptions []WithDiscoveryOption
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   bool
	}{
		{
			name: "discover it",
			fields: fields{
				drivers: map[string]PlcDriver{},
			},
			args: args{
				ctx: context.Background(),
				callback: func(event model.PlcDiscoveryItem) {
					// No-op
				},
				discoveryOptions: []WithDiscoveryOption{
					WithDiscoveryOptionTransport("test"),
					WithDiscoveryOptionProtocol("test"),
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				driver := NewMockPlcDriver(t)
				expect := driver.EXPECT()
				expect.GetProtocolName().Return("test")
				expect.SupportsDiscovery().Return(true)
				expect.DiscoverWithContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.drivers["test"] = driver
			},
		},
		{
			name: "discover it with an error",
			fields: fields{
				drivers: map[string]PlcDriver{},
			},
			args: args{
				ctx: context.Background(),
				callback: func(event model.PlcDiscoveryItem) {
					// No-op
				},
				discoveryOptions: []WithDiscoveryOption{
					WithDiscoveryOptionTransport("test"),
					WithDiscoveryOptionProtocol("test"),
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				driver := NewMockPlcDriver(t)
				expect := driver.EXPECT()
				expect.GetProtocolName().Return("test")
				expect.SupportsDiscovery().Return(true)
				expect.DiscoverWithContext(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Uh no"))
				fields.drivers["test"] = driver
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			if err := m.DiscoverWithContext(tt.args.ctx, tt.args.callback, tt.args.discoveryOptions...); (err != nil) != tt.wantErr {
				t.Errorf("DiscoverWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_plcDriverManger_GetConnection(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		connectionString string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		mockSetup    func(t *testing.T, fields *fields, args *args)
		wantVerifier func(t *testing.T, results <-chan PlcConnectionConnectResult) bool
	}{
		{
			name: "get one with wrong url",
			args: args{
				connectionString: "~:/?#[]@!$&'()*+,;=\n",
			},
			wantVerifier: func(t *testing.T, results <-chan PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(3 * time.Second)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.NotNil(t, result)
					assert.Nil(t, result.GetConnection())
					assert.NotNil(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "get one without a driver",
			wantVerifier: func(t *testing.T, results <-chan PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(3 * time.Second)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.NotNil(t, result)
					assert.Nil(t, result.GetConnection())
					assert.NotNil(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "get one with a driver",
			fields: fields{
				drivers: map[string]PlcDriver{},
			},
			args: args{
				connectionString: "test://something",
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				driver := NewMockPlcDriver(t)
				expect := driver.EXPECT()
				expect.GetProtocolName().Return("test")
				expect.GetDefaultTransport().Return("test")
				results := make(chan PlcConnectionConnectResult, 1)
				result := NewMockPlcConnectionConnectResult(t)
				result.EXPECT().GetConnection().Return(nil)
				result.EXPECT().GetErr().Return(nil)
				results <- result
				expect.GetConnection(mock.Anything, mock.Anything, mock.Anything).Return(results)
				fields.drivers["test"] = driver
			},
			wantVerifier: func(t *testing.T, results <-chan PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(3 * time.Second)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case <-timeout.C:
					t.Error("timeout")
				case result := <-results:
					assert.NotNil(t, result)
					assert.Nil(t, result.GetConnection())
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
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			if got := m.GetConnection(tt.args.connectionString); !tt.wantVerifier(t, got) {
				t.Errorf("GetConnection() = %v", got)
			}
		})
	}
}

func Test_plcDriverManger_GetDriver(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		driverName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    PlcDriver
		wantErr bool
	}{
		{
			name:    "get it (not there)",
			wantErr: true,
		},
		{
			name: "get it",
			fields: fields{
				drivers: map[string]PlcDriver{
					"": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			got, err := m.GetDriver(tt.args.driverName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("GetDriver() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plcDriverManger_GetTransport(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		transportName string
		in1           string
		in2           map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    transports.Transport
		wantErr bool
	}{
		{
			name:    "get it (no transport)",
			wantErr: true,
		},
		{
			name: "get it",
			fields: fields{
				transports: map[string]transports.Transport{
					"": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			got, err := m.GetTransport(tt.args.transportName, tt.args.in1, tt.args.in2)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, got, tt.want) {
				t.Errorf("GetTransport() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plcDriverManger_ListDriverNames(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "list em",
		},
		{
			name: "list em all",
			fields: fields{
				drivers: map[string]PlcDriver{
					"test": nil,
				},
			},
			want: []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			if got := m.ListDriverNames(); !assert.Equal(t, got, tt.want) {
				t.Errorf("ListDriverNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plcDriverManger_ListTransportNames(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "list em",
		},
		{
			name: "list em all",
			fields: fields{
				transports: map[string]transports.Transport{
					"test": nil,
				},
			},
			want: []string{"test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			if got := m.ListTransportNames(); !assert.Equal(t, got, tt.want) {
				t.Errorf("ListTransportNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_plcDriverManger_RegisterDriver(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		driver PlcDriver
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
	}{
		{
			name: "register it (already registered)",
			fields: fields{
				drivers: map[string]PlcDriver{
					"test": nil,
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				driver := NewMockPlcDriver(t)
				expect := driver.EXPECT()
				expect.GetProtocolName().Return("test")
				expect.GetProtocolCode().Return("test")
				args.driver = driver
			},
		},
		{
			name: "register it",
			fields: fields{
				drivers: map[string]PlcDriver{},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				driver := NewMockPlcDriver(t)
				expect := driver.EXPECT()
				expect.GetProtocolName().Return("test")
				expect.GetProtocolCode().Return("test")
				args.driver = driver
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			m.RegisterDriver(tt.args.driver)
		})
	}
}

func Test_plcDriverManger_RegisterTransport(t *testing.T) {
	type fields struct {
		drivers    map[string]PlcDriver
		transports map[string]transports.Transport
	}
	type args struct {
		transport transports.Transport
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
	}{
		{
			name: "register it (already registered)",
			fields: fields{
				transports: map[string]transports.Transport{
					"test": nil,
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				transport := NewMockTransport(t)
				transport.EXPECT().GetTransportName().Return("test")
				transport.EXPECT().GetTransportCode().Return("test")
				args.transport = transport
			},
		},
		{
			name: "register it",
			fields: fields{
				transports: map[string]transports.Transport{},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				transport := NewMockTransport(t)
				transport.EXPECT().GetTransportName().Return("test")
				transport.EXPECT().GetTransportCode().Return("test")
				args.transport = transport
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := &plcDriverManger{
				drivers:    tt.fields.drivers,
				transports: tt.fields.transports,
			}
			m.log = produceTestLogger(t)
			m.RegisterTransport(tt.args.transport)
		})
	}
}

func Test_withDiscoveryOption_isDiscoveryOption(t *testing.T) {
	type fields struct {
		WithDiscoveryOption options.WithDiscoveryOption
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "it is",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := withDiscoveryOption{
				WithDiscoveryOption: tt.fields.WithDiscoveryOption,
			}
			if got := w.isDiscoveryOption(); got != tt.want {
				t.Errorf("isDiscoveryOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

// note: we can't use testutils here due to import cycle
func produceTestLogger(t *testing.T) zerolog.Logger {
	return zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t),
		func(w *zerolog.ConsoleWriter) {
			// TODO: this is really an issue with go-junit-report not sanitizing output before dumping into xml...
			onJenkins := os.Getenv("JENKINS_URL") != ""
			onGithubAction := os.Getenv("GITHUB_ACTIONS") != ""
			onCI := os.Getenv("CI") != ""
			if onJenkins || onGithubAction || onCI {
				w.NoColor = true
			}
		}))
}
