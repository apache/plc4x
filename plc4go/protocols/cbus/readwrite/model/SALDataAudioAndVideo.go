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

// SALDataAudioAndVideo is the corresponding interface of SALDataAudioAndVideo
type SALDataAudioAndVideo interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	SALData
	// GetAudioVideoData returns AudioVideoData (property field)
	GetAudioVideoData() LightingData
	// IsSALDataAudioAndVideo is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSALDataAudioAndVideo()
}

// _SALDataAudioAndVideo is the data-structure of this message
type _SALDataAudioAndVideo struct {
	SALDataContract
	AudioVideoData LightingData
}

var _ SALDataAudioAndVideo = (*_SALDataAudioAndVideo)(nil)
var _ SALDataRequirements = (*_SALDataAudioAndVideo)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SALDataAudioAndVideo) GetApplicationId() ApplicationId {
	return ApplicationId_AUDIO_AND_VIDEO
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SALDataAudioAndVideo) GetParent() SALDataContract {
	return m.SALDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SALDataAudioAndVideo) GetAudioVideoData() LightingData {
	return m.AudioVideoData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewSALDataAudioAndVideo factory function for _SALDataAudioAndVideo
func NewSALDataAudioAndVideo(audioVideoData LightingData, salData SALData) *_SALDataAudioAndVideo {
	if audioVideoData == nil {
		panic("audioVideoData of type LightingData for SALDataAudioAndVideo must not be nil")
	}
	_result := &_SALDataAudioAndVideo{
		SALDataContract: NewSALData(salData),
		AudioVideoData:  audioVideoData,
	}
	_result.SALDataContract.(*_SALData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastSALDataAudioAndVideo(structType any) SALDataAudioAndVideo {
	if casted, ok := structType.(SALDataAudioAndVideo); ok {
		return casted
	}
	if casted, ok := structType.(*SALDataAudioAndVideo); ok {
		return *casted
	}
	return nil
}

func (m *_SALDataAudioAndVideo) GetTypeName() string {
	return "SALDataAudioAndVideo"
}

func (m *_SALDataAudioAndVideo) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.SALDataContract.(*_SALData).getLengthInBits(ctx))

	// Simple field (audioVideoData)
	lengthInBits += m.AudioVideoData.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_SALDataAudioAndVideo) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_SALDataAudioAndVideo) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_SALData, applicationId ApplicationId) (__sALDataAudioAndVideo SALDataAudioAndVideo, err error) {
	m.SALDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SALDataAudioAndVideo"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SALDataAudioAndVideo")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	audioVideoData, err := ReadSimpleField[LightingData](ctx, "audioVideoData", ReadComplex[LightingData](LightingDataParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'audioVideoData' field"))
	}
	m.AudioVideoData = audioVideoData

	if closeErr := readBuffer.CloseContext("SALDataAudioAndVideo"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SALDataAudioAndVideo")
	}

	return m, nil
}

func (m *_SALDataAudioAndVideo) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SALDataAudioAndVideo) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SALDataAudioAndVideo"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SALDataAudioAndVideo")
		}

		if err := WriteSimpleField[LightingData](ctx, "audioVideoData", m.GetAudioVideoData(), WriteComplex[LightingData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'audioVideoData' field")
		}

		if popErr := writeBuffer.PopContext("SALDataAudioAndVideo"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SALDataAudioAndVideo")
		}
		return nil
	}
	return m.SALDataContract.(*_SALData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SALDataAudioAndVideo) IsSALDataAudioAndVideo() {}

func (m *_SALDataAudioAndVideo) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
