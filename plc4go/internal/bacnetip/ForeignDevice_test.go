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

package bacnetip

import (
	"context"
	"net"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	model "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

// fakeBBMDRegistration is one Register-Foreign-Device the fake BBMD received.
type fakeBBMDRegistration struct {
	ttl uint16
	src *net.UDPAddr
}

// fakeBBMD answers each Register-Foreign-Device with the BVLC-Result code it was given, or not
// at all when silent.
type fakeBBMD struct {
	conn   *net.UDPConn
	code   model.BVLCResultCode
	silent bool
	log    zerolog.Logger
	wg     sync.WaitGroup

	mu            sync.Mutex
	registrations []fakeBBMDRegistration
}

func startFakeBBMD(t *testing.T, log zerolog.Logger, code model.BVLCResultCode, silent bool) *fakeBBMD {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	b := &fakeBBMD{conn: conn, code: code, silent: silent, log: log}
	b.wg.Go(b.serve)
	t.Cleanup(func() {
		_ = b.conn.Close()
		b.wg.Wait()
	})
	return b
}

func (b *fakeBBMD) address() string { return b.conn.LocalAddr().String() }

func (b *fakeBBMD) serve() {
	buf := make([]byte, 1500)
	for {
		n, src, err := b.conn.ReadFromUDP(buf)
		if err != nil {
			return // socket closed
		}
		bvlc, err := model.BVLCParse[model.BVLC](context.Background(), buf[:n])
		if err != nil {
			b.log.Warn().Err(err).Msg("fake BBMD: parse BVLC")
			continue
		}
		register, ok := bvlc.(model.BVLCRegisterForeignDevice)
		if !ok {
			b.log.Warn().Type("bvlc", bvlc).Msg("fake BBMD: not a Register-Foreign-Device")
			continue
		}
		b.mu.Lock()
		b.registrations = append(b.registrations, fakeBBMDRegistration{ttl: register.GetTtl(), src: src})
		b.mu.Unlock()
		if b.silent {
			continue
		}
		answer, err := model.NewBVLCResult(b.code).Serialize()
		if err != nil {
			b.log.Error().Err(err).Msg("fake BBMD: serialize BVLC-Result")
			continue
		}
		if _, err := b.conn.WriteToUDP(answer, src); err != nil {
			b.log.Error().Err(err).Msg("fake BBMD: write BVLC-Result")
		}
	}
}

func (b *fakeBBMD) received() []fakeBBMDRegistration {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]fakeBBMDRegistration(nil), b.registrations...)
}

// connectForeignDeviceTest connects the driver to device with the given connection options on an
// ephemeral local port.
func connectForeignDeviceTest(t *testing.T, log zerolog.Logger, device *fakeBacnetDevice, connectionOptions map[string][]string) *Connection {
	t.Helper()
	driver := NewDriver(options.WithCustomLogger(log))
	t.Cleanup(func() { _ = driver.Close() })
	connectionOptions["local-port"] = []string{"0"}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	connection, err := driver.GetConnection(ctx,
		url.URL{Scheme: "udp", Host: "127.0.0.1:" + strconv.Itoa(device.port())},
		map[string]transports.Transport{"udp": udp.NewTransport(options.WithCustomLogger(log))},
		connectionOptions,
	)
	require.NoError(t, err)
	return connection.(*Connection)
}

// TestForeignDevice_RegistersOnConnectFromTheConnectionSocket pins that a connection with a BBMD
// registers with the configured TTL right after connecting, from the socket it also talks to the
// device on, and that the device is still reachable through that now unconnected socket.
func TestForeignDevice_RegistersOnConnectFromTheConnectionSocket(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeBacnetDevice(t, log)
	t.Cleanup(device.stop)
	bbmd := startFakeBBMD(t, log, model.BVLCResultCode_SUCCESSFUL_COMPLETION, false)

	connection := connectForeignDeviceTest(t, log, device, map[string][]string{
		"ForeignDeviceBBMD": {bbmd.address()},
		"ForeignDeviceTTL":  {"42"},
	})
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })

	require.Eventually(t, func() bool { return len(bbmd.received()) > 0 }, 5*time.Second, 10*time.Millisecond)
	registration := bbmd.received()[0]
	assert.Equal(t, uint16(42), registration.ttl)

	code, value, ok := readPresentValue(t, connection)
	require.True(t, ok, "the device must still answer through the unconnected socket")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, float32(23.5), value)
	device.mu.Lock()
	deviceSrc := device.lastSrc
	device.mu.Unlock()
	require.NotNil(t, deviceSrc)
	assert.Equal(t, deviceSrc.Port, registration.src.Port,
		"the registration must come from the socket the BBMD forwards broadcasts to: the connection's")
}

// TestForeignDevice_Renews pins the renewal at TTL/2 (500 ms for a TTL of 1 s), also after a
// NAK or no answer: a failed registration is retried on the next tick, not in a busy loop.
func TestForeignDevice_Renews(t *testing.T) {
	for _, tc := range []struct {
		name   string
		code   model.BVLCResultCode
		silent bool
	}{
		{name: "accepted", code: model.BVLCResultCode_SUCCESSFUL_COMPLETION},
		{name: "nak", code: model.BVLCResultCode_REGISTER_FOREIGN_DEVICE_NAK},
		{name: "no answer", silent: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			log := testutils.ProduceTestingLogger(t)
			device := startFakeBacnetDevice(t, log)
			t.Cleanup(device.stop)
			bbmd := startFakeBBMD(t, log, tc.code, tc.silent)

			connection := connectForeignDeviceTest(t, log, device, map[string][]string{
				"ForeignDeviceBBMD": {bbmd.address()},
				"ForeignDeviceTTL":  {"1"},
				"ApduTimeoutMs":     {"200"},
			})
			t.Cleanup(func() { assert.NoError(t, connection.Close()) })

			require.Eventually(t, func() bool { return len(bbmd.received()) >= 3 }, 5*time.Second, 10*time.Millisecond)
			registrations := bbmd.received()
			for _, registration := range registrations {
				assert.Equal(t, uint16(1), registration.ttl)
			}
			// Three registrations take two intervals; a busy loop would have sent far more.
			assert.Less(t, len(registrations), 10)
		})
	}
}

// TestForeignDevice_CloseStopsRenewal pins that Close ends the registration loop and returns only
// after it has ended, and that nothing is sent to the BBMD afterwards.
func TestForeignDevice_CloseStopsRenewal(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeBacnetDevice(t, log)
	t.Cleanup(device.stop)
	bbmd := startFakeBBMD(t, log, model.BVLCResultCode_SUCCESSFUL_COMPLETION, false)

	connection := connectForeignDeviceTest(t, log, device, map[string][]string{
		"ForeignDeviceBBMD": {bbmd.address()},
		"ForeignDeviceTTL":  {"1"},
	})
	require.NotNil(t, connection.foreignDeviceDone)
	require.Eventually(t, func() bool { return len(bbmd.received()) > 0 }, 5*time.Second, 10*time.Millisecond)

	require.NoError(t, connection.Close())
	select {
	case <-connection.foreignDeviceDone:
	default:
		t.Fatal("Close returned while the foreign-device registration loop was still running")
	}
	sent := len(bbmd.received())
	time.Sleep(3 * foreignDeviceRenewalInterval(1))
	assert.Equal(t, sent, len(bbmd.received()), "no registration after Close")
}

// TestForeignDevice_InvalidatedConnectionStopsRenewal pins that the registration loop ends with the
// connection, also when the connection ends without Close: a fatal transport error invalidates it
// and closes its transport, but nobody may ever call Close.
func TestForeignDevice_InvalidatedConnectionStopsRenewal(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeBacnetDevice(t, log)
	t.Cleanup(device.stop)
	bbmd := startFakeBBMD(t, log, model.BVLCResultCode_SUCCESSFUL_COMPLETION, false)

	connection := connectForeignDeviceTest(t, log, device, map[string][]string{
		"ForeignDeviceBBMD": {bbmd.address()},
		"ForeignDeviceTTL":  {"1"},
	})
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })
	require.NotNil(t, connection.foreignDeviceDone)
	require.Eventually(t, func() bool { return len(bbmd.received()) > 0 }, 5*time.Second, 10*time.Millisecond)

	invalidator, ok := connection.DefaultConnection.(interface{ Invalidate() })
	require.True(t, ok)
	invalidator.Invalidate()
	require.False(t, connection.IsConnected())

	select {
	case <-connection.foreignDeviceDone:
	case <-time.After(10 * foreignDeviceRenewalInterval(1)):
		t.Fatal("the foreign-device registration loop outlived the invalidated connection")
	}
}

// TestForeignDevice_NoBBMDNoRegistration pins that a connection without ForeignDeviceBBMD starts
// no registration and keeps its socket connected to the device.
func TestForeignDevice_NoBBMDNoRegistration(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startFakeBacnetDevice(t, log)
	t.Cleanup(device.stop)

	connection := connectForeignDeviceTest(t, log, device, map[string][]string{
		"ForeignDeviceTTL": {"1"},
	})
	t.Cleanup(func() { assert.NoError(t, connection.Close()) })

	assert.Nil(t, connection.foreignDeviceDone)
	instance, ok := connection.messageCodec.(*MessageCodec).GetTransportInstance().(*udp.TransportInstance)
	require.True(t, ok)
	assert.False(t, instance.Unconnected)
}

func TestSplitForeignDeviceBBMD(t *testing.T) {
	for _, tc := range []struct {
		bbmd    string
		host    string
		port    uint16
		wantErr bool
	}{
		{bbmd: "10.0.0.5", host: "10.0.0.5", port: 47808},
		{bbmd: "10.0.0.5:47809", host: "10.0.0.5", port: 47809},
		{bbmd: "bbmd.example:47808", host: "bbmd.example", port: 47808},
		{bbmd: "10.0.0.5:0", wantErr: true},
		{bbmd: "10.0.0.5:x", wantErr: true},
		{bbmd: ":47808", wantErr: true},
	} {
		t.Run(tc.bbmd, func(t *testing.T) {
			host, port, err := splitForeignDeviceBBMD(tc.bbmd)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.host, host)
			assert.Equal(t, tc.port, port)
		})
	}
}
