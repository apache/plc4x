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

package model

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConfirmedServiceRequestAddListElement is the corresponding interface of BACnetConfirmedServiceRequestAddListElement
type BACnetConfirmedServiceRequestAddListElement interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConfirmedServiceRequest
	// GetObjectIdentifier returns ObjectIdentifier (property field)
	GetObjectIdentifier() BACnetContextTagObjectIdentifier
	// GetPropertyIdentifier returns PropertyIdentifier (property field)
	GetPropertyIdentifier() BACnetPropertyIdentifierTagged
	// GetArrayIndex returns ArrayIndex (property field)
	GetArrayIndex() BACnetContextTagUnsignedInteger
	// GetListOfElements returns ListOfElements (property field)
	GetListOfElements() BACnetConstructedData
}

// BACnetConfirmedServiceRequestAddListElementExactly can be used when we want exactly this type and not a type which fulfills BACnetConfirmedServiceRequestAddListElement.
// This is useful for switch cases.
type BACnetConfirmedServiceRequestAddListElementExactly interface {
	BACnetConfirmedServiceRequestAddListElement
	isBACnetConfirmedServiceRequestAddListElement() bool
}

// _BACnetConfirmedServiceRequestAddListElement is the data-structure of this message
type _BACnetConfirmedServiceRequestAddListElement struct {
	*_BACnetConfirmedServiceRequest
	ObjectIdentifier   BACnetContextTagObjectIdentifier
	PropertyIdentifier BACnetPropertyIdentifierTagged
	ArrayIndex         BACnetContextTagUnsignedInteger
	ListOfElements     BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConfirmedServiceRequestAddListElement) GetServiceChoice() BACnetConfirmedServiceChoice {
	return BACnetConfirmedServiceChoice_ADD_LIST_ELEMENT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConfirmedServiceRequestAddListElement) InitializeParent(parent BACnetConfirmedServiceRequest) {
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetParent() BACnetConfirmedServiceRequest {
	return m._BACnetConfirmedServiceRequest
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestAddListElement) GetObjectIdentifier() BACnetContextTagObjectIdentifier {
	return m.ObjectIdentifier
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetPropertyIdentifier() BACnetPropertyIdentifierTagged {
	return m.PropertyIdentifier
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetArrayIndex() BACnetContextTagUnsignedInteger {
	return m.ArrayIndex
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetListOfElements() BACnetConstructedData {
	return m.ListOfElements
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConfirmedServiceRequestAddListElement factory function for _BACnetConfirmedServiceRequestAddListElement
func NewBACnetConfirmedServiceRequestAddListElement(objectIdentifier BACnetContextTagObjectIdentifier, propertyIdentifier BACnetPropertyIdentifierTagged, arrayIndex BACnetContextTagUnsignedInteger, listOfElements BACnetConstructedData, serviceRequestLength uint32) *_BACnetConfirmedServiceRequestAddListElement {
	_result := &_BACnetConfirmedServiceRequestAddListElement{
		ObjectIdentifier:               objectIdentifier,
		PropertyIdentifier:             propertyIdentifier,
		ArrayIndex:                     arrayIndex,
		ListOfElements:                 listOfElements,
		_BACnetConfirmedServiceRequest: NewBACnetConfirmedServiceRequest(serviceRequestLength),
	}
	_result._BACnetConfirmedServiceRequest._BACnetConfirmedServiceRequestChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConfirmedServiceRequestAddListElement(structType any) BACnetConfirmedServiceRequestAddListElement {
	if casted, ok := structType.(BACnetConfirmedServiceRequestAddListElement); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestAddListElement); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetTypeName() string {
	return "BACnetConfirmedServiceRequestAddListElement"
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (objectIdentifier)
	lengthInBits += m.ObjectIdentifier.GetLengthInBits(ctx)

	// Simple field (propertyIdentifier)
	lengthInBits += m.PropertyIdentifier.GetLengthInBits(ctx)

	// Optional Field (arrayIndex)
	if m.ArrayIndex != nil {
		lengthInBits += m.ArrayIndex.GetLengthInBits(ctx)
	}

	// Optional Field (listOfElements)
	if m.ListOfElements != nil {
		lengthInBits += m.ListOfElements.GetLengthInBits(ctx)
	}

	return lengthInBits
}

func (m *_BACnetConfirmedServiceRequestAddListElement) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConfirmedServiceRequestAddListElementParse(ctx context.Context, theBytes []byte, serviceRequestLength uint32) (BACnetConfirmedServiceRequestAddListElement, error) {
	return BACnetConfirmedServiceRequestAddListElementParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), serviceRequestLength)
}

func BACnetConfirmedServiceRequestAddListElementParseWithBufferProducer(serviceRequestLength uint32) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestAddListElement, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestAddListElement, error) {
		return BACnetConfirmedServiceRequestAddListElementParseWithBuffer(ctx, readBuffer, serviceRequestLength)
	}
}

func BACnetConfirmedServiceRequestAddListElementParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, serviceRequestLength uint32) (BACnetConfirmedServiceRequestAddListElement, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestAddListElement"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConfirmedServiceRequestAddListElement")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	objectIdentifier, err := ReadSimpleField[BACnetContextTagObjectIdentifier](ctx, "objectIdentifier", ReadComplex[BACnetContextTagObjectIdentifier](BACnetContextTagParseWithBufferProducer[BACnetContextTagObjectIdentifier]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_BACNET_OBJECT_IDENTIFIER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'objectIdentifier' field"))
	}

	propertyIdentifier, err := ReadSimpleField[BACnetPropertyIdentifierTagged](ctx, "propertyIdentifier", ReadComplex[BACnetPropertyIdentifierTagged](BACnetPropertyIdentifierTaggedParseWithBufferProducer((uint8)(uint8(1)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'propertyIdentifier' field"))
	}

	_arrayIndex, err := ReadOptionalField[BACnetContextTagUnsignedInteger](ctx, "arrayIndex", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(2)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'arrayIndex' field"))
	}
	var arrayIndex BACnetContextTagUnsignedInteger
	if _arrayIndex != nil {
		arrayIndex = *_arrayIndex
	}

	_listOfElements, err := ReadOptionalField[BACnetConstructedData](ctx, "listOfElements", ReadComplex[BACnetConstructedData](BACnetConstructedDataParseWithBufferProducer[BACnetConstructedData]((uint8)(uint8(3)), (BACnetObjectType)(objectIdentifier.GetObjectType()), (BACnetPropertyIdentifier)(propertyIdentifier.GetValue()), (BACnetTagPayloadUnsignedInteger)((CastBACnetTagPayloadUnsignedInteger(utils.InlineIf(bool((arrayIndex) != (nil)), func() any { return CastBACnetTagPayloadUnsignedInteger((arrayIndex).GetPayload()) }, func() any { return CastBACnetTagPayloadUnsignedInteger(nil) }))))), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'listOfElements' field"))
	}
	var listOfElements BACnetConstructedData
	if _listOfElements != nil {
		listOfElements = *_listOfElements
	}

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestAddListElement"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConfirmedServiceRequestAddListElement")
	}

	// Create a partially initialized instance
	_child := &_BACnetConfirmedServiceRequestAddListElement{
		_BACnetConfirmedServiceRequest: &_BACnetConfirmedServiceRequest{
			ServiceRequestLength: serviceRequestLength,
		},
		ObjectIdentifier:   objectIdentifier,
		PropertyIdentifier: propertyIdentifier,
		ArrayIndex:         arrayIndex,
		ListOfElements:     listOfElements,
	}
	_child._BACnetConfirmedServiceRequest._BACnetConfirmedServiceRequestChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConfirmedServiceRequestAddListElement) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConfirmedServiceRequestAddListElement) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestAddListElement"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConfirmedServiceRequestAddListElement")
		}

		if err := WriteSimpleField[BACnetContextTagObjectIdentifier](ctx, "objectIdentifier", m.GetObjectIdentifier(), WriteComplex[BACnetContextTagObjectIdentifier](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'objectIdentifier' field")
		}

		if err := WriteSimpleField[BACnetPropertyIdentifierTagged](ctx, "propertyIdentifier", m.GetPropertyIdentifier(), WriteComplex[BACnetPropertyIdentifierTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'propertyIdentifier' field")
		}

		if err := WriteOptionalField[BACnetContextTagUnsignedInteger](ctx, "arrayIndex", GetRef(m.GetArrayIndex()), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'arrayIndex' field")
		}

		if err := WriteOptionalField[BACnetConstructedData](ctx, "listOfElements", GetRef(m.GetListOfElements()), WriteComplex[BACnetConstructedData](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'listOfElements' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestAddListElement"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConfirmedServiceRequestAddListElement")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConfirmedServiceRequestAddListElement) isBACnetConfirmedServiceRequestAddListElement() bool {
	return true
}

func (m *_BACnetConfirmedServiceRequestAddListElement) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
