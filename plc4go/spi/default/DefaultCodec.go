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
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// DefaultCodecRequirements adds required methods to MessageCodec that are needed when using DefaultCodec
type DefaultCodecRequirements interface {
	GetCodec() spi.MessageCodec
	Send(message spi.Message) error
	Receive() (spi.Message, error)
}

// DefaultCodec is a default codec implementation which has so sensitive defaults for message handling and a built-in worker
type DefaultCodec interface {
	utils.Serializable
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

type CustomMessageHandler func(codec DefaultCodecRequirements, message spi.Message) bool

func WithCustomMessageHandler(customMessageHandler CustomMessageHandler) options.WithOption {
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

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=defaultCodec
type defaultCodec struct {
	DefaultCodecRequirements `ignore:"true"`

	transportInstance transports.TransportInstance

	expectations                  []spi.Expectation
	defaultIncomingMessageChannel chan spi.Message
	customMessageHandling         func(codec DefaultCodecRequirements, message spi.Message) bool

	expectationsChangeMutex sync.RWMutex
	running                 atomic.Bool
	stateChange             sync.Mutex
	activeWorker            sync.WaitGroup

	log zerolog.Logger `ignore:"true"`
}

func buildDefaultCodec(defaultCodecRequirements DefaultCodecRequirements, transportInstance transports.TransportInstance, _options ...options.WithOption) DefaultCodec {
	var customMessageHandler func(codec DefaultCodecRequirements, message spi.Message) bool
	var logger = options.ExtractCustomLogger(_options...)

	for _, option := range _options {
		switch option := option.(type) {
		case withCustomMessageHandler:
			customMessageHandler = option.customMessageHandler
		}
	}

	return &defaultCodec{
		DefaultCodecRequirements:      defaultCodecRequirements,
		transportInstance:             transportInstance,
		defaultIncomingMessageChannel: make(chan spi.Message),
		expectations:                  []spi.Expectation{},
		customMessageHandling:         customMessageHandler,
		log:                           logger,
	}
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (d *DefaultExpectation) GetContext() context.Context {
	return d.Context
}

func (d *DefaultExpectation) GetExpiration() time.Time {
	return d.Expiration
}

func (d *DefaultExpectation) GetAcceptsMessage() spi.AcceptsMessage {
	return d.AcceptsMessage
}

func (d *DefaultExpectation) GetHandleMessage() spi.HandleMessage {
	return d.HandleMessage
}

func (d *DefaultExpectation) GetHandleError() spi.HandleError {
	return d.HandleError
}

func (d *DefaultExpectation) String() string {
	return fmt.Sprintf("Expectation(expires at %v)", d.Expiration)
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
	m.stateChange.Lock()
	defer m.stateChange.Unlock()
	if m.running.Load() {
		return errors.New("already running")
	}
	m.log.Trace().Msg("Connecting")
	if !m.transportInstance.IsConnected() {
		if err := m.transportInstance.ConnectWithContext(ctx); err != nil {
			return err
		}
	} else {
		m.log.Info().Msg("Transport instance already connected")
	}

	m.log.Debug().Msg("Message codec currently not running, starting worker now")
	m.activeWorker.Add(1)
	go m.Work(m.DefaultCodecRequirements)
	m.running.Store(true)
	return nil
}

func (m *defaultCodec) Disconnect() error {
	m.stateChange.Lock()
	defer m.stateChange.Unlock()
	if !m.running.Load() {
		return errors.New("already disconnected")
	}
	m.log.Trace().Msg("Disconnecting")
	m.running.Store(false)
	if m.transportInstance != nil {
		if err := m.transportInstance.Close(); err != nil {
			return errors.Wrap(err, "error closing transport instance")
		}
	}
	m.log.Trace().Msg("Waiting for worker to shutdown")
	m.activeWorker.Wait()
	return nil
}

func (m *defaultCodec) IsRunning() bool {
	return m.running.Load()
}

func (m *defaultCodec) Expect(ctx context.Context, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	m.expectationsChangeMutex.Lock()
	defer m.expectationsChangeMutex.Unlock()
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
	m.log.Trace().Msg("Sending request")
	// Send the actual message
	err := m.Send(message)
	if err != nil {
		return errors.Wrap(err, "Error sending the request")
	}
	return m.Expect(ctx, acceptsMessage, handleMessage, handleError, ttl)
}

func (m *defaultCodec) TimeoutExpectations(now time.Time) {
	m.expectationsChangeMutex.Lock() // TODO: Note: would be nice if this is a read mutex which can be upgraded
	defer m.expectationsChangeMutex.Unlock()
	for i := 0; i < len(m.expectations); i++ {
		expectation := m.expectations[i]
		// Check if this expectation has expired.
		if now.After(expectation.GetExpiration()) {
			// Remove this expectation from the list.
			m.expectations = append(m.expectations[:i], m.expectations[i+1:]...)
			i--
			// Call the error handler.
			go func(expectation spi.Expectation) {
				if err := expectation.GetHandleError()(utils.NewTimeoutError(now.Sub(expectation.GetExpiration()))); err != nil {
					m.log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
			}(expectation)
			continue
		}
		if err := expectation.GetContext().Err(); err != nil {
			// Remove this expectation from the list.
			m.expectations = append(m.expectations[:i], m.expectations[i+1:]...)
			i--
			go func(expectation spi.Expectation) {
				if err := expectation.GetHandleError()(err); err != nil {
					m.log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
			}(expectation)
			continue
		}
	}
}

func (m *defaultCodec) HandleMessages(message spi.Message) bool {
	m.expectationsChangeMutex.Lock() // TODO: Note: would be nice if this is a read mutex which can be upgraded
	defer m.expectationsChangeMutex.Unlock()
	messageHandled := false
	m.log.Trace().Msgf("Current number of expectations: %d", len(m.expectations))
	for index, expectation := range m.expectations {
		m.log.Trace().Msgf("Checking expectation %s", expectation)
		// Check if the current message matches the expectations
		// If it does, let it handle the message.
		if accepts := expectation.GetAcceptsMessage()(message); accepts {
			m.log.Debug().Stringer("expectation", expectation).Msg("accepts message")
			// TODO: decouple from worker thread
			if err := expectation.GetHandleMessage()(message); err != nil {
				m.log.Debug().Stringer("expectation", expectation).Err(err).Msg("errored handling the message")
				// Pass the error to the error handler.
				// TODO: decouple from worker thread
				if err := expectation.GetHandleError()(err); err != nil {
					m.log.Error().Err(err).Msg("Got an error handling error on expectation")
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
		} else {
			m.log.Trace().Stringer("expectation", expectation).Msg("doesn't accept message")
		}
	}
	return messageHandled
}

func (m *defaultCodec) Work(codec DefaultCodecRequirements) {
	defer m.activeWorker.Done()
	workerLog := m.log.With().Logger()
	if !config.TraceDefaultMessageCodecWorker {
		workerLog = zerolog.Nop()
	}

	defer func(workerLog zerolog.Logger) {
		if err := recover(); err != nil {
			// TODO: If this is an error, cast it to an error and log it with "Err(err)"
			m.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
		}
		if m.running.Load() {
			workerLog.Warn().Msg("Keep running")
			m.activeWorker.Add(1)
			go m.Work(codec)
		} else {
			workerLog.Info().Msg("Worker terminated")
		}
	}(workerLog)

	// Start an endless loop
mainLoop:
	for m.running.Load() {
		workerLog.Trace().Msg("Working")
		// Check for any expired expectations.
		// (Doing this outside the loop lets us expire expectations even if no input is coming in)
		now := time.Now()

		// Guard against empty expectations
		m.expectationsChangeMutex.RLock()
		numberOfExpectations := len(m.expectations)
		m.expectationsChangeMutex.RUnlock()
		if numberOfExpectations <= 0 && m.customMessageHandling == nil {
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
			}
			workerLog.Trace().Msg("Custom handling didn't handle the message")
		}

		workerLog.Trace().Msg("Handle message")
		// Go through all expectations
		messageHandled := m.HandleMessages(message)

		// If the message has not been handled and a default handler is provided, call this ...
		if !messageHandled {
			workerLog.Trace().Msg("Message was not handled")
			m.passToDefaultIncomingMessageChannel(workerLog, message)
		}
	}
}

func (m *defaultCodec) passToDefaultIncomingMessageChannel(workerLog zerolog.Logger, message spi.Message) {
	timeout := time.NewTimer(time.Millisecond * 40)
	defer utils.CleanupTimer(timeout)
	select {
	case m.defaultIncomingMessageChannel <- message:
	case <-timeout.C:
		timeout.Stop()
		workerLog.Warn().Msgf("Message discarded\n%s", message)
	}
}
