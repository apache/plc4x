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

[discriminatedType 'KnxNetIpMessage'
    [implicit      uint 8  'headerLength'    '6']
    [const         uint 8  'protocolVersion' '0x10']
    [discriminator uint 16 'msgType']
    [implicit      uint 16 'totalLength'     'lengthInBytes']
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
            [simple   uint 8 'communicationChannelId']
            [simple   Status 'status']
            [optional HPAIDataEndpoint            'hpaiDataEndpoint'            'status == Status.NO_ERROR']
            [optional ConnectionResponseDataBlock 'connectionResponseDataBlock' 'status == Status.NO_ERROR']
        ]
        ['0x0207' ConnectionStateRequest
            [simple   uint 8 'communicationChannelId']
            [reserved uint 8 '0x00']
            [simple   HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x0208' ConnectionStateResponse
            [simple uint 8 'communicationChannelId']
            [simple Status 'status']
        ]
        ['0x0209' DisconnectRequest
            [simple   uint 8 'communicationChannelId']
            [reserved uint 8 '0x00']
            [simple   HPAIControlEndpoint 'hpaiControlEndpoint']
        ]
        ['0x020A' DisconnectResponse
            [simple uint 8 'communicationChannelId']
            [simple Status 'status']
        ]
        ['0x020B' UnknownMessage [uint 16 'totalLength']
            [array int 8 'unknownData' count 'totalLength - 6']
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
    [implicit uint 8           'structureLength' 'lengthInBytes']
    [simple   HostProtocolCode 'hostProtocolCode']
    [simple   IPAddress        'ipAddress']
    [simple   uint 16          'ipPort']
]

[type 'HPAIControlEndpoint'
    [implicit uint 8           'structureLength' 'lengthInBytes']
    [simple   HostProtocolCode 'hostProtocolCode']
    [simple   IPAddress        'ipAddress']
    [simple   uint 16          'ipPort']
]

[type 'DIBDeviceInfo'
    [implicit uint 8       'structureLength' 'lengthInBytes']
    [simple   uint 8       'descriptionType']
    [simple   KnxMedium    'knxMedium']
    [simple   DeviceStatus 'deviceStatus']
    [simple   KnxAddress   'knxAddress']
    [simple   ProjectInstallationIdentifier 'projectInstallationIdentifier']
    [array    int 8        'knxNetIpDeviceSerialNumber' count '6']
    [simple   IPAddress    'knxNetIpDeviceMulticastAddress']
    [simple   MACAddress   'knxNetIpDeviceMacAddress']
    [array    int 8        'deviceFriendlyName'         count '30']
]

[type 'DIBSuppSvcFamilies'
    [implicit uint 8       'structureLength' 'lengthInBytes']
    [simple   uint 8       'descriptionType']
    [array    ServiceId    'serviceIds' count '3']
]

[type 'HPAIDataEndpoint'
    [implicit uint 8           'structureLength' 'lengthInBytes']
    [simple   HostProtocolCode 'hostProtocolCode']
    [simple   IPAddress        'ipAddress']
    [simple   uint 16          'ipPort']
]

[discriminatedType 'ConnectionRequestInformation'
    [implicit      uint 8    'structureLength' 'lengthInBytes']
    [discriminator uint 8    'connectionType']
    [typeSwitch 'connectionType'
        ['0x03' ConnectionRequestInformationDeviceManagement
        ]
        ['0x04' ConnectionRequestInformationTunnelConnection
            [simple   KnxLayer  'knxLayer']
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
            [simple KnxAddress 'knxAddress']
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
    [simple   Status 'status']
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
    [simple   Status 'status']
]

[type 'IPAddress'
    [array int 8 'addr' count '4']
]

[type 'MACAddress'
    [array int 8 'addr' count '6']
]

[type 'KnxAddress'
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
        // TODO: Check if this shouldn't be KnxNetIp instead of KnxNet
        ['0x06' KnxNetRemoteLogging
            [simple uint 8 'version']
        ]
        // TODO: Check if this shouldn't be KnxNetIp instead of KnxNet
        ['0x07' KnxNetRemoteConfigurationAndDiagnosis
            [simple uint 8 'version']
        ]
        // TODO: Check if this shouldn't be KnxNetIp instead of KnxNet
        ['0x08' KnxNetObjectServer
            [simple uint 8 'version']
        ]
    ]
]

// The CEMI part is described in the document
// "03_06_03 EMI_IMI v01.03.03 AS" Page 6ff
// NOTE: When inspecting traffic in WireShark it seems they got the
// standard/extended frame thing wrong. When comparing to the spec most
// normal traffic is actually extended frames.
[discriminatedType 'CEMI' [uint 8 'size']
    [discriminator uint 8 'messageCode']
    [typeSwitch 'messageCode'
        ['0x2B' LBusmonInd
            [simple   uint 8                    'additionalInformationLength']
            [array    CEMIAdditionalInformation 'additionalInformation' length 'additionalInformationLength']
            [simple   LDataFrame                'dataFrame']
            [optional uint 8                    'crc'                   'dataFrame.notAckFrame']
        ]

        // Page 72ff
        ['0x11' LDataReq
            [simple   uint 8                    'additionalInformationLength']
            [array    CEMIAdditionalInformation 'additionalInformation' length 'additionalInformationLength']
            [simple   LDataFrame                'dataFrame']
        ]
        ['0x29' LDataInd
            [simple   uint 8                    'additionalInformationLength']
            [array    CEMIAdditionalInformation 'additionalInformation' length 'additionalInformationLength']
            [simple   LDataFrame                'dataFrame']
        ]
        ['0x2E' LDataCon
            [simple   uint 8                    'additionalInformationLength']
            [array    CEMIAdditionalInformation 'additionalInformation' length 'additionalInformationLength']
            [simple   LDataFrame                'dataFrame']
        ]

        ['0x10' LRawReq
        ]
        ['0x2D' LRawInd
        ]
        ['0x2F' LRawCon
        ]

        ['0x13' LPollDataReq
        ]
        ['0x25' LPollDataCon
        ]

        ['0x41' TDataConnectedReq
        ]
        ['0x89' TDataConnectedInd
        ]

        ['0x4A' TDataIndividualReq
        ]
        ['0x94' TDataIndividualInd
        ]

        ['0xFC' MPropReadReq
            [simple uint 16 'interfaceObjectType']
            [simple uint  8 'objectInstance']
            [simple uint  8 'propertyId']
            [simple uint  4 'numberOfElements']
            [simple uint 12 'startIndex']
        ]
        ['0xFB' MPropReadCon
            [simple uint 16 'interfaceObjectType']
            [simple uint  8 'objectInstance']
            [simple uint  8 'propertyId']
            [simple uint  4 'numberOfElements']
            [simple uint 12 'startIndex']
            [simple uint 16 'unknown']
        ]

        ['0xF6' MPropWriteReq
        ]
        ['0xF5' MPropWriteCon
        ]

        ['0xF7' MPropInfoInd
        ]

        ['0xF8' MFuncPropCommandReq
        ]
        ['0xF9' MFuncPropStateReadReq
        ]
        ['0xFA' MFuncPropCon
        ]

        ['0xF1' MResetReq
        ]
        ['0xF0' MResetInd
        ]
    ]
]

[discriminatedType 'CEMIAdditionalInformation'
    [discriminator uint 8 'additionalInformationType']
    [typeSwitch 'additionalInformationType'
        ['0x03' CEMIAdditionalInformationBusmonitorInfo
            [const     uint 8 'len' '1']
            [simple    bit    'frameErrorFlag']
            [simple    bit    'bitErrorFlag']
            [simple    bit    'parityErrorFlag']
            [simple    bit    'unknownFlag']
            [simple    bit    'lostFlag']
            [simple    uint 3 'sequenceNumber']
        ]
        ['0x04' CEMIAdditionalInformationRelativeTimestamp
            [const    uint 8            'len' '2']
            [simple   RelativeTimestamp 'relativeTimestamp']
        ]
    ]
]

// The CEMI part is described in the document "03_06_03 EMI_IMI v01.03.03 AS" Page 73
// "03_02_02 Communication Medium TP1 v01.02.02 AS" Page 27
[discriminatedType 'LDataFrame'
    [discriminator bit          'extendedFrame']
    [discriminator bit          'polling']
    [simple        bit          'repeated']
    [simple        bit          'notAckFrame']
    [enum          CEMIPriority 'priority']
    [simple        bit          'acknowledgeRequested']
    [simple        bit          'errorFlag']
    [typeSwitch 'extendedFrame','polling'
       // Page 28ff
        ['false','false' LDataFrameData
            [simple   KnxAddress   'sourceAddress']
            [array    int 8        'destinationAddress' count '2']
            [simple   bit          'groupAddress']
            [simple   uint 3       'hopCount']
            [simple   Apdu         'apdu']
        ]
        // Page 29ff
        ['true','false' LDataFrameDataExt
            [simple   bit          'groupAddress']
            [simple   uint 3       'hopCount']
            [simple   uint 4       'extendedFrameFormat']
            [simple   KnxAddress   'sourceAddress']
            [array    int 8        'destinationAddress' count '2']
            [simple   Apdu         'apdu']
        ]
        // Page 31ff
        ['true','true' LDataFramePollingData
            [simple   KnxAddress   'sourceAddress']
            [array    int 8        'targetAddress' count '2']
            [reserved uint 4       '0x00']
            [simple   uint 6       'numberExpectedPollData']
        ]
        // Page 31ff
        ['false','true' LDataFramePollingData
            [simple   KnxAddress   'sourceAddress']
            [array    int 8        'targetAddress' count '2']
            [reserved uint 4       '0x00']
            [simple   uint 6       'numberExpectedPollData']
        ]
    ]
]

[discriminatedType 'Apdu'
    [simple   uint 8      'dataLength']
    // 10_01 Logical Tag Extended v01.02.01 AS.pdf Page 74ff
    [discriminator uint 1 'control']
    [simple        bit    'numbered']
    [simple        uint 4 'counter']
    [typeSwitch 'control'
        ['1' ApduControlContainer
            [simple ApduControl 'controlApdu']
        ]
        ['0' ApduDataContainer [uint 8 'dataLength']
            [simple ApduData 'dataApdu' ['dataLength']]
        ]
    ]
]

[discriminatedType 'ApduControl'
    [discriminator uint 2 'controlType']
    [typeSwitch 'controlType'
        ['0x0' ApduControlConnect
        ]
        ['0x1' ApduControlDisconnect
        ]
        ['0x2' ApduControlAck
        ]
        ['0x3' ApduControlNack
        ]
    ]
]

[discriminatedType 'ApduData' [uint 8 'dataLength']
    [discriminator uint 4 'apciType']
    // 03_03_07 Application Layer v01.06.02 AS Page 9ff
    [typeSwitch 'apciType'
        ['0x0' ApduDataGroupValueRead
        ]
        ['0x1' ApduDataGroupValueResponse
        ]
        ['0x2' ApduDataGroupValueWrite [uint 8 'dataLength']
            [simple int 6 'dataFirstByte']
            [array  int 8 'data' count  '(dataLength < 1) ? 0 : dataLength - 1']
        ]
        ['0x3' ApduDataIndividualAddressWrite
        ]
        ['0x4' ApduDataIndividualAddressRead
        ]
        ['0x5' ApduDataIndividualAddressResponse
        ]
        ['0x6' ApduDataAdcRead
        ]
        // In case of this type the following 6 bits contain more detailed information
        ['0x7' ApduDataAdcResponse
        ]
        ['0x8' ApduDataMemoryRead
        ]
        ['0x9' ApduDataMemoryResponse
        ]
        ['0xA' ApduDataMemoryWrite
        ]
        // In case of this type the following 6 bits contain more detailed information
        ['0xB' ApduDataUserMessage
        ]
        ['0xC' ApduDataDeviceDescriptorRead
            [simple uint 6 'descriptorType']
        ]
        ['0xD' ApduDataDeviceDescriptorResponse [uint 8 'dataLength']
            [simple uint 6 'descriptorType']
            [array  int 8 'data' count  '(dataLength < 1) ? 0 : dataLength - 1']
        ]
        ['0xE' ApduDataRestart
        ]
        ['0xF' ApduDataOther [uint 8 'dataLength']
            [simple ApduDataExt 'extendedApdu' ['dataLength']]
        ]
    ]
]

// 03_03_07 Application Layer v01.06.02 AS Page 9ff
[discriminatedType 'ApduDataExt' [uint 8 'length']
    [discriminator uint 6 'extApciType']
    [typeSwitch 'extApciType'
        ['0x00' ApduDataExtOpenRoutingTableRequest
        ]
        ['0x01' ApduDataExtReadRoutingTableRequest
        ]
        ['0x02' ApduDataExtReadRoutingTableResponse
        ]
        ['0x03' ApduDataExtWriteRoutingTableRequest
        ]
        ['0x08' ApduDataExtReadRouterMemoryRequest
        ]
        ['0x09' ApduDataExtReadRouterMemoryResponse
        ]
        ['0x0A' ApduDataExtWriteRouterMemoryRequest
        ]
        ['0x0D' ApduDataExtReadRouterStatusRequest
        ]
        ['0x0E' ApduDataExtReadRouterStatusResponse
        ]
        ['0x0F' ApduDataExtWriteRouterStatusRequest
        ]

        ['0x10' ApduDataExtMemoryBitWrite
        ]

        ['0x11' ApduDataExtAuthorizeRequest
        ]
        ['0x12' ApduDataExtAuthorizeResponse
        ]
        ['0x13' ApduDataExtKeyWrite
        ]
        ['0x14' ApduDataExtKeyResponse
        ]

        ['0x15' ApduDataExtPropertyValueRead
            [simple uint 8  'objectIndex']
            [simple uint 8  'propertyId']
            [simple uint 4  'count']
            [simple uint 12 'index']
        ]
        ['0x16' ApduDataExtPropertyValueResponse [uint 8 'length']
            [simple uint 8  'objectIndex']
            [simple uint 8  'propertyId']
            [simple uint 4  'count']
            [simple uint 12 'index']
            [array  uint 8 'data' count 'length - 5']
        ]
        ['0x17' ApduDataExtPropertyValueWrite
        ]
        ['0x18' ApduDataExtPropertyDescriptionRead
        ]
        ['0x19' ApduDataExtPropertyDescriptionResponse
        ]

        ['0x1A' ApduDataExtNetworkParameterRead
        ]
        ['0x1B' ApduDataExtNetworkParameterResponse
        ]

        ['0x1C' ApduDataExtIndividualAddressSerialNumberRead
        ]
        ['0x1D' ApduDataExtIndividualAddressSerialNumberResponse
        ]
        ['0x1E' ApduDataExtIndividualAddressSerialNumberWrite
        ]

        ['0x20' ApduDataExtDomainAddressWrite
        ]
        ['0x21' ApduDataExtDomainAddressRead
        ]
        ['0x22' ApduDataExtDomainAddressResponse
        ]
        ['0x23' ApduDataExtDomainAddressSelectiveRead
        ]

        ['0x24' ApduDataExtNetworkParameterWrite
        ]

        ['0x25' ApduDataExtLinkRead
        ]
        ['0x26' ApduDataExtLinkResponse
        ]
        ['0x27' ApduDataExtLinkWrite
        ]

        ['0x28' ApduDataExtGroupPropertyValueRead
        ]
        ['0x29' ApduDataExtGroupPropertyValueResponse
        ]
        ['0x2A' ApduDataExtGroupPropertyValueWrite
        ]
        ['0x2B' ApduDataExtGroupPropertyValueInfoReport
        ]

        ['0x2C' ApduDataExtDomainAddressSerialNumberRead
        ]
        ['0x2D' ApduDataExtDomainAddressSerialNumberResponse
        ]
        ['0x2E' ApduDataExtDomainAddressSerialNumberWrite
        ]

        ['0x30' ApduDataExtFileStreamInfoReport
        ]

    ]
]

[type 'RelativeTimestamp'
    [simple   uint 16 'timestamp']
]

[discriminatedType 'KnxGroupAddress' [uint 2 'numLevels']
    [typeSwitch 'numLevels'
        ['1' KnxGroupAddressFreeLevel
            [simple uint 16 'subGroup']
        ]
        ['2' KnxGroupAddress2Level
            [simple uint 5  'mainGroup']
            [simple uint 11 'subGroup']
        ]
        ['3' KnxGroupAddress3Level
            [simple uint 5 'mainGroup']
            [simple uint 3 'middleGroup']
            [simple uint 8 'subGroup']
        ]
    ]
]

[enum uint 2 'CEMIPriority'
    ['0x0' SYSTEM]
    ['0x1' NORMAL]
    ['0x2' URGENT]
    ['0x3' LOW]
]

[enum uint 8 'Status'
    ['0x00' NO_ERROR]
    ['0x01' PROTOCOL_TYPE_NOT_SUPPORTED]
    ['0x02' UNSUPPORTED_PROTOCOL_VERSION]
    ['0x04' OUT_OF_ORDER_SEQUENCE_NUMBER]
    ['0x21' INVALID_CONNECTION_ID]
    ['0x22' CONNECTION_TYPE_NOT_SUPPORTED]
    ['0x23' CONNECTION_OPTION_NOT_SUPPORTED]
    ['0x24' NO_MORE_CONNECTIONS]
    ['0x25' NO_MORE_UNIQUE_CONNECTIONS]
    ['0x26' DATA_CONNECTION]
    ['0x27' KNX_CONNECTION]
    ['0x29' TUNNELLING_LAYER_NOT_SUPPORTED]
]

[enum uint 8 'HostProtocolCode'
    ['0x01' IPV4_UDP]
    ['0x02' IPV4_TCP]
]

// The mode in which the connection should be established:
// TUNNEL_LINK_LAYER The gateway assigns a unique KNX address to the client.
//                   The client can then actively participate in communicating
//                   with other KNX devices.
// TUNNEL_RAW        The gateway will just pass along the packets and not
//                   automatically generate Ack frames for the packets it
//                   receives for a given client.
// TUNNEL_BUSMONITOR The client becomes a passive participant and all frames
//                   on the KNX bus get forwarded to the client. Only one
//                   Busmonitor connection is allowed at any given time.
[enum uint 8 'KnxLayer'
    ['0x02' TUNNEL_LINK_LAYER]
    ['0x04' TUNNEL_RAW]
    ['0x80' TUNNEL_BUSMONITOR]
]

[enum uint 8 'KnxMedium'
    ['0x01' MEDIUM_RESERVED_1]
    ['0x02' MEDIUM_TP1]
    ['0x04' MEDIUM_PL110]
    ['0x08' MEDIUM_RESERVED_2]
    ['0x10' MEDIUM_RF]
    ['0x20' MEDIUM_KNX_IP]
]

[enum uint 8 'SupportedPhysicalMedia' [string '-1' 'description',                                                    bit 'knxSupport']
    ['0x00' OTHER                     ['used_for_undefined_physical_medium',                                    'true']]
    ['0x01' OIL_METER                 ['measures_volume_of_oil',                                                'true']]
    ['0x02' ELECTRICITY_METER         ['measures_electric_energy',                                              'true']]
    ['0x03' GAS_METER                 ['measures_volume_of_gaseous_energy',                                     'true']]
    ['0x04' HEAT_METER                ['heat_energy_measured_in_outlet_pipe',                                   'true']]
    ['0x05' STEAM_METER               ['measures_weight_of_hot_steam',                                          'true']]
    ['0x06' WARM_WATER_METER          ['measured_heated_water_volume',                                          'true']]
    ['0x07' WATER_METER               ['measured_water_volume',                                                 'true']]
    ['0x08' HEAT_COST_ALLOCATOR       ['measured_relative_cumulated_heat_consumption',                          'true']]
    ['0x09' COMPRESSED_AIR            ['measures_weight_of_compressed_air',                                     'false']]
    ['0x0A' COOLING_LOAD_METER_INLET  ['cooling_energy_measured_in_inlet_pipe',                                 'true']]
    ['0x0B' COOLING_LOAD_METER_OUTLET ['cooling_energy_measured_in_outlet_pipe',                                'true']]
    ['0x0C' HEAT_INLET                ['heat_energy_measured_in_inlet_pipe',                                    'true']]
    ['0x0D' HEAT_AND_COOL             ['measures_both_heat_and_cool',                                           'true']]
    ['0x0E' BUS_OR_SYSTEM             ['no_meter',                                                              'false']]
    ['0x0F' UNKNOWN_DEVICE_TYPE       ['used_for_undefined_physical_medium',                                    'false']]
    ['0x20' BREAKER                   ['status_of_electric_energy_supply',                                      'true']]
    ['0x21' VALVE                     ['status_of_supply_of_Gas_or_water',                                      'true']]
    ['0x28' WASTE_WATER_METER         ['measured_volume_of_disposed_water',                                     'true']]
    ['0x29' GARBAGE                   ['measured_weight_of_disposed_rubbish',                                   'true']]
    ['0x37' RADIO_CONVERTER           ['enables_the_radio_transmission_of_a_meter_without_a_radio_interface',   'false']]
]

// The definition of the constants for medium type in the device descriptor differs from that of the other parts
// 03_05_01 Resources v01.09.03 AS.pdf Page 22
[enum uint 4 'DeviceDescriptorMediumType'
    ['0x0' TP1      ]
    ['0x1' PL110    ]
    ['0x2' RF       ]
    ['0x3' TP0      ]
    ['0x4' PL132    ]
    ['0x5' KNX_IP   ]
]

// 03_05_01 Resources v01.09.03 AS.pdf Page 22
[enum uint 4 'FirmwareType' [uint 8 'code']
    ['0x1' NONE                      ['0xAF']]
    ['0x2' BCU_1                     ['0x00']]
    ['0x3' BCU_1_SYSTEM_1            ['0x01']]
    ['0x4' BCU_2_SYSTEM_2            ['0x02']]
    ['0x5' BIM_M112                  ['0x70']]
    ['0x6' SYSTEM_B                  ['0x7B']]
    ['0x7' IR_DECODER                ['0x81']]
    ['0x8' MEDIA_COUPLER_PL_TP       ['0x90']]
    ['0x9' COUPLER                   ['0x91']]
    ['0xA' RF_BI_DIRECTIONAL_DEVICES ['0x01']]
    ['0xB' RF_UNI_DIRECTIONAL_DEVICES['0x11']]
    ['0xC' SYSTEM_300                ['0x30']]
    ['0xD' SYSTEM_7                  ['0x70']]
]

// Helper enum that binds the combinations of medium type and firmware
// type to the pre-defined constants the spec defines
// 03_05_01 Resources v01.09.03 AS.pdf Page 22
[enum uint 16 'DeviceDescriptorType0'   [DeviceDescriptorMediumType 'mediumType',   FirmwareType 'firmwareType'               ]
    ['0x0010' TP1_BCU_1_SYSTEM_1_0      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_1_SYSTEM_1'            ]]
    ['0x0011' TP1_BCU_1_SYSTEM_1_1      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_1_SYSTEM_1'            ]]
    ['0x0012' TP1_BCU_1_SYSTEM_1_2      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_1_SYSTEM_1'            ]]
    ['0x0013' TP1_BCU_1_SYSTEM_1_3      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_1_SYSTEM_1'            ]]
    ['0x0020' TP1_BCU_2_SYSTEM_2_0      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_2_SYSTEM_2'            ]]
    ['0x0021' TP1_BCU_2_SYSTEM_2_1      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_2_SYSTEM_2'            ]]
    ['0x0025' TP1_BCU_2_SYSTEM_2_5      ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BCU_2_SYSTEM_2'            ]]
    ['0x0300' TP1_SYSTEM_300            ['DeviceDescriptorMediumType.TP1',          'FirmwareType.SYSTEM_300'                ]]
    ['0x0700' TP1_BIM_M112_0            ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BIM_M112'                  ]]
    ['0x0701' TP1_BIM_M112_1            ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BIM_M112'                  ]]
    ['0x0705' TP1_BIM_M112_5            ['DeviceDescriptorMediumType.TP1',          'FirmwareType.BIM_M112'                  ]]
    ['0x07B0' TP1_SYSTEM_B              ['DeviceDescriptorMediumType.TP1',          'FirmwareType.SYSTEM_B'                  ]]
    ['0x0810' TP1_IR_DECODER_0          ['DeviceDescriptorMediumType.TP1',          'FirmwareType.IR_DECODER'                ]]
    ['0x0811' TP1_IR_DECODER_1          ['DeviceDescriptorMediumType.TP1',          'FirmwareType.IR_DECODER'                ]]
    ['0x0910' TP1_COUPLER_0             ['DeviceDescriptorMediumType.TP1',          'FirmwareType.COUPLER'                   ]]
    ['0x0911' TP1_COUPLER_1             ['DeviceDescriptorMediumType.TP1',          'FirmwareType.COUPLER'                   ]]
    ['0x0912' TP1_COUPLER_2             ['DeviceDescriptorMediumType.TP1',          'FirmwareType.COUPLER'                   ]]
    ['0x091A' TP1_KNXNETIP_ROUTER       ['DeviceDescriptorMediumType.TP1',          'FirmwareType.COUPLER'                   ]]
    ['0x0AFD' TP1_NONE_D                ['DeviceDescriptorMediumType.TP1',          'FirmwareType.NONE'                      ]]
    ['0x0AFE' TP1_NONE_E                ['DeviceDescriptorMediumType.TP1',          'FirmwareType.NONE'                      ]]
    ['0x1012' PL110_BCU_1_2             ['DeviceDescriptorMediumType.PL110',        'FirmwareType.BCU_1_SYSTEM_1'            ]]
    ['0x1013' PL110_BCU_1_3             ['DeviceDescriptorMediumType.PL110',        'FirmwareType.BCU_1_SYSTEM_1'            ]]
    ['0x17B0' PL110_SYSTEM_B            ['DeviceDescriptorMediumType.PL110',        'FirmwareType.SYSTEM_B'                  ]]
    ['0x1900' PL110_MEDIA_COUPLER_PL_TP ['DeviceDescriptorMediumType.PL110',        'FirmwareType.MEDIA_COUPLER_PL_TP'       ]]
    ['0x2010' RF_BI_DIRECTIONAL_DEVICES ['DeviceDescriptorMediumType.RF',           'FirmwareType.RF_BI_DIRECTIONAL_DEVICES' ]]
    ['0x2110' RF_UNI_DIRECTIONAL_DEVICES['DeviceDescriptorMediumType.RF',           'FirmwareType.RF_UNI_DIRECTIONAL_DEVICES']]
    ['0x3012' TP0_BCU_1                 ['DeviceDescriptorMediumType.TP0',          'FirmwareType.BCU_1'                     ]]
    ['0x4012' PL132_BCU_1               ['DeviceDescriptorMediumType.PL132',        'FirmwareType.BCU_1'                     ]]
    ['0x5705' KNX_IP_SYSTEM7            ['DeviceDescriptorMediumType.KNX_IP',       'FirmwareType.SYSTEM_7'                  ]]
]

// 03_05_01 Resources v01.09.03 AS.pdf Page 23ff
[type 'DeviceDescriptorType2'
    // Same manufacturer id as used elsewhere (Assigned by KNX Association)
    [simple uint 16            'manufacturerId' ]
    // Manufacturer specific device type id
    [simple uint 16            'deviceType'     ]
    // Manufacturer specific device type version
    [simple uint 8             'version'        ]
    // Indicates the Network Management procedures based on A_Link_Read-service are supported
    [simple bit                'readSupported'  ]
    // Indicates the Network Management procedures based on A_Link_Write-service are supported
    [simple bit                'writeSupported' ]
    [simple uint 6             'logicalTagBase' ]
    [simple ChannelInformation 'channelInfo1'   ]
    [simple ChannelInformation 'channelInfo2'   ]
    [simple ChannelInformation 'channelInfo3'   ]
    [simple ChannelInformation 'channelInfo4'   ]
]

// 03_05_01 Resources v01.09.03 AS.pdf Page 24
[type 'ChannelInformation'
    [simple uint 3  'numChannels']
    [simple uint 13 'channelCode']
]

// Switch constants taken from the generated KnxPropertyDataType type
// Representation looked up in:
// - 03_05_01 Resources v01.09.03 AS.pdf
// - 03_07_03 Standardized Identifier Tables v01.03.01 AS.pdf
// - 03_07_02 Datapoint Types v01.08.02 AS.pdf
[dataIo 'KnxProperty' [KnxPropertyDataType 'propertyType', uint 8 'dataLength']
    [typeSwitch 'propertyType'
        ['KnxPropertyDataType.PDT_CONTROL' BOOL
            [reserved uint 7        '0x00']
            [simple   bit           'value']
        ]
        ['KnxPropertyDataType.PDT_CHAR' SINT
            [simple   int 8         'value']
        ]
        ['KnxPropertyDataType.PDT_UNSIGNED_CHAR' USINT
            [simple   uint 8        'value']
        ]
        ['KnxPropertyDataType.PDT_INT' INT
            [simple   int 16        'value']
        ]
        ['KnxPropertyDataType.PDT_UNSIGNED_INT' UINT
            [simple   uint 16       'value']
        ]
        ['KnxPropertyDataType.PDT_KNX_FLOAT' REAL
            [simple   float 4.11    'value']
        ]
        ['KnxPropertyDataType.PDT_DATE' Struct
            [reserved uint 3 '0x00']
            [simple uint 5 'dayOfMonth']
            [reserved uint 4 '0x00']
            [simple uint 4 'month']
            [reserved uint 1 '0x00']
            [simple uint 7 'year']        ]
        ['KnxPropertyDataType.PDT_TIME' Struct
            [simple uint 3 'day']
            [simple uint 5 'hour']
            [reserved uint 2 '0x00']
            [simple uint 6 'minutes']
            [reserved uint 2 '0x00']
            [simple uint 6 'seconds']
        ]
        ['KnxPropertyDataType.PDT_LONG' DINT
            [simple   int 32        'value']
        ]
        ['KnxPropertyDataType.PDT_UNSIGNED_LONG' UDINT
            [simple   uint 32       'value']
        ]
        ['KnxPropertyDataType.PDT_FLOAT' REAL
            [simple   float 8.23    'value']
        ]
        ['KnxPropertyDataType.PDT_DOUBLE' LREAL
            [simple   float 11.52   'value']
        ]
        ['KnxPropertyDataType.PDT_CHAR_BLOCK' List
            [array uint 8           'value' count '10']
        ]
        ['KnxPropertyDataType.PDT_POLL_GROUP_SETTINGS' Struct
            [array    uint 8        'groupAddress' count '2']
            [simple   bit           'disable']
            [reserved uint 3        '0x0']
            [simple   uint 4        'pollingSoftNr']
        ]
        ['KnxPropertyDataType.PDT_SHORT_CHAR_BLOCK' List
            [array uint 8           'value' count '5']
        ]
        ['KnxPropertyDataType.PDT_DATE_TIME' Struct
            [simple uint 8 'year']
            [reserved uint 4 '0x00']
            [simple uint 4 'month']
            [reserved uint 3 '0x00']
            [simple uint 5 'dayofmonth']
            [simple uint 3 'dayofweek']
            [simple uint 5 'hourofday']
            [reserved uint 2 '0x00']
            [simple uint 6 'minutes']
            [reserved uint 2 '0x00']
            [simple uint 6 'seconds']
            [simple bit 'fault']
            [simple bit 'workingDay']
            [simple bit 'noWd']
            [simple bit 'noYear']
            [simple bit 'noDate']
            [simple bit 'noDayOfWeek']
            [simple bit 'noTime']
            [simple bit 'standardSummerTime']
            [simple bit 'qualityOfClock']
            [reserved uint 7 '0x00']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_01' List
            [array uint 8           'value' count '1']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_02' List
            [array uint 8           'value' count '2']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_03' List
            [array uint 8           'value' count '3']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_04' List
            [array uint 8           'value' count '4']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_05' List
            [array uint 8           'value' count '5']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_06' List
            [array uint 8           'value' count '6']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_07' List
            [array uint 8           'value' count '7']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_08' List
            [array uint 8           'value' count '8']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_09' List
            [array uint 8           'value' count '9']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_10' List
            [array uint 8           'value' count '10']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_11' List
            [array uint 8           'value' count '11']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_12' List
            [array uint 8           'value' count '12']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_13' List
            [array uint 8           'value' count '13']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_14' List
            [array uint 8           'value' count '14']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_15' List
            [array uint 8           'value' count '15']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_16' List
            [array uint 8           'value' count '16']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_17' List
            [array uint 8           'value' count '17']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_18' List
            [array uint 8           'value' count '18']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_19' List
            [array uint 8           'value' count '19']
        ]
        ['KnxPropertyDataType.PDT_GENERIC_20' List
            [array uint 8           'value' count '20']
        ]
        // Defaults to PDT_VARIABLE_LENGTH
        //['KnxPropertyDataType.PDT_UTF_8'
        //]
        ['KnxPropertyDataType.PDT_VERSION' Struct
            [simple uint 5 'magicNumber']
            [simple uint 5 'versionNumber']
            [simple uint 6 'revisionNumber']
        ]
        ['KnxPropertyDataType.PDT_ALARM_INFO' Struct
            [simple uint 8 'logNumber']
            [simple uint 8 'alarmPriority']
            [simple uint 8 'applicationArea']
            [simple uint 8 'errorClass']
            [reserved uint 4 '0x00']
            [simple bit 'errorcodeSup']
            [simple bit 'alarmtextSup']
            [simple bit 'timestampSup']
            [simple bit 'ackSup']
            [reserved uint 5 '0x00']
            [simple bit 'locked']
            [simple bit 'alarmunack']
            [simple bit 'inalarm']
        ]
        ['KnxPropertyDataType.PDT_BINARY_INFORMATION' BOOL
            [reserved uint 7        '0x00']
            [simple   bit           'value']
        ]
        ['KnxPropertyDataType.PDT_BITSET8' List
            [array    bit           'value' count '8']
        ]
        ['KnxPropertyDataType.PDT_BITSET16' List
            [array    bit           'value' count '16']
        ]
        ['KnxPropertyDataType.PDT_ENUM8' USINT
            [simple uint 8 'value']
        ]
        ['KnxPropertyDataType.PDT_SCALING' USINT
            [simple uint 8 'value']
        ]
        // Defaults to PDT_VARIABLE_LENGTH
        //['KnxPropertyDataType.PDT_NE_VL'
        //]
        // Defaults to PDT_VARIABLE_LENGTH
        //['KnxPropertyDataType.PDT_NE_FL'
        //]
        // Defaults to PDT_VARIABLE_LENGTH
        //['KnxPropertyDataType.PDT_FUNCTION'
        //]
        // Defaults to PDT_VARIABLE_LENGTH
        //['KnxPropertyDataType.PDT_ESCAPE'
        //]
        // 'KnxPropertyDataType.PDT_VARIABLE_LENGTH' == Catch all
        [ List [uint 8 'dataLength']
            [array uint 8 'value' count 'dataLength']
        ]
    ]
]