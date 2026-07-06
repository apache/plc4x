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

package serial

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
)

const (
	sharedSubscriberBufferCapacity = 64 << 10
	sharedReaderSlice              = 100 * time.Millisecond
	sharedReaderChunk              = 4096
)

// sharedPortConfig is the full tuple that must be identical for two
// connections to share one physical port. It is comparable by design.
type sharedPortConfig struct {
	port            serialport.Config
	dtr             bool
	rts             bool
	interframeDelay time.Duration
}

// sharedPortRegistry tracks the shared ports of one Transport. It is not
// global: sharing works within one driver manager; separate managers
// opening the same device remain unsupported (as without sharing).
type sharedPortRegistry struct {
	mu             sync.Mutex
	ports          map[string]*sharedPort
	bufferCapacity int // per-subscriber ring size; overridable in tests
	log            zerolog.Logger
}

func newSharedPortRegistry(log zerolog.Logger) *sharedPortRegistry {
	return &sharedPortRegistry{
		ports:          map[string]*sharedPort{},
		bufferCapacity: sharedSubscriberBufferCapacity,
		log:            log,
	}
}

// acquire joins the existing shared port (config must match exactly) or
// opens the device and starts its reader.
func (r *sharedPortRegistry) acquire(portName string, cfg sharedPortConfig) (serialport.Port, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.ports[portName]; ok {
		if existing.cfg != cfg {
			return nil, errors.Errorf(
				"serial port %s is already shared with a different configuration (existing: %+v, requested: %+v)",
				portName, existing.cfg, cfg)
		}
		if sub, ok := existing.subscribe(r.bufferCapacity); ok {
			return sub, nil
		}
		// The entry closed concurrently; drop it and open fresh below.
		delete(r.ports, portName)
	}
	port, err := serialport.Open(portName, cfg.port)
	if err != nil {
		return nil, errors.Wrap(err, "error opening shared serial port")
	}
	applySharedModemLines(port, cfg, r.log)
	sp := &sharedPort{
		registry: r,
		portName: portName,
		cfg:      cfg,
		port:     port,
		pacer:    newPacer(cfg.interframeDelay),
		log:      r.log,
	}
	r.ports[portName] = sp
	sub, ok := sp.subscribe(r.bufferCapacity)
	if !ok {
		// Cannot happen: the port was created closed==false and nothing
		// else references it yet.
		delete(r.ports, portName)
		_ = port.Close()
		return nil, errors.New("serialport: freshly opened shared port already closed")
	}
	go sp.readLoop()
	return sub, nil
}

// remove deletes the registry entry if it still maps to sp.
func (r *sharedPortRegistry) remove(sp *sharedPort) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ports[sp.portName] == sp {
		delete(r.ports, sp.portName)
	}
}

// applySharedModemLines asserts DTR/RTS once at open, mirroring the
// dedicated path's warn-don't-fail semantics.
func applySharedModemLines(port serialport.Port, cfg sharedPortConfig, log zerolog.Logger) {
	controlPort, ok := port.(serialport.ControlPort)
	if !ok {
		return
	}
	if cfg.dtr {
		if err := controlPort.SetDTR(true); err != nil {
			log.Warn().Err(err).Msg("could not assert DTR on shared port")
		}
	}
	if cfg.rts {
		if err := controlPort.SetRTS(true); err != nil {
			log.Warn().Err(err).Msg("could not assert RTS on shared port")
		}
	}
}

// sharedPort owns one physical port: a broadcast reader, the subscriber
// list, and the serialized paced write path.
//
// Locking: stateMu guards subs/closed; writeMu serializes writers. writeMu may
// acquire stateMu momentarily (via isClosed) but never the reverse — the nesting
// is strictly one-directional, so no cycle is possible. Never hold stateMu while
// taking registry.mu; on last-unsubscribe the physical port is closed BEFORE the
// registry entry is removed, so a concurrent acquire either joins a live port or
// re-opens a fully released device.
type sharedPort struct {
	registry *sharedPortRegistry
	portName string
	cfg      sharedPortConfig
	port     serialport.Port
	pacer    *pacer
	log      zerolog.Logger

	stateMu sync.Mutex
	subs    []*byteRing
	closed  bool

	writeMu sync.Mutex
}

// subscribe adds a subscriber; ok is false when the port already closed.
func (sp *sharedPort) subscribe(capacity int) (*subscription, bool) {
	sp.stateMu.Lock()
	defer sp.stateMu.Unlock()
	if sp.closed {
		return nil, false
	}
	ring := newByteRing(capacity)
	sp.subs = append(sp.subs, ring)
	sub := &subscription{sp: sp, ring: ring}
	sub.readDeadline.Store(time.Time{})
	sub.writeDeadline.Store(time.Time{})
	return sub, true
}

func (sp *sharedPort) isClosed() bool {
	sp.stateMu.Lock()
	defer sp.stateMu.Unlock()
	return sp.closed
}

// unsubscribe detaches ring; the last subscriber closes the port.
func (sp *sharedPort) unsubscribe(ring *byteRing) error {
	ring.close()
	sp.stateMu.Lock()
	for i, s := range sp.subs {
		if s == ring {
			sp.subs = append(sp.subs[:i], sp.subs[i+1:]...)
			break
		}
	}
	last := len(sp.subs) == 0 && !sp.closed
	if last {
		sp.closed = true
	}
	sp.stateMu.Unlock()
	if !last {
		return nil
	}
	err := sp.port.Close() // unblocks the reader; close before unmapping
	sp.registry.remove(sp)
	return err
}

// fail propagates a fatal port error to all subscribers and evicts the
// port so a reconnect re-opens the device.
func (sp *sharedPort) fail(err error) {
	sp.log.Warn().Err(err).Str("port", sp.portName).Msg("shared serial port failed")
	sp.stateMu.Lock()
	alreadyClosed := sp.closed
	sp.closed = true
	subs := sp.subs
	sp.subs = nil
	sp.stateMu.Unlock()
	for _, ring := range subs {
		ring.fail(err)
	}
	if !alreadyClosed {
		_ = sp.port.Close()
	}
	sp.registry.remove(sp)
}

// readLoop broadcasts everything the port delivers into every
// subscriber's ring. Deadline slices keep shutdown prompt.
func (sp *sharedPort) readLoop() {
	buf := make([]byte, sharedReaderChunk)
	warnedDrops := map[*byteRing]uint64{}
	for {
		if sp.isClosed() {
			return
		}
		_ = sp.port.SetReadDeadline(time.Now().Add(sharedReaderSlice))
		n, err := sp.port.Read(buf)
		if n > 0 {
			sp.pacer.noteActivity()
			sp.stateMu.Lock()
			subs := append([]*byteRing(nil), sp.subs...)
			sp.stateMu.Unlock()
			for _, ring := range subs {
				if dropped := ring.append(buf[:n]); dropped > 0 {
					total := ring.droppedTotal()
					if last := warnedDrops[ring]; last == 0 || total >= 2*last {
						warnedDrops[ring] = total
						sp.log.Warn().
							Int("dropped", dropped).
							Uint64("totalDropped", total).
							Str("port", sp.portName).
							Msg("shared serial subscriber buffer overflow, oldest bytes dropped")
					}
				}
			}
		}
		if err != nil && !transports.ErrorIs(err, os.ErrDeadlineExceeded) {
			if sp.isClosed() {
				return // expected: the last unsubscribe closed the port
			}
			sp.fail(err)
			return
		}
	}
}

// write serializes and paces writes from all subscribers.
func (sp *sharedPort) write(ring *byteRing, p []byte, deadline time.Time) (int, error) {
	sp.writeMu.Lock()
	defer sp.writeMu.Unlock()
	// The subscription may close while a Write is waiting on writeMu; recheck under the lock
	// so a closed connection's bytes never reach the shared line.
	if ring.isClosed() {
		return 0, os.ErrClosed
	}
	if sp.isClosed() {
		return 0, os.ErrClosed
	}
	if err := sp.pacer.waitTurn(deadline); err != nil {
		return 0, err
	}
	if err := sp.port.SetWriteDeadline(deadline); err != nil {
		return 0, err
	}
	n, err := sp.port.Write(p)
	sp.pacer.noteActivity()
	return n, err
}

// subscription is one connection's view of a shared port. It implements
// serialport.Port (and deliberately NOT serialport.ControlPort — modem
// line changes would cross-talk between connections).
type subscription struct {
	sp            *sharedPort
	ring          *byteRing
	readDeadline  atomic.Value // time.Time
	writeDeadline atomic.Value // time.Time
	closeOnce     sync.Once
	closeErr      error
}

var _ serialport.Port = (*subscription)(nil)

func (s *subscription) Read(p []byte) (int, error) {
	deadline, _ := s.readDeadline.Load().(time.Time)
	return s.ring.read(p, deadline)
}

func (s *subscription) Write(p []byte) (int, error) {
	if s.ring.isClosed() {
		return 0, os.ErrClosed
	}
	deadline, _ := s.writeDeadline.Load().(time.Time)
	return s.sp.write(s.ring, p, deadline)
}

func (s *subscription) SetReadDeadline(t time.Time) error {
	if s.ring.isClosed() {
		return os.ErrClosed
	}
	s.readDeadline.Store(t)
	return nil
}

func (s *subscription) SetWriteDeadline(t time.Time) error {
	if s.ring.isClosed() {
		return os.ErrClosed
	}
	s.writeDeadline.Store(t)
	return nil
}

func (s *subscription) Close() error {
	s.closeOnce.Do(func() {
		s.closeErr = s.sp.unsubscribe(s.ring)
	})
	return s.closeErr
}
