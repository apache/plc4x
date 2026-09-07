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

package ads

import (
	"context"
	"io"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/protocols/ads/discovery/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// maxAdsDiscoveryFrameSize is the maximum total discovery frame size (header plus payload)
// this codec accepts. Discovery responses are small; anything approaching the UDP datagram
// limit is bogus, and the 32-bit wire length must be bounded before sizing any allocation.
const maxAdsDiscoveryFrameSize = 64 * 1024

type DiscoveryMessageCodec struct {
	_default.DefaultCodec

	passLogToModel bool
	log            zerolog.Logger
}

func NewDiscoveryMessageCodec(transportInstance transports.TransportInstance, _options ...options.WithOption) *DiscoveryMessageCodec {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	codec := &DiscoveryMessageCodec{
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _options...)
	return codec
}

func (m *DiscoveryMessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *DiscoveryMessageCodec) classifyTransportError(err error) transports.TransportErrorKind {
	if err == nil {
		return transports.TransportErrorUnknown
	}
	if transport := m.GetTransportInstance(); transport != nil {
		return transport.ClassifyError(err)
	}
	return transports.TransportErrorUnknown
}

func (m *DiscoveryMessageCodec) isFatalTransportError(err error) bool {
	if err == nil || transports.ErrorIs(err, io.EOF) {
		return false
	}
	return m.classifyTransportError(err) == transports.TransportErrorFatal
}

func (m *DiscoveryMessageCodec) wrapFatalTransportError(err error, msg string) error {
	if err == nil {
		return nil
	}
	// TODO: Any additional context?
	return transports.NewTransportError(transports.TransportErrorFatal, errors.Wrap(err, msg))
}

func (m *DiscoveryMessageCodec) Send(ctx context.Context, interactionInfo string, message spi.Message) error {
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending message")
	// Cast the message to the correct type of struct
	tcpPaket := message.(model.AdsDiscovery)
	// Serialize the request
	bytes, err := tcpPaket.Serialize()
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.GetTransportInstance().Write(ctx, bytes)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *DiscoveryMessageCodec) Receive(ctx context.Context) (spi.Message, error) {
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 6) {
		m.log.Debug().Uint32("num", num).Msg("we got num readable bytes")
		data, err := m.GetTransportInstance().PeekReadableBytes(ctx, 6)
		if err != nil {
			m.log.Warn().Err(err).Msg("error peeking")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error peeking header")
			}
			return nil, nil
		}
		// Get the size of the entire packet little endian plus size of header.
		// Compute in uint64 so a wire length close to 0xFFFFFFFF cannot wrap the
		// packet size to zero (which would make this codec spin forever without
		// consuming a single byte).
		discoveryLength := (uint64(data[5]) << 24) + (uint64(data[4]) << 16) + (uint64(data[3]) << 8) + (uint64(data[2]))
		if discoveryLength == 0 || discoveryLength+6 > maxAdsDiscoveryFrameSize {
			// A zero length can never frame a valid discovery packet and an oversize
			// length would force a huge allocation; both are unrecoverable framing errors.
			return nil, transports.NewTransportError(transports.TransportErrorFatal,
				errors.Errorf("invalid ADS discovery frame length %d (must be > 0 and at most %d)", discoveryLength, maxAdsDiscoveryFrameSize-6))
		}
		packetSize := uint32(discoveryLength) + 6
		if num < packetSize {
			m.log.Debug().Uint32("num", num).Uint32("packetSize", packetSize).Msg("Not enough bytes. Got: num Need: packetSize")
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(ctx, packetSize)
		if err != nil {
			m.log.Warn().Err(err).Msg("error reading packet data")
			if m.isFatalTransportError(err) {
				return nil, m.wrapFatalTransportError(err, "error reading packet data")
			}
			return nil, nil
		}
		ctxForModel := options.GetLoggerContextForModel(ctx, m.log, options.WithPassLoggerToModel(m.passLogToModel))
		tcpPacket, err := model.AdsDiscoveryParse(ctxForModel, data)
		if err != nil {
			m.log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return tcpPacket, nil
	} else if err != nil {
		m.log.Warn().Err(err).Msg("Got error reading")
		if m.isFatalTransportError(err) {
			return nil, m.wrapFatalTransportError(err, "error getting readable bytes")
		}
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}
