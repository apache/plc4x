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
    "github.com/apache/plc4x/plc4go/internal/plc4go/bacnetip/readwrite/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
)

type BacnetipXmlParserHelper struct {
}

func (m BacnetipXmlParserHelper) Parse(typeName string, xmlString string) (spi.Message, error) {
    switch typeName {
    case "APDU":
        var obj *model.APDU
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetTag":
        var obj *model.BACnetTag
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetTagWithContent":
        var obj *model.BACnetTagWithContent
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetError":
        var obj *model.BACnetError
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "NLM":
        var obj *model.NLM
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetConfirmedServiceRequest":
        var obj *model.BACnetConfirmedServiceRequest
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetAddress":
        var obj *model.BACnetAddress
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetConfirmedServiceACK":
        var obj *model.BACnetConfirmedServiceACK
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetUnconfirmedServiceRequest":
        var obj *model.BACnetUnconfirmedServiceRequest
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BACnetServiceAck":
        var obj *model.BACnetServiceAck
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "BVLC":
        var obj *model.BVLC
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    case "NPDU":
        var obj *model.NPDU
        err := xml.Unmarshal([]byte(xmlString), &obj)
        if err != nil {
            return nil, errors.New("error unmarshalling xml: " + err.Error())
        }
        return obj, nil
    }
    return nil, errors.New("Unsupported type " + typeName)
}
