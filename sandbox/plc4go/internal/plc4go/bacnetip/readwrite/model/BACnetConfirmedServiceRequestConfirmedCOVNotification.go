//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
	"strconv"
)

// Constant values.
const BACnetConfirmedServiceRequestConfirmedCOVNotification_SUBSCRIBERPROCESSIDENTIFIERHEADER uint8 = 0x09
const BACnetConfirmedServiceRequestConfirmedCOVNotification_MONITOREDOBJECTIDENTIFIERHEADER uint8 = 0x1C
const BACnetConfirmedServiceRequestConfirmedCOVNotification_ISSUECONFIRMEDNOTIFICATIONSHEADER uint8 = 0x2C
const BACnetConfirmedServiceRequestConfirmedCOVNotification_LIFETIMEHEADER uint8 = 0x07
const BACnetConfirmedServiceRequestConfirmedCOVNotification_LISTOFVALUESOPENINGTAG uint8 = 0x4E
const BACnetConfirmedServiceRequestConfirmedCOVNotification_LISTOFVALUESCLOSINGTAG uint8 = 0x4F

// The data-structure of this message
type BACnetConfirmedServiceRequestConfirmedCOVNotification struct {
	subscriberProcessIdentifier               uint8
	monitoredObjectType                       uint16
	monitoredObjectInstanceNumber             uint32
	issueConfirmedNotificationsType           uint16
	issueConfirmedNotificationsInstanceNumber uint32
	lifetimeLength                            uint8
	lifetimeSeconds                           []int8
	notifications                             []IBACnetTagWithContent
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestConfirmedCOVNotification interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) ServiceChoice() uint8 {
	return 0x01
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestConfirmedCOVNotification(subscriberProcessIdentifier uint8, monitoredObjectType uint16, monitoredObjectInstanceNumber uint32, issueConfirmedNotificationsType uint16, issueConfirmedNotificationsInstanceNumber uint32, lifetimeLength uint8, lifetimeSeconds []int8, notifications []IBACnetTagWithContent) BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestConfirmedCOVNotification{subscriberProcessIdentifier: subscriberProcessIdentifier, monitoredObjectType: monitoredObjectType, monitoredObjectInstanceNumber: monitoredObjectInstanceNumber, issueConfirmedNotificationsType: issueConfirmedNotificationsType, issueConfirmedNotificationsInstanceNumber: issueConfirmedNotificationsInstanceNumber, lifetimeLength: lifetimeLength, lifetimeSeconds: lifetimeSeconds, notifications: notifications}
}

func CastIBACnetConfirmedServiceRequestConfirmedCOVNotification(structType interface{}) IBACnetConfirmedServiceRequestConfirmedCOVNotification {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestConfirmedCOVNotification {
		if iBACnetConfirmedServiceRequestConfirmedCOVNotification, ok := typ.(IBACnetConfirmedServiceRequestConfirmedCOVNotification); ok {
			return iBACnetConfirmedServiceRequestConfirmedCOVNotification
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestConfirmedCOVNotification(structType interface{}) BACnetConfirmedServiceRequestConfirmedCOVNotification {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestConfirmedCOVNotification {
		if sBACnetConfirmedServiceRequestConfirmedCOVNotification, ok := typ.(BACnetConfirmedServiceRequestConfirmedCOVNotification); ok {
			return sBACnetConfirmedServiceRequestConfirmedCOVNotification
		}
		return BACnetConfirmedServiceRequestConfirmedCOVNotification{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	// Const Field (subscriberProcessIdentifierHeader)
	lengthInBits += 8

	// Simple field (subscriberProcessIdentifier)
	lengthInBits += 8

	// Const Field (monitoredObjectIdentifierHeader)
	lengthInBits += 8

	// Simple field (monitoredObjectType)
	lengthInBits += 10

	// Simple field (monitoredObjectInstanceNumber)
	lengthInBits += 22

	// Const Field (issueConfirmedNotificationsHeader)
	lengthInBits += 8

	// Simple field (issueConfirmedNotificationsType)
	lengthInBits += 10

	// Simple field (issueConfirmedNotificationsInstanceNumber)
	lengthInBits += 22

	// Const Field (lifetimeHeader)
	lengthInBits += 5

	// Simple field (lifetimeLength)
	lengthInBits += 3

	// Array field
	if len(m.lifetimeSeconds) > 0 {
		lengthInBits += 8 * uint16(len(m.lifetimeSeconds))
	}

	// Const Field (listOfValuesOpeningTag)
	lengthInBits += 8

	// Array field
	if len(m.notifications) > 0 {
		for _, element := range m.notifications {
			lengthInBits += element.LengthInBits()
		}
	}

	// Const Field (listOfValuesClosingTag)
	lengthInBits += 8

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestConfirmedCOVNotificationParse(io spi.ReadBuffer, len uint16) (BACnetConfirmedServiceRequestInitializer, error) {

	// Const Field (subscriberProcessIdentifierHeader)
	subscriberProcessIdentifierHeader, _subscriberProcessIdentifierHeaderErr := io.ReadUint8(8)
	if _subscriberProcessIdentifierHeaderErr != nil {
		return nil, errors.New("Error parsing 'subscriberProcessIdentifierHeader' field " + _subscriberProcessIdentifierHeaderErr.Error())
	}
	if subscriberProcessIdentifierHeader != BACnetConfirmedServiceRequestConfirmedCOVNotification_SUBSCRIBERPROCESSIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestConfirmedCOVNotification_SUBSCRIBERPROCESSIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(subscriberProcessIdentifierHeader)))
	}

	// Simple Field (subscriberProcessIdentifier)
	subscriberProcessIdentifier, _subscriberProcessIdentifierErr := io.ReadUint8(8)
	if _subscriberProcessIdentifierErr != nil {
		return nil, errors.New("Error parsing 'subscriberProcessIdentifier' field " + _subscriberProcessIdentifierErr.Error())
	}

	// Const Field (monitoredObjectIdentifierHeader)
	monitoredObjectIdentifierHeader, _monitoredObjectIdentifierHeaderErr := io.ReadUint8(8)
	if _monitoredObjectIdentifierHeaderErr != nil {
		return nil, errors.New("Error parsing 'monitoredObjectIdentifierHeader' field " + _monitoredObjectIdentifierHeaderErr.Error())
	}
	if monitoredObjectIdentifierHeader != BACnetConfirmedServiceRequestConfirmedCOVNotification_MONITOREDOBJECTIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestConfirmedCOVNotification_MONITOREDOBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(monitoredObjectIdentifierHeader)))
	}

	// Simple Field (monitoredObjectType)
	monitoredObjectType, _monitoredObjectTypeErr := io.ReadUint16(10)
	if _monitoredObjectTypeErr != nil {
		return nil, errors.New("Error parsing 'monitoredObjectType' field " + _monitoredObjectTypeErr.Error())
	}

	// Simple Field (monitoredObjectInstanceNumber)
	monitoredObjectInstanceNumber, _monitoredObjectInstanceNumberErr := io.ReadUint32(22)
	if _monitoredObjectInstanceNumberErr != nil {
		return nil, errors.New("Error parsing 'monitoredObjectInstanceNumber' field " + _monitoredObjectInstanceNumberErr.Error())
	}

	// Const Field (issueConfirmedNotificationsHeader)
	issueConfirmedNotificationsHeader, _issueConfirmedNotificationsHeaderErr := io.ReadUint8(8)
	if _issueConfirmedNotificationsHeaderErr != nil {
		return nil, errors.New("Error parsing 'issueConfirmedNotificationsHeader' field " + _issueConfirmedNotificationsHeaderErr.Error())
	}
	if issueConfirmedNotificationsHeader != BACnetConfirmedServiceRequestConfirmedCOVNotification_ISSUECONFIRMEDNOTIFICATIONSHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestConfirmedCOVNotification_ISSUECONFIRMEDNOTIFICATIONSHEADER)) + " but got " + strconv.Itoa(int(issueConfirmedNotificationsHeader)))
	}

	// Simple Field (issueConfirmedNotificationsType)
	issueConfirmedNotificationsType, _issueConfirmedNotificationsTypeErr := io.ReadUint16(10)
	if _issueConfirmedNotificationsTypeErr != nil {
		return nil, errors.New("Error parsing 'issueConfirmedNotificationsType' field " + _issueConfirmedNotificationsTypeErr.Error())
	}

	// Simple Field (issueConfirmedNotificationsInstanceNumber)
	issueConfirmedNotificationsInstanceNumber, _issueConfirmedNotificationsInstanceNumberErr := io.ReadUint32(22)
	if _issueConfirmedNotificationsInstanceNumberErr != nil {
		return nil, errors.New("Error parsing 'issueConfirmedNotificationsInstanceNumber' field " + _issueConfirmedNotificationsInstanceNumberErr.Error())
	}

	// Const Field (lifetimeHeader)
	lifetimeHeader, _lifetimeHeaderErr := io.ReadUint8(5)
	if _lifetimeHeaderErr != nil {
		return nil, errors.New("Error parsing 'lifetimeHeader' field " + _lifetimeHeaderErr.Error())
	}
	if lifetimeHeader != BACnetConfirmedServiceRequestConfirmedCOVNotification_LIFETIMEHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestConfirmedCOVNotification_LIFETIMEHEADER)) + " but got " + strconv.Itoa(int(lifetimeHeader)))
	}

	// Simple Field (lifetimeLength)
	lifetimeLength, _lifetimeLengthErr := io.ReadUint8(3)
	if _lifetimeLengthErr != nil {
		return nil, errors.New("Error parsing 'lifetimeLength' field " + _lifetimeLengthErr.Error())
	}

	// Array field (lifetimeSeconds)
	var lifetimeSeconds []int8
	// Count array
	{
		lifetimeSeconds := make([]int8, lifetimeLength)
		for curItem := uint16(0); curItem < uint16(lifetimeLength); curItem++ {

			_lifetimeSecondsVal, _err := io.ReadInt8(8)
			if _err != nil {
				return nil, errors.New("Error parsing 'lifetimeSeconds' field " + _err.Error())
			}
			lifetimeSeconds = append(lifetimeSeconds, _lifetimeSecondsVal)
		}
	}

	// Const Field (listOfValuesOpeningTag)
	listOfValuesOpeningTag, _listOfValuesOpeningTagErr := io.ReadUint8(8)
	if _listOfValuesOpeningTagErr != nil {
		return nil, errors.New("Error parsing 'listOfValuesOpeningTag' field " + _listOfValuesOpeningTagErr.Error())
	}
	if listOfValuesOpeningTag != BACnetConfirmedServiceRequestConfirmedCOVNotification_LISTOFVALUESOPENINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestConfirmedCOVNotification_LISTOFVALUESOPENINGTAG)) + " but got " + strconv.Itoa(int(listOfValuesOpeningTag)))
	}

	// Array field (notifications)
	var notifications []IBACnetTagWithContent
	// Length array
	_notificationsLength := uint16(len) - uint16(uint16(18))
	_notificationsEndPos := io.GetPos() + uint16(_notificationsLength)
	for io.GetPos() < _notificationsEndPos {
		_message, _err := BACnetTagWithContentParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'notifications' field " + _err.Error())
		}
		var _item IBACnetTagWithContent
		_item, _ok := _message.(IBACnetTagWithContent)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to BACnetTagWithContent")
		}
		notifications = append(notifications, _item)
	}

	// Const Field (listOfValuesClosingTag)
	listOfValuesClosingTag, _listOfValuesClosingTagErr := io.ReadUint8(8)
	if _listOfValuesClosingTagErr != nil {
		return nil, errors.New("Error parsing 'listOfValuesClosingTag' field " + _listOfValuesClosingTagErr.Error())
	}
	if listOfValuesClosingTag != BACnetConfirmedServiceRequestConfirmedCOVNotification_LISTOFVALUESCLOSINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestConfirmedCOVNotification_LISTOFVALUESCLOSINGTAG)) + " but got " + strconv.Itoa(int(listOfValuesClosingTag)))
	}

	// Create the instance
	return NewBACnetConfirmedServiceRequestConfirmedCOVNotification(subscriberProcessIdentifier, monitoredObjectType, monitoredObjectInstanceNumber, issueConfirmedNotificationsType, issueConfirmedNotificationsInstanceNumber, lifetimeLength, lifetimeSeconds, notifications), nil
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) Serialize(io spi.WriteBuffer) {

	// Const Field (subscriberProcessIdentifierHeader)
	io.WriteUint8(8, 0x09)

	// Simple Field (subscriberProcessIdentifier)
	subscriberProcessIdentifier := uint8(m.subscriberProcessIdentifier)
	io.WriteUint8(8, (subscriberProcessIdentifier))

	// Const Field (monitoredObjectIdentifierHeader)
	io.WriteUint8(8, 0x1C)

	// Simple Field (monitoredObjectType)
	monitoredObjectType := uint16(m.monitoredObjectType)
	io.WriteUint16(10, (monitoredObjectType))

	// Simple Field (monitoredObjectInstanceNumber)
	monitoredObjectInstanceNumber := uint32(m.monitoredObjectInstanceNumber)
	io.WriteUint32(22, (monitoredObjectInstanceNumber))

	// Const Field (issueConfirmedNotificationsHeader)
	io.WriteUint8(8, 0x2C)

	// Simple Field (issueConfirmedNotificationsType)
	issueConfirmedNotificationsType := uint16(m.issueConfirmedNotificationsType)
	io.WriteUint16(10, (issueConfirmedNotificationsType))

	// Simple Field (issueConfirmedNotificationsInstanceNumber)
	issueConfirmedNotificationsInstanceNumber := uint32(m.issueConfirmedNotificationsInstanceNumber)
	io.WriteUint32(22, (issueConfirmedNotificationsInstanceNumber))

	// Const Field (lifetimeHeader)
	io.WriteUint8(5, 0x07)

	// Simple Field (lifetimeLength)
	lifetimeLength := uint8(m.lifetimeLength)
	io.WriteUint8(3, (lifetimeLength))

	// Array Field (lifetimeSeconds)
	if m.lifetimeSeconds != nil {
		for _, _element := range m.lifetimeSeconds {
			io.WriteInt8(8, _element)
		}
	}

	// Const Field (listOfValuesOpeningTag)
	io.WriteUint8(8, 0x4E)

	// Array Field (notifications)
	if m.notifications != nil {
		for _, _element := range m.notifications {
			_element.Serialize(io)
		}
	}

	// Const Field (listOfValuesClosingTag)
	io.WriteUint8(8, 0x4F)
}
