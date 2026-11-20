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

package cbus

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/cache"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/tcp"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// peerClosableTransportInstance is a test transport instance whose peer can close the connection: from then on
// every read ends with EOF while the instance still reports itself connected, which is what the TCP transport
// does once the C-Bus interface closed the socket.
type peerClosableTransportInstance struct {
	*test.TransportInstance
	peerClosed atomic.Bool
}

func (p *peerClosableTransportInstance) closeByPeer() {
	p.peerClosed.Store(true)
}

func (p *peerClosableTransportInstance) FillBuffer(ctx context.Context, until func(pos uint, currentByte byte, reader transports.ExtendedReader) (keepGoing bool)) error {
	if p.peerClosed.Load() {
		return fmt.Errorf("Error while peeking: %w", io.EOF)
	}
	return p.TransportInstance.FillBuffer(ctx, until)
}

func (p *peerClosableTransportInstance) PeekReadableBytes(ctx context.Context, numBytes uint32) ([]byte, error) {
	if p.peerClosed.Load() {
		return nil, io.EOF
	}
	return p.TransportInstance.PeekReadableBytes(ctx, numBytes)
}

func (p *peerClosableTransportInstance) Read(ctx context.Context, numBytes uint32) ([]byte, error) {
	if p.peerClosed.Load() {
		return nil, io.EOF
	}
	return p.TransportInstance.Read(ctx, numBytes)
}

// peerClosableTransport is the test transport, handing out peerClosableTransportInstances that answer the
// connection setup the way a C-Bus interface does.
type peerClosableTransport struct {
	*test.Transport

	lock      sync.Mutex
	instances []*peerClosableTransportInstance
}

func (p *peerClosableTransport) CreateTransportInstance(transportUrl url.URL, _options map[string][]string, withOptions ...options.WithOption) (transports.TransportInstance, error) {
	ti, err := p.Transport.CreateTransportInstance(transportUrl, _options, withOptions...)
	if err != nil {
		return nil, err
	}
	instance := &peerClosableTransportInstance{TransportInstance: ti.(*test.TransportInstance)}
	answerConnectionSetup(instance.TransportInstance)
	p.lock.Lock()
	defer p.lock.Unlock()
	p.instances = append(p.instances, instance)
	return instance, nil
}

func (p *peerClosableTransport) getInstances() []*peerClosableTransportInstance {
	p.lock.Lock()
	defer p.lock.Unlock()
	return append([]*peerClosableTransportInstance(nil), p.instances...)
}

// answerConnectionSetup answers the requests of the connection setup (reset, application filters, interface
// options) with the echo and confirmation each of them waits for.
func answerConnectionSetup(ti *test.TransportInstance) {
	replies := [][]string{
		{"~~~\r"},
		{"@A32100FF\r", "322100AD\r\n"},
		{"@A32200FF\r", "322200AC\r\n"},
		{"@A342000A\r", "3242008C\r\n"},
		{"@A3410079\r", "3241008D\r\n"},
		{"@A3300079\r", "3230009E\r\n"},
	}
	var lock sync.Mutex
	next := 0
	ti.SetWriteInterceptor(func(transportInstance *test.TransportInstance, _ []byte) {
		lock.Lock()
		defer lock.Unlock()
		if next >= len(replies) {
			return
		}
		for _, reply := range replies[next] {
			transportInstance.FillReadBuffer([]byte(reply))
		}
		next++
	})
}

// A C-Bus interface that closes the connection leaves the TCP transport "connected", answering every read with
// EOF. The codec used to log and drop that error and kept running, so the connection stayed connected, nothing
// invalidated it, and a connection cache handed out the dead connection forever: no reconnect, and no live
// data (SAL/MMI) any more. The EOF has to stop the codec as a fatal transport error, so the connection is
// reported dead and the cache's next lease gets a new connection, whether the dead one sat idle in the cache or
// was leased out when the interface closed it.
func TestConnection_peerClosedConnectionIsReplacedByTheCache(t *testing.T) {
	for _, leased := range []bool{false, true} {
		name := "idle connection"
		if leased {
			name = "leased connection"
		}
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				_options := testutils.EnrichOptionsWithOptionsForTesting(t)
				transport := &peerClosableTransport{Transport: test.NewTransport(_options...)}
				driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				driverManager.RegisterDriver(NewDriver(_options...))
				driverManager.(spi.TransportAware).RegisterTransport(transport)
				connectionCache := cache.NewPlcConnectionCache(driverManager, cache.WithCustomLogger(testutils.ProduceTestingLogger(t)))
				// Also on a failed assertion: the bubble's goroutines have to end.
				closeAll := sync.OnceFunc(func() {
					assert.NoError(t, connectionCache.Close())
					assert.NoError(t, driverManager.Close())
				})
				t.Cleanup(closeAll)
				// 1 ms instead of the test transport's 10 ms: the cache reconnects within lease.Close in 5 s, too short for
				// the connection setup at 10 ms.
				const connectionString = "c-bus:test://interface?simulatedLatency=1ms"

				lease, err := connectionCache.GetConnection(t.Context(), connectionString)
				require.NoError(t, err)
				require.True(t, lease.IsConnected())
				instances := transport.getInstances()
				require.Len(t, instances, 1)
				first := instances[0]
				if !leased {
					require.NoError(t, lease.Close()) // back to the cache, idle
				}

				first.closeByPeer()
				// The codec ends with a fatal transport error, which closes the transport.
				for deadline := time.Now().Add(time.Minute); first.IsConnected(); time.Sleep(10 * time.Millisecond) {
					require.True(t, time.Now().Before(deadline), "the connection whose peer closed it still runs a minute later")
				}
				synctest.Wait() // the fatal error's handler invalidated the connection
				if leased {
					assert.False(t, lease.IsConnected(), "the lease reports the closed connection as connected")
					assert.Error(t, lease.Ping(t.Context()), "the lease pings the closed connection fine")
					require.NoError(t, lease.Close()) // invalid: the cache connects again
				}

				lease, err = connectionCache.GetConnection(t.Context(), connectionString)
				require.NoError(t, err)
				instances = transport.getInstances()
				require.Len(t, instances, 2, "the cache handed out the closed connection again")
				assert.True(t, instances[1].IsConnected())
				assert.True(t, lease.IsConnected())
				assert.NoError(t, lease.Ping(t.Context()))
				require.NoError(t, lease.Close())

				closeAll()
				assert.False(t, instances[1].IsConnected(), "closing the cache left the new connection open")
			})
		})
	}
}

// A fill error is a lost connection when the TCP transport classifies it as fatal (EOF, a closed socket), but
// not the timeouts and context errors a fill on an idle line ends with, nor a full read buffer.
func TestMessageCodec_isConnectionLost(t *testing.T) {
	readTimeout := &net.OpError{Op: "read", Net: "tcp", Err: os.ErrDeadlineExceeded}
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "EOF", err: errors.Wrap(io.EOF, "Error while peeking"), want: true},
		{name: "closed socket", err: errors.Wrap(&net.OpError{Op: "read", Net: "tcp", Err: net.ErrClosed}, "Error while peeking"), want: true},
		{name: "read deadline", err: errors.Wrap(readTimeout, "Error while peeking")},
		{name: "fill timeout", err: errors.Wrap(context.DeadlineExceeded, "Timeout while filling buffer")},
		{name: "cancelled", err: context.Canceled},
		{name: "full read buffer", err: errors.Wrap(bufio.ErrBufferFull, "Error while peeking")},
	}
	codec := &MessageCodec{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, codec.isConnectionLost(&tcp.TransportInstance{}, tt.err))
		})
	}
}

// A connection whose codec stopped is not connected, even before the fatal transport error's handler, which runs
// on its own goroutine, invalidated it: a connection cache asking in between must not hand it out.
func TestConnection_IsConnected_falseOnceTheCodecStopped(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		_options := testutils.EnrichOptionsWithOptionsForTesting(t)
		ti, err := test.NewTransport(_options...).CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
		require.NoError(t, err)
		codec := NewMessageCodec(ti, _options...)
		c := &Connection{messageCodec: codec, log: testutils.ProduceTestingLogger(t)}
		c.DefaultConnection = _default.NewDefaultConnection(c, _options...)
		require.NoError(t, codec.Connect(t.Context()))
		c.SetConnected(true)
		assert.True(t, c.IsConnected())
		assert.NoError(t, c.Ping(t.Context()))

		require.NoError(t, codec.Disconnect()) // stands in for the fatal stop: nothing invalidates the connection
		assert.False(t, c.IsInvalidated())
		assert.False(t, c.IsConnected())
		assert.Error(t, c.Ping(t.Context()))
	})
}
