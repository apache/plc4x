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
    [implicit      uint 8  'headerLength' 'lengthInBytes']
    [const         uint 8  'protocolVersion' '0x10']
    [discriminator uint 16 'msgType']
    [implicit      uint 16 'totalLength' 'lengthInBytes']
    [typeSwitch 'msgType'
        ['0x0201' SearchRequest
            [field HPAIDiscoveryEndpoint 'hpaiIDiscoveryEndpoint']
        ]
        ['0x0202' SearchResponse
            [field HPAIControlEndpoint 'hpaiControlEndpoint']
            [field DIBDeviceInfo       'dibDeviceInfo']
            [field DIBSuppSvcFamilies  'dibSuppSvcFamilies']
        ]
        ['0x0203' DescriptionRequest
            [field HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x0204' DescriptionResponse
            [field DIBDeviceInfo       'dibDeviceInfo']
            [field DIBSuppSvcFamilies  'dibSuppSvcFamilies']
        ]
        ['0x0205' ConnectionRequest
            [field HPAIDiscoveryEndpoint        'hpaiDiscoveryEndpoint']
            [field HPAIDataEndpoint             'hpaiDataEndpoint']
            [field ConnectionRequestInformation 'connectionRequestInformation']
        ]
        ['0x0206' ConnectionResponse
            [field uint 8 'communicationChannelId']
            [field uint 8 'status']
            [field HPAIDataEndpoint            'hpaiDataEndpoint']
            [field ConnectionResponseDataBlock 'connectionResponseDataBlock']
        ]
        ['0x0207' ConnectionStateRequest
            [field    uint 8 'communicationChannelId']
            [reserved uint 8 '0x00']
            [field HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x0208' ConnectionStateResponse
            [field uint 8 'communicationChannelId']
            [field uint 8 'status']
        ]
        ['0x0209' DisconnectRequest
            [field    uint 8 'communicationChannelId']
            [reserved uint 8 '0x00']
            [field HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x020A' DisconnectResponse
            [field uint 8 'communicationChannelId']
            [field uint 8 'status']
        ]
        ['0x0310' DeviceConfigurationRequest
            [implicit uint 8 'structureLength' 'lengthInBytes']
            [field    uint 8 'communicationChannelId']
            [field    uint 8 'sequenceCounter']
            [reserved uint 8 '0x00']
            [field    CEMI 'cEMI']
        ]
        ['0x0311' DeviceConfigurationResponse
            [implicit uint 8 'structureLength' 'lengthInBytes']
            [field    uint 8 'communicationChannelId']
            [field    uint 8 'sequenceCounter']
            [field    uint 8 'status']
        ]
        ['0x0420' TunnelingRequest
            [implicit uint 8 'structureLength' 'lengthInBytes']
            [field    uint 8 'communicationChannelId']
            [field    uint 8 'sequenceCounter']
            [reserved uint 8 '0x00']
            [field    CEMI 'cEMI']
        ]
        ['0x0421' TunnelingResponse
            [implicit uint 8 'structureLength' 'lengthInBytes']
            [field    uint 8 'communicationChannelId']
            [field    uint 8 'sequenceCounter']
            [field    uint 8 'status']
        ]
    ]
]

[type 'HPAIDiscoveryEndpoint'
    [implicit uint 8    'structureLength' 'lengthInBytes']
    [field    uint 8    'hostProtocolCode']
    [field    IPAddress 'ipAddress']
    [field    uint 16   'ipPort']
]

[type 'HPAIControlEndpoint'
    [implicit uint 8    'structureLength' 'lengthInBytes']
    [field    uint 8    'hostProtocolCode']
    [field    IPAddress 'ipAddress']
    [field    uint 16   'ipPort']
]

[type 'DIBDeviceInfo'
    [implicit   uint 8       'structureLength' 'lengthInBytes']
    [field      uint 8       'descriptionType']
    [field      uint 8       'knxMedium']
    [field      DeviceStatus 'deviceStatus']
    [field      KNXAddress   'knxAddress']
    [field      ProjectInstallationIdentifier 'projectInstallationIdentifier']
    [arrayField uint 8       'knxNetIpDeviceSerialNumber' count '6']
    [field      IPAddress    'knxNetIpDeviceMulticastAddress']
    [field      MACAddress   'knxNetIpDeviceMacAddress']
    [arrayField uint 8       'deviceFriendlyName'         count '30']
]

[type 'DIBSuppSvcFamilies'
    [implicit   uint 8       'structureLength' 'lengthInBytes']
    [field      uint 8       'descriptionType']
    [arrayField ServiceId    'serviceIds' length '3']
]

[type 'HPAIDataEndpoint'
    [implicit uint 8    'structureLength' 'lengthInBytes']
    [field    uint 8    'hostProtocolCode']
    [field    IPAddress 'ipAddress']
    [field    uint 16   'ipPort']
]

[discriminatedType 'ConnectionRequestInformation'
    [implicit      uint 8    'structureLength' 'lengthInBytes']
    [discriminator uint 8    'connectionType']
    [typeSwitch 'connectionType'
        ['0x03' ConnectionRequestInformationDeviceManagement
        ]
        ['0x04' ConnectionRequestInformationTunnelConnection
            [field    uint 8    'knxLayer']
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
            [field         KNXAddress 'knxAddress']
        ]
    ]
]

[type 'IPAddress'
    [arrayField uint 8 'addr' count '4']
]

[type 'MACAddress'
    [arrayField uint 8 'addr' count '6']
]

[type 'KNXAddress'
    [field uint 4 'mainGroup']
    [field uint 4 'middleGroup']
    [field uint 8 'subGroup']
]

[type 'DeviceStatus'
    [reserved uint 7 '0x00']
    [field    bit  1 'programMode']
]

[type 'ProjectInstallationIdentifier'
    [field uint 8 'projectNumber']
    [field uint 8 'installationNumber']
]

[discriminatedType 'ServiceId'
    [discriminator uint 8 'serviceType']
    [typeSwitch 'serviceType'
        ['0x02' KnxNetIpCore
            [field uint 8 'version']
        ]
        ['0x03' KnxNetIpDeviceManagement
            [field uint 8 'version']
        ]
        ['0x04' KnxNetIpTunneling
            [field uint 8 'version']
        ]
    ]
]

[type 'CEMI'
    [field uint 8 'test']
]

