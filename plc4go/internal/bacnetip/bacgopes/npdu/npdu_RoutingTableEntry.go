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

package npdu

import (
	"bytes"
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

type RoutingTableEntry struct {
	*DebugContents `ignore:"true"`
	rtDNET         uint16
	rtPortId       uint8
	rtPortInfo     []byte
}

func NewRoutingTableEntry(options ...Option) *RoutingTableEntry {
	r := &RoutingTableEntry{}
	r.DebugContents = NewDebugContents(r, "rtDNET", "rtPortID", "rtPortInfo")
	ApplyAppliers(options, r)
	return r
}

func (r *RoutingTableEntry) GetDebugAttr(attr string) any {
	switch attr {
	case "rtDNET":
		return r.rtDNET
	case "rtPortID":
		return r.rtPortId
	case "rtPortInfo":
		return r.rtPortInfo
	default:
		return nil
	}
}

func WithRoutingTableEntryDestinationNetworkAddress(dnet uint16) GenericApplier[*RoutingTableEntry] {
	return WrapGenericApplier(func(r *RoutingTableEntry) { r.rtDNET = dnet })
}

func WithRoutingTableEntryPortId(id uint8) GenericApplier[*RoutingTableEntry] {
	return WrapGenericApplier(func(r *RoutingTableEntry) { r.rtPortId = id })
}

func WithRoutingTableEntryPortInfo(portInfo []byte) GenericApplier[*RoutingTableEntry] {
	return WrapGenericApplier(func(r *RoutingTableEntry) { r.rtPortInfo = portInfo })
}

func (r *RoutingTableEntry) tuple() (destinationNetworkAddress uint16, portId uint8, portInfoLength uint8, portInfo []byte) {
	return r.rtDNET, r.rtPortId, uint8(len(r.rtPortInfo)), r.rtPortInfo
}

func (r *RoutingTableEntry) Equals(other any) bool {
	if r == nil && other == nil {
		return true
	}
	if r == nil {
		return false
	}
	otherEntry, ok := other.(*RoutingTableEntry)
	if !ok {
		return false
	}
	return r.rtDNET == otherEntry.rtDNET &&
		r.rtPortId == otherEntry.rtPortId &&
		bytes.Equal(r.rtPortInfo, otherEntry.rtPortInfo)
}

func (r *RoutingTableEntry) String() string {
	return fmt.Sprintf("RoutingTableEntry{rtDNET: %d, rtPortId: %d, rtPortInfo: %d}", r.rtDNET, r.rtPortId, r.rtPortInfo)
}
