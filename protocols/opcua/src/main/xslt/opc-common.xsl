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
                xmlns:xs="http://www.w3.org/2001/XMLSchema"
                xmlns:opc="http://opcfoundation.org/BinarySchema/"
                xmlns:plc4x="https://plc4x.apache.org/"
                xmlns:map="http://www.w3.org/2005/xpath-functions/map"
                xmlns:ua="http://opcfoundation.org/UA/"
                xmlns:tns="http://opcfoundation.org/UA/"
                xmlns:node="http://opcfoundation.org/UA/2011/03/UANodeSet.xsd">

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />

    <xsl:param name="services"/>
    <xsl:param name="file" select="document($services)"/>

    <xsl:variable name="originaldoc" select="/"/>

    <xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz'"/>
    <xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ'"/>

    <xsl:variable name="dataTypeLength" as="map(xs:string, xs:int)">
        <xsl:map>
            <xsl:for-each select="//opc:EnumeratedType">
                <xsl:choose>
                    <xsl:when test="@Name != '' or @LengthInBits != ''">
                        <xsl:map-entry key="concat('ua:', xs:string(@Name))" select="xs:int(@LengthInBits)"/>
                    </xsl:when>
                </xsl:choose>
            </xsl:for-each>
        </xsl:map>
    </xsl:variable>

    <xsl:template match="node:UADataType">
        <xsl:variable name="servicesRoot" select=".."/>
        <xsl:variable name="browseName">
            <xsl:value-of select='@BrowseName'/>
        </xsl:variable>
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-datatype-string">
                <xsl:with-param name="text" select="@BrowseName"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="extensionId">
            <xsl:call-template name="plc4x:getExtensionId">
                <xsl:with-param name="dataTypeName" select="$browseName" />
                <xsl:with-param name="servicesRoot" select="$servicesRoot" />
            </xsl:call-template>
        </xsl:variable>
        <xsl:choose>
            <xsl:when test="$extensionId != ''"><xsl:text>
        </xsl:text><xsl:value-of select="concat('[''', $extensionId, ''' ', $objectTypeId)" /><xsl:text>
            </xsl:text>
                <xsl:call-template name="plc4x:parseFields">
                    <xsl:with-param name="baseNode" select="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]"/>
                    <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                    <xsl:with-param name="currentNodePosition">1</xsl:with-param>
                    <xsl:with-param name="currentBytePosition">0</xsl:with-param>
                    <xsl:with-param name="currentBitPosition">0</xsl:with-param>
                </xsl:call-template><xsl:text>
        ]</xsl:text>
            </xsl:when>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="node:UAVariable">
        <xsl:variable name="browseName">
            <xsl:value-of select='@BrowseName'/>
        </xsl:variable>
        <xsl:choose>
            <xsl:when test="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]">
                <xsl:choose>
                    <xsl:when test="not(@BrowseName='Vector') and not(substring(@BrowseName,1,1) = '&lt;') and not(number(substring(@BrowseName,1,1)))">
                        [type <xsl:value-of select='@BrowseName'/>
                        <xsl:apply-templates select="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$browseName]"/>]
                    </xsl:when>
                </xsl:choose>
            </xsl:when>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="opc:EnumeratedType">
        <xsl:message>[INFO] Parsing Enumerated Datatype - <xsl:value-of select="@Name"/></xsl:message>
        <xsl:choose>
            <xsl:when test="opc:EnumeratedValue"><xsl:text>

</xsl:text>[enum uint <xsl:value-of select="@LengthInBits"/><xsl:text> </xsl:text><xsl:value-of select="@Name"/><xsl:text>
</xsl:text>
        <xsl:apply-templates select="opc:Documentation"/><xsl:text>
    </xsl:text>
        <xsl:apply-templates select="opc:EnumeratedValue"/>
]
            </xsl:when>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="opc:Documentation">// <xsl:value-of select="."/></xsl:template>

    <xsl:template match="opc:EnumeratedValue">
        <xsl:message>[INFO] Parsing Enumerated Value - <xsl:value-of select="@Name"/></xsl:message>
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
                <xsl:with-param name="switchField" select="../@Name"/>
                <xsl:with-param name="switchValue" select="1"/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Value"/>' <xsl:value-of select="$objectTypeId"/>]
    </xsl:template>

    <xsl:template match="opc:OpaqueType[not(@Name = 'Duration')]">
        <xsl:message>[INFO] Parsing Opaque Datatype - <xsl:value-of select="@Name"/></xsl:message>
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
                <xsl:with-param name="switchField" select="@SwitchField"/>
                <xsl:with-param name="switchValue" select="@SwitchValue"/>
            </xsl:call-template>
        </xsl:variable>[type <xsl:value-of select="@Name"/><xsl:text>
    </xsl:text>
        <xsl:apply-templates select="opc:Documentation"/>
        <xsl:choose>
            <xsl:when test="@LengthInBits != ''">
                [simple uint <xsl:value-of select="@LengthInBits"/><xsl:text> </xsl:text><xsl:value-of select="$objectTypeId"/>]</xsl:when>
        </xsl:choose>
]
    </xsl:template>

    <xsl:template match="opc:StructuredType[not (@BaseType)]">
        <xsl:message>[INFO] Parsing Top Level Structured Datatype - <xsl:value-of select="@Name"/></xsl:message>
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-datatype-string">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>[type <xsl:value-of select="$objectTypeId"/><xsl:text>
    </xsl:text>
        <xsl:apply-templates select="opc:Documentation"/><xsl:text>
    </xsl:text>
        <xsl:call-template name="plc4x:parseFields">
            <xsl:with-param name="baseNode" select="."/>
            <xsl:with-param name="servicesRoot"/>
            <xsl:with-param name="currentNodePosition">1</xsl:with-param>
            <xsl:with-param name="currentBytePosition">0</xsl:with-param>
            <xsl:with-param name="currentBitPosition">0</xsl:with-param>
        </xsl:call-template>
        ]
    </xsl:template>


    <xsl:template match="opc:Field">
        <xsl:param name="servicesRoot"/>
        <xsl:param name="fieldCount"/>
        <xsl:message>[INFO] Parsing Field - <xsl:value-of select="@Name"/> of <xsl:value-of select="../@Name" /></xsl:message>
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
        <xsl:variable name="lowerCaseArrayField">
            <!-- turn NoOfDiagnosticInfos into diagnosticInfos -->
            <xsl:call-template name="lowerCaseLeadingChar">
                <xsl:with-param name="text" select="substring(@Name, 5)"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="dataType">
            <xsl:call-template name="plc4x:getDataTypeField">
                <xsl:with-param name="datatype" select="@TypeName"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="dataTypeLength"><xsl:value-of select="@Length"/></xsl:variable>
        <xsl:variable name="mspecType">
            <xsl:call-template name="plc4x:getMspecName">
                <xsl:with-param name="field" select="."/>
                <xsl:with-param name="name" select="$lowerCaseName"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="lowerCaseSwitchField">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@SwitchField"/>
                <xsl:with-param name="switchField" select="@SwitchField"/>
                <xsl:with-param name="switchValue" select="@SwitchValue"/>
            </xsl:call-template>
        </xsl:variable>
        <!-- Depending on what kind of mspec variable it is, we have to include different arguments -->
        <xsl:choose>
            <xsl:when test="@LengthField">
                <xsl:choose>
                    <xsl:when test="$dataType = 'ExtensionObjectDefinition'">
                        <xsl:variable name="browseName" select="substring-after(@TypeName,':')"/>
                        <xsl:variable name="id" select="number(substring-after($file/node:UANodeSet/node:UADataType[@BrowseName=$browseName]/@NodeId, '=')) + 2"/><xsl:value-of select="concat('[array ', $dataType, '(''&quot;', $id, '&quot;'') ', $lowerCaseName, ' count ''', $lowerCaseLengthField, ''']')" /></xsl:when>
                    <xsl:when test="$dataType = 'ExtensionObject'"><xsl:value-of select="concat('[array ', $dataType, '(''true'') ', $lowerCaseName, ' count ''', $lowerCaseLengthField, ''']')" /></xsl:when>
                    <xsl:otherwise>
                        <xsl:variable name="extensionIdArgs">
                            <xsl:call-template name="plc4x:getExtensionIdArgs">
                                <xsl:with-param name="dataTypeName" select="$dataType"/>
                                <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                            </xsl:call-template>
                        </xsl:variable><xsl:value-of select="concat('[array ', $dataType, $extensionIdArgs, ' ', $lowerCaseName, ' count ''', $lowerCaseLengthField, ''']')" /></xsl:otherwise>
                </xsl:choose>
            </xsl:when>
            <xsl:when test="$mspecType = 'reserved'">
                <xsl:choose>
                    <xsl:when test="xs:int(@Length) gt 1"><xsl:value-of select="concat('[', $mspecType, ' uint ', $dataTypeLength, ' ''0x00'']')" /></xsl:when>
                    <xsl:otherwise><xsl:value-of select="concat('[', $mspecType, ' ', $dataType, ' ''false'']')" /></xsl:otherwise>
                </xsl:choose>
            </xsl:when>
            <xsl:when test="$mspecType = 'optional'"><xsl:value-of select="concat('[', $mspecType, ' ', $dataType, ' ', $lowerCaseName, ' ''', $lowerCaseSwitchField, ''']')" /></xsl:when>
            <xsl:when test="$dataType = 'ExtensionObjectDefinition'">
                <xsl:variable name="browseName" select="substring-after(@TypeName,':')"/>
                <xsl:variable name="id" select="number(substring-after($file/node:UANodeSet/node:UADataType[@BrowseName=$browseName]/@NodeId, '=')) + 2"/>
                <xsl:value-of select="concat('[', $mspecType, ' ', $dataType, '(''&quot;', $id, '&quot;'') ', $lowerCaseName, ']')" /></xsl:when>
            <xsl:when test="$dataType = 'ExtensionObject'"><xsl:value-of select="concat('[', $mspecType, ' ', $dataType, '(''true'') ', $lowerCaseName, ']')" /></xsl:when>
            <xsl:otherwise>
                <xsl:choose>
                    <!-- Injected field used in rare cases when ExtensionObject is serialized to binary form -->
                    <xsl:when test="$lowerCaseName = 'lengthInBytes'"><xsl:value-of select="concat('[implicit ', $dataType, ' size ''lengthInBytes - 4'']')" /> // minus type information</xsl:when>
                    <!-- Array length -->
                    <xsl:when test="$mspecType = 'implicit'"><xsl:value-of select="concat('[implicit int 32 ', $lowerCaseName, ' ''', $lowerCaseArrayField, ' == null ? -1 : COUNT(', $lowerCaseArrayField, ')'']')" /></xsl:when>
                    <!-- Basic field -->
                    <xsl:otherwise>
                        <xsl:variable name="extensionIdArgs">
                            <xsl:call-template name="plc4x:getExtensionIdArgs">
                                <xsl:with-param name="dataTypeName" select="$dataType"/>
                                <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                            </xsl:call-template>
                        </xsl:variable><xsl:value-of select="concat('[', $mspecType, ' ', $dataType, $extensionIdArgs, ' ', $lowerCaseName, ']')" /></xsl:otherwise>
                </xsl:choose>
            </xsl:otherwise>
        </xsl:choose><!-- // fieldCount: <xsl:value-of select="$fieldCount" />, position <xsl:value-of select="position()" /> -->
        <xsl:if test="position() + 1 lt $fieldCount"><xsl:text>
            </xsl:text>
        </xsl:if>
    </xsl:template>

    <!-- Get the Mspec type simple/reserved/implicit/virtual/etc... -->
    <xsl:template name="plc4x:getMspecName">
        <xsl:param name="field"/>
        <xsl:param name="name"/>
        <xsl:message>[INFO] Parsing mspec name - <xsl:value-of select="$field/@TypeName"/> - <xsl:value-of select="$name"/> - <xsl:value-of select="$field/@SwitchField"/></xsl:message>
        <xsl:choose>
            <xsl:when test="starts-with($name, 'noOf')">implicit</xsl:when><!-- All noOfXyz (names lowercased by caller) fields in OPC UA are implicit array sizes -->
            <xsl:when test="$name = 'lengthInBytes'">implicit</xsl:when>
            <xsl:when test="starts-with($name, 'reserved')">reserved</xsl:when>
            <xsl:when test="$field/@SwitchField != ''">optional</xsl:when>
            <xsl:otherwise>simple</xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <!-- Convert a Data Type name so that it doesn't clash with mspec key words -->
    <xsl:template name="clean-datatype-string">
        <xsl:param name="text"/>
        <xsl:choose>
            <xsl:when test="$text = 'Vector'">OpcuaVector</xsl:when>
            <xsl:when test="$text = 'Vector'">OpcuaVector</xsl:when>
            <xsl:otherwise><xsl:value-of select="$text"/></xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <!-- Convert a variable name so that it doesn't clash with mspec key words -->
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

    <!-- Change the first character in string to lower case -->
    <xsl:template name="lowerCaseLeadingChar">
        <xsl:param name="text"/>
        <xsl:value-of select="concat(translate(substring($text, 1, 1), $uppercase, $lowercase), substring($text, 2))"/>
    </xsl:template>

    <!-- Convert the OPCUA data types to mspec data types. -->
    <xsl:template name="plc4x:getDataTypeField">
        <xsl:param name="datatype"/>
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
            <xsl:when test="$datatype = 'opc:Float'">float 32</xsl:when>
            <xsl:when test="$datatype = 'opc:Double'">float 64</xsl:when>
            <xsl:when test="$datatype = 'opc:Char'">string 8</xsl:when>
            <xsl:when test="$datatype = 'opc:CharArray'">PascalString</xsl:when>
            <xsl:when test="$datatype = 'opc:Guid'">GuidValue</xsl:when>
            <xsl:when test="$datatype = 'opc:ByteString'">PascalByteString</xsl:when>
            <xsl:when test="$datatype = 'opc:DateTime'">int 64</xsl:when>
            <xsl:when test="$datatype = 'opc:String'">PascalString</xsl:when>
            <!--<xsl:when test="not(starts-with($datatype, 'opc:'))">
                <xsl:value-of select="substring-after($datatype, ':')" />
                <xsl:choose>
                    <xsl:when test="starts-with($datatype, 'tns')"><xsl:value-of select="substring-after($datatype, ':')" /></xsl:when>
                    <xsl:otherwise><xsl:value-of select="" /></xsl:otherwise>
                </xsl:choose>
                <xsl:variable name="parent" select="substring-after($originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=substring-after($datatype,':')]/@BaseType, ':')"/>
                <xsl:choose>
                    <xsl:when test="$parent != ''"><xsl:variable name="id" select="substring-after($file/node:UANodeSet/node:UADataType[@BrowseName=substring-after($datatype,':')]/@BaseType, ':')"/><xsl:value-of select="$id" />/**/<xsl:value-of select="$parent" /> /*yolo*/
                        <xsl:choose>
                            <xsl:when test="substring-after($parent,':') = 'ExtensionObject'">ExtensionObjectDefinition</xsl:when>
                            <xsl:otherwise><xsl:value-of select="substring-after($parent,':')"/></xsl:otherwise>
                        </xsl:choose></xsl:when>
                    <xsl:otherwise><xsl:value-of select="substring-after($datatype,':')"/></xsl:otherwise>
                </xsl:choose>
            </xsl:when>
            -->
            <xsl:otherwise><xsl:value-of select="substring-after($datatype,':')"/></xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <!-- Gets the length in bits of a data type -->
    <xsl:function name="plc4x:getDataTypeLength" as="xs:integer">
        <xsl:param name="lengthMap" as="map(xs:string, xs:int)"/>
        <xsl:param name="datatype"/>
        <xsl:choose>
            <xsl:when test="map:contains($lengthMap, xs:string($datatype/[@TypeName]))">
                <xsl:message>[DEBUG] Bit Length <xsl:value-of select="$lengthMap(xs:string($datatype/[@TypeName]))"/></xsl:message>
                <xsl:value-of select="map:get($lengthMap, xs:string($datatype/[@TypeName]))"/>
            </xsl:when>
            <xsl:when test="($datatype/[@TypeName] = 'opc:Bit') or ($datatype/[@TypeName] = 'opc:Boolean')">
                <xsl:choose>
                    <xsl:when test="$datatype/[@Length] != ''">
                        <xsl:value-of select="xs:int($datatype/[@Length])"/>
                    </xsl:when>
                    <xsl:otherwise>1</xsl:otherwise>
                </xsl:choose>
            </xsl:when>
            <xsl:otherwise>8</xsl:otherwise>
        </xsl:choose>
    </xsl:function>

    <!-- Parse the fields for each type, rearranging all of the bit based fields so their order matches that of the PLC4X mspec -->
    <xsl:template name="plc4x:parseFields">
        <xsl:param name="baseNode"/>
        <xsl:param name="servicesRoot"/>
        <xsl:param name="currentNodePosition" as="xs:int"/>
        <xsl:param name="currentBitPosition" as="xs:int"/>
        <xsl:param name="currentBytePosition" as="xs:int"/>
        <xsl:choose>
            <xsl:when test="$currentNodePosition > count($baseNode/opc:Field)">
                <xsl:choose>
                    <xsl:when test="$currentBitPosition != 0">
                        <!-- Add a reserved field if we are halfway through a Byte.  -->
                        <xsl:call-template name="plc4x:parseFields">
                            <xsl:with-param name="baseNode">
                                <xsl:copy-of select="$baseNode/opc:Field[position() lt ($currentNodePosition - $currentBytePosition)]"/>
                                <xsl:element name="opc:Field">
                                    <xsl:attribute name="Name">ReservedX</xsl:attribute>
                                    <xsl:attribute name="TypeName">opc:Bit</xsl:attribute>
                                    <xsl:attribute name="Length"><xsl:value-of select="8-$currentBitPosition"/></xsl:attribute>
                                </xsl:element>
                                <xsl:copy-of select="$baseNode/opc:Field[(position() gt ($currentNodePosition - $currentBytePosition - 1))]"/>
                            </xsl:with-param>
                            <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                            <xsl:with-param name="currentNodePosition">
                                <xsl:value-of select="$currentNodePosition + 2"/>
                            </xsl:with-param>
                            <xsl:with-param name="currentBitPosition">
                                <xsl:value-of select="0"/>
                            </xsl:with-param>
                            <xsl:with-param name="currentBytePosition">
                                <xsl:value-of select="0"/>
                            </xsl:with-param>
                        </xsl:call-template>
                    </xsl:when>
                    <xsl:otherwise>
                        <!-- Return the rearranged nodes -->
                        <xsl:apply-templates select="$baseNode/opc:Field">
                            <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                            <xsl:with-param name="fieldCount" select="$currentNodePosition"/>
                        </xsl:apply-templates>
                    </xsl:otherwise>
                </xsl:choose>
            </xsl:when>
            <xsl:otherwise>
                <xsl:choose>
                    <xsl:when test="plc4x:getDataTypeLength($dataTypeLength, $baseNode/opc:Field[$currentNodePosition][@TypeName]) lt 8">
                        <xsl:choose>
                            <xsl:when test="$currentBitPosition=0">
                                <!-- Put node into current position -->
                                <xsl:call-template name="plc4x:parseFields">
                                    <xsl:with-param name="baseNode">
                                        <xsl:copy-of select="$baseNode/opc:Field"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                                    <xsl:with-param name="currentNodePosition">
                                        <xsl:value-of select="$currentNodePosition + 1"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBitPosition">
                                        <xsl:value-of select="plc4x:getDataTypeLength($dataTypeLength, $baseNode/opc:Field[position() = $currentNodePosition][@TypeName]) + $currentBitPosition"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBytePosition">
                                        <xsl:value-of select="$currentBytePosition + 1"/>
                                    </xsl:with-param>
                                </xsl:call-template>
                            </xsl:when>
                            <xsl:otherwise>
                                <!-- Put node into correct position based on bit and byte position -->
                                <xsl:call-template name="plc4x:parseFields">
                                    <xsl:with-param name="baseNode">
                                        <xsl:copy-of select="$baseNode/opc:Field[position() lt ($currentNodePosition - $currentBytePosition)]"/>
                                        <xsl:copy-of select="$baseNode/opc:Field[position() = $currentNodePosition]"/>
                                        <xsl:copy-of select="$baseNode/opc:Field[(position() gt ($currentNodePosition - $currentBytePosition - 1)) and (position() lt ($currentNodePosition))]"/>
                                        <xsl:copy-of select="$baseNode/opc:Field[position() gt $currentNodePosition]"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                                    <xsl:with-param name="currentNodePosition">
                                        <xsl:value-of select="$currentNodePosition + 1"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBitPosition">
                                        <xsl:value-of select="plc4x:getDataTypeLength($dataTypeLength, $baseNode/opc:Field[position() = $currentNodePosition][@TypeName]) + $currentBitPosition"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBytePosition">
                                        <xsl:value-of select="$currentBytePosition + 1"/>
                                    </xsl:with-param>
                                </xsl:call-template>
                            </xsl:otherwise>
                        </xsl:choose>
                    </xsl:when>
                    <xsl:otherwise>
                        <xsl:choose>
                            <xsl:when test="$currentBitPosition != 0 and $currentBitPosition lt 8">
                                <!-- Add a reserved field if we are halfway through a Byte.  -->
                                <xsl:call-template name="plc4x:parseFields">
                                    <xsl:with-param name="baseNode">
                                        <xsl:copy-of select="$baseNode/opc:Field[position() lt ($currentNodePosition - $currentBytePosition)]"/>
                                        <xsl:element name="opc:Field">
                                            <xsl:attribute name="Name">ReservedX</xsl:attribute>
                                            <xsl:attribute name="TypeName">opc:Bit</xsl:attribute>
                                            <xsl:attribute name="Length"><xsl:value-of select="8-$currentBitPosition"/></xsl:attribute>
                                        </xsl:element>
                                        <xsl:copy-of select="$baseNode/opc:Field[(position() gt ($currentNodePosition - $currentBytePosition - 1))]"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                                    <xsl:with-param name="currentNodePosition">
                                        <xsl:value-of select="$currentNodePosition + 2"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBitPosition">
                                        <xsl:value-of select="0"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBytePosition">
                                        <xsl:value-of select="0"/>
                                    </xsl:with-param>
                                </xsl:call-template>
                            </xsl:when>
                            <xsl:otherwise>
                                <!-- Put node into current position -->
                                <xsl:call-template name="plc4x:parseFields">
                                    <xsl:with-param name="baseNode">
                                        <xsl:copy-of select="$baseNode/opc:Field"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="servicesRoot" select="$servicesRoot"/>
                                    <xsl:with-param name="currentNodePosition">
                                        <xsl:value-of select="$currentNodePosition + 1"/>
                                    </xsl:with-param>
                                    <xsl:with-param name="currentBitPosition">0</xsl:with-param>
                                    <xsl:with-param name="currentBytePosition">0</xsl:with-param>
                                </xsl:call-template>
                            </xsl:otherwise>
                        </xsl:choose>
                    </xsl:otherwise>
                </xsl:choose>
            </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template name="plc4x:getExtensionIdArgs">
        <xsl:param name="dataTypeName"/>
        <xsl:param name="servicesRoot"/>
        <xsl:variable name="extensionId">
            <xsl:call-template name="plc4x:getExtensionId">
                <xsl:with-param name="dataTypeName" select="$dataTypeName" />
                <xsl:with-param name="servicesRoot" select="$servicesRoot" />
            </xsl:call-template>
        </xsl:variable>
        <xsl:if test="$extensionId != ''">('<xsl:value-of select="$extensionId" />')</xsl:if>
    </xsl:template>

    <xsl:template name="plc4x:getExtensionId">
        <xsl:param name="dataTypeName"/>
        <xsl:param name="servicesRoot"/>
        <xsl:if test="$originaldoc/opc:TypeDictionary/opc:StructuredType[@Name=$dataTypeName]/@BaseType">
            <xsl:variable name="datatypeDefinition" select="$servicesRoot/node:UADataType[@BrowseName=$dataTypeName]"/>
            <xsl:comment>The +2 comes from offset between type definition and binary encoding id.</xsl:comment>
            <xsl:value-of select="xs:int(substring($datatypeDefinition/@NodeId, 3)) + 2" />
        </xsl:if>
    </xsl:template>

</xsl:stylesheet>