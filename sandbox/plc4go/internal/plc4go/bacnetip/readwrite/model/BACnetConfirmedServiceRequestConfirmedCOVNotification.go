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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
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
    SubscriberProcessIdentifier uint8
    MonitoredObjectType uint16
    MonitoredObjectInstanceNumber uint32
    IssueConfirmedNotificationsType uint16
    IssueConfirmedNotificationsInstanceNumber uint32
    LifetimeLength uint8
    LifetimeSeconds []int8
    Notifications []IBACnetTagWithContent
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestConfirmedCOVNotification interface {
    IBACnetConfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) ServiceChoice() uint8 {
    return 0x01
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestConfirmedCOVNotification(subscriberProcessIdentifier uint8, monitoredObjectType uint16, monitoredObjectInstanceNumber uint32, issueConfirmedNotificationsType uint16, issueConfirmedNotificationsInstanceNumber uint32, lifetimeLength uint8, lifetimeSeconds []int8, notifications []IBACnetTagWithContent) BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestConfirmedCOVNotification{SubscriberProcessIdentifier: subscriberProcessIdentifier, MonitoredObjectType: monitoredObjectType, MonitoredObjectInstanceNumber: monitoredObjectInstanceNumber, IssueConfirmedNotificationsType: issueConfirmedNotificationsType, IssueConfirmedNotificationsInstanceNumber: issueConfirmedNotificationsInstanceNumber, LifetimeLength: lifetimeLength, LifetimeSeconds: lifetimeSeconds, Notifications: notifications}
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
        if sBACnetConfirmedServiceRequestConfirmedCOVNotification, ok := typ.(*BACnetConfirmedServiceRequestConfirmedCOVNotification); ok {
            return *sBACnetConfirmedServiceRequestConfirmedCOVNotification
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
    if len(m.LifetimeSeconds) > 0 {
        lengthInBits += 8 * uint16(len(m.LifetimeSeconds))
    }

    // Const Field (listOfValuesOpeningTag)
    lengthInBits += 8

    // Array field
    if len(m.Notifications) > 0 {
        for _, element := range m.Notifications {
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

func BACnetConfirmedServiceRequestConfirmedCOVNotificationParse(io *utils.ReadBuffer, len uint16) (BACnetConfirmedServiceRequestInitializer, error) {

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
    // Count array
    lifetimeSeconds := make([]int8, lifetimeLength)
    for curItem := uint16(0); curItem < uint16(lifetimeLength); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'lifetimeSeconds' field " + _err.Error())
        }
        lifetimeSeconds[curItem] = _item
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
    // Length array
    notifications := make([]IBACnetTagWithContent, 0)
    _notificationsLength := uint16(len) - uint16(uint16(18))
    _notificationsEndPos := io.GetPos() + uint16(_notificationsLength)
    for ;io.GetPos() < _notificationsEndPos; {
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

func (m BACnetConfirmedServiceRequestConfirmedCOVNotification) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Const Field (subscriberProcessIdentifierHeader)
    _subscriberProcessIdentifierHeaderErr := io.WriteUint8(8, 0x09)
    if _subscriberProcessIdentifierHeaderErr != nil {
        return errors.New("Error serializing 'subscriberProcessIdentifierHeader' field " + _subscriberProcessIdentifierHeaderErr.Error())
    }

    // Simple Field (subscriberProcessIdentifier)
    subscriberProcessIdentifier := uint8(m.SubscriberProcessIdentifier)
    _subscriberProcessIdentifierErr := io.WriteUint8(8, (subscriberProcessIdentifier))
    if _subscriberProcessIdentifierErr != nil {
        return errors.New("Error serializing 'subscriberProcessIdentifier' field " + _subscriberProcessIdentifierErr.Error())
    }

    // Const Field (monitoredObjectIdentifierHeader)
    _monitoredObjectIdentifierHeaderErr := io.WriteUint8(8, 0x1C)
    if _monitoredObjectIdentifierHeaderErr != nil {
        return errors.New("Error serializing 'monitoredObjectIdentifierHeader' field " + _monitoredObjectIdentifierHeaderErr.Error())
    }

    // Simple Field (monitoredObjectType)
    monitoredObjectType := uint16(m.MonitoredObjectType)
    _monitoredObjectTypeErr := io.WriteUint16(10, (monitoredObjectType))
    if _monitoredObjectTypeErr != nil {
        return errors.New("Error serializing 'monitoredObjectType' field " + _monitoredObjectTypeErr.Error())
    }

    // Simple Field (monitoredObjectInstanceNumber)
    monitoredObjectInstanceNumber := uint32(m.MonitoredObjectInstanceNumber)
    _monitoredObjectInstanceNumberErr := io.WriteUint32(22, (monitoredObjectInstanceNumber))
    if _monitoredObjectInstanceNumberErr != nil {
        return errors.New("Error serializing 'monitoredObjectInstanceNumber' field " + _monitoredObjectInstanceNumberErr.Error())
    }

    // Const Field (issueConfirmedNotificationsHeader)
    _issueConfirmedNotificationsHeaderErr := io.WriteUint8(8, 0x2C)
    if _issueConfirmedNotificationsHeaderErr != nil {
        return errors.New("Error serializing 'issueConfirmedNotificationsHeader' field " + _issueConfirmedNotificationsHeaderErr.Error())
    }

    // Simple Field (issueConfirmedNotificationsType)
    issueConfirmedNotificationsType := uint16(m.IssueConfirmedNotificationsType)
    _issueConfirmedNotificationsTypeErr := io.WriteUint16(10, (issueConfirmedNotificationsType))
    if _issueConfirmedNotificationsTypeErr != nil {
        return errors.New("Error serializing 'issueConfirmedNotificationsType' field " + _issueConfirmedNotificationsTypeErr.Error())
    }

    // Simple Field (issueConfirmedNotificationsInstanceNumber)
    issueConfirmedNotificationsInstanceNumber := uint32(m.IssueConfirmedNotificationsInstanceNumber)
    _issueConfirmedNotificationsInstanceNumberErr := io.WriteUint32(22, (issueConfirmedNotificationsInstanceNumber))
    if _issueConfirmedNotificationsInstanceNumberErr != nil {
        return errors.New("Error serializing 'issueConfirmedNotificationsInstanceNumber' field " + _issueConfirmedNotificationsInstanceNumberErr.Error())
    }

    // Const Field (lifetimeHeader)
    _lifetimeHeaderErr := io.WriteUint8(5, 0x07)
    if _lifetimeHeaderErr != nil {
        return errors.New("Error serializing 'lifetimeHeader' field " + _lifetimeHeaderErr.Error())
    }

    // Simple Field (lifetimeLength)
    lifetimeLength := uint8(m.LifetimeLength)
    _lifetimeLengthErr := io.WriteUint8(3, (lifetimeLength))
    if _lifetimeLengthErr != nil {
        return errors.New("Error serializing 'lifetimeLength' field " + _lifetimeLengthErr.Error())
    }

    // Array Field (lifetimeSeconds)
    if m.LifetimeSeconds != nil {
        for _, _element := range m.LifetimeSeconds {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'lifetimeSeconds' field " + _elementErr.Error())
            }
        }
    }

    // Const Field (listOfValuesOpeningTag)
    _listOfValuesOpeningTagErr := io.WriteUint8(8, 0x4E)
    if _listOfValuesOpeningTagErr != nil {
        return errors.New("Error serializing 'listOfValuesOpeningTag' field " + _listOfValuesOpeningTagErr.Error())
    }

    // Array Field (notifications)
    if m.Notifications != nil {
        for _, _element := range m.Notifications {
            _elementErr := _element.Serialize(io)
            if _elementErr != nil {
                return errors.New("Error serializing 'notifications' field " + _elementErr.Error())
            }
        }
    }

    // Const Field (listOfValuesClosingTag)
    _listOfValuesClosingTagErr := io.WriteUint8(8, 0x4F)
    if _listOfValuesClosingTagErr != nil {
        return errors.New("Error serializing 'listOfValuesClosingTag' field " + _listOfValuesClosingTagErr.Error())
    }

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
