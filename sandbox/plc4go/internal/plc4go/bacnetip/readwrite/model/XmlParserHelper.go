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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
)

type BacnetipXmlParserHelper struct {
}

func (m BacnetipXmlParserHelper) Parse(typeName string, xmlString string) (spi.Message, error) {
    switch typeName {
    case "APDU":
        var obj APDU
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetTag":
        var obj BACnetTag
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetTagWithContent":
        var obj BACnetTagWithContent
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetError":
        var obj BACnetError
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "NLM":
        var obj NLM
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetConfirmedServiceRequest":
        var obj BACnetConfirmedServiceRequest
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetAddress":
        var obj BACnetAddress
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetConfirmedServiceACK":
        var obj BACnetConfirmedServiceACK
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetUnconfirmedServiceRequest":
        var obj BACnetUnconfirmedServiceRequest
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetServiceAck":
        var obj BACnetServiceAck
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BVLC":
        var obj BVLC
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "NPDU":
        var obj NPDU
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    }
    return nil, errors.New("Unsupported type " + typeName)
}
