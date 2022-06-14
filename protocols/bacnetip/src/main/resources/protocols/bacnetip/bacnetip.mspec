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
            [optional BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                                            listOfElements              ]
        ]
        ['REMOVE_LIST_ELEMENT' BACnetConfirmedServiceRequestRemoveListElement
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')      objectIdentifier            ]
            [simple   BACnetPropertyIdentifierTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')                 propertyIdentifier          ]
            [optional BACnetContextTagUnsignedInteger('2', 'BACnetDataType.UNSIGNED_INTEGER')               arrayIndex                  ]
            [optional BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                                            listOfElements              ]
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
            [simple   BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
                                                                                                            propertyValue               ]
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
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
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
        ['0' BACnetRecipientDevice
            [simple   BACnetContextTagObjectIdentifier('0', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
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
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                                       serviceParameters            ]
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
            [optional   BACnetConstructedData('3', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(arrayIndex!=null?arrayIndex.payload:null)')
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
            [optional BACnetConstructedData('5', 'objectIdentifier.objectType', 'propertyIdentifier.value', '(propertyArrayIndex!=null?propertyArrayIndex.payload:null)')
                                                                                                            itemData            ]
            [optional BACnetContextTagUnsignedInteger('6', 'BACnetDataType.UNSIGNED_INTEGER')               firstSequenceNumber ]
        ]
        //
        ////


        ////
        // Remote Device Management Services

        ['CONFIRMED_PRIVATE_TRANSFER' BACnetServiceAckConfirmedPrivateTransfer
            [simple   BACnetVendorIdTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')                       vendorId                    ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')           serviceNumber               ]
            [optional BACnetConstructedData('2', 'BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                                        resultBlock                 ]
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
            [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                            commandValue                                            ]
            [simple   BACnetStatusFlagsTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                            statusFlags                                             ]
            [simple   BACnetConstructedData('2', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
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
            [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
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
           [simple   BACnetConstructedData('0', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
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
        [*, 'ABSENTEE_LIMIT', '2'                       BACnetConstructedDataAbsenteeLimit
            [simple   BACnetApplicationTagUnsignedInteger                     absenteeLimit                             ]
        ]
        [*, 'ACCEPTED_MODES'                            BACnetConstructedDataAcceptedModes
            [array    BACnetLifeSafetyModeTagged('0', 'TagClass.APPLICATION_TAGS')
                            acceptedModes              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_ALARM_EVENTS', '9'                  BACnetConstructedDataAccessAlarmEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                    accessAlarmEvents
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_DOORS'                              BACnetConstructedDataAccessDoors
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectReference
                                accessDoors
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCESS_EVENT', '9'                         BACnetConstructedDataAccessEvent
            [simple   BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS') accessEvent]
        ]
        [*, 'ACCESS_EVENT_AUTHENTICATION_FACTOR'        BACnetConstructedDataAccessEventAuthenticationFactor
            [simple   BACnetAuthenticationFactor  accessEventAuthenticationFactor               ]
        ]
        [*, 'ACCESS_EVENT_CREDENTIAL'                   BACnetConstructedDataAccessEventCredential
            [simple   BACnetDeviceObjectReference       accessEventCredential                                           ]
        ]
        [*, 'ACCESS_EVENT_TAG', '2'                     BACnetConstructedDataAccessEventTag
            [simple BACnetApplicationTagUnsignedInteger                     accessEventTag                              ]
        ]
        [*, 'ACCESS_EVENT_TIME'                         BACnetConstructedDataAccessEventTime
            [simple   BACnetTimeStamp                                                accessEventTime                    ]
        ]
        [*, 'ACCESS_TRANSACTION_EVENTS', '9'            BACnetConstructedDataAccessTransactionEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                        accessTransactionEvents
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACCOMPANIMENT'                             BACnetConstructedDataAccompaniment
            [simple   BACnetDeviceObjectReference       accompaniment                                                   ]
        ]
        [*, 'ACCOMPANIMENT_TIME', '2'                   BACnetConstructedDataAccompanimentTime
            [simple   BACnetApplicationTagUnsignedInteger                               accompanimentTime               ]
        ]
        [*, 'ACK_REQUIRED', '9'                         BACnetConstructedDataAckRequired
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS')
                                                                                        ackRequired                     ]
        ]
        [*, 'ACKED_TRANSITIONS', '9'                    BACnetConstructedDataAckedTransitions
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS') ackedTransitions                ]
        ]
        ['LOOP', 'ACTION', '9'                          BACnetConstructedDataLoopAction
            [simple   BACnetActionTagged('0', 'TagClass.APPLICATION_TAGS') action]
        ]
        ['COMMAND', 'ACTION'                            BACnetConstructedDataCommandAction
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetActionList
                            actionLists
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTION'                                    BACnetConstructedDataAction
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetActionList
                            actionLists
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTION_TEXT', '7'                          BACnetConstructedDataActionText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                    actionText
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ACTIVATION_TIME'                           BACnetConstructedDataActivationTime
            [simple   BACnetDateTime    activationTime                                                                  ]
        ]
        [*, 'ACTIVE_AUTHENTICATION_POLICY', '2'         BACnetConstructedDataActiveAuthenticationPolicy
            [simple   BACnetApplicationTagUnsignedInteger                               activeAuthenticationPolicy      ]
        ]
        [*, 'ACTIVE_COV_MULTIPLE_SUBSCRIPTIONS'         BACnetConstructedDataActiveCOVMultipleSubscriptions
            [array    BACnetCOVMultipleSubscription
                                activeCOVMultipleSubscriptions
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTIVE_COV_SUBSCRIPTIONS'                  BACnetConstructedDataActiveCOVSubscriptions
            [array    BACnetCOVSubscription
                            activeCOVSubscriptions
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'ACTIVE_TEXT', '7'                          BACnetConstructedDataActiveText
            [simple   BACnetApplicationTagCharacterString             activeText                                        ]
        ]
        [*, 'ACTIVE_VT_SESSIONS'                        BACnetConstructedDataActiveVTSessions
            [array    BACnetVTSession
                                activeVTSession
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ACTUAL_SHED_LEVEL'                         BACnetConstructedDataActualShedLevel
            [simple  BACnetShedLevel            actualShedLevel                                                         ]
        ]
        ['ACCESS_ZONE', 'ADJUST_VALUE', '3'             BACnetConstructedDataAccessZoneAdjustValue
            [simple   BACnetApplicationTagSignedInteger                               adjustValue                       ]
        ]
        ['PULSE_CONVERTER', 'ADJUST_VALUE', '4'         BACnetConstructedDataPulseConverterAdjustValue
            [simple   BACnetApplicationTagReal                                        adjustValue                       ]
        ]
        [*, 'ADJUST_VALUE', '3'                         BACnetConstructedDataAdjustValue
            [simple   BACnetApplicationTagSignedInteger                               adjustValue                       ]
        ]
        // TODO: pretty sure we need to catch a generic application tag here
        [*, 'ALARM_VALUE', '9'                          BACnetConstructedDataAlarmValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS') binaryPv                                   ]
        ]
        ['ACCESS_DOOR', 'ALARM_VALUES'                  BACnetConstructedDataAccessDoorAlarmValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['ACCESS_ZONE', 'ALARM_VALUES'                  BACnetConstructedDataAccessZoneAlarmValues
            [array    BACnetAccessZoneOccupancyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['BITSTRING_VALUE', 'ALARM_VALUES', '8'         BACnetConstructedDataBitStringValueAlarmValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagBitString
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['CHARACTERSTRING_VALUE', 'ALARM_VALUES'             BACnetConstructedDataCharacterStringValueAlarmValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_POINT', 'ALARM_VALUES'            BACnetConstructedDataLifeSafetyPointAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_ZONE', 'ALARM_VALUES'             BACnetConstructedDataLifeSafetyZoneAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_INPUT', 'ALARM_VALUES', '2'       BACnetConstructedDataMultiStateInputAlarmValues
            [array    BACnetApplicationTagUnsignedInteger
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_VALUE', 'ALARM_VALUES', '2'       BACnetConstructedDataMultiStateValueAlarmValues
            [array    BACnetApplicationTagUnsignedInteger
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['TIMER', 'ALARM_VALUES'                        BACnetConstructedDataTimerAlarmValues
            [array    BACnetTimerStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ALARM_VALUES'                              BACnetConstructedDataAlarmValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ALIGN_INTERVALS', '1'                      BACnetConstructedDataAlignIntervals
            [simple BACnetApplicationTagBoolean                               alignIntervals                            ]
        ]

        /////
        // All property implementations for every object

        ['ACCESS_CREDENTIAL'     , 'ALL'                BACnetConstructedDataAccessCredentialAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_DOOR'           , 'ALL'                BACnetConstructedDataAccessDoorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_POINT'          , 'ALL'                BACnetConstructedDataAccessPointAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_RIGHTS'         , 'ALL'                BACnetConstructedDataAccessRightsAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_USER'           , 'ALL'                BACnetConstructedDataAccessUserAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCESS_ZONE'           , 'ALL'                BACnetConstructedDataAccessZoneAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ACCUMULATOR'           , 'ALL'                BACnetConstructedDataAccumulatorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ALERT_ENROLLMENT'      , 'ALL'                BACnetConstructedDataAlertEnrollmentAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_INPUT'          , 'ALL'                BACnetConstructedDataAnalogInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_OUTPUT'         , 'ALL'                BACnetConstructedDataAnalogOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ANALOG_VALUE'          , 'ALL'                BACnetConstructedDataAnalogValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['AVERAGING'             , 'ALL'                BACnetConstructedDataAveragingAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_INPUT'          , 'ALL'                BACnetConstructedDataBinaryInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'ALL'                BACnetConstructedDataBinaryLightingOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_OUTPUT'         , 'ALL'                BACnetConstructedDataBinaryOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BINARY_VALUE'          , 'ALL'                BACnetConstructedDataBinaryValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['BITSTRING_VALUE'       , 'ALL'                BACnetConstructedDataBitstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CALENDAR'              , 'ALL'                BACnetConstructedDataCalendarAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CHANNEL'               , 'ALL'                BACnetConstructedDataChannelAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CHARACTERSTRING_VALUE' , 'ALL'                BACnetConstructedDataCharacterstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['COMMAND'               , 'ALL'                BACnetConstructedDataCommandAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['CREDENTIAL_DATA_INPUT' , 'ALL'                BACnetConstructedDataCredentialDataInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATEPATTERN_VALUE'     , 'ALL'                BACnetConstructedDataDatepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATE_VALUE'            , 'ALL'                BACnetConstructedDataDateValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATETIMEPATTERN_VALUE' , 'ALL'                BACnetConstructedDataDatetimepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DATETIME_VALUE'        , 'ALL'                BACnetConstructedDataDatetimeValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['DEVICE'                , 'ALL'                BACnetConstructedDataDeviceAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ELEVATOR_GROUP'        , 'ALL'                BACnetConstructedDataElevatorGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['ESCALATOR'             , 'ALL'                BACnetConstructedDataEscalatorAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['EVENT_ENROLLMENT'      , 'ALL'                BACnetConstructedDataEventEnrollmentAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['EVENT_LOG'             , 'ALL'                BACnetConstructedDataEventLogAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['FILE'                  , 'ALL'                BACnetConstructedDataFileAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['GLOBAL_GROUP'          , 'ALL'                BACnetConstructedDataGlobalGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['GROUP'                 , 'ALL'                BACnetConstructedDataGroupAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['INTEGER_VALUE'         , 'ALL'                BACnetConstructedDataIntegerValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LARGE_ANALOG_VALUE'    , 'ALL'                BACnetConstructedDataLargeAnalogValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFE_SAFETY_POINT'     , 'ALL'                BACnetConstructedDataLifeSafetyPointAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFE_SAFETY_ZONE'      , 'ALL'                BACnetConstructedDataLifeSafetyZoneAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIFT'                  , 'ALL'                BACnetConstructedDataLiftAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LIGHTING_OUTPUT'       , 'ALL'                BACnetConstructedDataLightingOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LOAD_CONTROL'          , 'ALL'                BACnetConstructedDataLoadControlAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['LOOP'                  , 'ALL'                BACnetConstructedDataLoopAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_INPUT'     , 'ALL'                BACnetConstructedDataMultiStateInputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_OUTPUT'    , 'ALL'                BACnetConstructedDataMultiStateOutputAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['MULTI_STATE_VALUE'     , 'ALL'                BACnetConstructedDataMultiStateValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NETWORK_PORT'          , 'ALL'                BACnetConstructedDataNetworkPortAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NETWORK_SECURITY'      , 'ALL'                BACnetConstructedDataNetworkSecurityAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NOTIFICATION_CLASS'    , 'ALL'                BACnetConstructedDataNotificationClassAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['NOTIFICATION_FORWARDER', 'ALL'                BACnetConstructedDataNotificationForwarderAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['OCTETSTRING_VALUE'     , 'ALL'                BACnetConstructedDataOctetstringValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['POSITIVE_INTEGER_VALUE', 'ALL'                BACnetConstructedDataPositiveIntegerValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['PROGRAM'               , 'ALL'                BACnetConstructedDataProgramAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['PULSE_CONVERTER'       , 'ALL'                BACnetConstructedDataPulseConverterAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['SCHEDULE'              , 'ALL'                BACnetConstructedDataScheduleAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['STRUCTURED_VIEW'       , 'ALL'                BACnetConstructedDataStructuredViewAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIMEPATTERN_VALUE'     , 'ALL'                BACnetConstructedDataTimepatternValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIME_VALUE'            , 'ALL'                BACnetConstructedDataTimeValueAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TIMER'                 , 'ALL'                BACnetConstructedDataTimerAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TREND_LOG'             , 'ALL'                BACnetConstructedDataTrendLogAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        ['TREND_LOG_MULTIPLE'    , 'ALL'                BACnetConstructedDataTrendLogMultipleAll
            [validation '1==2' "All should never occur in context of constructed data. If it does please report"]
        ]
        //
        /////

        [*, 'ALL_WRITES_SUCCESSFUL', '1'                BACnetConstructedDataAllWritesSuccessful
            [simple   BACnetApplicationTagBoolean                               allWritesSuccessful                     ]
        ]
        [*, 'ALLOW_GROUP_DELAY_INHIBIT', '1'            BACnetConstructedDataAllowGroupDelayInhibit
            [simple   BACnetApplicationTagBoolean                               allowGroupDelayInhibit                  ]
        ]
        [*, 'APDU_LENGTH', '2'                          BACnetConstructedDataAPDULength
            [simple BACnetApplicationTagUnsignedInteger                               apduLength                        ]
        ]
        [*, 'APDU_SEGMENT_TIMEOUT', '2'                 BACnetConstructedDataAPDUSegmentTimeout
            [simple   BACnetApplicationTagUnsignedInteger                               apduSegmentTimeout              ]
        ]
        [*, 'APDU_TIMEOUT', '2'                         BACnetConstructedDataAPDUTimeout
            [simple   BACnetApplicationTagUnsignedInteger                               apduTimeout                     ]
        ]
        [*, 'APPLICATION_SOFTWARE_VERSION', '7'         BACnetConstructedDataApplicationSoftwareVersion
            [simple   BACnetApplicationTagCharacterString                               applicationSoftwareVersion      ]
        ]
        [*, 'ARCHIVE', '1'                              BACnetConstructedDataArchive
            [simple BACnetApplicationTagBoolean                                          archive                        ]
        ]
        [*, 'ASSIGNED_ACCESS_RIGHTS'                    BACnetConstructedDataAssignedAccessRights
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAssignedAccessRights
                                        assignedAccessRights
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ASSIGNED_LANDING_CALLS'                    BACnetConstructedDataAssignedLandingCalls
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAssignedLandingCalls
                                        assignedLandingCalls
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'ATTEMPTED_SAMPLES', '2'                    BACnetConstructedDataAttemptedSamples
            [simple   BACnetApplicationTagUnsignedInteger                               attemptedSamples                ]
        ]
        [*, 'AUTHENTICATION_FACTORS'                    BACnetConstructedDataAuthenticationFactors
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetCredentialAuthenticationFactor
                            authenticationFactors
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'AUTHENTICATION_POLICY_LIST'                BACnetConstructedDataAuthenticationPolicyList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAuthenticationPolicy
                            authenticationPolicyList
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'AUTHENTICATION_POLICY_NAMES', '7'          BACnetConstructedDataAuthenticationPolicyNames
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                        authenticationPolicyNames
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHENTICATION_STATUS', '9'                BACnetConstructedDataAuthenticationStatus
            [simple   BACnetAuthenticationStatusTagged('0', 'TagClass.APPLICATION_TAGS') authenticationStatus           ]
        ]
        [*, 'AUTHORIZATION_EXEMPTIONS', '9'             BACnetConstructedDataAuthorizationExemptions
            [array    BACnetAuthorizationExemptionTagged('0', 'TagClass.APPLICATION_TAGS')
                                        authorizationExemption
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'AUTHORIZATION_MODE', '9'                   BACnetConstructedDataAuthorizationMode
            [simple   BACnetAuthorizationModeTagged('0', 'TagClass.APPLICATION_TAGS') authorizationMode                 ]
        ]
        [*, 'AUTO_SLAVE_DISCOVERY', '1'                 BACnetConstructedDataAutoSlaveDiscovery
            [simple   BACnetApplicationTagBoolean                                       autoSlaveDiscovery              ]
        ]
        [*, 'AVERAGE_VALUE', '4'                        BACnetConstructedDataAverageValue
            [simple BACnetApplicationTagReal                                averageValue                                ]
        ]
        [*, 'BACKUP_AND_RESTORE_STATE', '9'             BACnetConstructedDataBackupAndRestoreState
            [simple   BACnetBackupStateTagged('0', 'TagClass.APPLICATION_TAGS')         backupAndRestoreState           ]
        ]
        [*, 'BACKUP_FAILURE_TIMEOUT', '2'               BACnetConstructedDataBackupFailureTimeout
            [simple   BACnetApplicationTagUnsignedInteger                               backupFailureTimeout            ]
        ]
        [*, 'BACKUP_PREPARATION_TIME', '2'              BACnetConstructedDataBackupPreparationTime
            [simple   BACnetApplicationTagUnsignedInteger                               backupPreparationTime           ]
        ]
        [*, 'BACNET_IP_GLOBAL_ADDRESS'                  BACnetConstructedDataBACnetIPGlobalAddress
            [simple   BACnetHostNPort                   bacnetIpGlobalAddress                                           ]
        ]
        [*, 'BACNET_IP_MODE', '9'                       BACnetConstructedDataBACnetIPMode
            [simple BACnetIPModeTagged('0', 'TagClass.APPLICATION_TAGS') bacnetIpMode                                   ]
        ]
        [*, 'BACNET_IP_MULTICAST_ADDRESS', '6'          BACnetConstructedDataBACnetIPMulticastAddress
            [simple   BACnetApplicationTagOctetString   ipMulticastAddress                                              ]
        ]
        [*, 'BACNET_IP_NAT_TRAVERSAL', '1'              BACnetConstructedDataBACnetIPNATTraversal
            [simple BACnetApplicationTagBoolean                               bacnetIPNATTraversal                      ]
        ]
        [*, 'BACNET_IP_UDP_PORT', '2'                   BACnetConstructedDataBACnetIPUDPPort
            [simple BACnetApplicationTagUnsignedInteger                               ipUdpPort                         ]
        ]
        [*, 'BACNET_IPV6_MODE', '9'                     BACnetConstructedDataBACnetIPv6Mode
            [simple   BACnetIPModeTagged('0', 'TagClass.APPLICATION_TAGS') bacnetIpv6Mode                               ]
        ]
        [*, 'BACNET_IPV6_UDP_PORT', '2'                 BACnetConstructedDataBACnetIPv6UDPPort
            [simple BACnetApplicationTagUnsignedInteger                               ipv6UdpPort                       ]
        ]
        [*, 'BACNET_IPV6_MULTICAST_ADDRESS', '6'        BACnetConstructedDataBACnetIPv6MulticastAddress
            [simple   BACnetApplicationTagOctetString   ipv6MulticastAddress                                            ]
        ]
        [*, 'BASE_DEVICE_SECURITY_POLICY', '9'          BACnetConstructedDataBaseDeviceSecurityPolicy
            [simple BACnetSecurityLevelTagged('0', 'TagClass.APPLICATION_TAGS')     baseDeviceSecurityPolicy            ]
        ]
        [*, 'BBMD_ACCEPT_FD_REGISTRATIONS', '1'         BACnetConstructedDataBBMDAcceptFDRegistrations
            [simple BACnetApplicationTagBoolean                               bbmdAcceptFDRegistrations                 ]
        ]
        [*, 'BBMD_BROADCAST_DISTRIBUTION_TABLE'         BACnetConstructedDataBBMDBroadcastDistributionTable
            [array    BACnetBDTEntry
                                bbmdBroadcastDistributionTable
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'BBMD_FOREIGN_DEVICE_TABLE'                 BACnetConstructedDataBBMDForeignDeviceTable
            [array    BACnetBDTEntry
                                bbmdForeignDeviceTable
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'BELONGS_TO'                                BACnetConstructedDataBelongsTo
            [simple   BACnetDeviceObjectReference       belongsTo                                                       ]
        ]
        [*, 'BIAS', '4'                                 BACnetConstructedDataBias
            [simple   BACnetApplicationTagReal                                          bias                            ]
        ]
        [*, 'BIT_MASK', '8'                             BACnetConstructedDataBitMask
            [simple   BACnetApplicationTagBitString     bitString                                                       ]
        ]
        [*, 'BIT_TEXT', '7'                             BACnetConstructedDataBitText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                        bitText
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'BLINK_WARN_ENABLE', '1'                    BACnetConstructedDataBlinkWarnEnable
            [simple BACnetApplicationTagBoolean                               blinkWarnEnable                           ]
        ]
        [*, 'BUFFER_SIZE', '2'                          BACnetConstructedDataBufferSize
            [simple BACnetApplicationTagUnsignedInteger                       bufferSize                                ]
        ]
        [*, 'CAR_ASSIGNED_DIRECTION', '9'               BACnetConstructedDataCarAssignedDirection
            [simple   BACnetLiftCarDirectionTagged('0', 'TagClass.APPLICATION_TAGS')             assignedDirection      ]
        ]
        [*, 'CAR_DOOR_COMMAND', '9'                     BACnetConstructedDataCarDoorCommand
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLiftCarDoorCommandTagged('0', 'TagClass.APPLICATION_TAGS')
                            carDoorCommand
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'CAR_DOOR_STATUS', '9'                      BACnetConstructedDataCarDoorStatus
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDoorStatusTagged('0', 'TagClass.APPLICATION_TAGS')
                            carDoorStatus
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'CAR_DOOR_TEXT', '7'                        BACnetConstructedDataCarDoorText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                carDoorText
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CAR_DOOR_ZONE', '1'                        BACnetConstructedDataCarDoorZone
            [simple   BACnetApplicationTagBoolean               carDoorZone                                             ]
        ]
        [*, 'CAR_DRIVE_STATUS', '9'                     BACnetConstructedDataCarDriveStatus
            [simple   BACnetLiftCarDriveStatusTagged('0', 'TagClass.APPLICATION_TAGS')                  carDriveStatus  ]
        ]
        [*, 'CAR_LOAD', '4'                             BACnetConstructedDataCarLoad
            [simple   BACnetApplicationTagReal                                          carLoad                         ]
        ]
        [*, 'CAR_LOAD_UNITS', '9'                       BACnetConstructedDataCarLoadUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'CAR_MODE', '9'                             BACnetConstructedDataCarMode
            [simple   BACnetLiftCarModeTagged('0', 'TagClass.APPLICATION_TAGS')                    carMode              ]
        ]
        [*, 'CAR_MOVING_DIRECTION', '9'                 BACnetConstructedDataCarMovingDirection
            [simple   BACnetLiftCarDirectionTagged('0', 'TagClass.APPLICATION_TAGS')             carMovingDirection     ]
        ]
        [*, 'CAR_POSITION', '2'                         BACnetConstructedDataCarPosition
            [simple   BACnetApplicationTagUnsignedInteger                     carPosition                               ]
        ]
        [*, 'CHANGE_OF_STATE_COUNT', '2'                BACnetConstructedDataChangeOfStateCount
            [simple   BACnetApplicationTagUnsignedInteger                     changeIfStateCount                        ]
        ]
        [*, 'CHANGE_OF_STATE_TIME'                      BACnetConstructedDataChangeOfStateTime
            [simple   BACnetDateTime                                          changeOfStateTime                         ]
        ]
        [*, 'CHANGES_PENDING', '1'                      BACnetConstructedDataChangesPending
            [simple   BACnetApplicationTagBoolean                                       changesPending                  ]
        ]
        [*, 'CHANNEL_NUMBER', '2'                       BACnetConstructedDataChannelNumber
            [simple   BACnetApplicationTagUnsignedInteger                     channelNumber                             ]
        ]
        [*, 'CLIENT_COV_INCREMENT'                      BACnetConstructedDataClientCOVIncrement
            [simple   BACnetClientCOV                 covIncrement                                                      ]
        ]
        [*, 'COMMAND', '9'                              BACnetConstructedDataCommand
            [simple   BACnetNetworkPortCommandTagged('0', 'TagClass.APPLICATION_TAGS')                  command         ]
        ]
        [*, 'COMMAND_TIME_ARRAY'                        BACnetConstructedDataCommandTimeArray
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
        [*, 'CONFIGURATION_FILES', '12'                 BACnetConstructedDataConfigurationFiles
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                                configurationFiles
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CONTROL_GROUPS', '2'                       BACnetConstructedDataControlGroups
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                controlGroups
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CONTROLLED_VARIABLE_REFERENCE'             BACnetConstructedDataControlledVariableReference
            [simple   BACnetObjectPropertyReference                                     controlledVariableReference     ]
        ]
        [*, 'CONTROLLED_VARIABLE_UNITS', '9'            BACnetConstructedDataControlledVariableUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'CONTROLLED_VARIABLE_VALUE', '4'            BACnetConstructedDataControlledVariableValue
            [simple   BACnetApplicationTagReal                                          controlledVariableValue         ]
        ]
        [*, 'COUNT', '2'                                BACnetConstructedDataCount
            [simple BACnetApplicationTagUnsignedInteger                               count                             ]
        ]
        [*, 'COUNT_BEFORE_CHANGE', '2'                  BACnetConstructedDataCountBeforeChange
            [simple BACnetApplicationTagUnsignedInteger                               countBeforeChange                 ]
        ]
        [*, 'COUNT_CHANGE_TIME'                         BACnetConstructedDataCountChangeTime
            [simple   BACnetDateTime                                        countChangeTime                             ]
        ]
        ['INTEGER_VALUE', 'COV_INCREMENT', '2'          BACnetConstructedDataIntegerValueCOVIncrement
            [simple   BACnetApplicationTagUnsignedInteger                               covIncrement                    ]
        ]
        ['LARGE_ANALOG_VALUE', 'COV_INCREMENT', '5'     BACnetConstructedDataLargeAnalogValueCOVIncrement
            [simple   BACnetApplicationTagDouble                                        covIncrement                    ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'COV_INCREMENT', '2' BACnetConstructedDataPositiveIntegerValueCOVIncrement
            [simple   BACnetApplicationTagUnsignedInteger                               covIncrement                    ]
        ]
        [*, 'COV_INCREMENT', '4'                        BACnetConstructedDataCOVIncrement
            [simple   BACnetApplicationTagReal                                          covIncrement                    ]
        ]
        [*, 'COV_PERIOD', '2'                           BACnetConstructedDataCOVPeriod
            [simple   BACnetApplicationTagUnsignedInteger                               covPeriod                       ]
        ]
        [*, 'COV_RESUBSCRIPTION_INTERVAL', '2'          BACnetConstructedDataCOVResubscriptionInterval
            [simple   BACnetApplicationTagUnsignedInteger                               covResubscriptionInterval       ]
        ]
        [*, 'COVU_PERIOD', '2'                          BACnetConstructedDataCOVUPeriod
            [simple   BACnetApplicationTagUnsignedInteger                               covuPeriod                      ]
        ]
        [*, 'COVU_RECIPIENTS'                           BACnetConstructedDataCOVURecipients
            [array    BACnetRecipient
                                        covuRecipients
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'CREDENTIAL_DISABLE', '9'                   BACnetConstructedDataCredentialDisable
            [simple   BACnetAccessCredentialDisableTagged('0', 'TagClass.APPLICATION_TAGS')             credentialDisable]
        ]
        [*, 'CREDENTIAL_STATUS', '9'                    BACnetConstructedDataCredentialStatus
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    binaryPv                    ]
        ]
        [*, 'CREDENTIALS'                               BACnetConstructedDataCredentials
            [array    BACnetDeviceObjectReference
                        credentials
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'CREDENTIALS_IN_ZONE'                       BACnetConstructedDataCredentialsInZone
            [array    BACnetDeviceObjectReference
                        credentialsInZone
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'CURRENT_COMMAND_PRIORITY'                  BACnetConstructedDataCurrentCommandPriority
            [simple   BACnetOptionalUnsigned                currentCommandPriority                                      ]
        ]
        [*, 'DATABASE_REVISION', '2'                    BACnetConstructedDataDatabaseRevision
            [simple   BACnetApplicationTagUnsignedInteger                     databaseRevision                          ]
        ]
        [*, 'DATE_LIST'                                 BACnetConstructedDataDateList
            [array    BACnetCalendarEntry
                        dateList
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'DAYLIGHT_SAVINGS_STATUS', '1'              BACnetConstructedDataDaylightSavingsStatus
            [simple   BACnetApplicationTagBoolean                               daylightSavingsStatus                   ]
        ]
        [*, 'DAYS_REMAINING', '3'                       BACnetConstructedDataDaysRemaining
            [simple   BACnetApplicationTagSignedInteger                     daysRemaining                               ]
        ]
        ['INTEGER_VALUE', 'DEADBAND', '2'               BACnetConstructedDataIntegerValueDeadband
            [simple   BACnetApplicationTagUnsignedInteger                               deadband                        ]
        ]
        ['LARGE_ANALOG_VALUE', 'DEADBAND', '5'          BACnetConstructedDataLargeAnalogValueDeadband
            [simple   BACnetApplicationTagDouble                                        deadband                        ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'DEADBAND', '2'      BACnetConstructedDataPositiveIntegerValueDeadband
            [simple   BACnetApplicationTagUnsignedInteger                               deadband                        ]
        ]
        [*, 'DEADBAND', '4'                             BACnetConstructedDataDeadband
            [simple   BACnetApplicationTagReal                                          deadband                        ]
        ]
        [*, 'DEFAULT_FADE_TIME', '2'                    BACnetConstructedDataDefaultFadeTime
            [simple   BACnetApplicationTagUnsignedInteger                     defaultFadeTime                           ]
        ]
        [*, 'DEFAULT_RAMP_RATE', '4'                    BACnetConstructedDataDefaultRampRate
            [simple   BACnetApplicationTagReal                                          defaultRampRate                 ]
        ]
        [*, 'DEFAULT_STEP_INCREMENT', '4'               BACnetConstructedDataDefaultStepIncrement
            [simple   BACnetApplicationTagReal                                          defaultStepIncrement            ]
        ]
        [*, 'DEFAULT_SUBORDINATE_RELATIONSHIP', '9'     BACnetConstructedDataDefaultSubordinateRelationship
            [simple   BACnetRelationshipTagged('0', 'TagClass.APPLICATION_TAGS') defaultSubordinateRelationship         ]
        ]
        [*, 'DEFAULT_TIMEOUT', '2'                      BACnetConstructedDataDefaultTimeout
            [simple   BACnetApplicationTagUnsignedInteger                     defaultTimeout                            ]
        ]
        [*, 'DEPLOYED_PROFILE_LOCATION', '7'            BACnetConstructedDataDeployedProfileLocation
           [simple   BACnetApplicationTagCharacterString deployedProfileLocation                                        ]
        ]
        [*, 'DERIVATIVE_CONSTANT', '4'                  BACnetConstructedDataDerivativeConstant
            [simple   BACnetApplicationTagReal                                          derivativeConstant              ]
        ]
        [*, 'DERIVATIVE_CONSTANT_UNITS', '9'            BACnetConstructedDataDerivativeConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'DESCRIPTION', '7'                          BACnetConstructedDataDescription
            [simple   BACnetApplicationTagCharacterString                               description                     ]
        ]
        [*, 'DESCRIPTION_OF_HALT', '7'                  BACnetConstructedDataDescriptionOfHalt
            [simple   BACnetApplicationTagCharacterString   descriptionForHalt                                          ]
        ]
        [*, 'DEVICE_ADDRESS_BINDING'                    BACnetConstructedDataDeviceAddressBinding
            [array    BACnetAddressBinding
                                deviceAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'DEVICE_TYPE', '7'                          BACnetConstructedDataDeviceType
            [simple   BACnetApplicationTagCharacterString                               deviceType                      ]
        ]
        [*, 'DIRECT_READING', '4'                       BACnetConstructedDataDirectReading
            [simple   BACnetApplicationTagReal                                          directReading                   ]
        ]
        [*, 'DISTRIBUTION_KEY_REVISION', '2'            BACnetConstructedDataDistributionKeyRevision
            [simple BACnetApplicationTagUnsignedInteger                       distributionKeyRevision                   ]
        ]
        [*, 'DO_NOT_HIDE', '1'                          BACnetConstructedDataDoNotHide
            [simple   BACnetApplicationTagBoolean                                       doNotHide                       ]
        ]
        [*, 'DOOR_ALARM_STATE', '9'                     BACnetConstructedDataDoorAlarmState
            [simple BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            doorAlarmState  ]
        ]
        [*, 'DOOR_EXTENDED_PULSE_TIME', '2'             BACnetConstructedDataDoorExtendedPulseTime
            [simple BACnetApplicationTagUnsignedInteger                     doorExtendedPulseTime                       ]
        ]
        [*, 'DOOR_MEMBERS'                              BACnetConstructedDataDoorMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectReference
                        doorMembers
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'DOOR_OPEN_TOO_LONG_TIME', '2'              BACnetConstructedDataDoorOpenTooLongTime
            [simple BACnetApplicationTagUnsignedInteger                     doorOpenTooLongTime                         ]
        ]
        [*, 'DOOR_PULSE_TIME', '2'                      BACnetConstructedDataDoorPulseTime
            [simple BACnetApplicationTagUnsignedInteger                     doorPulseTime                               ]
        ]
        [*, 'DOOR_STATUS', '9'                          BACnetConstructedDataDoorStatus
            [simple   BACnetDoorStatusTagged('0', 'TagClass.APPLICATION_TAGS')
                                            doorStatus                      ]
        ]
        [*, 'DOOR_UNLOCK_DELAY_TIME', '2'               BACnetConstructedDataDoorUnlockDelayTime
            [simple   BACnetApplicationTagUnsignedInteger                     doorUnlockDelayTime                       ]
        ]
        [*, 'DUTY_WINDOW', '2'                          BACnetConstructedDataDutyWindow
            [simple   BACnetApplicationTagUnsignedInteger                     dutyWindow                                ]
        ]
        [*, 'EFFECTIVE_PERIOD'                          BACnetConstructedDataEffectivePeriod
            [simple   BACnetDateRange               dateRange   ]
        ]
        [*, 'EGRESS_ACTIVE', '1'                        BACnetConstructedDataEgressActive
            [simple   BACnetApplicationTagBoolean                               egressActive                            ]
        ]
        [*, 'EGRESS_TIME', '2'                          BACnetConstructedDataEgressTime
            [simple   BACnetApplicationTagUnsignedInteger                     egressTime                                ]
        ]
        [*, 'ELAPSED_ACTIVE_TIME', '2'                  BACnetConstructedDataElapsedActiveTime
            [simple   BACnetApplicationTagUnsignedInteger                     elapsedActiveTime                         ]
        ]
        [*, 'ELEVATOR_GROUP', '12'                      BACnetConstructedDataElevatorGroup
            [simple   BACnetApplicationTagObjectIdentifier                    elevatorGroup                             ]
        ]
        [*, 'ENABLE', '1'                               BACnetConstructedDataEnable
            [simple   BACnetApplicationTagBoolean                               enable                                  ]
        ]
        [*, 'ENERGY_METER', '4'                         BACnetConstructedDataEnergyMeter
            [simple   BACnetApplicationTagReal                                          energyMeter                     ]
        ]
        [*, 'ENERGY_METER_REF'                          BACnetConstructedDataEnergyMeterRef
            [simple   BACnetDeviceObjectReference                                       energyMeterRef                  ]
        ]
        [*, 'ENTRY_POINTS'                              BACnetConstructedDataEntryPoints
            [array    BACnetDeviceObjectReference
                                entryPoints
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'ERROR_LIMIT', '4'                          BACnetConstructedDataErrorLimit
            [simple   BACnetApplicationTagReal                                          errorLimit                      ]
        ]
        [*, 'ESCALATOR_MODE', '9'                       BACnetConstructedDataEscalatorMode
            [simple   BACnetEscalatorModeTagged('0', 'TagClass.APPLICATION_TAGS')       escalatorMode                   ]
        ]
        [*, 'EVENT_ALGORITHM_INHIBIT', '1'              BACnetConstructedDataEventAlgorithmInhibit
            [simple   BACnetApplicationTagBoolean                                       eventAlgorithmInhibit           ]
        ]
        [*, 'EVENT_ALGORITHM_INHIBIT_REF'               BACnetConstructedDataEventAlgorithmInhibitRef
            [simple   BACnetObjectPropertyReference                                     eventAlgorithmInhibitRef        ]
        ]
        [*, 'EVENT_DETECTION_ENABLE', '1'               BACnetConstructedDataEventDetectionEnable
            [simple   BACnetApplicationTagBoolean                                       eventDetectionEnable            ]
        ]
        [*, 'EVENT_ENABLE', '8'                         BACnetConstructedDataEventEnable
            [simple   BACnetEventTransitionBitsTagged('0', 'TagClass.APPLICATION_TAGS') eventEnable                     ]
        ]
        [*, 'EVENT_MESSAGE_TEXTS'                       BACnetConstructedDataEventMessageTexts
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                                        eventMessageTexts
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [virtual BACnetOptionalCharacterString    toOffnormalText    'COUNT(eventMessageTexts)==3?eventMessageTexts[0]:null'            ]
            [virtual BACnetOptionalCharacterString    toFaultText        'COUNT(eventMessageTexts)==3?eventMessageTexts[1]:null'            ]
            [virtual BACnetOptionalCharacterString    toNormalText       'COUNT(eventMessageTexts)==3?eventMessageTexts[2]:null'            ]
            [validation 'arrayIndexArgument!=null || COUNT(eventMessageTexts) == 3'
                                    "eventMessageTexts should have exactly 3 values"                                    ]
        ]
        [*, 'EVENT_MESSAGE_TEXTS_CONFIG'                BACnetConstructedDataEventMessageTextsConfig
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                                        eventMessageTextsConfig
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [virtual BACnetOptionalCharacterString    toOffnormalTextConfig    'COUNT(eventMessageTextsConfig)==3?eventMessageTextsConfig[0]:null']
            [virtual BACnetOptionalCharacterString    toFaultTextConfig        'COUNT(eventMessageTextsConfig)==3?eventMessageTextsConfig[1]:null']
            [virtual BACnetOptionalCharacterString    toNormalTextConfig       'COUNT(eventMessageTextsConfig)==3?eventMessageTextsConfig[2]:null']
            [validation 'arrayIndexArgument!=null || COUNT(eventMessageTextsConfig) == 3'
                        "eventMessageTextsConfig should have exactly 3 values"                                          ]
        ]
        [*, 'EVENT_PARAMETERS'                          BACnetConstructedDataEventParameters
            [simple   BACnetEventParameter                  eventParameter                                              ]
        ]
        [*, 'EVENT_STATE', '9'                          BACnetConstructedDataEventState
            [simple   BACnetEventStateTagged('0', 'TagClass.APPLICATION_TAGS')          eventState                      ]
        ]
        [*, 'EVENT_TIME_STAMPS'                         BACnetConstructedDataEventTimeStamps
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetTimeStamp
                                        eventTimeStamps
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [virtual BACnetTimeStamp    toOffnormal    'COUNT(eventTimeStamps)==3?eventTimeStamps[0]:null'                  ]
            [virtual BACnetTimeStamp    toFault        'COUNT(eventTimeStamps)==3?eventTimeStamps[1]:null'                  ]
            [virtual BACnetTimeStamp    toNormal       'COUNT(eventTimeStamps)==3?eventTimeStamps[2]:null'                  ]
            [validation 'arrayIndexArgument!=null || COUNT(eventTimeStamps) == 3'
                        "eventTimeStamps should have exactly 3 values"                                                  ]
        ]
        [*, 'EVENT_TYPE', '9'                           BACnetConstructedDataEventType
            [simple  BACnetEventTypeTagged('0', 'TagClass.APPLICATION_TAGS')        eventType                           ]
        ]
        [*, 'EXCEPTION_SCHEDULE'                        BACnetConstructedDataExceptionSchedule
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetSpecialEvent
                            exceptionSchedule
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'EXECUTION_DELAY', '2'                      BACnetConstructedDataExecutionDelay
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                        executionDelay
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'EXIT_POINTS'                               BACnetConstructedDataExitPoints
            [array    BACnetDeviceObjectReference
                        exitPoints
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'EXPECTED_SHED_LEVEL'                       BACnetConstructedDataExpectedShedLevel
            [simple  BACnetShedLevel            expectedShedLevel                                                       ]
        ]
        [*, 'EXPIRATION_TIME'                           BACnetConstructedDataExpirationTime
            [simple   BACnetDateTime            expirationTime                                                          ]
        ]
        [*, 'EXTENDED_TIME_ENABLE', '1'                 BACnetConstructedDataExtendedTimeEnable
            [simple   BACnetApplicationTagBoolean                               extendedTimeEnable                      ]
        ]
        [*, 'FAILED_ATTEMPT_EVENTS'                     BACnetConstructedDataFailedAttemptEvents
            [array    BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')
                                        failedAttemptEvents
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAILED_ATTEMPTS', '2'                      BACnetConstructedDataFailedAttempts
            [simple   BACnetApplicationTagUnsignedInteger                               failedAttempts                  ]
        ]
        [*, 'FAILED_ATTEMPTS_TIME', '2'                 BACnetConstructedDataFailedAttemptsTime
            [simple   BACnetApplicationTagUnsignedInteger                               failedAttemptsTime              ]
        ]
        ['ACCUMULATOR', 'FAULT_HIGH_LIMIT', '2'         BACnetConstructedDataAccumulatorFaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               faultHighLimit                  ]
        ]
        ['ANALOG_INPUT', 'FAULT_HIGH_LIMIT', '4'        BACnetConstructedDataAnalogInputFaultHighLimit
            [simple   BACnetApplicationTagReal                                          faultHighLimit                  ]
        ]
        ['ANALOG_VALUE', 'FAULT_HIGH_LIMIT', '4'        BACnetConstructedDataAnalogValueFaultHighLimit
            [simple   BACnetApplicationTagReal                                          faultHighLimit                  ]
        ]
        ['INTEGER_VALUE', 'FAULT_HIGH_LIMIT', '3'       BACnetConstructedDataIntegerValueFaultHighLimit
            [simple   BACnetApplicationTagSignedInteger                                 faultHighLimit                  ]
        ]
        ['LARGE_ANALOG_VALUE', 'FAULT_HIGH_LIMIT', '5'  BACnetConstructedDataLargeAnalogValueFaultHighLimit
            [simple   BACnetApplicationTagDouble                                        faultHighLimit                  ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'FAULT_HIGH_LIMIT', '2'   BACnetConstructedDataPositiveIntegerValueFaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               faultHighLimit                  ]
        ]
        [*, 'FAULT_HIGH_LIMIT', '2'                     BACnetConstructedDataFaultHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               faultHighLimit                  ]
        ]
        ['ACCUMULATOR', 'FAULT_LOW_LIMIT', '2'          BACnetConstructedDataAccumulatorFaultLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                               faultLowLimit                   ]
        ]
        ['ANALOG_INPUT', 'FAULT_LOW_LIMIT', '4'         BACnetConstructedDataAnalogInputFaultLowLimit
            [simple   BACnetApplicationTagReal                                          faultLowLimit                   ]
        ]
        ['ANALOG_VALUE', 'FAULT_LOW_LIMIT', '4'         BACnetConstructedDataAnalogValueFaultLowLimit
            [simple   BACnetApplicationTagReal                                          faultLowLimit                   ]
        ]
        ['LARGE_ANALOG_VALUE', 'FAULT_LOW_LIMIT', '5'   BACnetConstructedDataLargeAnalogValueFaultLowLimit
            [simple   BACnetApplicationTagDouble                                        faultLowLimit                   ]
        ]
        ['INTEGER_VALUE', 'FAULT_LOW_LIMIT', '3'        BACnetConstructedDataIntegerValueFaultLowLimit
            [simple   BACnetApplicationTagSignedInteger                                 faultLowLimit                   ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'FAULT_LOW_LIMIT', '2'    BACnetConstructedDataPositiveIntegerValueFaultLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                               faultLowLimit                   ]
        ]
        [*, 'FAULT_LOW_LIMIT', '4'                      BACnetConstructedDataFaultLowLimit
            [simple   BACnetApplicationTagReal                                          faultLowLimit                   ]
        ]
        [*, 'FAULT_PARAMETERS'                          BACnetConstructedDataFaultParameters
            [simple   BACnetFaultParameter              faultParameters                                                 ]
        ]
        ['ESCALATOR', 'FAULT_SIGNALS'                   BACnetConstructedDataEscalatorFaultSignals
            [array   BACnetEscalatorFaultTagged('0', 'TagClass.APPLICATION_TAGS')
                                    faultSignals
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFT', 'FAULT_SIGNALS'                        BACnetConstructedDataLiftFaultSignals
            [array   BACnetLiftFaultTagged('0', 'TagClass.APPLICATION_TAGS')
                                    faultSignals
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAULT_SIGNALS'                             BACnetConstructedDataFaultSignals
            [array   BACnetLiftFaultTagged('0', 'TagClass.APPLICATION_TAGS')
                                    faultSignals
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAULT_TYPE', '9'                           BACnetConstructedDataFaultType
            [simple   BACnetFaultTypeTagged('0', 'TagClass.APPLICATION_TAGS')           faultType                       ]
        ]
        ['ACCESS_DOOR', 'FAULT_VALUES'        BACnetConstructedDataAccessDoorFaultValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['CHARACTERSTRING_VALUE', 'FAULT_VALUES'        BACnetConstructedDataCharacterStringValueFaultValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetOptionalCharacterString
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_POINT', 'FAULT_VALUES'            BACnetConstructedDataLifeSafetyPointFaultValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['LIFE_SAFETY_ZONE', 'FAULT_VALUES'             BACnetConstructedDataLifeSafetyZoneFaultValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_INPUT', 'FAULT_VALUES', '2'       BACnetConstructedDataMultiStateInputFaultValues
            [array    BACnetApplicationTagUnsignedInteger
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        ['MULTI_STATE_VALUE', 'FAULT_VALUES', '2'       BACnetConstructedDataMultiStateValueFaultValues
            [array    BACnetApplicationTagUnsignedInteger
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FAULT_VALUES'                              BACnetConstructedDataFaultValues
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            faultValues
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'FD_BBMD_ADDRESS'                           BACnetConstructedDataFDBBMDAddress
            [simple   BACnetHostNPort          fDBBMDAddress]
        ]
        [*, 'FD_SUBSCRIPTION_LIFETIME', '2'             BACnetConstructedDataFDSubscriptionLifetime
            [simple   BACnetApplicationTagUnsignedInteger                               fdSubscriptionLifetime          ]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'FEEDBACK_VALUE', '9'    BACnetConstructedDataBinaryLightingOutputFeedbackValue
            [simple   BACnetBinaryLightingPVTagged('0', 'TagClass.APPLICATION_TAGS')    feedbackValue                   ]
        ]
        ['BINARY_OUTPUT', 'FEEDBACK_VALUE', '9'         BACnetConstructedDataBinaryOutputFeedbackValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')    feedbackValue                           ]
        ]
        ['LIGHTING_OUTPUT', 'FEEDBACK_VALUE', '4'       BACnetConstructedDataLightingOutputFeedbackValue
            [simple  BACnetApplicationTagReal                                  feedbackValue                            ]
        ]
        ['MULTI_STATE_OUTPUT', 'FEEDBACK_VALUE', '2'    BACnetConstructedDataMultiStateOutputFeedbackValue
            [simple  BACnetApplicationTagUnsignedInteger                       feedbackValue                            ]
        ]
        // TODO: similar to ALARM_VALUE either we catch it or we just exlude
        //[*, 'FEEDBACK_VALUE'                          BACnetConstructedDataFeedbackValue [validation    '1 == 2'    "TODO: implement me FEEDBACK_VALUE BACnetConstructedDataFeedbackValue"]]
        [*, 'FILE_ACCESS_METHOD', '9'                   BACnetConstructedDataFileAccessMethod
            [simple   BACnetFileAccessMethodTagged('0', 'TagClass.APPLICATION_TAGS')     fileAccessMethod               ]
        ]
        [*, 'FILE_SIZE', '2'                            BACnetConstructedDataFileSize
            [simple   BACnetApplicationTagUnsignedInteger                                fileSize                       ]
        ]
        [*, 'FILE_TYPE', '7'                            BACnetConstructedDataFileType
            [simple   BACnetApplicationTagCharacterString                                   fileType                    ]
        ]
        [*, 'FIRMWARE_REVISION', '7'                    BACnetConstructedDataFirmwareRevision
           [simple   BACnetApplicationTagCharacterString    firmwareRevision                                            ]
        ]
        [*, 'FLOOR_TEXT', '7'                           BACnetConstructedDataFloorText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'FULL_DUTY_BASELINE', '4'                   BACnetConstructedDataFullDutyBaseline
            [simple   BACnetApplicationTagReal                                          fullDutyBaseLine                ]
        ]
        [*, 'GLOBAL_IDENTIFIER', '2'                    BACnetConstructedDataGlobalIdentifier
            [simple   BACnetApplicationTagUnsignedInteger                           globalIdentifier                    ]
        ]
        [*, 'GROUP_ID', '2'                             BACnetConstructedDataGroupID
            [simple   BACnetApplicationTagUnsignedInteger                           groupId                             ]
        ]
        [*, 'GROUP_MEMBER_NAMES', '7'                   BACnetConstructedDataGroupMemberNames
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                            groupMemberNames
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['GLOBAL_GROUP', 'GROUP_MEMBERS'                BACnetConstructedDataGlobalGroupGroupMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectPropertyReference
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['ELEVATOR_GROUP', 'GROUP_MEMBERS', '12'        BACnetConstructedDataElevatorGroupGroupMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'GROUP_MEMBERS', '12'                       BACnetConstructedDataGroupMembers
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'GROUP_MODE', '9'                           BACnetConstructedDataGroupMode
            [simple   BACnetLiftGroupModeTagged('0', 'TagClass.APPLICATION_TAGS')       groupMode                       ]
        ]
        ['ACCUMULATOR', 'HIGH_LIMIT', '2'               BACnetConstructedDataAccumulatorHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               highLimit                       ]
        ]
        ['LARGE_ANALOG_VALUE', 'HIGH_LIMIT', '5'        BACnetConstructedDataLargeAnalogValueHighLimit
            [simple   BACnetApplicationTagDouble                                        highLimit                       ]
        ]
        ['INTEGER_VALUE', 'HIGH_LIMIT', '3'             BACnetConstructedDataIntegerValueHighLimit
            [simple   BACnetApplicationTagSignedInteger                                 highLimit                       ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'HIGH_LIMIT', '2'    BACnetConstructedDataPositiveIntegerValueHighLimit
            [simple   BACnetApplicationTagUnsignedInteger                               highLimit                       ]
        ]
        [*, 'HIGH_LIMIT', '4'                           BACnetConstructedDataHighLimit
            [simple   BACnetApplicationTagReal                                          highLimit                       ]
        ]
        [*, 'HIGHER_DECK', '12'                         BACnetConstructedDataHigherDeck
            [simple   BACnetApplicationTagObjectIdentifier              higherDeck                                      ]
        ]
        [*, 'IN_PROCESS', '1'                           BACnetConstructedDataInProcess
            [simple   BACnetApplicationTagBoolean                               inProcess                               ]
        ]
        [*, 'IN_PROGRESS', '9'                          BACnetConstructedDataInProgress
            [simple   BACnetLightingInProgressTagged('0', 'TagClass.APPLICATION_TAGS')         inProgress               ]
        ]
        [*, 'INACTIVE_TEXT', '7'                        BACnetConstructedDataInactiveText
            [simple   BACnetApplicationTagCharacterString                               inactiveText                    ]
        ]
        [*, 'INITIAL_TIMEOUT', '2'                      BACnetConstructedDataInitialTimeout
            [simple   BACnetApplicationTagUnsignedInteger                               initialTimeout                  ]
        ]
        [*, 'INPUT_REFERENCE'                           BACnetConstructedDataInputReference
            [simple   BACnetObjectPropertyReference                                     inputReference                  ]
        ]
        [*, 'INSTALLATION_ID', '2'                      BACnetConstructedDataInstallationID
            [simple   BACnetApplicationTagUnsignedInteger                               installationId                  ]
        ]
        [*, 'INSTANCE_OF', '7'                          BACnetConstructedDataInstanceOf
            [simple   BACnetApplicationTagCharacterString                               instanceOf                      ]
        ]
        [*, 'INSTANTANEOUS_POWER', '4'                  BACnetConstructedDataInstantaneousPower
            [simple   BACnetApplicationTagReal                                          instantaneousPower              ]
        ]
        [*, 'INTEGRAL_CONSTANT', '4'                    BACnetConstructedDataIntegralConstant
             [simple   BACnetApplicationTagReal                                integralConstant                         ]
        ]
        [*, 'INTEGRAL_CONSTANT_UNITS', '9'              BACnetConstructedDataIntegralConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        ['ANALOG_INPUT', 'INTERFACE_VALUE'              BACnetConstructedDataAnalogInputInterfaceValue
            [simple   BACnetOptionalREAL                            interfaceValue                                      ]
        ]
        ['ANALOG_OUTPUT', 'INTERFACE_VALUE'             BACnetConstructedDataAnalogOutputInterfaceValue
            [simple   BACnetOptionalREAL                            interfaceValue                                      ]
        ]
        ['BINARY_INPUT', 'INTERFACE_VALUE'              BACnetConstructedDataBinaryInputInterfaceValue
            [simple   BACnetOptionalBinaryPV                        interfaceValue                                      ]
        ]
        ['BINARY_OUTPUT', 'INTERFACE_VALUE'             BACnetConstructedDataBinaryOutputInterfaceValue
            [simple   BACnetOptionalBinaryPV                        interfaceValue                                      ]
        ]
        ['MULTI_STATE_INPUT', 'INTERFACE_VALUE'         BACnetConstructedDataMultiStateInputInterfaceValue
            [simple   BACnetOptionalBinaryPV                        interfaceValue                                      ]
        ]
        ['MULTI_STATE_OUTPUT', 'INTERFACE_VALUE'        BACnetConstructedDataMultiStateOutputInterfaceValue
            [simple   BACnetOptionalBinaryPV                        interfaceValue                                      ]
        ]
        // TODO: unlikely that we have a common type so maybe check that
        //[*, 'INTERFACE_VALUE'                         BACnetConstructedDataInterfaceValue [validation    '1 == 2'    "TODO: implement me INTERFACE_VALUE BACnetConstructedDataInterfaceValue"]]
        [*, 'INTERVAL_OFFSET', '2'                      BACnetConstructedDataIntervalOffset
            [simple   BACnetApplicationTagUnsignedInteger                     intervalOffset                            ]
        ]
        [*, 'IP_ADDRESS', '6'                           BACnetConstructedDataIPAddress
            [simple   BACnetApplicationTagOctetString   ipAddress                                                       ]
        ]
        [*, 'IP_DEFAULT_GATEWAY', '6'                   BACnetConstructedDataIPDefaultGateway
            [simple   BACnetApplicationTagOctetString   ipDefaultGateway                                                ]
        ]
        [*, 'IP_DHCP_ENABLE', '1'                       BACnetConstructedDataIPDHCPEnable
            [simple   BACnetApplicationTagBoolean       ipDhcpEnable                                                    ]
        ]
        [*, 'IP_DHCP_LEASE_TIME', '2'                   BACnetConstructedDataIPDHCPLeaseTime
            [simple BACnetApplicationTagUnsignedInteger                               ipDhcpLeaseTime                   ]
        ]
        [*, 'IP_DHCP_LEASE_TIME_REMAINING', '2'         BACnetConstructedDataIPDHCPLeaseTimeRemaining
            [simple BACnetApplicationTagUnsignedInteger                               ipDhcpLeaseTimeRemaining          ]
        ]
        [*, 'IP_DHCP_SERVER', '6'                       BACnetConstructedDataIPDHCPServer
            [simple   BACnetApplicationTagOctetString   dhcpServer                                                      ]
        ]
        [*, 'IP_DNS_SERVER', '6'                        BACnetConstructedDataIPDNSServer
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagOctetString
                                        ipDnsServer
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'IP_SUBNET_MASK', '6'                       BACnetConstructedDataIPSubnetMask
            [simple   BACnetApplicationTagOctetString   ipSubnetMask                                                    ]
        ]
        [*, 'IPV6_ADDRESS', '6'                         BACnetConstructedDataIPv6Address
            [simple   BACnetApplicationTagOctetString   ipv6Address                                                     ]
        ]
        [*, 'IPV6_AUTO_ADDRESSING_ENABLE', '1'          BACnetConstructedDataIPv6AutoAddressingEnable
            [simple   BACnetApplicationTagBoolean                                       autoAddressingEnable            ]
        ]
        [*, 'IPV6_DEFAULT_GATEWAY', '6'                 BACnetConstructedDataIPv6DefaultGateway
            [simple   BACnetApplicationTagOctetString   ipv6DefaultGateway                                              ]
        ]
        [*, 'IPV6_DHCP_LEASE_TIME', '2'                 BACnetConstructedDataIPv6DHCPLeaseTime
            [simple   BACnetApplicationTagUnsignedInteger                               ipv6DhcpLeaseTime               ]
        ]
        [*, 'IPV6_DHCP_LEASE_TIME_REMAINING', '2'       BACnetConstructedDataIPv6DHCPLeaseTimeRemaining
            [simple   BACnetApplicationTagUnsignedInteger                               ipv6DhcpLeaseTimeRemaining      ]
        ]
        [*, 'IPV6_DHCP_SERVER', '6'                     BACnetConstructedDataIPv6DHCPServer
            [simple   BACnetApplicationTagOctetString   dhcpServer                                                      ]
        ]
        [*, 'IPV6_DNS_SERVER', '6'                      BACnetConstructedDataIPv6DNSServer
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagOctetString
                                        ipv6DnsServer
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'IPV6_PREFIX_LENGTH', '2'                   BACnetConstructedDataIPv6PrefixLength
            [simple   BACnetApplicationTagUnsignedInteger                               ipv6PrefixLength                ]
        ]
        [*, 'IPV6_ZONE_INDEX', '7'                      BACnetConstructedDataIPv6ZoneIndex
            [simple   BACnetApplicationTagCharacterString ipv6ZoneIndex                                                 ]
        ]
        [*, 'IS_UTC', '1'                               BACnetConstructedDataIsUTC
            [simple   BACnetApplicationTagBoolean                                       isUtc                           ]
        ]
        [*, 'KEY_SETS'                                  BACnetConstructedDataKeySets
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetSecurityKeySet
                                keySets
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
            [validation 'arrayIndexArgument!=null || COUNT(keySets) == 2' "keySets should have exactly 2 values"        ]
        ]
        [*, 'LANDING_CALL_CONTROL'                      BACnetConstructedDataLandingCallControl
            [simple   BACnetLandingCallStatus             landingCallControl                                            ]
        ]
        [*, 'LANDING_CALLS'                             BACnetConstructedDataLandingCalls
            [array    BACnetLandingCallStatus
                                landingCallStatus
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LANDING_DOOR_STATUS'                       BACnetConstructedDataLandingDoorStatus
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLandingDoorStatus
                                landingDoorStatus
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LAST_ACCESS_EVENT', '9'                    BACnetConstructedDataLastAccessEvent
            [simple   BACnetAccessEventTagged('0', 'TagClass.APPLICATION_TAGS')   lastAccessEvent                       ]
        ]
        [*, 'LAST_ACCESS_POINT'                         BACnetConstructedDataLastAccessPoint
            [simple   BACnetDeviceObjectReference       lastAccessPoint                                                 ]
        ]
        [*, 'LAST_COMMAND_TIME'                         BACnetConstructedDataLastCommandTime
            [simple   BACnetTimeStamp                                                lastCommandTime                    ]
        ]
        [*, 'LAST_CREDENTIAL_ADDED'                     BACnetConstructedDataLastCredentialAdded
            [simple   BACnetDeviceObjectReference       lastCredentialAdded                                             ]
        ]
        [*, 'LAST_CREDENTIAL_ADDED_TIME'                BACnetConstructedDataLastCredentialAddedTime
            [simple   BACnetDateTime    lastCredentialAddedTime                                                         ]
        ]
        [*, 'LAST_CREDENTIAL_REMOVED'                   BACnetConstructedDataLastCredentialRemoved
            [simple   BACnetDeviceObjectReference       lastCredentialRemoved                                           ]
        ]
        [*, 'LAST_CREDENTIAL_REMOVED_TIME'              BACnetConstructedDataLastCredentialRemovedTime
            [simple   BACnetDateTime    lastCredentialRemovedTime                                                       ]
        ]
        [*, 'LAST_KEY_SERVER'                           BACnetConstructedDataLastKeyServer
            [simple   BACnetAddressBinding              lastKeyServer                                                   ]
        ]
        [*, 'LAST_NOTIFY_RECORD', '2'                   BACnetConstructedDataLastNotifyRecord
            [simple   BACnetApplicationTagUnsignedInteger                               lastNotifyRecord                ]
        ]
        [*, 'LAST_PRIORITY', '2'                        BACnetConstructedDataLastPriority
            [simple   BACnetApplicationTagUnsignedInteger                               lastPriority                    ]
        ]
        [*, 'LAST_RESTART_REASON', '9'                  BACnetConstructedDataLastRestartReason
            [simple   BACnetRestartReasonTagged('0', 'TagClass.APPLICATION_TAGS')   lastRestartReason                   ]
        ]
        [*, 'LAST_RESTORE_TIME'                         BACnetConstructedDataLastRestoreTime
            [simple   BACnetTimeStamp                                                     lastRestoreTime               ]
        ]
        [*, 'LAST_STATE_CHANGE', '9'                    BACnetConstructedDataLastStateChange
            [simple   BACnetTimerTransitionTagged('0', 'TagClass.APPLICATION_TAGS')             lastStateChange         ]
        ]
        [*, 'LAST_USE_TIME'                             BACnetConstructedDataLastUseTime
            [simple   BACnetDateTime                                        lastUseTime                                 ]
        ]
        [*, 'LIFE_SAFETY_ALARM_VALUES', '9'             BACnetConstructedDataLifeSafetyAlarmValues
            [array    BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')
                            alarmValues              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LIGHTING_COMMAND'                          BACnetConstructedDataLightingCommand
            [simple   BACnetLightingCommand     lightingCommand                                                         ]
        ]
        [*, 'LIGHTING_COMMAND_DEFAULT_PRIORITY', '2'    BACnetConstructedDataLightingCommandDefaultPriority
            [simple BACnetApplicationTagUnsignedInteger                       lightingCommandDefaultPriority            ]
        ]
        [*, 'LIMIT_ENABLE','8'                          BACnetConstructedDataLimitEnable
            [simple   BACnetLimitEnableTagged('0', 'TagClass.APPLICATION_TAGS')         limitEnable                     ]
        ]
        [*, 'LIMIT_MONITORING_INTERVAL', '2'            BACnetConstructedDataLimitMonitoringInterval
            [simple BACnetApplicationTagUnsignedInteger                       limitMonitoringInterval                   ]
        ]
        [*, 'LINK_SPEED', '4'                           BACnetConstructedDataLinkSpeed
            [simple   BACnetApplicationTagReal                                          linkSpeed                       ]
        ]
        [*, 'LINK_SPEED_AUTONEGOTIATE', '1'             BACnetConstructedDataLinkSpeedAutonegotiate
            [simple   BACnetApplicationTagBoolean                                       linkSpeedAutonegotiate          ]
        ]
        [*, 'LINK_SPEEDS', '4'                          BACnetConstructedDataLinkSpeeds
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagReal
                                        linkSpeeds
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'LIST_OF_GROUP_MEMBERS'                     BACnetConstructedDataListOfGroupMembers
            [array    BACnetReadAccessSpecification
                                        listOfGroupMembers
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        ['CHANNEL', 'LIST_OF_OBJECT_PROPERTY_REFERENCES'    BACnetConstructedDataChannelListOfObjectPropertyReferences
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectPropertyReference
                            references              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LIST_OF_OBJECT_PROPERTY_REFERENCES'        BACnetConstructedDataListOfObjectPropertyReferences
            [array    BACnetDeviceObjectPropertyReference
                            references              terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'LOCAL_DATE', '10'                          BACnetConstructedDataLocalDate
            [simple   BACnetApplicationTagDate          localDate                                                       ]
        ]
        [*, 'LOCAL_FORWARDING_ONLY', '1'                BACnetConstructedDataLocalForwardingOnly
            [simple BACnetApplicationTagBoolean                               localForwardingOnly                       ]
        ]
        [*, 'LOCAL_TIME', '11'                          BACnetConstructedDataLocalTime
            [simple BACnetApplicationTagTime                                  localTime                                 ]
        ]
        [*, 'LOCATION', '7'                             BACnetConstructedDataLocation
            [simple BACnetApplicationTagCharacterString     location                                                    ]
        ]
        [*, 'LOCK_STATUS', '9'                          BACnetConstructedDataLockStatus
            [simple   BACnetLockStatusTagged('0', 'TagClass.APPLICATION_TAGS')      lockStatus                          ]
        ]
        [*, 'LOCKOUT', '1'                              BACnetConstructedDataLockout
            [simple BACnetApplicationTagBoolean                               lockout                                   ]
        ]
        [*, 'LOCKOUT_RELINQUISH_TIME', '2'              BACnetConstructedDataLockoutRelinquishTime
            [simple BACnetApplicationTagUnsignedInteger                       lockoutRelinquishTime                     ]
        ]
        ['EVENT_LOG', 'LOG_BUFFER'                      BACnetConstructedDataEventLogLogBuffer
            [array    BACnetEventLogRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['TREND_LOG', 'LOG_BUFFER'                      BACnetConstructedDataTrendLogLogBuffer
            [array    BACnetLogRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['TREND_LOG_MULTIPLE', 'LOG_BUFFER'             BACnetConstructedDataTrendLogMultipleLogBuffer
            [array    BACnetLogMultipleRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'LOG_BUFFER'                                BACnetConstructedDataLogBuffer
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLogRecord
                            floorText
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['TREND_LOG', 'LOG_DEVICE_OBJECT_PROPERTY'      BACnetConstructedDataTrendLogLogDeviceObjectProperty
            [simple   BACnetDeviceObjectPropertyReference   logDeviceObjectProperty                                     ]
        ]
        ['TREND_LOG_MULTIPLE', 'LOG_DEVICE_OBJECT_PROPERTY' BACnetConstructedDataTrendLogMultipleLogDeviceObjectProperty
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectPropertyReference
                            groupMembers
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'LOG_DEVICE_OBJECT_PROPERTY'                BACnetConstructedDataLogDeviceObjectProperty
            [simple   BACnetDeviceObjectPropertyReference   logDeviceObjectProperty                                     ]
        ]
        [*, 'LOG_INTERVAL', '2'                         BACnetConstructedDataLogInterval
            [simple   BACnetApplicationTagUnsignedInteger                               logInterval                     ]
        ]
        [*, 'LOGGING_OBJECT', '12'                      BACnetConstructedDataLoggingObject
            [simple   BACnetApplicationTagObjectIdentifier                              loggingObject                   ]
        ]
        [*, 'LOGGING_RECORD'                            BACnetConstructedDataLoggingRecord
            [simple   BACnetAccumulatorRecord                                           loggingRecord                   ]
        ]
        [*, 'LOGGING_TYPE', '9'                         BACnetConstructedDataLoggingType
            [simple   BACnetLoggingTypeTagged('0', 'TagClass.APPLICATION_TAGS')         loggingType                     ]
        ]
        [*, 'LOW_DIFF_LIMIT'                            BACnetConstructedDataLowDiffLimit
            [simple   BACnetOptionalREAL                            lowDiffLimit                                        ]
        ]
        ['ACCUMULATOR', 'LOW_LIMIT', '2'                BACnetConstructedDataAccumulatorLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                               lowLimit                        ]
        ]
        ['LARGE_ANALOG_VALUE', 'LOW_LIMIT', '5'         BACnetConstructedDataLargeAnalogValueLowLimit
            [simple   BACnetApplicationTagDouble                                        lowLimit                        ]
        ]
        ['INTEGER_VALUE', 'LOW_LIMIT', '3'              BACnetConstructedDataIntegerValueLowLimit
            [simple   BACnetApplicationTagSignedInteger                                 lowLimit                        ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'LOW_LIMIT', '2'     BACnetConstructedDataPositiveIntegerValueLowLimit
            [simple   BACnetApplicationTagUnsignedInteger                                 lowLimit                      ]
        ]
        [*, 'LOW_LIMIT', '4'                            BACnetConstructedDataLowLimit
            [simple   BACnetApplicationTagReal                                          lowLimit                        ]
        ]
        [*, 'LOWER_DECK', '12'                          BACnetConstructedDataLowerDeck
            [simple   BACnetApplicationTagObjectIdentifier              lowerDeck                                       ]
        ]
        [*, 'MAC_ADDRESS', '6'                          BACnetConstructedDataMACAddress
            [simple   BACnetApplicationTagOctetString   macAddress                                                      ]
        ]
        [*, 'MACHINE_ROOM_ID', '12'                     BACnetConstructedDataMachineRoomID
            [simple   BACnetApplicationTagObjectIdentifier          machineRoomId                                       ]
        ]
        ['LIFE_SAFETY_ZONE', 'MAINTENANCE_REQUIRED', '1' BACnetConstructedDataLifeSafetyZoneMaintenanceRequired
            [simple   BACnetApplicationTagBoolean                                   maintenanceRequired                 ]
        ]
        [*, 'MAINTENANCE_REQUIRED', '9'                 BACnetConstructedDataMaintenanceRequired
            [simple   BACnetMaintenanceTagged('0', 'TagClass.APPLICATION_TAGS')     maintenanceRequired                 ]
        ]
        [*, 'MAKING_CAR_CALL'                           BACnetConstructedDataMakingCarCall
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                        makingCarCall
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'MANIPULATED_VARIABLE_REFERENCE'            BACnetConstructedDataManipulatedVariableReference
            [simple   BACnetObjectPropertyReference                                     manipulatedVariableReference    ]
        ]
        [*, 'MANUAL_SLAVE_ADDRESS_BINDING'              BACnetConstructedDataManualSlaveAddressBinding
            [array    BACnetAddressBinding
                                manualSlaveAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MASKED_ALARM_VALUES', '9'                  BACnetConstructedDataMaskedAlarmValues
            [array    BACnetDoorAlarmStateTagged('0', 'TagClass.APPLICATION_TAGS')
                                maskedAlarmValues
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MAX_ACTUAL_VALUE', '4'                     BACnetConstructedDataMaxActualValue
            [simple   BACnetApplicationTagReal                                          maxActualValue                  ]
        ]
        [*, 'MAX_APDU_LENGTH_ACCEPTED', '2'             BACnetConstructedDataMaxAPDULengthAccepted
            [simple BACnetApplicationTagUnsignedInteger                     maxApduLengthAccepted                       ]
        ]
        [*, 'MAX_FAILED_ATTEMPTS', '2'                  BACnetConstructedDataMaxFailedAttempts
            [simple BACnetApplicationTagUnsignedInteger                     maxFailedAttempts                           ]
        ]
        ['DEVICE', 'MAX_INFO_FRAMES', '2'               BACnetConstructedDataDeviceMaxInfoFrames
            [simple BACnetApplicationTagUnsignedInteger                     maxInfoFrames                               ]
        ]
        ['NETWORK_PORT', 'MAX_INFO_FRAMES', '2'         BACnetConstructedDataNetworkPortMaxInfoFrames
            [simple BACnetApplicationTagUnsignedInteger                     maxInfoFrames                               ]
        ]
        [*, 'MAX_INFO_FRAMES', '2'                      BACnetConstructedDataMaxInfoFrames
            [simple BACnetApplicationTagUnsignedInteger                     maxInfoFrames                               ]
        ]
        ['DEVICE', 'MAX_MASTER', '2'                    BACnetConstructedDataDeviceMaxMaster
            [simple BACnetApplicationTagUnsignedInteger                     maxMaster                                   ]
        ]
        ['NETWORK_PORT', 'MAX_MASTER', '2'              BACnetConstructedDataNetworkPortMaxMaster
            [simple BACnetApplicationTagUnsignedInteger                     maxMaster                                   ]
        ]
        [*, 'MAX_MASTER', '2'                           BACnetConstructedDataMaxMaster
            [simple BACnetApplicationTagUnsignedInteger                     maxMaster                                   ]
        ]
        ['ACCUMULATOR', 'MAX_PRES_VALUE', '2'           BACnetConstructedDataAccumulatorMaxPresValue
            [simple BACnetApplicationTagUnsignedInteger                     maxPresValue                                ]
        ]
        ['ANALOG_INPUT', 'MAX_PRES_VALUE', '4'          BACnetConstructedDataAnalogInputMaxPresValue
            [simple BACnetApplicationTagReal                                maxPresValue                                ]
        ]
        ['ANALOG_OUTPUT', 'MAX_PRES_VALUE', '4'         BACnetConstructedDataAnalogOutputMaxPresValue
            [simple BACnetApplicationTagReal                                maxPresValue                                ]
        ]
        ['ANALOG_VALUE', 'MAX_PRES_VALUE', '4'          BACnetConstructedDataAnalogValueMaxPresValue
            [simple BACnetApplicationTagReal                                maxPresValue                                ]
        ]
        ['LARGE_ANALOG_VALUE', 'MAX_PRES_VALUE', '5'    BACnetConstructedDataLargeAnalogValueMaxPresValue
            [simple BACnetApplicationTagDouble                              maxPresValue                                ]
        ]
        ['INTEGER_VALUE', 'MAX_PRES_VALUE', '3'         BACnetConstructedDataIntegerValueMaxPresValue
            [simple BACnetApplicationTagSignedInteger                       maxPresValue                                ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'MAX_PRES_VALUE', '2' BACnetConstructedDataPositiveIntegerValueMaxPresValue
            [simple BACnetApplicationTagUnsignedInteger                     maxPresValue                                ]
        ]
        ['TIMER', 'MAX_PRES_VALUE', '2'                 BACnetConstructedDataTimerMaxPresValue
            [simple BACnetApplicationTagUnsignedInteger                     maxPresValue                                ]
        ]
        [*, 'MAX_PRES_VALUE', '4'                       BACnetConstructedDataMaxPresValue
            [simple BACnetApplicationTagReal                                maxPresValue                                ]
        ]
        [*, 'MAX_SEGMENTS_ACCEPTED', '2'                BACnetConstructedDataMaxSegmentsAccepted
            [simple BACnetApplicationTagUnsignedInteger                               maxSegmentsAccepted               ]
        ]
        [*, 'MAXIMUM_OUTPUT', '4'                       BACnetConstructedDataMaximumOutput
            [simple BACnetApplicationTagReal                                maximumOutput                               ]
        ]
        [*, 'MAXIMUM_VALUE', '4'                        BACnetConstructedDataMaximumValue
            [simple BACnetApplicationTagReal                                maximumValue                                ]
        ]
        [*, 'MAXIMUM_VALUE_TIMESTAMP'                   BACnetConstructedDataMaximumValueTimestamp
            [simple   BACnetDateTime                                        maximumValueTimestamp                       ]
        ]
        [*, 'MEMBER_OF' BACnetConstructedDataMemberOf
            [array    BACnetDeviceObjectReference
                    zones
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'MEMBER_STATUS_FLAGS', '8'                  BACnetConstructedDataMemberStatusFlags
            [simple   BACnetStatusFlagsTagged('0', 'TagClass.APPLICATION_TAGS')         statusFlags                     ]
        ]
        [*, 'MEMBERS'                                   BACnetConstructedDataMembers
            [array    BACnetDeviceObjectReference
                                    members
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'MIN_ACTUAL_VALUE', '4'                     BACnetConstructedDataMinActualValue
            [simple   BACnetApplicationTagReal                                          minActualValue                  ]
        ]
        ['ACCUMULATOR', 'MIN_PRES_VALUE', '2'           BACnetConstructedDataAccumulatorMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                     minPresValue                              ]
        ]
        ['INTEGER_VALUE', 'MIN_PRES_VALUE', '3'         BACnetConstructedDataIntegerValueMinPresValue
            [simple   BACnetApplicationTagSignedInteger                       minPresValue                              ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'MIN_PRES_VALUE', '2'    BACnetConstructedDataPositiveIntegerValueMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                     minPresValue                              ]
        ]
        ['LARGE_ANALOG_VALUE', 'MIN_PRES_VALUE', '5'    BACnetConstructedDataLargeAnalogValueMinPresValue
            [simple   BACnetApplicationTagDouble                              minPresValue                              ]
        ]
        ['TIMER', 'MIN_PRES_VALUE', '2'                 BACnetConstructedDataTimerMinPresValue
            [simple   BACnetApplicationTagUnsignedInteger                     minPresValue                              ]
        ]
        [*, 'MIN_PRES_VALUE', '4'                       BACnetConstructedDataMinPresValue
            [simple   BACnetApplicationTagReal                                minPresValue                              ]
        ]
        [*, 'MINIMUM_OFF_TIME', '2'                     BACnetConstructedDataMinimumOffTime
            [simple   BACnetApplicationTagUnsignedInteger                     minimumOffTime                            ]
        ]
        [*, 'MINIMUM_ON_TIME', '2'                      BACnetConstructedDataMinimumOnTime
            [simple   BACnetApplicationTagUnsignedInteger                     minimumOnTime                            ]
        ]
        [*, 'MINIMUM_OUTPUT', '4'                       BACnetConstructedDataMinimumOutput
            [simple   BACnetApplicationTagReal                                minimumOutput                             ]
        ]
        [*, 'MINIMUM_VALUE', '4'                        BACnetConstructedDataMinimumValue
            [simple   BACnetApplicationTagReal                              minimumValue                                ]
        ]
        [*, 'MINIMUM_VALUE_TIMESTAMP'                   BACnetConstructedDataMinimumValueTimestamp
            [simple   BACnetDateTime                                        minimumValueTimestamp                       ]
        ]
        [*, 'MODE', '9'                                 BACnetConstructedDataMode
            [simple   BACnetLifeSafetyModeTagged('0', 'TagClass.APPLICATION_TAGS')              mode                    ]
        ]
        [*, 'MODEL_NAME', '7'                           BACnetConstructedDataModelName
           [simple   BACnetApplicationTagCharacterString               modelName                                        ]
        ]
        [*, 'MODIFICATION_DATE'                         BACnetConstructedDataModificationDate
            [simple   BACnetDateTime                                        modificationDate                            ]
        ]
        [*, 'MUSTER_POINT', '1'                         BACnetConstructedDataMusterPoint
            [simple   BACnetApplicationTagBoolean                                       musterPoint                     ]
        ]
        [*, 'NEGATIVE_ACCESS_RULES'                     BACnetConstructedDataNegativeAccessRules
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAccessRule
                            negativeAccessRules
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'NETWORK_ACCESS_SECURITY_POLICIES'          BACnetConstructedDataNetworkAccessSecurityPolicies
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNetworkSecurityPolicy
                            networkAccessSecurityPolicies
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'NETWORK_INTERFACE_NAME', '7'               BACnetConstructedDataNetworkInterfaceName
            [simple BACnetApplicationTagCharacterString                  networkInterfaceName                           ]
        ]
        [*, 'NETWORK_NUMBER', '2'                       BACnetConstructedDataNetworkNumber
            [simple BACnetApplicationTagUnsignedInteger                               networkNumber                     ]
        ]
        [*, 'NETWORK_NUMBER_QUALITY', '9'               BACnetConstructedDataNetworkNumberQuality
            [simple   BACnetNetworkNumberQualityTagged('0', 'TagClass.APPLICATION_TAGS')    networkNumberQuality        ]
        ]
        [*, 'NETWORK_TYPE', '9'                         BACnetConstructedDataNetworkType
            [simple   BACnetNetworkTypeTagged('0', 'TagClass.APPLICATION_TAGS')    networkType                          ]
        ]
        [*, 'NEXT_STOPPING_FLOOR', '2'                  BACnetConstructedDataNextStoppingFloor
            [simple   BACnetApplicationTagUnsignedInteger                     nextStoppingFloor                         ]
        ]
        [*, 'NODE_SUBTYPE', '7'                         BACnetConstructedDataNodeSubtype
            [simple BACnetApplicationTagCharacterString                  nodeSubType                                    ]
        ]
        [*, 'NODE_TYPE', '9'                            BACnetConstructedDataNodeType
            [simple BACnetNodeTypeTagged('0', 'TagClass.APPLICATION_TAGS')  nodeType                                    ]
        ]
        [*, 'NOTIFICATION_CLASS', '2'                   BACnetConstructedDataNotificationClass
            [simple BACnetApplicationTagUnsignedInteger                               notificationClass                 ]
        ]
        [*, 'NOTIFICATION_THRESHOLD', '2'               BACnetConstructedDataNotificationThreshold
            [simple BACnetApplicationTagUnsignedInteger                               notificationThreshold             ]
        ]
        [*, 'NOTIFY_TYPE', '9'                          BACnetConstructedDataNotifyType
            [simple BACnetNotifyTypeTagged('0', 'TagClass.APPLICATION_TAGS')          notifyType                        ]
        ]
        [*, 'NUMBER_OF_APDU_RETRIES', '2'               BACnetConstructedDataNumberOfAPDURetries
            [simple   BACnetApplicationTagUnsignedInteger                               numberOfApduRetries             ]
        ]
        [*, 'NUMBER_OF_AUTHENTICATION_POLICIES', '2'    BACnetConstructedDataNumberOfAuthenticationPolicies
            [simple   BACnetApplicationTagUnsignedInteger                               numberOfAuthenticationPolicies  ]
        ]
        [*, 'NUMBER_OF_STATES', '2'                     BACnetConstructedDataNumberOfStates
            [simple BACnetApplicationTagUnsignedInteger                                 numberOfState                   ]
        ]
        [*, 'OBJECT_IDENTIFIER', '12'                   BACnetConstructedDataObjectIdentifier
            [simple   BACnetApplicationTagObjectIdentifier                              objectIdentifier                ]
        ]
        [*, 'OBJECT_LIST'                               BACnetConstructedDataObjectList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                                objectList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'OBJECT_NAME', '7'                          BACnetConstructedDataObjectName
            [simple   BACnetApplicationTagCharacterString                               objectName                      ]
        ]
        [*, 'OBJECT_PROPERTY_REFERENCE'                 BACnetConstructedDataObjectPropertyReference
            [simple   BACnetDeviceObjectPropertyReference                               propertyReference               ]
        ]
        [*, 'OBJECT_TYPE', '9'                          BACnetConstructedDataObjectType
            [simple   BACnetObjectTypeTagged('0', 'TagClass.APPLICATION_TAGS')          objectType                      ]
        ]
        [*, 'OCCUPANCY_COUNT', '2'                      BACnetConstructedDataOccupancyCount
            [simple BACnetApplicationTagUnsignedInteger                       occupancyCount                            ]
        ]
        [*, 'OCCUPANCY_COUNT_ADJUST', '1'               BACnetConstructedDataOccupancyCountAdjust
            [simple BACnetApplicationTagBoolean                               occupancyCountAdjust                      ]
        ]
        [*, 'OCCUPANCY_COUNT_ENABLE', '1'               BACnetConstructedDataOccupancyCountEnable
            [simple BACnetApplicationTagBoolean                               occupancyCountEnable                      ]
        ]
        [*, 'OCCUPANCY_LOWER_LIMIT', '2'                BACnetConstructedDataOccupancyLowerLimit
            [simple BACnetApplicationTagUnsignedInteger                       occupancyLowerLimit                       ]
        ]
        [*, 'OCCUPANCY_LOWER_LIMIT_ENFORCED', '1'       BACnetConstructedDataOccupancyLowerLimitEnforced
            [simple BACnetApplicationTagBoolean                               occupancyLowerLimitEnforced               ]
        ]
        [*, 'OCCUPANCY_STATE', '9'                      BACnetConstructedDataOccupancyState
            [simple BACnetAccessZoneOccupancyStateTagged('0', 'TagClass.APPLICATION_TAGS') occupancyState               ]
        ]
        [*, 'OCCUPANCY_UPPER_LIMIT', '2'                BACnetConstructedDataOccupancyUpperLimit
            [simple BACnetApplicationTagUnsignedInteger                       occupancyUpperLimit                       ]
        ]
        [*, 'OCCUPANCY_UPPER_LIMIT_ENFORCED', '1'       BACnetConstructedDataOccupancyUpperLimitEnforced
            [simple BACnetApplicationTagBoolean                               occupancyUpperLimitEnforced               ]
        ]
        [*, 'OPERATION_DIRECTION', '9'                  BACnetConstructedDataOperationDirection
            [simple BACnetEscalatorOperationDirectionTagged('0', 'TagClass.APPLICATION_TAGS')   operationDirection      ]
        ]
        [*, 'OPERATION_EXPECTED', '9'                   BACnetConstructedDataOperationExpected
            [simple   BACnetLifeSafetyOperationTagged('0', 'TagClass.APPLICATION_TAGS')         lifeSafetyOperations    ]
        ]
        [*, 'OPTIONAL'                                  BACnetConstructedDataOptional
            [validation    '1 == 2'    "An property identified by OPTIONAL should never occur in the wild"]
        ]
        [*, 'OUT_OF_SERVICE', '1'                       BACnetConstructedDataOutOfService
            [simple   BACnetApplicationTagBoolean                                                       outOfService    ]
        ]
        [*, 'OUTPUT_UNITS', '9'                         BACnetConstructedDataOutputUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'PACKET_REORDER_TIME', '2'                  BACnetConstructedDataPacketReorderTime
            [simple BACnetApplicationTagUnsignedInteger                       packetReorderTime                         ]
        ]
        [*, 'PASSBACK_MODE'                             BACnetConstructedDataPassbackMode
            [simple BACnetAccessPassbackModeTagged('0', 'TagClass.APPLICATION_TAGS')             passbackMode           ]
        ]
        [*, 'PASSBACK_TIMEOUT', '2'                     BACnetConstructedDataPassbackTimeout
            [simple BACnetApplicationTagUnsignedInteger                       passbackTimeout                           ]
        ]
        [*, 'PASSENGER_ALARM', '1'                      BACnetConstructedDataPassengerAlarm
            [simple   BACnetApplicationTagBoolean                             passengerAlarm                            ]
        ]
        [*, 'POLARITY', '9'                             BACnetConstructedDataPolarity
            [simple   BACnetPolarityTagged('0', 'TagClass.APPLICATION_TAGS')  polarity                                  ]
        ]
        [*, 'PORT_FILTER'                               BACnetConstructedDataPortFilter
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetPortPermission
                            portFilter
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'POSITIVE_ACCESS_RULES'                     BACnetConstructedDataPositiveAccessRules
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAccessRule
                            positiveAccessRules
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'POWER', '4'                                BACnetConstructedDataPower
            [simple BACnetApplicationTagReal                                power                                       ]
        ]
        [*, 'POWER_MODE', '1'                           BACnetConstructedDataPowerMode
            [simple   BACnetApplicationTagBoolean               powerMode                                               ]
        ]
        [*, 'PRESCALE'                                  BACnetConstructedDataPrescale
            [simple   BACnetPrescale                            prescale                                                ]
        ]
        ['ACCESS_DOOR', 'PRESENT_VALUE', '9'            BACnetConstructedDataAccessDoorPresentValue
            [simple   BACnetDoorValueTagged('0', 'TagClass.APPLICATION_TAGS')           presentValue                    ]
        ]
        ['ALERT_ENROLLMENT', 'PRESENT_VALUE', '12'      BACnetConstructedDataAlertEnrollmentPresentValue
            [simple   BACnetApplicationTagObjectIdentifier                              presentValue                    ]
        ]
        ['ANALOG_INPUT', 'PRESENT_VALUE', '4'           BACnetConstructedDataAnalogInputPresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        ['ANALOG_OUTPUT', 'PRESENT_VALUE', '4'          BACnetConstructedDataAnalogOutputPresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        ['ANALOG_VALUE', 'PRESENT_VALUE', '4'           BACnetConstructedDataAnalogValuePresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        ['BINARY_INPUT', 'PRESENT_VALUE', '9'           BACnetConstructedDataBinaryInputPresentValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')            presentValue                    ]
        ]
        ['BINARY_OUTPUT', 'PRESENT_VALUE', '9'          BACnetConstructedDataBinaryOutputPresentValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')            presentValue                    ]
        ]
        ['BINARY_VALUE', 'PRESENT_VALUE', '9'           BACnetConstructedDataBinaryValuePresentValue
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')            presentValue                    ]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'PRESENT_VALUE', '9' BACnetConstructedDataBinaryLightingOutputPresentValue
            [simple   BACnetBinaryLightingPVTagged('0', 'TagClass.APPLICATION_TAGS')    presentValue                    ]
        ]
        ['BITSTRING_VALUE', 'PRESENT_VALUE', '8'        BACnetConstructedDataBitStringValuePresentValue
            [simple   BACnetApplicationTagBitString                                     presentValue                    ]
        ]
        ['CALENDAR', 'PRESENT_VALUE', '1'               BACnetConstructedDataCalendarPresentValue
            [simple   BACnetApplicationTagBoolean                                       presentValue                    ]
        ]
        ['CHANNEL', 'PRESENT_VALUE'                     BACnetConstructedDataChannelPresentValue
            [simple   BACnetChannelValue                                                presentValue                    ]
        ]
        ['CHARACTERSTRING_VALUE', 'PRESENT_VALUE','7'   BACnetConstructedDataCharacterStringValuePresentValue
            [simple   BACnetApplicationTagCharacterString                               presentValue                    ]
        ]
        ['CREDENTIAL_DATA_INPUT', 'PRESENT_VALUE'       BACnetConstructedDataCredentialDataInputPresentValue
            [simple   BACnetAuthenticationFactor                                        presentValue                    ]
        ]
        ['DATE_VALUE', 'PRESENT_VALUE', '10'            BACnetConstructedDataDateValuePresentValue
            [simple   BACnetApplicationTagDate                                          presentValue                    ]
        ]
        ['DATEPATTERN_VALUE', 'PRESENT_VALUE', '10'     BACnetConstructedDataDatePatternValuePresentValue
            [simple   BACnetApplicationTagDate                                          presentValue                    ]
        ]
        ['DATETIME_VALUE', 'PRESENT_VALUE', '11'        BACnetConstructedDataDateTimeValuePresentValue
            [simple   BACnetDateTime                                                    presentValue                    ]
        ]
        ['DATETIMEPATTERN_VALUE', 'PRESENT_VALUE', '11' BACnetConstructedDataDateTimePatternValuePresentValue
            [simple   BACnetDateTime                                                    presentValue                    ]
        ]
        ['INTEGER_VALUE', 'PRESENT_VALUE', '3'          BACnetConstructedDataIntegerValuePresentValue
            [simple   BACnetApplicationTagSignedInteger                                 presentValue                    ]
        ]
        ['LARGE_ANALOG_VALUE', 'PRESENT_VALUE', '5'     BACnetConstructedDataLargeAnalogValuePresentValue
            [simple   BACnetApplicationTagDouble                                        presentValue                    ]
        ]
        ['LIGHTING_OUTPUT', 'PRESENT_VALUE', '4'        BACnetConstructedDataLightingOutputPresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        ['LIFE_SAFETY_POINT', 'PRESENT_VALUE', '9'      BACnetConstructedDataLifeSafetyPointPresentValue
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')     presentValue                    ]
        ]
        ['LIFE_SAFETY_ZONE', 'PRESENT_VALUE', '9'       BACnetConstructedDataLifeSafetyZonePresentValue
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS')     presentValue                    ]
        ]
        ['LOAD_CONTROL', 'PRESENT_VALUE', '9'           BACnetConstructedDataLoadControlPresentValue
            [simple   BACnetShedStateTagged('0', 'TagClass.APPLICATION_TAGS')           presentValue                    ]
        ]
        ['LOOP', 'PRESENT_VALUE', '4'                   BACnetConstructedDataLoopPresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        ['PULSE_CONVERTER', 'PRESENT_VALUE', '4'        BACnetConstructedDataPulseConverterPresentValue
            [simple   BACnetApplicationTagReal                                          presentValue                    ]
        ]
        ['GROUP', 'PRESENT_VALUE'                       BACnetConstructedDataGroupPresentValue
            [array    BACnetReadAccessResult
                          presentValue
                                  terminated
                                  'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['GLOBAL_GROUP', 'PRESENT_VALUE'                BACnetConstructedDataGlobalGroupPresentValue
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetPropertyAccessResult
                          presentValue
                                  terminated
                                  'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['OCTETSTRING_VALUE', 'PRESENT_VALUE', '6'      BACnetConstructedDataOctetStringValuePresentValue
            [simple   BACnetApplicationTagOctetString                                   presentValue                    ]
        ]
        ['SCHEDULE', 'PRESENT_VALUE'                    BACnetConstructedDataSchedulePresentValue
            [simple   BACnetConstructedDataElement('BACnetObjectType.VENDOR_PROPRIETARY_VALUE', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
                                                                                        presentValue                    ]
        ]
        ['TIME_VALUE', 'PRESENT_VALUE', '11'            BACnetConstructedDataTimeValuePresentValue
            [simple   BACnetApplicationTagTime                                          presentValue                    ]
        ]
        ['TIMEPATTERN_VALUE', 'PRESENT_VALUE', '11'     BACnetConstructedDataTimePatternValuePresentValue
            [simple   BACnetApplicationTagTime                                          presentValue                    ]
        ]
        [*, 'PRESENT_VALUE', '2'                        BACnetConstructedDataPresentValue
            [simple   BACnetApplicationTagUnsignedInteger                     presentValue                              ]
        ]
        [*, 'PRIORITY'                                  BACnetConstructedDataPriority
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
        [*, 'PRIORITY_ARRAY'                            BACnetConstructedDataPriorityArray
            [simple   BACnetPriorityArray('objectTypeArgument', 'tagNumber', 'arrayIndexArgument')   priorityArray      ]
        ]
        [*, 'PRIORITY_FOR_WRITING', '2'                 BACnetConstructedDataPriorityForWriting
            [simple   BACnetApplicationTagUnsignedInteger                     priorityForWriting                        ]
        ]
        [*, 'PROCESS_IDENTIFIER', '2'                   BACnetConstructedDataProcessIdentifier
            [simple   BACnetApplicationTagUnsignedInteger                               processIdentifier               ]
        ]
        [*, 'PROCESS_IDENTIFIER_FILTER'                 BACnetConstructedDataProcessIdentifierFilter
            [simple   BACnetProcessIdSelection                          processIdentifierFilter                         ]
        ]
        [*, 'PROFILE_LOCATION', '7'                     BACnetConstructedDataProfileLocation
            [simple   BACnetApplicationTagCharacterString                               profileLocation                 ]
        ]
        [*, 'PROFILE_NAME', '7'                         BACnetConstructedDataProfileName
            [simple   BACnetApplicationTagCharacterString                               profileName                     ]
        ]
        [*, 'PROGRAM_CHANGE', '9'                       BACnetConstructedDataProgramChange
            [simple   BACnetProgramRequestTagged('0', 'TagClass.APPLICATION_TAGS')           programChange              ]
        ]
        [*, 'PROGRAM_LOCATION', '7'                     BACnetConstructedDataProgramLocation
            [simple   BACnetApplicationTagCharacterString                               programLocation                 ]
        ]
        [*, 'PROGRAM_STATE', '9'                        BACnetConstructedDataProgramState
            [simple   BACnetProgramStateTagged('0', 'TagClass.APPLICATION_TAGS')           programState                 ]
        ]
        [*, 'PROPERTY_LIST', '9'                        BACnetConstructedDataPropertyList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetPropertyIdentifierTagged('0', 'TagClass.APPLICATION_TAGS')
                                propertyList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'PROPORTIONAL_CONSTANT', '4'                BACnetConstructedDataProportionalConstant
            [simple BACnetApplicationTagReal                                proportionalConstant                        ]
        ]
        [*, 'PROPORTIONAL_CONSTANT_UNITS', '9'          BACnetConstructedDataProportionalConstantUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'PROTOCOL_LEVEL', '9'                       BACnetConstructedDataProtocolLevel
            [simple   BACnetProtocolLevelTagged('0', 'TagClass.APPLICATION_TAGS')                    protocolLevel      ]
        ]
        //[*, 'PROTOCOL_CONFORMANCE_CLASS'              BACnetConstructedDataProtocolConformanceClass [validation    '1 == 2'    "TODO: implement me PROTOCOL_CONFORMANCE_CLASS BACnetConstructedDataProtocolConformanceClass"]]
        [*, 'PROTOCOL_OBJECT_TYPES_SUPPORTED', '8'      BACnetConstructedDataProtocolObjectTypesSupported
            [simple   BACnetObjectTypesSupportedTagged('0', 'TagClass.APPLICATION_TAGS')            protocolObjectTypesSupported         ]
        ]
        [*, 'PROTOCOL_REVISION', '2'                    BACnetConstructedDataProtocolRevision
            [simple   BACnetApplicationTagUnsignedInteger                                           protocolRevision    ]
        ]
        [*, 'PROTOCOL_SERVICES_SUPPORTED','8'           BACnetConstructedDataProtocolServicesSupported
            [simple   BACnetServicesSupportedTagged('0', 'TagClass.APPLICATION_TAGS')   protocolServicesSupported       ]
        ]
        [*, 'PROTOCOL_VERSION', '2'                     BACnetConstructedDataProtocolVersion
            [simple   BACnetApplicationTagUnsignedInteger                                           protocolVersion     ]
        ]
        [*, 'PULSE_RATE', '2'                           BACnetConstructedDataPulseRate
            [simple   BACnetApplicationTagUnsignedInteger                                           pulseRate           ]
        ]
        [*, 'READ_ONLY', '1'                            BACnetConstructedDataReadOnly
            [simple BACnetApplicationTagBoolean                                readOnly                                 ]
        ]
        [*, 'REASON_FOR_DISABLE', '9'                   BACnetConstructedDataReasonForDisable
            [array    BACnetAccessCredentialDisableReasonTagged('0', 'TagClass.APPLICATION_TAGS')
                                            reasonForDisable
                                                    terminated
                                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'REASON_FOR_HALT', '9'                      BACnetConstructedDataReasonForHalt
            [simple   BACnetProgramErrorTagged('0', 'TagClass.APPLICATION_TAGS')           programError                 ]
        ]
        [*, 'RECIPIENT_LIST'                            BACnetConstructedDataRecipientList
            [array    BACnetDestination
                                recipientList
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        ['FILE', 'RECORD_COUNT', '2'                    BACnetConstructedDataFileRecordCount
            [simple BACnetApplicationTagUnsignedInteger                     recordCount                                 ]
        ]
        [*, 'RECORD_COUNT', '2'                         BACnetConstructedDataRecordCount
            [simple BACnetApplicationTagUnsignedInteger                     recordCount                                 ]
        ]
        [*, 'RECORDS_SINCE_NOTIFICATION', '2'           BACnetConstructedDataRecordsSinceNotification
            [simple BACnetApplicationTagUnsignedInteger                     recordsSinceNotifications                   ]
        ]
        [*, 'REFERENCE_PORT', '2'                       BACnetConstructedDataReferencePort
            [simple BACnetApplicationTagUnsignedInteger                     referencePort                               ]
        ]
        [*, 'REGISTERED_CAR_CALL'                       BACnetConstructedDataRegisteredCarCall
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetLiftCarCallList
                                            registeredCarCall
                                                    terminated
                                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'RELIABILITY', '9'                          BACnetConstructedDataReliability
            [simple   BACnetReliabilityTagged('0', 'TagClass.APPLICATION_TAGS')                         reliability     ]
        ]
        [*, 'RELIABILITY_EVALUATION_INHIBIT', '1'       BACnetConstructedDataReliabilityEvaluationInhibit
            [simple   BACnetApplicationTagBoolean                                       reliabilityEvaluationInhibit    ]
        ]
        ['ACCESS_DOOR', 'RELINQUISH_DEFAULT', '9'       BACnetConstructedDataAccessDoorRelinquishDefault
            [simple   BACnetDoorValueTagged('0', 'TagClass.APPLICATION_TAGS')           relinquishDefault               ]
        ]
        ['ANALOG_OUTPUT', 'RELINQUISH_DEFAULT', '4'     BACnetConstructedDataAnalogOutputRelinquishDefault
            [simple   BACnetApplicationTagReal                                          relinquishDefault               ]
        ]
        ['ANALOG_VALUE', 'RELINQUISH_DEFAULT', '4'      BACnetConstructedDataAnalogValueRelinquishDefault
            [simple   BACnetApplicationTagReal                                          relinquishDefault               ]
        ]
        ['BINARY_OUTPUT', 'RELINQUISH_DEFAULT', '9'    BACnetConstructedDataBinaryOutputRelinquishDefault
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')            relinquishDefault               ]
        ]
        ['BINARY_VALUE', 'RELINQUISH_DEFAULT', '9'    BACnetConstructedDataBinaryValueRelinquishDefault
            [simple   BACnetBinaryPVTagged('0', 'TagClass.APPLICATION_TAGS')            relinquishDefault               ]
        ]
        ['BINARY_LIGHTING_OUTPUT', 'RELINQUISH_DEFAULT', '9'    BACnetConstructedDataBinaryLightingOutputRelinquishDefault
            [simple   BACnetBinaryLightingPVTagged('0', 'TagClass.APPLICATION_TAGS')    relinquishDefault               ]
        ]
        ['BITSTRING_VALUE', 'RELINQUISH_DEFAULT', '8'   BACnetConstructedDataBitStringValueRelinquishDefault
            [simple   BACnetApplicationTagBitString                                     relinquishDefault               ]
        ]
        ['CHARACTERSTRING_VALUE', 'RELINQUISH_DEFAULT', '7' BACnetConstructedDataCharacterStringValueRelinquishDefault
            [simple   BACnetApplicationTagCharacterString                               relinquishDefault               ]
        ]
        ['DATE_VALUE', 'RELINQUISH_DEFAULT', '10'       BACnetConstructedDataDateValueRelinquishDefault
            [simple   BACnetApplicationTagDate                                          relinquishDefault               ]
        ]
        ['DATEPATTERN_VALUE', 'RELINQUISH_DEFAULT', '10'     BACnetConstructedDataDatePatternValueRelinquishDefault
            [simple   BACnetApplicationTagDate                                          relinquishDefault               ]
        ]
        ['DATETIME_VALUE', 'RELINQUISH_DEFAULT'         BACnetConstructedDataDateTimeValueRelinquishDefault
            [simple   BACnetDateTime                                                    relinquishDefault               ]
        ]
        ['DATETIMEPATTERN_VALUE', 'RELINQUISH_DEFAULT'  BACnetConstructedDataDateTimePatternValueRelinquishDefault
            [simple   BACnetDateTime                                                    relinquishDefault               ]
        ]
        ['LARGE_ANALOG_VALUE', 'RELINQUISH_DEFAULT', '5'     BACnetConstructedDataLargeAnalogValueRelinquishDefault
            [simple   BACnetApplicationTagDouble                                        relinquishDefault               ]
        ]
        ['LIGHTING_OUTPUT', 'RELINQUISH_DEFAULT', '4'   BACnetConstructedDataLightingOutputRelinquishDefault
            [simple   BACnetApplicationTagReal                                          relinquishDefault               ]
        ]
        ['TIMEPATTERN_VALUE', 'RELINQUISH_DEFAULT', '11'    BACnetConstructedDataTimePatternValueRelinquishDefault
            [simple   BACnetApplicationTagTime                                          relinquishDefault               ]
        ]
        ['TIME_VALUE', 'RELINQUISH_DEFAULT', '11'       BACnetConstructedDataTimeValueRelinquishDefault
            [simple   BACnetApplicationTagTime                                          relinquishDefault               ]
        ]
        ['INTEGER_VALUE', 'RELINQUISH_DEFAULT', '3'     BACnetConstructedDataIntegerValueRelinquishDefault
            [simple   BACnetApplicationTagSignedInteger                                 relinquishDefault               ]
        ]
        ['OCTETSTRING_VALUE', 'RELINQUISH_DEFAULT', '6' BACnetConstructedDataOctetStringValueRelinquishDefault
            [simple   BACnetApplicationTagSignedInteger                                 relinquishDefault               ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'RELINQUISH_DEFAULT', '2'    BACnetConstructedDataPositiveIntegerValueRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                               relinquishDefault               ]
        ]
        ['MULTI_STATE_OUTPUT', 'RELINQUISH_DEFAULT', '2'    BACnetConstructedDataMultiStateOutputRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                               relinquishDefault               ]
        ]
        ['MULTI_STATE_VALUE', 'RELINQUISH_DEFAULT', '2' BACnetConstructedDataMultiStateValueRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                               relinquishDefault               ]
        ]
        [*, 'RELINQUISH_DEFAULT', '2'                   BACnetConstructedDataRelinquishDefault
            [simple   BACnetApplicationTagUnsignedInteger                               relinquishDefault               ]
        ]
        [*, 'REPRESENTS'                                BACnetConstructedDataRepresents
            [simple   BACnetDeviceObjectReference       represents                                                      ]
        ]
        [*, 'REQUESTED_SHED_LEVEL'                      BACnetConstructedDataRequestedShedLevel
            [simple  BACnetShedLevel            requestedShedLevel                                                      ]
        ]
        [*, 'REQUESTED_UPDATE_INTERVAL', '2'            BACnetConstructedDataRequestedUpdateInterval
            [simple BACnetApplicationTagUnsignedInteger                       requestedUpdateInterval                   ]
        ]
        [*, 'REQUIRED'                                BACnetConstructedDataRequired
            [validation    '1 == 2'    "An property identified by REQUIRED should never occur in the wild"]
        ]
        ['LARGE_ANALOG_VALUE', 'RESOLUTION', '5'        BACnetConstructedDataLargeAnalogValueResolution
            [simple   BACnetApplicationTagDouble                                        resolution                      ]
        ]
        ['INTEGER_VALUE', 'RESOLUTION', '3'             BACnetConstructedDataIntegerValueResolution
            [simple   BACnetApplicationTagSignedInteger                                 resolution                      ]
        ]
        ['POSITIVE_INTEGER_VALUE', 'RESOLUTION', '2'    BACnetConstructedDataPositiveIntegerValueResolution
            [simple   BACnetApplicationTagUnsignedInteger                               resolution                      ]
        ]
        ['TIMER', 'RESOLUTION', '2'                     BACnetConstructedDataTimerResolution
            [simple   BACnetApplicationTagUnsignedInteger                               resolution                      ]
        ]
        [*, 'RESOLUTION', '4'                           BACnetConstructedDataResolution
            [simple   BACnetApplicationTagReal                                          resolution                      ]
        ]
        [*, 'RESTART_NOTIFICATION_RECIPIENTS'           BACnetConstructedDataRestartNotificationRecipients
            [array    BACnetRecipient
                                restartNotificationRecipients
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'RESTORE_COMPLETION_TIME', '2'              BACnetConstructedDataRestoreCompletionTime
            [simple BACnetApplicationTagUnsignedInteger                       completionTime                            ]
        ]
        [*, 'RESTORE_PREPARATION_TIME', '2'             BACnetConstructedDataRestorePreparationTime
            [simple   BACnetApplicationTagUnsignedInteger               restorePreparationTime                          ]
        ]
        [*, 'ROUTING_TABLE'                             BACnetConstructedDataRoutingTable
            [array    BACnetRouterEntry
                                routingTable
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SCALE'                                     BACnetConstructedDataScale
            [simple   BACnetScale                       scale                                                           ]
        ]
        [*, 'SCALE_FACTOR', '4'                         BACnetConstructedDataScaleFactor
            [simple   BACnetApplicationTagReal                                          scaleFactor                     ]
        ]
        [*, 'SCHEDULE_DEFAULT'                          BACnetConstructedDataScheduleDefault
            [simple   BACnetConstructedDataElement('objectTypeArgument', 'propertyIdentifierArgument', 'null')  scheduleDefault  ]
        ]
        [*, 'SECURED_STATUS', '9'                       BACnetConstructedDataSecuredStatus
            [simple   BACnetDoorSecuredStatusTagged('0', 'TagClass.APPLICATION_TAGS')         securedStatus             ]
        ]
        [*, 'SECURITY_PDU_TIMEOUT', '2'                 BACnetConstructedDataSecurityPDUTimeout
            [simple BACnetApplicationTagUnsignedInteger                       securityPduTimeout                        ]
        ]
        [*, 'SECURITY_TIME_WINDOW', '2'                 BACnetConstructedDataSecurityTimeWindow
            [simple BACnetApplicationTagUnsignedInteger                       securityTimeWindow                        ]
        ]
        [*, 'SEGMENTATION_SUPPORTED', '9'               BACnetConstructedDataSegmentationSupported
            [simple BACnetSegmentationTagged('0', 'TagClass.APPLICATION_TAGS')  segmentationSupported                   ]
        ]
        [*, 'SERIAL_NUMBER', '7'                        BACnetConstructedDataSerialNumber
            [simple   BACnetApplicationTagCharacterString                         serialNumber                          ]
        ]
        [*, 'SETPOINT', '4'                             BACnetConstructedDataSetpoint
            [simple   BACnetApplicationTagReal                                          setpoint                        ]
        ]
        [*, 'SETPOINT_REFERENCE'                        BACnetConstructedDataSetpointReference
            [simple   BACnetSetpointReference                                           setpointReference               ]
        ]
        [*, 'SETTING', '2'                              BACnetConstructedDataSetting
            [simple   BACnetApplicationTagUnsignedInteger                     setting                                   ]
        ]
        [*, 'SHED_DURATION', '2'                        BACnetConstructedDataShedDuration
            [simple BACnetApplicationTagUnsignedInteger                       shedDuration                              ]
        ]
        [*, 'SHED_LEVEL_DESCRIPTIONS', '7'              BACnetConstructedDataShedLevelDescriptions
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                shedLevelDescriptions
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SHED_LEVELS', '2'                          BACnetConstructedDataShedLevels
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                shedLevels
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SILENCED', '9'                             BACnetConstructedDataSilenced
            [simple   BACnetSilencedStateTagged('0', 'TagClass.APPLICATION_TAGS')       silenced                        ]
        ]
        [*, 'SLAVE_ADDRESS_BINDING'                     BACnetConstructedDataSlaveAddressBinding
            [array    BACnetAddressBinding
                                slaveAddressBinding
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SLAVE_PROXY_ENABLE', '1'                   BACnetConstructedDataSlaveProxyEnable
            [simple   BACnetApplicationTagBoolean                                       slaveProxyEnable                ]
        ]
        [*, 'START_TIME'                                BACnetConstructedDataStartTime
            [simple   BACnetDateTime                                          startTime                                 ]
        ]
        [*, 'STATE_CHANGE_VALUES'                       BACnetConstructedDataStateChangeValues
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
        [*, 'STATE_DESCRIPTION', '7'                    BACnetConstructedDataStateDescription
            [simple   BACnetApplicationTagCharacterString                               stateDescription                ]
        ]
        [*, 'STATE_TEXT', '7'                           BACnetConstructedDataStateText
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                                stateText
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'STATUS_FLAGS', '8'                         BACnetConstructedDataStatusFlags
            [simple   BACnetStatusFlagsTagged('0', 'TagClass.APPLICATION_TAGS')         statusFlags                     ]
        ]
        [*, 'STOP_TIME'                                 BACnetConstructedDataStopTime
            [simple   BACnetDateTime                                          stopTime                                  ]
        ]
        [*, 'STOP_WHEN_FULL', '1'                       BACnetConstructedDataStopWhenFull
            [simple BACnetApplicationTagBoolean                                          stopWhenFull                   ]
        ]
        [*, 'STRIKE_COUNT', '2'                         BACnetConstructedDataStrikeCount
            [simple BACnetApplicationTagUnsignedInteger                       strikeCount                               ]
        ]
        [*, 'STRUCTURED_OBJECT_LIST'                    BACnetConstructedDataStructuredObjectList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagObjectIdentifier
                                structuredObjectList
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'SUBORDINATE_ANNOTATIONS'                   BACnetConstructedDataSubordinateAnnotations
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagCharacterString
                    subordinateAnnotations
                            terminated
                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_LIST'                          BACnetConstructedDataSubordinateList
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetDeviceObjectReference
                        subordinateList
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'SUBORDINATE_NODE_TYPES', '9'               BACnetConstructedDataSubordinateNodeTypes
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNodeTypeTagged('0', 'TagClass.APPLICATION_TAGS')
                        subordinateNodeTypes
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'SUBORDINATE_RELATIONSHIPS', '9'            BACnetConstructedDataSubordinateRelationships
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetRelationshipTagged('0', 'TagClass.APPLICATION_TAGS')
                                    subordinateRelationships
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBORDINATE_TAGS'                          BACnetConstructedDataSubordinateTags
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNameValueCollection('0')
                        subordinateList
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUBSCRIBED_RECIPIENTS'                     BACnetConstructedDataSubscribedRecipients
            [array    BACnetEventNotificationSubscription
                            subscribedRecipients
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'            ]
        ]
        [*, 'SUPPORTED_FORMAT_CLASSES'                  BACnetConstructedDataSupportedFormatClasses
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetApplicationTagUnsignedInteger
                                        supportedFormats
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]

        ]
        [*, 'SUPPORTED_FORMATS'                         BACnetConstructedDataSupportedFormats
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetAuthenticationFactorFormat
                                        supportedFormats
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'SUPPORTED_SECURITY_ALGORITHMS'             BACnetConstructedDataSupportedSecurityAlgorithms
            [array    BACnetApplicationTagUnsignedInteger
                                        supportedSecurityAlgorithms
                                                terminated
                                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'SYSTEM_STATUS', '9'                        BACnetConstructedDataSystemStatus
            [simple   BACnetDeviceStatusTagged('0', 'TagClass.APPLICATION_TAGS')    systemStatus                        ]
        ]
        [*, 'TAGS'                                      BACnetConstructedDataTags
            // TODO: uint 64 ---> big int in java == boom
            [virtual  uint 64   zero           '0'  ]
            [optional BACnetApplicationTagUnsignedInteger numberOfDataElements 'arrayIndexArgument!=null && arrayIndexArgument.actualValue == zero']
            [array    BACnetNameValue
                            tags
                                terminated
                                'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'        ]
        ]
        [*, 'THREAT_AUTHORITY'                          BACnetConstructedDataThreatAuthority
            [simple   BACnetAccessThreatLevel                                           threatAuthority                 ]
        ]
        [*, 'THREAT_LEVEL'                              BACnetConstructedDataThreatLevel
            [simple   BACnetAccessThreatLevel                                           threatLevel                     ]
        ]
        [*, 'TIME_DELAY', '2'                           BACnetConstructedDataTimeDelay
            [simple   BACnetApplicationTagUnsignedInteger                               timeDelay                       ]
        ]
        [*, 'TIME_DELAY_NORMAL', '2'                    BACnetConstructedDataTimeDelayNormal
            [simple   BACnetApplicationTagUnsignedInteger                               timeDelayNormal                 ]
        ]
        [*, 'TIME_OF_ACTIVE_TIME_RESET'                 BACnetConstructedDataTimeOfActiveTimeReset
            [simple   BACnetDateTime                                          timeOfActiveTimeReset                     ]
        ]
        [*, 'TIME_OF_DEVICE_RESTART'                    BACnetConstructedDataTimeOfDeviceRestart
            [simple   BACnetTimeStamp                                                    timeOfDeviceRestart            ]
        ]
        [*, 'TIME_OF_STATE_COUNT_RESET'                 BACnetConstructedDataTimeOfStateCountReset
            [simple   BACnetDateTime                                          timeOfStateCountReset                     ]
        ]
        [*, 'TIME_OF_STRIKE_COUNT_RESET'                BACnetConstructedDataTimeOfStrikeCountReset
            [simple   BACnetDateTime                                        timeOfStrikeCountReset                      ]
        ]
        [*, 'TIME_SYNCHRONIZATION_INTERVAL', '2'        BACnetConstructedDataTimeSynchronizationInterval
            [simple   BACnetApplicationTagUnsignedInteger                               timeSynchronization             ]
        ]
        [*, 'TIME_SYNCHRONIZATION_RECIPIENTS'           BACnetConstructedDataTimeSynchronizationRecipients
            [array    BACnetRecipient
                            timeSynchronizationRecipients
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'TIMER_RUNNING', '1'                        BACnetConstructedDataTimerRunning
            [simple   BACnetApplicationTagBoolean                               timerRunning                            ]
        ]
        [*, 'TIMER_STATE', '9'                          BACnetConstructedDataTimerState
            [simple  BACnetTimerStateTagged('0', 'TagClass.APPLICATION_TAGS')   timerState] 
        ]
        [*, 'TOTAL_RECORD_COUNT', '2'                   BACnetConstructedDataTotalRecordCount
            [simple   BACnetApplicationTagUnsignedInteger                     totalRecordCount                          ]
        ]
        [*, 'TRACE_FLAG', '1'                           BACnetConstructedDataTraceFlag
            [simple   BACnetApplicationTagBoolean                               traceFlag                               ]
        ]
        ['LIGHTING_OUTPUT','TRACKING_VALUE', '4'        BACnetConstructedDataLightingOutputTrackingValue
            [simple   BACnetApplicationTagReal                                  trackingValue                           ]
        ]
        ['LIGHTING_OUTPUT','TRACKING_VALUE', '4'        BACnetConstructedDataLightingOutputTrackingValue
            [simple   BACnetApplicationTagReal                                  trackingValue                           ]
        ]
        [*, 'TRACKING_VALUE', '9'                       BACnetConstructedDataTrackingValue
            [simple   BACnetLifeSafetyStateTagged('0', 'TagClass.APPLICATION_TAGS') trackingValue                       ]
        ]
        [*, 'TRANSACTION_NOTIFICATION_CLASS', '2'       BACnetConstructedDataTransactionNotificationClass
            [simple BACnetApplicationTagUnsignedInteger                     transactionNotificationClass                ]
        ]
        [*, 'TRANSITION', '9'                           BACnetConstructedDataTransition
            [simple BACnetLightingTransitionTagged('0', 'TagClass.APPLICATION_TAGS')                    transition      ]
        ]
        [*, 'TRIGGER', '1'                              BACnetConstructedDataTrigger
            [simple BACnetApplicationTagBoolean                               trigger                                   ]
        ]
        [*, 'UNITS', '9'                                BACnetConstructedDataUnits
            [simple   BACnetEngineeringUnitsTagged('0', 'TagClass.APPLICATION_TAGS')                    units           ]
        ]
        [*, 'UPDATE_INTERVAL', '2'                      BACnetConstructedDataUpdateInterval
            [simple   BACnetApplicationTagUnsignedInteger                               updateInterval                  ]
        ]
        [*, 'UPDATE_KEY_SET_TIMEOUT', '2'               BACnetConstructedDataUpdateKeySetTimeout
            [simple BACnetApplicationTagUnsignedInteger                       updateKeySetTimeout                       ]
        ]
        ['CREDENTIAL_DATA_INPUT', 'UPDATE_TIME'         BACnetConstructedDataCredentialDataInputUpdateTime
            [simple   BACnetTimeStamp                                          updateTime                               ]
        ]
        [*, 'UPDATE_TIME'                               BACnetConstructedDataUpdateTime
            [simple   BACnetDateTime                                          updateTime                                ]
        ]
        [*, 'USER_EXTERNAL_IDENTIFIER', '7'             BACnetConstructedDataUserExternalIdentifier
            [simple   BACnetApplicationTagCharacterString   userExternalIdentifier                                      ]
        ]
        [*, 'USER_INFORMATION_REFERENCE', '7'           BACnetConstructedDataUserInformationReference
            [simple   BACnetApplicationTagCharacterString   userInformationReference                                    ]
        ]
        [*, 'USER_NAME', '7'                            BACnetConstructedDataUserName
            [simple   BACnetApplicationTagCharacterString   userName                                                    ]
        ]
        [*, 'USER_TYPE', '9'                            BACnetConstructedDataUserType
            [simple   BACnetAccessUserTypeTagged('0', 'TagClass.APPLICATION_TAGS')      userType                        ]
        ]
        [*, 'USES_REMAINING', '3'                       BACnetConstructedDataUsesRemaining
            [simple   BACnetApplicationTagSignedInteger                     usesRemaining                               ]
        ]
        [*, 'UTC_OFFSET', '3'                           BACnetConstructedDataUTCOffset
            [simple   BACnetApplicationTagSignedInteger                     utcOffset                                   ]
        ]
        [*, 'UTC_TIME_SYNCHRONIZATION_RECIPIENTS'       BACnetConstructedDataUTCTimeSynchronizationRecipients
            [array    BACnetRecipient
                                utcTimeSynchronizationRecipients
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'VALID_SAMPLES', '2'                        BACnetConstructedDataValidSamples
            [simple   BACnetApplicationTagUnsignedInteger                               validSamples                    ]
        ]
        [*, 'VALUE_BEFORE_CHANGE', '2'                  BACnetConstructedDataValueBeforeChange
            [simple   BACnetApplicationTagUnsignedInteger                               valuesBeforeChange              ]
        ]
        [*, 'VALUE_CHANGE_TIME'                         BACnetConstructedDataValueChangeTime
            [simple   BACnetDateTime                                        valueChangeTime                             ]
        ]
        [*, 'VALUE_SET', '2'                            BACnetConstructedDataValueSet
            [simple   BACnetApplicationTagUnsignedInteger                               valueSet                        ]
        ]
        [*, 'VALUE_SOURCE'                              BACnetConstructedDataValueSource
            [simple   BACnetValueSource                 valueSource                                                     ]
        ]
        [*, 'VALUE_SOURCE_ARRAY'                        BACnetConstructedDataValueSourceArray
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
        [*, 'VARIANCE_VALUE', '4'                       BACnetConstructedDataVarianceValue
            [simple BACnetApplicationTagReal                                varianceValue                               ]
        ]
        // Note: checking 2 here is no accident as vendor-id is usually represented by unsigned not enumerated...
        //       the enum is a addition from plc4x
        [*, 'VENDOR_IDENTIFIER', '2'                    BACnetConstructedDataVendorIdentifier
            [simple   BACnetVendorIdTagged('0', 'TagClass.APPLICATION_TAGS') vendorIdentifier                           ]
        ]
        [*, 'VENDOR_NAME', '7'                          BACnetConstructedDataVendorName
           [simple   BACnetApplicationTagCharacterString    vendorName                                                  ]
        ]
        [*, 'VERIFICATION_TIME', '3'                    BACnetConstructedDataVerificationTime
            [simple   BACnetApplicationTagSignedInteger                     verificationTime                            ]
        ]
        [*, 'VIRTUAL_MAC_ADDRESS_TABLE'                 BACnetConstructedDataVirtualMACAddressTable
            [array    BACnetVMACEntry
                                virtualMacAddressTable
                                    terminated
                                    'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)'    ]
        ]
        [*, 'VT_CLASSES_SUPPORTED', '9'                 BACnetConstructedDataVTClassesSupported
            [array    BACnetVTClassTagged('0', 'TagClass.APPLICATION_TAGS')
                                vtClassesSupported
                                        terminated
                                        'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, tagNumber)']
        ]
        [*, 'WEEKLY_SCHEDULE'                           BACnetConstructedDataWeeklySchedule
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
        [*, 'WINDOW_INTERVAL', '2'                      BACnetConstructedDataWindowInterval
            [simple   BACnetApplicationTagUnsignedInteger                               windowInterval                  ]
        ]
        [*, 'WINDOW_SAMPLES', '2'                       BACnetConstructedDataWindowSamples
            [simple   BACnetApplicationTagUnsignedInteger                               windowSamples                   ]
        ]
        [*, 'WRITE_STATUS', '9'                         BACnetConstructedDataWriteStatus
            [simple   BACnetWriteStatusTagged('0', 'TagClass.APPLICATION_TAGS') writeStatus                             ]
        ]
        [*, 'ZONE_FROM'                                 BACnetConstructedDataZoneFrom
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
            [simple   BACnetConstructedData('1', 'objectTypeArgument', 'BACnetPropertyIdentifier.VENDOR_PROPRIETARY_VALUE', 'null')
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
            [simple   BACnetCalendarEntryEnclosed('0')
                                    calendarEntry                                    ]
        ]
        ['1' BACnetSpecialEventPeriodCalendarReference
            [simple   BACnetContextTagObjectIdentifier('1', 'BACnetDataType.BACNET_OBJECT_IDENTIFIER')
                                   calendarReference                                ]
        ]
    ]
]

[type BACnetCalendarEntry
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
        ['0' BACnetProcessIdSelectionNull
            [simple   BACnetApplicationTagNull
                            nullValue                   ]
        ]
        [BACnetProcessIdSelectionValue
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
        ['0' BACnetHostAddressNull
            [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                                        none                                                                ]
        ]
        ['1' BACnetHostAddressIpAddress
            [simple   BACnetContextTagOctetString('1', 'BACnetDataType.OCTET_STRING')
                                        ipAddress                                                           ]
        ]
        ['2' BACnetHostAddressName
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
        ['0' BACnetShedLevelPercent
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        percent                                                             ]
        ]
        ['1' BACnetShedLevelLevel
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        level                                                               ]
        ]
        ['2' BACnetShedLevelAmount
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
        ['0' BACnetFaultParameterNone
            [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                                                none                                                    ]
        ]
        ['1' BACnetFaultParameterFaultCharacterString
            [simple   BACnetOpeningTag('1')
                                         openingTag                                                     ]
            [simple   BACnetFaultParameterFaultCharacterStringListOfFaultValues('0')
                                        listOfFaultValues                                               ]
            [simple   BACnetClosingTag('1')
                                         closingTag                                                     ]
        ]
        ['2' BACnetFaultParameterFaultExtended
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
        ['3' BACnetFaultParameterFaultLifeSafety
            [simple   BACnetOpeningTag('3')
                                        openingTag                                                      ]
            [simple   BACnetFaultParameterFaultLifeSafetyListOfFaultValues('0')
                                        listOfFaultValues                                               ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        modePropertyReference                                           ]
            [simple   BACnetClosingTag('3')
                                        closingTag                                                      ]
        ]
        ['4' BACnetFaultParameterFaultState
            [simple   BACnetOpeningTag('4')
                                        openingTag                                                      ]
            [simple   BACnetFaultParameterFaultStateListOfFaultValues('0')
                                        listOfFaultValues                                               ]
            [simple   BACnetClosingTag('4')
                                        closingTag                                                      ]
        ]
        ['5' BACnetFaultParameterFaultStatusFlags
            [simple   BACnetOpeningTag('5')
                                        openingTag                                                      ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        statusFlagsReference                                            ]
            [simple   BACnetClosingTag('5')
                                        closingTag                                                      ]
        ]
        ['6' BACnetFaultParameterFaultOutOfRange
            [simple   BACnetOpeningTag('6')
                                        openingTag                                                      ]
            [simple   BACnetFaultParameterFaultOutOfRangeMinNormalValue('0')
                                        minNormalValue                                                  ]
            [simple   BACnetFaultParameterFaultOutOfRangeMaxNormalValue('0')
                                        maxNormalValue                                                  ]
            [simple   BACnetClosingTag('6')
                                        closingTag                                                      ]
        ]
        ['7' BACnetFaultParameterFaultListed
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'       ]
    [virtual bit        peekedIsContextTag  'peekedTagHeader.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS']
    [validation '(!peekedIsContextTag) || (peekedIsContextTag && peekedTagHeader.lengthValueType != 0x6 && peekedTagHeader.lengthValueType != 0x7)'
                "unexpected opening or closing tag"                                 ]
    [typeSwitch peekedTagNumber, peekedIsContextTag
       ['0x0', 'false' BACnetFaultParameterFaultExtendedParametersEntryNull
           [simple  BACnetApplicationTagNull
                            nullValue                                                   ]
       ]
       ['0x4', 'false' BACnetFaultParameterFaultExtendedParametersEntryReal
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x2', 'false' BACnetFaultParameterFaultExtendedParametersEntryUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x1', 'false' BACnetFaultParameterFaultExtendedParametersEntryBoolean
           [simple   BACnetApplicationTagBoolean
                            booleanValue                                                ]
       ]
       ['0x3', 'false' BACnetFaultParameterFaultExtendedParametersEntryInteger
           [simple   BACnetApplicationTagSignedInteger
                            integerValue                                                ]
       ]
       ['0x5', 'false' BACnetFaultParameterFaultExtendedParametersEntryDouble
           [simple  BACnetApplicationTagDouble
                                doubleValue                                             ]
       ]
       ['0x6', 'false' BACnetFaultParameterFaultExtendedParametersEntryOctetString
           [simple   BACnetApplicationTagOctetString
                            octetStringValue                                            ]
       ]
       ['0x7', 'false' BACnetFaultParameterFaultExtendedParametersEntryCharacterString
           [simple   BACnetApplicationTagCharacterString
                            characterStringValue                                        ]
       ]
       ['0x8', 'false' BACnetFaultParameterFaultExtendedParametersEntryBitString
           [simple   BACnetApplicationTagBitString
                            bitStringValue                                              ]
       ]
       ['0x9', 'false' BACnetFaultParameterFaultExtendedParametersEntryEnumerated
           [simple   BACnetApplicationTagEnumerated
                            enumeratedValue                                             ]
       ]
       ['0xA', 'false' BACnetFaultParameterFaultExtendedParametersEntryDate
           [simple   BACnetApplicationTagDate
                            dateValue                                                   ]
       ]
       ['0xB', 'false' BACnetFaultParameterFaultExtendedParametersEntryTime
           [simple   BACnetApplicationTagTime
                            timeValue                                                   ]
       ]
       ['0xC', 'false' BACnetFaultParameterFaultExtendedParametersEntryObjectidentifier
           [simple   BACnetApplicationTagObjectIdentifier
                            objectidentifierValue                                       ]
       ]
       ['0', 'true' BACnetFaultParameterFaultExtendedParametersEntryReference
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'           ]
    [validation 'peekedTagHeader.tagClass == TagClass.APPLICATION_TAGS'
                "only application tags allowed"                                         ]
    [typeSwitch peekedTagNumber
       ['0x4' BACnetFaultParameterFaultOutOfRangeMinNormalValueReal
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x2' BACnetFaultParameterFaultOutOfRangeMinNormalValueUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x5' BACnetFaultParameterFaultOutOfRangeMinNormalValueDouble
           [simple  BACnetApplicationTagDouble
                            doubleValue                                                 ]
       ]
       ['0x3' BACnetFaultParameterFaultOutOfRangeMinNormalValueInteger
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
    [virtual uint 8     peekedTagNumber     'peekedTagHeader.actualTagNumber'           ]
    [validation 'peekedTagHeader.tagClass == TagClass.APPLICATION_TAGS'
                "only application tags allowed"                                         ]
    [typeSwitch peekedTagNumber
       ['0x4' BACnetFaultParameterFaultOutOfRangeMaxNormalValueReal
           [simple  BACnetApplicationTagReal
                            realValue                                                   ]
       ]
       ['0x2' BACnetFaultParameterFaultOutOfRangeMaxNormalValueUnsigned
           [simple   BACnetApplicationTagUnsignedInteger
                            unsignedValue                                               ]
       ]
       ['0x5' BACnetFaultParameterFaultOutOfRangeMaxNormalValueDouble
           [simple  BACnetApplicationTagDouble
                            doubleValue                                                 ]
       ]
       ['0x3' BACnetFaultParameterFaultOutOfRangeMaxNormalValueInteger
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
        ['0' BACnetEventParameterChangeOfBitstring
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
        ['1' BACnetEventParameterChangeOfState
            [simple   BACnetOpeningTag('1')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfStateListOfValues('1')
                                        listOfValues                                    ]
            [simple   BACnetClosingTag('1')
                                        closingTag                                      ]
        ]
        ['2' BACnetEventParameterChangeOfValue
            [simple   BACnetOpeningTag('2')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfValueCivCriteria('1')
                                        covCriteria                                     ]
            [simple   BACnetClosingTag('2')
                                        closingTag                                      ]
        ]
        ['3' BACnetEventParameterCommandFailure
            [simple   BACnetOpeningTag('3')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        feedbackPropertyReference                       ]
            [simple   BACnetClosingTag('3')
                                        closingTag                                      ]
        ]
        ['4' BACnetEventParameterFloatingLimit
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
        ['5' BACnetEventParameterOutOfRange
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
        ['8' BACnetEventParameterChangeOfLifeSavety
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
        ['9' BACnetEventParameterExtended
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
        ['10' BACnetEventParameterBufferReady
            [simple   BACnetOpeningTag('10')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        notificationThreshold                           ]
            [simple   BACnetContextTagUnsignedInteger('1', 'BACnetDataType.UNSIGNED_INTEGER')
                                        previousNotificationCount                       ]
            [simple   BACnetClosingTag('10')
                                        closingTag                                      ]
        ]
        ['11' BACnetEventParameterUnsignedRange
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
        ['13' BACnetEventParameterAccessEvent
            [simple   BACnetOpeningTag('13')
                                         openingTag                                     ]
            [simple   BACnetEventParameterAccessEventListOfAccessEvents('0')
                                        listOfAccessEvents                              ]
            [simple   BACnetDeviceObjectPropertyReferenceEnclosed('1')
                                        accessEventTimeReference                        ]
            [simple   BACnetClosingTag('13')
                                        closingTag                                      ]
        ]
        ['14' BACnetEventParameterDoubleOutOfRange
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
        ['15' BACnetEventParameterSignedOutOfRange
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
        ['16' BACnetEventParameterUnsignedOutOfRange
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
        ['17' BACnetEventParameterChangeOfCharacterString
            [simple   BACnetOpeningTag('17')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetEventParameterChangeOfCharacterStringListOfAlarmValues('1')
                                        listOfAlarmValues                               ]
            [simple   BACnetClosingTag('17')
                                        closingTag                                      ]
        ]
        ['18' BACnetEventParameterChangeOfStatusFlags
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
        ['20' BACnetEventParameterNone
            [simple   BACnetContextTagNull('20', 'BACnetDataType.NULL')
                                        none                                            ]
        ]
        ['21' BACnetEventParameterChangeOfDiscreteValue
            [simple   BACnetOpeningTag('21')
                                         openingTag                                     ]
            [simple   BACnetContextTagUnsignedInteger('0', 'BACnetDataType.UNSIGNED_INTEGER')
                                        timeDelay                                       ]
            [simple   BACnetClosingTag('21')
                                        closingTag                                      ]
        ]
        ['22' BACnetEventParameterChangeOfTimer
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
        ['0' BACnetEventParameterChangeOfValueCivCriteriaBitmask
            [simple   BACnetContextTagBitString('0', 'BACnetDataType.BIT_STRING')
                                        bitmask                                         ]
        ]
        ['1' BACnetEventParameterChangeOfValueCivCriteriaReferencedPropertyIncrement
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
        ['1' BACnetLandingCallStatusCommandDirection
            [simple BACnetLiftCarDirectionTagged('1', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                            direction                               ]
        ]
        ['2' BACnetLandingCallStatusCommandDestination
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
        ['0' BACnetScaleFloatScale
            [simple   BACnetContextTagReal('0', 'BACnetDataType.REAL')
                                        floatScale                                  ]
        ]
        ['1' BACnetScaleIntegerScale
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
        ['0' BACnetValueSourceNone
            [simple   BACnetContextTagNull('0', 'BACnetDataType.NULL')
                                        none                                        ]
        ]
        ['1' BACnetValueSourceObject
             [simple   BACnetDeviceObjectReferenceEnclosed('1')
                                        object                                      ]
        ]
        ['2' BACnetValueSourceAddress
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
        ['0x4' BACnetClientCOVObject
            [simple   BACnetApplicationTagReal
                                        realIncrement                               ]
        ]
        ['0x0' BACnetClientCOVNone
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
        ['0' BACnetLogDataLogStatus
            [simple   BACnetLogStatusTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        logStatus                                   ]
        ]
        ['1' BACnetLogDataLogData
            [simple   BACnetOpeningTag('1')
                                        innerOpeningTag                             ]
            [array    BACnetLogDataLogDataEntry
                                        logData
                                            terminated
                                            'STATIC_CALL("isBACnetConstructedDataClosingTag", readBuffer, false, 1)']
            [simple   BACnetClosingTag('1')
                                        innerClosingTag                             ]
        ]
        ['2' BACnetLogDataLogDataTimeChange
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
        ['0' BACnetLogDataLogDataEntryBooleanValue
            [simple   BACnetContextTagBoolean('0', 'BACnetDataType.BOOLEAN')
                                        booleanValue                                ]
        ]
        ['1' BACnetLogDataLogDataEntryRealValue
            [simple   BACnetContextTagReal('1', 'BACnetDataType.REAL')
                                        realValue                                   ]
        ]
        ['2' BACnetLogDataLogDataEntryEnumeratedValue
            [simple   BACnetContextTagEnumerated('2', 'BACnetDataType.ENUMERATED')
                                        enumeratedValue                             ]
        ]
        ['3' BACnetLogDataLogDataEntryUnsignedValue
            [simple   BACnetContextTagUnsignedInteger('3', 'BACnetDataType.UNSIGNED_INTEGER')
                                        unsignedValue                               ]
        ]
        ['4' BACnetLogDataLogDataEntryIntegerValue
            [simple   BACnetContextTagSignedInteger('4', 'BACnetDataType.SIGNED_INTEGER')
                                        integerValue                                ]
        ]
        ['5' BACnetLogDataLogDataEntryBitStringValue
            [simple   BACnetContextTagBitString('5', 'BACnetDataType.BIT_STRING')
                                        bitStringValue                              ]
        ]
        ['6' BACnetLogDataLogDataEntryNullValue
            [simple   BACnetContextTagNull('6', 'BACnetDataType.NULL')
                                        nullValue                                   ]
        ]
        ['7' BACnetLogDataLogDataEntryFailure
            [simple   ErrorEnclosed('7')
                                        failure                                     ]
        ]
        ['8' BACnetLogDataLogDataEntryAnyValue
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
        ['0' BACnetLogRecordLogDatumLogStatus
            [simple   BACnetLogStatusTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        logStatus                                   ]
        ]
        ['1' BACnetLogRecordLogDatumBooleanValue
            [simple   BACnetContextTagBoolean('1', 'BACnetDataType.BOOLEAN')
                                        booleanValue                                ]
        ]
        ['2' BACnetLogRecordLogDatumRealValue
            [simple   BACnetContextTagReal('2', 'BACnetDataType.REAL')
                                        realValue                                   ]
        ]
        ['3' BACnetLogRecordLogDatumEnumeratedValue
            [simple   BACnetContextTagEnumerated('3', 'BACnetDataType.ENUMERATED')
                                        enumeratedValue                             ]
        ]
        ['4' BACnetLogRecordLogDatumUnsignedValue
            [simple   BACnetContextTagUnsignedInteger('4', 'BACnetDataType.UNSIGNED_INTEGER')
                                        unsignedValue                               ]
        ]
        ['5' BACnetLogRecordLogDatumIntegerValue
            [simple   BACnetContextTagSignedInteger('5', 'BACnetDataType.SIGNED_INTEGER')
                                        integerValue                                ]
        ]
        ['6' BACnetLogRecordLogDatumBitStringValue
            [simple   BACnetContextTagBitString('6', 'BACnetDataType.BIT_STRING')
                                        bitStringValue                              ]
        ]
        ['7' BACnetLogRecordLogDatumNullValue
            [simple   BACnetContextTagNull('7', 'BACnetDataType.NULL')
                                        nullValue                                   ]
        ]
        ['8' BACnetLogRecordLogDatumFailure
            [simple   ErrorEnclosed('8')
                                        failure                                     ]
        ]
        ['9' BACnetLogRecordLogDatumTimeChange
            [simple   BACnetContextTagReal('9', 'BACnetDataType.REAL')
                                        timeChange                                  ]
        ]
        ['10' BACnetLogRecordLogDatumAnyValue
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
        ['0' BACnetEventLogRecordLogDatumLogStatus
            [simple   BACnetLogStatusTagged('0', 'TagClass.CONTEXT_SPECIFIC_TAGS')
                                        logStatus                                   ]
        ]
        ['1' BACnetEventLogRecordLogDatumNotification
            [simple   BACnetOpeningTag('1')
                            innerOpeningTag                                         ]
            //TODO this below slurps to much because of the service choice... :( find workaround we might need fragments for that...
            [simple   ConfirmedEventNotificationRequest
                                        notification                                ]
            [simple   BACnetClosingTag('tagNumber')
                            innerClosingTag                                         ]
        ]
        ['2' BACnetEventLogRecordLogDatumTimeChange
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
        ['4' BACnetPropertyAccessResultAccessResultPropertyValue
            [simple   BACnetConstructedData('4', 'objectTypeArgument', 'propertyIdentifierArgument', 'propertyArrayIndexArgument')
                                        propertyValue                               ]
        ]
        ['5' BACnetPropertyAccessResultAccessResultPropertyAccessError
            [simple   ErrorEnclosed('5')
                                        propertyAccessError                         ]
        ]
    ]
]
