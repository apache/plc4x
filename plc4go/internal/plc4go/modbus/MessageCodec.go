//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package modbus

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Expectation struct {
	expiration     time.Time
	acceptsMessage spi.AcceptsMessage
	handleMessage  spi.HandleMessage
	handleError    spi.HandleError
}

func (m Expectation) String() string {
	return fmt.Sprintf("Expectation(expires at %v)", m.expiration)
}

type MessageCodec struct {
	expectationCounter            int32
	transportInstance             transports.TransportInstance
	defaultIncomingMessageChannel chan interface{}
	expectations                  []Expectation
	running                       bool
}

func NewMessageCodec(transportInstance transports.TransportInstance, defaultIncomingMessageChannel chan interface{}) *MessageCodec {
	codec := &MessageCodec{
		expectationCounter:            1,
		transportInstance:             transportInstance,
		defaultIncomingMessageChannel: defaultIncomingMessageChannel,
		expectations:                  []Expectation{},
		running:                       true,
	}
	// TODO: should we better move this go func into Connect(). If not a better explanation why we start the worker so early
	// Start a worker that handles processing of responses
	go work(codec)
	return codec
}

func (m *MessageCodec) Connect() error {
	log.Info().Msg("Connecting")
	err := m.transportInstance.Connect()
	if err == nil {
		m.running = true
	}
	return err
}

func (m *MessageCodec) Disconnect() error {
	log.Info().Msg("Disconnecting")
	m.running = false
	return m.transportInstance.Close()
}

func (m *MessageCodec) Send(message interface{}) error {
	log.Trace().Msg("Sending message")
	// Cast the message to the correct type of struct
	adu := model.CastModbusTcpADU(message)
	// Serialize the request
	wb := utils.NewWriteBuffer()
	err := adu.Serialize(*wb)
	if err != nil {
		return errors.Wrap(err, "error serializing request")
	}

	// Send it to the PLC
	err = m.transportInstance.Write(wb.GetBytes())
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	return nil
}

func (m *MessageCodec) Expect(acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	expectation := Expectation{
		expiration:     time.Now().Add(ttl),
		acceptsMessage: acceptsMessage,
		handleMessage:  handleMessage,
		handleError:    handleError,
	}
	m.expectations = append(m.expectations, expectation)
	return nil
}

func (m *MessageCodec) SendRequest(message interface{}, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	log.Trace().Msg("Sending request")
	// Send the actual message
	err := m.Send(message)
	if err != nil {
		return errors.Wrap(err, "Error sending the request")
	}
	return m.Expect(acceptsMessage, handleMessage, handleError, ttl)
}

func (m *MessageCodec) GetDefaultIncomingMessageChannel() chan interface{} {
	return m.defaultIncomingMessageChannel
}

func (m *MessageCodec) receive() (interface{}, error) {
	log.Trace().Msg("receiving")
	// We need at least 6 bytes in order to know how big the packet is in total
	if num, err := m.transportInstance.GetNumReadableBytes(); (err == nil) && (num >= 6) {
		log.Debug().Msgf("we got %d readable bytes", num)
		data, err := m.transportInstance.PeekReadableBytes(6)
		if err != nil {
			log.Warn().Err(err).Msg("error peeking")
			// TODO: Possibly clean up ...
			return nil, nil
		}
		// Get the size of the entire packet
		packetSize := (uint32(data[4]) << 8) + uint32(data[5]) + 6
		if num >= packetSize {
			data, err = m.transportInstance.Read(packetSize)
			if err != nil {
				// TODO: Possibly clean up ...
				return nil, nil
			}
			rb := utils.NewReadBuffer(data)
			adu, err := model.ModbusTcpADUParse(rb, true)
			if err != nil {
				log.Warn().Err(err).Msg("error parsing")
				// TODO: Possibly clean up ...
				return nil, nil
			}
			return adu, nil
		}
	}
	// TODO: maybe we return here a not enough error error
	return nil, nil
}

func work(m *MessageCodec) {
	// Start an endless loop
	for m.running {
		log.Trace().Msg("working")
		if len(m.expectations) <= 0 {
			// Sleep for 10ms
			time.Sleep(time.Millisecond * 10)
			continue
		}
		message, err := m.receive()
		if err != nil {
			log.Error().Err(err).Msg("got an error reading from transport")
			continue
		}
		if message == nil {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		now := time.Now()
		messageHandled := false
		// Go through all expectations
		for index, expectation := range m.expectations {
			// Check if this expectation has expired.
			if now.After(expectation.expiration) {
				log.Debug().Stringer("expectation", expectation).Msg("expired")
				// Remove this expectation from the list.
				m.expectations = append(m.expectations[:index], m.expectations[index+1:]...)
				break
			}

			// Check if the current message matches the expectations
			// If it does, let it handle the message.
			if accepts := expectation.acceptsMessage(message); accepts {
				log.Debug().Stringer("expectation", expectation).Msg("accepts message")
				err = expectation.handleMessage(message)
				if err == nil {
					messageHandled = true
					// Remove this expectation from the list.
					m.expectations = append(m.expectations[:index], m.expectations[index+1:]...)
				} else {
					log.Error().Err(err).Msg("Error handling message")
				}
				break
			}
		}

		// If the message has not been handled and a default handler is provided, call this ...
		if !messageHandled {
			if m.defaultIncomingMessageChannel != nil {
				m.defaultIncomingMessageChannel <- message
			} else {
				log.Warn().Msgf("No handler registered for handling message %s", message)
			}
		}
	}
}

func (m MessageCodec) GetTransportInstance() transports.TransportInstance {
	return m.transportInstance
}
