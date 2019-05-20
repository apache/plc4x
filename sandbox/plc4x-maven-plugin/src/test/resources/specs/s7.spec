////////////////////////////////////////////////////////////////
// IsoOnTcp/TPKT
////////////////////////////////////////////////////////////////

[type TPKTPacket
    [const    uint8  'protocolId' '0x3']
    [reserved uint8  '0x0']
    [implicit uint16 'length' 'payload.size + 4']
    [embedded 'payload' {length: 'length - 4'}]
]

////////////////////////////////////////////////////////////////
// COTP
////////////////////////////////////////////////////////////////

[discriminatedType COTPPacket [length]
    [implicit uint8 'headerLength' 'this.size - (payload.size + 1)']
    [discriminator uint8 'tpduCode']
    [typeSwitch 'tpduCode'
        ['0xF0' Data
            [field bit 'eot']
            [field uint7 'tpduRef']
        ]
        ['0xE0' ConnectionRequest
            [field uint16 'destinationReference']
            [field uint16 'sourceReference']
            [field uint8  'protocolClass']
        ]
        ['0xD0' ConnectionResponse
            [field uint16 'destinationReference']
            [field uint16 'sourceReference']
            [field uint8  'protocolClass']
        ]
        ['0x80' DisconnectRequest
            [field uint16 'destinationReference']
            [field uint16 'sourceReference']
            [field unit8  'protocolClass']
        ]
        ['0xC0' DisconnectResponse
            [field uint16 'destinationReference']
            [field uint16 'sourceReference']
        ]
        ['0x70' TpduError
            [field uint16 'destinationReference']
            [field uint8  'rejectCause']
        ]
    ]
    [arrayField COTPParameter 'parameters' length '(headerLength + 1) - cur']
    [embedded 'payload' {length: 'length - (headerLength + 1)'}]
]

[discriminatedType COTPParameter
    [discriminator uint8 'type']
    [typeSwitch 'type'
        ['0xC0' TpduSize
            [field uint8 'tpduSize']
        ]
        ['0xC1' CallingTsap
            [field uint16 'tsapId']
        ]
        ['0xC2' CalledTsap
            [field uint16 'tsapId']
        ]
        ['0xC3' Checksum
            [field uint8 'checksum']
        ]
        ['0xE0' DisconnectAdditionalInformation
            [arrayField uint8 'data' count 'rest']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// S7
////////////////////////////////////////////////////////////////

[discriminatedType S7Message
    [const    uint8  'protocolId' '0x32']
    [discriminator uint8 'type']
    [reserved uint16 '0x00']
    [field    uint16 'tpduReference']
    [implicit uint16 'parameterLength' 'parameter.size']
    [implicit uint16 'payloadLength' 'payload.size']
    [typeSwitch 'type'
        ['0x01' Request
            [context 'messageType' 'request']
        ]
        ['0x03' Response
            [context 'messageType' 'response']
            [field uint8 'errorClass']
            [field uint8 'errorCode']
        ]
        ['0x07' UserData
            [context 'messageType' 'userData']
        ]
    ]
    [optionalField S7Parameter 'parameter' 'parameterLength > 0' {messageType: messageType}]
    [optionalField S7Payload 'payload' 'payloadLength > 0' {messageType: messageType, parameter: parameter}]
]

////////////////////////////////////////////////////////////////
// Parameters

[discriminatedType S7Parameter [messageType]
    [discriminator uint8 'type']
    [typeSwitch 'type,messageType'
        ['0xF0' SetupCommunication
            [reserved uint8 '0x0']
            [field uint16 'maxAmqCaller']
            [field uint16 'maxAmqCallee']
            [field uint16 'pduLength']
        ]
        ['0x04','request' ReadVarRequest
            [implicit uint8 'numItems' 'items.count']
            [arrayField S7VarRequestParameterItem 'items' count 'numItems']
        ]
        ['0x04','response' ReadVarResponse
            [field uint8 'numItems']
        ]
        ['0x05','request' WriteVarRequest
            [implicit uint8 'numItems' 'items.count']
            [arrayField S7VarRequestParameterItem 'items' count 'numItems']
        ]
        ['0x05','response' WriteVarResponse
            [field uint8 'numItems']
        ]
        ['0x00','userData' UserData
            [implicit uint8 'items.count']
            [arrayField UserDataItem 'items' count 'numItems']
        ]
    ]
]

[discriminatedType S7VarRequestParameterItem
    [discriminator uint8 'type']
    [typeSwitch 'type'
        ['0x12' Address
            [implicit uint8 'addressLength' 'address.size']
            [field S7Address 'address']
        ]
    ]
]

[discriminatedType S7Address
    [discriminator uint8 'type']
    [typeSwitch 'type'
        ['0x10' Any
            [field uint8  'transportSize']
            [field uint16 'length']
            [field uint8  'dbNumber']
            [field uint8  'area']
            [reserved uint5]
            [field uint16 'byteAddress']
            [field uint3  'bitAddress']
        ]
    ]
]

[discriminatorType UserDataItem
    [discriminator uint8 'type']
    [typeSwitch type
        ['0x12' CPUFunctions
            [internal      uint8 'parameterLength']
            [field         uint16 'type']
            [field         unit8 'subFunctionGroup']
            [field         uint8 'sequenceNumber']
            [optionalField uint8 'dataUnitReferenceNumber' 'parameterLength == 8']
            [optionalField uint8 'lastDataUnit' 'parameterLength == 8']
            [optionalField uint8 'errorCode' 'parameterLength == 8']

            // TODO: CPUFunctions still need some love ...

        ]
    ]
]

////////////////////////////////////////////////////////////////
// Payloads

[discriminatorType S7Payload [messageType, parameter]
    [typeSwitch 'parameter.type, messageType'
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

[type S7VarPayloadDataItem
    [field uint8 'returnCode']
    [field uint8 'transportSize']
    [field uint16 'length']
    [arrayField uint8 'data' count 'length']
]

[type S7VarPayloadStatusItem
    [field uint8 'returnCode']
]