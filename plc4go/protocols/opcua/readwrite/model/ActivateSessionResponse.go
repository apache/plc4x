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

// ActivateSessionResponse is the corresponding interface of ActivateSessionResponse
type ActivateSessionResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ExtensionObjectDefinition
	// GetServerNonce returns ServerNonce (property field)
	GetServerNonce() PascalByteString
	// GetNoOfResults returns NoOfResults (property field)
	GetNoOfResults() int32
	// GetResults returns Results (property field)
	GetResults() []StatusCode
	// GetNoOfDiagnosticInfos returns NoOfDiagnosticInfos (property field)
	GetNoOfDiagnosticInfos() int32
	// GetDiagnosticInfos returns DiagnosticInfos (property field)
	GetDiagnosticInfos() []DiagnosticInfo
}

// ActivateSessionResponseExactly can be used when we want exactly this type and not a type which fulfills ActivateSessionResponse.
// This is useful for switch cases.
type ActivateSessionResponseExactly interface {
	ActivateSessionResponse
	isActivateSessionResponse() bool
}

// _ActivateSessionResponse is the data-structure of this message
type _ActivateSessionResponse struct {
	*_ExtensionObjectDefinition
	ResponseHeader      ExtensionObjectDefinition
	ServerNonce         PascalByteString
	NoOfResults         int32
	Results             []StatusCode
	NoOfDiagnosticInfos int32
	DiagnosticInfos     []DiagnosticInfo
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ActivateSessionResponse) GetIdentifier() string {
	return "470"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ActivateSessionResponse) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_ActivateSessionResponse) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ActivateSessionResponse) GetResponseHeader() ExtensionObjectDefinition {
	return m.ResponseHeader
}

func (m *_ActivateSessionResponse) GetServerNonce() PascalByteString {
	return m.ServerNonce
}

func (m *_ActivateSessionResponse) GetNoOfResults() int32 {
	return m.NoOfResults
}

func (m *_ActivateSessionResponse) GetResults() []StatusCode {
	return m.Results
}

func (m *_ActivateSessionResponse) GetNoOfDiagnosticInfos() int32 {
	return m.NoOfDiagnosticInfos
}

func (m *_ActivateSessionResponse) GetDiagnosticInfos() []DiagnosticInfo {
	return m.DiagnosticInfos
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewActivateSessionResponse factory function for _ActivateSessionResponse
func NewActivateSessionResponse(responseHeader ExtensionObjectDefinition, serverNonce PascalByteString, noOfResults int32, results []StatusCode, noOfDiagnosticInfos int32, diagnosticInfos []DiagnosticInfo) *_ActivateSessionResponse {
	_result := &_ActivateSessionResponse{
		ResponseHeader:             responseHeader,
		ServerNonce:                serverNonce,
		NoOfResults:                noOfResults,
		Results:                    results,
		NoOfDiagnosticInfos:        noOfDiagnosticInfos,
		DiagnosticInfos:            diagnosticInfos,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastActivateSessionResponse(structType any) ActivateSessionResponse {
	if casted, ok := structType.(ActivateSessionResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ActivateSessionResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ActivateSessionResponse) GetTypeName() string {
	return "ActivateSessionResponse"
}

func (m *_ActivateSessionResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	// Simple field (serverNonce)
	lengthInBits += m.ServerNonce.GetLengthInBits(ctx)

	// Simple field (noOfResults)
	lengthInBits += 32

	// Array field
	if len(m.Results) > 0 {
		for _curItem, element := range m.Results {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.Results), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (noOfDiagnosticInfos)
	lengthInBits += 32

	// Array field
	if len(m.DiagnosticInfos) > 0 {
		for _curItem, element := range m.DiagnosticInfos {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.DiagnosticInfos), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_ActivateSessionResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func ActivateSessionResponseParse(ctx context.Context, theBytes []byte, identifier string) (ActivateSessionResponse, error) {
	return ActivateSessionResponseParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func ActivateSessionResponseParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (ActivateSessionResponse, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (ActivateSessionResponse, error) {
		return ActivateSessionResponseParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func ActivateSessionResponseParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (ActivateSessionResponse, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ActivateSessionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ActivateSessionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "responseHeader", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("394")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}

	serverNonce, err := ReadSimpleField[PascalByteString](ctx, "serverNonce", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serverNonce' field"))
	}

	noOfResults, err := ReadSimpleField(ctx, "noOfResults", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfResults' field"))
	}

	results, err := ReadCountArrayField[StatusCode](ctx, "results", ReadComplex[StatusCode](StatusCodeParseWithBuffer, readBuffer), uint64(noOfResults))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'results' field"))
	}

	noOfDiagnosticInfos, err := ReadSimpleField(ctx, "noOfDiagnosticInfos", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfDiagnosticInfos' field"))
	}

	diagnosticInfos, err := ReadCountArrayField[DiagnosticInfo](ctx, "diagnosticInfos", ReadComplex[DiagnosticInfo](DiagnosticInfoParseWithBuffer, readBuffer), uint64(noOfDiagnosticInfos))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'diagnosticInfos' field"))
	}

	if closeErr := readBuffer.CloseContext("ActivateSessionResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ActivateSessionResponse")
	}

	// Create a partially initialized instance
	_child := &_ActivateSessionResponse{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		ResponseHeader:             responseHeader,
		ServerNonce:                serverNonce,
		NoOfResults:                noOfResults,
		Results:                    results,
		NoOfDiagnosticInfos:        noOfDiagnosticInfos,
		DiagnosticInfos:            diagnosticInfos,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_ActivateSessionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ActivateSessionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ActivateSessionResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ActivateSessionResponse")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "serverNonce", m.GetServerNonce(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'serverNonce' field")
		}

		if err := WriteSimpleField[int32](ctx, "noOfResults", m.GetNoOfResults(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfResults' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "results", m.GetResults(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'results' field")
		}

		if err := WriteSimpleField[int32](ctx, "noOfDiagnosticInfos", m.GetNoOfDiagnosticInfos(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfDiagnosticInfos' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "diagnosticInfos", m.GetDiagnosticInfos(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'diagnosticInfos' field")
		}

		if popErr := writeBuffer.PopContext("ActivateSessionResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ActivateSessionResponse")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ActivateSessionResponse) isActivateSessionResponse() bool {
	return true
}

func (m *_ActivateSessionResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
