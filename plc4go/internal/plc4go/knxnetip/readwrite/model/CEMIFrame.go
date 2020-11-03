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
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type CEMIFrame struct {
    Repeated bool
    Priority ICEMIPriority
    AcknowledgeRequested bool
    ErrorFlag bool

}

// The corresponding interface
type ICEMIFrame interface {
    spi.Message
    NotAckFrame() bool
    Polling() bool
    StandardFrame() bool
    Serialize(io utils.WriteBuffer) error
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
        if sCEMIFrame, ok := typ.(*CEMIFrame); ok {
            return *sCEMIFrame
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

func CEMIFrameParse(io *utils.ReadBuffer) (spi.Message, error) {

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

func CEMIFrameSerialize(io utils.WriteBuffer, m CEMIFrame, i ICEMIFrame, childSerialize func() error) error {

    // Discriminator Field (standardFrame) (Used as input to a switch field)
    standardFrame := bool(i.StandardFrame())
    _standardFrameErr := io.WriteBit((standardFrame))
    if _standardFrameErr != nil {
        return errors.New("Error serializing 'standardFrame' field " + _standardFrameErr.Error())
    }

    // Discriminator Field (polling) (Used as input to a switch field)
    polling := bool(i.Polling())
    _pollingErr := io.WriteBit((polling))
    if _pollingErr != nil {
        return errors.New("Error serializing 'polling' field " + _pollingErr.Error())
    }

    // Simple Field (repeated)
    repeated := bool(m.Repeated)
    _repeatedErr := io.WriteBit((repeated))
    if _repeatedErr != nil {
        return errors.New("Error serializing 'repeated' field " + _repeatedErr.Error())
    }

    // Discriminator Field (notAckFrame) (Used as input to a switch field)
    notAckFrame := bool(i.NotAckFrame())
    _notAckFrameErr := io.WriteBit((notAckFrame))
    if _notAckFrameErr != nil {
        return errors.New("Error serializing 'notAckFrame' field " + _notAckFrameErr.Error())
    }

    // Enum field (priority)
    priority := CastCEMIPriority(m.Priority)
    _priorityErr := priority.Serialize(io)
    if _priorityErr != nil {
        return errors.New("Error serializing 'priority' field " + _priorityErr.Error())
    }

    // Simple Field (acknowledgeRequested)
    acknowledgeRequested := bool(m.AcknowledgeRequested)
    _acknowledgeRequestedErr := io.WriteBit((acknowledgeRequested))
    if _acknowledgeRequestedErr != nil {
        return errors.New("Error serializing 'acknowledgeRequested' field " + _acknowledgeRequestedErr.Error())
    }

    // Simple Field (errorFlag)
    errorFlag := bool(m.ErrorFlag)
    _errorFlagErr := io.WriteBit((errorFlag))
    if _errorFlagErr != nil {
        return errors.New("Error serializing 'errorFlag' field " + _errorFlagErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *CEMIFrame) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "repeated":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Repeated = data
            case "priority":
                var data *CEMIPriority
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Priority = data
            case "acknowledgeRequested":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.AcknowledgeRequested = data
            case "errorFlag":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ErrorFlag = data
            }
        }
    }
}

func (m CEMIFrame) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.CEMIFrame"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Repeated, xml.StartElement{Name: xml.Name{Local: "repeated"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Priority, xml.StartElement{Name: xml.Name{Local: "priority"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.AcknowledgeRequested, xml.StartElement{Name: xml.Name{Local: "acknowledgeRequested"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ErrorFlag, xml.StartElement{Name: xml.Name{Local: "errorFlag"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

