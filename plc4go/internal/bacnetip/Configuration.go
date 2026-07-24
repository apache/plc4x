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
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// Configuration captures driver-level options parsed from the connection URL.
// String/uint/bool fields are populated reflectively from driverOptions keyed by
// title-cased field name (e.g. "LocalDeviceId").
//
//go:generate go tool plc4xGenerator -type=Configuration
type Configuration struct {
	// LocalDeviceId is the BACnet device-instance number this driver claims when issuing
	// IAm and identifying itself in confirmed-request initiators. Range 0..4194303.
	LocalDeviceId uint32

	// LocalNetworkNumber is the BACnet network number for this driver's NPDU source
	// address. 0 means "same network as the destination" (i.e. local IP segment).
	LocalNetworkNumber uint16

	// MaxApduLengthAccepted is the largest APDU we will accept from peers, in octets.
	// Default 1476 (the largest fragment that fits a standard Ethernet MTU). Used to
	// derive MaxApduLengthAccepted enum constants at runtime.
	MaxApduLengthAccepted uint16

	// SegmentationSupported declares this driver's segmentation capabilities. One of
	// "segmented-both", "segmented-transmit", "segmented-receive", "no-segmentation".
	SegmentationSupported string

	// MaxSegmentsAccepted is the largest number of APDU segments we will accept in a
	// segmented response. Default 16.
	MaxSegmentsAccepted uint8

	// VendorId is the BACnet Vendor Id this driver advertises. Apache PLC4X uses 0x4D4D
	// ("MM") by convention until a real allocation is registered with ASHRAE.
	VendorId uint16

	// ForeignDeviceBBMD, when non-empty (host:port), enables BBMD foreign-device
	// registration at connect time so the driver can interoperate across IP subnets.
	ForeignDeviceBBMD string

	// ForeignDeviceTTL is the BBMD registration TTL in seconds. The driver re-registers
	// at TTL/2 to keep the foreign-device entry alive.
	ForeignDeviceTTL uint16

	// ApduTimeoutMs is the per-request timeout in milliseconds before the transaction
	// manager surfaces a timeout error.
	ApduTimeoutMs uint32

	// ApduRetries is the number of transparent retransmits attempted on confirmed
	// services before failing the request.
	ApduRetries uint8

	// CovLifetimeSeconds is the default SubscribeCOV lifetime (and refresh interval will
	// be lifetime/2). 0 means "indefinite" per BACnet spec.
	CovLifetimeSeconds uint32

	// DiscoveryTimeoutSeconds bounds the WhoIs broadcast wait time. Replaces the older
	// hardcoded 60-second wait inherited from the alpha discoverer.
	DiscoveryTimeoutSeconds uint32

	// StaticDevices lets users seed the DeviceInfoCache with peers that won't respond to
	// WhoIs (e.g. routed devices on a network with no BBMD). Comma-separated entries
	// of the form "<deviceId>@<network>:<host>:<port>".
	StaticDevices string

	// RemoteNetwork is the BACnet network number of this connection's target device
	// when it sits behind a BACnet router (ASHRAE 135 clause 6). The connection's
	// transport host is then the ROUTER's IP; every outgoing NPDU carries a
	// destination specifier (DNET=RemoteNetwork, DADR=RemoteAddress, hop count 255)
	// so the router forwards it onto that network, and replies arrive with the
	// device's source specifier mirrored back. 0 (default) means the target is on
	// the local segment and NPDUs stay specifier-free.
	RemoteNetwork uint16

	// RemoteAddress is the target device's MAC address on RemoteNetwork, required
	// when RemoteNetwork is set. For BACnet/IP-to-BACnet/IP routing use
	// "<ip>:<port>" (encoded as the 6-byte B/IP DADR); for other datalinks a hex
	// string ("0x0C") supplies the raw MAC octets.
	RemoteAddress string

	// PeerMaxApduLengthAccepted is the target device's MaxApduLengthAccepted in
	// octets, as the device declared it in its I-Am. When a confirmed request's
	// APDU would exceed this, the driver sends it as a segmented request
	// (ASHRAE 135 clause 5.4) — provided PeerSegmentationSupported allows it.
	// 0 (default) means unknown: requests are never segmented and an oversized
	// request fails fast instead of provoking an abort from the device.
	PeerMaxApduLengthAccepted uint16

	// PeerSegmentationSupported is the target device's segmentation capability
	// from its I-Am: "segmented-both", "segmented-transmit", "segmented-receive"
	// or "no-segmentation". Segmented requests are only sent when the peer can
	// RECEIVE segments ("segmented-both"/"segmented-receive"). Empty (default)
	// means unknown and is treated as "no-segmentation".
	PeerSegmentationSupported string
}

// ParseFromOptions populates a Configuration from the connection-URL query options,
// applying defaults for any field the user did not set explicitly. Field names
// are matched case-insensitively so users can write "localDeviceId", "LocalDeviceId",
// or "localdeviceid" interchangeably in the connection string.
func ParseFromOptions(log zerolog.Logger, optionsMap map[string][]string) (Configuration, error) {
	configuration := createDefaultConfiguration()
	rv := reflect.ValueOf(&configuration).Elem()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		key := field.Name
		optionValue := getFromOptions(log, optionsMap, key)
		if optionValue == "" {
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			rv.FieldByName(key).SetString(optionValue)
		case reflect.Bool:
			parsed, err := strconv.ParseBool(optionValue)
			if err != nil {
				return Configuration{}, errors.Wrapf(err, "Error parsing %s", key)
			}
			rv.FieldByName(key).SetBool(parsed)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bits := field.Type.Bits()
			parsed, err := strconv.ParseUint(optionValue, 0, bits)
			if err != nil {
				return Configuration{}, errors.Wrapf(err, "Error parsing %s", key)
			}
			rv.FieldByName(key).SetUint(parsed)
		default:
			return configuration, errors.Errorf("Configuration field kind %s not supported", field.Type.Kind())
		}
	}
	return configuration, nil
}

// Apache PLC4X uses 0x4D4D ("MM") as its placeholder BACnet Vendor Id; ASHRAE
// reserves the official allocation list.
const placeholderVendorId uint16 = 0x4D4D

func createDefaultConfiguration() Configuration {
	return Configuration{
		LocalDeviceId:           260001,
		LocalNetworkNumber:      0,
		MaxApduLengthAccepted:   1476,
		SegmentationSupported:   "segmented-both",
		MaxSegmentsAccepted:     16,
		VendorId:                placeholderVendorId,
		ApduTimeoutMs:           3000,
		ApduRetries:             3,
		CovLifetimeSeconds:      600,
		DiscoveryTimeoutSeconds: 5,
	}
}

// getFromOptions returns the first value associated with key, matching the
// optionsMap key case-insensitively (BACnet field names are CamelCase; users
// commonly type lowercase in URLs).
func getFromOptions(localLog zerolog.Logger, optionsMap map[string][]string, key string) string {
	target := strings.ToLower(key)
	for k, optionValues := range optionsMap {
		if strings.ToLower(k) != target {
			continue
		}
		if len(optionValues) == 0 {
			return ""
		}
		if len(optionValues) > 1 {
			localLog.Warn().Str("key", k).Msg("Options key must be unique")
		}
		return optionValues[0]
	}
	return ""
}
