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
                xmlns:fn="http://www.w3.org/2005/xpath-functions"
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
[enum uint 16 'KnxInterfaceObjectType' [string 'text']
    ['65535' OT_GENERAL ['"General Object"']]
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:InterfaceObjectTypes/knx:InterfaceObjectType"/>
]

[enum uint 32 'KnxInterfaceObjectProperty' [uint16 'propertyId', KnxInterfaceObjectType 'objectType', KnxPropertyDataType 'propertyDataType', string 'text']
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:InterfaceObjectProperties/knx:InterfaceObjectProperty"/>
]

[enum uint 8 'KnxPropertyDataType' [uint 8 'sizeInBytes']
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:PropertyDataTypes/knx:PropertyDataType"/>
]

[enum uint 16 'KnxDatapointType' [string 'text', uint 8 'sizeInBits']
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:DatapointTypes/knx:DatapointType"/>
]

[enum uint 32 'KnxDatapointSubtype' [KnxDatapointType 'datapointType', string 'text']
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:DatapointTypes/knx:DatapointType/knx:DatapointSubtypes/knx:DatapointSubtype"/>
]

[enum uint 16 'KnxManufacturer' [string 'text']
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:Manufacturers/knx:Manufacturer"/>
]
    </xsl:template>

    <xsl:template match="knx:InterfaceObjectType">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="getIdFromText">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Number"/>' <xsl:value-of select="$objectTypeId"/> ['"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:InterfaceObjectProperty">
        <xsl:variable name="objectTypeNumber">
            <xsl:choose>
                <xsl:when test="fn:starts-with(@Id, 'PID-G')">65535</xsl:when>
                <xsl:otherwise>
                    <xsl:value-of select="fn:substring-before(fn:substring-after(@Id, 'PID-'), '-')"/>
                </xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="objectPropertyValue" select="format-number(fn:number(concat($objectTypeNumber, format-number(@Number, '000'))), '#####0')"/>
        <xsl:variable name="objectPropertyId">
            <xsl:variable name="objectTypeId" select="@ObjectType"/>
            <xsl:choose>
                <xsl:when test="fn:starts-with(@Id, 'PID-G')">PID_GENERAL_<xsl:value-of select="fn:substring-after(@Name, 'PID_')"/></xsl:when>
                <xsl:otherwise>PID_<xsl:value-of select="fn:substring-after(//knx:InterfaceObjectType[@Id = $objectTypeId]/@Name, 'OT_')"/>_<xsl:value-of select="fn:replace(fn:substring-after(@Name, '_'), '%', 'PERCENT')"/><xsl:if test="@Id = 'PID-417-120'">_2</xsl:if></xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="objectType">
            <xsl:variable name="objectTypeId" select="@ObjectType"/>
            <xsl:choose>
                <xsl:when test="@ObjectType">KnxInterfaceObjectType.<xsl:value-of select="//knx:InterfaceObjectType[@Id = $objectTypeId]/@Name"/></xsl:when>
                <xsl:otherwise>KnxInterfaceObjectType.OT_GENERAL</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>['<xsl:value-of select="$objectPropertyValue"/>' <xsl:value-of select="$objectPropertyId"/> ['<xsl:value-of select="@Number"/>', '<xsl:value-of select="$objectType"/>', 'KnxPropertyDataType.PDT_CONTROL', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:PropertyDataType">
        <xsl:variable name="sizeInBytes">
            <xsl:choose>
                <xsl:when test="@Size"><xsl:value-of select="@Size"/></xsl:when>
                <xsl:otherwise>0</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>['<xsl:value-of select="@Number"/>' <xsl:value-of select="fn:replace(@Name, '-', '_')"/> ['<xsl:value-of select="$sizeInBytes"/>']]
    </xsl:template>

    <xsl:template match="knx:DatapointType">
        <xsl:variable name="datapointTypeId">
            <xsl:call-template name="getDatapointTypeId">
                <xsl:with-param name="contextNode" select="."/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="@Number"/>' <xsl:value-of select="$datapointTypeId"/> ['"<xsl:value-of select="@Text"/>"', '<xsl:value-of select="@SizeInBit"/>']]
    </xsl:template>

    <xsl:template match="knx:DatapointSubtype">
        <xsl:variable name="datapointSubtypeValue" select="concat(../../@Number, format-number(@Number, '0000'))"/>
        <xsl:variable name="datapointSubtypeId">
            <xsl:choose>
                <xsl:when test="fn:starts-with(@Name, 'DPT')">DPST_<xsl:value-of select="fn:substring-after(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(@Name, '\[', '_'), '\]', ''), '&#x00B3;', '_3'), '&#xB5;', 'y'), '/', ''), '-', '_'), 'DPT_')"/></xsl:when>
                <xsl:otherwise>DPST_<xsl:for-each select="tokenize(@Name, ' ')"><xsl:value-of select="concat(upper-case(substring(.,1,1)), substring(., 2))"/><xsl:if test="position()!=last()">_</xsl:if></xsl:for-each></xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="datapointTypeId">
            <xsl:call-template name="getDatapointTypeId">
                <xsl:with-param name="contextNode" select="../.."/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="$datapointSubtypeValue"/>' <xsl:value-of select="$datapointSubtypeId"/> ['KnxDatapointType.<xsl:value-of select="$datapointTypeId"/>', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:Manufacturer">
        <xsl:variable name="manufacturerId">
            <xsl:choose>
                <xsl:when test="@Name = 'Simon'">SIMON_<xsl:value-of select="@KnxManufacturerId"/></xsl:when>
                <xsl:when test="@Name = 'Not Assigned'">NOT_ASSIGNED_<xsl:value-of select="@KnxManufacturerId"/></xsl:when>
                <xsl:otherwise>
                    <xsl:call-template name="getIdFromText">
                        <xsl:with-param name="text" select="@Name"/>
                    </xsl:call-template>
                </xsl:otherwise>
            </xsl:choose>
        </xsl:variable>['<xsl:value-of select="@KnxManufacturerId"/>' <xsl:value-of select="$manufacturerId"/> ['"<xsl:value-of select="fn:replace(normalize-space(@Name),'&#xa0;', ' ')"/>"']]
    </xsl:template>

    <xsl:template name="getDatatypeId">
        <xsl:param name="contextNode"/>
        <xsl:value-of select="fn:replace(fn:replace(fn:replace(fn:replace(fn:upper-case($contextNode/@Text), ' ', '_'), '-', '_'), '&amp;', 'AND'), '/', '')"/>
    </xsl:template>

    <xsl:template name="getDatapointTypeId">
        <xsl:param name="contextNode"/>
        <xsl:choose>
            <xsl:when test="$contextNode/@Number = 26">DPT_8_BIT_SET_2</xsl:when>
            <xsl:when test="$contextNode/@Number = 237">DPT_CONFIGURATION_DIAGNOSTICS_16_BIT</xsl:when>
            <xsl:when test="$contextNode/@Number = 238">DPT_CONFIGURATION_DIAGNOSTICS_8_BIT</xsl:when>
            <xsl:when test="$contextNode/@Number = 241">DPT_STATUS_32_BIT</xsl:when>
            <xsl:when test="$contextNode/@Number = 242">DPT_STATUS_48_BIT</xsl:when>
            <xsl:when test="$contextNode/@Number = 250">DPT_STATUS_24_BIT</xsl:when>
            <xsl:otherwise>DPT_<xsl:value-of select="fn:replace(fn:replace(fn:replace(fn:replace(fn:upper-case($contextNode/@Text), ' ', '_'), '-', '_'), '&amp;', 'AND'), '/', '')"/></xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template name="getIdFromText">
        <xsl:param name="text"/>
        <xsl:variable name="cleanedText" select="fn:replace(fn:replace(fn:replace(fn:replace(fn:upper-case($text), '/', ''), '\(', ''), '\)', ''), '&#x2013;', '_')"/>
        <xsl:variable name="cleanedText2" select="fn:replace(fn:replace($cleanedText, '&#x00B3;', '_3'), '&#xC9;', 'E')"/>
        <xsl:value-of select="fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(normalize-space($cleanedText2),'&#xa0;', '_'), '&amp;', 'AND'), '-', '_'), ' ', '_'), '\.', '_'), ',', '_'), '\+', 'Plus'), '/', '_'), 'Ä', 'AE'), 'Ö', 'OE'), 'Ü', 'UE'), 'ß', 'SS'), ':', '_')"/>
    </xsl:template>

</xsl:stylesheet>