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
const BACnetUnconfirmedServiceRequestIAm_OBJECTIDENTIFIERHEADER uint8 = 0xC4
const BACnetUnconfirmedServiceRequestIAm_MAXIMUMAPDULENGTHACCEPTEDHEADER uint8 = 0x04
const BACnetUnconfirmedServiceRequestIAm_SEGMENTATIONSUPPORTEDHEADER uint8 = 0x91
const BACnetUnconfirmedServiceRequestIAm_VENDORIDHEADER uint8 = 0x21

// The data-structure of this message
type BACnetUnconfirmedServiceRequestIAm struct {
    ObjectType uint16
    ObjectInstanceNumber uint32
    MaximumApduLengthAcceptedLength uint8
    MaximumApduLengthAccepted []int8
    SegmentationSupported uint8
    VendorId uint8
    Parent *BACnetUnconfirmedServiceRequest
    IBACnetUnconfirmedServiceRequestIAm
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestIAm interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *BACnetUnconfirmedServiceRequestIAm) ServiceChoice() uint8 {
    return 0x00
}


func (m *BACnetUnconfirmedServiceRequestIAm) InitializeParent(parent *BACnetUnconfirmedServiceRequest) {
}

func NewBACnetUnconfirmedServiceRequestIAm(objectType uint16, objectInstanceNumber uint32, maximumApduLengthAcceptedLength uint8, maximumApduLengthAccepted []int8, segmentationSupported uint8, vendorId uint8, ) *BACnetUnconfirmedServiceRequest {
    child := &BACnetUnconfirmedServiceRequestIAm{
        ObjectType: objectType,
        ObjectInstanceNumber: objectInstanceNumber,
        MaximumApduLengthAcceptedLength: maximumApduLengthAcceptedLength,
        MaximumApduLengthAccepted: maximumApduLengthAccepted,
        SegmentationSupported: segmentationSupported,
        VendorId: vendorId,
        Parent: NewBACnetUnconfirmedServiceRequest(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastBACnetUnconfirmedServiceRequestIAm(structType interface{}) BACnetUnconfirmedServiceRequestIAm {
    castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestIAm {
        if casted, ok := typ.(BACnetUnconfirmedServiceRequestIAm); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetUnconfirmedServiceRequestIAm); ok {
            return *casted
        }
        if casted, ok := typ.(BACnetUnconfirmedServiceRequest); ok {
            return CastBACnetUnconfirmedServiceRequestIAm(casted.Child)
        }
        if casted, ok := typ.(*BACnetUnconfirmedServiceRequest); ok {
            return CastBACnetUnconfirmedServiceRequestIAm(casted.Child)
        }
        return BACnetUnconfirmedServiceRequestIAm{}
    }
    return castFunc(structType)
}

func (m *BACnetUnconfirmedServiceRequestIAm) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Const Field (objectIdentifierHeader)
    lengthInBits += 8

    // Simple field (objectType)
    lengthInBits += 10

    // Simple field (objectInstanceNumber)
    lengthInBits += 22

    // Const Field (maximumApduLengthAcceptedHeader)
    lengthInBits += 5

    // Simple field (maximumApduLengthAcceptedLength)
    lengthInBits += 3

    // Array field
    if len(m.MaximumApduLengthAccepted) > 0 {
        lengthInBits += 8 * uint16(len(m.MaximumApduLengthAccepted))
    }

    // Const Field (segmentationSupportedHeader)
    lengthInBits += 8

    // Simple field (segmentationSupported)
    lengthInBits += 8

    // Const Field (vendorIdHeader)
    lengthInBits += 8

    // Simple field (vendorId)
    lengthInBits += 8

    return lengthInBits
}

func (m *BACnetUnconfirmedServiceRequestIAm) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestIAmParse(io *utils.ReadBuffer) (*BACnetUnconfirmedServiceRequest, error) {

    // Const Field (objectIdentifierHeader)
    objectIdentifierHeader, _objectIdentifierHeaderErr := io.ReadUint8(8)
    if _objectIdentifierHeaderErr != nil {
        return nil, errors.New("Error parsing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
    }
    if objectIdentifierHeader != BACnetUnconfirmedServiceRequestIAm_OBJECTIDENTIFIERHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
    }

    // Simple Field (objectType)
    objectType, _objectTypeErr := io.ReadUint16(10)
    if _objectTypeErr != nil {
        return nil, errors.New("Error parsing 'objectType' field " + _objectTypeErr.Error())
    }

    // Simple Field (objectInstanceNumber)
    objectInstanceNumber, _objectInstanceNumberErr := io.ReadUint32(22)
    if _objectInstanceNumberErr != nil {
        return nil, errors.New("Error parsing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
    }

    // Const Field (maximumApduLengthAcceptedHeader)
    maximumApduLengthAcceptedHeader, _maximumApduLengthAcceptedHeaderErr := io.ReadUint8(5)
    if _maximumApduLengthAcceptedHeaderErr != nil {
        return nil, errors.New("Error parsing 'maximumApduLengthAcceptedHeader' field " + _maximumApduLengthAcceptedHeaderErr.Error())
    }
    if maximumApduLengthAcceptedHeader != BACnetUnconfirmedServiceRequestIAm_MAXIMUMAPDULENGTHACCEPTEDHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_MAXIMUMAPDULENGTHACCEPTEDHEADER)) + " but got " + strconv.Itoa(int(maximumApduLengthAcceptedHeader)))
    }

    // Simple Field (maximumApduLengthAcceptedLength)
    maximumApduLengthAcceptedLength, _maximumApduLengthAcceptedLengthErr := io.ReadUint8(3)
    if _maximumApduLengthAcceptedLengthErr != nil {
        return nil, errors.New("Error parsing 'maximumApduLengthAcceptedLength' field " + _maximumApduLengthAcceptedLengthErr.Error())
    }

    // Array field (maximumApduLengthAccepted)
    // Count array
    maximumApduLengthAccepted := make([]int8, maximumApduLengthAcceptedLength)
    for curItem := uint16(0); curItem < uint16(maximumApduLengthAcceptedLength); curItem++ {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'maximumApduLengthAccepted' field " + _err.Error())
        }
        maximumApduLengthAccepted[curItem] = _item
    }

    // Const Field (segmentationSupportedHeader)
    segmentationSupportedHeader, _segmentationSupportedHeaderErr := io.ReadUint8(8)
    if _segmentationSupportedHeaderErr != nil {
        return nil, errors.New("Error parsing 'segmentationSupportedHeader' field " + _segmentationSupportedHeaderErr.Error())
    }
    if segmentationSupportedHeader != BACnetUnconfirmedServiceRequestIAm_SEGMENTATIONSUPPORTEDHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_SEGMENTATIONSUPPORTEDHEADER)) + " but got " + strconv.Itoa(int(segmentationSupportedHeader)))
    }

    // Simple Field (segmentationSupported)
    segmentationSupported, _segmentationSupportedErr := io.ReadUint8(8)
    if _segmentationSupportedErr != nil {
        return nil, errors.New("Error parsing 'segmentationSupported' field " + _segmentationSupportedErr.Error())
    }

    // Const Field (vendorIdHeader)
    vendorIdHeader, _vendorIdHeaderErr := io.ReadUint8(8)
    if _vendorIdHeaderErr != nil {
        return nil, errors.New("Error parsing 'vendorIdHeader' field " + _vendorIdHeaderErr.Error())
    }
    if vendorIdHeader != BACnetUnconfirmedServiceRequestIAm_VENDORIDHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestIAm_VENDORIDHEADER)) + " but got " + strconv.Itoa(int(vendorIdHeader)))
    }

    // Simple Field (vendorId)
    vendorId, _vendorIdErr := io.ReadUint8(8)
    if _vendorIdErr != nil {
        return nil, errors.New("Error parsing 'vendorId' field " + _vendorIdErr.Error())
    }

    // Create a partially initialized instance
    _child := &BACnetUnconfirmedServiceRequestIAm{
        ObjectType: objectType,
        ObjectInstanceNumber: objectInstanceNumber,
        MaximumApduLengthAcceptedLength: maximumApduLengthAcceptedLength,
        MaximumApduLengthAccepted: maximumApduLengthAccepted,
        SegmentationSupported: segmentationSupported,
        VendorId: vendorId,
        Parent: &BACnetUnconfirmedServiceRequest{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *BACnetUnconfirmedServiceRequestIAm) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Const Field (objectIdentifierHeader)
    _objectIdentifierHeaderErr := io.WriteUint8(8, 0xC4)
    if _objectIdentifierHeaderErr != nil {
        return errors.New("Error serializing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
    }

    // Simple Field (objectType)
    objectType := uint16(m.ObjectType)
    _objectTypeErr := io.WriteUint16(10, (objectType))
    if _objectTypeErr != nil {
        return errors.New("Error serializing 'objectType' field " + _objectTypeErr.Error())
    }

    // Simple Field (objectInstanceNumber)
    objectInstanceNumber := uint32(m.ObjectInstanceNumber)
    _objectInstanceNumberErr := io.WriteUint32(22, (objectInstanceNumber))
    if _objectInstanceNumberErr != nil {
        return errors.New("Error serializing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
    }

    // Const Field (maximumApduLengthAcceptedHeader)
    _maximumApduLengthAcceptedHeaderErr := io.WriteUint8(5, 0x04)
    if _maximumApduLengthAcceptedHeaderErr != nil {
        return errors.New("Error serializing 'maximumApduLengthAcceptedHeader' field " + _maximumApduLengthAcceptedHeaderErr.Error())
    }

    // Simple Field (maximumApduLengthAcceptedLength)
    maximumApduLengthAcceptedLength := uint8(m.MaximumApduLengthAcceptedLength)
    _maximumApduLengthAcceptedLengthErr := io.WriteUint8(3, (maximumApduLengthAcceptedLength))
    if _maximumApduLengthAcceptedLengthErr != nil {
        return errors.New("Error serializing 'maximumApduLengthAcceptedLength' field " + _maximumApduLengthAcceptedLengthErr.Error())
    }

    // Array Field (maximumApduLengthAccepted)
    if m.MaximumApduLengthAccepted != nil {
        for _, _element := range m.MaximumApduLengthAccepted {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'maximumApduLengthAccepted' field " + _elementErr.Error())
            }
        }
    }

    // Const Field (segmentationSupportedHeader)
    _segmentationSupportedHeaderErr := io.WriteUint8(8, 0x91)
    if _segmentationSupportedHeaderErr != nil {
        return errors.New("Error serializing 'segmentationSupportedHeader' field " + _segmentationSupportedHeaderErr.Error())
    }

    // Simple Field (segmentationSupported)
    segmentationSupported := uint8(m.SegmentationSupported)
    _segmentationSupportedErr := io.WriteUint8(8, (segmentationSupported))
    if _segmentationSupportedErr != nil {
        return errors.New("Error serializing 'segmentationSupported' field " + _segmentationSupportedErr.Error())
    }

    // Const Field (vendorIdHeader)
    _vendorIdHeaderErr := io.WriteUint8(8, 0x21)
    if _vendorIdHeaderErr != nil {
        return errors.New("Error serializing 'vendorIdHeader' field " + _vendorIdHeaderErr.Error())
    }

    // Simple Field (vendorId)
    vendorId := uint8(m.VendorId)
    _vendorIdErr := io.WriteUint8(8, (vendorId))
    if _vendorIdErr != nil {
        return errors.New("Error serializing 'vendorId' field " + _vendorIdErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *BACnetUnconfirmedServiceRequestIAm) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "objectType":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ObjectType = data
            case "objectInstanceNumber":
                var data uint32
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ObjectInstanceNumber = data
            case "maximumApduLengthAcceptedLength":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MaximumApduLengthAcceptedLength = data
            case "maximumApduLengthAccepted":
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.MaximumApduLengthAccepted = utils.ByteToInt8(_decoded[0:_len])
            case "segmentationSupported":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SegmentationSupported = data
            case "vendorId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.VendorId = data
            }
        }
    }
}

func (m *BACnetUnconfirmedServiceRequestIAm) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestIAm"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ObjectType, xml.StartElement{Name: xml.Name{Local: "objectType"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ObjectInstanceNumber, xml.StartElement{Name: xml.Name{Local: "objectInstanceNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MaximumApduLengthAcceptedLength, xml.StartElement{Name: xml.Name{Local: "maximumApduLengthAcceptedLength"}}); err != nil {
        return err
    }
    _encodedMaximumApduLengthAccepted := make([]byte, base64.StdEncoding.EncodedLen(len(m.MaximumApduLengthAccepted)))
    base64.StdEncoding.Encode(_encodedMaximumApduLengthAccepted, utils.Int8ToByte(m.MaximumApduLengthAccepted))
    if err := e.EncodeElement(_encodedMaximumApduLengthAccepted, xml.StartElement{Name: xml.Name{Local: "maximumApduLengthAccepted"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SegmentationSupported, xml.StartElement{Name: xml.Name{Local: "segmentationSupported"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.VendorId, xml.StartElement{Name: xml.Name{Local: "vendorId"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

