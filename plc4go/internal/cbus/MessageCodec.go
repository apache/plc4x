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
	"context"
	"sync/atomic"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"hash/crc32"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=MessageCodec
type MessageCodec struct {
	_default.DefaultCodec

	requestContext readWriteModel.RequestContext
	cbusOptions    readWriteModel.CBusOptions

	monitoredMMIs   chan readWriteModel.CALReply
	monitoredSALs   chan readWriteModel.MonitoredSAL
	lastPackageHash atomic.Uint32
	hashEncountered atomic.Uint64

	currentlyReportedServerErrors atomic.Uint64

	log zerolog.Logger `ignore:"true"`
}

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	codec := &MessageCodec{
		requestContext: readWriteModel.NewRequestContext(false),
		cbusOptions:    readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false),
		monitoredMMIs:  make(chan readWriteModel.CALReply, 100),
		monitoredSALs:  make(chan readWriteModel.MonitoredSAL, 100),
		log:            options.ExtractCustomLogger(_options...),
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, append(_options, _default.WithCustomMessageHandler(extractMMIAndSAL(codec.log)))...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	cbusMessage, ok := message.(readWriteModel.CBusMessage)
	if !ok {
		return errors.Errorf("Invalid message type %T", message)
	}

	// Set the right request context
	m.requestContext = CreateRequestContext(cbusMessage)
	m.log.Debug().Msgf("Created request context\n%s", m.requestContext)

	// Serialize the request
	theBytes, err := cbusMessage.Serialize()
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
	confirmation := false
	// Fill the buffer
	{
		if err := ti.FillBuffer(func(pos uint, currentByte byte, reader *bufio.Reader) bool {
			m.log.Trace().Uint8("byte", currentByte).Msg("current byte")
			switch currentByte {
			case
				readWriteModel.ResponseTermination_CR,
				readWriteModel.ResponseTermination_LF:
				return false
			case byte(readWriteModel.ConfirmationType_CONFIRMATION_SUCCESSFUL):
				confirmation = true
				// In case we have directly more data in the buffer after a confirmation
				_, err := reader.Peek(int(pos + 1))
				return err == nil
			case
				byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_TO_MANY_RE_TRANSMISSIONS),
				byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_CORRUPTION),
				byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_SYNC_LOSS),
				byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_TOO_LONG),
				byte(readWriteModel.ConfirmationType_CHECKSUM_FAILURE):
				confirmation = true
				return false
			default:
				return true
			}
		}); err != nil {
			m.log.Debug().Err(err).Msg("Error filling buffer")
		}
	}
	m.log.Trace().Msg("Buffer filled")

	// Check how many readable bytes we have
	var readableBytes uint32
	{
		numBytesAvailableInBuffer, err := ti.GetNumBytesAvailableInBuffer()
		if err != nil {
			m.log.Warn().Err(err).Msg("Got error reading")
			return nil, nil
		}
		if numBytesAvailableInBuffer == 0 {
			m.log.Trace().Msg("Nothing to read")
			return nil, nil
		}
		readableBytes = numBytesAvailableInBuffer
	}
	m.log.Trace().Msgf("%d bytes available in buffer", readableBytes)

	// Check for an isolated error
	if bytes, err := ti.PeekReadableBytes(1); err == nil && (bytes[0] == byte(readWriteModel.ConfirmationType_CHECKSUM_FAILURE)) {
		_, _ = ti.Read(1)
		// Report one Error at a time
		return readWriteModel.CBusMessageParse(context.TODO(), bytes, true, m.requestContext, m.cbusOptions)
	}

	peekedBytes, err := ti.PeekReadableBytes(readableBytes)
	pciResponse, requestToPci := false, false
	indexOfCR := -1
	indexOfLF := -1
	indexOfConfirmation := -1
lookingForTheEnd:
	for i, peekedByte := range peekedBytes {
		switch peekedByte {
		case
			byte(readWriteModel.ConfirmationType_CONFIRMATION_SUCCESSFUL),
			byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_TO_MANY_RE_TRANSMISSIONS),
			byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_CORRUPTION),
			byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_SYNC_LOSS),
			byte(readWriteModel.ConfirmationType_NOT_TRANSMITTED_TOO_LONG),
			byte(readWriteModel.ConfirmationType_CHECKSUM_FAILURE):
			if indexOfConfirmation < 0 {
				indexOfConfirmation = i
			}
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
	m.log.Trace().Msgf("indexOfCR %d,indexOfLF %d,indexOfConfirmation %d", indexOfCR, indexOfLF, indexOfConfirmation)
	if indexOfCR < 0 && indexOfLF >= 0 {
		// This means that the package is garbage as a lf is always prefixed with a cr
		m.log.Debug().Err(err).Msg("Error reading")
		garbage, err := ti.Read(readableBytes)
		m.log.Warn().Bytes("garbage", garbage).Msg("Garbage bytes")
		return nil, err
	}
	if indexOfCR+1 == indexOfLF {
		m.log.Trace().Msg("pci response for sure")
		// This means a <cr> is directly followed by a <lf> which means that we know for sure this is a response
		pciResponse = true
	} else if indexOfCR >= 0 && int(readableBytes) >= indexOfCR+2 && peekedBytes[+indexOfCR+1] != '\n' {
		m.log.Trace().Msg("pci request for sure")
		// We got a request to pci for sure because the cr is followed by something else than \n
		requestToPci = true
	}
	const numberOfCyclesToWait = 15
	const estimatedElapsedTime = numberOfCyclesToWait * 10
	if !pciResponse && !requestToPci && indexOfLF < 0 {
		// To be sure we might receive that package later we hash the bytes and check if we might receive one
		hash := crc32.NewIEEE()
		_, _ = hash.Write(peekedBytes)
		newPackageHash := hash.Sum32()
		if newPackageHash == m.lastPackageHash.Load() {
			m.hashEncountered.Add(1)
		}
		m.log.Trace().Msgf("new hash %x, last hash %x, seen %d times", newPackageHash, m.lastPackageHash.Load(), m.hashEncountered.Load())
		m.lastPackageHash.Store(newPackageHash)
		if m.hashEncountered.Load() < numberOfCyclesToWait {
			m.log.Trace().Msg("Waiting for more data")
			return nil, nil
		} else {
			m.log.Trace().Msgf("stopping after ~%dms (%d cycles)", estimatedElapsedTime, numberOfCyclesToWait)
			// after numberOfCyclesToWait*10 ms we give up finding a lf
			m.lastPackageHash.Store(0)
			m.hashEncountered.Store(0)
			if indexOfCR >= 0 {
				m.log.Trace().Msg("setting requestToPci")
				requestToPci = true
			}
		}
	}
	if !pciResponse && !requestToPci && !confirmation {
		// Apparently we have not found any message yet
		m.log.Trace().Msg("no message found yet")
		return nil, nil
	}

	// Build length
	packetLength := indexOfCR + 1
	if pciResponse {
		packetLength = indexOfLF + 1
	}
	if !pciResponse && !requestToPci {
		packetLength = indexOfConfirmation + 1
	}

	// Sanity check
	if pciResponse && requestToPci {
		return nil, errors.New("Invalid state... Can not be response and request at the same time")
	}

	// We need to ensure that there is no ! till the first /r
	{
		peekedBytes, err := ti.PeekReadableBytes(readableBytes)
		if err != nil {
			return nil, err
		}
		// We check in the current stream for reported errors
		foundErrors := uint64(0)
		for _, peekedByte := range peekedBytes {
			if peekedByte == '!' {
				foundErrors++
			}
			if peekedByte == '\r' {
				// We only look for errors within
			}
		}
		// Now we report the errors one by one so for every request we get a proper rejection
		if foundErrors > m.currentlyReportedServerErrors.Load() {
			m.log.Debug().Msgf("We found %d errors in the current message. We have %d reported already", foundErrors, m.currentlyReportedServerErrors.Load())
			m.currentlyReportedServerErrors.Add(1)
			return readWriteModel.CBusMessageParse(context.TODO(), []byte{'!'}, true, m.requestContext, m.cbusOptions)
		}
		if foundErrors > 0 {
			m.log.Debug().Msgf("We should have reported all errors by now (%d in total which we reported %d), so we resetting the count", foundErrors, m.currentlyReportedServerErrors.Load())
			m.currentlyReportedServerErrors.Store(0)
		}
		m.log.Trace().Msgf("currentlyReportedServerErrors %d should be 0", m.currentlyReportedServerErrors.Load())
	}

	var rawInput []byte
	{
		m.log.Trace().Msgf("Read packet length %d", packetLength)
		read, err := ti.Read(uint32(packetLength))
		if err != nil {
			return nil, errors.Wrap(err, "Invalid state... If we have peeked that before we should be able to read that now")
		}
		rawInput = read
	}
	var sanitizedInput []byte
	// We remove every error marker we find
	{
		for _, b := range rawInput {
			if b != '!' {
				sanitizedInput = append(sanitizedInput, b)
			}
		}
	}
	m.log.Debug().Msgf("Parsing %q", sanitizedInput)
	cBusMessage, err := readWriteModel.CBusMessageParse(context.TODO(), sanitizedInput, pciResponse, m.requestContext, m.cbusOptions)
	if err != nil {
		m.log.Debug().Err(err).Msg("First Parse Failed")
		{ // Try SAL
			requestContext := readWriteModel.NewRequestContext(false)
			cBusMessage, secondErr := readWriteModel.CBusMessageParse(context.TODO(), sanitizedInput, pciResponse, requestContext, m.cbusOptions)
			if secondErr == nil {
				m.log.Trace().Msgf("Parsed message as SAL:\n%s", cBusMessage)
				return cBusMessage, nil
			} else {
				m.log.Debug().Err(secondErr).Msg("SAL parse failed too")
			}
		}
		{ // Try MMI
			requestContext := readWriteModel.NewRequestContext(false)
			cbusOptions := readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
			cBusMessage, secondErr := readWriteModel.CBusMessageParse(context.TODO(), sanitizedInput, true, requestContext, cbusOptions)
			if secondErr == nil {
				m.log.Trace().Msgf("Parsed message as MMI:\n%s", cBusMessage)
				return cBusMessage, nil
			} else {
				m.log.Debug().Err(secondErr).Msg("CAL parse failed too")
			}
		}

		m.log.Warn().Err(err).Msg("error parsing")
		return nil, nil
	}

	m.log.Trace().Msgf("Parsed message:\n%s", cBusMessage)
	return cBusMessage, nil
}

func extractMMIAndSAL(log zerolog.Logger) _default.CustomMessageHandler {
	return func(codec _default.DefaultCodecRequirements, message spi.Message) bool {
		log.Trace().Msgf("Custom handling message:\n%s", message)
		switch message := message.(type) {
		case readWriteModel.CBusMessageToClientExactly:
			switch reply := message.GetReply().(type) {
			case readWriteModel.ReplyOrConfirmationReplyExactly:
				switch reply := reply.GetReply().(type) {
				case readWriteModel.ReplyEncodedReplyExactly:
					switch encodedReply := reply.GetEncodedReply().(type) {
					case readWriteModel.MonitoredSALReplyExactly:
						log.Trace().Msg("Feed to monitored SALs")
						codec.(*MessageCodec).monitoredSALs <- encodedReply.GetMonitoredSAL()
					case readWriteModel.EncodedReplyCALReplyExactly:
						calData := encodedReply.GetCalReply().GetCalData()
						switch calData.(type) {
						case readWriteModel.CALDataStatusExactly, readWriteModel.CALDataStatusExtendedExactly:
							log.Trace().Msg("Feed to monitored MMIs")
							codec.(*MessageCodec).monitoredMMIs <- encodedReply.GetCalReply()
						default:
							log.Trace().Msgf("Not a CALDataStatusExactly or CALDataStatusExtendedExactly. Actual type %T", calData)
						}
					default:
						log.Trace().Msgf("Not a MonitoredSALReply or EncodedReplyCALReply. Actual type %T", encodedReply)
					}
				default:
					log.Trace().Msgf("Not a ReplyEncodedReply. Actual type %T", reply)
				}
			default:
				log.Trace().Msgf("Not a ReplyOrConfirmationReply. Actual type %T", reply)
			}
		default:
			log.Trace().Msgf("Not a CBusMessageToClient. Actual type %T", message)
		}
		// We never handle mmi or sal here as we might want to read them in a read-request too
		return false
	}
}
