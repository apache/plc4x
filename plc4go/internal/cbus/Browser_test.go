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
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBrowser(t *testing.T) {
	assert.NotNil(t, NewBrowser(nil))
}

func TestBrowser_BrowseQuery(t *testing.T) {
	type fields struct {
		DefaultBrowser  _default.DefaultBrowser
		connection      plc4go.PlcConnection
		sequenceCounter uint8
	}
	type args struct {
		ctx         context.Context
		interceptor func(result apiModel.PlcBrowseItem) bool
		queryName   string
		query       apiModel.PlcQuery
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		setup            func(t *testing.T, fields *fields)
		wantResponseCode apiModel.PlcResponseCode
		wantQueryResults []apiModel.PlcBrowseItem
	}{
		{
			name:             "invalid address",
			wantResponseCode: apiModel.PlcResponseCode_INVALID_ADDRESS,
		},
		{
			name: "non responding browse",
			args: args{
				ctx: testutils.TestContext(t),
				interceptor: func(result apiModel.PlcBrowseItem) bool {
					// No-OP
					return true
				},
				queryName: "testQuery",
				query:     NewUnitInfoQuery(readWriteModel.NewUnitAddress(2), nil, 1),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, transportInstance.Close())
				})
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
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					t.Logf("reacting to\n%s", hex.Dump(data))
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
						t.Log("Connection dance done")

						t.Log("Dispatching 3 MMI segments")
						transportInstance.FillReadBuffer([]byte("86020200F900FF0094120006000000000000000008000000000000000000CA\r\n"))
						transportInstance.FillReadBuffer([]byte("86020200F900FF580000000000000000000000000000000000000000000026\r\n"))
						transportInstance.FillReadBuffer([]byte("86020200F700FFB00000000000000000000000000000000000000000D0\r\n"))

						t.Log("Dispatching manufacturer")
						transportInstance.FillReadBuffer([]byte("g.890050435F434E49454422\r\n"))
					}
				})
				err = transport.AddPreregisteredInstances(transportUrl, transportInstance)
				require.NoError(t, err)
				driver := NewDriver(_options...)
				t.Cleanup(func() {
					assert.NoError(t, driver.Close())
				})
				connectionConnectResult := <-driver.GetConnection(transportUrl, map[string]transports.Transport{"test": transport}, map[string][]string{})
				if err := connectionConnectResult.GetErr(); err != nil {
					t.Fatal(err)
				}
				fields.connection = connectionConnectResult.GetConnection()
				t.Cleanup(func() {
					timer := time.NewTimer(10 * time.Second)
					t.Cleanup(func() {
						utils.CleanupTimer(timer)
					})
					select {
					case <-fields.connection.Close():
					case <-timer.C:
						t.Error("timeout")
					}
				})
			},
			wantResponseCode: apiModel.PlcResponseCode_OK,
			wantQueryResults: []apiModel.PlcBrowseItem{
				&spiModel.DefaultPlcBrowseItem{
					Tag:      NewCALIdentifyTag(readWriteModel.NewUnitAddress(2), nil, readWriteModel.Attribute_Manufacturer, 1),
					Name:     "testQuery",
					Readable: true,
					Options: map[string]apiValues.PlcValue{
						"CurrentValue": spiValues.NewPlcSTRING("PC_CNIED"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			m := Browser{
				DefaultBrowser:  tt.fields.DefaultBrowser,
				connection:      tt.fields.connection,
				sequenceCounter: tt.fields.sequenceCounter,
				log:             testutils.ProduceTestingLogger(t),
			}
			got, got1 := m.BrowseQuery(tt.args.ctx, tt.args.interceptor, tt.args.queryName, tt.args.query)
			assert.Equalf(t, tt.wantResponseCode, got, "BrowseQuery(%v, func(), %v,\n%v\n)", tt.args.ctx, tt.args.queryName, tt.args.query)
			assert.Equalf(t, tt.wantQueryResults, got1, "BrowseQuery(%v, func(), %v, \n%v\n)", tt.args.ctx, tt.args.queryName, tt.args.query)
			if m.connection != nil && m.connection.IsConnected() {
				t.Log("Closing connection")
				<-m.connection.Close()
			}
		})
	}
}

func TestBrowser_browseUnitInfo(t *testing.T) {
	type fields struct {
		DefaultBrowser  _default.DefaultBrowser
		connection      plc4go.PlcConnection
		sequenceCounter uint8
	}
	type args struct {
		ctx         context.Context
		interceptor func(result apiModel.PlcBrowseItem) bool
		queryName   string
		query       *unitInfoQuery
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		setup            func(t *testing.T, fields *fields)
		wantResponseCode apiModel.PlcResponseCode
		wantQueryResults []apiModel.PlcBrowseItem
	}{
		{
			name: "non responding browse",
			args: args{
				ctx: testutils.TestContext(t),
				interceptor: func(result apiModel.PlcBrowseItem) bool {
					// No-OP
					return true
				},
				queryName: "testQuery",
				query:     NewUnitInfoQuery(readWriteModel.NewUnitAddress(2), nil, 1).(*unitInfoQuery),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, transportInstance.Close())
				})
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
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					t.Logf("reacting to\n%s", hex.Dump(data))
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
						t.Log("Connection dance done")

						t.Log("Dispatching 3 MMI segments")
						transportInstance.FillReadBuffer([]byte("86020200F900FF0094120006000000000000000008000000000000000000CA\r\n"))
						transportInstance.FillReadBuffer([]byte("86020200F900FF580000000000000000000000000000000000000000000026\r\n"))
						transportInstance.FillReadBuffer([]byte("86020200F700FFB00000000000000000000000000000000000000000D0\r\n"))

						t.Log("Dispatching manufacturer")
						transportInstance.FillReadBuffer([]byte("g.890050435F434E49454422\r\n"))
					}
				})
				err = transport.AddPreregisteredInstances(transportUrl, transportInstance)
				require.NoError(t, err)
				driver := NewDriver(_options...)
				t.Cleanup(func() {
					assert.NoError(t, driver.Close())
				})
				connectionConnectResult := <-driver.GetConnection(transportUrl, map[string]transports.Transport{"test": transport}, map[string][]string{})
				if err := connectionConnectResult.GetErr(); err != nil {
					t.Error(err)
					t.FailNow()
				}
				fields.connection = connectionConnectResult.GetConnection()
				t.Cleanup(func() {
					timer := time.NewTimer(10 * time.Second)
					t.Cleanup(func() {
						utils.CleanupTimer(timer)
					})
					select {
					case <-fields.connection.Close():
					case <-timer.C:
						t.Error("timeout")
					}
				})
			},
			wantResponseCode: apiModel.PlcResponseCode_OK,
			wantQueryResults: []apiModel.PlcBrowseItem{
				&spiModel.DefaultPlcBrowseItem{
					Tag:      NewCALIdentifyTag(readWriteModel.NewUnitAddress(2), nil, readWriteModel.Attribute_Manufacturer, 1),
					Name:     "testQuery",
					Readable: true,
					Options: map[string]apiValues.PlcValue{
						"CurrentValue": spiValues.NewPlcSTRING("PC_CNIED"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
				t.Log("Setup done")
			}
			m := &Browser{
				DefaultBrowser:  tt.fields.DefaultBrowser,
				connection:      tt.fields.connection,
				sequenceCounter: tt.fields.sequenceCounter,
				log:             testutils.ProduceTestingLogger(t),
			}
			gotResponseCode, gotQueryResults := m.browseUnitInfo(tt.args.ctx, tt.args.interceptor, tt.args.queryName, tt.args.query)
			assert.Equalf(t, tt.wantResponseCode, gotResponseCode, "browseUnitInfo(%v, %v, %v, %v)", tt.args.ctx, tt.args.interceptor, tt.args.queryName, tt.args.query)
			assert.Equalf(t, tt.wantQueryResults, gotQueryResults, "browseUnitInfo(%v, %v, %v, %v)", tt.args.ctx, tt.args.interceptor, tt.args.queryName, tt.args.query)
			if m.connection != nil && m.connection.IsConnected() {
				t.Log("Closing connection")
				<-m.connection.Close()
			}
		})
	}
}

func TestBrowser_extractUnits(t *testing.T) {
	type fields struct {
		DefaultBrowser  _default.DefaultBrowser
		connection      plc4go.PlcConnection
		sequenceCounter uint8
	}
	type args struct {
		ctx                          context.Context
		query                        *unitInfoQuery
		getInstalledUnitAddressBytes func(ctx context.Context) (map[byte]any, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []readWriteModel.UnitAddress
		want1   bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "one unit",
			args: args{
				ctx: testutils.TestContext(t),
				query: &unitInfoQuery{
					unitAddress: readWriteModel.NewUnitAddress(2),
				},
			},
			want:    []readWriteModel.UnitAddress{readWriteModel.NewUnitAddress(2)},
			want1:   false,
			wantErr: assert.NoError,
		},
		{
			name: "all units error",
			args: args{
				ctx:   testutils.TestContext(t),
				query: &unitInfoQuery{},
				getInstalledUnitAddressBytes: func(ctx context.Context) (map[byte]any, error) {
					return nil, errors.New("not today")
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "all units",
			args: args{
				ctx:   testutils.TestContext(t),
				query: &unitInfoQuery{},
				getInstalledUnitAddressBytes: func(ctx context.Context) (map[byte]any, error) {
					return map[byte]any{0xAF: true, 0xFE: true}, nil
				},
			},
			want:    []readWriteModel.UnitAddress{readWriteModel.NewUnitAddress(0xAF), readWriteModel.NewUnitAddress(0xFE)},
			want1:   true,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Browser{
				DefaultBrowser:  tt.fields.DefaultBrowser,
				connection:      tt.fields.connection,
				sequenceCounter: tt.fields.sequenceCounter,
				log:             testutils.ProduceTestingLogger(t),
			}
			got, got1, err := m.extractUnits(tt.args.ctx, tt.args.query, tt.args.getInstalledUnitAddressBytes)
			if !tt.wantErr(t, err, fmt.Sprintf("extractUnits(%v, \n%v, func())", tt.args.ctx, tt.args.query)) {
				return
			}
			assert.Equalf(t, tt.want, got, "extractUnits(%v, \n%v, func())", tt.args.ctx, tt.args.query)
			assert.Equalf(t, tt.want1, got1, "extractUnits(%v, \n%v, func())", tt.args.ctx, tt.args.query)
		})
	}
}

func TestBrowser_extractAttributes(t *testing.T) {
	type fields struct {
		DefaultBrowser  _default.DefaultBrowser
		connection      plc4go.PlcConnection
		sequenceCounter uint8
	}
	type args struct {
		query *unitInfoQuery
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []readWriteModel.Attribute
		want1  bool
	}{
		{
			name: "one attribute",
			args: args{
				query: &unitInfoQuery{attribute: func() *readWriteModel.Attribute {
					attributeType := readWriteModel.Attribute_Type
					return &attributeType
				}()},
			},
			want:  []readWriteModel.Attribute{readWriteModel.Attribute_Type},
			want1: false,
		},
		{
			name: "all attributes",
			args: args{
				query: &unitInfoQuery{},
			},
			want:  readWriteModel.AttributeValues,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Browser{
				DefaultBrowser:  tt.fields.DefaultBrowser,
				connection:      tt.fields.connection,
				sequenceCounter: tt.fields.sequenceCounter,
				log:             testutils.ProduceTestingLogger(t),
			}
			got, got1 := m.extractAttributes(tt.args.query)
			assert.Equalf(t, tt.want, got, "extractAttributes(\n%v)", tt.args.query)
			assert.Equalf(t, tt.want1, got1, "extractAttributes(\n%v)", tt.args.query)
		})
	}
}

func TestBrowser_getInstalledUnitAddressBytes(t *testing.T) {
	type fields struct {
		DefaultBrowser  _default.DefaultBrowser
		connection      plc4go.PlcConnection
		sequenceCounter uint8
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields)
		want    map[byte]any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "get units",
			args: args{
				ctx: testutils.TestContext(t),
			},
			setup: func(t *testing.T, fields *fields) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				transportUrl := url.URL{Scheme: "test"}
				transportInstance, err := transport.CreateTransportInstance(transportUrl, nil, _options...)
				require.NoError(t, err)
				t.Cleanup(func() {
					assert.NoError(t, transportInstance.Close())
				})
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
				transportInstance.(*test.TransportInstance).SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
					t.Logf("reacting to\n%s", hex.Dump(data))
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
						t.Log("Connection dance done")

						t.Log("Dispatching 3 MMI segments")
						transportInstance.FillReadBuffer([]byte("86020200F900FF0094120006000000000000000008000000000000000000CA\r\n"))
						transportInstance.FillReadBuffer([]byte("86020200F900FF580000000000000000000000000000000000000000000026\r\n"))
						transportInstance.FillReadBuffer([]byte("86020200F700FFB00000000000000000000000000000000000000000D0\r\n"))
					}
				})
				err = transport.AddPreregisteredInstances(transportUrl, transportInstance)
				require.NoError(t, err)
				driver := NewDriver(_options...)
				t.Cleanup(func() {
					assert.NoError(t, driver.Close())
				})
				connectionConnectResult := <-driver.GetConnection(transportUrl, map[string]transports.Transport{"test": transport}, map[string][]string{})
				if err := connectionConnectResult.GetErr(); err != nil {
					t.Error(err)
					t.FailNow()
				}
				fields.connection = connectionConnectResult.GetConnection()
				t.Cleanup(func() {
					timer := time.NewTimer(6 * time.Second)
					t.Cleanup(func() {
						utils.CleanupTimer(timer)
					})
					select {
					case <-fields.connection.Close():
					case <-timer.C:
						t.Error("timeout waiting for connection close")
					}
				})
			},
			want: map[byte]any{
				1:  true,
				2:  true,
				3:  true,
				4:  true,
				6:  true,
				12: true,
				13: true,
				49: true,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			m := Browser{
				DefaultBrowser:  tt.fields.DefaultBrowser,
				connection:      tt.fields.connection,
				sequenceCounter: tt.fields.sequenceCounter,
				log:             testutils.ProduceTestingLogger(t),
			}
			got, err := m.getInstalledUnitAddressBytes(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("getInstalledUnitAddressBytes(%v)", tt.args.ctx)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getInstalledUnitAddressBytes(%v)", tt.args.ctx)
		})
	}
}
