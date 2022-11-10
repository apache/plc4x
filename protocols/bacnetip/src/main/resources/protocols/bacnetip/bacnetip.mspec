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

[type BacnetConstants
    [const    uint 16     bacnetUdpDefaultPort 47808]
]

[discriminatedType BVLC byteOrder='BIG_ENDIAN'
    [const         uint 8   bacnetType   0x81                   ]
    [discriminator uint 8   bvlcFunction                        ]
    [implicit      uint 16  bvlcLength          'lengthInBytes' ]
    [virtual       uint 16  bvlcPayloadLength   'bvlcLength-4'  ]
    [typeSwitch bvlcFunction
        ['0x00' *Result
            [simple   BVLCResultCode
                            code                                ]
        ]
        ['0x01' *WriteBroadcastDistributionTable(uint 16 bvlcPayloadLength)
            [array BVLCBroadcastDistributionTableEntry
                            table
                                length 'bvlcPayloadLength'      ]
        ]
        ['0x02' *ReadBroadcastDistributionTable
        ]
        ['0x03' *ReadBroadcastDistributionTableAck(uint 16 bvlcPayloadLength)
            [array BVLCBroadcastDistributionTableEntry
                            table
                                length 'bvlcPayloadLength'      ]
        ]
        ['0x04' *ForwardedNPDU(uint 16 bvlcPayloadLength)
            [array    uint 8  ip    count '4'                     ]
            [simple   uint 16 port                                ]
            [simple   NPDU('bvlcPayloadLength - 6')
                            npdu                                ]
        ]
        ['0x05' *RegisterForeignDevice
            [simple   uint 16 ttl]
        ]
        ['0x06' *ReadForeignDeviceTable
        ]
        ['0x07' *ReadForeignDeviceTableAck(uint 16 bvlcPayloadLength)
            [array BVLCForeignDeviceTableEntry
                            table
                                length 'bvlcPayloadLength'      ]
        ]
        ['0x08' *DeleteForeignDeviceTableEntry
            [array    uint 8  ip  count '4'                       ]
            [simple   uint 16 port                                ]
        ]
        ['0x09' *DistributeBroadcastToNetwork(uint 16 bvlcPayloadLength)
            [simple   NPDU('bvlcPayloadLength')
                            npdu                                ]
        ]
        ['0x0A' *OriginalUnicastNPDU(uint 16 bvlcPayloadLength)
            [simple   NPDU('bvlcPayloadLength')
                            npdu                                ]
        ]
        ['0x0B' *OriginalBroadcastNPDU(uint 16 bvlcPayloadLength)
            [simple   NPDU('bvlcPayloadLength')
                            npdu                                ]
        ]
        ['0x0C' *SecureBVLL(uint 16 bvlcPayloadLength)
            [array byte     securityWrapper
                            length 'bvlcPayloadLength'          ]
        ]
    ]
]

[type BVLCBroadcastDistributionTableEntry
    [array    uint 8      ip                          count '4'   ]
    [simple   uint 16     port                                    ]
    [array    uint 8      broadcastDistributionMap    count '4'   ]
]

[type BVLCForeignDeviceTableEntry
    [array    uint 8      ip                          count '4'   ]
    [simple   uint 16     port                                    ]
    [simple   uint 16     ttl                                     ]
    [simple   uint 16     secondRemainingBeforePurge              ]
]

[type NPDU(uint 16 npduLength)
    [simple   uint 8      protocolVersionNumber                                                                   ]
    [simple   NPDUControl control                                                                                 ]
    [optional uint 16     destinationNetworkAddress   'control.destinationSpecified'                              ]
    [optional uint 8      destinationLength           'control.destinationSpecified'                              ]
    [array    uint 8      destinationAddress count    'control.destinationSpecified ? destinationLength : 0'      ]
                                                        // (destinationNetworkAddress(16bit) + destinationLength(8bit) + destinationLength)?
    [virtual  uint 16     destinationLengthAddon      'control.destinationSpecified ? (3 + destinationLength) : 0'  ]
    [optional uint 16     sourceNetworkAddress        'control.sourceSpecified'                                   ]
    [optional uint 8      sourceLength                'control.sourceSpecified'                                   ]
    [array    uint 8      sourceAddress count         'control.sourceSpecified ? sourceLength : 0'                ]
                                                        // (sourceNetworkAddress(16bit) + sourceLength(8bit) + sourceLength)?
    [virtual  uint 16     sourceLengthAddon           'control.sourceSpecified ? (3 + sourceLength) : 0'            ]
    [optional uint 8      hopCount                    'control.destinationSpecified'                              ]
                                                        // protocolVersionNumber(8bit) + control(8bit) + sourceLengthAddon + destinationLengthAddon + hopcount
    [virtual  uint 16     payloadSubtraction          '2 + (sourceLengthAddon + destinationLengthAddon + ((control.destinationSpecified) ? 1 : 0))'     ]
    [optional NLM('npduLength - payloadSubtraction')
                            nlm
                                                        'control.messageTypeFieldPresent'                           ]
    [optional APDU('npduLength - payloadSubtraction')
                            apdu
                                                        '!control.messageTypeFieldPresent'                          ]
    [validation    'nlm != null || apdu != null'        "something is wrong here... apdu and nlm not set"           ]
]

[type NPDUControl
    [simple   bit         messageTypeFieldPresent         ]
    [reserved   uint 1      '0'                           ]
    [simple   bit         destinationSpecified            ]
    [reserved   uint 1      '0'                           ]
    [simple   bit         sourceSpecified                 ]
    [simple   bit         expectingReply                  ]
    [simple   NPDUNetworkPriority
                            networkPriority               ]
]

[discriminatedType NLM(uint 16 apduLength)
    [discriminator uint 8   messageType                   ]
    [virtual       bit      isVendorProprietaryMessage 'messageType >= 128']
    [typeSwitch messageType, isVendorProprietaryMessage
        ['0x00' *WhoIsRouterToNetwork
            [array      uint 16     destinationNetworkAddress length 'apduLength - 1']
        ]
        ['0x01' *IAmRouterToNetwork
            [array      uint 16     destinationNetworkAddress length 'apduLength - 1']
        ]
        ['0x02' *ICouldBeRouterToNetwork
            [simple   uint 16     destinationNetworkAddress   ]
            [simple   uint 8      performanceIndex            ]
        ]
        ['0x03' *RejectRouterToNetwork
            [simple   NLMRejectRouterToNetworkRejectReason
                                    rejectReason              ]
            [simple   uint 16     destinationNetworkAddress   ]
        ]
        ['0x04' *RouterBusyToNetwork
            [array    uint 16     destinationNetworkAddress length 'apduLength - 1']
        ]
        ['0x05' *RouterAvailableToNetwork
            [array    uint 16     destinationNetworkAddress length 'apduLength - 1']
        ]
        ['0x06' *InitalizeRoutingTable
            [simple   uint 8      numberOfPorts               ]
            [array    NLMInitalizeRoutingTablePortMapping
                                    portMappings
                        count 'numberOfPorts'                 ]
        ]
        ['0x07' *InitalizeRoutingTableAck
            [simple   uint 8      numberOfPorts               ]
            [array    NLMInitalizeRoutingTablePortMapping
                                    portMappings
                        count 'numberOfPorts'                 ]
        ]
        ['0x08' *EstablishConnectionToNetwork
            [simple   uint 16     destinationNetworkAddress   ]
            [simple   uint 8      terminationTime             ]
        ]
        ['0x09' *DisconnectConnectionToNetwork
            [simple   uint 16     destinationNetworkAddress   ]
        ]
        ['0x0A' *ChallengeRequest
            [simple   byte        messageChallenge            ]
            [simple   uint 32     originalMessageId           ]
            [simple   uint 32     originalTimestamp           ]
        ]
        ['0x0B' *SecurityPayload
            [simple   uint 16     payloadLength               ]
            [array    byte        payload length 'payloadLength']
        ]
        ['0x0C' *SecurityResponse
            [simple   SecurityResponseCode
                                  responseCode                ]
            [simple   uint 32     originalMessageId           ]
            [simple   uint 32     originalTimestamp           ]
            // TODO: type out variable parameters
            [array    byte      variableParameters length 'apduLength-(1+1+4+4)'            ]
        ]
        ['0x0D' *RequestKeyUpdate
            [simple   byte      set1KeyRevision              ]
            [simple   uint 32   set1ActivationTime           ]
            [simple   uint 32   set1ExpirationTime           ]
            [simple   byte      set2KeyRevision              ]
            [simple   uint 32   set2ActivationTime           ]
            [simple   uint 32   set2ExpirationTime           ]
            [simple   byte      distributionKeyRevision      ]
        ]
        ['0x0E' *UpdateKeyUpdate
            [simple   NLMUpdateKeyUpdateControlFlags
                                controlFlags                 ]
            [optional byte      set1KeyRevision     'controlFlags.set1KeyRevisionActivationTimeExpirationTimePresent'   ]
            [optional uint 32   set1ActivationTime  'controlFlags.set1KeyRevisionActivationTimeExpirationTimePresent'   ]
            [optional uint 32   set1ExpirationTime  'controlFlags.set1KeyCountKeyParametersPresent'                     ]
            [optional uint 8    set1KeyCount        'controlFlags.set1KeyCountKeyParametersPresent'                     ]
            [array    NLMUpdateKeyUpdateKeyEntry
                                set1Keys count 'set1KeyCount!=null?set1KeyCount:0'                                      ]
            [optional byte      set2KeyRevision     'controlFlags.set1KeyRevisionActivationTimeExpirationTimePresent'   ]
            [optional uint 32   set2ActivationTime  'controlFlags.set1KeyRevisionActivationTimeExpirationTimePresent'   ]
            [optional uint 32   set2ExpirationTime  'controlFlags.set1KeyCountKeyParametersPresent'                     ]
            [optional uint 8    set2KeyCount        'controlFlags.set1KeyCountKeyParametersPresent'                     ]
            [array    NLMUpdateKeyUpdateKeyEntry
                                set2Keys count 'set1KeyCount!=null?set1KeyCount:0'                                      ]
        ]
        ['0x0F' *UpdateKeyDistributionKey
            [simple   byte      keyRevision                 ]
            [simple   NLMUpdateKeyUpdateKeyEntry
                                key                         ]
        ]
        ['0x10' *RequestMasterKey
            [simple   uint 8    numberOfSupportedKeyAlgorithms  ]
            [array    byte      encryptionAndSignatureAlgorithms
                                    length  'apduLength-2'      ] // TODO: type those
        ]
        ['0x11' *SetMasterKey
            [simple   NLMUpdateKeyUpdateKeyEntry
                                key                         ]
        ]
        ['0x12' *WhatIsNetworkNumber
            // No content
        ]
        ['0x13' *NetworkNumberIs
            [simple   uint 16   networkNumber               ]
            [reserved uint 7    '0'                         ]
            [simple   bit       networkNumberConfigured     ]
        ]
        [*,'false' *Reserved
            [array    byte      unknownBytes length '(apduLength>0)?(apduLength-1):0'       ]
        ]
        [* *VendorProprietaryMessage
            [simple   BACnetVendorId
                                vendorId]
            [array    byte      proprietaryMessage length '(apduLength>0)?(apduLength-3):0' ]
        ]
    ]
]

[type NLMInitalizeRoutingTablePortMapping
    [simple   uint 16     destinationNetworkAddress       ]
    [simple   uint 8      portId                          ]
    [simple   uint 8      portInfoLength                  ]
    [array    byte        portInfo count 'portInfoLength' ]
]

[type NLMUpdateKeyUpdateControlFlags
    [simple   bit       set1KeyRevisionActivationTimeExpirationTimePresent  ]
    [simple   bit       set1KeyCountKeyParametersPresent                    ]
    [simple   bit       set1ShouldBeCleared                                 ]
    [simple   bit       set2KeyRevisionActivationTimeExpirationTimePresent  ]
    [simple   bit       set2KeyCountKeyParametersPresent                    ]
    [simple   bit       set2ShouldBeCleared                                 ]
    [simple   bit       moreMessagesToBeExpected                            ]
    [simple   bit       removeAllKeys                                       ]
]

[type NLMUpdateKeyUpdateKeyEntry
    [simple   uint 16   keyIdentifier                   ]
    [simple   uint 8    keySize                         ]
    [array    byte      key length 'keySize'            ]
]

[discriminatedType APDU(uint 16 apduLength)
    [discriminator ApduType apduType]
    [typeSwitch apduType
        ['CONFIRMED_REQUEST_PDU' *ConfirmedRequest
            [simple   bit       segmentedMessage                         ]
            [simple   bit       moreFollows                              ]
            [simple   bit       segmentedResponseAccepted                ]
            [reserved uint 2    '0'                                      ]
            [simple   MaxSegmentsAccepted
                                maxSegmentsAccepted                      ]
            [simple   MaxApduLengthAccepted
                                maxApduLengthAccepted                    ]
            [simple   uint 8    invokeId                                 ]
            [optional uint 8    sequenceNumber       'segmentedMessage'  ]
            [optional uint 8    proposedWindowSize   'segmentedMessage'  ]
            [virtual  uint 16   apduHeaderReduction
                                    // apduType(4bit)+bits(3bit)+reserved(2bits)+maxSegmentsAccepted(3bit)+maxApduLengthAccepted(4bit)+originalInvokeId(8bit)+(sequenceNumber(8bit)+proposedWindowSize(8bit))?
                                    '3 + (segmentedMessage ? 2 : 0)'        ]
            [optional BACnetConfirmedServiceRequest('apduLength - apduHeaderReduction')
                                serviceRequest       '!segmentedMessage'    ]
            [validation '(!segmentedMessage && serviceRequest != null) || segmentedMessage'
                        "service request should be set"                     ]
            // When we read the first segment we want the service choice to be part of the bytes so we only read it > 0
            [optional BACnetConfirmedServiceChoice
                                segmentServiceChoice 'segmentedMessage && sequenceNumber != 0']
            [virtual  uint 16   segmentReduction
                                    '(segmentServiceChoice != null)?(apduHeaderReduction+1):apduHeaderReduction'       ]
            [array    byte      segment
                                    length
                                    'segmentedMessage?((apduLength>0)?(apduLength - segmentReduction):0):0'             ]
        ]
        ['UNCONFIRMED_REQUEST_PDU' *UnconfirmedRequest
            [reserved uint 4                          '0'               ]
            [simple   BACnetUnconfirmedServiceRequest('apduLength - 1')
                                serviceRequest                          ]
        ]
        ['SIMPLE_ACK_PDU' *SimpleAck
            [reserved uint 4    '0'                                     ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetConfirmedServiceChoice
                                serviceChoice                           ]
        ]
        ['COMPLEX_ACK_PDU' *ComplexAck
            [simple   bit       segmentedMessage                        ]
            [simple   bit       moreFollows                             ]
            [reserved uint 2    '0'                                     ]
            [simple   uint 8    originalInvokeId                        ]
            [optional uint 8    sequenceNumber     'segmentedMessage'   ]
            [optional uint 8    proposedWindowSize 'segmentedMessage'   ]
            [virtual  uint 16   apduHeaderReduction
                                    // apduType(4bit)+bits(2bit)+reserved(2bits)+originalInvokeId(8bit)+(sequenceNumber(8bit)+proposedWindowSize(8bit))?
                                    '2 + (segmentedMessage ? 2 : 0)'    ]
            [optional BACnetServiceAck('apduLength - apduHeaderReduction')
                                serviceAck         '!segmentedMessage'  ]
            [validation '(!segmentedMessage && serviceAck != null) || segmentedMessage'
                        "service ack should be set"                     ]
            // When we read the first segment we want the service choice to be part of the bytes so we only read it > 0
            [optional BACnetConfirmedServiceChoice
                                segmentServiceChoice 'segmentedMessage && sequenceNumber != 0']
            [virtual  uint 16   segmentReduction
                                    '(segmentServiceChoice != null)?(apduHeaderReduction+1):apduHeaderReduction'
                                                                        ]
            [array    byte      segment
                                    length
                                    'segmentedMessage?((apduLength>0)?(apduLength - segmentReduction):0):0'
                                                                        ]
        ]
        ['SEGMENT_ACK_PDU' *SegmentAck
            [reserved uint 2    '0x00'                                  ]
            [simple   bit       negativeAck                             ]
            [simple   bit       server                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   uint 8    sequenceNumber                          ]
            [simple   uint 8    proposedWindowSize                      ]
        ]
        ['ERROR_PDU' *Error
            [reserved uint 4    '0x00'                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetConfirmedServiceChoice
                                errorChoice                             ]
            [simple   BACnetError('errorChoice')
                                error                                   ]
        ]
        ['REJECT_PDU' *Reject
            [reserved uint 4    '0x00'                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetRejectReasonTagged('1')
                                rejectReason                            ]
        ]
        ['ABORT_PDU' *Abort
            [reserved uint 3    '0x00'                                  ]
            [simple   bit       server                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetAbortReasonTagged('1')
                                abortReason                             ]
        ]
        [* *Unknown
            [simple   uint 4    unknownTypeRest                         ]
            [array    byte      unknownBytes length '(apduLength>0)?apduLength:0'    ]
        ]
    ]
]

[enum uint 4 ApduType
  ['0x0' CONFIRMED_REQUEST_PDU           ]
  ['0x1' UNCONFIRMED_REQUEST_PDU         ]
  ['0x2' SIMPLE_ACK_PDU                  ]
  ['0x3' COMPLEX_ACK_PDU                 ]
  ['0x4' SEGMENT_ACK_PDU                 ]
  ['0x5' ERROR_PDU                       ]
  ['0x6' REJECT_PDU                      ]
  ['0x7' ABORT_PDU                       ]
  /////
  // plc4x definitions to not fall back to 0x0 in case one of those is parsed

  ['0x8' APDU_UNKNOWN_8                  ]
  ['0x9' APDU_UNKNOWN_9                  ]
  ['0xA' APDU_UNKNOWN_A                  ]
  ['0xB' APDU_UNKNOWN_B                  ]
  ['0xC' APDU_UNKNOWN_C                  ]
  ['0xD' APDU_UNKNOWN_D                  ]
  ['0xE' APDU_UNKNOWN_E                  ]
  ['0xF' APDU_UNKNOWN_F                  ]
  //
  /////
]

// Not really tagged as it has no header but is consistent with naming schema enum+Tagged
[type BACnetRejectReasonTagged(uint 32 actualLength)
    [manual   BACnetRejectReason
                    value
                        'STATIC_CALL("readEnumGeneric", readBuffer, actualLength, BACnetRejectReason.VENDOR_PROPRIETARY_VALUE)'
                        'STATIC_CALL("writeEnumGeneric", writeBuffer, value)'
                        '_value.isProprietary?0:(actualLength * 8)'                     ]
    [virtual  bit   isProprietary
                        'value == BACnetRejectReason.VENDOR_PROPRIETARY_VALUE'          ]
    [manual   uint 32
                    proprietaryValue
                        'STATIC_CALL("readProprietaryEnumGeneric", readBuffer, actualLength, isProprietary)'
                        'STATIC_CALL("writeProprietaryEnumGeneric", writeBuffer, proprietaryValue, isProprietary)'
                        '_value.isProprietary?(actualLength * 8):0'                     ]
]

// Not really tagged as it has no header but is consistent with naming schema enum+Tagged
[type BACnetAbortReasonTagged(uint 32 actualLength)
    [manual   BACnetAbortReason
                    value
                        'STATIC_CALL("readEnumGeneric", readBuffer, actualLength, BACnetAbortReason.VENDOR_PROPRIETARY_VALUE)'
                        'STATIC_CALL("writeEnumGeneric", writeBuffer, value)'
                        '_value.isProprietary?0:(actualLength * 8)'                     ]
    [virtual  bit   isProprietary
                        'value == BACnetAbortReason.VENDOR_PROPRIETARY_VALUE'           ]
    [manual   uint 32
                    proprietaryValue
                        'STATIC_CALL("readProprietaryEnumGeneric", readBuffer, actualLength, isProprietary)'
                        'STATIC_CALL("writeProprietaryEnumGeneric", writeBuffer, proprietaryValue, isProprietary)'
                        '_value.isProprietary?(actualLength * 8):0'                     ]
]

[discriminatedType BACnetConfirmedServiceRequest(uint 32 serviceRequestLength)
    [discriminator BACnetConfirmedServiceChoice serviceChoice]
    // we substract serviceChoice from our payload
    [virtual       uint 32  serviceRequestPayloadLength '(serviceRequestLength>0)?(serviceRequestLength - 1):0'    ]
    [typeSwitch serviceChoice
        ////
        // Alarm and Event Services

        ['ACKNOWLEDGE_ALARM' *AcknowledgeAlarm
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           acknowledgingProcessIdentifier ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  eventObjectIdentifier          ]
            [simple   BACnetEventStateTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')                     eventStateAcknowledged         ]
            [simple   BACnetTimeStampEnclosed('3')                                                      timestamp                      ]
            [simple   BACnetContextTagCharacterString('4', 'BACnetDataType.CHARACTER_STRING')           acknowledgmentSource           ]
            [simple   BACnetTimeStampEnclosed('5')                                                      timeOfAcknowledgment           ]
        ]
        ['CONFIRMED_COV_NOTIFICATION' *ConfirmedCOVNotification
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           subscriberProcessIdentifier    ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier     ]
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  monitoredObjectIdentifier      ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')           lifetimeInSeconds              ]
            [simple   BACnetPropertyValues('4', 'monitoredObjectIdentifier.objectType')                 listOfValues                   ]
        ]
        ['CONFIRMED_COV_NOTIFICATION_MULTIPLE' *ConfirmedCOVNotificationMultiple
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           subscriberProcessIdentifier    ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier     ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')           timeRemaining                  ]
            [optional BACnetTimeStampEnclosed('3')                                                      timestamp                      ]
            [simple   ListOfCovNotificationsList('4')                                                   listOfCovNotifications         ]
        ]
        ['CONFIRMED_EVENT_NOTIFICATION' *ConfirmedEventNotification
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           processIdentifier              ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier     ]
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  eventObjectIdentifier          ]
            [simple   BACnetTimeStampEnclosed('3')                                                      timestamp                      ]
            [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')           notificationClass              ]
            [simple   BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')           priority                       ]
            [simple   BACnetEventTypeTagged('6', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      eventType                      ]
            [optional BACnetContextTagCharacterString('7', 'BACnetDataType.CHARACTER_STRING')           messageText                    ]
            [simple   BACnetNotifyTypeTagged('8', 'TagClass.CONTEXT_SPECIFIC_TAGS')                     notifyType                     ]
            [optional BACnetContextTagBoolean('9', 'BACnetDataType.BOOLEAN')                            ackRequired                    ]
            [optional BACnetEventStateTagged('10', 'TagClass.CONTEXT_SPECIFIC_TAGS')                    fromState                      ]
            [simple   BACnetEventStateTagged('11', 'TagClass.CONTEXT_SPECIFIC_TAGS')                    toState                        ]
            [optional BACnetNotificationParameters('12', 'eventObjectIdentifier.objectType')            eventValues                    ]
        ]
        ['GET_ENROLLMENT_SUMMARY' *GetEnrollmentSummary
            [simple   BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                        acknowledgmentFilter           ]
            [optional BACnetRecipientProcessEnclosed('1')                                               enrollmentFilter               ]
            [optional BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                        eventStateFilter               ]
            [optional BACnetEventTypeTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      eventTypeFilter                ]
            [optional BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter('4')              priorityFilter                 ]
            [optional BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')           notificationClassFilter        ]
        ]
        ['GET_EVENT_INFORMATION' *GetEventInformation
            [optional BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  lastReceivedObjectIdentifier   ]
        ]
        ['LIFE_SAFETY_OPERATION' *LifeSafetyOperation
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           requestingProcessIdentifier    ]
            [simple   BACnetContextTagCharacterString('1', 'BACnetDataType.CHARACTER_STRING')           requestingSource               ]
            [simple   BACnetLifeSafetyOperationTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')            request                        ]
            [optional BACnetContextTagObjectIdentifier('3', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  objectIdentifier               ]
        ]
        ['SUBSCRIBE_COV' *SubscribeCOV
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')            subscriberProcessIdentifier  ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')   monitoredObjectIdentifier    ]
            [optional BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')                             issueConfirmed               ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')            lifetimeInSeconds            ]
        ]
        ['SUBSCRIBE_COV_PROPERTY' *SubscribeCOVProperty
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')            subscriberProcessIdentifier  ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')   monitoredObjectIdentifier    ]
            [optional BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')                             issueConfirmedNotifications  ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')            lifetime                     ]
            [simple   BACnetPropertyReferenceEnclosed('4')                                               monitoredPropertyIdentifier  ]
            [optional BACnetContextTagReal('5', 'BACnetDataType.REAL')                                   covIncrement                 ]
        ]
        ['SUBSCRIBE_COV_PROPERTY_MULTIPLE' *SubscribeCOVPropertyMultiple
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')            subscriberProcessIdentifier  ]
            [optional BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')                             issueConfirmedNotifications  ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')            lifetime                     ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')            maxNotificationDelay         ]
            [simple   BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleListOfCovSubscriptionSpecificationsList('4')
                                                                                                         listOfCovSubscriptionSpecifications ]
        ]
        //
        ////

        ////
        // File Access Services

        ['ATOMIC_READ_FILE' *AtomicReadFile
            [simple   BACnetApplicationTagObjectIdentifier                                               fileIdentifier               ]
            [simple   BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord                          accessMethod                 ]
        ]
        ['ATOMIC_WRITE_FILE' *AtomicWriteFile
            [simple   BACnetApplicationTagObjectIdentifier                                               deviceIdentifier             ]
            [optional BACnetOpeningTag('0')                                                              openingTag                   ]
            [simple   BACnetApplicationTagSignedInteger                                                  fileStartPosition            ]
            [simple   BACnetApplicationTagOctetString                                                    fileData                     ]
            [optional BACnetClosingTag('0')                                                              closingTag                   ]
        ]
        //
        ////

        ////
        // Object Access Services
        ['ADD_LIST_ELEMENT' *AddListElement
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [optional BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                                            listOfElements              ]
        ]
        ['REMOVE_LIST_ELEMENT' *RemoveListElement
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [optional BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                                            listOfElements              ]
        ]
        ['CREATE_OBJECT' *CreateObject
            [simple   BACnetConfirmedServiceRequestCreateObjectObjectSpecifier('0')                         objectSpecifier             ]
            [optional BACnetPropertyValues('1', 'objectSpecifier.isObjectType?objectSpecifier.objectType:objectSpecifier.objectIdentifier.objectType')
                                                                                                            listOfValues                ]
        ]
        ['DELETE_OBJECT' *DeleteObject
            [simple   BACnetApplicationTagObjectIdentifier                                                  objectIdentifier            ]
        ]
        ['READ_PROPERTY' *ReadProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
        ]
        ['READ_PROPERTY_MULTIPLE' *ReadPropertyMultiple(uint 32 serviceRequestPayloadLength)
            [array    BACnetReadAccessSpecification                                                         data
                            length 'serviceRequestPayloadLength'                                                                        ]
        ]
        ['READ_RANGE' *ReadRange
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               propertyArrayIndex          ]
            // TODO: this attribute should be named range but this is a keyword in golang (so at this point we should build a language translator which makes keywords safe)
            [optional BACnetConfirmedServiceRequestReadRangeRange                                           readRange                   ]
        ]
        ['WRITE_PROPERTY' *WriteProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [simple   BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                                            propertyValue               ]
            [optional BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')               priority                    ]
        ]
        ['WRITE_PROPERTY_MULTIPLE' *WritePropertyMultiple(uint 32 serviceRequestPayloadLength)
            [array    BACnetWriteAccessSpecification                                                        data
                            length 'serviceRequestPayloadLength'                                                                        ]
        ]
        //
        ////

        ////
        // Remote Device Management Services

        ['DEVICE_COMMUNICATION_CONTROL' *DeviceCommunicationControl
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')               timeDuration                ]
            [simple   BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisableTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            enableDisable               ]
            [optional BACnetContextTagCharacterString('2', 'BACnetDataType.CHARACTER_STRING')               password                    ]

        ]
        ['CONFIRMED_PRIVATE_TRANSFER' *ConfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                           vendorId                    ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')               serviceNumber               ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                                            serviceParameters           ]
        ]
        ['CONFIRMED_TEXT_MESSAGE' *ConfirmedTextMessage
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      textMessageSourceDevice     ]
            [optional BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass('1')                    messageClass                ]
            [simple   BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            messagePriority             ]
            [simple   BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')               message                     ]
        ]
        ['REINITIALIZE_DEVICE' *ReinitializeDevice
            [simple   BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            reinitializedStateOfDevice  ]
            [optional BACnetContextTagCharacterString('1', 'BACnetDataType.CHARACTER_STRING')
                                                                                                            password                    ]
        ]

        ////
        //  Virtual Terminal Services

        ['VT_OPEN' *VTOpen
            [simple   BACnetVTClassTagged('0', 'TagClass.APPLICATION_TAGS')                                 vtClass                     ]
            [simple   BACnetApplicationTagUnsignedInteger                                                   localVtSessionIdentifier    ]
        ]
        ['VT_CLOSE' *VTClose(uint 32 serviceRequestPayloadLength)
            [array    BACnetApplicationTagUnsignedInteger                                                   listOfRemoteVtSessionIdentifiers
                                                               length 'serviceRequestPayloadLength'                                     ]
        ]
        ['VT_DATA' *VTData
            [simple   BACnetApplicationTagUnsignedInteger                                                   vtSessionIdentifier         ]
            [simple   BACnetApplicationTagOctetString                                                       vtNewData                   ]
            [simple   BACnetApplicationTagUnsignedInteger                                                   vtDataFlag                  ]
        ]
        //
        ////

        ////
        //  Removed Services

        ['AUTHENTICATE' *Authenticate(uint 32 serviceRequestPayloadLength)
            [array    byte                                                                                  bytesOfRemovedService
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
        ['REQUEST_KEY' *RequestKey(uint 32 serviceRequestPayloadLength)
            [array    byte                                                                                  bytesOfRemovedService
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
        ['READ_PROPERTY_CONDITIONAL' *ReadPropertyConditional(uint 32 serviceRequestPayloadLength)
            [array    byte                                                                                  bytesOfRemovedService
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
        //
        ////

        [* *Unknown(uint 32 serviceRequestPayloadLength)
            [array    byte                                                                                  unknownBytes
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
    ]
]

[type BACnetConfirmedServiceRequestCreateObjectObjectSpecifier(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                     openingTag                                                                         ]
    [optional   BACnetContextTagEnumerated('0', 'BACnetDataType.ENUMERATED')
                     rawObjectType                                                                      ]
    [virtual    bit  isObjectType   'rawObjectType != null'                                             ]
    [virtual    BACnetObjectType
                     objectType     'STATIC_CALL("mapBACnetObjectType", rawObjectType)'                 ]
    [optional   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                     objectIdentifier                                                                   ]
    [virtual    bit  isObjectIdentifier   'objectIdentifier != null'                                    ]
    [validation 'isObjectType || isObjectIdentifier' "either we need a objectType or a objectIdentifier"]
    [simple   BACnetClosingTag('tagNumber')
                     closingTag                                                                         ]
]

[type ListOfCovNotificationsList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                     openingTag                                                                         ]
    [array    ListOfCovNotifications
                     specifications
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                        closingTag                                                                      ]
]

[type ListOfCovNotifications
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        monitoredObjectIdentifier                                                       ]
    [simple   BACnetOpeningTag('1')
                        openingTag                                                                      ]
    [array      ListOfCovNotificationsValue('monitoredObjectIdentifier.objectType')
                        listOfValues
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)'    ]
    [simple   BACnetClosingTag('1')
                        closingTag                                                                      ]
]

[type ListOfCovNotificationsValue(BACnetObjectType objectTypeArgument)
    [simple   BACnetPropertyIdentifierTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                propertyIdentifier                                                      ]
    [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                arrayIndex                                                              ]
    [simple   BACnetConstructedData('2', 'objectTypeArgument', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                propertyValue                                                           ]
    [optional BACnetContextTagTime('3', 'BACnetDataType.TIME')
                                timeOfChange                                                            ]
]

[type BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleListOfCovSubscriptionSpecificationsList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                     openingTag                                                                         ]
    [array    BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleListOfCovSubscriptionSpecifications
                     specifications
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                     closingTag                                                                         ]
]

[type BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleListOfCovSubscriptionSpecifications
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    monitoredObjectIdentifier                                                           ]
    [simple   BACnetOpeningTag('1')
                    openingTag                                                                          ]
    [array    BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleListOfCovSubscriptionSpecificationsReference
                    listOfCovReferences
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)'        ]
    [simple   BACnetClosingTag('1')
                    closingTag                                                                          ]
]

[type BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleListOfCovSubscriptionSpecificationsReference
    [simple   BACnetPropertyReferenceEnclosed('1')
                    monitoredProperty                                                                   ]
    [optional BACnetContextTagReal('1', 'BACnetDataType.REAL')
                    covIncrement                                                                        ]
    [simple   BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')
                    timestamped                                                                         ]
]

[type BACnetReadAccessSpecification
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    objectIdentifier                                                                    ]
    [simple   BACnetOpeningTag('1')
                     openingTag                                                                         ]
    [array    BACnetPropertyReference
                    listOfPropertyReferences
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)'        ]
    [simple   BACnetClosingTag('1')
                     closingTag                                                                         ]
]

[type BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                    minPriority                 ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                    maxPriority                 ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type BACnetRecipientProcessEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [simple   BACnetRecipientProcess
                    recipientProcess            ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type BACnetRecipientProcess
    [simple   BACnetRecipientEnclosed('0')
                    recipient                   ]
    [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                    processIdentifier           ]
]

[type BACnetRecipientEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [simple   BACnetRecipient
                    recipient                   ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type BACnetRecipient
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Device
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                            deviceValue                      ]
        ]
        ['1' *Address
            [simple   BACnetAddressEnclosed('1')
                            addressValue                     ]
        ]
    ]
]

[type BACnetPropertyReferenceEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [simple   BACnetPropertyReference
                    reference                   ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type BACnetPropertyReference
    [simple   BACnetPropertyIdentifierTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    propertyIdentifier              ]
    [optional   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                    arrayIndex                      ]
]

[type BACnetObjectPropertyReferenceEnclosed(uint 8 tagNumber)
   [simple   BACnetOpeningTag('tagNumber')
                   openingTag                  ]
   [simple   BACnetObjectPropertyReference
                   objectPropertyReference     ]
   [simple   BACnetClosingTag('tagNumber')
                   closingTag                  ]
]

[type BACnetObjectPropertyReference
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    objectIdentifier                ]
    [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    propertyIdentifier              ]
    [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                    arrayIndex                      ]
]

[type BACnetWriteAccessSpecification
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    objectIdentifier                ]
    [simple   BACnetOpeningTag('1')
                     openingTag                     ]
    [array      BACnetPropertyWriteDefinition('objectIdentifier.objectType')
                    listOfPropertyWriteDefinition
                    terminated
                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)'
    ]
    [simple   BACnetClosingTag('1')
                     closingTag                     ]
]

[type BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                  openingTag                  ]
    [peek       BACnetTagHeader
                        peekedTagHeader                 ]

    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0'  *Numeric
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER') numericValue             ]
        ]
        ['1'  *Character
            [simple   BACnetContextTagCharacterString('1', 'BACnetDataType.CHARACTER_STRING') characterValue           ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                  closingTag                  ]
]

[type BACnetConfirmedServiceRequestReadRangeRange
    [peek       BACnetTagHeader
                    peekedTagHeader                 ]
    [simple   BACnetOpeningTag('peekedTagHeader.actualTagNumber')
                     openingTag                     ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0x3'  *ByPosition
            [simple   BACnetApplicationTagUnsignedInteger                   referenceIndex            ]
            [simple   BACnetApplicationTagSignedInteger                     count                     ]
        ]
        ['0x6'  *BySequenceNumber
            [simple   BACnetApplicationTagUnsignedInteger                   referenceSequenceNumber   ]
            [simple   BACnetApplicationTagSignedInteger                     count                     ]
        ]
        ['0x7'  *ByTime
            [simple   BACnetDateTime                                        referenceTime             ]
            [simple   BACnetApplicationTagSignedInteger                     count                     ]
        ]
    ]
    [simple   BACnetClosingTag('peekedTagHeader.actualTagNumber')
                     closingTag
    ]
]

[type BACnetPropertyWriteDefinition(BACnetObjectType objectTypeArgument)
    [simple   BACnetPropertyIdentifierTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    propertyIdentifier              ]
    [optional   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                    arrayIndex                      ]
    [optional   BACnetConstructedData('2', 'objectTypeArgument', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                    propertyValue                   ]
    [optional   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                    priority                        ]
]

[type BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord
    [peek       BACnetTagHeader
                    peekedTagHeader                 ]
    [simple   BACnetOpeningTag('peekedTagHeader.actualTagNumber')
                     openingTag                     ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0x0' BACnetConfirmedServiceRequestAtomicReadFileStream
            [simple   BACnetApplicationTagSignedInteger                     fileStartPosition   ]
            [simple   BACnetApplicationTagUnsignedInteger                   requestOctetCount   ]
        ]
        ['0x1' BACnetConfirmedServiceRequestAtomicReadFileRecord
            [simple   BACnetApplicationTagSignedInteger                     fileStartRecord     ]
            [simple   BACnetApplicationTagUnsignedInteger                   requestRecordCount  ]
        ]
    ]
    [simple   BACnetClosingTag('peekedTagHeader.actualTagNumber')
                     closingTag                     ]
]

[discriminatedType BACnetUnconfirmedServiceRequest(uint 16 serviceRequestLength)
    [discriminator BACnetUnconfirmedServiceChoice serviceChoice]
    [typeSwitch serviceChoice
        ['I_AM'  *IAm
            [simple   BACnetApplicationTagObjectIdentifier                        deviceIdentifier                ]
            [simple   BACnetApplicationTagUnsignedInteger                         maximumApduLengthAcceptedLength ]
            [simple   BACnetSegmentationTagged('0', 'TagClass.APPLICATION_TAGS')  segmentationSupported           ]
            [simple   BACnetVendorIdTagged('2', 'TagClass.APPLICATION_TAGS')      vendorId                        ]
        ]
        ['I_HAVE'  *IHave
            [simple   BACnetApplicationTagObjectIdentifier                        deviceIdentifier    ]
            [simple   BACnetApplicationTagObjectIdentifier                        objectIdentifier    ]
            [simple   BACnetApplicationTagCharacterString                         objectName          ]
        ]
        ['UNCONFIRMED_COV_NOTIFICATION'  *UnconfirmedCOVNotification
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')          subscriberProcessIdentifier ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER') initiatingDeviceIdentifier  ]
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER') monitoredObjectIdentifier   ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')          lifetimeInSeconds           ]
            [simple   BACnetPropertyValues('4', 'monitoredObjectIdentifier.objectType')                listOfValues                ]
        ]
        ['UNCONFIRMED_EVENT_NOTIFICATION'  *UnconfirmedEventNotification
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')          processIdentifier            ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER') initiatingDeviceIdentifier   ]
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER') eventObjectIdentifier        ]
            [simple   BACnetTimeStampEnclosed('3')                                                     timestamp                    ]
            [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')          notificationClass            ]
            [simple   BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')          priority                     ]
            [simple   BACnetEventTypeTagged('6', 'TagClass.CONTEXT_SPECIFIC_TAGS')                     eventType                    ]
            [optional BACnetContextTagCharacterString('7', 'BACnetDataType.CHARACTER_STRING')          messageText                  ]
            [simple   BACnetNotifyTypeTagged('8', 'TagClass.CONTEXT_SPECIFIC_TAGS')                    notifyType                   ]
            [optional BACnetContextTagBoolean('9', 'BACnetDataType.BOOLEAN')                           ackRequired                  ]
            [optional BACnetEventStateTagged('10', 'TagClass.CONTEXT_SPECIFIC_TAGS')                   fromState                    ]
            [simple   BACnetEventStateTagged('11', 'TagClass.CONTEXT_SPECIFIC_TAGS')                   toState                      ]
            [optional BACnetNotificationParameters('12', 'eventObjectIdentifier.objectType')           eventValues                  ]
        ]
        ['UNCONFIRMED_PRIVATE_TRANSFER'  *UnconfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      vendorId                     ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')          serviceNumber                ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                                       serviceParameters            ]
        ]
        ['UNCONFIRMED_TEXT_MESSAGE'  *UnconfirmedTextMessage
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      textMessageSourceDevice     ]
            [optional BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass('1')                    messageClass                ] // Note we reuse the once from confirmed here
            [simple   BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            messagePriority             ] // Note we reuse the once from confirmed here
            [simple   BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')               message                     ]
        ]
        ['TIME_SYNCHRONIZATION'  *TimeSynchronization
            [simple   BACnetApplicationTagDate synchronizedDate]
            [simple   BACnetApplicationTagTime synchronizedTime]
        ]
        ['WHO_HAS'  *WhoHas
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeLowLimit                                         ]
            [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeHighLimit  'deviceInstanceRangeLowLimit != null' ]
            [simple   BACnetUnconfirmedServiceRequestWhoHasObject                                       object                                                              ]
        ]
        ['WHO_IS'  *WhoIs
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeLowLimit                                                 ]
            [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeHighLimit  'deviceInstanceRangeLowLimit != null'         ]
        ]
        ['UTC_TIME_SYNCHRONIZATION'  *UTCTimeSynchronization
            [simple   BACnetApplicationTagDate synchronizedDate]
            [simple   BACnetApplicationTagTime synchronizedTime]
        ]
        ['WRITE_GROUP'  *WriteGroup
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           groupNumber                 ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           writePriority               ]
            [simple   BACnetGroupChannelValueList('2')                                                  changeList                  ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')           inhibitDelay                ]
        ]
        ['UNCONFIRMED_COV_NOTIFICATION_MULTIPLE'  *UnconfirmedCOVNotificationMultiple
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           subscriberProcessIdentifier ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier  ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')           timeRemaining               ]
            [optional BACnetTimeStampEnclosed('3')                                                      timestamp                   ]
            [simple   ListOfCovNotificationsList('4')                                                   listOfCovNotifications      ]
        ]
        [* *Unknown
            [array    byte    unknownBytes length '(serviceRequestLength>0)?(serviceRequestLength - 1):0']
        ]
    ]
]

[type BACnetUnconfirmedServiceRequestWhoHasObject
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['2'  *Identifier
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                        objectIdentifier                            ]
        ]
        ['3'  *Name
            [simple   BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')
                                        objectName                                  ]
        ]
    ]
]

[discriminatedType BACnetServiceAck(uint 32 serviceAckLength)
    [discriminator   BACnetConfirmedServiceChoice
                        serviceChoice                   ]
    // we substract serviceChoice from our payload
    [virtual       uint 32  serviceAckPayloadLength '(serviceAckLength>0)?(serviceAckLength - 1):0'    ]
    [typeSwitch serviceChoice
        ////
        // Alarm and Event Services

        ['GET_ALARM_SUMMARY'  *GetAlarmSummary
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                ]
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')  eventState                      ]
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS')
                                                                                acknowledgedTransitions         ]
        ]
        ['GET_ENROLLMENT_SUMMARY'  *GetEnrollmentSummary
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                ]
            [simple   BACnetEventTypeTagged('0', 'TagClass.APPLICATION_TAGS')   eventType                       ]
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')  eventState                      ]
            [simple   BACnetApplicationTagUnsignedInteger                       priority                        ]
            [optional BACnetApplicationTagUnsignedInteger                       notificationClass               ]
        ]
        ['GET_EVENT_INFORMATION'  *GetEventInformation
            [simple   BACnetEventSummariesList('0')                             listOfEventSummaries            ]
            [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')    moreEvents                      ]
        ]
        //
        ////

        ////
        // File Access Services

        ['ATOMIC_READ_FILE'  *AtomicReadFile
            [simple   BACnetApplicationTagBoolean                               endOfFile                       ]
            [simple   BACnetServiceAckAtomicReadFileStreamOrRecord              accessMethod                    ]
        ]
        ['ATOMIC_WRITE_FILE'  *AtomicWriteFile
            [simple   BACnetContextTagSignedInteger('0', 'BACnetDataType.SIGNED_INTEGER') fileStartPosition     ]
        ]
        //
        ////

        ////
        // Object Access Services
        ['CREATE_OBJECT'  *CreateObject
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                ]
        ]
        ['READ_PROPERTY'  *ReadProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                                                                objectIdentifier                ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                propertyIdentifier              ]
            [optional   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                                                                arrayIndex                      ]
            [optional   BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                values                          ]
        ]
        ['READ_PROPERTY_MULTIPLE'  *ReadPropertyMultiple(uint 32 serviceAckPayloadLength)
            [array    BACnetReadAccessResult                                    data
                            length 'serviceAckPayloadLength'                                                    ]
        ]
        ['READ_RANGE'  *ReadRange
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier    ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier  ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               propertyArrayIndex  ]
            [simple   BACnetResultFlagsTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')                        resultFlags         ]
            [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')               itemCount           ]
            [optional BACnetConstructedData('5', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(propertyArrayIndex!=null?propertyArrayIndex.payload:null)')
                                                                                                            itemData            ]
            [optional BACnetContextTagUnsignedInteger('6', 'BACnetDataType.UNSIGNED_INTEGER')               firstSequenceNumber ]
        ]
        //
        ////


        ////
        // Remote Device Management Services

        ['CONFIRMED_PRIVATE_TRANSFER'  *ConfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                       vendorId                    ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           serviceNumber               ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                                        resultBlock                 ]
        ]
        //
        ////

        ////
        //  Virtual Terminal Services

        ['VT_OPEN'  *VTOpen
            [simple   BACnetApplicationTagUnsignedInteger                       remoteVtSessionIdentifier                        ]
        ]
        ['VT_DATA'  *VTData
            [simple   BACnetApplicationTagUnsignedInteger                       vtSessionIdentifier                              ]
            [simple   BACnetApplicationTagOctetString                           vtNewData                                        ]
            [simple   BACnetApplicationTagUnsignedInteger                       vtDataFlag                                       ]
        ]
        //
        ////


        ////
        //  Removed Services

        ['AUTHENTICATE'  *Authenticate(uint 32 serviceAckPayloadLength)
            [array    byte                                                      bytesOfRemovedService
                        length 'serviceAckPayloadLength'                                                                ]
        ]
        ['REQUEST_KEY'  *RequestKey(uint 32 serviceAckPayloadLength)
            [array    byte                                                      bytesOfRemovedService
                        length 'serviceAckPayloadLength'                                                                ]
        ]
        ['READ_PROPERTY_CONDITIONAL'  *ReadPropertyConditional(uint 32 serviceAckPayloadLength)
            [array    byte                                                      bytesOfRemovedService
                        length 'serviceAckPayloadLength'                                                                ]
        ]
        //
        ////
    ]
]

[type BACnetEventSummariesList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                     openingTag                     ]
    [array    BACnetEventSummary
                         listOfEventSummaries
                         terminated
                         'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
    ]
    [simple   BACnetClosingTag('tagNumber')
                     closingTag                     ]
]

[type BACnetEventSummary
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    objectIdentifier                ]
    [simple   BACnetEventStateTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    eventState                      ]
    [simple   BACnetEventTransitionBitsTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    acknowledgedTransitions         ]
    [simple   BACnetEventTimestampsEnclosed('3')
                    eventTimestamps                 ]
    [simple   BACnetNotifyTypeTagged('4', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    notifyType                      ]
    [simple   BACnetEventTransitionBitsTagged('5', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    eventEnable                     ]
    [simple   BACnetEventPriorities('6')
                    eventPriorities                 ]
]

[type BACnetEventTimestamps
    [simple  BACnetTimeStamp
                    toOffnormal                     ]
    [simple  BACnetTimeStamp
                    toFault                         ]
    [simple  BACnetTimeStamp
                    toNormal                        ]
]

[type BACnetEventTimestampsEnclosed(uint 8 tagNumber)
    [simple  BACnetOpeningTag('tagNumber')
                    openingTag                      ]
    [simple  BACnetEventTimestamps
                    eventTimestamps                 ]
    [simple  BACnetClosingTag('tagNumber')
                    closingTag                      ]
]

[type BACnetEventPriorities(uint 8 tagNumber)
    [simple  BACnetOpeningTag('tagNumber')
                    openingTag
    ]
    [simple  BACnetApplicationTagUnsignedInteger
                    toOffnormal                     ]
    [simple  BACnetApplicationTagUnsignedInteger
                    toFault                         ]
    [simple  BACnetApplicationTagUnsignedInteger
                    toNormal                        ]
    [simple  BACnetClosingTag('tagNumber')
                    closingTag
    ]
]

[type BACnetGroupChannelValueList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                     openingTag                     ]
    [array    BACnetEventSummary
                         listOfEventSummaries
                         terminated
                         'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
    ]
    [simple   BACnetClosingTag('tagNumber')
                     closingTag                     ]
]

[type BACnetGroupChannelValue
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')   channel                         ]
    [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')   overridingPriority              ]
    [simple   BACnetChannelValue                                                        value                           ]
]

[type BACnetChannelValue
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual  bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' *Null
           [simple  BACnetApplicationTagNull
                        nullValue                                                  ]
       ]
       ['0x4', 'false' *Real
           [simple  BACnetApplicationTagReal
                        realValue                                                  ]
       ]
       ['0x9', 'false' *Enumerated
           [simple   BACnetApplicationTagEnumerated
                       enumeratedValue                                             ]
       ]
       ['0x2', 'false' *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                       unsignedValue                                               ]
       ]
       ['0x1', 'false' *Boolean
           [simple   BACnetApplicationTagBoolean
                       booleanValue                                                ]
       ]
       ['0x3', 'false' *Integer
           [simple   BACnetApplicationTagSignedInteger
                       integerValue                                                ]
       ]
       ['0x5', 'false' *Double
           [simple  BACnetApplicationTagDouble
                        doubleValue                                                ]
       ]
       ['0xB', 'false' *Time
           [simple   BACnetApplicationTagTime
                       timeValue                                                   ]
       ]
       ['0x7', 'false' *CharacterString
           [simple   BACnetApplicationTagCharacterString
                       characterStringValue                                        ]
       ]
       ['0x6', 'false' *OctetString
           [simple   BACnetApplicationTagOctetString
                       octetStringValue                                            ]
       ]
       ['0x8', 'false' *BitString
           [simple   BACnetApplicationTagBitString
                       bitStringValue                                              ]
       ]
       ['0xA', 'false' *Date
           [simple   BACnetApplicationTagDate
                       dateValue                                                   ]
       ]
       ['0xC', 'false' *Objectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                       objectidentifierValue                                       ]
       ]
       ['0', 'true' *LightingCommand
           [simple   BACnetLightingCommandEnclosed('0')
                       ligthingCommandValue                                        ]
       ]
    ]
]

[type BACnetLightingCommandEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag          ]
    [simple   BACnetLightingCommand
                    lightingCommand     ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag          ]
]

[type BACnetLightingCommand
    [simple    BACnetLightingOperationTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                        lightningOperation              ]
    [optional  BACnetContextTagReal('1', 'BACnetDataType.REAL')
                        targetLevel                     ]
    [optional  BACnetContextTagReal('2', 'BACnetDataType.REAL')
                        rampRate                        ]
    [optional  BACnetContextTagReal('3', 'BACnetDataType.REAL')
                        stepIncrement                   ]
    [optional  BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')
                        fadeTime                        ]
    [optional  BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')
                        priority                        ]
]

[type BACnetReadAccessResult
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    objectIdentifier                  ]
    [optional BACnetReadAccessResultListOfResults('1', 'objectIdentifier.objectType')
                    listOfResults                     ]
]

[type BACnetReadAccessResultListOfResults(uint 8 tagNumber, BACnetObjectType objectTypeArgument)
    [simple   BACnetOpeningTag('tagNumber')
                     openingTag                                                                 ]
    [array    BACnetReadAccessProperty('objectTypeArgument')
                    listOfReadAccessProperty
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                     closingTag                                                                 ]
]

[type BACnetReadAccessProperty(BACnetObjectType objectTypeArgument)
    [simple   BACnetPropertyIdentifierTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    propertyIdentifier                                                          ]
    [optional   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                    arrayIndex                                                                  ]
    [optional   BACnetReadAccessPropertyReadResult('objectTypeArgument', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                    readResult                                                                  ]
]

[type BACnetReadAccessPropertyReadResult(BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument, BACnetTagPayloadUnsignedInteger arrayIndexArgument)
    [peek       BACnetTagHeader
                            peekedTagHeader                                                     ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber'               ]
    [optional   BACnetConstructedData('4', 'objectTypeArgument', 'propertyIdentifierArgument', 'arrayIndexArgument')
                    propertyValue           'peekedTagNumber == 4'                              ]
    [validation    '(peekedTagNumber == 4 && propertyValue != null) || peekedTagNumber != 4 '
                   "failure parsing field 4"                                                    ]
    [optional   ErrorEnclosed('5')
                    propertyAccessError     'peekedTagNumber == 5'                              ]
    [validation    '(peekedTagNumber == 5 && propertyAccessError != null) || peekedTagNumber != 5'
                   "failure parsing field 5"                                                    ]
    [validation    'peekedTagNumber == 4 || peekedTagNumber == 5'
                   "should be either 4 or 5"
                   shouldFail=false                                                             ]
]

[type BACnetServiceAckAtomicReadFileStreamOrRecord
    [peek       BACnetTagHeader
                            peekedTagHeader
    ]
    [simple   BACnetOpeningTag('peekedTagHeader.actualTagNumber')
                     openingTag
    ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0x0' BACnetServiceAckAtomicReadFileStream
            [simple   BACnetApplicationTagSignedInteger
                            fileStartPosition           ]
            [simple   BACnetApplicationTagOctetString
                            fileData                    ]
        ]
        ['0x1' BACnetServiceAckAtomicReadFileRecord
            [simple   BACnetApplicationTagSignedInteger
                            fileStartRecord             ]
            [simple   BACnetApplicationTagUnsignedInteger
                            returnedRecordCount         ]
            [array    BACnetApplicationTagOctetString
                            fileRecordData
                            count
                            'returnedRecordCount.payload.actualValue'   ]
        ]
    ]
    [simple   BACnetClosingTag('peekedTagHeader.actualTagNumber')
                     closingTag
    ]
]

[discriminatedType BACnetError(BACnetConfirmedServiceChoice errorChoice)
    [typeSwitch errorChoice
        ['SUBSCRIBE_COV_PROPERTY_MULTIPLE'  SubscribeCOVPropertyMultipleError
            [simple   ErrorEnclosed('0')
                        errorType                           ]
            [simple   SubscribeCOVPropertyMultipleErrorFirstFailedSubscription('1')
                        firstFailedSubscription             ]
        ]
        ['ADD_LIST_ELEMENT'                 ChangeListAddError
            [simple   ErrorEnclosed('0')
                            errorType                       ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            firstFailedElementNumber        ]
        ]
        ['REMOVE_LIST_ELEMENT'              ChangeListRemoveError
            [simple   ErrorEnclosed('0')
                            errorType                       ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            firstFailedElementNumber        ]
        ]
        ['CREATE_OBJECT'                    CreateObjectError
            [simple   ErrorEnclosed('0')
                            errorType                       ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            firstFailedElementNumber        ]
        ]
        ['WRITE_PROPERTY_MULTIPLE'          WritePropertyMultipleError
            [simple   ErrorEnclosed('0')
                            errorType                   ]
            [simple   BACnetObjectPropertyReferenceEnclosed('1')
                        firstFailedWriteAttempt             ]
        ]
        ['CONFIRMED_PRIVATE_TRANSFER'       ConfirmedPrivateTransferError
            [simple   ErrorEnclosed('0')
                            errorType                   ]
            [simple   BACnetVendorIdTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            vendorId                    ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            serviceNumber               ]
            [optional BACnetConstructedData('3', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            errorParameters             ]
        ]
        ['VT_CLOSE'                         VTCloseError
            [simple   ErrorEnclosed('0')
                            errorType                   ]
            [optional VTCloseErrorListOfVTSessionIdentifiers('1')
                            listOfVtSessionIdentifiers  ]
        ]
        [* *General
            [simple   Error
                            error               ]
        ]
    ]
]

[type VTCloseErrorListOfVTSessionIdentifiers(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [array      BACnetApplicationTagUnsignedInteger
                    listOfVtSessionIdentifiers
                             terminated
                             'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)'
                                                ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type SubscribeCOVPropertyMultipleErrorFirstFailedSubscription(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    monitoredObjectIdentifier   ]
    [simple   BACnetPropertyReferenceEnclosed('1')
                    monitoredPropertyReference  ]
    [simple   ErrorEnclosed('2')
                    errorType                   ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type ErrorEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag          ]
    [simple   Error
                    error               ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag          ]
]

[type Error
    [simple   ErrorClassTagged('0', 'TagClass.APPLICATION_TAGS') errorClass           ]
    [simple   ErrorCodeTagged('0', 'TagClass.APPLICATION_TAGS')  errorCode            ]
]

[type BACnetNotificationParameters(uint 8 tagNumber, BACnetObjectType objectTypeArgument)
    [simple   BACnetOpeningTag('tagNumber')
                            openingTag                                              ]
    [peek       BACnetTagHeader
                            peekedTagHeader                                         ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber'   ]
    [typeSwitch peekedTagNumber
        ['0'  *ChangeOfBitString(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagBitString('0', 'BACnetDataType.BIT_STRING')
                            changeOfBitString                                       ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['1'  *ChangeOfState(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetPropertyStatesEnclosed('0')
                            changeOfState                                           ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['2'  *ChangeOfValue(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetNotificationParametersChangeOfValueNewValue('0')
                            newValue                                                ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['3'  *CommandFailure(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            commandValue                                            ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetConstructedData('2', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            feedbackValue                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['4'  *FloatingLimit(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagReal('0', 'BACnetDataType.REAL')
                            referenceValue                                          ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                            setPointValue                                           ]
            [simple   BACnetContextTagReal('3', 'BACnetDataType.REAL')
                            errorLimit                                              ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['5'  *OutOfRange(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagReal('0', 'BACnetDataType.REAL')
                            exceedingValue                                          ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                            deadband                                                ]
            [simple   BACnetContextTagReal('3', 'BACnetDataType.REAL')
                            exceededLimit                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['6'  *ComplexEventType(uint 8 peekedTagNumber)
            [simple   BACnetPropertyValues('peekedTagNumber', 'objectTypeArgument')
                            listOfValues                                            ]
        ]
        // 7 is deprecated
        ['8'  *ChangeOfLifeSafety(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            newState                                                ]
            [simple   BACnetLifeSafetyModeTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            newMode                                                 ]
            [simple   BACnetStatusFlagsTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetLifeSafetyOperationTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            operationExpected                                       ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['9'  *Extended(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            vendorId                                                ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            extendedEventType                                       ]
            [simple   BACnetNotificationParametersExtendedParameters('2')
                            parameters                                              ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['10'  *BufferReady(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('0')
                            bufferProperty                                          ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            previousNotification                                    ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            currentNotification                                     ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['11'  *UnsignedRange(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                            sequenceNumber                                          ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            exceededLimit                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        // 12 is reserved
        ['13'  *AccessEvent(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetAccessEventTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            accessEvent                                             ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            accessEventTag                                          ]
            [simple   BACnetTimeStampEnclosed('3')
                            accessEventTime                                         ]
            [simple   BACnetDeviceObjectReferenceEnclosed('4')
                            accessCredential                                        ]
            [optional BACnetAuthenticationFactorTypeTagged('5', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            authenticationFactor                                    ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['14'  *DoubleOutOfRange(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagDouble('0', 'BACnetDataType.DOUBLE')
                            exceedingValue                                          ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagDouble('2', 'BACnetDataType.DOUBLE')
                            deadband                                                ]
            [simple   BACnetContextTagDouble('3', 'BACnetDataType.DOUBLE')
                            exceededLimit                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['15'  *SignedOutOfRange(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagSignedInteger('0', 'BACnetDataType.SIGNED_INTEGER')
                            exceedingValue                                          ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            deadband                                                ]
            [simple   BACnetContextTagSignedInteger('3', 'BACnetDataType.SIGNED_INTEGER')
                            exceededLimit                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['16'  *UnsignedOutOfRange(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                            exceedingValue                                          ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            deadband                                                ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                            exceededLimit                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['17'  *ChangeOfCharacterString(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagCharacterString('0', 'BACnetDataType.CHARACTER_STRING')
                            changedValue                                            ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetContextTagCharacterString('2', 'BACnetDataType.CHARACTER_STRING')
                            alarmValue                                              ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['18'  *ChangeOfStatusFlags(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            presentValue                                            ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            referencedFlags                                         ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['19'  *ChangeOfReliability(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetReliabilityTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            reliability                                             ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetPropertyValues('2', 'objectTypeArgument')
                            propertyValues                                          ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        // 20 is not used
        ['21'  *ChangeOfDiscreteValue(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetNotificationParametersChangeOfDiscreteValueNewValue('0')
                            newValue                                                ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['22'  *ChangeOfTimer(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetTimerStateTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            newValue                                                ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetDateTimeEnclosed('2')
                            updateTime                                              ]
            [optional BACnetTimerTransitionTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            lastStateChange                                         ]
            [optional BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')
                            initialTimeout                                          ]
            [optional BACnetDateTimeEnclosed('5')
                            expirationTime                                          ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                              ]
]

[type BACnetNotificationParametersChangeOfValueNewValue(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                        openingTag                                                  ]
    [peek       BACnetTagHeader
                        peekedTagHeader
    ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0'  *ChangedBits(uint 8 peekedTagNumber)
            [simple   BACnetContextTagBitString('0', 'BACnetDataType.BIT_STRING')
                        changedBits                                                 ]
        ]
        ['1'  *ChangedValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagReal('1', 'BACnetDataType.REAL')
                        changedValue                                                ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                        closingTag                                                  ]
]

[type BACnetNotificationParametersExtendedParameters(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                        openingTag                                                  ]
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual  bit       isOpeningTag        'peekedTagHeader.lengthValueType == 0x6']
    [virtual  bit       isClosingTag        'peekedTagHeader.lengthValueType == 0x7']
    [optional BACnetApplicationTagNull
                    nullValue
                        'peekedTagNumber == 0x0 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagReal
                    realValue
                        'peekedTagNumber == 0x4 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagUnsignedInteger
                    unsignedValue
                        'peekedTagNumber == 0x2 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagBoolean
                    booleanValue
                        'peekedTagNumber == 0x1 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagSignedInteger
                    integerValue
                        'peekedTagNumber == 0x3 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagDouble
                    doubleValue
                        'peekedTagNumber == 0x5 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagOctetString
                    octetStringValue
                        'peekedTagNumber == 0x6 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagCharacterString
                    characterStringValue
                        'peekedTagNumber == 0x7 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagBitString
                    bitStringValue
                        'peekedTagNumber == 0x8 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagEnumerated
                    enumeratedValue
                        'peekedTagNumber == 0x9 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagDate
                    dateValue
                        'peekedTagNumber == 0xA && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagTime
                    timeValue
                        'peekedTagNumber == 0xB && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagObjectIdentifier
                    objectIdentifier
                        'peekedTagNumber == 0xC && !isOpeningTag'                       ]
    [optional BACnetDeviceObjectPropertyReferenceEnclosed('0')
                    reference
                        'isOpeningTag && !isClosingTag'                                 ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetPropertyValues(uint 8 tagNumber, BACnetObjectType objectTypeArgument)
    [simple  BACnetOpeningTag('tagNumber')
                    innerOpeningTag                                                     ]
    [array    BACnetPropertyValue('objectTypeArgument')
                    data
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                        ]
    [simple  BACnetClosingTag('tagNumber')
                    innerClosingTag                                                     ]
]

[type BACnetPropertyValue(BACnetObjectType objectTypeArgument)
    [simple   BACnetPropertyIdentifierTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier  ]
    [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')               propertyArrayIndex  ]
    [optional BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifier.value', '(propertyArrayIndex!=null?propertyArrayIndex.payload:null)')
                                                                                                    propertyValue       ]
    [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')               priority            ]
]

[type BACnetDeviceObjectPropertyReferenceEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')                                                         openingTag          ]
    [simple   BACnetDeviceObjectPropertyReference                                                   value               ]
    [simple   BACnetClosingTag('tagNumber')                                                         closingTag          ]
]

[type BACnetNotificationParametersChangeOfDiscreteValueNewValue(uint 8 tagNumber)
   [simple   BACnetOpeningTag('tagNumber')
                       openingTag                                                  ]
   [peek       BACnetTagHeader
                       peekedTagHeader
   ]
   [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
   [virtual  bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
   [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
   [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x1', 'false'  *Boolean
           [simple   BACnetApplicationTagBoolean
                       booleanValue                                                ]
       ]
       ['0x2', 'false'  *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                       unsignedValue                                               ]
       ]
       ['0x3', 'false'  *Integer
           [simple   BACnetApplicationTagSignedInteger
                       integerValue                                                ]
       ]
       ['0x9', 'false'  *Enumerated
           [simple   BACnetApplicationTagEnumerated
                       enumeratedValue                                             ]
       ]
       ['0x7', 'false'  *CharacterString
           [simple   BACnetApplicationTagCharacterString
                       characterStringValue                                        ]
       ]
       ['0x6', 'false'  *OctetString
           [simple   BACnetApplicationTagOctetString
                       octetStringValue                                            ]
       ]
       ['0xA', 'false'  *OctetDate
           [simple   BACnetApplicationTagDate
                       dateValue                                                   ]
       ]
       ['0xB', 'false'  *OctetTime
           [simple   BACnetApplicationTagTime
                       timeValue                                                   ]
       ]
       ['0xC', 'false'  *Objectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                       objectidentifierValue                                       ]
       ]
       ['0', 'true'  *Datetime
           [simple   BACnetDateTimeEnclosed('0')
                       dateTimeValue                                               ]
       ]
   ]
   [simple   BACnetClosingTag('tagNumber')
                       closingTag                                                  ]
]

[type BACnetActionList
    [simple   BACnetOpeningTag('0')
                    innerOpeningTag                                                             ]
    [array    BACnetActionCommand
                    action
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 0)']
    [simple   BACnetClosingTag('0')
                    innerClosingTag                                                             ]
]

[type BACnetActionCommand
    [optional   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        deviceIdentifier                                                        ]
    [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        objectIdentifier                                                        ]
    [simple   BACnetPropertyIdentifierTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                        propertyIdentifier                                                      ]
    [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                        arrayIndex                                                              ]
    [optional BACnetConstructedData('4', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                        propertyValue                                                           ]
    [optional BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')
                        priority                                                                ]
    [optional BACnetContextTagBoolean('6', 'BACnetDataType.BOOLEAN')
                        postDelay                                                               ]
    [simple   BACnetContextTagBoolean('7', 'BACnetDataType.BOOLEAN')
                        quitOnFailure                                                           ]
    [simple   BACnetContextTagBoolean('8', 'BACnetDataType.BOOLEAN')
                        writeSuccessful                                                         ]
]

[type BACnetPriorityArray(BACnetObjectType objectTypeArgument, uint 8 tagNumber, BACnetTagPayloadUnsignedInteger arrayIndexArgument)
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
    [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
    [array    BACnetPriorityValue('objectTypeArgument')
                            data
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [virtual  BACnetPriorityValue   priorityValue01         'COUNT(data)>0?data[0]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue02         'COUNT(data)>1?data[1]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue03         'COUNT(data)>2?data[2]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue04         'COUNT(data)>3?data[3]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue05         'COUNT(data)>4?data[4]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue06         'COUNT(data)>5?data[5]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue07         'COUNT(data)>6?data[6]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue08         'COUNT(data)>7?data[7]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue09         'COUNT(data)>8?data[8]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue10         'COUNT(data)>9?data[9]:null'        ]
    [virtual  BACnetPriorityValue   priorityValue11         'COUNT(data)>10?data[10]:null'      ]
    [virtual  BACnetPriorityValue   priorityValue12         'COUNT(data)>11?data[11]:null'      ]
    [virtual  BACnetPriorityValue   priorityValue13         'COUNT(data)>12?data[12]:null'      ]
    [virtual  BACnetPriorityValue   priorityValue14         'COUNT(data)>13?data[13]:null'      ]
    [virtual  BACnetPriorityValue   priorityValue15         'COUNT(data)>14?data[14]:null'      ]
    [virtual  BACnetPriorityValue   priorityValue16         'COUNT(data)>15?data[15]:null'      ]
    [validation 'arrayIndexArgument != null || COUNT(data) == 16' "Either indexed access or lenght 16 expected" ]
    [virtual  bit                   isIndexedAccess         'COUNT(data) == 1'                  ]
    [virtual  BACnetPriorityValue   indexEntry              'priorityValue01'                   ]
]

[type BACnetPriorityValue(BACnetObjectType objectTypeArgument)
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual  bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false'  *Null
           [simple  BACnetApplicationTagNull
                            nullValue                                                   ]
       ]
       ['0x4', 'false'  *Real
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x9', 'false'  *Enumerated
           [simple   BACnetApplicationTagEnumerated
                            enumeratedValue                                             ]
       ]
       ['0x2', 'false'  *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x1', 'false'  *Boolean
           [simple   BACnetApplicationTagBoolean
                            booleanValue                                                ]
       ]
       ['0x3', 'false'  *Integer
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
       ['0x5', 'false'  *Double
           [simple  BACnetApplicationTagDouble
                                doubleValue                                             ]
       ]
       ['0xB', 'false'  *Time
           [simple   BACnetApplicationTagTime
                            timeValue                                                   ]
       ]
       ['0x7', 'false'  *CharacterString
           [simple   BACnetApplicationTagCharacterString
                            characterStringValue                                        ]
       ]
       ['0x6', 'false'  *OctetString
           [simple   BACnetApplicationTagOctetString
                            octetStringValue                                            ]
       ]
       ['0x8', 'false'  *BitString
           [simple   BACnetApplicationTagBitString
                            bitStringValue                                              ]
       ]
       ['0xA', 'false'  *Date
           [simple   BACnetApplicationTagDate
                            dateValue                                                   ]
       ]
       ['0xC', 'false'  *Objectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                            objectidentifierValue                                       ]
       ]
       ['0', 'true'  *ConstructedValue
           [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            constructedValue                                            ]
       ]
       ['1', 'true'  *DateTime
            [simple   BACnetDateTimeEnclosed('1')
                            dateTimeValue                                               ]
       ]
    ]
]

[type BACnetPropertyStatesEnclosed(uint 8 tagNumber)
    [simple  BACnetOpeningTag('tagNumber')
                    openingTag                                  ]
    [simple  BACnetPropertyStates
                    propertyState                               ]
    [simple  BACnetClosingTag('tagNumber')
                    closingTag                                  ]
]

[type BACnetPropertyStates
    [peek    BACnetTagHeader
                    peekedTagHeader                             ]
    [virtual  uint 8 peekedTagNumber
                        'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0'  *Boolean(uint 8 peekedTagNumber)
            [simple   BACnetContextTagBoolean('peekedTagNumber', 'BACnetDataType.BOOLEAN')
                                booleanValue                    ]
        ]
        ['1'  *BinaryValue(uint 8 peekedTagNumber)
            [simple   BACnetBinaryPVTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                binaryValue                     ]
        ]
        ['2'  *EventType(uint 8 peekedTagNumber)
            [simple   BACnetEventTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                eventType                       ]
        ]
        ['3'  *Polarity(uint 8 peekedTagNumber)
            [simple   BACnetPolarityTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                polarity                        ]
        ]
        ['4'  *ProgramChange(uint 8 peekedTagNumber)
            [simple   BACnetProgramRequestTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                programChange                   ]
        ]
        ['5'  *ProgramChange(uint 8 peekedTagNumber)
            [simple   BACnetProgramStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                programState                    ]
        ]
        ['6'  *ReasonForHalt(uint 8 peekedTagNumber)
            [simple   BACnetProgramErrorTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                reasonForHalt                   ]
        ]
        ['7'  *Reliability(uint 8 peekedTagNumber)
            [simple   BACnetReliabilityTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                reliability                     ]
        ]
        ['8'  *State(uint 8 peekedTagNumber)
            [simple   BACnetEventStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                state                           ]
        ]
        ['9'  *SystemStatus(uint 8 peekedTagNumber)
            [simple   BACnetDeviceStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                systemStatus                    ]
        ]
        ['10'  *Units(uint 8 peekedTagNumber)
            [simple   BACnetEngineeringUnitsTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                units                           ]
        ]
        ['11'  *ExtendedValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagUnsignedInteger('peekedTagNumber', 'BACnetDataType.UNSIGNED_INTEGER')
                                unsignedValue                   ]
        ]
        ['12'  *LifeSafetyMode(uint 8 peekedTagNumber)
            [simple   BACnetLifeSafetyModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lifeSafetyMode                  ]
        ]
        ['13'  *LifeSafetyState(uint 8 peekedTagNumber)
            [simple   BACnetLifeSafetyStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lifeSafetyState                 ]
        ]
        ['14'  *RestartReason(uint 8 peekedTagNumber)
            [simple   BACnetRestartReasonTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                restartReason                   ]
        ]
        ['15'  *DoorAlarmState(uint 8 peekedTagNumber)
            [simple   BACnetDoorAlarmStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorAlarmState                  ]
        ]
        ['16'  *Action(uint 8 peekedTagNumber)
            [simple   BACnetActionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                action                          ]
        ]
        ['17'  *DoorSecuredStatus(uint 8 peekedTagNumber)
            [simple   BACnetDoorSecuredStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorSecuredStatus               ]
        ]
        ['18'  *DoorStatus(uint 8 peekedTagNumber)
            [simple   BACnetDoorStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorStatus                      ]
        ]
        ['19'  *DoorValue(uint 8 peekedTagNumber)
            [simple   BACnetDoorValueTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorValue                       ]
        ]
        ['20'  *FileAccessMethod(uint 8 peekedTagNumber)
            [simple   BACnetFileAccessMethodTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                fileAccessMethod                ]
        ]
        ['21'  *LockStatus(uint 8 peekedTagNumber)
            [simple   BACnetLockStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lockStatus                      ]
        ]
        ['22'  *LifeSafetyOperations(uint 8 peekedTagNumber)
            [simple   BACnetLifeSafetyOperationTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lifeSafetyOperations            ]
        ]
        ['23'  *Maintenance(uint 8 peekedTagNumber)
            [simple   BACnetMaintenanceTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                maintenance                     ]
        ]
        ['24'  *NodeType(uint 8 peekedTagNumber)
            [simple   BACnetNodeTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                nodeType                        ]
        ]
        ['25'  *NotifyType(uint 8 peekedTagNumber)
            [simple   BACnetNotifyTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                notifyType                      ]
        ]
        ['26'  *SecurityLevel(uint 8 peekedTagNumber)
            [simple   BACnetSecurityLevelTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                securityLevel                   ]
        ]
        ['27'  *ShedState(uint 8 peekedTagNumber)
            [simple   BACnetShedStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                shedState                       ]
        ]
        ['28'  *SilencedState(uint 8 peekedTagNumber)
            [simple   BACnetSilencedStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                silencedState                   ]
        ]
        //['29'  *Reserved(uint 8 peekedTagNumber) ]
        ['30'  *AccessEvent(uint 8 peekedTagNumber)
            [simple   BACnetAccessEventTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accessEvent                     ]
        ]
        ['31'  *ZoneOccupanyState(uint 8 peekedTagNumber)
            [simple   BACnetAccessZoneOccupancyStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                zoneOccupanyState               ]
        ]
        ['32'  *AccessCredentialDisableReason(uint 8 peekedTagNumber)
            [simple   BACnetAccessCredentialDisableReasonTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accessCredentialDisableReason   ]
        ]
        ['33'  *AccessCredentialDisable(uint 8 peekedTagNumber)
            [simple   BACnetAccessCredentialDisableTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accessCredentialDisable         ]
        ]
        ['34'  *AuthenticationStatus(uint 8 peekedTagNumber)
            [simple   BACnetAuthenticationStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                authenticationStatus            ]
        ]
        // 35 is undefined
        ['36'  *BackupState(uint 8 peekedTagNumber)
            [simple    BACnetBackupStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                backupState                     ]
        ]
        ['37'  *WriteStatus(uint 8 peekedTagNumber)
            [simple    BACnetWriteStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                writeStatus                     ]
        ]
        ['38'  *LightningInProgress(uint 8 peekedTagNumber)
            [simple    BACnetLightingInProgressTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lightningInProgress             ]
        ]
        ['39'  *LightningOperation(uint 8 peekedTagNumber)
            [simple    BACnetLightingOperationTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lightningOperation              ]
        ]
        ['40'  *LightningTransition(uint 8 peekedTagNumber)
            [simple    BACnetLightingTransitionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lightningTransition             ]
        ]
        ['41'  *IntegerValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagSignedInteger('peekedTagNumber', 'BACnetDataType.SIGNED_INTEGER')
                                integerValue                    ]
        ]
        ['42'  *BinaryLightningValue(uint 8 peekedTagNumber)
            [simple   BACnetBinaryLightingPVTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                binaryLightningValue            ]
        ]
        ['43'  *TimerState(uint 8 peekedTagNumber)
            [simple   BACnetTimerStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                timerState                      ]
        ]
        ['44'  *TimerTransition(uint 8 peekedTagNumber)
            [simple   BACnetTimerTransitionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                timerTransition                 ]
        ]
        ['45'  *BacnetIpMode(uint 8 peekedTagNumber)
            [simple   BACnetIPModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                bacnetIpMode                    ]
        ]
        ['46'  *NetworkPortCommand(uint 8 peekedTagNumber)
            [simple   BACnetNetworkPortCommandTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                networkPortCommand              ]
        ]
        ['47'  *NetworkType(uint 8 peekedTagNumber)
            [simple   BACnetNetworkTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                networkType                     ]
        ]
        ['48'  *NetworkNumberQuality(uint 8 peekedTagNumber)
            [simple   BACnetNetworkNumberQualityTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                networkNumberQuality            ]
        ]
        ['49'  *EscalatorOperationDirection(uint 8 peekedTagNumber)
            [simple   BACnetEscalatorOperationDirectionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                escalatorOperationDirection     ]
        ]
        ['50'  *EscalatorFault(uint 8 peekedTagNumber)
            [simple   BACnetEscalatorFaultTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                escalatorFault                  ]
        ]
        ['51'  *EscalatorMode(uint 8 peekedTagNumber)
            [simple   BACnetEscalatorModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                escalatorMode                   ]
        ]
        ['52'  *LiftCarDirection(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarDirectionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarDirection                ]
        ]
        ['53'  *LiftCarDoorCommand(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarDoorCommandTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarDoorCommand              ]
        ]
        ['54'  *LiftCarDriveStatus(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarDriveStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarDriveStatus              ]
        ]
        ['55'  *LiftCarMode(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarMode                     ]
        ]
        ['56'  *LiftGroupMode(uint 8 peekedTagNumber)
            [simple   BACnetLiftGroupModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftGroupMode                   ]
        ]
        ['57'  *LiftFault(uint 8 peekedTagNumber)
            [simple   BACnetLiftFaultTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftFault                       ]
        ]
        ['58'  *ProtocolLevel(uint 8 peekedTagNumber)
            [simple   BACnetProtocolLevelTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                protocolLevel                   ]
        ]
        // 59 undefined
        // 60 undefined
        // 61 undefined
        // 62 undefined
        ['63'  *ExtendedValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagUnsignedInteger('peekedTagNumber', 'BACnetDataType.UNSIGNED_INTEGER')
                                extendedValue                   ]
        ]
        [BACnetPropertyStateActionUnknown(uint 8 peekedTagNumber)
            [simple   BACnetContextTagUnknown('peekedTagNumber', 'BACnetDataType.UNKNOWN')
                                unknownValue                    ]
        ]
    ]
]

[type BACnetTimeStamp
    [peek    BACnetTagHeader
                        peekedTagHeader                         ]
    [virtual  uint 8     peekedTagNumber
                            'peekedTagHeader.actualTagNumber'   ]
    [typeSwitch peekedTagNumber
        ['0'  *Time
            [simple   BACnetContextTagTime('0', 'BACnetDataType.TIME')
                            timeValue                           ]
        ]
        ['1'  *Sequence
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            sequenceNumber                      ]
        ]
        ['2'  *DateTime
            [simple   BACnetDateTimeEnclosed('2')
                            dateTimeValue                       ]
        ]
    ]
]

[type BACnetTimeStampEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag          ]
    [simple   BACnetTimeStamp
                    timestamp           ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag          ]
]

[type BACnetTimeStampsEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag ]
    [array      BACnetTimeStamp
                        timestamps
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                            ]
    [simple   BACnetClosingTag('tagNumber')
                        closingTag          ]
]

[type BACnetDateTime
    [simple   BACnetApplicationTagDate
                        dateValue           ]
    [simple   BACnetApplicationTagTime
                        timeValue           ]
]

[type BACnetDateTimeEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                        openingTag          ]
    [simple   BACnetDateTime
                        dateTimeValue       ]
    [simple   BACnetClosingTag('tagNumber')
                        closingTag          ]
]

[type BACnetAddress
    [simple   BACnetApplicationTagUnsignedInteger
                        networkNumber       ]
    // TODO: uint 64 ---> big int in java == boom
    [virtual  uint 64   zero           '0'  ]
    [virtual  bit   isLocalNetwork  'networkNumber.actualValue == zero']
    [simple   BACnetApplicationTagOctetString
                        macAddress          ]
    [virtual  bit   isBroadcast  'macAddress.actualLength == 0']
]

[type BACnetAddressEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                  ]
    [simple   BACnetAddress
                    address                     ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type BACnetConstructedData(uint 8 tagNumber, BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument, BACnetTagPayloadUnsignedInteger arrayIndexArgument)
    [simple   BACnetOpeningTag('tagNumber')
                        openingTag                                                                              ]
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch objectTypeArgument, propertyIdentifierArgument, peekedTagNumber
        [*, 'ABSENTEE_LIMIT', '2'                       *AbsenteeLimit
            [simple   BACnetApplicationTagUnsignedInteger                     absenteeLimit                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                     actualValue       'absenteeLimit'         ]
        ]
        [*, 'ACCEPTED_MODES'                            *AcceptedModes
            [array    BACnetLifeSafetyModeTagged('0', 'TagClass.APPLICATION_TAGS')
                            acceptedModes              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_ALARM_EVENTS', '9'                  *AccessAlarmEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                    accessAlarmEvents
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_DOORS'                              *AccessDoors
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectReference
                                accessDoors
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_EVENT', '9'                         *AccessEvent
            [simple   BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS') accessEvent                             ]
            [virtual  BACnetAccessEventTagged                                   actualValue       'accessEvent'         ]
        ]
        [*, 'ACCESS_EVENT_AUTHENTICATION_FACTOR'        *AccessEventAuthenticationFactor
            [simple   BACnetAuthenticationFactor                                accessEventAuthenticationFactor         ]
            [virtual  BACnetAuthenticationFactor                                actualValue       'accessEventAuthenticationFactor'         ]
        ]
        [*, 'ACCESS_EVENT_CREDENTIAL'                   *AccessEventCredential
            [simple   BACnetDeviceObjectReference                               accessEventCredential                   ]
            [virtual  BACnetDeviceObjectReference                               actualValue       'accessEventCredential'         ]
        ]
        [*, 'ACCESS_EVENT_TAG', '2'                     *AccessEventTag
            [simple   BACnetApplicationTagUnsignedInteger                       accessEventTag                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue       'accessEventTag'      ]
        ]
        [*, 'ACCESS_EVENT_TIME'                         *AccessEventTime
            [simple   BACnetTimeStamp                                           accessEventTime                         ]
            [virtual  BACnetTimeStamp                                           actualValue       'accessEventTime'     ]
        ]
        [*, 'ACCESS_TRANSACTION_EVENTS', '9'            *AccessTransactionEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                        accessTransactionEvents
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCOMPANIMENT'                             *Accompaniment
            [simple   BACnetDeviceObjectReference                               accompaniment                           ]
            [virtual  BACnetDeviceObjectReference                               actualValue       'accompaniment'       ]
        ]
        [*, 'ACCOMPANIMENT_TIME', '2'                   *AccompanimentTime
            [simple   BACnetApplicationTagUnsignedInteger                       accompanimentTime                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue       'accompanimentTime'   ]
        ]
        [*, 'ACK_REQUIRED', '9'                         *AckRequired
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS')
                                                                                ackRequired                             ]
            [virtual  BACnetEventTransitionBitsTagged                           actualValue       'ackRequired'         ]
        ]
        [*, 'ACKED_TRANSITIONS', '9'                    *AckedTransitions
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS') ackedTransitions                ]
            [virtual  BACnetEventTransitionBitsTagged                           actualValue       'ackedTransitions'    ]
        ]
        ['LOOP', 'ACTION', '9'                          *LoopAction
            [simple   BACnetActionTagged('0', 'TagClass.APPLICATION_TAGS')      action                                  ]
            [virtual  BACnetActionTagged                                        actualValue       'action'              ]
        ]
        ['COMMAND', 'ACTION'                            *CommandAction
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetActionList
                            actionLists
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTION'                                    *Action
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetActionList
                            actionLists
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTION_TEXT', '7'                          *ActionText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                    actionText
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ACTIVATION_TIME'                           *ActivationTime
            [simple   BACnetDateTime                                            activationTime                          ]
            [virtual  BACnetDateTime                                            actualValue       'activationTime'      ]
        ]
        [*, 'ACTIVE_AUTHENTICATION_POLICY', '2'         *ActiveAuthenticationPolicy
            [simple   BACnetApplicationTagUnsignedInteger                               activeAuthenticationPolicy      ]
            [virtual  BACnetApplicationTagUnsignedInteger actualValue 'activeAuthenticationPolicy']
        ]
        [*, 'ACTIVE_COV_MULTIPLE_SUBSCRIPTIONS'         *ActiveCOVMultipleSubscriptions
            [array    BACnetCOVMultipleSubscription
                                activeCOVMultipleSubscriptions
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTIVE_COV_SUBSCRIPTIONS'                  *ActiveCOVSubscriptions
            [array    BACnetCOVSubscription
                            activeCOVSubscriptions
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'ACTIVE_TEXT', '7'                          *ActiveText
            [simple   BACnetApplicationTagCharacterString             activeText                                        ]
            [virtual  BACnetApplicationTagCharacterString actualValue 'activeText']
        ]
        [*, 'ACTIVE_VT_SESSIONS'                        *ActiveVTSessions
            [array    BACnetVTSession
                                activeVTSession
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTUAL_SHED_LEVEL'                         *ActualShedLevel
            [simple   BACnetShedLevel                                           actualShedLevel                         ]
            [virtual  BACnetShedLevel                                           actualValue       'actualShedLevel'     ]
        ]
        ['ACCESS_ZONE', 'ADJUST_VALUE', '3'             *AccessZoneAdjustValue
            [simple   BACnetApplicationTagSignedInteger                         adjustValue                             ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue         'adjustValue'       ]
        ]
        ['PULSE_CONVERTER', 'ADJUST_VALUE', '4'         *PulseConverterAdjustValue
            [simple   BACnetApplicationTagReal                                  adjustValue                             ]
            [virtual  BACnetApplicationTagReal                                  actualValue         'adjustValue'       ]
        ]
        [*, 'ADJUST_VALUE', '3'                         *AdjustValue
            [simple   BACnetApplicationTagSignedInteger                         adjustValue                             ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'adjustValue'               ]
        ]
        // TODO: pretty sure we need to catch a generic application tag here
        [*, 'ALARM_VALUE', '9'                          *AlarmValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    binaryPv                                ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'binaryPv'                  ]
        ]
        ['ACCESS_DOOR', 'ALARM_VALUES'                  *AccessDoorAlarmValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['ACCESS_ZONE', 'ALARM_VALUES'                  *AccessZoneAlarmValues
            [array    BACnetAccessZoneOccupancyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['BITSTRING_VALUE', 'ALARM_VALUES', '8'         *BitStringValueAlarmValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagBitString
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['CHARACTERSTRING_VALUE', 'ALARM_VALUES'             *CharacterStringValueAlarmValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_POINT', 'ALARM_VALUES'            *LifeSafetyPointAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_ZONE', 'ALARM_VALUES'             *LifeSafetyZoneAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_INPUT', 'ALARM_VALUES', '2'       *MultiStateInputAlarmValues
            [array    BACnetApplicationTagUnsignedInteger
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_VALUE', 'ALARM_VALUES', '2'       *MultiStateValueAlarmValues
            [array    BACnetApplicationTagUnsignedInteger
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['TIMER', 'ALARM_VALUES'                        *TimerAlarmValues
            [array    BACnetTimerStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ALARM_VALUES'                              *AlarmValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ALIGN_INTERVALS', '1'                      *AlignIntervals
            [simple   BACnetApplicationTagBoolean                                   alignIntervals                      ]
            [virtual  BACnetApplicationTagBoolean                                   actualValue 'alignIntervals'        ]
        ]

        /////
        // All property implementations for every object

        ['ACCESS_CREDENTIAL'     , 'ALL'                *AccessCredentialAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_DOOR'           , 'ALL'                *AccessDoorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_POINT'          , 'ALL'                *AccessPointAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_RIGHTS'         , 'ALL'                *AccessRightsAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_USER'           , 'ALL'                *AccessUserAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_ZONE'           , 'ALL'                *AccessZoneAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCUMULATOR'           , 'ALL'                *AccumulatorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ALERT_ENROLLMENT'      , 'ALL'                *AlertEnrollmentAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_INPUT'          , 'ALL'                *AnalogInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_OUTPUT'         , 'ALL'                *AnalogOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_VALUE'          , 'ALL'                *AnalogValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['AVERAGING'             , 'ALL'                *AveragingAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_INPUT'          , 'ALL'                *BinaryInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'ALL'                *BinaryLightingOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_OUTPUT'         , 'ALL'                *BinaryOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_VALUE'          , 'ALL'                *BinaryValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BITSTRING_VALUE'       , 'ALL'                *BitstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CALENDAR'              , 'ALL'                *CalendarAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CHANNEL'               , 'ALL'                *ChannelAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CHARACTERSTRING_VALUE' , 'ALL'                *CharacterstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['COMMAND'               , 'ALL'                *CommandAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CREDENTIAL_DATA_INPUT' , 'ALL'                *CredentialDataInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATEPATTERN_VALUE'     , 'ALL'                *DatepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATE_VALUE'            , 'ALL'                *DateValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATETIMEPATTERN_VALUE' , 'ALL'                *DatetimepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATETIME_VALUE'        , 'ALL'                *DatetimeValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DEVICE'                , 'ALL'                *DeviceAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ELEVATOR_GROUP'        , 'ALL'                *ElevatorGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ESCALATOR'             , 'ALL'                *EscalatorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['EVENT_ENROLLMENT'      , 'ALL'                *EventEnrollmentAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['EVENT_LOG'             , 'ALL'                *EventLogAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['FILE'                  , 'ALL'                *FileAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['GLOBAL_GROUP'          , 'ALL'                *GlobalGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['GROUP'                 , 'ALL'                *GroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['INTEGER_VALUE'         , 'ALL'                *IntegerValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LARGE_ANALOG_VALUE'    , 'ALL'                *LargeAnalogValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFE_SAFETY_POINT'     , 'ALL'                *LifeSafetyPointAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFE_SAFETY_ZONE'      , 'ALL'                *LifeSafetyZoneAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFT'                  , 'ALL'                *LiftAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIGHTING_OUTPUT'       , 'ALL'                *LightingOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LOAD_CONTROL'          , 'ALL'                *LoadControlAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LOOP'                  , 'ALL'                *LoopAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_INPUT'     , 'ALL'                *MultiStateInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_OUTPUT'    , 'ALL'                *MultiStateOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_VALUE'     , 'ALL'                *MultiStateValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NETWORK_PORT'          , 'ALL'                *NetworkPortAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NETWORK_SECURITY'      , 'ALL'                *NetworkSecurityAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NOTIFICATION_CLASS'    , 'ALL'                *NotificationClassAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NOTIFICATION_FORWARDER', 'ALL'                *NotificationForwarderAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['OCTETSTRING_VALUE'     , 'ALL'                *OctetstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['POSITIVE_INTEGER_VALUE', 'ALL'                *PositiveIntegerValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['PROGRAM'               , 'ALL'                *ProgramAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['PULSE_CONVERTER'       , 'ALL'                *PulseConverterAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['SCHEDULE'              , 'ALL'                *ScheduleAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['STRUCTURED_VIEW'       , 'ALL'                *StructuredViewAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIMEPATTERN_VALUE'     , 'ALL'                *TimepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIME_VALUE'            , 'ALL'                *TimeValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIMER'                 , 'ALL'                *TimerAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TREND_LOG'             , 'ALL'                *TrendLogAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TREND_LOG_MULTIPLE'    , 'ALL'                *TrendLogMultipleAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        //
        /////

        [*, 'ALL_WRITES_SUCCESSFUL', '1'                *AllWritesSuccessful
            [simple   BACnetApplicationTagBoolean                               allWritesSuccessful                     ]
            [virtual  BACnetApplicationTagBoolean actualValue 'allWritesSuccessful']
        ]
        [*, 'ALLOW_GROUP_DELAY_INHIBIT', '1'            *AllowGroupDelayInhibit
            [simple   BACnetApplicationTagBoolean                               allowGroupDelayInhibit                  ]
            [virtual  BACnetApplicationTagBoolean                               actualValue     'allowGroupDelayInhibit']
        ]
        [*, 'APDU_LENGTH', '2'                          *APDULength
            [simple   BACnetApplicationTagUnsignedInteger                       apduLength                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue         'apduLength'        ]
        ]
        [*, 'APDU_SEGMENT_TIMEOUT', '2'                 *APDUSegmentTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       apduSegmentTimeout                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'apduSegmentTimeout'        ]
        ]
        [*, 'APDU_TIMEOUT', '2'                         *APDUTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       apduTimeout                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'apduTimeout'               ]
        ]
        [*, 'APPLICATION_SOFTWARE_VERSION', '7'         *ApplicationSoftwareVersion
            [simple   BACnetApplicationTagCharacterString                       applicationSoftwareVersion              ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'applicationSoftwareVersion']
        ]
        [*, 'ARCHIVE', '1'                              *Archive
            [simple   BACnetApplicationTagBoolean                               archive                                 ]
            [virtual  BACnetApplicationTagBoolean                               actualValue         'archive'           ]
        ]
        [*, 'ASSIGNED_ACCESS_RIGHTS'                    *AssignedAccessRights
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAssignedAccessRights
                                        assignedAccessRights
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ASSIGNED_LANDING_CALLS'                    *AssignedLandingCalls
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAssignedLandingCalls
                                        assignedLandingCalls
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ATTEMPTED_SAMPLES', '2'                    *AttemptedSamples
            [simple   BACnetApplicationTagUnsignedInteger                       attemptedSamples                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'attemptedSamples'          ]
        ]
        [*, 'AUTHENTICATION_FACTORS'                    *AuthenticationFactors
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetCredentialAuthenticationFactor
                            authenticationFactors
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'AUTHENTICATION_POLICY_LIST'                *AuthenticationPolicyList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAuthenticationPolicy
                            authenticationPolicyList
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'AUTHENTICATION_POLICY_NAMES', '7'          *AuthenticationPolicyNames
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                        authenticationPolicyNames
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHENTICATION_STATUS', '9'                *AuthenticationStatus
            [simple   BACnetAuthenticationStatusTagged('0', 'TagClass.APPLICATION_TAGS') authenticationStatus           ]
            [virtual  BACnetAuthenticationStatusTagged                               actualValue  'authenticationStatus']
        ]
        [*, 'AUTHORIZATION_EXEMPTIONS', '9'             *AuthorizationExemptions
            [array    BACnetAuthorizationExemptionTagged('0', 'TagClass.APPLICATION_TAGS')
                                        authorizationExemption
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHORIZATION_MODE', '9'                   *AuthorizationMode
            [simple   BACnetAuthorizationModeTagged('0', 'TagClass.APPLICATION_TAGS') authorizationMode                 ]
            [virtual  BACnetAuthorizationModeTagged                             actualValue 'authorizationMode'         ]
        ]
        [*, 'AUTO_SLAVE_DISCOVERY', '1'                 *AutoSlaveDiscovery
            [simple   BACnetApplicationTagBoolean                               autoSlaveDiscovery                      ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'autoSlaveDiscovery'        ]
        ]
        [*, 'AVERAGE_VALUE', '4'                        *AverageValue
            [simple   BACnetApplicationTagReal                                  averageValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue  'averageValue'             ]
        ]
        [*, 'BACKUP_AND_RESTORE_STATE', '9'             *BackupAndRestoreState
            [simple   BACnetBackupStateTagged('0', 'TagClass.APPLICATION_TAGS') backupAndRestoreState                   ]
            [virtual  BACnetBackupStateTagged                                   actualValue  'backupAndRestoreState'    ]
        ]
        [*, 'BACKUP_FAILURE_TIMEOUT', '2'               *BackupFailureTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       backupFailureTimeout                    ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'backupFailureTimeout'      ]
        ]
        [*, 'BACKUP_PREPARATION_TIME', '2'              *BackupPreparationTime
            [simple   BACnetApplicationTagUnsignedInteger                       backupPreparationTime                   ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'backupPreparationTime'     ]
        ]
        [*, 'BACNET_IP_GLOBAL_ADDRESS'                  *BACnetIPGlobalAddress
            [simple   BACnetHostNPort                                           bacnetIpGlobalAddress                   ]
            [virtual  BACnetHostNPort                                           actualValue 'bacnetIpGlobalAddress'     ]
        ]
        [*, 'BACNET_IP_MODE', '9'                       *BACnetIPMode
            [simple   BACnetIPModeTagged('0', 'TagClass.APPLICATION_TAGS')      bacnetIpMode                            ]
            [virtual  BACnetIPModeTagged                                        actualValue 'bacnetIpMode'              ]
        ]
        [*, 'BACNET_IP_MULTICAST_ADDRESS', '6'          *BACnetIPMulticastAddress
            [simple   BACnetApplicationTagOctetString                           ipMulticastAddress                      ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipMulticastAddress'        ]
        ]
        [*, 'BACNET_IP_NAT_TRAVERSAL', '1'              *BACnetIPNATTraversal
            [simple   BACnetApplicationTagBoolean                               bacnetIPNATTraversal                    ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'bacnetIPNATTraversal'      ]
        ]
        [*, 'BACNET_IP_UDP_PORT', '2'                   *BACnetIPUDPPort
            [simple   BACnetApplicationTagUnsignedInteger                       ipUdpPort                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipUdpPort'                 ]
        ]
        [*, 'BACNET_IPV6_MODE', '9'                     *BACnetIPv6Mode
            [simple   BACnetIPModeTagged('0', 'TagClass.APPLICATION_TAGS')      bacnetIpv6Mode                          ]
            [virtual  BACnetIPModeTagged                                        actualValue 'bacnetIpv6Mode'            ]
        ]
        [*, 'BACNET_IPV6_UDP_PORT', '2'                 *BACnetIPv6UDPPort
            [simple   BACnetApplicationTagUnsignedInteger                       ipv6UdpPort                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipv6UdpPort'               ]
        ]
        [*, 'BACNET_IPV6_MULTICAST_ADDRESS', '6'        *BACnetIPv6MulticastAddress
            [simple   BACnetApplicationTagOctetString                           ipv6MulticastAddress                    ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipv6MulticastAddress'      ]
        ]
        [*, 'BASE_DEVICE_SECURITY_POLICY', '9'          *BaseDeviceSecurityPolicy
            [simple   BACnetSecurityLevelTagged('0', 'TagClass.APPLICATION_TAGS') baseDeviceSecurityPolicy              ]
            [virtual  BACnetSecurityLevelTagged                                 actualValue 'baseDeviceSecurityPolicy'  ]
        ]
        [*, 'BBMD_ACCEPT_FD_REGISTRATIONS', '1'         *BBMDAcceptFDRegistrations
            [simple   BACnetApplicationTagBoolean                               bbmdAcceptFDRegistrations               ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'bbmdAcceptFDRegistrations' ]
        ]
        [*, 'BBMD_BROADCAST_DISTRIBUTION_TABLE'         *BBMDBroadcastDistributionTable
            [array    BACnetBDTEntry
                                bbmdBroadcastDistributionTable
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'BBMD_FOREIGN_DEVICE_TABLE'                 *BBMDForeignDeviceTable
            [array    BACnetBDTEntry
                                bbmdForeignDeviceTable
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'BELONGS_TO'                                *BelongsTo
            [simple   BACnetDeviceObjectReference                               belongsTo                               ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'belongsTo'                 ]
        ]
        [*, 'BIAS', '4'                                 *Bias
            [simple   BACnetApplicationTagReal                                  bias                                    ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'bias'                      ]
        ]
        [*, 'BIT_MASK', '8'                             *BitMask
            [simple   BACnetApplicationTagBitString                             bitString                               ]
            [virtual  BACnetApplicationTagBitString                             actualValue 'bitString'                 ]
        ]
        [*, 'BIT_TEXT', '7'                             *BitText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                        bitText
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'BLINK_WARN_ENABLE', '1'                    *BlinkWarnEnable
            [simple   BACnetApplicationTagBoolean                               blinkWarnEnable                         ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'blinkWarnEnable'           ]
        ]
        [*, 'BUFFER_SIZE', '2'                          *BufferSize
            [simple   BACnetApplicationTagUnsignedInteger                       bufferSize                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'bufferSize'                ]
        ]
        [*, 'CAR_ASSIGNED_DIRECTION', '9'               *CarAssignedDirection
            [simple   BACnetLiftCarDirectionTagged('0', 'TagClass.APPLICATION_TAGS')             assignedDirection      ]
            [virtual  BACnetLiftCarDirectionTagged                              actualValue 'assignedDirection'         ]
        ]
        [*, 'CAR_DOOR_COMMAND', '9'                     *CarDoorCommand
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLiftCarDoorCommandTagged('0', 'TagClass.APPLICATION_TAGS')
                            carDoorCommand
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'CAR_DOOR_STATUS', '9'                      *CarDoorStatus
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDoorStatusTagged('0', 'TagClass.APPLICATION_TAGS')
                            carDoorStatus
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'CAR_DOOR_TEXT', '7'                        *CarDoorText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                carDoorText
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CAR_DOOR_ZONE', '1'                        *CarDoorZone
            [simple   BACnetApplicationTagBoolean                               carDoorZone                             ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'carDoorZone'               ]
        ]
        [*, 'CAR_DRIVE_STATUS', '9'                     *CarDriveStatus
            [simple   BACnetLiftCarDriveStatusTagged('0', 'TagClass.APPLICATION_TAGS')  carDriveStatus                  ]
            [virtual  BACnetLiftCarDriveStatusTagged                            actualValue 'carDriveStatus'            ]
        ]
        [*, 'CAR_LOAD', '4'                             *CarLoad
            [simple   BACnetApplicationTagReal                                  carLoad                                 ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'carLoad'                   ]
        ]
        [*, 'CAR_LOAD_UNITS', '9'                       *CarLoadUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        [*, 'CAR_MODE', '9'                             *CarMode
            [simple   BACnetLiftCarModeTagged('0', 'TagClass.APPLICATION_TAGS') carMode                                 ]
            [virtual  BACnetLiftCarModeTagged                                   actualValue 'carMode'                   ]
        ]
        [*, 'CAR_MOVING_DIRECTION', '9'                 *CarMovingDirection
            [simple   BACnetLiftCarDirectionTagged('0', 'TagClass.APPLICATION_TAGS')             carMovingDirection     ]
            [virtual  BACnetLiftCarDirectionTagged                              actualValue 'carMovingDirection'        ]
        ]
        [*, 'CAR_POSITION', '2'                         *CarPosition
            [simple   BACnetApplicationTagUnsignedInteger                       carPosition                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'carPosition'               ]
        ]
        [*, 'CHANGE_OF_STATE_COUNT', '2'                *ChangeOfStateCount
            [simple   BACnetApplicationTagUnsignedInteger                       changeIfStateCount                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'changeIfStateCount'        ]
        ]
        [*, 'CHANGE_OF_STATE_TIME'                      *ChangeOfStateTime
            [simple   BACnetDateTime                                            changeOfStateTime                       ]
            [virtual  BACnetDateTime                                            actualValue 'changeOfStateTime'         ]
        ]
        [*, 'CHANGES_PENDING', '1'                      *ChangesPending
            [simple   BACnetApplicationTagBoolean                               changesPending                          ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'changesPending'            ]
        ]
        [*, 'CHANNEL_NUMBER', '2'                       *ChannelNumber
            [simple   BACnetApplicationTagUnsignedInteger                       channelNumber                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'channelNumber'             ]
        ]
        [*, 'CLIENT_COV_INCREMENT'                      *ClientCOVIncrement
            [simple   BACnetClientCOV                                           covIncrement                            ]
            [virtual  BACnetClientCOV                                           actualValue 'covIncrement'              ]
        ]
        [*, 'COMMAND', '9'                              *Command
            [simple   BACnetNetworkPortCommandTagged('0', 'TagClass.APPLICATION_TAGS')                  command         ]
            [virtual  BACnetNetworkPortCommandTagged                            actualValue 'command'                   ]
        ]
        [*, 'COMMAND_TIME_ARRAY'                        *CommandTimeArray
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetTimeStamp
                            commandTimeArray
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
            [validation 'arrayIndexArgument!=null || COUNT(commandTimeArray) == 16'
                        "commandTimeArray should have exactly 16 values"                                                ]
        ]
        [*, 'CONFIGURATION_FILES', '12'                 *ConfigurationFiles
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                                configurationFiles
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CONTROL_GROUPS', '2'                       *ControlGroups
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                controlGroups
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CONTROLLED_VARIABLE_REFERENCE'             *ControlledVariableReference
            [simple   BACnetObjectPropertyReference                             controlledVariableReference             ]
            [virtual  BACnetObjectPropertyReference                             actualValue 'controlledVariableReference']
        ]
        [*, 'CONTROLLED_VARIABLE_UNITS', '9'            *ControlledVariableUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')    units                           ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        [*, 'CONTROLLED_VARIABLE_VALUE', '4'            *ControlledVariableValue
            [simple   BACnetApplicationTagReal                                          controlledVariableValue         ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'controlledVariableValue'   ]
        ]
        [*, 'COUNT', '2'                                *Count
            [simple   BACnetApplicationTagUnsignedInteger                       count                                   ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'count'                     ]
        ]
        [*, 'COUNT_BEFORE_CHANGE', '2'                  *CountBeforeChange
            [simple   BACnetApplicationTagUnsignedInteger                       countBeforeChange                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'countBeforeChange'         ]
        ]
        [*, 'COUNT_CHANGE_TIME'                         *CountChangeTime
            [simple   BACnetDateTime                                            countChangeTime                         ]
            [virtual  BACnetDateTime                                            actualValue 'countChangeTime'           ]
        ]
        ['INTEGER_VALUE', 'COV_INCREMENT', '2'          *IntegerValueCOVIncrement
            [simple   BACnetApplicationTagUnsignedInteger                       covIncrement                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'covIncrement'              ]
        ]
        ['LARGE_ANALOG_VALUE', 'COV_INCREMENT', '5'     *LargeAnalogValueCOVIncrement
            [simple   BACnetApplicationTagDouble                                covIncrement                            ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'covIncrement'              ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'COV_INCREMENT', '2' *PositiveIntegerValueCOVIncrement
            [simple   BACnetApplicationTagUnsignedInteger                       covIncrement                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'covIncrement'              ]
        ]
        [*, 'COV_INCREMENT', '4'                        *COVIncrement
            [simple   BACnetApplicationTagReal                                  covIncrement                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'covIncrement'              ]
        ]
        [*, 'COV_PERIOD', '2'                           *COVPeriod
            [simple   BACnetApplicationTagUnsignedInteger                       covPeriod                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'covPeriod'                 ]
        ]
        [*, 'COV_RESUBSCRIPTION_INTERVAL', '2'          *COVResubscriptionInterval
            [simple   BACnetApplicationTagUnsignedInteger                       covResubscriptionInterval               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'covResubscriptionInterval' ]
        ]
        [*, 'COVU_PERIOD', '2'                          *COVUPeriod
            [simple   BACnetApplicationTagUnsignedInteger                       covuPeriod                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'covuPeriod'                ]
        ]
        [*, 'COVU_RECIPIENTS'                           *COVURecipients
            [array    BACnetRecipient
                                        covuRecipients
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CREDENTIAL_DISABLE', '9'                   *CredentialDisable
            [simple   BACnetAccessCredentialDisableTagged('0', 'TagClass.APPLICATION_TAGS') credentialDisable           ]
            [virtual  BACnetAccessCredentialDisableTagged                       actualValue 'credentialDisable'         ]
        ]
        [*, 'CREDENTIAL_STATUS', '9'                    *CredentialStatus
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    binaryPv                                ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'binaryPv'                  ]
        ]
        [*, 'CREDENTIALS'                               *Credentials
            [array    BACnetDeviceObjectReference
                        credentials
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'CREDENTIALS_IN_ZONE'                       *CredentialsInZone
            [array    BACnetDeviceObjectReference
                        credentialsInZone
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'CURRENT_COMMAND_PRIORITY'                  *CurrentCommandPriority
            [simple   BACnetOptionalUnsigned                                    currentCommandPriority                  ]
            [virtual  BACnetOptionalUnsigned                                    actualValue 'currentCommandPriority'    ]
        ]
        [*, 'DATABASE_REVISION', '2'                    *DatabaseRevision
            [simple   BACnetApplicationTagUnsignedInteger                       databaseRevision                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'databaseRevision'          ]
        ]
        [*, 'DATE_LIST'                                 *DateList
            [array    BACnetCalendarEntry
                        dateList
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'DAYLIGHT_SAVINGS_STATUS', '1'              *DaylightSavingsStatus
            [simple   BACnetApplicationTagBoolean                               daylightSavingsStatus                   ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'daylightSavingsStatus'     ]
        ]
        [*, 'DAYS_REMAINING', '3'                       *DaysRemaining
            [simple   BACnetApplicationTagSignedInteger                         daysRemaining                           ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue     'daysRemaining'         ]
        ]
        ['INTEGER_VALUE', 'DEADBAND', '2'               *IntegerValueDeadband
            [simple   BACnetApplicationTagUnsignedInteger                       deadband                                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'deadband'                  ]
        ]
        ['LARGE_ANALOG_VALUE', 'DEADBAND', '5'          *LargeAnalogValueDeadband
            [simple   BACnetApplicationTagDouble                                deadband                                ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'deadband'                  ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'DEADBAND', '2'      *PositiveIntegerValueDeadband
            [simple   BACnetApplicationTagUnsignedInteger                       deadband                                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'deadband'                  ]
        ]
        [*, 'DEADBAND', '4'                             *Deadband
            [simple   BACnetApplicationTagReal                                  deadband                                ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'deadband'                  ]
        ]
        [*, 'DEFAULT_FADE_TIME', '2'                    *DefaultFadeTime
            [simple   BACnetApplicationTagUnsignedInteger                       defaultFadeTime                         ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'defaultFadeTime'           ]
        ]
        [*, 'DEFAULT_RAMP_RATE', '4'                    *DefaultRampRate
            [simple   BACnetApplicationTagReal                                  defaultRampRate                         ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'defaultRampRate'           ]
        ]
        [*, 'DEFAULT_STEP_INCREMENT', '4'               *DefaultStepIncrement
            [simple   BACnetApplicationTagReal                                  defaultStepIncrement                    ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'defaultStepIncrement'      ]
        ]
        [*, 'DEFAULT_SUBORDINATE_RELATIONSHIP', '9'     *DefaultSubordinateRelationship
            [simple   BACnetRelationshipTagged('0', 'TagClass.APPLICATION_TAGS') defaultSubordinateRelationship         ]
            [virtual  BACnetRelationshipTagged                                  actualValue 'defaultSubordinateRelationship'    ]
        ]
        [*, 'DEFAULT_TIMEOUT', '2'                      *DefaultTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       defaultTimeout                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'defaultTimeout'            ]
        ]
        [*, 'DEPLOYED_PROFILE_LOCATION', '7'            *DeployedProfileLocation
            [simple   BACnetApplicationTagCharacterString                       deployedProfileLocation                 ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'deployedProfileLocation'   ]
        ]
        [*, 'DERIVATIVE_CONSTANT', '4'                  *DerivativeConstant
            [simple   BACnetApplicationTagReal                                  derivativeConstant                      ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'derivativeConstant'        ]
        ]
        [*, 'DERIVATIVE_CONSTANT_UNITS', '9'            *DerivativeConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')    units                           ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        [*, 'DESCRIPTION', '7'                          *Description
            [simple   BACnetApplicationTagCharacterString                       description                             ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'description'               ]
        ]
        [*, 'DESCRIPTION_OF_HALT', '7'                  *DescriptionOfHalt
            [simple   BACnetApplicationTagCharacterString                       descriptionForHalt                      ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'descriptionForHalt'        ]
        ]
        [*, 'DEVICE_ADDRESS_BINDING'                    *DeviceAddressBinding
            [array    BACnetAddressBinding
                                deviceAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'DEVICE_TYPE', '7'                          *DeviceType
            [simple   BACnetApplicationTagCharacterString                       deviceType                              ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'deviceType'                ]
        ]
        [*, 'DIRECT_READING', '4'                       *DirectReading
            [simple   BACnetApplicationTagReal                                  directReading                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'directReading'             ]
        ]
        [*, 'DISTRIBUTION_KEY_REVISION', '2'            *DistributionKeyRevision
            [simple   BACnetApplicationTagUnsignedInteger                       distributionKeyRevision                 ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'distributionKeyRevision'   ]
        ]
        [*, 'DO_NOT_HIDE', '1'                          *DoNotHide
            [simple   BACnetApplicationTagBoolean                               doNotHide                               ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'doNotHide'                 ]
        ]
        [*, 'DOOR_ALARM_STATE', '9'                     *DoorAlarmState
            [simple   BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS') doorAlarmState                       ]
            [virtual  BACnetDoorAlarmStateTagged                                actualValue 'doorAlarmState'            ]
        ]
        [*, 'DOOR_EXTENDED_PULSE_TIME', '2'             *DoorExtendedPulseTime
            [simple   BACnetApplicationTagUnsignedInteger                       doorExtendedPulseTime                   ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'doorExtendedPulseTime'     ]
        ]
        [*, 'DOOR_MEMBERS'                              *DoorMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectReference
                        doorMembers
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'DOOR_OPEN_TOO_LONG_TIME', '2'              *DoorOpenTooLongTime
            [simple   BACnetApplicationTagUnsignedInteger                       doorOpenTooLongTime                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'doorOpenTooLongTime'       ]
        ]
        [*, 'DOOR_PULSE_TIME', '2'                      *DoorPulseTime
            [simple   BACnetApplicationTagUnsignedInteger                       doorPulseTime                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'doorPulseTime'             ]
        ]
        [*, 'DOOR_STATUS', '9'                          *DoorStatus
            [simple   BACnetDoorStatusTagged('0', 'TagClass.APPLICATION_TAGS')  doorStatus                              ]
            [virtual  BACnetDoorStatusTagged                                    actualValue 'doorStatus'                ]
        ]
        [*, 'DOOR_UNLOCK_DELAY_TIME', '2'               *DoorUnlockDelayTime
            [simple   BACnetApplicationTagUnsignedInteger                       doorUnlockDelayTime                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'doorUnlockDelayTime'       ]
        ]
        [*, 'DUTY_WINDOW', '2'                          *DutyWindow
            [simple   BACnetApplicationTagUnsignedInteger                       dutyWindow                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'dutyWindow'                ]
        ]
        [*, 'EFFECTIVE_PERIOD'                          *EffectivePeriod
            [simple   BACnetDateRange                                           dateRange                               ]
            [virtual  BACnetDateRange                                           actualValue 'dateRange'                 ]
        ]
        [*, 'EGRESS_ACTIVE', '1'                        *EgressActive
            [simple   BACnetApplicationTagBoolean                               egressActive                            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'egressActive'              ]
        ]
        [*, 'EGRESS_TIME', '2'                          *EgressTime
            [simple   BACnetApplicationTagUnsignedInteger                       egressTime                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'egressTime'                ]
        ]
        [*, 'ELAPSED_ACTIVE_TIME', '2'                  *ElapsedActiveTime
            [simple   BACnetApplicationTagUnsignedInteger                       elapsedActiveTime                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'elapsedActiveTime'         ]
        ]
        [*, 'ELEVATOR_GROUP', '12'                      *ElevatorGroup
            [simple   BACnetApplicationTagObjectIdentifier                      elevatorGroup                           ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'elevatorGroup'             ]
        ]
        [*, 'ENABLE', '1'                               *Enable
            [simple   BACnetApplicationTagBoolean                               enable                                  ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'enable'                    ]
        ]
        [*, 'ENERGY_METER', '4'                         *EnergyMeter
            [simple   BACnetApplicationTagReal                                  energyMeter                             ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'energyMeter'               ]
        ]
        [*, 'ENERGY_METER_REF'                          *EnergyMeterRef
            [simple   BACnetDeviceObjectReference                               energyMeterRef                          ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'energyMeterRef'            ]
        ]
        [*, 'ENTRY_POINTS'                              *EntryPoints
            [array    BACnetDeviceObjectReference
                                entryPoints
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ERROR_LIMIT', '4'                          *ErrorLimit
            [simple   BACnetApplicationTagReal                                  errorLimit                              ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'errorLimit'                ]
        ]
        [*, 'ESCALATOR_MODE', '9'                       *EscalatorMode
            [simple   BACnetEscalatorModeTagged('0', 'TagClass.APPLICATION_TAGS')   escalatorMode                       ]
            [virtual  BACnetEscalatorModeTagged                                 actualValue 'escalatorMode'             ]
        ]
        [*, 'EVENT_ALGORITHM_INHIBIT', '1'              *EventAlgorithmInhibit
            [simple   BACnetApplicationTagBoolean                               eventAlgorithmInhibit                   ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'eventAlgorithmInhibit'     ]
        ]
        [*, 'EVENT_ALGORITHM_INHIBIT_REF'               *EventAlgorithmInhibitRef
            [simple   BACnetObjectPropertyReference                             eventAlgorithmInhibitRef                ]
            [virtual  BACnetObjectPropertyReference                             actualValue 'eventAlgorithmInhibitRef'  ]
        ]
        [*, 'EVENT_DETECTION_ENABLE', '1'               *EventDetectionEnable
            [simple   BACnetApplicationTagBoolean                               eventDetectionEnable                    ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'eventDetectionEnable'      ]
        ]
        [*, 'EVENT_ENABLE', '8'                         *EventEnable
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS') eventEnable                     ]
            [virtual  BACnetEventTransitionBitsTagged                           actualValue 'eventEnable'               ]
        ]
        [*, 'EVENT_MESSAGE_TEXTS'                       *EventMessageTexts
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                                        eventMessageTexts
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [virtual  BACnetOptionalCharacterString    toOffnormalText    'COUNT(eventMessageTexts)==3?eventMessageTexts[0]:null'            ]
            [virtual  BACnetOptionalCharacterString    toFaultText        'COUNT(eventMessageTexts)==3?eventMessageTexts[1]:null'            ]
            [virtual  BACnetOptionalCharacterString    toNormalText       'COUNT(eventMessageTexts)==3?eventMessageTexts[2]:null'            ]
            [validation 'arrayIndexArgument!=null || COUNT(eventMessageTexts) == 3'
                                    "eventMessageTexts should have exactly 3 values"                                    ]
        ]
        [*, 'EVENT_MESSAGE_TEXTS_CONFIG'                *EventMessageTextsConfig
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                                        eventMessageTextsConfig
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [virtual  BACnetOptionalCharacterString    toOffnormalTextConfig    'COUNT(eventMessageTextsConfig)==3?eventMessageTextsConfig[0]:null']
            [virtual  BACnetOptionalCharacterString    toFaultTextConfig        'COUNT(eventMessageTextsConfig)==3?eventMessageTextsConfig[1]:null']
            [virtual  BACnetOptionalCharacterString    toNormalTextConfig       'COUNT(eventMessageTextsConfig)==3?eventMessageTextsConfig[2]:null']
            [validation 'arrayIndexArgument!=null || COUNT(eventMessageTextsConfig) == 3'
                        "eventMessageTextsConfig should have exactly 3 values"                                          ]
        ]
        [*, 'EVENT_PARAMETERS'                          *EventParameters
            [simple   BACnetEventParameter                                      eventParameter                          ]
            [virtual  BACnetEventParameter                                      actualValue 'eventParameter'            ]
        ]
        [*, 'EVENT_STATE', '9'                          *EventState
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')  eventState                              ]
            [virtual  BACnetEventStateTagged                                    actualValue 'eventState'                ]
        ]
        [*, 'EVENT_TIME_STAMPS'                         *EventTimeStamps
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetTimeStamp
                                        eventTimeStamps
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [virtual  BACnetTimeStamp    toOffnormal    'COUNT(eventTimeStamps)==3?eventTimeStamps[0]:null'             ]
            [virtual  BACnetTimeStamp    toFault        'COUNT(eventTimeStamps)==3?eventTimeStamps[1]:null'             ]
            [virtual  BACnetTimeStamp    toNormal       'COUNT(eventTimeStamps)==3?eventTimeStamps[2]:null'             ]
            [validation 'arrayIndexArgument!=null || COUNT(eventTimeStamps) == 3'
                        "eventTimeStamps should have exactly 3 values"                                                  ]
        ]
        [*, 'EVENT_TYPE', '9'                           *EventType
            [simple   BACnetEventTypeTagged('0', 'TagClass.APPLICATION_TAGS')   eventType                               ]
            [virtual  BACnetEventTypeTagged                                     actualValue 'eventType'                 ]
        ]
        [*, 'EXCEPTION_SCHEDULE'                        *ExceptionSchedule
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetSpecialEvent
                            exceptionSchedule
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'EXECUTION_DELAY', '2'                      *ExecutionDelay
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                        executionDelay
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'EXIT_POINTS'                               *ExitPoints
            [array    BACnetDeviceObjectReference
                        exitPoints
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'EXPECTED_SHED_LEVEL'                       *ExpectedShedLevel
            [simple   BACnetShedLevel                                           expectedShedLevel                       ]
            [virtual  BACnetShedLevel                                           actualValue 'expectedShedLevel'         ]
        ]
        [*, 'EXPIRATION_TIME'                           *ExpirationTime
            [simple   BACnetDateTime                                            expirationTime                          ]
            [virtual  BACnetDateTime                                            actualValue 'expirationTime'            ]
        ]
        [*, 'EXTENDED_TIME_ENABLE', '1'                 *ExtendedTimeEnable
            [simple   BACnetApplicationTagBoolean                               extendedTimeEnable                      ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'extendedTimeEnable'        ]
        ]
        [*, 'FAILED_ATTEMPT_EVENTS'                     *FailedAttemptEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                        failedAttemptEvents
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAILED_ATTEMPTS', '2'                      *FailedAttempts
            [simple   BACnetApplicationTagUnsignedInteger                       failedAttempts                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'failedAttempts'            ]
        ]
        [*, 'FAILED_ATTEMPTS_TIME', '2'                 *FailedAttemptsTime
            [simple   BACnetApplicationTagUnsignedInteger                       failedAttemptsTime                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'failedAttemptsTime'        ]
        ]
        ['ACCUMULATOR', 'FAULT_HIGH_LIMIT', '2'         *AccumulatorFaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                       faultHighLimit                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'faultHighLimit'            ]
        ]
        ['ANALOG_INPUT', 'FAULT_HIGH_LIMIT', '4'        *AnalogInputFaultHighLimit
            [simple   BACnetApplicationTagReal                                  faultHighLimit                          ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'faultHighLimit'            ]
        ]
        ['ANALOG_VALUE', 'FAULT_HIGH_LIMIT', '4'        *AnalogValueFaultHighLimit
            [simple   BACnetApplicationTagReal                                  faultHighLimit                          ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'faultHighLimit'            ]
        ]
        ['INTEGER_VALUE', 'FAULT_HIGH_LIMIT', '3'       *IntegerValueFaultHighLimit
            [simple   BACnetApplicationTagSignedInteger                         faultHighLimit                          ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'faultHighLimit'            ]
        ]
        ['LARGE_ANALOG_VALUE', 'FAULT_HIGH_LIMIT', '5'  *LargeAnalogValueFaultHighLimit
            [simple   BACnetApplicationTagDouble                                faultHighLimit                          ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'faultHighLimit'            ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'FAULT_HIGH_LIMIT', '2'   *PositiveIntegerValueFaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                       faultHighLimit                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'faultHighLimit'            ]
        ]
        [*, 'FAULT_HIGH_LIMIT', '2'                     *FaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                       faultHighLimit                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'faultHighLimit'            ]
        ]
        ['ACCUMULATOR', 'FAULT_LOW_LIMIT', '2'          *AccumulatorFaultLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                       faultLowLimit                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'faultLowLimit'             ]
        ]
        ['ANALOG_INPUT', 'FAULT_LOW_LIMIT', '4'         *AnalogInputFaultLowLimit
            [simple   BACnetApplicationTagReal                                  faultLowLimit                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'faultLowLimit'             ]
        ]
        ['ANALOG_VALUE', 'FAULT_LOW_LIMIT', '4'         *AnalogValueFaultLowLimit
            [simple   BACnetApplicationTagReal                                  faultLowLimit                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'faultLowLimit'             ]
        ]
        ['LARGE_ANALOG_VALUE', 'FAULT_LOW_LIMIT', '5'   *LargeAnalogValueFaultLowLimit
            [simple   BACnetApplicationTagDouble                                faultLowLimit                           ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'faultLowLimit'             ]
        ]
        ['INTEGER_VALUE', 'FAULT_LOW_LIMIT', '3'        *IntegerValueFaultLowLimit
            [simple   BACnetApplicationTagSignedInteger                         faultLowLimit                           ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'faultLowLimit'             ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'FAULT_LOW_LIMIT', '2'    *PositiveIntegerValueFaultLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                       faultLowLimit                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'faultLowLimit'             ]
        ]
        [*, 'FAULT_LOW_LIMIT', '4'                      *FaultLowLimit
            [simple   BACnetApplicationTagReal                                  faultLowLimit                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'faultLowLimit'             ]
        ]
        [*, 'FAULT_PARAMETERS'                          *FaultParameters
            [simple   BACnetFaultParameter                                      faultParameters                         ]
            [virtual  BACnetFaultParameter                                      actualValue 'faultParameters'           ]
        ]
        ['ESCALATOR', 'FAULT_SIGNALS'                   *EscalatorFaultSignals
            [array   BACnetEscalatorFaultTagged('0', 'TagClass.APPLICATION_TAGS')
                                    faultSignals
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFT', 'FAULT_SIGNALS'                        *LiftFaultSignals
            [array   BACnetLiftFaultTagged('0', 'TagClass.APPLICATION_TAGS')
                                    faultSignals
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAULT_SIGNALS'                             *FaultSignals
            [array   BACnetLiftFaultTagged('0', 'TagClass.APPLICATION_TAGS')
                                    faultSignals
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAULT_TYPE', '9'                           *FaultType
            [simple   BACnetFaultTypeTagged('0', 'TagClass.APPLICATION_TAGS')   faultType                               ]
            [virtual  BACnetFaultTypeTagged                                     actualValue 'faultType'                 ]
        ]
        ['ACCESS_DOOR', 'FAULT_VALUES'        *AccessDoorFaultValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['CHARACTERSTRING_VALUE', 'FAULT_VALUES'        *CharacterStringValueFaultValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_POINT', 'FAULT_VALUES'            *LifeSafetyPointFaultValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_ZONE', 'FAULT_VALUES'             *LifeSafetyZoneFaultValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_INPUT', 'FAULT_VALUES', '2'       *MultiStateInputFaultValues
            [array    BACnetApplicationTagUnsignedInteger
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_VALUE', 'FAULT_VALUES', '2'       *MultiStateValueFaultValues
            [array    BACnetApplicationTagUnsignedInteger
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAULT_VALUES'                              *FaultValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FD_BBMD_ADDRESS'                           *FDBBMDAddress
            [simple   BACnetHostNPort                                           fDBBMDAddress                           ]
            [virtual  BACnetHostNPort                                           actualValue 'fDBBMDAddress'             ]
        ]
        [*, 'FD_SUBSCRIPTION_LIFETIME', '2'             *FDSubscriptionLifetime
            [simple   BACnetApplicationTagUnsignedInteger                       fdSubscriptionLifetime                  ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'fdSubscriptionLifetime'    ]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'FEEDBACK_VALUE', '9'    *BinaryLightingOutputFeedbackValue
            [simple   BACnetBinaryLightingPVTagged('0', 'TagClass.APPLICATION_TAGS')    feedbackValue                   ]
            [virtual  BACnetBinaryLightingPVTagged                              actualValue 'feedbackValue'             ]
        ]
        ['BINARY_OUTPUT', 'FEEDBACK_VALUE', '9'         *BinaryOutputFeedbackValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    feedbackValue                           ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'feedbackValue'             ]
        ]
        ['LIGHTING_OUTPUT', 'FEEDBACK_VALUE', '4'       *LightingOutputFeedbackValue
            [simple   BACnetApplicationTagReal                                  feedbackValue                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'feedbackValue'             ]
        ]
        ['MULTI_STATE_OUTPUT', 'FEEDBACK_VALUE', '2'    *MultiStateOutputFeedbackValue
            [simple   BACnetApplicationTagUnsignedInteger                       feedbackValue                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'feedbackValue'             ]
        ]
        // TODO: similar to ALARM_VALUE either we catch it or we just exlude
        //[*, 'FEEDBACK_VALUE'                          *FeedbackValue [validation    '1 == 2'    "TODO: implement me FEEDBACK_VALUE *FeedbackValue"]]
        [*, 'FILE_ACCESS_METHOD', '9'                   *FileAccessMethod
            [simple   BACnetFileAccessMethodTagged('0', 'TagClass.APPLICATION_TAGS')     fileAccessMethod               ]
            [virtual  BACnetFileAccessMethodTagged                              actualValue 'fileAccessMethod'          ]
        ]
        [*, 'FILE_SIZE', '2'                            *FileSize
            [simple   BACnetApplicationTagUnsignedInteger                       fileSize                                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'fileSize'                  ]
        ]
        [*, 'FILE_TYPE', '7'                            *FileType
            [simple   BACnetApplicationTagCharacterString                       fileType                                ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'fileType'                  ]
        ]
        [*, 'FIRMWARE_REVISION', '7'                    *FirmwareRevision
            [simple   BACnetApplicationTagCharacterString                       firmwareRevision                        ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'firmwareRevision'          ]
        ]
        [*, 'FLOOR_TEXT', '7'                           *FloorText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'FULL_DUTY_BASELINE', '4'                   *FullDutyBaseline
            [simple   BACnetApplicationTagReal                                  fullDutyBaseLine                        ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'fullDutyBaseLine'          ]
        ]
        [*, 'GLOBAL_IDENTIFIER', '2'                    *GlobalIdentifier
            [simple   BACnetApplicationTagUnsignedInteger                       globalIdentifier                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'globalIdentifier'          ]
        ]
        [*, 'GROUP_ID', '2'                             *GroupID
            [simple   BACnetApplicationTagUnsignedInteger                       groupId                                 ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'groupId'                   ]
        ]
        [*, 'GROUP_MEMBER_NAMES', '7'                   *GroupMemberNames
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                            groupMemberNames
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['GLOBAL_GROUP', 'GROUP_MEMBERS'                *GlobalGroupGroupMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectPropertyReference
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['ELEVATOR_GROUP', 'GROUP_MEMBERS', '12'        *ElevatorGroupGroupMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'GROUP_MEMBERS', '12'                       *GroupMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'GROUP_MODE', '9'                           *GroupMode
            [simple   BACnetLiftGroupModeTagged('0', 'TagClass.APPLICATION_TAGS')       groupMode                       ]
            [virtual  BACnetLiftGroupModeTagged                                 actualValue 'groupMode'                 ]
        ]
        ['ACCUMULATOR', 'HIGH_LIMIT', '2'               *AccumulatorHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                       highLimit                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'highLimit'                 ]
        ]
        ['LARGE_ANALOG_VALUE', 'HIGH_LIMIT', '5'        *LargeAnalogValueHighLimit
            [simple   BACnetApplicationTagDouble                                highLimit                               ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'highLimit'                 ]
        ]
        ['INTEGER_VALUE', 'HIGH_LIMIT', '3'             *IntegerValueHighLimit
            [simple   BACnetApplicationTagSignedInteger                         highLimit                               ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'highLimit'                 ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'HIGH_LIMIT', '2'    *PositiveIntegerValueHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                       highLimit                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'highLimit'                 ]
        ]
        [*, 'HIGH_LIMIT', '4'                           *HighLimit
            [simple   BACnetApplicationTagReal                                  highLimit                               ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'highLimit'                 ]
        ]
        [*, 'HIGHER_DECK', '12'                         *HigherDeck
            [simple   BACnetApplicationTagObjectIdentifier                      higherDeck                              ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'higherDeck'                ]
        ]
        [*, 'IN_PROCESS', '1'                           *InProcess
            [simple   BACnetApplicationTagBoolean                               inProcess                               ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'inProcess'                 ]
        ]
        [*, 'IN_PROGRESS', '9'                          *InProgress
            [simple   BACnetLightingInProgressTagged('0', 'TagClass.APPLICATION_TAGS')         inProgress               ]
            [virtual  BACnetLightingInProgressTagged                            actualValue 'inProgress'                ]
        ]
        [*, 'INACTIVE_TEXT', '7'                        *InactiveText
            [simple   BACnetApplicationTagCharacterString                       inactiveText                            ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'inactiveText'              ]
        ]
        [*, 'INITIAL_TIMEOUT', '2'                      *InitialTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       initialTimeout                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'initialTimeout'            ]
        ]
        [*, 'INPUT_REFERENCE'                           *InputReference
            [simple   BACnetObjectPropertyReference                             inputReference                          ]
            [virtual  BACnetObjectPropertyReference                             actualValue 'inputReference'            ]
        ]
        [*, 'INSTALLATION_ID', '2'                      *InstallationID
            [simple   BACnetApplicationTagUnsignedInteger                       installationId                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'installationId'            ]
        ]
        [*, 'INSTANCE_OF', '7'                          *InstanceOf
            [simple   BACnetApplicationTagCharacterString                       instanceOf                              ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'instanceOf'                ]
        ]
        [*, 'INSTANTANEOUS_POWER', '4'                  *InstantaneousPower
            [simple   BACnetApplicationTagReal                                  instantaneousPower                      ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'instantaneousPower'        ]
        ]
        [*, 'INTEGRAL_CONSTANT', '4'                    *IntegralConstant
            [simple   BACnetApplicationTagReal                                  integralConstant                        ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'integralConstant'          ]
        ]
        [*, 'INTEGRAL_CONSTANT_UNITS', '9'              *IntegralConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')            units                   ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        ['ANALOG_INPUT', 'INTERFACE_VALUE'              *AnalogInputInterfaceValue
            [simple   BACnetOptionalREAL                                        interfaceValue                          ]
            [virtual  BACnetOptionalREAL                                        actualValue 'interfaceValue'            ]
        ]
        ['ANALOG_OUTPUT', 'INTERFACE_VALUE'             *AnalogOutputInterfaceValue
            [simple   BACnetOptionalREAL                                        interfaceValue                          ]
            [virtual  BACnetOptionalREAL                                        actualValue 'interfaceValue'            ]
        ]
        ['BINARY_INPUT', 'INTERFACE_VALUE'              *BinaryInputInterfaceValue
            [simple   BACnetOptionalBinaryPV                                    interfaceValue                          ]
            [virtual  BACnetOptionalBinaryPV                                    actualValue 'interfaceValue'            ]
        ]
        ['BINARY_OUTPUT', 'INTERFACE_VALUE'             *BinaryOutputInterfaceValue
            [simple   BACnetOptionalBinaryPV                                    interfaceValue                          ]
            [virtual  BACnetOptionalBinaryPV                                    actualValue 'interfaceValue'            ]
        ]
        ['MULTI_STATE_INPUT', 'INTERFACE_VALUE'         *MultiStateInputInterfaceValue
            [simple   BACnetOptionalBinaryPV                                    interfaceValue                          ]
            [virtual  BACnetOptionalBinaryPV                                    actualValue 'interfaceValue'            ]
        ]
        ['MULTI_STATE_OUTPUT', 'INTERFACE_VALUE'        *MultiStateOutputInterfaceValue
            [simple   BACnetOptionalBinaryPV                                    interfaceValue                          ]
            [virtual  BACnetOptionalBinaryPV                                    actualValue 'interfaceValue'            ]
        ]
        // TODO: unlikely that we have a common type so maybe check that
        //[*, 'INTERFACE_VALUE'                         *InterfaceValue [validation    '1 == 2'    "TODO: implement me INTERFACE_VALUE *InterfaceValue"]]
        [*, 'INTERVAL_OFFSET', '2'                      *IntervalOffset
            [simple   BACnetApplicationTagUnsignedInteger                       intervalOffset                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'intervalOffset'            ]
        ]
        [*, 'IP_ADDRESS', '6'                           *IPAddress
            [simple   BACnetApplicationTagOctetString                           ipAddress                               ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipAddress'                 ]
        ]
        [*, 'IP_DEFAULT_GATEWAY', '6'                   *IPDefaultGateway
            [simple   BACnetApplicationTagOctetString                           ipDefaultGateway                        ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipDefaultGateway'          ]
        ]
        [*, 'IP_DHCP_ENABLE', '1'                       *IPDHCPEnable
            [simple   BACnetApplicationTagBoolean                               ipDhcpEnable                            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'ipDhcpEnable'              ]
        ]
        [*, 'IP_DHCP_LEASE_TIME', '2'                   *IPDHCPLeaseTime
            [simple   BACnetApplicationTagUnsignedInteger                       ipDhcpLeaseTime                         ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipDhcpLeaseTime'           ]
        ]
        [*, 'IP_DHCP_LEASE_TIME_REMAINING', '2'         *IPDHCPLeaseTimeRemaining
            [simple   BACnetApplicationTagUnsignedInteger                       ipDhcpLeaseTimeRemaining                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipDhcpLeaseTimeRemaining'  ]
        ]
        [*, 'IP_DHCP_SERVER', '6'                       *IPDHCPServer
            [simple   BACnetApplicationTagOctetString                           dhcpServer                              ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'dhcpServer'                ]
        ]
        [*, 'IP_DNS_SERVER', '6'                        *IPDNSServer
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagOctetString
                                        ipDnsServer
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'IP_SUBNET_MASK', '6'                       *IPSubnetMask
            [simple   BACnetApplicationTagOctetString                           ipSubnetMask                            ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipSubnetMask'              ]
        ]
        [*, 'IPV6_ADDRESS', '6'                         *IPv6Address
            [simple   BACnetApplicationTagOctetString                           ipv6Address                             ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipv6Address'               ]
        ]
        [*, 'IPV6_AUTO_ADDRESSING_ENABLE', '1'          *IPv6AutoAddressingEnable
            [simple   BACnetApplicationTagBoolean                               autoAddressingEnable                    ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'autoAddressingEnable'      ]
        ]
        [*, 'IPV6_DEFAULT_GATEWAY', '6'                 *IPv6DefaultGateway
            [simple   BACnetApplicationTagOctetString                           ipv6DefaultGateway                      ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'ipv6DefaultGateway'        ]
        ]
        [*, 'IPV6_DHCP_LEASE_TIME', '2'                 *IPv6DHCPLeaseTime
            [simple   BACnetApplicationTagUnsignedInteger                       ipv6DhcpLeaseTime                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipv6DhcpLeaseTime'         ]
        ]
        [*, 'IPV6_DHCP_LEASE_TIME_REMAINING', '2'       *IPv6DHCPLeaseTimeRemaining
            [simple   BACnetApplicationTagUnsignedInteger                       ipv6DhcpLeaseTimeRemaining              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipv6DhcpLeaseTimeRemaining']
        ]
        [*, 'IPV6_DHCP_SERVER', '6'                     *IPv6DHCPServer
            [simple   BACnetApplicationTagOctetString                           dhcpServer                              ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'dhcpServer'                ]
        ]
        [*, 'IPV6_DNS_SERVER', '6'                      *IPv6DNSServer
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagOctetString
                                        ipv6DnsServer
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'IPV6_PREFIX_LENGTH', '2'                   *IPv6PrefixLength
            [simple   BACnetApplicationTagUnsignedInteger                       ipv6PrefixLength                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'ipv6PrefixLength'          ]
        ]
        [*, 'IPV6_ZONE_INDEX', '7'                      *IPv6ZoneIndex
            [simple   BACnetApplicationTagCharacterString                       ipv6ZoneIndex                           ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'ipv6ZoneIndex'             ]
        ]
        [*, 'IS_UTC', '1'                               *IsUTC
            [simple   BACnetApplicationTagBoolean                               isUtc                                   ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'isUtc'                     ]
        ]
        [*, 'KEY_SETS'                                  *KeySets
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetSecurityKeySet
                                keySets
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [validation 'arrayIndexArgument!=null || COUNT(keySets) == 2' "keySets should have exactly 2 values"        ]
        ]
        [*, 'LANDING_CALL_CONTROL'                      *LandingCallControl
            [simple   BACnetLandingCallStatus                                   landingCallControl                      ]
            [virtual  BACnetLandingCallStatus                                   actualValue 'landingCallControl'        ]
        ]
        [*, 'LANDING_CALLS'                             *LandingCalls
            [array    BACnetLandingCallStatus
                                landingCallStatus
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LANDING_DOOR_STATUS'                       *LandingDoorStatus
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLandingDoorStatus
                                landingDoorStatus
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LAST_ACCESS_EVENT', '9'                    *LastAccessEvent
            [simple   BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS') lastAccessEvent                         ]
            [virtual  BACnetAccessEventTagged                                   actualValue 'lastAccessEvent'           ]
        ]
        [*, 'LAST_ACCESS_POINT'                         *LastAccessPoint
            [simple   BACnetDeviceObjectReference                               lastAccessPoint                         ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'lastAccessPoint'           ]
        ]
        [*, 'LAST_COMMAND_TIME'                         *LastCommandTime
            [simple   BACnetTimeStamp                                           lastCommandTime                         ]
            [virtual  BACnetTimeStamp                                           actualValue 'lastCommandTime'           ]
        ]
        [*, 'LAST_CREDENTIAL_ADDED'                     *LastCredentialAdded
            [simple   BACnetDeviceObjectReference                               lastCredentialAdded                     ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'lastCredentialAdded'       ]
        ]
        [*, 'LAST_CREDENTIAL_ADDED_TIME'                *LastCredentialAddedTime
            [simple   BACnetDateTime                                            lastCredentialAddedTime                 ]
            [virtual  BACnetDateTime                                            actualValue 'lastCredentialAddedTime'   ]
        ]
        [*, 'LAST_CREDENTIAL_REMOVED'                   *LastCredentialRemoved
            [simple   BACnetDeviceObjectReference                               lastCredentialRemoved                   ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'lastCredentialRemoved'     ]
        ]
        [*, 'LAST_CREDENTIAL_REMOVED_TIME'              *LastCredentialRemovedTime
            [simple   BACnetDateTime                                            lastCredentialRemovedTime               ]
            [virtual  BACnetDateTime                                            actualValue 'lastCredentialRemovedTime' ]
        ]
        [*, 'LAST_KEY_SERVER'                           *LastKeyServer
            [simple   BACnetAddressBinding                                      lastKeyServer                           ]
            [virtual  BACnetAddressBinding                                      actualValue 'lastKeyServer'             ]
        ]
        [*, 'LAST_NOTIFY_RECORD', '2'                   *LastNotifyRecord
            [simple   BACnetApplicationTagUnsignedInteger                       lastNotifyRecord                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'lastNotifyRecord'          ]
        ]
        [*, 'LAST_PRIORITY', '2'                        *LastPriority
            [simple   BACnetApplicationTagUnsignedInteger                       lastPriority                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'lastPriority'              ]
        ]
        [*, 'LAST_RESTART_REASON', '9'                  *LastRestartReason
            [simple   BACnetRestartReasonTagged('0', 'TagClass.APPLICATION_TAGS')   lastRestartReason                   ]
            [virtual  BACnetRestartReasonTagged                                 actualValue 'lastRestartReason'         ]
        ]
        [*, 'LAST_RESTORE_TIME'                         *LastRestoreTime
            [simple   BACnetTimeStamp                                           lastRestoreTime                         ]
            [virtual  BACnetTimeStamp                                           actualValue 'lastRestoreTime'           ]
        ]
        [*, 'LAST_STATE_CHANGE', '9'                    *LastStateChange
            [simple   BACnetTimerTransitionTagged('0', 'TagClass.APPLICATION_TAGS')             lastStateChange         ]
            [virtual  BACnetTimerTransitionTagged                               actualValue 'lastStateChange'           ]
        ]
        [*, 'LAST_USE_TIME'                             *LastUseTime
            [simple   BACnetDateTime                                            lastUseTime                             ]
            [virtual  BACnetDateTime                                            actualValue 'lastUseTime'               ]
        ]
        [*, 'LIFE_SAFETY_ALARM_VALUES', '9'             *LifeSafetyAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LIGHTING_COMMAND'                          *LightingCommand
            [simple   BACnetLightingCommand                                     lightingCommand                         ]
            [virtual  BACnetLightingCommand                                     actualValue 'lightingCommand'           ]
        ]
        [*, 'LIGHTING_COMMAND_DEFAULT_PRIORITY', '2'    *LightingCommandDefaultPriority
            [simple   BACnetApplicationTagUnsignedInteger                       lightingCommandDefaultPriority          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'lightingCommandDefaultPriority'           ]
        ]
        [*, 'LIMIT_ENABLE','8'                          *LimitEnable
            [simple   BACnetLimitEnableTagged('0', 'TagClass.APPLICATION_TAGS') limitEnable                             ]
            [virtual  BACnetLimitEnableTagged                                   actualValue 'limitEnable'               ]
        ]
        [*, 'LIMIT_MONITORING_INTERVAL', '2'            *LimitMonitoringInterval
            [simple   BACnetApplicationTagUnsignedInteger                       limitMonitoringInterval                 ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'limitMonitoringInterval'   ]
        ]
        [*, 'LINK_SPEED', '4'                           *LinkSpeed
            [simple   BACnetApplicationTagReal                                  linkSpeed                               ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'linkSpeed'                 ]
        ]
        [*, 'LINK_SPEED_AUTONEGOTIATE', '1'             *LinkSpeedAutonegotiate
            [simple   BACnetApplicationTagBoolean                               linkSpeedAutonegotiate                  ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'linkSpeedAutonegotiate'    ]
        ]
        [*, 'LINK_SPEEDS', '4'                          *LinkSpeeds
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagReal
                                        linkSpeeds
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'LIST_OF_GROUP_MEMBERS'                     *ListOfGroupMembers
            [array    BACnetReadAccessSpecification
                                        listOfGroupMembers
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        ['CHANNEL', 'LIST_OF_OBJECT_PROPERTY_REFERENCES'    *ChannelListOfObjectPropertyReferences
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectPropertyReference
                            references              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LIST_OF_OBJECT_PROPERTY_REFERENCES'        *ListOfObjectPropertyReferences
            [array    BACnetDeviceObjectPropertyReference
                            references              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LOCAL_DATE', '10'                          *LocalDate
            [simple   BACnetApplicationTagDate                                  localDate                               ]
            [virtual  BACnetApplicationTagDate                                  actualValue 'localDate'                 ]
        ]
        [*, 'LOCAL_FORWARDING_ONLY', '1'                *LocalForwardingOnly
            [simple   BACnetApplicationTagBoolean                               localForwardingOnly                     ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'localForwardingOnly'       ]
        ]
        [*, 'LOCAL_TIME', '11'                          *LocalTime
            [simple   BACnetApplicationTagTime                                  localTime                               ]
            [virtual  BACnetApplicationTagTime                                  actualValue 'localTime'                 ]
        ]
        [*, 'LOCATION', '7'                             *Location
            [simple   BACnetApplicationTagCharacterString                       location                                ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'location'                  ]
        ]
        [*, 'LOCK_STATUS', '9'                          *LockStatus
            [simple   BACnetLockStatusTagged('0', 'TagClass.APPLICATION_TAGS')  lockStatus                              ]
            [virtual  BACnetLockStatusTagged                                    actualValue 'lockStatus'                ]
        ]
        [*, 'LOCKOUT', '1'                              *Lockout
            [simple   BACnetApplicationTagBoolean                               lockout                                 ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'lockout'                   ]
        ]
        [*, 'LOCKOUT_RELINQUISH_TIME', '2'              *LockoutRelinquishTime
            [simple   BACnetApplicationTagUnsignedInteger                       lockoutRelinquishTime                   ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'lockoutRelinquishTime'     ]
        ]
        ['EVENT_LOG', 'LOG_BUFFER'                      *EventLogLogBuffer
            [array    BACnetEventLogRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['TREND_LOG', 'LOG_BUFFER'                      *TrendLogLogBuffer
            [array    BACnetLogRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['TREND_LOG_MULTIPLE', 'LOG_BUFFER'             *TrendLogMultipleLogBuffer
            [array    BACnetLogMultipleRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'LOG_BUFFER'                                *LogBuffer
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLogRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['TREND_LOG', 'LOG_DEVICE_OBJECT_PROPERTY'      *TrendLogLogDeviceObjectProperty
            [simple   BACnetDeviceObjectPropertyReference                       logDeviceObjectProperty                 ]
            [virtual  BACnetDeviceObjectPropertyReference                       actualValue 'logDeviceObjectProperty'   ]
        ]
        ['TREND_LOG_MULTIPLE', 'LOG_DEVICE_OBJECT_PROPERTY' *TrendLogMultipleLogDeviceObjectProperty
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectPropertyReference
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'LOG_DEVICE_OBJECT_PROPERTY'                *LogDeviceObjectProperty
            [simple   BACnetDeviceObjectPropertyReference                       logDeviceObjectProperty                 ]
            [virtual  BACnetDeviceObjectPropertyReference                       actualValue 'logDeviceObjectProperty'   ]
        ]
        [*, 'LOG_INTERVAL', '2'                         *LogInterval
            [simple   BACnetApplicationTagUnsignedInteger                       logInterval                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'logInterval'               ]
        ]
        [*, 'LOGGING_OBJECT', '12'                      *LoggingObject
            [simple   BACnetApplicationTagObjectIdentifier                      loggingObject                           ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'loggingObject'             ]
        ]
        [*, 'LOGGING_RECORD'                            *LoggingRecord
            [simple   BACnetAccumulatorRecord                                   loggingRecord                           ]
            [virtual  BACnetAccumulatorRecord                                   actualValue 'loggingRecord'             ]
        ]
        [*, 'LOGGING_TYPE', '9'                         *LoggingType
            [simple   BACnetLoggingTypeTagged('0', 'TagClass.APPLICATION_TAGS') loggingType                             ]
            [virtual  BACnetLoggingTypeTagged                                   actualValue 'loggingType'               ]
        ]
        [*, 'LOW_DIFF_LIMIT'                            *LowDiffLimit
            [simple   BACnetOptionalREAL                                        lowDiffLimit                            ]
            [virtual  BACnetOptionalREAL                                        actualValue 'lowDiffLimit'              ]
        ]
        ['ACCUMULATOR', 'LOW_LIMIT', '2'                *AccumulatorLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                       lowLimit                                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'lowLimit'                  ]
        ]
        ['LARGE_ANALOG_VALUE', 'LOW_LIMIT', '5'         *LargeAnalogValueLowLimit
            [simple   BACnetApplicationTagDouble                                lowLimit                                ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'lowLimit'                  ]
        ]
        ['INTEGER_VALUE', 'LOW_LIMIT', '3'              *IntegerValueLowLimit
            [simple   BACnetApplicationTagSignedInteger                         lowLimit                                ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'lowLimit'                  ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'LOW_LIMIT', '2'     *PositiveIntegerValueLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                       lowLimit                                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'lowLimit'                  ]
        ]
        [*, 'LOW_LIMIT', '4'                            *LowLimit
            [simple   BACnetApplicationTagReal                                  lowLimit                                ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'lowLimit'                  ]
        ]
        [*, 'LOWER_DECK', '12'                          *LowerDeck
            [simple   BACnetApplicationTagObjectIdentifier                      lowerDeck                               ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'lowerDeck'                 ]
        ]
        [*, 'MAC_ADDRESS', '6'                          *MACAddress
            [simple   BACnetApplicationTagOctetString                           macAddress                              ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'macAddress'                ]
        ]
        [*, 'MACHINE_ROOM_ID', '12'                     *MachineRoomID
            [simple   BACnetApplicationTagObjectIdentifier                      machineRoomId                           ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'machineRoomId'             ]
        ]
        ['LIFE_SAFETY_ZONE', 'MAINTENANCE_REQUIRED', '1' *LifeSafetyZoneMaintenanceRequired
            [simple   BACnetApplicationTagBoolean                               maintenanceRequired                     ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'maintenanceRequired'       ]
        ]
        [*, 'MAINTENANCE_REQUIRED', '9'                 *MaintenanceRequired
            [simple   BACnetMaintenanceTagged('0', 'TagClass.APPLICATION_TAGS') maintenanceRequired                     ]
            [virtual  BACnetMaintenanceTagged                                   actualValue 'maintenanceRequired'       ]
        ]
        [*, 'MAKING_CAR_CALL'                           *MakingCarCall
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                        makingCarCall
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'MANIPULATED_VARIABLE_REFERENCE'            *ManipulatedVariableReference
            [simple   BACnetObjectPropertyReference                             manipulatedVariableReference            ]
            [virtual  BACnetObjectPropertyReference                             actualValue 'manipulatedVariableReference' ]
        ]
        [*, 'MANUAL_SLAVE_ADDRESS_BINDING'              *ManualSlaveAddressBinding
            [array    BACnetAddressBinding
                                manualSlaveAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MASKED_ALARM_VALUES', '9'                  *MaskedAlarmValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                                maskedAlarmValues
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MAX_ACTUAL_VALUE', '4'                     *MaxActualValue
            [simple   BACnetApplicationTagReal                                  maxActualValue                          ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maxActualValue'            ]
        ]
        [*, 'MAX_APDU_LENGTH_ACCEPTED', '2'             *MaxAPDULengthAccepted
            [simple   BACnetApplicationTagUnsignedInteger                       maxApduLengthAccepted                   ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxApduLengthAccepted'     ]
        ]
        [*, 'MAX_FAILED_ATTEMPTS', '2'                  *MaxFailedAttempts
            [simple   BACnetApplicationTagUnsignedInteger                       maxFailedAttempts                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxFailedAttempts'         ]
        ]
        ['DEVICE', 'MAX_INFO_FRAMES', '2'               *DeviceMaxInfoFrames
            [simple   BACnetApplicationTagUnsignedInteger                       maxInfoFrames                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxInfoFrames'             ]
        ]
        ['NETWORK_PORT', 'MAX_INFO_FRAMES', '2'         *NetworkPortMaxInfoFrames
            [simple   BACnetApplicationTagUnsignedInteger                       maxInfoFrames                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxInfoFrames'             ]
        ]
        [*, 'MAX_INFO_FRAMES', '2'                      *MaxInfoFrames
            [simple   BACnetApplicationTagUnsignedInteger                       maxInfoFrames                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxInfoFrames'             ]
        ]
        ['DEVICE', 'MAX_MASTER', '2'                    *DeviceMaxMaster
            [simple   BACnetApplicationTagUnsignedInteger                       maxMaster                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxMaster'                 ]
        ]
        ['NETWORK_PORT', 'MAX_MASTER', '2'              *NetworkPortMaxMaster
            [simple   BACnetApplicationTagUnsignedInteger                       maxMaster                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxMaster'                 ]
        ]
        [*, 'MAX_MASTER', '2'                           *MaxMaster
            [simple   BACnetApplicationTagUnsignedInteger                       maxMaster                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxMaster'                 ]
        ]
        ['ACCUMULATOR', 'MAX_PRES_VALUE', '2'           *AccumulatorMaxPresValue
            [simple   BACnetApplicationTagUnsignedInteger                       maxPresValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxPresValue'              ]
        ]
        ['ANALOG_INPUT', 'MAX_PRES_VALUE', '4'          *AnalogInputMaxPresValue
            [simple   BACnetApplicationTagReal                                  maxPresValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maxPresValue'              ]
        ]
        ['ANALOG_OUTPUT', 'MAX_PRES_VALUE', '4'         *AnalogOutputMaxPresValue
            [simple   BACnetApplicationTagReal                                  maxPresValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maxPresValue'              ]
        ]
        ['ANALOG_VALUE', 'MAX_PRES_VALUE', '4'          *AnalogValueMaxPresValue
            [simple   BACnetApplicationTagReal                                  maxPresValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maxPresValue'              ]
        ]
        ['LARGE_ANALOG_VALUE', 'MAX_PRES_VALUE', '5'    *LargeAnalogValueMaxPresValue
            [simple   BACnetApplicationTagDouble                                maxPresValue                            ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'maxPresValue'              ]
        ]
        ['INTEGER_VALUE', 'MAX_PRES_VALUE', '3'         *IntegerValueMaxPresValue
            [simple   BACnetApplicationTagSignedInteger                         maxPresValue                            ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'maxPresValue'              ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'MAX_PRES_VALUE', '2' *PositiveIntegerValueMaxPresValue
            [simple   BACnetApplicationTagUnsignedInteger                       maxPresValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxPresValue'              ]
        ]
        ['TIMER', 'MAX_PRES_VALUE', '2'                 *TimerMaxPresValue
            [simple   BACnetApplicationTagUnsignedInteger                       maxPresValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxPresValue'              ]
        ]
        [*, 'MAX_PRES_VALUE', '4'                       *MaxPresValue
            [simple   BACnetApplicationTagReal                                  maxPresValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maxPresValue'              ]
        ]
        [*, 'MAX_SEGMENTS_ACCEPTED', '2'                *MaxSegmentsAccepted
            [simple   BACnetApplicationTagUnsignedInteger                       maxSegmentsAccepted                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'maxSegmentsAccepted'       ]
        ]
        [*, 'MAXIMUM_OUTPUT', '4'                       *MaximumOutput
            [simple   BACnetApplicationTagReal                                  maximumOutput                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maximumOutput'             ]
        ]
        [*, 'MAXIMUM_VALUE', '4'                        *MaximumValue
            [simple   BACnetApplicationTagReal                                  maximumValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'maximumValue'              ]
        ]
        [*, 'MAXIMUM_VALUE_TIMESTAMP'                   *MaximumValueTimestamp
            [simple   BACnetDateTime                                            maximumValueTimestamp                   ]
            [virtual  BACnetDateTime                                            actualValue 'maximumValueTimestamp'     ]
        ]
        [*, 'MEMBER_OF' *MemberOf
            [array    BACnetDeviceObjectReference
                    zones
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MEMBER_STATUS_FLAGS', '8'                  *MemberStatusFlags
            [simple   BACnetStatusFlagsTagged('0', 'TagClass.APPLICATION_TAGS') statusFlags                             ]
            [virtual  BACnetStatusFlagsTagged                                   actualValue 'statusFlags'               ]
        ]
        [*, 'MEMBERS'                                   *Members
            [array    BACnetDeviceObjectReference
                                    members
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'MIN_ACTUAL_VALUE', '4'                     *MinActualValue
            [simple   BACnetApplicationTagReal                                  minActualValue                          ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'minActualValue'            ]
        ]
        ['ACCUMULATOR', 'MIN_PRES_VALUE', '2'           *AccumulatorMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                       minPresValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'minPresValue'              ]
        ]
        ['INTEGER_VALUE', 'MIN_PRES_VALUE', '3'         *IntegerValueMinPresValue
            [simple   BACnetApplicationTagSignedInteger                         minPresValue                            ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'minPresValue'              ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'MIN_PRES_VALUE', '2'    *PositiveIntegerValueMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                       minPresValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'minPresValue'              ]
        ]
        ['LARGE_ANALOG_VALUE', 'MIN_PRES_VALUE', '5'    *LargeAnalogValueMinPresValue
            [simple   BACnetApplicationTagDouble                                minPresValue                            ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'minPresValue'              ]
        ]
        ['TIMER', 'MIN_PRES_VALUE', '2'                 *TimerMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                       minPresValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'minPresValue'              ]
        ]
        [*, 'MIN_PRES_VALUE', '4'                       *MinPresValue
            [simple   BACnetApplicationTagReal                                  minPresValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'minPresValue'              ]
        ]
        [*, 'MINIMUM_OFF_TIME', '2'                     *MinimumOffTime
            [simple   BACnetApplicationTagUnsignedInteger                       minimumOffTime                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'minimumOffTime'            ]
        ]
        [*, 'MINIMUM_ON_TIME', '2'                      *MinimumOnTime
            [simple   BACnetApplicationTagUnsignedInteger                       minimumOnTime                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'minimumOnTime'             ]
        ]
        [*, 'MINIMUM_OUTPUT', '4'                       *MinimumOutput
            [simple   BACnetApplicationTagReal                                  minimumOutput                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'minimumOutput'             ]
        ]
        [*, 'MINIMUM_VALUE', '4'                        *MinimumValue
            [simple   BACnetApplicationTagReal                                  minimumValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'minimumValue'              ]
        ]
        [*, 'MINIMUM_VALUE_TIMESTAMP'                   *MinimumValueTimestamp
            [simple   BACnetDateTime                                            minimumValueTimestamp                   ]
            [virtual  BACnetDateTime                                            actualValue 'minimumValueTimestamp'     ]
        ]
        [*, 'MODE', '9'                                 *Mode
            [simple   BACnetLifeSafetyModeTagged('0', 'TagClass.APPLICATION_TAGS')              mode                    ]
            [virtual  BACnetLifeSafetyModeTagged                                actualValue 'mode'                      ]
        ]
        [*, 'MODEL_NAME', '7'                           *ModelName
            [simple   BACnetApplicationTagCharacterString                       modelName                               ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'modelName'                 ]
        ]
        [*, 'MODIFICATION_DATE'                         *ModificationDate
            [simple   BACnetDateTime                                            modificationDate                        ]
            [virtual  BACnetDateTime                                            actualValue 'modificationDate'          ]
        ]
        [*, 'MUSTER_POINT', '1'                         *MusterPoint
            [simple   BACnetApplicationTagBoolean                               musterPoint                             ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'musterPoint'               ]
        ]
        [*, 'NEGATIVE_ACCESS_RULES'                     *NegativeAccessRules
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAccessRule
                            negativeAccessRules
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'NETWORK_ACCESS_SECURITY_POLICIES'          *NetworkAccessSecurityPolicies
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNetworkSecurityPolicy
                            networkAccessSecurityPolicies
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'NETWORK_INTERFACE_NAME', '7'               *NetworkInterfaceName
            [simple   BACnetApplicationTagCharacterString                       networkInterfaceName                    ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'networkInterfaceName'      ]
        ]
        [*, 'NETWORK_NUMBER', '2'                       *NetworkNumber
            [simple   BACnetApplicationTagUnsignedInteger                       networkNumber                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'networkNumber'             ]
        ]
        [*, 'NETWORK_NUMBER_QUALITY', '9'               *NetworkNumberQuality
            [simple   BACnetNetworkNumberQualityTagged('0', 'TagClass.APPLICATION_TAGS')    networkNumberQuality        ]
            [virtual  BACnetNetworkNumberQualityTagged                          actualValue 'networkNumberQuality'      ]
        ]
        [*, 'NETWORK_TYPE', '9'                         *NetworkType
            [simple   BACnetNetworkTypeTagged('0', 'TagClass.APPLICATION_TAGS') networkType                             ]
            [virtual  BACnetNetworkTypeTagged                                   actualValue 'networkType'               ]
        ]
        [*, 'NEXT_STOPPING_FLOOR', '2'                  *NextStoppingFloor
            [simple   BACnetApplicationTagUnsignedInteger                       nextStoppingFloor                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'nextStoppingFloor'         ]
        ]
        [*, 'NODE_SUBTYPE', '7'                         *NodeSubtype
            [simple   BACnetApplicationTagCharacterString                       nodeSubType                             ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'nodeSubType'               ]
        ]
        [*, 'NODE_TYPE', '9'                            *NodeType
            [simple   BACnetNodeTypeTagged('0', 'TagClass.APPLICATION_TAGS')    nodeType                                ]
            [virtual  BACnetNodeTypeTagged                                      actualValue 'nodeType'                  ]
        ]
        [*, 'NOTIFICATION_CLASS', '2'                   *NotificationClass
            [simple   BACnetApplicationTagUnsignedInteger                       notificationClass                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'notificationClass'         ]
        ]
        [*, 'NOTIFICATION_THRESHOLD', '2'               *NotificationThreshold
            [simple   BACnetApplicationTagUnsignedInteger                       notificationThreshold                   ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'notificationThreshold'     ]
        ]
        [*, 'NOTIFY_TYPE', '9'                          *NotifyType
            [simple   BACnetNotifyTypeTagged('0', 'TagClass.APPLICATION_TAGS')  notifyType                              ]
            [virtual  BACnetNotifyTypeTagged                                    actualValue 'notifyType'                ]
        ]
        [*, 'NUMBER_OF_APDU_RETRIES', '2'               *NumberOfAPDURetries
            [simple   BACnetApplicationTagUnsignedInteger                       numberOfApduRetries                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'numberOfApduRetries'       ]
        ]
        [*, 'NUMBER_OF_AUTHENTICATION_POLICIES', '2'    *NumberOfAuthenticationPolicies
            [simple   BACnetApplicationTagUnsignedInteger                       numberOfAuthenticationPolicies          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'numberOfAuthenticationPolicies']
        ]
        [*, 'NUMBER_OF_STATES', '2'                     *NumberOfStates
            [simple   BACnetApplicationTagUnsignedInteger                       numberOfState                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'numberOfState'             ]
        ]
        [*, 'OBJECT_IDENTIFIER', '12'                   *ObjectIdentifier
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                        ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'objectIdentifier'          ]
        ]
        [*, 'OBJECT_LIST'                               *ObjectList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                                objectList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'OBJECT_NAME', '7'                          *ObjectName
            [simple   BACnetApplicationTagCharacterString                       objectName                              ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'objectName'                ]
        ]
        [*, 'OBJECT_PROPERTY_REFERENCE'                 *ObjectPropertyReference
            [simple   BACnetDeviceObjectPropertyReference                       propertyReference                       ]
            [virtual  BACnetDeviceObjectPropertyReference                       actualValue 'propertyReference'         ]
        ]
        [*, 'OBJECT_TYPE', '9'                          *ObjectType
            [simple   BACnetObjectTypeTagged('0', 'TagClass.APPLICATION_TAGS')  objectType                              ]
            [virtual  BACnetObjectTypeTagged                                    actualValue 'objectType'                ]
        ]
        [*, 'OCCUPANCY_COUNT', '2'                      *OccupancyCount
            [simple   BACnetApplicationTagUnsignedInteger                       occupancyCount                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'occupancyCount'            ]
        ]
        [*, 'OCCUPANCY_COUNT_ADJUST', '1'               *OccupancyCountAdjust
            [simple   BACnetApplicationTagBoolean                               occupancyCountAdjust                    ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'occupancyCountAdjust'      ]
        ]
        [*, 'OCCUPANCY_COUNT_ENABLE', '1'               *OccupancyCountEnable
            [simple   BACnetApplicationTagBoolean                               occupancyCountEnable                    ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'occupancyCountEnable'      ]
        ]
        [*, 'OCCUPANCY_LOWER_LIMIT', '2'                *OccupancyLowerLimit
            [simple   BACnetApplicationTagUnsignedInteger                       occupancyLowerLimit                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'occupancyLowerLimit'       ]
        ]
        [*, 'OCCUPANCY_LOWER_LIMIT_ENFORCED', '1'       *OccupancyLowerLimitEnforced
            [simple   BACnetApplicationTagBoolean                               occupancyLowerLimitEnforced             ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'occupancyLowerLimitEnforced']
        ]
        [*, 'OCCUPANCY_STATE', '9'                      *OccupancyState
            [simple   BACnetAccessZoneOccupancyStateTagged('0', 'TagClass.APPLICATION_TAGS') occupancyState             ]
            [virtual  BACnetAccessZoneOccupancyStateTagged                      actualValue 'occupancyState'            ]
        ]
        [*, 'OCCUPANCY_UPPER_LIMIT', '2'                *OccupancyUpperLimit
            [simple   BACnetApplicationTagUnsignedInteger                       occupancyUpperLimit                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'occupancyUpperLimit'       ]
        ]
        [*, 'OCCUPANCY_UPPER_LIMIT_ENFORCED', '1'       *OccupancyUpperLimitEnforced
            [simple   BACnetApplicationTagBoolean                               occupancyUpperLimitEnforced             ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'occupancyUpperLimitEnforced']
        ]
        [*, 'OPERATION_DIRECTION', '9'                  *OperationDirection
            [simple   BACnetEscalatorOperationDirectionTagged('0', 'TagClass.APPLICATION_TAGS')   operationDirection    ]
            [virtual  BACnetEscalatorOperationDirectionTagged                   actualValue 'operationDirection'        ]
        ]
        [*, 'OPERATION_EXPECTED', '9'                   *OperationExpected
            [simple   BACnetLifeSafetyOperationTagged('0', 'TagClass.APPLICATION_TAGS')         lifeSafetyOperations    ]
            [virtual  BACnetLifeSafetyOperationTagged                           actualValue 'lifeSafetyOperations'      ]
        ]
        [*, 'OPTIONAL'                                  *Optional
            [validation    '1 == 2'    "An property identified by OPTIONAL should never occur in the wild"]
        ]
        [*, 'OUT_OF_SERVICE', '1'                       *OutOfService
            [simple   BACnetApplicationTagBoolean                               outOfService                            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'outOfService'              ]
        ]
        [*, 'OUTPUT_UNITS', '9'                         *OutputUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        [*, 'PACKET_REORDER_TIME', '2'                  *PacketReorderTime
            [simple   BACnetApplicationTagUnsignedInteger                       packetReorderTime                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'packetReorderTime'         ]
        ]
        [*, 'PASSBACK_MODE'                             *PassbackMode
            [simple   BACnetAccessPassbackModeTagged('0', 'TagClass.APPLICATION_TAGS')             passbackMode         ]
            [virtual  BACnetAccessPassbackModeTagged                            actualValue 'passbackMode'              ]
        ]
        [*, 'PASSBACK_TIMEOUT', '2'                     *PassbackTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       passbackTimeout                         ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'passbackTimeout'           ]
        ]
        [*, 'PASSENGER_ALARM', '1'                      *PassengerAlarm
            [simple   BACnetApplicationTagBoolean                               passengerAlarm                          ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'passengerAlarm'            ]
        ]
        [*, 'POLARITY', '9'                             *Polarity
            [simple   BACnetPolarityTagged('0', 'TagClass.APPLICATION_TAGS')    polarity                                ]
            [virtual  BACnetPolarityTagged                                      actualValue 'polarity'                  ]
        ]
        [*, 'PORT_FILTER'                               *PortFilter
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetPortPermission
                            portFilter
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'POSITIVE_ACCESS_RULES'                     *PositiveAccessRules
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAccessRule
                            positiveAccessRules
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'POWER', '4'                                *Power
            [simple   BACnetApplicationTagReal                                  power                                   ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'power'                     ]
        ]
        [*, 'POWER_MODE', '1'                           *PowerMode
            [simple   BACnetApplicationTagBoolean                               powerMode                               ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'powerMode'                 ]
        ]
        [*, 'PRESCALE'                                  *Prescale
            [simple   BACnetPrescale                                            prescale                                ]
            [virtual  BACnetPrescale                                            actualValue 'prescale'                  ]
        ]
        ['ACCESS_DOOR', 'PRESENT_VALUE', '9'            *AccessDoorPresentValue
            [simple   BACnetDoorValueTagged('0', 'TagClass.APPLICATION_TAGS')   presentValue                            ]
            [virtual  BACnetDoorValueTagged                                     actualValue 'presentValue'              ]
        ]
        ['ALERT_ENROLLMENT', 'PRESENT_VALUE', '12'      *AlertEnrollmentPresentValue
            [simple   BACnetApplicationTagObjectIdentifier                      presentValue                            ]
            [virtual  BACnetApplicationTagObjectIdentifier                      actualValue 'presentValue'              ]
        ]
        ['ANALOG_INPUT', 'PRESENT_VALUE', '4'           *AnalogInputPresentValue
            [simple   BACnetApplicationTagReal                                  presentValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'presentValue'              ]
        ]
        ['ANALOG_OUTPUT', 'PRESENT_VALUE', '4'          *AnalogOutputPresentValue
            [simple   BACnetApplicationTagReal                                  presentValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'presentValue'              ]
        ]
        ['ANALOG_VALUE', 'PRESENT_VALUE', '4'           *AnalogValuePresentValue
            [simple   BACnetApplicationTagReal                                  presentValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'presentValue'              ]
        ]
        ['BINARY_INPUT', 'PRESENT_VALUE', '9'           *BinaryInputPresentValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    presentValue                            ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'presentValue'              ]
        ]
        ['BINARY_OUTPUT', 'PRESENT_VALUE', '9'          *BinaryOutputPresentValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    presentValue                            ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'presentValue'              ]
        ]
        ['BINARY_VALUE', 'PRESENT_VALUE', '9'           *BinaryValuePresentValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    presentValue                            ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'presentValue'              ]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'PRESENT_VALUE', '9' *BinaryLightingOutputPresentValue
            [simple   BACnetBinaryLightingPVTagged('0', 'TagClass.APPLICATION_TAGS')    presentValue                    ]
            [virtual  BACnetBinaryLightingPVTagged                              actualValue 'presentValue'              ]
        ]
        ['BITSTRING_VALUE', 'PRESENT_VALUE', '8'        *BitStringValuePresentValue
            [simple   BACnetApplicationTagBitString                             presentValue                            ]
            [virtual  BACnetApplicationTagBitString                             actualValue 'presentValue'              ]
        ]
        ['CALENDAR', 'PRESENT_VALUE', '1'               *CalendarPresentValue
            [simple   BACnetApplicationTagBoolean                               presentValue                            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'presentValue'              ]
        ]
        ['CHANNEL', 'PRESENT_VALUE'                     *ChannelPresentValue
            [simple   BACnetChannelValue                                        presentValue                            ]
            [virtual  BACnetChannelValue                                        actualValue 'presentValue'              ]
        ]
        ['CHARACTERSTRING_VALUE', 'PRESENT_VALUE','7'   *CharacterStringValuePresentValue
            [simple   BACnetApplicationTagCharacterString                       presentValue                            ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'presentValue'              ]
        ]
        ['CREDENTIAL_DATA_INPUT', 'PRESENT_VALUE'       *CredentialDataInputPresentValue
            [simple   BACnetAuthenticationFactor                                presentValue                            ]
            [virtual  BACnetAuthenticationFactor                                actualValue 'presentValue'              ]
        ]
        ['DATE_VALUE', 'PRESENT_VALUE', '10'            *DateValuePresentValue
            [simple   BACnetApplicationTagDate                                  presentValue                            ]
            [virtual  BACnetApplicationTagDate                                  actualValue 'presentValue'              ]
        ]
        ['DATEPATTERN_VALUE', 'PRESENT_VALUE', '10'     *DatePatternValuePresentValue
            [simple   BACnetApplicationTagDate                                  presentValue                            ]
            [virtual  BACnetApplicationTagDate                                  actualValue 'presentValue'              ]
        ]
        ['DATETIME_VALUE', 'PRESENT_VALUE', '11'        *DateTimeValuePresentValue
            [simple   BACnetDateTime                                            presentValue                            ]
            [virtual  BACnetDateTime                                            actualValue 'presentValue'              ]
        ]
        ['DATETIMEPATTERN_VALUE', 'PRESENT_VALUE', '11' *DateTimePatternValuePresentValue
            [simple   BACnetDateTime                                            presentValue                            ]
            [virtual  BACnetDateTime                                            actualValue 'presentValue'              ]
        ]
        ['INTEGER_VALUE', 'PRESENT_VALUE', '3'          *IntegerValuePresentValue
            [simple   BACnetApplicationTagSignedInteger                         presentValue                            ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'presentValue'              ]
        ]
        ['LARGE_ANALOG_VALUE', 'PRESENT_VALUE', '5'     *LargeAnalogValuePresentValue
            [simple   BACnetApplicationTagDouble                                presentValue                            ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'presentValue'              ]
        ]
        ['LIGHTING_OUTPUT', 'PRESENT_VALUE', '4'        *LightingOutputPresentValue
            [simple   BACnetApplicationTagReal                                  presentValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'presentValue'              ]
        ]
        ['LIFE_SAFETY_POINT', 'PRESENT_VALUE', '9'      *LifeSafetyPointPresentValue
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')     presentValue                    ]
            [virtual  BACnetLifeSafetyStateTagged                               actualValue 'presentValue'              ]
        ]
        ['LIFE_SAFETY_ZONE', 'PRESENT_VALUE', '9'       *LifeSafetyZonePresentValue
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')     presentValue                    ]
            [virtual  BACnetLifeSafetyStateTagged                               actualValue 'presentValue'              ]
        ]
        ['LOAD_CONTROL', 'PRESENT_VALUE', '9'           *LoadControlPresentValue
            [simple   BACnetShedStateTagged('0', 'TagClass.APPLICATION_TAGS')           presentValue                    ]
            [virtual  BACnetShedStateTagged                                     actualValue 'presentValue'              ]
        ]
        ['LOOP', 'PRESENT_VALUE', '4'                   *LoopPresentValue
            [simple   BACnetApplicationTagReal                                  presentValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'presentValue'              ]
        ]
        ['PULSE_CONVERTER', 'PRESENT_VALUE', '4'        *PulseConverterPresentValue
            [simple   BACnetApplicationTagReal                                  presentValue                            ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'presentValue'              ]
        ]
        ['GROUP', 'PRESENT_VALUE'                       *GroupPresentValue
            [array    BACnetReadAccessResult
                          presentValue
                                  terminated
                                  'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'      ]
        ]
        ['GLOBAL_GROUP', 'PRESENT_VALUE'                *GlobalGroupPresentValue
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetPropertyAccessResult
                          presentValue
                                  terminated
                                  'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'      ]
        ]
        ['OCTETSTRING_VALUE', 'PRESENT_VALUE', '6'      *OctetStringValuePresentValue
            [simple   BACnetApplicationTagOctetString                           presentValue                            ]
            [virtual  BACnetApplicationTagOctetString                           actualValue 'presentValue'              ]
        ]
        ['SCHEDULE', 'PRESENT_VALUE'                    *SchedulePresentValue
            [simple   BACnetConstructedDataElement('BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                presentValue                            ]
            [virtual  BACnetConstructedDataElement                              actualValue 'presentValue'              ]
        ]
        ['TIME_VALUE', 'PRESENT_VALUE', '11'            *TimeValuePresentValue
            [simple   BACnetApplicationTagTime                                  presentValue                            ]
            [virtual  BACnetApplicationTagTime                                  actualValue 'presentValue'              ]
        ]
        ['TIMEPATTERN_VALUE', 'PRESENT_VALUE', '11'     *TimePatternValuePresentValue
            [simple   BACnetApplicationTagTime                                  presentValue                            ]
            [virtual  BACnetApplicationTagTime                                  actualValue 'presentValue'              ]
        ]
        [*, 'PRESENT_VALUE', '2'                        *PresentValue
            [simple   BACnetApplicationTagUnsignedInteger                       presentValue                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'presentValue'              ]
        ]
        [*, 'PRIORITY'                                  *Priority
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                            priority
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
            [validation 'arrayIndexArgument!=null || COUNT(priority) == 3'
                        "priority should have exactly 3 values"                                                         ]
        ]
        [*, 'PRIORITY_ARRAY'                            *PriorityArray
            [simple   BACnetPriorityArray('objectTypeArgument', 'tagNumber', 'arrayIndexArgument')   priorityArray      ]
            [virtual  BACnetPriorityArray                                       actualValue 'priorityArray'             ]
        ]
        [*, 'PRIORITY_FOR_WRITING', '2'                 *PriorityForWriting
            [simple   BACnetApplicationTagUnsignedInteger                       priorityForWriting                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'priorityForWriting'        ]
        ]
        [*, 'PROCESS_IDENTIFIER', '2'                   *ProcessIdentifier
            [simple   BACnetApplicationTagUnsignedInteger                       processIdentifier                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'processIdentifier'         ]
        ]
        [*, 'PROCESS_IDENTIFIER_FILTER'                 *ProcessIdentifierFilter
            [simple   BACnetProcessIdSelection                                  processIdentifierFilter                 ]
            [virtual  BACnetProcessIdSelection                                  actualValue 'processIdentifierFilter'   ]
        ]
        [*, 'PROFILE_LOCATION', '7'                     *ProfileLocation
            [simple   BACnetApplicationTagCharacterString                       profileLocation                         ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'profileLocation'           ]
        ]
        [*, 'PROFILE_NAME', '7'                         *ProfileName
            [simple   BACnetApplicationTagCharacterString                       profileName                             ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'profileName'               ]
        ]
        [*, 'PROGRAM_CHANGE', '9'                       *ProgramChange
            [simple   BACnetProgramRequestTagged('0', 'TagClass.APPLICATION_TAGS')           programChange              ]
            [virtual  BACnetProgramRequestTagged                                actualValue 'programChange'             ]
        ]
        [*, 'PROGRAM_LOCATION', '7'                     *ProgramLocation
            [simple   BACnetApplicationTagCharacterString                       programLocation                         ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'programLocation'           ]
        ]
        [*, 'PROGRAM_STATE', '9'                        *ProgramState
            [simple   BACnetProgramStateTagged('0', 'TagClass.APPLICATION_TAGS') programState                           ]
            [virtual  BACnetProgramStateTagged                                  actualValue 'programState'              ]
        ]
        [*, 'PROPERTY_LIST', '9'                        *PropertyList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetPropertyIdentifierTagged('0', 'TagClass.APPLICATION_TAGS')
                                propertyList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'PROPORTIONAL_CONSTANT', '4'                *ProportionalConstant
            [simple   BACnetApplicationTagReal                                  proportionalConstant                    ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'proportionalConstant'      ]
        ]
        [*, 'PROPORTIONAL_CONSTANT_UNITS', '9'          *ProportionalConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        [*, 'PROTOCOL_LEVEL', '9'                       *ProtocolLevel
            [simple   BACnetProtocolLevelTagged('0', 'TagClass.APPLICATION_TAGS')                    protocolLevel      ]
            [virtual  BACnetProtocolLevelTagged                                 actualValue 'protocolLevel'             ]
        ]
        //[*, 'PROTOCOL_CONFORMANCE_CLASS'              *ProtocolConformanceClass [validation    '1 == 2'    "TODO: implement me PROTOCOL_CONFORMANCE_CLASS *ProtocolConformanceClass"]]
        [*, 'PROTOCOL_OBJECT_TYPES_SUPPORTED', '8'      *ProtocolObjectTypesSupported
            [simple   BACnetObjectTypesSupportedTagged('0', 'TagClass.APPLICATION_TAGS')            protocolObjectTypesSupported         ]
            [virtual  BACnetObjectTypesSupportedTagged                          actualValue 'protocolObjectTypesSupported' ]
        ]
        [*, 'PROTOCOL_REVISION', '2'                    *ProtocolRevision
            [simple   BACnetApplicationTagUnsignedInteger                       protocolRevision                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'protocolRevision'          ]
        ]
        [*, 'PROTOCOL_SERVICES_SUPPORTED','8'           *ProtocolServicesSupported
            [simple   BACnetServicesSupportedTagged('0', 'TagClass.APPLICATION_TAGS')   protocolServicesSupported       ]
            [virtual  BACnetServicesSupportedTagged                             actualValue 'protocolServicesSupported' ]
        ]
        [*, 'PROTOCOL_VERSION', '2'                     *ProtocolVersion
            [simple   BACnetApplicationTagUnsignedInteger                       protocolVersion                         ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'protocolVersion'           ]
        ]
        [*, 'PULSE_RATE', '2'                           *PulseRate
            [simple   BACnetApplicationTagUnsignedInteger                       pulseRate                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'pulseRate'                 ]
        ]
        [*, 'READ_ONLY', '1'                            *ReadOnly
            [simple   BACnetApplicationTagBoolean                               readOnly                                ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'readOnly'                  ]
        ]
        [*, 'REASON_FOR_DISABLE', '9'                   *ReasonForDisable
            [array    BACnetAccessCredentialDisableReasonTagged('0', 'TagClass.APPLICATION_TAGS')
                                            reasonForDisable
                                                    terminated
                                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'REASON_FOR_HALT', '9'                      *ReasonForHalt
            [simple   BACnetProgramErrorTagged('0', 'TagClass.APPLICATION_TAGS')           programError                 ]
            [virtual  BACnetProgramErrorTagged                                  actualValue 'programError'              ]
        ]
        [*, 'RECIPIENT_LIST'                            *RecipientList
            [array    BACnetDestination
                                recipientList
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['FILE', 'RECORD_COUNT', '2'                    *FileRecordCount
            [simple   BACnetApplicationTagUnsignedInteger                       recordCount                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'recordCount'               ]
        ]
        [*, 'RECORD_COUNT', '2'                         *RecordCount
            [simple   BACnetApplicationTagUnsignedInteger                       recordCount                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'recordCount'               ]
        ]
        [*, 'RECORDS_SINCE_NOTIFICATION', '2'           *RecordsSinceNotification
            [simple   BACnetApplicationTagUnsignedInteger                       recordsSinceNotifications               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'recordsSinceNotifications' ]
        ]
        [*, 'REFERENCE_PORT', '2'                       *ReferencePort
            [simple   BACnetApplicationTagUnsignedInteger                       referencePort                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'referencePort'             ]
        ]
        [*, 'REGISTERED_CAR_CALL'                       *RegisteredCarCall
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLiftCarCallList
                                            registeredCarCall
                                                    terminated
                                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'RELIABILITY', '9'                          *Reliability
            [simple   BACnetReliabilityTagged('0', 'TagClass.APPLICATION_TAGS') reliability                             ]
            [virtual  BACnetReliabilityTagged                                   actualValue 'reliability'               ]
        ]
        [*, 'RELIABILITY_EVALUATION_INHIBIT', '1'       *ReliabilityEvaluationInhibit
            [simple   BACnetApplicationTagBoolean                               reliabilityEvaluationInhibit            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'reliabilityEvaluationInhibit']
        ]
        ['ACCESS_DOOR', 'RELINQUISH_DEFAULT', '9'       *AccessDoorRelinquishDefault
            [simple   BACnetDoorValueTagged('0', 'TagClass.APPLICATION_TAGS')   relinquishDefault                       ]
            [virtual  BACnetDoorValueTagged                                     actualValue 'relinquishDefault'         ]
        ]
        ['ANALOG_OUTPUT', 'RELINQUISH_DEFAULT', '4'     *AnalogOutputRelinquishDefault
            [simple   BACnetApplicationTagReal                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'relinquishDefault'         ]
        ]
        ['ANALOG_VALUE', 'RELINQUISH_DEFAULT', '4'      *AnalogValueRelinquishDefault
            [simple   BACnetApplicationTagReal                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'relinquishDefault'         ]
        ]
        ['BINARY_OUTPUT', 'RELINQUISH_DEFAULT', '9'    *BinaryOutputRelinquishDefault
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    relinquishDefault                       ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'relinquishDefault'         ]
        ]
        ['BINARY_VALUE', 'RELINQUISH_DEFAULT', '9'    *BinaryValueRelinquishDefault
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    relinquishDefault                       ]
            [virtual  BACnetBinaryPVTagged                                      actualValue 'relinquishDefault'         ]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'RELINQUISH_DEFAULT', '9'    *BinaryLightingOutputRelinquishDefault
            [simple   BACnetBinaryLightingPVTagged('0', 'TagClass.APPLICATION_TAGS')    relinquishDefault               ]
            [virtual  BACnetBinaryLightingPVTagged                              actualValue 'relinquishDefault'         ]
        ]
        ['BITSTRING_VALUE', 'RELINQUISH_DEFAULT', '8'   *BitStringValueRelinquishDefault
            [simple   BACnetApplicationTagBitString                             relinquishDefault                       ]
            [virtual  BACnetApplicationTagBitString                             actualValue 'relinquishDefault'         ]
        ]
        ['CHARACTERSTRING_VALUE', 'RELINQUISH_DEFAULT', '7' *CharacterStringValueRelinquishDefault
            [simple   BACnetApplicationTagCharacterString                       relinquishDefault                       ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'relinquishDefault'         ]
        ]
        ['DATE_VALUE', 'RELINQUISH_DEFAULT', '10'       *DateValueRelinquishDefault
            [simple   BACnetApplicationTagDate                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagDate                                  actualValue 'relinquishDefault'         ]
        ]
        ['DATEPATTERN_VALUE', 'RELINQUISH_DEFAULT', '10'     *DatePatternValueRelinquishDefault
            [simple   BACnetApplicationTagDate                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagDate                                  actualValue 'relinquishDefault'         ]
        ]
        ['DATETIME_VALUE', 'RELINQUISH_DEFAULT'         *DateTimeValueRelinquishDefault
            [simple   BACnetDateTime                                            relinquishDefault                       ]
            [virtual  BACnetDateTime                                            actualValue 'relinquishDefault'         ]
        ]
        ['DATETIMEPATTERN_VALUE', 'RELINQUISH_DEFAULT'  *DateTimePatternValueRelinquishDefault
            [simple   BACnetDateTime                                            relinquishDefault                       ]
            [virtual  BACnetDateTime                                            actualValue 'relinquishDefault'         ]
        ]
        ['LARGE_ANALOG_VALUE', 'RELINQUISH_DEFAULT', '5'     *LargeAnalogValueRelinquishDefault
            [simple   BACnetApplicationTagDouble                                relinquishDefault                       ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'relinquishDefault'         ]
        ]
        ['LIGHTING_OUTPUT', 'RELINQUISH_DEFAULT', '4'   *LightingOutputRelinquishDefault
            [simple   BACnetApplicationTagReal                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'relinquishDefault'         ]
        ]
        ['TIMEPATTERN_VALUE', 'RELINQUISH_DEFAULT', '11'    *TimePatternValueRelinquishDefault
            [simple   BACnetApplicationTagTime                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagTime                                  actualValue 'relinquishDefault'         ]
        ]
        ['TIME_VALUE', 'RELINQUISH_DEFAULT', '11'       *TimeValueRelinquishDefault
            [simple   BACnetApplicationTagTime                                  relinquishDefault                       ]
            [virtual  BACnetApplicationTagTime                                  actualValue 'relinquishDefault'         ]
        ]
        ['INTEGER_VALUE', 'RELINQUISH_DEFAULT', '3'     *IntegerValueRelinquishDefault
            [simple   BACnetApplicationTagSignedInteger                         relinquishDefault                       ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'relinquishDefault'         ]
        ]
        ['OCTETSTRING_VALUE', 'RELINQUISH_DEFAULT', '6' *OctetStringValueRelinquishDefault
            [simple   BACnetApplicationTagSignedInteger                         relinquishDefault                       ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'relinquishDefault'         ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'RELINQUISH_DEFAULT', '2'    *PositiveIntegerValueRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                       relinquishDefault                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'relinquishDefault'         ]
        ]
        ['MULTI_STATE_OUTPUT', 'RELINQUISH_DEFAULT', '2'    *MultiStateOutputRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                       relinquishDefault                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'relinquishDefault'         ]
        ]
        ['MULTI_STATE_VALUE', 'RELINQUISH_DEFAULT', '2' *MultiStateValueRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                       relinquishDefault                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'relinquishDefault'         ]
        ]
        [*, 'RELINQUISH_DEFAULT', '2'                   *RelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                       relinquishDefault                       ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'relinquishDefault'         ]
        ]
        [*, 'REPRESENTS'                                *Represents
            [simple   BACnetDeviceObjectReference                               represents                              ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'represents'                ]
        ]
        [*, 'REQUESTED_SHED_LEVEL'                      *RequestedShedLevel
            [simple   BACnetShedLevel                                           requestedShedLevel                      ]
            [virtual  BACnetShedLevel                                           actualValue 'requestedShedLevel'        ]
        ]
        [*, 'REQUESTED_UPDATE_INTERVAL', '2'            *RequestedUpdateInterval
            [simple   BACnetApplicationTagUnsignedInteger                       requestedUpdateInterval                 ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'requestedUpdateInterval'   ]
        ]
        [*, 'REQUIRED'                                *Required
            [validation    '1 == 2'    "An property identified by REQUIRED should never occur in the wild"]
        ]
        ['LARGE_ANALOG_VALUE', 'RESOLUTION', '5'        *LargeAnalogValueResolution
            [simple   BACnetApplicationTagDouble                                resolution                              ]
            [virtual  BACnetApplicationTagDouble                                actualValue 'resolution'                ]
        ]
        ['INTEGER_VALUE', 'RESOLUTION', '3'             *IntegerValueResolution
            [simple   BACnetApplicationTagSignedInteger                         resolution                              ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'resolution'                ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'RESOLUTION', '2'    *PositiveIntegerValueResolution
            [simple   BACnetApplicationTagUnsignedInteger                       resolution                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'resolution'                ]
        ]
        ['TIMER', 'RESOLUTION', '2'                     *TimerResolution
            [simple   BACnetApplicationTagUnsignedInteger                       resolution                              ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'resolution'                ]
        ]
        [*, 'RESOLUTION', '4'                           *Resolution
            [simple   BACnetApplicationTagReal                                  resolution                              ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'resolution'                ]
        ]
        [*, 'RESTART_NOTIFICATION_RECIPIENTS'           *RestartNotificationRecipients
            [array    BACnetRecipient
                                restartNotificationRecipients
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'RESTORE_COMPLETION_TIME', '2'              *RestoreCompletionTime
            [simple   BACnetApplicationTagUnsignedInteger                       completionTime                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'completionTime'            ]
        ]
        [*, 'RESTORE_PREPARATION_TIME', '2'             *RestorePreparationTime
            [simple   BACnetApplicationTagUnsignedInteger                       restorePreparationTime                  ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'restorePreparationTime'    ]
        ]
        [*, 'ROUTING_TABLE'                             *RoutingTable
            [array    BACnetRouterEntry
                                routingTable
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SCALE'                                     *Scale
            [simple   BACnetScale                                               scale                                   ]
            [virtual  BACnetScale                                               actualValue 'scale'                     ]
        ]
        [*, 'SCALE_FACTOR', '4'                         *ScaleFactor
            [simple   BACnetApplicationTagReal                                  scaleFactor                             ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'scaleFactor'               ]
        ]
        [*, 'SCHEDULE_DEFAULT'                          *ScheduleDefault
            [simple   BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifierArgument', 'null')  scheduleDefault  ]
            [virtual  BACnetConstructedDataElement                              actualValue 'scheduleDefault'           ]
        ]
        [*, 'SECURED_STATUS', '9'                       *SecuredStatus
            [simple   BACnetDoorSecuredStatusTagged('0', 'TagClass.APPLICATION_TAGS')         securedStatus             ]
            [virtual  BACnetDoorSecuredStatusTagged                             actualValue 'securedStatus'             ]
        ]
        [*, 'SECURITY_PDU_TIMEOUT', '2'                 *SecurityPDUTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       securityPduTimeout                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'securityPduTimeout'        ]
        ]
        [*, 'SECURITY_TIME_WINDOW', '2'                 *SecurityTimeWindow
            [simple   BACnetApplicationTagUnsignedInteger                       securityTimeWindow                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'securityTimeWindow'        ]
        ]
        [*, 'SEGMENTATION_SUPPORTED', '9'               *SegmentationSupported
            [simple   BACnetSegmentationTagged('0', 'TagClass.APPLICATION_TAGS')  segmentationSupported                 ]
            [virtual  BACnetSegmentationTagged                                  actualValue 'segmentationSupported'     ]
        ]
        [*, 'SERIAL_NUMBER', '7'                        *SerialNumber
            [simple   BACnetApplicationTagCharacterString                       serialNumber                            ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'serialNumber'              ]
        ]
        [*, 'SETPOINT', '4'                             *Setpoint
            [simple   BACnetApplicationTagReal                                  setpoint                                ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'setpoint'                  ]
        ]
        [*, 'SETPOINT_REFERENCE'                        *SetpointReference
            [simple   BACnetSetpointReference                                   setpointReference                       ]
            [virtual  BACnetSetpointReference                                   actualValue 'setpointReference'         ]
        ]
        [*, 'SETTING', '2'                              *Setting
            [simple   BACnetApplicationTagUnsignedInteger                       setting                                 ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'setting'                   ]
        ]
        [*, 'SHED_DURATION', '2'                        *ShedDuration
            [simple   BACnetApplicationTagUnsignedInteger                       shedDuration                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'shedDuration'              ]
        ]
        [*, 'SHED_LEVEL_DESCRIPTIONS', '7'              *ShedLevelDescriptions
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                shedLevelDescriptions
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SHED_LEVELS', '2'                          *ShedLevels
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                shedLevels
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SILENCED', '9'                             *Silenced
            [simple   BACnetSilencedStateTagged('0', 'TagClass.APPLICATION_TAGS')       silenced                        ]
            [virtual  BACnetSilencedStateTagged                                 actualValue 'silenced'                  ]
        ]
        [*, 'SLAVE_ADDRESS_BINDING'                     *SlaveAddressBinding
            [array    BACnetAddressBinding
                                slaveAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SLAVE_PROXY_ENABLE', '1'                   *SlaveProxyEnable
            [simple   BACnetApplicationTagBoolean                               slaveProxyEnable                        ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'slaveProxyEnable'          ]
        ]
        [*, 'START_TIME'                                *StartTime
            [simple   BACnetDateTime                                            startTime                               ]
            [virtual  BACnetDateTime                                            actualValue 'startTime'                 ]
        ]
        [*, 'STATE_CHANGE_VALUES'                       *StateChangeValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetTimerStateChangeValue('objectTypeArgument')
                                stateChangeValues
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [validation 'arrayIndexArgument!=null || COUNT(stateChangeValues) == 7'
                        "stateChangeValues should have exactly 7 values"                                                ]
        ]
        [*, 'STATE_DESCRIPTION', '7'                    *StateDescription
            [simple   BACnetApplicationTagCharacterString                       stateDescription                        ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'stateDescription'          ]
        ]
        [*, 'STATE_TEXT', '7'                           *StateText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                stateText
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'STATUS_FLAGS', '8'                         *StatusFlags
            [simple   BACnetStatusFlagsTagged('0', 'TagClass.APPLICATION_TAGS') statusFlags                             ]
            [virtual  BACnetStatusFlagsTagged                                   actualValue 'statusFlags'               ]
        ]
        [*, 'STOP_TIME'                                 *StopTime
            [simple   BACnetDateTime                                            stopTime                                ]
            [virtual  BACnetDateTime                                            actualValue 'stopTime'                  ]
        ]
        [*, 'STOP_WHEN_FULL', '1'                       *StopWhenFull
            [simple   BACnetApplicationTagBoolean                               stopWhenFull                            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'stopWhenFull'              ]
        ]
        [*, 'STRIKE_COUNT', '2'                         *StrikeCount
            [simple   BACnetApplicationTagUnsignedInteger                       strikeCount                             ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'strikeCount'               ]
        ]
        [*, 'STRUCTURED_OBJECT_LIST'                    *StructuredObjectList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                                structuredObjectList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SUBORDINATE_ANNOTATIONS'                   *SubordinateAnnotations
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                    subordinateAnnotations
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_LIST'                          *SubordinateList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectReference
                        subordinateList
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'SUBORDINATE_NODE_TYPES', '9'               *SubordinateNodeTypes
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNodeTypeTagged('0', 'TagClass.APPLICATION_TAGS')
                        subordinateNodeTypes
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'SUBORDINATE_RELATIONSHIPS', '9'            *SubordinateRelationships
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetRelationshipTagged('0', 'TagClass.APPLICATION_TAGS')
                                    subordinateRelationships
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_TAGS'                          *SubordinateTags
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNameValueCollection('0')
                        subordinateList
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBSCRIBED_RECIPIENTS'                     *SubscribedRecipients
            [array    BACnetEventNotificationSubscription
                            subscribedRecipients
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUPPORTED_FORMAT_CLASSES'                  *SupportedFormatClasses
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                        supportedFormats
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]

        ]
        [*, 'SUPPORTED_FORMATS'                         *SupportedFormats
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAuthenticationFactorFormat
                                        supportedFormats
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'SUPPORTED_SECURITY_ALGORITHMS'             *SupportedSecurityAlgorithms
            [array    BACnetApplicationTagUnsignedInteger
                                        supportedSecurityAlgorithms
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'SYSTEM_STATUS', '9'                        *SystemStatus
            [simple   BACnetDeviceStatusTagged('0', 'TagClass.APPLICATION_TAGS')        systemStatus                    ]
            [virtual  BACnetDeviceStatusTagged                                  actualValue 'systemStatus'              ]
        ]
        [*, 'TAGS'                                      *Tags
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNameValue
                            tags
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'THREAT_AUTHORITY'                          *ThreatAuthority
            [simple   BACnetAccessThreatLevel                                   threatAuthority                         ]
            [virtual  BACnetAccessThreatLevel                                   actualValue 'threatAuthority'           ]
        ]
        [*, 'THREAT_LEVEL'                              *ThreatLevel
            [simple   BACnetAccessThreatLevel                                   threatLevel                             ]
            [virtual  BACnetAccessThreatLevel                                   actualValue 'threatLevel'               ]
        ]
        [*, 'TIME_DELAY', '2'                           *TimeDelay
            [simple   BACnetApplicationTagUnsignedInteger                       timeDelay                               ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'timeDelay'                 ]
        ]
        [*, 'TIME_DELAY_NORMAL', '2'                    *TimeDelayNormal
            [simple   BACnetApplicationTagUnsignedInteger                       timeDelayNormal                         ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'timeDelayNormal'           ]
        ]
        [*, 'TIME_OF_ACTIVE_TIME_RESET'                 *TimeOfActiveTimeReset
            [simple   BACnetDateTime                                            timeOfActiveTimeReset                   ]
            [virtual  BACnetDateTime                                            actualValue 'timeOfActiveTimeReset'     ]
        ]
        [*, 'TIME_OF_DEVICE_RESTART'                    *TimeOfDeviceRestart
            [simple   BACnetTimeStamp                                           timeOfDeviceRestart                     ]
            [virtual  BACnetTimeStamp                                           actualValue 'timeOfDeviceRestart'       ]
        ]
        [*, 'TIME_OF_STATE_COUNT_RESET'                 *TimeOfStateCountReset
            [simple   BACnetDateTime                                            timeOfStateCountReset                   ]
            [virtual  BACnetDateTime                                            actualValue 'timeOfStateCountReset'     ]
        ]
        [*, 'TIME_OF_STRIKE_COUNT_RESET'                *TimeOfStrikeCountReset
            [simple   BACnetDateTime                                            timeOfStrikeCountReset                  ]
            [virtual  BACnetDateTime                                            actualValue 'timeOfStrikeCountReset'    ]
        ]
        [*, 'TIME_SYNCHRONIZATION_INTERVAL', '2'        *TimeSynchronizationInterval
            [simple   BACnetApplicationTagUnsignedInteger                       timeSynchronization                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'timeSynchronization'       ]
        ]
        [*, 'TIME_SYNCHRONIZATION_RECIPIENTS'           *TimeSynchronizationRecipients
            [array    BACnetRecipient
                            timeSynchronizationRecipients
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'TIMER_RUNNING', '1'                        *TimerRunning
            [simple   BACnetApplicationTagBoolean                               timerRunning                            ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'timerRunning'              ]
        ]
        [*, 'TIMER_STATE', '9'                          *TimerState
            [simple   BACnetTimerStateTagged('0', 'TagClass.APPLICATION_TAGS')  timerState]
            [virtual  BACnetTimerStateTagged                                    actualValue 'timerState'                ]
        ]
        [*, 'TOTAL_RECORD_COUNT', '2'                   *TotalRecordCount
            [simple   BACnetApplicationTagUnsignedInteger                       totalRecordCount                        ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'totalRecordCount'          ]
        ]
        [*, 'TRACE_FLAG', '1'                           *TraceFlag
            [simple   BACnetApplicationTagBoolean                               traceFlag                               ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'traceFlag'                 ]
        ]
        ['LIGHTING_OUTPUT','TRACKING_VALUE', '4'        *LightingOutputTrackingValue
            [simple   BACnetApplicationTagReal                                  trackingValue                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'trackingValue'             ]
        ]
        ['LIGHTING_OUTPUT','TRACKING_VALUE', '4'        *LightingOutputTrackingValue
            [simple   BACnetApplicationTagReal                                  trackingValue                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'trackingValue'             ]
        ]
        [*, 'TRACKING_VALUE', '9'                       *TrackingValue
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS') trackingValue                       ]
            [virtual  BACnetLifeSafetyStateTagged                               actualValue 'trackingValue'             ]
        ]
        [*, 'TRANSACTION_NOTIFICATION_CLASS', '2'       *TransactionNotificationClass
            [simple   BACnetApplicationTagUnsignedInteger                       transactionNotificationClass            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'transactionNotificationClass'  ]
        ]
        [*, 'TRANSITION', '9'                           *Transition
            [simple   BACnetLightingTransitionTagged('0', 'TagClass.APPLICATION_TAGS')                    transition    ]
            [virtual  BACnetLightingTransitionTagged                            actualValue 'transition'                ]
        ]
        [*, 'TRIGGER', '1'                              *Trigger
            [simple   BACnetApplicationTagBoolean                               trigger                                 ]
            [virtual  BACnetApplicationTagBoolean                               actualValue 'trigger'                   ]
        ]
        [*, 'UNITS', '9'                                *Units
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
            [virtual  BACnetEngineeringUnitsTagged                              actualValue 'units'                     ]
        ]
        [*, 'UPDATE_INTERVAL', '2'                      *UpdateInterval
            [simple   BACnetApplicationTagUnsignedInteger                       updateInterval                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'updateInterval'            ]
        ]
        [*, 'UPDATE_KEY_SET_TIMEOUT', '2'               *UpdateKeySetTimeout
            [simple   BACnetApplicationTagUnsignedInteger                       updateKeySetTimeout                     ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'updateKeySetTimeout'       ]
        ]
        ['CREDENTIAL_DATA_INPUT', 'UPDATE_TIME'         *CredentialDataInputUpdateTime
            [simple   BACnetTimeStamp                                           updateTime                              ]
            [virtual  BACnetTimeStamp                                           actualValue 'updateTime'                ]
        ]
        [*, 'UPDATE_TIME'                               *UpdateTime
            [simple   BACnetDateTime                                            updateTime                              ]
            [virtual  BACnetDateTime                                            actualValue 'updateTime'                ]
        ]
        [*, 'USER_EXTERNAL_IDENTIFIER', '7'             *UserExternalIdentifier
            [simple   BACnetApplicationTagCharacterString                       userExternalIdentifier                  ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'userExternalIdentifier'    ]
        ]
        [*, 'USER_INFORMATION_REFERENCE', '7'           *UserInformationReference
            [simple   BACnetApplicationTagCharacterString                       userInformationReference                ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'userInformationReference'  ]
        ]
        [*, 'USER_NAME', '7'                            *UserName
            [simple   BACnetApplicationTagCharacterString                       userName                                ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'userName'                  ]
        ]
        [*, 'USER_TYPE', '9'                            *UserType
            [simple   BACnetAccessUserTypeTagged('0', 'TagClass.APPLICATION_TAGS')      userType                        ]
            [virtual  BACnetAccessUserTypeTagged                                actualValue 'userType'                  ]
        ]
        [*, 'USES_REMAINING', '3'                       *UsesRemaining
            [simple   BACnetApplicationTagSignedInteger                         usesRemaining                           ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'usesRemaining'             ]
        ]
        [*, 'UTC_OFFSET', '3'                           *UTCOffset
            [simple   BACnetApplicationTagSignedInteger                         utcOffset                               ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'utcOffset'                 ]
        ]
        [*, 'UTC_TIME_SYNCHRONIZATION_RECIPIENTS'       *UTCTimeSynchronizationRecipients
            [array    BACnetRecipient
                                utcTimeSynchronizationRecipients
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'VALID_SAMPLES', '2'                        *ValidSamples
            [simple   BACnetApplicationTagUnsignedInteger                       validSamples                            ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'validSamples'              ]
        ]
        [*, 'VALUE_BEFORE_CHANGE', '2'                  *ValueBeforeChange
            [simple   BACnetApplicationTagUnsignedInteger                       valuesBeforeChange                      ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'valuesBeforeChange'        ]
        ]
        [*, 'VALUE_CHANGE_TIME'                         *ValueChangeTime
            [simple   BACnetDateTime                                            valueChangeTime                         ]
            [virtual  BACnetDateTime                                            actualValue 'valueChangeTime'           ]
        ]
        [*, 'VALUE_SET', '2'                            *ValueSet
            [simple   BACnetApplicationTagUnsignedInteger                       valueSet                                ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'valueSet'                  ]
        ]
        [*, 'VALUE_SOURCE'                              *ValueSource
            [simple   BACnetValueSource                                         valueSource                             ]
            [virtual  BACnetValueSource                                         actualValue 'valueSource'               ]
        ]
        [*, 'VALUE_SOURCE_ARRAY'                        *ValueSourceArray
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetValueSource
                                vtClassesSupported
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [validation 'arrayIndexArgument!=null || COUNT(vtClassesSupported) == 16'
                        "vtClassesSupported should have exactly 16 values"                                              ]
        ]
        [*, 'VARIANCE_VALUE', '4'                       *VarianceValue
            [simple   BACnetApplicationTagReal                                  varianceValue                           ]
            [virtual  BACnetApplicationTagReal                                  actualValue 'varianceValue'             ]
        ]
        // Note: checking 2 here is no accident as vendor-id is usually represented by unsigned not enumerated...
        //       the enum is a addition from plc4x
        [*, 'VENDOR_IDENTIFIER', '2'                    *VendorIdentifier
            [simple   BACnetVendorIdTagged('0', 'TagClass.APPLICATION_TAGS')    vendorIdentifier                        ]
            [virtual  BACnetVendorIdTagged                                      actualValue 'vendorIdentifier'          ]
        ]
        [*, 'VENDOR_NAME', '7'                          *VendorName
            [simple   BACnetApplicationTagCharacterString                       vendorName                              ]
            [virtual  BACnetApplicationTagCharacterString                       actualValue 'vendorName'                ]
        ]
        [*, 'VERIFICATION_TIME', '3'                    *VerificationTime
            [simple   BACnetApplicationTagSignedInteger                         verificationTime                        ]
            [virtual  BACnetApplicationTagSignedInteger                         actualValue 'verificationTime'          ]
        ]
        [*, 'VIRTUAL_MAC_ADDRESS_TABLE'                 *VirtualMACAddressTable
            [array    BACnetVMACEntry
                                virtualMacAddressTable
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'VT_CLASSES_SUPPORTED', '9'                 *VTClassesSupported
            [array    BACnetVTClassTagged('0', 'TagClass.APPLICATION_TAGS')
                                vtClassesSupported
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'WEEKLY_SCHEDULE'                           *WeeklySchedule
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDailySchedule
                                weeklySchedule
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
            [validation 'arrayIndexArgument!=null || COUNT(weeklySchedule) == 7'
                        "weeklySchedule should have exactly 7 values"                                                   ]
        ]
        [*, 'WINDOW_INTERVAL', '2'                      *WindowInterval
            [simple   BACnetApplicationTagUnsignedInteger                       windowInterval                          ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'windowInterval'            ]
        ]
        [*, 'WINDOW_SAMPLES', '2'                       *WindowSamples
            [simple   BACnetApplicationTagUnsignedInteger                       windowSamples                           ]
            [virtual  BACnetApplicationTagUnsignedInteger                       actualValue 'windowSamples'             ]
        ]
        [*, 'WRITE_STATUS', '9'                         *WriteStatus
            [simple   BACnetWriteStatusTagged('0', 'TagClass.APPLICATION_TAGS') writeStatus                             ]
            [virtual  BACnetWriteStatusTagged                                   actualValue 'writeStatus'               ]
        ]
        [*, 'ZONE_FROM'                                 *ZoneFrom
            [simple   BACnetDeviceObjectReference                               zoneFrom                                ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'zoneFrom'                  ]
        ]
        [*, 'ZONE_MEMBERS'                              *ZoneMembers
            [array    BACnetDeviceObjectReference
                    members
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ZONE_TO'                                 *ZoneTo
            [simple   BACnetDeviceObjectReference                               zoneTo                                  ]
            [virtual  BACnetDeviceObjectReference                               actualValue 'zoneTo'                    ]
        ]
        // BACnetConstructedDataUnspecified is used for unmapped properties
        [* *Unspecified
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifierArgument', 'arrayIndexArgument')
                            data                    terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
    ]
    [simple       BACnetClosingTag('tagNumber')
                        closingTag                                                                              ]
]

[type BACnetDeviceObjectReference
    [optional BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        deviceIdentifier                                                                        ]
    [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        objectIdentifier                                                                        ]
]

[type BACnetDeviceObjectReferenceEnclosed(uint 8 tagNumber)
   [simple   BACnetOpeningTag('tagNumber')
                   openingTag                   ]
   [simple   BACnetDeviceObjectReference
                   objectReference              ]
   [simple   BACnetClosingTag('tagNumber')
                   closingTag                   ]
]

[type BACnetDeviceObjectPropertyReference
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        objectIdentifier                                                                        ]
    [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                        propertyIdentifier                                                                      ]
    [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                        arrayIndex                                                                              ]
    [optional BACnetContextTagObjectIdentifier('3', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                        deviceIdentifier                                                                        ]
]

[type BACnetConstructedDataElement(BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument, BACnetTagPayloadUnsignedInteger arrayIndexArgument)
    [peek       BACnetTagHeader
                            peekedTagHeader                                                                     ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [virtual    bit         isApplicationTag    'peekedTagHeader.tagClass == TagClass.APPLICATION_TAGS'         ]
    [virtual    bit         isConstructedData   '!isApplicationTag && peekedTagHeader.lengthValueType == 0x6'   ]
    [virtual    bit         isContextTag        '!isConstructedData && !isApplicationTag'                       ]
    [validation '!isContextTag || (isContextTag && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected closing tag"                                                                        ]
    [optional   BACnetApplicationTag
                            applicationTag      'isApplicationTag'                                              ]
    [optional   BACnetContextTag('peekedTagNumber', 'BACnetDataType.UNKNOWN')
                            contextTag          'isContextTag'                                                  ]
    [optional   BACnetConstructedData('peekedTagNumber', 'objectTypeArgument', 'propertyIdentifierArgument', 'arrayIndexArgument')
                            constructedData     'isConstructedData'                                             ]
    [validation '(isApplicationTag && applicationTag != null) || (isContextTag && contextTag != null) || (isConstructedData && constructedData != null)'
                "BACnetConstructedDataElement could not parse anything"                                         ]
]

[type BACnetOptionalBinaryPV
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Null
            [simple   BACnetApplicationTagNull
                            nullValue                 ]
        ]
        [* *Value
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')
                            binaryPv                    ]
        ]
    ]
]

[type BACnetOptionalCharacterString
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Null
            [simple   BACnetApplicationTagNull
                                nullValue               ]
        ]
        [* *Value
            [simple   BACnetApplicationTagCharacterString
                            characterstring             ]
        ]
    ]
]

[type BACnetOptionalREAL
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Null
            [simple   BACnetApplicationTagNull
                                nullValue               ]
        ]
        [* *Value
            [simple   BACnetApplicationTagReal
                            realValue                   ]
        ]
    ]
]

[type BACnetOptionalUnsigned
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Null
            [simple   BACnetApplicationTagNull
                                nullValue               ]
        ]
        [* *Value
            [simple   BACnetApplicationTagUnsignedInteger
                                unsignedValue           ]
        ]
    ]
]

[type BACnetNameValue
    [simple   BACnetContextTagCharacterString('0', 'BACnetDataType.CHARACTER_STRING')
                            name                                                                                        ]
    [optional BACnetConstructedData('1', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            value                                                                                       ]
]

[type BACnetNameValueCollection(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                            openingTag                                                                                  ]
    [array    BACnetNameValue
                            members
                                   terminated
                                   'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'     ]
    [simple   BACnetClosingTag('tagNumber')
                            closingTag                                                                                  ]
]

[type BACnetAddressBinding
    [simple   BACnetApplicationTagObjectIdentifier
                            deviceIdentifier                                                                            ]
    [simple   BACnetAddress
                            deviceAddress                                                                               ]
]

[type BACnetAccessRule
    [simple   BACnetAccessRuleTimeRangeSpecifierTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            timeRangeSpecifier                                                                          ]
    [optional BACnetDeviceObjectPropertyReferenceEnclosed('1')
                            timeRange           'timeRangeSpecifier!=null'                                              ]
    [simple   BACnetAccessRuleLocationSpecifierTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            locationSpecifier                                                                           ]
    [optional BACnetDeviceObjectReferenceEnclosed('3')
                            location           'locationSpecifier!=null'                                                ]
    [simple   BACnetContextTagBoolean('4', 'BACnetDataType.BOOLEAN')
                            enable                                                                                      ]
]

[type BACnetCredentialAuthenticationFactor
    [simple   BACnetAccessAuthenticationFactorDisableTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            disable                                                                                     ]
    [simple   BACnetAuthenticationFactorEnclosed('1')
                            authenticationFactor                                                                        ]
]

[type BACnetAuthenticationFactor
    [simple   BACnetAuthenticationFactorTypeTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            formatType                                                                                  ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            formatClass                                                                                 ]
    [simple   BACnetContextTagOctetString('2', 'BACnetDataType.OCTET_STRING')
                            value                                                                                       ]
]

[type BACnetAuthenticationFactorFormat
    [simple   BACnetAuthenticationFactorTypeTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            formatType                                                                                  ]
    [optional BACnetVendorIdTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            vendorId                                                                                    ]
    [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                            vendorFormat                                                                                ]
]

[type BACnetAuthenticationFactorEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                                                          ]
    [simple   BACnetAuthenticationFactor
                    authenticationFactor                                                                                ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                                                          ]
]

[type BACnetAssignedAccessRights
    [simple   BACnetDeviceObjectReferenceEnclosed('0')
                    assignedAccessRights                                                                                ]
    [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                    enable                                                                                              ]
]

[type BACnetAccessThreatLevel
    [simple   BACnetApplicationTagUnsignedInteger
                    threatLevel                                                                                         ]
]

[type BACnetAuthenticationPolicy
    [simple   BACnetAuthenticationPolicyList('0')
                    policy                                                                                              ]
    [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                    orderEnforced                                                                                       ]
    [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                    timeout                                                                                             ]
]

[type BACnetAuthenticationPolicyList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                                                          ]
    [array    BACnetAuthenticationPolicyListEntry
                    entries
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                                                          ]
]

[type BACnetAuthenticationPolicyListEntry
    [simple   BACnetDeviceObjectReferenceEnclosed('0')
                    credentialDataInput                                                                                 ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                    index                                                                                               ]
]

[type BACnetTimerStateChangeValue(BACnetObjectType objectTypeArgument)
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual  bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' *Null
           [simple  BACnetApplicationTagNull
                            nullValue                                                   ]
       ]
       ['0x1', 'false' *Boolean
           [simple   BACnetApplicationTagBoolean
                            booleanValue                                                ]
       ]
       ['0x2', 'false' *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x3', 'false' *Integer
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
       ['0x4', 'false' *Real
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x5', 'false' *Double
           [simple  BACnetApplicationTagDouble
                                doubleValue                                             ]
       ]
       ['0x6', 'false' *OctetString
           [simple   BACnetApplicationTagOctetString
                            octetStringValue                                            ]
       ]
       ['0x7', 'false' *CharacterString
           [simple   BACnetApplicationTagCharacterString
                            characterStringValue                                        ]
       ]
       ['0x8', 'false' *BitString
           [simple   BACnetApplicationTagBitString
                            bitStringValue                                              ]
       ]
       ['0x9', 'false' *Enumerated
           [simple   BACnetApplicationTagEnumerated
                            enumeratedValue                                             ]
       ]
       ['0xA', 'false' *Date
           [simple   BACnetApplicationTagDate
                            dateValue                                                   ]
       ]
       ['0xB', 'false' *Time
           [simple   BACnetApplicationTagTime
                            timeValue                                                   ]
       ]
       ['0xC', 'false' *Objectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                            objectidentifierValue                                       ]
       ]
       ['0', 'true' *NoValue
           [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                            noValue                                                     ]
       ]
       ['1', 'true' *ConstructedValue
            [simple   BACnetConstructedData('1', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                        constructedValue                                ]
       ]
       ['2', 'true' *DateTime
            [simple   BACnetDateTimeEnclosed('2')
                            dateTimeValue                                               ]
       ]
       ['3', 'true' *LightingCommand
           [simple   BACnetLightingCommandEnclosed('3')
                       ligthingCommandValue                                             ]
       ]
    ]
]

[type BACnetSpecialEvent
    [simple   BACnetSpecialEventPeriod
                        period                                      ]
    [simple   BACnetSpecialEventListOfTimeValues('2')
                        listOfTimeValues                            ]
    [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                        eventPriority                               ]
]

[type BACnetSpecialEventPeriod
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [validation         'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [typeSwitch peekedTagNumber
        ['0' *CalendarEntry
            [simple   BACnetCalendarEntryEnclosed('0')
                                    calendarEntry                                    ]
        ]
        ['1' *CalendarReference
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                   calendarReference                                ]
        ]
    ]
]

[type BACnetCalendarEntry
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [validation         'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [typeSwitch peekedTagNumber
        ['0' *Date
            [simple   BACnetContextTagDate('0', 'BACnetDataType.DATE')
                                        dateValue                                   ]
        ]
        ['1' *DateRange
            [simple   BACnetDateRangeEnclosed('1')
                                        dateRange                                   ]
        ]
        ['2' *WeekNDay
            [simple   BACnetWeekNDayTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        weekNDay                                    ]
        ]
    ]
]

[type BACnetCalendarEntryEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                openingTag                                                                              ]
    [simple   BACnetCalendarEntry
                calendarEntry                                                                           ]
    [simple   BACnetClosingTag('tagNumber')
                closingTag                                                                              ]
]

[type BACnetDateRange
    [simple   BACnetApplicationTagDate  startDate   ]
    [simple   BACnetApplicationTagDate  endDate     ]
]

[type BACnetDateRangeEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                openingTag                                                                              ]
    [simple   BACnetDateRange
                dateRange                                                                               ]
    [simple   BACnetClosingTag('tagNumber')
                closingTag                                                                              ]
]

[type BACnetWeekNDay
    [validation '1==2' "Unusable type. Exits only for consistency. Use BACnetWeekNDayTagged"]
]

[type BACnetWeekNDayTagged(uint 8 tagNumber, TagClass tagClass)
    [simple   BACnetTagHeader
                            header                                                                              ]
    [validation    'header.tagClass == tagClass'    "tag class doesn't match"                                   ]
    [validation    '(header.tagClass == TagClass.APPLICATION_TAGS) || (header.actualTagNumber == tagNumber)'
                                                    "tagnumber doesn't match" shouldFail=false                  ]
    [validation    'header.actualLength == 3' "We should have at least 3 octets"                                ]
    // TODO: once we progress in codegen var enough that we can detect the source for array access we can use that again
    // ... at the moment in java this produces a .get(0) call and this doesn't work with byte arrays
    //[simple        BACnetTagPayloadOctetString  payload                                                         ]
    // TODO see comment above
    //[virtual       uint 8 month                                     'payload.octets[0]'   ]
    // TODO: temporary
    [simple        uint 8 month]
    [virtual       bit    oddMonths                                 'month == 13'         ]
    [virtual       bit    evenMonths                                'month == 14'         ]
    [virtual       bit    anyMonth                                  'month == 0xFF'       ]
    // TODO see comment above
    //[virtual       uint 8 weekOfMonth                               'payload.octets[1]'   ]
    // TODO: temporary
    [simple        uint 8 weekOfMonth]
    [virtual       bit    days1to7                                  'weekOfMonth == 1'    ]
    [virtual       bit    days8to14                                 'weekOfMonth == 2'    ]
    [virtual       bit    days15to21                                'weekOfMonth == 3'    ]
    [virtual       bit    days22to28                                'weekOfMonth == 4'    ]
    [virtual       bit    days29to31                                'weekOfMonth == 5'    ]
    [virtual       bit    last7DaysOfThisMonth                      'weekOfMonth == 6'    ]
    [virtual       bit    any7DaysPriorToLast7DaysOfThisMonth       'weekOfMonth == 7'    ]
    [virtual       bit    any7DaysPriorToLast14DaysOfThisMonth      'weekOfMonth == 8'    ]
    [virtual       bit    any7DaysPriorToLast21DaysOfThisMonth      'weekOfMonth == 9'    ]
    [virtual       bit    anyWeekOfthisMonth                        'weekOfMonth == 0xFF' ]
    // TODO see comment above
    //[virtual       uint 8 dayOfWeek                                 'payload.octets[2]'   ]
    // TODO: temporary
    [simple        uint 8 dayOfWeek]
    [virtual       bit    anyDayOfWeek                              'dayOfWeek == 0xFF' ]
]

[type BACnetSpecialEventListOfTimeValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                                              ]
    [array    BACnetTimeValue
                    listOfTimeValues
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                                              ]
]

[type BACnetTimeValue
    [simple   BACnetApplicationTagTime
                    timeValue                                                                               ]
    [simple BACnetConstructedDataElement('BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                    value                                                                                   ]
]

[type BACnetDailySchedule
    [simple   BACnetOpeningTag('0')
                        openingTag                                                                          ]
    [array    BACnetTimeValue
                    daySchedule
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 0)'            ]
    [simple   BACnetClosingTag('0')
                    closingTag                                                                              ]
]

[type BACnetEventNotificationSubscription
    [simple   BACnetRecipientEnclosed('0')
                                recipient                                                                   ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                processIdentifier                                                           ]
    [optional BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')
                                issueConfirmedNotifications                                                 ]
    [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                timeRemaining                                                               ]
]

[type BACnetPortPermission
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                port                                                                        ]
    [optional BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                                enable                                                                      ]
]

[type BACnetProcessIdSelection
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Null
            [simple   BACnetApplicationTagNull
                            nullValue                   ]
        ]
        [* *Value
            [simple   BACnetApplicationTagUnsignedInteger
                            processIdentifier           ]
        ]
    ]
]

[type BACnetNetworkSecurityPolicy
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                portId                                                                      ]
    [simple   BACnetSecurityPolicyTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                securityLevel                                                               ]
]

[type BACnetSecurityKeySet
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                keyRevision                                                                 ]
    [simple   BACnetDateTimeEnclosed('1')
                                activationTime                                                              ]
    [simple   BACnetDateTimeEnclosed('2')
                                expirationTime                                                              ]
    [simple   BACnetSecurityKeySetKeyIds('3')
                                keyIds                                                                      ]
]

[type BACnetSecurityKeySetKeyIds(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                            openingTag                                                                                  ]
    [array    BACnetKeyIdentifier
                            keyIds
                                   terminated
                                   'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'     ]
    [simple   BACnetClosingTag('tagNumber')
                            closingTag                                                                                  ]
]

[type BACnetKeyIdentifier
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                algorithm                                                                   ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                keyId                                                                       ]
]

[type BACnetRouterEntry
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                networkNumber                                                               ]
    [simple   BACnetContextTagOctetString('1', 'BACnetDataType.OCTET_STRING')
                                macAddress                                                                  ]
    [simple   BACnetRouterEntryStatusTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                status                                                                      ]
    [optional BACnetContextTagOctetString('3', 'BACnetDataType.OCTET_STRING')
                                performanceIndex                                                            ]
]

[type BACnetHostNPort
    [simple   BACnetHostAddressEnclosed('0')
                                host                                                                        ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                port                                                                        ]
]

[type BACnetHostNPortEnclosed(uint 8 tagNumber)
     [simple   BACnetOpeningTag('tagNumber')
                                 openingTag                                                                  ]
     [simple   BACnetHostNPort
                                 bacnetHostNPort                                                             ]
     [simple   BACnetClosingTag('tagNumber')
                                 closingTag                                                                  ]
 ]

[type BACnetHostAddress
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Null
            [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                                        none                                                                ]
        ]
        ['1' *IpAddress
            [simple   BACnetContextTagOctetString('1', 'BACnetDataType.OCTET_STRING')
                                        ipAddress                                                           ]
        ]
        ['2' *Name
            [simple   BACnetContextTagCharacterString('2', 'BACnetDataType.CHARACTER_STRING')
                                        name                                                                ]
        ]
    ]
]

[type BACnetHostAddressEnclosed(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                                openingTag                                                                  ]
    [simple   BACnetHostAddress
                                hostAddress                                                                 ]
    [simple   BACnetClosingTag('tagNumber')
                                closingTag                                                                  ]
]

[type BACnetBDTEntry
    [simple   BACnetHostNPortEnclosed('0')
                                bbmdAddress                                                                 ]
    [optional BACnetContextTagOctetString('1', 'BACnetDataType.OCTET_STRING')
                                broadcastMask                                                               ]
]

[type BACnetVMACEntry
    [optional BACnetContextTagOctetString('0', 'BACnetDataType.OCTET_STRING')
                                virtualMacAddress                                                           ]
    [optional BACnetContextTagOctetString('1', 'BACnetDataType.OCTET_STRING')
                                nativeMacAddress                                                            ]
]

[type BACnetSetpointReference
    [optional BACnetObjectPropertyReferenceEnclosed('0')
                                setPointReference                                                           ]
]

[type BACnetShedLevel
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *Percent
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        percent                                                             ]
        ]
        ['1' *Level
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        level                                                               ]
        ]
        ['2' *Amount
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        amount                                                              ]
        ]
    ]
]

[type BACnetLiftCarCallList
    [simple BACnetLiftCarCallListFloorList('0') floorNumbers                                            ]
]

[type BACnetLiftCarCallListFloorList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                             openingTag                                                                 ]
    [array    BACnetApplicationTagUnsignedInteger
                            floorNumbers
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                             closingTag                                                                 ]
]

[type BACnetAssignedLandingCalls
    [simple BACnetAssignedLandingCallsLandingCallsList('0')  landingCalls                               ]
]

[type BACnetAssignedLandingCallsLandingCallsList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                             openingTag                                                                 ]
    [array    BACnetAssignedLandingCallsLandingCallsListEntry
                            landingCalls
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                             closingTag                                                                 ]
]

[type BACnetAssignedLandingCallsLandingCallsListEntry
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                            floorNumber                                                                 ]
    [simple   BACnetLiftCarDirectionTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            direction                                                                   ]
]

[type BACnetLandingDoorStatus
    [simple BACnetLandingDoorStatusLandingDoorsList('0')  landingDoors                                  ]
]

[type BACnetLandingDoorStatusLandingDoorsList(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                             openingTag                                                                 ]
    [array    BACnetLandingDoorStatusLandingDoorsListEntry
                            landingDoors
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                             closingTag                                                                 ]
]

[type BACnetLandingDoorStatusLandingDoorsListEntry
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                            floorNumber                                                                 ]
    [simple   BACnetDoorStatusTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            doorStatus                                                                  ]
]

[type BACnetFaultParameter
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' *None
            [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                                                none                                                    ]
        ]
        ['1' *FaultCharacterString
            [simple   BACnetOpeningTag('1')
                                         openingTag                                                     ]
            [simple   BACnetFaultParameterFaultCharacterStringListOfFaultValues('0')
                                        listOfFaultValues                                               ]
            [simple   BACnetClosingTag('1')
                                         closingTag                                                     ]
        ]
        ['2' *FaultExtended
            [simple   BACnetOpeningTag('2')
                                         openingTag                                                     ]
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        vendorId                                                        ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        extendedFaultType                                               ]
            [simple   BACnetFaultParameterFaultExtendedParameters('2')
                                        parameters                                                      ]
            [simple   BACnetClosingTag('2')
                                         closingTag                                                     ]
        ]
        ['3' *FaultLifeSafety
            [simple   BACnetOpeningTag('3')
                                        openingTag                                                      ]
            [simple   BACnetFaultParameterFaultLifeSafetyListOfFaultValues('0')
                                        listOfFaultValues                                               ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        modePropertyReference                                           ]
            [simple   BACnetClosingTag('3')
                                        closingTag                                                      ]
        ]
        ['4' *FaultState
            [simple   BACnetOpeningTag('4')
                                        openingTag                                                      ]
            [simple   BACnetFaultParameterFaultStateListOfFaultValues('0')
                                        listOfFaultValues                                               ]
            [simple   BACnetClosingTag('4')
                                        closingTag                                                      ]
        ]
        ['5' *FaultStatusFlags
            [simple   BACnetOpeningTag('5')
                                        openingTag                                                      ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        statusFlagsReference                                            ]
            [simple   BACnetClosingTag('5')
                                        closingTag                                                      ]
        ]
        ['6' *FaultOutOfRange
            [simple   BACnetOpeningTag('6')
                                        openingTag                                                      ]
            [simple   BACnetFaultParameterFaultOutOfRangeMinNormalValue('0')
                                        minNormalValue                                                  ]
            [simple   BACnetFaultParameterFaultOutOfRangeMaxNormalValue('0')
                                        maxNormalValue                                                  ]
            [simple   BACnetClosingTag('6')
                                        closingTag                                                      ]
        ]
        ['7' *FaultListed
            [simple   BACnetOpeningTag('7')
                                        openingTag                                                      ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('0')
                                        faultListReference                                              ]
            [simple   BACnetClosingTag('7')
                                        closingTag                                                      ]
        ]
    ]
]

[type BACnetFaultParameterFaultCharacterStringListOfFaultValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                                         openingTag                                                     ]
    [array    BACnetApplicationTagCharacterString
                                listOfFaultValues
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
    [simple   BACnetClosingTag('tagNumber')
                                         closingTag                                                     ]
]

[type BACnetFaultParameterFaultExtendedParameters(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                                         openingTag                                                     ]
    [array    BACnetFaultParameterFaultExtendedParametersEntry
                                parameters
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
    [simple   BACnetClosingTag('tagNumber')
                                         closingTag                                                     ]
]

[type BACnetFaultParameterFaultExtendedParametersEntry
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual  bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' *Null
           [simple  BACnetApplicationTagNull
                            nullValue                                                   ]
       ]
       ['0x4', 'false' *Real
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x2', 'false' *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x1', 'false' *Boolean
           [simple   BACnetApplicationTagBoolean
                            booleanValue                                                ]
       ]
       ['0x3', 'false' *Integer
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
       ['0x5', 'false' *Double
           [simple  BACnetApplicationTagDouble
                                doubleValue                                             ]
       ]
       ['0x6', 'false' *OctetString
           [simple   BACnetApplicationTagOctetString
                            octetStringValue                                            ]
       ]
       ['0x7', 'false' *CharacterString
           [simple   BACnetApplicationTagCharacterString
                            characterStringValue                                        ]
       ]
       ['0x8', 'false' *BitString
           [simple   BACnetApplicationTagBitString
                            bitStringValue                                              ]
       ]
       ['0x9', 'false' *Enumerated
           [simple   BACnetApplicationTagEnumerated
                            enumeratedValue                                             ]
       ]
       ['0xA', 'false' *Date
           [simple   BACnetApplicationTagDate
                            dateValue                                                   ]
       ]
       ['0xB', 'false' *Time
           [simple   BACnetApplicationTagTime
                            timeValue                                                   ]
       ]
       ['0xC', 'false' *Objectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                            objectidentifierValue                                       ]
       ]
       ['0', 'true' *Reference
           [simple   BACnetDeviceObjectPropertyReferenceEnclosed('0')
                                       reference                                        ]
       ]
    ]
]

[type BACnetFaultParameterFaultLifeSafetyListOfFaultValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                            openingTag                                                     ]
    [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            listIfFaultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                            closingTag                                                     ]
]

[type BACnetFaultParameterFaultStateListOfFaultValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                            openingTag                                                     ]
    [array    BACnetPropertyStates
                            listIfFaultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
    [simple   BACnetClosingTag('tagNumber')
                            closingTag                                                     ]
]

[type BACnetFaultParameterFaultOutOfRangeMinNormalValue(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                                         openingTag                                     ]
    [peek       BACnetTagHeader
                           peekedTagHeader                                              ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'           ]
    [validation 'peekedTagHeader.tagClass == TagClass.APPLICATION_TAGS'
                "only application tags allowed"                                         ]
    [typeSwitch peekedTagNumber
       ['0x4' *Real
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x2' *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x5' *Double
           [simple  BACnetApplicationTagDouble
                            doubleValue                                                 ]
       ]
       ['0x3' *Integer
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                                closingTag                                              ]
]

[type BACnetFaultParameterFaultOutOfRangeMaxNormalValue(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                                         openingTag                                     ]
    [peek       BACnetTagHeader
                           peekedTagHeader                                              ]
    [virtual  uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'           ]
    [validation 'peekedTagHeader.tagClass == TagClass.APPLICATION_TAGS'
                "only application tags allowed"                                         ]
    [typeSwitch peekedTagNumber
       ['0x4' *Real
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x2' *Unsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x5' *Double
           [simple  BACnetApplicationTagDouble
                            doubleValue                                                 ]
       ]
       ['0x3' *Integer
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                                closingTag                                              ]
]

[type BACnetEventParameter
    [peek     BACnetTagHeader
                        peekedTagHeader                                                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'           ]
    [typeSwitch peekedTagNumber
        ['0' *ChangeOfBitstring
            [simple   BACnetOpeningTag('0')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetContextTagBitString('1', 'BACnetDataType.BIT_STRING')
                                        bitmask                                         ]
            [simple   BACnetEventParameterChangeOfBitstringListOfBitstringValues('2')
                                        listOfBitstringValues                           ]
            [simple   BACnetClosingTag('0')
                                        closingTag                                      ]
        ]
        ['1' *ChangeOfState
            [simple   BACnetOpeningTag('1')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfStateListOfValues('1')
                                        listOfValues                                    ]
            [simple   BACnetClosingTag('1')
                                        closingTag                                      ]
        ]
        ['2' *ChangeOfValue
            [simple   BACnetOpeningTag('2')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfValueCivCriteria('1')
                                        covCriteria                                     ]
            [simple   BACnetClosingTag('2')
                                        closingTag                                      ]
        ]
        ['3' *CommandFailure
            [simple   BACnetOpeningTag('3')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        feedbackPropertyReference                       ]
            [simple   BACnetClosingTag('3')
                                        closingTag                                      ]
        ]
        ['4' *FloatingLimit
            [simple   BACnetOpeningTag('4')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        setpointReference                               ]
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        lowDiffLimit                                    ]
            [simple   BACnetContextTagReal('3', 'BACnetDataType.REAL')
                                        highDiffLimit                                   ]
            [simple   BACnetContextTagReal('4', 'BACnetDataType.REAL')
                                        deadband                                        ]
            [simple   BACnetClosingTag('4')
                                        closingTag                                      ]
        ]
        ['5' *OutOfRange
            [simple   BACnetOpeningTag('5')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetContextTagReal('1', 'BACnetDataType.REAL')
                                        lowDiffLimit                                    ]
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        highDiffLimit                                   ]
            [simple   BACnetContextTagReal('3', 'BACnetDataType.REAL')
                                        deadband                                        ]
            [simple   BACnetClosingTag('5')
                                        closingTag                                      ]
        ]
        // 6 is undefined
        // 7 is deprecated
        ['8' *ChangeOfLifeSavety
            [simple   BACnetOpeningTag('8')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfLifeSavetyListOfLifeSavetyAlarmValues('1')
                                        listOfLifeSavetyAlarmValues                     ]
            [simple   BACnetEventParameterChangeOfLifeSavetyListOfAlarmValues('2')
                                        listOfAlarmValues                               ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('4')
                                        modePropertyReference                           ]
            [simple   BACnetClosingTag('8')
                                        closingTag                                      ]
        ]
        ['9' *Extended
            [simple   BACnetOpeningTag('9')
                                         openingTag                                     ]
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            vendorId                                                    ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            extendedEventType                                           ]
            [simple   BACnetEventParameterExtendedParameters('2')
                            parameters                                                  ]
            [simple   BACnetClosingTag('9')
                                        closingTag                                      ]
        ]
        ['10' *BufferReady
            [simple   BACnetOpeningTag('10')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        notificationThreshold                           ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        previousNotificationCount                       ]
            [simple   BACnetClosingTag('10')
                                        closingTag                                      ]
        ]
        ['11' *UnsignedRange
            [simple   BACnetOpeningTag('11')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        lowLimit                                       ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                        highLimit                                      ]
            [simple   BACnetClosingTag('11')
                                        closingTag                                      ]
        ]
        // 12 is reserved for future addenda
        ['13' *AccessEvent
            [simple   BACnetOpeningTag('13')
                                         openingTag                                     ]
            [simple   BACnetEventParameterAccessEventListOfAccessEvents('0')
                                        listOfAccessEvents                              ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        accessEventTimeReference                        ]
            [simple   BACnetClosingTag('13')
                                        closingTag                                      ]
        ]
        ['14' *DoubleOutOfRange
            [simple   BACnetOpeningTag('14')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetContextTagDouble('1', 'BACnetDataType.DOUBLE')
                                        lowLimit                                        ]
            [simple   BACnetContextTagDouble('2', 'BACnetDataType.DOUBLE')
                                        highLimit                                       ]
            [simple   BACnetContextTagDouble('3', 'BACnetDataType.DOUBLE')
                                        deadband                                        ]
            [simple   BACnetClosingTag('14')
                                        closingTag                                      ]
        ]
        ['15' *SignedOutOfRange
            [simple   BACnetOpeningTag('15')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetContextTagSignedInteger('1', 'BACnetDataType.SIGNED_INTEGER')
                                        lowLimit                                        ]
            [simple   BACnetContextTagSignedInteger('2', 'BACnetDataType.SIGNED_INTEGER')
                                        highLimit                                       ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                        deadband                                        ]
            [simple   BACnetClosingTag('15')
                                        closingTag                                      ]
        ]
        ['16' *UnsignedOutOfRange
            [simple   BACnetOpeningTag('16')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        lowLimit                                        ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                        highLimit                                       ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                        deadband                                        ]
            [simple   BACnetClosingTag('16')
                                        closingTag                                      ]
        ]
        ['17' *ChangeOfCharacterString
            [simple   BACnetOpeningTag('17')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfCharacterStringListOfAlarmValues('1')
                                        listOfAlarmValues                               ]
            [simple   BACnetClosingTag('17')
                                        closingTag                                      ]
        ]
        ['18' *ChangeOfStatusFlags
            [simple   BACnetOpeningTag('18')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        selectedFlags                                   ]
            [simple   BACnetClosingTag('18')
                                        closingTag                                      ]
        ]
        // 19 is not used
        ['20' *None
            [simple   BACnetContextTagNull('20', 'BACnetDataType.NULL')
                                        none                                            ]
        ]
        ['21' *ChangeOfDiscreteValue
            [simple   BACnetOpeningTag('21')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetClosingTag('21')
                                        closingTag                                      ]
        ]
        ['22' *ChangeOfTimer
            [simple   BACnetOpeningTag('22')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfTimerAlarmValue('1')
                                        alarmValues                                     ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('2')
                                        updateTimeReference                             ]
            [simple   BACnetClosingTag('22')
                                        closingTag                                      ]
        ]
    ]
]

[type BACnetEventParameterChangeOfBitstringListOfBitstringValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]
    [array    BACnetApplicationTagBitString
                    listOfBitstringValues
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                        ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetEventParameterChangeOfStateListOfValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]
    [array    BACnetPropertyStates
                    listOfValues
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                        ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetEventParameterChangeOfValueCivCriteria(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]
    [peek     BACnetTagHeader
                    peekedTagHeader                                                     ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'           ]
    [typeSwitch peekedTagNumber
        ['0' *Bitmask
            [simple   BACnetContextTagBitString('0', 'BACnetDataType.BIT_STRING')
                                        bitmask                                         ]
        ]
        ['1' *ReferencedPropertyIncrement
            [simple   BACnetContextTagReal('1', 'BACnetDataType.REAL')
                                        referencedPropertyIncrement                     ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetEventParameterChangeOfLifeSavetyListOfLifeSavetyAlarmValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]
    [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                    listOfLifeSavetyAlarmValues
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                        ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetEventParameterChangeOfLifeSavetyListOfAlarmValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]
    [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                    listOfAlarmValues
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                        ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetEventParameterExtendedParameters(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                        openingTag                                                  ]
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual  bit       isOpeningTag        'peekedTagHeader.lengthValueType == 0x6']
    [virtual  bit       isClosingTag        'peekedTagHeader.lengthValueType == 0x7']
    [optional BACnetApplicationTagNull
                    nullValue
                        'peekedTagNumber == 0x0 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagReal
                    realValue
                        'peekedTagNumber == 0x4 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagUnsignedInteger
                    unsignedValue
                        'peekedTagNumber == 0x2 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagBoolean
                    booleanValue
                        'peekedTagNumber == 0x1 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagSignedInteger
                    integerValue
                        'peekedTagNumber == 0x3 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagDouble
                    doubleValue
                        'peekedTagNumber == 0x5 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagOctetString
                    octetStringValue
                        'peekedTagNumber == 0x6 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagCharacterString
                    characterStringValue
                        'peekedTagNumber == 0x7 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagBitString
                    bitStringValue
                        'peekedTagNumber == 0x8 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagEnumerated
                    enumeratedValue
                        'peekedTagNumber == 0x9 && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagDate
                    dateValue
                        'peekedTagNumber == 0xA && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagTime
                    timeValue
                        'peekedTagNumber == 0xB && !isOpeningTag && !isClosingTag'      ]
    [optional BACnetApplicationTagObjectIdentifier
                    objectIdentifier
                        'peekedTagNumber == 0xC && !isOpeningTag'                       ]
    [optional BACnetDeviceObjectPropertyReferenceEnclosed('0')
                    reference
                        'isOpeningTag && !isClosingTag'                                 ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetEventParameterAccessEventListOfAccessEvents(uint 8 tagNumber)
   [simple   BACnetOpeningTag('tagNumber')
                   openingTag                                                          ]
   [array    BACnetDeviceObjectPropertyReference
                   listOfAccessEvents
                       terminated
                       'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                       ]
   [simple   BACnetClosingTag('tagNumber')
                   closingTag                                                          ]
]


[type BACnetEventParameterChangeOfCharacterStringListOfAlarmValues(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                   openingTag                                                          ]
    [array    BACnetApplicationTagCharacterString
                   listOfAlarmValues
                       terminated
                       'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                       ]
    [simple   BACnetClosingTag('tagNumber')
                   closingTag                                                          ]
]


[type BACnetEventParameterChangeOfTimerAlarmValue(uint 8 tagNumber)
   [simple   BACnetOpeningTag('tagNumber')
                   openingTag                                                          ]
   [array    BACnetTimerStateTagged('0', 'TagClass.APPLICATION_TAGS')
                   alarmValues
                       terminated
                       'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'
                                                                                       ]
   [simple   BACnetClosingTag('tagNumber')
                   closingTag                                                          ]
]

[type BACnetLandingCallStatus
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        floorNumber                                     ]
    [simple   BACnetLandingCallStatusCommand
                                        command                                         ]
    [optional BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')
                                        floorText                                       ]
]

[type BACnetLandingCallStatusCommand
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['1' *Direction
            [simple BACnetLiftCarDirectionTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                            direction                               ]
        ]
        ['2' *Destination
             [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                            destination                             ]
        ]
    ]
]

[type BACnetCOVMultipleSubscription
    [simple   BACnetRecipientProcessEnclosed('0')
                                        recipient                                       ]
    [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                                        issueConfirmedNotifications                     ]
    [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeRemaining                                   ]
    [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                        maxNotificationDelay                            ]
    [simple   BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification('4')
                                        listOfCovSubscriptionSpecification              ]
]

[type BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]

    [array    BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationEntry
                    listOfCovSubscriptionSpecificationEntry
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationEntry
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                    monitoredObjectIdentifier                                           ]
    [simple   BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationEntryListOfCovReferences('1')
                    listOfCovReferences                                                 ]
]

[type BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationEntryListOfCovReferences(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                          ]

    [array    BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationEntryListOfCovReferencesEntry
                    listOfCovReferences
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                          ]
]

[type BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationEntryListOfCovReferencesEntry
    [simple   BACnetPropertyReferenceEnclosed('0')
                    monitoredProperty                                                   ]
    [optional BACnetContextTagReal('1', 'BACnetDataType.REAL')
                    covIncrement                                                        ]
    [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                    timestamped                                                         ]
]

[type BACnetVTSession
    [simple   BACnetApplicationTagUnsignedInteger
                    localVtSessionId                                                    ]
    [simple   BACnetApplicationTagUnsignedInteger
                    removeVtSessionId                                                   ]
    [simple   BACnetAddress
                    remoteVtAddress                                                     ]
]

[type BACnetCOVSubscription
    [simple   BACnetRecipientProcessEnclosed('0')
                                        recipient                                       ]
    [simple   BACnetObjectPropertyReferenceEnclosed('1')
                                        monitoredPropertyReference                      ]
    [simple   BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')
                                        issueConfirmedNotifications                     ]
    [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeRemaining                                   ]
    [optional BACnetContextTagReal('4', 'BACnetDataType.REAL')
                                        covIncrement                                    ]
]

[type BACnetPrescale
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        multiplier                                  ]
    [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        moduloDivide                                ]
]

[type BACnetScale
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' *FloatScale
            [simple   BACnetContextTagReal('0', 'BACnetDataType.REAL')
                                        floatScale                                  ]
        ]
        ['1' *IntegerScale
             [simple   BACnetContextTagSignedInteger('1', 'BACnetDataType.SIGNED_INTEGER')
                                        integerScale                                ]
        ]
    ]
]

[type BACnetAccumulatorRecord
    [simple   BACnetDateTimeEnclosed('0')
                                timestamp                                           ]
    [simple   BACnetContextTagSignedInteger('1', 'BACnetDataType.SIGNED_INTEGER')
                                presentValue                                        ]
    [simple   BACnetContextTagSignedInteger('2', 'BACnetDataType.SIGNED_INTEGER')
                                accumulatedValue                                    ]
    [simple   BACnetAccumulatorRecordAccumulatorStatusTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accumulatorStatus                                   ]
]

[type BACnetValueSource
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' *None
            [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                                        none                                        ]
        ]
        ['1' *Object
             [simple   BACnetDeviceObjectReferenceEnclosed('1')
                                        object                                      ]
        ]
        ['2' *Address
             [simple   BACnetAddressEnclosed('2')
                                        address                                     ]
        ]
    ]
]

[type BACnetClientCOV
    [peek     BACnetTagHeader
                        peekedTagHeader                                             ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0x4' *Object
            [simple   BACnetApplicationTagReal
                                        realIncrement                               ]
        ]
        ['0x0' *None
            [simple   BACnetApplicationTagNull
                                        defaultIncrement                            ]
        ]
    ]
]

[type BACnetDestination
    [simple BACnetDaysOfWeekTagged('0', 'TagClass.APPLICATION_TAGS')
                    validDays                                                       ]
    [simple BACnetApplicationTagTime
                    fromTime                                                        ]
    [simple BACnetApplicationTagTime
                    toTime                                                          ]
    [simple BACnetRecipient
                    recipient                                                       ]
    [simple BACnetApplicationTagUnsignedInteger
                    processIdentifier                                               ]
    [simple BACnetApplicationTagBoolean
                    issueConfirmedNotifications                                     ]
    [simple BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS')
                    transitions                                                     ]
]

[type BACnetLogMultipleRecord
    [simple   BACnetDateTimeEnclosed('0')
                    timestamp                                                       ]
    [simple   BACnetLogData('1')
                    logData                                                         ]
]

[type BACnetLogData(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                      ]
    [peek     BACnetTagHeader
                    peekedTagHeader                                                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' *LogStatus
            [simple   BACnetLogStatusTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        logStatus                                   ]
        ]
        ['1' *LogData
            [simple   BACnetOpeningTag('1')
                                        innerOpeningTag                             ]
            [array    BACnetLogDataLogDataEntry
                                        logData
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)']
            [simple   BACnetClosingTag('1')
                                        innerClosingTag                             ]
        ]
        ['2' *LogDataTimeChange
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        timeChange                                  ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                      ]
]

[type BACnetLogDataLogDataEntry
    [peek     BACnetTagHeader
                    peekedTagHeader                                                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' *BooleanValue
            [simple   BACnetContextTagBoolean('0', 'BACnetDataType.BOOLEAN')
                                        booleanValue                                ]
        ]
        ['1' *RealValue
            [simple   BACnetContextTagReal('1', 'BACnetDataType.REAL')
                                        realValue                                   ]
        ]
        ['2' *EnumeratedValue
            [simple   BACnetContextTagEnumerated('2', 'BACnetDataType.ENUMERATED')
                                        enumeratedValue                             ]
        ]
        ['3' *UnsignedValue
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                        unsignedValue                               ]
        ]
        ['4' *IntegerValue
            [simple   BACnetContextTagSignedInteger('4', 'BACnetDataType.SIGNED_INTEGER')
                                        integerValue                                ]
        ]
        ['5' *BitStringValue
            [simple   BACnetContextTagBitString('5', 'BACnetDataType.BIT_STRING')
                                        bitStringValue                              ]
        ]
        ['6' *NullValue
            [simple   BACnetContextTagNull('6', 'BACnetDataType.NULL')
                                        nullValue                                   ]
        ]
        ['7' *Failure
            [simple   ErrorEnclosed('7')
                                        failure                                     ]
        ]
        ['8' *AnyValue
            [optional BACnetConstructedData('8', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                        anyValue                                    ]
        ]
    ]
]

[type BACnetLogRecord
    [simple   BACnetDateTimeEnclosed('0')
                    timestamp                                                       ]
    [simple   BACnetLogRecordLogDatum('1')
                    logDatum                                                        ]
    [optional BACnetStatusFlagsTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    statusFlags                                                     ]
]

[type BACnetLogRecordLogDatum(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                      ]
    [peek     BACnetTagHeader
                    peekedTagHeader                                                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' *LogStatus
            [simple   BACnetLogStatusTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        logStatus                                   ]
        ]
        ['1' *BooleanValue
            [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                                        booleanValue                                ]
        ]
        ['2' *RealValue
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        realValue                                   ]
        ]
        ['3' *EnumeratedValue
            [simple   BACnetContextTagEnumerated('3', 'BACnetDataType.ENUMERATED')
                                        enumeratedValue                             ]
        ]
        ['4' *UnsignedValue
            [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')
                                        unsignedValue                               ]
        ]
        ['5' *IntegerValue
            [simple   BACnetContextTagSignedInteger('5', 'BACnetDataType.SIGNED_INTEGER')
                                        integerValue                                ]
        ]
        ['6' *BitStringValue
            [simple   BACnetContextTagBitString('6', 'BACnetDataType.BIT_STRING')
                                        bitStringValue                              ]
        ]
        ['7' *NullValue
            [simple   BACnetContextTagNull('7', 'BACnetDataType.NULL')
                                        nullValue                                   ]
        ]
        ['8' *Failure
            [simple   ErrorEnclosed('8')
                                        failure                                     ]
        ]
        ['9' *TimeChange
            [simple   BACnetContextTagReal('9', 'BACnetDataType.REAL')
                                        timeChange                                  ]
        ]
        ['10' *AnyValue
            [optional BACnetConstructedData('10', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                        anyValue                                    ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                      ]
]

[type BACnetEventLogRecord
    [simple   BACnetDateTimeEnclosed('0')
                    timestamp                                                       ]
    [simple   BACnetEventLogRecordLogDatum('1')
                    logDatum                                                        ]
]

[type BACnetEventLogRecordLogDatum(uint 8 tagNumber)
    [simple   BACnetOpeningTag('tagNumber')
                    openingTag                                                      ]
    [peek     BACnetTagHeader
                    peekedTagHeader                                                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' *LogStatus
            [simple   BACnetLogStatusTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        logStatus                                   ]
        ]
        ['1' *Notification
            [simple   BACnetOpeningTag('1')
                            innerOpeningTag                                         ]
            //TODO this below slurps to much because of the service choice... :( find workaround we might need fragments for that...
            [simple   ConfirmedEventNotificationRequest
                                        notification                                ]
            [simple   BACnetClosingTag('tagNumber')
                            innerClosingTag                                         ]
        ]
        ['2' *TimeChange
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        timeChange                                  ]
        ]
    ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                                                      ]
]

// TODO: this is copy paste from BACnetConfirmedServiceRequestConfirmedEventNotification for now... at the end this should be a seperate type like here and above use like a fragment
[type ConfirmedEventNotificationRequest
    [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           processIdentifier              ]
    [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier     ]
    [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  eventObjectIdentifier          ]
    [simple   BACnetTimeStampEnclosed('3')                                                      timestamp                      ]
    [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')           notificationClass              ]
    [simple   BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')           priority                       ]
    [simple   BACnetEventTypeTagged('6', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      eventType                      ]
    [optional BACnetContextTagCharacterString('7', 'BACnetDataType.CHARACTER_STRING')           messageText                    ]
    [simple   BACnetNotifyTypeTagged('8', 'TagClass.CONTEXT_SPECIFIC_TAGS')                     notifyType                     ]
    [optional BACnetContextTagBoolean('9', 'BACnetDataType.BOOLEAN')                            ackRequired                    ]
    [optional BACnetEventStateTagged('10', 'TagClass.CONTEXT_SPECIFIC_TAGS')                    fromState                      ]
    [simple   BACnetEventStateTagged('11', 'TagClass.CONTEXT_SPECIFIC_TAGS')                    toState                        ]
    [optional BACnetNotificationParameters('12', 'eventObjectIdentifier.objectType')            eventValues                    ]
]

[type BACnetPropertyAccessResult
    [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                objectIdentifier                                        ]
    [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                propertyIdentifier                                      ]
    [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                propertyArrayIndex                                      ]
    [optional BACnetContextTagObjectIdentifier('3', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                deviceIdentifier                                        ]
    [simple   BACnetPropertyAccessResultAccessResult('objectIdentifier.objectType', 'propertyIdentifier.value', '(propertyArrayIndex!=null?propertyArrayIndex.payload:null)')
                                accessResult                                            ]
]

[type BACnetPropertyAccessResultAccessResult(BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument, BACnetTagPayloadUnsignedInteger propertyArrayIndexArgument)
    [peek     BACnetTagHeader
                    peekedTagHeader                                                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['4' *PropertyValue
            [simple   BACnetConstructedData('4', 'objectTypeArgument', 'propertyIdentifierArgument', 'propertyArrayIndexArgument')
                                        propertyValue                               ]
        ]
        ['5' *PropertyAccessError
            [simple   ErrorEnclosed('5')
                                        propertyAccessError                         ]
        ]
    ]
]
