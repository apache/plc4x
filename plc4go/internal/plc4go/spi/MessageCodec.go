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
package spi

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/plcerrors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type Expectation interface {
	GetExpiration() time.Time
	GetAcceptsMessage() AcceptsMessage
	GetHandleMessage() HandleMessage
	GetHandleError() HandleError
	fmt.Stringer
}

// If this function returns true, the message is forwarded to the message handler
type AcceptsMessage func(message interface{}) bool

// Function for handling the message, returns an error if anything goes wrong
type HandleMessage func(message interface{}) error

// Function for handling the message, returns an error if anything goes wrong
type HandleError func(err error) error

type DefaultExpectation struct {
	Expiration     time.Time
	AcceptsMessage AcceptsMessage
	HandleMessage  HandleMessage
	HandleError    HandleError
}

func (m *DefaultExpectation) GetExpiration() time.Time {
	return m.Expiration
}

func (m *DefaultExpectation) GetAcceptsMessage() AcceptsMessage {
	return m.AcceptsMessage
}

func (m *DefaultExpectation) GetHandleMessage() HandleMessage {
	return m.HandleMessage
}

func (m *DefaultExpectation) GetHandleError() HandleError {
	return m.HandleError
}

func (m *DefaultExpectation) String() string {
	return fmt.Sprintf("Expectation(expires at %v)", m.Expiration)
}

type MessageCodec interface {
	Connect() error
	Disconnect() error

	// Sends a given message
	Send(message interface{}) error
	// Wait for a given timespan for a message to come in, which returns 'true' for 'acceptMessage'
	// and is then forwarded to the 'handleMessage' function
	Expect(acceptsMessage AcceptsMessage, handleMessage HandleMessage, handleError HandleError, ttl time.Duration) error
	// A combination that sends a message first and then waits for a response
	SendRequest(message interface{}, acceptsMessage AcceptsMessage, handleMessage HandleMessage, handleError HandleError, ttl time.Duration) error

	GetDefaultIncomingMessageChannel() chan interface{}
}

// DefaultCodecRequiredInterface adds required methods to MessageCodec that are needed when using DefaultCodec
type DefaultCodecRequiredInterface interface {
	MessageCodec
	TimeoutExpectations(now time.Time)
	HandleMessages(message interface{}) bool
	Receive() (interface{}, error)
}

// DefaultCodec is a default codec implementation which has so sensitive defaults for message handling and a built-in worker
type DefaultCodec struct {
	DefaultCodecRequiredInterface
	TransportInstance             transports.TransportInstance
	DefaultIncomingMessageChannel chan interface{}
	Expectations                  []Expectation
	Running                       bool
	CustomWorkLoop                func(codec *DefaultCodecRequiredInterface)
	CustomMessageHandling         func(codec *DefaultCodecRequiredInterface, message interface{}) bool
}

func NewDefaultCodec(transportInstance transports.TransportInstance) *DefaultCodec {
	return &DefaultCodec{
		TransportInstance:             transportInstance,
		DefaultIncomingMessageChannel: make(chan interface{}),
		Expectations:                  []Expectation{},
		Running:                       false,
	}
}

func (m *DefaultCodec) GetTransportInstance() transports.TransportInstance {
	return m.TransportInstance
}

func (m *DefaultCodec) GetDefaultIncomingMessageChannel() chan interface{} {
	return m.DefaultIncomingMessageChannel
}

func (m *DefaultCodec) Connect() error {
	log.Info().Msg("Connecting")
	err := m.TransportInstance.Connect()
	if err == nil {
		if !m.Running {
			log.Debug().Msg("Message codec currently not running")
			if m.CustomWorkLoop != nil {
				log.Info().Msg("Starting with custom loop")
				go m.CustomWorkLoop(&m.DefaultCodecRequiredInterface)
			} else {
				go m.Work(&m.DefaultCodecRequiredInterface)
			}

		}
		m.Running = true
	}
	return err
}

func (m *DefaultCodec) Disconnect() error {
	log.Info().Msg("Disconnecting")
	m.Running = false
	return m.TransportInstance.Close()
}

func (m *DefaultCodec) Expect(acceptsMessage AcceptsMessage, handleMessage HandleMessage, handleError HandleError, ttl time.Duration) error {
	expectation := &DefaultExpectation{
		Expiration:     time.Now().Add(ttl),
		AcceptsMessage: acceptsMessage,
		HandleMessage:  handleMessage,
		HandleError:    handleError,
	}
	m.Expectations = append(m.Expectations, expectation)
	return nil
}

func (m *DefaultCodec) SendRequest(message interface{}, acceptsMessage AcceptsMessage, handleMessage HandleMessage, handleError HandleError, ttl time.Duration) error {
	log.Trace().Msg("Sending request")
	// Send the actual message
	err := m.Send(message)
	if err != nil {
		return errors.Wrap(err, "Error sending the request")
	}
	return m.Expect(acceptsMessage, handleMessage, handleError, ttl)
}

func (m *DefaultCodec) TimeoutExpectations(now time.Time) {
	for index, expectation := range m.Expectations {
		// Check if this expectation has expired.
		if now.After(expectation.GetExpiration()) {
			// Remove this expectation from the list.
			m.Expectations = append(m.Expectations[:index], m.Expectations[index+1:]...)
			// Call the error handler.
			// TODO: decouple from worker thread
			err := expectation.GetHandleError()(plcerrors.NewTimeoutError(now.Sub(expectation.GetExpiration())))
			if err != nil {
				log.Error().Err(err).Msg("Got an error handling error on expectation")
			}
		}
	}
}

func (m *DefaultCodec) HandleMessages(message interface{}) bool {
	messageHandled := false
	for index, expectation := range m.Expectations {
		// Check if the current message matches the expectations
		// If it does, let it handle the message.
		if accepts := expectation.GetAcceptsMessage()(message); accepts {
			log.Debug().Stringer("expectation", expectation).Msg("accepts message")
			// TODO: decouple from worker thread
			err := expectation.GetHandleMessage()(message)
			if err != nil {
				// Pass the error to the error handler.
				// TODO: decouple from worker thread
				err := expectation.GetHandleError()(err)
				if err != nil {
					log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
				continue
			}
			messageHandled = true
			// If this is the last element of the list remove it differently than if it's before that
			if (index + 1) == len(m.Expectations) {
				m.Expectations = m.Expectations[:index]
			} else if (index + 1) < len(m.Expectations) {
				m.Expectations = append(m.Expectations[:index], m.Expectations[index+1:]...)
			}
		}
	}
	return messageHandled
}

func (m *DefaultCodec) Work(codec *DefaultCodecRequiredInterface) {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msgf("recovered from %v", err)
		}
		log.Info().Msg("Keep running")
		m.Work(codec)
	}()

	workerLog := log.With().Logger()
	if !config.TraceDefaultMessageCodecWorker {
		workerLog = zerolog.Nop()
	}
	// Start an endless loop
mainLoop:
	for m.Running {
		workerLog.Trace().Msg("Working")
		// Check for any expired expectations.
		// (Doing this outside the loop lets us expire expectations even if no input is coming in)
		now := time.Now()

		// Guard against empty expectations
		if len(m.Expectations) <= 0 {
			workerLog.Trace().Msg("no available expectations")
			// Sleep for 10ms
			time.Sleep(time.Millisecond * 10)
			continue mainLoop
		}
		m.TimeoutExpectations(now)

		workerLog.Trace().Msg("Receiving message")
		// Check for incoming messages.
		message, err := m.Receive()
		if err != nil {
			workerLog.Error().Err(err).Msg("got an error reading from transport")
			time.Sleep(time.Millisecond * 10)
			continue mainLoop
		}
		if message == nil {
			workerLog.Trace().Msg("Not enough data yet")
			// Sleep for 10ms before checking again, in order to not
			// consume 100% CPU Power.
			time.Sleep(time.Millisecond * 10)
			continue mainLoop
		}

		if m.CustomMessageHandling != nil {
			workerLog.Trace().Msg("Executing custom handling")
			if m.CustomMessageHandling(codec, message) {
				continue mainLoop
			}
		}

		workerLog.Trace().Msg("Handle message")
		// Go through all expectations
		messageHandled := m.HandleMessages(message)

		// If the message has not been handled and a default handler is provided, call this ...
		if !messageHandled {
			workerLog.Trace().Msg("Message was not handled")
			// TODO: how do we prevent endless blocking if there is no reader on this channel?
			select {
			case m.DefaultIncomingMessageChannel <- message:
			default:
				workerLog.Warn().Msg("Message discarded")
			}
		}
	}
}
