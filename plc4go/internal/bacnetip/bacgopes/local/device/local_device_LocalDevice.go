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

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/object"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/object"
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

type _LocalDeviceObject struct {
	*CurrentPropertyListMixIn
	*DeviceObject
	*DefaultRFormatter

	properties        []Property
	defaultProperties map[string]any

	App any `ignore:"true"` // TODO: is *Application but creates a circular dependency. figure out what is a smart way
}

func NewLocalDeviceObject(args Args, kwArgs KWArgs, options ...Option) (LocalDeviceObject, error) {
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
	l.CurrentPropertyListMixIn, err = NewCurrentPropertyListMixIn(kwArgs, options...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating mixing")
	}
	return l, nil
}

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
		return ""
	}
	return attr.(string)
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
