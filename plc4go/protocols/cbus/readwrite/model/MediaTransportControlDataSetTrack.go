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

// MediaTransportControlDataSetTrack is the corresponding interface of MediaTransportControlDataSetTrack
type MediaTransportControlDataSetTrack interface {
	utils.LengthAware
	utils.Serializable
	MediaTransportControlData
	// GetTrackMSB returns TrackMSB (property field)
	GetTrackMSB() byte
	// GetTrackMMSB returns TrackMMSB (property field)
	GetTrackMMSB() byte
	// GetTrackMLSB returns TrackMLSB (property field)
	GetTrackMLSB() byte
	// GetTrackLSB returns TrackLSB (property field)
	GetTrackLSB() byte
}

// MediaTransportControlDataSetTrackExactly can be used when we want exactly this type and not a type which fulfills MediaTransportControlDataSetTrack.
// This is useful for switch cases.
type MediaTransportControlDataSetTrackExactly interface {
	MediaTransportControlDataSetTrack
	isMediaTransportControlDataSetTrack() bool
}

// _MediaTransportControlDataSetTrack is the data-structure of this message
type _MediaTransportControlDataSetTrack struct {
	*_MediaTransportControlData
	TrackMSB  byte
	TrackMMSB byte
	TrackMLSB byte
	TrackLSB  byte
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_MediaTransportControlDataSetTrack) InitializeParent(parent MediaTransportControlData, commandTypeContainer MediaTransportControlCommandTypeContainer, mediaLinkGroup byte) {
	m.CommandTypeContainer = commandTypeContainer
	m.MediaLinkGroup = mediaLinkGroup
}

func (m *_MediaTransportControlDataSetTrack) GetParent() MediaTransportControlData {
	return m._MediaTransportControlData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_MediaTransportControlDataSetTrack) GetTrackMSB() byte {
	return m.TrackMSB
}

func (m *_MediaTransportControlDataSetTrack) GetTrackMMSB() byte {
	return m.TrackMMSB
}

func (m *_MediaTransportControlDataSetTrack) GetTrackMLSB() byte {
	return m.TrackMLSB
}

func (m *_MediaTransportControlDataSetTrack) GetTrackLSB() byte {
	return m.TrackLSB
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewMediaTransportControlDataSetTrack factory function for _MediaTransportControlDataSetTrack
func NewMediaTransportControlDataSetTrack(trackMSB byte, trackMMSB byte, trackMLSB byte, trackLSB byte, commandTypeContainer MediaTransportControlCommandTypeContainer, mediaLinkGroup byte) *_MediaTransportControlDataSetTrack {
	_result := &_MediaTransportControlDataSetTrack{
		TrackMSB:                   trackMSB,
		TrackMMSB:                  trackMMSB,
		TrackMLSB:                  trackMLSB,
		TrackLSB:                   trackLSB,
		_MediaTransportControlData: NewMediaTransportControlData(commandTypeContainer, mediaLinkGroup),
	}
	_result._MediaTransportControlData._MediaTransportControlDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastMediaTransportControlDataSetTrack(structType interface{}) MediaTransportControlDataSetTrack {
	if casted, ok := structType.(MediaTransportControlDataSetTrack); ok {
		return casted
	}
	if casted, ok := structType.(*MediaTransportControlDataSetTrack); ok {
		return *casted
	}
	return nil
}

func (m *_MediaTransportControlDataSetTrack) GetTypeName() string {
	return "MediaTransportControlDataSetTrack"
}

func (m *_MediaTransportControlDataSetTrack) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_MediaTransportControlDataSetTrack) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (trackMSB)
	lengthInBits += 8

	// Simple field (trackMMSB)
	lengthInBits += 8

	// Simple field (trackMLSB)
	lengthInBits += 8

	// Simple field (trackLSB)
	lengthInBits += 8

	return lengthInBits
}

func (m *_MediaTransportControlDataSetTrack) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func MediaTransportControlDataSetTrackParse(readBuffer utils.ReadBuffer) (MediaTransportControlDataSetTrack, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("MediaTransportControlDataSetTrack"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for MediaTransportControlDataSetTrack")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (trackMSB)
	_trackMSB, _trackMSBErr := readBuffer.ReadByte("trackMSB")
	if _trackMSBErr != nil {
		return nil, errors.Wrap(_trackMSBErr, "Error parsing 'trackMSB' field of MediaTransportControlDataSetTrack")
	}
	trackMSB := _trackMSB

	// Simple Field (trackMMSB)
	_trackMMSB, _trackMMSBErr := readBuffer.ReadByte("trackMMSB")
	if _trackMMSBErr != nil {
		return nil, errors.Wrap(_trackMMSBErr, "Error parsing 'trackMMSB' field of MediaTransportControlDataSetTrack")
	}
	trackMMSB := _trackMMSB

	// Simple Field (trackMLSB)
	_trackMLSB, _trackMLSBErr := readBuffer.ReadByte("trackMLSB")
	if _trackMLSBErr != nil {
		return nil, errors.Wrap(_trackMLSBErr, "Error parsing 'trackMLSB' field of MediaTransportControlDataSetTrack")
	}
	trackMLSB := _trackMLSB

	// Simple Field (trackLSB)
	_trackLSB, _trackLSBErr := readBuffer.ReadByte("trackLSB")
	if _trackLSBErr != nil {
		return nil, errors.Wrap(_trackLSBErr, "Error parsing 'trackLSB' field of MediaTransportControlDataSetTrack")
	}
	trackLSB := _trackLSB

	if closeErr := readBuffer.CloseContext("MediaTransportControlDataSetTrack"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for MediaTransportControlDataSetTrack")
	}

	// Create a partially initialized instance
	_child := &_MediaTransportControlDataSetTrack{
		_MediaTransportControlData: &_MediaTransportControlData{},
		TrackMSB:                   trackMSB,
		TrackMMSB:                  trackMMSB,
		TrackMLSB:                  trackMLSB,
		TrackLSB:                   trackLSB,
	}
	_child._MediaTransportControlData._MediaTransportControlDataChildRequirements = _child
	return _child, nil
}

func (m *_MediaTransportControlDataSetTrack) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MediaTransportControlDataSetTrack"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for MediaTransportControlDataSetTrack")
		}

		// Simple Field (trackMSB)
		trackMSB := byte(m.GetTrackMSB())
		_trackMSBErr := writeBuffer.WriteByte("trackMSB", (trackMSB))
		if _trackMSBErr != nil {
			return errors.Wrap(_trackMSBErr, "Error serializing 'trackMSB' field")
		}

		// Simple Field (trackMMSB)
		trackMMSB := byte(m.GetTrackMMSB())
		_trackMMSBErr := writeBuffer.WriteByte("trackMMSB", (trackMMSB))
		if _trackMMSBErr != nil {
			return errors.Wrap(_trackMMSBErr, "Error serializing 'trackMMSB' field")
		}

		// Simple Field (trackMLSB)
		trackMLSB := byte(m.GetTrackMLSB())
		_trackMLSBErr := writeBuffer.WriteByte("trackMLSB", (trackMLSB))
		if _trackMLSBErr != nil {
			return errors.Wrap(_trackMLSBErr, "Error serializing 'trackMLSB' field")
		}

		// Simple Field (trackLSB)
		trackLSB := byte(m.GetTrackLSB())
		_trackLSBErr := writeBuffer.WriteByte("trackLSB", (trackLSB))
		if _trackLSBErr != nil {
			return errors.Wrap(_trackLSBErr, "Error serializing 'trackLSB' field")
		}

		if popErr := writeBuffer.PopContext("MediaTransportControlDataSetTrack"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for MediaTransportControlDataSetTrack")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_MediaTransportControlDataSetTrack) isMediaTransportControlDataSetTrack() bool {
	return true
}

func (m *_MediaTransportControlDataSetTrack) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
