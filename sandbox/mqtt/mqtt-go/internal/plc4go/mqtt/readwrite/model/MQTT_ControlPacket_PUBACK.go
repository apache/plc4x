/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model


import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

	// Code generated by code-generation. DO NOT EDIT.


// The data-structure of this message
type MQTT_ControlPacket_PUBACK struct {
	*MQTT_ControlPacket
	RemainingLength uint8
	PacketIdentifier uint16
	ReasonCode *MQTT_ReasonCode
	PropertyLength *uint32
	Properties []*MQTT_Property
}

// The corresponding interface
type IMQTT_ControlPacket_PUBACK interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *MQTT_ControlPacket_PUBACK) PacketType() MQTT_ControlPacketType {
	return MQTT_ControlPacketType_PUBACK
}


func (m *MQTT_ControlPacket_PUBACK) InitializeParent(parent *MQTT_ControlPacket) {
}

func NewMQTT_ControlPacket_PUBACK(remainingLength uint8, packetIdentifier uint16, reasonCode *MQTT_ReasonCode, propertyLength *uint32, properties []*MQTT_Property) *MQTT_ControlPacket {
	child := &MQTT_ControlPacket_PUBACK{
		RemainingLength: remainingLength,
		PacketIdentifier: packetIdentifier,
		ReasonCode: reasonCode,
		PropertyLength: propertyLength,
		Properties: properties,
    	MQTT_ControlPacket: NewMQTT_ControlPacket(),
	}
	child.Child = child
	return child.MQTT_ControlPacket
}

func CastMQTT_ControlPacket_PUBACK(structType interface{}) *MQTT_ControlPacket_PUBACK {
	castFunc := func(typ interface{}) *MQTT_ControlPacket_PUBACK {
		if casted, ok := typ.(MQTT_ControlPacket_PUBACK); ok {
			return &casted
		}
		if casted, ok := typ.(*MQTT_ControlPacket_PUBACK); ok {
			return casted
		}
		if casted, ok := typ.(MQTT_ControlPacket); ok {
			return CastMQTT_ControlPacket_PUBACK(casted.Child)
		}
		if casted, ok := typ.(*MQTT_ControlPacket); ok {
			return CastMQTT_ControlPacket_PUBACK(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *MQTT_ControlPacket_PUBACK) GetTypeName() string {
	return "MQTT_ControlPacket_PUBACK"
}

func (m *MQTT_ControlPacket_PUBACK) LengthInBits() uint16 {
	return m.LengthInBitsConditional(false)
}

func (m *MQTT_ControlPacket_PUBACK) LengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.ParentLengthInBits())

	// Reserved Field (reserved)
	lengthInBits += 4

	// Simple field (remainingLength)
	lengthInBits += 8;

	// Simple field (packetIdentifier)
	lengthInBits += 16;

	// Optional Field (reasonCode)
	if m.ReasonCode != nil {
		lengthInBits += 8
	}

	// Optional Field (propertyLength)
	if m.PropertyLength != nil {
		lengthInBits += 32
	}

	// Array field
	if len(m.Properties) > 0 {
		for _, element := range m.Properties {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}


func (m *MQTT_ControlPacket_PUBACK) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func MQTT_ControlPacket_PUBACKParse(readBuffer utils.ReadBuffer) (*MQTT_ControlPacket, error) {
	if pullErr := readBuffer.PullContext("MQTT_ControlPacket_PUBACK"); pullErr != nil {
		return nil, pullErr
	}
	var startPos = readBuffer.GetPos()
	var curPos uint16

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := readBuffer.ReadUint8("reserved", 4)
		if _err != nil {
			return nil, errors.Wrap(_err, "Error parsing 'reserved' field")
		}
		if reserved != uint8(0x0) {
			log.Info().Fields(map[string]interface{}{
				"expected value": uint8(0x0),
				"got value": reserved,
			}).Msg("Got unexpected response.")
		}
	}

	// Simple Field (remainingLength)
_remainingLength, _remainingLengthErr := readBuffer.ReadUint8("remainingLength", 8)
	if _remainingLengthErr != nil {
		return nil, errors.Wrap(_remainingLengthErr, "Error parsing 'remainingLength' field")
	}
	remainingLength := _remainingLength

	// Simple Field (packetIdentifier)
_packetIdentifier, _packetIdentifierErr := readBuffer.ReadUint16("packetIdentifier", 16)
	if _packetIdentifierErr != nil {
		return nil, errors.Wrap(_packetIdentifierErr, "Error parsing 'packetIdentifier' field")
	}
	packetIdentifier := _packetIdentifier

	// Optional Field (reasonCode) (Can be skipped, if a given expression evaluates to false)
	curPos = readBuffer.GetPos() - startPos
	var reasonCode *MQTT_ReasonCode = nil
	if bool(((remainingLength) - (curPos)) < ((3))) {
		if pullErr := readBuffer.PullContext("reasonCode"); pullErr != nil {
			return nil, pullErr
		}
		_val, _err := MQTT_ReasonCodeParse(readBuffer)
		if _err != nil {
			return nil, errors.Wrap(_err, "Error parsing 'reasonCode' field")
		}
		reasonCode = &_val
		if closeErr := readBuffer.CloseContext("reasonCode"); closeErr != nil {
			return nil, closeErr
		}
	}

	// Optional Field (propertyLength) (Can be skipped, if a given expression evaluates to false)
	curPos = readBuffer.GetPos() - startPos
	var propertyLength *uint32 = nil
	if bool(((remainingLength) - (curPos)) < ((4))) {
		_val, _err := readBuffer.ReadUint32("propertyLength", 32)
		if _err != nil {
			return nil, errors.Wrap(_err, "Error parsing 'propertyLength' field")
		}
		propertyLength = &_val
	}

	// Array field (properties)
	if pullErr := readBuffer.PullContext("properties", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, pullErr
	}
	// Length array
	properties := make([]*MQTT_Property, 0)
	{
		_propertiesLength := utils.InlineIf(bool(bool(((propertyLength)) != (nil))), func() interface{} {return uint16((*propertyLength))}, func() interface{} {return uint16(uint16(0))}).(uint16)
		_propertiesEndPos := readBuffer.GetPos() + uint16(_propertiesLength)
		for ;readBuffer.GetPos() < _propertiesEndPos; {
			_item, _err := MQTT_PropertyParse(readBuffer)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'properties' field")
			}
			properties = append(properties, _item)
		}
	}
	if closeErr := readBuffer.CloseContext("properties", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, closeErr
	}

	if closeErr := readBuffer.CloseContext("MQTT_ControlPacket_PUBACK"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &MQTT_ControlPacket_PUBACK{
		RemainingLength: remainingLength,
		PacketIdentifier: packetIdentifier,
		ReasonCode: reasonCode,
		PropertyLength: propertyLength,
		Properties: properties,
        MQTT_ControlPacket: &MQTT_ControlPacket{},
	}
	_child.MQTT_ControlPacket.Child = _child
	return _child.MQTT_ControlPacket, nil
}

func (m *MQTT_ControlPacket_PUBACK) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MQTT_ControlPacket_PUBACK"); pushErr != nil {
			return pushErr
		}

	// Reserved Field (reserved)
	{
		_err := writeBuffer.WriteUint8("reserved", 4, uint8(0x0))
		if _err != nil {
			return errors.Wrap(_err, "Error serializing 'reserved' field")
		}
	}

	// Simple Field (remainingLength)
	remainingLength := uint8(m.RemainingLength)
	_remainingLengthErr := writeBuffer.WriteUint8("remainingLength", 8, (remainingLength))
	if _remainingLengthErr != nil {
		return errors.Wrap(_remainingLengthErr, "Error serializing 'remainingLength' field")
	}

	// Simple Field (packetIdentifier)
	packetIdentifier := uint16(m.PacketIdentifier)
	_packetIdentifierErr := writeBuffer.WriteUint16("packetIdentifier", 16, (packetIdentifier))
	if _packetIdentifierErr != nil {
		return errors.Wrap(_packetIdentifierErr, "Error serializing 'packetIdentifier' field")
	}

	// Optional Field (reasonCode) (Can be skipped, if the value is null)
	var reasonCode *MQTT_ReasonCode = nil
	if m.ReasonCode != nil {
		if pushErr := writeBuffer.PushContext("reasonCode"); pushErr != nil {
			return pushErr
		}
		reasonCode = m.ReasonCode
		_reasonCodeErr := reasonCode.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("reasonCode"); popErr != nil {
			return popErr
		}
		if _reasonCodeErr != nil {
			return errors.Wrap(_reasonCodeErr, "Error serializing 'reasonCode' field")
		}
	}

	// Optional Field (propertyLength) (Can be skipped, if the value is null)
	var propertyLength *uint32 = nil
	if m.PropertyLength != nil {
		propertyLength = m.PropertyLength
		_propertyLengthErr := writeBuffer.WriteUint32("propertyLength", 32, *(propertyLength))
		if _propertyLengthErr != nil {
			return errors.Wrap(_propertyLengthErr, "Error serializing 'propertyLength' field")
		}
	}

	// Array Field (properties)
	if m.Properties != nil {
		if pushErr := writeBuffer.PushContext("properties", utils.WithRenderAsList(true)); pushErr != nil {
			return pushErr
		}
		for _, _element := range m.Properties {
			_elementErr := _element.Serialize(writeBuffer)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'properties' field")
			}
		}
		if popErr := writeBuffer.PopContext("properties", utils.WithRenderAsList(true)); popErr != nil {
			return popErr
		}
	}

		if popErr := writeBuffer.PopContext("MQTT_ControlPacket_PUBACK"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *MQTT_ControlPacket_PUBACK) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	m.Serialize(buffer)
	return buffer.GetBox().String()
}



