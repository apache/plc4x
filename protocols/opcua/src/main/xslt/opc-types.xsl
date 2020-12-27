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
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>[enum int 32 '<xsl:value-of select="$objectTypeId"/>'
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
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Value"/>' <xsl:value-of select="$objectTypeId"/>]
    </xsl:template>

    <xsl:template match="opc:OpaqueType[not(@Name = 'Duration')]">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>[type '<xsl:value-of select="$objectTypeId"/>'
        <xsl:apply-templates select="opc:Documentation"/>
        [simple int 32 'testing']
  ]
    </xsl:template>

    <xsl:template match="opc:StructuredType[not(@Name = 'Vector')]">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>[type '<xsl:value-of select="$objectTypeId"/>'
    <xsl:apply-templates select="opc:Documentation"/>
    <xsl:apply-templates select="opc:Field"/>
  ]
    </xsl:template>

    <xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz'" />
    <xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ'" />

    <xsl:template match="opc:Field">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="lowerCaseName">
            <xsl:value-of select="concat(translate(substring(@Name, 1, 1), $uppercase, $lowercase), substring(@Name, 2))"/>
        </xsl:variable>
        <xsl:variable name="dataType">
            <xsl:choose>
                <xsl:when test="@TypeName = 'opc:Bit'">bit</xsl:when>
                <xsl:when test="@TypeName = 'opc:Boolean'">bit</xsl:when>
                <xsl:when test="@TypeName = 'opc:Byte'">uint 8</xsl:when>
                <xsl:when test="@TypeName = 'opc:SByte'">int 8</xsl:when>
                <xsl:when test="@TypeName = 'opc:Int16'">int 16</xsl:when>
                <xsl:when test="@TypeName = 'opc:UInt16'">uint 16</xsl:when>
                <xsl:when test="@TypeName = 'opc:Int32'">int 32</xsl:when>
                <xsl:when test="@TypeName = 'opc:UInt32'">uint 32</xsl:when>
                <xsl:when test="@TypeName = 'opc:Int64'">int 64</xsl:when>
                <xsl:when test="@TypeName = 'opc:UInt64'">uint 64</xsl:when>
                <xsl:when test="@TypeName = 'opc:Float'">float 8.23</xsl:when>
                <xsl:when test="@TypeName = 'opc:Double'">float 11.52</xsl:when>
                <xsl:when test="@TypeName = 'opc:Char'">string '<xsl:value-of select="concat(translate(substring(@LengthField, 1, 1), $uppercase, $lowercase), substring(@LengthField, 2))"/>'</xsl:when>
                <xsl:when test="@TypeName = 'opc:CharArray'">string '-1'</xsl:when>
                <xsl:when test="@TypeName = 'opc:Guid'">string '-1'</xsl:when>
                <xsl:when test="@TypeName = 'opc:ByteString'">string '-1'</xsl:when>
                <xsl:when test="@TypeName = 'opc:DateTime'">string '-1'</xsl:when>
                <xsl:when test="@TypeName = 'opc:String'">string</xsl:when>
                <xsl:otherwise><xsl:value-of select="substring-after(@TypeName,':')"/></xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
            <xsl:choose>
                <xsl:when test="@LengthField and @TypeName = 'opc:String'">[array <xsl:value-of select="$dataType"/> '<xsl:value-of select="$lowerCaseName"/>' '<xsl:value-of select="$lowerCaseName"/>' count '<xsl:value-of select="concat(translate(substring(@LengthField, 1, 1), $uppercase, $lowercase), substring(@LengthField, 2))"/>']
                </xsl:when>
                <xsl:when test="@LengthField">[array <xsl:value-of select="$dataType"/>  '<xsl:value-of select="$lowerCaseName"/>' count '<xsl:value-of select="concat(translate(substring(@LengthField, 1, 1), $uppercase, $lowercase), substring(@LengthField, 2))"/>']
                </xsl:when>
                <xsl:when test="@TypeName = 'opc:String'">[simple int 32 '<xsl:value-of select="$lowerCaseName"/>Size']
                [simple <xsl:value-of select="$dataType"/> '<xsl:value-of select="$lowerCaseName"/>Size' '<xsl:value-of select="$objectTypeId"/>']
                </xsl:when>
                <xsl:otherwise>[simple <xsl:value-of select="$dataType"/> '<xsl:value-of select="$lowerCaseName"/>']
                </xsl:otherwise>
            </xsl:choose>
    </xsl:template>

    <xsl:template name="clean-id-string">
        <xsl:param name="text"/>
        <xsl:call-template name="string-replace-all">
            <xsl:with-param name="text" select="$text"/>
            <xsl:with-param name="replace" select="'%'"/>
            <xsl:with-param name="with" select="'PERCENT'"/>
        </xsl:call-template>
    </xsl:template>

    <xsl:template name="string-replace-all">
        <xsl:param name="text"/>
        <xsl:param name="replace"/>
        <xsl:param name="with"/>
        <xsl:choose>
            <xsl:when test="$text = '' or $replace = ''or not($replace)">
                <xsl:value-of select="$text"/>
            </xsl:when>
            <xsl:when test="contains($text, $replace)">
                <xsl:value-of select="substring-before($text,$replace)"/>
                <xsl:value-of select="$with"/>
                <xsl:call-template name="string-replace-all">
                    <xsl:with-param name="text" select="substring-after($text,$replace)"/>
                    <xsl:with-param name="replace" select="$replace"/>
                    <xsl:with-param name="with" select="$with"/>
                </xsl:call-template>
            </xsl:when>
            <xsl:otherwise>
                <xsl:value-of select="$text"/>
            </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

</xsl:stylesheet>
