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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataCOVResubscriptionInterval is the corresponding interface of BACnetConstructedDataCOVResubscriptionInterval
type BACnetConstructedDataCOVResubscriptionInterval interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetCovResubscriptionInterval returns CovResubscriptionInterval (property field)
	GetCovResubscriptionInterval() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
}

// BACnetConstructedDataCOVResubscriptionIntervalExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataCOVResubscriptionInterval.
// This is useful for switch cases.
type BACnetConstructedDataCOVResubscriptionIntervalExactly interface {
	BACnetConstructedDataCOVResubscriptionInterval
	isBACnetConstructedDataCOVResubscriptionInterval() bool
}

// _BACnetConstructedDataCOVResubscriptionInterval is the data-structure of this message
type _BACnetConstructedDataCOVResubscriptionInterval struct {
	*_BACnetConstructedData
	CovResubscriptionInterval BACnetApplicationTagUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_COV_RESUBSCRIPTION_INTERVAL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataCOVResubscriptionInterval) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetCovResubscriptionInterval() BACnetApplicationTagUnsignedInteger {
	return m.CovResubscriptionInterval
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetCovResubscriptionInterval())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataCOVResubscriptionInterval factory function for _BACnetConstructedDataCOVResubscriptionInterval
func NewBACnetConstructedDataCOVResubscriptionInterval(covResubscriptionInterval BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataCOVResubscriptionInterval {
	_result := &_BACnetConstructedDataCOVResubscriptionInterval{
		CovResubscriptionInterval: covResubscriptionInterval,
		_BACnetConstructedData:    NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataCOVResubscriptionInterval(structType interface{}) BACnetConstructedDataCOVResubscriptionInterval {
	if casted, ok := structType.(BACnetConstructedDataCOVResubscriptionInterval); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataCOVResubscriptionInterval); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetTypeName() string {
	return "BACnetConstructedDataCOVResubscriptionInterval"
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (covResubscriptionInterval)
	lengthInBits += m.CovResubscriptionInterval.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataCOVResubscriptionIntervalParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataCOVResubscriptionInterval, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataCOVResubscriptionInterval"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataCOVResubscriptionInterval")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (covResubscriptionInterval)
	if pullErr := readBuffer.PullContext("covResubscriptionInterval"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for covResubscriptionInterval")
	}
	_covResubscriptionInterval, _covResubscriptionIntervalErr := BACnetApplicationTagParse(readBuffer)
	if _covResubscriptionIntervalErr != nil {
		return nil, errors.Wrap(_covResubscriptionIntervalErr, "Error parsing 'covResubscriptionInterval' field of BACnetConstructedDataCOVResubscriptionInterval")
	}
	covResubscriptionInterval := _covResubscriptionInterval.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("covResubscriptionInterval"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for covResubscriptionInterval")
	}

	// Virtual field
	_actualValue := covResubscriptionInterval
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataCOVResubscriptionInterval"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataCOVResubscriptionInterval")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataCOVResubscriptionInterval{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		CovResubscriptionInterval: covResubscriptionInterval,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataCOVResubscriptionInterval"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataCOVResubscriptionInterval")
		}

		// Simple Field (covResubscriptionInterval)
		if pushErr := writeBuffer.PushContext("covResubscriptionInterval"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for covResubscriptionInterval")
		}
		_covResubscriptionIntervalErr := writeBuffer.WriteSerializable(m.GetCovResubscriptionInterval())
		if popErr := writeBuffer.PopContext("covResubscriptionInterval"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for covResubscriptionInterval")
		}
		if _covResubscriptionIntervalErr != nil {
			return errors.Wrap(_covResubscriptionIntervalErr, "Error serializing 'covResubscriptionInterval' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataCOVResubscriptionInterval"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataCOVResubscriptionInterval")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) isBACnetConstructedDataCOVResubscriptionInterval() bool {
	return true
}

func (m *_BACnetConstructedDataCOVResubscriptionInterval) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
