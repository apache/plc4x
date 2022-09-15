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

package cache

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/simulated"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/viney-shih/go-lock"
	"testing"
)

func Test_connectionContainer_String(t1 *testing.T) {
	type fields struct {
		cacheLog         *zerolog.Logger
		lock             lock.RWMutex
		connectionString string
		driverManager    plc4go.PlcDriverManager
		tracerEnabled    bool
		connection       spi.PlcConnection
		leaseCounter     uint32
		closed           bool
		state            cachedPlcConnectionState
		queue            []chan plc4go.PlcConnectionConnectResult
		listeners        []connectionListener
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &connectionContainer{
				cacheLog:         tt.fields.cacheLog,
				lock:             tt.fields.lock,
				connectionString: tt.fields.connectionString,
				driverManager:    tt.fields.driverManager,
				tracerEnabled:    tt.fields.tracerEnabled,
				connection:       tt.fields.connection,
				leaseCounter:     tt.fields.leaseCounter,
				closed:           tt.fields.closed,
				state:            tt.fields.state,
				queue:            tt.fields.queue,
				listeners:        tt.fields.listeners,
			}
			assert.Equalf(t1, tt.want, t.String(), "String()")
		})
	}
}

func Test_connectionContainer_addListener(t1 *testing.T) {
	type fields struct {
		cacheLog         *zerolog.Logger
		lock             lock.RWMutex
		connectionString string
		driverManager    plc4go.PlcDriverManager
		tracerEnabled    bool
		connection       spi.PlcConnection
		leaseCounter     uint32
		closed           bool
		state            cachedPlcConnectionState
		queue            []chan plc4go.PlcConnectionConnectResult
		listeners        []connectionListener
	}
	type args struct {
		listener connectionListener
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &connectionContainer{
				cacheLog:         tt.fields.cacheLog,
				lock:             tt.fields.lock,
				connectionString: tt.fields.connectionString,
				driverManager:    tt.fields.driverManager,
				tracerEnabled:    tt.fields.tracerEnabled,
				connection:       tt.fields.connection,
				leaseCounter:     tt.fields.leaseCounter,
				closed:           tt.fields.closed,
				state:            tt.fields.state,
				queue:            tt.fields.queue,
				listeners:        tt.fields.listeners,
			}
			t.addListener(tt.args.listener)
		})
	}
}

func Test_connectionContainer_connect(t1 *testing.T) {
	type fields struct {
		cacheLog         *zerolog.Logger
		lock             lock.RWMutex
		connectionString string
		driverManager    plc4go.PlcDriverManager
		tracerEnabled    bool
		connection       spi.PlcConnection
		leaseCounter     uint32
		closed           bool
		state            cachedPlcConnectionState
		queue            []chan plc4go.PlcConnectionConnectResult
		listeners        []connectionListener
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "connect fresh",
			fields: fields{
				cacheLog: &log.Logger,
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
				connectionString: "simulated://1.2.3.4:42",
				lock:             lock.NewCASMutex(),
				queue:            []chan plc4go.PlcConnectionConnectResult{},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &connectionContainer{
				cacheLog:         tt.fields.cacheLog,
				lock:             tt.fields.lock,
				connectionString: tt.fields.connectionString,
				driverManager:    tt.fields.driverManager,
				tracerEnabled:    tt.fields.tracerEnabled,
				connection:       tt.fields.connection,
				leaseCounter:     tt.fields.leaseCounter,
				closed:           tt.fields.closed,
				state:            tt.fields.state,
				queue:            tt.fields.queue,
				listeners:        tt.fields.listeners,
			}
			t.connect()
		})
	}
}

func Test_connectionContainer_lease(t1 *testing.T) {
	type fields struct {
		cacheLog         *zerolog.Logger
		lock             lock.RWMutex
		connectionString string
		driverManager    plc4go.PlcDriverManager
		tracerEnabled    bool
		connection       spi.PlcConnection
		leaseCounter     uint32
		closed           bool
		state            cachedPlcConnectionState
		queue            []chan plc4go.PlcConnectionConnectResult
		listeners        []connectionListener
	}
	tests := []struct {
		name       string
		fields     fields
		wantNotNil bool
	}{
		{
			name: "lease fresh",
			fields: fields{
				cacheLog: &log.Logger,
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
				connectionString: "simulated://1.2.3.4:42",
				lock:             lock.NewCASMutex(),
				queue:            []chan plc4go.PlcConnectionConnectResult{},
			},
			wantNotNil: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &connectionContainer{
				cacheLog:         tt.fields.cacheLog,
				lock:             tt.fields.lock,
				connectionString: tt.fields.connectionString,
				driverManager:    tt.fields.driverManager,
				tracerEnabled:    tt.fields.tracerEnabled,
				connection:       tt.fields.connection,
				leaseCounter:     tt.fields.leaseCounter,
				closed:           tt.fields.closed,
				state:            tt.fields.state,
				queue:            tt.fields.queue,
				listeners:        tt.fields.listeners,
			}
			assert.True(t1, tt.wantNotNil, t.lease(), "lease()")
		})
	}
}

func Test_connectionContainer_returnConnection(t1 *testing.T) {
	type fields struct {
		cacheLog         *zerolog.Logger
		lock             lock.RWMutex
		connectionString string
		driverManager    plc4go.PlcDriverManager
		tracerEnabled    bool
		connection       spi.PlcConnection
		leaseCounter     uint32
		closed           bool
		state            cachedPlcConnectionState
		queue            []chan plc4go.PlcConnectionConnectResult
		listeners        []connectionListener
	}
	type args struct {
		state cachedPlcConnectionState
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "return connection fresh",
			fields: fields{
				cacheLog: &log.Logger,
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
				connectionString: "simulated://1.2.3.4:42",
				lock:             lock.NewCASMutex(),
				queue:            []chan plc4go.PlcConnectionConnectResult{},
			},
			args: args{
				state: StateInitialized,
			},
			wantErr: assert.NoError,
		},
		{
			name: "return unconnected connection",
			fields: fields{
				cacheLog: &log.Logger,
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
				connectionString: "simulated://1.2.3.4:42",
				lock:             lock.NewCASMutex(),
				queue:            []chan plc4go.PlcConnectionConnectResult{},
			},
			args: args{
				state: StateInUse,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &connectionContainer{
				cacheLog:         tt.fields.cacheLog,
				lock:             tt.fields.lock,
				connectionString: tt.fields.connectionString,
				driverManager:    tt.fields.driverManager,
				tracerEnabled:    tt.fields.tracerEnabled,
				connection:       tt.fields.connection,
				leaseCounter:     tt.fields.leaseCounter,
				closed:           tt.fields.closed,
				state:            tt.fields.state,
				queue:            tt.fields.queue,
				listeners:        tt.fields.listeners,
			}
			tt.wantErr(t1, t.returnConnection(tt.args.state), fmt.Sprintf("returnConnection(%v)", tt.args.state))
		})
	}
}
