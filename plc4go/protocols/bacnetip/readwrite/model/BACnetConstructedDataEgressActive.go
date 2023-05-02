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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataEgressActive is the corresponding interface of BACnetConstructedDataEgressActive
type BACnetConstructedDataEgressActive interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetEgressActive returns EgressActive (property field)
	GetEgressActive() BACnetApplicationTagBoolean
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagBoolean
}

// BACnetConstructedDataEgressActiveExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataEgressActive.
// This is useful for switch cases.
type BACnetConstructedDataEgressActiveExactly interface {
	BACnetConstructedDataEgressActive
	isBACnetConstructedDataEgressActive() bool
}

// _BACnetConstructedDataEgressActive is the data-structure of this message
type _BACnetConstructedDataEgressActive struct {
	*_BACnetConstructedData
	EgressActive BACnetApplicationTagBoolean
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataEgressActive) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataEgressActive) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_EGRESS_ACTIVE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataEgressActive) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataEgressActive) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataEgressActive) GetEgressActive() BACnetApplicationTagBoolean {
	return m.EgressActive
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataEgressActive) GetActualValue() BACnetApplicationTagBoolean {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagBoolean(m.GetEgressActive())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataEgressActive factory function for _BACnetConstructedDataEgressActive
func NewBACnetConstructedDataEgressActive(egressActive BACnetApplicationTagBoolean, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataEgressActive {
	_result := &_BACnetConstructedDataEgressActive{
		EgressActive:           egressActive,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataEgressActive(structType any) BACnetConstructedDataEgressActive {
	if casted, ok := structType.(BACnetConstructedDataEgressActive); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataEgressActive); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataEgressActive) GetTypeName() string {
	return "BACnetConstructedDataEgressActive"
}

func (m *_BACnetConstructedDataEgressActive) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (egressActive)
	lengthInBits += m.EgressActive.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataEgressActive) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConstructedDataEgressActiveParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataEgressActive, error) {
	return BACnetConstructedDataEgressActiveParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataEgressActiveParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataEgressActive, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataEgressActive"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataEgressActive")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (egressActive)
	if pullErr := readBuffer.PullContext("egressActive"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for egressActive")
	}
	_egressActive, _egressActiveErr := BACnetApplicationTagParseWithBuffer(ctx, readBuffer)
	if _egressActiveErr != nil {
		return nil, errors.Wrap(_egressActiveErr, "Error parsing 'egressActive' field of BACnetConstructedDataEgressActive")
	}
	egressActive := _egressActive.(BACnetApplicationTagBoolean)
	if closeErr := readBuffer.CloseContext("egressActive"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for egressActive")
	}

	// Virtual field
	_actualValue := egressActive
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataEgressActive"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataEgressActive")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataEgressActive{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		EgressActive: egressActive,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataEgressActive) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataEgressActive) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataEgressActive"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataEgressActive")
		}

		// Simple Field (egressActive)
		if pushErr := writeBuffer.PushContext("egressActive"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for egressActive")
		}
		_egressActiveErr := writeBuffer.WriteSerializable(ctx, m.GetEgressActive())
		if popErr := writeBuffer.PopContext("egressActive"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for egressActive")
		}
		if _egressActiveErr != nil {
			return errors.Wrap(_egressActiveErr, "Error serializing 'egressActive' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataEgressActive"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataEgressActive")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataEgressActive) isBACnetConstructedDataEgressActive() bool {
	return true
}

func (m *_BACnetConstructedDataEgressActive) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
