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
	"github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
)

type KnxnetipXmlParserHelper struct {
}

func (m KnxnetipXmlParserHelper) Parse(typeName string, xmlString string) (spi.Message, error) {
	switch typeName {
	case "HPAIControlEndpoint":
		var obj *model.HPAIControlEndpoint
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "TunnelingResponseDataBlock":
		var obj *model.TunnelingResponseDataBlock
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "DeviceConfigurationAckDataBlock":
		var obj *model.DeviceConfigurationAckDataBlock
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "ConnectionRequestInformation":
		var obj *model.ConnectionRequestInformation
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "HPAIDiscoveryEndpoint":
		var obj *model.HPAIDiscoveryEndpoint
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "ProjectInstallationIdentifier":
		var obj *model.ProjectInstallationIdentifier
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "ServiceId":
		var obj *model.ServiceId
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "HPAIDataEndpoint":
		var obj *model.HPAIDataEndpoint
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "RelativeTimestamp":
		var obj *model.RelativeTimestamp
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "CEMI":
		var obj *model.CEMI
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "KnxNetIpMessage":
		var obj *model.KnxNetIpMessage
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "DeviceStatus":
		var obj *model.DeviceStatus
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "IPAddress":
		var obj *model.IPAddress
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "CEMIAdditionalInformation":
		var obj *model.CEMIAdditionalInformation
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "KnxAddress":
		var obj *model.KnxAddress
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "ConnectionResponseDataBlock":
		var obj *model.ConnectionResponseDataBlock
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "TunnelingRequestDataBlock":
		var obj *model.TunnelingRequestDataBlock
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "DIBDeviceInfo":
		var obj *model.DIBDeviceInfo
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "DeviceConfigurationRequestDataBlock":
		var obj *model.DeviceConfigurationRequestDataBlock
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "DIBSuppSvcFamilies":
		var obj *model.DIBSuppSvcFamilies
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "LDataFrame":
		var obj *model.LDataFrame
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "KnxGroupAddress":
		var obj *model.KnxGroupAddress
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	case "MACAddress":
		var obj *model.MACAddress
		err := xml.Unmarshal([]byte(xmlString), &obj)
		if err != nil {
			return nil, errors.New("error unmarshalling xml: " + err.Error())
		}
		return obj, nil
	}
	return nil, errors.New("Unsupported type " + typeName)
}
