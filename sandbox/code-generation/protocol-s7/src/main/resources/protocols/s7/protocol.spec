////////////////////////////////////////////////////////////////
// IsoOnTcp/TPKT
////////////////////////////////////////////////////////////////

[type 'TPKTPacket'
    [const    uint 8     'protocolId' '0x3']
    [reserved uint 8     '0x00']
    [implicit uint 16    'len'        'payload.size + 4']
    [field    COTPPacket 'payload'    {payloadLength: 'len - 4'}]
]

////////////////////////////////////////////////////////////////
// COTP
////////////////////////////////////////////////////////////////

[discriminatedType 'COTPPacket' ['payloadLength']
    [implicit      uint 8 'headerLength' 'this.size - (payload.size + 1)']
    [discriminator uint 8 'tpduCode']
    [typeSwitch 'tpduCode'
        ['0xF0' Data
            [field bit    'eot']
            [field uint 7 'tpduRef']
        ]
        ['0xE0' ConnectionRequest
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
            [field uint 8  'protocolClass']
        ]
        ['0xD0' ConnectionResponse
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
            [field uint 8  'protocolClass']
        ]
        ['0x80' DisconnectRequest
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
            [field uint 8  'protocolClass']
        ]
        ['0xC0' DisconnectResponse
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
        ]
        ['0x70' TpduError
            [field uint 16 'destinationReference']
            [field uint 8  'rejectCause']
        ]
    ]
    [arrayField COTPParameter 'parameters' length '(headerLength + 1) - cur']
    [field      S7Message     'payload'    {payloadLength: 'payloadLength - (headerLength + 1)'}]
]

[discriminatedType 'COTPParameter'
    [discriminator uint 8 'parameterType']
    [typeSwitch 'parameterType'
        ['0xC0' TpduSize
            [field uint 8 'tpduSize']
        ]
        ['0xC1' CallingTsap
            [field uint 16 'tsapId']
        ]
        ['0xC2' CalledTsap
            [field uint 16 'tsapId']
        ]
        ['0xC3' Checksum
            [field uint 8 'checksum']
        ]
        ['0xE0' DisconnectAdditionalInformation
            [arrayField uint 8 'data' count 'rest']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// S7
////////////////////////////////////////////////////////////////

[discriminatedType 'S7Message' ['payloadLength']
    [const         uint 8  'protocolId'      '0x32']
    [discriminator uint 8  'messageType']
    [reserved      uint 16 '0x0000']
    [field         uint 16 'tpduReference']
    [implicit      uint 16 'parameterLength' 'parameters.size']
    [implicit      uint 16 'payloadLength'   'payloads.size']
    [typeSwitch 'messageType'
        ['0x01' Request
            [context string 'messageType' 'request']
        ]
        ['0x03' Response
            [context string 'messageType' 'response']
            [field uint 8 'errorClass']
            [field uint 8 'errorCode']
        ]
        ['0x07' UserData
            [context string 'messageType' 'userData']
        ]
    ]
    [field S7Parameter 'parameter' {messageType: 'messageType'}]
    [field S7Payload 'payload' {messageType: 'messageType', parameter: 'parameter'}]
]

////////////////////////////////////////////////////////////////
// Parameters

[discriminatedType 'S7Parameter' ['messageType']
    [discriminator uint 8 'parameterType']
    [typeSwitch 'parameterType','messageType'
        ['0xF0' SetupCommunication
            [reserved uint 8  '0x00']
            [field    uint 16 'maxAmqCaller']
            [field    uint 16 'maxAmqCallee']
            [field    uint 16 'pduLength']
        ]
        ['0x04','request' ReadVarRequest
            [implicit   uint 8                    'numItems' 'items.size']
            [arrayField S7VarRequestParameterItem 'items'    count 'numItems']
        ]
        ['0x04','response' ReadVarResponse
            [field uint 8 'numItems']
        ]
        ['0x05','request' WriteVarRequest
            [implicit   uint 8                    'numItems' 'items.size']
            [arrayField S7VarRequestParameterItem 'items'    count 'numItems']
        ]
        ['0x05','response' WriteVarResponse
            [field uint 8 'numItems']
        ]
        ['0x00','userData' UserData
            [implicit   uint 8       'numItems' 'items.size']
            [arrayField UserDataItem 'items' count 'numItems']
        ]
    ]
]

[discriminatedType 'S7VarRequestParameterItem'
    [discriminator uint 8 'parameterItemType']
    [typeSwitch 'parameterItemType'
        ['0x12' Address
            [implicit uint 8    'addressLength' 'address.size']
            [field    S7Address 'address']
        ]
    ]
]

[discriminatedType 'S7Address'
    [discriminator uint 8 'addressType']
    [typeSwitch 'addressType'
        ['0x10' Any
            [field    uint 8  'transportSize']
            [field    uint 16 'numberOfElements']
            [field    uint 8  'dbNumber']
            [field    uint 8  'area']
            [reserved uint 5  '0x00']
            [field    uint 16 'byteAddress']
            [field    uint 3  'bitAddress']
        ]
    ]
]

// TODO: CPUFunctions still need some love ...
[discriminatedType 'UserDataItem'
    [discriminator uint 8 'itemType']
    [typeSwitch 'itemType'
        ['0x12' CPUFunctions
            [implicit      uint 8  'parameterLength' 'size']
            [field         uint 16 'cpuFunctionType']
            [field         uint 8  'subFunctionGroup']
            [field         uint 8  'sequenceNumber']
            [optionalField uint 8  'dataUnitReferenceNumber' 'parameterLength == 8']
            [optionalField uint 8  'lastDataUnit' 'parameterLength == 8']
            [optionalField uint 8  'errorCode' 'parameterLength == 8']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// Payloads

[discriminatedType 'S7Payload' ['messageType', 'parameter']
    [typeSwitch 'parameter.parameterType', 'messageType'
        ['0x04','response' ReadVarResponse
            [arrayField S7VarPayloadDataItem 'items' count 'parameter.numItems']
        ]
        ['0x05','request' WriteVarRequest
            [arrayField S7VarPayloadDataItem 'items' count 'parameter.numItems']
        ]
        ['0x05','response' WriteVarResponse
            [arrayField S7VarPayloadStatusItem 'items' count 'parameter.numItems']
        ]
        ['0x00','userData' UserData
        ]
    ]
]

[type 'S7VarPayloadDataItem'
    [field      uint 8  'returnCode']
    [field      uint 8  'transportSize']
    [field      uint 16 'dataLength']
    [arrayField uint 8  'data' count 'dataLength']
]

[type 'S7VarPayloadStatusItem'
    [field uint 8 'returnCode']
]