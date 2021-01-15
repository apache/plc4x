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
                xmlns:opc="http://opcfoundation.org/BinarySchema/"
                xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                xmlns:ua="http://opcfoundation.org/UA/"
                xmlns:tns="http://opcfoundation.org/UA/"
                xmlns:node="http://opcfoundation.org/UA/2011/03/UANodeSet.xsd">

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />

    <xsl:param name="services"></xsl:param>
    <xsl:param name="statusCodes"></xsl:param>
    <xsl:param name="servicesEnum"></xsl:param>

    <xsl:variable name="originaldoc" select="/"/>

    <xsl:param name="file" select="document($services)"/>
    <xsl:param name="statusCodeFile" select="unparsed-text($statusCodes)"/>
    <xsl:param name="servicesEnumFile" select="unparsed-text($servicesEnum)"/>

    <xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz'"/>
    <xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ'"/>

    <xsl:template match="/">
[discriminatedType 'OpcuaMessage'
    [simple         int 8   'OPCUAnodeIdEncodingMask' ]
    [simple         int 8   'OPCUAnodeIdNamespaceIndex' ]
    [discriminator  int 16   'OPCUAnodeId' ]
    [typeSwitch 'OPCUAnodeId'
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='OpenSecureChannelRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='OpenSecureChannelResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CreateSessionRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CreateSessionResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CreateSubscriptionRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CreateSubscriptionResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CreateMonitoredItemsRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CreateMonitoredItemsRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='DeleteSubscriptionsRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='ActivateSessionRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='ActivateSessionResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='ReadRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='ReadResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='WriteRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='WriteResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='BrowseRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='BrowseResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='ServiceFault']"/>

    ['473' CloseSessionRequest
        [simple RequestHeader 'requestHeader']
        [reserved uint 7 '0x00']
        [simple bit 'deleteSubscriptions']
    ]
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CloseSessionResponse']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CloseSecureChannelRequest']"/>
        <xsl:apply-templates select="$file/node:UANodeSet/node:UADataType[@BrowseName='CloseSecureChannelResponse']"/>
    ]
]

[type 'RequestHeader'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='RequestHeader']"/>
]
[enum int 32 'SecurityTokenRequestType'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='SecurityTokenRequestType']"/>
]
[enum int 32 'MessageSecurityMode'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='MessageSecurityMode']"/>
]
[type 'ResponseHeader'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='ResponseHeader']"/>
]
[type 'ChannelSecurityToken'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='ChannelSecurityToken']"/>
]

[type 'MonitoredItemCreateRequest'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='MonitoredItemCreateRequest']"/>
]

[type 'BrowseResult'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='BrowseResult']"/>
]

[type 'ViewDescription'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='ViewDescription']"/>
]

[type 'BrowseDescription'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='BrowseDescription']"/>
]

[type 'ReferenceDescription'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='ReferenceDescription']"/>
]

[enum int 32 'MonitoringMode'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='MonitoringMode']"/>
]

[type 'MonitoringParameters'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='MonitoringParameters']"/>
]

[enum int 32 'BrowseDirection'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='BrowseDirection']"/>
]

[type 'ReferenceDescription'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='ReferenceDescription']"/>
]

[enum int 32 'NodeClass'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='NodeClass']"/>
]

[type 'UserNameIdentityToken'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='UserNameIdentityToken']"/>
]

[type 'DiagnosticInfo'
    [simple bit 'symbolicIdSpecified']
    [simple bit 'namespaceURISpecified']
    [simple bit 'localizedTextSpecified']
    [simple bit 'localeSpecified']
    [simple bit 'additionalInfoSpecified']
    [simple bit 'innerStatusCodeSpecified']
    [simple bit 'innerDiagnosticInfoSpecified']
    [simple bit 'reserved1']
    [optional int 32 'symbolicId' 'symbolicIdSpecified']
    [optional int 32 'namespaceURI' 'namespaceURISpecified']
    [optional int 32 'locale' 'localizedTextSpecified']
    [optional int 32 'localizedText' 'localeSpecified']
    [optional PascalString 'additionalInfo' 'additionalInfoSpecified']
    [optional StatusCode 'innerStatusCode' 'innerStatusCodeSpecified']
    [optional DiagnosticInfo 'innerDiagnosticInfo' 'innerDiagnosticInfoSpecified']
]

[type 'StatusCode'
    [simple int 32 'statusCode']
]

[type 'XmlElement'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='XmlElement']"/>
]

[type 'DataValue'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='DataValue']"/>
]

[enum int 6 'NodeIdType'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='NodeIdType']"/>
]

[type 'TwoByteNodeId'
        [simple uint 8 'identifier']
]

[type 'FourByteNodeId'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='FourByteNodeId']"/>
]

[type 'NumericNodeId'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='NumericNodeId']"/>
]

[type 'StringNodeId'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='StringNodeId']"/>
]

[type 'GuidNodeId'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='GuidNodeId']"/>
]

[type 'ByteStringNodeId'
    [simple uint 16 'namespaceIndex']
    [simple int 32 'bodyLength']
    [array int 8 'body' count 'bodyLength == -1 ? 0 : bodyLength']
]

[type 'DataValue'
    //A value with an associated timestamp, and quality..
    [simple bit 'reserved2']
    [simple bit 'reserved1']
    [simple bit 'serverPicosecondsSpecified']
    [simple bit 'sourcePicosecondsSpecified']
    [simple bit 'serverTimestampSpecified']
    [simple bit 'sourceTimestampSpecified']
    [simple bit 'statusCodeSpecified']
    [simple bit 'valueSpecified']
    [optional Variant 'value' 'valueSpecified']
    [optional StatusCode 'statusCode' 'statusCodeSpecified']
    [optional int 64 'sourceTimestamp' 'sourceTimestampSpecified']
    [optional uint 16 'sourcePicoseconds' 'sourcePicosecondsSpecified']
    [optional int 64 'serverTimestamp' 'serverTimestampSpecified']
    [optional uint 16 'serverPicoseconds' 'serverPicosecondsSpecified']
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
            [array ExtensionObject 'value' count 'arrayLength == null ? 1 : arrayLength']
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



[discriminatedType 'NodeId'
    [reserved int 2 '0x00']
    [simple NodeIdType 'nodeIdType']
    [typeSwitch 'nodeIdType'
        ['NodeIdType.nodeIdTypeTwoByte' NodeIdTwoByte
            [simple TwoByteNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeFourByte' NodeIdFourByte
            [simple FourByteNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeNumeric' NodeIdNumeric
            [simple NumericNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeString' NodeIdString
            [simple StringNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeGuid' NodeIdGuid
            [simple GuidNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeByteString' NodeIdByteString
            [simple ByteStringNodeId 'id']
        ]
    ]
]

[discriminatedType 'ExpandedNodeId'
    [simple bit 'namespaceURISpecified']
    [simple bit 'serverIndexSpecified']
    [discriminator NodeIdType 'nodeIdType']
    [typeSwitch 'nodeIdType'
        ['NodeIdType.nodeIdTypeTwoByte' ExpandedNodeIdTwoByte
            [simple TwoByteNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeFourByte' ExpandedNodeIdFourByte
            [simple FourByteNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeNumeric' ExpandedNodeIdNumeric
            [simple NumericNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeString' ExpandedNodeIdString
            [simple StringNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeGuid' ExpandedNodeIdGuid
            [simple GuidNodeId 'id']
        ]
        ['NodeIdType.nodeIdTypeByteString' ExpandedNodeIdByteString
            [simple ByteStringNodeId 'id']
        ]
    ]
    [optional PascalString 'namespaceURI' 'namespaceURISpecified']
    [optional uint 32 'serverIndex' 'serverIndexSpecified']
]

[discriminatedType 'ExtensionObject'
    //A serialized object prefixed with its data type identifier.
    [simple ExpandedNodeId 'nodeId']
    [simple uint 8 'encodingMask']
    [optional int 32 'bodyLength' 'encodingMask > 0']
    [array int 8 'body' count 'bodyLength == null ? 0 : bodyLength']
]

[type 'PascalString'
    [simple int 32 'stringLength']
    [optional string 'stringLength == -1 ? 0 : stringLength * 8' 'UTF-8' 'stringValue' 'stringLength >= 0']
]

[type 'PascalByteString'
    [simple int 32 'stringLength']
    [array int 8 'stringValue' count 'stringLength == -1 ? 0 : stringLength']
]

[type 'LocalizedText'
    [simple uint 6 'reserved1']
    [simple bit 'localeSpecified']
    [simple bit 'textSpecified']
    [optional PascalString 'Locale' 'localeSpecified']
    [optional PascalString 'Text' 'textSpecified']
]

[type 'QualifiedName'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='QualifiedName']"/>
]

[type 'ApplicationDescription'
    [simple PascalString 'applicationUri']
    [simple PascalString 'productUri']
    [simple LocalizedText 'applicationName']
    [simple ApplicationType 'applicationType']
    [simple PascalString 'gatewayServerUri']
    [simple PascalString 'discoveryProfileUri']
    [simple int 32 'noOfDiscoveryUrls']
    [optional PascalString 'discoveryUrls' 'noOfDiscoveryUrls > 0']
]

[type 'EndpointDescription'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='EndpointDescription']"/>
]

[type 'SignedSoftwareCertificate'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='SignedSoftwareCertificate']"/>
]


[type 'SignatureData'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='SignatureData']"/>
]


[enum int 32 'ApplicationType'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='ApplicationType']"/>
]

[type 'UserTokenPolicy'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='UserTokenPolicy']"/>
]

[enum int 32 'UserTokenType'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='UserTokenType']"/>
]

[enum int 32 'TimestampsToReturn'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:EnumeratedType[@Name='TimestampsToReturn']"/>
]

[type 'ReadValueId'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='ReadValueId']"/>
]

[type 'WriteValue'
    <xsl:apply-templates select="/opc:TypeDictionary/opc:StructuredType[@Name='WriteValue']"/>
]

        <xsl:call-template name="statusCodeParsing"/>
        <xsl:call-template name="servicesEnumParsing"/>

    </xsl:template>

    <xsl:template match="node:UAVariable">
        <xsl:variable name="browseName">
            <xsl:value-of select='@BrowseName'/>
        </xsl:variable>
        <xsl:choose>
            <xsl:when test="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]">
                <xsl:choose>
                    <xsl:when test="not(@BrowseName='Vector') and not(substring(@BrowseName,1,1) = '&lt;') and not(number(substring(@BrowseName,1,1)))">
    [type '<xsl:value-of select='@BrowseName'/>'
        <xsl:apply-templates select="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]"/>]
                    </xsl:when>
                </xsl:choose>
            </xsl:when>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="node:UADataType[not(Definition)]">
        <xsl:variable name="browseName">
            <xsl:value-of select='@BrowseName'/>
        </xsl:variable>
        <xsl:choose>
            <xsl:when test="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]">
                <xsl:choose>
                    <xsl:when test="not(Definition) and not(@BrowseName = 'Duration') and not(number(substring(@BrowseName,1,1))) and not(@IsAbstract) and number(substring(@NodeId,3)) &gt; 29">
    ['<xsl:value-of select="number(substring(@NodeId,3)) + 2"/><xsl:text>' </xsl:text><xsl:value-of select='@BrowseName'/><xsl:text>
                </xsl:text>
                        <xsl:apply-templates select="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]"/>]
                    </xsl:when>
                </xsl:choose>
            </xsl:when>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="opc:EnumeratedType">
        <xsl:apply-templates select="opc:Documentation"/>
        <xsl:apply-templates select="opc:EnumeratedValue"/>
    </xsl:template>

    <xsl:template match="opc:Documentation">
    //<xsl:value-of select="."/>.
    </xsl:template>

    <xsl:template match="opc:EnumeratedValue">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
                <xsl:with-param name="switchField" select="../@Name"/>
                <xsl:with-param name="switchValue" select="1"/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Value"/>' <xsl:value-of select="$objectTypeId"/>]
    </xsl:template>

    <xsl:template match="opc:OpaqueType[not(@Name = 'Duration')]">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
                <xsl:with-param name="switchField" select="@SwitchField"/>
                <xsl:with-param name="switchValue" select="@SwitchValue"/>
            </xsl:call-template>
        </xsl:variable>[type '<xsl:value-of select="@Name"/>'
        <xsl:apply-templates select="opc:Documentation"/>
  ]
    </xsl:template>

    <xsl:template match="opc:StructuredType[not(@Name = 'Vector')]">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
                <xsl:with-param name="switchField" select="@SwitchField"/>
                <xsl:with-param name="switchValue" select="@SwitchValue"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:apply-templates select="opc:Documentation"/>
        <xsl:apply-templates select="opc:Field"/>
    </xsl:template>

    <xsl:template match="opc:Field">
        <xsl:variable name="objectTypeId">
            <xsl:value-of select="@Name"/>
        </xsl:variable>
        <xsl:variable name="lowerCaseName">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
                <xsl:with-param name="switchField" select="@SwitchField"/>
                <xsl:with-param name="switchValue" select="@SwitchValue"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="lowerCaseLengthField">
            <xsl:call-template name="lowerCaseLeadingChar">
                <xsl:with-param name="text" select="@LengthField"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="mspecType">simple</xsl:variable>
        <xsl:variable name="dataType">
            <xsl:call-template name="translateDataType">
                <xsl:with-param name="datatype" select="@TypeName"/>
                <xsl:with-param name="name" select="-1"/>
            </xsl:call-template>
        </xsl:variable>

        <xsl:choose>
            <xsl:when test="@LengthField">[array <xsl:value-of select="$dataType"/>  '<xsl:value-of select="$lowerCaseName"/>' count '<xsl:value-of select="$lowerCaseLengthField"/>']
    </xsl:when>
            <xsl:otherwise>[<xsl:value-of select="$mspecType"/><xsl:text> </xsl:text><xsl:value-of select="$dataType"/> '<xsl:value-of select="$lowerCaseName"/>']
    </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template name="clean-id-string">
        <xsl:param name="text"/>
        <xsl:param name="switchField"/>
        <xsl:param name="switchValue"/>
        <xsl:choose>
            <xsl:when test="$switchValue">
                <xsl:call-template name="lowerCaseLeadingChar">
                    <xsl:with-param name="text" select="concat($switchField, $text)"/>
                </xsl:call-template>
            </xsl:when>
            <xsl:otherwise>
                <xsl:call-template name="lowerCaseLeadingChar">
                    <xsl:with-param name="text" select="$text"/>
                </xsl:call-template>
            </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template name="lowerCaseLeadingChar">
        <xsl:param name="text"/>
        <xsl:value-of select="concat(translate(substring($text, 1, 1), $uppercase, $lowercase), substring($text, 2))"/>
    </xsl:template>

    <xsl:template name="translateDataType">
        <xsl:param name="datatype"/>
        <xsl:param name="name"/>
        <xsl:choose>
            <xsl:when test="$datatype = 'opc:Bit'">bit</xsl:when>
            <xsl:when test="$datatype = 'opc:Boolean'">bit</xsl:when>
            <xsl:when test="$datatype = 'opc:Byte'">uint 8</xsl:when>
            <xsl:when test="$datatype = 'opc:SByte'">int 8</xsl:when>
            <xsl:when test="$datatype = 'opc:Int16'">int 16</xsl:when>
            <xsl:when test="$datatype = 'opc:UInt16'">uint 16</xsl:when>
            <xsl:when test="$datatype = 'opc:Int32'">int 32</xsl:when>
            <xsl:when test="$datatype = 'opc:UInt32'">uint 32</xsl:when>
            <xsl:when test="$datatype = 'opc:Int64'">int 64</xsl:when>
            <xsl:when test="$datatype = 'opc:UInt64'">uint 64</xsl:when>
            <xsl:when test="$datatype = 'opc:Float'">float 8.23</xsl:when>
            <xsl:when test="$datatype = 'opc:Double'">float 11.52</xsl:when>
            <xsl:when test="$datatype = 'opc:Char'">string '1'</xsl:when>
            <xsl:when test="$datatype = 'opc:CharArray'">PascalString</xsl:when>
            <xsl:when test="$datatype = 'opc:Guid'">GuidValue</xsl:when>
            <xsl:when test="$datatype = 'opc:ByteString'">PascalByteString</xsl:when>
            <xsl:when test="$datatype = 'opc:DateTime'">int 64</xsl:when>
            <xsl:when test="$datatype = 'opc:String'">PascalString</xsl:when>
            <xsl:otherwise><xsl:value-of select="substring-after($datatype,':')"/></xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template name="statusCodeParsing" >
        <xsl:variable name="tokenizedLine" select="tokenize($statusCodeFile, '\r\n|\r|\n')" />
[enum int 32 'OpcuaStatusCodes'
<xsl:for-each select="$tokenizedLine">
    <xsl:variable select="tokenize(., ',')" name="values" />    ['<xsl:value-of select="$values[2]"/>'  <xsl:value-of select="$values[1]"/>]
</xsl:for-each>
]
</xsl:template>

    <xsl:template name="servicesEnumParsing" >
        <xsl:variable name="tokenizedLine" select="tokenize($servicesEnumFile, '\r\n|\r|\n')" />
[enum int 32 'OpcuaNodeIdServices'
        <xsl:for-each select="$tokenizedLine">
            <xsl:variable select="tokenize(., ',')" name="values" />
            <xsl:choose>
                <xsl:when test="$values[2]">['<xsl:value-of select="$values[2]"/>'  <xsl:value-of select="$values[1]"/>]
    </xsl:when>
            </xsl:choose>
        </xsl:for-each>
]
    </xsl:template>

</xsl:stylesheet>
