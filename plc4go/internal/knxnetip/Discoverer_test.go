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

package knxnetip

import (
	"context"
	"net"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// closeTimeout is the budget a Discoverer.Close gets in the tests below. The
// discovery window itself is 5s per scan, but Close cancels all in-flight scans,
// so this is generous by design: it only has to catch a Close which hangs
// forever (which is what an unbalanced WaitGroup used to produce).
const closeTimeout = 30 * time.Second

func testDiscoverer(t *testing.T) *Discoverer {
	// Deliberately NOT testutils.EnrichOptionsWithOptionsForTesting: that one
	// spins up an extra worker pool which would show up in the goroutine counts
	// these tests compare.
	return NewDiscoverer(options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
}

// requireCloseWithin closes the discoverer and fails the test if Close doesn't
// return within the given budget instead of blocking the test binary forever.
func requireCloseWithin(t *testing.T, d *Discoverer, timeout time.Duration) {
	t.Helper()
	closed := make(chan error, 1)
	go func() {
		closed <- d.Close()
	}()
	select {
	case err := <-closed:
		require.NoError(t, err)
	case <-time.After(timeout):
		dumpGoroutines(t)
		t.Fatalf("Close did not return within %s", timeout)
	}
}

// assertNoGoroutineLeak checks that the number of running goroutines settles
// back at (or below) the baseline. Goroutines the runtime tears down lazily
// make a single sample flaky, hence the retries.
func assertNoGoroutineLeak(t *testing.T, baseline int) {
	t.Helper()
	current := runtime.NumGoroutine()
	for i := 0; i < 100 && current > baseline; i++ {
		time.Sleep(50 * time.Millisecond)
		current = runtime.NumGoroutine()
	}
	if current > baseline {
		dumpGoroutines(t)
		t.Fatalf("goroutine leak detected: baseline %d, now %d", baseline, current)
	}
}

func dumpGoroutines(t *testing.T) {
	t.Helper()
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	t.Logf("goroutine dump:\n%s", buf[:n])
}

// TestDiscovererDiscoverWithCanceledContext pins that a discovery started on an
// already canceled context neither leaks goroutines nor makes Close hang.
func TestDiscovererDiscoverWithCanceledContext(t *testing.T) {
	baseline := runtime.NumGoroutine()

	d := testDiscoverer(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var callbackCalls atomic.Int32
	if err := d.Discover(ctx, func(_ apiModel.PlcDiscoveryItem) { callbackCalls.Add(1) }); err != nil {
		assert.ErrorIs(t, err, context.Canceled)
	}
	requireCloseWithin(t, d, closeTimeout)

	assert.Zero(t, callbackCalls.Load(), "no discovery item can be found on a canceled context")
	assertNoGoroutineLeak(t, baseline)
}

// TestDiscovererDiscoverWithExpiringContext runs a couple of discoveries whose
// context expires while the scans are being set up. This is the shape which used
// to strand the scan goroutines (work submitted to the worker pool was silently
// dropped once the context was done, so the WaitGroup tracking it never reached
// zero and Close blocked forever).
func TestDiscovererDiscoverWithExpiringContext(t *testing.T) {
	baseline := runtime.NumGoroutine()

	for range 5 {
		d := testDiscoverer(t)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		if err := d.Discover(ctx, func(_ apiModel.PlcDiscoveryItem) {}); err != nil {
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		}
		requireCloseWithin(t, d, closeTimeout)
		cancel()
	}

	assertNoGoroutineLeak(t, baseline)
}

// requireScannableInterface skips the test if this machine has no interface a scan could
// run on. Without one the test would silently degrade into "Close with nothing to cancel",
// which the other tests cover already.
func requireScannableInterface(t *testing.T) {
	t.Helper()
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Skipf("cannot enumerate the network interfaces: %v", err)
	}
	for _, netInterface := range interfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ipv4Addr net.IP
			switch typedAddr := addr.(type) {
			case *net.IPNet:
				ipv4Addr = typedAddr.IP.To4()
			case *net.IPAddr:
				ipv4Addr = typedAddr.IP.To4()
			}
			if ipv4Addr != nil && !ipv4Addr.IsLoopback() {
				return
			}
		}
	}
	t.Skip("no non-loopback ipv4 interface to scan on")
}

// TestDiscovererCloseCancelsRunningDiscovery pins that Close aborts scans which
// are still within their discovery window instead of waiting it out.
//
// Unlike the other tests here this one really opens udp sockets and sends multicast
// SearchRequests, which is the only way to get a scan into its discovery window, so it is
// skipped in short mode and on machines without a scannable interface.
func TestDiscovererCloseCancelsRunningDiscovery(t *testing.T) {
	if testing.Short() {
		t.Skip("this test does real network i/o")
	}
	requireScannableInterface(t)

	baseline := runtime.NumGoroutine()

	d := testDiscoverer(t)
	require.NoError(t, d.Discover(context.Background(), func(_ apiModel.PlcDiscoveryItem) {}))

	start := time.Now()
	requireCloseWithin(t, d, closeTimeout)
	assert.Less(t, time.Since(start), discoveryWindow, "Close should cancel the running scans, not wait out the discovery window")

	assertNoGoroutineLeak(t, baseline)
}

// TestDiscovererCloseWithoutDiscover pins that Close is safe (and quick) even if
// no discovery ever ran.
func TestDiscovererCloseWithoutDiscover(t *testing.T) {
	baseline := runtime.NumGoroutine()
	d := testDiscoverer(t)
	requireCloseWithin(t, d, closeTimeout)
	assertNoGoroutineLeak(t, baseline)
}

func buildTestSearchResponse(deviceName string, addr []byte, port uint16) driverModel.SearchResponse {
	friendlyName := make([]byte, 30)
	copy(friendlyName, deviceName)
	return driverModel.NewSearchResponse(
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			driverModel.NewIPAddress(addr),
			port,
		),
		driverModel.NewDIBDeviceInfo(
			8,
			driverModel.KnxMedium_MEDIUM_TP1,
			driverModel.NewDeviceStatus(false),
			driverModel.NewKnxAddress(1, 1, 0),
			driverModel.NewProjectInstallationIdentifier(0, 1),
			[]byte{0, 0, 0, 0, 0, 0},
			driverModel.NewIPAddress([]byte{224, 0, 23, 12}),
			driverModel.NewMACAddress([]byte{0, 0, 0, 0, 0, 0}),
			friendlyName,
		),
		driverModel.NewDIBSuppSvcFamilies(2, nil),
	)
}

// TestDiscovererDiscoveryItemFromMessage pins the (formerly unchecked) type
// assertion on messages coming out of the codec: anything which isn't a
// SearchResponse has to be skipped instead of panicking the scan goroutine.
func TestDiscovererDiscoveryItemFromMessage(t *testing.T) {
	tests := []struct {
		name           string
		message        spi.Message
		wantOk         bool
		wantName       string
		wantRemoteHost string
	}{
		{
			name:    "nil message",
			message: nil,
			wantOk:  false,
		},
		{
			name: "not a search response",
			message: driverModel.NewSearchRequest(
				driverModel.NewHPAIDiscoveryEndpoint(
					driverModel.HostProtocolCode_IPV4_UDP,
					driverModel.NewIPAddress([]byte{192, 168, 42, 11}),
					3671,
				),
			),
			wantOk: false,
		},
		{
			name:           "search response",
			message:        buildTestSearchResponse("KNX IP Router", []byte{192, 168, 42, 11}, 3671),
			wantOk:         true,
			wantName:       "KNX IP Router",
			wantRemoteHost: "192.168.42.11:3671",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := testDiscoverer(t)
			t.Cleanup(func() {
				require.NoError(t, d.Close())
			})
			got, ok := d.discoveryItemFromMessage(tt.message)
			require.Equal(t, tt.wantOk, ok)
			if !tt.wantOk {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tt.wantName, got.GetName())
			assert.Equal(t, "knxnet-ip", got.GetProtocolCode())
			assert.Equal(t, "udp", got.GetTransportCode())
			assert.Equal(t, tt.wantRemoteHost, got.GetTransportUrl().Host)
		})
	}
}
