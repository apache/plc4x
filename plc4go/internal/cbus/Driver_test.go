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
	"fmt"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestDriver_DiscoverWithContext(t *testing.T) {
	type fields struct {
		DefaultDriver           _default.DefaultDriver
		tm                      transactions.RequestTransactionManager
		awaitSetupComplete      bool
		awaitDisconnectComplete bool
	}
	type args struct {
		ctx              context.Context
		callback         func(event apiModel.PlcDiscoveryItem)
		discoveryOptions []options.WithDiscoveryOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "localhost discovery",
			args: args{
				callback: func(event apiModel.PlcDiscoveryItem) {
					t.Log(event)
				},
				discoveryOptions: []options.WithDiscoveryOption{options.WithDiscoveryOptionLocalAddress("localhost")},
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				ctx, cancelFunc := context.WithCancel(context.Background())
				t.Cleanup(func() {
					cancelFunc()
					// We give it on second to settle, so it can stop everything
					time.Sleep(200 * time.Millisecond)
				})
				args.ctx = ctx
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &Driver{
				DefaultDriver:           tt.fields.DefaultDriver,
				tm:                      tt.fields.tm,
				awaitSetupComplete:      tt.fields.awaitSetupComplete,
				awaitDisconnectComplete: tt.fields.awaitDisconnectComplete,
			}
			m.log = testutils.ProduceTestingLogger(t)
			tt.wantErr(t, m.DiscoverWithContext(tt.args.ctx, tt.args.callback, tt.args.discoveryOptions...), fmt.Sprintf("DiscoverWithContext(%v, func()*, %v)", tt.args.ctx, tt.args.discoveryOptions))
		})
	}
}

func TestDriver_GetConnectionWithContext(t *testing.T) {
	type fields struct {
		DefaultDriver           _default.DefaultDriver
		tm                      transactions.RequestTransactionManager
		awaitSetupComplete      bool
		awaitDisconnectComplete bool
	}
	type args struct {
		ctx          context.Context
		transportUrl url.URL
		transports   map[string]transports.Transport
		options      map[string][]string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantVerifier func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "get connection transport not found",
			fields: fields{
				DefaultDriver:           _default.NewDefaultDriver(nil, "test", "test", "test", NewTagHandler()),
				tm:                      nil,
				awaitSetupComplete:      false,
				awaitDisconnectComplete: false,
			},
			args: args{
				ctx: context.Background(),
				transportUrl: url.URL{
					Scheme: "test",
				},
				transports: map[string]transports.Transport{},
				options:    map[string][]string{},
			},
			wantVerifier: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(20 * time.Millisecond)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
					return false
				case result := <-results:
					assert.Error(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "get connection invalid options for transport",
			fields: fields{
				DefaultDriver:           _default.NewDefaultDriver(nil, "test", "test", "test", NewTagHandler()),
				tm:                      nil,
				awaitSetupComplete:      false,
				awaitDisconnectComplete: false,
			},
			args: args{
				ctx: context.Background(),
				transportUrl: url.URL{
					Scheme: "test",
				},
				transports: map[string]transports.Transport{
					"test": test.NewTransport(),
				},
				options: map[string][]string{
					"failTestTransport": {"yesSir"},
				},
			},
			wantVerifier: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(20 * time.Millisecond)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
					return false
				case result := <-results:
					assert.Error(t, result.GetErr())
					assert.Equal(t, "couldn't initialize transport configuration for given transport url test:: test transport failed on purpose", result.GetErr().Error())
				}
				return true
			},
		},
		{
			name: "get connection invalid options for driver",
			fields: fields{
				DefaultDriver:           _default.NewDefaultDriver(nil, "test", "test", "test", NewTagHandler()),
				tm:                      nil,
				awaitSetupComplete:      false,
				awaitDisconnectComplete: false,
			},
			args: args{
				ctx: context.Background(),
				transportUrl: url.URL{
					Scheme: "test",
				},
				transports: map[string]transports.Transport{
					"test": test.NewTransport(),
				},
				options: map[string][]string{
					"MonitoredApplication1": {"pineapple"},
				},
			},
			wantVerifier: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(20 * time.Millisecond)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
					return false
				case result := <-results:
					assert.Error(t, result.GetErr())
				}
				return true
			},
		},
		{
			name: "get connection",
			fields: fields{
				DefaultDriver:           _default.NewDefaultDriver(nil, "test", "test", "test", NewTagHandler()),
				tm:                      nil,
				awaitSetupComplete:      false,
				awaitDisconnectComplete: false,
			},
			args: args{
				ctx: context.Background(),
				transportUrl: url.URL{
					Scheme: "test",
				},
				transports: map[string]transports.Transport{
					"test": test.NewTransport(),
				},
				options: map[string][]string{},
			},
			wantVerifier: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(20 * time.Millisecond)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
					return false
				case result := <-results:
					assert.NoError(t, result.GetErr())
					assert.NotNil(t, result.GetConnection())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Driver{
				DefaultDriver:           tt.fields.DefaultDriver,
				tm:                      tt.fields.tm,
				awaitSetupComplete:      tt.fields.awaitSetupComplete,
				awaitDisconnectComplete: tt.fields.awaitDisconnectComplete,
			}
			assert.Truef(t, tt.wantVerifier(t, m.GetConnectionWithContext(tt.args.ctx, tt.args.transportUrl, tt.args.transports, tt.args.options)), "GetConnectionWithContext(%v, %v, %v, %v)", tt.args.ctx, tt.args.transportUrl, tt.args.transports, tt.args.options)
		})
	}
}

func TestDriver_SetAwaitDisconnectComplete(t *testing.T) {
	NewDriver().(*Driver).SetAwaitDisconnectComplete(true)
}

func TestDriver_SetAwaitSetupComplete(t *testing.T) {
	NewDriver().(*Driver).SetAwaitSetupComplete(true)
}

func TestDriver_SupportsDiscovery(t *testing.T) {
	NewDriver().(*Driver).SupportsDiscovery()
}

func TestNewDriver(t *testing.T) {
	tests := []struct {
		name string
		want plc4go.PlcDriver
	}{
		{
			name: "create",
			want: NewDriver(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDriver(), "NewDriver()")
		})
	}
}

func TestDriver_reportError(t *testing.T) {
	type fields struct {
		DefaultDriver           _default.DefaultDriver
		tm                      transactions.RequestTransactionManager
		awaitSetupComplete      bool
		awaitDisconnectComplete bool
	}
	type args struct {
		err error
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantAsserter func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool
	}{
		{
			name: "report it",
			args: args{
				err: errors.New("No no no no no"),
			},
			wantAsserter: func(t *testing.T, results <-chan plc4go.PlcConnectionConnectResult) bool {
				timeout := time.NewTimer(20 * time.Millisecond)
				defer utils.CleanupTimer(timeout)
				select {
				case <-timeout.C:
					t.Error("timeout")
					return false
				case result := <-results:
					assert.Error(t, result.GetErr())
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Driver{
				DefaultDriver:           tt.fields.DefaultDriver,
				tm:                      tt.fields.tm,
				awaitSetupComplete:      tt.fields.awaitSetupComplete,
				awaitDisconnectComplete: tt.fields.awaitDisconnectComplete,
			}
			assert.Truef(t, tt.wantAsserter(t, m.reportError(tt.args.err)), "reportError(%v)", tt.args.err)
		})
	}
}
