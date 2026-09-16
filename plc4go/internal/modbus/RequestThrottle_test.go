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

package modbus

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// throttleCodec stands in for a device: it answers every request after a while, and records how
// many requests were on the wire at the same time while doing so. A device that is really a
// TCP<->RTU gateway is the one this exists for - it can only carry one request at a time, and
// answers a flood by resetting the connection.
type throttleCodec struct {
	// panicSends makes the next n sends panic, modelling a codec or transport that blows up
	// mid-send rather than returning an error.
	panicSends int
	t          *testing.T

	wg sync.WaitGroup

	mutex sync.Mutex
	// hold is how long the "device" takes to answer, which is the window a second request would
	// have to overlap with the first one in.
	hold time.Duration
	// failSends is how many of the next sends fail before sending starts working.
	failSends int
	// attempts counts every request handed to the codec, the failed sends included; sends counts
	// the ones that reached the wire.
	attempts int
	sends    int
	// inFlight is how many requests are on the wire right now, and maxInFlight the most there ever
	// were at once. The whole point of the transaction manager is that the latter stays at one.
	inFlight    int
	maxInFlight int
}

func newThrottleCodec(t *testing.T, hold time.Duration) *throttleCodec {
	codec := &throttleCodec{t: t, hold: hold}
	t.Cleanup(codec.wg.Wait)
	return codec
}

func (c *throttleCodec) Connect(context.Context) error { return nil }
func (c *throttleCodec) Disconnect() error             { return nil }
func (c *throttleCodec) IsRunning() bool               { return true }
func (c *throttleCodec) Send(context.Context, string, spi.Message) error {
	return nil
}
func (c *throttleCodec) Expect(context.Context, string, spi.AcceptsMessage, spi.HandleMessage, spi.HandleError) {
}
func (c *throttleCodec) GetDefaultIncomingMessageChannel() chan spi.Message { return nil }

func (c *throttleCodec) SendRequest(ctx context.Context, _ string, message spi.Message, _ spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError) error {
	// The real codec refuses a request whose context is already gone rather than registering an
	// expectation for it (spi/default.defaultCodec.SendRequest).
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "Not sending message as context is aborted")
	}
	c.mutex.Lock()
	c.attempts++
	if c.panicSends > 0 {
		c.panicSends--
		c.mutex.Unlock()
		panic("a codec that panics mid-send")
	}
	if c.failSends > 0 {
		c.failSends--
		c.mutex.Unlock()
		return errors.New("send failed: broken pipe")
	}
	c.sends++
	c.inFlight++
	if c.inFlight > c.maxInFlight {
		c.maxInFlight = c.inFlight
	}
	hold := c.hold
	c.mutex.Unlock()

	c.wg.Go(func() {
		select {
		case <-time.After(hold):
			c.offTheWire()
			_ = handleMessage(answerFor(c.t, message))
		case <-ctx.Done():
			// What the real codec does with an expectation whose context ran out
			// (spi/default.defaultCodec.TimeoutExpectations).
			c.offTheWire()
			_ = handleError(ctx.Err())
		}
	})
	return nil
}

// offTheWire records that the request that was on the wire has been answered. It runs before the
// answer is handed on, so that "in flight" means what it says: the span in which a second request
// would have been on the wire next to the first one.
func (c *throttleCodec) offTheWire() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inFlight--
}

func (c *throttleCodec) counts() (attempts int, sends int, maxInFlight int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.attempts, c.sends, c.maxInFlight
}

// setHold changes how long the "device" takes over the requests it is sent from here on.
func (c *throttleCodec) setHold(hold time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.hold = hold
}

// failNextSends makes the next n sends fail the way a broken pipe does: nothing reaches the wire
// and no expectation is left behind, so the only thing that can give the permit back is the
// driver itself.
func (c *throttleCodec) failNextSends(n int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.failSends = n
}

// panicNextSends makes the next n sends panic instead of returning.
func (c *throttleCodec) panicNextSends(n int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.panicSends = n
}

// answerFor is the response the device sends back for a request. Only the requests these tests
// provoke are covered; anything else is a test that means something other than what it says.
func answerFor(t *testing.T, request spi.Message) spi.Message {
	t.Helper()
	adu, ok := request.(readWriteModel.ModbusTcpADU)
	require.True(t, ok, "expected a ModbusTcpADU, got %T", request)
	var responsePdu readWriteModel.ModbusPDU
	switch pdu := adu.GetPdu().(type) {
	case readWriteModel.ModbusPDUReadHoldingRegistersRequest:
		responsePdu = readWriteModel.NewModbusPDUReadHoldingRegistersResponse(make([]byte, int(pdu.GetQuantity())*2))
	case readWriteModel.ModbusPDUWriteSingleRegisterRequest:
		responsePdu = readWriteModel.NewModbusPDUWriteSingleRegisterResponse(pdu.GetAddress(), pdu.GetValue())
	case readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest:
		responsePdu = readWriteModel.NewModbusPDUWriteMultipleHoldingRegistersResponse(pdu.GetStartingAddress(), pdu.GetQuantity())
	default:
		require.Failf(t, "unexpected request", "no canned answer for a %T", adu.GetPdu())
	}
	return readWriteModel.NewModbusTcpADU(adu.GetTransactionIdentifier(), adu.GetUnitIdentifier(), responsePdu)
}

func throttleTestConnection(t *testing.T, configuration Configuration, codec *throttleCodec) *Connection {
	t.Helper()
	connection := NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	return connection
}

// writeOf builds a write request the way a caller does, through the connection - which is the only
// way to get the interceptor that cuts it into one sub-write per tag.
func writeOf(t *testing.T, connection *Connection, tagCount int) apiModel.PlcWriteRequest {
	t.Helper()
	builder := connection.WriteRequestBuilder()
	for i := 0; i < tagCount; i++ {
		builder.AddTagAddress(fmt.Sprintf("tag%d", i), fmt.Sprintf("4x%05d:UINT", i+1), uint16(i))
	}
	request, err := builder.Build()
	require.NoError(t, err)
	return request
}

func readOf(t *testing.T, connection *Connection, address string) apiModel.PlcReadRequest {
	t.Helper()
	request, err := connection.ReadRequestBuilder().AddTagAddress("tag", address).Build()
	require.NoError(t, err)
	return request
}

// A write request of several tags is cut into one request per tag, and those used to go onto the
// wire all at once, spaced by nothing but a 4ms sleep. They queue now: modbus keeps one PDU in
// flight per connection, the way all three of plc4j's modbus connections do.
func TestConnection_aMultiTagWriteSendsOneRequestAtATime(t *testing.T) {
	codec := newThrottleCodec(t, 20*time.Millisecond)
	connection := throttleTestConnection(t, DefaultConfiguration(), codec)
	request := writeOf(t, connection, 5)

	select {
	case result := <-request.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		for _, tagName := range request.GetTagNames() {
			assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode(tagName), tagName)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("the write never finished")
	}

	_, sends, maxInFlight := codec.counts()
	assert.Equal(t, 5, sends, "every tag is written")
	assert.Equal(t, 1, maxInFlight, "a modbus connection has one request on the wire at a time")
}

// The bound belongs to the connection, not to a single request: reads and writes started
// independently queue behind each other too, because they share the one wire.
func TestConnection_concurrentReadsAndWritesDoNotOverlap(t *testing.T) {
	codec := newThrottleCodec(t, 20*time.Millisecond)
	connection := throttleTestConnection(t, DefaultConfiguration(), codec)
	ctx := testutils.TestContext(t)

	var requests sync.WaitGroup
	for i := 0; i < 3; i++ {
		address := fmt.Sprintf("4x%05d:UINT", i+1)
		read := readOf(t, connection, address)
		write := writeOf(t, connection, 1)
		requests.Go(func() {
			result := <-read.Execute(ctx)
			assert.NoError(t, result.GetErr())
		})
		requests.Go(func() {
			result := <-write.Execute(ctx)
			assert.NoError(t, result.GetErr())
		})
	}
	requests.Wait()

	_, sends, maxInFlight := codec.counts()
	assert.Equal(t, 6, sends, "every request reaches the wire")
	assert.Equal(t, 1, maxInFlight, "reads and writes share the connection's single slot")
}

// The ping shares the wire with everything else, so it queues rather than slipping past a read -
// and it must not deadlock behind the permit it is waiting for.
func TestConnection_pingQueuesBehindARead(t *testing.T) {
	codec := newThrottleCodec(t, 20*time.Millisecond)
	connection := throttleTestConnection(t, DefaultConfiguration(), codec)
	ctx := testutils.TestContext(t)

	var requests sync.WaitGroup
	read := readOf(t, connection, "4x00001:UINT")
	requests.Go(func() {
		result := <-read.Execute(ctx)
		assert.NoError(t, result.GetErr())
	})
	requests.Go(func() { assert.NoError(t, connection.Ping(ctx)) })
	requests.Wait()

	_, sends, maxInFlight := codec.counts()
	assert.Equal(t, 2, sends)
	assert.Equal(t, 1, maxInFlight, "the ping waits its turn like everything else")
}

// A send that fails has to give the permit back. A permit that is never given back is worse than
// no bound at all: the connection stops sending anything, forever, and nothing says why.
func TestConnection_aFailedSendGivesThePermitBack(t *testing.T) {
	codec := newThrottleCodec(t, time.Millisecond)
	configuration := DefaultConfiguration()
	// Short, so that a leaked permit shows up as a failed request rather than as a test that takes
	// five seconds per request to admit it.
	configuration.requestTimeout = 500 * time.Millisecond
	connection := throttleTestConnection(t, configuration, codec)
	ctx := testutils.TestContext(t)

	codec.failNextSends(2)
	for i := 0; i < 2; i++ {
		result := <-writeOf(t, connection, 1).Execute(ctx)
		require.Error(t, result.GetErr(), "the send failed, so the write did")
	}

	result := <-writeOf(t, connection, 1).Execute(ctx)
	require.NoError(t, result.GetErr(), "the failed sends held onto the connection's only permit")
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("tag0"))

	attempts, sends, _ := codec.counts()
	assert.Equal(t, 3, attempts)
	assert.Equal(t, 1, sends)
}

// A request the caller gave up on has to give the permit back as well, whether it was already on
// the wire or still waiting its turn.
func TestConnection_aCancelledRequestGivesThePermitBack(t *testing.T) {
	// Long enough that the request is still unanswered when it is cancelled.
	codec := newThrottleCodec(t, time.Minute)
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 500 * time.Millisecond
	connection := throttleTestConnection(t, configuration, codec)

	abandoned, cancel := context.WithCancel(testutils.TestContext(t))
	first := readOf(t, connection, "4x00001:UINT").Execute(abandoned)
	// Wait until it is really on the wire, so that the permit is really taken.
	require.Eventually(t, func() bool { _, sends, _ := codec.counts(); return sends == 1 },
		2*time.Second, time.Millisecond, "the first read never reached the wire")
	cancel()
	select {
	case <-first:
	case <-time.After(2 * time.Second):
		t.Fatal("the cancelled read never returned")
	}

	// The device answers again, so a permit that came back shows as an answered request.
	codec.setHold(time.Millisecond)
	result := <-readOf(t, connection, "4x00002:UINT").Execute(testutils.TestContext(t))
	require.NoError(t, result.GetErr(), "the cancelled read held onto the connection's only permit")
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("tag"))
}

// TestConnection_aPanicMidSendGivesThePermitBack is the difference between a throttle that
// degrades and one that wedges.
//
// The permit is deliberately NOT released when the operation returns - it has to outlive that,
// until the response handler fires - so a panic between the two would keep it forever. With one
// permit per connection, forever means the connection never sends anything again. Worse, the
// worker pool recovers the panic, so nothing upstream would ever report why.
func TestConnection_aPanicMidSendGivesThePermitBack(t *testing.T) {
	codec := newThrottleCodec(t, time.Millisecond)
	configuration := DefaultConfiguration()
	// Short, so a leaked permit shows up as a failed request rather than a five second wait.
	configuration.requestTimeout = 500 * time.Millisecond
	connection := throttleTestConnection(t, configuration, codec)
	ctx := testutils.TestContext(t)

	codec.panicNextSends(1)
	result := <-writeOf(t, connection, 1).Execute(ctx)
	require.Error(t, result.GetErr(), "a request whose send panicked cannot have succeeded")

	// A healthy request afterwards is the whole point: the permit has to have come back.
	result = <-writeOf(t, connection, 1).Execute(ctx)
	require.NoError(t, result.GetErr(), "the panicking request kept the connection's only permit")
}
