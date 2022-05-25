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
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xs="http://www.w3.org/2001/XMLSchema"
    xmlns:knx20="http://knx.org/xml/project/20"
    exclude-result-prefixes="xs"
    version="2.0">
    
    <xsl:output method="text"/>
    
    <xsl:variable name="knxMasterData" select="document('https://update.knx.org/data/XML/project-20/knx_master.xml')/*"/>
    
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
package org.apache.plc4x.java.knxnetip;
        
import org.apache.plc4x.java.knxnetip.readwrite.KnxDatapointType;
        
import java.util.Map;
import java.util.HashMap;

public class ComObjectTypeResolver {

    protected Integer getIntParameter() {
        return (int) (Math.random() * 100);
    }
        
    public Map&lt;Integer, KnxDatapointType> getComObjectTypes(byte[] configurationMemory) {
        Map&lt;Integer, KnxDatapointType> comObjectTypes = new HashMap&lt;>();

        // Declare any parameters, that might be used. 
        <xsl:for-each-group select="//Parameter" group-by="@Name">
            <xsl:call-template name="ParameterDelaration">
                <xsl:with-param name="param" select="current-group()[1]"/>
            </xsl:call-template>
        </xsl:for-each-group>       
        // Output the logic.
        <xsl:apply-templates select="Dynamic"/>       
        return comObjectTypes;
    }
}            
    </xsl:template>
    
    <xsl:template match="Dynamic">
        <xsl:apply-templates select="*"/>
    </xsl:template>
    
    <xsl:template name="ParameterDelaration">
        <xsl:param name="param"/>
        <xsl:variable name="parameterType">
            <xsl:call-template name="getVariableType">
                <xsl:with-param name="parameter" select="$param"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="parmeterVariableName">
            <xsl:call-template name="toVariableName">
                <xsl:with-param name="text" select="$param/@Name"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:value-of select="$parameterType"/><xsl:text> </xsl:text><xsl:value-of select="$parmeterVariableName"/> = null; 
    </xsl:template>
    
    <xsl:template match="Parameter">
        <xsl:variable name="parmeterVariableName">
            <xsl:call-template name="toVariableName">
                <xsl:with-param name="text" select="@Name"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="parameterGetValueCall">
            <xsl:call-template name="getGetVariableMethodCall">
                <xsl:with-param name="parameter" select="."/>
            </xsl:call-template>
        </xsl:variable><xsl:value-of select="$parmeterVariableName"/> = <xsl:value-of select="$parameterGetValueCall"/>;
    </xsl:template>
        
    <xsl:template match="ComObject">
        <xsl:variable name="comObjectType">
            <xsl:call-template name="getComObjectType">
                <xsl:with-param name="comObject" select="."/>
            </xsl:call-template>
        </xsl:variable>comObjectTypes.put(<xsl:value-of select="@Number"/>, KnxDatapointType.<xsl:value-of select="$comObjectType"/>);
    </xsl:template>

    <xsl:template match="switch">
        <xsl:variable name="parmeterVariableName">
            <xsl:call-template name="toVariableName">
                <xsl:with-param name="text" select="@on"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:variable name="hasDefault" select="case[@default = true()]"/><xsl:if test="$hasDefault">
            {
                boolean switchGuard<xsl:value-of select="$parmeterVariableName"/> = false;
        </xsl:if><xsl:for-each select="case">
            <xsl:choose>
                <!-- TODO: implement the default=true stuff<xsl:when test=""></xsl:when>-->
                <xsl:when test="@default = true()">
                    if(switchGuard<xsl:value-of select="$parmeterVariableName"/> == false) {</xsl:when>
                <xsl:when test="starts-with(@value, '!=')">
                    if((<xsl:value-of select="$parmeterVariableName"/> != null) &amp;&amp; (<xsl:value-of select="$parmeterVariableName"/> != <xsl:value-of select="@value"/>)) {
                    <xsl:if test="$hasDefault">
                        switchGuard<xsl:value-of select="$parmeterVariableName"/> = true;
                    </xsl:if></xsl:when>
                <xsl:when test="starts-with(@value, '&lt;=')">
                    if((<xsl:value-of select="$parmeterVariableName"/> != null) &amp;&amp; (<xsl:value-of select="$parmeterVariableName"/> &lt;= <xsl:value-of select="@value"/>)) {
                    <xsl:if test="$hasDefault">
                        switchGuard<xsl:value-of select="$parmeterVariableName"/> = true;                    
                    </xsl:if></xsl:when>
                <xsl:when test="starts-with(@value, '>=')">
                    if((<xsl:value-of select="$parmeterVariableName"/> != null) &amp;&amp; ((<xsl:value-of select="$parmeterVariableName"/> >= <xsl:value-of select="@value"/>)) {
                    <xsl:if test="$hasDefault">
                        switchGuard<xsl:value-of select="$parmeterVariableName"/> = true;                   
                    </xsl:if></xsl:when>
                <xsl:when test="starts-with(@value, '&lt;')">
                    if((<xsl:value-of select="$parmeterVariableName"/> != null) &amp;&amp; (<xsl:value-of select="$parmeterVariableName"/> &lt; <xsl:value-of select="@value"/>)) {
                    <xsl:if test="$hasDefault">
                        switchGuard<xsl:value-of select="$parmeterVariableName"/> = true;                    
                    </xsl:if></xsl:when>
                <xsl:when test="starts-with(@value, '>')">
                    if((<xsl:value-of select="$parmeterVariableName"/> != null) &amp;&amp; (<xsl:value-of select="$parmeterVariableName"/> > <xsl:value-of select="@value"/>)) {
                    <xsl:if test="$hasDefault">
                        switchGuard<xsl:value-of select="$parmeterVariableName"/> = true;                    
                    </xsl:if></xsl:when>
                <xsl:otherwise>
                    if((<xsl:value-of select="$parmeterVariableName"/> != null) &amp;&amp; (<xsl:for-each select="tokenize(@value, ' ')">(<xsl:value-of select="$parmeterVariableName"/> == <xsl:value-of select="."/>)<xsl:if test="position() != last()"> or </xsl:if></xsl:for-each>)) {
                    <xsl:if test="$hasDefault">
                        switchGuard<xsl:value-of select="$parmeterVariableName"/> = true;                    
                    </xsl:if></xsl:otherwise>
            </xsl:choose>
            <xsl:apply-templates select="*"/>
            }</xsl:for-each><xsl:if test="$hasDefault">
            }
        </xsl:if>
        
    </xsl:template>



    <xsl:template name="getVariableType">
        <xsl:param name="parameter"/>
        <xsl:choose>
            <xsl:when test="($parameter/ParameterType/@SizeInBit &lt; 31) and ($parameter/ParameterType/@Type = 'unsignedInt')">Integer</xsl:when>
            <xsl:when test="($parameter/ParameterType/@SizeInBit &lt; 32) and ($parameter/ParameterType/@Type = 'signedInt')">Integer</xsl:when>
            <xsl:when test="($parameter/ParameterType/@SizeInBit &lt; 31) and ($parameter/ParameterType/@Type = 'none')">Integer</xsl:when>
            <xsl:otherwise>hurz</xsl:otherwise>
        </xsl:choose>
    </xsl:template>
    
    <xsl:template name="getGetVariableMethodCall">
        <xsl:param name="parameter"/>
        <xsl:choose>
            <xsl:when test="($parameter/ParameterType/@SizeInBit &lt; 31) and ($parameter/ParameterType/@Type = 'unsignedInt')">getIntParameter()</xsl:when>
            <xsl:when test="($parameter/ParameterType/@SizeInBit &lt; 32) and ($parameter/ParameterType/@Type = 'signedInt')">getIntParameter()</xsl:when>
            <xsl:when test="($parameter/ParameterType/@SizeInBit &lt; 31) and ($parameter/ParameterType/@Type = 'none')">getIntParameter()</xsl:when>
            <xsl:otherwise>hurz</xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template name="toVariableName">
        <xsl:param name="text"/>             
        <xsl:variable name="filteredText" select="replace(replace(replace($text, '[,\-&quot;,\(, \)]', ''), '[:,\.,+]', '_'), '/', 'Or')"/>
        <xsl:variable name="total">
            <xsl:call-template name="toCamelCaseRest">
                <xsl:with-param name="text" select="$filteredText"/>
            </xsl:call-template>
        </xsl:variable>
        <xsl:value-of select="concat('var', $total)"/>
    </xsl:template>

    <xsl:template name="toCamelCaseRest">
        <xsl:param name="text"/>
        <xsl:variable name="uppercase">ABCDEFGHIJKLMNOPQRSTUVWXYZ</xsl:variable>
        <xsl:variable name="lowercase">abcdefghijklmnopqrstuvwxyz</xsl:variable>
        <xsl:variable name="word" select="substring-before(concat($text, ' '), ' ')" />
        <xsl:variable name="firstLetter" select="translate(substring($word, 1, 1), $lowercase, $uppercase)"/>
        <xsl:variable name="restWord" select="translate(substring($word, 2), $uppercase, $lowercase)"/>
        <xsl:variable name="restExpression">
            <xsl:if test="contains($text, ' ')">
                <xsl:call-template name="toCamelCaseRest">
                    <xsl:with-param name="text" select="substring-after($text, ' ')"/>
                </xsl:call-template>
            </xsl:if>       
        </xsl:variable>
        <xsl:value-of select="concat($firstLetter, $restWord, $restExpression)"/>
    </xsl:template>

    <xsl:template name="getComObjectType">
        <xsl:param name="comObject"/>
        <xsl:choose>
            <!-- If an explicit datapoint type is defined, use that -->
            <xsl:when test="$comObject/@DatapointType and (string-length($comObject/@DatapointType) > 0) and $knxMasterData//knx20:DatapointSubtype[@Id=substring-before(concat($comObject/@DatapointType, ' '), ' ')]">
                <!-- Lookup the value in @DatapointType in knx_master_data.xml and use the @Name -->
                <xsl:value-of select="$knxMasterData//knx20:DatapointSubtype[@Id=substring-before(concat($comObject/@DatapointType, ' '), ' ')]/@Name"/>
            </xsl:when>
            <!-- In this case we're referencing a data point type and not a subtype -->
            <xsl:when test="$comObject/@DatapointType and (string-length($comObject/@DatapointType) > 0) and $knxMasterData//knx20:DatapointType[@Id=substring-before(concat($comObject/@DatapointType, ' '), ' ')]">
                <xsl:value-of select="$knxMasterData//knx20:DatapointType[@Id=substring-before(concat($comObject/@DatapointType, ' '), ' ')]/knx20:DatapointSubtypes/knx20:DatapointSubtype[1]/@Name"/>
            </xsl:when>
            <!-- Otherwise derive one from the object-size -->
            <xsl:when test="$comObject/@ObjectSize = '1 Bit'">DPT_Switch</xsl:when>
            <xsl:when test="$comObject/@ObjectSize = '4 Bit'">DPT_Control_Dimming</xsl:when>
            <xsl:when test="$comObject/@ObjectSize = '1 Byte'">DPT_Scaling</xsl:when>
            <xsl:when test="$comObject/@ObjectSize = '2 Bytes'">DPT_Value_2_Ucount</xsl:when>
            <xsl:when test="$comObject/@ObjectSize = '4 Bytes'">DPT_Value_4_Ucount</xsl:when>
            <xsl:otherwise> // TODO: Add case for ObjectSize='<xsl:value-of select="$comObject/@ObjectSize"/>'</xsl:otherwise>
        </xsl:choose>
    </xsl:template>

</xsl:stylesheet>