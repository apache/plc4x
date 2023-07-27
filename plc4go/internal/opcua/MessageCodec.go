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

package opcua

import (
	"context"
	"encoding/binary"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	channel *SecureChannel

	stateChange sync.Mutex

	passLogToModel bool           `ignore:"true"`
	log            zerolog.Logger `ignore:"true"`
}

func NewMessageCodec(transportInstance transports.TransportInstance, channel *SecureChannel, _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		channel:        channel,
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Msgf("Sending message\n%s", message)
	// Cast the message to the correct type of struct
	messagePdu, ok := message.(readWriteModel.MessagePDU)
	if !ok {
		return errors.Errorf("Invalid message type %T", message)
	}

	// Serialize the request
	theBytes, err := messagePdu.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(theBytes)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive() (spi.Message, error) {
	m.log.Trace().Msg("Receive")
	ti := m.GetTransportInstance()
	if !ti.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}

	if err := ti.FillBuffer(
		func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
			numBytesAvailable, err := ti.GetNumBytesAvailableInBuffer()
			if err != nil {
				return false
			}
			return numBytesAvailable < 8
		}); err != nil {
		m.log.Warn().Err(err).Msg("error filling buffer")
	}

	data, err := ti.PeekReadableBytes(8)
	if err != nil {
		m.log.Debug().Err(err).Msg("error peeking")
		return nil, nil
	}
	numberOfBytesToRead := binary.LittleEndian.Uint32(data[:4])
	readBytes, err := ti.Read(numberOfBytesToRead)
	if err != nil {
		return nil, errors.Wrapf(err, "could not read %d bytes", readBytes)
	}
	ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
	messagePdu, err := readWriteModel.MessagePDUParse(ctxForModel, readBytes, true)
	if err != nil {
		return nil, errors.New("Could not parse pdu")
	}

	return messagePdu, nil
}
