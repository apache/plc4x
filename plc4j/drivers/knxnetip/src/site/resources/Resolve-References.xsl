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
                xmlns:knx="http://knx.org/xml/project/11"
                exclude-result-prefixes="xs"
                version="2.0">
    <xsl:output indent="yes"/>

    <xsl:template match="/">
        <Dynamic>
            <xsl:apply-templates select="//knx:Dynamic/knx:Channel"/>
        </Dynamic>
    </xsl:template>

    <xsl:template match="knx:Channel">
        <xsl:apply-templates select="knx:ParameterBlock | knx:choose"/>
    </xsl:template>

    <xsl:template match="knx:ParameterBlock">
        <xsl:apply-templates select="knx:ParameterBlock | knx:ParameterRefRef | knx:ComObjectRef | knx:choose"/>
    </xsl:template>

    <xsl:template match="knx:ParameterRefRef">
        <xsl:variable name="refId" select="@RefId"/>
        <xsl:variable name="refElement" select="//knx:ParameterRef[@Id=$refId]"/>
        <xsl:variable name="elementId" select="$refElement/@RefId"/>
        <xsl:variable name="element" select="//knx:Parameter[@Id=$elementId]"/>
        <Parameter>
            <!-- First start with all attributes of the ParameterRef (except Id and RefId of the ref-element -->
            <xsl:for-each select="$refElement/@*">
                <xsl:choose>
                    <xsl:when test="name() = 'Id'">
                        <xsl:attribute name="ParameterRefId"><xsl:value-of select="."/></xsl:attribute>
                    </xsl:when>
                    <xsl:when test="(name() != 'RefId') and (name() != 'ParameterType')">
                        <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                    </xsl:when>
                </xsl:choose>
            </xsl:for-each>
            <!-- Then add all attributes of the Parameter (except Id and RefId of the element, if the ref-element didn't contain them) -->
            <xsl:for-each select="$element/@*">
                <xsl:variable name="attributeName" select="name()"/>
                <xsl:choose>
                    <xsl:when test="name() = 'Id'">
                        <xsl:attribute name="ParameterId"><xsl:value-of select="."/></xsl:attribute>
                    </xsl:when>
                    <xsl:when test="not($refElement/@*[name() = $attributeName]) and (name() != 'ParameterType')">
                        <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                    </xsl:when>
                </xsl:choose>
            </xsl:for-each>
            <!-- This will generally just copy the memory and other stuff -->
            <xsl:apply-templates select="$element/*"/>
            <!-- Get the ParameterType from either the ParamerterRef or the Parameter -->
            <xsl:variable name="parameterType" select="(//knx:ParameterType[@Id = $refElement/@ParameterType] | //knx:ParameterType[@Id = $element/@ParameterType])[1]"/>
            <ParameterType>
                <!-- Copy all attributes, except the Id attributes) -->
                <xsl:for-each select="$parameterType/@*">
                    <xsl:if test="name() != 'Id'">
                        <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                    </xsl:if>
                </xsl:for-each>
                <!-- Get the size in bits from the next lower element -->
                <xsl:choose>
                    <xsl:when test="$parameterType/*/@SizeInBit">
                        <xsl:attribute name="SizeInBit" select="$parameterType/*/@SizeInBit"/>
                    </xsl:when>
                    <xsl:otherwise>
                        <xsl:attribute name="SizeInBit" select="0"/>
                    </xsl:otherwise>
                </xsl:choose>
                <!-- Get the parameter type -->
                <xsl:choose>
                    <xsl:when test="$parameterType/*/@Type">
                        <xsl:attribute name="Type" select="$parameterType/*/@Type"/>
                    </xsl:when>
                    <xsl:when test="$parameterType/knx:TypeTestriction">
                        <xsl:attribute name="Type" select="'unsignedInt'"/>
                    </xsl:when>
                    <xsl:otherwise>
                        <xsl:attribute name="Type" select="'none'"/>
                    </xsl:otherwise>
                </xsl:choose>
                <xsl:choose>
                    <!-- If this is an enum -->
                    <xsl:when test="$parameterType/knx:TypeRestriction">
                        <xsl:for-each select="$parameterType/knx:TypeRestriction/knx:Enumeration">
                            <Enumeration>
                                <xsl:copy-of select="@*[name() != 'Id']"></xsl:copy-of>
                            </Enumeration>
                        </xsl:for-each>
                    </xsl:when>
                    <!-- If this is a numeric value -->
                    <xsl:when test="$parameterType/knx:TypeNumber">
                        <xsl:copy-of select="$parameterType/knx:TypeNumber/@*[(name() != 'Type') and (name() != 'SizeInBit')]"/>
                    </xsl:when>
                </xsl:choose>
            </ParameterType>
        </Parameter>
    </xsl:template>

    <xsl:template match="knx:ComObjectRefRef">
        <xsl:variable name="refId" select="@RefId"/>
        <xsl:variable name="refElement" select="//knx:ComObjectRef[@Id=$refId]"/>
        <xsl:variable name="elementId" select="$refElement/@RefId"/>
        <xsl:variable name="element" select="//knx:ComObject[@Id=$elementId]"/>
        <ComObject>
            <!-- First start with all attributes (except Id and RefId of the ref-element -->
            <xsl:for-each select="$refElement/@*">
                <xsl:if test="(name() != 'Id') and (name() != 'RefId')">
                    <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                </xsl:if>
            </xsl:for-each>
            <!-- Then add all attributes (except Id and RefId of the element, if the ref-element didn't contain them) -->
            <xsl:for-each select="$element/@*">
                <xsl:variable name="attributeName" select="name()"/>
                <xsl:if test="not($refElement/@*[name() = $attributeName])">
                    <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                </xsl:if>
            </xsl:for-each>
            <xsl:apply-templates select="$element/*"/>
        </ComObject>
    </xsl:template>

    <xsl:template match="knx:choose">
        <xsl:choose>
            <xsl:when test="(count(knx:when) = 1) and (knx:when/@default = 'true')">
                <xsl:for-each select="knx:when">
                    <xsl:apply-templates select="knx:ParameterBlock | knx:ParameterRefRef | knx:ComObjectRefRef | knx:choose"/>
                </xsl:for-each>
            </xsl:when>
            <xsl:otherwise>
                <xsl:variable name="parameterRefId" select="@ParamRefId"/>
                <xsl:variable name="parameterName">
                    <xsl:for-each select="//knx:ParameterRef[@Id=$parameterRefId]">
                        <xsl:variable name="parmeterId" select="@RefId"/>
                        <xsl:for-each select="//knx:Parameter[@Id=$parmeterId]">
                            <xsl:variable name="parameterId">
                                <xsl:choose>
                                    <xsl:when test="contains(@Id, '_P-')">_P-<xsl:value-of select="substring-after(@Id, '_P-')"/></xsl:when>
                                    <xsl:when test="contains(@Id, '_UP-')">_UP-<xsl:value-of select="substring-after(@Id, '_UP-')"/></xsl:when>
                                </xsl:choose>
                            </xsl:variable>
                            <xsl:value-of select="@Name"/> <xsl:value-of select="$parameterId"/>
                        </xsl:for-each>
                    </xsl:for-each>
                </xsl:variable>
                <Switch on="{$parameterName}">
                    <xsl:for-each select="knx:when">
                        <Case value="{@test}">
                            <xsl:apply-templates select="knx:ParameterBlock | knx:ParameterRefRef | knx:ComObjectRefRef | knx:choose"/>
                        </Case>
                    </xsl:for-each>
                </Switch>
            </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <xsl:template match="knx:Memory">
        <xsl:variable name="memoryElement" select="."/>
        <xsl:variable name="codeSegmentId" select="@CodeSegment"/>
        <xsl:variable name="codeSegment" select="//*[@Id=$codeSegmentId]"/>
        <xsl:element name="{name($codeSegment)}">
            <!-- Copy all attributes of the Memory element -->
            <xsl:for-each select="$memoryElement/@*">
                <xsl:if test="name() != 'CodeSegment'">
                    <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                </xsl:if>
            </xsl:for-each>
            <!-- Now copy all attributes of the referenced memory block that weren't overwritten in the Memory element -->
            <xsl:for-each select="$codeSegment/@*">
                <xsl:variable name="attributeName" select="name()"/>
                <xsl:if test="not($memoryElement/@*[name() = $attributeName]) and (name() != 'Id')">
                    <xsl:attribute name="{name()}"><xsl:value-of select="."/></xsl:attribute>
                </xsl:if>
            </xsl:for-each>
        </xsl:element>
    </xsl:template>

</xsl:stylesheet>