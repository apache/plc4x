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
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestMessageCodec_Send(t *testing.T) {
	type fields struct {
		DefaultCodec   _default.DefaultCodec
		requestContext readWriteModel.RequestContext
		cbusOptions    readWriteModel.CBusOptions
		monitoredMMIs  chan readWriteModel.CALReply
		monitoredSALs  chan readWriteModel.MonitoredSAL
	}
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields, args *args)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "send nothing",
			wantErr: assert.Error,
		},
		{
			name: "a cbus message",
			args: args{
				message: readWriteModel.NewCBusMessageToClient(readWriteModel.NewReplyOrConfirmationConfirmation(
					0x00,
					readWriteModel.NewConfirmation(
						readWriteModel.NewAlpha('!'),
						nil,
						readWriteModel.ConfirmationType_CHECKSUM_FAILURE,
					),
					nil,
				)),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("@A62120\r@A62120\r"))
				codec := NewMessageCodec(ti, _options...)
				require.NoError(t, codec.Connect(t.Context()))
				t.Cleanup(func() {
					assert.NoError(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &MessageCodec{
				DefaultCodec:   tt.fields.DefaultCodec,
				requestContext: tt.fields.requestContext,
				cbusOptions:    tt.fields.cbusOptions,
				monitoredMMIs:  tt.fields.monitoredMMIs,
				monitoredSALs:  tt.fields.monitoredSALs,
			}
			tt.wantErr(t, m.Send(t.Context(), tt.args.message), fmt.Sprintf("Send(%v)", tt.args.message))
		})
	}
}

func TestMessageCodec_Receive(t *testing.T) {
	requestContext := readWriteModel.NewRequestContext(false)
	cbusOptions := readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)

	type fields struct {
		DefaultCodec   _default.DefaultCodec
		requestContext readWriteModel.RequestContext
		cbusOptions    readWriteModel.CBusOptions
		monitoredMMIs  chan readWriteModel.CALReply
		monitoredSALs  chan readWriteModel.MonitoredSAL
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		setup       func(*testing.T, *fields, *args)
		manipulator func(*testing.T, *MessageCodec)
		want        spi.Message
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "No data",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
				var cancelFunc context.CancelFunc
				args.ctx, cancelFunc = context.WithTimeout(args.ctx, 10*time.Millisecond)
				t.Cleanup(cancelFunc)
			},
			wantErr: assert.NoError,
		},
		{
			name: "checksum error",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			want: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewServerErrorReply(
					33,
				),
			),
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("!"))
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			wantErr: assert.NoError,
		},
		{
			name: "A21 echo",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("@A62120\r@A62120\r"))
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			wantErr: assert.NoError,
		},
		{
			name: "garbage",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("what on earth\n\r"))
				codec := NewMessageCodec(ti, _options...)

				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			wantErr: assert.NoError,
		},
		{
			name: "error encountered multiple time",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("AFFE!!!\r"))
				codec := NewMessageCodec(ti, _options...)

				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.hashEncountered.Store(9999)
			},
			want: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewServerErrorReply(
					33,
				),
			),
			wantErr: assert.NoError,
		},
		{
			name: "error encountered and reported multiple time",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.hashEncountered.Store(9999)
				messageCodec.currentlyReportedServerErrors.Store(9999)
			},
			want: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestDirectCommandAccess(
					64,
					nil,
					nil,
					readWriteModel.RequestType_DIRECT_COMMAND,
					readWriteModel.NewRequestTermination(),
					readWriteModel.NewCALDataRecall(
						nil,
						readWriteModel.CALCommandTypeContainer_CALCommandRecall,
						nil,
						readWriteModel.Parameter_UNKNOWN_33,
						1,
					),
					nil,
				),
			),
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("@1A2001!!!\r"))
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			wantErr: assert.NoError,
		},
		{
			name: "mmi",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("86040200F940380001000000000000000008000000000000000000000000FA\r\n"))
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.hashEncountered.Store(9999)
				messageCodec.currentlyReportedServerErrors.Store(9999)
			},
			want: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewReplyOrConfirmationReply(
					56,
					readWriteModel.NewReplyEncodedReply(
						56,
						cbusOptions,
						readWriteModel.NewEncodedReplyCALReply(
							requestContext,
							134,
							readWriteModel.NewCALReplyLong(
								134,
								readWriteModel.NewCALDataStatusExtended(
									requestContext,
									249,
									nil,
									64,
									56,
									0,
									[]readWriteModel.StatusByte{
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_ON,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_OFF,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
										readWriteModel.NewStatusByte(
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
											readWriteModel.GAVState_DOES_NOT_EXIST,
										),
									},
									nil,
								),
								262656,
								readWriteModel.NewUnitAddress(4),
								nil,
								readWriteModel.NewSerialInterfaceAddress(2),
								utils.ToPtr(byte(0)),
								nil,
							),
						),
						nil,
					),
					readWriteModel.NewResponseTermination(),
				),
			),
			wantErr: assert.NoError,
		},
		{
			name: "sal",
			fields: fields{
				requestContext: requestContext,
				cbusOptions:    cbusOptions,
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			args: args{
				ctx: t.Context(),
			},
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.hashEncountered.Store(9999)
				messageCodec.currentlyReportedServerErrors.Store(9999)
			},
			setup: func(t *testing.T, fields *fields, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)

				transport := test.NewTransport(_options...)
				ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
				require.NoError(t, err)

				require.NoError(t, ti.Connect(t.Context()))
				ti.(*test.TransportInstance).FillReadBuffer([]byte("0531AC0079042F0401430316000011\r\n"))
				codec := NewMessageCodec(ti, _options...)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			want: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewReplyOrConfirmationReply(
					48,
					readWriteModel.NewReplyEncodedReply(
						48,
						cbusOptions,
						readWriteModel.NewMonitoredSALReply(
							requestContext,
							5,
							readWriteModel.NewMonitoredSALLongFormSmartMode(
								5,
								3255296,
								readWriteModel.NewUnitAddress(49),
								nil,
								172,
								utils.ToPtr(byte(0)),
								nil,
								readWriteModel.NewSALDataAirConditioning(
									readWriteModel.NewSALDataAirConditioning(
										nil,
										readWriteModel.NewAirConditioningDataSetZoneHvacMode(
											readWriteModel.AirConditioningCommandTypeContainer_AirConditioningCommandSetZoneHvacMode,
											4,
											readWriteModel.NewHVACZoneList(false, false, false, false, false, false, false, true),
											readWriteModel.NewHVACModeAndFlags(true, false, false, false, readWriteModel.HVACModeAndFlagsMode_HEAT_AND_COOL),
											3,
											readWriteModel.NewHVACTemperature(5632),
											nil,
											readWriteModel.NewHVACAuxiliaryLevel(false, 0),
										),
									),
									readWriteModel.NewAirConditioningDataSetZoneGroupOn(
										readWriteModel.AirConditioningCommandTypeContainer_AirConditioningCommandSetZoneGroupOn,
										4,
									),
								),
							),
						),
						nil,
					),
					readWriteModel.NewResponseTermination(),
				),
			),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args)
			}
			m := &MessageCodec{
				DefaultCodec:   tt.fields.DefaultCodec,
				requestContext: tt.fields.requestContext,
				cbusOptions:    tt.fields.cbusOptions,
				monitoredMMIs:  tt.fields.monitoredMMIs,
				monitoredSALs:  tt.fields.monitoredSALs,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			got, err := m.Receive(tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("Receive()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Receive()")
		})
	}
}

func TestMessageCodec_Receive_Delayed_Response(t *testing.T) {
	t.Run("instant data", func(t *testing.T) {
		_options := testutils.EnrichOptionsWithOptionsForTesting(t)

		transport := test.NewTransport(_options...)
		ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
		require.NoError(t, err)
		require.NoError(t, ti.Connect(t.Context()))
		codec := NewMessageCodec(ti, _options...)
		t.Cleanup(func() {
			assert.Error(t, codec.Disconnect())
		})
		codec.requestContext = readWriteModel.NewRequestContext(true)

		timeoutCtx := func(timeout time.Duration) context.Context {
			withTimeout, cancelFunc := context.WithTimeout(t.Context(), timeout)
			t.Cleanup(cancelFunc)
			return withTimeout
		}

		var msg spi.Message
		msg, err = codec.Receive(timeoutCtx(1 * time.Second))
		// No data yet so this should return no error and no data
		assert.NoError(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		ti.(*test.TransportInstance).FillReadBuffer([]byte("i."))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive(timeoutCtx(1 * time.Second))
		assert.NoError(t, err)
		assert.Nil(t, msg)

		// Now we fill in the payload
		ti.(*test.TransportInstance).FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))

		// We should wait for more data, so no error
		msg, err = codec.Receive(timeoutCtx(2 * time.Second))
		assert.NoError(t, err)
		require.NotNil(t, msg)

		// The message should have a confirmation with an alpha
		assert.True(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
	t.Run("data after 6 times", func(t *testing.T) {
		_options := testutils.EnrichOptionsWithOptionsForTesting(t)

		transport := test.NewTransport(_options...)
		ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
		require.NoError(t, err)

		require.NoError(t, ti.Connect(t.Context()))
		codec := NewMessageCodec(ti, _options...)
		t.Cleanup(func() {
			assert.Error(t, codec.Disconnect())
		})
		codec.requestContext = readWriteModel.NewRequestContext(true)

		canceledCtx := func() context.Context {
			ctx, cancelFunc := context.WithCancel(t.Context())
			cancelFunc()
			return ctx
		}
		var msg spi.Message
		msg, err = codec.Receive(canceledCtx())
		// No data yet so this should return no error and no data
		assert.NoError(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		ti.(*test.TransportInstance).FillReadBuffer([]byte("i."))

		for i := 0; i < 8; i++ {
			t.Logf("%d try", i+1)
			// We should wait for more data, so no error, no message
			msg, err = codec.Receive(canceledCtx())
			assert.NoError(t, err)
			assert.Nil(t, msg)
		}

		// Now we fill in the payload
		ti.(*test.TransportInstance).FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive(t.Context())
		assert.NoError(t, err)
		assert.NotNil(t, msg)

		// The message should have a confirmation with an alpha
		assert.True(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
	t.Run("data after 15 times", func(t *testing.T) {
		_options := testutils.EnrichOptionsWithOptionsForTesting(t)

		transport := test.NewTransport(_options...)
		ti, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, map[string][]string{"simulatedLatency": {"1ms"}}, _options...)
		require.NoError(t, err)

		require.NoError(t, ti.Connect(t.Context()))
		codec := NewMessageCodec(ti, _options...)
		t.Cleanup(func() {
			assert.Error(t, codec.Disconnect())
		})
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		msg, err = codec.Receive(t.Context())
		// No data yet so this should return no error and no data
		assert.NoError(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		ti.(*test.TransportInstance).FillReadBuffer([]byte("i."))

		for i := 0; i <= 15; i++ {
			t.Logf("%d try", i+1)
			// We should wait for more data, so no error, no message
			msg, err = codec.Receive(t.Context())
			if i == 15 {
				require.NoError(t, err)
				require.NotNil(t, msg)
				// This should be the confirmation only ...
				reply := msg.(readWriteModel.CBusMessageToClient).GetReply()
				assert.True(t, reply.GetIsAlpha())
				// ... and no content
				assert.Nil(t, reply.(readWriteModel.ReplyOrConfirmationConfirmation).GetEmbeddedReply())
			} else {
				assert.NoError(t, err)
				assert.Nil(t, msg, "Got message at %d try", i+1)
			}
		}

		// Now we fill in the payload
		ti.(*test.TransportInstance).FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive(t.Context())
		assert.NoError(t, err)
		assert.NotNil(t, msg)

		// The message should have a confirmation without an alpha
		assert.False(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
}

func TestNewMessageCodec(t *testing.T) {
	type args struct {
		transportInstance transports.TransportInstance
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "create it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := NewMessageCodec(tt.args.transportInstance)
			t.Cleanup(func() {
				assert.Error(t, codec.Disconnect())
			})
			assert.NotNilf(t, codec, "NewMessageCodec(%v)", tt.args.transportInstance)
		})
	}
}

func TestMessageCodec_GetCodec(t *testing.T) {
	// just a useless test...
	(&MessageCodec{}).GetCodec()
}

func Test_extractMMIAndSAL(t *testing.T) {
	type args struct {
		codec   _default.DefaultCodecRequirements
		message spi.Message
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T, args *args)
		handled bool
	}{
		{
			name: "extract it",
		},
		{
			name: "monitored sal",
			args: args{
				message: readWriteModel.NewCBusMessageToClient(
					readWriteModel.NewReplyOrConfirmationReply(
						0,
						readWriteModel.NewReplyEncodedReply(
							0,
							readWriteModel.NewCBusOptionsBuilder().MustBuild(),
							readWriteModel.NewMonitoredSALReply(
								readWriteModel.NewRequestContext(false),
								0,
								readWriteModel.NewMonitoredSALShortFormBasicMode(
									0,
									0,
									nil,
									nil,
									nil,
									readWriteModel.ApplicationIdContainer_RESERVED_00,
									nil,
								),
							),
							nil,
						),
						readWriteModel.NewResponseTermination(),
					),
				),
			},
			setup: func(t *testing.T, args *args) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				codec := NewMessageCodec(nil, _options...)
				codec.monitoredSALs = make(chan readWriteModel.MonitoredSAL, 1)
				args.codec = codec
			},
			handled: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
			assert.Equalf(t, tt.handled, extractMMIAndSAL(testutils.ProduceTestingLogger(t))(t.Context(), tt.args.codec, tt.args.message), "extractMMIAndSAL(%v, %v) to be handled", tt.args.codec, tt.args.message)
		})
	}
}
