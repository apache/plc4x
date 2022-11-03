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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"math"
	"sync"
	"time"
)

type DeviceInventory struct {
	sync.RWMutex
	devices map[string]DeviceEntry
}

func (d *DeviceInventory) getEntryForDestination(destination []uint8) (DeviceEntry, error) {
	d.RLock()
	defer d.RUnlock()
	deviceKey := string(destination)
	deviceEntry, ok := d.devices[deviceKey]
	if !ok {
		return NoDeviceEntry, errors.Errorf("no entry found for device key %s", deviceKey)
	}
	return deviceEntry, nil
}

var NoDeviceEntry = DeviceEntry{
	MaximumApduLengthAcceptedLength: readWriteModel.NewBACnetTagPayloadUnsignedInteger(nil, nil, nil, nil, nil, nil, nil, func() *uint64 {
		var maxUint64 uint64 = math.MaxUint64
		return &maxUint64
	}(), 4),
}

type DeviceEntry struct {
	DeviceIdentifier                readWriteModel.BACnetTagPayloadObjectIdentifier
	MaximumApduLengthAcceptedLength readWriteModel.BACnetTagPayloadUnsignedInteger
	SegmentationSupported           bool
	VendorId                        readWriteModel.BACnetVendorId
	DeviceObjects                   []DeviceObject
}

func (d DeviceEntry) GetDeviceObjects(filter ...DeviceObjectFilter) []DeviceObject {
	var deviceObjects []DeviceObject
	for _, object := range d.DeviceObjects {
		shouldBeAdded := true
		for _, objectFilter := range filter {
			shouldBeAdded = shouldBeAdded && objectFilter(object)
		}
		if shouldBeAdded {
			deviceObjects = append(deviceObjects, object)
		}
	}
	return deviceObjects
}

type DeviceObjectFilter func(DeviceObject) bool

type DeviceObject struct {
	ObjectName        string
	ObjectIdentifier  readWriteModel.BACnetTagPayloadObjectIdentifier
	CachedObjectValue interface{}
	TimeOfCache       time.Time
}
