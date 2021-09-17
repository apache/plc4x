/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

[discriminatedType 'BVLC'
    [const         uint 8  'bacnetType'   '0x81'       ]
    [discriminator uint 8  'bvlcFunction'              ]
    [implicit      uint 16 'bvlcLength' 'lengthInBytes']
    [typeSwitch 'bvlcFunction'
        ['0x00' BVLCResult
        ]
        ['0x01' BVLCWideBroadcastDistributionTable
        ]
        ['0x02' BVLCReadBroadcastDistributionTable
        ]
        ['0x03' BVLCReadBroadcastDistributionTableAck
        ]
        ['0x04' BVLCForwardedNPDU [uint 16 'bvlcLength']
            [array  uint 8  'ip'    count '4'         ]
            [simple uint 16 'port'                    ]
            [simple NPDU    'npdu' ['bvlcLength - 10']]
        ]
        ['0x05' BVLCRegisterForeignDevice
        ]
        ['0x06' BVLCReadForeignDeviceTable
        ]
        ['0x07' BVLCReadForeignDeviceTableAck
        ]
        ['0x08' BVLCDeleteForeignDeviceTableEntry
        ]
        ['0x09' BVLCDistributeBroadcastToNetwork
        ]
        ['0x0A' BVLCOriginalUnicastNPDU [uint 16 'bvlcLength']
            [simple NPDU 'npdu' ['bvlcLength - 4']]
        ]
        ['0x0B' BVLCOriginalBroadcastNPDU [uint 16 'bvlcLength']
            [simple NPDU 'npdu' ['bvlcLength - 4']]
        ]
        ['0x0C' BVLCSecureBVLL
        ]
    ]
]

[type 'NPDU' [uint 16 'npduLength']
    [simple   uint 8        'protocolVersionNumber']
    [simple   bit           'messageTypeFieldPresent']
    [reserved uint 1        '0']
    [simple   bit           'destinationSpecified']
    [reserved uint 1        '0']
    [simple   bit           'sourceSpecified']
    [simple   bit           'expectingReply']
    [simple   uint 2        'networkPriority']
    [optional uint 16       'destinationNetworkAddress' 'destinationSpecified']
    [optional uint 8        'destinationLength'         'destinationSpecified']
    [array    uint 8        'destinationAddress' count  'destinationSpecified ? destinationLength : 0']
    [optional uint 16       'sourceNetworkAddress'      'sourceSpecified']
    [optional uint 8        'sourceLength'              'sourceSpecified']
    [array    uint 8        'sourceAddress' count       'sourceSpecified ? sourceLength : 0']
    [optional uint 8        'hopCount'                  'destinationSpecified']
    [optional NLM           'nlm'                       'messageTypeFieldPresent'  ['npduLength - (2 + (sourceSpecified ? 3 + sourceLength : 0) + (destinationSpecified ? 3 + destinationLength: 0) + ((destinationSpecified || sourceSpecified) ? 1 : 0))']]
    [optional APDU          'apdu'                      '!messageTypeFieldPresent' ['npduLength - (2 + (sourceSpecified ? 3 + sourceLength : 0) + (destinationSpecified ? 3 + destinationLength: 0) + ((destinationSpecified || sourceSpecified) ? 1 : 0))']]
]

[discriminatedType 'NLM' [uint 16 'apduLength']
    [discriminator uint 8  'messageType']
    [optional      uint 16 'vendorId' '(messageType >= 128) && (messageType <= 255)']
    [typeSwitch 'messageType'
        ['0x0' NLMWhoIsRouterToNetwork [uint 16 'apduLength', uint 8  'messageType']
            [array uint 16 'destinationNetworkAddress' length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
        ['0x1' NLMIAmRouterToNetwork [uint 16 'apduLength', uint 8  'messageType']
            [array uint 16 'destinationNetworkAddress' length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
    ]
]

[discriminatedType 'APDU' [uint 16 'apduLength']
    [discriminator uint 4 'apduType']
    [typeSwitch 'apduType'
        ['0x0' APDUConfirmedRequest [uint 16 'apduLength']
            [simple   bit    'segmentedMessage'                       ]
            [simple   bit    'moreFollows'                            ]
            [simple   bit    'segmentedResponseAccepted'              ]
            [reserved uint 2 '0'                                      ]
            [simple   uint 3 'maxSegmentsAccepted'                    ]
            [simple   uint 4 'maxApduLengthAccepted'                  ]
            [simple   uint 8 'invokeId'                               ]
            [optional uint 8 'sequenceNumber'       'segmentedMessage']
            [optional uint 8 'proposedWindowSize'   'segmentedMessage']
            [simple   BACnetConfirmedServiceRequest 'serviceRequest'  ['apduLength - (3 + (segmentedMessage ? 2 : 0))']]
        ]
        ['0x1' APDUUnconfirmedRequest [uint 16 'apduLength']
            [reserved uint 4                          '0'             ]
            [simple   BACnetUnconfirmedServiceRequest 'serviceRequest' ['apduLength - 1']]
        ]
        ['0x2' APDUSimpleAck
            [reserved uint 4 '0'               ]
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'serviceChoice'   ]
        ]
        ['0x3' APDUComplexAck
            [simple   bit               'segmentedMessage'                     ]
            [simple   bit               'moreFollows'                          ]
            [reserved uint 2            '0'                                    ]
            [simple   uint 8            'originalInvokeId'                     ]
            [optional uint 8            'sequenceNumber'     'segmentedMessage']
            [optional uint 8            'proposedWindowSize' 'segmentedMessage']
            [simple   BACnetServiceAck  'serviceAck'                           ]
        ]
        ['0x4' APDUSegmentAck
            [reserved uint 2 '0x00'              ]
            [simple   bit    'negativeAck'       ]
            [simple   bit    'server'            ]
            [simple   uint 8 'originalInvokeId'  ]
            [simple   uint 8 'sequenceNumber'    ]
            [simple   uint 8 'proposedWindowSize']
        ]
        ['0x5' APDUError
            [reserved uint 4      '0x00'            ]
            [simple   uint 8      'originalInvokeId']
            [simple   BACnetError 'error'           ]
        ]
        ['0x6' APDUReject
            [reserved uint 4 '0x00'            ]
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'rejectReason'    ]
        ]
        ['0x7' APDUAbort
            [reserved uint 3 '0x00'            ]
            [simple   bit    'server'          ]
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'abortReason'     ]
        ]
    ]
]

[discriminatedType 'BACnetConfirmedServiceRequest' [uint 16 'len']
    [discriminator uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x00' BACnetConfirmedServiceRequestAcknowledgeAlarm
        ]
        ['0x01' BACnetConfirmedServiceRequestConfirmedCOVNotification [uint 16 'len']
            [const  uint 8               'subscriberProcessIdentifierHeader'         '0x09'                 ]
            [simple uint 8               'subscriberProcessIdentifier'                                      ]
            [const  uint 8               'monitoredObjectIdentifierHeader'           '0x1C'                 ]
            [simple uint 10              'monitoredObjectType'                                              ]
            [simple uint 22              'monitoredObjectInstanceNumber'                                    ]
            [const  uint 8               'issueConfirmedNotificationsHeader'         '0x2C'                 ]
            [simple uint 10              'issueConfirmedNotificationsType'                                  ]
            [simple uint 22              'issueConfirmedNotificationsInstanceNumber'                        ]
            [const  uint 5               'lifetimeHeader'                            '0x07'                 ]
            [simple uint 3               'lifetimeLength'                                                   ]
            [array  int  8               'lifetimeSeconds'                           count  'lifetimeLength']
            [const  uint 8               'listOfValuesOpeningTag'                    '0x4E'                 ]
            [array  BACnetTagWithContent 'notifications'                             length 'len - 18'      ]
            [const  uint 8               'listOfValuesClosingTag'                    '0x4F'                 ]
        ]
        ['0x02' BACnetConfirmedServiceRequestConfirmedEventNotification
        ]

        ['0x04' BACnetConfirmedServiceRequestGetEnrollmentSummary
        ]
        ['0x05' BACnetConfirmedServiceRequestSubscribeCOV
            [const  uint 8  'subscriberProcessIdentifierHeader'   '0x09'                ]
            [simple uint 8  'subscriberProcessIdentifier'                               ]
            [const  uint 8  'monitoredObjectIdentifierHeader'     '0x1C'                ]
            [simple uint 10 'monitoredObjectType'                                       ]
            [simple uint 22 'monitoredObjectInstanceNumber'                             ]
            [const  uint 8  'issueConfirmedNotificationsHeader'   '0x29'                ]
            [const  uint 7  'issueConfirmedNotificationsSkipBits' '0x00'                ]
            [simple bit     'issueConfirmedNotifications'                               ]
            [const  uint 5  'lifetimeHeader'                      '0x07'                ]
            [simple uint 3  'lifetimeLength'                                            ]
            [array  int 8   'lifetimeSeconds'                     count 'lifetimeLength']
        ]

        ['0x06' BACnetConfirmedServiceRequestAtomicReadFile
        ]
        ['0x07' BACnetConfirmedServiceRequestAtomicWriteFile
        ]

        ['0x08' BACnetConfirmedServiceRequestAddListElement
        ]
        ['0x09' BACnetConfirmedServiceRequestRemoveListElement
        ]
        ['0x0A' BACnetConfirmedServiceRequestCreateObject
        ]
        ['0x0B' BACnetConfirmedServiceRequestDeleteObject
        ]
        ['0x0C' BACnetConfirmedServiceRequestReadProperty
            [const  uint 8  'objectIdentifierHeader'   '0x0C'                          ]
            [simple uint 10 'objectType'                                               ]
            [simple uint 22 'objectInstanceNumber'                                     ]
            [const  uint 5  'propertyIdentifierHeader' '0x03'                          ]
            [simple uint 3  'propertyIdentifierLength'                                 ]
            [array  int 8   'propertyIdentifier'       count 'propertyIdentifierLength']
        ]
        ['0x0E' BACnetConfirmedServiceRequestReadPropertyMultiple
        ]
        ['0x0F' BACnetConfirmedServiceRequestWriteProperty [uint 16 'len']
            [const    uint 8    'objectIdentifierHeader'    '0x0C'                          ]
            [simple   uint 10   'objectType'                                                ]
            [simple   uint 22   'objectInstanceNumber'                                      ]
            [const    uint 5    'propertyIdentifierHeader' '0x03'                           ]
            [simple   uint 3    'propertyIdentifierLength'                                  ]
            [array    int 8     'propertyIdentifier'        count 'propertyIdentifierLength']
            [const    uint 8    'openingTag'                '0x3E'                          ]
            [simple   BACnetTag 'value'                                                     ]
            [const    uint 8    'closingTag'                '0x3F'                          ]
            [optional BACnetTag 'priority'                  'curPos < (len - 1)'            ]
        ]
        ['0x10' BACnetConfirmedServiceRequestWritePropertyMultiple
        ]

        ['0x11' BACnetConfirmedServiceRequestDeviceCommunicationControl
        ]
        ['0x12' BACnetConfirmedServiceRequestConfirmedPrivateTransfer
        ]
        ['0x13' BACnetConfirmedServiceRequestConfirmedTextMessage
        ]
        ['0x14' BACnetConfirmedServiceRequestReinitializeDevice
        ]

        ['0x15' BACnetConfirmedServiceRequestVTOpen
        ]
        ['0x16' BACnetConfirmedServiceRequestVTClose
        ]
        ['0x17' BACnetConfirmedServiceRequestVTData
        ]

        ['0x18' BACnetConfirmedServiceRequestRemovedAuthenticate
        ]
        ['0x19' BACnetConfirmedServiceRequestRemovedRequestKey
        ]
        ['0x0D' BACnetConfirmedServiceRequestRemovedReadPropertyConditional
        ]

        ['0x1A' BACnetConfirmedServiceRequestReadRange
        ]
        ['0x1B' BACnetConfirmedServiceRequestLifeSafetyOperation
        ]
        ['0x1C' BACnetConfirmedServiceRequestSubscribeCOVProperty

        ]
        ['0x1D' BACnetConfirmedServiceRequestGetEventInformation
        ]

        ['0x1E' BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple
        ]
        ['0x1F' BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple
        ]
    ]
]

[discriminatedType 'BACnetUnconfirmedServiceRequest' [uint 16 'len']
    [discriminator uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x00' BACnetUnconfirmedServiceRequestIAm
            [const uint 8 'objectIdentifierHeader' '0xC4']
            [simple uint 10 'objectType']
            [simple uint 22 'objectInstanceNumber']
            [const uint 5 'maximumApduLengthAcceptedHeader' '0x04']
            [simple uint 3 'maximumApduLengthAcceptedLength']
            [array int 8 'maximumApduLengthAccepted' count 'maximumApduLengthAcceptedLength']
            [const uint 8 'segmentationSupportedHeader' '0x91']
            [simple uint 8 'segmentationSupported']
            [const uint 8 'vendorIdHeader' '0x21']
            [simple uint 8 'vendorId']
        ]
        ['0x01' BACnetUnconfirmedServiceRequestIHave
        ]
        ['0x02' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
        ]
        ['0x03' BACnetUnconfirmedServiceRequestUnconfirmedEventNotification
        ]
        ['0x04' BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer [uint 16 'len']
            [const uint 8 'vendorIdHeader' '0x09']
            [simple uint 8 'vendorId']
            [const uint 8 'serviceNumberHeader' '0x1A']
            [simple uint 16 'serviceNumber']
            [const uint 8 'listOfValuesOpeningTag' '0x2E']
            [array int 8 'values' length 'len - 8']
            [const uint 8 'listOfValuesClosingTag' '0x2F']
        ]
        ['0x05' BACnetUnconfirmedServiceRequestUnconfirmedTextMessage
        ]
        ['0x06' BACnetUnconfirmedServiceRequestTimeSynchronization
        ]
        ['0x07' BACnetUnconfirmedServiceRequestWhoHas
            [const uint 8 'deviceInstanceLowLimitHeader' '0x0B'         ]
            [simple uint 24 'deviceInstanceLowLimit'                    ]
            [const uint 8 'deviceInstanceHighLimitHeader' '0x1B'        ]
            [simple uint 24 'deviceInstanceHighLimit'                   ]
            [const uint 8 'objectNameHeader' '0x3D'                     ]
            [implicit uint 8 'objectNameLength' 'COUNT(objectName) + 1' ]
            [simple uint 8 'objectNameCharacterSet'                     ]
            [array int 8 'objectName' length 'objectNameLength - 1'     ]
        ]
        ['0x08' BACnetUnconfirmedServiceRequestWhoIs
            // TODO: here we need proper bacnet tags (like a dicriminator etc... see line 494 BACnetTag)
            [const uint 5 'deviceInstanceRangeLowLimitHeader' '0x01']
            [simple uint 3 'deviceInstanceRangeLowLimitLength']
            [array int 8 'deviceInstanceRangeLowLimit' count 'deviceInstanceRangeLowLimitLength']
            [const uint 5 'deviceInstanceRangeHighLimitHeader' '0x03']
            [simple uint 3 'deviceInstanceRangeHighLimitLength']
            [array int 8 'deviceInstanceRangeHighLimit' count 'deviceInstanceRangeHighLimitLength']
        ]
        ['0x09' BACnetUnconfirmedServiceRequestUTCTimeSynchronization
        ]
        ['0x0A' BACnetUnconfirmedServiceRequestWriteGroup
        ]
        ['0x0B' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple
        ]
    ]
]

[discriminatedType 'BACnetServiceAck'
    [discriminator   uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x03' BACnetServiceAckGetAlarmSummary

        ]
        ['0x04' BACnetServiceAckGetEnrollmentSummary

        ]
        ['0x1D' BACnetServiceAckGetEventInformation

        ]

        ['0x06' BACnetServiceAckAtomicReadFile

        ]
        ['0x07' BACnetServiceAckAtomicWriteFile

        ]

        ['0x0A' BACnetServiceAckCreateObject

        ]
        ['0x0C' BACnetServiceAckReadProperty
            [const uint 8 'objectIdentifierHeader' '0x0C']
            [simple uint 10 'objectType']
            [simple uint 22 'objectInstanceNumber']
            [const uint 5 'propertyIdentifierHeader' '0x03']
            [simple uint 3 'propertyIdentifierLength']
            [array int 8 'propertyIdentifier' count 'propertyIdentifierLength']
            [const uint 8 'openingTag' '0x3E']
            [simple BACnetTag 'value']
            [const uint 8 'closingTag' '0x3F']
        ]
        ['0x0E' BACnetServiceAckReadPropertyMultiple

        ]
        ['0x1A' BACnetServiceAckReadRange

        ]

        ['0x12' BACnetServiceAckConfirmedPrivateTransfer

        ]

        ['0x15' BACnetServiceAckVTOpen

        ]
        ['0x17' BACnetServiceAckVTData

        ]

        ['0x18' BACnetServiceAckRemovedAuthenticate

        ]
        ['0x0D' BACnetServiceAckRemovedReadPropertyConditional

        ]
    ]
]

[discriminatedType 'BACnetConfirmedServiceACK'
    [discriminator uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x03' BACnetConfirmedServiceACKGetAlarmSummary
        ]
        ['0x04' BACnetConfirmedServiceACKGetEnrollmentSummary
        ]
        ['0x1D' BACnetConfirmedServiceACKGetEventInformation
        ]

        ['0x06' BACnetConfirmedServiceACKAtomicReadFile
        ]
        ['0x07' BACnetConfirmedServiceACKAtomicWriteFile
        ]

        ['0x0A' BACnetConfirmedServiceACKCreateObject
        ]
        ['0x0C' BACnetConfirmedServiceACKReadProperty
        ]
        ['0x0E' BACnetConfirmedServiceACKReadPropertyMultiple
        ]
        ['0x1A' BACnetConfirmedServiceACKReadRange
        ]

        ['0x12' BACnetConfirmedServiceACKConfirmedPrivateTransfer
        ]

        ['0x15' BACnetConfirmedServiceACKVTOpen
        ]
        ['0x17' BACnetConfirmedServiceACKVTData
        ]

        ['0x18' BACnetConfirmedServiceACKRemovedAuthenticate
        ]
        ['0x0D' BACnetConfirmedServiceACKRemovedReadPropertyConditional
        ]
    ]
]

[discriminatedType 'BACnetError'
    [discriminator uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x03' BACnetErrorGetAlarmSummary
        ]
        ['0x04' BACnetErrorGetEnrollmentSummary
        ]
        ['0x1D' BACnetErrorGetEventInformation
        ]

        ['0x06' BACnetErrorAtomicReadFile
        ]
        ['0x07' BACnetErrorAtomicWriteFile
        ]

        ['0x0A' BACnetErrorCreateObject
        ]
        ['0x0C' BACnetErrorReadProperty
            [const uint 5 'errorClassHeader' '0x12']
            [simple uint 3 'errorClassLength']
            [array int 8 'errorClass' count 'errorClassLength']
            [const uint 5 'errorCodeHeader' '0x12']
            [simple uint 3 'errorCodeLength']
            [array int 8 'errorCode' count 'errorCodeLength']
        ]
        ['0x0E' BACnetErrorReadPropertyMultiple
        ]
        ['0x1A' BACnetErrorReadRange
        ]

        ['0x12' BACnetErrorConfirmedPrivateTransfer
        ]

        ['0x15' BACnetErrorVTOpen
        ]
        ['0x17' BACnetErrorVTData
        ]

        ['0x18' BACnetErrorRemovedAuthenticate
        ]
        ['0x0D' BACnetErrorRemovedReadPropertyConditional
        ]
    ]
]

[type 'BACnetAddress'
    [array  uint 8 'address' count '4']
    [simple uint 16 'port']
]

[type 'BACnetTagWithContent'
    [simple        uint 4    'tagNumber'       ]
    [simple        TagClass  'tagClass'        ]
    [simple        uint 3    'lengthValueType' ]
    [optional      uint 8    'extTagNumber'        'tagNumber == 15'     ]
    [optional      uint 8    'extLength'           'lengthValueType == 5']
    [array         uint 8    'propertyIdentifier'  length          '(lengthValueType == 5) ? extLength : lengthValueType']
    [const         uint 8    'openTag'             '0x2e']
    [simple        BACnetTag 'value'                     ]
    [const         uint 8    'closingTag'          '0x2f']
]

[discriminatedType 'BACnetTag'
    [simple        uint 4   'tagNumber'                                                  ]
    [discriminator TagClass 'tagClass'                                                   ]
    [simple        uint 3   'lengthValueType'                                            ]
    [virtual       bit      'isPrimitiveAndNotBoolean' '!(tagClass == TagClass.CONTEXT_SPECIFIC_TAGS && lengthValueType == 6) && tagNumber != 1']
    [optional      uint 8   'extTagNumber'    'isPrimitiveAndNotBoolean && tagNumber == 15'                          ]
    [optional      uint 8   'extLength'       'isPrimitiveAndNotBoolean && lengthValueType == 5'                     ]
    [optional      uint 16  'extExtLength'    'isPrimitiveAndNotBoolean && lengthValueType == 5 && extLength == 254' ]
    [optional      uint 32  'extExtExtLength' 'isPrimitiveAndNotBoolean && lengthValueType == 5 && extLength == 255' ]
    [virtual       uint 32  'actualLength'    'lengthValueType == 5 && extLength == 255 ? extExtExtLength : (lengthValueType == 5 && extLength == 254 ? extExtLength : (lengthValueType == 5 ? extLength : (isPrimitiveAndNotBoolean ? lengthValueType : 0)))']
    [typeSwitch 'tagClass','tagNumber'
        ['APPLICATION_TAGS','0x0' BACnetTagApplicationNull
        ]
        ['APPLICATION_TAGS','0x1' BACnetTagApplicationBoolean
        ]
        ['APPLICATION_TAGS','0x2' BACnetTagApplicationUnsignedInteger [uint 3 'lengthValueType', uint 8 'extLength']
            [array int 8 'data' length '(lengthValueType == 5) ? extLength : lengthValueType']
        ]
        ['APPLICATION_TAGS','0x3' BACnetTagApplicationSignedInteger [uint 3 'lengthValueType', uint 8 'extLength']
            [array int 8 'data' length '(lengthValueType == 5) ? extLength : lengthValueType']
        ]
        ['APPLICATION_TAGS','0x4' BACnetTagApplicationReal [uint 3 'lengthValueType', uint 8 'extLength']
            [simple float 8.23 'value']
        ]
        ['APPLICATION_TAGS','0x5' BACnetTagApplicationDouble [uint 3 'lengthValueType', uint 8 'extLength']
            [simple float 11.52 'value']
        ]
        ['APPLICATION_TAGS','0x6' BACnetTagApplicationOctetString
        ]
        ['APPLICATION_TAGS','0x7' BACnetTagApplicationCharacterString
        ]
        ['APPLICATION_TAGS','0x8' BACnetTagApplicationBitString [uint 3 'lengthValueType', uint 8 'extLength']
            [simple uint 8 'unusedBits']
            [array int 8 'data' length '(lengthValueType == 5) ? (extLength - 1) : (lengthValueType - 1)']
        ]
        ['APPLICATION_TAGS','0x9' BACnetTagApplicationEnumerated [uint 3 'lengthValueType', uint 8 'extLength']
            [array int 8 'data' length '(lengthValueType == 5) ? extLength : lengthValueType']
        ]
        ['APPLICATION_TAGS','0xA' BACnetTagApplicationDate
        ]
        ['APPLICATION_TAGS','0xB' BACnetTagApplicationTime
        ]
        ['APPLICATION_TAGS','0xC' BACnetTagApplicationObjectIdentifier
        ]
        ['CONTEXT_SPECIFIC_TAGS' BACnetTagContext [uint 4 'tagNumber', uint 8 'extTagNumber', uint 3 'lengthValueType', uint 8 'extLength']
            [array int 8 'data' length '(lengthValueType == 5) ? extLength : lengthValueType']
        ]
    ]
]

[enum uint 1 'TagClass'
    ['0x0' APPLICATION_TAGS]
    ['0x1' CONTEXT_SPECIFIC_TAGS]
]

[enum int 4 'ApplicationTag'
    ['0x0' NULL]
    ['0x1' BOOLEAN]
    ['0x2' UNSIGNED_INTEGER]
    ['0x3' SIGNED_INTEGER]
    ['0x4' REAL]
    ['0x5' DOUBLE]
    ['0x6' OCTET_STRING]
    ['0x7' CHARACTER_STRING]
    ['0x8' BIT_STRING]
    ['0x9' ENUMERATED]
    ['0xA' DATE]
    ['0xB' TIME]
    ['0xC' BACNET_OBJECT_IDENTIFIER]
]

[enum uint 4 'BACnetNetworkType'
    ['0x0' ETHERNET]
    ['0x1' ARCNET]
    ['0x2' MSTP]
    ['0x3' PTP]
    ['0x4' LONTALK]
    ['0x5' IPV4]
    ['0x6' ZIGBEE]
    ['0x7' VIRTUAL]
    ['0x8' REMOVED_NON_BACNET]
    ['0x9' IPV6]
    ['0xA' SERIAL]
]

[enum uint 8 'BACnetNodeType'
    ['0x00' UNKNOWN]
    ['0x01' SYSTEM]
    ['0x02' NETWORK]
    ['0x03' DEVICE]
    ['0x04' ORGANIZATIONAL]
    ['0x05' AREA]
    ['0x06' EQUIPMENT]
    ['0x07' POINT]
    ['0x08' COLLECTION]
    ['0x09' PROPERTY]
    ['0x0A' FUNCTIONAL]
    ['0x0B' OTHER]
    ['0x0C' SUBSYSTEM]
    ['0x0D' BUILDING]
    ['0x0E' FLOOR]
    ['0x0F' SECTION]
    ['0x10' MODULE]
    ['0x11' TREE]
    ['0x12' MEMBER]
    ['0x13' PROTOCOL]
    ['0x14' ROOM]
    ['0x15' ZONE]
]

[enum uint 4 'BACnetNotifyType'
    ['0x0' ALARM]
    ['0x1' EVENT]
    ['0x2' ACK_NOTIFICATION]
]

[enum uint 10 'BACnetObjectType'
    ['0x000' ANALOG_INPUT]
    ['0x001' ANALOG_OUTPUT]
    ['0x002' ANALOG_VALUE]

    ['0x003' BINARY_INPUT]
    ['0x004' BINARY_OUTPUT]
    ['0x005' BINARY_VALUE]

    ['0x00D' MULTISTATE_INPUT]
    ['0x00E' MULTISTATE_OUTPUT]
    ['0x013' MULTISTATE_VALUE]

    ['0x011' SCHEDULE]

    ['0x008' DEVICE]

    ['0x006' CALENDAR]
    ['0x007' COMMAND]
    ['0x009' EVENT_ENROLLMENT]
    ['0x00A' FILE]
    ['0x00B' GROUP]
    ['0x00C' LOOP]
    ['0x00F' NOTIFICATION_CLASS]
    ['0x010' PROGRAM]
    ['0x012' AVERAGING]
    ['0x014' TREND_LOG]
    ['0x015' LIFE_SAFETY_POINT]
    ['0x016' LIFE_SAFETY_ZONE]
    ['0x017' ACCUMULATOR]
    ['0x018' PULSE_CONVERTER]
    ['0x019' EVENT_LOG]
    ['0x01A' GLOBAL_GROUP]
    ['0x01B' TREND_LOG_MULTIPLE]
    ['0x01C' LOAD_CONTROL]
    ['0x01D' STRUCTURED_VIEW]
    ['0x01E' ACCESS_DOOR]
    ['0x01F' TIMER]
    ['0x020' ACCESS_CREDENTIAL]
    ['0x021' ACCESS_POINT]
    ['0x022' ACCESS_RIGHTS]
    ['0x023' ACCESS_USER]
    ['0x024' ACCESS_ZONE]
    ['0x025' CREDENTIAL_DATA_INPUT]
    ['0x026' NETWORK_SECURITY]
    ['0x027' BITSTRING_VALUE]
    ['0x028' CHARACTERSTRING_VALUE]
    ['0x029' DATEPATTERN_VALUE]
    ['0x02A' DATE_VALUE]
    ['0x02B' DATETIMEPATTERN_VALUE]
    ['0x02C' DATETIME_VALUE]
    ['0x02D' INTEGER_VALUE]
    ['0x02E' LARGE_ANALOG_VALUE]
    ['0x02F' OCTETSTRING_VALUE]
    ['0x030' POSITIVE_INTEGER_VALUE]
    ['0x031' TIMEPATTERN_VALUE]
    ['0x032' TIME_VALUE]
    ['0x033' NOTIFICATION_FORWARDER]
    ['0x034' ALERT_ENROLLMENT]
    ['0x035' CHANNEL]
    ['0x036' LIGHTING_OUTPUT]
    ['0x037' BINARY_LIGHTING_OUTPUT]
    ['0x038' NETWORK_PORT]
    ['0x039' ELEVATOR_GROUP]
    ['0x03A' ESCALATOR]
]
