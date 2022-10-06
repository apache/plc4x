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

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAction
    ['0'    DIRECT  ]
    ['1'    REVERSE ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAccessPassbackMode
    ['0'    PASSBACK_OFF                            ]
    ['1'    HARD_PASSBACK                           ]
    ['2'    SOFT_PASSBACK                           ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAccessRuleLocationSpecifier
    ['0'    SPECIFIED                               ]
    ['1'    ALL                                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAccessRuleTimeRangeSpecifier
    ['0'    SPECIFIED                               ]
    ['1'    ALWAYS                                  ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAccumulatorRecordAccumulatorStatus
    ['0'    NORMAL                                  ]
    ['1'    STARTING                                ]
    ['2'    RECOVERED                               ]
    ['3'    ABNORMAL                                ]
    ['4'    FAILED                                  ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAuthenticationFactorType
    ['0'    UNDEFINED                               ]
    ['1'    ERROR                                   ]
    ['2'    CUSTOM                                  ]
    ['3'    SIMPLE_NUMBER16                         ]
    ['4'    SIMPLE_NUMBER32                         ]
    ['5'    SIMPLE_NUMBER56                         ]
    ['6'    SIMPLE_ALPHA_NUMERIC                    ]
    ['7'    ABA_TRACK2                              ]
    ['8'    WIEGAND26                               ]
    ['9'    WIEGAND37                               ]
    ['10'   WIEGAND37_FACILITY                      ]
    ['11'   FACILITY16_CARD32                       ]
    ['12'   FACILITY32_CARD32                       ]
    ['13'   FASC_N                                  ]
    ['14'   FASC_N_BCD                              ]
    ['15'   FASC_N_LARGE                            ]
    ['16'   FASC_N_LARGE_BCD                        ]
    ['17'   GSA75                                   ]
    ['18'   CHUID                                   ]
    ['19'   CHUID_FULL                              ]
    ['20'   GUID                                    ]
    ['21'   CBEFF_A                                 ]
    ['22'   CBEFF_B                                 ]
    ['23'   CBEFF_C                                 ]
    ['24'   USER_PASSWORD                           ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetAuthenticationStatus
    ['0'    NOT_READY                               ]
    ['1'    READY                                   ]
    ['2'    DISABLED                                ]
    ['3'    WAITING_FOR_AUTHENTICATION_FACTOR       ]
    ['4'    WAITING_FOR_ACCOMPANIMENT               ]
    ['5'    WAITING_FOR_VERIFICATION                ]
    ['6'    IN_PROGRESS                             ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetBackupState
    ['0'    IDLE                                    ]
    ['1'    PREPARING_FOR_BACKUP                    ]
    ['2'    PREPARING_FOR_RESTORE                   ]
    ['3'    PERFORMING_A_BACKUP                     ]
    ['4'    PERFORMING_A_RESTORE                    ]
    ['5'    BACKUP_FAILURE                          ]
    ['6'    RESTORE_FAILURE                         ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetBinaryPV
    ['0'    INACTIVE    ]
    ['1'    ACTIVE      ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetConfirmedServiceChoice
    ////
    // Alarm and Event Services

    ['0x00' ACKNOWLEDGE_ALARM                       ]
    ['0x01' CONFIRMED_COV_NOTIFICATION              ]
    ['0x1F' CONFIRMED_COV_NOTIFICATION_MULTIPLE     ]
    ['0x02' CONFIRMED_EVENT_NOTIFICATION            ]
    ['0x03' GET_ALARM_SUMMARY                       ]
    ['0x04' GET_ENROLLMENT_SUMMARY                  ]
    ['0x1D' GET_EVENT_INFORMATION                   ]
    ['0x1B' LIFE_SAFETY_OPERATION                   ]
    ['0x05' SUBSCRIBE_COV                           ]
    ['0x1C' SUBSCRIBE_COV_PROPERTY                  ]
    ['0x1E' SUBSCRIBE_COV_PROPERTY_MULTIPLE         ]
    //
    ////

    ////
    // File Access Services

    ['0x06' ATOMIC_READ_FILE                        ]
    ['0x07' ATOMIC_WRITE_FILE                       ]
    //
    ////

    ////
    // Object Access Services

    ['0x08' ADD_LIST_ELEMENT                        ]
    ['0x09' REMOVE_LIST_ELEMENT                     ]
    ['0x0A' CREATE_OBJECT                           ]
    ['0x0B' DELETE_OBJECT                           ]
    ['0x0C' READ_PROPERTY                           ]
    ['0x0E' READ_PROPERTY_MULTIPLE                  ]
    ['0x1A' READ_RANGE                              ]
    ['0x0F' WRITE_PROPERTY                          ]
    ['0x10' WRITE_PROPERTY_MULTIPLE                 ]
    //
    ////

    ////
    // Remote Device Management Services

    ['0x11' DEVICE_COMMUNICATION_CONTROL            ]
    ['0x12' CONFIRMED_PRIVATE_TRANSFER              ]
    ['0x13' CONFIRMED_TEXT_MESSAGE                  ]
    ['0x14' REINITIALIZE_DEVICE                     ]
    //
    ////

    ////
    //  Virtual Terminal Services

    ['0x15' VT_OPEN                                 ]
    ['0x16' VT_CLOSE                                ]
    ['0x17' VT_DATA                                 ]
    //
    ////

    ////
    //  Removed Services

    ['0x18' AUTHENTICATE                            ]
    ['0x19' REQUEST_KEY                             ]
    ['0x0D' READ_PROPERTY_CONDITIONAL               ]
    //
    ////
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriority
    ['0' NORMAL                                 ]
    ['1' URGENT                                 ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable
    ['0' ENABLE                                 ]
    ['1' DISABLE                                ]
    ['2' DISABLE_INITIATION                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter
    ['0' ALL                                    ]
    ['1' ACKED                                  ]
    ['2' NOT_ACKED                              ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilter
    ['0' OFFNORMAL                              ]
    ['1' FAULT                                  ]
    ['2' NORMAL                                 ]
    ['3' ALL                                    ]
    ['4' ACTIVE                                 ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice
    ['0x0' COLDSTART                            ]
    ['0x1' WARMSTART                            ]
    ['0x2' ACTIVATE_CHANGES                     ]
    ['0x3' STARTBACKUP                          ]
    ['0x4' ENDBACKUP                            ]
    ['0x5' STARTRESTORE                         ]
    ['0x6' ENDRESTORE                           ]
    ['0x7' ABORTRESTORE                         ]

    // This state should never occur as this is fixed. The generic approach however demands a fallback enum
    ['0xFF' VENDOR_PROPRIETARY_VALUE            ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetDoorSecuredStatus
    ['0'  SECURED                                   ]
    ['1'  UNSECURED                                 ]
    ['2'  UNKNOWN                                   ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetDoorValue
    ['0'  LOCK                                      ]
    ['1'  UNLOCK                                    ]
    ['2'  PULSE_UNLOCK                              ]
    ['3'  EXTENDED_PULSE_UNLOCK                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetFaultType
    ['0'  NONE                                      ]
    ['1'  FAULT_CHARACTERSTRING                     ]
    ['2'  FAULT_EXTENDED                            ]
    ['3'  FAULT_LIFE_SAFETY                         ]
    ['4'  FAULT_STATE                               ]
    ['5'  FAULT_STATUS_FLAGS                        ]
    ['6'  FAULT_OUT_OF_RANGE                        ]
    ['7'  FAULT_LISTED                              ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetFileAccessMethod
    ['0'  RECORD_ACCESS                             ]
    ['1'  STREAM_ACCESS                             ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetIPMode
    ['0'  NORMAL                                    ]
    ['1'  FOREIGN                                   ]
    ['2'  BBMD                                      ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetLiftCarDoorCommand
    ['0'  NONE                                      ]
    ['1'  OPEN                                      ]
    ['2'  CLOSE                                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetLiftGroupMode
    ['0'  UNKNOWN                                   ]
    ['1'  NORMAL                                    ]
    ['2'  DOWN_PEAK                                 ]
    ['3'  TWO_WAY                                   ]
    ['4'  FOUR_WAY                                  ]
    ['5'  EMERGENCY_POWER                           ]
    ['6'  UP_PEAK                                   ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetLightingInProgress
    ['0'  IDLE                                      ]
    ['1'  FADE_ACTIVE                               ]
    ['2'  RAMP_ACTIVE                               ]
    ['3'  NOT_CONTROLLED                            ]
    ['4'  OTHER                                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetLockStatus
    ['0'  LOCKED                                    ]
    ['1'  UNLOCKED                                  ]
    ['2'  LOCK_FAULT                                ]
    ['3'  UNUSED                                    ]
    ['4'  UNKNOWN                                   ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetNetworkNumberQuality
    ['0'  UNKNOWN                               ]
    ['1'  LEARNED                               ]
    ['2'  LEARNED_CONFIGURED                    ]
    ['3'  CONFIGURED                            ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetNodeType
    ['0x00' UNKNOWN                             ]
    ['0x01' SYSTEM                              ]
    ['0x02' NETWORK                             ]
    ['0x03' DEVICE                              ]
    ['0x04' ORGANIZATIONAL                      ]
    ['0x05' AREA                                ]
    ['0x06' EQUIPMENT                           ]
    ['0x07' POINT                               ]
    ['0x08' COLLECTION                          ]
    ['0x09' PROPERTY                            ]
    ['0x0A' FUNCTIONAL                          ]
    ['0x0B' OTHER                               ]
    ['0x0C' SUBSYSTEM                           ]
    ['0x0D' BUILDING                            ]
    ['0x0E' FLOOR                               ]
    ['0x0F' SECTION                             ]
    ['0x10' MODULE                              ]
    ['0x11' TREE                                ]
    ['0x12' MEMBER                              ]
    ['0x13' PROTOCOL                            ]
    ['0x14' ROOM                                ]
    ['0x15' ZONE                                ]
]

[enum uint 8 BACnetNotifyType
    ['0x0' ALARM                                ]
    ['0x1' EVENT                                ]
    ['0x2' ACK_NOTIFICATION                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetPolarity
    ['0'  NORMAL                                ]
    ['1'  REVERSE                               ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetProgramRequest
    ['0'  READY                                 ]
    ['1'  LOAD                                  ]
    ['2'  RUN                                   ]
    ['3'  HALT                                  ]
    ['4'  RESTART                               ]
    ['5'  UNLOAD                                ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetProgramState
    ['0'  IDLE                                  ]
    ['1'  LOADING                               ]
    ['2'  RUNNING                               ]
    ['3'  WAITING                               ]
    ['4'  HALTED                                ]
    ['5'  UNLOADING                             ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetProtocolLevel
    ['0'  PHYSICAL                              ]
    ['1'  PROTOCOL                              ]
    ['2'  BACNET_APPLICATION                    ]
    ['3'  NON_BACNET_APPLICATION                ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetRouterEntryStatus
    ['0'  AVAILABLE                              ]
    ['1'  BUSY                                   ]
    ['2'  DISCONNECTED                           ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetSecurityLevel
    ['0'  INCAPABLE                             ] // indicates that the device is configured to not use security
    ['1'  PLAIN                                 ]
    ['2'  SIGNED                                ]
    ['3'  ENCRYPTED                             ]
    ['4'  SIGNED_END_TO_END                     ]
    ['5'  ENCRYPTED_END_TO_END                  ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetSecurityPolicy
    ['0'  PLAIN_NON_TRUSTED                     ]
    ['1'  PLAIN_TRUSTED                         ]
    ['2'  SIGNED_TRUSTED                        ]
    ['3'  ENCRYPTED_TRUSTED                     ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetSegmentation
    ['0' SEGMENTED_BOTH                         ]
    ['1' SEGMENTED_TRANSMIT                     ]
    ['2' SEGMENTED_RECEIVE                      ]
    ['3' NO_SEGMENTATION                        ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetShedState
    ['0'  SHED_INACTIVE                         ]
    ['1'  SHED_REQUEST_PENDING                  ]
    ['2'  SHED_COMPLIANT                        ]
    ['3'  SHED_NON_COMPLIANT                    ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetTimerState
    ['0'  IDLE                                  ]
    ['1'  RUNNING                               ]
    ['2'  EXPIRED                               ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetTimerTransition
    ['0'  NONE                                  ]
    ['1'  IDLE_TO_RUNNING                       ]
    ['2'  RUNNING_TO_IDLE                       ]
    ['3'  RUNNING_TO_RUNNING                    ]
    ['4'  RUNNING_TO_EXPIRED                    ]
    ['5'  FORCED_TO_EXPIRED                     ]
    ['6'  EXPIRED_TO_IDLE                       ]
    ['7'  EXPIRED_TO_RUNNING                    ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetUnconfirmedServiceChoice
    ['0x00' I_AM                                        ]
    ['0x01' I_HAVE                                      ]
    ['0x02' UNCONFIRMED_COV_NOTIFICATION                ]
    ['0x03' UNCONFIRMED_EVENT_NOTIFICATION              ]
    ['0x04' UNCONFIRMED_PRIVATE_TRANSFER                ]
    ['0x05' UNCONFIRMED_TEXT_MESSAGE                    ]
    ['0x06' TIME_SYNCHRONIZATION                        ]
    ['0x07' WHO_HAS                                     ]
    ['0x08' WHO_IS                                      ]
    ['0x09' UTC_TIME_SYNCHRONIZATION                    ]
    ['0x0A' WRITE_GROUP                                 ]
    ['0x0B' UNCONFIRMED_COV_NOTIFICATION_MULTIPLE       ]
]

// Attention: No proprietary extension allowed
[enum uint 8 BACnetWriteStatus
    ['0'  IDLE                                          ]
    ['1'  IN_PROGRESS                                   ]
    ['2'  SUCCESSFUL                                    ]
    ['3'  FAILED                                        ]
]

// Attention: No proprietary extension allowed
[enum uint 16 BVLCResultCode
    ['0x0000' SUCCESSFUL_COMPLETION                             ]
    ['0x0010' WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK            ]
    ['0x0020' READ_BROADCAST_DISTRIBUTION_TABLE_NAK             ]
    ['0x0030' REGISTER_FOREIGN_DEVICE_NAK                       ]
    ['0x0040' READ_FOREIGN_DEVICE_TABLE_NAK                     ]
    ['0x0050' DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK             ]
    ['0x0060' DISTRIBUTE_BROADCAST_TO_NETWORK_NAK               ]
]

// Attention: No proprietary extension allowed
[enum uint 8  NLMRejectRouterToNetworkRejectReason
    ['0'    OTHER                                           ]
    ['1'    NOT_DIRECTLY_CONNECTED                          ]
    ['2'    BUSY                                            ]
    ['3'    UNKNOWN_NLMT                                    ]
    ['4'    TOO_LONG                                        ]
    ['5'    SECURITY_ERROR                                  ]
    ['6'    ADDRESSING_ERROR                                ]
]

// Attention: No proprietary extension allowed
[enum uint 2 NPDUNetworkPriority
    ['3' LIFE_SAVETY_MESSAGE                                ]
    ['2' CRITICAL_EQUIPMENT_MESSAGE                         ]
    ['1' URGENT_MESSAGE                                     ]
    ['0' NORMAL_MESSAGE                                     ]
]

// Attention: No proprietary extension allowed
[enum uint 3 MaxSegmentsAccepted
    ['0x0' UNSPECIFIED              ]
    ['0x1' NUM_SEGMENTS_02          ]
    ['0x2' NUM_SEGMENTS_04          ]
    ['0x3' NUM_SEGMENTS_08          ]
    ['0x4' NUM_SEGMENTS_16          ]
    ['0x5' NUM_SEGMENTS_32          ]
    ['0x6' NUM_SEGMENTS_64          ]
    ['0x7' MORE_THAN_64_SEGMENTS    ]
]

// Attention: No proprietary extension allowed
[enum uint 4 MaxApduLengthAccepted
    ['0x0' MINIMUM_MESSAGE_SIZE     ] // 50 octets
    ['0x1' NUM_OCTETS_128           ]
    ['0x2' NUM_OCTETS_206           ] // fits in a LonTalk frame
    ['0x3' NUM_OCTETS_480           ] // fits in an ARCNET frame
    ['0x4' NUM_OCTETS_1024          ]
    ['0x5' NUM_OCTETS_1476          ] // fits in an Ethernet frame
    ['0x6' RESERVED_BY_ASHRAE_01    ]
    ['0x7' RESERVED_BY_ASHRAE_02    ]
    ['0x8' RESERVED_BY_ASHRAE_03    ]
    ['0x9' RESERVED_BY_ASHRAE_04    ]
    ['0xA' RESERVED_BY_ASHRAE_05    ]
    ['0xB' RESERVED_BY_ASHRAE_06    ]
    ['0xC' RESERVED_BY_ASHRAE_07    ]
    ['0xD' RESERVED_BY_ASHRAE_08    ]
    ['0xE' RESERVED_BY_ASHRAE_09    ]
    ['0xF' RESERVED_BY_ASHRAE_10    ]
]
