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
	"bytes"
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// slotTestConnection opens a connection of driver on its own test transport and hands back the transport
// instance, so the test can play the C-Bus interface behind it. The connection setup is skipped: its
// context is cancelled once the connection is there, so the background setup gives up at once and sends
// nothing that the test's interceptor would have to tell apart from the reads.
func slotTestConnection(t *testing.T, driver *Driver) (*Connection, *test.TransportInstance) {
	t.Helper()
	setupCtx, cancelSetup := context.WithCancel(testutils.TestContext(t))
	connection, err := driver.GetConnection(setupCtx, url.URL{Scheme: "test"}, map[string]transports.Transport{
		"test": test.NewTransport(testutils.EnrichOptionsWithOptionsForTesting(t)...),
	}, map[string][]string{})
	cancelSetup()
	require.NoError(t, err)
	c := connection.(*Connection)
	t.Cleanup(func() {
		assert.NoError(t, c.Close())
	})
	c.wg.Wait() // the setup has given up
	return c, c.messageCodec.GetTransportInstance().(*test.TransportInstance)
}

// identifyTypeRead reads the Type attribute of unit 2 over connection, bounded by timeout.
func identifyTypeRead(t *testing.T, ctx context.Context, connection *Connection) <-chan apiModel.PlcReadRequestResult {
	t.Helper()
	request, err := connection.ReadRequestBuilder().
		AddTag("type", NewCALIdentifyTag(readWriteModel.NewUnitAddress(2), nil, readWriteModel.Attribute_Type, 1)).
		Build()
	require.NoError(t, err)
	return request.Execute(ctx)
}

// isIdentifyRequest tells the identify read's request from anything else the connection writes.
func isIdentifyRequest(data []byte) bool {
	return bytes.HasPrefix(data, []byte(`\`)) && bytes.HasSuffix(data, []byte("\r"))
}

// Two connections of one driver are two C-Bus interfaces, each with its own bus. A read on one of them that
// the interface never answers must not hold up a read on the other: the request slot (one request on the
// wire at a time) belongs to the connection, so the second read is sent at once and answered well within
// its own timeout instead of waiting for the first one's slot and timing out queued.
func TestDriver_requestSlotPerConnection(t *testing.T) {
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...).(*Driver)
	driver.SetAwaitSetupComplete(false)
	t.Cleanup(func() {
		assert.NoError(t, driver.Close())
	})

	silent, silentTransport := slotTestConnection(t, driver)
	silentSent := make(chan struct{})
	silentTransport.SetWriteInterceptor(func(_ *test.TransportInstance, data []byte) {
		if isIdentifyRequest(data) {
			select {
			case <-silentSent:
			default:
				close(silentSent) // and never answer
			}
		}
	})
	answering, answeringTransport := slotTestConnection(t, driver)
	answeringTransport.SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
		if isIdentifyRequest(data) {
			transportInstance.FillReadBuffer([]byte("g.890150435F434E49454421\r\n"))
		}
	})

	// The silent interface's read holds its connection's request slot until its timeout ends.
	silentCtx, cancelSilent := context.WithTimeout(testutils.TestContext(t), 8*time.Second)
	defer cancelSilent()
	silentResult := identifyTypeRead(t, silentCtx, silent)
	select {
	case <-silentSent:
	case <-time.After(10 * time.Second):
		t.Fatal("the read on the silent interface was never sent")
	}

	answeringCtx, cancelAnswering := context.WithTimeout(testutils.TestContext(t), 3*time.Second)
	defer cancelAnswering()
	started := time.Now()
	result := <-identifyTypeRead(t, answeringCtx, answering)
	require.NoError(t, result.GetErr())
	response := result.GetResponse()
	require.NotNil(t, response)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("type"),
		"the read on the answering interface waited %s for the silent interface's request slot", time.Since(started))
	if value := response.GetValue("type"); assert.NotNil(t, value) {
		assert.Equal(t, "PC_CNIED", value.GetString())
	}

	// The silent read is still waiting for its reply, and ends with its timeout.
	select {
	case <-silentResult:
		t.Fatal("the read on the silent interface ended before its timeout")
	default:
	}
	select {
	case result := <-silentResult:
		if result.GetResponse() != nil {
			assert.NotEqual(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("type"))
		}
	case <-time.After(20 * time.Second):
		t.Fatal("the read on the silent interface did not end with its timeout")
	}
}

// Closing or invalidating a connection ends the request it has in flight, shuts its request slot, so nothing
// waits for a turn that never comes, and leaves every other connection of the driver able to send.
func TestConnection_closeShutsOnlyItsOwnRequestSlot(t *testing.T) {
	for name, end := range map[string]func(*Connection){
		"close":      func(c *Connection) { assert.NoError(t, c.Close()) },
		"invalidate": func(c *Connection) { c.Invalidate() },
	} {
		t.Run(name, func(t *testing.T) {
			driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...).(*Driver)
			driver.SetAwaitSetupComplete(false)
			t.Cleanup(func() {
				assert.NoError(t, driver.Close())
			})
			ended, endedTransport := slotTestConnection(t, driver)
			inFlightSent := make(chan struct{})
			endedTransport.SetWriteInterceptor(func(_ *test.TransportInstance, data []byte) {
				if isIdentifyRequest(data) {
					select {
					case <-inFlightSent:
					default:
						close(inFlightSent) // and never answer
					}
				}
			})
			other, otherTransport := slotTestConnection(t, driver)
			otherTransport.SetWriteInterceptor(func(transportInstance *test.TransportInstance, data []byte) {
				if isIdentifyRequest(data) {
					transportInstance.FillReadBuffer([]byte("g.890150435F434E49454421\r\n"))
				}
			})

			inFlightCtx, cancelInFlight := context.WithTimeout(testutils.TestContext(t), 30*time.Second)
			defer cancelInFlight()
			inFlight := identifyTypeRead(t, inFlightCtx, ended)
			select {
			case <-inFlightSent:
			case <-time.After(10 * time.Second):
				t.Fatal("the read on the connection to be ended was never sent")
			}

			end(ended)

			select {
			case result := <-inFlight:
				if result.GetResponse() != nil {
					assert.NotEqual(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("type"))
				}
			case <-time.After(10 * time.Second):
				t.Fatal("the read in flight did not end with its connection")
			}

			ctx, cancel := context.WithTimeout(testutils.TestContext(t), 3*time.Second)
			defer cancel()
			assert.ErrorContains(t, ended.tm.StartTransaction("after close").AwaitCompletion(ctx), "shutdown",
				"the ended connection still hands out request slots")

			result := <-identifyTypeRead(t, ctx, other)
			require.NoError(t, result.GetErr())
			require.NotNil(t, result.GetResponse())
			assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("type"),
				"ending one connection stopped another from sending")
		})
	}
}

// A read queued on a connection behind one that is still on the wire must not outlive the connection: closing it
// ends the queued read at once, as a failed one, not after its whole timeout or as a read without an answer.
func TestConnection_closeEndsTheQueuedReads(t *testing.T) {
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...).(*Driver)
	driver.SetAwaitSetupComplete(false)
	t.Cleanup(func() {
		assert.NoError(t, driver.Close())
	})
	connection, transport := slotTestConnection(t, driver)
	sent := make(chan struct{}, 10)
	transport.SetWriteInterceptor(func(_ *test.TransportInstance, data []byte) {
		if isIdentifyRequest(data) {
			sent <- struct{}{} // and never answer
		}
	})

	ctx, cancel := context.WithTimeout(testutils.TestContext(t), 8*time.Second)
	defer cancel()
	inFlight := identifyTypeRead(t, ctx, connection)
	select {
	case <-sent:
	case <-time.After(5 * time.Second):
		t.Fatal("the first read was never sent")
	}
	queued := identifyTypeRead(t, ctx, connection)
	// The second read queues behind the first one, which holds the connection's only request slot.
	time.Sleep(200 * time.Millisecond)
	select {
	case <-sent:
		t.Fatal("the second read was sent while the first one was still on the wire")
	default:
	}

	started := time.Now()
	require.NoError(t, connection.Close())
	for name, results := range map[string]<-chan apiModel.PlcReadRequestResult{"in flight": inFlight, "queued": queued} {
		select {
		case result := <-results:
			require.NoError(t, result.GetErr(), name)
			require.NotNil(t, result.GetResponse(), name)
			assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, result.GetResponse().GetResponseCode("type"), name)
		case <-time.After(5 * time.Second):
			t.Fatalf("the %s read did not end with its connection", name)
		}
	}
	assert.Less(t, time.Since(started), 2*time.Second, "the reads waited for their timeout instead of ending with the connection")
}
