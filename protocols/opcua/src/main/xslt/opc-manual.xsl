<?xml version="1.0" encoding="UTF-8"?>
<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->
<xsl:stylesheet version="2.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:opc="http://opcfoundation.org/BinarySchema/"
                xmlns:ua="http://opcfoundation.org/UA/"
                xmlns:node="http://opcfoundation.org/UA/2011/03/UANodeSet.xsd">

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />

    <xsl:import href="opc-common.xsl"/>

    <xsl:variable name="originaldoc" select="/"/>

    <xsl:param name="osType"/>

    <xsl:param name="services"/>

    <!-- Intermediate, to reformat the url on windows systems -->
    <xsl:param name="servicesUrl">
        <xsl:choose>
            <xsl:when test="$osType = 'win'">file:///<xsl:value-of select="replace($services, '\\', '/')"/></xsl:when>
            <xsl:otherwise><xsl:value-of select="$services"/></xsl:otherwise>
        </xsl:choose>
    </xsl:param>

    <xsl:param name="file" select="document($servicesUrl)"/>

    <xsl:template match="/">
// Remark: The different fields are encoded in Little-endian.

[type OpcuaAPU(bit response) byteOrder='LITTLE_ENDIAN'
    [simple MessagePDU('response') message]
]

[discriminatedType MessagePDU(bit response)
    [discriminator string 24            messageType]
    [typeSwitch messageType,response
        ['"HEL"','false'     OpcuaHelloRequest
            [simple          string 8           chunk]
            [implicit        int 32             messageSize 'lengthInBytes']
            [simple          int 32             version]
            [simple          int 32             receiveBufferSize]
            [simple          int 32             sendBufferSize]
            [simple          int 32             maxMessageSize]
            [simple          int 32             maxChunkCount]
            [simple          PascalString       endpoint]
        ]
        ['"ACK"','true'     OpcuaAcknowledgeResponse
            [simple          string 8           chunk]
            [implicit        int 32             messageSize 'lengthInBytes']
            [simple          int 32             version]
            [simple          int 32             receiveBufferSize]
            [simple          int 32             sendBufferSize]
            [simple          int 32             maxMessageSize]
            [simple          int 32             maxChunkCount]
        ]
        ['"OPN"','false'     OpcuaOpenRequest
            [simple          string 8           chunk]
            [implicit        int 32             messageSize 'lengthInBytes']
            [simple          int 32             secureChannelId]
            [simple          PascalString       endpoint]
            [simple          PascalByteString   senderCertificate]
            [simple          PascalByteString   receiverCertificateThumbprint]
            [simple          int 32             sequenceNumber]
            [simple          int 32             requestId]
            [array           byte               message count 'messageSize - (endpoint.stringLength == -1 ? 0 : endpoint.stringLength ) - (senderCertificate.stringLength == -1 ? 0 : senderCertificate.stringLength) - (receiverCertificateThumbprint.stringLength == -1 ? 0 : receiverCertificateThumbprint.stringLength) - 32']
       ]
       ['"OPN"','true'     OpcuaOpenResponse
           [simple          string 8           chunk]
           [implicit        int 32             messageSize 'lengthInBytes']
           [simple          int 32             secureChannelId]
           [simple          PascalString       securityPolicyUri]
           [simple          PascalByteString   senderCertificate]
           [simple          PascalByteString   receiverCertificateThumbprint]
           [simple          int 32             sequenceNumber]
           [simple          int 32             requestId]
           [array           byte               message count 'messageSize - (securityPolicyUri.stringLength == -1 ? 0 : securityPolicyUri.stringLength) - (senderCertificate.stringLength == -1 ? 0 : senderCertificate.stringLength) - (receiverCertificateThumbprint.stringLength == -1 ? 0 : receiverCertificateThumbprint.stringLength) - 32']
       ]
       ['"CLO"','false'     OpcuaCloseRequest
           [simple          string 8           chunk]
           [implicit        int 32             messageSize 'lengthInBytes']
           [simple          int 32             secureChannelId]
           [simple          int 32             secureTokenId]
           [simple          int 32             sequenceNumber]
           [simple          int 32             requestId]
           [simple          ExtensionObject('false')       message]
       ]
       ['"MSG"','false'     OpcuaMessageRequest
           [simple          string 8           chunk]
           [implicit        int 32             messageSize 'lengthInBytes']
           [simple          int 32             secureChannelId]
           [simple          int 32             secureTokenId]
           [simple          int 32             sequenceNumber]
           [simple          int 32             requestId]
           [array           byte               message count 'messageSize - 24']
       ]
       ['"MSG"','true'     OpcuaMessageResponse
           [simple          string 8           chunk]
           [implicit        int 32             messageSize 'lengthInBytes']
           [simple          int 32             secureChannelId]
           [simple          int 32             secureTokenId]
           [simple          int 32             sequenceNumber]
           [simple          int 32             requestId]
           [array           byte               message count 'messageSize - 24']
       ]
    ]
]

[type ByteStringArray
    [simple int 32 arrayLength]
    [array uint 8 value count 'arrayLength']
]

[type GuidValue
    [simple uint 32 data1]
    [simple uint 16 data2]
    [simple uint 16 data3]
    [array  byte    data4 count '2']
    [array  byte    data5 count '6']
]

[type ExpandedNodeId
    [simple bit namespaceURISpecified]
    [simple bit serverIndexSpecified]
    [simple NodeIdTypeDefinition nodeId]
    [virtual vstring '-1' identifier 'nodeId.identifier']
    [optional PascalString namespaceURI 'namespaceURISpecified']
    [optional uint 32 serverIndex 'serverIndexSpecified']
]

[type ExtensionHeader
    [reserved int 5 '0x00']
    [simple bit xmlbody]
    [simple bit binaryBody]
]

[type ExtensionObjectEncodingMask
    [reserved int 5 '0x00']
    [simple bit typeIdSpecified]
    [simple bit xmlbody]
    [simple bit binaryBody]
]

[type ExtensionObject(bit includeEncodingMask)
    //A serialized object prefixed with its data type identifier.
    [simple ExpandedNodeId typeId]
    [optional ExtensionObjectEncodingMask encodingMask 'includeEncodingMask']
    [virtual vstring '-1' identifier 'typeId.identifier']
    [simple ExtensionObjectDefinition('identifier') body]
]

[discriminatedType ExtensionObjectDefinition(vstring '-1' identifier)
    [typeSwitch identifier
        ['"0"' NullExtension
        ]

        <xsl:for-each select="/opc:TypeDictionary/opc:StructuredType[(@BaseType = 'ua:ExtensionObject') and not(@Name = 'UserIdentityToken') and not(@Name = 'PublishedDataSetDataType') and not(@Name = 'DataSetReaderDataType')]">
            <xsl:message><xsl:value-of select="@Name"/></xsl:message>
            <xsl:variable name="extensionName" select="@Name"/>
            <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName=$extensionName]"/>
        </xsl:for-each>

        ['"811"' DataChangeNotification
            [implicit int 32 notificationLength 'lengthInBytes']
            [simple int 32 noOfMonitoredItems]
            [array ExtensionObjectDefinition('"808"')  monitoredItems count 'noOfMonitoredItems']
            [simple int 32 noOfDiagnosticInfos]
            [array DiagnosticInfo  diagnosticInfos count 'noOfDiagnosticInfos']
        ]
        ['"916"' EventNotificationList
            [implicit int 32 notificationLength 'lengthInBytes']
            [simple int 32 noOfEvents]
            [array ExtensionObjectDefinition('"919"')  events count 'noOfEvents']
        ]
        ['"820"' StatusChangeNotification
            [implicit int 32 notificationLength 'lengthInBytes']
            [simple StatusCode status]
            [simple DiagnosticInfo diagnosticInfo]
        ]

        ['"316"' UserIdentityToken
            [implicit int 32 policyLength 'policyId.lengthInBytes  + userIdentityTokenDefinition.lengthInBytes']
            [simple PascalString policyId]
            [simple UserIdentityTokenDefinition('policyId.stringValue') userIdentityTokenDefinition]
        ]
    ]
]

[discriminatedType UserIdentityTokenDefinition(vstring '-1' identifier)
    [typeSwitch identifier
        ['"anonymous"' AnonymousIdentityToken
        ]
        ['"username"' UserNameIdentityToken
            [simple PascalString userName]
            [simple PascalByteString password]
            [simple PascalString encryptionAlgorithm]
        ]
        ['"certificate"' X509IdentityToken
            [simple PascalByteString certificateData]
        ]
        ['"identity"' IssuedIdentityToken
            [simple PascalByteString tokenData]
            [simple PascalString encryptionAlgorithm]
        ]
    ]
]


[discriminatedType Variant
    [simple bit arrayLengthSpecified]
    [simple bit arrayDimensionsSpecified]
    [discriminator uint 6 VariantType]
    [typeSwitch VariantType,arrayLengthSpecified
        ['1' VariantBoolean (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array byte value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['2' VariantSByte (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array byte value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['3' VariantByte (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array uint 8 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['4' VariantInt16 (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array int 16 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['5' VariantUInt16 (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array uint 16 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['6' VariantInt32 (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array int 32 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['7' VariantUInt32 (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array uint 32 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['8' VariantInt64 (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array int 64 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['9' VariantUInt64 (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array uint 64 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['10' VariantFloat (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array float 32 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['11' VariantDouble (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array float 64 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['12' VariantString (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array PascalString value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['13' VariantDateTime (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array int 64 value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['14' VariantGuid (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array GuidValue value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['15' VariantByteString (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array ByteStringArray value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['16' VariantXmlElement (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array PascalString value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['17' VariantNodeId (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array NodeId value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['18' VariantExpandedNodeId (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array ExpandedNodeId value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['19' VariantStatusCode (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array StatusCode value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['20' VariantQualifiedName (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array QualifiedName value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['21' VariantLocalizedText (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array LocalizedText value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['22' VariantExtensionObject (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array ExtensionObject('true') value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['23' VariantDataValue (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array DataValue value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['24' VariantVariant (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array Variant value count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['25' VariantDiagnosticInfo (bit arrayLengthSpecified)
            [optional int 32 arrayLength 'arrayLengthSpecified']
            [array DiagnosticInfo value count 'arrayLength == null ? 1 : arrayLength']
        ]
    ]
    [optional int 32 noOfArrayDimensions 'arrayDimensionsSpecified']
    [array bit arrayDimensions count 'noOfArrayDimensions == null ? 0 : noOfArrayDimensions']
]

[discriminatedType NodeIdTypeDefinition
    [abstract vstring '-1' identifier]
    [discriminator NodeIdType nodeType]
    [typeSwitch nodeType
        ['nodeIdTypeTwoByte' NodeIdTwoByte
            [simple uint 8 id]
            [virtual vstring '-1' identifier 'id']
        ]
        ['nodeIdTypeFourByte' NodeIdFourByte
            [simple uint 8 namespaceIndex]
            [simple uint 16 id]
            [virtual vstring '-1' identifier 'id']
        ]
        ['nodeIdTypeNumeric' NodeIdNumeric
            [simple uint 16 namespaceIndex]
            [simple uint 32 id]
            [virtual vstring '-1' identifier 'id']
        ]
        ['nodeIdTypeString' NodeIdString
            [simple uint 16 namespaceIndex]
            [simple PascalString id]
            [virtual vstring '-1' identifier 'id.stringValue']
        ]
        ['nodeIdTypeGuid' NodeIdGuid
            [simple uint 16 namespaceIndex]
            [array byte id count '16']
            [virtual vstring '-1' identifier 'id']
        ]
        ['nodeIdTypeByteString' NodeIdByteString
            [simple uint 16 namespaceIndex]
            [simple PascalByteString id]
            [virtual vstring '-1' identifier 'id.stringValue']
        ]
    ]
]

[type NodeId
    [reserved int 2 '0x00']
    [simple NodeIdTypeDefinition nodeId]
    [virtual vstring '-1' id 'nodeId.identifier']
]

[type PascalString
    [implicit int 32 sLength          'stringValue.length == 0 ? -1 : stringValue.length']
    [simple vstring 'sLength == -1 ? 0 : sLength * 8' stringValue]
    [virtual  int 32 stringLength     'stringValue.length == -1 ? 0 : stringValue.length']
]

[type PascalByteString
    [simple int 32 stringLength]
    [array byte stringValue count 'stringLength == -1 ? 0 : stringLength' ]
]

[type Structure

]

//[type DataTypeDefinition
//
//]

// StructuredTypes
<xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[(@Name != 'ExtensionObject') and (@Name != 'Variant') and (@Name != 'NodeId') and (@Name != 'ExpandedNodeId') and not(@BaseType)]"/>

// EnumeratedTypes
<xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType"/>

// OpaqueType
<xsl:apply-templates select="/opc:TypeDictionary/opc:OpaqueType"/>

[enum string 112 OpcuaDataType(uint 8 variantType)
    ['"NULL"' NULL ['0']]
    ['"BOOL"' BOOL ['1']]
    ['"BYTE"' BYTE ['3']]
    ['"SINT"' SINT ['2']]
    ['"INT"' INT ['4']]
    ['"DINT"' DINT ['6']]
    ['"LINT"' LINT ['8']]
    ['"USINT"' USINT ['3']]
    ['"UINT"' UINT ['5']]
    ['"UDINT"' UDINT ['7']]
    ['"ULINT"' ULINT ['9']]
    ['"REAL"' REAL ['10']]
    ['"LREAL"' LREAL ['11']]
    ['"TIME"' TIME ['1']]
    ['"LTIME"' LTIME ['1']]
    ['"DATE"' DATE ['1']]
    ['"LDATE"' LDATE ['1']]
    ['"TIME_OF_DAY"' TIME_OF_DAY ['1']]
    ['"LTIME_OF_DAY"' LTIME_OF_DAY ['1']]
    ['"DATE_AND_TIME"' DATE_AND_TIME ['13']]
    ['"LDATE_AND_TIME"' LDATE_AND_TIME ['1']]
    ['"CHAR"' CHAR ['1']]
    ['"WCHAR"' WCHAR ['1']]
    ['"STRING"' STRING ['12']]
]

[enum string 8 OpcuaIdentifierType
    ['"s"' STRING_IDENTIFIER]
    ['"i"' NUMBER_IDENTIFIER]
    ['"g"' GUID_IDENTIFIER]
    ['"b"' BINARY_IDENTIFIER]
]


    </xsl:template>
</xsl:stylesheet>
