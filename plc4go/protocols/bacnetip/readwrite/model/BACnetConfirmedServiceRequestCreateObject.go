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

// BACnetConfirmedServiceRequestCreateObject is the corresponding interface of BACnetConfirmedServiceRequestCreateObject
type BACnetConfirmedServiceRequestCreateObject interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConfirmedServiceRequest
	// GetObjectSpecifier returns ObjectSpecifier (property field)
	GetObjectSpecifier() BACnetConfirmedServiceRequestCreateObjectObjectSpecifier
	// GetListOfValues returns ListOfValues (property field)
	GetListOfValues() BACnetPropertyValues
}

// BACnetConfirmedServiceRequestCreateObjectExactly can be used when we want exactly this type and not a type which fulfills BACnetConfirmedServiceRequestCreateObject.
// This is useful for switch cases.
type BACnetConfirmedServiceRequestCreateObjectExactly interface {
	BACnetConfirmedServiceRequestCreateObject
	isBACnetConfirmedServiceRequestCreateObject() bool
}

// _BACnetConfirmedServiceRequestCreateObject is the data-structure of this message
type _BACnetConfirmedServiceRequestCreateObject struct {
	*_BACnetConfirmedServiceRequest
	ObjectSpecifier BACnetConfirmedServiceRequestCreateObjectObjectSpecifier
	ListOfValues    BACnetPropertyValues
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConfirmedServiceRequestCreateObject) GetServiceChoice() BACnetConfirmedServiceChoice {
	return BACnetConfirmedServiceChoice_CREATE_OBJECT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConfirmedServiceRequestCreateObject) InitializeParent(parent BACnetConfirmedServiceRequest) {
}

func (m *_BACnetConfirmedServiceRequestCreateObject) GetParent() BACnetConfirmedServiceRequest {
	return m._BACnetConfirmedServiceRequest
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestCreateObject) GetObjectSpecifier() BACnetConfirmedServiceRequestCreateObjectObjectSpecifier {
	return m.ObjectSpecifier
}

func (m *_BACnetConfirmedServiceRequestCreateObject) GetListOfValues() BACnetPropertyValues {
	return m.ListOfValues
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConfirmedServiceRequestCreateObject factory function for _BACnetConfirmedServiceRequestCreateObject
func NewBACnetConfirmedServiceRequestCreateObject(objectSpecifier BACnetConfirmedServiceRequestCreateObjectObjectSpecifier, listOfValues BACnetPropertyValues, serviceRequestLength uint32) *_BACnetConfirmedServiceRequestCreateObject {
	_result := &_BACnetConfirmedServiceRequestCreateObject{
		ObjectSpecifier:                objectSpecifier,
		ListOfValues:                   listOfValues,
		_BACnetConfirmedServiceRequest: NewBACnetConfirmedServiceRequest(serviceRequestLength),
	}
	_result._BACnetConfirmedServiceRequest._BACnetConfirmedServiceRequestChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConfirmedServiceRequestCreateObject(structType any) BACnetConfirmedServiceRequestCreateObject {
	if casted, ok := structType.(BACnetConfirmedServiceRequestCreateObject); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestCreateObject); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestCreateObject) GetTypeName() string {
	return "BACnetConfirmedServiceRequestCreateObject"
}

func (m *_BACnetConfirmedServiceRequestCreateObject) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (objectSpecifier)
	lengthInBits += m.ObjectSpecifier.GetLengthInBits(ctx)

	// Optional Field (listOfValues)
	if m.ListOfValues != nil {
		lengthInBits += m.ListOfValues.GetLengthInBits(ctx)
	}

	return lengthInBits
}

func (m *_BACnetConfirmedServiceRequestCreateObject) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConfirmedServiceRequestCreateObjectParse(ctx context.Context, theBytes []byte, serviceRequestLength uint32) (BACnetConfirmedServiceRequestCreateObject, error) {
	return BACnetConfirmedServiceRequestCreateObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), serviceRequestLength)
}

func BACnetConfirmedServiceRequestCreateObjectParseWithBufferProducer(serviceRequestLength uint32) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestCreateObject, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConfirmedServiceRequestCreateObject, error) {
		return BACnetConfirmedServiceRequestCreateObjectParseWithBuffer(ctx, readBuffer, serviceRequestLength)
	}
}

func BACnetConfirmedServiceRequestCreateObjectParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, serviceRequestLength uint32) (BACnetConfirmedServiceRequestCreateObject, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestCreateObject"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConfirmedServiceRequestCreateObject")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	objectSpecifier, err := ReadSimpleField[BACnetConfirmedServiceRequestCreateObjectObjectSpecifier](ctx, "objectSpecifier", ReadComplex[BACnetConfirmedServiceRequestCreateObjectObjectSpecifier](BACnetConfirmedServiceRequestCreateObjectObjectSpecifierParseWithBufferProducer((uint8)(uint8(0))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'objectSpecifier' field"))
	}

	_listOfValues, err := ReadOptionalField[BACnetPropertyValues](ctx, "listOfValues", ReadComplex[BACnetPropertyValues](BACnetPropertyValuesParseWithBufferProducer((uint8)(uint8(1)), (BACnetObjectType)(CastBACnetObjectType(utils.InlineIf(objectSpecifier.GetIsObjectType(), func() any { return CastBACnetObjectType(objectSpecifier.GetObjectType()) }, func() any { return CastBACnetObjectType(objectSpecifier.GetObjectIdentifier().GetObjectType()) })))), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'listOfValues' field"))
	}
	var listOfValues BACnetPropertyValues
	if _listOfValues != nil {
		listOfValues = *_listOfValues
	}

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestCreateObject"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConfirmedServiceRequestCreateObject")
	}

	// Create a partially initialized instance
	_child := &_BACnetConfirmedServiceRequestCreateObject{
		_BACnetConfirmedServiceRequest: &_BACnetConfirmedServiceRequest{
			ServiceRequestLength: serviceRequestLength,
		},
		ObjectSpecifier: objectSpecifier,
		ListOfValues:    listOfValues,
	}
	_child._BACnetConfirmedServiceRequest._BACnetConfirmedServiceRequestChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConfirmedServiceRequestCreateObject) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConfirmedServiceRequestCreateObject) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestCreateObject"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConfirmedServiceRequestCreateObject")
		}

		if err := WriteSimpleField[BACnetConfirmedServiceRequestCreateObjectObjectSpecifier](ctx, "objectSpecifier", m.GetObjectSpecifier(), WriteComplex[BACnetConfirmedServiceRequestCreateObjectObjectSpecifier](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'objectSpecifier' field")
		}

		if err := WriteOptionalField[BACnetPropertyValues](ctx, "listOfValues", GetRef(m.GetListOfValues()), WriteComplex[BACnetPropertyValues](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'listOfValues' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestCreateObject"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConfirmedServiceRequestCreateObject")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConfirmedServiceRequestCreateObject) isBACnetConfirmedServiceRequestCreateObject() bool {
	return true
}

func (m *_BACnetConfirmedServiceRequestCreateObject) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
