/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package _default

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/options"

	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/plcerrors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DefaultCodecRequirements adds required methods to MessageCodec that are needed when using DefaultCodec
type DefaultCodecRequirements interface {
	GetCodec() spi.MessageCodec
	Send(message interface{}) error
	Receive() (interface{}, error)
}

// DefaultCodec is a default codec implementation which has so sensitive defaults for message handling and a built-in worker
type DefaultCodec interface {
	spi.MessageCodec
	spi.TransportInstanceExposer
}

// NewDefaultCodec is the factory for a DefaultCodec
func NewDefaultCodec(requirements DefaultCodecRequirements, transportInstance transports.TransportInstance, options ...options.WithOption) DefaultCodec {
	return buildDefaultCodec(requirements, transportInstance, options...)
}

type DefaultExpectation struct {
	Expiration     time.Time
	AcceptsMessage spi.AcceptsMessage
	HandleMessage  spi.HandleMessage
	HandleError    spi.HandleError
}

func WithCustomMessageHandler(customMessageHandler func(codec *DefaultCodecRequirements, message interface{}) bool) options.WithOption {
	return withCustomMessageHandler{customMessageHandler: customMessageHandler}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withCustomMessageHandler struct {
	options.Option
	customMessageHandler func(codec *DefaultCodecRequirements, message interface{}) bool
}

type defaultCodec struct {
	DefaultCodecRequirements
	transportInstance             transports.TransportInstance
	defaultIncomingMessageChannel chan interface{}
	expectations                  []spi.Expectation
	running                       bool
	customMessageHandling         func(codec *DefaultCodecRequirements, message interface{}) bool
}

func buildDefaultCodec(defaultCodecRequirements DefaultCodecRequirements, transportInstance transports.TransportInstance, options ...options.WithOption) DefaultCodec {
	var customMessageHandler func(codec *DefaultCodecRequirements, message interface{}) bool

	for _, option := range options {
		switch option.(type) {
		case withCustomMessageHandler:
			customMessageHandler = option.(withCustomMessageHandler).customMessageHandler
			log.Debug()
		}
	}

	return &defaultCodec{
		defaultCodecRequirements,
		transportInstance,
		make(chan interface{}),
		[]spi.Expectation{},
		false,
		customMessageHandler,
	}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (m *DefaultExpectation) GetExpiration() time.Time {
	return m.Expiration
}

func (m *DefaultExpectation) GetAcceptsMessage() spi.AcceptsMessage {
	return m.AcceptsMessage
}

func (m *DefaultExpectation) GetHandleMessage() spi.HandleMessage {
	return m.HandleMessage
}

func (m *DefaultExpectation) GetHandleError() spi.HandleError {
	return m.HandleError
}

func (m *DefaultExpectation) String() string {
	return fmt.Sprintf("Expectation(expires at %v)", m.Expiration)
}

func (m *defaultCodec) GetTransportInstance() transports.TransportInstance {
	return m.transportInstance
}

func (m *defaultCodec) GetDefaultIncomingMessageChannel() chan interface{} {
	return m.defaultIncomingMessageChannel
}

func (m *defaultCodec) Connect() error {
	log.Trace().Msg("Connecting")
	err := m.transportInstance.Connect()
	if err == nil {
		if !m.running {
			log.Debug().Msg("Message codec currently not running, starting worker now")
			go m.Work(&m.DefaultCodecRequirements)
		}
		m.running = true
	}
	return err
}

func (m *defaultCodec) Disconnect() error {
	log.Info().Msg("Disconnecting")
	m.running = false
	return m.transportInstance.Close()
}

func (m *defaultCodec) IsRunning() bool {
	return m.running
}

func (m *defaultCodec) Expect(acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	expectation := &DefaultExpectation{
		Expiration:     time.Now().Add(ttl),
		AcceptsMessage: acceptsMessage,
		HandleMessage:  handleMessage,
		HandleError:    handleError,
	}
	m.expectations = append(m.expectations, expectation)
	return nil
}

func (m *defaultCodec) SendRequest(message interface{}, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	log.Trace().Msg("Sending request")
	// Send the actual message
	err := m.Send(message)
	if err != nil {
		return errors.Wrap(err, "Error sending the request")
	}
	return m.Expect(acceptsMessage, handleMessage, handleError, ttl)
}

func (m *defaultCodec) TimeoutExpectations(now time.Time) {
	for index, expectation := range m.expectations {
		// Check if this expectation has expired.
		if now.After(expectation.GetExpiration()) {
			// Remove this expectation from the list.
			m.expectations = append(m.expectations[:index], m.expectations[index+1:]...)
			// Call the error handler.
			// TODO: decouple from worker thread
			err := expectation.GetHandleError()(plcerrors.NewTimeoutError(now.Sub(expectation.GetExpiration())))
			if err != nil {
				log.Error().Err(err).Msg("Got an error handling error on expectation")
			}
		}
	}
}

func (m *defaultCodec) HandleMessages(message interface{}) bool {
	messageHandled := false
	for index, expectation := range m.expectations {
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
			if (index + 1) == len(m.expectations) {
				m.expectations = m.expectations[:index]
			} else if (index + 1) < len(m.expectations) {
				m.expectations = append(m.expectations[:index], m.expectations[index+1:]...)
			}
		}
	}
	return messageHandled
}

func (m *defaultCodec) Work(codec *DefaultCodecRequirements) {
	defer func() {
		if err := recover(); err != nil {
			// TODO: If this is an error, cast it to an error and log it with "Err(err)"
			log.Error().Msgf("recovered from: %#v at %s", err, string(debug.Stack()))
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
	for m.running {
		workerLog.Trace().Msg("Working")
		// Check for any expired expectations.
		// (Doing this outside the loop lets us expire expectations even if no input is coming in)
		now := time.Now()

		// Guard against empty expectations
		if len(m.expectations) <= 0 && m.customMessageHandling == nil {
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

		if m.customMessageHandling != nil {
			workerLog.Trace().Msg("Executing custom handling")
			if m.customMessageHandling(codec, message) {
				continue mainLoop
			}
		}

		workerLog.Trace().Msg("Handle message")
		// Go through all expectations
		messageHandled := m.HandleMessages(message)

		// If the message has not been handled and a default handler is provided, call this ...
		if !messageHandled {
			workerLog.Trace().Msg("Message was not handled")
			timeout := time.NewTimer(time.Millisecond * 40)
			select {
			case m.defaultIncomingMessageChannel <- message:
				if !timeout.Stop() {
					<-timeout.C
				}
			case <-timeout.C:
				timeout.Stop()
				workerLog.Warn().Msgf("Message discarded %s", message)
			}
		}
	}
}
