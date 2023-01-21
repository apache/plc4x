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
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"hash/fnv"
)

type DeviceInfo struct {
	DeviceIdentifier readWriteModel.BACnetTagPayloadObjectIdentifier
	Address          Address

	MaximumApduLengthAccepted *readWriteModel.MaxApduLengthAccepted
	SegmentationSupported     *readWriteModel.BACnetSegmentation
	MaxSegmentsAccepted       *readWriteModel.MaxSegmentsAccepted
	VendorId                  *readWriteModel.BACnetVendorId
	MaximumNpduLength         *uint

	_refCount int
	_cacheKey DeviceInfoCacheKey
}

func NewDeviceInfo(deviceIdentifier readWriteModel.BACnetTagPayloadObjectIdentifier, address Address) DeviceInfo {
	return DeviceInfo{
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

type DeviceInfoCache struct {
	cache map[uint32]DeviceInfo
}

func NewDeviceInfoCache() *DeviceInfoCache {
	return &DeviceInfoCache{
		cache: make(map[uint32]DeviceInfo),
	}
}

func (d *DeviceInfoCache) String() string {
	return fmt.Sprintf("DeviceInfoCache(%d)", len(d.cache))
}

// HasDeviceInfo Return true if cache has information about the device.
func (d *DeviceInfoCache) HasDeviceInfo(key DeviceInfoCacheKey) bool {
	_, ok := d.cache[key.HashKey()]
	return ok
}

// IAmDeviceInfo Create a device information record based on the contents of an IAmRequest and put it in the cache.
func (d *DeviceInfoCache) IAmDeviceInfo(iAm readWriteModel.BACnetUnconfirmedServiceRequestIAm, pduSource Address) {
	log.Debug().Msgf("IAmDeviceInfo\n%s", iAm)

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
	log.Debug().Msgf("GetDeviceInfo %s", key)

	// get the info if it's there
	deviceInfo, ok := d.cache[key.HashKey()]
	log.Debug().Msgf("deviceInfo: %#v", deviceInfo)

	return deviceInfo, ok
}

// UpdateDeviceInfo The application has updated one or more fields in the device information record and the cache needs
//        to be updated to reflect the changes.  If this is a cached version of a persistent record then this is the
//        opportunity to update the database.
func (d *DeviceInfoCache) UpdateDeviceInfo(deviceInfo DeviceInfo) {
	log.Debug().Msgf("UpdateDeviceInfo %#v", deviceInfo)

	// get the current key
	cacheKey := deviceInfo._cacheKey
	if cacheKey.Instance != nil && deviceInfo.DeviceIdentifier.GetInstanceNumber() != *cacheKey.Instance {
		instanceNumber := deviceInfo.DeviceIdentifier.GetInstanceNumber()
		cacheKey.Instance = &instanceNumber
		delete(d.cache, cacheKey.HashKey())
		d.cache[DeviceInfoCacheKey{Instance: &instanceNumber}.HashKey()] = deviceInfo
	}
	if !deviceInfo.Address.Equals(cacheKey.PduSource) {
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
//        machine.
func (d *DeviceInfoCache) Acquire(key DeviceInfoCacheKey) (DeviceInfo, bool) {
	log.Debug().Msgf("Acquire %#v", key)

	deviceInfo, ok := d.cache[key.HashKey()]
	if ok {
		deviceInfo._refCount++
		d.cache[key.HashKey()] = deviceInfo
	}

	return deviceInfo, ok
}

// Release This function is called by the segmentation state machine when it has finished with the device information.
func (d *DeviceInfoCache) Release(deviceInfo DeviceInfo) error {

	//this information record might be used by more than one SSM
	if deviceInfo._refCount == 0 {
		return errors.New("reference count")
	}

	// decrement the reference count
	deviceInfo._refCount--
	d.cache[deviceInfo._cacheKey.HashKey()] = deviceInfo
	return nil
}

type Application struct {
	*ApplicationServiceElement
	Collector

	objectName       map[string]*LocalDeviceObject
	objectIdentifier map[string]*LocalDeviceObject
	localDevice      *LocalDeviceObject
	deviceInfoCache  *DeviceInfoCache
	controllers      map[string]interface{}
	helpers          map[string]func(pdu _PDU) error

	_startupDisabled bool
}

func NewApplication(localDevice *LocalDeviceObject, deviceInfoCache *DeviceInfoCache, aseID *int) (*Application, error) {
	log.Debug().Msgf("NewApplication localDevice=%v deviceInfoCache=%s aseID=%d", localDevice, deviceInfoCache, aseID)
	a := &Application{}
	var err error
	a.ApplicationServiceElement, err = NewApplicationServiceElement(aseID, a)
	if err != nil {
		return nil, err
	}

	// local objects by ID and name
	a.objectName = map[string]*LocalDeviceObject{}
	a.objectIdentifier = map[string]*LocalDeviceObject{}

	// keep track of the local device
	if localDevice != nil {
		a.localDevice = localDevice

		// bind the device object to this application
		localDevice.App = a

		// local objects by ID and name
		a.objectName[localDevice.ObjectName] = localDevice
		a.objectName[localDevice.ObjectIdentifier] = localDevice
	}

	// use the provided cache or make a default one
	if deviceInfoCache == nil {
		deviceInfoCache = NewDeviceInfoCache()
	}
	a.deviceInfoCache = deviceInfoCache

	// controllers for managing confirmed requests as a client
	a.controllers = map[string]interface{}{}

	// now set up the rest of the capabilities
	a.Collector = Collector{}

	// if starting up is enabled, find all the startup functions
	if !a._startupDisabled {
		for _, fn := range a.CapabilityFunctions("startup") {
			log.Debug().Msgf("startup fn %t", fn != nil)
			Deferred(fn)
		}
	}
	return a, nil
}

// AddObject adds an object to the local collection
func (a *Application) AddObject(obj *LocalDeviceObject) error {
	log.Debug().Msgf("AddObject %v", obj)

	// extract the object name and identifier
	objectName := obj.ObjectName
	if objectName == "" {
		return errors.New("object name required")
	}
	objectIdentifier := obj.ObjectIdentifier
	if objectIdentifier == "" {
		return errors.New("object identifier required")
	}

	// make sure it hasn't already been defined
	if _, ok := a.objectName[objectName]; ok {
		return errors.Errorf("already an object with name %s", objectName)
	}
	if _, ok := a.objectIdentifier[objectIdentifier]; ok {
		return errors.Errorf("already an object with identifier %s", objectIdentifier)
	}

	// now put it in local dictionaries
	a.objectName[objectName] = obj
	a.objectIdentifier[objectIdentifier] = obj

	// append the new object's identifier to the local device's object list if there is one and has an object list property
	if a.localDevice != nil {
		a.localDevice.ObjectList = append(a.localDevice.ObjectList, objectIdentifier)
	}

	// let the object know which application stack it belongs to
	obj.App = a

	return nil
}

// DeleteObject deletes an object from the local collection
func (a *Application) DeleteObject(obj *LocalDeviceObject) error {
	log.Debug().Msgf("DeleteObject %v", obj)

	// extract the object name and identifier
	objectName := obj.ObjectName
	objectIdentifier := obj.ObjectIdentifier

	// delete it from the application
	delete(a.objectName, objectName)
	delete(a.objectIdentifier, objectIdentifier)

	// remove the object's identifier from the device's object list if there is one and has an object list property
	if a.localDevice != nil {
		foundIndex := -1
		for i, s := range a.localDevice.ObjectList {
			if s == objectIdentifier {
				foundIndex = i
			}
		}
		if foundIndex >= 0 {
			a.localDevice.ObjectList = append(a.localDevice.ObjectList[0:foundIndex], a.localDevice.ObjectList[foundIndex+1:]...)
		}
	}

	// make sure the object knows it's detached from an application
	obj.App = nil

	return nil
}

// GetObjectId returns a local object or None.
func (a *Application) GetObjectId(objectId string) *LocalDeviceObject {
	return a.objectIdentifier[objectId]
}

// GetObjectName returns a local object or None.
func (a *Application) GetObjectName(objectName string) *LocalDeviceObject {
	return a.objectName[objectName]
}

// IterObjects iterates over the objects
func (a *Application) IterObjects() []*LocalDeviceObject {
	localDeviceObjects := make([]*LocalDeviceObject, 0, len(a.objectIdentifier))
	for _, object := range a.objectIdentifier {
		localDeviceObjects = append(localDeviceObjects, object)
	}
	return localDeviceObjects
}

// GetServicesSupported returns a ServicesSupported bit string based in introspection, look for helper methods that
//        match confirmed and unconfirmed services.
// TODO: match that with readWriteModel.BACnetServicesSupported
func (a *Application) GetServicesSupported() []string {
	servicesSupported := make([]string, 0, len(a.helpers))
	for key, _ := range a.helpers {
		servicesSupported = append(servicesSupported, key)
	}
	return servicesSupported
}

func (a *Application) Request(apdu _PDU) error {
	log.Debug().Msgf("Request\n%s", apdu)

	// double-check the input is the right kind of APDU
	switch apdu.GetMessage().(type) {
	case readWriteModel.APDUUnconfirmedRequestExactly, readWriteModel.APDUConfirmedRequestExactly:
	default:
		return errors.New("APDU expected")
	}
	return a.ApplicationServiceElement.Request(apdu)
}

func (a *Application) Indication(apdu _PDU) error {
	log.Debug().Msgf("Indication\n%s", apdu)

	// get a helper function
	helperName := fmt.Sprintf("Do_%T", apdu)
	helperFn := a.helpers[helperName]
	log.Debug().Msgf("helperFn: %s == %t", helperName, helperFn != nil)

	// send back a reject for unrecognized services
	if helperFn == nil {
		if _, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); ok {
			return errors.Errorf("no function %s", helperName)
		}
		return nil
	}

	if err := helperFn(apdu); err != nil {
		log.Debug().Err(err).Msgf("err result")
		// TODO: do proper mapping
		if err := a.Response(NewPDU(readWriteModel.NewAPDUError(0, readWriteModel.BACnetConfirmedServiceChoice_CREATE_OBJECT, nil, 0))); err != nil {
			return err
		}
	}

	return nil
}

// TODO: finish
type ApplicationIOController struct {
	*IOController
	*Application
	queueByAddress map[string]SieveQueue
}

func NewApplicationIOController(localDevice *LocalDeviceObject, deviceInfoCache *DeviceInfoCache, aseID *int) (*ApplicationIOController, error) {
	a := &ApplicationIOController{
		// queues for each address
		queueByAddress: make(map[string]SieveQueue),
	}
	var err error
	a.IOController, err = NewIOController("", a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	a.Application, err = NewApplication(localDevice, deviceInfoCache, aseID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application")
	}
	return a, nil
}

func (a *ApplicationIOController) ProcessIO(iocb _IOCB) error {
	log.Debug().Msgf("ProcessIO %s", iocb)

	// get the destination address from the pdu
	destinationAddress := iocb.getDestination()
	log.Debug().Msgf("destinationAddress %s", destinationAddress)

	// look up the queue
	queue, ok := a.queueByAddress[destinationAddress.String()]
	if !ok {
		newQueue, _ := NewSieveQueue(a._AppRequest, destinationAddress)
		queue = *newQueue
		a.queueByAddress[destinationAddress.String()] = queue
	}
	log.Debug().Msgf("queue %v", queue)

	// ask the queue to process the request
	return queue.RequestIO(iocb)
}

func (a *ApplicationIOController) _AppComplete(address *Address, apdu _PDU) error {
	log.Debug().Msgf("_AppComplete %s\n%s", address, apdu)

	// look up the queue
	queue, ok := a.queueByAddress[address.String()]
	if !ok {
		log.Debug().Msgf("no queue for %s", address)
		return nil
	}
	log.Debug().Msgf("queue %v", queue)

	// make sure it has an active iocb
	if queue.activeIOCB == nil {
		log.Debug().Msgf("no active request for %s", address)
		return nil
	}

	// this request is complete
	switch apdu.GetMessage().(type) {
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUComplexAckExactly:
		if err := queue.CompleteIO(queue.activeIOCB, apdu); err != nil {
			return err
		}
	case readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly, readWriteModel.APDUAbortExactly:
		// TODO: extract error
		if err := queue.AbortIO(queue.activeIOCB, errors.Errorf("%s", apdu)); err != nil {
			return err
		}
	default:
		return errors.New("unrecognized APDU type")
	}
	log.Debug().Msg("controller finished")
	// if the queue is empty and idle, forget about the controller
	if len(queue.ioQueue.queue) == 0 && queue.activeIOCB == nil {
		delete(a.queueByAddress, address.String())
	}
	return nil
}

func (a *ApplicationIOController) _AppRequest(apdu _PDU) {
	log.Debug().Msgf("_AppRequest\n%s", apdu)

	// send it downstream, bypass the guard
	if err := a.Application.Request(apdu); err != nil {
		log.Error().Stack().Err(err).Msg("Uh oh")
		return
	}

	// if this was an unconfirmed request, it's complete, no message
	if _, ok := apdu.(readWriteModel.APDUUnconfirmedRequestExactly); ok {
		if err := a._AppComplete(apdu.GetPDUDestination(), apdu); err != nil {
			log.Error().Err(err).Msg("AppRequest failed")
			return
		}
	}
}

func (a *ApplicationIOController) Request(apdu _PDU) error {
	log.Debug().Msgf("Request\n%s", apdu)

	// if this is not unconfirmed request, tell the application to use the IOCB interface
	if _, ok := apdu.(readWriteModel.APDUUnconfirmedRequestExactly); !ok {
		return errors.New("use IOCB for confirmed requests")
	}

	// send it downstream
	return a.Application.Request(apdu)
}

func (a *ApplicationIOController) Confirmation(apdu _PDU) error {
	log.Debug().Msgf("Confirmation\n%s", apdu)

	// this is an ack, error, reject or abort
	return a._AppComplete(apdu.GetPDUSource(), apdu)
}

type BIPSimpleApplication struct {
	*ApplicationIOController
	*WhoIsIAmServices
	*ReadWritePropertyServices
	localAddress Address
	asap         *ApplicationServiceAccessPoint
	smap         *StateMachineAccessPoint
	nsap         *NetworkServiceAccessPoint
	nse          *NetworkServiceElement
	bip          *BIPSimple
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer
}

func NewBIPSimpleApplication(localDevice *LocalDeviceObject, localAddress Address, deviceInfoCache *DeviceInfoCache, aseID *int) (*BIPSimpleApplication, error) {
	b := &BIPSimpleApplication{}
	var err error
	b.ApplicationIOController, err = NewApplicationIOController(localDevice, deviceInfoCache, aseID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	b.WhoIsIAmServices, err = NewWhoIsIAmServices(b)
	if err != nil {
		return nil, errors.Wrap(err, "error WhoIs/IAm services")
	}
	b.ReadWritePropertyServices, err = NewReadWritePropertyServices()
	if err != nil {
		return nil, errors.Wrap(err, "error read write property services")
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
	b.nsap, err = NewNetworkServiceAccessPoint(nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = NewNetworkServiceElement(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}
	if err := bind(b.nse, b.nsap); err != nil {
		return nil, errors.Wrap(err, "error binding network stack")
	}

	// bind the top layers
	if err := bind(b, b.asap, b.smap, b.nsap); err != nil {
		return nil, errors.Wrap(err, "error binding top layers")
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
		return nil, errors.Wrap(err, "error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := b.nsap.bind(b.bip, nil, &b.localAddress); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BIPSimpleApplication) Close() error {
	log.Debug().Msg("close socket")

	// pass to the multiplexer, then down to the sockets
	return b.mux.Close()
}

type BIPForeignApplication struct {
	*ApplicationIOController
	*WhoIsIAmServices
	*ReadWritePropertyServices
	localAddress Address
	asap         *ApplicationServiceAccessPoint
	smap         *StateMachineAccessPoint
	nsap         *NetworkServiceAccessPoint
	nse          *NetworkServiceElement
	bip          *BIPForeign
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer
}

func NewBIPForeignApplication(localDevice *LocalDeviceObject, localAddress Address, bbmdAddress *Address, bbmdTTL *int, deviceInfoCache *DeviceInfoCache, aseID *int) (*BIPForeignApplication, error) {
	b := &BIPForeignApplication{}
	var err error
	b.ApplicationIOController, err = NewApplicationIOController(localDevice, deviceInfoCache, aseID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating io controller")
	}
	b.WhoIsIAmServices, err = NewWhoIsIAmServices(b)
	if err != nil {
		return nil, errors.Wrap(err, "error WhoIs/IAm services")
	}
	b.ReadWritePropertyServices, err = NewReadWritePropertyServices()
	if err != nil {
		return nil, errors.Wrap(err, "error read write property services")
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
	b.nsap, err = NewNetworkServiceAccessPoint(nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	b.nse, err = NewNetworkServiceElement(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}
	if err := bind(b.nse, b.nsap); err != nil {
		return nil, errors.Wrap(err, "error binding network stack")
	}

	// bind the top layers
	if err := bind(b, b.asap, b.smap, b.nsap); err != nil {
		return nil, errors.Wrap(err, "error binding top layers")
	}

	// create a generic BIP stack, bound to the Annex J server on the UDP multiplexer
	b.bip, err = NewBIPForeign(bbmdAddress, bbmdTTL, nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new bip")
	}
	b.annexj, err = NewAnnexJCodec(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	b.mux, err = NewUDPMultiplexer(b.localAddress, true)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new udp multiplexer")
	}

	// bind the bottom layers
	if err := bind(b.bip, b.annexj, b.mux.annexJ); err != nil {
		return nil, errors.Wrap(err, "error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := b.nsap.bind(b.bip, nil, &b.localAddress); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *BIPForeignApplication) Close() error {
	log.Debug().Msg("close socket")

	// pass to the multiplexer, then down to the sockets
	return b.mux.Close()
}

type BIPNetworkApplication struct {
	*NetworkServiceElement
	localAddress Address
	nsap         *NetworkServiceAccessPoint
	bip          any // BIPSimple or BIPForeign
	annexj       *AnnexJCodec
	mux          *UDPMultiplexer
}

func NewBIPNetworkApplication(localAddress Address, bbmdAddress *Address, bbmdTTL *int, eID *int) (*BIPNetworkApplication, error) {
	n := &BIPNetworkApplication{}
	var err error
	n.NetworkServiceElement, err = NewNetworkServiceElement(eID)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new network service element")
	}

	n.localAddress = localAddress

	// a network service access point will be needed
	n.nsap, err = NewNetworkServiceAccessPoint(nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating network service access point")
	}

	// give the NSAP a generic network layer service element
	if err := bind(n, n.nsap); err != nil {
		return nil, errors.New("error binding network layer")
	}

	// create a generic BIP stack, bound to the Annex J server
	// on the UDP multiplexer
	if bbmdAddress == nil && bbmdTTL == nil {
		n.bip, err = NewBIPSimple(nil, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error creating BIPSimple")
		}
	} else {
		n.bip, err = NewBIPForeign(bbmdAddress, bbmdTTL, nil, nil, nil)
		if err != nil {
			return nil, errors.Wrap(err, "error creating BIPForeign")
		}
	}
	n.annexj, err = NewAnnexJCodec(nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new annex j codec")
	}
	n.mux, err = NewUDPMultiplexer(n.localAddress, true)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new udp multiplexer")
	}

	// bind the bottom layers
	if err := bind(n.bip, n.annexj, n.mux.annexJ); err != nil {
		return nil, errors.Wrap(err, "error binding bottom layers")
	}

	// bind the BIP stack to the network, no network number
	if err := n.nsap.bind(n.bip.(_Server), nil, &n.localAddress); err != nil {
		return nil, err
	}

	return n, nil
}
