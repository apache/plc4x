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

package _default

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/apache/plc4x/plc4go/spi/options"

	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/plcerrors"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DefaultCodecRequirements adds required methods to MessageCodec that are needed when using DefaultCodec
type DefaultCodecRequirements interface {
	GetCodec() spi.MessageCodec
	Send(message spi.Message) error
	Receive() (spi.Message, error)
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
	Context        context.Context
	Expiration     time.Time
	AcceptsMessage spi.AcceptsMessage
	HandleMessage  spi.HandleMessage
	HandleError    spi.HandleError
}

func WithCustomMessageHandler(customMessageHandler func(codec DefaultCodecRequirements, message spi.Message) bool) options.WithOption {
	return withCustomMessageHandler{customMessageHandler: customMessageHandler}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withCustomMessageHandler struct {
	options.Option
	customMessageHandler func(codec DefaultCodecRequirements, message spi.Message) bool
}

type defaultCodec struct {
	DefaultCodecRequirements
	transportInstance             transports.TransportInstance
	defaultIncomingMessageChannel chan spi.Message
	expectations                  []spi.Expectation
	running                       bool
	customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool
}

func buildDefaultCodec(defaultCodecRequirements DefaultCodecRequirements, transportInstance transports.TransportInstance, options ...options.WithOption) DefaultCodec {
	var customMessageHandler func(codec DefaultCodecRequirements, message spi.Message) bool

	for _, option := range options {
		switch option.(type) {
		case withCustomMessageHandler:
			customMessageHandler = option.(withCustomMessageHandler).customMessageHandler
		}
	}

	return &defaultCodec{
		defaultCodecRequirements,
		transportInstance,
		make(chan spi.Message),
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

func (m *DefaultExpectation) GetContext() context.Context {
	return m.Context
}

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

func (m *defaultCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	return m.defaultIncomingMessageChannel
}

func (m *defaultCodec) Connect() error {
	return m.ConnectWithContext(context.Background())
}

func (m *defaultCodec) ConnectWithContext(ctx context.Context) error {
	log.Trace().Msg("Connecting")
	if !m.transportInstance.IsConnected() {
		if err := m.transportInstance.ConnectWithContext(ctx); err != nil {
			return err
		}
	} else {
		log.Info().Msg("Transport instance already connected")
	}

	if !m.running {
		log.Debug().Msg("Message codec currently not running, starting worker now")
		go m.Work(m.DefaultCodecRequirements)
	}
	m.running = true
	return nil
}

func (m *defaultCodec) Disconnect() error {
	log.Trace().Msg("Disconnecting")
	m.running = false
	return m.transportInstance.Close()
}

func (m *defaultCodec) IsRunning() bool {
	return m.running
}

func (m *defaultCodec) Expect(ctx context.Context, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	expectation := &DefaultExpectation{
		Context:        ctx,
		Expiration:     time.Now().Add(ttl),
		AcceptsMessage: acceptsMessage,
		HandleMessage:  handleMessage,
		HandleError:    handleError,
	}
	m.expectations = append(m.expectations, expectation)
	return nil
}

func (m *defaultCodec) SendRequest(ctx context.Context, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "Not sending message as context is aborted")
	}
	log.Trace().Msg("Sending request")
	// Send the actual message
	err := m.Send(message)
	if err != nil {
		return errors.Wrap(err, "Error sending the request")
	}
	return m.Expect(ctx, acceptsMessage, handleMessage, handleError, ttl)
}

func (m *defaultCodec) TimeoutExpectations(now time.Time) {
	for i := 0; i < len(m.expectations); i++ {
		expectation := m.expectations[i]
		// Check if this expectation has expired.
		if now.After(expectation.GetExpiration()) {
			// Remove this expectation from the list.
			m.expectations = append(m.expectations[:i], m.expectations[i+1:]...)
			i--
			// Call the error handler.
			go func(expectation spi.Expectation) {
				if err := expectation.GetHandleError()(plcerrors.NewTimeoutError(now.Sub(expectation.GetExpiration()))); err != nil {
					log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
			}(expectation)
		}
		if err := expectation.GetContext().Err(); err != nil {
			// Remove this expectation from the list.
			m.expectations = append(m.expectations[:i], m.expectations[i+1:]...)
			i--
			go func(expectation spi.Expectation) {
				if err := expectation.GetHandleError()(err); err != nil {
					log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
			}(expectation)
		}
	}
}

func (m *defaultCodec) HandleMessages(message spi.Message) bool {
	messageHandled := false
	log.Trace().Msgf("Current number of expectations: %d", len(m.expectations))
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

func (m *defaultCodec) Work(codec DefaultCodecRequirements) {
	workerLog := log.With().Logger()
	if !config.TraceDefaultMessageCodecWorker {
		workerLog = zerolog.Nop()
	}

	defer func(workerLog zerolog.Logger) {
		if err := recover(); err != nil {
			// TODO: If this is an error, cast it to an error and log it with "Err(err)"
			log.Error().Msgf("recovered from: %#v at %s", err, string(debug.Stack()))
		}
		if m.running {
			workerLog.Warn().Msg("Keep running")
			m.Work(codec)
		} else {
			workerLog.Info().Msg("Worker terminated")
		}
	}(workerLog)

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
				workerLog.Trace().Msg("Custom handling handled the message")
				continue mainLoop
			} else {
				workerLog.Trace().Msg("Custom handling didn't handle the message")
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
				workerLog.Warn().Msgf("Message discarded\n%s", message)
			}
		}
	}
}
