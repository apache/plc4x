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
                xmlns:knx="http://knx.org/xml/project/20">

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

[enum uint 8 'KnxObjectType' [string 'text']<xsl:apply-templates select="knx:KNX/knx:MasterData/knx:InterfaceObjectTypes/knx:InterfaceObjectType"/>
]

[enum uint 8 'KnxObjectProperties' [string 'name', string 'dataTypeId', string 'text']<xsl:apply-templates select="knx:KNX/knx:MasterData/knx:InterfaceObjectProperties/knx:InterfaceObjectProperty"/>
]

[enum uint 16 'KnxManufacturers' [string 'text']<xsl:apply-templates select="knx:KNX/knx:MasterData/knx:Manufacturers/knx:Manufacturer"/>
]
    </xsl:template>

    <xsl:template match="knx:InterfaceObjectType">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Number"/>' <xsl:value-of select="$objectTypeId"/> ['"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:InterfaceObjectProperty">
        <xsl:variable name="dataTypeId">
            <xsl:choose>
                <xsl:when test="contains(@PDT, ' ')">
                    <xsl:call-template name="clean-id-string">
                        <xsl:with-param name="text" select="substring-before(@PDT, ' ')"/>
                    </xsl:call-template>
                </xsl:when>
                <xsl:when test="string-length(@PDT) &gt; 0">
                    <xsl:call-template name="clean-id-string">
                        <xsl:with-param name="text" select="@PDT"/>
                    </xsl:call-template>
                </xsl:when>
                <xsl:otherwise>null</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="interfaceObjectPropertyId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text">
                    <xsl:call-template name="string-replace-all">
                        <xsl:with-param name="text" select="@Id"/>
                        <xsl:with-param name="replace" select="'-'"/>
                        <xsl:with-param name="with" select="'_'"/>
                    </xsl:call-template>
                </xsl:with-param>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Number"/>' <xsl:value-of select="$interfaceObjectPropertyId"/> ['<xsl:value-of select="@Name"/>', '<xsl:value-of select="$dataTypeId"/>', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:Manufacturer">
        <xsl:variable name="manufacturerId">
            <xsl:call-template name="clean-id-string">
                <xsl:with-param name="text">
                    <xsl:call-template name="string-replace-all">
                        <xsl:with-param name="text" select="@Id"/>
                        <xsl:with-param name="replace" select="'-'"/>
                        <xsl:with-param name="with" select="'_'"/>
                    </xsl:call-template>
                </xsl:with-param>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@KnxManufacturerId"/>' <xsl:value-of select="$manufacturerId"/> ['"<xsl:value-of select="@Name"/>"']]
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