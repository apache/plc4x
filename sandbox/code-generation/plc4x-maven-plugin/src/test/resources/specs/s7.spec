////////////////////////////////////////////////////////////////
// IsoOnTcp/TPKT
////////////////////////////////////////////////////////////////

[type 'TPKTPacket'
    [const    uint 8     'protocolId' '0x03']
    [reserved uint 8     '0x00']
    [implicit uint 16    'len'        'payload.size + 4']
    [field    COTPPacket 'payload']
]

////////////////////////////////////////////////////////////////
// COTP
////////////////////////////////////////////////////////////////

[discriminatedType 'COTPPacket'
    [implicit      uint 8 'headerLength' 'this.size - (payload.size + 1)']
    [discriminator uint 8 'tpduCode']
    [typeSwitch 'tpduCode'
        ['0xF0' COTPPacketData
            [field bit    'eot']
            [field uint 7 'tpduRef']
        ]
        ['0xE0' COTPPacketConnectionRequest
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
            [field uint 8  'protocolClass']
        ]
        ['0xD0' COTPPacketConnectionResponse
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
            [field uint 8  'protocolClass']
        ]
        ['0x80' COTPPacketDisconnectRequest
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
            [field uint 8  'protocolClass']
        ]
        ['0xC0' COTPPacketDisconnectResponse
            [field uint 16 'destinationReference']
            [field uint 16 'sourceReference']
        ]
        ['0x70' COTPPacketTpduError
            [field uint 16 'destinationReference']
            [field uint 8  'rejectCause']
        ]
    ]
    [arrayField COTPParameter 'parameters' length '(headerLength + 1) - cur']
    [field      S7Message     'payload']
]

[discriminatedType 'COTPParameter'
    [discriminator uint 8 'parameterType']
    [typeSwitch 'parameterType'
        ['0xC0' COTPParameterTpduSize
            [field uint 8 'tpduSize']
        ]
        ['0xC1' COTPParameterCallingTsap
            [field uint 16 'tsapId']
        ]
        ['0xC2' COTPParameterCalledTsap
            [field uint 16 'tsapId']
        ]
        ['0xC3' COTPParameterChecksum
            [field uint 8 'checksum']
        ]
        ['0xE0' COTPParameterDisconnectAdditionalInformation
            [arrayField uint 8 'data' count 'rest']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// S7
////////////////////////////////////////////////////////////////

[discriminatedType 'S7Message'
    [const         uint 8  'protocolId'      '0x32']
    [discriminator uint 8  'messageType']
    [reserved      uint 16 '0x0000']
    [field         uint 16 'tpduReference']
    [implicit      uint 16 'parameterLength' 'parameters.size']
    [implicit      uint 16 'payloadLength'   'payloads.size']
    [typeSwitch 'messageType'
        ['0x01' S7MessageRequest
        ]
        ['0x03' S7MessageResponse
            [field uint 8 'errorClass']
            [field uint 8 'errorCode']
        ]
        ['0x07' S7MessageUserData
        ]
    ]
    [field S7Parameter 'parameter' ['messageType']]
    [field S7Payload   'payload'   ['messageType', 'parameter']]
]

////////////////////////////////////////////////////////////////
// Parameters

[discriminatedType 'S7Parameter' [uint 8 'messageType']
    [discriminator uint 8 'parameterType']
    [typeSwitch 'parameterType','messageType'
        ['0xF0' SetupCommunication
            [reserved uint 8  '0x00']
            [field    uint 16 'maxAmqCaller']
            [field    uint 16 'maxAmqCallee']
            [field    uint 16 'pduLength']
        ]
        ['0x04','0x01' S7ParameterReadVarRequest
            [implicit   uint 8                    'numItems' 'items.size']
            [arrayField S7VarRequestParameterItem 'items'    count 'numItems']
        ]
        ['0x04','0x03' S7ParameterReadVarResponse
            [field uint 8 'numItems']
        ]
        ['0x05','0x01' S7ParameterWriteVarRequest
            [implicit   uint 8                    'numItems' 'items.size']
            [arrayField S7VarRequestParameterItem 'items'    count 'numItems']
        ]
        ['0x05','0x03' S7ParameterWriteVarResponse
            [field uint 8 'numItems']
        ]
        ['0x00','0x07' S7ParameterUserData
            [implicit   uint 8       'numItems' 'items.size']
            [arrayField UserDataItem 'items' count 'numItems']
        ]
    ]
]

[discriminatedType 'S7VarRequestParameterItem'
    [discriminator uint 8 'parameterItemType']
    [typeSwitch 'parameterItemType'
        ['0x12' S7VarRequestParameterItemAddress
            [implicit uint 8    'addressLength' 'address.size']
            [field    S7Address 'address']
        ]
    ]
]

[discriminatedType 'S7Address'
    [discriminator uint 8 'addressType']
    [typeSwitch 'addressType'
        ['0x10' S7AddressAny
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
        ['0x12' UserDataItemCPUFunctions
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

[discriminatedType 'S7Payload' [uint 8 'messageType', S7Parameter 'parameter']
    [typeSwitch 'parameter.parameterType', 'messageType'
        ['0x04','response' S7PayloadReadVarResponse
            [arrayField S7VarPayloadDataItem 'items' count 'parameter.numItems']
        ]
        ['0x05','request' S7PayloadWriteVarRequest
            [arrayField S7VarPayloadDataItem 'items' count 'parameter.numItems']
        ]
        ['0x05','response' S7PayloadWriteVarResponse
            [arrayField S7VarPayloadStatusItem 'items' count 'parameter.numItems']
        ]
        ['0x00','userData' S7PayloadUserData
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