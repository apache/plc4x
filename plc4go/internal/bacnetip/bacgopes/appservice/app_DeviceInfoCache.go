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

package appservice

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// NOTE: needs to reside here otherwise there is a circular dependency

// DeviceInfoCacheKey caches by either Instance, PduSource of both
type DeviceInfoCacheKey struct {
	Instance  *uint32
	PduSource *Address
}

func (k DeviceInfoCacheKey) HashKey() uint32 {
	h := fnv.New32a()
	if k.Instance != nil {
		_ = binary.Write(h, binary.BigEndian, *k.Instance)
	}
	_ = binary.Write(h, binary.BigEndian, k.PduSource.String())
	return h.Sum32()
}

func (k DeviceInfoCacheKey) String() string {
	return fmt.Sprintf("key: %d/%v", k.Instance, k.PduSource)
}

//go:generate go tool plc4xGenerator -type=DeviceInfoCache -prefix=app_
type DeviceInfoCache struct {
	cache map[uint32]DeviceInfo

	log zerolog.Logger
}

func NewDeviceInfoCache(localLog zerolog.Logger) *DeviceInfoCache {
	if _debug != nil {
		_debug("__init__")
	}
	return &DeviceInfoCache{
		cache: make(map[uint32]DeviceInfo),
		log:   localLog,
	}
}

// HasDeviceInfo Return true if cache has information about the device.
func (d *DeviceInfoCache) HasDeviceInfo(key DeviceInfoCacheKey) bool {
	if _debug != nil {
		_debug("has_device_info %r", key)
	}
	_, ok := d.cache[key.HashKey()]
	return ok
}

// IAmDeviceInfo Create a device information record based on the contents of an IAmRequest and put it in the cache.
func (d *DeviceInfoCache) IAmDeviceInfo(iAm readWriteModel.BACnetUnconfirmedServiceRequestIAm, pduSource Address) {
	d.log.Debug().Stringer("iAm", iAm).Msg("IAmDeviceInfo")
	if _debug != nil {
		_debug("iam_device_info %r", iAm)
	}

	deviceIdentifier := iAm.GetDeviceIdentifier()
	// Get the device instance
	deviceInstance := deviceIdentifier.GetInstanceNumber()

	// get the existing cache record if it exists
	deviceInfo, ok := d.cache[DeviceInfoCacheKey{&deviceInstance, nil}.HashKey()]

	// maybe there is a record for this address
	if !ok {
		deviceInfo, ok = d.cache[DeviceInfoCacheKey{nil, &pduSource}.HashKey()]
	}

	// make a new one using the class provided
	if !ok {
		deviceInfo = NewDeviceInfo(deviceIdentifier.GetPayload(), pduSource)
	}

	// jam in the correct values
	maximumApduLengthAccepted := readWriteModel.MaxApduLengthAccepted(iAm.GetMaximumApduLengthAcceptedLength().GetActualValue())
	deviceInfo.MaximumApduLengthAccepted = &maximumApduLengthAccepted
	segmentationSupported := iAm.GetSegmentationSupported().GetValue()
	deviceInfo.SegmentationSupported = &segmentationSupported
	vendorId := iAm.GetVendorId().GetValue()
	deviceInfo.VendorId = &vendorId

	// tell the cache this is an updated record
	d.UpdateDeviceInfo(deviceInfo)
}

// GetDeviceInfo gets a DeviceInfo from cache
func (d *DeviceInfoCache) GetDeviceInfo(key DeviceInfoCacheKey) (DeviceInfo, bool) {
	d.log.Debug().Stringer("key", key).Msg("GetDeviceInfo %s")
	if _debug != nil {
		_debug("get_device_info %r", key)
	}

	// get the info if it's there
	deviceInfo, ok := d.cache[key.HashKey()]
	d.log.Debug().Stringer("deviceInfo", &deviceInfo).Msg("deviceInfo")
	if _debug != nil {
		_debug("    - device_info: %r", deviceInfo)
	}

	return deviceInfo, ok
}

// UpdateDeviceInfo The application has updated one or more fields in the device information record and the cache needs
//
//	to be updated to reflect the changes.  If this is a cached version of a persistent record then this is the
//	opportunity to update the database.
func (d *DeviceInfoCache) UpdateDeviceInfo(deviceInfo DeviceInfo) {
	d.log.Debug().Stringer("deviceInfo", &deviceInfo).Msg("UpdateDeviceInfo")
	if _debug != nil {
		_debug("update_device_info %r", deviceInfo)
	}

	// get the current key
	cacheKey := deviceInfo._cacheKey
	if cacheKey.Instance != nil && deviceInfo.DeviceIdentifier.GetInstanceNumber() != *cacheKey.Instance {
		if _debug != nil {
			_debug("    - device identifier updated")
		}
		instanceNumber := deviceInfo.DeviceIdentifier.GetInstanceNumber()
		cacheKey.Instance = &instanceNumber
		delete(d.cache, cacheKey.HashKey())
		d.cache[DeviceInfoCacheKey{Instance: &instanceNumber}.HashKey()] = deviceInfo
	}
	if !deviceInfo.Address.Equals(cacheKey.PduSource) {
		if _debug != nil {
			_debug("    - device address updated")
		}
		cacheKey.PduSource = &deviceInfo.Address
		delete(d.cache, cacheKey.HashKey())
		d.cache[DeviceInfoCacheKey{PduSource: cacheKey.PduSource}.HashKey()] = deviceInfo
	}

	// update the key
	instanceNumber := deviceInfo.DeviceIdentifier.GetInstanceNumber()
	deviceInfo._cacheKey = DeviceInfoCacheKey{
		Instance:  &instanceNumber,
		PduSource: &deviceInfo.Address,
	}
	d.cache[deviceInfo._cacheKey.HashKey()] = deviceInfo
}

// Acquire Return the known information about the device and mark the record as being used by a segmentation state
//
//	machine.
func (d *DeviceInfoCache) Acquire(key DeviceInfoCacheKey) (DeviceInfo, bool) {
	d.log.Debug().Stringer("key", key).Msg("Acquire")
	if _debug != nil {
		_debug("acquire %r", key)
	}

	deviceInfo, ok := d.cache[key.HashKey()]
	if ok {
		if _debug != nil {
			_debug("    - reference bump")
		}
		deviceInfo._refCount++
		d.cache[key.HashKey()] = deviceInfo
	}

	if _debug != nil {
		_debug("    - device_info: %r", deviceInfo)
	}

	return deviceInfo, ok
}

// Release This function is called by the segmentation state machine when it has finished with the device information.
func (d *DeviceInfoCache) Release(deviceInfo DeviceInfo) error {
	if _debug != nil {
		_debug("release %r", deviceInfo)
	}

	//this information record might be used by more than one SSM
	if deviceInfo._refCount == 0 {
		return errors.New("reference count")
	}

	// decrement the reference count
	deviceInfo._refCount--
	d.cache[deviceInfo._cacheKey.HashKey()] = deviceInfo
	return nil
}
