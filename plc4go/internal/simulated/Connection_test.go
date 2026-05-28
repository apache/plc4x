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
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

func TestConnection_Connect(t *testing.T) {
	type fields struct {
		device       *Device
		fieldHandler spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		options      map[string][]string
		connected    bool
		invalidated  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
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
			wantErr: assert.NoError,
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
			wantErr: assert.Error,
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
			err := c.Connect(t.Context())
			tt.wantErr(t, err)
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
		invalidated  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
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
			wantErr: assert.NoError,
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
			wantErr: assert.Error,
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
			err := c.Close()
			tt.wantErr(t, err)
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
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
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
			args: args{
				ctx: t.Context(),
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
			args: args{
				ctx: t.Context(),
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
			args: args{
				ctx: t.Context(),
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
				var wg sync.WaitGroup
				t.Cleanup(wg.Wait)
				wg.Go(func() {
					t.Log(c.Close())
					ch <- true
				})
				return ch
			}
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
			case <-t.Context().Done():
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
			want: &_default.DefaultConnectionMetadata{
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
		invalidated  bool
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
		{
			name: "invalidated",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
				invalidated:  true,
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
			if tt.fields.invalidated {
				c.invalidated.Store(true)
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
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		prepare      func(*Connection)
		wantErr      assert.ErrorAssertionFunc
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
			args: args{
				ctx: t.Context(),
			},
			wantErr:      assert.NoError,
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
			args: args{
				ctx: t.Context(),
			},
			wantErr:      assert.NoError,
			delayAtLeast: 1000,
		},
		{
			name: "invalidated",
			fields: fields{
				device:       NewDevice("hurz"),
				fieldHandler: NewTagHandler(),
				valueHandler: NewValueHandler(),
				options:      map[string][]string{},
				connected:    true,
			},
			args: args{
				ctx: t.Context(),
			},
			prepare: func(c *Connection) {
				c.invalidated.Store(true)
			},
			wantErr:      assert.Error,
			delayAtLeast: 0,
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
			if tt.prepare != nil {
				prepare := tt.prepare
				prepare(c)
			}
			err := c.Ping(tt.args.ctx)
			tt.wantErr(t, err)
		})
	}
}

func TestConnection_Invalidate(t *testing.T) {
	conn := NewConnection(NewDevice("hurz"), NewTagHandler(), NewValueHandler(), map[string][]string{})
	require.NoError(t, conn.Connect(t.Context()))
	conn.Invalidate()
	assert.True(t, conn.IsInvalidated())
	assert.False(t, conn.IsConnected())
	assert.Error(t, conn.Ping(t.Context()))
	require.NoError(t, conn.Close())
	conn.Invalidate()
	assert.True(t, conn.IsInvalidated())
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
