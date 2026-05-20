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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestParseFromOptions_Defaults(t *testing.T) {
	got, err := ParseFromOptions(testutils.ProduceTestingLogger(t), map[string][]string{})
	require.NoError(t, err)
	assert.Equal(t, createDefaultConfiguration(), got)
}

func TestParseFromOptions_AllSet(t *testing.T) {
	opts := map[string][]string{
		"LocalDeviceId":           {"1234"},
		"LocalNetworkNumber":      {"5"},
		"MaxApduLengthAccepted":   {"480"},
		"SegmentationSupported":   {"no-segmentation"},
		"MaxSegmentsAccepted":     {"8"},
		"VendorId":                {"0xCAFE"},
		"ForeignDeviceBBMD":       {"10.0.0.5:47808"},
		"ForeignDeviceTTL":        {"30"},
		"ApduTimeoutMs":           {"5000"},
		"ApduRetries":             {"2"},
		"CovLifetimeSeconds":      {"120"},
		"DiscoveryTimeoutSeconds": {"2"},
		"StaticDevices":           {"42@1:10.0.0.7:47808"},
	}
	got, err := ParseFromOptions(testutils.ProduceTestingLogger(t), opts)
	require.NoError(t, err)

	want := createDefaultConfiguration()
	want.LocalDeviceId = 1234
	want.LocalNetworkNumber = 5
	want.MaxApduLengthAccepted = 480
	want.SegmentationSupported = "no-segmentation"
	want.MaxSegmentsAccepted = 8
	want.VendorId = 0xCAFE
	want.ForeignDeviceBBMD = "10.0.0.5:47808"
	want.ForeignDeviceTTL = 30
	want.ApduTimeoutMs = 5000
	want.ApduRetries = 2
	want.CovLifetimeSeconds = 120
	want.DiscoveryTimeoutSeconds = 2
	want.StaticDevices = "42@1:10.0.0.7:47808"
	assert.Equal(t, want, got)
}

func TestParseFromOptions_CaseInsensitive(t *testing.T) {
	got, err := ParseFromOptions(testutils.ProduceTestingLogger(t), map[string][]string{
		"localdeviceid": {"42"},
	})
	require.NoError(t, err)
	assert.Equal(t, uint32(42), got.LocalDeviceId)
}

func TestParseFromOptions_InvalidBool(t *testing.T) {
	// No bool fields today, but the dispatch path matters: a non-numeric value on a
	// numeric field should error rather than silently zeroing.
	_, err := ParseFromOptions(testutils.ProduceTestingLogger(t), map[string][]string{
		"LocalDeviceId": {"not-a-number"},
	})
	assert.Error(t, err)
}

func TestParseFromOptions_InvalidUint(t *testing.T) {
	// Overflow uint8 ApduRetries.
	_, err := ParseFromOptions(testutils.ProduceTestingLogger(t), map[string][]string{
		"ApduRetries": {"1000"},
	})
	assert.Error(t, err)
}

func TestCreateDefaultConfiguration(t *testing.T) {
	cfg := createDefaultConfiguration()
	assert.Equal(t, uint32(260001), cfg.LocalDeviceId)
	assert.Equal(t, uint16(1476), cfg.MaxApduLengthAccepted)
	assert.Equal(t, "segmented-both", cfg.SegmentationSupported)
	assert.Equal(t, uint8(16), cfg.MaxSegmentsAccepted)
	assert.Equal(t, uint16(0x4D4D), cfg.VendorId)
	assert.Equal(t, uint32(3000), cfg.ApduTimeoutMs)
	assert.Equal(t, uint8(3), cfg.ApduRetries)
	assert.Equal(t, uint32(600), cfg.CovLifetimeSeconds)
	assert.Equal(t, uint32(5), cfg.DiscoveryTimeoutSeconds)
}

func TestGetFromOptions(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	assert.Empty(t, getFromOptions(log, map[string][]string{}, "missing"))
	assert.Equal(t, "first", getFromOptions(log, map[string][]string{"present": {"first", "second"}}, "present"))
}
