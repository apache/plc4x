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
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// DeviceInfo captures everything we need to know to address a remote BACnet
// device for routed Read/Write/Subscribe requests. Most fields are populated
// from IAm responses (during discovery or in response to WhoIs); the seed
// entries from Configuration.StaticDevices supply Address + SourceNetwork only.
type DeviceInfo struct {
	// DeviceId is the BACnet device-instance number (0..4194303).
	DeviceId uint32

	// Address is the UDP endpoint reachable on the local IP segment that hosts
	// (or routes to) this device.
	Address *net.UDPAddr

	// SourceNetwork is the BACnet network number the device is on. 0 means
	// "local network" (no routing needed).
	SourceNetwork uint16

	// SourceAddress is the per-network MAC-equivalent address used in the NPDU
	// DADR field when routing across networks. For BACnet/IP devices on a
	// remote network it's typically the device's IPv4 address + port (6 bytes).
	SourceAddress []byte

	// SegmentationSupported is what the device declared in its IAm. We honor
	// this when deciding whether to send a segmented request to it.
	SegmentationSupported model.BACnetSegmentation

	// MaxApdu is the device's MaxApduLengthAccepted (already converted to bytes).
	MaxApdu uint16

	// VendorId is the BACnet Vendor Id from IAm.
	VendorId uint16

	// LastSeen is set whenever we receive any APDU from this device, used for
	// staleness pruning if/when we add a TTL.
	LastSeen time.Time
}

// DeviceInfoCache is the in-memory routing table. Safe for concurrent use.
type DeviceInfoCache struct {
	mu      sync.RWMutex
	entries map[uint32]*DeviceInfo
}

func NewDeviceInfoCache() *DeviceInfoCache {
	return &DeviceInfoCache{entries: make(map[uint32]*DeviceInfo)}
}

// Get returns a snapshot of the DeviceInfo for the given device id, or nil if
// the device isn't known. The returned pointer is a copy — callers may not
// mutate the cache through it.
func (c *DeviceInfoCache) Get(deviceId uint32) *DeviceInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry, ok := c.entries[deviceId]; ok {
		copy := *entry
		return &copy
	}
	return nil
}

// Put inserts or overwrites a device-info entry. LastSeen is bumped to now.
func (c *DeviceInfoCache) Put(info DeviceInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	info.LastSeen = time.Now()
	c.entries[info.DeviceId] = &info
}

// Touch updates only the LastSeen timestamp for the given device id. Used by
// the message-receive path to refresh entries we've already seen.
func (c *DeviceInfoCache) Touch(deviceId uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, ok := c.entries[deviceId]; ok {
		entry.LastSeen = time.Now()
	}
}

// All returns a copy of every entry. Useful for diagnostics; tests use it to
// assert eventual-consistency without exposing the internal map.
func (c *DeviceInfoCache) All() []DeviceInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]DeviceInfo, 0, len(c.entries))
	for _, e := range c.entries {
		out = append(out, *e)
	}
	return out
}

// LoadStatic seeds the cache from a Configuration.StaticDevices string. The
// format is comma-separated entries of the form "<deviceId>@<network>:<host>:<port>".
// Whitespace around separators is tolerated.
//
//	"1234@0:10.0.0.5:47808,42@1:10.0.0.6:47808"
//
// Entries that fail to parse are skipped and counted in the returned error;
// the cache always contains every entry that did parse cleanly.
func (c *DeviceInfoCache) LoadStatic(spec string) error {
	if spec == "" {
		return nil
	}
	var failures []string
	for _, raw := range strings.Split(spec, ",") {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		info, err := parseStaticDevice(entry)
		if err != nil {
			failures = append(failures, entry+": "+err.Error())
			continue
		}
		c.Put(info)
	}
	if len(failures) > 0 {
		return errors.Errorf("parsing static devices: %s", strings.Join(failures, "; "))
	}
	return nil
}

// parseStaticDevice parses a single entry of the form
// "<deviceId>@<network>:<host>:<port>".
func parseStaticDevice(entry string) (DeviceInfo, error) {
	atIdx := strings.IndexByte(entry, '@')
	if atIdx <= 0 || atIdx == len(entry)-1 {
		return DeviceInfo{}, errors.New(`expected "<deviceId>@<network>:<host>:<port>"`)
	}
	deviceIdPart := entry[:atIdx]
	remainder := entry[atIdx+1:]

	deviceId, err := strconv.ParseUint(deviceIdPart, 10, 32)
	if err != nil {
		return DeviceInfo{}, errors.Wrap(err, "deviceId")
	}

	netSep := strings.IndexByte(remainder, ':')
	if netSep <= 0 {
		return DeviceInfo{}, errors.New("missing ':' after network number")
	}
	networkPart := remainder[:netSep]
	hostPort := remainder[netSep+1:]
	network, err := strconv.ParseUint(networkPart, 10, 16)
	if err != nil {
		return DeviceInfo{}, errors.Wrap(err, "network")
	}
	addr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		return DeviceInfo{}, errors.Wrap(err, "host:port")
	}
	return DeviceInfo{
		DeviceId:      uint32(deviceId),
		SourceNetwork: uint16(network),
		Address:       addr,
	}, nil
}
