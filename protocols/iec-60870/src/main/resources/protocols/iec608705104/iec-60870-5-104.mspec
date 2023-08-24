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

// https://www.fit.vut.cz/research/publication-file/11570/TR-IEC104.pdf

[type IEC608705104Constants
    [const          uint 16     defaultPort 2404]
]

[discriminatedType APDU    byteOrder='LITTLE_ENDIAN'
    [const         uint 8  startByte    0x86               ]
    [implicit      uint 8  apciLength   'lengthInBytes - 2']
    [discriminator uint 16 command                         ]
    [typeSwitch command
        // U-Format Frames
        ['0x43' *UFormatTestFrameActivation
            [padding uint 8 pad '0x00' '2'                 ]
        ]
        ['0x83' *UFormatTestFrameConfirmation
            [padding uint 8 pad '0x00' '2'                 ]
        ]
        ['0x13' *UFormatStopDataTransferActivation
            [padding uint 8 pad '0x00' '2'                 ]
        ]
        ['0x23' *UFormatStopDataTransferConfirmation
            [padding uint 8 pad '0x00' '2'                 ]
        ]
        ['0x07' *UFormatStartDataTransferActivation
            [padding uint 8 pad '0x00' '2'                 ]
        ]
        ['0x0B' *UFormatStartDataTransferConfirmation
            [padding uint 8 pad '0x00' '2'                 ]
        ]

        // S-Format Frames
        ['0x01' *SFormat
            [simple  uint 16 receiveSequenceNo               ]
        ]

        // I-Format Frames (Catch-all for all other values)
        [        APDUIFormat
            // TODO: Fix this ...
            // [virtual uint 15 sendSequenceNo 'command >> 1']
            // TODO: Shift this right by one bit to make it a uint 15
            [simple  uint 16 receiveSequenceNo               ]
            // Payload
            [simple ASDU asdu]
        ]
    ]
]

[type ASDU
    [simple   TypeIdentification  typeIdentification                             ]
    [simple   bit                 structureQualifier                             ]
    [implicit uint 7              numberOfObjects     'COUNT(informationObjects)']
    [simple   bit                 test                                           ]
    [simple   bit                 negative                                       ]
    [simple   CauseOfTransmission causeOfTransmission                            ]
    [simple   uint 8              originatorAddress                              ]
    [simple   uint 16             asduAddressField                               ]
    [array    InformationObject   informationObjects  count 'numberOfObjects'    ]
]

// http://ijlalhaider.pbworks.com/w/file/fetch/64131148/Practical%20Modern%20SCADA%20Protocols.pdf

[discriminatedType InformationObject(TypeIdentification typeIdentification)
    [simple uint 24 address]
    [typeSwitch typeIdentification
        ['SINGLE_POINT_INFORMATION'                                 *_SINGLE_POINT_INFORMATION
// SIQ
        ]
        ['SINGLE_POINT_INFORMATION_WITH_TIME_TAG'                   *_SINGLE_POINT_INFORMATION_WITH_TIME_TAG
// SIQ + CP24Time2a
        ]
        ['DOUBLE_POINT_INFORMATION'                                 *_DOUBLE_POINT_INFORMATION
// DIQ
        ]
        ['DOUBLE_POINT_INFORMATION_WITH_TIME_TAG'                   *_DOUBLE_POINT_INFORMATION_WITH_TIME_TAG
// DIQ + CP24Time2a
        ]
        ['STEP_POSITION_INFORMATION'                                *_STEP_POSITION_INFORMATION
// VTI + QDS
        ]
        ['STEP_POSITION_INFORMATION_WITH_TIME_TAG'                  *_STEP_POSITION_INFORMATION_WITH_TIME_TAG
// VTI + QDS + CP24Time2a
        ]
        ['BITSTRING_OF_32_BIT'                                      *_BITSTRING_OF_32_BIT
// BSI + QDS
        ]
        ['BITSTRING_OF_32_BIT_WITH_TIME_TAG'                        *_BITSTRING_OF_32_BIT_WITH_TIME_TAG
// BSI + QDS + P24Time2a
        ]
        ['MEASURED_VALUE_NORMALISED_VALUE'                          *_MEASURED_VALUE_NORMALISED_VALUE
// NVA + QDS
        ]
        ['MEASURED_VALUE_NORMALIZED_VALUE_WITH_TIME_TAG'            *_MEASURED_VALUE_NORMALIZED_VALUE_WITH_TIME_TAG
// NVA + QDS + CP24Time2a
        ]
        ['MEASURED_VALUE_SCALED_VALUE'                              *_MEASURED_VALUE_SCALED_VALUE
// SVA + QDS
        ]
        ['MEASURED_VALUE_SCALED_VALUE_WIT_TIME_TAG'                 *_MEASURED_VALUE_SCALED_VALUE_WIT_TIME_TAG
// SVA + QDS + CP24Time2a
        ]
        ['MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER'               *_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER
// IEEE STD 754 + QDS
        ]
        ['MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG' *_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG
// IEEE STD 754 + QDS + CP24Time2a
        ]
        ['INTEGRATED_TOTALS'                                        *_INTEGRATED_TOTALS
// BCR
        ]
        ['INTEGRATED_TOTALS_WITH_TIME_TAG'                          *_INTEGRATED_TOTALS_WITH_TIME_TAG
// BCR + CP24Time2a
        ]
        ['EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG'              *_EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG
// CP16Time2a + CP24Time2a
        ]
        ['PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG' *_PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG
// SEP + QDP +CP16Time2a + CP24Time2a
        ]
        ['PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG' *_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG
// OCI + QDP + CP16Time2a + CP24Time2a
        ]
        ['PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION' *_PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION
// SCD + QDS
        ]
        ['MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR' *_MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR
// NVA
        ]
        ['SINGLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A'        *_SINGLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A
// SIQ + CP56Time2a
        ]
        ['DOUBLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A'        *_DOUBLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A
// DIQ + CP56Time2a
        ]
        ['STEP_POSITION_INFORMATION_WITH_TIME_TAG_CP56TIME2A'       *_STEP_POSITION_INFORMATION_WITH_TIME_TAG_CP56TIME2A
// VTI + QDS + CP56Time2a
        ]
        ['BITSTRING_OF_32_BIT_WITH_TIME_TAG_CP56TIME2A'             *_BITSTRING_OF_32_BIT_WITH_TIME_TAG_CP56TIME2A
// BSI + QDS + CP56Time2a
        ]
        ['MEASURED_VALUE_NORMALISED_VALUE_WITH_TIME_TAG_CP56TIME2A' *_MEASURED_VALUE_NORMALISED_VALUE_WITH_TIME_TAG_CP56TIME2A
// NVA + QDS + CP56Time2a
        ]
        ['MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG_CP56TIME2A'     *_MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG_CP56TIME2A
// SVA + QDS + CP56Time2a
        ]
        ['MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG_CP56TIME2A' *_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG_CP56TIME2A
// IEEE STD 754 + QDS + CP56Time2a
        ]
        ['INTEGRATED_TOTALS_WITH_TIME_TAG_CP56TIME2A'               *_INTEGRATED_TOTALS_WITH_TIME_TAG_CP56TIME2A
// BCR + CP56Time2a
        ]
        ['EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A'   *_EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A
// CP16Time2a + CP56Time2a
        ]
        ['PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A' *_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A
// SEP + QDP + CP16Time2a + CP56Time2a
        ]
        ['PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A' *_PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A
// OCI + QDP + CP16Time2a + CP56Time2a
        ]
        ['SINGLE_COMMAND'                                           *_SINGLE_COMMAND
// SCO
        ]
        ['DOUBLE_COMMAND'                                           *_DOUBLE_COMMAND
// DCO
        ]
        ['REGULATING_STEP_COMMAND'                                  *_REGULATING_STEP_COMMAND
// RCO
        ]
        ['SET_POINT_COMMAND_NORMALISED_VALUE'                       *_SET_POINT_COMMAND_NORMALISED_VALUE
// NVA + QOS
        ]
        ['SET_POINT_COMMAND_SCALED_VALUE'                           *_SET_POINT_COMMAND_SCALED_VALUE
// SVA + QOS
        ]
        ['SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER'            *_SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER
// 1 IEEE STD 754 + QOS
        ]
        ['BITSTRING_32_BIT_COMMAND'                                 *_BITSTRING_32_BIT_COMMAND
// BSI
        ]
        ['SINGLE_COMMAND_WITH_TIME_TAG_CP56TIME2A'                  *_SINGLE_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['DOUBLE_COMMAND_WITH_TIME_TAG_CP56TIME2A'                  *_DOUBLE_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['REGULATING_STEP_COMMAND_WITH_TIME_TAG_CP56TIME2A'         *_REGULATING_STEP_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['MEASURED_VALUE_NORMALISED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A' *_MEASURED_VALUE_NORMALISED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['MEASURED_VALUE_SCALED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A' *_MEASURED_VALUE_SCALED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_COMMAND_WITH_TIME_TAG_CP56TIME2A' *_MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['BITSTRING_OF_32_BIT_COMMAND_WITH_TIME_TAG_CP56TIME2A'     *_BITSTRING_OF_32_BIT_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['END_OF_INITIALISATION'                                    *_END_OF_INITIALISATION
// COI
        ]
        ['INTERROGATION_COMMAND'                                    *_INTERROGATION_COMMAND
// QOI
        ]
        ['COUNTER_INTERROGATION_COMMAND'                            *_COUNTER_INTERROGATION_COMMAND
// QCC
        ]
        ['READ_COMMAND'                                             *_READ_COMMAND
        ]
        ['CLOCK_SYNCHRONISATION_COMMAND'                            *_CLOCK_SYNCHRONISATION_COMMAND
// CP56Time2a
        ]
        ['TEST_COMMAND'                                             *_TEST_COMMAND
// FBP
        ]
        ['RESET_PROCESS_COMMAND'                                    *_RESET_PROCESS_COMMAND
// QRP
        ]
        ['DELAY_ACQUISITION_COMMAND'                                *_DELAY_ACQUISITION_COMMAND
// CP16Time2a
        ]
        ['TEST_COMMAND_WITH_TIME_TAG_CP56TIME2A'                    *_TEST_COMMAND_WITH_TIME_TAG_CP56TIME2A
        ]
        ['PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE'            *_PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE
// NVA + QPM
        ]
        ['PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE'                *_PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE
// SVA + QPM
        ]
        ['PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER' *_PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER
// IEEE STD 754 + QPM
        ]
        ['PARAMETER_ACTIVATION'                                     *_PARAMETER_ACTIVATION
// QPA
        ]
        ['FILE_READY'                                               *_FILE_READY
// NOF + LOF + FRQ
        ]
        ['SECTION_READY'                                            *_SECTION_READY
// NOF + NOS + LOF + SRQ
        ]
        ['CALL_DIRECTORY_SELECT_FILE_CALL_FILE_CALL_SECTION'        *_CALL_DIRECTORY_SELECT_FILE_CALL_FILE_CALL_SECTION
// NOF + NOS + SCQ
        ]
        ['LAST_SECTION_LAST_SEGMENT'                                *_LAST_SECTION_LAST_SEGMENT
// NOF + NOS + LSQ + CHS
        ]
        ['ACK_FILE_ACK_SECTION'                                     *_ACK_FILE_ACK_SECTION
// NOF + NOS + AFQ
        ]
        ['SEGMENT'                                                  *_SEGMENT
// NOF + NOS + LOS + segment
        ]
        ['DIRECTORY'                                                *_DIRECTORY
// NOF + LOF + SOF + CP56Time2a
        ]
    ]
]

//////////////////////////////////////////////////////////
// Process related information elements
//////////////////////////////////////////////////////////

// SIQ
[type SinglePointInformation
    [simple   bit    invalid            ]
    [simple   bit    notTopical         ]
    [simple   bit    substituted        ]
    [simple   bit    blocked            ]
    [reserved uint 3 '0'                ]
    [simple   bit    stausOn            ]
]

// DIQ
[type DoublePointInformation
    [simple   bit    invalid            ]
    [simple   bit    notTopical         ]
    [simple   bit    substituted        ]
    [simple   bit    blocked            ]
    [reserved uint 2 '0'                ]
    [simple   uint 2 dpiCode            ] // TODO: Possible ENUM
]

// BSI
[type BinaryStateInformation
    [simple uint 32 bits] // TODO: Possibly bit-string
]

// SCD
[type StatusChangeDetection
    [simple uint 32 bits] // TODO: Possibly bit-string
]

// QDS
[type QualityDescriptor
    [simple   bit    invalid            ]
    [simple   bit    notTopical         ]
    [simple   bit    substituted        ]
    [simple   bit    blocked            ]
    [reserved uint 3 '0'                ]
    [simple   bit    overflow           ]
]

// VTI
[type ValueWithTransientStateIndication
    [simple   bit    transientState     ]
    [simple   uint 7 value              ]
]

// NVA
[type NormalizedValue
    // TODO: F16
]

// SVA
[type ScaledValue
    [simple   int 16 value               ]
]

// R32
[type ShortFloatingPointNumber
   [simple  float 16 value               ] // TODO: Double-Check
]

// BCR
[type BinaryCounterReading
    [simple   uint 32 counterValue       ]
    [simple   bit     counterValid       ]
    [simple   bit     counterAdjusted    ]
    [simple   bit     carry              ]
    [simple   uint 5  sequenceNumber     ]
]

// SEP
[type SingleEventOfProtectionEquipment
    [simple   bit    invalid            ]
    [simple   bit    notTopical         ]
    [simple   bit    substituted        ]
    [simple   bit    blocked            ]
    [simple   bit    elapsedTimeInvalid ]
    [reserved uint 1 '0'                ]
    [simple   uint 2 eventState         ] // TODO: Possible ENUM
]

// SPE
[type StartEventsOfProtectionEquipment
    [reserved uint 2 '0'                    ]
    [simple   bit    startOfOperationInReverseDirection]
    [simple   bit    startOfOperationIE     ]
    [simple   bit    stateOfOperationPhaseL3]
    [simple   bit    stateOfOperationPhaseL2]
    [simple   bit    stateOfOperationPhaseL1]
    [simple   bit    generalStartOfOperation]
]

// OCI
[type OutputCircuitInformation
    [reserved uint 4 '0'                    ]
    [simple   bit    stateOfOperationPhaseL3]
    [simple   bit    stateOfOperationPhaseL2]
    [simple   bit    stateOfOperationPhaseL1]
    [simple   bit    generalStartOfOperation]
]

// QDP
[type QualityDescriptorForPointsOfProtectionEquipment
    [simple   bit    invalid            ]
    [simple   bit    notTopical         ]
    [simple   bit    substituted        ]
    [simple   bit    blocked            ]
    [simple   bit    elapsedTimeInvalid ]
    [reserved uint 3 '0'                ]
]

//////////////////////////////////////////////////////////
// Command information elements
//////////////////////////////////////////////////////////

// SCO
[type SingleCommand
    // TODO: Implement
]

// DCO
[type DoubleCommand
    // TODO: Implement
]

// RCO
[type RegulatingStepCommand
    // TODO: Implement
]

//////////////////////////////////////////////////////////
// Time information elements
//////////////////////////////////////////////////////////

// CP56Time2a
[type SevenOctetBinaryTime
    // TODO: Implement
]

// CP24Time2a
[type ThreeOctetBinaryTime
    // TODO: Implement
]

// CP16Time2av
[type TwoOctetBinaryTime
    // TODO: Implement
]

//////////////////////////////////////////////////////////
// Qualifier information elements
//////////////////////////////////////////////////////////

// QOI
[type QualifierOfInterrogation
    [simple   uint 8 qualifierOfCommand         ] // TODO: Possible ENUM
]

// QCC
[type QualifierOfCounterInterrogationCommand
    [simple   uint 2 freeze         ] // TODO: Possible ENUM
    [simple   uint 6 request        ] // TODO: Possible ENUM
]

// QPM
[type QualifierOfParameterOfMeasuredValues
    [simple   bit    parameterInOperation]
    [simple   bit    localParameterChange]
    [simple   uint 6 kindOfParameter     ] // TODO: Possible ENUM
]

// QPA
[type QualifierOfParameterActivation
    [simple   uint 8 qualifier         ] // TODO: Possible ENUM
]

// QOC
// TODO: Only 6 bit long !!!!!!
[type QualifierOfCommand
    [simple   bit    select]
    [simple   uint 5 qualifier         ] // TODO: Possible ENUM
]

// QRP
[type QualifierOfResetProcessCommand
    [simple   uint 8 qualifier         ] // TODO: Possible ENUM
]

// QOS
[type QualifierOfSetPointCommand
    [simple   bit    select]
    [simple   uint 7 qualifier         ] // TODO: Possible ENUM
]

//////////////////////////////////////////////////////////
// File transfer information elements
//////////////////////////////////////////////////////////

// FRQ
[type FileReadyQualifier
    // TODO: Implement
]

// SRQ
[type SectionReadyQualifier
    // TODO: Implement
]

// SCQ
[type SelectAndCallQualifier
    // TODO: Implement
]

// LSQ
[type LastSectionOrSegmentQualifier
    // TODO: Implement
]

// AFQ
[type AcknowledgeFileOrSectionQualifier
    // TODO: Implement
]

// NOF
[type NameOfFile
    // TODO: Implement
]

// NOS
[type NameOfSection
    // TODO: Implement
]

// LOF
[type LengthOfFile
    // TODO: Implement
]

// LOS
[type LengthOfSegment
    // TODO: Implement
]

// CHS
[type Checksum
    // TODO: Implement
]

// SOF
[type StatusOfFile
    // TODO: Implement
]

//////////////////////////////////////////////////////////
// Miscellaneous information elements
//////////////////////////////////////////////////////////

// COI
[type CauseOfInitialization
    [simple   bit    select]
    [simple   uint 7 qualifier         ] // TODO: Possible ENUM
]

// FBP
[type FixedTestBitPatternTwoOctet
    [simple uint 16 pattern] // TODO: Possibly bit-string
]

[enum uint 8 TypeIdentification
    ['0x00' NOT_USED                                                                          ]
    ['0x01' SINGLE_POINT_INFORMATION                                                          ]
    ['0x02' SINGLE_POINT_INFORMATION_WITH_TIME_TAG                                            ]
    ['0x03' DOUBLE_POINT_INFORMATION                                                          ]
    ['0x04' DOUBLE_POINT_INFORMATION_WITH_TIME_TAG                                            ]
    ['0x05' STEP_POSITION_INFORMATION                                                         ]
    ['0x06' STEP_POSITION_INFORMATION_WITH_TIME_TAG                                           ]
    ['0x07' BITSTRING_OF_32_BIT                                                               ]
    ['0x08' BITSTRING_OF_32_BIT_WITH_TIME_TAG                                                 ]
    ['0x09' MEASURED_VALUE_NORMALISED_VALUE                                                   ]
    ['0x0A' MEASURED_VALUE_NORMALIZED_VALUE_WITH_TIME_TAG                                     ]
    ['0x0B' MEASURED_VALUE_SCALED_VALUE                                                       ]
    ['0x0C' MEASURED_VALUE_SCALED_VALUE_WIT_TIME_TAG                                          ]
    ['0x0D' MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER                                        ]
    ['0x0E' MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG                          ]
    ['0x0F' INTEGRATED_TOTALS                                                                 ]
    ['0x10' INTEGRATED_TOTALS_WITH_TIME_TAG                                                   ]
    ['0x11' EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG                                       ]
    ['0x12' PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG                         ]
    ['0x13' PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG           ]
    ['0x14' PACKED_SINGLE_POINT_INFORMATION_WITH_STATUS_CHANGE_DETECTION                      ]
    ['0x15' MEASURED_VALUE_NORMALIZED_VALUE_WITHOUT_QUALITY_DESCRIPTOR                        ]
    ['0x1E' SINGLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A                                 ]
    ['0x1F' DOUBLE_POINT_INFORMATION_WITH_TIME_TAG_CP56TIME2A                                 ]
    ['0x20' STEP_POSITION_INFORMATION_WITH_TIME_TAG_CP56TIME2A                                ]
    ['0x21' BITSTRING_OF_32_BIT_WITH_TIME_TAG_CP56TIME2A                                      ]
    ['0x22' MEASURED_VALUE_NORMALISED_VALUE_WITH_TIME_TAG_CP56TIME2A                          ]
    ['0x23' MEASURED_VALUE_SCALED_VALUE_WITH_TIME_TAG_CP56TIME2A                              ]
    ['0x24' MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_WITH_TIME_TAG_CP56TIME2A               ]
    ['0x25' INTEGRATED_TOTALS_WITH_TIME_TAG_CP56TIME2A                                        ]
    ['0x26' EVENT_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A                            ]
    ['0x27' PACKED_START_EVENTS_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A              ]
    ['0x28' PACKED_OUTPUT_CIRCUIT_INFORMATION_OF_PROTECTION_EQUIPMENT_WITH_TIME_TAG_CP56TIME2A]
    ['0x2D' SINGLE_COMMAND                                                                    ]
    ['0x2E' DOUBLE_COMMAND                                                                    ]
    ['0x2F' REGULATING_STEP_COMMAND                                                           ]
    ['0x30' SET_POINT_COMMAND_NORMALISED_VALUE                                                ]
    ['0x31' SET_POINT_COMMAND_SCALED_VALUE                                                    ]
    ['0x32' SET_POINT_COMMAND_SHORT_FLOATING_POINT_NUMBER                                     ]
    ['0x33' BITSTRING_32_BIT_COMMAND                                                          ]
    ['0x3A' SINGLE_COMMAND_WITH_TIME_TAG_CP56TIME2A                                           ]
    ['0x3B' DOUBLE_COMMAND_WITH_TIME_TAG_CP56TIME2A                                           ]
    ['0x3C' REGULATING_STEP_COMMAND_WITH_TIME_TAG_CP56TIME2A                                  ]
    ['0x3D' MEASURED_VALUE_NORMALISED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A                  ]
    ['0x3E' MEASURED_VALUE_SCALED_VALUE_COMMAND_WITH_TIME_TAG_CP56TIME2A                      ]
    ['0x3F' MEASURED_VALUE_SHORT_FLOATING_POINT_NUMBER_COMMAND_WITH_TIME_TAG_CP56TIME2A       ]
    ['0x40' BITSTRING_OF_32_BIT_COMMAND_WITH_TIME_TAG_CP56TIME2A                              ]
    ['0x46' END_OF_INITIALISATION                                                             ]
    ['0x64' INTERROGATION_COMMAND                                                             ]
    ['0x65' COUNTER_INTERROGATION_COMMAND                                                     ]
    ['0x66' READ_COMMAND                                                                      ]
    ['0x67' CLOCK_SYNCHRONISATION_COMMAND                                                     ]
    ['0x68' TEST_COMMAND                                                                      ]
    ['0x69' RESET_PROCESS_COMMAND                                                             ]
    ['0x6A' DELAY_ACQUISITION_COMMAND                                                         ]
    ['0x6B' TEST_COMMAND_WITH_TIME_TAG_CP56TIME2A                                             ]
    ['0x6E' PARAMETER_OF_MEASURED_VALUES_NORMALIZED_VALUE                                     ]
    ['0x6F' PARAMETER_OF_MEASURED_VALUES_SCALED_VALUE                                         ]
    ['0x70' PARAMETER_OF_MEASURED_VALUES_SHORT_FLOATING_POINT_NUMBER                          ]
    ['0x71' PARAMETER_ACTIVATION                                                              ]
    ['0x78' FILE_READY                                                                        ]
    ['0x79' SECTION_READY                                                                     ]
    ['0x7A' CALL_DIRECTORY_SELECT_FILE_CALL_FILE_CALL_SECTION                                 ]
    ['0x7B' LAST_SECTION_LAST_SEGMENT                                                         ]
    ['0x7C' ACK_FILE_ACK_SECTION                                                              ]
    ['0x7D' SEGMENT                                                                           ]
    ['0x7E' DIRECTORY                                                                         ]
]

[enum uint 6 CauseOfTransmission
    ['0x01' PERIODIC_CYCLIC_PER_CYC                                  ]
    ['0x02' BACKGROUND_INTERROGATION_BACK                            ]
    ['0x03' SPONTANEOUS_SPONT                                        ]
    ['0x04' INITIALIZED_INIT                                         ]
    ['0x05' INTERROGATION_OR_INTERROGATED_REQ                        ]
    ['0x06' ACTIVATION_ACT                                           ]
    ['0x07' CONFIRMATION_ACTIVATION_ACTCON                           ]
    ['0x08' DEACTIVATION_DEACT                                       ]
    ['0x09' CONFIRMATION_DEACTIVATION_DEACTCON                       ]
    ['0x0A' TERMINATION_ACTIVATION_ACTTERM                           ]
    ['0x0B' FEEDBACK_CAUSED_BY_DISTANT_COMMAND_RETREM                ]
    ['0x0C' FEEDBACK_CAUSED_BY_LOCAL_COMMAND_RETLOC                  ]
    ['0x0D' DATA_TRANSMISSION_FILE                                   ]
    ['0x14' INTERROGATED_BY_GENERAL_INTERROGATION_INROGEN            ]
    ['0x15' INTERROGATED_BY_INTERROGATION_GROUP_1_INRO1              ]
    ['0x16' INTERROGATED_BY_INTERROGATION_GROUP_2_INRO2              ]
    ['0x17' INTERROGATED_BY_INTERROGATION_GROUP_3_INRO3              ]
    ['0x18' INTERROGATED_BY_INTERROGATION_GROUP_4_INRO4              ]
    ['0x19' INTERROGATED_BY_INTERROGATION_GROUP_5_INRO5              ]
    ['0x1A' INTERROGATED_BY_INTERROGATION_GROUP_6_INRO6              ]
    ['0x1B' INTERROGATED_BY_INTERROGATION_GROUP_7_INRO7              ]
    ['0x1C' INTERROGATED_BY_INTERROGATION_GROUP_8_INRO8              ]
    ['0x1D' INTERROGATED_BY_INTERROGATION_GROUP_9_INRO9              ]
    ['0x1E' INTERROGATED_BY_INTERROGATION_GROUP_10_INRO10            ]
    ['0x1F' INTERROGATED_BY_INTERROGATION_GROUP_11_INRO11            ]
    ['0x20' INTERROGATED_BY_INTERROGATION_GROUP_12_INRO12            ]
    ['0x21' INTERROGATED_BY_INTERROGATION_GROUP_13_INRO13            ]
    ['0x22' INTERROGATED_BY_INTERROGATION_GROUP_14_INRO14            ]
    ['0x23' INTERROGATED_BY_INTERROGATION_GROUP_15_INRO15            ]
    ['0x24' INTERROGATED_BY_INTERROGATION_GROUP_16_INRO16            ]
    ['0x25' INTERROGATED_BY_COUNTER_GENERAL_INTERROGATION_REQCOGEN   ]
    ['0x26' INTERROGATED_BY_INTERROGATION_COUNTER_GROUP_1_REQCO1     ]
    ['0x27' INTERROGATED_BY_INTERROGATION_COUNTER_GROUP_2_REQCO2     ]
    ['0x28' INTERROGATED_BY_INTERROGATION_COUNTER_GROUP_3_REQCO3     ]
    ['0x29' INTERROGATED_BY_INTERROGATION_COUNTER_GROUP_4_REQCO4     ]
    ['0x2C' TYPE_IDENTIFICATION_UNKNOWN_UKNOWN_TYPE                  ]
    ['0x2D' CAUSE_UNKNOWN_UKNOWN_CAUSE                               ]
    ['0x2E' ASDU_ADDRESS_UNKNOWN_UNKNOWN_ASDU_ADDRESS                ]
    ['0x2F' INFORMATION_OBJECT_ADDRESS_UNKNOWN_UNKNOWN_OBJECT_ADDRESS]
]