//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type CEMIFrame struct {
	repeated             bool
	priority             ICEMIPriority
	acknowledgeRequested bool
	errorFlag            bool
}

// The corresponding interface
type ICEMIFrame interface {
	spi.Message
	NotAckFrame() bool
	Polling() bool
	StandardFrame() bool
	Serialize(io spi.WriteBuffer)
}

type CEMIFrameInitializer interface {
	initialize(repeated bool, priority ICEMIPriority, acknowledgeRequested bool, errorFlag bool) spi.Message
}

func CEMIFrameNotAckFrame(m ICEMIFrame) bool {
	return m.NotAckFrame()
}

func CEMIFramePolling(m ICEMIFrame) bool {
	return m.Polling()
}

func CEMIFrameStandardFrame(m ICEMIFrame) bool {
	return m.StandardFrame()
}

func CastICEMIFrame(structType interface{}) ICEMIFrame {
	castFunc := func(typ interface{}) ICEMIFrame {
		if iCEMIFrame, ok := typ.(ICEMIFrame); ok {
			return iCEMIFrame
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIFrame(structType interface{}) CEMIFrame {
	castFunc := func(typ interface{}) CEMIFrame {
		if sCEMIFrame, ok := typ.(CEMIFrame); ok {
			return sCEMIFrame
		}
		return CEMIFrame{}
	}
	return castFunc(structType)
}

func (m CEMIFrame) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (standardFrame)
	lengthInBits += 1

	// Discriminator Field (polling)
	lengthInBits += 1

	// Simple field (repeated)
	lengthInBits += 1

	// Discriminator Field (notAckFrame)
	lengthInBits += 1

	// Enum Field (priority)
	lengthInBits += 2

	// Simple field (acknowledgeRequested)
	lengthInBits += 1

	// Simple field (errorFlag)
	lengthInBits += 1

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m CEMIFrame) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIFrameParse(io spi.ReadBuffer) (spi.Message, error) {

	// Discriminator Field (standardFrame) (Used as input to a switch field)
	standardFrame, _standardFrameErr := io.ReadBit()
	if _standardFrameErr != nil {
		return nil, errors.New("Error parsing 'standardFrame' field " + _standardFrameErr.Error())
	}

	// Discriminator Field (polling) (Used as input to a switch field)
	polling, _pollingErr := io.ReadBit()
	if _pollingErr != nil {
		return nil, errors.New("Error parsing 'polling' field " + _pollingErr.Error())
	}

	// Simple Field (repeated)
	repeated, _repeatedErr := io.ReadBit()
	if _repeatedErr != nil {
		return nil, errors.New("Error parsing 'repeated' field " + _repeatedErr.Error())
	}

	// Discriminator Field (notAckFrame) (Used as input to a switch field)
	notAckFrame, _notAckFrameErr := io.ReadBit()
	if _notAckFrameErr != nil {
		return nil, errors.New("Error parsing 'notAckFrame' field " + _notAckFrameErr.Error())
	}

	// Enum field (priority)
	priority, _priorityErr := CEMIPriorityParse(io)
	if _priorityErr != nil {
		return nil, errors.New("Error parsing 'priority' field " + _priorityErr.Error())
	}

	// Simple Field (acknowledgeRequested)
	acknowledgeRequested, _acknowledgeRequestedErr := io.ReadBit()
	if _acknowledgeRequestedErr != nil {
		return nil, errors.New("Error parsing 'acknowledgeRequested' field " + _acknowledgeRequestedErr.Error())
	}

	// Simple Field (errorFlag)
	errorFlag, _errorFlagErr := io.ReadBit()
	if _errorFlagErr != nil {
		return nil, errors.New("Error parsing 'errorFlag' field " + _errorFlagErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer CEMIFrameInitializer
	var typeSwitchError error
	switch {
	case notAckFrame == false:
		initializer, typeSwitchError = CEMIFrameAckParse(io)
	case notAckFrame == true && standardFrame == true && polling == false:
		initializer, typeSwitchError = CEMIFrameDataParse(io)
	case notAckFrame == true && standardFrame == true && polling == true:
		initializer, typeSwitchError = CEMIFramePollingDataParse(io)
	case notAckFrame == true && standardFrame == false && polling == false:
		initializer, typeSwitchError = CEMIFrameDataExtParse(io)
	case notAckFrame == true && standardFrame == false && polling == true:
		initializer, typeSwitchError = CEMIFramePollingDataExtParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(repeated, priority, acknowledgeRequested, errorFlag), nil
}

func CEMIFrameSerialize(io spi.WriteBuffer, m CEMIFrame, i ICEMIFrame, childSerialize func()) {

	// Discriminator Field (standardFrame) (Used as input to a switch field)
	standardFrame := bool(i.StandardFrame())
	io.WriteBit((bool)(standardFrame))

	// Discriminator Field (polling) (Used as input to a switch field)
	polling := bool(i.Polling())
	io.WriteBit((bool)(polling))

	// Simple Field (repeated)
	repeated := bool(m.repeated)
	io.WriteBit((bool)(repeated))

	// Discriminator Field (notAckFrame) (Used as input to a switch field)
	notAckFrame := bool(i.NotAckFrame())
	io.WriteBit((bool)(notAckFrame))

	// Enum field (priority)
	priority := CastCEMIPriority(m.priority)
	priority.Serialize(io)

	// Simple Field (acknowledgeRequested)
	acknowledgeRequested := bool(m.acknowledgeRequested)
	io.WriteBit((bool)(acknowledgeRequested))

	// Simple Field (errorFlag)
	errorFlag := bool(m.errorFlag)
	io.WriteBit((bool)(errorFlag))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	childSerialize()

}
