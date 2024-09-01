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

// BACnetConstructedDataAuthenticationStatus is the corresponding interface of BACnetConstructedDataAuthenticationStatus
type BACnetConstructedDataAuthenticationStatus interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetAuthenticationStatus returns AuthenticationStatus (property field)
	GetAuthenticationStatus() BACnetAuthenticationStatusTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetAuthenticationStatusTagged
}

// BACnetConstructedDataAuthenticationStatusExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataAuthenticationStatus.
// This is useful for switch cases.
type BACnetConstructedDataAuthenticationStatusExactly interface {
	BACnetConstructedDataAuthenticationStatus
	isBACnetConstructedDataAuthenticationStatus() bool
}

// _BACnetConstructedDataAuthenticationStatus is the data-structure of this message
type _BACnetConstructedDataAuthenticationStatus struct {
	*_BACnetConstructedData
	AuthenticationStatus BACnetAuthenticationStatusTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataAuthenticationStatus) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataAuthenticationStatus) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_AUTHENTICATION_STATUS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataAuthenticationStatus) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataAuthenticationStatus) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataAuthenticationStatus) GetAuthenticationStatus() BACnetAuthenticationStatusTagged {
	return m.AuthenticationStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataAuthenticationStatus) GetActualValue() BACnetAuthenticationStatusTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetAuthenticationStatusTagged(m.GetAuthenticationStatus())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataAuthenticationStatus factory function for _BACnetConstructedDataAuthenticationStatus
func NewBACnetConstructedDataAuthenticationStatus(authenticationStatus BACnetAuthenticationStatusTagged, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataAuthenticationStatus {
	_result := &_BACnetConstructedDataAuthenticationStatus{
		AuthenticationStatus:   authenticationStatus,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataAuthenticationStatus(structType any) BACnetConstructedDataAuthenticationStatus {
	if casted, ok := structType.(BACnetConstructedDataAuthenticationStatus); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataAuthenticationStatus); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataAuthenticationStatus) GetTypeName() string {
	return "BACnetConstructedDataAuthenticationStatus"
}

func (m *_BACnetConstructedDataAuthenticationStatus) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (authenticationStatus)
	lengthInBits += m.AuthenticationStatus.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataAuthenticationStatus) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConstructedDataAuthenticationStatusParse(ctx context.Context, theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataAuthenticationStatus, error) {
	return BACnetConstructedDataAuthenticationStatusParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataAuthenticationStatusParseWithBufferProducer(tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataAuthenticationStatus, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataAuthenticationStatus, error) {
		return BACnetConstructedDataAuthenticationStatusParseWithBuffer(ctx, readBuffer, tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
	}
}

func BACnetConstructedDataAuthenticationStatusParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataAuthenticationStatus, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataAuthenticationStatus"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataAuthenticationStatus")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	authenticationStatus, err := ReadSimpleField[BACnetAuthenticationStatusTagged](ctx, "authenticationStatus", ReadComplex[BACnetAuthenticationStatusTagged](BACnetAuthenticationStatusTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'authenticationStatus' field"))
	}

	actualValue, err := ReadVirtualField[BACnetAuthenticationStatusTagged](ctx, "actualValue", (*BACnetAuthenticationStatusTagged)(nil), authenticationStatus)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataAuthenticationStatus"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataAuthenticationStatus")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataAuthenticationStatus{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		AuthenticationStatus: authenticationStatus,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataAuthenticationStatus) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataAuthenticationStatus) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataAuthenticationStatus"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataAuthenticationStatus")
		}

		if err := WriteSimpleField[BACnetAuthenticationStatusTagged](ctx, "authenticationStatus", m.GetAuthenticationStatus(), WriteComplex[BACnetAuthenticationStatusTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'authenticationStatus' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataAuthenticationStatus"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataAuthenticationStatus")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataAuthenticationStatus) isBACnetConstructedDataAuthenticationStatus() bool {
	return true
}

func (m *_BACnetConstructedDataAuthenticationStatus) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
