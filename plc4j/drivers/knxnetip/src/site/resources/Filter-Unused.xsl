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
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xs="http://www.w3.org/2001/XMLSchema"
    exclude-result-prefixes="xs"
    version="2.0">
    
    <xsl:output indent="yes"/>
    
    <xsl:template match="Dynamic">
        <Dynamic>
            <xsl:apply-templates select="*"/>
        </Dynamic>
    </xsl:template>
    
    <xsl:template match="Parameter">
        <xsl:variable name="parameterName" select="@Name"/>
        <xsl:if test="//switch[@on = $parameterName]">
            <xsl:copy-of select="."/>
        </xsl:if>
    </xsl:template>
    
    <xsl:template match="ComObject">
        <xsl:copy-of select="."/>
    </xsl:template>
    
    <xsl:template match="switch">
        <xsl:variable name="switchContent">
            <xsl:apply-templates select="*"/>
        </xsl:variable>
        <xsl:if test="$switchContent//Parameter | $switchContent//ComObject">
            <switch>
                <xsl:copy-of select="@*"/>
                <xsl:copy-of select="$switchContent"/>
            </switch>
        </xsl:if>
    </xsl:template>
    
    <xsl:template match="case">
        <xsl:variable name="caseContent">
            <xsl:apply-templates select="*"/>
        </xsl:variable>
        <xsl:if test="$caseContent//Parameter | $caseContent//ComObject">
            <case>
                <xsl:copy-of select="@*"/>
                <xsl:copy-of select="$caseContent"/>
            </case>
        </xsl:if>
    </xsl:template>
    
</xsl:stylesheet>