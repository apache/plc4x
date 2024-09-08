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

// BACnetConstructedDataRecipientList is the corresponding interface of BACnetConstructedDataRecipientList
type BACnetConstructedDataRecipientList interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetRecipientList returns RecipientList (property field)
	GetRecipientList() []BACnetDestination
	// IsBACnetConstructedDataRecipientList is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataRecipientList()
}

// _BACnetConstructedDataRecipientList is the data-structure of this message
type _BACnetConstructedDataRecipientList struct {
	BACnetConstructedDataContract
	RecipientList []BACnetDestination
}

var _ BACnetConstructedDataRecipientList = (*_BACnetConstructedDataRecipientList)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataRecipientList)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataRecipientList) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataRecipientList) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_RECIPIENT_LIST
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataRecipientList) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataRecipientList) GetRecipientList() []BACnetDestination {
	return m.RecipientList
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataRecipientList factory function for _BACnetConstructedDataRecipientList
func NewBACnetConstructedDataRecipientList(recipientList []BACnetDestination, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataRecipientList {
	_result := &_BACnetConstructedDataRecipientList{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		RecipientList:                 recipientList,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataRecipientList(structType any) BACnetConstructedDataRecipientList {
	if casted, ok := structType.(BACnetConstructedDataRecipientList); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataRecipientList); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataRecipientList) GetTypeName() string {
	return "BACnetConstructedDataRecipientList"
}

func (m *_BACnetConstructedDataRecipientList) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).getLengthInBits(ctx))

	// Array field
	if len(m.RecipientList) > 0 {
		for _, element := range m.RecipientList {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_BACnetConstructedDataRecipientList) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataRecipientList) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataRecipientList BACnetConstructedDataRecipientList, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataRecipientList"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataRecipientList")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	recipientList, err := ReadTerminatedArrayField[BACnetDestination](ctx, "recipientList", ReadComplex[BACnetDestination](BACnetDestinationParseWithBuffer, readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'recipientList' field"))
	}
	m.RecipientList = recipientList

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataRecipientList"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataRecipientList")
	}

	return m, nil
}

func (m *_BACnetConstructedDataRecipientList) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataRecipientList) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataRecipientList"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataRecipientList")
		}

		if err := WriteComplexTypeArrayField(ctx, "recipientList", m.GetRecipientList(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'recipientList' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataRecipientList"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataRecipientList")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataRecipientList) IsBACnetConstructedDataRecipientList() {}

func (m *_BACnetConstructedDataRecipientList) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
