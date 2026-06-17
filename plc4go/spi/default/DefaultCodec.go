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
	stdErrors "errors"
	"runtime/debug"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// DefaultCodecRequirements adds required methods to MessageCodec that are needed when using DefaultCodec
type DefaultCodecRequirements interface {
	GetCodec() spi.MessageCodec
	Send(ctx context.Context, interactionInfo string, message spi.Message) error
	Receive(ctx context.Context) (spi.Message, error)
}

// DefaultCodec is a default codec implementation which has so sensitive defaults for message handling and a built-in worker
type DefaultCodec interface {
	utils.Serializable
	spi.MessageCodec
	spi.TransportInstanceExposer
	spi.TransportErrorHandlerSetter
}

// NewDefaultCodec is the factory for a DefaultCodec
func NewDefaultCodec(requirements DefaultCodecRequirements, transportInstance transports.TransportInstance, options ...options.WithOption) DefaultCodec {
	return buildDefaultCodec(requirements, transportInstance, options...)
}

type CustomMessageHandler func(ctx context.Context, codec DefaultCodecRequirements, message spi.Message) bool

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
	customMessageHandler CustomMessageHandler
}

//go:generate go tool plc4xGenerator -type=defaultCodec
type defaultCodec struct {
	DefaultCodecRequirements `ignore:"true"`

	transportInstance transports.TransportInstance

	expectations                  []spi.Expectation
	defaultIncomingMessageChannel chan spi.Message
	customMessageHandling         CustomMessageHandler

	expectationsChangeMutex sync.RWMutex
	running                 atomic.Bool
	stateChange             sync.Mutex
	activeWorker            sync.WaitGroup
	notifyExpireWorker      chan struct{} `ignore:"true"`
	notifyReceiveWorker     chan struct{} `ignore:"true"`

	receiveTimeout                 time.Duration
	traceDefaultMessageCodecWorker bool

	ctx       context.Context
	ctxCancel context.CancelFunc

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger

	transportErrorHandler transports.TransportErrorHandler
}

func buildDefaultCodec(defaultCodecRequirements DefaultCodecRequirements, transportInstance transports.TransportInstance, _options ...options.WithOption) DefaultCodec {
	var customMessageHandler CustomMessageHandler

	for _, option := range _options {
		switch option := option.(type) {
		case withCustomMessageHandler:
			customMessageHandler = option.customMessageHandler
		}
	}

	receiveTimeout, timeoutDefined := options.ExtractReceiveTimeout(_options...)
	if !timeoutDefined {
		receiveTimeout = 60 * time.Second
	}
	traceDefaultMessageCodecWorker, _ := options.ExtractTraceDefaultMessageCodecWorker(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	d := &defaultCodec{
		DefaultCodecRequirements:       defaultCodecRequirements,
		transportInstance:              transportInstance,
		defaultIncomingMessageChannel:  make(chan spi.Message, 100),
		expectations:                   []spi.Expectation{},
		customMessageHandling:          customMessageHandler,
		notifyExpireWorker:             make(chan struct{}),
		notifyReceiveWorker:            make(chan struct{}),
		receiveTimeout:                 receiveTimeout,
		traceDefaultMessageCodecWorker: traceDefaultMessageCodecWorker || config.TraceDefaultMessageCodecWorker,
		log:                            customLogger,
	}
	d.ctx, d.ctxCancel = context.WithCancel(context.Background())
	return d
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (m *defaultCodec) GetTransportInstance() transports.TransportInstance {
	return m.transportInstance
}

func (m *defaultCodec) SetTransportErrorHandler(handler transports.TransportErrorHandler) {
	m.transportErrorHandler = handler
}

func (m *defaultCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	return m.defaultIncomingMessageChannel
}

func (m *defaultCodec) Connect(ctx context.Context) error {
	m.stateChange.Lock()
	defer m.stateChange.Unlock()
	if m.running.Load() {
		return errors.New("already running")
	}
	m.log.Trace().Msg("connecting")
	if !m.transportInstance.IsConnected() {
		if err := m.transportInstance.Connect(ctx); err != nil {
			return err
		}
	} else {
		m.log.Debug().Msg("Transport instance already connected")
	}

	m.log.Debug().Msg("Message codec currently not running, starting worker now")
	m.startWorkers()
	m.running.Store(true)
	m.log.Trace().Msg("connected")
	return nil
}

func (m *defaultCodec) Disconnect() error {
	m.stateChange.Lock()
	defer m.stateChange.Unlock()
	if !m.running.Load() {
		return errors.New("already disconnected")
	}
	m.log.Trace().Msg("Disconnecting")
	m.ctxCancel()
	m.running.Store(false)
	close(m.notifyExpireWorker)
	close(m.notifyReceiveWorker)
	m.log.Trace().Msg("Waiting for worker to shutdown")
	m.activeWorker.Wait()
	m.log.Trace().Msg("worker shut down")
	if m.transportInstance != nil {
		m.log.Trace().Msg("closing transport instance")
		if err := m.transportInstance.Close(); err != nil {
			return errors.Wrap(err, "error closing transport instance")
		}
	}
	for _, expectation := range m.expectations {
		m.wg.Go(func() {
			err := errors.New("disconnected")
			expectation.Cancel(err)
			_ = expectation.GetHandleError()(err)
		})
	}
	m.wg.Wait()
	m.log.Trace().Msg("disconnected")
	return nil
}

func (m *defaultCodec) IsRunning() bool {
	return m.running.Load()
}

func (m *defaultCodec) Expect(ctx context.Context, interactionInfo string, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) {
	m.expectationsChangeMutex.Lock()
	defer m.expectationsChangeMutex.Unlock()
	ttl := m.receiveTimeout
	if deadline, ok := ctx.Deadline(); ok {
		ttl = time.Until(deadline)
	}
	expectation := newDefaultExpectation(ctx, interactionInfo, ttl, acceptsMessage, handleMessage, handleError)
	m.expectations = append(m.expectations, expectation)
	m.log.Debug().Str("interactionInfo", interactionInfo).Stringer("expectation", expectation).Msg("Added expectation")
	select {
	case m.notifyExpireWorker <- struct{}{}:
	default:
	}
	select {
	case m.notifyReceiveWorker <- struct{}{}:
	default:
	}
}

func (m *defaultCodec) SendRequest(ctx context.Context, interactionInfo string, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "Not sending message as context is aborted")
	}
	m.Expect(ctx, interactionInfo, acceptsMessage, handleMessage, handleError) // We register the expectation first to avoid getting a response between sending and adding the expect
	m.log.Trace().Str("interactionInfo", interactionInfo).Msg("Sending request")
	return m.Send(ctx, interactionInfo, message)
}

func (m *defaultCodec) TimeoutExpectations(now time.Time) time.Duration {
	m.expectationsChangeMutex.Lock() // TODO: Note: would be nice if this is a read mutex which can be upgraded
	defer m.expectationsChangeMutex.Unlock()
	m.expectations = slices.DeleteFunc(m.expectations, func(expectation spi.Expectation) bool {
		// Check if this expectation has expired.
		if now.After(expectation.GetExpiration()) {
			// Remove this expectation from the list.
			m.log.Debug().Interface("expectation", expectation).Msg("timeout expectation")
			// Call the error handler.
			m.wg.Go(func() {
				timeoutErr := utils.NewTimeoutError(expectation.GetExpiration().Sub(expectation.GetCreationTime()))
				expectation.Cancel(timeoutErr)
				if err := expectation.GetHandleError()(timeoutErr); err != nil {
					m.log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
			})
			return true
		}
		if err := expectation.GetContext().Err(); err != nil {
			m.log.Debug().Err(err).Interface("expectation", expectation).Msg("expectation canceled")
			// Remove this expectation from the list.
			m.wg.Go(func() {
				if err := expectation.GetHandleError()(err); err != nil {
					m.log.Error().Err(err).Msg("Got an error handling error on expectation")
				}
			})
			return true
		}
		return false
	})
	nextExpire := 30 * time.Second
	for _, expectation := range m.expectations {
		expiresIn := time.Until(expectation.GetExpiration())
		if expiresIn < nextExpire {
			nextExpire = expiresIn
		}
	}
	return nextExpire
}

func (m *defaultCodec) HandleMessages(message spi.Message) bool {
	m.expectationsChangeMutex.Lock() // TODO: Note: would be nice if this is a read mutex which can be upgraded
	defer m.expectationsChangeMutex.Unlock()
	messageHandled := false
	m.log.Trace().Int("nExpectations", len(m.expectations)).Msg("Current number of expectations")
	m.expectations = slices.DeleteFunc(m.expectations, func(expectation spi.Expectation) bool {
		expectationLog := m.log.With().Stringer("expectation", expectation).Logger()
		expectationLog.Trace().Msg("Checking expectation")
		// Check if the current message matches the expectations
		// If it does, let it handle the message.
		if accepts := expectation.GetAcceptsMessage()(message); accepts {
			expectationLog.Trace().Interface("handleMessage", message).Msg("accepts message")
			// TODO: decouple from worker thread
			if err := expectation.GetHandleMessage()(message); err != nil {
				expectationLog.Debug().Err(err).Msg("errored handling the message")
				// Pass the error to the error handler.
				m.wg.Go(func() {
					if err := expectation.GetHandleError()(err); err != nil {
						m.log.Error().Err(err).Msg("Got an error handling error on expectation")
					}
				})
				return false
			}
			m.log.Trace().Msg("message handled")
			messageHandled = true
			return true
		} else {
			expectationLog.Trace().Interface("handleMessage", message).Msg("doesn't accept message")
			return false
		}
	})
	m.log.Trace().Bool("messageHandled", messageHandled).Msg("finished message handling")
	return messageHandled
}

func (m *defaultCodec) startWorkers() {
	m.log.Trace().Msg("starting workers")
	m.startExpire()
	m.startReceive()
}

func (m *defaultCodec) startExpire() {
	m.log.Trace().Msg("starting expire worker")
	m.activeWorker.Go(m.ExpireWork)
}

func (m *defaultCodec) startReceive() {
	m.log.Trace().Msg("starting receive worker")
	m.activeWorker.Go(m.ReceiveWork)
}

func (m *defaultCodec) ExpireWork() {
	workerLog := m.log.With().Logger()
	if !m.traceDefaultMessageCodecWorker {
		workerLog = zerolog.Nop()
	}
	workerLog = workerLog.With().Str("workerId", uuid.NewString()).Logger()
	workerLog.Trace().Msg("Starting expire work")
	defer workerLog.Trace().Msg("expire work ended")

	defer func() {
		if err := recover(); err != nil {
			m.log.Error().
				Str("stack", string(debug.Stack())).
				Interface("err", err).
				Msg("panic-ed")
		}
		if m.running.Load() {
			workerLog.Warn().Msg("Keep running")
			m.startExpire()
		} else {
			workerLog.Info().Msg("expire worker terminated")
		}
	}()

	// Start an endless loop
mainLoop:
	for m.running.Load() {
		workerLog.Trace().Msg("expire mainloop cycle")
		now := time.Now()

		// Guard against empty expectations
		m.expectationsChangeMutex.RLock()
		numberOfExpectations := len(m.expectations)
		m.expectationsChangeMutex.RUnlock()
		if numberOfExpectations <= 0 && m.customMessageHandling == nil {
			workerLog.Trace().Msg("no available expectations")
			timer := time.NewTimer(30 * time.Second)
			select {
			case <-m.notifyExpireWorker:
				workerLog.Trace().Msg("waking up because of notification")
			case <-m.ctx.Done():
				workerLog.Trace().Msg("context done, exiting expire work")
				return
			case <-timer.C:
				workerLog.Trace().Msg("waking up for next expire")
			}
			continue mainLoop
		}
		nextExpire := m.TimeoutExpectations(now)
		m.expectationsChangeMutex.RLock()
		numberOfExpectations = len(m.expectations)
		m.expectationsChangeMutex.RUnlock()
		workerLog.Debug().Dur("nextExpire", nextExpire).Int("numberOfExpectations", numberOfExpectations).Msg("waiting for next expire")
		timer := time.NewTimer(nextExpire)
		select {
		case <-m.notifyExpireWorker:
			workerLog.Trace().Msg("waking up because of notification")
		case <-m.ctx.Done():
			workerLog.Trace().Msg("context done, exiting expire work")
			return
		case <-timer.C:
			workerLog.Trace().Msg("waking up for next expire")
		}
	}
}

func (m *defaultCodec) ReceiveWork() {
	workerLog := m.log.With().Logger()
	if !m.traceDefaultMessageCodecWorker {
		workerLog = zerolog.Nop()
	}
	workerLog = workerLog.With().Str("workerId", uuid.NewString()).Logger()
	workerLog.Trace().Msg("Starting work")
	defer workerLog.Trace().Msg("work ended")

	defer func() {
		if err := recover(); err != nil {
			m.log.Error().
				Str("stack", string(debug.Stack())).
				Interface("err", err).
				Msg("panic-ed")
		}
		if m.running.Load() {
			workerLog.Warn().Msg("Keep running")
			m.startReceive()
		} else {
			workerLog.Info().Msg("receive worker terminated")
		}
	}()

	lastLoopTime := time.Now()
	// Start an endless loop
mainLoop:
	for m.running.Load() {
		const cycleTime = 10 * time.Millisecond
		if processingTime := time.Since(lastLoopTime); processingTime < cycleTime {
			// Ensure that we leave at least 10ms between loops to not burn cycles
			sleepTime := cycleTime - processingTime
			workerLog.Trace().Dur("sleepTime", sleepTime).Msg("sleeping")
			time.Sleep(sleepTime)
		} else {
			workerLog.Debug().Dur("processingTime", processingTime).Dur("cycleTime", cycleTime).Msg("no need to sleep")
		}
		workerLog.Trace().Msg("receive mainloop cycle")
		lastLoopTime = time.Now()

		// Guard against empty expectations
		m.expectationsChangeMutex.RLock()
		numberOfExpectations := len(m.expectations)
		m.expectationsChangeMutex.RUnlock()
		if numberOfExpectations <= 0 && m.customMessageHandling == nil {
			workerLog.Trace().Msg("no available expectations")
			timer := time.NewTimer(30 * time.Second)
			select {
			case <-m.notifyReceiveWorker:
				workerLog.Trace().Msg("waking up because of notification")
			case <-m.ctx.Done():
				workerLog.Trace().Msg("context done, exiting receive work")
				return
			case <-timer.C:
				workerLog.Trace().Msg("waking up for next receive")
				continue mainLoop
			}
		}

		workerLog.Trace().Msg("Receiving message")
		// Check for incoming messages.
		var message spi.Message
		var err error
		{
			syncer := make(chan struct{})
			m.wg.Go(func() {
				defer close(syncer)
				if !m.running.Load() {
					err = errors.New("not running")
					return
				}
				message, err = m.Receive(m.ctx)
			})
			timeoutTimer := time.NewTimer(m.receiveTimeout)
			select {
			case <-syncer:
			// nothing
			case <-m.ctx.Done():
				workerLog.Trace().Msg("context done, exiting receive work")
				return
			case <-timeoutTimer.C:
				workerLog.Error().Dur("receiveTimeout", m.receiveTimeout).Msg("receive timeout")
				continue mainLoop
			}
		}
		if err != nil {
			workerLog.Error().Err(err).Msg("got an error reading from transport")
			if m.handleTransportError(workerLog, err) {
				continue mainLoop
			}
			workerLog.Debug().Msg("transport error requested worker shutdown")
			return
		}
		if message == nil {
			workerLog.Trace().Msg("Not enough data yet")
			continue mainLoop
		}
		workerLog.Trace().Interface("message", message).Msg("got message")

		if m.customMessageHandling != nil {
			workerLog.Trace().Msg("Executing custom handling")
			start := time.Now()
			handled := m.customMessageHandling(m.ctx, m.DefaultCodecRequirements, message)
			workerLog.Trace().TimeDiff("elapsedTime", time.Now(), start).Bool("handled", handled).Msg("custom handling took elapsedTime")
			if handled {
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
	select {
	case m.defaultIncomingMessageChannel <- message:
	default:
		workerLog.Warn().Interface("message", message).Msg("Message discarded")
	}
}

func (m *defaultCodec) handleTransportError(workerLog zerolog.Logger, err error) bool {
	if err == nil {
		return true
	}
	if stdErrors.Is(err, context.Canceled) {
		workerLog.Debug().Msg("receive aborted due to context cancellation")
		return false
	}

	kind := transports.TransportErrorUnknown
	if stdErrors.Is(err, context.DeadlineExceeded) {
		kind = transports.TransportErrorRetryable
	} else if m.transportInstance != nil {
		kind = m.transportInstance.ClassifyError(err)
	}
	if kind == transports.TransportErrorUnknown {
		workerLog.Warn().Err(err).Msg("transport error classified as unknown; treating as fatal")
		kind = transports.TransportErrorFatal
	}

	switch kind {
	case transports.TransportErrorTransient:
		workerLog.Debug().Err(err).Msg("transient transport error; keeping worker alive")
		m.emitTransportError(kind, transports.NewTransportError(kind, err))
		return true
	case transports.TransportErrorRetryable:
		workerLog.Warn().Err(err).Msg("retryable transport error; resetting transport instance")
		if m.transportInstance != nil {
			defer func() {
				if recoverErr := recover(); recoverErr != nil {
					workerLog.Error().Interface("panic", recoverErr).Msg("panic while resetting transport instance")
				}
			}()
			m.transportInstance.Reset()
		}
		m.emitTransportError(kind, transports.NewTransportError(kind, err))
		return true
	case transports.TransportErrorFatal:
		workerLog.Error().Err(err).Msg("fatal transport error; shutting down codec")
		wrappedErr := transports.NewTransportError(kind, err)
		m.failAllExpectations(wrappedErr)
		if m.transportInstance != nil {
			if closeErr := m.transportInstance.Close(); closeErr != nil {
				workerLog.Warn().Err(closeErr).Msg("error closing transport after fatal condition")
			}
		}
		if m.ctxCancel != nil {
			m.ctxCancel()
		}
		m.running.Store(false)
		m.emitTransportError(kind, wrappedErr)
		return false
	default:
		workerLog.Error().Err(err).Msg("unexpected transport error classification; treating as fatal")
		wrappedErr := transports.NewTransportError(transports.TransportErrorFatal, err)
		m.failAllExpectations(wrappedErr)
		if m.transportInstance != nil {
			if closeErr := m.transportInstance.Close(); closeErr != nil {
				workerLog.Warn().Err(closeErr).Msg("error closing transport after unexpected classification")
			}
		}
		if m.ctxCancel != nil {
			m.ctxCancel()
		}
		m.running.Store(false)
		m.emitTransportError(transports.TransportErrorFatal, wrappedErr)
		return false
	}
}

func (m *defaultCodec) emitTransportError(kind transports.TransportErrorKind, err error) {
	if m.transportErrorHandler != nil {
		m.transportErrorHandler(kind, err)
	}
}

func (m *defaultCodec) failAllExpectations(err error) {
	m.expectationsChangeMutex.Lock()
	expectations := slices.Clone(m.expectations)
	m.expectations = nil
	m.expectationsChangeMutex.Unlock()

	for _, expectation := range expectations {
		expectation := expectation
		expectation.Cancel(err)
		if handleErr := expectation.GetHandleError(); handleErr != nil {
			m.wg.Go(func() {
				if handlerErr := handleErr(err); handlerErr != nil {
					m.log.Error().Err(handlerErr).Msg("error returned by expectation error handler")
				}
			})
		}
	}
}
