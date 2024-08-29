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
	"hash/crc32"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
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

	stateChange sync.Mutex

	passLogToModel bool           `ignore:"true"`
	log            zerolog.Logger `ignore:"true"`
}

func NewMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *MessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &MessageCodec{
		requestContext: readWriteModel.NewRequestContext(false),
		cbusOptions:    readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false),
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, append(_options, _default.WithCustomMessageHandler(extractMMIAndSAL(codec.log)))...)
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Connect() error {
	return m.ConnectWithContext(context.Background())
}

func (m *MessageCodec) ConnectWithContext(ctx context.Context) error {
	m.stateChange.Lock()
	defer m.stateChange.Unlock()
	if m.IsRunning() {
		return errors.New("already running")
	}
	m.log.Trace().Msg("building channels")
	m.monitoredMMIs = make(chan readWriteModel.CALReply, 100)
	m.monitoredSALs = make(chan readWriteModel.MonitoredSAL, 100)
	return m.DefaultCodec.ConnectWithContext(ctx)
}

func (m *MessageCodec) Disconnect() error {
	m.stateChange.Lock()
	defer m.stateChange.Unlock()
	if !m.IsRunning() {
		return errors.New("already disconnected")
	}
	err := m.DefaultCodec.Disconnect()
	m.log.Trace().Msg("closing channels")
	close(m.monitoredMMIs)
	close(m.monitoredSALs)
	return err
}

func (m *MessageCodec) Send(message spi.Message) error {
	m.log.Trace().Stringer("message", message).Msg("Sending message")
	// Cast the message to the correct type of struct
	cbusMessage, ok := message.(readWriteModel.CBusMessage)
	if !ok {
		return errors.Errorf("Invalid message type %T", message)
	}

	// Set the right request context
	m.requestContext = CreateRequestContext(cbusMessage)
	m.log.Debug().Stringer("requestContext", m.requestContext).Msg("Created request context")

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
	if !ti.IsConnected() {
		return nil, errors.New("Transport instance not connected")
	}
	confirmation := false
	// Fill the buffer
	{
		if err := ti.FillBuffer(func(pos uint, currentByte byte, reader transports.ExtendedReader) bool {
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
	m.log.Trace().Uint32("readableBytes", readableBytes).Msg("readableBytes bytes available in buffer")

	// Check for an isolated error
	if bytes, err := ti.PeekReadableBytes(1); err == nil && (bytes[0] == byte(readWriteModel.ConfirmationType_CHECKSUM_FAILURE)) {
		_, _ = ti.Read(1)
		// Report one Error at a time
		ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
		return readWriteModel.CBusMessageParse(ctxForModel, bytes, true, m.requestContext, m.cbusOptions)
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
	m.log.Trace().
		Int("indexOfCR", indexOfCR).
		Int("indexOfLF", indexOfLF).
		Int("indexOfConfirmation", indexOfConfirmation).
		Msg("working with indexes")
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
	const estimatedElapsedTime = numberOfCyclesToWait * 10 * time.Millisecond
	if !pciResponse && !requestToPci && indexOfLF < 0 {
		// To be sure we might receive that package later we hash the bytes and check if we might receive one
		hash := crc32.NewIEEE()
		_, _ = hash.Write(peekedBytes)
		newPackageHash := hash.Sum32()
		lastPackageHash := m.lastPackageHash.Load()
		if newPackageHash == lastPackageHash {
			m.hashEncountered.Add(1)
		}
		hasEncountered := m.hashEncountered.Load()
		m.log.Trace().
			Uint32("newPackageHash", newPackageHash).
			Uint32("lastPackageHash", lastPackageHash).
			Uint64("hashEncountered", hasEncountered).
			Msg("new hash newPackageHash, last hash lastPackageHash, seen hashEncountered times")
		m.lastPackageHash.Store(newPackageHash)
		if hasEncountered < numberOfCyclesToWait {
			m.log.Trace().Msg("Waiting for more data")
			return nil, nil
		} else {
			m.log.Trace().
				Dur("estimatedElapsedTime", estimatedElapsedTime).
				Int("numberOfCyclesToWait", numberOfCyclesToWait).
				Msg("stopping after estimatedElapsedTime ms (numberOfCyclesToWait cycles)")
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
		currentlyReportedServerErrors := m.currentlyReportedServerErrors.Load()
		if foundErrors > currentlyReportedServerErrors {
			m.log.Debug().
				Uint64("foundErrors", foundErrors).
				Uint64("currentlyReportedServerErrors", currentlyReportedServerErrors).
				Msg("We found foundErrors errors in the current message. We have currentlyReportedServerErrors reported already")
			m.currentlyReportedServerErrors.Add(1)
			ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
			return readWriteModel.CBusMessageParse(ctxForModel, []byte{'!'}, true, m.requestContext, m.cbusOptions)
		}
		if foundErrors > 0 {
			m.log.Debug().
				Uint64("foundErrors", foundErrors).
				Uint64("currentlyReportedServerErrors", currentlyReportedServerErrors).
				Msg("We should have reported all errors by now (foundErrors in total which we reported currentlyReportedServerErrors), so we resetting the count")
			m.currentlyReportedServerErrors.Store(0)
		}
		m.log.Trace().
			Uint64("currentlyReportedServerErrors", currentlyReportedServerErrors).
			Msg("currentlyReportedServerErrors should be 0")
	}

	var rawInput []byte
	{
		m.log.Trace().Int("packetLength", packetLength).Msg("Read packet length")
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
	m.log.Debug().Bytes("sanitizedInput", sanitizedInput).Msg("Parsing")
	ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
	start := time.Now()
	cBusMessage, err := readWriteModel.CBusMessageParse(ctxForModel, sanitizedInput, pciResponse, m.requestContext, m.cbusOptions)
	m.log.Trace().TimeDiff("elapsedTime", time.Now(), start).Msg("Parsing took elapsedTime")
	if err != nil {
		m.log.Debug().Err(err).Msg("First Parse Failed")
		{ // Try SAL
			m.log.Trace().Msg("try SAL")
			requestContext := readWriteModel.NewRequestContext(false)
			cBusMessage, secondErr := readWriteModel.CBusMessageParse(ctxForModel, sanitizedInput, pciResponse, requestContext, m.cbusOptions)
			if secondErr == nil {
				m.log.Trace().Msg("Parsed message as SAL")
				return cBusMessage, nil
			} else {
				m.log.Debug().Err(secondErr).Msg("SAL parse failed too")
			}
		}
		{ // Try MMI
			m.log.Trace().Msg("try MMI")
			requestContext := readWriteModel.NewRequestContext(false)
			cbusOptions := readWriteModel.NewCBusOptions(false, false, false, false, false, false, false, false, false)
			cBusMessage, secondErr := readWriteModel.CBusMessageParse(ctxForModel, sanitizedInput, true, requestContext, cbusOptions)
			if secondErr == nil {
				m.log.Trace().Msg("Parsed message as MMI")
				return cBusMessage, nil
			} else {
				m.log.Debug().Err(secondErr).Msg("CAL parse failed too")
			}
		}

		m.log.Warn().Err(err).Msg("error parsing")
		return nil, nil
	}

	return cBusMessage, nil
}

func extractMMIAndSAL(log zerolog.Logger) _default.CustomMessageHandler {
	return func(codec _default.DefaultCodecRequirements, message spi.Message) bool {
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
							log.Trace().
								Type("actualType", calData).
								Msg("Not a CALDataStatusExactly or CALDataStatusExtendedExactly")
						}
					default:
						log.Trace().
							Type("actualType", encodedReply).
							Msg("Not a MonitoredSALReply or EncodedReplyCALReply")
					}
				default:
					log.Trace().
						Type("actualType", reply).
						Msg("Not a ReplyEncodedReply")
				}
			default:
				log.Trace().
					Type("actualType", reply).
					Msg("Not a ReplyOrConfirmationReply")
			}
		default:
			log.Trace().
				Type("actualType", message).
				Msg("Not a CBusMessageToClient")
		}
		// We never handle mmi or sal here as we might want to read them in a read-request too
		return false
	}
}
