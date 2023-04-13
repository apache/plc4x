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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/nettest"
	"net"
	"sync"
	"sync/atomic"
	"testing"
)

func TestDiscoverer_Discover(t *testing.T) {
	type fields struct {
		transportInstanceCreationQueue utils.Executor
		deviceScanningQueue            utils.Executor
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
		setup    func() (params []interface{})
		teardown func(params []interface{})
	}{
		{
			name: "discover unknown device",
			fields: fields{
				transportInstanceCreationQueue: utils.NewFixedSizeExecutor(50, 100),
				deviceScanningQueue:            utils.NewFixedSizeExecutor(50, 100),
			},
			args: args{
				ctx: context.Background(),
				callback: func(_ apiModel.PlcDiscoveryItem) {
				},
				discoveryOptions: []options.WithDiscoveryOption{
					options.WithDiscoveryOptionDeviceName("blub"),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "test with loopback",
			fields: fields{
				transportInstanceCreationQueue: utils.NewFixedSizeExecutor(50, 100),
				deviceScanningQueue:            utils.NewFixedSizeExecutor(50, 100),
			},
			args: args{
				ctx: context.Background(),
				callback: func(_ apiModel.PlcDiscoveryItem) {
				},
				discoveryOptions: []options.WithDiscoveryOption{
					options.WithDiscoveryOptionDeviceName("blub"),
				},
			},
			wantErr: assert.NoError,
			setup: func() (params []interface{}) {
				oldaddressProviderRetriever := addressProviderRetriever
				addressProviderRetriever = func(_ []string) ([]addressProvider, error) {
					loopbackInterface, err := nettest.LoopbackInterface()
					if err != nil {
						return nil, err
					}
					return []addressProvider{&wrappedInterface{loopbackInterface}}, nil
				}
				return []interface{}{oldaddressProviderRetriever}
			},
			teardown: func(params []interface{}) {
				addressProviderRetriever = params[0].(func(deviceNames []string) ([]addressProvider, error))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				transportInstanceCreationQueue: tt.fields.transportInstanceCreationQueue,
				deviceScanningQueue:            tt.fields.deviceScanningQueue,
			}
			var params []interface{}
			if tt.setup != nil {
				params = tt.setup()
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
		transportInstanceCreationWorkItemId atomic.Int32
		transportInstanceCreationQueue      utils.Executor
		deviceScanningWorkItemId            atomic.Int32
		deviceScanningQueue                 utils.Executor
	}
	type args struct {
		tcpTransportInstance *tcp.TransportInstance
		callback             func(event apiModel.PlcDiscoveryItem)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   utils.Runnable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				transportInstanceCreationWorkItemId: tt.fields.transportInstanceCreationWorkItemId,
				transportInstanceCreationQueue:      tt.fields.transportInstanceCreationQueue,
				deviceScanningWorkItemId:            tt.fields.deviceScanningWorkItemId,
				deviceScanningQueue:                 tt.fields.deviceScanningQueue,
			}
			assert.Equalf(t, tt.want, d.createDeviceScanDispatcher(tt.args.tcpTransportInstance, tt.args.callback), "createDeviceScanDispatcher(%v, %v)", tt.args.tcpTransportInstance, tt.args.callback)
		})
	}
}

func TestDiscoverer_createTransportInstanceDispatcher(t *testing.T) {
	type fields struct {
		transportInstanceCreationWorkItemId atomic.Int32
		transportInstanceCreationQueue      utils.Executor
		deviceScanningWorkItemId            atomic.Int32
		deviceScanningQueue                 utils.Executor
	}
	type args struct {
		ctx                context.Context
		wg                 *sync.WaitGroup
		ip                 net.IP
		tcpTransport       *tcp.Transport
		transportInstances chan transports.TransportInstance
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   utils.Runnable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				transportInstanceCreationWorkItemId: tt.fields.transportInstanceCreationWorkItemId,
				transportInstanceCreationQueue:      tt.fields.transportInstanceCreationQueue,
				deviceScanningWorkItemId:            tt.fields.deviceScanningWorkItemId,
				deviceScanningQueue:                 tt.fields.deviceScanningQueue,
			}
			assert.Equalf(t, tt.want, d.createTransportInstanceDispatcher(tt.args.ctx, tt.args.wg, tt.args.ip, tt.args.tcpTransport, tt.args.transportInstances), "createTransportInstanceDispatcher(%v, %v, %v, %v, %v)", tt.args.ctx, tt.args.wg, tt.args.ip, tt.args.tcpTransport, tt.args.transportInstances)
		})
	}
}

func TestNewDiscoverer(t *testing.T) {
	tests := []struct {
		name string
		want *Discoverer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDiscoverer(), "NewDiscoverer()")
		})
	}
}
