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
                xmlns:ua="http://opcfoundation.org/UA/"
>

    <xsl:output
        method="text"
        indent="no"
        encoding="utf-8"
    />

    <xsl:import href="opc-common.xsl"/>

    <xsl:variable name="originaldoc" select="/"/>

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



    <xsl:template match="/">

    </xsl:template>
</xsl:stylesheet>
