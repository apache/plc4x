//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

[discriminatedType 'KNXNetIPMessage'
    [implicit      uint 8  'headerLength'    '6']
    [const         uint 8  'protocolVersion' '0x10']
    [discriminator uint 16 'msgType']
    [implicit      uint 16 'totalLength' 'lengthInBytes']
    [typeSwitch 'msgType'
        ['0x0201' SearchRequest
            [simple HPAIDiscoveryEndpoint 'hpaiIDiscoveryEndpoint']
        ]
        ['0x0202' SearchResponse
            [simple HPAIControlEndpoint 'hpaiControlEndpoint']
            [simple DIBDeviceInfo       'dibDeviceInfo']
            [simple DIBSuppSvcFamilies  'dibSuppSvcFamilies']
        ]
        ['0x0203' DescriptionRequest
            [simple HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x0204' DescriptionResponse
            [simple DIBDeviceInfo       'dibDeviceInfo']
            [simple DIBSuppSvcFamilies  'dibSuppSvcFamilies']
        ]
        ['0x0205' ConnectionRequest
            [simple HPAIDiscoveryEndpoint        'hpaiDiscoveryEndpoint']
            [simple HPAIDataEndpoint             'hpaiDataEndpoint']
            [simple ConnectionRequestInformation 'connectionRequestInformation']
        ]
        ['0x0206' ConnectionResponse
            [simple uint 8 'communicationChannelId']
            [simple uint 8 'status']
            [simple HPAIDataEndpoint            'hpaiDataEndpoint']
            [simple ConnectionResponseDataBlock 'connectionResponseDataBlock']
        ]
        ['0x0207' ConnectionStateRequest
            [simple   uint 8 'communicationChannelId']
            [reserved uint 8 '0x00']
            [simple   HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x0208' ConnectionStateResponse
            [simple uint 8 'communicationChannelId']
            [simple uint 8 'status']
        ]
        ['0x0209' DisconnectRequest
            [simple   uint 8 'communicationChannelId']
            [reserved uint 8 '0x00']
            [simple   HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x020A' DisconnectResponse
            [simple uint 8 'communicationChannelId']
            [simple uint 8 'status']
        ]
        ['0x0310' DeviceConfigurationRequest [uint 16 'totalLength']
            [simple DeviceConfigurationRequestDataBlock 'deviceConfigurationRequestDataBlock']
            [simple CEMI                                'cemi' ['totalLength - (6 + deviceConfigurationRequestDataBlock.lengthInBytes)']]
        ]
        ['0x0311' DeviceConfigurationAck
            [simple DeviceConfigurationAckDataBlock 'deviceConfigurationAckDataBlock']
        ]
        ['0x0420' TunnelingRequest [uint 16 'totalLength']
            [simple TunnelingRequestDataBlock 'tunnelingRequestDataBlock']
            [simple CEMI                      'cemi' ['totalLength - (6 + tunnelingRequestDataBlock.lengthInBytes)']]
        ]
        ['0x0421' TunnelingResponse
            [simple TunnelingResponseDataBlock 'tunnelingResponseDataBlock']
        ]
        ['0x0530' RoutingIndication
        ]
    ]
]

[type 'HPAIDiscoveryEndpoint'
    [implicit uint 8    'structureLength' 'lengthInBytes']
    [simple   uint 8    'hostProtocolCode']
    [simple   IPAddress 'ipAddress']
    [simple   uint 16   'ipPort']
]

[type 'HPAIControlEndpoint'
    [implicit uint 8    'structureLength' 'lengthInBytes']
    [simple   uint 8    'hostProtocolCode']
    [simple   IPAddress 'ipAddress']
    [simple   uint 16   'ipPort']
]

[type 'DIBDeviceInfo'
    [implicit uint 8       'structureLength' 'lengthInBytes']
    [simple   uint 8       'descriptionType']
    [simple   uint 8       'knxMedium']
    [simple   DeviceStatus 'deviceStatus']
    [simple   KNXAddress   'knxAddress']
    [simple   ProjectInstallationIdentifier 'projectInstallationIdentifier']
    [array    uint 8       'knxNetIpDeviceSerialNumber' count '6']
    [simple   IPAddress    'knxNetIpDeviceMulticastAddress']
    [simple   MACAddress   'knxNetIpDeviceMacAddress']
    [array    uint 8       'deviceFriendlyName'         count '30']
]

[type 'DIBSuppSvcFamilies'
    [implicit uint 8       'structureLength' 'lengthInBytes']
    [simple   uint 8       'descriptionType']
    [array    ServiceId    'serviceIds' count '3']
]

[type 'HPAIDataEndpoint'
    [implicit uint 8    'structureLength' 'lengthInBytes']
    [simple   uint 8    'hostProtocolCode']
    [simple   IPAddress 'ipAddress']
    [simple   uint 16   'ipPort']
]

[discriminatedType 'ConnectionRequestInformation'
    [implicit      uint 8    'structureLength' 'lengthInBytes']
    [discriminator uint 8    'connectionType']
    [typeSwitch 'connectionType'
        ['0x03' ConnectionRequestInformationDeviceManagement
        ]
        ['0x04' ConnectionRequestInformationTunnelConnection
            [simple   uint 8    'knxLayer']
            [reserved uint 8    '0x00']
        ]
    ]
]

[discriminatedType 'ConnectionResponseDataBlock'
    [implicit      uint 8     'structureLength' 'lengthInBytes']
    [discriminator uint 8     'connectionType']
    [typeSwitch 'connectionType'
        ['0x03' ConnectionResponseDataBlockDeviceManagement
        ]
        ['0x04' ConnectionResponseDataBlockTunnelConnection
            [simple KNXAddress 'knxAddress']
        ]
    ]
]

[type 'DeviceConfigurationRequestDataBlock'
    [implicit uint 8 'structureLength' 'lengthInBytes']
    [simple   uint 8 'communicationChannelId']
    [simple   uint 8 'sequenceCounter']
    [reserved uint 8 '0x00']
]

[type 'DeviceConfigurationAckDataBlock'
    [implicit uint 8 'structureLength' 'lengthInBytes']
    [simple   uint 8 'communicationChannelId']
    [simple   uint 8 'sequenceCounter']
    [simple   uint 8 'status']
]

[type 'TunnelingRequestDataBlock'
    [implicit uint 8 'structureLength' 'lengthInBytes']
    [simple   uint 8 'communicationChannelId']
    [simple   uint 8 'sequenceCounter']
    [reserved uint 8 '0x00']
]

[type 'TunnelingResponseDataBlock'
    [implicit uint 8 'structureLength' 'lengthInBytes']
    [simple   uint 8 'communicationChannelId']
    [simple   uint 8 'sequenceCounter']
    [simple   uint 8 'status']
]

[type 'IPAddress'
    [array uint 8 'addr' count '4']
]

[type 'MACAddress'
    [array uint 8 'addr' count '6']
]

[type 'KNXAddress'
    [simple uint 4 'mainGroup']
    [simple uint 4 'middleGroup']
    [simple uint 8 'subGroup']
]

[type 'DeviceStatus'
    [reserved uint 7 '0x00']
    [simple   bit    'programMode']
]

[type 'ProjectInstallationIdentifier'
    [simple uint 8 'projectNumber']
    [simple uint 8 'installationNumber']
]

[discriminatedType 'ServiceId'
    [discriminator uint 8 'serviceType']
    [typeSwitch 'serviceType'
        ['0x02' KnxNetIpCore
            [simple uint 8 'version']
        ]
        ['0x03' KnxNetIpDeviceManagement
            [simple uint 8 'version']
        ]
        ['0x04' KnxNetIpTunneling
            [simple uint 8 'version']
        ]
    ]
]

[discriminatedType 'CEMI' [uint 8 'size']
    [discriminator uint 8 'messageCode']
    [typeSwitch 'messageCode'
        ['0x10' CEMILRawReq
        ]
        ['0x11' CEMILDataReq
        ]
        ['0x13' CEMILPollDataReq
        ]

        ['0x25' CEMILPollDataCon
        ]
        ['0x29' CEMILDataInd
        ]
        ['0x2B' CEMILBusmonInd
            [simple uint 8                    'additionalInformationLength']
            [array  CEMIAdditionalInformation 'additionalInformation' length 'additionalInformationLength']
            [array  uint 8                    'rawFrame'              count  'size - (additionalInformationLength + 2)']
        ]
        ['0x2D' CEMILRawInd
        ]
        ['0x2E' CEMILDataCon
        ]
        ['0x2F' CEMILRawCon
        ]

        ['0xFC' CEMIMPropReadReq
            [simple uint 16 'interfaceObjectType']
            [simple uint  8 'objectInstance']
            [simple uint  8 'propertyId']
            [simple uint  4 'numberOfElements']
            [simple uint 12 'startIndex']
        ]
        ['0xFB' CEMIMPropReadCon
            [simple uint 16 'interfaceObjectType']
            [simple uint  8 'objectInstance']
            [simple uint  8 'propertyId']
            [simple uint  4 'numberOfElements']
            [simple uint 12 'startIndex']
        ]
    ]
]

[discriminatedType 'CEMIAdditionalInformation'
    [discriminator uint 8 'additionalInformationType']
    [typeSwitch 'additionalInformationType'
        ['0x03' CEMIAdditionalInformationBusmonitorInfo
            [implicit uint 8 'len' '1']
            [simple    bit    'frameErrorFlag']
            [simple    bit    'bitErrorFlag']
            [simple    bit    'parityErrorFlag']
            [simple    bit    'unknownFlag']
            [simple    bit    'lostFlag']
            [simple    uint 3 'sequenceNumber']
        ]
        ['0x04' CEMIAdditionalInformationRelativeTimestamp
            [implicit uint 8  'len' '2']
            [simple   RelativeTimestamp 'relativeTimestamp']
        ]
    ]
]

[type 'CEMIControlField1'
    [simple   bit    'standardFrame']
    [reserved uint 1 '0x00']
    [simple   bit    'doNotRepeat']
    [simple   bit    'broadcast']
    [simple   uint 2 'priority']
    [simple   bit    'ackRequested']
    [simple   bit    'error']
]

[type 'CEMIControlField2'
    [simple   bit    'groupAddress']
    [simple   uint 3 'hopCount']
    [simple   uint 3 'extendedFrameFormat']
]

[type 'RelativeTimestamp'
    [simple   uint 16 'timestamp']
]
