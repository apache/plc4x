/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/////////////////////////////////////////////////////////////////////////////////////////
//
//   PROFINET IO
//
// CM: Context Manager
//
/////////////////////////////////////////////////////////////////////////////////////////

// TODO: Check if it's really Little Endian
// 5.1.2
// 5.5.2.2
[discriminatedType PnIoCm_Packet(DceRpc_PacketType packetType)
    [typeSwitch packetType
        ['REQUEST' PnIoCm_Packet_Req
            [simple uint 32      argsMaximum                          ]
            [implicit uint 32    argsLength       'lengthInBytes - 20']
            [simple uint 32      arrayMaximumCount                    ]
            [simple uint 32      arrayOffset                          ]
            [implicit uint 32    arrayActualCount  'lengthInBytes - 20']
            [array  PnIoCm_Block blocks            length 'argsLength']
        ]
        ['PING' PnIoCm_Packet_Ping
        ]
        ['RESPONSE' PnIoCm_Packet_Res
            [simple uint 8       errorCode2                           ]
            [simple uint 8       errorCode1                           ]
            [simple uint 8       errorDecode                          ]
            [simple uint 8       errorCode                            ]
            [implicit uint 32    argsLength       'lengthInBytes - 1 - 1 - 1 - 1 - 4 - 4 - 4 - 4']
            [simple uint 32      arrayMaximumCount                    ]
            [simple uint 32      arrayOffset                          ]
            [implicit uint 32    arrayActualCount  'lengthInBytes - 1 - 1 - 1 - 1 - 4 - 4 - 4 - 4'    ]
            [array  PnIoCm_Block blocks            length 'argsLength']
        ]
        ['FAULT' PnIoCm_Packet_Fault
            [simple uint 32      status                               ]
        ]
        ['NO_CALL' PnIoCm_Packet_NoCall
        ]
        ['REJECT'   PnIoCm_Packet_Rej
            [simple uint 32      status                               ]
        ]
    ]
]

[type UserData(uint 32 recordDataLength) byteOrder='BIG_ENDIAN'
    [array              byte      data count         'recordDataLength'       ]
]

[type PnIo_CyclicServiceDataUnit(int 16 dataUnitLength)
    [array    byte   data       count 'dataUnitLength'                 ]
]

// Big Endian
[discriminatedType PnIoCm_Block byteOrder='BIG_ENDIAN'
    [discriminator PnIoCm_BlockType blockType                           ]
    [typeSwitch blockType
        ['IOD_WRITE_REQUEST_HEADER' IODWriteRequestHeader
            [implicit      uint 16          blockLength      'index < 0x8000 ? lengthInBytes - 4 - recordDataLength : lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   uint 16                         sequenceNumber                                         ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 32                         api                                                    ]
            [simple   uint 16                         slotNumber                                             ]
            [simple   uint 16                         subSlotNumber                                          ]
            [const    uint 16                         padField                  0x0000                       ]
            [simple   uint 16                         index                                                  ]
            [simple   uint 32                         recordDataLength                                       ]
            [padding  uint 8      pad '0x00'          'index < 0x8000 ? 64 - 6 - 2 - 16 - 4 - 2 - 2 - 2 - 2 - 4 : 64 - 6 - 2 - 16 - 4 - 2 - 2 - 2 - 2 - 4']
            [optional UserData('recordDataLength')     userData          'index < 0x8000'                    ]
        ]
        ['IOD_WRITE_RESPONSE_HEADER' IODWriteResponseHeader
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   uint 16                         sequenceNumber                                         ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 32                         api                                                    ]
            [simple   uint 16                         slotNumber                                             ]
            [simple   uint 16                         subSlotNumber                                          ]
            [const    uint 16                         padField                  0x0000                       ]
            [simple   uint 16                         index                                                  ]
            [simple   uint 32                         recordDataLength                                       ]
            [padding  uint 8      pad '0x00'          '64 - 6 - 2 - 16 - 4 - 2 - 2 - 2 - 2 - 4']
        ]
        ['PD_INTERFACE_ADJUST' PDInterfaceAdjust
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [const    uint 16                         padField                  0x0000                       ]
            [const    uint 16                         multipleInterfaceModeReserved2                  0x0000 ]
            [const    uint 15                         multipleInterfaceModeReserved1                  0x0000 ]
            [simple   MultipleInterfaceModeNameOfDevice multipleInterfaceModeNameOfDevice                    ]
        ]
        ['PD_PORT_DATA_CHECK' PDPortDataCheck
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [const    uint 16                         padField                  0x0000                       ]
            [simple   uint 16                         slotNumber                                             ]
            [simple   uint 16                         subSlotNumber                                          ]
            [simple   PnIoCm_Block                    checkPeers                                             ]
        ]
        ['CHECK_PEERS'  CheckPeers
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [const    uint 8                          noOfPeers                 0x01                         ]
            [simple   PascalString                    peerPortId                                             ]
            [simple   PascalString                    peerChassisId                                          ]
        ]
        ['AR_BLOCK_REQ' PnIoCm_Block_ArReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_ArType                   arType                                                 ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [simple   MacAddress                      cmInitiatorMacAddr                                     ]
            [simple   DceRpc_ObjectUuid               cmInitiatorObjectUuid                                  ]
            // Begin ARProperties
            [simple   bit                             pullModuleAlarmAllowed                                 ]
            [simple   bit                             nonLegacyStartupMode                                   ]
            [simple   bit                             combinedObjectContainerUsed                            ]
            [reserved uint 17                         '0x00000'                                              ]
            [simple   bit                             acknowledgeCompanionAr                                 ]
            [simple   PnIoCm_CompanionArType          companionArType                                        ]
            [simple   bit                             deviceAccess                                           ]
            [reserved uint 3                          '0x0'                                                  ]
            [simple   bit                             cmInitiator                                            ]
            [simple   bit                             supervisorTakeoverAllowed                              ]
            [simple   PnIoCm_State                    state                                                  ]
            // End ARProperties
            [simple   uint 16                         cmInitiatorActivityTimeoutFactor                       ]
            [simple   uint 16                         cmInitiatorUdpRtPort                                   ]
            [implicit uint 16                         stationNameLength     'STR_LEN(cmInitiatorStationName)']
            [simple   vstring 'stationNameLength * 8' cmInitiatorStationName                                 ]
        ]
        ['AR_BLOCK_RES' PnIoCm_Block_ArRes
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_ArType          arType                                                          ]
            [simple   Uuid                   arUuid                                                          ]
            [simple   uint 16                sessionKey                                                      ]
            [simple   MacAddress             cmResponderMacAddr                                              ]
            [simple   uint 16                responderUDPRTPort                                              ]
        ]
        ['IOD_CONTROL_REQ' PnIoCm_Control_Request
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [reserved uint 16                         '0x0000'                                               ]
        ]
        ['IOX_BLOCK_REQ'    PnIoCM_Block_Request
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [simple   uint 16                         controlBlockProperties                                 ]
        ]
        ['IOX_BLOCK_RES'    PnIoCM_Block_Response
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [simple   uint 16                         controlBlockProperties                                 ]
        ]
        ['IOD_CONTROL_RES' PnIoCm_Control_Response
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [reserved uint 16                         '0x0000'                                               ]
        ]
        ['IO_CR_BLOCK_REQ' PnIoCm_Block_IoCrReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple PnIoCm_IoCrType          ioCrType                                               ]
            [simple uint 16                  ioCrReference                                          ]
            [simple uint 16                  lt                                                     ]
            // Begin IOCRProperties
            [simple   bit                    fullSubFrameStructure                                  ]
            [simple   bit                    distributedSubFrameWatchDog                            ]
            [simple   bit                    fastForwardingMacAdr                                   ]
            [reserved uint 17                '0x0000'                                               ]
            [simple   bit                    mediaRedundancy                                        ]
            [reserved uint 7                 '0x00'                                                 ]
            [simple   PnIoCm_RtClass         rtClass                                                ]
            // End IOCRProperties
            [simple   uint 16                dataLength                                             ]
            [simple   uint 16                frameId                                                ]
            [simple   uint 16                sendClockFactor                                        ]
            [simple   uint 16                reductionRatio                                         ]
            [simple   uint 16                phase                                                  ]
            [simple   uint 16                sequence                                               ]
            [simple   uint 32                frameSendOffset                                        ]
            [simple   uint 16                watchDogFactor                                         ]
            [simple   uint 16                dataHoldFactor                                         ]
            [simple   uint 16                ioCrTagHeader                                          ]
            [simple   MacAddress             ioCrMulticastMacAdr                                    ]
            [implicit uint 16                numberOfApis        'COUNT(apis)'                      ]
            [array    PnIoCm_IoCrBlockReqApi apis                count         'numberOfApis'       ]
        ]
        ['IO_CR_BLOCK_RES' PnIoCm_Block_IoCrRes
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple PnIoCm_IoCrType          ioCrType                                               ]
            [simple uint 16                  ioCrReference                                          ]
            [simple   uint 16                frameId                                                ]
        ]
        ['ALARM_CR_BLOCK_REQ' PnIoCm_Block_AlarmCrReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_AlarmCrType     alarmType                                              ]
            [simple   uint 16                lt                                                     ]
            // Begin AlarmCrProperties
            [reserved uint 30                '0x00000000'                                           ]
            [simple   bit                    transport                                              ]
            [simple   bit                    priority                                               ]
            // End AlarmCrProperties
            [simple   uint 16                rtaTimeoutFactor                                       ]
            [simple   uint 16                rtaRetries                                             ]
            [simple   uint 16                localAlarmReference                                    ]
            [simple   uint 16                maxAlarmDataLength                                     ]
            [simple   uint 16                alarmCtrTagHeaderHigh                                  ]
            [simple   uint 16                alarmCtrTagHeaderLow                                   ]
        ]
        ['ALARM_CR_BLOCK_RES' PnIoCm_Block_AlarmCrRes
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_AlarmCrType     alarmType                                              ]
            [simple   uint 16                localAlarmReference                                    ]
            [simple   uint 16                maxAlarmDataLength                                     ]
        ]
        ['EXPECTED_SUBMODULE_BLOCK_REQ' PnIoCm_Block_ExpectedSubmoduleReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [implicit uint 16                numberOfApis         'COUNT(apis)'                     ]
            [array    PnIoCm_ExpectedSubmoduleBlockReqApi apis   count         'numberOfApis'       ]
        ]
        ['MODULE_DIFF_BLOCK' PnIoCm_Block_ModuleDiff
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [implicit uint 16                numberOfApis         'COUNT(apis)'                     ]
            [array    PnIoCm_ModuleDiffBlockApi apis              count         'numberOfApis'      ]
        ]
        ['AR_SERVER_BLOCK' PnIoCm_Block_ArServer
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PascalString16BitLength            stationName            ]
            [padding  uint 8      pad '0x00'          '20 - 6 - (stationName.stringLength)'              ]
        ]
    ]
]

[type PnIoCm_IoCrBlockReqApi byteOrder='BIG_ENDIAN'
    [const    uint 32             api              0x00000000             ]
    [implicit uint 16             numIoDataObjects 'COUNT(ioDataObjects)'   ]
    [array    PnIoCm_IoDataObject ioDataObjects    count 'numIoDataObjects' ]
    [implicit uint 16             numIoCss         'COUNT(ioCss)'           ]
    [array    PnIoCm_IoCs         ioCss            count 'numIoCss'         ]
]

[type PnIoCm_IoDataObject byteOrder='BIG_ENDIAN'
    [simple   uint 16 slotNumber             ]
    [simple   uint 16 subSlotNumber          ]
    [simple   uint 16 ioDataObjectFrameOffset]
]

[type PnIoCm_IoCs byteOrder='BIG_ENDIAN'
    [simple   uint 16 slotNumber   ]
    [simple   uint 16 subSlotNumber]
    [simple   uint 16 ioFrameOffset]
]

[type PnIoCm_DataUnitIoCs byteOrder='BIG_ENDIAN'
    [simple   bit               dataState]
    [simple   uint 2            instance ]
    [reserved uint 4            '0x00'   ]
    [simple   bit               extension]
]

[type PnIoCm_DataUnitDataObject(uint 16 dataObjectLength) byteOrder='BIG_ENDIAN'
    [array    byte              dataState   count  'dataObjectLength']
    [simple   PnIoCm_DataUnitIoCs iops   ]
]

[type PnIoCm_ExpectedSubmoduleBlockReqApi byteOrder='BIG_ENDIAN'
    [const    uint 32          api               0x00000000                       ]
    [simple   uint 16          slotNumber                                           ]
    [simple   uint 32          moduleIdentNumber                                    ]
    [simple   uint 16          moduleProperties                                     ]
    [implicit uint 16          numSubmodules     'COUNT(submodules)'                ]
    [array    PnIoCm_Submodule submodules        count               'numSubmodules']
]

[type PnIoCm_ModuleDiffBlockApi byteOrder='BIG_ENDIAN'
    [const    uint 32                          api        0x00000000                    ]
    [implicit uint 16                          numModules 'COUNT(modules)'                ]
    [array    PnIoCm_ModuleDiffBlockApi_Module modules    count               'numModules']
]

[type PnIoCm_ModuleDiffBlockApi_Module byteOrder='BIG_ENDIAN'
    [simple   uint 16                             slotNumber                                           ]
    [simple   uint 32                             moduleIdentNumber                                    ]
    [simple   PnIoCm_ModuleState                  moduleState                                          ]
    [implicit uint 16                             numSubmodules     'COUNT(submodules)'                ]
    [array    PnIoCm_ModuleDiffBlockApi_Submodule submodules        count               'numSubmodules']
]

[type PnIoCm_ModuleDiffBlockApi_Submodule byteOrder='BIG_ENDIAN'
    [simple uint 16          subslotNumber       ]
    [simple uint 32          submoduleIdentNumber]
    [simple bit              codingUsesIdentInfo ]
    [simple PnIoCm_IdentInfo identInfo           ]
    [simple PnIoCm_ArInfo    arInfo              ]
    [simple bit              diagInfoAvailable   ]
    [simple bit              maintenanceDemanded ]
    [simple bit              maintenanceRequired ]
    [simple bit              qualifiedInfo       ]
    [simple PnIoCm_AddInfo   addInfo             ]
]

[discriminatedType PnIoCm_Submodule byteOrder='BIG_ENDIAN'
    [simple        uint 16                slotNumber                    ]
    [simple        uint 32                submoduleIdentNumber          ]
    // Begin SubmoduleProperties
    [reserved      uint 10                '0x000'                       ]
    [simple        bit                    discardIoxs                   ]
    [simple        bit                    reduceOutputModuleDataLength  ]
    [simple        bit                    reduceInputModuleDataLength   ]
    [simple        bit                    sharedInput                   ]
    [discriminator PnIoCm_SubmoduleType   submoduleType                 ]
    // End SubmoduleProperties
    [typeSwitch submoduleType
        ['NO_INPUT_NO_OUTPUT_DATA' PnIoCm_Submodule_NoInputNoOutputData
            [const    uint 16             dataDescription       0x0001]
            [const    uint 16             submoduleDataLength   0x0000]
            [const    uint 8              lengthIoCs            0x01  ]
            [const    uint 8              lengthIoPs            0x01  ]
        ]
        ['INPUT_DATA' PnIoCm_Submodule_InputData
            [const    uint 16             inputDataDescription  0x0001]
            [simple   uint 16             inputSubmoduleDataLength      ]
            [simple   uint 8              inputLengthIoCs               ]
            [simple   uint 8              inputLengthIoPs               ]
        ]
        ['OUTPUT_DATA' PnIoCm_Submodule_OutputData
            [const    uint 16             inputDataDescription  0x0002]
            [simple   uint 16             inputSubmoduleDataLength      ]
            [simple   uint 8              inputLengthIoCs               ]
            [simple   uint 8              inputLengthIoPs               ]
        ]
        ['INPUT_AND_OUTPUT_DATA' PnIoCm_Submodule_InputAndOutputData
            [const    uint 16             inputDataDescription  0x0001]
            [simple   uint 16             inputSubmoduleDataLength      ]
            [simple   uint 8              inputLengthIoCs               ]
            [simple   uint 8              inputLengthIoPs               ]
            [const    uint 16             outputDataDescription 0x0002]
            [simple   uint 16             outputSubmoduleDataLength     ]
            [simple   uint 8              outputLengthIoCs              ]
            [simple   uint 8              outputLengthIoPs              ]
        ]
    ]
]

[enum uint 16 PnIoCm_BlockType
    ['0x0008' IOD_WRITE_REQUEST_HEADER    ]
    ['0x0101' AR_BLOCK_REQ                ]
    ['0x0102' IO_CR_BLOCK_REQ             ]
    ['0x0103' ALARM_CR_BLOCK_REQ          ]
    ['0x0104' EXPECTED_SUBMODULE_BLOCK_REQ]
    ['0x0110' IOD_CONTROL_REQ             ]
    ['0x0112' IOX_BLOCK_REQ               ]
    ['0x0200' PD_PORT_DATA_CHECK          ]
    ['0x020a' CHECK_PEERS                 ]
    ['0x0250' PD_INTERFACE_ADJUST         ]
    ['0x8008' IOD_WRITE_RESPONSE_HEADER    ]
    ['0x8101' AR_BLOCK_RES                ]
    ['0x8102' IO_CR_BLOCK_RES             ]
    ['0x8103' ALARM_CR_BLOCK_RES          ]
    ['0x8104' MODULE_DIFF_BLOCK           ]
    ['0x8106' AR_SERVER_BLOCK             ]
    ['0x8110' IOD_CONTROL_RES             ]
    ['0x8112' IOX_BLOCK_RES               ]
]

[enum uint 8 ProfinetDeviceState
    ['0x00'     IDLE]
    ['0x01'     STARTUP]
    ['0x02'     PREMED]
    ['0x03'     WAITAPPLRDY]
    ['0x04'     APPLRDY]
    ['0x05'     CYCLICDATA]
    ['0x06'     SET_IP]
    ['0xFF'     ABORT]
]

[enum uint 16 PnIoCm_ArType
    ['0x0001' IO_CONTROLLER]
]

[enum uint 2 PnIoCm_CompanionArType
    ['0x0' SINGLE_AR]
]

[enum uint 3 PnIoCm_State
    ['0x1' ACTIVE]
]

[enum uint 16 PnIoCm_IoCrType
    ['0x0001' INPUT_CR]
    ['0x0002' OUTPUT_CR]
]

[enum uint 4 PnIoCm_RtClass
    ['0x2' RT_CLASS_2]
]

[enum uint 16 PnIoCm_AlarmCrType
    ['0x0001' ALARM_CR]
]

[enum uint 16 PnIoCm_ModuleState
    ['0x0002' PROPER_MODULE]
]

[enum uint 2 PnIoCm_SubmoduleType
    ['0x0' NO_INPUT_NO_OUTPUT_DATA]
    ['0x1' INPUT_DATA]
    ['0x2' OUTPUT_DATA]
    ['0x3' INPUT_AND_OUTPUT_DATA]
]

[enum bit MultipleInterfaceModeNameOfDevice
    ['false' PORT_PROVIDED_BY_LLDP]
    ['true'  NAME_PROVIDED_BY_LLDP]
]

[enum uint 16 PnIoCm_DescriptionType
    ['0x0001' INPUT]
]

[enum uint 4 PnIoCm_IdentInfo
    ['0x0' OK]
]

[enum uint 4 PnIoCm_ArInfo
    ['0x0' OWN]
]

[enum uint 3 PnIoCm_AddInfo
    ['0x0' NONE]
]

[enum uint 4 IntegerEncoding
    ['0x0' BIG_ENDIAN]
    ['0x1' LITTLE_ENDIAN]
]

[enum uint 4 CharacterEncoding
    ['0x0' ASCII]
    ['0x1' EBCDIC]
]

[enum uint 8 FloatingPointEncoding
    ['0x00' IEEE]
    ['0x01' VAX ]
    ['0x02' CRAY]
    ['0x03' IBM ]
]


[dataIo DataItem(ProfinetDataType dataType, uint 16 numberOfValues)
    [typeSwitch dataType,numberOfValues
        ['BOOL','1'  BOOL
            [simple   bit     value                            ]
        ]
        ['BOOL'      List
            [array    bit     value count 'numberOfValues'     ]
        ]
        ['OCTETSTRING','1'  USINT
            [simple uint 8 value]
        ]
        ['OCTETSTRING' List
            [array    uint 8     value count 'numberOfValues' ]
        ]
        ['BYTE','1'  BYTE
            [simple uint 8 value]
        ]
        ['BYTE' List
            [array    bit     value count 'numberOfValues * 8' ]
        ]
        ['WORD'      WORD
            [simple   uint 16 value]
        ]
        ['DWORD'     DWORD
            [simple   uint 32 value]
        ]
        ['LWORD'     LWORD
            [simple   uint 64 value]
        ]
        ['SINT','1' SINT
            [simple   int 8   value ]
        ]
        ['SINT' List
            [array int 8 value count 'numberOfValues']
        ]
        ['INT','1' INT
            [simple int 16 value]
        ]
        ['INT' List
            [array int 16 value count 'numberOfValues']
        ]
        ['DINT','1' DINT
            [simple int 32 value]
        ]
        ['DINT' List
            [array int 32 value count 'numberOfValues']
        ]
        ['LINT','1' LINT
            [simple int 64 value]
        ]
        ['LINT' List
            [array int 64 value count 'numberOfValues']
        ]
        ['USINT','1' USINT
            [simple   uint 8 value ]
        ]
        ['USINT' List
            [array uint 8 value count 'numberOfValues']
        ]
        ['UINT','1' UINT
            [simple uint 16 value]
        ]
        ['UINT' List
            [array uint 16 value count 'numberOfValues']
        ]
        ['UDINT','1' UDINT
            [simple uint 32 value]
        ]
        ['UDINT' List
            [array uint 32 value count 'numberOfValues']
        ]
        ['ULINT','1' ULINT
            [simple uint 64 value]
        ]
        ['ULINT' List
            [array uint 64 value count 'numberOfValues']
        ]
        ['REAL','1' REAL
            [simple float 32  value]
        ]
        ['REAL' List
            [array float 32 value count 'numberOfValues']
        ]
        ['LREAL','1' LREAL
            [simple float 64  value]
        ]
        ['LREAL' List
            [array float 64 value count 'numberOfValues']
        ]
        ['CHAR','1' CHAR
            [simple string 8 value encoding='"UTF-8"']
        ]
        ['CHAR' List
            [array string 8 value count 'numberOfValues' encoding='"UTF-8"']
        ]
        ['WCHAR','1' WCHAR
            [simple string 16 value encoding='"UTF-16"']
        ]
        ['WCHAR' List
            [array string 16 value count 'numberOfValues' encoding='"UTF-16"']
        ]
        ['UNICODESTRING8','1' CHAR
            [simple string 8 value encoding='"UTF-8"']
        ]
        ['UNICODESTRING8' List
            [array string 8 value count 'numberOfValues' encoding='"UTF-8"']
        ]
        ['WSTRING','1' CHAR
            [simple string 16 value encoding='"UTF-16"']
        ]
        ['WSTRING' List
            [array string 16 value count 'numberOfValues' encoding='"UTF-16"']
        ]
        ['VISIBLESTRING','1' CHAR
            [simple string 8 value encoding='"UTF-8"']
        ]
        ['VISIBLESTRING' List
            [array string 8 value count 'numberOfValues' encoding='"UTF-8"']
        ]
        ['F_MESSAGETRAILER4BYTE','1'  List
            [array    uint 8     value count '4 * 8' ]
        ]
        ['F_MESSAGETRAILER4BYTE' List
            [array    uint 8     value count 'numberOfValues * 32' ]
        ]
    ]
]

[enum uint 8 ProfinetDataType(uint 8 dataTypeSize, string 16 conversion)
    ['1' BOOL ['1','BOOLEAN']]
    ['2' BYTE ['1','BYTE']]
    ['3' WORD ['2','WORD']]
    ['4' DWORD ['4','DWORD']]
    ['5' LWORD ['8','LWORD']]
    ['6' SINT ['1','SIGNED8']]
    ['7' INT ['2','SIGNED16']]
    ['8' DINT ['4','SIGNED32']]
    ['9' LINT ['8','SIGNED64']]
    ['10' USINT ['1','UNSIGNED8']]
    ['11' UINT ['2','UNSIGNED16']]
    ['12' UDINT ['4','UNSIGNED32']]
    ['13' ULINT ['8','UNSIGNED64']]
    ['14' REAL ['4','FLOAT32']]
    ['15' LREAL ['8','FLOAT64']]
    ['16' TIME ['8','TIME']]
    ['17' LTIME ['8','LTIME']]
    ['18' DATE ['8','DATE']]
    ['19' LDATE ['8','LDATE']]
    ['20' TIME_OF_DAY ['8','TIME_OF_DAY']]
    ['21' LTIME_OF_DAY ['8','LTIME_OF_DAY']]
    ['22' DATE_AND_TIME ['8','DATE_AND_TIME']]
    ['23' LDATE_AND_TIME ['8','LDATE_AND_TIME']]
    ['24' CHAR ['1','CHAR']]
    ['25' WCHAR ['2','WCHAR']]
    ['26' UNICODESTRING8 ['1','UNICODESTRING8']]
    ['27' WSTRING ['2','WSTRING']]
    ['28' VISIBLESTRING ['1','VISIBLESTRING']]
    ['29' F_MESSAGETRAILER4BYTE ['4','F_MESSAGETRAILER4BYTE']]
    ['30' TIMESTAMP ['12','TIMESTAMP']]
    ['31' TIMESTAMPDIFFERENCE ['12','TIMESTAMPDIFFERENCE']]
    ['32' TIMESTAMPDIFFERENCESHORT ['8','TIMESTAMPDIFFERENCESHORT']]
    ['33' OCTETSTRING ['1','OCTETSTRING']]
]
