//go:build integration

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

// Package bacnetip_test contains end-to-end integration tests for the BACnet/IP
// driver. They require a running bacpypes3 simulator on localhost:47808 (see
// the README for docker compose instructions) and are gated behind both a
// build tag (`integration`) and an environment variable (`BACNET_IT`) so the
// default `make test` does not try to run them.
//

package bacnetip_test

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/drivers"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/transports"
)

const (
	envSimulatorHost = "BACNET_IT_HOST" // optional override; defaults to 127.0.0.1
	envEnabled       = "BACNET_IT"      // master switch (must be set, e.g. BACNET_IT=1)
	// deviceInstance must mirror BACNET_INSTANCE in docker-compose / device.py.
	// Assertions on the discovered name use this so a stray BACnet packet
	// from another device on the bridge would not pass the test.
	deviceInstance = "1234"
)

func skipIfDisabled(t *testing.T) {
	t.Helper()
	if os.Getenv(envEnabled) == "" {
		t.Skip("BACNET_IT not set; skipping. Run via docker compose + BACNET_IT=1.")
	}
}

func simulatorHost() string {
	if h := os.Getenv(envSimulatorHost); h != "" {
		return h
	}
	return "127.0.0.1"
}

func newDriverManager(t *testing.T) plc4go.PlcDriverManager {
	t.Helper()
	dm := plc4go.NewPlcDriverManager()
	drivers.RegisterBacnetDriver(dm)
	transports.RegisterUdpTransport(dm)
	return dm
}

func connect(t *testing.T, ctx context.Context) plc4go.PlcConnection {
	t.Helper()
	dm := newDriverManager(t)
	conn, err := dm.GetConnection(ctx, "bacnet-ip:udp://"+simulatorHost()+":47808")
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

// readPlcReadResult drains a one-shot result channel with a timeout. The
// per-request channel is closed by the driver after exactly one send, so a
// successful receive plus a timeout case is sufficient.
func readPlcReadResult(t *testing.T, ch <-chan apiModel.PlcReadRequestResult) apiModel.PlcReadRequestResult {
	t.Helper()
	select {
	case r := <-ch:
		require.NotNil(t, r)
		require.NoError(t, r.GetErr())
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("read request timed out")
	}
	return nil
}

func readPlcWriteResult(t *testing.T, ch <-chan apiModel.PlcWriteRequestResult) apiModel.PlcWriteRequestResult {
	t.Helper()
	select {
	case r := <-ch:
		require.NotNil(t, r)
		require.NoError(t, r.GetErr())
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("write request timed out")
	}
	return nil
}

func readPlcSubscriptionResult(t *testing.T, ch <-chan apiModel.PlcSubscriptionRequestResult) apiModel.PlcSubscriptionRequestResult {
	t.Helper()
	select {
	case r := <-ch:
		require.NotNil(t, r)
		require.NoError(t, r.GetErr())
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("subscribe request timed out")
	}
	return nil
}

func readPlcUnsubscriptionResult(t *testing.T, ch <-chan apiModel.PlcUnsubscriptionRequestResult) apiModel.PlcUnsubscriptionRequestResult {
	t.Helper()
	select {
	case r := <-ch:
		require.NotNil(t, r)
		require.NoError(t, r.GetErr())
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("unsubscribe request timed out")
	}
	return nil
}

func TestIT_Discover_FindsSimulator(t *testing.T) {
	skipIfDisabled(t)
	dm := newDriverManager(t)

	ctx, cancel := context.WithTimeout(t.Context(), 12*time.Second)
	defer cancel()

	// Discoverer formats the event name as "device DEVICE:<instance>" — assert
	// the instance matches so an unrelated IAm broadcast can't pass the test.
	found := make(chan apiModel.PlcDiscoveryItem, 8)
	err := dm.Discover(ctx, func(event apiModel.PlcDiscoveryItem) {
		select {
		case found <- event:
		default:
		}
	}, plc4go.WithDiscoveryOptionProtocol("bacnet-ip"))
	require.NoError(t, err)

	for {
		select {
		case ev := <-found:
			if strings.Contains(ev.GetName(), deviceInstance) {
				assert.Equal(t, "bacnet-ip", ev.GetProtocolCode())
				assert.Equal(t, "udp", ev.GetTransportCode())
				return
			}
			t.Logf("discovered non-matching device %q, waiting for instance %s", ev.GetName(), deviceInstance)
		case <-ctx.Done():
			t.Fatalf("discovery timed out without finding device instance %s", deviceInstance)
		}
	}
}

func TestIT_Read_AnalogValueRoundTrip(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("av0", "ANALOG_VALUE,0/PRESENT_VALUE")
	req, err := rrb.Build()
	require.NoError(t, err)

	r := readPlcReadResult(t, req.Execute(ctx))
	resp := r.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("av0"))
	val := resp.GetValue("av0")
	require.NotNil(t, val)
	assert.InDelta(t, 0.0, val.GetFloat32(), 0.01)
}

func TestIT_Read_BinaryValue(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("bv0", "BINARY_VALUE,0/PRESENT_VALUE")
	req, err := rrb.Build()
	require.NoError(t, err)

	r := readPlcReadResult(t, req.Execute(ctx))
	resp := r.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("bv0"))
	val := resp.GetValue("bv0")
	require.NotNil(t, val)
	// BV PresentValue is Enumerated (0=inactive, 1=active). Decoder returns
	// it as PlcUDINT. device.py initializes BV-0 with presentValue="inactive".
	assert.Equal(t, uint32(0), val.GetUint32(), "BV-0 should be inactive (0)")
}

func TestIT_Read_MultiStateValue(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("msv0", "MULTI_STATE_VALUE,0/PRESENT_VALUE")
	req, err := rrb.Build()
	require.NoError(t, err)

	r := readPlcReadResult(t, req.Execute(ctx))
	resp := r.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("msv0"))
	val := resp.GetValue("msv0")
	require.NotNil(t, val)
	// MSV PresentValue is Unsigned, initialized to state 1 in device.py.
	assert.Equal(t, uint64(1), val.GetUint64(), "MSV-0 should be in state 1")
}

func TestIT_Read_UnknownObject_ReturnsNotFound(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	rrb := conn.ReadRequestBuilder()
	// AV.99 isn't defined in device.py — bacpypes3 returns
	// APDUError(unknown-object) which the driver maps to NOT_FOUND.
	rrb.AddTagAddress("ghost", "ANALOG_VALUE,99/PRESENT_VALUE")
	req, err := rrb.Build()
	require.NoError(t, err)

	r := readPlcReadResult(t, req.Execute(ctx))
	resp := r.GetResponse()
	assert.Equal(t, apiModel.PlcResponseCode_NOT_FOUND, resp.GetResponseCode("ghost"),
		"reading a non-existent AV should map to NOT_FOUND")
}

func TestIT_Write_AnalogValueRoundTrip(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	wrb := conn.WriteRequestBuilder()
	wrb.AddTagAddress("av1", "ANALOG_VALUE,1/PRESENT_VALUE", float32(42.5))
	wreq, err := wrb.Build()
	require.NoError(t, err)
	w := readPlcWriteResult(t, wreq.Execute(ctx))
	assert.Equal(t, apiModel.PlcResponseCode_OK, w.GetResponse().GetResponseCode("av1"))

	// Re-read to confirm.
	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("av1", "ANALOG_VALUE,1/PRESENT_VALUE")
	rreq, err := rrb.Build()
	require.NoError(t, err)
	r := readPlcReadResult(t, rreq.Execute(ctx))
	assert.InDelta(t, 42.5, r.GetResponse().GetValue("av1").GetFloat32(), 0.001)
}

func TestIT_WritePropertyMultiple_ThreeAnalogValues(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	// >1 tag in a single WriteRequest triggers WritePropertyMultiple in Writer.go.
	wrb := conn.WriteRequestBuilder()
	wrb.AddTagAddress("av2", "ANALOG_VALUE,2/PRESENT_VALUE", float32(11.25))
	wrb.AddTagAddress("av3", "ANALOG_VALUE,3/PRESENT_VALUE", float32(22.5))
	wrb.AddTagAddress("av4", "ANALOG_VALUE,4/PRESENT_VALUE", float32(33.75))
	wreq, err := wrb.Build()
	require.NoError(t, err)
	w := readPlcWriteResult(t, wreq.Execute(ctx))
	resp := w.GetResponse()
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("av2"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("av3"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("av4"))

	// Re-read all three via ReadPropertyMultiple to confirm they actually
	// committed on the server — a SimpleAck alone doesn't prove the writes
	// landed in the right slots.
	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("av2", "ANALOG_VALUE,2/PRESENT_VALUE")
	rrb.AddTagAddress("av3", "ANALOG_VALUE,3/PRESENT_VALUE")
	rrb.AddTagAddress("av4", "ANALOG_VALUE,4/PRESENT_VALUE")
	rreq, err := rrb.Build()
	require.NoError(t, err)
	r := readPlcReadResult(t, rreq.Execute(ctx))
	rresp := r.GetResponse()
	assert.InDelta(t, 11.25, rresp.GetValue("av2").GetFloat32(), 0.001)
	assert.InDelta(t, 22.5, rresp.GetValue("av3").GetFloat32(), 0.001)
	assert.InDelta(t, 33.75, rresp.GetValue("av4").GetFloat32(), 0.001)
}

func TestIT_Subscribe_AnalogInput_InitialNotification(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 8*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	srb := conn.SubscriptionRequestBuilder()
	srb.AddChangeOfStateTagAddress("ai0", "ANALOG_INPUT,0/PRESENT_VALUE")
	sreq, err := srb.Build()
	require.NoError(t, err)

	sr := readPlcSubscriptionResult(t, sreq.Execute(ctx))
	sresp := sr.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, sresp.GetResponseCode("ai0"))

	// bacpypes3 sends an initial UnconfirmedCOVNotification right after the
	// subscription is accepted — that's what this test waits for. Verifies
	// the unsolicited-message path (codec.HandleMessages fallthrough →
	// defaultIncomingMessageChannel → Connection.routeIncomingMessage →
	// Subscriber.HandleUnconfirmedCOVNotification → consumer).
	notif := make(chan apiModel.PlcSubscriptionEvent, 4)
	for _, h := range sresp.GetSubscriptionHandles() {
		h.Register(func(event apiModel.PlcSubscriptionEvent) {
			select {
			case notif <- event:
			default:
			}
		})
	}

	select {
	case ev := <-notif:
		require.Equal(t, apiModel.PlcResponseCode_OK, ev.GetResponseCode("ai0"))
		val := ev.GetValue("ai0")
		require.NotNil(t, val, "subscription event should carry a value")
		// The initial notification carries whatever AI-0 currently reads.
		// The simulator's sawtooth has been ticking since container start —
		// docker compose builds run for tens of seconds — so the value is
		// usually large (we've seen ~35.0). Assert that it decoded as a
		// finite, non-negative Real, which is the actual property type
		// guarantee from device.py. The follow-on Sawtooth test verifies
		// that subsequent ticks generate fresh notifications.
		f := val.GetFloat32()
		assert.False(t, math.IsNaN(float64(f)), "AI-0 should decode as a real number")
		assert.GreaterOrEqual(t, f, float32(0.0), "sawtooth values are non-negative")
		assert.Less(t, f, float32(100.0), "sawtooth wraps at 100, so values stay in [0,100)")
		t.Logf("initial COV notification: AI-0 = %v", f)
	case <-ctx.Done():
		t.Fatal("did not receive the initial COV notification within timeout")
	}
}

func TestIT_Subscribe_ReceivesSawtoothChange(t *testing.T) {
	skipIfDisabled(t)
	// Sawtooth ticks every 2s; covIncrement=0.5 means every tick fires a COV.
	// We need ≥one full sawtooth period (2s) after the initial notification,
	// plus headroom for startup variance.
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	srb := conn.SubscriptionRequestBuilder()
	srb.AddChangeOfStateTagAddress("ai0", "ANALOG_INPUT,0/PRESENT_VALUE")
	sreq, err := srb.Build()
	require.NoError(t, err)

	sr := readPlcSubscriptionResult(t, sreq.Execute(ctx))
	sresp := sr.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, sresp.GetResponseCode("ai0"))

	notif := make(chan apiModel.PlcSubscriptionEvent, 8)
	for _, h := range sresp.GetSubscriptionHandles() {
		h.Register(func(event apiModel.PlcSubscriptionEvent) {
			select {
			case notif <- event:
			default:
			}
		})
	}

	// Capture the initial notification. bacpypes3 sends one immediately on
	// subscribe with the *current* AI-0 value — which has been ticking since
	// the simulator started, so it's already non-zero. Just accepting any
	// non-zero value (which my earlier attempt did) would pass on this
	// initial echo alone and prove nothing about sawtooth-driven COV.
	var initialVal float32
	select {
	case ev := <-notif:
		require.Equal(t, apiModel.PlcResponseCode_OK, ev.GetResponseCode("ai0"))
		initialVal = ev.GetValue("ai0").GetFloat32()
		t.Logf("initial COV notification (will discard): AI-0 = %v", initialVal)
	case <-ctx.Done():
		t.Fatal("never received the initial COV notification")
	}

	// Now wait for a CHANGE: a subsequent COV with a value different from
	// initialVal. Sawtooth bumps by 1.0 every 2s and covIncrement=0.5, so
	// the next tick after initial should trigger one.
	deadline := time.After(6 * time.Second)
	for {
		select {
		case ev := <-notif:
			require.Equal(t, apiModel.PlcResponseCode_OK, ev.GetResponseCode("ai0"))
			val := ev.GetValue("ai0").GetFloat32()
			t.Logf("subsequent COV: AI-0 = %v (initial was %v)", val, initialVal)
			if val != initialVal {
				return
			}
		case <-deadline:
			t.Fatalf("AI-0 never changed after subscribe; stuck at %v — sawtooth or covIncrement filter not firing", initialVal)
		case <-ctx.Done():
			t.Fatalf("context cancelled while waiting for sawtooth-driven COV after initialVal=%v", initialVal)
		}
	}
}

// ── Tier 2: extended integration coverage ──────────────────────────────────

func TestIT_Read_DeviceObjectName(t *testing.T) {
	// Exercises the CharacterString decode path. device.py sets
	// DEVICE:1234/OBJECT_NAME = "plc4x-it".
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("name", "DEVICE,1234/OBJECT_NAME")
	req, err := rrb.Build()
	require.NoError(t, err)

	r := readPlcReadResult(t, req.Execute(ctx))
	resp := r.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("name"))
	val := resp.GetValue("name")
	require.NotNil(t, val)
	assert.Equal(t, "plc4x-it", val.GetString())
}

func TestIT_Read_MultipleTags_ReadPropertyMultiple(t *testing.T) {
	// Hits Reader.go's >1-tag branch which builds a
	// BACnetConfirmedServiceRequestReadPropertyMultiple. The earlier
	// WPM-readback only proved this works when all tags are the same
	// object-type; here we mix AV / BV / MSV to force the multi-typed
	// dispatch in the response decoder.
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("av0", "ANALOG_VALUE,0/PRESENT_VALUE")
	rrb.AddTagAddress("bv0", "BINARY_VALUE,0/PRESENT_VALUE")
	rrb.AddTagAddress("msv0", "MULTI_STATE_VALUE,0/PRESENT_VALUE")
	req, err := rrb.Build()
	require.NoError(t, err)

	r := readPlcReadResult(t, req.Execute(ctx))
	resp := r.GetResponse()
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("av0"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("bv0"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, resp.GetResponseCode("msv0"))
	assert.InDelta(t, 0.0, resp.GetValue("av0").GetFloat32(), 0.01)
	assert.Equal(t, uint32(0), resp.GetValue("bv0").GetUint32())
	assert.Equal(t, uint64(1), resp.GetValue("msv0").GetUint64())
}

func TestIT_Write_BinaryValue(t *testing.T) {
	// Exercises plcValueToApplicationTag's Boolean → Enumerated hint path,
	// the BV setter on bacpypes3, and the readback's BACnetBinaryPVTagged
	// decode through taggedEnumToPlcValue.
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	wrb := conn.WriteRequestBuilder()
	wrb.AddTagAddress("bv1", "BINARY_VALUE,1/PRESENT_VALUE", true)
	wreq, err := wrb.Build()
	require.NoError(t, err)
	w := readPlcWriteResult(t, wreq.Execute(ctx))
	require.Equal(t, apiModel.PlcResponseCode_OK, w.GetResponse().GetResponseCode("bv1"))

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("bv1", "BINARY_VALUE,1/PRESENT_VALUE")
	rreq, err := rrb.Build()
	require.NoError(t, err)
	r := readPlcReadResult(t, rreq.Execute(ctx))
	assert.Equal(t, uint32(1), r.GetResponse().GetValue("bv1").GetUint32(),
		"BV-1 should read back as 1 (active) after Write(true)")
}

func TestIT_Write_MultiStateValue(t *testing.T) {
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	wrb := conn.WriteRequestBuilder()
	// MSV.0 has numberOfStates=4; state 3 is valid.
	wrb.AddTagAddress("msv0", "MULTI_STATE_VALUE,0/PRESENT_VALUE", uint32(3))
	wreq, err := wrb.Build()
	require.NoError(t, err)
	w := readPlcWriteResult(t, wreq.Execute(ctx))
	require.Equal(t, apiModel.PlcResponseCode_OK, w.GetResponse().GetResponseCode("msv0"))

	rrb := conn.ReadRequestBuilder()
	rrb.AddTagAddress("msv0", "MULTI_STATE_VALUE,0/PRESENT_VALUE")
	rreq, err := rrb.Build()
	require.NoError(t, err)
	r := readPlcReadResult(t, rreq.Execute(ctx))
	assert.Equal(t, uint64(3), r.GetResponse().GetValue("msv0").GetUint64(),
		"MSV-0 should read back as 3 after the write")
}

func TestIT_Write_ReadOnlyProperty_Fails(t *testing.T) {
	// OBJECT_TYPE is universally read-only — bacpypes3 returns
	// APDUError(property, write-access-denied), which the driver should
	// map to ACCESS_DENIED. (AI.PresentValue would be a more obvious
	// candidate, but bacpypes3 doesn't enforce outOfService=False for
	// writes on AI in 0.0.102 — so we'd silently get an OK.) Confirms
	// the error mapping on the write path (counterpart to read NOT_FOUND).
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	wrb := conn.WriteRequestBuilder()
	wrb.AddTagAddress("av0type", "ANALOG_VALUE,0/OBJECT_TYPE", uint32(2))
	wreq, err := wrb.Build()
	require.NoError(t, err)
	w := readPlcWriteResult(t, wreq.Execute(ctx))
	code := w.GetResponse().GetResponseCode("av0type")
	// We accept either ACCESS_DENIED (the canonical mapping for
	// write-access-denied) or INVALID_DATA (some BACnet stacks REJECT
	// the request outright before producing an APDUError). The key is
	// that it MUST NOT be OK.
	assert.NotEqual(t, apiModel.PlcResponseCode_OK, code,
		"writing to read-only OBJECT_TYPE should not succeed; got %v", code)
	t.Logf("write to read-only OBJECT_TYPE returned response code %v", code)
}

func TestIT_ConcurrentReads_AllSucceed(t *testing.T) {
	// Fire 5 reads in parallel on the same connection. The transaction
	// manager must serialize them onto the wire and dispatch responses
	// to the right caller. A race in invoke-id allocation or expectation
	// matching surfaces as a deadlock or cross-talk.
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	const N = 5
	results := make([]<-chan apiModel.PlcReadRequestResult, N)
	for i := range N {
		rrb := conn.ReadRequestBuilder()
		rrb.AddTagAddress(fmt.Sprintf("av%d", i), fmt.Sprintf("ANALOG_VALUE,%d/PRESENT_VALUE", i))
		req, err := rrb.Build()
		require.NoError(t, err)
		results[i] = req.Execute(ctx)
	}
	for i, ch := range results {
		r := readPlcReadResult(t, ch)
		key := fmt.Sprintf("av%d", i)
		require.Equal(t, apiModel.PlcResponseCode_OK, r.GetResponse().GetResponseCode(key),
			"AV.%d should read OK; transaction manager may have crossed wires", i)
	}
}

func TestIT_Unsubscribe_StopsNotifications(t *testing.T) {
	// Subscribe, capture one notification, unsubscribe, then verify no more
	// notifications arrive in the next 3s (longer than sawtooth period).
	// Exercises Subscriber.Unsubscribe (sends SubscribeCOV with lifetime=0).
	skipIfDisabled(t)
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	conn := connect(t, ctx)

	srb := conn.SubscriptionRequestBuilder()
	srb.AddChangeOfStateTagAddress("ai0", "ANALOG_INPUT,0/PRESENT_VALUE")
	sreq, err := srb.Build()
	require.NoError(t, err)
	sr := readPlcSubscriptionResult(t, sreq.Execute(ctx))
	sresp := sr.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, sresp.GetResponseCode("ai0"))

	notif := make(chan apiModel.PlcSubscriptionEvent, 16)
	handles := sresp.GetSubscriptionHandles()
	require.NotEmpty(t, handles)
	for _, h := range handles {
		h.Register(func(event apiModel.PlcSubscriptionEvent) {
			select {
			case notif <- event:
			default:
			}
		})
	}

	// Drain the initial notification.
	select {
	case <-notif:
	case <-time.After(5 * time.Second):
		t.Fatal("never got the initial notification — subscribe likely failed silently")
	}

	// Unsubscribe.
	ureq, err := conn.UnsubscriptionRequestBuilder().AddHandles(handles...).Build()
	require.NoError(t, err)
	uresp := readPlcUnsubscriptionResult(t, ureq.Execute(ctx))
	require.NotNil(t, uresp.GetResponse())

	// Drain any in-flight notifications that snuck in between unsubscribe
	// send and bacpypes3 processing it.
	drainDeadline := time.After(500 * time.Millisecond)
drain:
	for {
		select {
		case <-notif:
		case <-drainDeadline:
			break drain
		}
	}

	// Now wait 3s — sawtooth ticks every 2s, so a working subscription
	// would deliver at least one new notification in this window. Zero
	// means unsubscribe took effect.
	select {
	case ev := <-notif:
		t.Fatalf("got a notification after unsubscribe: AI-0=%v", ev.GetValue("ai0").GetFloat32())
	case <-time.After(3 * time.Second):
		// Expected — no more notifications.
	}
}
