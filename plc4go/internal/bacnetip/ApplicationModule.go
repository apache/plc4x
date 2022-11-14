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
	"bytes"
	"encoding/binary"
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"hash/fnv"
)

type DeviceInfo struct {
	DeviceIdentifier readWriteModel.BACnetTagPayloadObjectIdentifier
	Address          []byte

	MaximumApduLengthAccepted *readWriteModel.MaxApduLengthAccepted
	SegmentationSupported     *readWriteModel.BACnetSegmentation
	MaxSegmentsAccepted       *readWriteModel.MaxSegmentsAccepted
	VendorId                  *readWriteModel.BACnetVendorId
	MaximumNpduLength         *uint

	_refCount int
	_cacheKey DeviceInfoCacheKey
}

func NewDeviceInfo(deviceIdentifier readWriteModel.BACnetTagPayloadObjectIdentifier, address []byte) *DeviceInfo {
	return &DeviceInfo{
		DeviceIdentifier: deviceIdentifier,
		Address:          address,

		MaximumApduLengthAccepted: func() *readWriteModel.MaxApduLengthAccepted {
			octets1024 := readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1024
			return &octets1024
		}(),
		SegmentationSupported: func() *readWriteModel.BACnetSegmentation {
			noSegmentation := readWriteModel.BACnetSegmentation_NO_SEGMENTATION
			return &noSegmentation
		}(),
	}
}

// DeviceInfoCacheKey caches by either Instance, PduSource of both
type DeviceInfoCacheKey struct {
	Instance  *uint32
	PduSource []byte
}

func (k DeviceInfoCacheKey) HashKey() uint32 {
	h := fnv.New32a()
	if k.Instance != nil {
		_ = binary.Write(h, binary.BigEndian, *k.Instance)
	}
	_, _ = h.Write(k.PduSource)
	return h.Sum32()
}

func (k DeviceInfoCacheKey) String() string {
	return fmt.Sprintf("key: %d/%x", k.Instance, k.PduSource)
}

type DeviceInfoCache struct {
	cache map[uint32]DeviceInfo
}

func NewDeviceInfoCache() *DeviceInfoCache {
	return &DeviceInfoCache{
		cache: make(map[uint32]DeviceInfo),
	}
}

// HasDeviceInfo Return true if cache has information about the device.
func (i *DeviceInfoCache) HasDeviceInfo(key DeviceInfoCacheKey) bool {
	_, ok := i.cache[key.HashKey()]
	return ok
}

// IAmDeviceInfo Create a device information record based on the contents of an IAmRequest and put it in the cache.
func (i *DeviceInfoCache) IAmDeviceInfo(iAm readWriteModel.BACnetUnconfirmedServiceRequestIAm, pduSource []byte) {
	log.Debug().Msgf("IAmDeviceInfo\n%s", iAm)

	deviceIdentifier := iAm.GetDeviceIdentifier()
	// Get the device instance
	deviceInstance := deviceIdentifier.GetInstanceNumber()

	// get the existing cache record if it exists
	deviceInfo, ok := i.cache[DeviceInfoCacheKey{&deviceInstance, nil}.HashKey()]

	// maybe there is a record for this address
	if !ok {
		deviceInfo, ok = i.cache[DeviceInfoCacheKey{nil, pduSource}.HashKey()]
	}

	// make a new one using the class provided
	if !ok {
		deviceInfo = DeviceInfo{
			DeviceIdentifier: deviceIdentifier.GetPayload(),
			Address:          pduSource,
		}
	}

	// jam in the correct values
	maximumApduLengthAccepted := readWriteModel.MaxApduLengthAccepted(iAm.GetMaximumApduLengthAcceptedLength().GetActualValue())
	deviceInfo.MaximumApduLengthAccepted = &maximumApduLengthAccepted
	sementationSupported := iAm.GetSegmentationSupported().GetValue()
	deviceInfo.SegmentationSupported = &sementationSupported
	vendorId := iAm.GetVendorId().GetValue()
	deviceInfo.VendorId = &vendorId

	// tell the cache this is an updated record
	i.UpdateDeviceInfo(deviceInfo)
}

// GetDeviceInfo gets a DeviceInfo from cache
func (i *DeviceInfoCache) GetDeviceInfo(key DeviceInfoCacheKey) (DeviceInfo, bool) {
	log.Debug().Msgf("GetDeviceInfo %s", key)

	// get the info if it's there
	deviceInfo, ok := i.cache[key.HashKey()]
	log.Debug().Msgf("deviceInfo: %#v", deviceInfo)

	return deviceInfo, ok
}

// UpdateDeviceInfo The application has updated one or more fields in the device information record and the cache needs
//        to be updated to reflect the changes.  If this is a cached version of a persistent record then this is the
//        opportunity to update the database.
func (i *DeviceInfoCache) UpdateDeviceInfo(deviceInfo DeviceInfo) {
	log.Debug().Msgf("UpdateDeviceInfo %#v", deviceInfo)

	// get the current key
	cacheKey := deviceInfo._cacheKey
	if cacheKey.Instance != nil && deviceInfo.DeviceIdentifier.GetInstanceNumber() != *cacheKey.Instance {
		instanceNumber := deviceInfo.DeviceIdentifier.GetInstanceNumber()
		cacheKey.Instance = &instanceNumber
		delete(i.cache, cacheKey.HashKey())
		i.cache[DeviceInfoCacheKey{Instance: &instanceNumber}.HashKey()] = deviceInfo
	}
	if bytes.Compare(deviceInfo.Address, cacheKey.PduSource) != 0 {
		cacheKey.PduSource = deviceInfo.Address
		delete(i.cache, cacheKey.HashKey())
		i.cache[DeviceInfoCacheKey{PduSource: cacheKey.PduSource}.HashKey()] = deviceInfo
	}

	// update the key
	instanceNumber := deviceInfo.DeviceIdentifier.GetInstanceNumber()
	deviceInfo._cacheKey = DeviceInfoCacheKey{
		Instance:  &instanceNumber,
		PduSource: deviceInfo.Address,
	}
	i.cache[deviceInfo._cacheKey.HashKey()] = deviceInfo
}

// Acquire Return the known information about the device and mark the record as being used by a segmentation state
//        machine.
func (i *DeviceInfoCache) Acquire(key DeviceInfoCacheKey) (DeviceInfo, bool) {
	log.Debug().Msgf("Acquire %#v", key)

	deviceInfo, ok := i.cache[key.HashKey()]
	if ok {
		deviceInfo._refCount++
		i.cache[key.HashKey()] = deviceInfo
	}

	return deviceInfo, ok
}

// Release This function is called by the segmentation state machine when it has finished with the device information.
func (i *DeviceInfoCache) Release(deviceInfo DeviceInfo) error {

	//this information record might be used by more than one SSM
	if deviceInfo._refCount == 0 {
		return errors.New("reference count")
	}

	// decrement the reference count
	deviceInfo._refCount--
	i.cache[deviceInfo._cacheKey.HashKey()] = deviceInfo
	return nil
}

// TODO: implement
type Application struct {
	ApplicationServiceElement
	Collector
}

// TODO: implement
type IOController struct {
}

// TODO: implement
type ApplicationIOController struct {
	IOController
	Application
}

func NewApplicationIOController(interface{}, interface{}, interface{}, *int) (*ApplicationIOController, error) {
	return &ApplicationIOController{}, nil
}

type BIPSimpleApplication struct {
	*ApplicationIOController
	*WhoIsIAmServices
	*ReadWritePropertyServices
	localAddress interface{}
	asap         *ApplicationServiceAccessPoint
	smap         *StateMachineAccessPoint
	nsap         *NetworkServiceAccessPoint
	nse          *NetworkServiceElement
	bip          *BIPSimple
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer
}

func NewBIPSimpleApplication(localDevice LocalDeviceObject, localAddress, deviceInfoCache *DeviceInfoCache, aseID *int) (*BIPSimpleApplication, error) {
	b := &BIPSimpleApplication{}
	var err error
	b.ApplicationIOController, err = NewApplicationIOController(localDevice, localAddress, deviceInfoCache, aseID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}

	b.localAddress = localAddress

	// include a application decoder
	b.asap, err = NewApplicationServiceAccessPoint(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service access point")
	}

	// pass the device object to the state machine access point, so it can know if it should support segmentation
	b.smap, err = NewStateMachineAccessPoint(localDevice, deviceInfoCache, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating state machine access point")
	}

	// pass the device object to the state machine access point so it # can know if it should support segmentation
	// Note: deviceInfoCache already passed above, so we don't need to do it again here

	// a network service access point will be needed
	b.nsap, err = NewNetworkServiceAccessPoint()
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = NewNetworkServiceElement()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}
	if err := bind(b.nse, b.nsap); err != nil {
		return nil, errors.New("error binding network stack")
	}

	// bind the top layers
	if err := bind(b, b.asap, b.smap, b.nsap); err != nil {
		return nil, errors.New("error binding top layers")
	}

	// create a generic BIP stack, bound to the Annex J server on the UDP multiplexer
	b.bip, err = NewBIPSimple(nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new bip")
	}
	b.annexj, err = NewAnnexJCodec(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	b.mux, err = NewUDPMultiplexer(b.localAddress, false)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new udp multiplexer")
	}

	// bind the bottom layers
	if err := bind(b.bip, b.annexj, b.mux.annexJ); err != nil {
		return nil, errors.New("error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := b.nsap.bind(b.bip, nil, b.localAddress); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BIPSimpleApplication) Close() error {
	log.Debug().Msg("close socket")
	// pass to the multiplexer, then down to the sockets
	return b.mux.Close()
}
