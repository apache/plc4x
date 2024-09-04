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

package bacgopes

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ApplicationRequirements interface {
	ApplicationServiceElementRequirements
}

type Application struct {
	ApplicationServiceElementContract
	Collector

	objectName       map[string]*LocalDeviceObject
	objectIdentifier map[string]*LocalDeviceObject
	localDevice      *LocalDeviceObject
	deviceInfoCache  *DeviceInfoCache
	controllers      map[string]any
	helpers          map[string]func(pdu PDU) error

	_startupDisabled bool

	// pass through args
	argAseID *int

	log zerolog.Logger
}

func NewApplication(localLog zerolog.Logger, localDevice *LocalDeviceObject, opts ...func(*Application)) (*Application, error) {
	a := &Application{
		log: localLog,
	}
	for _, opt := range opts {
		opt(a)
	}
	localLog.Debug().
		Interface("localDevice", localDevice).
		Interface("deviceInfoCache", a.deviceInfoCache).
		Interface("aseID", a.argAseID).
		Msg("NewApplication")
	var err error
	a.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, func(element *applicationServiceElement) {
		element.elementID = a.argAseID
	})
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
	if a.deviceInfoCache == nil {
		a.deviceInfoCache = NewDeviceInfoCache(localLog)
	}

	// controllers for managing confirmed requests as a client
	a.controllers = map[string]any{}

	// now set up the rest of the capabilities
	a.Collector = Collector{}

	// if starting up is enabled, find all the startup functions
	if !a._startupDisabled {
		for _, fn := range a.CapabilityFunctions("startup") {
			localLog.Debug().Interface("fn", fn).Msg("startup fn")
			Deferred(fn, NoArgs, NoKWArgs)
		}
	}
	return a, nil
}

func WithApplicationAseID(aseID int) func(*Application) {
	return func(a *Application) {
		a.argAseID = &aseID
	}
}

func WithApplicationDeviceInfoCache(deviceInfoCache *DeviceInfoCache) func(*Application) {
	return func(a *Application) {
		a.deviceInfoCache = deviceInfoCache
	}
}

func (a *Application) GetDeviceInfoCache() *DeviceInfoCache {
	return a.deviceInfoCache
}

func (a *Application) String() string {
	return fmt.Sprintf("Application(TBD...)") // TODO: fill some info here
}

// AddObject adds an object to the local collection
func (a *Application) AddObject(obj *LocalDeviceObject) error {
	a.log.Debug().Stringer("obj", obj).Msg("AddObject")

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
	a.log.Debug().Stringer("obj", obj).Msg("DeleteObject")

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
//
//	match confirmed and unconfirmed services.
//
// TODO: match that with readWriteModel.BACnetServicesSupported
func (a *Application) GetServicesSupported() []string {
	servicesSupported := make([]string, 0, len(a.helpers))
	for key := range a.helpers {
		servicesSupported = append(servicesSupported, key)
	}
	return servicesSupported
}

func (a *Application) Request(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")
	apdu := args.Get0PDU()

	// double-check the input is the right kind of APDU
	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUUnconfirmedRequest, readWriteModel.APDUConfirmedRequest:
	default:
		return errors.New("APDU expected")
	}
	return a.ApplicationServiceElementContract.Request(args, kwargs)
}

func (a *Application) Indication(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	apdu := args.Get0PDU()

	// get a helper function
	helperName := fmt.Sprintf("Do_%T", apdu)
	helperFn := a.helpers[helperName]
	a.log.Debug().
		Str("helperName", helperName).
		Interface("helperFn", helperFn).
		Msg("working with helper")

	// send back a reject for unrecognized services
	if helperFn == nil {
		if _, ok := apdu.(readWriteModel.APDUConfirmedRequest); ok {
			return errors.Errorf("no function %s", helperName)
		}
		return nil
	}

	if err := helperFn(apdu); err != nil {
		a.log.Debug().Err(err).Msg("err result")
		// TODO: do proper mapping
		if err := a.Response(NewArgs(NewPDU(readWriteModel.NewAPDUError(0, readWriteModel.BACnetConfirmedServiceChoice_CREATE_OBJECT, nil, 0))), kwargs); err != nil {
			return err
		}
	}

	return nil
}
