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

[discriminatedType OpenProtocolMessage
    [implicit      uint 32              length               'lengthInBytes'      encoding="AsciiUint"]
    [discriminator MID                  mid                                                           ]
    [simple        OpenProtocolRevision revision                                                      ]
    [simple        uint 8               noAckFlag                                                     ]
    [simple        uint 16              stationId                                 encoding="AsciiUint"]
    [simple        uint 16              spindleId                                 encoding="AsciiUint"]
    [simple        uint 16              sequenceNumber                            encoding="AsciiUint"]
    [simple        uint 8               numberOfMessageParts                      encoding="AsciiUint"]
    [simple        uint 8               messagePartNumber                         encoding="AsciiUint"]
    [typeSwitch mid
        ['ApplicationCommunicationStart' *ApplicationCommunicationStart
        ]
        ['ApplicationCommunicationStartAcknowledge' *ApplicationCommunicationStartAcknowledge(revision)
            [array ApplicationCommunicationStartAcknowledgeBlock blocks           count 'revision.numCommunicationStartAcknowledgeBlocks']
        ]
        ['ApplicationCommunicationStop' *ApplicationCommunicationStop
        ]
        ['ApplicationCommandError' *ApplicationCommandError
            [simple MID                  requestMid]
            [simple Error                error     ]
        ]
        ['ApplicationCommandAccepted' *ApplicationCommandAccepted
            [simple MID                  requestMid]
        ]
        ['ApplicationGenericDataRequest' *ApplicationGenericDataRequest
            [simple   MID                  requestMid                             ]
            [simple   OpenProtocolRevision wantedRevision                         ]
            [implicit uint 16              extraDataLength 'LEN(extraData)'       ]
            [array    byte                 extraData       count 'extraDataLength']
        ]
    ]
    [const         uint 8  end                  0x00                                     ]
]

[discriminatedType ApplicationCommunicationStartAcknowledgeBlock
    [discriminator uint 16 blockType encoding="AsciiUint"]
    [typeSwitch blockType
        // Revision 1
        ['01' *CellId
            [simple   uint 32    cellId                    encoding="AsciiUint"]
        ]
        ['02' *ChannelId
            [simple   uint 16    channelId                 encoding="AsciiUint"]
        ]
        ['03' *ControllerName
            [simple   string 200 controllerName            encoding="ASCII"    ]
        ]

        // Additional Blocks for Revision 2
        ['04' *SupplierCode
            [simple   uint 24    supplierCode              encoding="AsciiUint"]
        ]

        // Additional Blocks for Revision 3
        ['05' *OpenProtocolVersion
            [simple   string 152 openProtocolVersion       encoding="ASCII"    ]
        ]
        ['06' *ControllerSoftwareVersion
            [simple   string 152 controllerSoftwareVersion encoding="ASCII"    ]
        ]
        ['07' *ToolSoftwareVersion
            [simple   string 152 toolSoftwareVersion       encoding="ASCII"    ]
        ]

        // Additional Blocks for Revision 4
        ['08' *RbuType
            [simple   string 192 rbuType                   encoding="ASCII"    ]
        ]
        ['09' *ControllerSerialNumber
            [simple   string 80  controllerSerialNumber    encoding="ASCII"    ]
        ]

        // Additional Blocks for Revision 5
        ['10' *SystemType
            [simple   string 24  systemType                encoding="ASCII"    ]
        ]
        ['11' *SystemSubtype
            [simple   string 24  systemSubtype             encoding="ASCII"    ]
        ]

        // Additional Blocks for Revision 6
        ['12' *SequenceNumberSupport
            [reserved uint 7     0                                             ]
            [simple   bit        sequenceNumberSupport                         ]
        ]
        ['13' *LinkingHandlingSupport
            [reserved uint 7     0                                             ]
            [simple   bit        linkingHandlingSupport                        ]
        ]
        ['14' *StationId
            [simple   string 80  stationId                 encoding="ASCII"    ]
        ]
        ['15' *StationName
            [simple   string 200 stationName               encoding="ASCII"    ]
        ]
        ['16' *ClientId
            [simple   uint 8     clientId                  encoding="AsciiUint"]
        ]
    ]
]

// Depending on the revision of the device, a different number of blocks are supported.
[enum uint 24 OpenProtocolRevision(uint 8 numCommunicationStartAcknowledgeBlocks) encoding="AsciiUint"
    ['1' Revision1                (       3                                     )]
    ['2' Revision2                (       4                                     )]
    ['3' Revision3                (       7                                     )]
    ['4' Revision4                (       9                                     )]
    ['5' Revision5                (       11                                    )]
    ['6' Revision6                (       16                                    )]
]

[enum MidTypes
    [JobMessage                         ] //  600 -  699
    [ToolMessage                        ] //  700 -  799
    [VinMessage                         ] //  800 -  899
    [TighteningResultMessage            ] //  900 -  999
    [AlarmMessage                       ] // 1000 - 1099
    [TimeMessage                        ] // 1100 - 1199
    [MultiSpindleStatusMessage          ] // 1200 - 1299
    [MultiSpindleResultMessage          ] // 1300 - 1399
    [UserInterfaceMessage               ] // 1400 - 1499
    [JobMessageAdvanced                 ] // 1500 - 1599
    [MultipleIdentifiersMessage         ] // 1600 - 1699
    [IOInterfaceMessage                 ] // 1700 - 1799
    [PlcUserDataMessage                 ] // 1800 - 1899
    [SelectorMessage                    ] // 1900 - 1999
    [ToolLocationSystemMessage          ] // 2000 - 2099
    [ControllerMessage                  ] // 2100 - 2199
    [StatisticMessage                   ] // 2200 - 2299
    [AutomaticManualModeMessage         ] // 2300 - 2399
    [OpenProtocolCommandsDisabledMessage] // 2400 - 2499
    [ParameterSetMessage                ] // 2500 - 2599
    [NewGroupsMessage                   ] // 2600 - 9999
]

[enum uint 32 Mid   encoding="AsciiUint"
    ['0001' ApplicationCommunicationStart               ] // *
    ['0002' ApplicationCommunicationStartAcknowledge    ] // *
    ['0003' ApplicationCommunicationStop                ]
    ['0004' ApplicationCommandError                     ] // *
    ['0005' ApplicationCommandAccepted                  ] // *
    ['0006' ApplicationGenericDataRequest               ]
    ['0007' Reserved                                    ]
    ['0008' ApplicationGenericSubscription              ] // *
    ['0009' ApplicationGenericUnsubscribe               ] // *
    ['0010' ParameterSetIdUploadRequest                 ] // *
    ['0011' ParameterSetIdUploadReply                   ] // *
    ['0012' ParameterSetDataUploadRequest               ] // *
    ['0013' ParameterSetDataUploadReply                 ] // *
    ['0014' ParameterSetSelectedSubscribe               ]
    ['0015' ParameterSetSelected                        ]
    ['0016' ParameterSetSelectedAcknowledge             ]
    ['0017' ParameterSetSelectedUnsubscribe             ]
    ['0018' SelectParameterSet                          ] // *
    ['0019' SetParameterSetBatchSize                    ]
    ['0020' ResetParameterSetBatchCounter               ]
    ['0021' LockAtBatchDoneSubscribe                    ]
    ['0022' LockAtBatchDoneUpload                       ]
    ['0023' LockAtBatchDoneUploadAcknowledge            ]
    ['0024' LockAtBatchDoneUnsubscribe                  ]
    ['0025' ReservedForFord                             ]

    ['0030' JobIdUploadRequest                          ]
    ['0031' JobIdUploadReply                            ]
    ['0032' JobDataUploadRequest                        ]
    ['0033' JobDataUploadReply                          ]
    ['0034' JobInfoSubscribe                            ]
    ['0035' JobInfo                                     ]
    ['0036' JobInfoAcknowledge                          ]
    ['0037' JobInfoUnsubscribe                          ]
    ['0038' SelectJob                                   ]
    ['0039' JobRestart                                  ]
    ['0040' ToolDataUploadRequest                       ]
    ['0041' ToolDataUploadReply                         ]
    ['0042' DisableTool                                 ] // *
    ['0043' EnableTool                                  ] // *
    ['0044' DisconnectToolRequest                       ]
    ['0045' SetCalibrationValueRequest                  ]
    ['0046' SetPrimaryToolRequest                       ]
    ['0047' PairingHandling                             ]
    ['0048' PairingStatus                               ]
    ['0049' PairingStatusAcknowledge                    ]
    ['0050' VehicleIdNumberDownloadRequest              ]
    ['0051' VehicleIdNumberSubscribe                    ]
    ['0052' VehicleIdNumber                             ]
    ['0053' VehicleIdNumberAcknowledge                  ]
    ['0054' VehicleIdNumberUnsubscribe                  ]

    ['0060' LastTighteningResultDataSubscribe           ] // *
    ['0061' LastTighteningResultData                    ] // *
    ['0062' LastTighteningResultDataAcknowledge         ]
    ['0063' LastTighteningResultDataUnsubscribe         ]
    ['0064' OldTighteningResultUploadRequest            ]
    ['0065' OldTighteningResultUploadReply              ]

    ['0070' AlarmSubscribe                              ] // *
    ['0071' Alarm                                       ] // *
    ['0072' AlarmAcknowledge                            ] // *
    ['0073' AlarmUnsubscribe                            ] // *
    ['0074' AlarmAcknowledgedOnController               ]
    ['0075' AlarmAcknowledgedOnControllerAcknowledge    ]
    ['0076' AlarmStatus                                 ] // *
    ['0077' AlarmStatusAcknowledge                      ] // *
    ['0078' AcknowledgeAlarmRemotelyOnController        ]

    ['0080' ReadTimeUploadRequest                       ]
    ['0081' ReadTimeUploadReply                         ]
    ['0082' SetTime                                     ] // *

    ['0090' MultiSpindleStatusSubscribe                 ]
    ['0091' MultiSpindleStatus                          ]
    ['0092' MultiSpindleStatusAcknowledge               ]
    ['0093' MultiSpindleStatusUnsubscribe               ]

    ['0100' MultiSpindleResultSubscribe                 ]
    ['0101' MultiSpindleResult                          ]
    ['0102' MultiSpindleResultAcknowledge               ]
    ['0103' MultiSpindleResultUnsubscribe               ]

    ['0105' LastPowerMacsTighteningResultDataSubscribe  ]
    ['0106' LastPowerMacsTighteningResultStationData    ]
    ['0107' LastPowerMacsTighteningResultBoltData       ]
    ['0108' LastPowerMacsTighteningResultDataAcknowledge]
    ['0109' LastPowerMacsTighteningResultDataUnsubscribe]
    ['0110' DisplayUserTextOnCompact                    ]
    ['0111' DisplayUserTextOnGraph                      ]

    ['0113' FlashGreenLightOnTool                       ]

    ['0120' JobLineControlInfoSubscribe                 ]
    ['0121' JobLineControlStarted                       ]
    ['0122' JobLineControlAlert1                        ]
    ['0123' JobLineControlAlert2                        ]
    ['0124' JobLineControlDone                          ]
    ['0125' JobLineControlInfoAcknowledge               ]
    ['0126' JobLineControlUnsubscribe                   ]
    ['0127' AbortJob                                    ]
    ['0128' JobBatchIncrement                           ]
    ['0129' JobBatchDecrement                           ]
    ['0130' JobOff                                      ]
    ['0131' SetJobLineControlStart                      ]
    ['0132' SetJobLineControlAlert1                     ]
    ['0133' SetJobLineControlAlert2                     ]

    ['0140' ExecuteDynamicJobRequest                    ] // *

    ['0150' IdentifierDownloadRequest                   ] // *
    ['0151' MultipleIdentifiersWorkOrderSubscribe       ]
    ['0152' MultipleIdentifiersWorkOrder                ]
    ['0153' MultipleIdentifiersWorkOrderAcknowledge     ]
    ['0154' MultipleIdentifiersWorkOrderUnsubscribe     ]
    ['0155' BypassIdentifier                            ]
    ['0156' ResetLatestIdentifier                       ]
    ['0157' ResetAllIdentifiers                         ]

    ['0200' SetExternalControlledRelays                 ]

    ['0210' StatusExternalMonitoredInputsSubscribe      ]
    ['0211' StatusExternalMonitoredInputs               ]
    ['0212' StatusExternalMonitoredInputsAcknowledge    ]
    ['0213' StatusExternalMonitoredInputsUnsubscribe    ]
    ['0214' IoDeviceStatusRequest                       ]
    ['0215' IoDeviceStatusReply                         ]
    ['0216' RelayFunctionSubscribe                      ]
    ['0217' RelayFunction                               ]
    ['0218' RelayFunctionAcknowledge                    ]
    ['0219' RelayFunctionUnsubscribe                    ]
    ['0220' DigitalInputFunctionSubscribe               ]
    ['0221' DigitalInputFunction                        ]
    ['0222' DigitalInputFunctionAcknowledge             ]
    ['0223' DigitalInputFunctionUnsubscribe             ]
    ['0224' SetDigitalInputFunction                     ]
    ['0225' ResetDigitalInputFunction                   ]

    ['0240' UserDataDownload                            ]
    ['0241' UserDataSubscribe                           ]
    ['0242' UserData                                    ]
    ['0243' UserDataAcknowledge                         ]
    ['0244' UserDataUnsubscribe                         ]
    ['0245' UserDataDownloadWithOffset                  ]

    ['0250' SelectorSocketInfoSubscribe                 ]
    ['0251' SelectorSocketInfo                          ]
    ['0252' SelectorSocketInfoAcknowledge               ]
    ['0253' SelectorSocketInfoUnsubscribe               ]
    ['0254' SelectorControlGreenLights                  ]
    ['0255' SelectorControlRedLights                    ]

    ['0260' ToolTagIdRequest                            ]
    ['0261' ToolTagIdSubscribe                          ]
    ['0262' ToolTagId                                   ]
    ['0263' ToolTagIdAcknowledge                        ]
    ['0264' ToolTagIdUnsubscribe                        ]

    ['0270' ControllerRebootRequest                     ]

    ['0300' HistogramUploadRequest                      ]
    ['0301' HistogramUploadReply                        ]

    ['0400' AutomaticManualModeSubscribe                ]
    ['0401' AutomaticManualMode                         ]
    ['0402' AutomaticManualModeAcknowledge              ]
    ['0403' AutomaticManualModeUnsubscribe              ]

    ['0410' AutoDisableSettingsRequest                  ]
    ['0411' AutoDisableSettingsReply                    ]

    ['0420' OpenProtocolCommandsDisabledSubscribe       ]
    ['0421' OpenProtocolCommandsDisabled                ]
    ['0422' OpenProtocolCommandsDisabledAcknowledge     ]
    ['0423' OpenProtocolCommandsDisabledUnsubscribe     ]

    ['0500' MotorTuningResultDataSubscribe              ]
    ['0501' MotorTuningResultData                       ]
    ['0502' MotorTuningResultDataAcknowledge            ]
    ['0503' MotorTuningResultDataUnsubscribe            ]
    ['0504' MotorTuningRequest                          ]

    ['0700' TighteningDataDownloadStatusForRadioTools   ]

    ['0900' ResultTracesCurve                           ] // *
    ['0901' ResultTracesCurvePlotData                   ] // *

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

[enum uint 16 Error encoding="AsciiUint"
    ['00' NoError]
    ['01' InvalidData]
    ['02' ParameterSetIdNotPresent]
    ['03' ParameterSetCanNotBeSet]
    ['04' ParameterSetNotRunning]

    ['06' VinUploadSubscriptionAlreadyExists]
    ['07' VinUploadSubscriptionDoesNotExist]
    ['08' VinInputSourceNotGranted]
    ['09' LastTighteningResultSubscriptionAlreadyExists]
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