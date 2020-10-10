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
	"strconv"
)

// Constant values.
const BACnetConfirmedServiceRequestSubscribeCOV_SUBSCRIBERPROCESSIDENTIFIERHEADER uint8 = 0x09
const BACnetConfirmedServiceRequestSubscribeCOV_MONITOREDOBJECTIDENTIFIERHEADER uint8 = 0x1C
const BACnetConfirmedServiceRequestSubscribeCOV_ISSUECONFIRMEDNOTIFICATIONSHEADER uint8 = 0x29
const BACnetConfirmedServiceRequestSubscribeCOV_ISSUECONFIRMEDNOTIFICATIONSSKIPBITS uint8 = 0x00
const BACnetConfirmedServiceRequestSubscribeCOV_LIFETIMEHEADER uint8 = 0x07

// The data-structure of this message
type BACnetConfirmedServiceRequestSubscribeCOV struct {
	subscriberProcessIdentifier   uint8
	monitoredObjectType           uint16
	monitoredObjectInstanceNumber uint32
	issueConfirmedNotifications   bool
	lifetimeLength                uint8
	lifetimeSeconds               []int8
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestSubscribeCOV interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestSubscribeCOV) ServiceChoice() uint8 {
	return 0x05
}

func (m BACnetConfirmedServiceRequestSubscribeCOV) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestSubscribeCOV(subscriberProcessIdentifier uint8, monitoredObjectType uint16, monitoredObjectInstanceNumber uint32, issueConfirmedNotifications bool, lifetimeLength uint8, lifetimeSeconds []int8) BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestSubscribeCOV{subscriberProcessIdentifier: subscriberProcessIdentifier, monitoredObjectType: monitoredObjectType, monitoredObjectInstanceNumber: monitoredObjectInstanceNumber, issueConfirmedNotifications: issueConfirmedNotifications, lifetimeLength: lifetimeLength, lifetimeSeconds: lifetimeSeconds}
}

func (m BACnetConfirmedServiceRequestSubscribeCOV) LengthInBits() uint16 {
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

	// Const Field (issueConfirmedNotificationsSkipBits)
	lengthInBits += 7

	// Simple field (issueConfirmedNotifications)
	lengthInBits += 1

	// Const Field (lifetimeHeader)
	lengthInBits += 5

	// Simple field (lifetimeLength)
	lengthInBits += 3

	// Array field
	if len(m.lifetimeSeconds) > 0 {
		lengthInBits += 8 * uint16(len(m.lifetimeSeconds))
	}

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestSubscribeCOV) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestSubscribeCOVParse(io spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Const Field (subscriberProcessIdentifierHeader)
	var subscriberProcessIdentifierHeader uint8 = io.ReadUint8(8)
	if subscriberProcessIdentifierHeader != BACnetConfirmedServiceRequestSubscribeCOV_SUBSCRIBERPROCESSIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestSubscribeCOV_SUBSCRIBERPROCESSIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(subscriberProcessIdentifierHeader)))
	}

	// Simple Field (subscriberProcessIdentifier)
	var subscriberProcessIdentifier uint8 = io.ReadUint8(8)

	// Const Field (monitoredObjectIdentifierHeader)
	var monitoredObjectIdentifierHeader uint8 = io.ReadUint8(8)
	if monitoredObjectIdentifierHeader != BACnetConfirmedServiceRequestSubscribeCOV_MONITOREDOBJECTIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestSubscribeCOV_MONITOREDOBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(monitoredObjectIdentifierHeader)))
	}

	// Simple Field (monitoredObjectType)
	var monitoredObjectType uint16 = io.ReadUint16(10)

	// Simple Field (monitoredObjectInstanceNumber)
	var monitoredObjectInstanceNumber uint32 = io.ReadUint32(22)

	// Const Field (issueConfirmedNotificationsHeader)
	var issueConfirmedNotificationsHeader uint8 = io.ReadUint8(8)
	if issueConfirmedNotificationsHeader != BACnetConfirmedServiceRequestSubscribeCOV_ISSUECONFIRMEDNOTIFICATIONSHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestSubscribeCOV_ISSUECONFIRMEDNOTIFICATIONSHEADER)) + " but got " + strconv.Itoa(int(issueConfirmedNotificationsHeader)))
	}

	// Const Field (issueConfirmedNotificationsSkipBits)
	var issueConfirmedNotificationsSkipBits uint8 = io.ReadUint8(7)
	if issueConfirmedNotificationsSkipBits != BACnetConfirmedServiceRequestSubscribeCOV_ISSUECONFIRMEDNOTIFICATIONSSKIPBITS {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestSubscribeCOV_ISSUECONFIRMEDNOTIFICATIONSSKIPBITS)) + " but got " + strconv.Itoa(int(issueConfirmedNotificationsSkipBits)))
	}

	// Simple Field (issueConfirmedNotifications)
	var issueConfirmedNotifications bool = io.ReadBit()

	// Const Field (lifetimeHeader)
	var lifetimeHeader uint8 = io.ReadUint8(5)
	if lifetimeHeader != BACnetConfirmedServiceRequestSubscribeCOV_LIFETIMEHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestSubscribeCOV_LIFETIMEHEADER)) + " but got " + strconv.Itoa(int(lifetimeHeader)))
	}

	// Simple Field (lifetimeLength)
	var lifetimeLength uint8 = io.ReadUint8(3)

	// Array field (lifetimeSeconds)
	var lifetimeSeconds []int8
	// Count array
	{
		lifetimeSeconds := make([]int8, lifetimeLength)
		for curItem := uint16(0); curItem < uint16(lifetimeLength); curItem++ {

			lifetimeSeconds = append(lifetimeSeconds, io.ReadInt8(8))
		}
	}

	// Create the instance
	return NewBACnetConfirmedServiceRequestSubscribeCOV(subscriberProcessIdentifier, monitoredObjectType, monitoredObjectInstanceNumber, issueConfirmedNotifications, lifetimeLength, lifetimeSeconds), nil
}

func (m BACnetConfirmedServiceRequestSubscribeCOV) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IBACnetConfirmedServiceRequestSubscribeCOV); ok {

			// Const Field (subscriberProcessIdentifierHeader)
			io.WriteUint8(8, 0x09)

			// Simple Field (subscriberProcessIdentifier)
			var subscriberProcessIdentifier uint8 = m.subscriberProcessIdentifier
			io.WriteUint8(8, (subscriberProcessIdentifier))

			// Const Field (monitoredObjectIdentifierHeader)
			io.WriteUint8(8, 0x1C)

			// Simple Field (monitoredObjectType)
			var monitoredObjectType uint16 = m.monitoredObjectType
			io.WriteUint16(10, (monitoredObjectType))

			// Simple Field (monitoredObjectInstanceNumber)
			var monitoredObjectInstanceNumber uint32 = m.monitoredObjectInstanceNumber
			io.WriteUint32(22, (monitoredObjectInstanceNumber))

			// Const Field (issueConfirmedNotificationsHeader)
			io.WriteUint8(8, 0x29)

			// Const Field (issueConfirmedNotificationsSkipBits)
			io.WriteUint8(7, 0x00)

			// Simple Field (issueConfirmedNotifications)
			var issueConfirmedNotifications bool = m.issueConfirmedNotifications
			io.WriteBit((bool)(issueConfirmedNotifications))

			// Const Field (lifetimeHeader)
			io.WriteUint8(5, 0x07)

			// Simple Field (lifetimeLength)
			var lifetimeLength uint8 = m.lifetimeLength
			io.WriteUint8(3, (lifetimeLength))

			// Array Field (lifetimeSeconds)
			if m.lifetimeSeconds != nil {
				for _, _element := range m.lifetimeSeconds {
					io.WriteInt8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
