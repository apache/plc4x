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

package app

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/appservice"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/capability"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/core"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type ApplicationRequirements interface {
	ApplicationServiceElementRequirements
}

//go:generate go tool plc4xGenerator -type=Application -prefix=app_
type Application struct {
	ApplicationServiceElementContract
	Collector

	objectName       map[string]LocalDeviceObject
	objectIdentifier map[string]LocalDeviceObject
	localDevice      LocalDeviceObject
	localAddress     *Address
	deviceInfoCache  *appservice.DeviceInfoCache
	controllers      map[string]any
	helpers          map[string]func(apdu APDU) error `ignore:"true"`

	_startupDisabled bool

	log zerolog.Logger
}

func NewApplication(localLog zerolog.Logger, options ...Option) (*Application, error) {
	a := &Application{
		log:     localLog,
		helpers: map[string]func(apdu APDU) error{},
	}
	ApplyAppliers(options, a)
	optionsForParent := AddLeafTypeIfAbundant(options, a)
	var err error
	a.ApplicationServiceElementContract, err = NewApplicationServiceElement(localLog, optionsForParent...)
	if err != nil {
		return nil, err
	}
	localLog.Debug().
		Stringer("localDevice", a.localDevice).
		Stringer("deviceInfoCache", a.deviceInfoCache).
		Interface("aseID", a.GetElementId()).
		Msg("NewApplication")
	if _debug != nil {
		_debug("__init__ %r %r deviceInfoCache=%r aseID=%r", a.localDevice, a.localAddress, a.deviceInfoCache, a.GetElementId())
	}

	// local objects by ID and name
	a.objectName = map[string]LocalDeviceObject{}
	a.objectIdentifier = map[string]LocalDeviceObject{}

	// keep track of the local device
	if a.localDevice != nil {
		// bind the device object to this application
		a.localDevice.SetApp(a)

		// local objects by ID and name
		a.objectName[a.localDevice.GetObjectName()] = a.localDevice
		a.objectName[a.localDevice.GetObjectIdentifier()] = a.localDevice
	}

	// use the provided cache or make a default one
	if a.deviceInfoCache == nil {
		a.deviceInfoCache = appservice.NewDeviceInfoCache(localLog)
	}

	// controllers for managing confirmed requests as a client
	a.controllers = map[string]any{}

	// now set up the rest of the capabilities
	var init func()
	a.Collector, init = NewCollector(a.log) // TODO: pass a if it has stuff to collect
	init()

	// if starting up is enabled, find all the startup functions
	if !a._startupDisabled {
		for fn := range a.CapabilityFunctions("startup") {
			if _debug != nil {
				_debug("    - startup fn: %t", fn != nil)
			}
			localLog.Debug().Interface("fn", fn).Msg("startup fn")
			Deferred(fn, NoArgs, NoKWArgs())
		}
	}
	return a, nil
}

func WithApplicationLocalDeviceObject(localDevice LocalDeviceObject) GenericApplier[*Application] {
	return WrapGenericApplier(func(a *Application) { a.localDevice = localDevice })
}

func WithApplicationDeviceInfoCache(deviceInfoCache *appservice.DeviceInfoCache) GenericApplier[*Application] {
	return WrapGenericApplier(func(a *Application) { a.deviceInfoCache = deviceInfoCache })
}

func (a *Application) GetDeviceInfoCache() *appservice.DeviceInfoCache {
	return a.deviceInfoCache
}

// AddObject adds an object to the local collection
func (a *Application) AddObject(obj LocalDeviceObject) error {
	a.log.Debug().Stringer("obj", obj).Msg("AddObject")
	if _debug != nil {
		_debug("add_object %r", obj)
	}

	// extract the object name and identifier
	objectName := obj.GetObjectName()
	if objectName == "" {
		return errors.New("object name required")
	}
	objectIdentifier := obj.GetObjectIdentifier()
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
		a.localDevice.SetObjectList(append(a.localDevice.GetObjectList(), objectIdentifier))
	}

	// let the object know which application stack it belongs to
	obj.SetApp(a)

	return nil
}

// DeleteObject deletes an object from the local collection
func (a *Application) DeleteObject(obj LocalDeviceObject) error {
	a.log.Debug().Stringer("obj", obj).Msg("DeleteObject")
	if _debug != nil {
		_debug("delete_object %r", obj)
	}

	// extract the object name and identifier
	objectName := obj.GetObjectName()
	objectIdentifier := obj.GetObjectIdentifier()

	// delete it from the application
	delete(a.objectName, objectName)
	delete(a.objectIdentifier, objectIdentifier)

	// remove the object's identifier from the device's object list if there is one and has an object list property
	if a.localDevice != nil {
		foundIndex := -1
		for i, s := range a.localDevice.GetObjectList() {
			if s == objectIdentifier {
				foundIndex = i
			}
		}
		if foundIndex >= 0 {
			a.localDevice.SetObjectList(append(a.localDevice.GetObjectList()[0:foundIndex], a.localDevice.GetObjectList()[foundIndex+1:]...))
		}
	}

	// make sure the object knows it's detached from an application
	obj.SetApp(nil)

	return nil
}

// GetObjectId returns a local object or None.
func (a *Application) GetObjectId(objectId string) LocalDeviceObject {
	return a.objectIdentifier[objectId]
}

// GetObjectName returns a local object or None.
func (a *Application) GetObjectName(objectName string) LocalDeviceObject {
	return a.objectName[objectName]
}

// IterObjects iterates over the objects
func (a *Application) IterObjects() []LocalDeviceObject {
	localDeviceObjects := make([]LocalDeviceObject, 0, len(a.objectIdentifier))
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
	if _debug != nil {
		_debug("get_services_supported")
	}
	servicesSupported := make([]string, 0, len(a.helpers))
	for key := range a.helpers {
		servicesSupported = append(servicesSupported, key)
	}
	return servicesSupported
}

func (a *Application) Request(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Request")
	apdu := GA[APDU](args, 0)
	if _debug != nil {
		_debug("request %r", apdu)
	}

	// double-check the input is the right kind of APDU
	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUUnconfirmedRequest, readWriteModel.APDUConfirmedRequest:
	default:
		return errors.New("APDU expected")
	}
	return a.ApplicationServiceElementContract.Request(args, kwArgs)
}

func (a *Application) Indication(args Args, kwArgs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	apdu := GA[APDU](args, 0)
	if _debug != nil {
		_debug("indication %r", apdu)
	}

	// get a helper function
	helperName := fmt.Sprintf("Do_%T", apdu)
	helperFn := a.helpers[helperName]
	a.log.Debug().
		Str("helperName", helperName).
		Bool("helperFn", helperFn != nil).
		Msg("working with helper")
	if _debug != nil {
		_debug("    - helperFn: %p", helperFn)
	}

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
		if _, ok := apdu.(readWriteModel.APDUConfirmedRequest); ok {
			resp, err := NewError(NoArgs, NKW(KWErrorClass, "device", KWErrorCode, "operationalProblem", KWContext, apdu))
			if err != nil {
				return errors.Wrap(err, "error creating error")
			}
			if err := a.Response(NA(resp), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error sending response")
			}
		}
	}

	return nil
}

func (a *Application) RegisterHelperFn(name string, fn func(apdu APDU) error) error {
	if _, ok := a.helpers[name]; ok {
		return errors.Errorf("helper %s already registered", name)
	}
	a.helpers[name] = fn
	return nil
}
