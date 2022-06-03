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
    [const         uint 8   bacnetType   0x81                   ]
    [discriminator uint 8   bvlcFunction                        ]
    [implicit      uint 16  bvlcLength          'lengthInBytes' ]
    [virtual       uint 16  bvlcPayloadLength   'bvlcLength-4'  ]
    [typeSwitch bvlcFunction
        ['0x00' BVLCResult
            [simple   BVLCResultCode
                            code                                ]
        ]
        ['0x01' BVLCWriteBroadcastDistributionTable(uint 16 bvlcPayloadLength)
            [array BVLCBroadcastDistributionTableEntry
                            table
                                length 'bvlcPayloadLength'      ]
        ]
        ['0x02' BVLCReadBroadcastDistributionTable
        ]
        ['0x03' BVLCReadBroadcastDistributionTableAck(uint 16 bvlcPayloadLength)
            [array BVLCBroadcastDistributionTableEntry
                            table
                                length 'bvlcPayloadLength'      ]
        ]
        ['0x04' BVLCForwardedNPDU(uint 16 bvlcPayloadLength)
            [array    uint 8  ip    count '4'                     ]
            [simple   uint 16 port                                ]
            [simple   NPDU('bvlcPayloadLength - 6')
                            npdu                                ]
        ]
        ['0x05' BVLCRegisterForeignDevice
            [simple   uint 16 ttl]
        ]
        ['0x06' BVLCReadForeignDeviceTable
        ]
        ['0x07' BVLCReadForeignDeviceTableAck(uint 16 bvlcPayloadLength)
            [array BVLCForeignDeviceTableEntry
                            table
                                length 'bvlcPayloadLength'      ]
        ]
        ['0x08' BVLCDeleteForeignDeviceTableEntry
            [array    uint 8  ip  count '4'                       ]
            [simple   uint 16 port                                ]
        ]
        ['0x09' BVLCDistributeBroadcastToNetwork(uint 16 bvlcPayloadLength)
            [simple   NPDU('bvlcPayloadLength')
                            npdu                                ]
        ]
        ['0x0A' BVLCOriginalUnicastNPDU(uint 16 bvlcPayloadLength)
            [simple   NPDU('bvlcPayloadLength')
                            npdu                                ]
        ]
        ['0x0B' BVLCOriginalBroadcastNPDU(uint 16 bvlcPayloadLength)
            [simple   NPDU('bvlcPayloadLength')
                            npdu                                ]
        ]
        ['0x0C' BVLCSecureBVLL(uint 16 bvlcPayloadLength)
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
    [optional      BACnetVendorId
                            vendorId '(messageType >= 128) && (messageType <= 255)']
    [typeSwitch messageType
        ['0x00' NLMWhoIsRouterToNetwork(uint 8 messageType)
            [array      uint 16     destinationNetworkAddress length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
        ['0x01' NLMIAmRouterToNetwork(uint 8 messageType)
            [array      uint 16     destinationNetworkAddress length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
        ['0x02' NLMICouldBeRouterToNetwork(uint 8 messageType)
            [simple   uint 16     destinationNetworkAddress   ]
            [simple   uint 8      performanceIndex            ]
        ]
        ['0x03' NLMRejectRouterToNetwork(uint 8 messageType)
            [simple   NLMRejectRouterToNetworkRejectReason
                                    rejectReason              ]
            [simple   uint 16     destinationNetworkAddress   ]
        ]
        ['0x04' NLMRouterBusyToNetwork(uint 8 messageType)
            [array      uint 16     destinationNetworkAddress length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
        ['0x05' NLMRouterAvailableToNetwork(uint 8 messageType)
            [array      uint 16     destinationNetworkAddress length 'apduLength - (((messageType >= 128) && (messageType <= 255)) ? 3 : 1)']
        ]
        ['0x06' NLMInitalizeRoutingTable(uint 8 messageType)
            [simple   uint 8      numberOfPorts               ]
            [array      NLMInitalizeRoutingTablePortMapping
                                    portMappings
                        count 'numberOfPorts'                 ]
        ]
        ['0x07' NLMInitalizeRoutingTableAck(uint 8 messageType)
            [simple   uint 8      numberOfPorts               ]
            [array      NLMInitalizeRoutingTablePortMapping
                                    portMappings
                        count 'numberOfPorts'                 ]
        ]
        ['0x08' NLMEstablishConnectionToNetwork(uint 8 messageType)
            [simple   uint 16     destinationNetworkAddress   ]
            [simple   uint 8      terminationTime             ]
        ]
        ['0x09' NLMDisconnectConnectionToNetwork(uint 8 messageType)
            [simple   uint 16     destinationNetworkAddress   ]
        ]
    ]
]

[type NLMInitalizeRoutingTablePortMapping
    [simple   uint 16     destinationNetworkAddress       ]
    [simple   uint 8      portId                          ]
    [simple   uint 8      portInfoLength                  ]
    [array      byte        portInfo count 'portInfoLength' ]
]

[discriminatedType APDU(uint 16 apduLength)
    [discriminator ApduType apduType]
    [typeSwitch apduType
        ['CONFIRMED_REQUEST_PDU' APDUConfirmedRequest
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
            [optional uint 8    segmentServiceChoice 'segmentedMessage && sequenceNumber != 0']
            [virtual  uint 16   segmentReduction
                                    '(segmentServiceChoice != null)?(apduHeaderReduction+1):apduHeaderReduction'       ]
            [array    byte      segment
                                    length
                                    'segmentedMessage?((apduLength>0)?(apduLength - segmentReduction):0):0'             ]
        ]
        ['UNCONFIRMED_REQUEST_PDU' APDUUnconfirmedRequest
            [reserved uint 4                          '0'               ]
            [simple   BACnetUnconfirmedServiceRequest('apduLength - 1')
                                serviceRequest                          ]
        ]
        ['SIMPLE_ACK_PDU' APDUSimpleAck
            [reserved uint 4    '0'                                     ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   uint 8    serviceChoice                           ]
        ]
        ['COMPLEX_ACK_PDU' APDUComplexAck
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
            [optional uint 8    segmentServiceChoice 'segmentedMessage && sequenceNumber != 0']
            [virtual  uint 16   segmentReduction
                                    '(segmentServiceChoice != null)?(apduHeaderReduction+1):apduHeaderReduction'
                                                                        ]
            [array    byte      segment
                                    length
                                    'segmentedMessage?((apduLength>0)?(apduLength - segmentReduction):0):0'
                                                                        ]
        ]
        ['SEGMENT_ACK_PDU' APDUSegmentAck
            [reserved uint 2    '0x00'                                  ]
            [simple   bit       negativeAck                             ]
            [simple   bit       server                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   uint 8    sequenceNumber                          ]
            [simple   uint 8    proposedWindowSize                      ]
        ]
        ['ERROR_PDU' APDUError
            [reserved uint 4    '0x00'                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetConfirmedServiceChoice
                                errorChoice                             ]
            [simple   BACnetError('errorChoice')
                                error                                   ]
        ]
        ['REJECT_PDU' APDUReject
            [reserved uint 4    '0x00'                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetRejectReasonTagged('1')
                                rejectReason                            ]
        ]
        ['ABORT_PDU' APDUAbort
            [reserved uint 3    '0x00'                                  ]
            [simple   bit       server                                  ]
            [simple   uint 8    originalInvokeId                        ]
            [simple   BACnetAbortReasonTagged('1')
                                abortReason                             ]
        ]
        [APDUUnknown
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

[discriminatedType BACnetConfirmedServiceRequest(uint 16 serviceRequestLength)
    [discriminator BACnetConfirmedServiceChoice serviceChoice]
    // we substract serviceChoice from our payload
    [virtual       uint 16  serviceRequestPayloadLength '(serviceRequestLength>0)?(serviceRequestLength - 1):0'    ]
    [typeSwitch serviceChoice
        ////
        // Alarm and Event Services

        ['ACKNOWLEDGE_ALARM' BACnetConfirmedServiceRequestAcknowledgeAlarm
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           acknowledgingProcessIdentifier ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  eventObjectIdentifier          ]
            [simple   BACnetEventStateTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')                     eventStateAcknowledged         ]
            [simple   BACnetTimeStampEnclosed('3')                                                      timestamp                      ]
            [simple   BACnetContextTagCharacterString('4', 'BACnetDataType.CHARACTER_STRING')           acknowledgmentSource           ]
            [simple   BACnetTimeStampEnclosed('5')                                                      timeOfAcknowledgment           ]
        ]
        ['CONFIRMED_COV_NOTIFICATION' BACnetConfirmedServiceRequestConfirmedCOVNotification
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           subscriberProcessIdentifier    ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier     ]
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  monitoredObjectIdentifier      ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')           lifetimeInSeconds              ]
            [simple   BACnetPropertyValues('4', 'monitoredObjectIdentifier.objectType')                 listOfValues                   ]
        ]
        ['CONFIRMED_COV_NOTIFICATION_MULTIPLE' BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           subscriberProcessIdentifier    ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier     ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')           timeRemaining                  ]
            [optional BACnetTimeStampEnclosed('3')                                                      timestamp                      ]
            [simple   ListOfCovNotificationsList('4')                                                   listOfCovNotifications         ]
        ]
        ['CONFIRMED_EVENT_NOTIFICATION' BACnetConfirmedServiceRequestConfirmedEventNotification
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
        ['GET_ENROLLMENT_SUMMARY' BACnetConfirmedServiceRequestGetEnrollmentSummary
            [simple   BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                        acknowledgmentFilter           ]
            [optional BACnetRecipientProcessEnclosed('1')                                               enrollmentFilter               ]
            [optional BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                        eventStateFilter               ]
            [optional BACnetEventTypeTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      eventTypeFilter                ]
            [optional BACnetConfirmedServiceRequestGetEnrollmentSummaryPriorityFilter('4')              priorityFilter                 ]
            [optional BACnetContextTagUnsignedInteger('5', 'BACnetDataType.UNSIGNED_INTEGER')           notificationClassFilter        ]
        ]
        ['GET_EVENT_INFORMATION' BACnetConfirmedServiceRequestGetEventInformation
            [optional BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  lastReceivedObjectIdentifier   ]
        ]
        ['LIFE_SAFETY_OPERATION' BACnetConfirmedServiceRequestLifeSafetyOperation
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           requestingProcessIdentifier    ]
            [simple   BACnetContextTagCharacterString('1', 'BACnetDataType.CHARACTER_STRING')           requestingSource               ]
            [simple   BACnetLifeSafetyOperationTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')            request                        ]
            [optional BACnetContextTagObjectIdentifier('3', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  objectIdentifier               ]
        ]
        ['SUBSCRIBE_COV' BACnetConfirmedServiceRequestSubscribeCOV
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')            subscriberProcessIdentifier  ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')   monitoredObjectIdentifier    ]
            [optional BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')                             issueConfirmed               ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')            lifetimeInSeconds            ]
        ]
        ['SUBSCRIBE_COV_PROPERTY' BACnetConfirmedServiceRequestSubscribeCOVProperty
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')            subscriberProcessIdentifier  ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')   monitoredObjectIdentifier    ]
            [optional BACnetContextTagBoolean('2', 'BACnetDataType.BOOLEAN')                             issueConfirmedNotifications  ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')            lifetime                     ]
            [simple   BACnetPropertyReferenceEnclosed('4')                                               monitoredPropertyIdentifier  ]
            [optional BACnetContextTagReal('5', 'BACnetDataType.REAL')                                   covIncrement                 ]
        ]
        ['SUBSCRIBE_COV_PROPERTY_MULTIPLE' BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple
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

        ['ATOMIC_READ_FILE' BACnetConfirmedServiceRequestAtomicReadFile
            [simple   BACnetApplicationTagObjectIdentifier                                               fileIdentifier               ]
            [simple   BACnetConfirmedServiceRequestAtomicReadFileStreamOrRecord                          accessMethod                 ]
        ]
        ['ATOMIC_WRITE_FILE' BACnetConfirmedServiceRequestAtomicWriteFile
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
        ['ADD_LIST_ELEMENT' BACnetConfirmedServiceRequestAddListElement
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [optional BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value') listOfElements              ]
        ]
        ['REMOVE_LIST_ELEMENT' BACnetConfirmedServiceRequestRemoveListElement
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [optional BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value') listOfElements              ]
        ]
        ['CREATE_OBJECT' BACnetConfirmedServiceRequestCreateObject
            [simple   BACnetConfirmedServiceRequestCreateObjectObjectSpecifier('0')                         objectSpecifier             ]
            [optional BACnetPropertyValues('1', 'objectSpecifier.isObjectType?objectSpecifier.objectType:objectSpecifier.objectIdentifier.objectType')
                                                                                                            listOfValues                ]
        ]
        ['DELETE_OBJECT' BACnetConfirmedServiceRequestDeleteObject
            [simple   BACnetApplicationTagObjectIdentifier                                                  objectIdentifier            ]
        ]
        ['READ_PROPERTY' BACnetConfirmedServiceRequestReadProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
        ]
        ['READ_PROPERTY_MULTIPLE' BACnetConfirmedServiceRequestReadPropertyMultiple(uint 16 serviceRequestPayloadLength)
            [array    BACnetReadAccessSpecification                                                         data
                            length 'serviceRequestPayloadLength'                                                                        ]
        ]
        ['READ_RANGE' BACnetConfirmedServiceRequestReadRange
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               propertyArrayIndex          ]
            // TODO: this attribute should be named range but this is a keyword in golang (so at this point we should build a language translator which makes keywords safe)
            [optional BACnetConfirmedServiceRequestReadRangeRange                                           readRange                   ]
        ]
        ['WRITE_PROPERTY' BACnetConfirmedServiceRequestWriteProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [simple   BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value') propertyValue               ]
            [optional BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')               priority                    ]
        ]
        ['WRITE_PROPERTY_MULTIPLE' BACnetConfirmedServiceRequestWritePropertyMultiple(uint 16 serviceRequestPayloadLength)
            [array    BACnetWriteAccessSpecification                                                        data
                            length 'serviceRequestPayloadLength'                                                                        ]
        ]
        //
        ////

        ////
        // Remote Device Management Services

        ['DEVICE_COMMUNICATION_CONTROL' BACnetConfirmedServiceRequestDeviceCommunicationControl
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')               timeDuration                ]
            [simple   BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisableTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            enableDisable               ]
            [optional BACnetContextTagCharacterString('2', 'BACnetDataType.CHARACTER_STRING')               password                    ]

        ]
        ['CONFIRMED_PRIVATE_TRANSFER' BACnetConfirmedServiceRequestConfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                           vendorId                    ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')               serviceNumber               ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                                                                                                            serviceParameters           ]
        ]
        ['CONFIRMED_TEXT_MESSAGE' BACnetConfirmedServiceRequestConfirmedTextMessage
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      textMessageSourceDevice     ]
            [optional BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass('1')                    messageClass                ]
            [simple   BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            messagePriority             ]
            [simple   BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')               message                     ]
        ]
        ['REINITIALIZE_DEVICE' BACnetConfirmedServiceRequestReinitializeDevice
            [simple   BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            reinitializedStateOfDevice  ]
            [optional BACnetContextTagCharacterString('1', 'BACnetDataType.CHARACTER_STRING')
                                                                                                            password                    ]
        ]

        ////
        //  Virtual Terminal Services

        ['VT_OPEN' BACnetConfirmedServiceRequestVTOpen
            [simple   BACnetVTClassTagged('0', 'TagClass.APPLICATION_TAGS')                                 vtClass                     ]
            [simple   BACnetApplicationTagUnsignedInteger                                                   localVtSessionIdentifier    ]
        ]
        ['VT_CLOSE' BACnetConfirmedServiceRequestVTClose(uint 16 serviceRequestPayloadLength)
            [array    BACnetApplicationTagUnsignedInteger                                                   listOfRemoteVtSessionIdentifiers
                                                               length 'serviceRequestPayloadLength'                                     ]
        ]
        ['VT_DATA' BACnetConfirmedServiceRequestVTData
            [simple   BACnetApplicationTagUnsignedInteger                                                   vtSessionIdentifier         ]
            [simple   BACnetApplicationTagOctetString                                                       vtNewData                   ]
            [simple   BACnetApplicationTagUnsignedInteger                                                   vtDataFlag                  ]
        ]
        //
        ////

        ////
        //  Removed Services

        ['AUTHENTICATE' BACnetConfirmedServiceRequestAuthenticate(uint 16 serviceRequestPayloadLength)
            [array    byte                                                                                  bytesOfRemovedService
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
        ['REQUEST_KEY' BACnetConfirmedServiceRequestRequestKey(uint 16 serviceRequestPayloadLength)
            [array    byte                                                                                  bytesOfRemovedService
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
        ['READ_PROPERTY_CONDITIONAL' BACnetConfirmedServiceRequestReadPropertyConditional(uint 16 serviceRequestPayloadLength)
            [array    byte                                                                                  bytesOfRemovedService
                        length 'serviceRequestPayloadLength'                                                                            ]
        ]
        //
        ////

        [BACnetConfirmedServiceRequestConfirmedUnknown(uint 16 serviceRequestPayloadLength)
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
    [simple   BACnetConstructedData('2', 'objectTypeArgument', 'propertyIdentifier.value')
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
        ['0' BACnetRecipientDevice
            [simple   BACnetApplicationTagObjectIdentifier
                            deviceValue                      ]
        ]
        ['1' BACnetRecipientAddress
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
        ['0' BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassNumeric
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER') numericValue             ]
        ]
        ['1' BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassCharacter
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
        ['0x3' BACnetConfirmedServiceRequestReadRangeRangeByPosition
            [simple   BACnetApplicationTagUnsignedInteger                   referenceIndex            ]
            [simple   BACnetApplicationTagSignedInteger                     count                     ]
        ]
        ['0x6' BACnetConfirmedServiceRequestReadRangeRangeBySequenceNumber
            [simple   BACnetApplicationTagUnsignedInteger                   referenceSequenceNumber   ]
            [simple   BACnetApplicationTagSignedInteger                     count                     ]
        ]
        ['0x7' BACnetConfirmedServiceRequestReadRangeRangeByTime
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
    [optional   BACnetConstructedData('2', 'objectTypeArgument', 'propertyIdentifier.value')
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
        ['I_AM' BACnetUnconfirmedServiceRequestIAm
            [simple   BACnetApplicationTagObjectIdentifier                        deviceIdentifier                ]
            [simple   BACnetApplicationTagUnsignedInteger                         maximumApduLengthAcceptedLength ]
            [simple   BACnetSegmentationTagged('0', 'TagClass.APPLICATION_TAGS')  segmentationSupported           ]
            [simple   BACnetVendorIdTagged('2', 'TagClass.APPLICATION_TAGS')      vendorId                        ]
        ]
        ['I_HAVE' BACnetUnconfirmedServiceRequestIHave
            [simple   BACnetApplicationTagObjectIdentifier                        deviceIdentifier    ]
            [simple   BACnetApplicationTagObjectIdentifier                        objectIdentifier    ]
            [simple   BACnetApplicationTagCharacterString                         objectName          ]
        ]
        ['UNCONFIRMED_COV_NOTIFICATION' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')          subscriberProcessIdentifier ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER') initiatingDeviceIdentifier  ]
            [simple   BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER') monitoredObjectIdentifier   ]
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')          lifetimeInSeconds           ]
            [simple   BACnetPropertyValues('4', 'monitoredObjectIdentifier.objectType')                listOfValues                ]
        ]
        ['UNCONFIRMED_EVENT_NOTIFICATION' BACnetUnconfirmedServiceRequestUnconfirmedEventNotification
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
        ['UNCONFIRMED_PRIVATE_TRANSFER' BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      vendorId                     ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')          serviceNumber                ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE') serviceParameters           ]
        ]
        ['UNCONFIRMED_TEXT_MESSAGE' BACnetUnconfirmedServiceRequestUnconfirmedTextMessage
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      textMessageSourceDevice     ]
            [optional BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass('1')                    messageClass                ] // Note we reuse the once from confirmed here
            [simple   BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                                            messagePriority             ] // Note we reuse the once from confirmed here
            [simple   BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')               message                     ]
        ]
        ['TIME_SYNCHRONIZATION' BACnetUnconfirmedServiceRequestTimeSynchronization
            [simple   BACnetApplicationTagDate synchronizedDate]
            [simple   BACnetApplicationTagTime synchronizedTime]
        ]
        ['WHO_HAS' BACnetUnconfirmedServiceRequestWhoHas
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeLowLimit                                         ]
            [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeHighLimit  'deviceInstanceRangeLowLimit != null' ]
            [optional BACnetContextTagObjectIdentifier('2', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  objectIdentifier                                                    ]
            [optional BACnetContextTagCharacterString('3', 'BACnetDataType.CHARACTER_STRING')           objectName                    'objectIdentifier == null'            ]
        ]
        ['WHO_IS' BACnetUnconfirmedServiceRequestWhoIs
            [optional BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeLowLimit                                                 ]
            [optional BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           deviceInstanceRangeHighLimit  'deviceInstanceRangeLowLimit != null'         ]
        ]
        ['UTC_TIME_SYNCHRONIZATION' BACnetUnconfirmedServiceRequestUTCTimeSynchronization
            [simple   BACnetApplicationTagDate synchronizedDate]
            [simple   BACnetApplicationTagTime synchronizedTime]
        ]
        ['WRITE_GROUP' BACnetUnconfirmedServiceRequestWriteGroup
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           groupNumber                 ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           writePriority               ]
            [simple   BACnetGroupChannelValueList('2')                                                  changeList                  ]
            [optional BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')           inhibitDelay                ]
        ]
        ['UNCONFIRMED_COV_NOTIFICATION_MULTIPLE' BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')           subscriberProcessIdentifier ]
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')  initiatingDeviceIdentifier  ]
            [simple   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')           timeRemaining               ]
            [optional BACnetTimeStampEnclosed('3')                                                      timestamp                   ]
            [simple   ListOfCovNotificationsList('4')                                                   listOfCovNotifications      ]
        ]
        [BACnetUnconfirmedServiceRequestUnconfirmedUnknown
            [array    byte    unknownBytes length '(serviceRequestLength>0)?(serviceRequestLength - 1):0']
        ]
    ]
]

[discriminatedType BACnetServiceAck(uint 16 serviceAckLength)
    [discriminator   BACnetConfirmedServiceChoice
                        serviceChoice                   ]
    // we substract serviceChoice from our payload
    [virtual       uint 16  serviceAckPayloadLength '(serviceAckLength>0)?(serviceAckLength - 1):0'    ]
    [typeSwitch serviceChoice
        ////
        // Alarm and Event Services

        ['GET_ALARM_SUMMARY' BACnetServiceAckGetAlarmSummary
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                ]
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')  eventState                      ]
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS')
                                                                                acknowledgedTransitions         ]
        ]
        ['GET_ENROLLMENT_SUMMARY' BACnetServiceAckGetEnrollmentSummary
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                ]
            [simple   BACnetEventTypeTagged('0', 'TagClass.APPLICATION_TAGS')   eventType                       ]
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')  eventState                      ]
            [simple   BACnetApplicationTagUnsignedInteger                       priority                        ]
            [optional BACnetApplicationTagUnsignedInteger                       notificationClass               ]
        ]
        ['GET_EVENT_INFORMATION' BACnetServiceAckGetEventInformation
            [simple   BACnetEventSummariesList('0')                             listOfEventSummaries            ]
            [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')    moreEvents                      ]
        ]
        //
        ////

        ////
        // File Access Services

        ['ATOMIC_READ_FILE' BACnetServiceAckAtomicReadFile
            [simple   BACnetApplicationTagBoolean                               endOfFile                       ]
            [simple   BACnetServiceAckAtomicReadFileStreamOrRecord              accessMethod                    ]
        ]
        ['ATOMIC_WRITE_FILE' BACnetServiceAckAtomicWriteFile
            [simple   BACnetContextTagSignedInteger('0', 'BACnetDataType.SIGNED_INTEGER') fileStartPosition     ]
        ]
        //
        ////

        ////
        // Object Access Services
        ['CREATE_OBJECT' BACnetServiceAckCreateObject
            [simple   BACnetApplicationTagObjectIdentifier                      objectIdentifier                ]
        ]
        ['READ_PROPERTY' BACnetServiceAckReadProperty
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                                                                objectIdentifier                ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                                                                propertyIdentifier              ]
            [optional   BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')
                                                                                arrayIndex                      ]
            [optional   BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value')
                                                                                values                          ]
        ]
        ['READ_PROPERTY_MULTIPLE' BACnetServiceAckReadPropertyMultiple(uint 16 serviceAckPayloadLength)
            [array    BACnetReadAccessResult                                    data
                            length 'serviceAckPayloadLength'                                                    ]
        ]
        ['READ_RANGE' BACnetServiceAckReadRange
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier    ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier  ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               propertyArrayIndex  ]
            [simple   BACnetResultFlagsTagged('3', 'TagClass.CONTEXT_SPECIFIC_TAGS')                        resultFlags         ]
            [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')               itemCount           ]
            [optional BACnetConstructedData('5', 'objectIdentifier.objectType', 'propertyIdentifier.value') itemData            ]
            [optional BACnetContextTagUnsignedInteger('6', 'BACnetDataType.UNSIGNED_INTEGER')               firstSequenceNumber ]
        ]
        //
        ////


        ////
        // Remote Device Management Services

        ['CONFIRMED_PRIVATE_TRANSFER' BACnetServiceAckConfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                      vendorId                    ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')          serviceNumber               ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE') resultBlock                 ]
        ]
        //
        ////

        ////
        //  Virtual Terminal Services

        ['VT_OPEN' BACnetServiceAckVTOpen
            [simple   BACnetApplicationTagUnsignedInteger                       remoteVtSessionIdentifier                        ]
        ]
        ['VT_DATA' BACnetServiceAckVTData
            [simple   BACnetApplicationTagUnsignedInteger                       vtSessionIdentifier                              ]
            [simple   BACnetApplicationTagOctetString                           vtNewData                                        ]
            [simple   BACnetApplicationTagUnsignedInteger                       vtDataFlag                                       ]
        ]
        //
        ////


        ////
        //  Removed Services

        ['AUTHENTICATE' BACnetServiceAckAuthenticate(uint 16 serviceAckPayloadLength)
            [array    byte                                                      bytesOfRemovedService
                        length 'serviceAckPayloadLength'                                                                ]
        ]
        ['REQUEST_KEY' BACnetServiceAckRequestKey(uint 16 serviceAckPayloadLength)
            [array    byte                                                      bytesOfRemovedService
                        length 'serviceAckPayloadLength'                                                                ]
        ]
        ['READ_PROPERTY_CONDITIONAL' BACnetServiceAckReadPropertyConditional(uint 16 serviceAckPayloadLength)
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

[type BACnetEventMessageTexts
    [simple  BACnetApplicationTagCharacterString
                    toOffnormalText                 ]
    [simple  BACnetApplicationTagCharacterString
                    toFaultText                     ]
    [simple  BACnetApplicationTagCharacterString
                    toNormalText                    ]
]

[type BACnetEventMessageTextsConfig
    [simple  BACnetApplicationTagCharacterString
                    toOffnormalTextConfig           ]
    [simple  BACnetApplicationTagCharacterString
                    toFaultTextConfig               ]
    [simple  BACnetApplicationTagCharacterString
                    toNormalTextConfig              ]
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' BACnetChannelValueNull
           [simple  BACnetApplicationTagNull
                        nullValue                                                  ]
       ]
       ['0x4', 'false' BACnetChannelValueReal
           [simple  BACnetApplicationTagReal
                        realValue                                                  ]
       ]
       ['0x9', 'false' BACnetChannelValueEnumerated
           [simple   BACnetApplicationTagEnumerated
                       enumeratedValue                                             ]
       ]
       ['0x2', 'false' BACnetChannelValueUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                       unsignedValue                                               ]
       ]
       ['0x1', 'false' BACnetChannelValueBoolean
           [simple   BACnetApplicationTagBoolean
                       booleanValue                                                ]
       ]
       ['0x3', 'false' BACnetChannelValueInteger
           [simple   BACnetApplicationTagSignedInteger
                       integerValue                                                ]
       ]
       ['0x5', 'false' BACnetChannelValueDouble
           [simple  BACnetApplicationTagDouble
                        doubleValue                                                ]
       ]
       ['0xB', 'false' BACnetChannelValueTime
           [simple   BACnetApplicationTagTime
                       timeValue                                                   ]
       ]
       ['0x7', 'false' BACnetChannelValueCharacterString
           [simple   BACnetApplicationTagCharacterString
                       characterStringValue                                        ]
       ]
       ['0x6', 'false' BACnetChannelValueOctetString
           [simple   BACnetApplicationTagOctetString
                       octetStringValue                                            ]
       ]
       ['0x8', 'false' BACnetChannelValueBitString
           [simple   BACnetApplicationTagBitString
                       bitStringValue                                              ]
       ]
       ['0xA', 'false' BACnetChannelValueDate
           [simple   BACnetApplicationTagDate
                       dateValue                                                   ]
       ]
       ['0xC', 'false' BACnetChannelValueObjectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                       objectidentifierValue                                       ]
       ]
       ['0', 'true' BACnetChannelValueLightingCommand
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
                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)'    ]
    [simple   BACnetClosingTag('tagNumber')
                     closingTag                                                                 ]
]

[type BACnetReadAccessProperty(BACnetObjectType objectTypeArgument)
    [simple   BACnetPropertyIdentifierTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                    propertyIdentifier                                                          ]
    [optional   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                    arrayIndex                                                                  ]
    [optional   BACnetReadAccessPropertyReadResult('objectTypeArgument', 'propertyIdentifier.value')
                    readResult                                                                  ]
]

[type BACnetReadAccessPropertyReadResult(BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument)
    [peek       BACnetTagHeader
                            peekedTagHeader                                                     ]
    [virtual    uint 8      peekedTagNumber     'peekedTagHeader.actualTagNumber'               ]
    [optional   BACnetConstructedData('4', 'objectTypeArgument', 'propertyIdentifierArgument')
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
            [optional BACnetConstructedData('3', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                            errorParameters             ]
        ]
        ['VT_CLOSE'                         VTCloseError
            [simple   ErrorEnclosed('0')
                            errorType                   ]
            [optional VTCloseErrorListOfVTSessionIdentifiers('1')
                            listOfVtSessionIdentifiers  ]
        ]
        [BACnetErrorGeneral
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
        ['0' BACnetNotificationParametersChangeOfBitString(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetContextTagBitString('0', 'BACnetDataType.BIT_STRING')
                            changeOfBitString                                       ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['1' BACnetNotificationParametersChangeOfState(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetPropertyStatesEnclosed('0')
                            changeOfState                                           ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['2' BACnetNotificationParametersChangeOfValue(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetNotificationParametersChangeOfValueNewValue('0')
                            newValue                                                ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['3' BACnetNotificationParametersCommandFailure(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                            commandValue                                            ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetConstructedData('2', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                            feedbackValue                                           ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['4' BACnetNotificationParametersFloatingLimit(uint 8 peekedTagNumber)
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
        ['5' BACnetNotificationParametersOutOfRange(uint 8 peekedTagNumber)
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
        ['6' BACnetNotificationParametersComplexEventType(uint 8 peekedTagNumber)
            [simple   BACnetPropertyValues('peekedTagNumber', 'objectTypeArgument')
                            listOfValues                                            ]
        ]
        // 7 is deprecated
        ['8' BACnetNotificationParametersChangeOfLifeSafety(uint 8 peekedTagNumber)
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
        ['9' BACnetNotificationParametersExtended(uint 8 peekedTagNumber)
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
        ['10' BACnetNotificationParametersBufferReady(uint 8 peekedTagNumber)
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
        ['11' BACnetNotificationParametersUnsignedRange(uint 8 peekedTagNumber)
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
        ['13' BACnetNotificationParametersAccessEvent(uint 8 peekedTagNumber)
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
        ['14' BACnetNotificationParametersDoubleOutOfRange(uint 8 peekedTagNumber)
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
        ['15' BACnetNotificationParametersSignedOutOfRange(uint 8 peekedTagNumber)
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
        ['16' BACnetNotificationParametersUnsignedOutOfRange(uint 8 peekedTagNumber)
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
        ['17' BACnetNotificationParametersChangeOfCharacterString(uint 8 peekedTagNumber)
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
        ['18' BACnetNotificationParametersChangeOfStatusFlags(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                            presentValue                                            ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            referencedFlags                                         ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['19' BACnetNotificationParametersChangeOfReliability(uint 8 peekedTagNumber)
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
        ['21' BACnetNotificationParametersChangeOfDiscreteValue(uint 8 peekedTagNumber)
            [simple   BACnetOpeningTag('peekedTagNumber')
                            innerOpeningTag                                         ]
            [simple   BACnetNotificationParametersChangeOfDiscreteValueNewValue('0')
                            newValue                                                ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetClosingTag('peekedTagNumber')
                            innerClosingTag                                         ]
        ]
        ['22' BACnetNotificationParametersChangeOfTimer(uint 8 peekedTagNumber)
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' BACnetNotificationParametersChangeOfValueNewValueChangedBits(uint 8 peekedTagNumber)
            [simple   BACnetContextTagBitString('0', 'BACnetDataType.BIT_STRING')
                        changedBits                                                 ]
        ]
        ['1' BACnetNotificationParametersChangeOfValueNewValueChangedValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagReal('0', 'BACnetDataType.REAL')
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
    [optional BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifier.value')                propertyValue       ]
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
   [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
   [virtual bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
   [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
   [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x1', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueBoolean
           [simple   BACnetApplicationTagBoolean
                       booleanValue                                                ]
       ]
       ['0x2', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                       unsignedValue                                               ]
       ]
       ['0x3', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger
           [simple   BACnetApplicationTagSignedInteger
                       integerValue                                                ]
       ]
       ['0x9', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueEnumerated
           [simple   BACnetApplicationTagEnumerated
                       enumeratedValue                                             ]
       ]
       ['0x7', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueCharacterString
           [simple   BACnetApplicationTagCharacterString
                       characterStringValue                                        ]
       ]
       ['0x6', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueOctetString
           [simple   BACnetApplicationTagOctetString
                       octetStringValue                                            ]
       ]
       ['0xA', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueOctetDate
           [simple   BACnetApplicationTagDate
                       dateValue                                                   ]
       ]
       ['0xB', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueOctetTime
           [simple   BACnetApplicationTagTime
                       timeValue                                                   ]
       ]
       ['0xC', 'false' BACnetNotificationParametersChangeOfDiscreteValueNewValueObjectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                       objectidentifierValue                                       ]
       ]
       ['0', 'true' BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime
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
    [optional BACnetConstructedData('4', 'objectIdentifier.objectType', 'propertyIdentifier.value')
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

// Note per spec this should be 16 but we reuse this for index access and non conformant transmission
[type BACnetPriorityArray(BACnetObjectType objectTypeArgument, uint 8 tagNumber)
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
    [virtual  bit                   isValidPriorityArray    'COUNT(data) == 16'                 ]
    [virtual  bit                   isIndexedAccess         'COUNT(data) == 1'                  ]
    [virtual  BACnetPriorityValue   indexEntry              'priorityValue01'                   ]
]

[type BACnetPriorityValue(BACnetObjectType objectTypeArgument)
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' BACnetPriorityValueNull
           [simple  BACnetApplicationTagNull
                            nullValue                                                   ]
       ]
       ['0x4', 'false' BACnetPriorityValueReal
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x9', 'false' BACnetPriorityValueEnumerated
           [simple   BACnetApplicationTagEnumerated
                            enumeratedValue                                             ]
       ]
       ['0x2', 'false' BACnetPriorityValueUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x1', 'false' BACnetPriorityValueBoolean
           [simple   BACnetApplicationTagBoolean
                            booleanValue                                                ]
       ]
       ['0x3', 'false' BACnetPriorityValueInteger
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
       ['0x5', 'false' BACnetPriorityValueDouble
           [simple  BACnetApplicationTagDouble
                                doubleValue                                             ]
       ]
       ['0xB', 'false' BACnetPriorityValueTime
           [simple   BACnetApplicationTagTime
                            timeValue                                                   ]
       ]
       ['0x7', 'false' BACnetPriorityValueCharacterString
           [simple   BACnetApplicationTagCharacterString
                            characterStringValue                                        ]
       ]
       ['0x6', 'false' BACnetPriorityValueOctetString
           [simple   BACnetApplicationTagOctetString
                            octetStringValue                                            ]
       ]
       ['0x8', 'false' BACnetPriorityValueBitString
           [simple   BACnetApplicationTagBitString
                            bitStringValue                                              ]
       ]
       ['0xA', 'false' BACnetPriorityValueDate
           [simple   BACnetApplicationTagDate
                            dateValue                                                   ]
       ]
       ['0xC', 'false' BACnetPriorityValueObjectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                            objectidentifierValue                                       ]
       ]
       ['0', 'true' BACnetPriorityValueConstructedValue
           [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                            constructedValue                                            ]
       ]
       ['1', 'true' BACnetPriorityValueDateTime
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
    [virtual uint 8 peekedTagNumber
                        'peekedTagHeader.actualTagNumber'       ]
    [typeSwitch peekedTagNumber
        ['0' BACnetPropertyStatesBoolean(uint 8 peekedTagNumber)
            [simple   BACnetContextTagBoolean('peekedTagNumber', 'BACnetDataType.BOOLEAN')
                                booleanValue                    ]
        ]
        ['1' BACnetPropertyStatesBinaryValue(uint 8 peekedTagNumber)
            [simple   BACnetBinaryPVTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                binaryValue                     ]
        ]
        ['2' BACnetPropertyStatesEventType(uint 8 peekedTagNumber)
            [simple   BACnetEventTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                eventType                       ]
        ]
        ['3' BACnetPropertyStatesPolarity(uint 8 peekedTagNumber)
            [simple   BACnetPolarityTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                polarity                        ]
        ]
        ['4' BACnetPropertyStatesProgramChange(uint 8 peekedTagNumber)
            [simple   BACnetProgramRequestTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                programChange                   ]
        ]
        ['5' BACnetPropertyStatesProgramChange(uint 8 peekedTagNumber)
            [simple   BACnetProgramStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                programState                    ]
        ]
        ['6' BACnetPropertyStatesReasonForHalt(uint 8 peekedTagNumber)
            [simple   BACnetProgramErrorTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                reasonForHalt                   ]
        ]
        ['7' BACnetPropertyStatesReliability(uint 8 peekedTagNumber)
            [simple   BACnetReliabilityTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                reliability                     ]
        ]
        ['8' BACnetPropertyStatesState(uint 8 peekedTagNumber)
            [simple   BACnetEventStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                state                           ]
        ]
        ['9' BACnetPropertyStatesSystemStatus(uint 8 peekedTagNumber)
            [simple   BACnetDeviceStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                systemStatus                    ]
        ]
        ['10' BACnetPropertyStatesUnits(uint 8 peekedTagNumber)
            [simple   BACnetEngineeringUnitsTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                units                           ]
        ]
        ['11' BACnetPropertyStatesExtendedValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagUnsignedInteger('peekedTagNumber', 'BACnetDataType.UNSIGNED_INTEGER')
                                unsignedValue                   ]
        ]
        ['12' BACnetPropertyStatesLifeSafetyMode(uint 8 peekedTagNumber)
            [simple   BACnetLifeSafetyModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lifeSafetyMode                  ]
        ]
        ['13' BACnetPropertyStatesLifeSafetyState(uint 8 peekedTagNumber)
            [simple   BACnetLifeSafetyStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lifeSafetyState                 ]
        ]
        ['14' BACnetPropertyStatesRestartReason(uint 8 peekedTagNumber)
            [simple   BACnetRestartReasonTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                restartReason                   ]
        ]
        ['15' BACnetPropertyStatesDoorAlarmState(uint 8 peekedTagNumber)
            [simple   BACnetDoorAlarmStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorAlarmState                  ]
        ]
        ['16' BACnetPropertyStatesAction(uint 8 peekedTagNumber)
            [simple   BACnetActionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                action                          ]
        ]
        ['17' BACnetPropertyStatesDoorSecuredStatus(uint 8 peekedTagNumber)
            [simple   BACnetDoorSecuredStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorSecuredStatus               ]
        ]
        ['18' BACnetPropertyStatesDoorStatus(uint 8 peekedTagNumber)
            [simple   BACnetDoorStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorStatus                      ]
        ]
        ['19' BACnetPropertyStatesDoorValue(uint 8 peekedTagNumber)
            [simple   BACnetDoorValueTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                doorValue                       ]
        ]
        ['20' BACnetPropertyStatesFileAccessMethod(uint 8 peekedTagNumber)
            [simple   BACnetFileAccessMethodTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                fileAccessMethod                ]
        ]
        ['21' BACnetPropertyStatesLockStatus(uint 8 peekedTagNumber)
            [simple   BACnetLockStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lockStatus                      ]
        ]
        ['22' BACnetPropertyStatesLifeSafetyOperations(uint 8 peekedTagNumber)
            [simple   BACnetLifeSafetyOperationTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lifeSafetyOperations            ]
        ]
        ['23' BACnetPropertyStatesMaintenance(uint 8 peekedTagNumber)
            [simple   BACnetMaintenanceTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                maintenance                     ]
        ]
        ['24' BACnetPropertyStatesNodeType(uint 8 peekedTagNumber)
            [simple   BACnetNodeTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                nodeType                        ]
        ]
        ['25' BACnetPropertyStatesNotifyType(uint 8 peekedTagNumber)
            [simple   BACnetNotifyTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                notifyType                      ]
        ]
        ['26' BACnetPropertyStatesSecurityLevel(uint 8 peekedTagNumber)
            [simple   BACnetSecurityLevelTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                securityLevel                   ]
        ]
        ['27' BACnetPropertyStatesShedState(uint 8 peekedTagNumber)
            [simple   BACnetShedStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                shedState                       ]
        ]
        ['28' BACnetPropertyStatesSilencedState(uint 8 peekedTagNumber)
            [simple   BACnetSilencedStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                silencedState                   ]
        ]
        //['29' BACnetPropertyStatesReserved(uint 8 peekedTagNumber) ]
        ['30' BACnetPropertyStatesAccessEvent(uint 8 peekedTagNumber)
            [simple   BACnetAccessEventTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accessEvent                     ]
        ]
        ['31' BACnetPropertyStatesZoneOccupanyState(uint 8 peekedTagNumber)
            [simple   BACnetAccessZoneOccupancyStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                zoneOccupanyState               ]
        ]
        ['32' BACnetPropertyStatesAccessCredentialDisableReason(uint 8 peekedTagNumber)
            [simple   BACnetAccessCredentialDisableReasonTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accessCredentialDisableReason   ]
        ]
        ['33' BACnetPropertyStatesAccessCredentialDisable(uint 8 peekedTagNumber)
            [simple   BACnetAccessCredentialDisableTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                accessCredentialDisable         ]
        ]
        ['34' BACnetPropertyStatesAuthenticationStatus(uint 8 peekedTagNumber)
            [simple   BACnetAuthenticationStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                authenticationStatus            ]
        ]
        // 35 is undefined
        ['36' BACnetPropertyStatesBackupState(uint 8 peekedTagNumber)
            [simple    BACnetBackupStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                backupState                     ]
        ]
        ['37' BACnetPropertyStatesWriteStatus(uint 8 peekedTagNumber)
            [simple    BACnetWriteStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                writeStatus                     ]
        ]
        ['38' BACnetPropertyStatesLightningInProgress(uint 8 peekedTagNumber)
            [simple    BACnetLightingInProgressTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lightningInProgress             ]
        ]
        ['39' BACnetPropertyStatesLightningOperation(uint 8 peekedTagNumber)
            [simple    BACnetLightingOperationTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lightningOperation              ]
        ]
        ['40' BACnetPropertyStatesLightningTransition(uint 8 peekedTagNumber)
            [simple    BACnetLightingTransitionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                lightningTransition             ]
        ]
        ['41' BACnetPropertyStatesIntegerValue(uint 8 peekedTagNumber)
            [simple   BACnetContextTagSignedInteger('peekedTagNumber', 'BACnetDataType.SIGNED_INTEGER')
                                integerValue                    ]
        ]
        ['42' BACnetPropertyStatesBinaryLightningValue(uint 8 peekedTagNumber)
            [simple   BACnetBinaryLightingPVTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                binaryLightningValue            ]
        ]
        ['43' BACnetPropertyStatesTimerState(uint 8 peekedTagNumber)
            [simple   BACnetTimerStateTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                timerState                      ]
        ]
        ['44' BACnetPropertyStatesTimerTransition(uint 8 peekedTagNumber)
            [simple   BACnetTimerTransitionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                timerTransition                 ]
        ]
        ['45' BACnetPropertyStatesBacnetIpMode(uint 8 peekedTagNumber)
            [simple   BACnetIPModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                bacnetIpMode                    ]
        ]
        ['46' BACnetPropertyStatesNetworkPortCommand(uint 8 peekedTagNumber)
            [simple   BACnetNetworkPortCommandTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                networkPortCommand              ]
        ]
        ['47' BACnetPropertyStatesNetworkType(uint 8 peekedTagNumber)
            [simple   BACnetNetworkTypeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                networkType                     ]
        ]
        ['48' BACnetPropertyStatesNetworkNumberQuality(uint 8 peekedTagNumber)
            [simple   BACnetNetworkNumberQualityTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                networkNumberQuality            ]
        ]
        ['49' BACnetPropertyStatesEscalatorOperationDirection(uint 8 peekedTagNumber)
            [simple   BACnetEscalatorOperationDirectionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                escalatorOperationDirection     ]
        ]
        ['50' BACnetPropertyStatesEscalatorFault(uint 8 peekedTagNumber)
            [simple   BACnetEscalatorFaultTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                escalatorFault                  ]
        ]
        ['51' BACnetPropertyStatesEscalatorMode(uint 8 peekedTagNumber)
            [simple   BACnetEscalatorModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                escalatorMode                   ]
        ]
        ['52' BACnetPropertyStatesLiftCarDirection(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarDirectionTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarDirection                ]
        ]
        ['53' BACnetPropertyStatesLiftCarDoorCommand(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarDoorCommandTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarDoorCommand              ]
        ]
        ['54' BACnetPropertyStatesLiftCarDriveStatus(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarDriveStatusTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarDriveStatus              ]
        ]
        ['55' BACnetPropertyStatesLiftCarMode(uint 8 peekedTagNumber)
            [simple   BACnetLiftCarModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftCarMode                     ]
        ]
        ['56' BACnetPropertyStatesLiftGroupMode(uint 8 peekedTagNumber)
            [simple   BACnetLiftGroupModeTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftGroupMode                   ]
        ]
        ['57' BACnetPropertyStatesLiftFault(uint 8 peekedTagNumber)
            [simple   BACnetLiftFaultTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                liftFault                       ]
        ]
        ['58' BACnetPropertyStatesProtocolLevel(uint 8 peekedTagNumber)
            [simple   BACnetProtocolLevelTagged('peekedTagNumber', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                protocolLevel                   ]
        ]
        // 59 undefined
        // 60 undefined
        // 61 undefined
        // 62 undefined
        ['63' BACnetPropertyStatesExtendedValue(uint 8 peekedTagNumber)
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
    [virtual uint 8     peekedTagNumber
                            'peekedTagHeader.actualTagNumber'   ]
    [typeSwitch peekedTagNumber
        ['0' BACnetTimeStampTime
            [simple   BACnetContextTagTime('0', 'BACnetDataType.TIME')
                            timeValue                           ]
        ]
        ['1' BACnetTimeStampSequence
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                            sequenceNumber                      ]
        ]
        ['2' BACnetTimeStampDateTime
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
    [simple   BACnetRecipient
                    recipient                   ]
    [simple   BACnetClosingTag('tagNumber')
                    closingTag                  ]
]

[type BACnetConstructedData(uint 8 tagNumber, BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument)
    [simple   BACnetOpeningTag('tagNumber')
                        openingTag                                                                              ]
    [typeSwitch objectTypeArgument, propertyIdentifierArgument
        [*, 'ABSENTEE_LIMIT'                          BACnetConstructedDataAbsenteeLimit
            [simple   BACnetApplicationTagUnsignedInteger                     absenteeLimit                             ]
        ]
        [*, 'ACCEPTED_MODES'                          BACnetConstructedDataAcceptedModes
            [array    BACnetLifeSafetyModeTagged('0', 'TagClass.APPLICATION_TAGS')
                            acceptedModes              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_ALARM_EVENTS'                     BACnetConstructedDataAccessAlarmEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                    accessAlarmEvents
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_DOORS'                            BACnetConstructedDataAccessDoors
            [array    BACnetDeviceObjectReference
                                accessDoors
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_EVENT'                            BACnetConstructedDataAccessEvent
            [simple   BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS') accessEvent]
        ]
        [*, 'ACCESS_EVENT_AUTHENTICATION_FACTOR'      BACnetConstructedDataAccessEventAuthenticationFactor
            [simple   BACnetAuthenticationFactor  accessEventAuthenticationFactor               ]
        ]
        [*, 'ACCESS_EVENT_CREDENTIAL'                 BACnetConstructedDataAccessEventCredential
            [simple   BACnetDeviceObjectReference       accessEventCredential                                           ]
        ]
        [*, 'ACCESS_EVENT_TAG'                        BACnetConstructedDataAccessEventTag
            [simple BACnetApplicationTagUnsignedInteger                     accessEventTag                              ]
        ]
        [*, 'ACCESS_EVENT_TIME'                       BACnetConstructedDataAccessEventTime
            [simple   BACnetTimeStamp                                                accessEventTime                    ]
        ]
        [*, 'ACCESS_TRANSACTION_EVENTS'               BACnetConstructedDataAccessTransactionEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                        accessTransactionEvents
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCOMPANIMENT'                           BACnetConstructedDataAccompaniment
            [simple   BACnetDeviceObjectReference       accompaniment                                                   ]
        ]
        [*, 'ACCOMPANIMENT_TIME'                      BACnetConstructedDataAccompanimentTime
            [simple   BACnetApplicationTagUnsignedInteger                               accompanimentTime               ]
        ]
        //[*, 'ACK_REQUIRED'                            BACnetConstructedDataAckRequired [validation    '1 == 2'    "TODO: implement me ACK_REQUIRED BACnetConstructedDataAckRequired"]]
        [*, 'ACKED_TRANSITIONS'                       BACnetConstructedDataAckedTransitions
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS') ackedTransitions                ]
        ]
        ['LOOP', 'ACTION'                             BACnetConstructedDataLoopAction
            [simple   BACnetActionTagged('0', 'TagClass.APPLICATION_TAGS') action]
        ]
        [*, 'ACTION'                                  BACnetConstructedDataAction
            [array    BACnetActionList
                            actionLists
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'ACTION_TEXT'                             BACnetConstructedDataActionText [validation    '1 == 2'    "TODO: implement me ACTION_TEXT BACnetConstructedDataActionText"]]
        [*, 'ACTIVATION_TIME'                         BACnetConstructedDataActivationTime
            [simple   BACnetDateTime    activationTime                                                                  ]
        ]
        [*, 'ACTIVE_AUTHENTICATION_POLICY'            BACnetConstructedDataActiveAuthenticationPolicy
            [simple   BACnetApplicationTagUnsignedInteger                               activeAuthenticationPolicy      ]
        ]
        //[*, 'ACTIVE_COV_MULTIPLE_SUBSCRIPTIONS'       BACnetConstructedDataActiveCOVMultipleSubscriptions [validation    '1 == 2'    "TODO: implement me ACTIVE_COV_MULTIPLE_SUBSCRIPTIONS BACnetConstructedDataActiveCOVMultipleSubscriptions"]]
        //[*, 'ACTIVE_COV_SUBSCRIPTIONS'                BACnetConstructedDataActiveCOVSubscriptions [validation    '1 == 2'    "TODO: implement me ACTIVE_COV_SUBSCRIPTIONS BACnetConstructedDataActiveCOVSubscriptions"]]
        //[*, 'ACTIVE_TEXT'                             BACnetConstructedDataActiveText [validation    '1 == 2'    "TODO: implement me ACTIVE_TEXT BACnetConstructedDataActiveText"]]
        //[*, 'ACTIVE_VT_SESSIONS'                      BACnetConstructedDataActiveVTSessions [validation    '1 == 2'    "TODO: implement me ACTIVE_VT_SESSIONS BACnetConstructedDataActiveVTSessions"]]
        //[*, 'ACTUAL_SHED_LEVEL'                       BACnetConstructedDataActualShedLevel [validation    '1 == 2'    "TODO: implement me ACTUAL_SHED_LEVEL BACnetConstructedDataActualShedLevel"]]
        //[*, 'ADJUST_VALUE'                            BACnetConstructedDataAdjustValue [validation    '1 == 2'    "TODO: implement me ADJUST_VALUE BACnetConstructedDataAdjustValue"]]
        //[*, 'ALARM_VALUE'                             BACnetConstructedDataAlarmValue [validation    '1 == 2'    "TODO: implement me ALARM_VALUE BACnetConstructedDataAlarmValue"]]
        ['LIFE_SAFETY_POINT', 'ALARM_VALUES'                            BACnetConstructedDataLifeSafetyPointAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'ALARM_VALUES'                            BACnetConstructedDataAlarmValues [validation    '1 == 2'    "TODO: implement me ALARM_VALUES BACnetConstructedDataAlarmValues"]]
        //[*, 'ALIGN_INTERVALS'                         BACnetConstructedDataAlignIntervals [validation    '1 == 2'    "TODO: implement me ALIGN_INTERVALS BACnetConstructedDataAlignIntervals"]]

        /////
        // All property implementations for every object

        ['ACCESS_CREDENTIAL'     , 'ALL'              BACnetConstructedDataAccessCredentialAl
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_DOOR'           , 'ALL'              BACnetConstructedDataAccessDoorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_POINT'          , 'ALL'              BACnetConstructedDataAccessPointAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_RIGHTS'         , 'ALL'              BACnetConstructedDataAccessRightsAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_USER'           , 'ALL'              BACnetConstructedDataAccessUserAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_ZONE'           , 'ALL'              BACnetConstructedDataAccessZoneAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCUMULATOR'           , 'ALL'              BACnetConstructedDataAccumulatorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ALERT_ENROLLMENT'      , 'ALL'              BACnetConstructedDataAlertEnrollmentAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_INPUT'          , 'ALL'              BACnetConstructedDataAnalogInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_OUTPUT'         , 'ALL'              BACnetConstructedDataAnalogOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_VALUE'          , 'ALL'              BACnetConstructedDataAnalogValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['AVERAGING'             , 'ALL'              BACnetConstructedDataAveragingAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_INPUT'          , 'ALL'              BACnetConstructedDataBinaryInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'ALL'              BACnetConstructedDataBinaryLightingOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_OUTPUT'         , 'ALL'              BACnetConstructedDataBinaryOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_VALUE'          , 'ALL'              BACnetConstructedDataBinaryValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BITSTRING_VALUE'       , 'ALL'              BACnetConstructedDataBitstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CALENDAR'              , 'ALL'              BACnetConstructedDataCalendarAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CHANNEL'               , 'ALL'              BACnetConstructedDataChannelAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CHARACTERSTRING_VALUE' , 'ALL'              BACnetConstructedDataCharacterstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['COMMAND'               , 'ALL'              BACnetConstructedDataCommandAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CREDENTIAL_DATA_INPUT' , 'ALL'              BACnetConstructedDataCredentialDataInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATEPATTERN_VALUE'     , 'ALL'              BACnetConstructedDataDatepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATE_VALUE'            , 'ALL'              BACnetConstructedDataDateValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATETIMEPATTERN_VALUE' , 'ALL'              BACnetConstructedDataDatetimepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATETIME_VALUE'        , 'ALL'              BACnetConstructedDataDatetimeValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DEVICE'                , 'ALL'              BACnetConstructedDataDeviceAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ELEVATOR_GROUP'        , 'ALL'              BACnetConstructedDataElevatorGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ESCALATOR'             , 'ALL'              BACnetConstructedDataEscalatorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['EVENT_ENROLLMENT'      , 'ALL'              BACnetConstructedDataEventEnrollmentAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['EVENT_LOG'             , 'ALL'              BACnetConstructedDataEventLogAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['FILE'                  , 'ALL'              BACnetConstructedDataFileAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['GLOBAL_GROUP'          , 'ALL'              BACnetConstructedDataGlobalGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['GROUP'                 , 'ALL'              BACnetConstructedDataGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['INTEGER_VALUE'         , 'ALL'              BACnetConstructedDataIntegerValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LARGE_ANALOG_VALUE'    , 'ALL'              BACnetConstructedDataLargeAnalogValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFE_SAFETY_POINT'     , 'ALL'              BACnetConstructedDataLifeSafetyPointAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFE_SAFETY_ZONE'      , 'ALL'              BACnetConstructedDataLifeSafetyZoneAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFT'                  , 'ALL'              BACnetConstructedDataLiftAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIGHTING_OUTPUT'       , 'ALL'              BACnetConstructedDataLightingOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LOAD_CONTROL'          , 'ALL'              BACnetConstructedDataLoadControlAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LOOP'                  , 'ALL'              BACnetConstructedDataLoopAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_INPUT'     , 'ALL'              BACnetConstructedDataMultiStateInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_OUTPUT'    , 'ALL'              BACnetConstructedDataMultiStateOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_VALUE'     , 'ALL'              BACnetConstructedDataMultiStateValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NETWORK_PORT'          , 'ALL'              BACnetConstructedDataNetworkPortAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NETWORK_SECURITY'      , 'ALL'              BACnetConstructedDataNetworkSecurityAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NOTIFICATION_CLASS'    , 'ALL'              BACnetConstructedDataNotificationClassAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NOTIFICATION_FORWARDER', 'ALL'              BACnetConstructedDataNotificationForwarderAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['OCTETSTRING_VALUE'     , 'ALL'              BACnetConstructedDataOctetstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['POSITIVE_INTEGER_VALUE', 'ALL'              BACnetConstructedDataPositiveIntegerValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['PROGRAM'               , 'ALL'              BACnetConstructedDataProgramAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['PULSE_CONVERTER'       , 'ALL'              BACnetConstructedDataPulseConverterAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['SCHEDULE'              , 'ALL'              BACnetConstructedDataScheduleAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['STRUCTURED_VIEW'       , 'ALL'              BACnetConstructedDataStructuredViewAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIMEPATTERN_VALUE'     , 'ALL'              BACnetConstructedDataTimepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIME_VALUE'            , 'ALL'              BACnetConstructedDataTimeValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIMER'                 , 'ALL'              BACnetConstructedDataTimerAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TREND_LOG'             , 'ALL'              BACnetConstructedDataTrendLogAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TREND_LOG_MULTIPLE'    , 'ALL'              BACnetConstructedDataTrendLogMultipleAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        //
        /////

        //[*, 'ALL_WRITES_SUCCESSFUL'                   BACnetConstructedDataAllWritesSuccessful [validation    '1 == 2'    "TODO: implement me ALL_WRITES_SUCCESSFUL BACnetConstructedDataAllWritesSuccessful"]]
        //[*, 'ALLOW_GROUP_DELAY_INHIBIT'               BACnetConstructedDataAllowGroupDelayInhibit [validation    '1 == 2'    "TODO: implement me ALLOW_GROUP_DELAY_INHIBIT BACnetConstructedDataAllowGroupDelayInhibit"]]
        //[*, 'APDU_LENGTH'                             BACnetConstructedDataAPDULength [validation    '1 == 2'    "TODO: implement me APDU_LENGTH BACnetConstructedDataAPDULength"]]
        //[*, 'APDU_SEGMENT_TIMEOUT'                    BACnetConstructedDataApduSegmentTimeout [validation    '1 == 2'    "TODO: implement me APDU_SEGMENT_TIMEOUT BACnetConstructedDataApduSegmentTimeout"]]
        //[*, 'APDU_TIMEOUT'                            BACnetConstructedDataAPDUTimeout [validation    '1 == 2'    "TODO: implement me APDU_TIMEOUT BACnetConstructedDataAPDUTimeout"]]
        //[*, 'APPLICATION_SOFTWARE_VERSION'            BACnetConstructedDataApplicationSoftwareVersion [validation    '1 == 2'    "TODO: implement me APPLICATION_SOFTWARE_VERSION BACnetConstructedDataApplicationSoftwareVersion"]]
        //[*, 'ARCHIVE'                                 BACnetConstructedDataArchive [validation    '1 == 2'    "TODO: implement me ARCHIVE BACnetConstructedDataArchive"]]
        [*, 'ASSIGNED_ACCESS_RIGHTS'                  BACnetConstructedDataAssignedAccessRights
            [array    BACnetAssignedAccessRights
                                        assignedAccessRights
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'ASSIGNED_LANDING_CALLS'                  BACnetConstructedDataAssignedLandingCalls [validation    '1 == 2'    "TODO: implement me ASSIGNED_LANDING_CALLS BACnetConstructedDataAssignedLandingCalls"]]
        //[*, 'ATTEMPTED_SAMPLES'                       BACnetConstructedDataAttemptedSamples [validation    '1 == 2'    "TODO: implement me ATTEMPTED_SAMPLES BACnetConstructedDataAttemptedSamples"]]
        [*, 'AUTHENTICATION_FACTORS'                  BACnetConstructedDataAuthenticationFactors
            [array    BACnetCredentialAuthenticationFactor
                            authenticationFactors
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHENTICATION_POLICY_LIST'              BACnetConstructedDataAuthenticationPolicyList
            [array    BACnetAuthenticationPolicy
                            authenticationPolicyList
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHENTICATION_POLICY_NAMES'             BACnetConstructedDataAuthenticationPolicyNames
            [array    BACnetApplicationTagCharacterString
                                        authenticationPolicyNames
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHENTICATION_STATUS'                   BACnetConstructedDataAuthenticationStatus
            [simple   BACnetAuthenticationStatusTagged('0', 'TagClass.APPLICATION_TAGS') authenticationStatus           ]
        ]
        [*, 'AUTHORIZATION_EXEMPTIONS'                BACnetConstructedDataAuthorizationExemptions
            [array    BACnetAuthorizationExemptionTagged('0', 'TagClass.APPLICATION_TAGS')
                                        authorizationExemption
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHORIZATION_MODE'                      BACnetConstructedDataAuthorizationMode
            [simple   BACnetAuthorizationModeTagged('0', 'TagClass.APPLICATION_TAGS') authorizationMode                 ]
        ]
        //[*, 'AUTO_SLAVE_DISCOVERY'                    BACnetConstructedDataAutoSlaveDiscovery [validation    '1 == 2'    "TODO: implement me AUTO_SLAVE_DISCOVERY BACnetConstructedDataAutoSlaveDiscovery"]]
        //[*, 'AVERAGE_VALUE'                           BACnetConstructedDataAverageValue [validation    '1 == 2'    "TODO: implement me AVERAGE_VALUE BACnetConstructedDataAverageValue"]]
        //[*, 'BACKUP_AND_RESTORE_STATE'                BACnetConstructedDataBackupAndRestoreState [validation    '1 == 2'    "TODO: implement me BACKUP_AND_RESTORE_STATE BACnetConstructedDataBackupAndRestoreState"]]
        //[*, 'BACKUP_FAILURE_TIMEOUT'                  BACnetConstructedDataBackupFailureTimeout [validation    '1 == 2'    "TODO: implement me BACKUP_FAILURE_TIMEOUT BACnetConstructedDataBackupFailureTimeout"]]
        //[*, 'BACKUP_PREPARATION_TIME'                 BACnetConstructedDataBackupPreparationTime [validation    '1 == 2'    "TODO: implement me BACKUP_PREPARATION_TIME BACnetConstructedDataBackupPreparationTime"]]
        //[*, 'BACNET_IP_GLOBAL_ADDRESS'                BACnetConstructedDataBACnetIpGlobalAddress [validation    '1 == 2'    "TODO: implement me BACNET_IP_GLOBAL_ADDRESS BACnetConstructedDataBACnetIpGlobalAddress"]]
        //[*, 'BACNET_IP_MODE'                          BACnetConstructedDataBACnetIpMode [validation    '1 == 2'    "TODO: implement me BACNET_IP_MODE BACnetConstructedDataBACnetIpMode"]]
        //[*, 'BACNET_IP_MULTICAST_ADDRESS'             BACnetConstructedDataBACnetIpMulticastAddress [validation    '1 == 2'    "TODO: implement me BACNET_IP_MULTICAST_ADDRESS BACnetConstructedDataBACnetIpMulticastAddress"]]
        //[*, 'BACNET_IP_NAT_TRAVERSAL'                 BACnetConstructedDataBACnetIpNatTraversal [validation    '1 == 2'    "TODO: implement me BACNET_IP_NAT_TRAVERSAL BACnetConstructedDataBACnetIpNatTraversal"]]
        //[*, 'BACNET_IP_UDP_PORT'                      BACnetConstructedDataBACnetIpUdpPort [validation    '1 == 2'    "TODO: implement me BACNET_IP_UDP_PORT BACnetConstructedDataBACnetIpUdpPort"]]
        //[*, 'BACNET_IPV6_MODE'                        BACnetConstructedDataBACnetIpV6Mode [validation    '1 == 2'    "TODO: implement me BACNET_IPV6_MODE BACnetConstructedDataBACnetIpV6Mode"]]
        //[*, 'BACNET_IPV6_UDP_PORT'                    BACnetConstructedDataBACnetIpV6UdpPort [validation    '1 == 2'    "TODO: implement me BACNET_IPV6_UDP_PORT BACnetConstructedDataBACnetIpV6UdpPort"]]
        //[*, 'BACNET_IPV6_MULTICAST_ADDRESS'           BACnetConstructedDataBACnetIpV6MulticastAddress [validation    '1 == 2'    "TODO: implement me BACNET_IPV6_MULTICAST_ADDRESS BACnetConstructedDataBACnetIpV6MulticastAddress"]]
        //[*, 'BASE_DEVICE_SECURITY_POLICY'             BACnetConstructedDataBaseDeviceSecurityPolicy [validation    '1 == 2'    "TODO: implement me BASE_DEVICE_SECURITY_POLICY BACnetConstructedDataBaseDeviceSecurityPolicy"]]
        [*, 'BBMD_ACCEPT_FD_REGISTRATIONS'            BACnetConstructedDataBBMDAcceptFDRegistrations
            [simple BACnetApplicationTagBoolean                               bbmdAcceptFDRegistrations                 ]
        ]
        //[*, 'BBMD_BROADCAST_DISTRIBUTION_TABLE'       BACnetConstructedDataBBMDBroadcastDistributionTable [validation    '1 == 2'    "TODO: implement me BBMD_BROADCAST_DISTRIBUTION_TABLE BACnetConstructedDataBBMDBroadcastDistributionTable"]]
        //[*, 'BBMD_FOREIGN_DEVICE_TABLE'               BACnetConstructedDataBBMDForeignDeviceTable [validation    '1 == 2'    "TODO: implement me BBMD_FOREIGN_DEVICE_TABLE BACnetConstructedDataBBMDForeignDeviceTable"]]
        [*, 'BELONGS_TO'                              BACnetConstructedDataBelongsTo
            [simple   BACnetDeviceObjectReference       belongsTo                                                      ]
        ]
        //[*, 'BIAS'                                    BACnetConstructedDataBias [validation    '1 == 2'    "TODO: implement me BIAS BACnetConstructedDataBias"]]
        [*, 'BIT_MASK'                                BACnetConstructedDataBitMask
            [simple   BACnetApplicationTagBitString     bitString                                                       ]
        ]
        [*, 'BIT_TEXT'                                BACnetConstructedDataBitText
            [array    BACnetApplicationTagCharacterString
                    bitText
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'BLINK_WARN_ENABLE'                       BACnetConstructedDataBlinkWarnEnable [validation    '1 == 2'    "TODO: implement me BLINK_WARN_ENABLE BACnetConstructedDataBlinkWarnEnable"]]
        //[*, 'BUFFER_SIZE'                             BACnetConstructedDataBufferSize [validation    '1 == 2'    "TODO: implement me BUFFER_SIZE BACnetConstructedDataBufferSize"]]
        //[*, 'CAR_ASSIGNED_DIRECTION'                  BACnetConstructedDataCarAssignedDirection [validation    '1 == 2'    "TODO: implement me CAR_ASSIGNED_DIRECTION BACnetConstructedDataCarAssignedDirection"]]
        //[*, 'CAR_DOOR_COMMAND'                        BACnetConstructedDataCarDoorCommand [validation    '1 == 2'    "TODO: implement me CAR_DOOR_COMMAND BACnetConstructedDataCarDoorCommand"]]
        //[*, 'CAR_DOOR_STATUS'                         BACnetConstructedDataCarDoorStatus [validation    '1 == 2'    "TODO: implement me CAR_DOOR_STATUS BACnetConstructedDataCarDoorStatus"]]
        //[*, 'CAR_DOOR_TEXT'                           BACnetConstructedDataCarDoorText [validation    '1 == 2'    "TODO: implement me CAR_DOOR_TEXT BACnetConstructedDataCarDoorText"]]
        //[*, 'CAR_DOOR_ZONE'                           BACnetConstructedDataCarDoorZone [validation    '1 == 2'    "TODO: implement me CAR_DOOR_ZONE BACnetConstructedDataCarDoorZone"]]
        //[*, 'CAR_DRIVE_STATUS'                        BACnetConstructedDataCarDriveStatus [validation    '1 == 2'    "TODO: implement me CAR_DRIVE_STATUS BACnetConstructedDataCarDriveStatus"]]
        //[*, 'CAR_LOAD'                                BACnetConstructedDataCarLoad [validation    '1 == 2'    "TODO: implement me CAR_LOAD BACnetConstructedDataCarLoad"]]
        [*, 'CAR_LOAD_UNITS'                          BACnetConstructedDataCarLoadUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        //[*, 'CAR_MODE'                                BACnetConstructedDataCarMode [validation    '1 == 2'    "TODO: implement me CAR_MODE BACnetConstructedDataCarMode"]]
        //[*, 'CAR_MOVING_DIRECTION'                    BACnetConstructedDataCarMovingDirection [validation    '1 == 2'    "TODO: implement me CAR_MOVING_DIRECTION BACnetConstructedDataCarMovingDirection"]]
        //[*, 'CAR_POSITION'                            BACnetConstructedDataCarPosition [validation    '1 == 2'    "TODO: implement me CAR_POSITION BACnetConstructedDataCarPosition"]]
        //[*, 'CHANGE_OF_STATE_COUNT'                   BACnetConstructedDataChangeOfStateCount [validation    '1 == 2'    "TODO: implement me CHANGE_OF_STATE_COUNT BACnetConstructedDataChangeOfStateCount"]]
        //[*, 'CHANGE_OF_STATE_TIME'                    BACnetConstructedDataChangeOfStateTime [validation    '1 == 2'    "TODO: implement me CHANGE_OF_STATE_TIME BACnetConstructedDataChangeOfStateTime"]]
        //[*, 'CHANGES_PENDING'                         BACnetConstructedDataChangesPending [validation    '1 == 2'    "TODO: implement me CHANGES_PENDING BACnetConstructedDataChangesPending"]]
        //[*, 'CHANNEL_NUMBER'                          BACnetConstructedDataChannelNumber [validation    '1 == 2'    "TODO: implement me CHANNEL_NUMBER BACnetConstructedDataChannelNumber"]]
        //[*, 'CLIENT_COV_INCREMENT'                    BACnetConstructedDataClientCovIncrement [validation    '1 == 2'    "TODO: implement me CLIENT_COV_INCREMENT BACnetConstructedDataClientCovIncrement"]]
        //[*, 'COMMAND'                                 BACnetConstructedDataCommand [validation    '1 == 2'    "TODO: implement me COMMAND BACnetConstructedDataCommand"]]
        //[*, 'COMMAND_TIME_ARRAY'                      BACnetConstructedDataCommandTimeArray [validation    '1 == 2'    "TODO: implement me COMMAND_TIME_ARRAY BACnetConstructedDataCommandTimeArray"]]
        //[*, 'CONFIGURATION_FILES'                     BACnetConstructedDataConfigurationFiles [validation    '1 == 2'    "TODO: implement me CONFIGURATION_FILES BACnetConstructedDataConfigurationFiles"]]
        //[*, 'CONTROL_GROUPS'                          BACnetConstructedDataControlGroups [validation    '1 == 2'    "TODO: implement me CONTROL_GROUPS BACnetConstructedDataControlGroups"]]
        //[*, 'CONTROLLED_VARIABLE_REFERENCE'           BACnetConstructedDataControlledVariableReference [validation    '1 == 2'    "TODO: implement me CONTROLLED_VARIABLE_REFERENCE BACnetConstructedDataControlledVariableReference"]]
        [*, 'CONTROLLED_VARIABLE_UNITS'               BACnetConstructedDataControlledVariableUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'CONTROLLED_VARIABLE_VALUE'               BACnetConstructedDataControlledVariableValue
            [simple   BACnetApplicationTagReal                                                          controlledVariableValue         ]
        ]
        //[*, 'COUNT'                                   BACnetConstructedDataCount [validation    '1 == 2'    "TODO: implement me COUNT BACnetConstructedDataCount"]]
        //[*, 'COUNT_BEFORE_CHANGE'                     BACnetConstructedDataCountBeforeChange [validation    '1 == 2'    "TODO: implement me COUNT_BEFORE_CHANGE BACnetConstructedDataCountBeforeChange"]]
        //[*, 'COUNT_CHANGE_TIME'                       BACnetConstructedDataCountChangeTime [validation    '1 == 2'    "TODO: implement me COUNT_CHANGE_TIME BACnetConstructedDataCountChangeTime"]]
        [*, 'COV_INCREMENT'                           BACnetConstructedDataCOVIncrement
            [simple   BACnetApplicationTagReal                                          covIncrement                    ]
        ]
        //[*, 'COV_PERIOD'                              BACnetConstructedDataCOVPeriod [validation    '1 == 2'    "TODO: implement me COV_PERIOD BACnetConstructedDataCOVPeriod"]]
        //[*, 'COV_RESUBSCRIPTION_INTERVAL'             BACnetConstructedDataCOVResubscriptionInterval [validation    '1 == 2'    "TODO: implement me COV_RESUBSCRIPTION_INTERVAL BACnetConstructedDataCOVResubscriptionInterval"]]
        //[*, 'COVU_PERIOD'                             BACnetConstructedDataCOVUPeriod [validation    '1 == 2'    "TODO: implement me COVU_PERIOD BACnetConstructedDataCOVUPeriod"]]
        //[*, 'COVU_RECIPIENTS'                         BACnetConstructedDataCOVURecipients [validation    '1 == 2'    "TODO: implement me COVU_RECIPIENTS BACnetConstructedDataCOVURecipients"]]
        [*, 'CREDENTIAL_DISABLE'                      BACnetConstructedDataCredentialDisable
            [simple   BACnetAccessCredentialDisableTagged('0', 'TagClass.APPLICATION_TAGS')             credentialDisable]
        ]
        [*, 'CREDENTIAL_STATUS'                       BACnetConstructedDataCredentialStatus
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    binaryPv                    ]
        ]
        [*, 'CREDENTIALS'                             BACnetConstructedDataCredentials
            [array    BACnetDeviceObjectReference
                        credentials
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'CREDENTIALS_IN_ZONE'                     BACnetConstructedDataCredentialsInZone [validation    '1 == 2'    "TODO: implement me CREDENTIALS_IN_ZONE BACnetConstructedDataCredentialsInZone"]]
        //[*, 'CURRENT_COMMAND_PRIORITY'                BACnetConstructedDataCurrentCommandPriority [validation    '1 == 2'    "TODO: implement me CURRENT_COMMAND_PRIORITY BACnetConstructedDataCurrentCommandPriority"]]
        //[*, 'DATABASE_REVISION'                       BACnetConstructedDataDatabaseRevision [validation    '1 == 2'    "TODO: implement me DATABASE_REVISION BACnetConstructedDataDatabaseRevision"]]
        //[*, 'DATE_LIST'                               BACnetConstructedDataDateList [validation    '1 == 2'    "TODO: implement me DATE_LIST BACnetConstructedDataDateList"]]
        //[*, 'DAYLIGHT_SAVINGS_STATUS'                 BACnetConstructedDataDaylightSavingsStatus [validation    '1 == 2'    "TODO: implement me DAYLIGHT_SAVINGS_STATUS BACnetConstructedDataDaylightSavingsStatus"]]
        [*, 'DAYS_REMAINING'                          BACnetConstructedDataDaysRemaining
            [simple   BACnetApplicationTagSignedInteger                     daysRemaining                               ]
        ]
        [*, 'DEADBAND'                                BACnetConstructedDataDeadband
            [simple   BACnetApplicationTagReal                                          deadband                        ]
        ]
        //[*, 'DEFAULT_FADE_TIME'                       BACnetConstructedDataDefaultFadeTime [validation    '1 == 2'    "TODO: implement me DEFAULT_FADE_TIME BACnetConstructedDataDefaultFadeTime"]]
        //[*, 'DEFAULT_RAMP_RATE'                       BACnetConstructedDataDefaultRampRate [validation    '1 == 2'    "TODO: implement me DEFAULT_RAMP_RATE BACnetConstructedDataDefaultRampRate"]]
        //[*, 'DEFAULT_STEP_INCREMENT'                  BACnetConstructedDataDefaultStepIncrement [validation    '1 == 2'    "TODO: implement me DEFAULT_STEP_INCREMENT BACnetConstructedDataDefaultStepIncrement"]]
        [*, 'DEFAULT_SUBORDINATE_RELATIONSHIP'        BACnetConstructedDataDefaultSubordinateRelationship
            [simple   BACnetRelationshipTagged('0', 'TagClass.APPLICATION_TAGS') defaultSubordinateRelationship         ]
        ]
        [*, 'DEFAULT_TIMEOUT'                         BACnetConstructedDataDefaultTimeout
            [simple   BACnetApplicationTagUnsignedInteger                     defaultTimeout                            ]
        ]
        //[*, 'DEPLOYED_PROFILE_LOCATION'               BACnetConstructedDataDeployedProfileLocation [validation    '1 == 2'    "TODO: implement me DEPLOYED_PROFILE_LOCATION BACnetConstructedDataDeployedProfileLocation"]]
        //[*, 'DERIVATIVE_CONSTANT'                     BACnetConstructedDataDerivativeConstant [validation    '1 == 2'    "TODO: implement me DERIVATIVE_CONSTANT BACnetConstructedDataDerivativeConstant"]]
        [*, 'DERIVATIVE_CONSTANT_UNITS'               BACnetConstructedDataDerivativeConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'DESCRIPTION'                             BACnetConstructedDataDescription
            [simple   BACnetApplicationTagCharacterString                               description]
        ]
        //[*, 'DESCRIPTION_OF_HALT'                     BACnetConstructedDataDescriptionOfHalt [validation    '1 == 2'    "TODO: implement me DESCRIPTION_OF_HALT BACnetConstructedDataDescriptionOfHalt"]]
        //[*, 'DEVICE_ADDRESS_BINDING'                  BACnetConstructedDataDeviceAddressBinding [validation    '1 == 2'    "TODO: implement me DEVICE_ADDRESS_BINDING BACnetConstructedDataDeviceAddressBinding"]]
        [*, 'DEVICE_TYPE'                             BACnetConstructedDataDeviceType
            [simple   BACnetApplicationTagCharacterString                               deviceType  ]
        ]
        //[*, 'DIRECT_READING'                          BACnetConstructedDataDirectReading [validation    '1 == 2'    "TODO: implement me DIRECT_READING BACnetConstructedDataDirectReading"]]
        //[*, 'DISTRIBUTION_KEY_REVISION'               BACnetConstructedDataDistributionKeyRevision [validation    '1 == 2'    "TODO: implement me DISTRIBUTION_KEY_REVISION BACnetConstructedDataDistributionKeyRevision"]]
        //[*, 'DO_NOT_HIDE'                             BACnetConstructedDataDoNotHide [validation    '1 == 2'    "TODO: implement me DO_NOT_HIDE BACnetConstructedDataDoNotHide"]]
        [*, 'DOOR_ALARM_STATE'                        BACnetConstructedDataDoorAlarmState
            [simple BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            doorAlarmState  ]
        ]
        [*, 'DOOR_EXTENDED_PULSE_TIME'                BACnetConstructedDataDoorExtendedPulseTime
            [simple BACnetApplicationTagUnsignedInteger                     doorExtendedPulseTime                       ]
        ]
        [*, 'DOOR_MEMBERS'                            BACnetConstructedDataDoorMembers
            [array    BACnetDeviceObjectReference
                        doorMembers
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'DOOR_OPEN_TOO_LONG_TIME'                 BACnetConstructedDataDoorOpenTooLongTime
            [simple BACnetApplicationTagUnsignedInteger                     doorOpenTooLongTime                         ]
        ]
        [*, 'DOOR_PULSE_TIME'                         BACnetConstructedDataDoorPulseTime
            [simple BACnetApplicationTagUnsignedInteger                     doorPulseTime                               ]
        ]
        [*, 'DOOR_STATUS'                             BACnetConstructedDataDoorStatus
            [simple   BACnetDoorStatusTagged('0', 'TagClass.APPLICATION_TAGS')
                                            doorStatus                      ]
        ]
        [*, 'DOOR_UNLOCK_DELAY_TIME'                  BACnetConstructedDataDoorUnlockDelayTime
            [simple   BACnetApplicationTagUnsignedInteger                     doorUnlockDelayTime                         ]
        ]
        //[*, 'DUTY_WINDOW'                             BACnetConstructedDataDutyWindow [validation    '1 == 2'    "TODO: implement me DUTY_WINDOW BACnetConstructedDataDutyWindow"]]
        [*, 'EFFECTIVE_PERIOD'                        BACnetConstructedDataEffectivePeriod
            [simple   BACnetDateRange               dateRange   ]
        ]
        //[*, 'EGRESS_ACTIVE'                           BACnetConstructedDataEgressActive [validation    '1 == 2'    "TODO: implement me EGRESS_ACTIVE BACnetConstructedDataEgressActive"]]
        //[*, 'EGRESS_TIME'                             BACnetConstructedDataEgressTime [validation    '1 == 2'    "TODO: implement me EGRESS_TIME BACnetConstructedDataEgressTime"]]
        //[*, 'ELAPSED_ACTIVE_TIME'                     BACnetConstructedDataElapsedActiveTime [validation    '1 == 2'    "TODO: implement me ELAPSED_ACTIVE_TIME BACnetConstructedDataElapsedActiveTime"]]
        //[*, 'ELEVATOR_GROUP'                          BACnetConstructedDataElevatorGroup [validation    '1 == 2'    "TODO: implement me ELEVATOR_GROUP BACnetConstructedDataElevatorGroup"]]
        //[*, 'ENABLE'                                  BACnetConstructedDataEnable [validation    '1 == 2'    "TODO: implement me ENABLE BACnetConstructedDataEnable"]]
        //[*, 'ENERGY_METER'                            BACnetConstructedDataEnergyMeter [validation    '1 == 2'    "TODO: implement me ENERGY_METER BACnetConstructedDataEnergyMeter"]]
        //[*, 'ENERGY_METER_REF'                        BACnetConstructedDataEnergyMeterRef [validation    '1 == 2'    "TODO: implement me ENERGY_METER_REF BACnetConstructedDataEnergyMeterRef"]]
        //[*, 'ENTRY_POINTS'                            BACnetConstructedDataEntryPoints [validation    '1 == 2'    "TODO: implement me ENTRY_POINTS BACnetConstructedDataEntryPoints"]]
        //[*, 'ERROR_LIMIT'                             BACnetConstructedDataErrorLimit [validation    '1 == 2'    "TODO: implement me ERROR_LIMIT BACnetConstructedDataErrorLimit"]]
        //[*, 'ESCALATOR_MODE'                          BACnetConstructedDataEscalatorMode [validation    '1 == 2'    "TODO: implement me ESCALATOR_MODE BACnetConstructedDataEscalatorMode"]]
        [*, 'EVENT_ALGORITHM_INHIBIT'                 BACnetConstructedDataEventAlgorithmInhibit
            [simple   BACnetApplicationTagBoolean                                       eventAlgorithmInhibit           ]
        ]
        [*, 'EVENT_ALGORITHM_INHIBIT_REF'             BACnetConstructedDataEventAlgorithmInhibitRef
            [simple   BACnetObjectPropertyReference                                     eventAlgorithmInhibitRef        ]
        ]
        [*, 'EVENT_DETECTION_ENABLE'                  BACnetConstructedDataEventDetectionEnable
            [simple   BACnetApplicationTagBoolean                                       eventDetectionEnable            ]
        ]
        [*, 'EVENT_ENABLE'                            BACnetConstructedDataEventEnable
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS') eventEnable                     ]
        ]
        [*, 'EVENT_MESSAGE_TEXTS'                     BACnetConstructedDataEventMessageTexts
            [simple   BACnetEventMessageTexts               eventMessageTexts                                           ]
        ]
        [*, 'EVENT_MESSAGE_TEXTS_CONFIG'              BACnetConstructedDataEventMessageTextsConfig
            [simple   BACnetEventMessageTextsConfig          eventStateConfig                                           ]
        ]
        //[*, 'EVENT_PARAMETERS'                        BACnetConstructedDataEventParameters [validation    '1 == 2'    "TODO: implement me EVENT_PARAMETERS BACnetConstructedDataEventParameters"]]
        [*, 'EVENT_STATE'                             BACnetConstructedDataEventState
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')          eventState                      ]
        ]
        [*, 'EVENT_TIME_STAMPS'                         BACnetConstructedDataEventTimestamps
            [simple  BACnetEventTimestamps                        eventTimeStamps                                       ]
        ]
        //[*, 'EVENT_TYPE'                              BACnetConstructedDataEventType [validation    '1 == 2'    "TODO: implement me EVENT_TYPE BACnetConstructedDataEventType"]]
        [*, 'EXCEPTION_SCHEDULE'                      BACnetConstructedDataExceptionSchedule
            [array    BACnetSpecialEvent
                            exceptionSchedule
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        //[*, 'EXECUTION_DELAY'                         BACnetConstructedDataExecutionDelay [validation    '1 == 2'    "TODO: implement me EXECUTION_DELAY BACnetConstructedDataExecutionDelay"]]
        //[*, 'EXIT_POINTS'                             BACnetConstructedDataExitPoints [validation    '1 == 2'    "TODO: implement me EXIT_POINTS BACnetConstructedDataExitPoints"]]
        //[*, 'EXPECTED_SHED_LEVEL'                     BACnetConstructedDataExpectedShedLevel [validation    '1 == 2'    "TODO: implement me EXPECTED_SHED_LEVEL BACnetConstructedDataExpectedShedLevel"]]
        //[*, 'EXPIRATION_TIME'                         BACnetConstructedDataExpirationTime [validation    '1 == 2'    "TODO: implement me EXPIRATION_TIME BACnetConstructedDataExpirationTime"]]
        [*, 'EXTENDED_TIME_ENABLE'                    BACnetConstructedDataExtendedTimeEnable
            [simple   BACnetApplicationTagBoolean                               extendedTimeEnable                      ]
        ]
        [*, 'FAILED_ATTEMPT_EVENTS'                   BACnetConstructedDataFailedAttemptEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                        failedAttemptEvents
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAILED_ATTEMPTS'                         BACnetConstructedDataFailedAttempts
            [simple   BACnetApplicationTagUnsignedInteger                               failedAttempts                  ]
        ]
        [*, 'FAILED_ATTEMPTS_TIME'                    BACnetConstructedDataFailedAttemptsTime
            [simple   BACnetApplicationTagUnsignedInteger                               failedAttemptsTime              ]
        ]
        [*, 'FAULT_HIGH_LIMIT'                        BACnetConstructedDataFaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               faultHighLimit                  ]
        ]
        [*, 'FAULT_LOW_LIMIT'                         BACnetConstructedDataFaultLowLimit
            [simple   BACnetApplicationTagReal                                          faultLowLimit                   ]
        ]
        //[*, 'FAULT_PARAMETERS'                        BACnetConstructedDataFaultParameters [validation    '1 == 2'    "TODO: implement me FAULT_PARAMETERS BACnetConstructedDataFaultParameters"]]
        //[*, 'FAULT_SIGNALS'                           BACnetConstructedDataFaultSignals [validation    '1 == 2'    "TODO: implement me FAULT_SIGNALS BACnetConstructedDataFaultSignals"]]
        //[*, 'FAULT_TYPE'                              BACnetConstructedDataFaultType [validation    '1 == 2'    "TODO: implement me FAULT_TYPE BACnetConstructedDataFaultType"]]
        ['LIFE_SAFETY_POINT', 'FAULT_VALUES'                            BACnetConstructedDataLifeSafetyPointFaultValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'FAULT_VALUES'                            BACnetConstructedDataFaultValues [validation    '1 == 2'    "TODO: implement me FAULT_VALUES BACnetConstructedDataFaultValues"]]
        //[*, 'FD_BBMD_ADDRESS'                         BACnetConstructedDataFdBbmdAddress [validation    '1 == 2'    "TODO: implement me FD_BBMD_ADDRESS BACnetConstructedDataFdBbmdAddress"]]
        //[*, 'FD_SUBSCRIPTION_LIFETIME'                BACnetConstructedDataFdSubscriptionLifetime [validation    '1 == 2'    "TODO: implement me FD_SUBSCRIPTION_LIFETIME BACnetConstructedDataFdSubscriptionLifetime"]]
        //[*, 'FEEDBACK_VALUE'                          BACnetConstructedDataFeedbackValue [validation    '1 == 2'    "TODO: implement me FEEDBACK_VALUE BACnetConstructedDataFeedbackValue"]]
        //[*, 'FILE_ACCESS_METHOD'                      BACnetConstructedDataFileAccessMethod [validation    '1 == 2'    "TODO: implement me FILE_ACCESS_METHOD BACnetConstructedDataFileAccessMethod"]]
        //[*, 'FILE_SIZE'                               BACnetConstructedDataFileSize [validation    '1 == 2'    "TODO: implement me FILE_SIZE BACnetConstructedDataFileSize"]]
        //[*, 'FILE_TYPE'                               BACnetConstructedDataFileType [validation    '1 == 2'    "TODO: implement me FILE_TYPE BACnetConstructedDataFileType"]]
        //[*, 'FIRMWARE_REVISION'                       BACnetConstructedDataFirmwareRevision [validation    '1 == 2'    "TODO: implement me FIRMWARE_REVISION BACnetConstructedDataFirmwareRevision"]]
        //[*, 'FLOOR_TEXT'                              BACnetConstructedDataFloorText [validation    '1 == 2'    "TODO: implement me FLOOR_TEXT BACnetConstructedDataFloorText"]]
        //[*, 'FULL_DUTY_BASELINE'                      BACnetConstructedDataFullDutyBaseline [validation    '1 == 2'    "TODO: implement me FULL_DUTY_BASELINE BACnetConstructedDataFullDutyBaseline"]]
        //[*, 'GLOBAL_IDENTIFIER'                       BACnetConstructedDataGlobalIdentifier [validation    '1 == 2'    "TODO: implement me GLOBAL_IDENTIFIER BACnetConstructedDataGlobalIdentifier"]]
        //[*, 'GROUP_ID'                                BACnetConstructedDataGroupId [validation    '1 == 2'    "TODO: implement me GROUP_ID BACnetConstructedDataGroupId"]]
        //[*, 'GROUP_MEMBER_NAMES'                      BACnetConstructedDataGroupMemberNames [validation    '1 == 2'    "TODO: implement me GROUP_MEMBER_NAMES BACnetConstructedDataGroupMemberNames"]]
        //[*, 'GROUP_MEMBERS'                           BACnetConstructedDataGroupMembers [validation    '1 == 2'    "TODO: implement me GROUP_MEMBERS BACnetConstructedDataGroupMembers"]]
        //[*, 'GROUP_MODE'                              BACnetConstructedDataGroupMode [validation    '1 == 2'    "TODO: implement me GROUP_MODE BACnetConstructedDataGroupMode"]]
        ['ACCUMULATOR', 'HIGH_LIMIT'                    BACnetConstructedDataAccumulatorHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               highLimit                       ]
        ]
        [*, 'HIGH_LIMIT'                                BACnetConstructedDataHighLimit
            [simple   BACnetApplicationTagReal                                          highLimit                       ]
        ]
        //[*, 'HIGHER_DECK'                             BACnetConstructedDataHigherDeck [validation    '1 == 2'    "TODO: implement me HIGHER_DECK BACnetConstructedDataHigherDeck"]]
        //[*, 'IN_PROCESS'                              BACnetConstructedDataInProcess [validation    '1 == 2'    "TODO: implement me IN_PROCESS BACnetConstructedDataInProcess"]]
        //[*, 'IN_PROGRESS'                             BACnetConstructedDataInProgress [validation    '1 == 2'    "TODO: implement me IN_PROGRESS BACnetConstructedDataInProgress"]]
        //[*, 'INACTIVE_TEXT'                           BACnetConstructedDataInactiveText [validation    '1 == 2'    "TODO: implement me INACTIVE_TEXT BACnetConstructedDataInactiveText"]]
        [*, 'INITIAL_TIMEOUT'                         BACnetConstructedDataInitialTimeout
            [simple   BACnetApplicationTagUnsignedInteger                               initialTimeout                  ]
        ]
        //[*, 'INPUT_REFERENCE'                         BACnetConstructedDataInputReference [validation    '1 == 2'    "TODO: implement me INPUT_REFERENCE BACnetConstructedDataInputReference"]]
        //[*, 'INSTALLATION_ID'                         BACnetConstructedDataInstallationId [validation    '1 == 2'    "TODO: implement me INSTALLATION_ID BACnetConstructedDataInstallationId"]]
        //[*, 'INSTANCE_OF'                             BACnetConstructedDataInstanceOf [validation    '1 == 2'    "TODO: implement me INSTANCE_OF BACnetConstructedDataInstanceOf"]]
        //[*, 'INSTANTANEOUS_POWER'                     BACnetConstructedDataInstantaneousPower [validation    '1 == 2'    "TODO: implement me INSTANTANEOUS_POWER BACnetConstructedDataInstantaneousPower"]]
        //[*, 'INTEGRAL_CONSTANT'                       BACnetConstructedDataIntegralConstant [validation    '1 == 2'    "TODO: implement me INTEGRAL_CONSTANT BACnetConstructedDataIntegralConstant"]]
        [*, 'INTEGRAL_CONSTANT_UNITS'                 BACnetConstructedDataIntegralConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        ['ANALOG_INPUT', 'INTERFACE_VALUE'                         BACnetConstructedDataAnalogInputInterfaceValue
            [optional BACnetOptionalREAL                            interfaceValue                                      ]
        ]
        //[*, 'INTERFACE_VALUE'                         BACnetConstructedDataInterfaceValue [validation    '1 == 2'    "TODO: implement me INTERFACE_VALUE BACnetConstructedDataInterfaceValue"]]
        //[*, 'INTERVAL_OFFSET'                         BACnetConstructedDataIntervalOffset [validation    '1 == 2'    "TODO: implement me INTERVAL_OFFSET BACnetConstructedDataIntervalOffset"]]
        //[*, 'IP_ADDRESS'                              BACnetConstructedDataIpAddress [validation    '1 == 2'    "TODO: implement me IP_ADDRESS BACnetConstructedDataIpAddress"]]
        //[*, 'IP_DEFAULT_GATEWAY'                      BACnetConstructedDataIpDefaultGateway [validation    '1 == 2'    "TODO: implement me IP_DEFAULT_GATEWAY BACnetConstructedDataIpDefaultGateway"]]
        //[*, 'IP_DHCP_ENABLE'                          BACnetConstructedDataIpDhcpEnable [validation    '1 == 2'    "TODO: implement me IP_DHCP_ENABLE BACnetConstructedDataIpDhcpEnable"]]
        //[*, 'IP_DHCP_LEASE_TIME'                      BACnetConstructedDataIPDHCLeaseTime [validation    '1 == 2'    "TODO: implement me IP_DHCP_LEASE_TIME BACnetConstructedDataIpDhcpLeaseTime"]]
        [*, 'IP_DHCP_LEASE_TIME_REMAINING'            BACnetConstructedDataIPDHCPLeaseTimeRemaining
            [simple BACnetApplicationTagUnsignedInteger                               ipDhcpLeaseTimeRemaining        ]
        ]
        //[*, 'IP_DHCP_SERVER'                          BACnetConstructedDataIpDhcpServer [validation    '1 == 2'    "TODO: implement me IP_DHCP_SERVER BACnetConstructedDataIpDhcpServer"]]
        //[*, 'IP_DNS_SERVER'                           BACnetConstructedDataIpDnsServer [validation    '1 == 2'    "TODO: implement me IP_DNS_SERVER BACnetConstructedDataIpDnsServer"]]
        //[*, 'IP_SUBNET_MASK'                          BACnetConstructedDataIpSubnetMask [validation    '1 == 2'    "TODO: implement me IP_SUBNET_MASK BACnetConstructedDataIpSubnetMask"]]
        //[*, 'IPV6_ADDRESS'                            BACnetConstructedDataIPv6Address [validation    '1 == 2'    "TODO: implement me IPV6_ADDRESS BACnetConstructedDataIpv6Address"]]
        //[*, 'IPV6_AUTO_ADDRESSING_ENABLE'             BACnetConstructedDataIPv6AutoAddressingEnable [validation    '1 == 2'    "TODO: implement me IPV6_AUTO_ADDRESSING_ENABLE BACnetConstructedDataIpv6AutoAddressingEnable"]]
        //[*, 'IPV6_DEFAULT_GATEWAY'                    BACnetConstructedDataIPv6DefaultGateway [validation    '1 == 2'    "TODO: implement me IPV6_DEFAULT_GATEWAY BACnetConstructedDataIpv6DefaultGateway"]]
        //[*, 'IPV6_DHCP_LEASE_TIME'                    BACnetConstructedDataIPv6DhcpLeaseTime [validation    '1 == 2'    "TODO: implement me IPV6_DHCP_LEASE_TIME BACnetConstructedDataIpv6DhcpLeaseTime"]]
        [*, 'IPV6_DHCP_LEASE_TIME_REMAINING'          BACnetConstructedDataIPv6DHCPLeaseTimeRemaining
            [simple BACnetApplicationTagUnsignedInteger                               ipv6DhcpLeaseTimeRemaining        ]
        ]
        //[*, 'IPV6_DHCP_SERVER'                        BACnetConstructedDataIPv6DhcpServer [validation    '1 == 2'    "TODO: implement me IPV6_DHCP_SERVER BACnetConstructedDataIpv6DhcpServer"]]
        //[*, 'IPV6_DNS_SERVER'                         BACnetConstructedDataIPv6DnsServer [validation    '1 == 2'    "TODO: implement me IPV6_DNS_SERVER BACnetConstructedDataIpv6DnsServer"]]
        //[*, 'IPV6_PREFIX_LENGTH'                      BACnetConstructedDataIPv6PrefixLength [validation    '1 == 2'    "TODO: implement me IPV6_PREFIX_LENGTH BACnetConstructedDataIpv6PrefixLength"]]
        //[*, 'IPV6_ZONE_INDEX'                         BACnetConstructedDataIPv6ZoneIndex [validation    '1 == 2'    "TODO: implement me IPV6_ZONE_INDEX BACnetConstructedDataIpv6ZoneIndex"]]
        //[*, 'IS_UTC'                                  BACnetConstructedDataIsUtc [validation    '1 == 2'    "TODO: implement me IS_UTC BACnetConstructedDataIsUtc"]]
        //[*, 'KEY_SETS'                                BACnetConstructedDataKeySets [validation    '1 == 2'    "TODO: implement me KEY_SETS BACnetConstructedDataKeySets"]]
        //[*, 'LANDING_CALL_CONTROL'                    BACnetConstructedDataLandingCallControl [validation    '1 == 2'    "TODO: implement me LANDING_CALL_CONTROL BACnetConstructedDataLandingCallControl"]]
        //[*, 'LANDING_CALLS'                           BACnetConstructedDataLandingCalls [validation    '1 == 2'    "TODO: implement me LANDING_CALLS BACnetConstructedDataLandingCalls"]]
        //[*, 'LANDING_DOOR_STATUS'                     BACnetConstructedDataLandingDoorStatus [validation    '1 == 2'    "TODO: implement me LANDING_DOOR_STATUS BACnetConstructedDataLandingDoorStatus"]]
        [*, 'LAST_ACCESS_EVENT'                       BACnetConstructedDataLastAccessEvent
            [simple   BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')   lastAccessEvent                       ]
        ]
        [*, 'LAST_ACCESS_POINT'                       BACnetConstructedDataLastAccessPoint
            [simple   BACnetDeviceObjectReference       lastAccessPoint                                                 ]
        ]
        //[*, 'LAST_COMMAND_TIME'                       BACnetConstructedDataLastCommandTime [validation    '1 == 2'    "TODO: implement me LAST_COMMAND_TIME BACnetConstructedDataLastCommandTime"]]
        //[*, 'LAST_CREDENTIAL_ADDED'                   BACnetConstructedDataLastCredentialAdded [validation    '1 == 2'    "TODO: implement me LAST_CREDENTIAL_ADDED BACnetConstructedDataLastCredentialAdded"]]
        //[*, 'LAST_CREDENTIAL_ADDED_TIME'              BACnetConstructedDataLastCredentialAddedTime [validation    '1 == 2'    "TODO: implement me LAST_CREDENTIAL_ADDED_TIME BACnetConstructedDataLastCredentialAddedTime"]]
        //[*, 'LAST_CREDENTIAL_REMOVED'                 BACnetConstructedDataLastCredentialRemoved [validation    '1 == 2'    "TODO: implement me LAST_CREDENTIAL_REMOVED BACnetConstructedDataLastCredentialRemoved"]]
        //[*, 'LAST_CREDENTIAL_REMOVED_TIME'            BACnetConstructedDataLastCredentialRemovedTime [validation    '1 == 2'    "TODO: implement me LAST_CREDENTIAL_REMOVED_TIME BACnetConstructedDataLastCredentialRemovedTime"]]
        //[*, 'LAST_KEY_SERVER'                         BACnetConstructedDataLastKeyServer [validation    '1 == 2'    "TODO: implement me LAST_KEY_SERVER BACnetConstructedDataLastKeyServer"]]
        //[*, 'LAST_NOTIFY_RECORD'                      BACnetConstructedDataLastNotifyRecord [validation    '1 == 2'    "TODO: implement me LAST_NOTIFY_RECORD BACnetConstructedDataLastNotifyRecord"]]
        //[*, 'LAST_PRIORITY'                           BACnetConstructedDataLastPriority [validation    '1 == 2'    "TODO: implement me LAST_PRIORITY BACnetConstructedDataLastPriority"]]
        //[*, 'LAST_RESTART_REASON'                     BACnetConstructedDataLastRestartReason [validation    '1 == 2'    "TODO: implement me LAST_RESTART_REASON BACnetConstructedDataLastRestartReason"]]
        //[*, 'LAST_RESTORE_TIME'                       BACnetConstructedDataLastRestoreTime [validation    '1 == 2'    "TODO: implement me LAST_RESTORE_TIME BACnetConstructedDataLastRestoreTime"]]
        [*, 'LAST_STATE_CHANGE'                       BACnetConstructedDataLastStateChange
            [simple   BACnetTimerTransitionTagged('0', 'TagClass.APPLICATION_TAGS')             lastStateChange         ]
        ]
        [*, 'LAST_USE_TIME'                           BACnetConstructedDataLastUseTime
            [simple   BACnetDateTime                                        lastUseTime             ]
        ]
        [*, 'LIFE_SAFETY_ALARM_VALUES'                BACnetConstructedDataLifeSafetyAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'LIGHTING_COMMAND'                        BACnetConstructedDataLightingCommand [validation    '1 == 2'    "TODO: implement me LIGHTING_COMMAND BACnetConstructedDataLightingCommand"]]
        //[*, 'LIGHTING_COMMAND_DEFAULT_PRIORITY'       BACnetConstructedDataLightingCommandDefaultPriority [validation    '1 == 2'    "TODO: implement me LIGHTING_COMMAND_DEFAULT_PRIORITY BACnetConstructedDataLightingCommandDefaultPriority"]]
        [*, 'LIMIT_ENABLE'                            BACnetConstructedDataLimitEnable
            [simple   BACnetLimitEnableTagged('0', 'TagClass.APPLICATION_TAGS')         limitEnable                     ]
        ]
        //[*, 'LIMIT_MONITORING_INTERVAL'               BACnetConstructedDataLimitMonitoringInterval [validation    '1 == 2'    "TODO: implement me LIMIT_MONITORING_INTERVAL BACnetConstructedDataLimitMonitoringInterval"]]
        //[*, 'LINK_SPEED'                              BACnetConstructedDataLinkSpeed [validation    '1 == 2'    "TODO: implement me LINK_SPEED BACnetConstructedDataLinkSpeed"]]
        //[*, 'LINK_SPEED_AUTONEGOTIATE'                BACnetConstructedDataLinkSpeedAutonegotiate [validation    '1 == 2'    "TODO: implement me LINK_SPEED_AUTONEGOTIATE BACnetConstructedDataLinkSpeedAutonegotiate"]]
        //[*, 'LINK_SPEEDS'                             BACnetConstructedDataLinkSpeeds [validation    '1 == 2'    "TODO: implement me LINK_SPEEDS BACnetConstructedDataLinkSpeeds"]]
        //[*, 'LIST_OF_GROUP_MEMBERS'                   BACnetConstructedDataListOfGroupMembers [validation    '1 == 2'    "TODO: implement me LIST_OF_GROUP_MEMBERS BACnetConstructedDataListOfGroupMembers"]]
        [*, 'LIST_OF_OBJECT_PROPERTY_REFERENCES'        BACnetConstructedDataListOfObjectPropertyReferences
            [array    BACnetDeviceObjectPropertyReference
                            references              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'LOCAL_DATE'                              BACnetConstructedDataLocalDate [validation    '1 == 2'    "TODO: implement me LOCAL_DATE BACnetConstructedDataLocalDate"]]
        //[*, 'LOCAL_FORWARDING_ONLY'                   BACnetConstructedDataLocalForwardingOnly [validation    '1 == 2'    "TODO: implement me LOCAL_FORWARDING_ONLY BACnetConstructedDataLocalForwardingOnly"]]
        //[*, 'LOCAL_TIME'                              BACnetConstructedDataLocalTime [validation    '1 == 2'    "TODO: implement me LOCAL_TIME BACnetConstructedDataLocalTime"]]
        //[*, 'LOCATION'                                BACnetConstructedDataLocation [validation    '1 == 2'    "TODO: implement me LOCATION BACnetConstructedDataLocation"]]
        [*, 'LOCK_STATUS'                             BACnetConstructedDataLockStatus
            [simple   BACnetLockStatusTagged('0', 'TagClass.APPLICATION_TAGS')      lockStatus  ]
        ]
        [*, 'LOCKOUT'                                 BACnetConstructedDataLockout
            [simple BACnetApplicationTagBoolean                               lockout                 ]
        ]
        [*, 'LOCKOUT_RELINQUISH_TIME'                 BACnetConstructedDataLockoutRelinquishTime
            [simple BACnetApplicationTagUnsignedInteger                       lockoutRelinquishTime                     ]
        ]
        //[*, 'LOG_BUFFER'                              BACnetConstructedDataLogBuffer [validation    '1 == 2'    "TODO: implement me LOG_BUFFER BACnetConstructedDataLogBuffer"]]
        //[*, 'LOG_DEVICE_OBJECT_PROPERTY'              BACnetConstructedDataLogDeviceObjectProperty [validation    '1 == 2'    "TODO: implement me LOG_DEVICE_OBJECT_PROPERTY BACnetConstructedDataLogDeviceObjectProperty"]]
        //[*, 'LOG_INTERVAL'                            BACnetConstructedDataLogInterval [validation    '1 == 2'    "TODO: implement me LOG_INTERVAL BACnetConstructedDataLogInterval"]]
        //[*, 'LOGGING_OBJECT'                          BACnetConstructedDataLoggingObject [validation    '1 == 2'    "TODO: implement me LOGGING_OBJECT BACnetConstructedDataLoggingObject"]]
        //[*, 'LOGGING_RECORD'                          BACnetConstructedDataLoggingRecord [validation    '1 == 2'    "TODO: implement me LOGGING_RECORD BACnetConstructedDataLoggingRecord"]]
        //[*, 'LOGGING_TYPE'                            BACnetConstructedDataLoggingType [validation    '1 == 2'    "TODO: implement me LOGGING_TYPE BACnetConstructedDataLoggingType"]]
        //[*, 'LOW_DIFF_LIMIT'                          BACnetConstructedDataLowDiffLimit [validation    '1 == 2'    "TODO: implement me LOW_DIFF_LIMIT BACnetConstructedDataLowDiffLimit"]]
        ['ACCUMULATOR', 'LOW_LIMIT'                     BACnetConstructedDataAccumulatorLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                               lowLimit                        ]
        ]
        [*, 'LOW_LIMIT'                                 BACnetConstructedDataLowLimit
            [simple   BACnetApplicationTagReal                                          lowLimit                        ]
        ]
        //[*, 'LOWER_DECK'                              BACnetConstructedDataLowerDeck [validation    '1 == 2'    "TODO: implement me LOWER_DECK BACnetConstructedDataLowerDeck"]]
        //[*, 'MAC_ADDRESS'                             BACnetConstructedDataMacAddress [validation    '1 == 2'    "TODO: implement me MAC_ADDRESS BACnetConstructedDataMacAddress"]]
        //[*, 'MACHINE_ROOM_ID'                         BACnetConstructedDataMachineRoomId [validation    '1 == 2'    "TODO: implement me MACHINE_ROOM_ID BACnetConstructedDataMachineRoomId"]]
        //[*, 'MAINTENANCE_REQUIRED'                    BACnetConstructedDataMaintenanceRequired [validation    '1 == 2'    "TODO: implement me MAINTENANCE_REQUIRED BACnetConstructedDataMaintenanceRequired"]]
        //[*, 'MAKING_CAR_CALL'                         BACnetConstructedDataMakingCarCall [validation    '1 == 2'    "TODO: implement me MAKING_CAR_CALL BACnetConstructedDataMakingCarCall"]]
        //[*, 'MANIPULATED_VARIABLE_REFERENCE'          BACnetConstructedDataManipulatedVariableReference [validation    '1 == 2'    "TODO: implement me MANIPULATED_VARIABLE_REFERENCE BACnetConstructedDataManipulatedVariableReference"]]
        [*, 'MANUAL_SLAVE_ADDRESS_BINDING'            BACnetConstructedDataManualSlaveAddressBinding
            [array    BACnetAddressBinding
                                manualSlaveAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MASKED_ALARM_VALUES'                     BACnetConstructedDataMaskedAlarmValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                                maskedAlarmValues
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'MAX_ACTUAL_VALUE'                        BACnetConstructedDataMaxActualValue [validation    '1 == 2'    "TODO: implement me MAX_ACTUAL_VALUE BACnetConstructedDataMaxActualValue"]]
        //[*, 'MAX_APDU_LENGTH_ACCEPTED'                BACnetConstructedDataMaxApduLengthAccepted [validation    '1 == 2'    "TODO: implement me MAX_APDU_LENGTH_ACCEPTED BACnetConstructedDataMaxApduLengthAccepted"]]
        [*, 'MAX_FAILED_ATTEMPTS'                     BACnetConstructedDataMaxFailedAttempts
            [simple BACnetApplicationTagUnsignedInteger                     maxFailedAttempts                           ]
        ]
        //[*, 'MAX_INFO_FRAMES'                         BACnetConstructedDataMaxInfoFrames [validation    '1 == 2'    "TODO: implement me MAX_INFO_FRAMES BACnetConstructedDataMaxInfoFrames"]]
        //[*, 'MAX_MASTER'                              BACnetConstructedDataMaxMaster [validation    '1 == 2'    "TODO: implement me MAX_MASTER BACnetConstructedDataMaxMaster"]]
        ['ACCUMULATOR', 'MAX_PRES_VALUE'               BACnetConstructedDataAccumulatorMaxPresValue
            [simple BACnetApplicationTagUnsignedInteger                     maxPresValue                                 ]
        ]
        [*, 'MAX_PRES_VALUE'                            BACnetConstructedDataMaxPresValue
            [simple BACnetApplicationTagReal                                maxPresValue                                ]
        ]
        //[*, 'MAX_SEGMENTS_ACCEPTED'                   BACnetConstructedDataMaxSegmentsAccepted [validation    '1 == 2'    "TODO: implement me MAX_SEGMENTS_ACCEPTED BACnetConstructedDataMaxSegmentsAccepted"]]
        //[*, 'MAXIMUM_OUTPUT'                          BACnetConstructedDataMaximumOutput [validation    '1 == 2'    "TODO: implement me MAXIMUM_OUTPUT BACnetConstructedDataMaximumOutput"]]
        //[*, 'MAXIMUM_VALUE'                           BACnetConstructedDataMaximumValue [validation    '1 == 2'    "TODO: implement me MAXIMUM_VALUE BACnetConstructedDataMaximumValue"]]
        //[*, 'MAXIMUM_VALUE_TIMESTAMP'                 BACnetConstructedDataMaximumValueTimestamp [validation    '1 == 2'    "TODO: implement me MAXIMUM_VALUE_TIMESTAMP BACnetConstructedDataMaximumValueTimestamp"]]
        [*, 'MEMBER_OF' BACnetConstructedDataMemberOf
            [array    BACnetDeviceObjectReference
                    zones
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'MEMBER_STATUS_FLAGS'                     BACnetConstructedDataMemberStatusFlags [validation    '1 == 2'    "TODO: implement me MEMBER_STATUS_FLAGS BACnetConstructedDataMemberStatusFlags"]]
        [*, 'MEMBERS'                                 BACnetConstructedDataMembers
            [array    BACnetDeviceObjectReference
                                    members
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'MIN_ACTUAL_VALUE'                        BACnetConstructedDataMinActualValue [validation    '1 == 2'    "TODO: implement me MIN_ACTUAL_VALUE BACnetConstructedDataMinActualValue"]]
        ['ACCUMULATOR', 'MIN_PRES_VALUE'               BACnetConstructedDataAccumulatorMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                     minPresValue                              ]
        ]
        [*, 'MIN_PRES_VALUE'                            BACnetConstructedDataMinPresValue
            [simple   BACnetApplicationTagReal                                minPresValue                              ]
        ]
        //[*, 'MINIMUM_OFF_TIME'                        BACnetConstructedDataMinimumOffTime [validation    '1 == 2'    "TODO: implement me MINIMUM_OFF_TIME BACnetConstructedDataMinimumOffTime"]]
        //[*, 'MINIMUM_ON_TIME'                         BACnetConstructedDataMinimumOnTime [validation    '1 == 2'    "TODO: implement me MINIMUM_ON_TIME BACnetConstructedDataMinimumOnTime"]]
        //[*, 'MINIMUM_OUTPUT'                          BACnetConstructedDataMinimumOutput [validation    '1 == 2'    "TODO: implement me MINIMUM_OUTPUT BACnetConstructedDataMinimumOutput"]]
        [*, 'MINIMUM_VALUE'                           BACnetConstructedDataMinimumValue
            [simple   BACnetApplicationTagReal                              minimumValue                                ]
        ]
        //[*, 'MINIMUM_VALUE_TIMESTAMP'                 BACnetConstructedDataMinimumValueTimestamp [validation    '1 == 2'    "TODO: implement me MINIMUM_VALUE_TIMESTAMP BACnetConstructedDataMinimumValueTimestamp"]]
        //[*, 'MODE'                                    BACnetConstructedDataMode [validation    '1 == 2'    "TODO: implement me MODE BACnetConstructedDataMode"]]
        //[*, 'MODEL_NAME'                              BACnetConstructedDataModelName [validation    '1 == 2'    "TODO: implement me MODEL_NAME BACnetConstructedDataModelName"]]
        //[*, 'MODIFICATION_DATE'                       BACnetConstructedDataModificationDate [validation    '1 == 2'    "TODO: implement me MODIFICATION_DATE BACnetConstructedDataModificationDate"]]
        [*, 'MUSTER_POINT'                            BACnetConstructedDataMusterPoint
            [simple   BACnetApplicationTagBoolean                                       musterPoint                     ]
        ]
        [*, 'NEGATIVE_ACCESS_RULES'                   BACnetConstructedDataNegativeAccessRules
            [array    BACnetAccessRule
                            negativeAccessRules
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        //[*, 'NETWORK_ACCESS_SECURITY_POLICIES'        BACnetConstructedDataNetworkAccessSecurityPolicies [validation    '1 == 2'    "TODO: implement me NETWORK_ACCESS_SECURITY_POLICIES BACnetConstructedDataNetworkAccessSecurityPolicies"]]
        //[*, 'NETWORK_INTERFACE_NAME'                  BACnetConstructedDataNetworkInterfaceName [validation    '1 == 2'    "TODO: implement me NETWORK_INTERFACE_NAME BACnetConstructedDataNetworkInterfaceName"]]
        //[*, 'NETWORK_NUMBER'                          BACnetConstructedDataNetworkNumber [validation    '1 == 2'    "TODO: implement me NETWORK_NUMBER BACnetConstructedDataNetworkNumber"]]
        //[*, 'NETWORK_NUMBER_QUALITY'                  BACnetConstructedDataNetworkNumberQuality [validation    '1 == 2'    "TODO: implement me NETWORK_NUMBER_QUALITY BACnetConstructedDataNetworkNumberQuality"]]
        //[*, 'NETWORK_TYPE'                            BACnetConstructedDataNetworkType [validation    '1 == 2'    "TODO: implement me NETWORK_TYPE BACnetConstructedDataNetworkType"]]
        //[*, 'NEXT_STOPPING_FLOOR'                     BACnetConstructedDataNextStoppingFloor [validation    '1 == 2'    "TODO: implement me NEXT_STOPPING_FLOOR BACnetConstructedDataNextStoppingFloor"]]
        [*, 'NODE_SUBTYPE'                            BACnetConstructedDataNodeSubtype
            [simple BACnetApplicationTagCharacterString                  nodeSubType                                    ]
        ]
        [*, 'NODE_TYPE'                               BACnetConstructedDataNodeType
            [simple BACnetNodeTypeTagged('0', 'TagClass.APPLICATION_TAGS')  nodeType                                    ]
        ]
        [*, 'NOTIFICATION_CLASS'                      BACnetConstructedDataNotificationClass
            [simple BACnetApplicationTagUnsignedInteger                               notificationClass                 ]
        ]
        //[*, 'NOTIFICATION_THRESHOLD'                  BACnetConstructedDataNotificationThreshold [validation    '1 == 2'    "TODO: implement me NOTIFICATION_THRESHOLD BACnetConstructedDataNotificationThreshold"]]
        [*, 'NOTIFY_TYPE'                             BACnetConstructedDataNotifyType
            [simple BACnetNotifyTypeTagged('0', 'TagClass.APPLICATION_TAGS')          notifyType                        ]
        ]
        //[*, 'NUMBER_OF_APDU_RETRIES'                  BACnetConstructedDataNumberOfApduRetries [validation    '1 == 2'    "TODO: implement me NUMBER_OF_APDU_RETRIES BACnetConstructedDataNumberOfApduRetries"]]
        [*, 'NUMBER_OF_AUTHENTICATION_POLICIES'       BACnetConstructedDataNumberOfAuthenticationPolicies
            [simple   BACnetApplicationTagUnsignedInteger                               numberOfAuthenticationPolicies  ]
        ]
        //[*, 'NUMBER_OF_STATES'                        BACnetConstructedDataNumberOfStates [validation    '1 == 2'    "TODO: implement me NUMBER_OF_STATES BACnetConstructedDataNumberOfStates"]]
        [*, 'OBJECT_IDENTIFIER'                       BACnetConstructedDataObjectIdentifier
            [simple   BACnetApplicationTagObjectIdentifier                              objectIdentifier                ]
        ]
        //[*, 'OBJECT_LIST'                             BACnetConstructedDataObjectList [validation    '1 == 2'    "TODO: implement me OBJECT_LIST BACnetConstructedDataObjectList"]]
        [*, 'OBJECT_NAME'                             BACnetConstructedDataObjectName
            [simple   BACnetApplicationTagCharacterString                               objectName                      ]
        ]
        //[*, 'OBJECT_PROPERTY_REFERENCE'               BACnetConstructedDataObjectPropertyReference [validation    '1 == 2'    "TODO: implement me OBJECT_PROPERTY_REFERENCE BACnetConstructedDataObjectPropertyReference"]]
        [*, 'OBJECT_TYPE'                             BACnetConstructedDataObjectType
            [simple   BACnetObjectTypeTagged('0', 'TagClass.APPLICATION_TAGS')          objectType                      ]
        ]
        //[*, 'OCCUPANCY_COUNT'                         BACnetConstructedDataOccupancyCount [validation    '1 == 2'    "TODO: implement me OCCUPANCY_COUNT BACnetConstructedDataOccupancyCount"]]
        [*, 'OCCUPANCY_COUNT_ADJUST'                  BACnetConstructedDataOccupancyCountAdjust
            [simple BACnetApplicationTagBoolean                               occupancyCountAdjust                      ]
        ]
        //[*, 'OCCUPANCY_COUNT_ENABLE'                  BACnetConstructedDataOccupancyCountEnable [validation    '1 == 2'    "TODO: implement me OCCUPANCY_COUNT_ENABLE BACnetConstructedDataOccupancyCountEnable"]]
        //[*, 'OCCUPANCY_LOWER_LIMIT'                   BACnetConstructedDataOccupancyLowerLimit [validation    '1 == 2'    "TODO: implement me OCCUPANCY_LOWER_LIMIT BACnetConstructedDataOccupancyLowerLimit"]]
        [*, 'OCCUPANCY_LOWER_LIMIT_ENFORCED'          BACnetConstructedDataOccupancyLowerLimitEnforced
            [simple BACnetApplicationTagBoolean                               occupancyLowerLimitEnforced               ]
        ]
        //[*, 'OCCUPANCY_STATE'                         BACnetConstructedDataOccupancyState [validation    '1 == 2'    "TODO: implement me OCCUPANCY_STATE BACnetConstructedDataOccupancyState"]]
        //[*, 'OCCUPANCY_UPPER_LIMIT'                   BACnetConstructedDataOccupancyUpperLimit [validation    '1 == 2'    "TODO: implement me OCCUPANCY_UPPER_LIMIT BACnetConstructedDataOccupancyUpperLimit"]]
        [*, 'OCCUPANCY_UPPER_LIMIT_ENFORCED'          BACnetConstructedDataOccupancyUpperLimitEnforced
            [simple BACnetApplicationTagBoolean                               occupancyUpperLimitEnforced               ]
        ]
        //[*, 'OPERATION_DIRECTION'                     BACnetConstructedDataOperationDirection [validation    '1 == 2'    "TODO: implement me OPERATION_DIRECTION BACnetConstructedDataOperationDirection"]]
        //[*, 'OPERATION_EXPECTED'                      BACnetConstructedDataOperationExpected [validation    '1 == 2'    "TODO: implement me OPERATION_EXPECTED BACnetConstructedDataOperationExpected"]]
        //[*, 'OPTIONAL'                                BACnetConstructedDataOptional [validation    '1 == 2'    "TODO: implement me OPTIONAL BACnetConstructedDataOptional"]]
        [*, 'OUT_OF_SERVICE'                          BACnetConstructedDataOutOfService
            [simple   BACnetApplicationTagBoolean                                                       outOfService    ]
        ]
        [*, 'OUTPUT_UNITS'                            BACnetConstructedDataOutputUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        //[*, 'PACKET_REORDER_TIME'                     BACnetConstructedDataPacketReorderTime [validation    '1 == 2'    "TODO: implement me PACKET_REORDER_TIME BACnetConstructedDataPacketReorderTime"]]
        //[*, 'PASSBACK_MODE'                           BACnetConstructedDataPassbackMode [validation    '1 == 2'    "TODO: implement me PASSBACK_MODE BACnetConstructedDataPassbackMode"]]
        //[*, 'PASSBACK_TIMEOUT'                        BACnetConstructedDataPassbackTimeout [validation    '1 == 2'    "TODO: implement me PASSBACK_TIMEOUT BACnetConstructedDataPassbackTimeout"]]
        //[*, 'PASSENGER_ALARM'                         BACnetConstructedDataPassengerAlarm [validation    '1 == 2'    "TODO: implement me PASSENGER_ALARM BACnetConstructedDataPassengerAlarm"]]
        //[*, 'POLARITY'                                BACnetConstructedDataPolarity [validation    '1 == 2'    "TODO: implement me POLARITY BACnetConstructedDataPolarity"]]
        //[*, 'PORT_FILTER'                             BACnetConstructedDataPortFilter [validation    '1 == 2'    "TODO: implement me PORT_FILTER BACnetConstructedDataPortFilter"]]
        [*, 'POSITIVE_ACCESS_RULES'                     BACnetConstructedDataPositiveAccessRules
            [array    BACnetAccessRule
                            positiveAccessRules
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        //[*, 'POWER'                                   BACnetConstructedDataPower [validation    '1 == 2'    "TODO: implement me POWER BACnetConstructedDataPower"]]
        //[*, 'POWER_MODE'                              BACnetConstructedDataPowerMode [validation    '1 == 2'    "TODO: implement me POWER_MODE BACnetConstructedDataPowerMode"]]
        //[*, 'PRESCALE'                                BACnetConstructedDataPrescale [validation    '1 == 2'    "TODO: implement me PRESCALE BACnetConstructedDataPrescale"]]
        ['ANALOG_INPUT', 'PRESENT_VALUE'                           BACnetConstructedDataAnalogInputPresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        //[*, 'PRESENT_VALUE'                           BACnetConstructedDataPresentValue [validation    '1 == 2'    "TODO: implement me PRESENT_VALUE BACnetConstructedDataPresentValue"]]
        //[*, 'PRIORITY'                                BACnetConstructedDataPriority [validation    '1 == 2'    "TODO: implement me PRIORITY BACnetConstructedDataPriority"]]
        [*, 'PRIORITY_ARRAY'                          BACnetConstructedDataPriorityArray
            [simple   BACnetPriorityArray('objectTypeArgument', 'tagNumber')                            priorityArray   ]
        ]
        //[*, 'PRIORITY_FOR_WRITING'                    BACnetConstructedDataPriorityForWriting [validation    '1 == 2'    "TODO: implement me PRIORITY_FOR_WRITING BACnetConstructedDataPriorityForWriting"]]
        //[*, 'PROCESS_IDENTIFIER'                      BACnetConstructedDataProcessIdentifier [validation    '1 == 2'    "TODO: implement me PROCESS_IDENTIFIER BACnetConstructedDataProcessIdentifier"]]
        //[*, 'PROCESS_IDENTIFIER_FILTER'               BACnetConstructedDataProcessIdentifierFilter [validation    '1 == 2'    "TODO: implement me PROCESS_IDENTIFIER_FILTER BACnetConstructedDataProcessIdentifierFilter"]]
        [*, 'PROFILE_LOCATION'                        BACnetConstructedDataProfileLocation
            [simple   BACnetApplicationTagCharacterString                               profileLocation                 ]
        ]
        [*, 'PROFILE_NAME'                            BACnetConstructedDataProfileName
            [simple   BACnetApplicationTagCharacterString                               profileName                     ]
        ]
        //[*, 'PROGRAM_CHANGE'                          BACnetConstructedDataProgramChange [validation    '1 == 2'    "TODO: implement me PROGRAM_CHANGE BACnetConstructedDataProgramChange"]]
        //[*, 'PROGRAM_LOCATION'                        BACnetConstructedDataProgramLocation [validation    '1 == 2'    "TODO: implement me PROGRAM_LOCATION BACnetConstructedDataProgramLocation"]]
        //[*, 'PROGRAM_STATE'                           BACnetConstructedDataProgramState [validation    '1 == 2'    "TODO: implement me PROGRAM_STATE BACnetConstructedDataProgramState"]]
        [*, 'PROPERTY_LIST'                           BACnetConstructedDataPropertyList
            [array    BACnetPropertyIdentifierTagged('0', 'TagClass.APPLICATION_TAGS')
                                propertyList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        //[*, 'PROPORTIONAL_CONSTANT'                   BACnetConstructedDataProportionalConstant [validation    '1 == 2'    "TODO: implement me PROPORTIONAL_CONSTANT BACnetConstructedDataProportionalConstant"]]
        [*, 'PROPORTIONAL_CONSTANT_UNITS'             BACnetConstructedDataProportionalConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        //[*, 'PROTOCOL_LEVEL'                          BACnetConstructedDataProtocolLevel [validation    '1 == 2'    "TODO: implement me PROTOCOL_LEVEL BACnetConstructedDataProtocolLevel"]]
        //[*, 'PROTOCOL_CONFORMANCE_CLASS'              BACnetConstructedDataProtocolConformanceClass [validation    '1 == 2'    "TODO: implement me PROTOCOL_CONFORMANCE_CLASS BACnetConstructedDataProtocolConformanceClass"]]
        //[*, 'PROTOCOL_OBJECT_TYPES_SUPPORTED'         BACnetConstructedDataProtocolObjectTypesSupported [validation    '1 == 2'    "TODO: implement me PROTOCOL_OBJECT_TYPES_SUPPORTED BACnetConstructedDataProtocolObjectTypesSupported"]]
        //[*, 'PROTOCOL_REVISION'                       BACnetConstructedDataProtocolRevision [validation    '1 == 2'    "TODO: implement me PROTOCOL_REVISION BACnetConstructedDataProtocolRevision"]]
        //[*, 'PROTOCOL_SERVICES_SUPPORTED'             BACnetConstructedDataProtocolServicesSupported [validation    '1 == 2'    "TODO: implement me PROTOCOL_SERVICES_SUPPORTED BACnetConstructedDataProtocolServicesSupported"]]
        //[*, 'PROTOCOL_VERSION'                        BACnetConstructedDataProtocolVersion [validation    '1 == 2'    "TODO: implement me PROTOCOL_VERSION BACnetConstructedDataProtocolVersion"]]
        //[*, 'PULSE_RATE'                              BACnetConstructedDataPulseRate [validation    '1 == 2'    "TODO: implement me PULSE_RATE BACnetConstructedDataPulseRate"]]
        //[*, 'READ_ONLY'                               BACnetConstructedDataReadOnly [validation    '1 == 2'    "TODO: implement me READ_ONLY BACnetConstructedDataReadOnly"]]
        [*, 'REASON_FOR_DISABLE'                      BACnetConstructedDataReasonForDisable
            [array    BACnetAccessCredentialDisableReasonTagged('0', 'TagClass.APPLICATION_TAGS')
                                            reasonForDisable
                                                    terminated
                                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'REASON_FOR_HALT'                         BACnetConstructedDataReasonForHalt [validation    '1 == 2'    "TODO: implement me REASON_FOR_HALT BACnetConstructedDataReasonForHalt"]]
        //[*, 'RECIPIENT_LIST'                          BACnetConstructedDataRecipientList [validation    '1 == 2'    "TODO: implement me RECIPIENT_LIST BACnetConstructedDataRecipientList"]]
        //[*, 'RECORD_COUNT'                            BACnetConstructedDataRecordCount [validation    '1 == 2'    "TODO: implement me RECORD_COUNT BACnetConstructedDataRecordCount"]]
        //[*, 'RECORDS_SINCE_NOTIFICATION'              BACnetConstructedDataRecordsSinceNotification [validation    '1 == 2'    "TODO: implement me RECORDS_SINCE_NOTIFICATION BACnetConstructedDataRecordsSinceNotification"]]
        //[*, 'REFERENCE_PORT'                          BACnetConstructedDataReferencePort [validation    '1 == 2'    "TODO: implement me REFERENCE_PORT BACnetConstructedDataReferencePort"]]
        //[*, 'REGISTERED_CAR_CALL'                     BACnetConstructedDataRegisteredCarCall [validation    '1 == 2'    "TODO: implement me REGISTERED_CAR_CALL BACnetConstructedDataRegisteredCarCall"]]
        [*, 'RELIABILITY'                             BACnetConstructedDataReliability
            [simple   BACnetReliabilityTagged('0', 'TagClass.APPLICATION_TAGS')                         reliability     ]
        ]
        [*, 'RELIABILITY_EVALUATION_INHIBIT'          BACnetConstructedDataReliabilityEvaluationInhibit
            [simple   BACnetApplicationTagBoolean                                       reliabilityEvaluationInhibit    ]
        ]
        //[*, 'RELINQUISH_DEFAULT'                      BACnetConstructedDataRelinquishDefault [validation    '1 == 2'    "TODO: implement me RELINQUISH_DEFAULT BACnetConstructedDataRelinquishDefault"]]
        [*, 'REPRESENTS'                              BACnetConstructedDataRepresents
            [simple   BACnetDeviceObjectReference       represents                                                      ]
        ]
        //[*, 'REQUESTED_SHED_LEVEL'                    BACnetConstructedDataRequestedShedLevel [validation    '1 == 2'    "TODO: implement me REQUESTED_SHED_LEVEL BACnetConstructedDataRequestedShedLevel"]]
        //[*, 'REQUESTED_UPDATE_INTERVAL'               BACnetConstructedDataRequestedUpdateInterval [validation    '1 == 2'    "TODO: implement me REQUESTED_UPDATE_INTERVAL BACnetConstructedDataRequestedUpdateInterval"]]
        //[*, 'REQUIRED'                                BACnetConstructedDataRequired [validation    '1 == 2'    "TODO: implement me REQUIRED BACnetConstructedDataRequired"]]
        [*, 'RESOLUTION'                              BACnetConstructedDataResolution
            [simple   BACnetApplicationTagReal                                          resolution                      ]
        ]
        //[*, 'RESTART_NOTIFICATION_RECIPIENTS'         BACnetConstructedDataRestartNotificationRecipients [validation    '1 == 2'    "TODO: implement me RESTART_NOTIFICATION_RECIPIENTS BACnetConstructedDataRestartNotificationRecipients"]]
        //[*, 'RESTORE_COMPLETION_TIME'                 BACnetConstructedDataRestoreCompletionTime [validation    '1 == 2'    "TODO: implement me RESTORE_COMPLETION_TIME BACnetConstructedDataRestoreCompletionTime"]]
        //[*, 'RESTORE_PREPARATION_TIME'                BACnetConstructedDataRestorePreparationTime [validation    '1 == 2'    "TODO: implement me RESTORE_PREPARATION_TIME BACnetConstructedDataRestorePreparationTime"]]
        //[*, 'ROUTING_TABLE'                           BACnetConstructedDataRoutingTable [validation    '1 == 2'    "TODO: implement me ROUTING_TABLE BACnetConstructedDataRoutingTable"]]
        //[*, 'SCALE'                                   BACnetConstructedDataScale [validation    '1 == 2'    "TODO: implement me SCALE BACnetConstructedDataScale"]]
        //[*, 'SCALE_FACTOR'                            BACnetConstructedDataScaleFactor [validation    '1 == 2'    "TODO: implement me SCALE_FACTOR BACnetConstructedDataScaleFactor"]]
        [*, 'SCHEDULE_DEFAULT'                        BACnetConstructedDataScheduleDefault
            [simple   BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifierArgument') scheduleDefault  ]
        ]
        [*, 'SECURED_STATUS'                          BACnetConstructedDataSecuredStatus
            [simple   BACnetDoorSecuredStatusTagged('0', 'TagClass.APPLICATION_TAGS')         securedStatus             ]
        ]
        //[*, 'SECURITY_PDU_TIMEOUT'                    BACnetConstructedDataSecurityPduTimeout [validation    '1 == 2'    "TODO: implement me SECURITY_PDU_TIMEOUT BACnetConstructedDataSecurityPduTimeout"]]
        //[*, 'SECURITY_TIME_WINDOW'                    BACnetConstructedDataSecurityTimeWindow [validation    '1 == 2'    "TODO: implement me SECURITY_TIME_WINDOW BACnetConstructedDataSecurityTimeWindow"]]
        //[*, 'SEGMENTATION_SUPPORTED'                  BACnetConstructedDataSegmentationSupported [validation    '1 == 2'    "TODO: implement me SEGMENTATION_SUPPORTED BACnetConstructedDataSegmentationSupported"]]
        //[*, 'SERIAL_NUMBER'                           BACnetConstructedDataSerialNumber [validation    '1 == 2'    "TODO: implement me SERIAL_NUMBER BACnetConstructedDataSerialNumber"]]
        //[*, 'SETPOINT'                                BACnetConstructedDataSetpoint [validation    '1 == 2'    "TODO: implement me SETPOINT BACnetConstructedDataSetpoint"]]
        //[*, 'SETPOINT_REFERENCE'                      BACnetConstructedDataSetpointReference [validation    '1 == 2'    "TODO: implement me SETPOINT_REFERENCE BACnetConstructedDataSetpointReference"]]
        //[*, 'SETTING'                                 BACnetConstructedDataSetting [validation    '1 == 2'    "TODO: implement me SETTING BACnetConstructedDataSetting"]]
        //[*, 'SHED_DURATION'                           BACnetConstructedDataShedDuration [validation    '1 == 2'    "TODO: implement me SHED_DURATION BACnetConstructedDataShedDuration"]]
        //[*, 'SHED_LEVEL_DESCRIPTIONS'                 BACnetConstructedDataShedLevelDescriptions [validation    '1 == 2'    "TODO: implement me SHED_LEVEL_DESCRIPTIONS BACnetConstructedDataShedLevelDescriptions"]]
        //[*, 'SHED_LEVELS'                             BACnetConstructedDataShedLevels [validation    '1 == 2'    "TODO: implement me SHED_LEVELS BACnetConstructedDataShedLevels"]]
        //[*, 'SILENCED'                                BACnetConstructedDataSilenced [validation    '1 == 2'    "TODO: implement me SILENCED BACnetConstructedDataSilenced"]]
        //[*, 'SLAVE_ADDRESS_BINDING'                   BACnetConstructedDataSlaveAddressBinding [validation    '1 == 2'    "TODO: implement me SLAVE_ADDRESS_BINDING BACnetConstructedDataSlaveAddressBinding"]]
        //[*, 'SLAVE_PROXY_ENABLE'                      BACnetConstructedDataSlaveProxyEnable [validation    '1 == 2'    "TODO: implement me SLAVE_PROXY_ENABLE BACnetConstructedDataSlaveProxyEnable"]]
        //[*, 'START_TIME'                              BACnetConstructedDataStartTime [validation    '1 == 2'    "TODO: implement me START_TIME BACnetConstructedDataStartTime"]]
        [*, 'STATE_CHANGE_VALUES'                     BACnetConstructedDataStateChangeValues
            [array    BACnetTimerStateChangeValue('objectTypeArgument')
                                stateChangeValues
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
            [validation 'COUNT(stateChangeValues) == 7' "stateChangeValues should have exactly 7 values"                ]
        ]
        //[*, 'STATE_DESCRIPTION'                       BACnetConstructedDataStateDescription [validation    '1 == 2'    "TODO: implement me STATE_DESCRIPTION BACnetConstructedDataStateDescription"]]
        //[*, 'STATE_TEXT'                              BACnetConstructedDataStateText [validation    '1 == 2'    "TODO: implement me STATE_TEXT BACnetConstructedDataStateText"]]
        [*, 'STATUS_FLAGS'                            BACnetConstructedDataStatusFlags
            [simple   BACnetStatusFlagsTagged('0', 'TagClass.APPLICATION_TAGS')         statusFlags                     ]
        ]
        //[*, 'STOP_TIME'                               BACnetConstructedDataStopTime [validation    '1 == 2'    "TODO: implement me STOP_TIME BACnetConstructedDataStopTime"]]
        //[*, 'STOP_WHEN_FULL'                          BACnetConstructedDataStopWhenFull [validation    '1 == 2'    "TODO: implement me STOP_WHEN_FULL BACnetConstructedDataStopWhenFull"]]
        [*, 'STRIKE_COUNT'                            BACnetConstructedDataStrikeCount
            [simple BACnetApplicationTagUnsignedInteger                       strikeCount                     ]
        ]
        //[*, 'STRUCTURED_OBJECT_LIST'                  BACnetConstructedDataStructuredObjectList [validation    '1 == 2'    "TODO: implement me STRUCTURED_OBJECT_LIST BACnetConstructedDataStructuredObjectList"]]
        [*, 'SUBORDINATE_ANNOTATIONS'                 BACnetConstructedDataSubordinateAnnotations
            [array    BACnetApplicationTagCharacterString
                    subordinateAnnotations
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_LIST'                        BACnetConstructedDataSubordinateList
            [array    BACnetDeviceObjectReference
                        subordinateList
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_NODE_TYPES'                  BACnetConstructedDataSubordinateNodeTypes
            [array    BACnetNodeTypeTagged('0', 'TagClass.APPLICATION_TAGS')
                        subordinateNodeTypes
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_RELATIONSHIPS'               BACnetConstructedDataSubordinateRelationships
            [array    BACnetRelationshipTagged('0', 'TagClass.APPLICATION_TAGS')
                                    subordinateRelationships
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_TAGS'                        BACnetConstructedDataSubordinateTags
            [array    BACnetNameValueCollection
                        subordinateList
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        //[*, 'SUBSCRIBED_RECIPIENTS'                   BACnetConstructedDataSubscribedRecipients [validation    '1 == 2'    "TODO: implement me SUBSCRIBED_RECIPIENTS BACnetConstructedDataSubscribedRecipients"]]
        //[*, 'SUPPORTED_FORMAT_CLASSES'                BACnetConstructedDataSupportedFormatClasses [validation    '1 == 2'    "TODO: implement me SUPPORTED_FORMAT_CLASSES BACnetConstructedDataSupportedFormatClasses"]]
        //[*, 'SUPPORTED_FORMATS'                       BACnetConstructedDataSupportedFormats [validation    '1 == 2'    "TODO: implement me SUPPORTED_FORMATS BACnetConstructedDataSupportedFormats"]]
        //[*, 'SUPPORTED_SECURITY_ALGORITHMS'           BACnetConstructedDataSupportedSecurityAlgorithms [validation    '1 == 2'    "TODO: implement me SUPPORTED_SECURITY_ALGORITHMS BACnetConstructedDataSupportedSecurityAlgorithms"]]
        //[*, 'SYSTEM_STATUS'                           BACnetConstructedDataSystemStatus [validation    '1 == 2'    "TODO: implement me SYSTEM_STATUS BACnetConstructedDataSystemStatus"]]
        [*, 'TAGS'                                    BACnetConstructedDataTags
            [array    BACnetNameValue
                                tags
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'THREAT_AUTHORITY'                        BACnetConstructedDataThreatAuthority
            [simple   BACnetAccessThreatLevel                                           threatAuthority                 ]
        ]
        [*, 'THREAT_LEVEL'                            BACnetConstructedDataThreatLevel
            [simple   BACnetAccessThreatLevel                                           threatLevel                     ]
        ]
        [*, 'TIME_DELAY'                              BACnetConstructedDataTimeDelay
            [simple   BACnetApplicationTagUnsignedInteger                               timeDelay                       ]
        ]
        [*, 'TIME_DELAY_NORMAL'                       BACnetConstructedDataTimeDelayNormal
            [simple   BACnetApplicationTagUnsignedInteger                               timeDelayNormal                 ]
        ]
        //[*, 'TIME_OF_ACTIVE_TIME_RESET'               BACnetConstructedDataTimeOfActiveTimeReset [validation    '1 == 2'    "TODO: implement me TIME_OF_ACTIVE_TIME_RESET BACnetConstructedDataTimeOfActiveTimeReset"]]
        //[*, 'TIME_OF_DEVICE_RESTART'                  BACnetConstructedDataTimeOfDeviceRestart [validation    '1 == 2'    "TODO: implement me TIME_OF_DEVICE_RESTART BACnetConstructedDataTimeOfDeviceRestart"]]
        //[*, 'TIME_OF_STATE_COUNT_RESET'               BACnetConstructedDataTimeOfStateCountReset [validation    '1 == 2'    "TODO: implement me TIME_OF_STATE_COUNT_RESET BACnetConstructedDataTimeOfStateCountReset"]]
        [*, 'TIME_OF_STRIKE_COUNT_RESET'              BACnetConstructedDataTimeOfStrikeCountReset
            [simple   BACnetDateTime                                        timeOfStrikeCountReset                      ]
        ]
        //[*, 'TIME_SYNCHRONIZATION_INTERVAL'           BACnetConstructedDataTimeSynchronizationInterval [validation    '1 == 2'    "TODO: implement me TIME_SYNCHRONIZATION_INTERVAL BACnetConstructedDataTimeSynchronizationInterval"]]
        //[*, 'TIME_SYNCHRONIZATION_RECIPIENTS'         BACnetConstructedDataTimeSynchronizationRecipients [validation    '1 == 2'    "TODO: implement me TIME_SYNCHRONIZATION_RECIPIENTS BACnetConstructedDataTimeSynchronizationRecipients"]]
        [*, 'TIMER_RUNNING'                           BACnetConstructedDataTimerRunning
            [simple   BACnetApplicationTagBoolean                               timerRunning                 ]
        ]
        [*, 'TIMER_STATE'                             BACnetConstructedDataTimerState
            [simple  BACnetTimerStateTagged('0', 'TagClass.APPLICATION_TAGS')   timerState] 
        ]
        //[*, 'TOTAL_RECORD_COUNT'                      BACnetConstructedDataTotalRecordCount [validation    '1 == 2'    "TODO: implement me TOTAL_RECORD_COUNT BACnetConstructedDataTotalRecordCount"]]
        [*, 'TRACE_FLAG'                              BACnetConstructedDataTraceFlag
            [simple   BACnetApplicationTagBoolean                               traceFlag                               ]
        ]
        //[*, 'TRACKING_VALUE'                          BACnetConstructedDataTrackingValue [validation    '1 == 2'    "TODO: implement me TRACKING_VALUE BACnetConstructedDataTrackingValue"]]
        [*, 'TRANSACTION_NOTIFICATION_CLASS'          BACnetConstructedDataTransactionNotificationClass
            [simple BACnetApplicationTagUnsignedInteger                     transactionNotificationClass                ]
        ]
        //[*, 'TRANSITION'                              BACnetConstructedDataTransition [validation    '1 == 2'    "TODO: implement me TRANSITION BACnetConstructedDataTransition"]]
        //[*, 'TRIGGER'                                 BACnetConstructedDataTrigger [validation    '1 == 2'    "TODO: implement me TRIGGER BACnetConstructedDataTrigger"]]
        [*, 'UNITS'                                   BACnetConstructedDataUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'UPDATE_INTERVAL'                         BACnetConstructedDataUpdateInterval
            [simple   BACnetApplicationTagUnsignedInteger                               updateInterval                  ]
        ]
        //[*, 'UPDATE_KEY_SET_TIMEOUT'                  BACnetConstructedDataUpdateKeySetTimeout [validation    '1 == 2'    "TODO: implement me UPDATE_KEY_SET_TIMEOUT BACnetConstructedDataUpdateKeySetTimeout"]]
        //[*, 'UPDATE_TIME'                             BACnetConstructedDataUpdateTime [validation    '1 == 2'    "TODO: implement me UPDATE_TIME BACnetConstructedDataUpdateTime"]]
        [*, 'USER_EXTERNAL_IDENTIFIER'                BACnetConstructedDataUserExternalIdentifier
            [simple   BACnetApplicationTagCharacterString   userExternalIdentifier                                      ]
        ]
        [*, 'USER_INFORMATION_REFERENCE'              BACnetConstructedDataUserInformationReference
            [simple   BACnetApplicationTagCharacterString   userInformationReference                                    ]
        ]
        [*, 'USER_NAME'                               BACnetConstructedDataUserName
            [simple   BACnetApplicationTagCharacterString   userName                                                    ]
        ]
        [*, 'USER_TYPE'                               BACnetConstructedDataUserType
            [simple   BACnetAccessUserTypeTagged('0', 'TagClass.APPLICATION_TAGS')      userType                        ]
        ]
        [*, 'USES_REMAINING'                          BACnetConstructedDataUsesRemaining
            [simple   BACnetApplicationTagSignedInteger                     usesRemaining                               ]
        ]
        //[*, 'UTC_OFFSET'                              BACnetConstructedDataUtcOffset [validation    '1 == 2'    "TODO: implement me UTC_OFFSET BACnetConstructedDataUtcOffset"]]
        //[*, 'UTC_TIME_SYNCHRONIZATION_RECIPIENTS'     BACnetConstructedDataUtcTimeSynchronizationRecipients [validation    '1 == 2'    "TODO: implement me UTC_TIME_SYNCHRONIZATION_RECIPIENTS BACnetConstructedDataUtcTimeSynchronizationRecipients"]]
        //[*, 'VALID_SAMPLES'                           BACnetConstructedDataValidSamples [validation    '1 == 2'    "TODO: implement me VALID_SAMPLES BACnetConstructedDataValidSamples"]]
        //[*, 'VALUE_BEFORE_CHANGE'                     BACnetConstructedDataValueBeforeChange [validation    '1 == 2'    "TODO: implement me VALUE_BEFORE_CHANGE BACnetConstructedDataValueBeforeChange"]]
        //[*, 'VALUE_CHANGE_TIME'                       BACnetConstructedDataValueChangeTime [validation    '1 == 2'    "TODO: implement me VALUE_CHANGE_TIME BACnetConstructedDataValueChangeTime"]]
        //[*, 'VALUE_SET'                               BACnetConstructedDataValueSet [validation    '1 == 2'    "TODO: implement me VALUE_SET BACnetConstructedDataValueSet"]]
        //[*, 'VALUE_SOURCE'                            BACnetConstructedDataValueSource [validation    '1 == 2'    "TODO: implement me VALUE_SOURCE BACnetConstructedDataValueSource"]]
        //[*, 'VALUE_SOURCE_ARRAY'                      BACnetConstructedDataValueSourceArray [validation    '1 == 2'    "TODO: implement me VALUE_SOURCE_ARRAY BACnetConstructedDataValueSourceArray"]]
        //[*, 'VARIANCE_VALUE'                          BACnetConstructedDataVarianceValue [validation    '1 == 2'    "TODO: implement me VARIANCE_VALUE BACnetConstructedDataVarianceValue"]]
        //[*, 'VENDOR_IDENTIFIER'                       BACnetConstructedDataVendorIdentifier [validation    '1 == 2'    "TODO: implement me VENDOR_IDENTIFIER BACnetConstructedDataVendorIdentifier"]]
        //[*, 'VENDOR_NAME'                             BACnetConstructedDataVendorName [validation    '1 == 2'    "TODO: implement me VENDOR_NAME BACnetConstructedDataVendorName"]]
        [*, 'VERIFICATION_TIME'                       BACnetConstructedDataVerificationTime
            [simple   BACnetApplicationTagSignedInteger                     verificationTime                            ]
        ]
        //[*, 'VIRTUAL_MAC_ADDRESS_TABLE'               BACnetConstructedDataVirtualMacAddressTable [validation    '1 == 2'    "TODO: implement me VIRTUAL_MAC_ADDRESS_TABLE BACnetConstructedDataVirtualMacAddressTable"]]
        //[*, 'VT_CLASSES_SUPPORTED'                    BACnetConstructedDataVtClassesSupported [validation    '1 == 2'    "TODO: implement me VT_CLASSES_SUPPORTED BACnetConstructedDataVtClassesSupported"]]
        [*, 'WEEKLY_SCHEDULE'                         BACnetConstructedDataWeeklySchedule
            [array    BACnetDailySchedule
                                weeklySchedule
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
            [validation 'COUNT(weeklySchedule) == 7' "weeklySchedule should have exactly 7 values"
        ]
        //[*, 'WINDOW_INTERVAL'                         BACnetConstructedDataWindowInterval [validation    '1 == 2'    "TODO: implement me WINDOW_INTERVAL BACnetConstructedDataWindowInterval"]]
        //[*, 'WINDOW_SAMPLES'                          BACnetConstructedDataWindowSamples [validation    '1 == 2'    "TODO: implement me WINDOW_SAMPLES BACnetConstructedDataWindowSamples"]]
        //[*, 'WRITE_STATUS'                            BACnetConstructedDataWriteStatus [validation    '1 == 2'    "TODO: implement me WRITE_STATUS BACnetConstructedDataWriteStatus"]]
        [*, 'ZONE_FROM'                               BACnetConstructedDataZoneFrom
            [simple   BACnetDeviceObjectReference       zoneFrom                                                        ]
        ]
        [*, 'ZONE_MEMBERS'                              BACnetConstructedDataZoneMembers
            [array    BACnetDeviceObjectReference
                    members
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ZONE_TO'                                 BACnetConstructedDataZoneTo
            [simple   BACnetDeviceObjectReference       zoneTo                                                          ]
        ]
        // BACnetConstructedDataUnspecified is used for unmapped properties
        [BACnetConstructedDataUnspecified
            [array    BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifierArgument')
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

[type BACnetConstructedDataElement(BACnetObjectType objectTypeArgument, BACnetPropertyIdentifier propertyIdentifierArgument)
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
    [optional   BACnetConstructedData('peekedTagNumber', 'objectTypeArgument', 'propertyIdentifierArgument')
                            constructedData     'isConstructedData'                                             ]
    [validation '(isApplicationTag && applicationTag != null) || (isContextTag && contextTag != null) || (isConstructedData && constructedData != null)'
                "BACnetConstructedDataElement could not parse anything"                                         ]
]

[type BACnetOptionalBinaryPV
    [peek     BACnetTagHeader
                        peekedTagHeader                 ]
    [virtual  uint 8    peekedTagNumber     'peekedTagHeader.actualTagNumber']
    [typeSwitch peekedTagNumber
        ['0' BACnetOptionalBinaryPVNull
            [simple   BACnetApplicationTagNull
                            nullValue                 ]
        ]
        [BACnetOptionalBinaryPVValue
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
        ['0' BACnetOptionalCharacterStringNull
            [simple   BACnetApplicationTagNull
                                nullValue               ]
        ]
        [BACnetOptionalCharacterStringValue
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
        ['0' BACnetOptionalREALNull
            [simple   BACnetApplicationTagNull
                                nullValue               ]
        ]
        [BACnetOptionalREALValue
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
        ['0' BACnetOptionalUnsignedNull
            [simple   BACnetApplicationTagNull
                                nullValue               ]
        ]
        [BACnetOptionalUnsignedValue
            [simple   BACnetApplicationTagUnsignedInteger
                                unsignedValue           ]
        ]
    ]
]

[type BACnetNameValue
    [simple   BACnetContextTagCharacterString('0', 'BACnetDataType.CHARACTER_STRING')
                            name                                                                                        ]
    [optional BACnetConstructedData('1', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                            value                                                                                       ]
]

[type BACnetNameValueCollection
    [simple   BACnetOpeningTag('0')
                            openingTag                                                                                  ]
    [array    BACnetNameValue
                            members
                                   terminated
                                   'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 0)'             ]
    [simple   BACnetClosingTag('0')
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' BACnetTimerStateChangeValueNull
           [simple  BACnetApplicationTagNull
                            nullValue                                                   ]
       ]
       ['0x1', 'false' BACnetTimerStateChangeValueBoolean
           [simple   BACnetApplicationTagBoolean
                            booleanValue                                                ]
       ]
       ['0x2', 'false' BACnetTimerStateChangeValueUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x3', 'false' BACnetTimerStateChangeValueInteger
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
       ['0x4', 'false' BACnetTimerStateChangeValueReal
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x5', 'false' BACnetTimerStateChangeValueDouble
           [simple  BACnetApplicationTagDouble
                                doubleValue                                             ]
       ]
       ['0x6', 'false' BACnetTimerStateChangeValueOctetString
           [simple   BACnetApplicationTagOctetString
                            octetStringValue                                            ]
       ]
       ['0x7', 'false' BACnetTimerStateChangeValueCharacterString
           [simple   BACnetApplicationTagCharacterString
                            characterStringValue                                        ]
       ]
       ['0x8', 'false' BACnetTimerStateChangeValueBitString
           [simple   BACnetApplicationTagBitString
                            bitStringValue                                              ]
       ]
       ['0x9', 'false' BACnetTimerStateChangeValueEnumerated
           [simple   BACnetApplicationTagEnumerated
                            enumeratedValue                                             ]
       ]
       ['0xA', 'false' BACnetTimerStateChangeValueDate
           [simple   BACnetApplicationTagDate
                            dateValue                                                   ]
       ]
       ['0xB', 'false' BACnetTimerStateChangeValueTime
           [simple   BACnetApplicationTagTime
                            timeValue                                                   ]
       ]
       ['0xC', 'false' BACnetTimerStateChangeValueObjectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                            objectidentifierValue                                       ]
       ]
       ['0', 'true' BACnetTimerStateChangeValueNoValue
           [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                            noValue                                                     ]
       ]
       ['1', 'true' BACnetTimerStateChangeValueConstructedValue
            [simple   BACnetConstructedData('1', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                                        constructedValue                                ]
       ]
       ['2', 'true' BACnetTimerStateChangeValueDateTime
            [simple   BACnetDateTimeEnclosed('2')
                            dateTimeValue                                               ]
       ]
       ['3', 'true' BACnetTimerStateChangeValueLightingCommand
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [validation         'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [typeSwitch peekedTagNumber
        ['0' BACnetSpecialEventPeriodCalendarEntry
            [simple   BACnetCalendarEntry('0')
                                    calendarEntry                                    ]
        ]
        ['1' BACnetSpecialEventPeriodCalendarReference
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                   calendarReference                                ]
        ]
    ]
]

[type BACnetCalendarEntry(uint 8 tagNumber)
    [peek       BACnetTagHeader
                           peekedTagHeader                                          ]
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [validation         'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [typeSwitch peekedTagNumber
        ['0' BACnetCalendarEntryDate
            [simple   BACnetContextTagDate('0', 'BACnetDataType.DATE')
                                        dateValue                                   ]
        ]
        ['1' BACnetCalendarEntryDateRange
            [simple   BACnetDateRangeEnclosed('1')
                                        dateRange                                   ]
        ]
        ['2' BACnetCalendarEntryWeekNDay
            [simple   BACnetWeekNDayTagged('2', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        weekNDay                                    ]
        ]
    ]
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
    [simple BACnetConstructedDataElement('BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE')
                    value                                                                                   ]
]

[type BACnetDailySchedule
    [simple   BACnetOpeningTag('0')
                        openingTag                                                                              ]
    [array    BACnetTimeValue
                    daySchedule
                        terminated
                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 0)'    ]
    [simple   BACnetClosingTag('0')
                    closingTag                                                                              ]
]
