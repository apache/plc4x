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
    "encoding/base64"
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "strconv"
)

// Constant values.
const BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGELOWLIMITHEADER uint8 = 0x01
const BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGEHIGHLIMITHEADER uint8 = 0x03

// The data-structure of this message
type BACnetUnconfirmedServiceRequestWhoIs struct {
    DeviceInstanceRangeLowLimitLength uint8
    DeviceInstanceRangeLowLimit []int8
    DeviceInstanceRangeHighLimitLength uint8
    DeviceInstanceRangeHighLimit []int8
    Parent *BACnetUnconfirmedServiceRequest
    IBACnetUnconfirmedServiceRequestWhoIs
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestWhoIs interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *BACnetUnconfirmedServiceRequestWhoIs) ServiceChoice() uint8 {
    return 0x08
}


func (m *BACnetUnconfirmedServiceRequestWhoIs) InitializeParent(parent *BACnetUnconfirmedServiceRequest) {
}

func NewBACnetUnconfirmedServiceRequestWhoIs(deviceInstanceRangeLowLimitLength uint8, deviceInstanceRangeLowLimit []int8, deviceInstanceRangeHighLimitLength uint8, deviceInstanceRangeHighLimit []int8, ) *BACnetUnconfirmedServiceRequest {
    child := &BACnetUnconfirmedServiceRequestWhoIs{
        DeviceInstanceRangeLowLimitLength: deviceInstanceRangeLowLimitLength,
        DeviceInstanceRangeLowLimit: deviceInstanceRangeLowLimit,
        DeviceInstanceRangeHighLimitLength: deviceInstanceRangeHighLimitLength,
        DeviceInstanceRangeHighLimit: deviceInstanceRangeHighLimit,
        Parent: NewBACnetUnconfirmedServiceRequest(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastBACnetUnconfirmedServiceRequestWhoIs(structType interface{}) BACnetUnconfirmedServiceRequestWhoIs {
    castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestWhoIs {
        if casted, ok := typ.(BACnetUnconfirmedServiceRequestWhoIs); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetUnconfirmedServiceRequestWhoIs); ok {
            return *casted
        }
        if casted, ok := typ.(BACnetUnconfirmedServiceRequest); ok {
            return CastBACnetUnconfirmedServiceRequestWhoIs(casted.Child)
        }
        if casted, ok := typ.(*BACnetUnconfirmedServiceRequest); ok {
            return CastBACnetUnconfirmedServiceRequestWhoIs(casted.Child)
        }
        return BACnetUnconfirmedServiceRequestWhoIs{}
    }
    return castFunc(structType)
}

func (m *BACnetUnconfirmedServiceRequestWhoIs) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Const Field (deviceInstanceRangeLowLimitHeader)
    lengthInBits += 5

    // Simple field (deviceInstanceRangeLowLimitLength)
    lengthInBits += 3

    // Array field
    if len(m.DeviceInstanceRangeLowLimit) > 0 {
        lengthInBits += 8 * uint16(len(m.DeviceInstanceRangeLowLimit))
    }

    // Const Field (deviceInstanceRangeHighLimitHeader)
    lengthInBits += 5

    // Simple field (deviceInstanceRangeHighLimitLength)
    lengthInBits += 3

    // Array field
    if len(m.DeviceInstanceRangeHighLimit) > 0 {
        lengthInBits += 8 * uint16(len(m.DeviceInstanceRangeHighLimit))
    }

    return lengthInBits
}

func (m *BACnetUnconfirmedServiceRequestWhoIs) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestWhoIsParse(io *utils.ReadBuffer) (*BACnetUnconfirmedServiceRequest, error) {

    // Const Field (deviceInstanceRangeLowLimitHeader)
    deviceInstanceRangeLowLimitHeader, _deviceInstanceRangeLowLimitHeaderErr := io.ReadUint8(5)
    if _deviceInstanceRangeLowLimitHeaderErr != nil {
        return nil, errors.New("Error parsing 'deviceInstanceRangeLowLimitHeader' field " + _deviceInstanceRangeLowLimitHeaderErr.Error())
    }
    if deviceInstanceRangeLowLimitHeader != BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGELOWLIMITHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGELOWLIMITHEADER)) + " but got " + strconv.Itoa(int(deviceInstanceRangeLowLimitHeader)))
    }

    // Simple Field (deviceInstanceRangeLowLimitLength)
    deviceInstanceRangeLowLimitLength, _deviceInstanceRangeLowLimitLengthErr := io.ReadUint8(3)
    if _deviceInstanceRangeLowLimitLengthErr != nil {
        return nil, errors.New("Error parsing 'deviceInstanceRangeLowLimitLength' field " + _deviceInstanceRangeLowLimitLengthErr.Error())
    }

    // Array field (deviceInstanceRangeLowLimit)
    // Count array
    deviceInstanceRangeLowLimit := make([]int8, deviceInstanceRangeLowLimitLength)
    for curItem := uint16(0); curItem < uint16(deviceInstanceRangeLowLimitLength); curItem++ {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'deviceInstanceRangeLowLimit' field " + _err.Error())
        }
        deviceInstanceRangeLowLimit[curItem] = _item
    }

    // Const Field (deviceInstanceRangeHighLimitHeader)
    deviceInstanceRangeHighLimitHeader, _deviceInstanceRangeHighLimitHeaderErr := io.ReadUint8(5)
    if _deviceInstanceRangeHighLimitHeaderErr != nil {
        return nil, errors.New("Error parsing 'deviceInstanceRangeHighLimitHeader' field " + _deviceInstanceRangeHighLimitHeaderErr.Error())
    }
    if deviceInstanceRangeHighLimitHeader != BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGEHIGHLIMITHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGEHIGHLIMITHEADER)) + " but got " + strconv.Itoa(int(deviceInstanceRangeHighLimitHeader)))
    }

    // Simple Field (deviceInstanceRangeHighLimitLength)
    deviceInstanceRangeHighLimitLength, _deviceInstanceRangeHighLimitLengthErr := io.ReadUint8(3)
    if _deviceInstanceRangeHighLimitLengthErr != nil {
        return nil, errors.New("Error parsing 'deviceInstanceRangeHighLimitLength' field " + _deviceInstanceRangeHighLimitLengthErr.Error())
    }

    // Array field (deviceInstanceRangeHighLimit)
    // Count array
    deviceInstanceRangeHighLimit := make([]int8, deviceInstanceRangeHighLimitLength)
    for curItem := uint16(0); curItem < uint16(deviceInstanceRangeHighLimitLength); curItem++ {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'deviceInstanceRangeHighLimit' field " + _err.Error())
        }
        deviceInstanceRangeHighLimit[curItem] = _item
    }

    // Create a partially initialized instance
    _child := &BACnetUnconfirmedServiceRequestWhoIs{
        DeviceInstanceRangeLowLimitLength: deviceInstanceRangeLowLimitLength,
        DeviceInstanceRangeLowLimit: deviceInstanceRangeLowLimit,
        DeviceInstanceRangeHighLimitLength: deviceInstanceRangeHighLimitLength,
        DeviceInstanceRangeHighLimit: deviceInstanceRangeHighLimit,
        Parent: &BACnetUnconfirmedServiceRequest{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *BACnetUnconfirmedServiceRequestWhoIs) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Const Field (deviceInstanceRangeLowLimitHeader)
    _deviceInstanceRangeLowLimitHeaderErr := io.WriteUint8(5, 0x01)
    if _deviceInstanceRangeLowLimitHeaderErr != nil {
        return errors.New("Error serializing 'deviceInstanceRangeLowLimitHeader' field " + _deviceInstanceRangeLowLimitHeaderErr.Error())
    }

    // Simple Field (deviceInstanceRangeLowLimitLength)
    deviceInstanceRangeLowLimitLength := uint8(m.DeviceInstanceRangeLowLimitLength)
    _deviceInstanceRangeLowLimitLengthErr := io.WriteUint8(3, (deviceInstanceRangeLowLimitLength))
    if _deviceInstanceRangeLowLimitLengthErr != nil {
        return errors.New("Error serializing 'deviceInstanceRangeLowLimitLength' field " + _deviceInstanceRangeLowLimitLengthErr.Error())
    }

    // Array Field (deviceInstanceRangeLowLimit)
    if m.DeviceInstanceRangeLowLimit != nil {
        for _, _element := range m.DeviceInstanceRangeLowLimit {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'deviceInstanceRangeLowLimit' field " + _elementErr.Error())
            }
        }
    }

    // Const Field (deviceInstanceRangeHighLimitHeader)
    _deviceInstanceRangeHighLimitHeaderErr := io.WriteUint8(5, 0x03)
    if _deviceInstanceRangeHighLimitHeaderErr != nil {
        return errors.New("Error serializing 'deviceInstanceRangeHighLimitHeader' field " + _deviceInstanceRangeHighLimitHeaderErr.Error())
    }

    // Simple Field (deviceInstanceRangeHighLimitLength)
    deviceInstanceRangeHighLimitLength := uint8(m.DeviceInstanceRangeHighLimitLength)
    _deviceInstanceRangeHighLimitLengthErr := io.WriteUint8(3, (deviceInstanceRangeHighLimitLength))
    if _deviceInstanceRangeHighLimitLengthErr != nil {
        return errors.New("Error serializing 'deviceInstanceRangeHighLimitLength' field " + _deviceInstanceRangeHighLimitLengthErr.Error())
    }

    // Array Field (deviceInstanceRangeHighLimit)
    if m.DeviceInstanceRangeHighLimit != nil {
        for _, _element := range m.DeviceInstanceRangeHighLimit {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'deviceInstanceRangeHighLimit' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *BACnetUnconfirmedServiceRequestWhoIs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "deviceInstanceRangeLowLimitLength":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DeviceInstanceRangeLowLimitLength = data
            case "deviceInstanceRangeLowLimit":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.DeviceInstanceRangeLowLimit = utils.ByteToInt8(_decoded[0:_len])
            case "deviceInstanceRangeHighLimitLength":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DeviceInstanceRangeHighLimitLength = data
            case "deviceInstanceRangeHighLimit":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.DeviceInstanceRangeHighLimit = utils.ByteToInt8(_decoded[0:_len])
            }
        }
    }
}

func (m *BACnetUnconfirmedServiceRequestWhoIs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestWhoIs"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DeviceInstanceRangeLowLimitLength, xml.StartElement{Name: xml.Name{Local: "deviceInstanceRangeLowLimitLength"}}); err != nil {
        return err
    }
    _encodedDeviceInstanceRangeLowLimit := make([]byte, base64.StdEncoding.EncodedLen(len(m.DeviceInstanceRangeLowLimit)))
    base64.StdEncoding.Encode(_encodedDeviceInstanceRangeLowLimit, utils.Int8ToByte(m.DeviceInstanceRangeLowLimit))
    if err := e.EncodeElement(_encodedDeviceInstanceRangeLowLimit, xml.StartElement{Name: xml.Name{Local: "deviceInstanceRangeLowLimit"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DeviceInstanceRangeHighLimitLength, xml.StartElement{Name: xml.Name{Local: "deviceInstanceRangeHighLimitLength"}}); err != nil {
        return err
    }
    _encodedDeviceInstanceRangeHighLimit := make([]byte, base64.StdEncoding.EncodedLen(len(m.DeviceInstanceRangeHighLimit)))
    base64.StdEncoding.Encode(_encodedDeviceInstanceRangeHighLimit, utils.Int8ToByte(m.DeviceInstanceRangeHighLimit))
    if err := e.EncodeElement(_encodedDeviceInstanceRangeHighLimit, xml.StartElement{Name: xml.Name{Local: "deviceInstanceRangeHighLimit"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

