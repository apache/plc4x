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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/object"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type LocalDeviceObject interface {
	fmt.Stringer
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

func NewLocalDeviceObject(args Args, kwArgs KWArgs) LocalDeviceObject {
	return &localDeviceObject{}
}

type localDeviceObject struct {
	*CurrentPropertyListMixIn
	*object.DeviceObject

	// TODO: replace below...
	NumberOfAPDURetries       *uint
	APDUTimeout               *uint
	SegmentationSupported     *readWriteModel.BACnetSegmentation `directSerialize:"true"`
	APDUSegmentTimeout        *uint
	MaxSegmentsAccepted       *readWriteModel.MaxSegmentsAccepted   `directSerialize:"true"`
	MaximumApduLengthAccepted *readWriteModel.MaxApduLengthAccepted `directSerialize:"true"`
	App                       any                                   `ignore:"true"` // TODO: is *Application but creates a circular dependency. figure out what is a smart way
	ObjectName                string
	ObjectIdentifier          string
	VendorIdentifier          uint16
	ObjectList                []string
}

func (l *localDeviceObject) GetObjectIdentifier() string {
	return l.ObjectIdentifier
}

func (l *localDeviceObject) GetMaximumApduLengthAccepted() *readWriteModel.MaxApduLengthAccepted {
	return l.MaximumApduLengthAccepted
}

func (l *localDeviceObject) GetSegmentationSupported() *readWriteModel.BACnetSegmentation {
	return l.SegmentationSupported
}

func (l *localDeviceObject) GetVendorIdentifier() any {
	return l.VendorIdentifier
}

func (l *localDeviceObject) GetNumberOfAPDURetries() *uint {
	return l.NumberOfAPDURetries
}

func (l *localDeviceObject) GetAPDUTimeout() *uint {
	return l.APDUTimeout
}

func (l *localDeviceObject) SetApp(a any) {
	l.App = a
}

func (l *localDeviceObject) GetAPDUSegmentTimeout() *uint {
	return l.APDUSegmentTimeout
}

func (l *localDeviceObject) GetObjectName() string {
	return l.ObjectName
}

func (l *localDeviceObject) GetMaxSegmentsAccepted() *readWriteModel.MaxSegmentsAccepted {
	return l.MaxSegmentsAccepted
}

func (l *localDeviceObject) GetObjectList() []string {
	return l.ObjectList
}

func (l *localDeviceObject) SetObjectList(strings []string) {
	l.ObjectList = strings
}

func (l *localDeviceObject) String() string {
	panic("implementme")
}
