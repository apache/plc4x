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

[discriminatedType BVLC byteOrder='BIG_ENDIAN'
    [const         uint 8  bacnetType   0x81       ]
    [discriminator uint 8  bvlcFunction              ]
    [implicit      uint 16 bvlcLength 'lengthInBytes']
    [typeSwitch bvlcFunction
        ['0x00' BVLCResult
            [simple BVLCResultCode code]
        ]
        ['0x01' BVLCWideBroadcastDistributionTable
        ]
        ['0x02' BVLCReadBroadcastDistributionTable
        ]
        ['0x03' BVLCReadBroadcastDistributionTableAck
        ]
        ['0x04' BVLCForwardedNPDU(uint 16 bvlcLength)
            [array  uint 8  ip    count '4'         ]
            [simple uint 16 port                    ]
            [simple NPDU('bvlcLength - 10') npdu   ]
        ]
        ['0x05' BVLCRegisterForeignDevice
            [simple uint 16 ttl]
        ]
        ['0x06' BVLCReadForeignDeviceTable
        ]
        ['0x07' BVLCReadForeignDeviceTableAck
        ]
        ['0x08' BVLCDeleteForeignDeviceTableEntry
        ]
        ['0x09' BVLCDistributeBroadcastToNetwork(uint 16 bvlcLength)
            [simple NPDU('bvlcLength - 4') npdu]
        ]
        ['0x0A' BVLCOriginalUnicastNPDU(uint 16 bvlcLength)
            [simple NPDU('bvlcLength - 4') npdu]
        ]
        ['0x0B' BVLCOriginalBroadcastNPDU(uint 16 bvlcLength)
            [simple NPDU('bvlcLength - 4') npdu]
        ]
        ['0x0C' BVLCSecureBVLL
        ]
    ]
]

[type NPDU(uint 16 npduLength)
    [simple   uint 8        protocolVersionNumber]
    [simple   bit           messageTypeFieldPresent]
    [reserved uint 1        '0']
    [simple   bit           destinationSpecified]
    [reserved uint 1        '0']
    [simple   bit           sourceSpecified]
    [simple   bit           expectingReply]
    [simple   uint 2        networkPriority]
    [optional uint 16       destinationNetworkAddress 'destinationSpecified']
    [optional uint 8        destinationLength         'destinationSpecified']
    [array    uint 8        destinationAddress count  'destinationSpecified ? destinationLength : 0']
    [optional uint 16       sourceNetworkAddress      'sourceSpecified']
    [optional uint 8        sourceLength              'sourceSpecified']
    [array    uint 8        sourceAddress count       'sourceSpecified ? sourceLength : 0']
    [optional uint 8        hopCount                  'destinationSpecified']
    [optional NLM('npduLength - (2 + (sourceSpecified ? 3 + sourceLength : 0) + (destinationSpecified ? 3 + destinationLength: 0) + ((destinationSpecified || sourceSpecified) ? 1 : 0))')           nlm                       'messageTypeFieldPresent'  ]
    [optional APDU('npduLength - (2 + (sourceSpecified ? 3 + sourceLength : 0) + (destinationSpecified ? 3 + destinationLength: 0) + ((destinationSpecified || sourceSpecified) ? 1 : 0))')          apdu                      '!messageTypeFieldPresent' ]
]

[discriminatedType NLM(uint 16 apduLength)
    [discriminator uint 8  messageType]
    [optional      uint 16 vendorId '(messageType >= 128) && (messageType <= 255)']
    [typeSwitch messageType
        ['0x0' NLMWhoIsRouterToNetwork(uint 8 messageType)
            [array uint 16 destinationNetworkAddress length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
        ['0x1' NLMIAmRouterToNetwork(uint 8 messageType)
            [array uint 16 destinationNetworkAddress length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
    ]
]

[discriminatedType APDU(uint 16 apduLength)
    [discriminator uint 4 apduType]
    [typeSwitch apduType
        ['0x0' APDUConfirmedRequest
            [simple   bit    segmentedMessage                       ]
            [simple   bit    moreFollows                            ]
            [simple   bit    segmentedResponseAccepted              ]
            [reserved uint 2 '0'                                      ]
            [simple   uint 3 maxSegmentsAccepted                    ]
            [simple   uint 4 maxApduLengthAccepted                  ]
            [simple   uint 8 invokeId                               ]
            [optional uint 8 sequenceNumber       'segmentedMessage']
            [optional uint 8 proposedWindowSize   'segmentedMessage']
            [simple   BACnetConfirmedServiceRequest('apduLength - (3 + (segmentedMessage ? 2 : 0))') serviceRequest]
        ]
        ['0x1' APDUUnconfirmedRequest
            [reserved uint 4                          '0'             ]
            [simple   BACnetUnconfirmedServiceRequest('apduLength - 1') serviceRequest]
        ]
        ['0x2' APDUSimpleAck
            [reserved uint 4 '0'               ]
            [simple   uint 8 originalInvokeId]
            [simple   uint 8 serviceChoice   ]
        ]
        ['0x3' APDUComplexAck
            [simple   bit               segmentedMessage                     ]
            [simple   bit               moreFollows                          ]
            [reserved uint 2            '0'                                    ]
            [simple   uint 8            originalInvokeId                     ]
            [optional uint 8            sequenceNumber     'segmentedMessage']
            [optional uint 8            proposedWindowSize 'segmentedMessage']
            [simple   BACnetServiceAck  serviceAck                           ]
        ]
        ['0x4' APDUSegmentAck
            [reserved uint 2 '0x00'              ]
            [simple   bit    negativeAck       ]
            [simple   bit    server            ]
            [simple   uint 8 originalInvokeId  ]
            [simple   uint 8 sequenceNumber    ]
            [simple   uint 8 proposedWindowSize]
        ]
        ['0x5' APDUError
            [reserved uint 4      '0x00'            ]
            [simple   uint 8      originalInvokeId]
            [simple   BACnetError error           ]
        ]
        ['0x6' APDUReject
            [reserved uint 4 '0x00'            ]
            [simple   uint 8 originalInvokeId]
            [simple   uint 8 rejectReason    ]
        ]
        ['0x7' APDUAbort
            [reserved uint 3 '0x00'            ]
            [simple   bit    server          ]
            [simple   uint 8 originalInvokeId]
            [simple   uint 8 abortReason     ]
        ]
    ]
]

[discriminatedType BACnetConfirmedServiceRequest(uint 16 len)
    [discriminator uint 8 serviceChoice]
    [typeSwitch serviceChoice
        ['0x00' BACnetConfirmedServiceRequestAcknowledgeAlarm
        ]
        ['0x01' BACnetConfirmedServiceRequestConfirmedCOVNotification
            [const  uint 8               subscriberProcessIdentifierHeader         0x09                 ]
            [simple uint 8               subscriberProcessIdentifier                                      ]
            [const  uint 8               monitoredObjectIdentifierHeader           0x1C                 ]
            [simple uint 10              monitoredObjectType                                              ]
            [simple uint 22              monitoredObjectInstanceNumber                                    ]
            [const  uint 8               issueConfirmedNotificationsHeader         0x2C                 ]
            [simple uint 10              issueConfirmedNotificationsType                                  ]
            [simple uint 22              issueConfirmedNotificationsInstanceNumber                        ]
            [const  uint 5               lifetimeHeader                            0x07                 ]
            [simple uint 3               lifetimeLength                                                   ]
            [array  int  8               lifetimeSeconds                           count  'lifetimeLength']
            [const  uint 8               listOfValuesOpeningTag                    0x4E                 ]
            [array  BACnetTagWithContent notifications                             length 'len - 18'      ]
            [const  uint 8               listOfValuesClosingTag                    0x4F                 ]
        ]
        ['0x02' BACnetConfirmedServiceRequestConfirmedEventNotification
        ]

        ['0x04' BACnetConfirmedServiceRequestGetEnrollmentSummary
        ]
        ['0x05' BACnetConfirmedServiceRequestSubscribeCOV
            [const  uint 8  subscriberProcessIdentifierHeader   0x09                ]
            [simple uint 8  subscriberProcessIdentifier                               ]
            [const  uint 8  monitoredObjectIdentifierHeader     0x1C                ]
            [simple uint 10 monitoredObjectType                                       ]
            [simple uint 22 monitoredObjectInstanceNumber                             ]
            [const  uint 8  issueConfirmedNotificationsHeader   0x29                ]
            [const  uint 7  issueConfirmedNotificationsSkipBits 0x00                ]
            [simple bit     issueConfirmedNotifications                               ]
            [const  uint 5  lifetimeHeader                      0x07                ]
            [simple uint 3  lifetimeLength                                            ]
            [array  int 8   lifetimeSeconds                     count 'lifetimeLength']
        ]

        ['0x06' BACnetConfirmedServiceRequestAtomicReadFile
        ]
        ['0x07' BACnetConfirmedServiceRequestAtomicWriteFile
            [simple BACnetApplicationTagObjectIdentifier                deviceIdentifier    ]
            [optional BACnetContextTagNull('0', 'BACnetDataType.NULL')  openingTag          ]
            [simple BACnetApplicationTagSignedInteger                   fileStartPosition   ]
            [simple BACnetApplicationTagOctetString                     fileData            ]
            [optional BACnetContextTagNull('0', 'BACnetDataType.NULL')  closingTag          ]
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
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')     objectIdentifier                                                                    ]
            [simple   BACnetContextTagPropertyIdentifier('1', 'BACnetDataType.BACNET_PROPERTY_IDENTIFIER') propertyIdentifier                                                                  ]
            // TODO: check if this is the right identifier type and size
            [optional uint 32                                                                              arrayIndex 'propertyIdentifier.value == BACnetPropertyIdentifier.VALUE_SOURCE_ARRAY']
            // TODO: check if values are missing here?
        ]
        ['0x0E' BACnetConfirmedServiceRequestReadPropertyMultiple
        ]
        ['0x0F' BACnetConfirmedServiceRequestWriteProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')     objectIdentifier                                                                    ]
            [simple   BACnetContextTagPropertyIdentifier('1', 'BACnetDataType.BACNET_PROPERTY_IDENTIFIER') propertyIdentifier                                                                  ]
            // TODO: check if this is the right identifier type and size
            [optional uint 32 arrayIndex 'propertyIdentifier.value == BACnetPropertyIdentifier.VALUE_SOURCE_ARRAY']
             // TODO: possible revert again
            [simple   BACnetPropertyValue('propertyIdentifier.value')                                      propertyValue                                                                       ]
            [optional BACnetTag priority                  'curPos < (len - 1)'            ]
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
          [simple BACnetContextTagDeviceState('0', 'BACnetDataType.BACNET_DEVICE_STATE')     reinitializedStateOfDevice  ]
          [optional BACnetContextTagCharacterString('1', 'BACnetDataType.CHARACTER_STRING')  password                    ]
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

[discriminatedType BACnetUnconfirmedServiceRequest(uint 16 len)
    [discriminator uint 8 serviceChoice]
    [typeSwitch serviceChoice
        ['0x00' BACnetUnconfirmedServiceRequestIAm
            [simple BACnetApplicationTagObjectIdentifier    deviceIdentifier                ]
            [simple BACnetApplicationTagUnsignedInteger     maximumApduLengthAcceptedLength ]
            [simple BACnetApplicationTagEnumerated          segmentationSupported ] // TODO: map to enum
            [simple BACnetApplicationTagUnsignedInteger     vendorId ] // TODO: vendor list?
        ]
        ['0x01' BACnetUnconfirmedServiceRequestIHave
            [simple BACnetApplicationTagObjectIdentifier    deviceIdentifier    ]
            [simple BACnetApplicationTagObjectIdentifier    objectIdentifier    ]
            [simple BACnetApplicationTagCharacterString     objectName          ]
        ]
        ['0x02' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
        ]
        ['0x03' BACnetUnconfirmedServiceRequestUnconfirmedEventNotification
        ]
        ['0x04' BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
            [const uint 8 vendorIdHeader 0x09]
            [simple uint 8 vendorId]
            [const uint 8 serviceNumberHeader 0x1A]
            [simple uint 16 serviceNumber]
            [const uint 8 listOfValuesOpeningTag 0x2E]
            [array int 8 values length 'len - 8']
            [const uint 8 listOfValuesClosingTag 0x2F]
        ]
        ['0x05' BACnetUnconfirmedServiceRequestUnconfirmedTextMessage
        ]
        ['0x06' BACnetUnconfirmedServiceRequestTimeSynchronization
            [simple BACnetApplicationTagDate synchronizedDate]
            [simple BACnetApplicationTagTime synchronizedTime]
        ]
        ['0x07' BACnetUnconfirmedServiceRequestWhoHas
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeLowLimit                                         ]
            [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeHighLimit  'deviceInstanceRangeLowLimit != null' ]
            [optional BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  objectIdentifier                                                    ]
            [optional BACnetContextTagOctetString('3', 'BACnetDataType.OCTET_STRING')                   objectName                    'objectIdentifier == null'            ]
        ]
        ['0x08' BACnetUnconfirmedServiceRequestWhoIs
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')   deviceInstanceRangeLowLimit                                         ]
            [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')   deviceInstanceRangeHighLimit  'deviceInstanceRangeLowLimit != null' ]
        ]
        ['0x09' BACnetUnconfirmedServiceRequestUTCTimeSynchronization
        ]
        ['0x0A' BACnetUnconfirmedServiceRequestWriteGroup
        ]
        ['0x0B' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple
        ]
    ]
]

[discriminatedType BACnetServiceAck
    [discriminator   uint 8 serviceChoice]
    [typeSwitch serviceChoice
        ['0x03' BACnetServiceAckGetAlarmSummary

        ]
        ['0x04' BACnetServiceAckGetEnrollmentSummary

        ]
        ['0x1D' BACnetServiceAckGetEventInformation

        ]

        ['0x06' BACnetServiceAckAtomicReadFile

        ]
        ['0x07' BACnetServiceAckAtomicWriteFile
            [simple BACnetContextTagSignedInteger('0', 'BACnetDataType.SIGNED_INTEGER') fileStartPosition]
        ]

        ['0x0A' BACnetServiceAckCreateObject

        ]
        ['0x0C' BACnetServiceAckReadProperty
            [simple BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')     objectIdentifier   ]
            [simple BACnetContextTagPropertyIdentifier('1', 'BACnetDataType.BACNET_PROPERTY_IDENTIFIER') propertyIdentifier ]
            // TODO: check if this is the right identifier type and size
            [optional uint 32 arrayIndex 'propertyIdentifier.value == BACnetPropertyIdentifier.VALUE_SOURCE_ARRAY']
            [optional EnclosedTags('3', '7') enclosedTags]
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

[discriminatedType BACnetConfirmedServiceACK
    [discriminator uint 8 serviceChoice]
    [typeSwitch serviceChoice
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

[type EnclosedTags(uint 4 openingTagNumber, uint 4 closingTagNumber)
    [optional       BACnetContextTag('openingTagNumber', 'BACnetDataType.NULL') openingTag]
    [manualArray    BACnetTag data terminated 'STATIC_CALL("openingClosingTerminate", readBuffer, openingTag)' 'STATIC_CALL("parseTags", readBuffer)' 'STATIC_CALL("writeTags", writeBuffer, _value)' 'STATIC_CALL("tagsLength", data)']
    [optional       BACnetContextTag('closingTagNumber', 'BACnetDataType.NULL') closingTag]
]

[discriminatedType BACnetPropertyValue(BACnetPropertyIdentifier identifier)
    // TODO: this is a tag but there is currently no way to validate (e.h. validate expression)
    [const    uint 8    openingTag                0x3E                          ]
    [typeSwitch identifier
        ['OBJECT_TYPE'       BACnetPropertyValueObjectType
            [simple BACnetApplicationTagObjectIdentifier objectIdentifier]
        ]
        ['PRIORITY_ARRAY'       BACnetPropertyValuePriorityValue
            [array byte values count '16']
        ]
        ['PRESENT_VALUE'        BACnetPropertyValuePresentValue
            [simple BACnetTag value]
        ]
        ['RELINQUISH_DEFAULT'   BACnetPropertyValueRelinquishDefault
            [simple BACnetApplicationTagReal value]
        ]
    ]
    // TODO: this is a tag but there is currently no way to validate (e.h. validate expression)
    [const    uint 8    closingTag                0x3F                          ]
]

[discriminatedType BACnetError
    [discriminator uint 8 serviceChoice]
    [typeSwitch serviceChoice
        ['0x03' BACnetErrorGetAlarmSummary
        ]
        ['0x02' BACnetErrorConfirmedEventNotification
            [simple BACnetApplicationTagEnumerated errorClass]
            [simple BACnetApplicationTagEnumerated errorCode]
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
            [simple BACnetApplicationTagEnumerated errorClass]
            [simple BACnetApplicationTagEnumerated errorCode]
        ]
        ['0x0E' BACnetErrorReadPropertyMultiple
        ]
        ['0x0F' BACnetErrorWriteProperty
            [simple BACnetApplicationTagEnumerated errorClass]
            [simple BACnetApplicationTagEnumerated errorCode]
        ]
        ['0x1A' BACnetErrorReadRange
        ]

        ['0x12' BACnetErrorConfirmedPrivateTransfer
        ]
        ['0x14' BACnetErrorPasswordFailure
            [simple BACnetApplicationTagEnumerated errorClass]
            [simple BACnetApplicationTagEnumerated errorCode]
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

[type BACnetAddress
    [array  uint 8 address count '4']
    [simple uint 16 port]
]

[type BACnetTagWithContent
    [simple        uint 4    tagNumber       ]
    [simple        TagClass  tagClass        ]
    [simple        uint 3    lengthValueType ]
    [optional      uint 8    extTagNumber        'tagNumber == 15'     ]
    [optional      uint 8    extLength           'lengthValueType == 5']
    [array         uint 8    propertyIdentifier  length          '(lengthValueType == 5) ? extLength : lengthValueType']
    [const         uint 8    openTag             0x2e]
    [simple        BACnetTag value                     ]
    [const         uint 8    closingTag          0x2f]
]

[discriminatedType BACnetTag
    [simple        uint 4   tagNumber                                                   ]
    [discriminator TagClass tagClass                                                    ]
    [simple        uint 3   lengthValueType                                             ]
    [optional      uint 8   extTagNumber    'tagNumber == 15'                           ]
    [virtual       uint 8   actualTagNumber 'tagNumber < 15 ? tagNumber : extTagNumber' ]
    [virtual       bit      isPrimitiveAndNotBoolean '!(tagClass == TagClass.CONTEXT_SPECIFIC_TAGS && lengthValueType == 6) && tagNumber != 1']
    [optional      uint 8   extLength       'isPrimitiveAndNotBoolean && lengthValueType == 5'                     ]
    [optional      uint 16  extExtLength    'isPrimitiveAndNotBoolean && lengthValueType == 5 && extLength == 254' ]
    [optional      uint 32  extExtExtLength 'isPrimitiveAndNotBoolean && lengthValueType == 5 && extLength == 255' ]
    [virtual       uint 32  actualLength    'lengthValueType == 5 && extLength == 255 ? extExtExtLength : (lengthValueType == 5 && extLength == 254 ? extExtLength : (lengthValueType == 5 ? extLength : (isPrimitiveAndNotBoolean ? lengthValueType : 0)))']
    [typeSwitch tagClass, tagNumber
        ['APPLICATION_TAGS','0x0' BACnetApplicationTagNull
        ]
        ['APPLICATION_TAGS','0x1' BACnetApplicationTagBoolean(uint 32 actualLength)
            [virtual bit value   'actualLength == 1'    ]
            [virtual bit isTrue  'value'                ]
            [virtual bit isFalse '!value'               ]
        ]
        ['APPLICATION_TAGS','0x2' BACnetApplicationTagUnsignedInteger(uint 32 actualLength)
            [virtual    bit     isUint8     'actualLength == 1' ]
            [optional   uint  8 valueUint8  'isUint8'           ]
            [virtual    bit     isUint16    'actualLength == 2' ]
            [optional   uint 16 valueUint16 'isUint16'          ]
            [virtual    bit     isUint32    'actualLength == 3' ]
            [optional   uint 32 valueUint32 'isUint32'          ]
            // TODO: we only go up to uint32 till we have the BigInteger stuff in java solved
            [virtual    uint 32 actualValue 'isUint8?valueUint8:(isUint16?valueUint16:(isUint32?valueUint32:0))']
            /*
            [virtual    bit     isUint64    'actualLength == 4' ]
            [optional   uint 64 valueUint64 'isUint64'          ]
            [virtual    uint 64 actualValue 'isUint8?valueUint8:(isUint16?valueUint16:(isUint32?valueUint32:(isUint64?valueUint64:0)))']
            */
        ]
        ['APPLICATION_TAGS','0x3' BACnetApplicationTagSignedInteger(uint 32 actualLength)
            [virtual    bit     isInt8     'actualLength == 1'  ]
            [optional   int 8   valueInt8  'isInt8'             ]
            [virtual    bit     isInt16    'actualLength == 2'  ]
            [optional   int 16  valueInt16 'isInt16'            ]
            [virtual    bit     isInt32    'actualLength == 3'  ]
            [optional   int 32  valueInt32 'isInt32'            ]
            [virtual    bit     isInt64    'actualLength == 4'  ]
            [optional   int 64  valueInt64 'isInt64'            ]
            [virtual    uint 64 actualValue 'isInt8?valueInt8:(isInt16?valueInt16:(isInt64?valueInt64:0))']
        ]
        ['APPLICATION_TAGS','0x4' BACnetApplicationTagReal
            [simple float 32 value]
        ]
        ['APPLICATION_TAGS','0x5' BACnetApplicationTagDouble
            [simple float 64 value]
        ]
        ['APPLICATION_TAGS','0x6' BACnetApplicationTagOctetString(uint 32 actualLength)
            // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
            [virtual    uint     16                   actualLengthInBit 'actualLength * 8']
            [simple     vstring 'actualLengthInBit'  value encoding='"ASCII"']
        ]
        ['APPLICATION_TAGS','0x7' BACnetApplicationTagCharacterString(uint 32 actualLength)
            [simple     BACnetCharacterEncoding      encoding]
            // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
            [virtual    uint     16                  actualLengthInBit 'actualLength * 8 - 8']
            // TODO: call to string on encoding or add type conversion so we can use the enum above
            [simple     vstring 'actualLengthInBit'  value encoding='"UTF-8"']
        ]
        ['APPLICATION_TAGS','0x8' BACnetApplicationTagBitString(uint 32 actualLength)
            [simple uint 8 unusedBits]
            [array int 8 data length 'actualLength']
        ]
        ['APPLICATION_TAGS','0x9' BACnetApplicationTagEnumerated(uint 32 actualLength)
            [array int 8 data length 'actualLength']
        ]
        ['APPLICATION_TAGS','0xA' BACnetApplicationTagDate
            [virtual int  8 wildcard '0xFF']
            [simple  int  8 yearMinus1900]
            [virtual bit    yearIsWildcard 'yearMinus1900 == wildcard']
            [virtual int 16 year 'yearMinus1900 + 1900']
            [simple  int  8 month]
            [virtual bit    monthIsWildcard 'month == wildcard']
            [virtual bit    oddMonthWildcard 'month == 13']
            [virtual bit    evenMonthWildcard 'month == 14']
            [simple  int  8 dayOfMonth]
            [virtual bit    dayOfMonthIsWildcard 'dayOfMonth == wildcard']
            [virtual bit    lastDayOfMonthWildcard 'dayOfMonth == 32']
            [virtual bit    oddDayOfMonthWildcard 'dayOfMonth == 33']
            [virtual bit    evenDayOfMonthWildcard 'dayOfMonth == 34']
            [simple  int  8 dayOfWeek]
            [virtual bit    dayOfWeekIsWildcard 'dayOfWeek == wildcard']
        ]
        ['APPLICATION_TAGS','0xB' BACnetApplicationTagTime
            [virtual int  8 wildcard '0xFF']
            [simple  int  8 hour]
            [virtual bit    hourIsWildcard 'hour == wildcard']
            [simple  int  8 minute]
            [virtual bit    minuteIsWildcard 'minute == wildcard']
            [simple  int  8 second]
            [virtual bit    secondIsWildcard 'second == wildcard']
            [simple  int  8 fractional]
            [virtual bit    fractionalIsWildcard 'fractional == wildcard']
        ]
        ['APPLICATION_TAGS','0xC' BACnetApplicationTagObjectIdentifier
            [simple BACnetObjectType    objectType] // TODO: map to enum
            [simple uint 22             instanceNumber]
        ]
        ['CONTEXT_SPECIFIC_TAGS' BACnetContextTagWithoutContext
            // A Context tag here can be ignored as we are missing context here
        ]
    ]
]

[discriminatedType BACnetContextTag(uint 4 tagNumberArgument, BACnetDataType dataType)
    [assert        uint 4           tagNumber           'tagNumberArgument'                                           ]
    [assert        TagClass         tagClass            'TagClass.CONTEXT_SPECIFIC_TAGS'                              ]
    [simple        uint 3           lengthValueType                                                                   ]
    [optional      uint 8           extTagNumber        'tagNumber == 15'                                             ]
    [virtual       uint 8           actualTagNumber     'tagNumber < 15 ? tagNumber : extTagNumber'                   ]
    [optional      uint 8           extLength           'lengthValueType == 5'                                        ]
    [optional      uint 16          extExtLength        'lengthValueType == 5 && extLength == 254'                    ]
    [optional      uint 32          extExtExtLength     'lengthValueType == 5 && extLength == 255'                    ]
    [virtual       uint 32          actualLength        'lengthValueType == 5 && extLength == 255 ? extExtExtLength : (lengthValueType == 5 && extLength == 254 ? extExtLength : (lengthValueType == 5 ? extLength : lengthValueType))']
    [typeSwitch dataType
        ['NULL' BACnetContextTagNull
        ]
        ['BOOLEAN' BACnetContextTagBoolean(uint 32 actualLength)
            [virtual bit value   'actualLength == 1'    ]
            [virtual bit isTrue  'value'                ]
            [virtual bit isFalse '!value'               ]
        ]
        ['UNSIGNED_INTEGER' BACnetContextTagUnsignedInteger(uint 32 actualLength)
            [virtual    bit     isUint8     'actualLength == 1' ]
            [optional   uint  8 valueUint8  'isUint8'           ]
            [virtual    bit     isUint16    'actualLength == 2' ]
            [optional   uint 16 valueUint16 'isUint16'          ]
            [virtual    bit     isUint32    'actualLength == 3' ]
            [optional   uint 32 valueUint32 'isUint32'          ]
            // TODO: we only go up to uint32 till we have the BigInteger stuff in java solved
            [virtual    uint 32 actualValue 'isUint8?valueUint8:(isUint16?valueUint16:(isUint32?valueUint32:0))']
            /*
            [virtual    bit     isUint64    'actualLength == 4' ]
            [optional   uint 64 valueUint64 'isUint64'          ]
            [virtual    uint 64 actualValue 'isUint8?valueUint8:(isUint16?valueUint16:(isUint32?valueUint32:(isUint64?valueUint64:0)))']
            */
        ]
        ['SIGNED_INTEGER' BACnetContextTagSignedInteger(uint 32 actualLength)
            [virtual    bit     isInt8     'actualLength == 1'  ]
            [optional   int 8   valueInt8  'isInt8'             ]
            [virtual    bit     isInt16    'actualLength == 2'  ]
            [optional   int 16  valueInt16 'isInt16'            ]
            [virtual    bit     isInt32    'actualLength == 3'  ]
            [optional   int 32  valueInt32 'isInt32'            ]
            [virtual    bit     isInt64    'actualLength == 4'  ]
            [optional   int 64  valueInt64 'isInt64'            ]
            [virtual    uint 64 actualValue 'isInt8?valueInt8:(isInt16?valueInt16:(isInt64?valueInt64:0))']
        ]
        ['REAL' BACnetContextTagReal(uint 32 actualLength)
            [simple     float 32 value]
        ]
        ['DOUBLE' BACnetContextTagDouble(uint 32 actualLength)
            [simple     float 64 value]
        ]
        ['OCTET_STRING' BACnetContextTagOctetString(uint 32 actualLength)
            // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
            [virtual    uint     16                   actualLengthInBit 'actualLength * 8']
            [simple     vstring 'actualLengthInBit'  value encoding='"ASCII"']
        ]
        ['CHARACTER_STRING' BACnetContextTagCharacterString(uint 32 actualLength)
            [simple     BACnetCharacterEncoding      encoding]
            // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
            [virtual    uint     16                  actualLengthInBit 'actualLength * 8 - 8']
            // TODO: call to string on encoding or add type conversion so we can use the enum above
            [simple     vstring 'actualLengthInBit'  value encoding='"UTF-8"']
        ]
        ['BIT_STRING' BACnetContextTagBitString(uint 32 actualLength)
            // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
            [virtual    uint 16                   actualLengthInBit 'actualLength * 8']
            [simple     uint  8 unusedBits]
            [array      int   8 data length 'actualLengthInBit']
        ]
        ['ENUMERATED' BACnetContextTagEnumerated(uint 32 actualLength)
            // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
            [virtual    uint     16                   actualLengthInBit 'actualLength * 8']
            [array int 8 data length 'actualLengthInBit']
        ]
        ['DATE' BACnetContextTagDate
            [virtual int  8 wildcard '0xFF']
            [simple  int  8 yearMinus1900]
            [virtual bit    yearIsWildcard 'yearMinus1900 == wildcard']
            [simple  int  8 month]
            [virtual bit    monthIsWildcard 'month == wildcard']
            [virtual bit    oddMonthWildcard 'month == 13']
            [virtual bit    evenMonthWildcard 'month == 14']
            [simple  int  8 dayOfMonth]
            [virtual bit    dayOfMonthIsWildcard 'dayOfMonth == wildcard']
            [virtual bit    lastDayOfMonthWildcard 'dayOfMonth == 32']
            [virtual bit    oddDayOfMonthWildcard 'dayOfMonth == 33']
            [virtual bit    evenDayOfMonthWildcard 'dayOfMonth == 34']
            [simple  int  8 dayOfWeek]
            [virtual bit    dayOfWeekIsWildcard 'dayOfWeek == wildcard']
        ]
        ['TIME' BACnetContextTagTime
            [virtual int  8 wildcard '0xFF']
            [simple  int  8 hour]
            [virtual bit    hourIsWildcard 'hour == wildcard']
            [simple  int  8 minute]
            [virtual bit    minuteIsWildcard 'minute == wildcard']
            [simple  int  8 second]
            [virtual bit    secondIsWildcard 'second == wildcard']
            [simple  int  8 fractional]
            [virtual bit    fractionalIsWildcard 'fractional == wildcard']
        ]
        ['BACNET_OBJECT_IDENTIFIER' BACnetContextTagObjectIdentifier
            [simple BACnetObjectType    objectType      ]
            [simple uint 22             instanceNumber  ]
        ]
        ['BACNET_PROPERTY_IDENTIFIER' BACnetContextTagPropertyIdentifier(uint 32 actualLength)
            [manual  BACnetPropertyIdentifier   value   'STATIC_CALL("readPropertyIdentifier", readBuffer, actualLength)' 'STATIC_CALL("writePropertyIdentifier", writeBuffer, value)' '_value.actualLength']
            [manual  uint 32                    proprietaryValue   'STATIC_CALL("readProprietaryPropertyIdentifier", readBuffer, value, actualLength)' 'STATIC_CALL("writeProprietaryPropertyIdentifier", writeBuffer, value, proprietaryValue)' '_value.actualLength']
            [virtual bit                        isProprietary      'value == BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE']
        ]
        ['BACNET_DEVICE_STATE' BACnetContextTagDeviceState
            [simple BACnetDeviceState   state]
        ]
    ]
]

[enum uint 16 BVLCResultCode
    ['0x0000' SUCCESSFUL_COMPLETION]
    ['0x0010' WRITE_BROADCAST_DISTRIBUTION_TABLE_NAK]
    ['0x0020' READ_BROADCAST_DISTRIBUTION_TABLE_NAK]
    ['0x0030' REGISTER_FOREIGN_DEVICE_NAK]
    ['0x0040' READ_FOREIGN_DEVICE_TABLE_NAK]
    ['0x0050' DELETE_FOREIGN_DEVICE_TABLE_ENTRY_NAK]
    ['0x0060' DISTRIBUTE_BROADCAST_TO_NETWORK_NAK]
]

[enum uint 1 TagClass
    ['0x0' APPLICATION_TAGS]
    ['0x1' CONTEXT_SPECIFIC_TAGS]
]

[enum int 4 BACnetDataType
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
    ['0xD' BACNET_PROPERTY_IDENTIFIER]
    ['0xE' BACNET_DEVICE_STATE]
]

[enum byte BACnetCharacterEncoding
    ['0x0' ISO_10646] // UTF-8
    ['0x1' IBM_Microsoft_DBCS]
    ['0x2' JIS_X_0208]
    ['0x3' ISO_10646_4] // (UCS-4)
    ['0x4' ISO_10646_2] //(UCS-2)
    ['0x5' ISO_8859_1]
]

[enum uint 4 BACnetNetworkType
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

[enum uint 8 BACnetDeviceState
    ['0x0' COLDSTART]
    ['0x1' WARMSTART]
    ['0x2' ACTIVATE_CHANGES]
    ['0x3' STARTBACKUP]
    ['0x4' ENDBACKUP]
    ['0x5' STARTRESTORE]
    ['0x6' ENDRESTORE]
    ['0x7' ABORTRESTORE]
]

[enum uint 8 BACnetNodeType
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

[enum uint 4 BACnetNotifyType
    ['0x0' ALARM]
    ['0x1' EVENT]
    ['0x2' ACK_NOTIFICATION]
]

[enum uint 10 BACnetObjectType
    ['32' ACCESS_CREDENTIAL]
    ['30' ACCESS_DOOR]
    ['33' ACCESS_POINT]
    ['34' ACCESS_RIGHTS]
    ['35' ACCESS_USER]
    ['36' ACCESS_ZONE]
    ['23' ACCUMULATOR]
    ['52' ALERT_ENROLLMENT]
    ['0' ANALOG_INPUT]
    ['1' ANALOG_OUTPUT]
    ['2' ANALOG_VALUE]
    ['18' AVERAGING]
    ['3' BINARY_INPUT]
    ['55' BINARY_LIGHTING_OUTPUT]
    ['4' BINARY_OUTPUT]
    ['5' BINARY_VALUE]
    ['39' BITSTRING_VALUE]
    ['6' CALENDAR]
    ['53' CHANNEL]
    ['40' CHARACTERSTRING_VALUE]
    ['7' COMMAND]
    ['37' CREDENTIAL_DATA_INPUT]
    ['41' DATEPATTERN_VALUE]
    ['42' DATE_VALUE]
    ['43' DATETIMEPATTERN_VALUE]
    ['44' DATETIME_VALUE]
    ['8' DEVICE]
    ['57' ELEVATOR_GROUP]
    ['58' ESCALATOR]
    ['9' EVENT_ENROLLMENT]
    ['25' EVENT_LOG]
    ['10' FILE]
    ['26' GLOBAL_GROUP]
    ['11' GROUP]
    ['45' INTEGER_VALUE]
    ['46' LARGE_ANALOG_VALUE]
    ['21' LIFE_SAFETY_POINT]
    ['22' LIFE_SAFETY_ZONE]
    ['59' LIFT]
    ['54' LIGHTING_OUTPUT]
    ['28' LOAD_CONTROL]
    ['12' LOOP]
    ['13' MULTI_STATE_INPUT]
    ['14' MULTI_STATE_OUTPUT]
    ['19' MULTI_STATE_VALUE]
    ['56' NETWORK_PORT]
    ['38' NETWORK_SECURITY]
    ['15' NOTIFICATION_CLASS]
    ['51' NOTIFICATION_FORWARDER]
    ['47' OCTETSTRING_VALUE]
    ['48' POSITIVE_INTEGER_VALUE]
    ['16' PROGRAM]
    ['24' PULSE_CONVERTER]
    ['17' SCHEDULE]
    ['29' STRUCTURED_VIEW]
    ['49' TIMEPATTERN_VALUE]
    ['50' TIME_VALUE]
    ['31' TIMER]
    ['20' TREND_LOG]
    ['27' TREND_LOG_MULTIPLE]
]

[enum uint 32 BACnetPropertyIdentifier
    ['244' ABSENTEE_LIMIT]
    ['175' ACCEPTED_MODES]
    ['245' ACCESS_ALARM_EVENTS]
    ['246' ACCESS_DOORS]
    ['247' ACCESS_EVENT]
    ['248' ACCESS_EVENT_AUTHENTICATION_FACTOR]
    ['249' ACCESS_EVENT_CREDENTIAL]
    ['322' ACCESS_EVENT_TAG]
    ['250' ACCESS_EVENT_TIME]
    ['251' ACCESS_TRANSACTION_EVENTS]
    ['252' ACCOMPANIMENT]
    ['253' ACCOMPANIMENT_TIME]
    ['1'   ACK_REQUIRED]
    ['0'   ACKED_TRANSITIONS]
    ['2'   ACTION]
    ['3'   ACTION_TEXT]
    ['254' ACTIVATION_TIME]
    ['255' ACTIVE_AUTHENTICATION_POLICY]
    ['481' ACTIVE_COV_MULTIPLE_SUBSCRIPTIONS]
    ['152' ACTIVE_COV_SUBSCRIPTIONS]
    ['4'   ACTIVE_TEXT]
    ['5'   ACTIVE_VT_SESSIONS]
    ['212' ACTUAL_SHED_LEVEL]
    ['176' ADJUST_VALUE]
    ['6'   ALARM_VALUE]
    ['7'   ALARM_VALUES]
    ['193' ALIGN_INTERVALS]
    ['8'   ALL]
    ['9'   ALL_WRITES_SUCCESSFUL]
    ['365' ALLOW_GROUP_DELAY_INHIBIT]
    ['399' APDU_LENGTH]
    ['10'  APDU_SEGMENT_TIMEOUT]
    ['11'  APDU_TIMEOUT]
    ['12'  APPLICATION_SOFTWARE_VERSION]
    ['13'  ARCHIVE]
    ['256' ASSIGNED_ACCESS_RIGHTS]
    ['447' ASSIGNED_LANDING_CALLS]
    ['124' ATTEMPTED_SAMPLES]
    ['257' AUTHENTICATION_FACTORS]
    ['258' AUTHENTICATION_POLICY_LIST]
    ['259' AUTHENTICATION_POLICY_NAMES]
    ['260' AUTHENTICATION_STATUS]
    ['364' AUTHORIZATION_EXEMPTIONS]
    ['261' AUTHORIZATION_MODE]
    ['169' AUTO_SLAVE_DISCOVERY]
    ['125' AVERAGE_VALUE]
    ['338' BACKUP_AND_RESTORE_STATE]
    ['153' BACKUP_FAILURE_TIMEOUT]
    ['339' BACKUP_PREPARATION_TIME]
    ['407' BACNET_IP_GLOBAL_ADDRESS]
    ['408' BACNET_IP_MODE]
    ['409' BACNET_IP_MULTICAST_ADDRESS]
    ['410' BACNET_IP_NAT_TRAVERSAL]
    ['412' BACNET_IP_UDP_PORT]
    ['435' BACNET_IPV6_MODE]
    ['438' BACNET_IPV6_UDP_PORT]
    ['440' BACNET_IPV6_MULTICAST_ADDRESS]
    ['327' BASE_DEVICE_SECURITY_POLICY]
    ['413' BBMD_ACCEPT_FD_REGISTRATIONS]
    ['414' BBMD_BROADCAST_DISTRIBUTION_TABLE]
    ['415' BBMD_FOREIGN_DEVICE_TABLE]
    ['262' BELONGS_TO]
    ['14'  BIAS]
    ['342' BIT_MASK]
    ['343' BIT_TEXT]
    ['373' BLINK_WARN_ENABLE]
    ['126' BUFFER_SIZE]
    ['448' CAR_ASSIGNED_DIRECTION]
    ['449' CAR_DOOR_COMMAND]
    ['450' CAR_DOOR_STATUS]
    ['451' CAR_DOOR_TEXT]
    ['452' CAR_DOOR_ZONE]
    ['453' CAR_DRIVE_STATUS]
    ['454' CAR_LOAD]
    ['455' CAR_LOAD_UNITS]
    ['456' CAR_MODE]
    ['457' CAR_MOVING_DIRECTION]
    ['458' CAR_POSITION]
    ['15'  CHANGE_OF_STATE_COUNT]
    ['16'  CHANGE_OF_STATE_TIME]
    ['416' CHANGES_PENDING]
    ['366' CHANNEL_NUMBER]
    ['127' CLIENT_COV_INCREMENT]
    ['417' COMMAND]
    ['430' COMMAND_TIME_ARRAY]
    ['154' CONFIGURATION_FILES]
    ['367' CONTROL_GROUPS]
    ['19'  CONTROLLED_VARIABLE_REFERENCE]
    ['20'  CONTROLLED_VARIABLE_UNITS]
    ['21'  CONTROLLED_VARIABLE_VALUE]
    ['177' COUNT]
    ['178' COUNT_BEFORE_CHANGE]
    ['179' COUNT_CHANGE_TIME]
    ['22'  COV_INCREMENT]
    ['180' COV_PERIOD]
    ['128' COV_RESUBSCRIPTION_INTERVAL]
    ['349' COVU_PERIOD]
    ['350' COVU_RECIPIENTS]
    ['263' CREDENTIAL_DISABLE]
    ['264' CREDENTIAL_STATUS]
    ['265' CREDENTIALS]
    ['266' CREDENTIALS_IN_ZONE]
    ['431' CURRENT_COMMAND_PRIORITY]
    ['155' DATABASE_REVISION]
    ['23'  DATE_LIST]
    ['24'  DAYLIGHT_SAVINGS_STATUS]
    ['267' DAYS_REMAINING]
    ['25'  DEADBAND]
    ['374' DEFAULT_FADE_TIME]
    ['375' DEFAULT_RAMP_RATE]
    ['376' DEFAULT_STEP_INCREMENT]
    ['490' DEFAULT_SUBORDINATE_RELATIONSHIP]
    ['393' DEFAULT_TIMEOUT]
    ['484' DEPLOYED_PROFILE_LOCATION]
    ['26'  DERIVATIVE_CONSTANT]
    ['27'  DERIVATIVE_CONSTANT_UNITS]
    ['28'  DESCRIPTION]
    ['29'  DESCRIPTION_OF_HALT]
    ['30'  DEVICE_ADDRESS_BINDING]
    ['31'  DEVICE_TYPE]
    ['156' DIRECT_READING]
    ['328' DISTRIBUTION_KEY_REVISION]
    ['329' DO_NOT_HIDE]
    ['226' DOOR_ALARM_STATE]
    ['227' DOOR_EXTENDED_PULSE_TIME]
    ['228' DOOR_MEMBERS]
    ['229' DOOR_OPEN_TOO_LONG_TIME]
    ['230' DOOR_PULSE_TIME]
    ['231' DOOR_STATUS]
    ['232' DOOR_UNLOCK_DELAY_TIME]
    ['213' DUTY_WINDOW]
    ['32'  EFFECTIVE_PERIOD]
    ['386' EGRESS_ACTIVE]
    ['377' EGRESS_TIME]
    ['33'  ELAPSED_ACTIVE_TIME]
    ['459' ELEVATOR_GROUP]
    ['133' ENABLE]
    ['460' ENERGY_METER]
    ['461' ENERGY_METER_REF]
    ['268' ENTRY_POINTS]
    ['34'  ERROR_LIMIT]
    ['462' ESCALATOR_MODE]
    ['354' EVENT_ALGORITHM_INHIBIT]
    ['355' EVENT_ALGORITHM_INHIBIT_REF]
    ['353' EVENT_DETECTION_ENABLE]
    ['35'  EVENT_ENABLE]
    ['351' EVENT_MESSAGE_TEXTS]
    ['352' EVENT_MESSAGE_TEXTS_CONFIG]
    ['83'  EVENT_PARAMETERS]
    ['36'  EVENT_STATE]
    ['130' EVENT_TIME_STAMPS]
    ['37'  EVENT_TYPE]
    ['38'  EXCEPTION_SCHEDULE]
    ['368' EXECUTION_DELAY]
    ['269' EXIT_POINTS]
    ['214' EXPECTED_SHED_LEVEL]
    ['270' EXPIRATION_TIME]
    ['271' EXTENDED_TIME_ENABLE]
    ['272' FAILED_ATTEMPT_EVENTS]
    ['273' FAILED_ATTEMPTS]
    ['274' FAILED_ATTEMPTS_TIME]
    ['388' FAULT_HIGH_LIMIT]
    ['389' FAULT_LOW_LIMIT]
    ['358' FAULT_PARAMETERS]
    ['463' FAULT_SIGNALS]
    ['359' FAULT_TYPE]
    ['39'  FAULT_VALUES]
    ['418' FD_BBMD_ADDRESS]
    ['419' FD_SUBSCRIPTION_LIFETIME]
    ['40'  FEEDBACK_VALUE]
    ['41'  FILE_ACCESS_METHOD]
    ['42'  FILE_SIZE]
    ['43'  FILE_TYPE]
    ['44'  FIRMWARE_REVISION]
    ['464' FLOOR_TEXT]
    ['215' FULL_DUTY_BASELINE]
    ['323' GLOBAL_IDENTIFIER]
    ['465' GROUP_ID]
    ['346' GROUP_MEMBER_NAMES]
    ['345' GROUP_MEMBERS]
    ['467' GROUP_MODE]
    ['45'  HIGH_LIMIT]
    ['468' HIGHER_DECK]
    ['47'  IN_PROCESS]
    ['378' IN_PROGRESS]
    ['46'  INACTIVE_TEXT]
    ['394' INITIAL_TIMEOUT]
    ['181' INPUT_REFERENCE]
    ['469' INSTALLATION_ID]
    ['48'  INSTANCE_OF]
    ['379' INSTANTANEOUS_POWER]
    ['49'  INTEGRAL_CONSTANT]
    ['50'  INTEGRAL_CONSTANT_UNITS]
    ['387' INTERFACE_VALUE]
    ['195' INTERVAL_OFFSET]
    ['400' IP_ADDRESS]
    ['401' IP_DEFAULT_GATEWAY]
    ['402' IP_DHCP_ENABLE]
    ['403' IP_DHCP_LEASE_TIME]
    ['404' IP_DHCP_LEASE_TIME_REMAINING]
    ['405' IP_DHCP_SERVER]
    ['406' IP_DNS_SERVER]
    ['411' IP_SUBNET_MASK]
    ['436' IPV6_ADDRESS]
    ['442' IPV6_AUTO_ADDRESSING_ENABLE]
    ['439' IPV6_DEFAULT_GATEWAY]
    ['443' IPV6_DHCP_LEASE_TIME]
    ['444' IPV6_DHCP_LEASE_TIME_REMAINING]
    ['445' IPV6_DHCP_SERVER]
    ['441' IPV6_DNS_SERVER]
    ['437' IPV6_PREFIX_LENGTH]
    ['446' IPV6_ZONE_INDEX]
    ['344' IS_UTC]
    ['330' KEY_SETS]
    ['471' LANDING_CALL_CONTROL]
    ['470' LANDING_CALLS]
    ['472' LANDING_DOOR_STATUS]
    ['275' LAST_ACCESS_EVENT]
    ['276' LAST_ACCESS_POINT]
    ['432' LAST_COMMAND_TIME]
    ['277' LAST_CREDENTIAL_ADDED]
    ['278' LAST_CREDENTIAL_ADDED_TIME]
    ['279' LAST_CREDENTIAL_REMOVED]
    ['280' LAST_CREDENTIAL_REMOVED_TIME]
    ['331' LAST_KEY_SERVER]
    ['173' LAST_NOTIFY_RECORD]
    ['369' LAST_PRIORITY]
    ['196' LAST_RESTART_REASON]
    ['157' LAST_RESTORE_TIME]
    ['395' LAST_STATE_CHANGE]
    ['281' LAST_USE_TIME]
    ['166' LIFE_SAFETY_ALARM_VALUES]
    ['380' LIGHTING_COMMAND]
    ['381' LIGHTING_COMMAND_DEFAULT_PRIORITY]
    ['52'  LIMIT_ENABLE]
    ['182' LIMIT_MONITORING_INTERVAL]
    ['420' LINK_SPEED]
    ['422' LINK_SPEED_AUTONEGOTIATE]
    ['421' LINK_SPEEDS]
    ['53'  LIST_OF_GROUP_MEMBERS]
    ['54'  LIST_OF_OBJECT_PROPERTY_REFERENCES]
    ['56'  LOCAL_DATE]
    ['360' LOCAL_FORWARDING_ONLY]
    ['57'  LOCAL_TIME]
    ['58'  LOCATION]
    ['233' LOCK_STATUS]
    ['282' LOCKOUT]
    ['283' LOCKOUT_RELINQUISH_TIME]
    ['131' LOG_BUFFER]
    ['132' LOG_DEVICE_OBJECT_PROPERTY]
    ['134' LOG_INTERVAL]
    ['183' LOGGING_OBJECT]
    ['184' LOGGING_RECORD]
    ['197' LOGGING_TYPE]
    ['390' LOW_DIFF_LIMIT]
    ['59'  LOW_LIMIT]
    ['473' LOWER_DECK]
    ['423' MAC_ADDRESS]
    ['474' MACHINE_ROOM_ID]
    ['158' MAINTENANCE_REQUIRED]
    ['475' MAKING_CAR_CALL]
    ['60'  MANIPULATED_VARIABLE_REFERENCE]
    ['170' MANUAL_SLAVE_ADDRESS_BINDING]
    ['234' MASKED_ALARM_VALUES]
    ['382' MAX_ACTUAL_VALUE]
    ['62'  MAX_APDU_LENGTH_ACCEPTED]
    ['285' MAX_FAILED_ATTEMPTS]
    ['63'  MAX_INFO_FRAMES]
    ['64'  MAX_MASTER]
    ['65'  MAX_PRES_VALUE]
    ['167' MAX_SEGMENTS_ACCEPTED]
    ['61'  MAXIMUM_OUTPUT]
    ['135' MAXIMUM_VALUE]
    ['149' MAXIMUM_VALUE_TIMESTAMP]
    ['159' MEMBER_OF]
    ['347' MEMBER_STATUS_FLAGS]
    ['286' MEMBERS]
    ['383' MIN_ACTUAL_VALUE]
    ['69'  MIN_PRES_VALUE]
    ['66'  MINIMUM_OFF_TIME]
    ['67'  MINIMUM_ON_TIME]
    ['68'  MINIMUM_OUTPUT]
    ['136' MINIMUM_VALUE]
    ['150' MINIMUM_VALUE_TIMESTAMP]
    ['160' MODE]
    ['70'  MODEL_NAME]
    ['71'  MODIFICATION_DATE]
    ['287' MUSTER_POINT]
    ['288' NEGATIVE_ACCESS_RULES]
    ['332' NETWORK_ACCESS_SECURITY_POLICIES]
    ['424' NETWORK_INTERFACE_NAME]
    ['425' NETWORK_NUMBER]
    ['426' NETWORK_NUMBER_QUALITY]
    ['427' NETWORK_TYPE]
    ['476' NEXT_STOPPING_FLOOR]
    ['207' NODE_SUBTYPE]
    ['208' NODE_TYPE]
    ['17'  NOTIFICATION_CLASS]
    ['137' NOTIFICATION_THRESHOLD]
    ['72'  NOTIFY_TYPE]
    ['73'  NUMBER_OF_APDU_RETRIES]
    ['289' NUMBER_OF_AUTHENTICATION_POLICIES]
    ['74'  NUMBER_OF_STATES]
    ['75'  OBJECT_IDENTIFIER]
    ['76'  OBJECT_LIST]
    ['77'  OBJECT_NAME]
    ['78'  OBJECT_PROPERTY_REFERENCE]
    ['79'  OBJECT_TYPE]
    ['290' OCCUPANCY_COUNT]
    ['291' OCCUPANCY_COUNT_ADJUST]
    ['292' OCCUPANCY_COUNT_ENABLE]
    ['294' OCCUPANCY_LOWER_LIMIT]
    ['295' OCCUPANCY_LOWER_LIMIT_ENFORCED]
    ['296' OCCUPANCY_STATE]
    ['297' OCCUPANCY_UPPER_LIMIT]
    ['298' OCCUPANCY_UPPER_LIMIT_ENFORCED]
    ['477' OPERATION_DIRECTION]
    ['161' OPERATION_EXPECTED]
    ['80'  OPTIONAL]
    ['81'  OUT_OF_SERVICE]
    ['82'  OUTPUT_UNITS]
    ['333' PACKET_REORDER_TIME]
    ['300' PASSBACK_MODE]
    ['301' PASSBACK_TIMEOUT]
    ['478' PASSENGER_ALARM]
    ['84'  POLARITY]
    ['363' PORT_FILTER]
    ['302' POSITIVE_ACCESS_RULES]
    ['384' POWER]
    ['479' POWER_MODE]
    ['185' PRESCALE]
    ['85'  PRESENT_VALUE]
    ['86'  PRIORITY]
    ['87'  PRIORITY_ARRAY]
    ['88'  PRIORITY_FOR_WRITING]
    ['89'  PROCESS_IDENTIFIER]
    ['361' PROCESS_IDENTIFIER_FILTER]
    ['485' PROFILE_LOCATION]
    ['168' PROFILE_NAME]
    ['90'  PROGRAM_CHANGE]
    ['91'  PROGRAM_LOCATION]
    ['92'  PROGRAM_STATE]
    ['371' PROPERTY_LIST]
    ['93'  PROPORTIONAL_CONSTANT]
    ['94'  PROPORTIONAL_CONSTANT_UNITS]
    ['482' PROTOCOL_LEVEL]
    ['95'  PROTOCOL_CONFORMANCE_CLASS]
    ['96'  PROTOCOL_OBJECT_TYPES_SUPPORTED]
    ['139' PROTOCOL_REVISION]
    ['97'  PROTOCOL_SERVICES_SUPPORTED]
    ['98'  PROTOCOL_VERSION]
    ['186' PULSE_RATE]
    ['99'  READ_ONLY]
    ['303' REASON_FOR_DISABLE]
    ['100' REASON_FOR_HALT]
    ['102' RECIPIENT_LIST]
    ['141' RECORD_COUNT]
    ['140' RECORDS_SINCE_NOTIFICATION]
    ['483' REFERENCE_PORT]
    ['480' REGISTERED_CAR_CALL]
    ['103' RELIABILITY]
    ['357' RELIABILITY_EVALUATION_INHIBIT]
    ['104' RELINQUISH_DEFAULT]
    ['491' REPRESENTS]
    ['218' REQUESTED_SHED_LEVEL]
    ['348' REQUESTED_UPDATE_INTERVAL]
    ['105' REQUIRED]
    ['106' RESOLUTION]
    ['202' RESTART_NOTIFICATION_RECIPIENTS]
    ['340' RESTORE_COMPLETION_TIME]
    ['341' RESTORE_PREPARATION_TIME]
    ['428' ROUTING_TABLE]
    ['187' SCALE]
    ['188' SCALE_FACTOR]
    ['174' SCHEDULE_DEFAULT]
    ['235' SECURED_STATUS]
    ['334' SECURITY_PDU_TIMEOUT]
    ['335' SECURITY_TIME_WINDOW]
    ['107' SEGMENTATION_SUPPORTED]
    ['372' SERIAL_NUMBER]
    ['108' SETPOINT]
    ['109' SETPOINT_REFERENCE]
    ['162' SETTING]
    ['219' SHED_DURATION]
    ['220' SHED_LEVEL_DESCRIPTIONS]
    ['221' SHED_LEVELS]
    ['163' SILENCED]
    ['171' SLAVE_ADDRESS_BINDING]
    ['172' SLAVE_PROXY_ENABLE]
    ['142' START_TIME]
    ['396' STATE_CHANGE_VALUES]
    ['222' STATE_DESCRIPTION]
    ['110' STATE_TEXT]
    ['111' STATUS_FLAGS]
    ['143' STOP_TIME]
    ['144' STOP_WHEN_FULL]
    ['391' STRIKE_COUNT]
    ['209' STRUCTURED_OBJECT_LIST]
    ['210' SUBORDINATE_ANNOTATIONS]
    ['211' SUBORDINATE_LIST]
    ['487' SUBORDINATE_NODE_TYPES]
    ['489' SUBORDINATE_RELATIONSHIPS]
    ['488' SUBORDINATE_TAGS]
    ['362' SUBSCRIBED_RECIPIENTS]
    ['305' SUPPORTED_FORMAT_CLASSES]
    ['304' SUPPORTED_FORMATS]
    ['336' SUPPORTED_SECURITY_ALGORITHMS]
    ['112' SYSTEM_STATUS]
    ['486' TAGS]
    ['306' THREAT_AUTHORITY]
    ['307' THREAT_LEVEL]
    ['113' TIME_DELAY]
    ['356' TIME_DELAY_NORMAL]
    ['114' TIME_OF_ACTIVE_TIME_RESET]
    ['203' TIME_OF_DEVICE_RESTART]
    ['115' TIME_OF_STATE_COUNT_RESET]
    ['392' TIME_OF_STRIKE_COUNT_RESET]
    ['204' TIME_SYNCHRONIZATION_INTERVAL]
    ['116' TIME_SYNCHRONIZATION_RECIPIENTS]
    ['397' TIMER_RUNNING]
    ['398' TIMER_STATE]
    ['145' TOTAL_RECORD_COUNT]
    ['308' TRACE_FLAG]
    ['164' TRACKING_VALUE]
    ['309' TRANSACTION_NOTIFICATION_CLASS]
    ['385' TRANSITION]
    ['205' TRIGGER]
    ['117' UNITS]
    ['118' UPDATE_INTERVAL]
    ['337' UPDATE_KEY_SET_TIMEOUT]
    ['189' UPDATE_TIME]
    ['310' USER_EXTERNAL_IDENTIFIER]
    ['311' USER_INFORMATION_REFERENCE]
    ['317' USER_NAME]
    ['318' USER_TYPE]
    ['319' USES_REMAINING]
    ['119' UTC_OFFSET]
    ['206' UTC_TIME_SYNCHRONIZATION_RECIPIENTS]
    ['146' VALID_SAMPLES]
    ['190' VALUE_BEFORE_CHANGE]
    ['192' VALUE_CHANGE_TIME]
    ['191' VALUE_SET]
    ['433' VALUE_SOURCE]
    ['434' VALUE_SOURCE_ARRAY]
    ['151' VARIANCE_VALUE]
    ['120' VENDOR_IDENTIFIER]
    ['121' VENDOR_NAME]
    ['326' VERIFICATION_TIME]
    ['429' VIRTUAL_MAC_ADDRESS_TABLE]
    ['122' VT_CLASSES_SUPPORTED]
    ['123' WEEKLY_SCHEDULE]
    ['147' WINDOW_INTERVAL]
    ['148' WINDOW_SAMPLES]
    ['370' WRITE_STATUS]
    ['320' ZONE_FROM]
    ['165' ZONE_MEMBERS]
    ['321' ZONE_TO]
    ['9999' VENDOR_PROPRIETARY_VALUE]
]
