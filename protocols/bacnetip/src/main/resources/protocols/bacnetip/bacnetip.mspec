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

[discriminatedType 'BVLC'
    [const         uint 8 'bacnetType' '0x81']
    [discriminator uint 8 'bvlcFunction']
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
        ['0x04' BVLCForwardedNPDU
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
        ['0x0A' BVLCOriginalUnicastNPDU
        ]
        ['0x0B' BVLCOriginalBroadcastNPDU
        ]
        ['0x0C' BVLCSecureBVLL
        ]
    ]
]

[type 'NPDU'
    [simple   uint 8 'protocolVersionNumber']
    [simple   bit    'messageTypeFieldPresent']
    [reserved uint 1 '0']
    [simple   bit    'destinationSpecified']
    [reserved uint 1 '0']
    [simple   bit    'sourceSpecified']
    [reserved uint 1 '0']
    [simple   bit    'expectingReply']
    [simple   uint 2 'networkPriority']
    [optional uint 16 'destinationNetworkAddress' 'destinationSpecified']
    [optional uint 8  'destinationLength'         'destinationSpecified']
    [optional Address 'destinationAddress'        'destinationSpecified && (destinationLength > 0)']
    [optional uint 16 'sourceNetworkAddress'      'sourceSpecified']
    [optional uint 8  'sourceLength'              'sourceSpecified']
    [optional Address 'sourceAddress'             'sourceSpecified && (sourceLength > 0)']
    [optional uint 8  'hopCount'                  'destinationSpecified || sourceSpecified']
    [optional uint 8  'messageType'               'messageTypeFieldPresent']
    [optional uint 16 'vendorId'                  'messageTypeFieldPresent']
    [simple   APDU    'apdu']
]

[discriminatedType 'APDU'
    [discriminator uint 4 'apduType']
    [typeSwitch 'apduType'
        ['0x0' APDUConfirmedRequest
            [simple   bit    'segmentedMessage']
            [simple   bit    'moreFollows']
            [simple   bit    'segmentedResponseAccepted']
            [reserved uint 2 '0']
            [simple   uint 3 'maxSegmentsAccepted']
            [simple   uint 4 'maxApduLengthAccepted']
            [simple   uint 6 'invokeId']
            [optional uint 8 'sequenceNumber' 'segmentedMessage']
            [optional uint 8 'proposedWindowSize' 'segmentedMessage']
            [simple   BACnetConfirmedServiceRequest 'serviceRequest']
        ]
        ['0x1' APDUUnconfirmedRequest
            [reserved uint 4 '0']
            [simple   uint 8 'serviceChoice']
            [simple   BACnetUnconfirmedServiceRequest 'serviceRequest']
        ]
        ['0x2' APDUSimpleAck
            [reserved uint 4 '0']
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'serviceChoice']
        ]
        ['0x3' APDUComplexAck
            [simple   bit    'segmentedMessage']
            [simple   bit    'moreFollows']
            [reserved uint 2 '0']
            [simple   uint 8 'originalInvokeId']
            [optional uint 8 'sequenceNumber'     'segmentedMessage']
            [optional uint 8 'proposedWindowSize' 'segmentedMessage']
            [simple   uint 8 'serviceChoice']
            [simple   BACnetServiceAck 'serviceAck']
        ]
        ['0x4' APDUSegmentAck
            [reserved uint 2 '0x00']
            [simple   bit    'negativeAck']
            [simple   bit    'server']
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'sequenceNumber']
            [simple   uint 8 'proposedWindowSize']
        ]
        ['0x5' APDUError
            [reserved uint 4 '0x00']
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'errorChoice']
            [simple   BACnetError  'error']
        ]
        ['0x6' APDUReject
            [reserved uint 4 '0x00']
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'rejectReason']
        ]
        ['0x7' APDUAbort
            [reserved uint 3 '0x00']
            [simple   bit    'server']
            [simple   uint 8 'originalInvokeId']
            [simple   uint 8 'abortReason']
        ]
    ]
]

[discriminatedType 'BACnetConfirmedServiceRequest'
    [discriminator uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x00' BACnetConfirmedServiceRequestAcknowledgeAlarm
        ]
        ['0x01' BACnetConfirmedServiceRequestConfirmedCOVNotification
            subscriber-process-identifier (Unsigned32)
            initiating-device-identifier (BACnetObjectIdentifier)
            monitored-object-identifier (BACnetObjectIdentifier)
            time-remaining (Unsigned)
            list-of-values (SEQUENCE OF BACnetPropertyValue)
        ]
        ['0x1F' BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple
        ]
        ['0x02' BACnetConfirmedServiceRequestConfirmedEventNotification
        ]

        ['0x04' BACnetConfirmedServiceRequestGetEnrollmentSummary
        ]
        ['0x1D' BACnetConfirmedServiceRequestGetEventInformation
        ]
        ['0x1B' BACnetConfirmedServiceRequestLifeSafetyOperation
        ]
        ['0x05' BACnetConfirmedServiceRequestSubscribeCOV
        ]
        ['0x0C' BACnetConfirmedServiceRequestSubscribeCOVProperty
        ]
        ['0x0E' BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple
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
        ]
        ['0x0E' BACnetConfirmedServiceRequestReadPropertyMultiple
        ]
        ['0x1A' BACnetConfirmedServiceRequestReadRange
        ]
        ['0x0F' BACnetConfirmedServiceRequestWriteProperty
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

[discriminatedType 'BACnetUnconfirmedServiceRequest'
    [discriminator uint 8 'serviceChoice']
    [typeSwitch 'serviceChoice'
        ['0x00' BACnetUnconfirmedServiceRequestIAm
        ]
        ['0x01' BACnetUnconfirmedServiceRequestIHave
        ]
        ['0x02' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
        ]
        ['0x03' BACnetUnconfirmedServiceRequestUnconfirmedEventNotification
        ]
        ['0x04' BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
        ]
        ['0x05' BACnetUnconfirmedServiceRequestUnconfirmedTextMessage
        ]
        ['0x06' BACnetUnconfirmedServiceRequestTimeSynchronization
        ]
        ['0x07' BACnetUnconfirmedServiceRequestWhoHas
        ]
        ['0x08' BACnetUnconfirmedServiceRequestWhoIs
        ]
        ['0x09' BACnetUnconfirmedServiceRequestUTCTimeSynchronization
        ]
        ['0x0A' BACnetUnconfirmedServiceRequestWriteGroup
        ]
        ['0x0B' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple
        ]
    ]
]

[type 'BACnetServiceAck'
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

[type 'Address'
    [array uint 8 'address' count '4']
    [simple uint 16 'port']
]

[enum uint 4 'ApplicationTags'
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
