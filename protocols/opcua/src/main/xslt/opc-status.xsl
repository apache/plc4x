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
>

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />
    <xsl:import href="opc-common.xsl"/>

    <xsl:variable name="originaldoc" select="/"/>

    <xsl:param name="osType"/>

    <xsl:param name="statusCodes"/>

    <!-- Intermediate, to reformat the url on windows systems -->
    <xsl:param name="statusCodesUrl">
        <xsl:choose>
            <xsl:when test="$osType = 'win'">file:///<xsl:value-of select="replace($statusCodes, '\\', '/')"/></xsl:when>
            <xsl:otherwise><xsl:value-of select="$statusCodes"/></xsl:otherwise>
        </xsl:choose>
    </xsl:param>

    <xsl:param name="statusCodeFile" select="unparsed-text($statusCodesUrl)"/>

    <xsl:template match="/">
        <xsl:call-template name="statusCodeParsing"/>
    </xsl:template>

    <xsl:template name="statusCodeParsing" >
        <xsl:variable name="tokenizedLine" select="tokenize($statusCodeFile, '\r\n|\r|\n')" />
[enum uint 32 OpcuaStatusCode<xsl:text>
    </xsl:text>
        <xsl:for-each select="$tokenizedLine">
            <xsl:variable select="tokenize(., ',')" name="values" />    ['<xsl:value-of select="$values[2]"/>'  <xsl:value-of select="$values[1]"/>]<xsl:text>
    </xsl:text>
        </xsl:for-each>
]
    </xsl:template>
</xsl:stylesheet>