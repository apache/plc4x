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

package bacnetip

import (
	"context"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/local"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"time"
)

// ApplicationLayerMessageCodec is a wrapper for MessageCodec which takes care of segmentation, retries etc.
type ApplicationLayerMessageCodec struct {
	bipSimpleApplication *BIPSimpleApplication
	messageCode          *MessageCodec
	deviceInfoCache      DeviceInfoCache

	localAddress  *net.UDPAddr
	remoteAddress *net.UDPAddr
}

func NewApplicationLayerMessageCodec(udpTransport *udp.Transport, transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr, remoteAddress *net.UDPAddr) (*ApplicationLayerMessageCodec, error) {
	// Have the transport create a new transport-instance.
	transportInstance, err := udpTransport.CreateTransportInstanceForLocalAddress(transportUrl, options, localAddress)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transport instance")
	}
	_ = transportInstance
	a := &ApplicationLayerMessageCodec{
		localAddress:  localAddress,
		remoteAddress: remoteAddress,
	}
	address, err := NewAddress(localAddress)
	if err != nil {
		return nil, errors.Wrap(err, "error creating address")
	}
	application, err := NewBIPSimpleApplication(&local.LocalDeviceObject{}, *address, &a.deviceInfoCache, nil)
	if err != nil {
		return nil, err
	}
	a.bipSimpleApplication = application
	a.messageCode = NewMessageCodec(transportInstance)
	return a, nil
}

func (m *ApplicationLayerMessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *ApplicationLayerMessageCodec) Connect() error {
	return m.messageCode.Connect()
}

func (m *ApplicationLayerMessageCodec) ConnectWithContext(ctx context.Context) error {
	return m.messageCode.ConnectWithContext(ctx)
}

func (m *ApplicationLayerMessageCodec) Disconnect() error {
	if err := m.bipSimpleApplication.Close(); err != nil {
		log.Error().Err(err).Msg("error closing application")
	}
	return m.messageCode.Disconnect()
}

func (m *ApplicationLayerMessageCodec) IsRunning() bool {
	return m.messageCode.IsRunning()
}

func (m *ApplicationLayerMessageCodec) Send(message spi.Message) error {
	address, err2 := NewAddress(m.remoteAddress)
	if err2 != nil {
		panic(err2)
	}
	iocb, err := NewIOCB(NewPDU(message, WithPDUDestination(address)), m.remoteAddress)
	if err != nil {
		return errors.Wrap(err, "error creating IOCB")
	}
	go func() {
		go m.bipSimpleApplication.RequestIO(iocb)
		iocb.Wait()
		if iocb.ioError != nil {
			// TODO: handle error
			println(iocb.ioError)
		} else if iocb.ioResponse != nil {
			// TODO: response?
			println(iocb.ioResponse)
		} else {
			// TODO: what now?
		}
	}()
	return nil
}

func (m *ApplicationLayerMessageCodec) Expect(ctx context.Context, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	// TODO: implement me
	return nil
}

func (m *ApplicationLayerMessageCodec) SendRequest(ctx context.Context, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {

	// TODO: implement me
	m.Send(message)

	return nil
}

func (m *ApplicationLayerMessageCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	return m.messageCode.GetDefaultIncomingMessageChannel()
}

type MessageCodec struct {
	_default.DefaultCodec
}

func NewMessageCodec(transportInstance transports.TransportInstance) *MessageCodec {
	codec := &MessageCodec{}
	codec.DefaultCodec = _default.NewDefaultCodec(codec, transportInstance, _default.WithCustomMessageHandler(codec.handleCustomMessage))
	return codec
}

func (m *MessageCodec) GetCodec() spi.MessageCodec {
	return m
}

func (m *MessageCodec) Send(message spi.Message) error {
	log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	bvlcPacket := message.(model.BVLC)
	// Serialize the request
	theBytes, err := bvlcPacket.Serialize()
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
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.GetTransportInstance().GetNumBytesAvailableInBuffer(); (err == nil) && (num >= 4) {
		log.Debug().Msgf("we got %d readable bytes", num)
		data, err := m.GetTransportInstance().PeekReadableBytes(4)
		if err != nil {
			log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		packetSize := uint32((uint16(data[2]) << 8) + uint16(data[3]))
		if num < packetSize {
			log.Debug().Msgf("Not enough bytes. Got: %d Need: %d\n", num, packetSize)
			return nil, nil
		}
		data, err = m.GetTransportInstance().Read(packetSize)
		if err != nil {
			log.Debug().Err(err).Msg("Error reading")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		bvlcPacket, err := model.BVLCParse(data)
		if err != nil {
			log.Warn().Err(err).Msg("error parsing")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		return bvlcPacket, nil
	} else if err != nil {
		log.Warn().Err(err).Msg("Got error reading")
		return nil, nil
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}

func (m *MessageCodec) handleCustomMessage(_ _default.DefaultCodecRequirements, message spi.Message) bool {
	// For now, we just put them in the incoming channel
	m.GetDefaultIncomingMessageChannel() <- message
	return true
}
