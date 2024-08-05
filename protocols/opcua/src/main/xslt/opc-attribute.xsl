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

    <xsl:param name="fileName"/>

    <!-- Intermediate, to reformat the url on windows systems -->
    <xsl:param name="csvFileUrl">
        <xsl:choose>
            <xsl:when test="$osType = 'win'">file:///<xsl:value-of select="replace($fileName, '\\', '/')"/></xsl:when>
            <xsl:otherwise><xsl:value-of select="$fileName"/></xsl:otherwise>
        </xsl:choose>
    </xsl:param>

    <xsl:param name="csvFile" select="unparsed-text($csvFileUrl)"/>

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
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Code generated by code-generation. DO NOT EDIT.
        <xsl:call-template name="csvParser"/>
    </xsl:template>

    <xsl:template name="csvParser">
        <xsl:variable name="tokenizedLine" select="tokenize($csvFile, '\r\n|\r|\n')" />
[enum uint 32 AttributeId
        <xsl:for-each select="$tokenizedLine">
            <xsl:variable select="tokenize(., ',')" name="values" />
            <xsl:if test="$values[1] != ''"><xsl:text>
        </xsl:text><xsl:value-of select="concat('[''', $values[2], ''' ', $values[1], ']')" /></xsl:if>
        </xsl:for-each>
]
    </xsl:template>
</xsl:stylesheet>