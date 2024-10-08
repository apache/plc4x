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

package device

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/constructeddata"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/object"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/object"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type LocalDeviceObject interface {
	fmt.Stringer
	fmt.Formatter
	GetObjectIdentifier() string
	GetMaximumApduLengthAccepted() *readWriteModel.MaxApduLengthAccepted
	GetSegmentationSupported() *readWriteModel.BACnetSegmentation
	GetVendorIdentifier() any
	GetNumberOfAPDURetries() *uint
	GetAPDUTimeout() *uint
	SetApp(a any)
	GetAPDUSegmentTimeout() *uint
	GetObjectName() string
	GetMaxSegmentsAccepted() *readWriteModel.MaxSegmentsAccepted
	GetObjectList() []string
	SetObjectList([]string)
}

func NewLocalDeviceObject(_ Args, kwArgs KWArgs, options ...Option) (LocalDeviceObject, error) {
	if _debug != nil {
		_debug("__init__ %r", kwArgs)
	}
	l := &_LocalDeviceObject{
		DefaultRFormatter: NewDefaultRFormatter(),
		properties: []Property{
			NewCurrentLocalTime(),
			NewCurrentLocalDate(),
			NewCurrentProtocolServicesSupported(),
		},
		defaultProperties: map[string]any{
			"maxApduLengthAccepted": 1024,
			"segmentationSupported": "segmentedBoth",
			"maxSegmentsAccepted":   16,
			"apduSegmentTimeout":    5000,
			"apduTimeout":           3000,
			"numberOfApduRetries":   3,
		},
	}
	options = AddSharedSuperIfAbundant[Object](options)
	var err error
	l.CurrentPropertyListMixIn, err = NewCurrentPropertyListMixIn(options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating mixing")
	}
	l.DeviceObject, err = NewDeviceObject(options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating device object")
	}

	// start with an empty dictionary of device object properties
	initArgs := make(KWArgs)
	iniArg, _ := KWO[KWArgs](kwArgs, "ini", nil)
	if _debug != nil {
		_debug("    - ini_arg: %r", iniArg)
	}

	if _, ok := RegisteredObjectTypes[fmt.Sprintf("%T", l)]; !ok {
		if _debug != nil {
			_debug("    - unregistered", kwArgs)
		}
		vendorIdentifier, ok := KWO(kwArgs, KWVendorIdentifier, -1)
		if _debug != nil {
			_debug("    - keyword vendor identifier: %r", vendorIdentifier)
		}

		if !ok {
			vendorIdentifier, ok = KWO[int](iniArg, "vendorIdentifier", -1)
			_debug("    - INI vendor identifier: %r", vendorIdentifier)
		}

		if !ok {
			return nil, errors.New("vendorIdentifier required to auto-register the LocalDeviceObject class")
		}

		if _, err := RegisterObjectType(NKW(KWCls, l, KWVendorIdentifier, vendorIdentifier)); err != nil {
			return nil, errors.Wrap(err, "error register object type")
		}
	}

	// look for properties, fill in values from the keyword arguments or
	// the INI parameter (converted to a proper value) if it was provided
	for propid, prop := range l.DeviceObject.Get_Properties() {
		// special processing for object identifier
		if propid == "objectIdentifier" {
			continue
		}

		// use keyword argument if it was provided
		var propValue any
		var propDataType any
		if v, ok := kwArgs[KnownKey(propid)]; ok {
			propValue = v
		} else {
			propValue, ok = KWO[any](iniArg, KnownKey(strings.ToLower(propid)), nil)
			if !ok {
				continue
			}

			propDataType = prop.GetDataType()
			// TODO: convert
			_ = propDataType
			// at long last
			initArgs[KnownKey(propid)] = propValue
		}
	}
	// check for object identifier as a keyword parameter or in the INI file,
	// and it might be just an int, so make it a tuple if necessary
	var objectIdentifier any
	if v, ok := KWO[any](kwArgs, "objectIdentifier", nil); ok {
		objectIdentifier = v
		if vint, ok := v.(int); ok {
			objectIdentifier = ObjectIdentifierTuple{Left: "device", Right: vint}
		}
	} else if v, ok = KWO[any](iniArg, "objectidentifier", nil); ok {
		objectIdentifier = ObjectIdentifierTuple{Left: "device", Right: v.(int)}
	} else {
		return nil, errors.New("objectIdentifier required")
	}
	initArgs["objectIdentifier"] = objectIdentifier
	if _debug != nil {
		_debug("    - object identifier: %r", objectIdentifier)
	}

	// fill in default property values not in init_args
	for attr, value := range l.defaultProperties {
		if _, ok := initArgs[KnownKey(attr)]; !ok {
			initArgs[KnownKey(attr)] = value
		}
	}

	// check for properties this class implements
	if _, ok := initArgs[KnownKey("localDate")]; ok {
		return nil, errors.New("localDate is provided by LocalDeviceObject and cannot be overridden")
	}
	if _, ok := initArgs[KnownKey("localTime")]; ok {
		return nil, errors.New("localTime is provided by LocalDeviceObject and cannot be overridden")
	}
	if _, ok := initArgs[KnownKey("protocolServicesSupported")]; ok {
		return nil, errors.New("protocolServicesSupported is provided by LocalDeviceObject and cannot be overridden")
	}

	// the object list is provided
	if _, ok := initArgs[KnownKey("objectList")]; ok {
		return nil, errors.New("objectList is provided by LocalDeviceObject and cannot be overridden")
	}
	initArgs["objectList"] = ArrayOfEs[ObjectIdentifier](NewObjectIdentifier, 0, nil)

	// check for a minimum value
	if v, ok := KWO[int](kwArgs, "maxApduLengthAccepted", 0); ok && v < 50 {
		return nil, errors.New("invalid max APDU length accepted")
	}

	// dump the updated attributes
	if _debug != nil {
		_debug("    - init_args: %r", initArgs)
	}

	// proceed as usual
	if err := l.DeviceObject.Init(NoArgs, initArgs); err != nil {
		return nil, errors.Wrap(err, "error creating device object")
	}

	// pass along special property values that are not BACnet properties
	for key, value := range kwArgs {
		if strings.HasPrefix(string(key), "_") {
			l.DeviceObject.SetAttr(string(key), value)
		}
	}
	return l, nil
}

type _LocalDeviceObject struct {
	*CurrentPropertyListMixIn
	*DeviceObject
	*DefaultRFormatter

	properties        []Property
	defaultProperties map[string]any

	App any `ignore:"true"` // TODO: is *Application but creates a circular dependency. figure out what is a smart way
}

var _ Object = (*_LocalDeviceObject)(nil)

func (l *_LocalDeviceObject) GetObjectIdentifier() string {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("objectIdentifier")
	if !ok {
		return ""
	}
	return attr.(string)
}

func (l *_LocalDeviceObject) GetMaximumApduLengthAccepted() *readWriteModel.MaxApduLengthAccepted {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("maximumApduLengthAccepted")
	if !ok {
		return nil
	}
	return attr.(*readWriteModel.MaxApduLengthAccepted)
}

func (l *_LocalDeviceObject) GetSegmentationSupported() *readWriteModel.BACnetSegmentation {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("segmentationSupported")
	if !ok {
		return nil
	}
	return attr.(*readWriteModel.BACnetSegmentation)
}

func (l *_LocalDeviceObject) GetVendorIdentifier() any {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("vendorIdentifier")
	if !ok {
		return 0
	}
	return attr.(int)
}

func (l *_LocalDeviceObject) GetNumberOfAPDURetries() *uint {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("numberOfAPDURetries")
	if !ok {
		return nil
	}
	return attr.(*uint)
}

func (l *_LocalDeviceObject) GetAPDUTimeout() *uint {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("apduTimeout")
	if !ok {
		return nil
	}
	return attr.(*uint)
}

func (l *_LocalDeviceObject) SetApp(a any) {
	l.App = a
}

func (l *_LocalDeviceObject) GetAPDUSegmentTimeout() *uint {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("apduSegmentTimeout")
	if !ok {
		return nil
	}
	return attr.(*uint)
}

func (l *_LocalDeviceObject) GetObjectName() string {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("objectName")
	if !ok {
		return ""
	}
	return attr.(string)
}

func (l *_LocalDeviceObject) GetMaxSegmentsAccepted() *readWriteModel.MaxSegmentsAccepted {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("maxSegmentsAccepted")
	if !ok {
		return nil
	}
	return attr.(*readWriteModel.MaxSegmentsAccepted)
}

func (l *_LocalDeviceObject) GetObjectList() []string {
	attr, ok := l.CurrentPropertyListMixIn.GetAttr("objectList")
	if !ok {
		return nil
	}
	return attr.([]string)
}

func (l *_LocalDeviceObject) SetObjectList(strings []string) {
	l.CurrentPropertyListMixIn.SetAttr("objectList", strings)
}

func (l *_LocalDeviceObject) Init(args Args, kwArgs KWArgs) error {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.Init(args, kwArgs)
}

func (l *_LocalDeviceObject) GetObjectType() string {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.GetObjectType()
}

func (l *_LocalDeviceObject) GetAttr(name string) (any, bool) {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.GetAttr(name)
}

func (l *_LocalDeviceObject) SetAttr(name string, value any) {
	// ambiguous selector avoidance
	l.CurrentPropertyListMixIn.SetAttr(name, value)
}

func (l *_LocalDeviceObject) AddProperty(prop Property) {
	// ambiguous selector avoidance
	l.CurrentPropertyListMixIn.AddProperty(prop)
}

func (l *_LocalDeviceObject) ReadProperty(args Args, kwArgs KWArgs) error {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.ReadProperty(args, kwArgs)
}

func (l *_LocalDeviceObject) WriteProperty(args Args, kwArgs KWArgs) error {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.WriteProperty(args, kwArgs)
}

func (l *_LocalDeviceObject) DeleteProperty(prop string) {
	// ambiguous selector avoidance
	l.CurrentPropertyListMixIn.DeleteProperty(prop)
}

func (l *_LocalDeviceObject) GetProperties() []Property {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.GetProperties()
}

func (l *_LocalDeviceObject) Get_Properties() map[string]Property {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.Get_Properties()
}

func (l *_LocalDeviceObject) Get_PropertiesMonitors() map[string][]func(old, new any) {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.Get_PropertiesMonitors()
}

func (l *_LocalDeviceObject) Get_Values() map[string]any {
	// ambiguous selector avoidance
	return l.CurrentPropertyListMixIn.Get_Values()
}

func (l *_LocalDeviceObject) Set_Properties(_properties map[string]Property) {
	// ambiguous selector avoidance
	l.CurrentPropertyListMixIn.Set_Properties(_properties)
}

func (l *_LocalDeviceObject) PrintDebugContents(indent int, file io.Writer, _ids []uintptr) {
	// ambiguous selector avoidance
	l.CurrentPropertyListMixIn.PrintDebugContents(indent, file, _ids)
}
