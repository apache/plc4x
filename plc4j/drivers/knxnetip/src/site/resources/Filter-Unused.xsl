<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xs="http://www.w3.org/2001/XMLSchema"
    exclude-result-prefixes="xs"
    version="2.0">
    
    <xsl:output indent="yes"/>
    
    <xsl:key name="chooseParameterRefId" match="//choose" use="@ParamRefId"/>
    
    <xsl:template match="Dynamic">
        <Dynamic>
            <xsl:apply-templates select="*"/>
        </Dynamic>
    </xsl:template>
    
    <!-- Only copy a Parameter, if it is used in any switch expression -->
    <xsl:template match="Parameter">
        <xsl:variable name="paramRefIdTmp" select="@ParameterRefId"/>
        <xsl:if test="key('chooseParameterRefId', $paramRefIdTmp)">
            <xsl:copy-of select="."/>
        </xsl:if>
    </xsl:template>
    
    <!-- Copy all ComObjects -->
    <xsl:template match="ComObject">
        <xsl:copy-of select="."/>
    </xsl:template>
    
    
    <!-- Only copy a switch, if there is at least one ComObject reference or one Parameter inside -->
    <!--xsl:template match="choose">
        <xsl:variable name="chooseContent">
            <xsl:apply-templates select="*"/>
        </xsl:variable>
        <xsl:if test="$chooseContent//Parameter | $chooseContent//ComObject">
            <switch>
                <xsl:copy-of select="@*"/>
                <xsl:copy-of select="$chooseContent"/>
            </switch>
        </xsl:if>
    </xsl:template>
    
    <xsl:template match="when">
        <xsl:variable name="whenContent">
            <xsl:apply-templates select="*"/>
        </xsl:variable>
        <xsl:if test="$whenContent//Parameter | $whenContent//ComObject">
            <case>
                <xsl:copy-of select="@*"/>
                <xsl:copy-of select="$whenContent"/>
            </case>
        </xsl:if>
    </xsl:template-->
    
</xsl:stylesheet>