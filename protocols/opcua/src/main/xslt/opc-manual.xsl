<?xml version="1.0" encoding="UTF-8"?>
<!--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License"); you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing,
  software distributed under the License is distributed on an
  "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
  KIND, either express or implied.  See the License for the
  specific language governing permissions and limitations
  under the License.
-->
<xsl:stylesheet version="2.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:xs="http://www.w3.org/2001/XMLSchema"
                xmlns:opc="http://opcfoundation.org/BinarySchema/"
                xmlns:plc4x="https://plc4x.apache.org/"
                xmlns:map="http://www.w3.org/2005/xpath-functions/map"
                xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                xmlns:ua="http://opcfoundation.org/UA/"
                xmlns:tns="http://opcfoundation.org/UA/"
                xmlns:node="http://opcfoundation.org/UA/2011/03/UANodeSet.xsd">

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />

    <xsl:import href="opc-common.xsl"/>

    <xsl:variable name="originaldoc" select="/"/>

    <xsl:param name="services"></xsl:param>
    <xsl:param name="file" select="document($services)"/>

    <xsl:template match="/">
// Remark: The different fields are encoded in Little-endian.

[type 'OpcuaAPU' [bit 'response']
    [simple         MessagePDU   'message' ['response']]
]

[discriminatedType 'MessagePDU' [bit 'response']
    [discriminator string '24'          'messageType']
    [typeSwitch 'messageType','response'
        ['HEL','false'     OpcuaHelloRequest
            [simple          string '8'         'chunk']
            [implicit        int 32             'messageSize' 'lengthInBytes']
            [simple          int 32             'version']
            [simple          int 32             'receiveBufferSize']
            [simple          int 32             'sendBufferSize']
            [simple          int 32             'maxMessageSize']
            [simple          int 32             'maxChunkCount']
            [simple          PascalString       'endpoint']
        ]
        ['ACK','true'     OpcuaAcknowledgeResponse
            [simple          string '8'         'chunk']
            [implicit        int 32             'messageSize' 'lengthInBytes']
            [simple          int 32             'version']
            [simple          int 32             'receiveBufferSize']
            [simple          int 32             'sendBufferSize']
            [simple          int 32             'maxMessageSize']
            [simple          int 32             'maxChunkCount']
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
            [array           int 8              'message' count 'messageSize - (endpoint.stringLength == -1 ? 0 : endpoint.stringLength ) - (senderCertificate.stringLength == -1 ? 0 : senderCertificate.stringLength) - (receiverCertificateThumbprint.stringLength == -1 ? 0 : receiverCertificateThumbprint.stringLength) - 32']
       ]
       ['OPN','true'     OpcuaOpenResponse
           [simple          string '8'         'chunk']
           [implicit        int 32             'messageSize' 'lengthInBytes']
           [simple          int 32             'secureChannelId']
           [simple          PascalString       'securityPolicyUri']
           [simple          PascalByteString   'senderCertificate']
           [simple          PascalByteString   'receiverCertificateThumbprint']
           [simple          int 32             'sequenceNumber']
           [simple          int 32             'requestId']
           [array           int 8              'message' count 'messageSize - (securityPolicyUri.stringLength == -1 ? 0 : securityPolicyUri.stringLength) - (senderCertificate.stringLength == -1 ? 0 : senderCertificate.stringLength) - (receiverCertificateThumbprint.stringLength == -1 ? 0 : receiverCertificateThumbprint.stringLength) - 32']
       ]
       ['CLO','false'     OpcuaCloseRequest
           [simple          string '8'         'chunk']
           [implicit        int 32             'messageSize' 'lengthInBytes']
           [simple          int 32             'secureChannelId']
           [simple          int 32             'secureTokenId']
           [simple          int 32             'sequenceNumber']
           [simple          int 32             'requestId']
           [simple          ExtensionObject       'message' ['false']]
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

[type 'ByteStringArray'
    [simple int 32 'arrayLength']
    [array uint 8 'value' count 'arrayLength']
]

[type 'GuidValue'
    [simple uint 32 'data1']
    [simple uint 16 'data2']
    [simple uint 16 'data3']
    [array int 8 'data4' count '2']
    [array int 8 'data5' count '6']
]

[type 'ExpandedNodeId'
    [simple bit 'namespaceURISpecified']
    [simple bit 'serverIndexSpecified']
    [simple NodeIdTypeDefinition 'nodeId']
    [virtual string '-1' 'utf-8' 'identifier' 'nodeId.identifier']
    [optional PascalString 'namespaceURI' 'namespaceURISpecified']
    [optional uint 32 'serverIndex' 'serverIndexSpecified']
]

[type 'ExtensionHeader'
    [reserved int 5 '0x00']
    [simple bit 'xmlbody']
    [simple bit 'binaryBody]
]

[type 'ExtensionObjectEncodingMask'
    [reserved int 5 '0x00']
    [simple bit 'typeIdSpecified']
    [simple bit 'xmlbody']
    [simple bit 'binaryBody]
]

[type 'ExtensionObject' [bit 'includeEncodingMask']
    //A serialized object prefixed with its data type identifier.
    [simple ExpandedNodeId 'typeId']
    [optional ExtensionObjectEncodingMask 'encodingMask' 'includeEncodingMask']
    [virtual string '-1' 'identifier' 'typeId.identifier']
    [simple ExtensionObjectDefinition 'body' ['identifier']]
]

[discriminatedType 'ExtensionObjectDefinition' [string '-1' 'identifier']
    [typeSwitch 'identifier'
        ['0' NullExtension
        ]

        <xsl:for-each select="/opc:TypeDictionary/opc:StructuredType[(@BaseType = 'ua:ExtensionObject') and not(@Name = 'UserIdentityToken') and not(@Name = 'PublishedDataSetDataType') and not(@Name = 'DataSetReaderDataType')]">
            <xsl:message><xsl:value-of select="@Name"/></xsl:message>
            <xsl:variable name="extensionName" select="@Name"/>
            <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName=$extensionName]"/>
        </xsl:for-each>

        ['811' DataChangeNotification
            [implicit int 32 'notificationLength' 'lengthInBytes']
            [simple int 32 'noOfMonitoredItems']
            [array ExtensionObjectDefinition  'monitoredItems' count 'noOfMonitoredItems' ['808']]
            [simple int 32 'noOfDiagnosticInfos']
            [array DiagnosticInfo  'diagnosticInfos' count 'noOfDiagnosticInfos']
        ]
        ['916' EventNotificationList
            [implicit int 32 'notificationLength' 'lengthInBytes']
            [simple int 32 'noOfEvents']
            [array ExtensionObjectDefinition  'events' count 'noOfEvents' ['919']]
        ]
        ['820' StatusChangeNotification
            [implicit int 32 'notificationLength' 'lengthInBytes']
            [simple StatusCode 'status']
            [simple DiagnosticInfo 'diagnosticInfo']
        ]

        ['316' UserIdentityToken
            [implicit int 32 'policyLength' 'policyId.lengthInBytes']
            [simple PascalString 'policyId']
            [simple UserIdentityTokenDefinition 'userIdentityTokenDefinition' ['policyId.stringValue']]
        ]
    ]
]

[discriminatedType 'UserIdentityTokenDefinition' [string '-1' 'identifier']
    [typeSwitch 'identifier'
        ['anonymous' AnonymousIdentityToken
        ]
        ['username' UserNameIdentityToken
            [simple PascalString 'userName']
            [simple PascalByteString 'password']
            [simple PascalString 'encryptionAlgorithm']
        ]
        ['certificate' X509IdentityToken
            [simple PascalByteString 'certificateData']
        ]
        ['identity' IssuedIdentityToken
            [simple PascalByteString 'tokenData']
            [simple PascalString 'encryptionAlgorithm']
        ]
    ]
]


[discriminatedType 'Variant'
    [simple bit 'arrayLengthSpecified']
    [simple bit 'arrayDimensionsSpecified']
    [discriminator uint 6 'VariantType']
    [typeSwitch 'VariantType','arrayLengthSpecified'
        ['1' VariantBoolean [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array int 8 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['2' VariantSByte [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array int 8 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['3' VariantByte [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array uint 8 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['4' VariantInt16 [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array int 16 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['5' VariantUInt16 [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array uint 16 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['6' VariantInt32 [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array int 32 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['7' VariantUInt32 [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array uint 32 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['8' VariantInt64 [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array int 64 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['9' VariantUInt64 [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array uint 64 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['10' VariantFloat [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array float 8.23 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['11' VariantDouble [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array float 11.52 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['12' VariantString [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array PascalString 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['13' VariantDateTime [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array int 64 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['14' VariantGuid [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array GuidValue 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['15' VariantByteString [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array ByteStringArray 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['16' VariantXmlElement [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array PascalString 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['17' VariantNodeId [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array NodeId 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['18' VariantExpandedNodeId [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array ExpandedNodeId 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['19' VariantStatusCode [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array StatusCode 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['20' VariantQualifiedName [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array QualifiedName 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['21' VariantLocalizedText [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array LocalizedText 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['22' VariantExtensionObject [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array ExtensionObject 'value' count 'arrayLength == null ? 1 : arrayLength' ['true']]
        ]
        ['23' VariantDataValue [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array DataValue 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['24' VariantVariant [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array Variant 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
        ['25' VariantDiagnosticInfo [bit 'arrayLengthSpecified']
            [optional int 32 'arrayLength' 'arrayLengthSpecified']
            [array DiagnosticInfo 'value' count 'arrayLength == null ? 1 : arrayLength']
        ]
    ]
    [optional int 32 'noOfArrayDimensions' 'arrayDimensionsSpecified']
    [array bit 'arrayDimensions' count 'noOfArrayDimensions == null ? 0 : noOfArrayDimensions']
]

[discriminatedType 'NodeIdTypeDefinition'
    [abstract string '-1' 'identifier']
    [discriminator NodeIdType 'nodeType']
    [typeSwitch 'nodeType'
        ['nodeIdTypeTwoByte' NodeIdTwoByte
            [simple uint 8 'id']
            [virtual string '-1' 'identifier' 'id']
        ]
        ['nodeIdTypeFourByte' NodeIdFourByte
            [simple uint 8 'namespaceIndex']
            [simple uint 16 'id']
            [virtual string '-1' 'identifier' 'id']
        ]
        ['nodeIdTypeNumeric' NodeIdNumeric
            [simple uint 16 'namespaceIndex']
            [simple uint 32 'id']
            [virtual string '-1' 'identifier' 'id']
        ]
        ['nodeIdTypeString' NodeIdString
            [simple uint 16 'namespaceIndex']
            [simple PascalString 'id']
            [virtual string '-1' 'identifier' 'id.stringValue']
        ]
        ['nodeIdTypeGuid' NodeIdGuid
            [simple uint 16 'namespaceIndex']
            [array int 8 'id' count '16']
            [virtual string '-1' 'identifier' 'id']
        ]
        ['nodeIdTypeByteString' NodeIdByteString
            [simple uint 16 'namespaceIndex']
            [simple PascalByteString 'id']
            [virtual string '-1' 'identifier' 'id.stringValue']
        ]
    ]
]

[type 'NodeId'
    [reserved int 2 '0x00']
    [simple NodeIdTypeDefinition 'nodeId']
    [virtual string '-1' 'id' 'nodeId.identifier']
]

[type 'PascalString'
    [implicit int 32 'sLength'          'stringValue.length == 0 ? -1 : stringValue.length']
    [simple string 'sLength == -1 ? 0 : sLength * 8' 'UTF-8' 'stringValue']
    [virtual  int 32 'stringLength'     'stringValue.length == -1 ? 0 : stringValue.length']
]

[type 'PascalByteString'
    [simple int 32 'stringLength']
    [array int 8 'stringValue' count 'stringLength == -1 ? 0 : stringLength' ]
]

[type 'Structure'

]

[type 'DataTypeDefinition'

]

<xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[(@Name != 'ExtensionObject') and (@Name != 'Variant') and (@Name != 'NodeId') and (@Name != 'ExpandedNodeId') and not(@BaseType)]"/>
<xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType"/>
<xsl:apply-templates select="/opc:TypeDictionary/opc:OpaqueType"/>

[enum string '-1' 'OpcuaDataType' [uint 8 'variantType']
    ['NULL' NULL ['0']]
    ['BOOL' BOOL ['1']]
    ['BYTE' BYTE ['3']]
    ['SINT' SINT ['2']]
    ['INT' INT ['4']]
    ['DINT' DINT ['6']]
    ['LINT' LINT ['8']]
    ['USINT' USINT ['3']]
    ['UINT' UINT ['5']]
    ['UDINT' UDINT ['7']]
    ['ULINT' ULINT ['9']]
    ['REAL' REAL ['10']]
    ['LREAL' LREAL ['11']]
    ['TIME' TIME ['1']]
    ['LTIME' LTIME ['1']]
    ['DATE' DATE ['1']]
    ['LDATE' LDATE ['1']]
    ['TIME_OF_DAY' TIME_OF_DAY ['1']]
    ['LTIME_OF_DAY' LTIME_OF_DAY ['1']]
    ['DATE_AND_TIME' DATE_AND_TIME ['13']]
    ['LDATE_AND_TIME' LDATE_AND_TIME ['1']]
    ['CHAR' CHAR ['1']]
    ['WCHAR' WCHAR ['1']]
    ['STRING' STRING ['12']]
]

[enum string '-1' 'OpcuaIdentifierType'
    ['s' STRING_IDENTIFIER]
    ['i' NUMBER_IDENTIFIER]
    ['g' GUID_IDENTIFIER]
    ['b' BINARY_IDENTIFIER]
]


    </xsl:template>
</xsl:stylesheet>
