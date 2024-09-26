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

// ErrorReportingDataGeneric is the corresponding interface of ErrorReportingDataGeneric
type ErrorReportingDataGeneric interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ErrorReportingData
	// GetSystemCategory returns SystemCategory (property field)
	GetSystemCategory() ErrorReportingSystemCategory
	// GetMostRecent returns MostRecent (property field)
	GetMostRecent() bool
	// GetAcknowledge returns Acknowledge (property field)
	GetAcknowledge() bool
	// GetMostSevere returns MostSevere (property field)
	GetMostSevere() bool
	// GetSeverity returns Severity (property field)
	GetSeverity() ErrorReportingSeverity
	// GetDeviceId returns DeviceId (property field)
	GetDeviceId() uint8
	// GetErrorData1 returns ErrorData1 (property field)
	GetErrorData1() uint8
	// GetErrorData2 returns ErrorData2 (property field)
	GetErrorData2() uint8
	// GetIsMostSevereError returns IsMostSevereError (virtual field)
	GetIsMostSevereError() bool
	// GetIsMostRecentError returns IsMostRecentError (virtual field)
	GetIsMostRecentError() bool
	// GetIsMostRecentAndMostSevere returns IsMostRecentAndMostSevere (virtual field)
	GetIsMostRecentAndMostSevere() bool
	// IsErrorReportingDataGeneric is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsErrorReportingDataGeneric()
	// CreateBuilder creates a ErrorReportingDataGenericBuilder
	CreateErrorReportingDataGenericBuilder() ErrorReportingDataGenericBuilder
}

// _ErrorReportingDataGeneric is the data-structure of this message
type _ErrorReportingDataGeneric struct {
	ErrorReportingDataContract
	SystemCategory ErrorReportingSystemCategory
	MostRecent     bool
	Acknowledge    bool
	MostSevere     bool
	Severity       ErrorReportingSeverity
	DeviceId       uint8
	ErrorData1     uint8
	ErrorData2     uint8
}

var _ ErrorReportingDataGeneric = (*_ErrorReportingDataGeneric)(nil)
var _ ErrorReportingDataRequirements = (*_ErrorReportingDataGeneric)(nil)

// NewErrorReportingDataGeneric factory function for _ErrorReportingDataGeneric
func NewErrorReportingDataGeneric(commandTypeContainer ErrorReportingCommandTypeContainer, systemCategory ErrorReportingSystemCategory, mostRecent bool, acknowledge bool, mostSevere bool, severity ErrorReportingSeverity, deviceId uint8, errorData1 uint8, errorData2 uint8) *_ErrorReportingDataGeneric {
	if systemCategory == nil {
		panic("systemCategory of type ErrorReportingSystemCategory for ErrorReportingDataGeneric must not be nil")
	}
	_result := &_ErrorReportingDataGeneric{
		ErrorReportingDataContract: NewErrorReportingData(commandTypeContainer),
		SystemCategory:             systemCategory,
		MostRecent:                 mostRecent,
		Acknowledge:                acknowledge,
		MostSevere:                 mostSevere,
		Severity:                   severity,
		DeviceId:                   deviceId,
		ErrorData1:                 errorData1,
		ErrorData2:                 errorData2,
	}
	_result.ErrorReportingDataContract.(*_ErrorReportingData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ErrorReportingDataGenericBuilder is a builder for ErrorReportingDataGeneric
type ErrorReportingDataGenericBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(systemCategory ErrorReportingSystemCategory, mostRecent bool, acknowledge bool, mostSevere bool, severity ErrorReportingSeverity, deviceId uint8, errorData1 uint8, errorData2 uint8) ErrorReportingDataGenericBuilder
	// WithSystemCategory adds SystemCategory (property field)
	WithSystemCategory(ErrorReportingSystemCategory) ErrorReportingDataGenericBuilder
	// WithSystemCategoryBuilder adds SystemCategory (property field) which is build by the builder
	WithSystemCategoryBuilder(func(ErrorReportingSystemCategoryBuilder) ErrorReportingSystemCategoryBuilder) ErrorReportingDataGenericBuilder
	// WithMostRecent adds MostRecent (property field)
	WithMostRecent(bool) ErrorReportingDataGenericBuilder
	// WithAcknowledge adds Acknowledge (property field)
	WithAcknowledge(bool) ErrorReportingDataGenericBuilder
	// WithMostSevere adds MostSevere (property field)
	WithMostSevere(bool) ErrorReportingDataGenericBuilder
	// WithSeverity adds Severity (property field)
	WithSeverity(ErrorReportingSeverity) ErrorReportingDataGenericBuilder
	// WithDeviceId adds DeviceId (property field)
	WithDeviceId(uint8) ErrorReportingDataGenericBuilder
	// WithErrorData1 adds ErrorData1 (property field)
	WithErrorData1(uint8) ErrorReportingDataGenericBuilder
	// WithErrorData2 adds ErrorData2 (property field)
	WithErrorData2(uint8) ErrorReportingDataGenericBuilder
	// Build builds the ErrorReportingDataGeneric or returns an error if something is wrong
	Build() (ErrorReportingDataGeneric, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ErrorReportingDataGeneric
}

// NewErrorReportingDataGenericBuilder() creates a ErrorReportingDataGenericBuilder
func NewErrorReportingDataGenericBuilder() ErrorReportingDataGenericBuilder {
	return &_ErrorReportingDataGenericBuilder{_ErrorReportingDataGeneric: new(_ErrorReportingDataGeneric)}
}

type _ErrorReportingDataGenericBuilder struct {
	*_ErrorReportingDataGeneric

	parentBuilder *_ErrorReportingDataBuilder

	err *utils.MultiError
}

var _ (ErrorReportingDataGenericBuilder) = (*_ErrorReportingDataGenericBuilder)(nil)

func (b *_ErrorReportingDataGenericBuilder) setParent(contract ErrorReportingDataContract) {
	b.ErrorReportingDataContract = contract
}

func (b *_ErrorReportingDataGenericBuilder) WithMandatoryFields(systemCategory ErrorReportingSystemCategory, mostRecent bool, acknowledge bool, mostSevere bool, severity ErrorReportingSeverity, deviceId uint8, errorData1 uint8, errorData2 uint8) ErrorReportingDataGenericBuilder {
	return b.WithSystemCategory(systemCategory).WithMostRecent(mostRecent).WithAcknowledge(acknowledge).WithMostSevere(mostSevere).WithSeverity(severity).WithDeviceId(deviceId).WithErrorData1(errorData1).WithErrorData2(errorData2)
}

func (b *_ErrorReportingDataGenericBuilder) WithSystemCategory(systemCategory ErrorReportingSystemCategory) ErrorReportingDataGenericBuilder {
	b.SystemCategory = systemCategory
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithSystemCategoryBuilder(builderSupplier func(ErrorReportingSystemCategoryBuilder) ErrorReportingSystemCategoryBuilder) ErrorReportingDataGenericBuilder {
	builder := builderSupplier(b.SystemCategory.CreateErrorReportingSystemCategoryBuilder())
	var err error
	b.SystemCategory, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "ErrorReportingSystemCategoryBuilder failed"))
	}
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithMostRecent(mostRecent bool) ErrorReportingDataGenericBuilder {
	b.MostRecent = mostRecent
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithAcknowledge(acknowledge bool) ErrorReportingDataGenericBuilder {
	b.Acknowledge = acknowledge
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithMostSevere(mostSevere bool) ErrorReportingDataGenericBuilder {
	b.MostSevere = mostSevere
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithSeverity(severity ErrorReportingSeverity) ErrorReportingDataGenericBuilder {
	b.Severity = severity
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithDeviceId(deviceId uint8) ErrorReportingDataGenericBuilder {
	b.DeviceId = deviceId
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithErrorData1(errorData1 uint8) ErrorReportingDataGenericBuilder {
	b.ErrorData1 = errorData1
	return b
}

func (b *_ErrorReportingDataGenericBuilder) WithErrorData2(errorData2 uint8) ErrorReportingDataGenericBuilder {
	b.ErrorData2 = errorData2
	return b
}

func (b *_ErrorReportingDataGenericBuilder) Build() (ErrorReportingDataGeneric, error) {
	if b.SystemCategory == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'systemCategory' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ErrorReportingDataGeneric.deepCopy(), nil
}

func (b *_ErrorReportingDataGenericBuilder) MustBuild() ErrorReportingDataGeneric {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ErrorReportingDataGenericBuilder) Done() ErrorReportingDataBuilder {
	return b.parentBuilder
}

func (b *_ErrorReportingDataGenericBuilder) buildForErrorReportingData() (ErrorReportingData, error) {
	return b.Build()
}

func (b *_ErrorReportingDataGenericBuilder) DeepCopy() any {
	_copy := b.CreateErrorReportingDataGenericBuilder().(*_ErrorReportingDataGenericBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateErrorReportingDataGenericBuilder creates a ErrorReportingDataGenericBuilder
func (b *_ErrorReportingDataGeneric) CreateErrorReportingDataGenericBuilder() ErrorReportingDataGenericBuilder {
	if b == nil {
		return NewErrorReportingDataGenericBuilder()
	}
	return &_ErrorReportingDataGenericBuilder{_ErrorReportingDataGeneric: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ErrorReportingDataGeneric) GetParent() ErrorReportingDataContract {
	return m.ErrorReportingDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ErrorReportingDataGeneric) GetSystemCategory() ErrorReportingSystemCategory {
	return m.SystemCategory
}

func (m *_ErrorReportingDataGeneric) GetMostRecent() bool {
	return m.MostRecent
}

func (m *_ErrorReportingDataGeneric) GetAcknowledge() bool {
	return m.Acknowledge
}

func (m *_ErrorReportingDataGeneric) GetMostSevere() bool {
	return m.MostSevere
}

func (m *_ErrorReportingDataGeneric) GetSeverity() ErrorReportingSeverity {
	return m.Severity
}

func (m *_ErrorReportingDataGeneric) GetDeviceId() uint8 {
	return m.DeviceId
}

func (m *_ErrorReportingDataGeneric) GetErrorData1() uint8 {
	return m.ErrorData1
}

func (m *_ErrorReportingDataGeneric) GetErrorData2() uint8 {
	return m.ErrorData2
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_ErrorReportingDataGeneric) GetIsMostSevereError() bool {
	ctx := context.Background()
	_ = ctx
	return bool(m.GetMostSevere())
}

func (m *_ErrorReportingDataGeneric) GetIsMostRecentError() bool {
	ctx := context.Background()
	_ = ctx
	return bool(m.GetMostRecent())
}

func (m *_ErrorReportingDataGeneric) GetIsMostRecentAndMostSevere() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool(m.GetIsMostRecentError()) && bool(m.GetIsMostSevereError()))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastErrorReportingDataGeneric(structType any) ErrorReportingDataGeneric {
	if casted, ok := structType.(ErrorReportingDataGeneric); ok {
		return casted
	}
	if casted, ok := structType.(*ErrorReportingDataGeneric); ok {
		return *casted
	}
	return nil
}

func (m *_ErrorReportingDataGeneric) GetTypeName() string {
	return "ErrorReportingDataGeneric"
}

func (m *_ErrorReportingDataGeneric) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ErrorReportingDataContract.(*_ErrorReportingData).GetLengthInBits(ctx))

	// Simple field (systemCategory)
	lengthInBits += m.SystemCategory.GetLengthInBits(ctx)

	// Simple field (mostRecent)
	lengthInBits += 1

	// Simple field (acknowledge)
	lengthInBits += 1

	// Simple field (mostSevere)
	lengthInBits += 1

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// Simple field (severity)
	lengthInBits += 3

	// Simple field (deviceId)
	lengthInBits += 8

	// Simple field (errorData1)
	lengthInBits += 8

	// Simple field (errorData2)
	lengthInBits += 8

	return lengthInBits
}

func (m *_ErrorReportingDataGeneric) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ErrorReportingDataGeneric) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ErrorReportingData) (__errorReportingDataGeneric ErrorReportingDataGeneric, err error) {
	m.ErrorReportingDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ErrorReportingDataGeneric"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ErrorReportingDataGeneric")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	systemCategory, err := ReadSimpleField[ErrorReportingSystemCategory](ctx, "systemCategory", ReadComplex[ErrorReportingSystemCategory](ErrorReportingSystemCategoryParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'systemCategory' field"))
	}
	m.SystemCategory = systemCategory

	mostRecent, err := ReadSimpleField(ctx, "mostRecent", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'mostRecent' field"))
	}
	m.MostRecent = mostRecent

	acknowledge, err := ReadSimpleField(ctx, "acknowledge", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'acknowledge' field"))
	}
	m.Acknowledge = acknowledge

	mostSevere, err := ReadSimpleField(ctx, "mostSevere", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'mostSevere' field"))
	}
	m.MostSevere = mostSevere

	// Validation
	if !(bool(mostRecent) || bool(mostSevere)) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "Invalid Error condition"})
	}

	isMostSevereError, err := ReadVirtualField[bool](ctx, "isMostSevereError", (*bool)(nil), mostSevere)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isMostSevereError' field"))
	}
	_ = isMostSevereError

	isMostRecentError, err := ReadVirtualField[bool](ctx, "isMostRecentError", (*bool)(nil), mostRecent)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isMostRecentError' field"))
	}
	_ = isMostRecentError

	isMostRecentAndMostSevere, err := ReadVirtualField[bool](ctx, "isMostRecentAndMostSevere", (*bool)(nil), bool(isMostRecentError) && bool(isMostSevereError))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isMostRecentAndMostSevere' field"))
	}
	_ = isMostRecentAndMostSevere

	severity, err := ReadEnumField[ErrorReportingSeverity](ctx, "severity", "ErrorReportingSeverity", ReadEnum(ErrorReportingSeverityByValue, ReadUnsignedByte(readBuffer, uint8(3))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'severity' field"))
	}
	m.Severity = severity

	deviceId, err := ReadSimpleField(ctx, "deviceId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'deviceId' field"))
	}
	m.DeviceId = deviceId

	errorData1, err := ReadSimpleField(ctx, "errorData1", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'errorData1' field"))
	}
	m.ErrorData1 = errorData1

	errorData2, err := ReadSimpleField(ctx, "errorData2", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'errorData2' field"))
	}
	m.ErrorData2 = errorData2

	if closeErr := readBuffer.CloseContext("ErrorReportingDataGeneric"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ErrorReportingDataGeneric")
	}

	return m, nil
}

func (m *_ErrorReportingDataGeneric) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ErrorReportingDataGeneric) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ErrorReportingDataGeneric"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ErrorReportingDataGeneric")
		}

		if err := WriteSimpleField[ErrorReportingSystemCategory](ctx, "systemCategory", m.GetSystemCategory(), WriteComplex[ErrorReportingSystemCategory](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'systemCategory' field")
		}

		if err := WriteSimpleField[bool](ctx, "mostRecent", m.GetMostRecent(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'mostRecent' field")
		}

		if err := WriteSimpleField[bool](ctx, "acknowledge", m.GetAcknowledge(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'acknowledge' field")
		}

		if err := WriteSimpleField[bool](ctx, "mostSevere", m.GetMostSevere(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'mostSevere' field")
		}
		// Virtual field
		isMostSevereError := m.GetIsMostSevereError()
		_ = isMostSevereError
		if _isMostSevereErrorErr := writeBuffer.WriteVirtual(ctx, "isMostSevereError", m.GetIsMostSevereError()); _isMostSevereErrorErr != nil {
			return errors.Wrap(_isMostSevereErrorErr, "Error serializing 'isMostSevereError' field")
		}
		// Virtual field
		isMostRecentError := m.GetIsMostRecentError()
		_ = isMostRecentError
		if _isMostRecentErrorErr := writeBuffer.WriteVirtual(ctx, "isMostRecentError", m.GetIsMostRecentError()); _isMostRecentErrorErr != nil {
			return errors.Wrap(_isMostRecentErrorErr, "Error serializing 'isMostRecentError' field")
		}
		// Virtual field
		isMostRecentAndMostSevere := m.GetIsMostRecentAndMostSevere()
		_ = isMostRecentAndMostSevere
		if _isMostRecentAndMostSevereErr := writeBuffer.WriteVirtual(ctx, "isMostRecentAndMostSevere", m.GetIsMostRecentAndMostSevere()); _isMostRecentAndMostSevereErr != nil {
			return errors.Wrap(_isMostRecentAndMostSevereErr, "Error serializing 'isMostRecentAndMostSevere' field")
		}

		if err := WriteSimpleEnumField[ErrorReportingSeverity](ctx, "severity", "ErrorReportingSeverity", m.GetSeverity(), WriteEnum[ErrorReportingSeverity, uint8](ErrorReportingSeverity.GetValue, ErrorReportingSeverity.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 3))); err != nil {
			return errors.Wrap(err, "Error serializing 'severity' field")
		}

		if err := WriteSimpleField[uint8](ctx, "deviceId", m.GetDeviceId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'deviceId' field")
		}

		if err := WriteSimpleField[uint8](ctx, "errorData1", m.GetErrorData1(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'errorData1' field")
		}

		if err := WriteSimpleField[uint8](ctx, "errorData2", m.GetErrorData2(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'errorData2' field")
		}

		if popErr := writeBuffer.PopContext("ErrorReportingDataGeneric"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ErrorReportingDataGeneric")
		}
		return nil
	}
	return m.ErrorReportingDataContract.(*_ErrorReportingData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ErrorReportingDataGeneric) IsErrorReportingDataGeneric() {}

func (m *_ErrorReportingDataGeneric) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ErrorReportingDataGeneric) deepCopy() *_ErrorReportingDataGeneric {
	if m == nil {
		return nil
	}
	_ErrorReportingDataGenericCopy := &_ErrorReportingDataGeneric{
		m.ErrorReportingDataContract.(*_ErrorReportingData).deepCopy(),
		m.SystemCategory.DeepCopy().(ErrorReportingSystemCategory),
		m.MostRecent,
		m.Acknowledge,
		m.MostSevere,
		m.Severity,
		m.DeviceId,
		m.ErrorData1,
		m.ErrorData2,
	}
	m.ErrorReportingDataContract.(*_ErrorReportingData)._SubType = m
	return _ErrorReportingDataGenericCopy
}

func (m *_ErrorReportingDataGeneric) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
