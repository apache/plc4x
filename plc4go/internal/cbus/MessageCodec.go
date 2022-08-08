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
	"bufio"
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/default"
	"github.com/apache/plc4x/plc4go/internal/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	readwriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"hash/crc32"
)

type MessageCodec struct {
	_default.DefaultCodec

	requestContext readwriteModel.RequestContext
	cbusOptions    readwriteModel.CBusOptions

	monitoredSALs   chan readwriteModel.MonitoredSAL
	lastPackageHash uint32
	hashEncountered uint
}

func NewMessageCodec(transportInstance transports.TransportInstance) *MessageCodec {
	codec := &MessageCodec{
		requestContext: readwriteModel.NewRequestContext(false),
		cbusOptions:    readwriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false),
		monitoredSALs:  make(chan readwriteModel.MonitoredSAL, 100),
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _default.WithCustomMessageHandler(func(codec _default.DefaultCodecRequirements, message spi.Message) bool {
		switch message := message.(type) {
		case readwriteModel.CBusMessageToClientExactly:
			switch reply := message.GetReply().(type) {
			case readwriteModel.ReplyOrConfirmationReplyExactly:
				switch reply := reply.GetReply().(type) {
				case readwriteModel.ReplyEncodedReplyExactly:
					switch encodedReply := reply.GetEncodedReply().(type) {
					case readwriteModel.MonitoredSALReplyExactly:
						codec.(*MessageCodec).monitoredSALs <- encodedReply.GetMonitoredSAL()
						return true
					}
				}
			}
		}
		return false
	}))
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	cbusMessage := message.(readwriteModel.CBusMessage)

	// Set the right request context
	m.requestContext = CreateRequestContext(cbusMessage)

	// Serialize the request
	wb := utils.NewWriteBufferByteBased()
	err := cbusMessage.Serialize(wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(wb.GetBytes())
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Receive() (spi.Message, error) {
	log.Trace().Msg("receiving")

	ti := m.GetTransportInstance()
	if err := ti.FillBuffer(func(_ uint, currentByte byte, reader *bufio.Reader) bool {
		hitCr := currentByte == '\r'
		if hitCr {
			// Make sure we peek one more
			_, _ = reader.Peek(1)
			return false
		}
		return true
	}); err != nil {
		return nil, err
	}
	readableBytes, err := ti.GetNumBytesAvailableInBuffer()
	if err != nil {
		log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	if readableBytes == 0 {
		log.Trace().Msg("Nothing to read")
		return nil, nil
	}
	// TODO: we might get a simple confirmation like g# without anything other... so we might need to handle that

	peekedBytes, err := ti.PeekReadableBytes(readableBytes)
	pciResponse, requestToPci := false, false
	indexOfCR := -1
	indexOfLF := -1
lookingForTheEnd:
	for i, peekedByte := range peekedBytes {
		switch peekedByte {
		case '\r':
			if indexOfCR >= 0 {
				// We found the next <cr> so we know we have a package
				requestToPci = true
				break lookingForTheEnd
			}
			indexOfCR = i
		case '\n':
			indexOfLF = i
			// If we found a <nl> we know definitely that we hit the end of a message
			break lookingForTheEnd
		}
	}
	if indexOfCR < 0 && indexOfLF >= 0 {
		// This means that the package is garbage as a lf is always prefixed with a cr
		log.Debug().Err(err).Msg("Error reading")
		// TODO: Possibly clean up ...
		return nil, nil
	}
	if indexOfCR+1 == indexOfLF {
		// This means a <cr> is directly followed by a <lf> which means that we know for sure this is a response
		pciResponse = true
	}
	if !pciResponse && !requestToPci && indexOfLF < 0 {
		// To be sure we might receive that package later we hash the bytes and check if we might receive one
		hash := crc32.NewIEEE()
		_, _ = hash.Write(peekedBytes)
		newPackageHash := hash.Sum32()
		if newPackageHash == m.lastPackageHash {
			m.hashEncountered++
		}
		m.lastPackageHash = newPackageHash
		if m.hashEncountered < 9 {
			return nil, nil
		} else {
			// after 90ms we give up finding a lf
			m.lastPackageHash, m.hashEncountered = 0, 0
			if indexOfCR >= 0 {
				requestToPci = true
			}
		}
	}
	if !pciResponse && !requestToPci {
		// Apparently we have not found any message yet
		return nil, nil
	}

	packetLength := indexOfCR + 1
	if pciResponse {
		packetLength = indexOfLF + 1
	}

	// Sanity check
	if pciResponse && requestToPci {
		panic("Invalid state... Can not be response and request at the same time")
	}

	read, err := ti.Read(uint32(packetLength))
	if err != nil {
		panic("Invalid state... If we have peeked that before we should be able to read that now")
	}
	rb := utils.NewReadBufferByteBased(read)
	cBusMessage, err := readwriteModel.CBusMessageParse(rb, pciResponse, m.requestContext, m.cbusOptions)
	if err != nil {
		log.Debug().Err(err).Msg("First Parse Failed")
		{ // Try SAL
			rb := utils.NewReadBufferByteBased(read)
			cBusMessage, secondErr := readwriteModel.CBusMessageParse(rb, pciResponse, readwriteModel.NewRequestContext(false), m.cbusOptions)
			if secondErr == nil {
				return cBusMessage, nil
			} else {
				log.Debug().Err(secondErr).Msg("SAL parse failed too")
			}
		}
		{ // Try MMI
			requestContext := readwriteModel.NewRequestContext(false)
			cbusOptions := readwriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
			rb := utils.NewReadBufferByteBased(read)
			cBusMessage, secondErr := readwriteModel.CBusMessageParse(rb, true, requestContext, cbusOptions)
			if secondErr == nil {
				return cBusMessage, nil
			} else {
				log.Debug().Err(secondErr).Msg("CAL parse failed too")
			}
		}

		log.Warn().Err(err).Msg("error parsing")
		// TODO: Possibly clean up ...
		return nil, nil
	}
	return cBusMessage, nil
}
