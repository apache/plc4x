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

[discriminatedType MQTT_ControlPacket    byteOrder='BIG_ENDIAN'
    [discriminator MQTT_ControlPacketType packetType                                        ]
    [abstract uint 8 remainingLength                                                        ]
    [typeSwitch packetType
        ['CONNECT' MQTT_ControlPacket_CONNECT
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   MQTT_String            protocolName                                   ]
            [simple   uint 8                 protocolVersion                                ]
            // Connect flags start
            [simple   bit                    userNameFlagSet                                ]
            [simple   bit                    passwordFlagSet                                ]
            [simple   bit                    willRetainFlagSet                              ]
            [simple   uint 2                 willQosLevel                                   ]
            [simple   bit                    willFlagSet                                    ]
            // Actually a session start ...
            [simple   bit                    cleanStartFlagSet                              ]
            [reserved bit                    'false'                                        ]
            // Connect flags end
            [simple   uint 16                keepAlive                                      ]

            // Properties
            [simple   uint 32                propertyLength        encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length 'propertyLength'         ]

            // Payload
            [simple   MQTT_String            clientId                                       ]
            // TODO: If willFlagSet is true, the will properties come here (3.1.3.2)
            // TODO: If willFlagSet is true, the will topic comes here (3.1.3.3)
            // TODO: If willFlagSet is true, the will payload comes here (3.1.3.4)
            // If userNameFlagSet, here comes the username. (String)
            [optional MQTT_String            username       'userNameFlagSet'               ]
            // If passwordFlagSet, here comes the password. (String)
            [optional MQTT_String            password       'passwordFlagSet'               ]
        ]
        ['CONNACK' MQTT_ControlPacket_CONNACK
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            // Acknowledge flags start
            [reserved uint 7                 '0x00'                                         ]
            [simple   bit                    sessionPresentFlagSet                          ]
            // Acknowledge flags end
            [simple   MQTT_ReasonCode        reasonCode                                     ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']
        ]
        ['PUBLISH' MQTT_ControlPacket_PUBLISH
            // Fixed Header
            [simple   bit                    dup                                            ]
            [simple   MQTT_QOS               qos                                            ]
            [simple   bit                    retain                                         ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   MQTT_String            topicName                                      ]
            [optional uint 16                packetIdentifier 'qos != MQTT_QOS.AT_MOST_ONCE']

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']

            // Payload
            [array    byte                   payload      count 'remainingLength - curPos'  ]
        ]
        // Used if QOS = 1
        ['PUBACK' MQTT_ControlPacket_PUBACK
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   uint 16                packetIdentifier                               ]
            [optional MQTT_ReasonCode        reasonCode     'remainingLength - curPos < 3'  ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']
        ]
        // Used if QOS = 2
        ['PUBREC' MQTT_ControlPacket_PUBREC
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   uint 16                packetIdentifier                               ]
            [optional MQTT_ReasonCode        reasonCode     'remainingLength - curPos < 3'  ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']
        ]
        ['PUBREL' MQTT_ControlPacket_PUBREL
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   uint 16                packetIdentifier                               ]
            [optional MQTT_ReasonCode        reasonCode     'remainingLength - curPos < 3'  ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']
        ]
        ['PUBCOMP' MQTT_ControlPacket_PUBCOMP
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   uint 16                packetIdentifier                               ]
            [optional MQTT_ReasonCode        reasonCode     'remainingLength - curPos < 3'  ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']
        ]
        ['SUBSCRIBE' MQTT_ControlPacket_SUBSCRIBE
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   uint 16                packetIdentifier                               ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']

            // Payload
            [array    Filter                 filters      count 'remainingLength - curPos'  ]
        ]
        ['SUBACK' MQTT_ControlPacket_SUBACK
            // Fixed Header
            [reserved uint 4                 '0x0'                                          ]
            [simple   uint 8                 remainingLength                                ]
            // Variable Header
            [simple   uint 16                packetIdentifier                               ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']

            // Payload
            [array    MQTT_ReasonCode        results        count 'remainingLength - curPos'    ]
        ]
        ['UNSUBSCRIBE' MQTT_ControlPacket_UNSUBSCRIBE
            // Fixed Header
            [reserved uint 4                 '0x0'                                              ]
            [simple   uint 8                 remainingLength                                    ]
            // Variable Header
            [simple   uint 16                packetIdentifier                                   ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']

            // Payload
            [array    MQTT_String            filters        count 'remainingLength - curPos'    ]
        ]
        ['UNSUBACK' MQTT_ControlPacket_UNSUBACK
            // Fixed Header
            [reserved uint 4                 '0x0'                                              ]
            [simple   uint 8                 remainingLength                                    ]
            // Variable Header
            [simple   uint 16                packetIdentifier                                   ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']

            // Payload
            [array    MQTT_ReasonCode        results        count 'remainingLength - curPos'    ]
        ]
        ['PINGREQ' MQTT_ControlPacket_PINGREQ
            // Fixed Header
            [reserved uint 4                 '0x0'                                              ]
            [simple   uint 8                 remainingLength                                    ]
        ]
        ['PINGRESP' MQTT_ControlPacket_PINGRESP
            // Fixed Header
            [reserved uint 4                 '0x0'                                              ]
            [simple   uint 8                 remainingLength                                    ]
        ]
        ['DISCONNECT' MQTT_ControlPacket_DISCONNECT
            // Fixed Header
            [reserved uint 4                 '0x0'                                              ]
            [simple   uint 8                 remainingLength                                    ]
            // Variable Header
            [simple   MQTT_ReasonCode        reason                                             ]
        ]
        ['AUTH' MQTT_ControlPacket_AUTH
            // Fixed Header
            [reserved uint 4                 '0x0'                                              ]
            [simple   uint 8                 remainingLength                                    ]
            // Variable Header
            [simple   MQTT_ReasonCode        reason                                             ]

            // Properties
            [optional uint 32                propertyLength 'remainingLength - curPos < 4'       encoding='"varLenUint32"']
            [array    MQTT_Property          properties     length '(propertyLength != null) ? propertyLength : 0']
        ]
    ]
]

[discriminatedType MQTT_Property
    [simple MQTT_PropertyType propertyType]
    [typeSwitch propertyType
        ['PAYLOAD_FORMAT_INDICATOR'          MQTT_Property_PAYLOAD_FORMAT_INDICATOR
            [simple uint 8                   value                              ]
        ]
        ['MESSAGE_EXPIRY_INTERVAL'           MQTT_Property_MESSAGE_EXPIRY_INTERVAL
            [simple uint 32                  value                              ]
        ]
        ['CONTENT_TYPE'                      MQTT_Property_CONTENT_TYPE
            [simple   MQTT_String            value                              ]
        ]
        ['RESPONSE_TOPIC'                    MQTT_Property_RESPONSE_TOPIC
            [simple   MQTT_String            value                              ]
        ]
        ['CORRELATION_DATA'                  MQTT_Property_CORRELATION_DATA
            // TODO: Find out what "Binary Data" means ...
        ]
        ['SUBSCRIPTION_IDENTIFIER'           MQTT_Property_SUBSCRIPTION_IDENTIFIER
            [simple uint 32                  value     encoding='"varLenUint32"']
        ]
        ['SESSION_EXPIRY_INTERVAL'           MQTT_Property_EXPIRY_INTERVAL
            [simple uint 32                  value                              ]
        ]
        ['ASSIGNED_CLIENT_IDENTIFIER'        MQTT_Property_ASSIGNED_CLIENT_IDENTIFIER
            [simple   MQTT_String            value                              ]
        ]
        ['SERVER_KEEP_ALIVE'                 MQTT_Property_SERVER_KEEP_ALIVE
            [simple uint 16                  value                              ]
        ]
        ['AUTHENTICATION_METHOD'             MQTT_Property_AUTHENTICATION_METHOD
            [simple   MQTT_String            value                              ]
        ]
        ['AUTHENTICATION_DATA'               MQTT_Property_AUTHENTICATION_DATA
            // TODO: Find out what "Binary Data" means ...
        ]
        ['REQUEST_PROBLEM_INFORMATION'       MQTT_Property_REQUEST_PROBLEM_INFORMATION
            // TODO: Probably an enum
            [simple uint 8                   value                              ]
        ]
        ['WILL_DELAY_INTERVAL'               MQTT_Property_WILL_DELAY_INTERVAL
            [simple uint 32                  value                              ]
        ]
        ['REQUEST_RESPONSE_INFORMATION'      MQTT_Property_REQUEST_RESPONSE_INFORMATION
            // TODO: Probably an enum
            [simple uint 8                   value                              ]
        ]
        ['RESPONSE_INFORMATION'              MQTT_Property_RESPONSE_INFORMATION
            [simple   MQTT_String            value                              ]
        ]
        ['SERVER_REFERENCE'                  MQTT_Property_SERVER_REFERENCE
            [simple   MQTT_String            value                              ]
        ]
        ['REASON_STRING'                     MQTT_Property_REASON_STRING
            [simple   MQTT_String            value                              ]
        ]
        ['RECEIVE_MAXIMUM'                   MQTT_Property_RECEIVE_MAXIMUM
            [simple uint 16                  value                              ]
        ]
        ['TOPIC_ALIAS_MAXIMUM'               MQTT_Property_TOPIC_ALIAS_MAXIMUM
            [simple uint 16                  value                              ]
        ]
        ['TOPIC_ALIAS'                       MQTT_Property_TOPIC_ALIAS
            [simple uint 16                  value                              ]
        ]
        ['MAXIMUM_QOS'                       MQTT_Property_MAXIMUM_QOS
            [simple uint 8                   value                              ]
        ]
        ['RETAIN_AVAILABLE'                  MQTT_Property_RETAIN_AVAILABLE
            [simple uint 8                   value                              ]
        ]
        ['USER_PROPERTY'                     MQTT_Property_USER_PROPERTY
            [simple   MQTT_String            name                               ]
            [simple   MQTT_String            value                              ]
        ]
        ['MAXIMUM_PACKET_SIZE'               MQTT_Property_MAXIMUM_PACKET_SIZE
            [simple uint 32                  value                              ]
        ]
        ['WILDCARD_SUBSCRIPTION_AVAILABLE'   MQTT_Property_WILDCARD_SUBSCRIPTION_AVAILABLE
            [simple uint 8                   value                              ]
        ]
        ['SUBSCRIPTION_IDENTIFIER_AVAILABLE' MQTT_Property_SUBSCRIPTION_IDENTIFIER_AVAILABLE
            [simple uint 8                   value                              ]
        ]
        ['SHARED_SUBSCRIPTION_AVAILABLE'     MQTT_Property_SHARED_SUBSCRIPTION_AVAILABLE
            [simple uint 8                   value                              ]
        ]
    ]
]

[type Filter
    [simple   MQTT_String         filter        ]
    // Subscription Options Start
    [reserved uint 2              '0x0'         ]
    // Don't forward to clients with a client id equal to the sending one
    [simple   MQTT_RetainHandling retainHandling]
    [simple   bit                 retain        ]
    [simple   bit                 noLocal       ]
    [simple   MQTT_QOS            maxQos        ]
    // Subscription Options End
]

[type MQTT_String
    [implicit uint 16                stringLength 'STR_LEN(value)']
    [simple   vstring 'stringLength * 8' value                        ]
]

[enum uint 4 MQTT_ControlPacketType
    ['0x0' RESERVED   ]
    ['0x1' CONNECT    ]
    ['0x2' CONNACK    ]
    ['0x3' PUBLISH    ]
    ['0x4' PUBACK     ]
    ['0x5' PUBREC     ]
    ['0x6' PUBREL     ]
    ['0x7' PUBCOMP    ]
    ['0x8' SUBSCRIBE  ]
    ['0x9' SUBACK     ]
    ['0xA' UNSUBSCRIBE]
    ['0xB' UNSUBACK   ]
    ['0xC' PINGREQ    ]
    ['0xD' PINGRESP   ]
    ['0xE' DISCONNECT ]
    ['0xF' AUTH       ]
]

[enum uint 8 MQTT_PropertyType
    ['0x01' PAYLOAD_FORMAT_INDICATOR         ]
    ['0x02' MESSAGE_EXPIRY_INTERVAL          ]
    ['0x03' CONTENT_TYPE                     ]
    ['0x08' RESPONSE_TOPIC                   ]
    ['0x09' CORRELATION_DATA                 ]
    ['0x0B' SUBSCRIPTION_IDENTIFIER          ]
    ['0x11' SESSION_EXPIRY_INTERVAL          ]
    ['0x12' ASSIGNED_CLIENT_IDENTIFIER       ]
    ['0x13' SERVER_KEEP_ALIVE                ]
    ['0x15' AUTHENTICATION_METHOD            ]
    ['0x16' AUTHENTICATION_DATA              ]
    ['0x17' REQUEST_PROBLEM_INFORMATION      ]
    ['0x18' WILL_DELAY_INTERVAL              ]
    ['0x19' REQUEST_RESPONSE_INFORMATION     ]
    ['0x1A' RESPONSE_INFORMATION             ]
    ['0x1C' SERVER_REFERENCE                 ]
    ['0x1F' REASON_STRING                    ]
    ['0x21' RECEIVE_MAXIMUM                  ]
    ['0x22' TOPIC_ALIAS_MAXIMUM              ]
    ['0x23' TOPIC_ALIAS                      ]
    ['0x24' MAXIMUM_QOS                      ]
    ['0x25' RETAIN_AVAILABLE                 ]
    ['0x26' USER_PROPERTY                    ]
    ['0x27' MAXIMUM_PACKET_SIZE              ]
    ['0x28' WILDCARD_SUBSCRIPTION_AVAILABLE  ]
    ['0x29' SUBSCRIPTION_IDENTIFIER_AVAILABLE]
    ['0x2A' SHARED_SUBSCRIPTION_AVAILABLE    ]
]

[enum uint 8 MQTT_ReasonCode                       (bit connackResponse, bit pubackPubrecResponse, bit pubrelPubcompResponse, bit subackResponse, bit unsubackResponse, bit disconnectReason, bit authReason)
    ['0X00' SUCCESS                                ['true'             , 'true'                  , 'true'                   , 'true'            , 'true'              , 'true'              , 'true'       ]]
    ['0X01' GRANTED_QOS_1                          ['false'            , 'false'                 , 'false'                  , 'true'            , 'false'             , 'false'             , 'false'      ]]
    ['0X02' GRANTED_QOS_2                          ['false'            , 'false'                 , 'false'                  , 'true'            , 'false'             , 'false'             , 'false'      ]]
    ['0X04' DISCONNECT_WITH_WILL_MESSAGE           ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0x10' NO_MATCHING_SUBSCRIBERS                ['false'            , 'true'                  , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0x11' NO_SUBSCRIPTION_EXISTED                ['false'            , 'false'                 , 'false'                  , 'false'           , 'true'              , 'false'             , 'false'      ]]
    ['0X18' CONTINUE_AUTHENTICATION                ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'true'       ]]
    ['0X19' RE_AUTHENTICATE                        ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'true'       ]]
    ['0X80' UNSPECIFIED_ERROR                      ['true'             , 'true'                  , 'false'                  , 'true'            , 'true'              , 'true'              , 'false'      ]]
    ['0X81' MALFORMED_PACKET                       ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X82' PROTOCOL_ERROR                         ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X83' IMPLEMENTATION_SPECIFIC_ERROR          ['true'             , 'true'                  , 'false'                  , 'true'            , 'true'              , 'true'              , 'false'      ]]
    ['0X84' UNSUPPORTED_PROTOCOL_VERSION           ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X85' CLIENT_IDENTIFIER_NOT_VALID            ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X86' BAD_USER_NAME_OR_PASSWORD              ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X87' NOT_AUTHORIZED                         ['true'             , 'true'                  , 'false'                  , 'true'            , 'true'              , 'true'              , 'false'      ]]
    ['0X88' SERVER_UNAVAILABLE                     ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X89' SERVER_BUSY                            ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X8A' BANNED                                 ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X8B' SERVER_SHUTTING_DOWN                   ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X8C' BAD_AUTHENTICATION_METHOD              ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X8D' KEEP_ALIVE_TIMEOUT                     ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X8E' SESSION_TAKEN_OVER                     ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X8F' TOPIC_FILTER_INVALID                   ['false'            , 'false'                 , 'false'                  , 'true'            , 'true'              , 'true'              , 'false'      ]]
    ['0X90' TOPIC_NAME_INVALID                     ['true'             , 'true'                  , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X91' PACKET_IDENTIFIER_IN_USE               ['false'            , 'true'                  , 'false'                  , 'true'            , 'true'              , 'false'             , 'false'      ]]
    ['0X92' PACKET_IDENTIFIER_NOT_FOUND            ['false'            , 'false'                 , 'true'                   , 'false'           , 'false'             , 'false'             , 'false'      ]]
    ['0X93' RECEIVE_MAXIMUM_EXCEEDED               ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X94' TOPIC_ALIAS_INVALID                    ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X95' PACKET_TOO_LARGE                       ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X96' MESSAGE_RATE_TO_HIGH                   ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X97' QUOTA_EXCEEDED                         ['true'             , 'true'                  , 'false'                  , 'true'            , 'false'             , 'true'              , 'false'      ]]
    ['0X98' ADMINISTRATIVE_ACTION                  ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X99' PAYLOAD_FORMAT_INVALID                 ['true'             , 'true'                  , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X9A' RETAIN_NOT_SUPPORTED                   ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X9B' QOS_NOT_SUPPORTED                      ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X9C' USE_ANOTHER_SERVER                     ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X9D' SERVER_MOVED                           ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0X9E' SHARED_SUBSCRIPTIONS_NOT_SUPPORTED     ['false'            , 'false'                 , 'false'                  , 'true'            , 'false'             , 'true'              , 'false'      ]]
    ['0X9F' CONNECTION_RATE_EXCEEDED               ['true'             , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0XA0' MAXIMUM_CONNECT_TIME                   ['false'            , 'false'                 , 'false'                  , 'false'           , 'false'             , 'true'              , 'false'      ]]
    ['0XA1' SUBSCRIPTION_IDENTIFIERS_NOT_SUPPORTED ['false'            , 'false'                 , 'false'                  , 'true'            , 'false'             , 'true'              , 'false'      ]]
    ['0XA2' WILDCARD_SUBSCRIPTIONS_NOT_SUPPORTED   ['false'            , 'false'                 , 'false'                  , 'true'            , 'false'             , 'true'              , 'false'      ]]
]

[enum uint 2 MQTT_QOS
    ['0x0' AT_MOST_ONCE ]
    ['0x1' AT_LEAST_ONCE]
    ['0x2' EXACTLY_ONCE ]
]

[enum uint 2 MQTT_RetainHandling
    ['0x0' SEND_RETAINED_MESSAGES_AT_THE_TIME_OF_THE_SUBSCRIBE                                  ]
    ['0x1' SEND_RETAINED_MESSAGES_AT_SUBSCRIBE_ONLY_IF_THE_SUBSCRIPTION_DOES_NOT_CURRENTLY_EXIST]
    ['0x2' DO_NOT_SEND_RETAINED_MESSAGES_AT_THE_TIME_OF_SUBSCRIBE                               ]
]
