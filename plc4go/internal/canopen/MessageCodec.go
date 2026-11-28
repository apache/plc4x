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

package canopen

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// receivePollInterval is how long one read of the transport waits before the receive loop checks
// whether it was asked to stop. It is not a latency: a frame which arrives is handled at once.
const receivePollInterval = 200 * time.Millisecond

// MessageCodec moves CANopen frames between the bus and the connection.
//
// It does not use the SPI's DefaultCodec, because that codec's receive worker waits 10 ms between
// messages. That suits request/response protocols, but a CAN bus carries thousands of frames a
// second and PDOs arrive unasked; a 100 frames per second ceiling would silently fall behind. This
// codec runs one receive loop which reads frames back to back and hands each to, in order:
//
//  1. the first waiter whose predicate accepts it (an SDO transfer waiting for its answer), else
//  2. the frame handler (subscriptions), else
//  3. the default incoming message channel, if nobody drains that it is dropped.
//
// That is the structure of plc4j's CANOpenConnection, with its expect() listeners and dispatch.
type MessageCodec struct {
	transportInstance transports.TransportInstance
	format            frameFormat

	// handler gets every frame no waiter took. Set before Connect.
	handler func(frame readWriteModel.CANOpenFrame)
	// onFatal is told when the receive loop stops because the transport broke.
	onFatal func(err error)

	waitersMutex sync.Mutex
	waiters      []*waiter

	writeMutex sync.Mutex

	lifecycleMutex sync.Mutex
	cancelReceive  context.CancelFunc
	receiveDone    chan struct{}
	// stopped is closed when the current receive loop ends, which is what fails pending waiters.
	stopped chan struct{}

	defaultIncomingMessageChannel chan spi.Message

	wg sync.WaitGroup

	log zerolog.Logger
}

var (
	_ spi.MessageCodec             = (*MessageCodec)(nil)
	_ spi.TransportInstanceExposer = (*MessageCodec)(nil)
)

func NewMessageCodec(transportInstance transports.TransportInstance, compact bool, _options ...options.WithOption) *MessageCodec {
	var format frameFormat = rawFrameFormat{}
	if compact {
		format = compactFrameFormat{}
	}
	stopped := make(chan struct{})
	close(stopped)
	return &MessageCodec{
		transportInstance:             transportInstance,
		format:                        format,
		stopped:                       stopped,
		defaultIncomingMessageChannel: make(chan spi.Message, 100),
		log:                           options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

func (m *MessageCodec) GetTransportInstance() transports.TransportInstance {
	return m.transportInstance
}

func (m *MessageCodec) Connect(ctx context.Context) error {
	m.lifecycleMutex.Lock()
	defer m.lifecycleMutex.Unlock()
	if m.cancelReceive != nil {
		return errors.New("already connected")
	}
	if !m.transportInstance.IsConnected() {
		if err := m.transportInstance.Connect(ctx); err != nil {
			return errors.Wrap(err, "error connecting the transport")
		}
	}
	receiveCtx, cancel := context.WithCancel(context.Background())
	m.cancelReceive = cancel
	m.receiveDone = make(chan struct{})
	m.stopped = make(chan struct{})
	done, stopped := m.receiveDone, m.stopped
	m.wg.Go(func() {
		defer close(done)
		defer close(stopped)
		m.receive(receiveCtx)
	})
	return nil
}

func (m *MessageCodec) Disconnect() error {
	m.lifecycleMutex.Lock()
	cancel, done := m.cancelReceive, m.receiveDone
	m.cancelReceive, m.receiveDone = nil, nil
	m.lifecycleMutex.Unlock()
	if cancel == nil {
		// DefaultConnection.Close recognises exactly this text.
		return errors.New("already disconnected")
	}
	cancel()
	<-done
	return nil
}

func (m *MessageCodec) IsRunning() bool {
	m.lifecycleMutex.Lock()
	defer m.lifecycleMutex.Unlock()
	return m.cancelReceive != nil
}

// receive is the receive loop.
func (m *MessageCodec) receive(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			m.log.Error().Str("stack", string(debug.Stack())).Interface("err", err).Msg("CANopen receive loop panic-ed")
		}
	}()
	recordSize := m.format.recordSize()
	for ctx.Err() == nil {
		record, err := m.readRecord(ctx, recordSize)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			if !m.transportInstance.IsConnected() || m.transportInstance.ClassifyError(err) == transports.TransportErrorFatal {
				m.log.Warn().Err(err).Msg("CANopen receive loop stopped: the transport failed")
				if m.onFatal != nil {
					go m.onFatal(err)
				}
				return
			}
			continue
		}
		frame, ok, err := m.format.decode(ctx, record)
		if err != nil {
			m.log.Debug().Err(err).Hex("record", record).Msg("Dropping a CAN frame which doesn't decode as CANopen")
			continue
		}
		if !ok {
			continue
		}
		m.dispatch(frame)
	}
}

func (m *MessageCodec) readRecord(ctx context.Context, recordSize uint32) ([]byte, error) {
	readCtx, cancel := context.WithTimeout(ctx, receivePollInterval)
	defer cancel()
	if _, err := m.transportInstance.PeekReadableBytes(readCtx, recordSize); err != nil {
		return nil, err
	}
	return m.transportInstance.Read(readCtx, recordSize)
}

func (m *MessageCodec) dispatch(frame readWriteModel.CANOpenFrame) {
	m.waitersMutex.Lock()
	for i, w := range m.waiters {
		if w.accepts(frame) {
			m.waiters = append(m.waiters[:i:i], m.waiters[i+1:]...)
			m.waitersMutex.Unlock()
			w.frames <- frame
			return
		}
	}
	m.waitersMutex.Unlock()
	if m.handler != nil {
		m.handler(frame)
		return
	}
	select {
	case m.defaultIncomingMessageChannel <- frame:
	default:
	}
}

// waiter is a one-shot registration for the next frame a predicate accepts.
type waiter struct {
	codec   *MessageCodec
	accepts func(readWriteModel.CANOpenFrame) bool
	frames  chan readWriteModel.CANOpenFrame
	stopped chan struct{}
}

// expect registers a waiter. Register before sending the request, or a fast answer can arrive
// before anybody is waiting for it.
func (m *MessageCodec) expect(accepts func(readWriteModel.CANOpenFrame) bool) *waiter {
	m.lifecycleMutex.Lock()
	stopped := m.stopped
	m.lifecycleMutex.Unlock()
	w := &waiter{codec: m, accepts: accepts, frames: make(chan readWriteModel.CANOpenFrame, 1), stopped: stopped}
	m.waitersMutex.Lock()
	m.waiters = append(m.waiters, w)
	m.waitersMutex.Unlock()
	return w
}

// wait returns the accepted frame, or an error once the context ends or the codec stops.
func (w *waiter) wait(ctx context.Context) (readWriteModel.CANOpenFrame, error) {
	select {
	case frame := <-w.frames:
		return frame, nil
	case <-ctx.Done():
		w.cancel()
		return nil, ctx.Err()
	case <-w.stopped:
		w.cancel()
		return nil, errors.New("the CANopen connection is closed")
	}
}

// cancel drops the registration. Calling it after the waiter was served is harmless.
func (w *waiter) cancel() {
	m := w.codec
	m.waitersMutex.Lock()
	defer m.waitersMutex.Unlock()
	for i, registered := range m.waiters {
		if registered == w {
			m.waiters = append(m.waiters[:i:i], m.waiters[i+1:]...)
			return
		}
	}
}

// sendFrame puts one frame on the bus.
func (m *MessageCodec) sendFrame(ctx context.Context, frame readWriteModel.CANOpenFrame) error {
	record, err := m.format.encode(frame)
	if err != nil {
		return errors.Wrap(err, "error encoding the CANopen frame")
	}
	m.writeMutex.Lock()
	defer m.writeMutex.Unlock()
	if err := m.transportInstance.Write(ctx, record); err != nil {
		return errors.Wrap(err, "error sending the CANopen frame")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// spi.MessageCodec, for callers that use the generic interface
////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *MessageCodec) Send(ctx context.Context, _ string, message spi.Message) error {
	frame, ok := message.(readWriteModel.CANOpenFrame)
	if !ok {
		return errors.Errorf("message is not a CANOpenFrame, got %T", message)
	}
	return m.sendFrame(ctx, frame)
}

func (m *MessageCodec) Expect(ctx context.Context, _ string, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) {
	w := m.expect(func(frame readWriteModel.CANOpenFrame) bool { return acceptsMessage(frame) })
	m.wg.Go(func() {
		frame, err := w.wait(ctx)
		if err != nil {
			_ = handleError(err)
			return
		}
		if err := handleMessage(frame); err != nil {
			_ = handleError(err)
		}
	})
}

func (m *MessageCodec) SendRequest(ctx context.Context, interactionInfo string, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	m.Expect(ctx, interactionInfo, acceptsMessage, handleMessage, handleError)
	return m.Send(ctx, interactionInfo, message)
}

func (m *MessageCodec) GetDefaultIncomingMessageChannel() chan spi.Message {
	return m.defaultIncomingMessageChannel
}

func (m *MessageCodec) String() string {
	return "canopen.MessageCodec{" + m.transportInstance.String() + "}"
}
