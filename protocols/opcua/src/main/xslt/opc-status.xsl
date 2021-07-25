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

    <xsl:param name="statusCodes"></xsl:param>

    <xsl:param name="statusCodeFile" select="unparsed-text($statusCodes)"/>

    <xsl:template match="/">
        <xsl:call-template name="statusCodeParsing"/>
    </xsl:template>

    <xsl:template name="statusCodeParsing" >
        <xsl:variable name="tokenizedLine" select="tokenize($statusCodeFile, '\r\n|\r|\n')" />
[enum uint 32 'OpcuaStatusCode'<xsl:text>
    </xsl:text>
        <xsl:for-each select="$tokenizedLine">
            <xsl:variable select="tokenize(., ',')" name="values" />    ['<xsl:value-of select="$values[2]"/>L'  <xsl:value-of select="$values[1]"/>]<xsl:text>
    </xsl:text>
        </xsl:for-each>
]
    </xsl:template>
</xsl:stylesheet>