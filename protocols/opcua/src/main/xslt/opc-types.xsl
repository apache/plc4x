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
<xsl:stylesheet version="1.0"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:opc="http://opcfoundation.org/BinarySchema/"
                xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
                xmlns:ua="http://opcfoundation.org/UA/"
                xmlns:tns="http://opcfoundation.org/UA/">

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />

    <xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz'" />
    <xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ'" />

    <xsl:template match="/">
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.
//
        <xsl:apply-templates select="opc:TypeDictionary"/>
    </xsl:template>

    <xsl:template match="opc:EnumeratedType">
        [enum int 32 '<xsl:value-of select="@Name"/>'
        <xsl:apply-templates select="opc:Documentation"/>
        <xsl:apply-templates select="opc:EnumeratedValue"/>
  ]
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
        <xsl:choose>
            <xsl:when test="*[@SwitchField]">[discriminatedType '<xsl:value-of select="@Name"/>'</xsl:when>
            <xsl:when test="not(*[@SwitchField])">[type '<xsl:value-of select="@Name"/>'</xsl:when>
        </xsl:choose>
    <xsl:apply-templates select="opc:Documentation"/>
    <xsl:apply-templates select="opc:Field"/>
  ]
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
            <xsl:when test="@LengthField and @TypeName = 'opc:String'">[array <xsl:value-of select="$dataType"/> '-1' '<xsl:value-of select="$lowerCaseName"/>' count '<xsl:value-of select="$lowerCaseLengthField"/>']
            </xsl:when>
            <xsl:when test="@LengthField">[array <xsl:value-of select="$dataType"/>  '<xsl:value-of select="$lowerCaseName"/>' count '<xsl:value-of select="$lowerCaseLengthField"/>']
            </xsl:when>
            <xsl:when test="@TypeName = 'opc:String'">[simple int 32 '<xsl:value-of select="$lowerCaseName"/>Size']
            [simple <xsl:value-of select="$dataType"/> '<xsl:value-of select="$lowerCaseName"/>Size' '<xsl:value-of select="@Name"/>']
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
            <xsl:when test="$datatype = 'opc:Char'">string '-1'</xsl:when>
            <xsl:when test="$datatype = 'opc:CharArray'">string '-1'</xsl:when>
            <xsl:when test="$datatype = 'opc:Guid'">string '-1'</xsl:when>
            <xsl:when test="$datatype = 'opc:ByteString'">string '-1'</xsl:when>
            <xsl:when test="$datatype = 'opc:DateTime'">string '-1'</xsl:when>
            <xsl:when test="$datatype = 'opc:String'">string</xsl:when>
            <xsl:otherwise><xsl:value-of select="substring-after($datatype,':')"/></xsl:otherwise>
        </xsl:choose>
    </xsl:template>

</xsl:stylesheet>
