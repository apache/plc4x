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

package simulated

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestConnection_Connect(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name         string
		fields       fields
		want         plc4go.PlcConnectionConnectResult
		delayAtLeast time.Duration
		wantErr      bool
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    false,
			},
			want: _default.NewDefaultPlcConnectionConnectResult(&Connection{
				device:       NewDevice("hurz"),
				tagHandler:   NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			}, nil),
			delayAtLeast: 0,
			wantErr:      false,
		},
		// If the connection was already connected, the
		// connection should fail with an error.
		{
			name: "already connected",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			want: _default.NewDefaultPlcConnectionConnectResult(&Connection{
				device:       NewDevice("hurz"),
				tagHandler:   NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			}, nil),
			delayAtLeast: 0,
			wantErr:      true,
		},
		// If the connection should simulate a delay, make sure it doesn't
		// return immediately.
		{
			name: "delayed connected",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options: map[string][]string{
					"connectionDelay": {"1000"},
				},
				connected: false,
			},
			want: _default.NewDefaultPlcConnectionConnectResult(&Connection{
				device:       NewDevice("hurz"),
				tagHandler:   NewTagHandler(),
				valueHandler: NewValueHandler(),
				options: map[string][]string{
					"connectionDelay": {"1000"},
				},
				connected: true,
			}, nil),
			delayAtLeast: time.Second * 1,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			timeBeforeConnect := time.Now()
			connectionChan := c.Connect()
			timeout := time.NewTimer(3 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case connectResult := <-connectionChan:
				timeAfterConnect := time.Now()
				// If an expected delay was defined, check if connecting
				// took at least this long.
				if tt.delayAtLeast > 0 {
					connectionTime := timeAfterConnect.Sub(timeBeforeConnect)
					if connectionTime < tt.delayAtLeast {
						t.Errorf("TestConnection.Connect() connected too fast. Expected at least %v but connected after %v", tt.delayAtLeast, connectionTime)
					}
				}
				// If we wanted an error, but didn't get one or the other way around.
				if tt.wantErr != (connectResult.GetErr() != nil) {
					t.Errorf("TestConnection.Connect() hasErr= %v, wantErr %v", connectResult.GetErr() != nil, tt.wantErr)
				} else if !tt.wantErr {
					// Check if we're connected.
					if !assert.Equal(t, tt.want, connectResult) {
						t.Errorf("TestConnection.Connect() = %v, want %v", connectResult, tt.want)
					}
				}
			case <-timeout.C:
				t.Errorf("TestConnection.Connect() got timeout")
			}
		})
	}
}

func TestConnection_Close(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name         string
		fields       fields
		want         plc4go.PlcConnectionCloseResult
		delayAtLeast time.Duration
		wantErr      bool
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			want: _default.NewDefaultPlcConnectionCloseResult(&Connection{
				device:       NewDevice("hurz"),
				tagHandler:   NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    false,
			}, nil),
			delayAtLeast: 0,
			wantErr:      false,
		},
		{
			name: "not connected",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    false,
			},
			want: _default.NewDefaultPlcConnectionCloseResult(&Connection{
				device:       NewDevice("hurz"),
				tagHandler:   NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    false,
			}, nil),
			delayAtLeast: 0,
			wantErr:      true,
		},
		{
			name: "delayed close",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options: map[string][]string{
					"closingDelay": {"1000"},
				},
				connected: true,
			},
			want: _default.NewDefaultPlcConnectionCloseResult(&Connection{
				device:       NewDevice("hurz"),
				tagHandler:   NewTagHandler(),
				valueHandler: NewValueHandler(),
				options: map[string][]string{
					"closingDelay": {"1000"},
				},
				connected: false,
			}, nil),
			delayAtLeast: 1000,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			timeBeforeClose := time.Now()
			closeChan := c.Close()
			timeout := time.NewTimer(3 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case closeResult := <-closeChan:
				timeAfterClose := time.Now()
				// If an expected delay was defined, check if closing
				// took at least this long.
				if tt.delayAtLeast > 0 {
					connectionTime := timeAfterClose.Sub(timeBeforeClose)
					if connectionTime < tt.delayAtLeast {
						t.Errorf("TestConnection.Close() connected too fast. Expected at least %v but connected after %v", tt.delayAtLeast, connectionTime)
					}
				}
				// If we wanted an error, but didn't get one or the other way around.
				if tt.wantErr != (closeResult.GetErr() != nil) {
					t.Errorf("TestConnection.Close() hasErr= %v, wantErr %v", closeResult.GetErr() != nil, tt.wantErr)
				} else if !tt.wantErr {
					if !assert.Equal(t, tt.want, closeResult) {
						t.Errorf("TestConnection.Close() = %v, want %v", closeResult, tt.want)
					}
				}
			case <-timeout.C:
				t.Errorf("TestConnection.Close() got timeout")
			}
		})
	}
}

func TestConnection_BlockingClose(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name         string
		fields       fields
		delayAtLeast time.Duration
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			delayAtLeast: 0,
		},
		{
			name: "not connected",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    false,
			},
			delayAtLeast: 0,
		},
		{
			name: "delayed close",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options: map[string][]string{
					"closingDelay": {"1000"},
				},
				connected: true,
			},
			delayAtLeast: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			timeBeforeClose := time.Now()
			executor := func() <-chan bool {
				ch := make(chan bool)
				go func() {
					c.BlockingClose()
					ch <- true
				}()
				return ch
			}
			timeout := time.NewTimer(3 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case <-executor():
				timeAfterClose := time.Now()
				// If an expected delay was defined, check if closing
				// took at least this long.
				if tt.delayAtLeast > 0 {
					connectionTime := timeAfterClose.Sub(timeBeforeClose)
					if connectionTime < tt.delayAtLeast {
						t.Errorf("TestConnection.Close() connected too fast. Expected at least %v but connected after %v", tt.delayAtLeast, connectionTime)
					}
				}
			case <-timeout.C:
				t.Errorf("TestConnection.Close() got timeout")
			}
		})
	}
}

func TestConnection_GetMetadata(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcConnectionMetadata
	}{
		{
			name:   "simple",
			fields: fields{},
			want: _default.DefaultConnectionMetadata{
				ConnectionAttributes: map[string]string{
					"connectionDelay": "Delay applied when connecting",
					"closingDelay":    "Delay applied when closing the connection",
					"pingDelay":       "Delay applied when executing a ping operation",
					"readDelay":       "Delay applied when executing a read operation",
					"writeDelay":      "Delay applied when executing a write operation",
				},
				ProvidesReading:     true,
				ProvidesWriting:     true,
				ProvidesSubscribing: false,
				ProvidesBrowsing:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			if got := c.GetMetadata(); !assert.Equal(t, tt.want, got) {
				t.Errorf("GetMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnection_IsConnected(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			want: true,
		},
		{
			name: "not connected",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			if got := c.IsConnected(); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnection_Ping(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name         string
		fields       fields
		want         plc4go.PlcConnectionPingResult
		delayAtLeast time.Duration
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			want:         _default.NewDefaultPlcConnectionPingResult(nil),
			delayAtLeast: 0,
		},
		{
			name: "delayed ping",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options: map[string][]string{
					"pingDelay": {"1000"},
				},
				connected: true,
			},
			want:         _default.NewDefaultPlcConnectionPingResult(nil),
			delayAtLeast: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			timeBeforePing := time.Now()
			pingChan := c.Ping()
			timeout := time.NewTimer(3 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case pingResult := <-pingChan:
				timeAfterPing := time.Now()
				// If an expected delay was defined, check if closing
				// took at least this long.
				if tt.delayAtLeast > 0 {
					pingTime := timeAfterPing.Sub(timeBeforePing)
					if pingTime < tt.delayAtLeast {
						t.Errorf("TestConnection.Ping() completed too fast. Expected at least %v but returned after %v", tt.delayAtLeast, pingTime)
					}
				}
				if !assert.Equal(t, tt.want, pingResult) {
					t.Errorf("TestConnection.Ping() = %v, want %v", pingResult, tt.want)
				}
			case <-timeout.C:
				t.Errorf("TestConnection.Ping() got timeout")
			}
		})
	}
}

func TestConnection_BrowseRequestBuilder(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			defer func() {
				if r := recover(); tt.wantErr && r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			c.BrowseRequestBuilder()
		})
	}
}

func TestConnection_ReadRequestBuilder(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcReadRequestBuilder
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			want: spiModel.NewDefaultPlcReadRequestBuilder(NewTagHandler(), NewReader(NewDevice("hurz"), map[string][]string{}, nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			if got := c.ReadRequestBuilder(); !assert.Equal(t, tt.want, got) {
				t.Errorf("ReadRequestBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConnection_SubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			assert.NotNil(t, c.SubscriptionRequestBuilder())
		})
	}
}

func TestConnection_UnsubscriptionRequestBuilder(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			defer func() {
				if r := recover(); tt.wantErr && r == nil {
					t.Errorf("The code did not panic")
				}
			}()
			c.UnsubscriptionRequestBuilder()
		})
	}
}

func TestConnection_WriteRequestBuilder(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcWriteRequestBuilder
	}{
		{
			name: "simple",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			want: spiModel.NewDefaultPlcWriteRequestBuilder(NewTagHandler(), NewValueHandler(), NewWriter(NewDevice("hurz"), map[string][]string{}, nil)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Connection{
				device:       tt.fields.device,
				tagHandler:   tt.fields.fieldHandler,
				valueHandler: tt.fields.valueHandler,
				options:      tt.fields.options,
				connected:    tt.fields.connected,
			}
			if got := c.WriteRequestBuilder(); !assert.Equal(t, tt.want, got) {
				t.Errorf("WriteRequestBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
