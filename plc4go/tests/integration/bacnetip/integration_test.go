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
//go:build integration

package bacnetip_test

import (
	"context"
	"os"
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

func TestIT_Discover_FindsSimulator(t *testing.T) {
	skipIfDisabled(t)
	dm := newDriverManager(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	found := make(chan struct{}, 1)
	err := dm.Discover(ctx, func(event apiModel.PlcDiscoveryItem) {
		select {
		case found <- struct{}{}:
		default:
		}
	}, plc4go.WithDiscoveryOptionProtocol("bacnet-ip"))
	require.NoError(t, err)

	select {
	case <-found:
	case <-ctx.Done():
		t.Fatal("discovery timed out without finding the simulator")
	}
}

func TestIT_Read_AnalogValueRoundTrip(t *testing.T) {
	skipIfDisabled(t)
	dm := newDriverManager(t)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	conn, err := dm.GetConnection(ctx, "bacnet-ip:udp://"+simulatorHost()+":47808")
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

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

func TestIT_Write_AnalogValueRoundTrip(t *testing.T) {
	skipIfDisabled(t)
	dm := newDriverManager(t)

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	conn, err := dm.GetConnection(ctx, "bacnet-ip:udp://"+simulatorHost()+":47808")
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

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

func TestIT_Subscribe_AnalogInput(t *testing.T) {
	skipIfDisabled(t)
	dm := newDriverManager(t)

	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()

	conn, err := dm.GetConnection(ctx, "bacnet-ip:udp://"+simulatorHost()+":47808")
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	srb := conn.SubscriptionRequestBuilder()
	srb.AddChangeOfStateTagAddress("ai0", "ANALOG_INPUT,0/PRESENT_VALUE")
	sreq, err := srb.Build()
	require.NoError(t, err)

	sr := readPlcSubscriptionResult(t, sreq.Execute(ctx))
	sresp := sr.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, sresp.GetResponseCode("ai0"))

	notif := make(chan apiModel.PlcSubscriptionEvent, 4)
	for _, h := range sresp.GetSubscriptionHandles() {
		h.Register(func(event apiModel.PlcSubscriptionEvent) {
			select {
			case notif <- event:
			default:
			}
		})
	}

	// device.py bumps AI-0 every 2 seconds — allow up to 10s.
	select {
	case <-notif:
	case <-ctx.Done():
		t.Fatal("did not receive a COV notification within timeout")
	}
}
