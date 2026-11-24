//go:build linux

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

package profinet

// Connects to a real Profinet device when PROFINET_TEST_INTERFACE, PROFINET_TEST_DEVICE (its name
// of station) and PROFINET_TEST_GSD (a directory with its GSDML) are set. tools/pnet-peer runs the
// open-source p-net sample device for this in Docker; a Pi with p-net or any device works too.

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

func TestConnectToDevice(t *testing.T) {
	iface, name, gsd := os.Getenv("PROFINET_TEST_INTERFACE"), os.Getenv("PROFINET_TEST_DEVICE"), os.Getenv("PROFINET_TEST_GSD")
	if iface == "" || name == "" || gsd == "" {
		t.Skip("set PROFINET_TEST_INTERFACE, PROFINET_TEST_DEVICE and PROFINET_TEST_GSD to run against a device")
	}
	options := url.Values{"device-name": {name}, "gsd-directory": {gsd}, "reduction-ratio": {"32"}}
	if ip := os.Getenv("PROFINET_TEST_LOCAL_IP"); ip != "" {
		options.Set("local-ip", ip)
	}
	// PROFINET_TEST_ECHO_SLOT puts p-net's echo module (IDM_40: Float32 and Unsigned32 out, the
	// same times a gain parameter back in) in that slot and checks values make the round trip. The
	// gain is a parameter record whose GSDML default is 2 while p-net starts with 1, so doubled
	// values also show the parameterization reached the device.
	echoSlot := os.Getenv("PROFINET_TEST_ECHO_SLOT")
	if echoSlot != "" {
		options.Set("module."+echoSlot, "IDM_40")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	connect := func() plc4go.PlcConnection {
		connection, err := driver.GetConnection(ctx, url.URL{Scheme: rawsocket.TransportCode, Host: iface},
			map[string]transports.Transport{rawsocket.TransportCode: rawsocket.NewTransport()}, options)
		require.NoError(t, err)
		return connection
	}
	connection := connect()
	defer func() { _ = connection.Close() }()
	defer func() {
		if t.Failed() {
			return
		}
		// Close sends the Release call ("Profinet application relation released" in the debug log
		// when the device answers Done), and a new relation can follow on the same device.
		require.NoError(t, connection.Close())
		connection = connect()
		assert.True(t, connection.IsConnected(), "reconnects after a release")
	}()

	c := connection.(*Connection)
	for _, s := range c.ar.plan.Submodules {
		t.Logf("slot %d subslot 0x%04X module 0x%08X submodule 0x%08X in %d out %d (%s)", s.Slot, s.Subslot, s.ModuleIdent, s.SubmoduleIdent, s.Layout.InputLength, s.Layout.OutputLength, s.Resolved.ModuleID())
	}
	// Several watchdog periods: the relation must stay up on our cyclic frames.
	time.Sleep(5 * time.Second)
	assert.True(t, connection.IsConnected(), "the relation stays up")
	_, status, last := c.ar.inputImage()
	assert.WithinDuration(t, time.Now(), last, time.Second, "input frames keep arriving")
	t.Logf("input data status 0x%02X", status)

	if echoSlot == "" {
		return
	}
	write, err := connection.WriteRequestBuilder().
		AddTagAddress("f", echoSlot+".1.OUTPUT.0:REAL", float32(12.5)).
		AddTagAddress("u", echoSlot+".1.OUTPUT.1:UDINT", uint32(4242)).Build()
	require.NoError(t, err)
	written := <-write.Execute(ctx)
	require.NoError(t, written.GetErr())
	require.Equal(t, "OK", written.GetResponse().GetResponseCode("f").String())
	require.Eventually(t, func() bool {
		read, err := connection.ReadRequestBuilder().
			AddTagAddress("f", echoSlot+".1.INPUT.0:REAL").
			AddTagAddress("u", echoSlot+".1.INPUT.1:UDINT").Build()
		if err != nil {
			return false
		}
		result := <-read.Execute(ctx)
		if result.GetErr() != nil || result.GetResponse().GetResponseCode("f").String() != "OK" {
			return false
		}
		f, u := result.GetResponse().GetValue("f").GetFloat32(), result.GetResponse().GetValue("u").GetUint32()
		t.Logf("echo: %v %v", f, u)
		return f == 25 && u == 8484
	}, 5*time.Second, 100*time.Millisecond, "the echo module returns what was written times the GSDML default gain of 2")
}

// TestDeviceAlarms needs PROFINET_TEST_ALARM_BUTTON too: the file p-net's sample device reads as
// its button 2 (pn_dev -d <file>). Each press makes the device raise the next of its demo events
// on its first 8-bit input module, which goes in PROFINET_TEST_ALARM_SLOT (default 2): a process
// alarm, then cyclic data as backup and back, then a diagnosis. The driver must deliver each alarm to the subscription and acknowledge
// it, which the device confirms by acknowledging the AlarmAck frame.
func TestDeviceAlarms(t *testing.T) {
	iface, name, gsd, button := os.Getenv("PROFINET_TEST_INTERFACE"), os.Getenv("PROFINET_TEST_DEVICE"), os.Getenv("PROFINET_TEST_GSD"), os.Getenv("PROFINET_TEST_ALARM_BUTTON")
	if iface == "" || name == "" || gsd == "" || button == "" {
		t.Skip("set PROFINET_TEST_INTERFACE, PROFINET_TEST_DEVICE, PROFINET_TEST_GSD and PROFINET_TEST_ALARM_BUTTON to run against p-net's sample device")
	}
	slot := os.Getenv("PROFINET_TEST_ALARM_SLOT")
	if slot == "" {
		slot = "2"
	}
	options := url.Values{"device-name": {name}, "gsd-directory": {gsd}, "reduction-ratio": {"32"}, "module." + slot: {"IDM_30"}}
	if ip := os.Getenv("PROFINET_TEST_LOCAL_IP"); ip != "" {
		options.Set("local-ip", ip)
	}
	require.NoError(t, os.WriteFile(button, []byte("0\n"), 0o644))
	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()
	driver := NewDriver(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	connection, err := driver.GetConnection(ctx, url.URL{Scheme: rawsocket.TransportCode, Host: iface},
		map[string]transports.Transport{rawsocket.TransportCode: rawsocket.NewTransport()}, options)
	require.NoError(t, err)
	defer func() { _ = connection.Close() }()
	c := connection.(*Connection)

	request, err := connection.SubscriptionRequestBuilder().AddEventTagAddress("alarms", "ALARM:"+slot).Build()
	require.NoError(t, err)
	subscribed := <-request.Execute(ctx)
	require.NoError(t, subscribed.GetErr())
	require.Equal(t, "OK", subscribed.GetResponse().GetResponseCode("alarms").String())
	handle, err := subscribed.GetResponse().GetSubscriptionHandle("alarms")
	require.NoError(t, err)
	events := make(chan map[string]apiValues.PlcValue, 8)
	handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		events <- event.GetValue("alarms").GetStruct()
	})

	acknowledged := map[string]uint16{} // AlarmAcks the device acknowledged, by priority
	press := func() {
		// The device reads the file every 10 ms and acts on a 0 followed by a 1: hold each long enough.
		require.NoError(t, os.WriteFile(button, []byte("0\n"), 0o644))
		time.Sleep(300 * time.Millisecond)
		require.NoError(t, os.WriteFile(button, []byte("1\n"), 0o644))
		time.Sleep(300 * time.Millisecond)
	}
	expect := func(typeName string) {
		var alarm map[string]apiValues.PlcValue
		select {
		case alarm = <-events:
		case <-time.After(10 * time.Second):
			t.Fatalf("no %s alarm from the device", typeName)
		}
		priority := alarm["priority"].GetString()
		t.Logf("alarm: %s priority, %s, slot %d subslot %d, sequence %d, USI 0x%04X, data %X", priority, alarm["alarmTypeName"].GetString(),
			alarm["slot"].GetUint16(), alarm["subslot"].GetUint16(), alarm["sequence"].GetUint16(), alarm["usi"].GetUint16(), alarm["data"].GetRaw())
		require.Equal(t, typeName, alarm["alarmTypeName"].GetString())
		require.Equal(t, slot, strconv.Itoa(int(alarm["slot"].GetUint16())))
		ch := c.ar.alarmLow
		if priority == "HIGH" {
			ch = c.ar.alarmHigh
		}
		want := acknowledged[priority] & 0x7FFF
		acknowledged[priority]++
		require.Eventually(t, func() bool {
			ch.mutex.Lock()
			defer ch.mutex.Unlock()
			return !ch.waitingTack && ch.sendSeq == want
		}, 3*time.Second, 20*time.Millisecond, "the device acknowledges the AlarmAck of the %s alarm", typeName)
	}
	press()
	expect("Process")
	// The next two presses switch the device's cyclic data to backup and redundant and back, which
	// raise no alarm; the fourth adds a channel diagnosis, which does.
	press()
	press()
	time.Sleep(time.Second)
	assert.True(t, connection.IsConnected(), "the relation stays up while the device's data is backup")
	press()
	expect("Diagnosis")
	time.Sleep(time.Second)
	assert.True(t, connection.IsConnected(), "the relation stays up through the alarms")
}
