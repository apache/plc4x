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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// routedDestination is the NPDU destination specifier for a connection whose
// target device sits behind a BACnet router (Configuration.RemoteNetwork /
// RemoteAddress): DNET is the device's network number, DADR its MAC on that
// network. nil means local-segment addressing (no specifier), the default.
type routedDestination struct {
	dnet uint16
	dadr []uint8
}

// routedDestinationHopCount is the initial hop count of freshly originated
// routed NPDUs (ASHRAE 135 clause 6.2.2 recommends starting at 255).
const routedDestinationHopCount uint8 = 255

// parseRemoteAddress turns Configuration.RemoteAddress into DADR octets:
// "<ip>:<port>" becomes the 6-byte B/IP address (Annex J), and a "0x…" hex
// string supplies raw MAC octets for non-IP datalinks.
func parseRemoteAddress(remoteAddress string) ([]uint8, error) {
	if s := strings.TrimPrefix(remoteAddress, "0x"); s != remoteAddress {
		dadr, err := hex.DecodeString(s)
		if err != nil {
			return nil, errors.Wrapf(err, "RemoteAddress %q is not valid hex", remoteAddress)
		}
		if len(dadr) == 0 {
			return nil, errors.Errorf("RemoteAddress %q decodes to zero octets", remoteAddress)
		}
		return dadr, nil
	}
	ap, err := netip.ParseAddrPort(remoteAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "RemoteAddress %q is neither <ip>:<port> nor 0x-prefixed hex", remoteAddress)
	}
	if !ap.Addr().Is4() {
		return nil, errors.Errorf("RemoteAddress %q must be IPv4 for a B/IP DADR", remoteAddress)
	}
	ip4 := ap.Addr().As4()
	dadr := make([]uint8, 6)
	copy(dadr[0:4], ip4[:])
	binary.BigEndian.PutUint16(dadr[4:6], ap.Port())
	return dadr, nil
}

// routedDestinationFromConfiguration derives the connection-scoped destination
// specifier, or nil for local-segment connections. RemoteAddress without
// RemoteNetwork (and vice versa) is a configuration error.
func routedDestinationFromConfiguration(configuration Configuration) (*routedDestination, error) {
	if configuration.RemoteNetwork == 0 && configuration.RemoteAddress == "" {
		return nil, nil
	}
	if configuration.RemoteNetwork == 0 || configuration.RemoteAddress == "" {
		return nil, errors.New("RemoteNetwork and RemoteAddress must be set together")
	}
	dadr, err := parseRemoteAddress(configuration.RemoteAddress)
	if err != nil {
		return nil, err
	}
	return &routedDestination{dnet: configuration.RemoteNetwork, dadr: dadr}, nil
}

// wrapAPDU encapsulates an APDU in the BACnet/IP NPDU + BVLC layers expected
// by MessageCodec.Send (which type-asserts to model.BVLC). Without this
// wrapping the codec panics on the cast and the request silently dies.
//
// expectingReply is set for confirmed requests; unconfirmed requests pass
// false. With dest == nil the NPDU is local-only (no DNET/SNET); a non-nil
// dest emits the destination specifier (DNET/DLEN/DADR + hop count) so the
// first router on the connection's segment forwards the request onto the
// target's network (ASHRAE 135 clause 6).
func wrapAPDU(apdu model.APDU, expectingReply bool, dest *routedDestination) model.BVLC {
	control := model.NewNPDUControl(
		false,       // messageTypeFieldPresent
		dest != nil, // destinationSpecified
		false,       // sourceSpecified
		expectingReply,
		model.NPDUNetworkPriority_NORMAL_MESSAGE,
	)
	var destNet *uint16
	var destLen *uint8
	var destAddr []uint8
	var hopCount *uint8
	if dest != nil {
		dnet := dest.dnet
		destNet = &dnet
		dlen := uint8(len(dest.dadr))
		destLen = &dlen
		destAddr = dest.dadr
		hops := routedDestinationHopCount
		hopCount = &hops
	}
	npdu := model.NewNPDU(
		1, // protocolVersionNumber
		control,
		destNet,
		destLen,
		destAddr,
		nil, // sourceNetworkAddress
		nil, // sourceLength
		nil, // sourceAddress
		hopCount,
		nil, // nlm
		apdu,
	)
	return model.NewBVLCOriginalUnicastNPDU(npdu)
}

// routedOriginOptions extracts a routed frame's NPDU source specifier as
// connection options (RemoteNetwork/RemoteAddress — the keys the connection
// URL accepts), or nil for local frames. A 6-octet SADR is rendered as the
// B/IP "<ip>:<port>" form; other datalink MACs as 0x-prefixed hex.
func routedOriginOptions(npdu model.NPDU) map[string][]string {
	if npdu == nil || npdu.GetControl() == nil || !npdu.GetControl().GetSourceSpecified() {
		return nil
	}
	snet := npdu.GetSourceNetworkAddress()
	sadr := npdu.GetSourceAddress()
	if snet == nil || len(sadr) == 0 {
		return nil
	}
	var remoteAddress string
	if len(sadr) == 6 {
		remoteAddress = fmt.Sprintf("%d.%d.%d.%d:%d", sadr[0], sadr[1], sadr[2], sadr[3], binary.BigEndian.Uint16(sadr[4:6]))
	} else {
		remoteAddress = "0x" + hex.EncodeToString(sadr)
	}
	return map[string][]string{
		"RemoteNetwork": {strconv.FormatUint(uint64(*snet), 10)},
		"RemoteAddress": {remoteAddress},
	}
}
