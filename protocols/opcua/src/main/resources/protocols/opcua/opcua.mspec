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

// Remark: The different fields are encoded in Little-endian.

[type 'OpcuaAPU' [bit 'response']
    [simple         MessagePDU   'message' ['response']]
]

[discriminatedType 'MessagePDU' [bit 'response']
    [discriminator string '24'          'messageType']
    [typeSwitch 'messageType','response'
        ['HEL','false'     OpcuaHelloRequest
            [simple          string '1 * 8'            'chunk']
            [implicit          int 32             'messageSize' 'lengthInBytes']
            [simple          int 32             'version']
            [simple          int 32             'receiveBufferSize']
            [simple          int 32             'sendBufferSize']
            [simple          int 32             'maxMessageSize']
            [simple          int 32             'maxChunkCount']
            [simple          int 32             'stringLength']
            [simple          string 'stringLength * 8' 'endpoint']
        ]
        ['HEL','true'     OpcuaHelloResponse
        ]
        ['ACK','false'     OpcuaAcknowledgeRequest
        ]
        ['ACK','true'     OpcuaAcknowledgeResponse
            [simple          string '8'            'chunk']
            [implicit          int 32             'messageSize' 'lengthInBytes']
            [simple          int 32             'version']
            [simple          int 32             'receiveBufferSize']
            [simple          int 32             'sendBufferSize']
            [simple          int 32             'maxMessageSize']
            [simple          int 32             'maxChunkCount']
        ]
        ['ERR','false'     OpcuaErrorRequest
        ]
        ['ERR','true'     OpcuaErrorResponse
        ]
        ['OPN','false'     OpcuaOpenRequest
            [simple          string '8'         'chunk']
            [implicit        int 32             'messageSize' 'lengthInBytes']
            [simple          int 32             'secureChannelId']
            [simple          PascalString       'endpoint']
            [simple          PascalByteString   'senderCertificate']
            [simple          PascalByteString   'receiverCertificateThumbprint']
            [simple          int 32             'sequenceNumber']
            [simple          int 32             'requestId']
            [array           int 8              'message' count 'messageSize - endpoint.stringLength - senderCertificate.stringLength - receiverCertificateThumbprint.stringLength - 32']
       ]
       ['OPN','true'     OpcuaOpenResponse
           [simple          string '8'         'chunk']
           [implicit          int 32            'messageSize' 'lengthInBytes']
           [simple          int 32             'secureChannelId']
           [simple          int 32             'securityPolicyUriSize']
           [simple          string 'securityPolicyUriSize == -1 ? 0 : securityPolicyUriSize * 8'          'endpoint']
           [simple          int 32             'senderCertificateSize']
           [simple          string 'senderCertificateSize == -1 ? 0 : senderCertificateSize * 8'             'senderCertificate']
           [simple          int 32             'receiverCertificateThumbprintSize']
           [simple          string 'receiverCertificateThumbprintSize == -1 ? 0 : receiverCertificateThumbprintSize * 8'             'receiverCertificateThumbprint']
           [simple          int 32             'sequenceNumber']
           [simple          int 32             'requestId']
           [array           int 8              'message' count 'messageSize - (receiverCertificateThumbprintSize == -1 ? 0 : receiverCertificateThumbprintSize) - (senderCertificateSize == -1 ? 0 : senderCertificateSize) - (securityPolicyUriSize == -1 ? 0 : securityPolicyUriSize) - 32)']
       ]
       ['CLO','false'     OpcuaCloseRequest
           [simple          string '8'         'chunk']
           [implicit        int 32             'messageSize' 'lengthInBytes']
           [simple          int 32             'secureChannelId']
           [simple          int 32             'secureTokenId']
           [simple          int 32             'sequenceNumber']
           [simple          int 32             'requestId']
           [simple          OpcuaMessage       'message']
       ]
       ['CLO','true'     OpcuaCloseResponse
       ]
       ['MSG','false'     OpcuaMessageRequest
           [simple          string '8'         'chunk']
           [implicit        int 32             'messageSize' 'lengthInBytes']
           [simple          int 32             'secureChannelId']
           [simple          int 32             'secureTokenId']
           [simple          int 32             'sequenceNumber']
           [simple          int 32             'requestId']
           [array           int 8              'message' count 'messageSize - 24']
       ]
       ['MSG','true'     OpcuaMessageResponse
           [simple          string '8'         'chunk']
           [implicit        int 32             'messageSize' 'lengthInBytes']
           [simple          int 32             'secureChannelId']
           [simple          int 32             'secureTokenId']
           [simple          int 32             'sequenceNumber']
           [simple          int 32             'requestId']
           [array           int 8              'message' count 'messageSize - 24']
       ]
    ]
]


[enum string '-1' 'OpcuaDataType' [uint 8 'variantType']
    ['IEC61131_NULL' NULL ['0']]
    ['IEC61131_BOOL' BOOL ['1']]
    ['IEC61131_BYTE' BYTE ['3']]
    ['IEC61131_SINT' SINT ['2']]
    ['IEC61131_INT' INT ['4']]
    ['IEC61131_DINT' DINT ['6']]
    ['IEC61131_LINT' LINT ['8']]
    ['IEC61131_USINT' USINT ['3']]
    ['IEC61131_UINT' UINT ['5']]
    ['IEC61131_UDINT' UDINT ['7']]
    ['IEC61131_ULINT' ULINT ['9']]
    ['IEC61131_REAL' REAL ['10']]
    ['IEC61131_LREAL' LREAL ['11']]
    ['IEC61131_TIME' TIME ['1']]
    ['IEC61131_LTIME' LTIME ['1']]
    ['IEC61131_DATE' DATE ['1']]
    ['IEC61131_LDATE' LDATE ['1']]
    ['IEC61131_TIME_OF_DAY' TIME_OF_DAY ['1']]
    ['IEC61131_LTIME_OF_DAY' LTIME_OF_DAY ['1']]
    ['IEC61131_DATE_AND_TIME' DATE_AND_TIME ['13']]
    ['IEC61131_LDATE_AND_TIME' LDATE_AND_TIME ['1']]
    ['IEC61131_CHAR' CHAR ['1']]
    ['IEC61131_WCHAR' WCHAR ['1']]
    ['IEC61131_STRING' STRING ['12']]
]

[enum string '-1' 'OpcuaIdentifierType'
    ['s' STRING_IDENTIFIER]
    ['i' NUMBER_IDENTIFIER]
    ['g' GUID_IDENTIFIER]
    ['b' BINARY_IDENTIFIER]
]
