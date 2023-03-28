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
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageCodec_Send(t *testing.T) {
	type fields struct {
		DefaultCodec                  _default.DefaultCodec
		requestContext                readWriteModel.RequestContext
		cbusOptions                   readWriteModel.CBusOptions
		monitoredMMIs                 chan readWriteModel.CALReply
		monitoredSALs                 chan readWriteModel.MonitoredSAL
		lastPackageHash               uint32
		hashEncountered               uint
		currentlyReportedServerErrors uint
	}
	type args struct {
		message spi.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageCodec{
				DefaultCodec:                  tt.fields.DefaultCodec,
				requestContext:                tt.fields.requestContext,
				cbusOptions:                   tt.fields.cbusOptions,
				monitoredMMIs:                 tt.fields.monitoredMMIs,
				monitoredSALs:                 tt.fields.monitoredSALs,
				lastPackageHash:               tt.fields.lastPackageHash,
				hashEncountered:               tt.fields.hashEncountered,
				currentlyReportedServerErrors: tt.fields.currentlyReportedServerErrors,
			}
			tt.wantErr(t, m.Send(tt.args.message), fmt.Sprintf("Send(%v)", tt.args.message))
		})
	}
}

func TestMessageCodec_Receive(t *testing.T) {
	type fields struct {
		DefaultCodec                  _default.DefaultCodec
		requestContext                readWriteModel.RequestContext
		cbusOptions                   readWriteModel.CBusOptions
		monitoredMMIs                 chan readWriteModel.CALReply
		monitoredSALs                 chan readWriteModel.MonitoredSAL
		lastPackageHash               uint32
		hashEncountered               uint
		currentlyReportedServerErrors uint
	}
	tests := []struct {
		name    string
		fields  fields
		want    spi.Message
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "No data",
			fields: fields{
				DefaultCodec: NewMessageCodec(func() transports.TransportInstance {
					transport := test.NewTransport()
					transportInstance := test.NewTransportInstance(transport)
					return transportInstance
				}()),
				requestContext:                nil,
				cbusOptions:                   nil,
				monitoredMMIs:                 nil,
				monitoredSALs:                 nil,
				lastPackageHash:               0,
				hashEncountered:               0,
				currentlyReportedServerErrors: 0,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MessageCodec{
				DefaultCodec:                  tt.fields.DefaultCodec,
				requestContext:                tt.fields.requestContext,
				cbusOptions:                   tt.fields.cbusOptions,
				monitoredMMIs:                 tt.fields.monitoredMMIs,
				monitoredSALs:                 tt.fields.monitoredSALs,
				lastPackageHash:               tt.fields.lastPackageHash,
				hashEncountered:               tt.fields.hashEncountered,
				currentlyReportedServerErrors: tt.fields.currentlyReportedServerErrors,
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
		transport := test.NewTransport()
		transportInstance := test.NewTransportInstance(transport)
		codec := NewMessageCodec(transportInstance)
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		var err error
		msg, err = codec.Receive()
		// No data yet so this should error
		assert.Error(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		err = transportInstance.FillReadBuffer([]byte("i."))
		assert.NoError(t, err)

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
		assert.NoError(t, err)
		assert.Nil(t, msg)

		// Now we fill in the payload
		err = transportInstance.FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))
		assert.NoError(t, err)

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
		assert.NoError(t, err)
		assert.NotNil(t, msg)

		// The message should have a confirmation with an alpha
		assert.True(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
	t.Run("data after 6 times", func(t *testing.T) {
		transport := test.NewTransport()
		transportInstance := test.NewTransportInstance(transport)
		codec := NewMessageCodec(transportInstance)
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		var err error
		msg, err = codec.Receive()
		// No data yet so this should error
		assert.Error(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		err = transportInstance.FillReadBuffer([]byte("i."))
		assert.NoError(t, err)

		for i := 0; i < 8; i++ {
			t.Logf("%d try", i+1)
			// We should wait for more data, so no error, no message
			msg, err = codec.Receive()
			assert.NoError(t, err)
			assert.Nil(t, msg)
		}

		// Now we fill in the payload
		err = transportInstance.FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))
		assert.NoError(t, err)

		// We should wait for more data, so no error, no message
		msg, err = codec.Receive()
		assert.NoError(t, err)
		assert.NotNil(t, msg)

		// The message should have a confirmation with an alpha
		assert.True(t, msg.(readWriteModel.CBusMessageToClient).GetReply().GetIsAlpha())
	})
	t.Run("data after 16 times", func(t *testing.T) {
		transport := test.NewTransport()
		transportInstance := test.NewTransportInstance(transport)
		codec := NewMessageCodec(transportInstance)
		codec.requestContext = readWriteModel.NewRequestContext(true)

		var msg spi.Message
		var err error
		msg, err = codec.Receive()
		// No data yet so this should error
		assert.Error(t, err)
		assert.Nil(t, msg)
		// Now we add a confirmation
		err = transportInstance.FillReadBuffer([]byte("i."))
		assert.NoError(t, err)

		for i := 0; i < 16; i++ {
			t.Logf("%d try", i+1)
			// We should wait for more data, so no error, no message
			msg, err = codec.Receive()
			if i == 15 {
				assert.NoError(t, err)
				assert.NotNil(t, msg)
				// This should be the confirmation only ...
				reply := msg.(readWriteModel.CBusMessageToClient).GetReply()
				assert.True(t, reply.GetIsAlpha())
				// ... and no content
				assert.Nil(t, reply.(readWriteModel.ReplyOrConfirmationConfirmation).GetEmbeddedReply())
			} else {
				assert.NoError(t, err)
				assert.Nil(t, msg)
			}
		}

		// Now we fill in the payload
		err = transportInstance.FillReadBuffer([]byte("86FD0201078900434C495053414C20C2\r\n"))
		assert.NoError(t, err)

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
		want *MessageCodec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMessageCodec(tt.args.transportInstance), "NewMessageCodec(%v)", tt.args.transportInstance)
		})
	}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, extractMMIAndSAL(tt.args.codec, tt.args.message), "extractMMIAndSAL(%v, %v)", tt.args.codec, tt.args.message)
		})
	}
}
