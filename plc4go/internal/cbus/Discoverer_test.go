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
	"net"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/pool"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/nettest"
)

func TestNewDiscoverer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "just create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNilf(t, NewDiscoverer(), "NewDiscoverer()")
		})
	}
}

func TestDiscoverer_Discover(t *testing.T) {
	type fields struct {
		transportInstanceCreationQueue pool.Executor
		deviceScanningQueue            pool.Executor
	}
	type args struct {
		ctx              context.Context
		callback         func(event apiModel.PlcDiscoveryItem)
		discoveryOptions []options.WithDiscoveryOption
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  assert.ErrorAssertionFunc
		setup    func(t *testing.T, fields *fields, args *args) (params []any)
		teardown func(params []any)
	}{
		{
			name: "discover unknown device",
			args: args{
				callback: func(_ apiModel.PlcDiscoveryItem) {
				},
				discoveryOptions: []options.WithDiscoveryOption{
					options.WithDiscoveryOptionDeviceName("blub"),
				},
			},
			setup: func(t *testing.T, fields *fields, args *args) (params []any) {
				ctx, cancelFunc := context.WithCancel(context.Background())
				t.Cleanup(func() {
					cancelFunc()
					// We give it on second to settle, so it can stop everything
					time.Sleep(1 * time.Second)
				})
				args.ctx = ctx

				executor := pool.NewFixedSizeExecutor(50, 100, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				t.Cleanup(executor.Stop)
				fields.transportInstanceCreationQueue = executor
				fields.deviceScanningQueue = executor
				return nil
			},
			wantErr: assert.NoError,
		},
		{
			name: "test with loopback",
			args: args{
				callback: func(_ apiModel.PlcDiscoveryItem) {
				},
				discoveryOptions: []options.WithDiscoveryOption{
					options.WithDiscoveryOptionDeviceName("blub"),
				},
			},
			wantErr: assert.NoError,
			setup: func(t *testing.T, fields *fields, args *args) (params []any) {
				ctx, cancelFunc := context.WithCancel(context.Background())
				t.Cleanup(func() {
					cancelFunc()
					// We give it on second to settle, so it can stop everything
					time.Sleep(1 * time.Second)
				})
				args.ctx = ctx
				executor := pool.NewFixedSizeExecutor(50, 100, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				executor.Start()
				t.Cleanup(executor.Stop)
				fields.transportInstanceCreationQueue = executor
				fields.deviceScanningQueue = executor
				oldaddressProviderRetriever := addressProviderRetriever
				addressProviderRetriever = func(log zerolog.Logger, _ []string) ([]addressProvider, error) {
					loopbackInterface, err := nettest.LoopbackInterface()
					if err != nil {
						return nil, err
					}
					return []addressProvider{&wrappedInterface{loopbackInterface}}, nil
				}
				return []any{oldaddressProviderRetriever}
			},
			teardown: func(params []any) {
				addressProviderRetriever = params[0].(func(log zerolog.Logger, deviceNames []string) ([]addressProvider, error))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var params []any
			if tt.setup != nil {
				params = tt.setup(t, &tt.fields, &tt.args)
			}
			d := &Discoverer{
				transportInstanceCreationQueue: tt.fields.transportInstanceCreationQueue,
				deviceScanningQueue:            tt.fields.deviceScanningQueue,
				log:                            testutils.ProduceTestingLogger(t),
			}
			tt.wantErr(t, d.Discover(tt.args.ctx, tt.args.callback, tt.args.discoveryOptions...), fmt.Sprintf("Discover(%v, func(), %v)", tt.args.ctx, tt.args.discoveryOptions))
			if tt.teardown != nil {
				tt.teardown(params)
			}
		})
	}
}

func TestDiscoverer_createDeviceScanDispatcher(t *testing.T) {
	type fields struct {
		transportInstanceCreationQueue pool.Executor
		deviceScanningQueue            pool.Executor
	}
	type args struct {
		tcpTransportInstance *tcp.TransportInstance
		callback             func(t *testing.T, event apiModel.PlcDiscoveryItem)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args)
	}{
		{
			name: "create a dispatcher",
			args: args{
				callback: func(t *testing.T, event apiModel.PlcDiscoveryItem) {
					assert.NotNil(t, event)
				},
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				listen, err := net.Listen("tcp", "127.0.0.1:0")
				require.NoError(t, err)
				dispatchWg := sync.WaitGroup{}
				dispatchWg.Add(1)
				t.Cleanup(dispatchWg.Wait)
				go func() {
					defer dispatchWg.Done()
					conn, err := listen.Accept()
					if err != nil {
						t.Error(err)
						return
					}
					t.Logf("writing out")
					write, err := conn.Write([]byte("x.890050435F434E49454422\r\n"))
					if err != nil {
						t.Error(err)
						return
					}
					t.Logf("%d bytes written", write)
				}()
				t.Cleanup(func() {
					t.Log("close listener")
					if err := listen.Close(); err != nil {
						t.Error(err)
					}
				})

				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := tcp.NewTransport(_options...)
				parse, err := url.Parse("tcp://" + listen.Addr().String())
				require.NoError(t, err)
				instance, err := transport.CreateTransportInstance(*parse, nil, _options...)
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, instance.Close())
				})
				args.tcpTransportInstance = instance.(*tcp.TransportInstance)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			d := &Discoverer{
				transportInstanceCreationQueue: tt.fields.transportInstanceCreationQueue,
				deviceScanningQueue:            tt.fields.deviceScanningQueue,
				log:                            testutils.ProduceTestingLogger(t),
			}
			dispatcher := d.createDeviceScanDispatcher(tt.args.tcpTransportInstance, func(event apiModel.PlcDiscoveryItem) {
				tt.args.callback(t, event)
			})
			assert.NotNilf(t, dispatcher, "createDeviceScanDispatcher(%v, func())", tt.args.tcpTransportInstance)
			t.Log("Calling dispatcher now")
			dispatcher()
			t.Log("dispatching done")
		})
	}
}

func TestDiscoverer_createTransportInstanceDispatcher(t *testing.T) {
	type fields struct {
		transportInstanceCreationQueue pool.Executor
		deviceScanningQueue            pool.Executor
	}
	type args struct {
		ctx                context.Context
		wg                 *sync.WaitGroup
		ip                 net.IP
		tcpTransport       *tcp.Transport
		transportInstances chan transports.TransportInstance
		cBusPort           uint16
		addressLogger      zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args)
	}{
		{
			name: "create a dispatcher",
			args: args{
				ctx: testutils.TestContext(t),
				wg: func() *sync.WaitGroup {
					var wg sync.WaitGroup
					return &wg
				}(),
				ip:                 net.IPv4(127, 0, 0, 1),
				transportInstances: make(chan transports.TransportInstance, 1),
				cBusPort: func() uint16 {
					listen, err := net.Listen("tcp", "127.0.0.1:0")
					if err != nil {
						t.Fatal(err)
					}
					dispatchWg := sync.WaitGroup{}
					dispatchWg.Add(1)
					t.Cleanup(dispatchWg.Wait)
					go func() {
						defer dispatchWg.Done()
						conn, err := listen.Accept()
						if err != nil {
							t.Log(err)
							return
						}
						write, err := conn.Write([]byte("x.890050435F434E49454422\r\n"))
						if err != nil {
							t.Error(err)
							return
						}
						t.Logf("%d written", write)
					}()
					t.Cleanup(func() {
						if err := listen.Close(); err != nil {
							t.Error(err)
						}
					})
					parse, err := url.Parse("tcp://" + listen.Addr().String())
					if err != nil {
						t.Fatal(err)
					}
					port, err := strconv.ParseUint(parse.Port(), 10, 16)
					if err != nil {
						t.Fatal(err)
					}
					return uint16(port)
				}(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				args.tcpTransport = tcp.NewTransport(options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			d := &Discoverer{
				transportInstanceCreationQueue: tt.fields.transportInstanceCreationQueue,
				deviceScanningQueue:            tt.fields.deviceScanningQueue,
				log:                            testutils.ProduceTestingLogger(t),
			}
			dispatcher := d.createTransportInstanceDispatcher(tt.args.ctx, tt.args.wg, tt.args.ip, tt.args.tcpTransport, tt.args.transportInstances, tt.args.cBusPort, tt.args.addressLogger)
			assert.NotNilf(t, dispatcher, "createTransportInstanceDispatcher(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.wg, tt.args.ip, tt.args.tcpTransport, tt.args.transportInstances)
			dispatcher()
			timeout := time.NewTimer(2 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case <-timeout.C:
				t.Error("timeout")
			case ti := <-tt.args.transportInstances:
				timeout.Stop()
				assert.NotNil(t, ti)
			}
		})
	}
}

func TestDiscoverer_extractDeviceNames(t *testing.T) {
	type fields struct {
		transportInstanceCreationQueue pool.Executor
		deviceScanningQueue            pool.Executor
	}
	type args struct {
		discoveryOptions []options.WithDiscoveryOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "no options, no devices",
			want: []string{},
		},
		{
			name: "one device option",
			args: args{
				discoveryOptions: []options.WithDiscoveryOption{
					options.WithDiscoveryOptionDeviceName("blub"),
				},
			},
			want: []string{"blub"},
		},
		{
			name: "two device option",
			args: args{
				discoveryOptions: []options.WithDiscoveryOption{
					options.WithDiscoveryOptionDeviceName("blub"),
					options.WithDiscoveryOptionDeviceName("blab"),
				},
			},
			want: []string{"blub", "blab"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				transportInstanceCreationQueue: tt.fields.transportInstanceCreationQueue,
				deviceScanningQueue:            tt.fields.deviceScanningQueue,
				log:                            testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, d.extractDeviceNames(tt.args.discoveryOptions...), "extractDeviceNames(%v)", tt.args.discoveryOptions)
		})
	}
}

func Test_wrappedInterface_containedInterface(t *testing.T) {
	type fields struct {
		Interface *net.Interface
	}
	tests := []struct {
		name   string
		fields fields
		want   net.Interface
	}{
		{
			name: "get it",
			fields: fields{
				Interface: &net.Interface{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &wrappedInterface{
				Interface: tt.fields.Interface,
			}
			assert.Equalf(t, tt.want, w.containedInterface(), "containedInterface()")
		})
	}
}

func Test_wrappedInterface_name(t *testing.T) {
	type fields struct {
		Interface *net.Interface
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
			fields: fields{
				Interface: &net.Interface{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &wrappedInterface{
				Interface: tt.fields.Interface,
			}
			assert.Equalf(t, tt.want, w.name(), "name()")
		})
	}
}
