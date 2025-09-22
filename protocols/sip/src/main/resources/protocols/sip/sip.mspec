[type Constants
    [const byte SPACE 0x20]
    [const byte COLON 0x3a]
    [const byte SLASH 0x2f]
    [const byte R     0x0d]
    [const byte N     0x0a]
]

[type SipPDU (uint 16 len) byteOrder='BIG_ENDIAN'
    [simple SipRequestLine requestLine]
    [array Header headers length 'len - requestLine.lengthInBytes - 2']
    //[array byte data count 'len - requestLine.lengthInBytes - 2']
    [reserved byte 8 'Constants.R']
    [reserved byte 8 'Constants.N']
]

[type Header
    [manual vstring header      'STATIC_CALL("readStringTill", readBuffer, ":")'      'STATIC_CALL("writeStringTill", writeBuffer, header)'     '8 * STR_LEN(header)'     ]
    [reserved byte 8 'Constants.COLON']
    [reserved byte 8 'Constants.SPACE']
    [manual vstring value       'STATIC_CALL("readStringTill", readBuffer, "\r\n")'   'STATIC_CALL("writeStringTill", writeBuffer, value)'      '8 * STR_LEN(value)'      ]
    [const byte R     0x0d]
    [const byte N     0x0a]
]

[type SipRequestLine
    [manual vstring method      'STATIC_CALL("readStringTill", readBuffer, " ")'      'STATIC_CALL("writeStringTill", writeBuffer, method)'     '8 * STR_LEN(method)'     ]
    [reserved byte 8 'Constants.SPACE']
    [manual vstring proto       'STATIC_CALL("readStringTill", readBuffer, ":")'      'STATIC_CALL("writeStringTill", writeBuffer, proto)'      '8 * STR_LEN(proto)'      ]
    [reserved byte 8 'Constants.COLON']
    [manual vstring requestUri  'STATIC_CALL("readStringTill", readBuffer, " ")'      'STATIC_CALL("writeStringTill", writeBuffer, requestUri)' '8 * STR_LEN(requestUri)' ]
    [reserved byte 8 'Constants.SPACE']
    [manual vstring protocol    'STATIC_CALL("readStringTill", readBuffer, "/")'      'STATIC_CALL("writeStringTill", writeBuffer, protocol)'   '8 * STR_LEN(protocol)'   ]
    [reserved byte 8 'Constants.SLASH']
    [manual vstring version     'STATIC_CALL("readStringTill", readBuffer, "\r\n")'   'STATIC_CALL("writeStringTill", writeBuffer, version)'    '8 * STR_LEN(version)'    ]
    [reserved byte 8 'Constants.R']
    [reserved byte 8 'Constants.N']
]

