/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package pool

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/simulated"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	"testing"
	"time"
)

func TestPlcConnectionPool_GetConnection(t1 *testing.T) {
	type fields struct {
		driverManager plc4go.PlcDriverManager
	}
	type args struct {
		connectionString string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantTimeout bool
	}{
		{name: "simple",
			fields: fields{
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
			}, args: args{
				connectionString: "simulated://1.2.3.4:42",
			},
			wantErr:     false,
			wantTimeout: false,
		},
		{name: "simpleWithTimeout",
			fields: fields{
				driverManager: func() plc4go.PlcDriverManager {
					driverManager := plc4go.NewPlcDriverManager()
					driverManager.RegisterDriver(simulated.NewDriver())
					return driverManager
				}(),
			}, args: args{
				connectionString: "simulated://1.2.3.4:42?connectionDelay=5",
			},
			wantErr:     false,
			wantTimeout: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := NewPlcConnectionPool(tt.fields.driverManager)
			got := t.GetConnection(tt.args.connectionString)
			select {
			case connectResult := <-got:
				if tt.wantErr && (connectResult.GetErr() == nil) {
					t1.Errorf("PlcConnectionPool.GetConnection() = %v, wantErr %v", connectResult.GetErr(), tt.wantErr)
				} else if connectResult.GetErr() != nil {
					t1.Errorf("PlcConnectionPool.GetConnection() error = %v, wantErr %v", connectResult.GetErr(), tt.wantErr)
				}
			case <-time.After(3 * time.Second):
				if !tt.wantTimeout {
					t1.Errorf("PlcConnectionPool.GetConnection() got timeout")
				}
			}
		})
	}
}
