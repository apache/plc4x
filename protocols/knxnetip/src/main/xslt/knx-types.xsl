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
/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
[enum uint 16 KnxDatapointMainType(uint 16 number, uint 8 sizeInBits, string 8 name)
    ['0' DPT_UNKNOWN ['0', '0', '"Unknown Datapoint Type"']]
    // Begin: Some typed needed to support all IEC types on KNX (Which the standard generally doesn't support)
    ['1' DPT_64_BIT_SET ['0', '64', '"Unknown Datapoint Type"']]
    ['2' DPT_8_BYTE_UNSIGNED_VALUE ['0', '64', '"Unknown Datapoint Type"']]
    ['3' DPT_8_BYTE_SIGNED_VALUE ['0', '64', '"Unknown Datapoint Type"']]
    ['4' DPT_12_BYTE_SIGNED_VALUE ['0', '96', '"Unknown Datapoint Type"']]
    ['5' DPT_8_BYTE_FLOAT_VALUE ['0', '64', '"Unknown Datapoint Type"']]
    // End: Custom types
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:DatapointTypes/knx:DatapointType"/>
]

[enum uint 32 KnxDatapointType(uint 16 number, KnxDatapointMainType datapointMainType, string 8 name)
    ['0' DPT_UNKNOWN    ['0', 'DPT_UNKNOWN',               '"Unknown Datapoint Subtype"']]
    // Begin: Some typed needed to support all IEC types on KNX (Which the standard generally doesn't support)
    ['1' BOOL           ['0', 'DPT_1_BIT',                 '"BOOL"']]
    ['2' BYTE           ['0', 'DPT_8_BIT_SET',             '"BYTE"']]
    ['3' WORD           ['0', 'DPT_16_BIT_SET',            '"WORD"']]
    ['4' DWORD          ['0', 'DPT_32_BIT_SET',            '"DWORD"']]
    ['5' LWORD          ['0', 'DPT_64_BIT_SET',            '"LWORD"']]
    ['6' USINT          ['0', 'DPT_8_BIT_UNSIGNED_VALUE',  '"USINT"']]
    ['7' SINT           ['0', 'DPT_8_BIT_SIGNED_VALUE',    '"SINT"']]
    ['8' UINT           ['0', 'DPT_2_BYTE_UNSIGNED_VALUE', '"UINT"']]
    ['9' INT            ['0', 'DPT_2_BYTE_SIGNED_VALUE',   '"INT"']]
    ['10' UDINT         ['0', 'DPT_4_BYTE_UNSIGNED_VALUE', '"UDINT"']]
    ['11' DINT          ['0', 'DPT_4_BYTE_SIGNED_VALUE',   '"DINT"']]
    ['12' ULINT         ['0', 'DPT_8_BYTE_UNSIGNED_VALUE', '"ULINT"']]
    ['13' LINT          ['0', 'DPT_8_BYTE_SIGNED_VALUE',   '"LINT"']]
    ['14' REAL          ['0', 'DPT_4_BYTE_FLOAT_VALUE',    '"REAL"']]
    ['15' LREAL         ['0', 'DPT_8_BYTE_FLOAT_VALUE',    '"LREAL"']]
    ['16' CHAR          ['0', 'DPT_CHARACTER',             '"CHAR"']]
    ['17' WCHAR         ['0', 'DPT_2_BYTE_UNSIGNED_VALUE', '"WCHAR"']]
    ['18' STRING        ['0', 'DPT_UNKNOWN',               '"STRING"']]
    ['19' WSTRING       ['0', 'DPT_UNKNOWN',               '"WSTRING"']]
    ['20' TIME          ['0', 'DPT_4_BYTE_UNSIGNED_VALUE', '"TIME"']]
    ['21' LTIME         ['0', 'DPT_8_BYTE_UNSIGNED_VALUE', '"LTIME"']]
    ['22' DATE          ['0', 'DPT_2_BYTE_UNSIGNED_VALUE', '"DATE"']]
    ['23' TIME_OF_DAY   ['0', 'DPT_4_BYTE_UNSIGNED_VALUE', '"TIME_OF_DAY"']]
    ['24' TOD           ['0', 'DPT_4_BYTE_UNSIGNED_VALUE', '"TOD"']]
    ['25' DATE_AND_TIME ['0', 'DPT_12_BYTE_SIGNED_VALUE',  '"DATE_AND_TIME"']]
    ['26' DT            ['0', 'DPT_12_BYTE_SIGNED_VALUE',  '"DT"']]
    // End: Custom types
        <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:DatapointTypes/knx:DatapointType/knx:DatapointSubtypes/knx:DatapointSubtype"/>
]

[enum uint 16 KnxInterfaceObjectType(string 8 code, string 8 name)
    ['0' OT_UNKNOWN ['U', '"Unknown Interface Object Type"']]
    ['1' OT_GENERAL ['G', '"General Interface Object Type"']]
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:InterfaceObjectTypes/knx:InterfaceObjectType"/>
]

[enum uint 32 KnxInterfaceObjectProperty(uint 8 propertyId, KnxInterfaceObjectType objectType, KnxPropertyDataType propertyDataType, string 8 name)
    ['0' PID_UNKNOWN    ['0', 'OT_UNKNOWN', 'PDT_UNKNOWN', '"Unknown Interface Object Property"']]
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:InterfaceObjectProperties/knx:InterfaceObjectProperty"/>
]

[enum uint 8 KnxPropertyDataType(uint 8 number, uint 8 sizeInBytes, string 8 name)
    ['0' PDT_UNKNOWN    ['0', '0',  '"Unknown Property Data Type"']]
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:PropertyDataTypes/knx:PropertyDataType"/>
]

[enum uint 16 KnxManufacturer(uint 16 number, string 8 name)
    ['0' M_UNKNOWN ['0', '"Unknown Manufacturer"']]
    <xsl:apply-templates select="knx:KNX/knx:MasterData/knx:Manufacturers/knx:Manufacturer"/>
]

[dataIo KnxDatapoint(KnxDatapointType datapointType)
    [typeSwitch 'datapointType'
        ['BOOL' BOOL
            [reserved uint 7 '0x00']
            [simple   bit    value]
        ]
        ['BYTE' BYTE
            [simple   uint 8    value]
        ]
        ['WORD' WORD
            [simple   uint 16    value]
        ]
        ['DWORD' DWORD
            [simple   uint 32    value]
        ]
        ['LWORD' LWORD
            [simple   uint 64    value]
        ]
        ['USINT' USINT
            [simple   uint 8     value]
        ]
        ['SINT' SINT
            [simple   int 8      value]
        ]
        ['UINT' UINT
            [simple   uint 16    value]
        ]
        ['INT' INT
            [simple   int 16     value]
        ]
        ['UDINT' UDINT
            [simple   uint 32    value]
        ]
        ['DINT' DINT
            [simple   int 32     value]
        ]
        ['ULINT' ULINT
            [simple   uint 64    value]
        ]
        ['LINT' LINT
            [simple   int 64     value]
        ]
        ['REAL' REAL
            [simple   float 32 value]
        ]
        ['LREAL' LREAL
            [simple   float 64 value]
        ]
        ['CHAR' CHAR
            [simple   uint 8     value]
        ]
        ['WCHAR' WCHAR
            [simple   uint 16    value]
        ]
        //['STRING' STRING
        //]
        //['WSTRING' WSTRING
        //]
        ['TIME' TIME
            [simple uint 32 value]
        ]
        ['LTIME' LTIME
            [simple uint 64 value]
        ]
        ['DATE' DATE
            [simple uint 16 value]
        ]
        ['TIME_OF_DAY' TIME_OF_DAY
            [simple uint 32 value]
        ]
        ['TOD' TIME_OF_DAY
            [simple uint 32 value]
        ]
        ['DATE_AND_TIME' DATE_AND_TIME
            [simple uint 16 year]
            [simple uint 8  month]
            [simple uint 8  day]
            [simple uint 8  dayOfWeek]
            [simple uint 8  hour]
            [simple uint 8  minutes]
            [simple uint 8  seconds]
            [simple uint 32 nanos]
        ]
        ['DT' DATE_AND_TIME
            [simple uint 16 year]
            [simple uint 8  month]
            [simple uint 8  day]
            [simple uint 8  dayOfWeek]
            [simple uint 8  hour]
            [simple uint 8  minutes]
            [simple uint 8  seconds]
            [simple uint 32 nanos]
        ]

    <xsl:for-each select="knx:KNX/knx:MasterData/knx:DatapointTypes/knx:DatapointType/knx:DatapointSubtypes/knx:DatapointSubtype">
        <xsl:call-template name="generateDataIoEntry">
            <xsl:with-param name="datapointSubtype" select="."/>
        </xsl:call-template>
    </xsl:for-each>
    ]
]
    </xsl:template>

    <xsl:template match="knx:DatapointType">
        <xsl:variable name="datapointTypeId">
            <xsl:call-template name="getDatapointTypeId">
                <xsl:with-param name="contextNode" select="."/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="position() + 5"/>' <xsl:value-of select="$datapointTypeId"/> ['<xsl:value-of select="@Number"/>', '<xsl:value-of select="@SizeInBit"/>', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:DatapointSubtype">
        <xsl:variable name="datapointSubtypeId">
            <xsl:choose>
                <xsl:when test="fn:starts-with(@Name, 'DPT')">DPT_<xsl:value-of select="fn:substring-after(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(@Name, '\[', '_'), '\]', ''), '&#x00B3;', '_3'), '&#xB5;', 'y'), '/', ''), '-', '_'), 'DPT_')"/></xsl:when>
                <xsl:otherwise>DPT_<xsl:for-each select="tokenize(@Name, ' ')"><xsl:value-of select="concat(upper-case(substring(.,1,1)), substring(., 2))"/><xsl:if test="position()!=last()">_</xsl:if></xsl:for-each></xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="datapointTypeId">
            <xsl:call-template name="getDatapointTypeId">
                <xsl:with-param name="contextNode" select="../.."/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="position() + 26"/>' <xsl:value-of select="$datapointSubtypeId"/> ['<xsl:value-of select="@Number"/>', '<xsl:value-of select="$datapointTypeId"/>', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:InterfaceObjectType">
        <xsl:variable name="objectTypeId">
            <xsl:call-template name="getIdFromText">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>['<xsl:value-of select="position() + 1"/>' <xsl:value-of select="$objectTypeId"/> ['<xsl:value-of select="@Number"/>', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:InterfaceObjectProperty">
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
                <xsl:when test="@ObjectType"><xsl:value-of select="//knx:InterfaceObjectType[@Id = $objectTypeId]/@Name"/></xsl:when>
                <xsl:otherwise>OT_GENERAL</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="propertyDataTypeName">
            <xsl:choose>
                <xsl:when test="@PDT and not(fn:contains(@PDT, ' '))">
                    <xsl:variable name="pdtId" select="@PDT"/>
                    <xsl:value-of select="//knx:PropertyDataType[@Id=$pdtId]/@Name"/>
                </xsl:when>
                <xsl:when test="@PDT">
                    <xsl:variable name="pdtId" select="substring-before(@PDT, ' ')"/>
                    <xsl:value-of select="//knx:PropertyDataType[@Id=$pdtId]/@Name"/>
                </xsl:when>
                <xsl:otherwise>PDT_UNKNOWN</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>['<xsl:value-of select="position()"/>' <xsl:value-of select="$objectPropertyId"/> ['<xsl:value-of select="@Number"/>', '<xsl:value-of select="$objectType"/>', '<xsl:value-of select="$propertyDataTypeName"/>', '"<xsl:value-of select="@Text"/>"']]
    </xsl:template>

    <xsl:template match="knx:PropertyDataType">
        <xsl:variable name="sizeInBytes">
            <xsl:choose>
                <xsl:when test="@Size"><xsl:value-of select="@Size"/></xsl:when>
                <xsl:otherwise>0</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>['<xsl:value-of select="position()"/>' <xsl:value-of select="fn:replace(@Name, '-', '_')"/> ['<xsl:value-of select="@Number"/>', '<xsl:value-of select="$sizeInBytes"/>', '"<xsl:value-of select="@Name"/>"']]
    </xsl:template>

    <xsl:template match="knx:Manufacturer">
        <xsl:variable name="manufacturerId">
            <xsl:choose>
                <xsl:when test="@Name = '3ATEL'">THREEATEL</xsl:when>
                <xsl:when test="@Name = '1Home'">ONEHOME</xsl:when>
                <xsl:when test="@Name = 'Simon'">SIMON_<xsl:value-of select="@KnxManufacturerId"/></xsl:when>
                <xsl:when test="@Name = 'Not Assigned'">NOT_ASSIGNED_<xsl:value-of select="@KnxManufacturerId"/></xsl:when>
                <xsl:otherwise>
                    <xsl:call-template name="getIdFromText">
                        <xsl:with-param name="text" select="@Name"/>
                    </xsl:call-template>
                </xsl:otherwise>
            </xsl:choose>
        </xsl:variable>['<xsl:value-of select="position()"/>' M_<xsl:value-of select="$manufacturerId"/> ['<xsl:value-of select="@KnxManufacturerId"/>', '"<xsl:value-of select="fn:replace(normalize-space(@Name),'&#xa0;', ' ')"/>"']]
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

    <xsl:template name="generateDataIoEntry">
        <xsl:param name="datapointSubtype"/>
        <xsl:variable name="datapointSubtypeId">
            <xsl:choose>
                <xsl:when test="fn:starts-with(@Name, 'DPT')">DPT_<xsl:value-of select="fn:substring-after(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(fn:replace(@Name, '\[', '_'), '\]', ''), '&#x00B3;', '_3'), '&#xB5;', 'y'), '/', ''), '-', '_'), 'DPT_')"/></xsl:when>
                <xsl:otherwise>DPT_<xsl:for-each select="tokenize(@Name, ' ')"><xsl:value-of select="concat(upper-case(substring(.,1,1)), substring(., 2))"/><xsl:if test="position()!=last()">_</xsl:if></xsl:for-each></xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        <xsl:variable name="datapointValueType">
            <xsl:choose>
                <xsl:when test="count($datapointSubtype/knx:Format/knx:Bit|$datapointSubtype/knx:Format/knx:String|$datapointSubtype/knx:Format/knx:UnsignedInteger|$datapointSubtype/knx:Format/knx:SignedInteger|$datapointSubtype/knx:Format/knx:Float|$datapointSubtype/knx:Format/knx:Enumeration|$datapointSubtype/knx:Format/knx:RefType) &gt; 1">Struct</xsl:when>
                <xsl:when test="$datapointSubtype/knx:Format/knx:Bit">BOOL</xsl:when>
                <xsl:when test="$datapointSubtype/knx:Format/knx:String">STRING</xsl:when>
                <xsl:when test="$datapointSubtype/knx:Format/knx:UnsignedInteger">
                    <xsl:choose>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:UnsignedInteger/@Width) &lt;= 8">USINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:UnsignedInteger/@Width) &lt;= 16">UINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:UnsignedInteger/@Width) &lt;= 32">UDINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:UnsignedInteger/@Width) &lt;= 64">ULINT</xsl:when>
                    </xsl:choose>
                </xsl:when>
                <xsl:when test="$datapointSubtype/knx:Format/knx:SignedInteger">
                    <xsl:choose>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:SignedInteger/@Width) &lt;= 8">SINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:SignedInteger/@Width) &lt;= 16">INT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:SignedInteger/@Width) &lt;= 32">DINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:SignedInteger/@Width) &lt;= 64">LINT</xsl:when>
                    </xsl:choose>
                </xsl:when>
                <xsl:when test="$datapointSubtype/knx:Format/knx:Float">
                    <xsl:choose>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Float/@Width) &lt;= 16">REAL</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Float/@Width) &lt;= 32">REAL</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Float/@Width) &lt;= 64">LREAL</xsl:when>
                    </xsl:choose>
                </xsl:when>
                <xsl:when test="$datapointSubtype/knx:Format/knx:Enumeration">
                    <xsl:choose>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Enumeration/@Width) &lt;= 8">USINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Enumeration/@Width) &lt;= 16">UINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Enumeration/@Width) &lt;= 32">UDINT</xsl:when>
                        <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Enumeration/@Width) &lt;= 64">ULINT</xsl:when>
                    </xsl:choose>
                </xsl:when>
                <xsl:otherwise>Hurz:<xsl:value-of select="name($datapointSubtype/*)"/></xsl:otherwise>
            </xsl:choose>
        </xsl:variable>
        ['<xsl:value-of select="$datapointSubtypeId"/>' <xsl:value-of select="$datapointValueType"/>
        <xsl:choose>
            <xsl:when test="count($datapointSubtype/knx:Format/knx:Bit|$datapointSubtype/knx:Format/knx:String|$datapointSubtype/knx:Format/knx:UnsignedInteger|$datapointSubtype/knx:Format/knx:SignedInteger|$datapointSubtype/knx:Format/knx:Float|$datapointSubtype/knx:Format/knx:Enumeration|$datapointSubtype/knx:Format/knx:RefType) &gt; 1">
                <xsl:variable name="resolvedFields">
                    <xsl:call-template name="resolveTypeReferences">
                        <xsl:with-param name="fields" select="$datapointSubtype/knx:Format/*"/>
                    </xsl:call-template>
                </xsl:variable>
                <xsl:variable name="size">
                    <xsl:call-template name="fieldsSize">
                        <xsl:with-param name="fields" select="$resolvedFields/*"/>
                    </xsl:call-template>
                </xsl:variable>
                <xsl:if test="(($size mod 8) != 0) and (($size mod 8) &lt;= 6)">
            [reserved uint <xsl:value-of select="8 - ($size mod 8)"/> '0x00']
                </xsl:if>
                <xsl:for-each select="$resolvedFields/*">
                    <xsl:variable name="fieldType">
                        <xsl:choose>
                            <xsl:when test="name() = 'Reserved'">reserved</xsl:when>
                            <xsl:otherwise>simple</xsl:otherwise>
                        </xsl:choose>
                    </xsl:variable>
                    <xsl:variable name="fieldName">
                        <xsl:choose>
                            <xsl:when test="name() = 'Reserved'">'0x00'</xsl:when>
                            <!-- Special case for one enum -->
                            <xsl:when test="@Id = 'DPST-6-20_F-6'">mode</xsl:when>
                            <xsl:when test="(@Id = 'DPST-15-0_F-1') and fn:position() = 2">value1</xsl:when>
                            <xsl:when test="(@Id = 'DPST-15-0_F-1') and fn:position() = 3">value2</xsl:when>
                            <xsl:when test="(@Id = 'DPST-15-0_F-1') and fn:position() = 4">value3</xsl:when>
                            <xsl:when test="(@Id = 'DPST-15-0_F-1') and fn:position() = 5">value4</xsl:when>
                            <xsl:when test="(@Id = 'DPST-15-0_F-1') and fn:position() = 6">value5</xsl:when>
                            <xsl:when test="(@Id = 'DPST-15-0_F-1') and fn:position() = 7">value6</xsl:when>
                            <xsl:when test="@Id = 'DPST-1-2_F-1'">valueTrue</xsl:when>
                            <xsl:when test="@Name and not(@Id = 'DPST-244-600_F-2' or @Id = 'DPST-244-600_F-3')">
                                <xsl:call-template name="getFieldName">
                                    <xsl:with-param name="fieldName" select="@Name"/>
                                </xsl:call-template>
                            </xsl:when>
                            <xsl:when test="@Set">
                                <xsl:call-template name="getFieldName">
                                    <xsl:with-param name="fieldName" select="@Set"/>
                                </xsl:call-template>
                            </xsl:when>
                            <xsl:when test="@Unit">
                                <xsl:call-template name="getFieldName">
                                    <xsl:with-param name="fieldName" select="@Unit"/>
                                </xsl:call-template>
                            </xsl:when>
                            <xsl:otherwise>Hurz</xsl:otherwise>
                        </xsl:choose>
                    </xsl:variable>
                    <xsl:variable name="dataType">
                        <xsl:call-template name="getMspecType">
                            <xsl:with-param name="field" select="."/>
                        </xsl:call-template>
                    </xsl:variable>                    
                    [<xsl:value-of select="$fieldType"/><xsl:text disable-output-escaping="yes"> </xsl:text><xsl:value-of select="$dataType"/><xsl:text> </xsl:text><xsl:value-of select="fn:concat(fn:lower-case(fn:substring($fieldName, 1, 1)), fn:substring($fieldName, 2))"/>]
                </xsl:for-each>
            </xsl:when>
            <xsl:when test="$datapointSubtype/knx:Format/knx:Bit">
            [reserved uint 7 '0x00']
            [simple   bit    value]
            </xsl:when>
            <xsl:when test="$datapointSubtype/knx:Format/knx:String">
                <xsl:variable name="encoding">
                    <xsl:choose>
                        <xsl:when test="$datapointSubtype/knx:Format/knx:String/@Encoding = 'us-ascii'">ASCII</xsl:when>
                        <xsl:when test="$datapointSubtype/knx:Format/knx:String/@Encoding = 'iso-8859-1'">ISO-8859-1</xsl:when>
                    </xsl:choose>
                </xsl:variable>
            [simple   string <xsl:value-of select="$datapointSubtype/knx:Format/knx:String/@Width"/> value encoding='"<xsl:value-of select="$encoding"/>"']
            </xsl:when>
            <xsl:when test="$datapointSubtype/knx:Format/knx:UnsignedInteger">
                <xsl:choose>
                    <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:UnsignedInteger/@Width) &lt; 6">
            [reserved uint <xsl:value-of select="8 - fn:number($datapointSubtype/knx:Format/knx:UnsignedInteger/@Width)"/> '0x00']
                    </xsl:when>
                </xsl:choose>
            [simple   uint <xsl:value-of select="$datapointSubtype/knx:Format/knx:UnsignedInteger/@Width"/> value]
            </xsl:when>
            <xsl:when test="$datapointSubtype/knx:Format/knx:SignedInteger">
                <xsl:choose>
                    <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:SignedInteger/@Width) &lt; 6">
            [reserved uint <xsl:value-of select="8 - fn:number($datapointSubtype/knx:Format/knx:SignedInteger/@Width)"/> '0x00']
                    </xsl:when>
                </xsl:choose>
            [simple   int <xsl:value-of select="$datapointSubtype/knx:Format/knx:SignedInteger/@Width"/> value]
            </xsl:when>
            <xsl:when test="$datapointSubtype/knx:Format/knx:Float">
                <xsl:choose>
                    <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Float/@Width) = 16">
            [simple   float 16 value]
                    </xsl:when>
                    <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Float/@Width) = 32">
            [simple   float 32 value]
                    </xsl:when>
                    <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Float/@Width) = 64">
            [simple   float 64 value]
                    </xsl:when>
                </xsl:choose>
            </xsl:when>
            <xsl:when test="$datapointSubtype/knx:Format/knx:Enumeration">
                <xsl:choose>
                    <xsl:when test="fn:number($datapointSubtype/knx:Format/knx:Enumeration/@Width) &lt; 6">
            [reserved uint <xsl:value-of select="8 - fn:number($datapointSubtype/knx:Format/knx:Enumeration/@Width)"/> '0x00']
                    </xsl:when>
                </xsl:choose>
            [simple   uint <xsl:value-of select="$datapointSubtype/knx:Format/knx:Enumeration/@Width"/> value]
            </xsl:when>
        </xsl:choose>
        ]
    </xsl:template>

    <xsl:template name="resolveTypeReferences">
        <xsl:param name="fields"/>
        <xsl:for-each select="$fields">
            <xsl:choose>
                <xsl:when test="name(.) = 'RefType'">
                    <xsl:variable name="curNode" select="."/>
                    <xsl:copy-of select="//*[@Id = $curNode/@RefId]"/>
                </xsl:when>
                <xsl:otherwise>
                    <xsl:copy-of select="."/>
                </xsl:otherwise>
            </xsl:choose>
        </xsl:for-each>
    </xsl:template>

    <xsl:template name="getMspecType">
        <xsl:param name="field"/>
        <xsl:choose>
            <xsl:when test="name($field) = 'Bit'">bit</xsl:when>
            <xsl:when test="name($field) = 'String'">
                <xsl:variable name="encoding">
                    <xsl:choose>
                        <xsl:when test="$field/@Encoding = 'us-ascii'">ASCII</xsl:when>
                        <xsl:when test="$field/@Encoding = 'iso-8859-1'">ISO-8859-1</xsl:when>
                    </xsl:choose>
                </xsl:variable>
                string <xsl:value-of select="$field/@Width"/></xsl:when>
            <xsl:when test="name($field) = 'UnsignedInteger'">uint <xsl:value-of select="$field/@Width"/></xsl:when>
            <xsl:when test="name($field) = 'SignedInteger'">int <xsl:value-of select="$field/@Width"/></xsl:when>
            <xsl:when test="name($field) = 'Float'">
                <xsl:choose>
                    <xsl:when test="$field/@Width = 16">float 16</xsl:when>
                    <xsl:when test="$field/@Width = 32">float 32</xsl:when>
                    <xsl:when test="$field/@Width = 64">float 64</xsl:when>
                    <xsl:otherwise>hurz</xsl:otherwise>
                </xsl:choose>
            </xsl:when>
            <xsl:when test="name($field) = 'Enumeration'">uint <xsl:value-of select="$field/@Width"/></xsl:when>
            <xsl:when test="name($field) = 'Reserved'">uint <xsl:value-of select="$field/@Width"/></xsl:when>
        </xsl:choose>
    </xsl:template>
    
    <xsl:template name="getFieldName">
        <xsl:param name="fieldName"/>
        <xsl:variable name="cleanedName" select="fn:replace(fn:replace($fieldName, '%', 'Percent'), '[^a-zA-Z0-9]', ' ')"/>
        <xsl:for-each select="fn:tokenize($cleanedName, ' ')">
            <xsl:choose>
                <xsl:when test="fn:string-length(.) = 1"><xsl:value-of select="fn:upper-case(.)"/></xsl:when>
                <xsl:otherwise>
                    <xsl:variable name="firstLetter" select="fn:upper-case(fn:substring(., 1, 1))"/>
                    <xsl:variable name="rest" select="fn:lower-case(fn:substring(., 2))"/>
                    <xsl:value-of select="fn:concat($firstLetter, $rest)"/>                   
                </xsl:otherwise>
            </xsl:choose>
        </xsl:for-each>
    </xsl:template>
    
    <xsl:template name="fieldsSize">
        <xsl:param name="fields"/>
        <xsl:choose>
            <xsl:when test="count($fields) &gt; 1">
                <xsl:variable name="restSize">
                    <xsl:call-template name="fieldsSize">
                        <xsl:with-param name="fields" select="$fields[position() &gt; 1]"/>
                    </xsl:call-template>
                </xsl:variable>
                <xsl:choose>
                    <xsl:when test="name($fields[1]) = 'Bit'"><xsl:value-of select="1 + $restSize"/></xsl:when>
                    <xsl:otherwise><xsl:value-of select="number($fields[1]/@Width) + $restSize"/></xsl:otherwise>
                </xsl:choose>               
            </xsl:when>
            <xsl:otherwise>
                <xsl:choose>
                    <xsl:when test="name($fields[1]) = 'Bit'"><xsl:value-of select="1"/></xsl:when>
                    <xsl:otherwise><xsl:value-of select="number($fields[1]/@Width)"/></xsl:otherwise>
                </xsl:choose>
            </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

</xsl:stylesheet>