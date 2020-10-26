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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

type KnxnetipParserHelper struct {
}

func (m KnxnetipParserHelper) Parse(typeName string, arguments []string, io *utils.ReadBuffer) (spi.Message, error) {
    switch typeName {
    case "CEMIAdditionalInformation":
        return CEMIAdditionalInformationParse(io)
    case "HPAIControlEndpoint":
        return HPAIControlEndpointParse(io)
    case "TunnelingResponseDataBlock":
        return TunnelingResponseDataBlockParse(io)
    case "ConnectionResponseDataBlock":
        return ConnectionResponseDataBlockParse(io)
    case "TunnelingRequestDataBlock":
        return TunnelingRequestDataBlockParse(io)
    case "KNXGroupAddress":
        numLevels, err := utils.StrToUint8(arguments[0])
        if err != nil {
            return nil, err
        }
        return KNXGroupAddressParse(io, numLevels)
    case "DIBDeviceInfo":
        return DIBDeviceInfoParse(io)
    case "DeviceConfigurationRequestDataBlock":
        return DeviceConfigurationRequestDataBlockParse(io)
    case "DeviceConfigurationAckDataBlock":
        return DeviceConfigurationAckDataBlockParse(io)
    case "DIBSuppSvcFamilies":
        return DIBSuppSvcFamiliesParse(io)
    case "ConnectionRequestInformation":
        return ConnectionRequestInformationParse(io)
    case "HPAIDiscoveryEndpoint":
        return HPAIDiscoveryEndpointParse(io)
    case "ProjectInstallationIdentifier":
        return ProjectInstallationIdentifierParse(io)
    case "KNXAddress":
        return KNXAddressParse(io)
    case "CEMIDataFrame":
        return CEMIDataFrameParse(io)
    case "ServiceId":
        return ServiceIdParse(io)
    case "KNXNetIPMessage":
        return KNXNetIPMessageParse(io)
    case "HPAIDataEndpoint":
        return HPAIDataEndpointParse(io)
    case "RelativeTimestamp":
        return RelativeTimestampParse(io)
    case "CEMI":
        size, err := utils.StrToUint8(arguments[0])
        if err != nil {
            return nil, err
        }
        return CEMIParse(io, size)
    case "MACAddress":
        return MACAddressParse(io)
    case "CEMIFrame":
        return CEMIFrameParse(io)
    case "DeviceStatus":
        return DeviceStatusParse(io)
    case "IPAddress":
        return IPAddressParse(io)
    }
    return nil, errors.New("Unsupported type " + typeName)
}
