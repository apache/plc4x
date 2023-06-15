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
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
			args: args{message: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewReplyOrConfirmationConfirmation(
					readWriteModel.NewConfirmation(readWriteModel.NewAlpha('!'), nil, readWriteModel.ConfirmationType_CHECKSUM_FAILURE), nil, 0x00, nil, nil,
				), nil, nil,
			)},
			setup: func(t *testing.T, fields *fields, args *args) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				codec := NewMessageCodec(instance, loggerOption)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
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
			tt.wantErr(t, m.Send(tt.args.message), fmt.Sprintf("Send(%v)", tt.args.message))
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
	tests := []struct {
		name        string
		fields      fields
		setup       func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, messageCodec *MessageCodec)
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
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				codec := NewMessageCodec(instance, loggerOption)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
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
			want: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewServerErrorReply(
					33, cbusOptions, requestContext,
				),
				requestContext, cbusOptions,
			),
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("!"))
				codec := NewMessageCodec(instance, loggerOption)
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
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("@A62120\r@A62120\r"))
				codec := NewMessageCodec(instance, loggerOption)
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
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("what on earth\n\r"))
				codec := NewMessageCodec(instance, loggerOption)
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
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("AFFE!!!\r"))
				codec := NewMessageCodec(instance, loggerOption)
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
					33, cbusOptions, requestContext,
				),
				requestContext, cbusOptions,
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
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.hashEncountered.Store(9999)
				messageCodec.currentlyReportedServerErrors.Store(9999)
			},
			want: readWriteModel.NewCBusMessageToServer(
				readWriteModel.NewRequestDirectCommandAccess(
					readWriteModel.NewCALDataRecall(
						readWriteModel.Parameter_UNKNOWN_33,
						1,
						readWriteModel.CALCommandTypeContainer_CALCommandRecall,
						nil,
						nil,
					),
					nil,
					readWriteModel.RequestType_DIRECT_COMMAND,
					nil,
					nil,
					64,
					readWriteModel.NewRequestTermination(),
					cbusOptions,
				),
				requestContext, cbusOptions,
			),
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("@1A2001!!!\r"))
				codec := NewMessageCodec(instance, loggerOption)
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
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("86040200F940380001000000000000000008000000000000000000000000FA\r\n"))
				codec := NewMessageCodec(instance, loggerOption)
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
					readWriteModel.NewReplyEncodedReply(
						readWriteModel.NewEncodedReplyCALReply(
							readWriteModel.NewCALReplyLong(
								262656,
								readWriteModel.NewUnitAddress(4),
								nil,
								readWriteModel.NewSerialInterfaceAddress(2),
								func() *byte {
									var b byte = 0
									return &b
								}(),
								nil,
								134,
								readWriteModel.NewCALDataStatusExtended(
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
									249,
									nil,
									requestContext,
								),
								cbusOptions,
								requestContext,
							),
							134,
							cbusOptions,
							requestContext,
						),
						nil,
						56,
						cbusOptions,
						requestContext,
					),
					readWriteModel.NewResponseTermination(),
					56,
					cbusOptions,
					requestContext,
				),
				requestContext, cbusOptions,
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
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.hashEncountered.Store(9999)
				messageCodec.currentlyReportedServerErrors.Store(9999)
			},
			setup: func(t *testing.T, fields *fields) {
				// Setup logger
				logger := testutils.ProduceTestingLogger(t)

				loggerOption := options.WithCustomLogger(logger)

				transport := test.NewTransport(loggerOption)
				instance := test.NewTransportInstance(transport, loggerOption)
				instance.FillReadBuffer([]byte("0531AC0079042F0401430316000011\r\n"))
				codec := NewMessageCodec(instance, loggerOption)
				t.Cleanup(func() {
					assert.Error(t, codec.Disconnect())
				})
				fields.DefaultCodec = codec
			},
			want: readWriteModel.NewCBusMessageToClient(
				readWriteModel.NewReplyOrConfirmationReply(
					readWriteModel.NewReplyEncodedReply(
						readWriteModel.NewMonitoredSALReply(
							readWriteModel.NewMonitoredSALLongFormSmartMode(
								3255296,
								readWriteModel.NewUnitAddress(49),
								nil,
								172,
								func() *byte {
									var b byte = 0
									return &b
								}(),
								nil,
								readWriteModel.NewSALDataAirConditioning(
									readWriteModel.NewAirConditioningDataSetZoneGroupOn(
										4,
										readWriteModel.AirConditioningCommandTypeContainer_AirConditioningCommandSetZoneGroupOn,
									),
									readWriteModel.NewSALDataAirConditioning(
										readWriteModel.NewAirConditioningDataSetZoneHvacMode(
											4,
											readWriteModel.NewHVACZoneList(false, false, false, false, false, false, false, true),
											readWriteModel.NewHVACModeAndFlags(true, false, false, false, readWriteModel.HVACModeAndFlagsMode_HEAT_AND_COOL),
											3,
											readWriteModel.NewHVACTemperature(5632),
											nil,
											readWriteModel.NewHVACAuxiliaryLevel(false, 0),
											readWriteModel.AirConditioningCommandTypeContainer_AirConditioningCommandSetZoneHvacMode,
										),
										nil,
									),
								),
								5,
								cbusOptions,
							),
							5,
							cbusOptions,
							requestContext,
						),
						nil,
						48,
						cbusOptions,
						requestContext,
					),
					readWriteModel.NewResponseTermination(),
					48,
					cbusOptions,
					requestContext,
				),
				requestContext, cbusOptions,
			),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
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
			got, err := m.Receive()
			if !tt.wantErr(t, err, fmt.Sprintf("Receive()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Receive()")
		})
	}
}

func TestMessageCodec_Receive_Delayed_Response(t *testing.T) {
	t.Run("instant data", func(t *testing.T) {
		// Setup logger
		logger := testutils.ProduceTestingLogger(t)

		loggerOption := options.WithCustomLogger(logger)

		transport := test.NewTransport(loggerOption)
		transportInstance := test.NewTransportInstance(transport, loggerOption)
		codec := NewMessageCodec(transportInstance, loggerOption)
		t.Cleanup(func() {
			assert.Error(t, codec.Disconnect())
		})
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		var err error
		msg, err = codec.Receive()
		// No data yet so this should return no error and no data
		assert.NoError(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		transportInstance.FillReadBuffer([]byte("i."))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
		assert.NoError(t, err)
		assert.Nil(t, msg)

		// Now we fill in the payload
		transportInstance.FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
		assert.NoError(t, err)
		require.NotNil(t, msg)

		// The message should have a confirmation with an alpha
		assert.True(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
	t.Run("data after 6 times", func(t *testing.T) {
		// Setup logger
		logger := testutils.ProduceTestingLogger(t)

		loggerOption := options.WithCustomLogger(logger)

		transport := test.NewTransport(loggerOption)
		transportInstance := test.NewTransportInstance(transport, loggerOption)
		codec := NewMessageCodec(transportInstance, loggerOption)
		t.Cleanup(func() {
			assert.Error(t, codec.Disconnect())
		})
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		var err error
		msg, err = codec.Receive()
		// No data yet so this should return no error and no data
		assert.NoError(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		transportInstance.FillReadBuffer([]byte("i."))

		for i := 0; i < 8; i++ {
			t.Logf("%d try", i+1)
			// We should wait for more data, so no error, no message
			msg, err = codec.Receive()
			assert.NoError(t, err)
			assert.Nil(t, msg)
		}

		// Now we fill in the payload
		transportInstance.FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
		assert.NoError(t, err)
		assert.NotNil(t, msg)

		// The message should have a confirmation with an alpha
		assert.True(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
	t.Run("data after 15 times", func(t *testing.T) {
		// Setup logger
		logger := testutils.ProduceTestingLogger(t)

		loggerOption := options.WithCustomLogger(logger)

		transport := test.NewTransport(loggerOption)
		transportInstance := test.NewTransportInstance(transport, loggerOption)
		codec := NewMessageCodec(transportInstance, loggerOption)
		t.Cleanup(func() {
			assert.Error(t, codec.Disconnect())
		})
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		var err error
		msg, err = codec.Receive()
		// No data yet so this should return no error and no data
		assert.NoError(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		transportInstance.FillReadBuffer([]byte("i."))

		for i := 0; i <= 15; i++ {
			t.Logf("%d try", i+1)
			// We should wait for more data, so no error, no message
			msg, err = codec.Receive()
			if i == 15 {
				assert.NoError(t, err)
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
		transportInstance.FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
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
		name string
		args args
		want bool
	}{
		{
			name: "extract it",
		},
		{
			name: "monitored sal",
			args: args{
				codec: NewMessageCodec(nil),
				message: readWriteModel.NewCBusMessageToClient(
					readWriteModel.NewReplyOrConfirmationReply(
						readWriteModel.NewReplyEncodedReply(
							readWriteModel.NewMonitoredSALReply(
								nil,
								0,
								nil,
								nil,
							),
							nil,
							0,
							nil,
							nil,
						),
						nil,
						0,
						nil,
						nil,
					),
					nil,
					nil,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, extractMMIAndSAL(testutils.ProduceTestingLogger(t))(tt.args.codec, tt.args.message), "extractMMIAndSAL(%v, %v)", tt.args.codec, tt.args.message)
		})
	}
}

func TestMessageCodec_String(t *testing.T) {
	type fields struct {
		DefaultCodec   _default.DefaultCodec
		requestContext readWriteModel.RequestContext
		cbusOptions    readWriteModel.CBusOptions
		monitoredMMIs  chan readWriteModel.CALReply
		monitoredSALs  chan readWriteModel.MonitoredSAL
		log            zerolog.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, messageCodec *MessageCodec)
		want        string
	}{
		{
			name: "string it",
			fields: fields{
				DefaultCodec:   _default.NewDefaultCodec(nil, test.NewTransportInstance(test.NewTransport())),
				requestContext: readWriteModel.NewRequestContext(true),
				cbusOptions:    readWriteModel.NewCBusOptions(true, true, true, true, true, true, true, true, true),
				monitoredMMIs:  nil,
				monitoredSALs:  nil,
			},
			manipulator: func(t *testing.T, messageCodec *MessageCodec) {
				messageCodec.lastPackageHash.Store(2)
				messageCodec.hashEncountered.Store(3)
				messageCodec.currentlyReportedServerErrors.Store(4)
			},
			want: `
╔═MessageCodec════════════════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═defaultCodec═════════════════════════════════════════════════════════════════════════╗                             ║
║║╔═transportInstance╗╔═defaultIncomingMessageChannel╗╔═customMessageHandling╗╔═running╗║                             ║
║║║       test       ║║         0 element(s)         ║║       b0 false       ║║b0 false║║                             ║
║║╚══════════════════╝╚══════════════════════════════╝╚══════════════════════╝╚════════╝║                             ║
║╚══════════════════════════════════════════════════════════════════════════════════════╝                             ║
║╔═requestContext/RequestContext/sendIdentifyRequestBefore═════════════════════════════════════════════════════╗      ║
║║                                                   b1 true                                                   ║      ║
║╚═════════════════════════════════════════════════════════════════════════════════════════════════════════════╝      ║
║╔═cbusOptions/CBusOptions═══════════════════════════════════════════════════════════╗╔═monitoredMMIs╗╔═monitoredSALs╗║
║║╔═connect╗╔═smart═╗╔═idmon═╗╔═exstat╗╔═monitor╗╔═monall╗╔═pun═══╗╔═pcn═══╗╔═srchk═╗║║ 0 element(s) ║║ 0 element(s) ║║
║║║b1 true ║║b1 true║║b1 true║║b1 true║║b1 true ║║b1 true║║b1 true║║b1 true║║b1 true║║╚══════════════╝╚══════════════╝║
║║╚════════╝╚═══════╝╚═══════╝╚═══════╝╚════════╝╚═══════╝╚═══════╝╚═══════╝╚═══════╝║                                ║
║╚═══════════════════════════════════════════════════════════════════════════════════╝                                ║
║╔═lastPackageHash╗╔═hashEncountered════╗╔═currentlyReportedServerErrors╗                                             ║
║║  0x00000002 2  ║║0x0000000000000003 3║║     0x0000000000000004 4     ║                                             ║
║╚════════════════╝╚════════════════════╝╚══════════════════════════════╝                                             ║
╚═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageCodec{
				DefaultCodec:   tt.fields.DefaultCodec,
				requestContext: tt.fields.requestContext,
				cbusOptions:    tt.fields.cbusOptions,
				monitoredMMIs:  tt.fields.monitoredMMIs,
				monitoredSALs:  tt.fields.monitoredSALs,
				log:            tt.fields.log,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, m)
			}
			assert.Equalf(t, tt.want, m.String(), "String()")
		})
	}
}
