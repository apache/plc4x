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

// Spec available from here:
// https://de.scribd.com/document/428086428/OpenProtocol-Specification-R-2-8-0-9836-4415-01

[type Constants
    [const          uint 16     tcpDefaultPort 4545]
]

[discriminatedType OpenProtocolMessage(uint 24 revision)
    [implicit      uint 32              length               'lengthInBytes - 1'                encoding='"ASCII"'                      ]
    [discriminator Mid                  mid                                                     encoding='"ASCII"'                      ]
    [optional      uint 24              midRevision                                                encoding='"ASCII"' nullBytesHex='202020']
    [optional      uint 8               noAckFlag                                               encoding='"ASCII"' nullBytesHex='20'    ]
    [optional      uint 16              targetStationId                                         encoding='"ASCII"' nullBytesHex='2020'  ]
    [optional      uint 16              targetSpindleId                                         encoding='"ASCII"' nullBytesHex='2020'  ]
    [optional      uint 16              sequenceNumber                                          encoding='"ASCII"' nullBytesHex='2020'  ]
    [optional      uint 8               numberOfMessageParts                                    encoding='"ASCII"' nullBytesHex='20'    ]
    [optional      uint 8               messagePartNumber                                       encoding='"ASCII"' nullBytesHex='20'    ]
    [typeSwitch mid
        ['ApplicationCommunicationStart' *ApplicationCommunicationStart
        ]
        ['ApplicationCommunicationStartAcknowledge' *ApplicationCommunicationStartAcknowledge
            [typeSwitch revision
                ['1' *Rev1
                    [const    uint   16  blockIdCellId                    1    encoding='"ASCII"']
                    [simple   uint   32  cellId                                encoding='"ASCII"']
                    [const    uint   16  blockIdChannelId                 2    encoding='"ASCII"']
                    [simple   uint   16  channelId                             encoding='"ASCII"']
                    [const    uint   16  blockIdControllerName            3    encoding='"ASCII"']
                    [simple   string 200 controllerName                        encoding='"ASCII"']
                ]
                ['2' *Rev2
                    [const    uint   16  blockIdCellId                    1    encoding='"ASCII"']
                    [simple   uint   32  cellId                                encoding='"ASCII"']
                    [const    uint   16  blockIdChannelId                 2    encoding='"ASCII"']
                    [simple   uint   16  channelId                             encoding='"ASCII"']
                    [const    uint   16  blockIdControllerName            3    encoding='"ASCII"']
                    [simple   string 200 controllerName                        encoding='"ASCII"']
                    [const    uint   16  blockIdSupplierCode              4    encoding='"ASCII"']
                    [simple   uint   24  supplierCode                          encoding='"ASCII"']
                ]
                ['3' *Rev3
                    [const    uint   16  blockIdCellId                    1    encoding='"ASCII"']
                    [simple   uint   32  cellId                                encoding='"ASCII"']
                    [const    uint   16  blockIdChannelId                 2    encoding='"ASCII"']
                    [simple   uint   16  channelId                             encoding='"ASCII"']
                    [const    uint   16  blockIdControllerName            3    encoding='"ASCII"']
                    [simple   string 200 controllerName                        encoding='"ASCII"']
                    [const    uint   16  blockIdSupplierCode              4    encoding='"ASCII"']
                    [simple   uint   24  supplierCode                          encoding='"ASCII"']
                    [const    uint   16  blockIdOpenProtocolVersion       5    encoding='"ASCII"']
                    [simple   string 152 openProtocolVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSoftwareVersion 6    encoding='"ASCII"']
                    [simple   string 152 controllerSoftwareVersion             encoding='"ASCII"']
                    [const    uint   16  blockIdToolSoftwareVersion       7    encoding='"ASCII"']
                    [simple   string 152 toolSoftwareVersion                   encoding='"ASCII"']
                ]
                ['4' *Rev4
                    [const    uint   16  blockIdCellId                    1    encoding='"ASCII"']
                    [simple   uint   32  cellId                                encoding='"ASCII"']
                    [const    uint   16  blockIdChannelId                 2    encoding='"ASCII"']
                    [simple   uint   16  channelId                             encoding='"ASCII"']
                    [const    uint   16  blockIdControllerName            3    encoding='"ASCII"']
                    [simple   string 200 controllerName                        encoding='"ASCII"']
                    [const    uint   16  blockIdSupplierCode              4    encoding='"ASCII"']
                    [simple   uint   24  supplierCode                          encoding='"ASCII"']
                    [const    uint   16  blockIdOpenProtocolVersion       5    encoding='"ASCII"']
                    [simple   string 152 openProtocolVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSoftwareVersion 6    encoding='"ASCII"']
                    [simple   string 152 controllerSoftwareVersion             encoding='"ASCII"']
                    [const    uint   16  blockIdToolSoftwareVersion       7    encoding='"ASCII"']
                    [simple   string 152 toolSoftwareVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdRbuType                   8    encoding='"ASCII"']
                    [simple   string 192 rbuType                               encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSerialNumber    9    encoding='"ASCII"']
                    [simple   string 80  controllerSerialNumber                encoding='"ASCII"']
                ]
                ['5' *Rev5
                    [const    uint   16  blockIdCellId                    1    encoding='"ASCII"']
                    [simple   uint   32  cellId                                encoding='"ASCII"']
                    [const    uint   16  blockIdChannelId                 2    encoding='"ASCII"']
                    [simple   uint   16  channelId                             encoding='"ASCII"']
                    [const    uint   16  blockIdControllerName            3    encoding='"ASCII"']
                    [simple   string 200 controllerName                        encoding='"ASCII"']
                    [const    uint   16  blockIdSupplierCode              4    encoding='"ASCII"']
                    [simple   uint   24  supplierCode                          encoding='"ASCII"']
                    [const    uint   16  blockIdOpenProtocolVersion       5    encoding='"ASCII"']
                    [simple   string 152 openProtocolVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSoftwareVersion 6    encoding='"ASCII"']
                    [simple   string 152 controllerSoftwareVersion             encoding='"ASCII"']
                    [const    uint   16  blockIdToolSoftwareVersion       7    encoding='"ASCII"']
                    [simple   string 152 toolSoftwareVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdRbuType                   8    encoding='"ASCII"']
                    [simple   string 192 rbuType                               encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSerialNumber    9    encoding='"ASCII"']
                    [simple   string 80  controllerSerialNumber                encoding='"ASCII"']
                    [const    uint   16  blockIdSystemType                10   encoding='"ASCII"']
                    [simple   string 24  systemType                            encoding='"ASCII"']
                    [const    uint   16  blockIdSystemSubtype             11   encoding='"ASCII"']
                    [simple   string 24  systemSubtype                         encoding='"ASCII"']
                ]
                ['6' *Rev6
                    [const    uint   16  blockIdCellId                    1    encoding='"ASCII"']
                    [simple   uint   32  cellId                                encoding='"ASCII"']
                    [const    uint   16  blockIdChannelId                 2    encoding='"ASCII"']
                    [simple   uint   16  channelId                             encoding='"ASCII"']
                    [const    uint   16  blockIdControllerName            3    encoding='"ASCII"']
                    [simple   string 200 controllerName                        encoding='"ASCII"']
                    [const    uint   16  blockIdSupplierCode              4    encoding='"ASCII"']
                    [simple   uint   24  supplierCode                          encoding='"ASCII"']
                    [const    uint   16  blockIdOpenProtocolVersion       5    encoding='"ASCII"']
                    [simple   string 152 openProtocolVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSoftwareVersion 6    encoding='"ASCII"']
                    [simple   string 152 controllerSoftwareVersion             encoding='"ASCII"']
                    [const    uint   16  blockIdToolSoftwareVersion       7    encoding='"ASCII"']
                    [simple   string 152 toolSoftwareVersion                   encoding='"ASCII"']
                    [const    uint   16  blockIdRbuType                   8    encoding='"ASCII"']
                    [simple   string 192 rbuType                               encoding='"ASCII"']
                    [const    uint   16  blockIdControllerSerialNumber    9    encoding='"ASCII"']
                    [simple   string 80  controllerSerialNumber                encoding='"ASCII"']
                    [const    uint   16  blockIdSystemType                10   encoding='"ASCII"']
                    [simple   string 24  systemType                            encoding='"ASCII"']
                    [const    uint   16  blockIdSystemSubtype             11   encoding='"ASCII"']
                    [simple   string 24  systemSubtype                         encoding='"ASCII"']
                    [const    uint   16  blockIdSequenceNumberSupport     12   encoding='"ASCII"']
                    [reserved uint   7   '0x00'                                                  ]
                    [simple   bit        sequenceNumberSupport                                   ]
                    [const    uint   16  blockIdLinkingHandlingSupport    13   encoding='"ASCII"']
                    [reserved uint   7   '0x00'                                                  ]
                    [simple   bit        linkingHandlingSupport                                  ]
                    [const    uint   16  blockIdStationId                 14   encoding='"ASCII"']
                    [simple   string 80  stationId                             encoding='"ASCII"']
                    [const    uint   16  blockIdStationName               15   encoding='"ASCII"']
                    [simple   string 200 stationName                           encoding='"ASCII"']
                    [const    uint   16  blockIdClientId                  16   encoding='"ASCII"']
                    [simple   uint   8   clientId                              encoding='"ASCII"']
                ]
            ]
        ]
        ['ApplicationCommunicationStop' *ApplicationCommunicationStop
        ]
        ['ApplicationCommandError'      *ApplicationCommandError
            [typeSwitch revision
                ['1' *Rev1
                    [simple Mid            requestMid                                           encoding='"ASCII"']
                    [simple Error          error                                                encoding='"ASCII"']
                ]
            ]
        ]
        ['ApplicationCommandAccepted'   *ApplicationCommandAccepted
            [typeSwitch revision
                ['1' *Rev1
                    [simple Mid            midNumberAccepted                                    encoding='"ASCII"']
                ]
            ]
        ]
        ['ApplicationGenericDataRequest' *ApplicationGenericDataRequest
            [typeSwitch revision
                ['1' *Rev1
                    [simple   Mid                  requestMid                                           encoding='"ASCII"']
                    [simple   uint 24              wantedRevision                                       encoding='"ASCII"']
                    [implicit uint 16              extraDataLength 'COUNT(extraData)'                   encoding='"ASCII"']
                    [array    byte                 extraData       count 'extraDataLength'                                ]
                ]
            ]
        ]
        ['ApplicationGenericSubscription' *ApplicationGenericSubscription
            [typeSwitch revision
                ['1' *Rev1
                    [simple   Mid                  subscriptionMid                                      encoding='"ASCII"']
                    [simple   uint 24              wantedRevision                                       encoding='"ASCII"']
                    [implicit uint 16              extraDataLength 'COUNT(extraData)'                   encoding='"ASCII"']
                    [array    byte                 extraData       count 'extraDataLength'                                ]
                ]
            ]
        ]
        ['ApplicationGenericUnsubscribe' *ApplicationGenericUnsubscribe
            [typeSwitch revision
                ['1' *Rev1
                    [simple   Mid                  subscriptionMid                                      encoding='"ASCII"']
                    [simple   uint 24              extraDataRevision                                    encoding='"ASCII"']
                    [implicit uint 16              extraDataLength 'COUNT(extraData)'                   encoding='"ASCII"']
                    [array    byte                 extraData       count 'extraDataLength'                                ]
                ]
            ]
        ]
    ]
    [const         uint 8  end                  0x00                                                              ]
]

[enum uint 32 Mid   encoding='"ASCII"'
    ['1' ApplicationCommunicationStart               ] // *
    ['2' ApplicationCommunicationStartAcknowledge    ] // *
    ['3' ApplicationCommunicationStop                ]
    ['4' ApplicationCommandError                     ] // *
    ['5' ApplicationCommandAccepted                  ] // *
    ['6' ApplicationGenericDataRequest               ]
    ['7' Reserved                                    ]
    ['8' ApplicationGenericSubscription              ] // *
    ['9' ApplicationGenericUnsubscribe               ] // *
    ['10' ParameterSetIdUploadRequest                 ] // *
    ['11' ParameterSetIdUploadReply                   ] // *
    ['12' ParameterSetDataUploadRequest               ] // *
    ['13' ParameterSetDataUploadReply                 ] // *
    ['14' ParameterSetSelectedSubscribe               ]
    ['15' ParameterSetSelected                        ]
    ['16' ParameterSetSelectedAcknowledge             ]
    ['17' ParameterSetSelectedUnsubscribe             ]
    ['18' SelectParameterSet                          ] // *
    ['19' SetParameterSetBatchSize                    ]
    ['20' ResetParameterSetBatchCounter               ]
    ['21' LockAtBatchDoneSubscribe                    ]
    ['22' LockAtBatchDoneUpload                       ]
    ['23' LockAtBatchDoneUploadAcknowledge            ]
    ['24' LockAtBatchDoneUnsubscribe                  ]
    ['25' ReservedForFord                             ]

    ['30' JobIdUploadRequest                          ]
    ['31' JobIdUploadReply                            ]
    ['32' JobDataUploadRequest                        ]
    ['33' JobDataUploadReply                          ]
    ['34' JobInfoSubscribe                            ]
    ['35' JobInfo                                     ]
    ['36' JobInfoAcknowledge                          ]
    ['37' JobInfoUnsubscribe                          ]
    ['38' SelectJob                                   ]
    ['39' JobRestart                                  ]
    ['40' ToolDataUploadRequest                       ]
    ['41' ToolDataUploadReply                         ]
    ['42' DisableTool                                 ] // *
    ['43' EnableTool                                  ] // *
    ['44' DisconnectToolRequest                       ]
    ['45' SetCalibrationValueRequest                  ]
    ['46' SetPrimaryToolRequest                       ]
    ['47' PairingHandling                             ]
    ['48' PairingStatus                               ]
    ['49' PairingStatusAcknowledge                    ]
    ['50' VehicleIdNumberDownloadRequest              ]
    ['51' VehicleIdNumberSubscribe                    ]
    ['52' VehicleIdNumber                             ]
    ['53' VehicleIdNumberAcknowledge                  ]
    ['54' VehicleIdNumberUnsubscribe                  ]

    ['60' LastTighteningResultDataSubscribe           ] // *
    ['61' LastTighteningResultData                    ] // *
    ['62' LastTighteningResultDataAcknowledge         ]
    ['63' LastTighteningResultDataUnsubscribe         ]
    ['64' OldTighteningResultUploadRequest            ]
    ['65' OldTighteningResultUploadReply              ]

    ['70' AlarmSubscribe                              ] // *
    ['71' Alarm                                       ] // *
    ['72' AlarmAcknowledge                            ] // *
    ['73' AlarmUnsubscribe                            ] // *
    ['74' AlarmAcknowledgedOnController               ]
    ['75' AlarmAcknowledgedOnControllerAcknowledge    ]
    ['76' AlarmStatus                                 ] // *
    ['77' AlarmStatusAcknowledge                      ] // *
    ['78' AcknowledgeAlarmRemotelyOnController        ]

    ['80' ReadTimeUploadRequest                       ]
    ['81' ReadTimeUploadReply                         ]
    ['82' SetTime                                     ] // *

    ['90' MultiSpindleStatusSubscribe                 ]
    ['91' MultiSpindleStatus                          ]
    ['92' MultiSpindleStatusAcknowledge               ]
    ['93' MultiSpindleStatusUnsubscribe               ]

    ['100' MultiSpindleResultSubscribe                 ]
    ['101' MultiSpindleResult                          ]
    ['102' MultiSpindleResultAcknowledge               ]
    ['103' MultiSpindleResultUnsubscribe               ]

    ['105' LastPowerMacsTighteningResultDataSubscribe  ]
    ['106' LastPowerMacsTighteningResultStationData    ]
    ['107' LastPowerMacsTighteningResultBoltData       ]
    ['108' LastPowerMacsTighteningResultDataAcknowledge]
    ['109' LastPowerMacsTighteningResultDataUnsubscribe]
    ['110' DisplayUserTextOnCompact                    ]
    ['111' DisplayUserTextOnGraph                      ]

    ['113' FlashGreenLightOnTool                       ]

    ['120' JobLineControlInfoSubscribe                 ]
    ['121' JobLineControlStarted                       ]
    ['122' JobLineControlAlert1                        ]
    ['123' JobLineControlAlert2                        ]
    ['124' JobLineControlDone                          ]
    ['125' JobLineControlInfoAcknowledge               ]
    ['126' JobLineControlUnsubscribe                   ]
    ['127' AbortJob                                    ]
    ['128' JobBatchIncrement                           ]
    ['129' JobBatchDecrement                           ]
    ['130' JobOff                                      ]
    ['131' SetJobLineControlStart                      ]
    ['132' SetJobLineControlAlert1                     ]
    ['133' SetJobLineControlAlert2                     ]

    ['140' ExecuteDynamicJobRequest                    ] // *

    ['150' IdentifierDownloadRequest                   ] // *
    ['151' MultipleIdentifiersWorkOrderSubscribe       ]
    ['152' MultipleIdentifiersWorkOrder                ]
    ['153' MultipleIdentifiersWorkOrderAcknowledge     ]
    ['154' MultipleIdentifiersWorkOrderUnsubscribe     ]
    ['155' BypassIdentifier                            ]
    ['156' ResetLatestIdentifier                       ]
    ['157' ResetAllIdentifiers                         ]

    ['200' SetExternalControlledRelays                 ]

    ['210' StatusExternalMonitoredInputsSubscribe      ]
    ['211' StatusExternalMonitoredInputs               ]
    ['212' StatusExternalMonitoredInputsAcknowledge    ]
    ['213' StatusExternalMonitoredInputsUnsubscribe    ]
    ['214' IoDeviceStatusRequest                       ]
    ['215' IoDeviceStatusReply                         ]
    ['216' RelayFunctionSubscribe                      ]
    ['217' RelayFunction                               ]
    ['218' RelayFunctionAcknowledge                    ]
    ['219' RelayFunctionUnsubscribe                    ]
    ['220' DigitalInputFunctionSubscribe               ]
    ['221' DigitalInputFunction                        ]
    ['222' DigitalInputFunctionAcknowledge             ]
    ['223' DigitalInputFunctionUnsubscribe             ]
    ['224' SetDigitalInputFunction                     ]
    ['225' ResetDigitalInputFunction                   ]

    ['240' UserDataDownload                            ]
    ['241' UserDataSubscribe                           ]
    ['242' UserData                                    ]
    ['243' UserDataAcknowledge                         ]
    ['244' UserDataUnsubscribe                         ]
    ['245' UserDataDownloadWithOffset                  ]

    ['250' SelectorSocketInfoSubscribe                 ]
    ['251' SelectorSocketInfo                          ]
    ['252' SelectorSocketInfoAcknowledge               ]
    ['253' SelectorSocketInfoUnsubscribe               ]
    ['254' SelectorControlGreenLights                  ]
    ['255' SelectorControlRedLights                    ]

    ['260' ToolTagIdRequest                            ]
    ['261' ToolTagIdSubscribe                          ]
    ['262' ToolTagId                                   ]
    ['263' ToolTagIdAcknowledge                        ]
    ['264' ToolTagIdUnsubscribe                        ]

    ['270' ControllerRebootRequest                     ]

    ['300' HistogramUploadRequest                      ]
    ['301' HistogramUploadReply                        ]

    ['400' AutomaticManualModeSubscribe                ]
    ['401' AutomaticManualMode                         ]
    ['402' AutomaticManualModeAcknowledge              ]
    ['403' AutomaticManualModeUnsubscribe              ]

    ['410' AutoDisableSettingsRequest                  ]
    ['411' AutoDisableSettingsReply                    ]

    ['420' OpenProtocolCommandsDisabledSubscribe       ]
    ['421' OpenProtocolCommandsDisabled                ]
    ['422' OpenProtocolCommandsDisabledAcknowledge     ]
    ['423' OpenProtocolCommandsDisabledUnsubscribe     ]

    ['500' MotorTuningResultDataSubscribe              ]
    ['501' MotorTuningResultData                       ]
    ['502' MotorTuningResultDataAcknowledge            ]
    ['503' MotorTuningResultDataUnsubscribe            ]
    ['504' MotorTuningRequest                          ]

    ['700' TighteningDataDownloadStatusForRadioTools   ]

    ['900' ResultTracesCurve                           ] // *
    ['901' ResultTracesCurvePlotData                   ] // *

    ['1201' LastOperationResultOverallData              ]
    ['1202' LastOperationResultObjectData               ]
    ['1203' LastOperationResultDataAcknowledge          ]

    ['2100' DeviceCommand                               ]

    ['2500' ProgramDataDownload                         ]
    ['2501' ProgramDataUploadReply                      ]
    ['2502' PasswordRequest                             ]
    ['2503' PasswordReply                               ] // PasswordResponse in the documentation
    ['2504' ProgramPsetSelectionInDynamicJob            ]
    ['2505' DynamicPsetSelection                        ]

    ['2600' ModeIdUploadRequest                         ]
    ['2601' ModeIdUploadReply                           ]
    ['2602' ModeDataUploadRequest                       ]
    ['2603' ModeDataUploadReply                         ]
    ['2604' ModeSelected                                ]
    ['2605' ModeSelectedAcknowledge                     ]
    ['2606' SelectMode                                  ]

    ['8000' AudiEmergencyStatusSubscribe                ]
    ['8001' AudiEmergencyStatus                         ]
    ['8002' AudiEmergencyStatusAcknowledge              ]
    ['8003' AudiEmergencyStatusUnsubscribe              ]

    ['9997' LinkLevelPositiveAcknowledge                ]
    ['9998' LinkLevelNegativeAcknowledge                ]
    ['9999' KeepAliveOpenProtocolCommunication          ] // *
]

[enum uint 16 Error encoding='"ASCII"'
    ['0'  NoError]
    ['1'  InvalidData]
    ['2'  ParameterSetIdNotPresent]
    ['3'  ParameterSetCanNotBeSet]
    ['4'  ParameterSetNotRunning]

    ['6'  VinUploadSubscriptionAlreadyExists]
    ['7'  VinUploadSubscriptionDoesNotExist]
    ['8'  VinInputSourceNotGranted]
    ['9'  LastTighteningResultSubscriptionAlreadyExists]
    ['10' LastTighteningResultSubscriptionDoesNowExist]
    ['11' AlarmSubscriptionAlreadyExists]
    ['12' AlarmSubscriptionDoesNotExist]
    ['13' ParameterSetSubscriptionAlreadyExists]
    ['14' ParameterSetSubscriptionDoesNotExist]
    ['15' TighteningIdRequestedNotFound]
    ['16' ConnectionRejectedProtocolBusy]
    ['17' JobIdNotPresent]
    ['18' JobInfoSubscriptionAlreadyExists]
    ['19' JobInfoSubscriptionDoesNotExist]
    ['20' JobCanNotBeSet]
    ['21' JobNotRunning]
    ['22' NotPossibleToExecuteDynamicJobRequest]
    ['23' JobBatchDecrementFailed]
    ['24' NotPossibleToCreatePset]
    ['25' ProgrammingControlNotGranted]
    ['26' WrongToolTypeToPsetDownloadConnected]
    ['27' ToolIsInaccessible]
    ['28' JobAbortionIsInProgress]

    ['30' ControllerIdNotASyncMasterOrStationController]
    ['31' MultiSpindleStatusSubscriptionAlreadyExists]
    ['32' MultiSpindleStatusSubscriptionDoesNotExist]
    ['33' MultiSpindleResultSubscriptionAlreadyExists]
    ['34' MultiSpindleResultSubscriptionDoesNotExist]
    ['35' OtherMasterClientAlreadyConnected]

    ['40' JobLineControlInfoSubscriptionAlreadyExists]
    ['41' JobLineControlInfoSubscriptionDoesNotExist]
    ['42' IdentifierInputSourceNotGranted]
    ['43' MultipleIdentifiersWorkOrderSubscriptionAlreadyExists]
    ['44' MultipleIdentifiersWorkOrderSubscriptionDoesNotExist]

    ['50' StatusExternalMonitoredInputsSubscriptionAlreadyExists]
    ['51' StatusExternalMonitoredInputsSubscriptionDoesNotExist]
    ['52' IoDeviceNotConnected]
    ['53' FaultyIoDeviceId]
    ['54' ToolTagIdUnknown]
    ['55' ToolTagIdSubscriptionAlreadyExists]
    ['56' ToolTagIdSubscriptionDoesNotExist]
    ['57' ToolMotorTuningFailed]
    ['58' NoAlarmPresent]
    ['59' ToolCurrentlyInUse]
    ['60' NoHistogramAvailable]
    ['61' PairingFailed]
    ['62' PairingDenied]
    ['63' PairingOrPairingAbortionAttemptOrWrongToolType]
    ['64' PairingAbortionDenied]
    ['65' PairingAbortionFailed]
    ['66' PairingDisconnectionFailed]
    ['67' PairingInProgressOrAlreadyDone]
    ['68' PairingDeniedNoProgramControl]
    ['69' UnsupportedExtraDataRevision]
    ['70' CalibrationFailed]
    ['71' SubscriptionAlreadyExists]
    ['72' SubscriptionDoesNotExist]
    ['73' SubscribedMidUnsupported]
    ['74' SubscribedMidRevisionUnsupported]
    ['75' RequestedMidUnsupported]
    ['76' RequestedMidRevisionUnsupported]
    ['77' RequestedOnSpecificDataNotSupported]
    ['78' SubscriptionOnSpecificDataNotSupported]
    ['79' CommandFailed]
    ['80' AudiEmergencyStatusSubscriptionAlreadyExists] // AudiEmergencyStatusSubscriptionExists
    ['81' AudiEmergencyStatusSubscriptionDoesNotExist]
    ['82' AutomaticOrManualModeSubscriptionAlreadyExists] // AutomaticOrManualModeSubscribeAlreadyExists
    ['83' AutomaticOrManualModeSubscriptionDoesNotExist] // AutomaticOrManualModeSubscribeDoesNotExist
    ['84' TheRelayFunctionSubscriptionAlreadyExists]
    ['85' TheRelayFunctionSubscriptionDoesNotExist]
    ['86' TheSelectorSocketInfoSubscriptionAlreadyExists]
    ['87' TheSelectorSocketInfoSubscriptionDoesNotExist]
    ['88' TheDiginInfoSubscriptionAlreadyExists]
    ['89' TheDiginInfoSubscriptionDoesNotExist]
    ['90' LockAtBatchDoneSubscriptionAlreadyExists]
    ['91' LockAtBatchDoneSubscriptionDoesNotExist]
    ['92' OpenProtocolCommendsDisabled]
    ['93' OpenProtocolCommendsDisabledSubscriptionAlreadyExists]
    ['94' OpenProtocolCommendsDisabledSubscriptionDoesNotExist]
    ['95' RejectRequestPowerMacsIsInManualMode]
    ['96' RejectConnectionClientAlreadyConnected]
    ['97' MidRevisionUnsupported]
    ['98' ControllerInternalRequestTimeout]
    ['99' UnknownMid]
]