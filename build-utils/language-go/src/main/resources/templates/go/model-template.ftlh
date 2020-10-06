<#--
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
<#-- Prevent freemarker from escaping stuff -->
<#outputformat "undefined">
<#-- Declare the name and type of variables passed in to the template -->
<#-- @ftlvariable name="languageName" type="java.lang.String" -->
<#-- @ftlvariable name="protocolName" type="java.lang.String" -->
<#-- @ftlvariable name="outputFlavor" type="java.lang.String" -->
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.go.GoLanguageTemplateHelper" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
<#-- Declare the name and type of variables declared locally inside the template -->
<#-- @ftlvariable name="arrayField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ArrayField" -->
<#-- @ftlvariable name="checksumField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ChecksumField" -->
<#-- @ftlvariable name="constField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ConstField" -->
<#-- @ftlvariable name="discriminatorField" type="org.apache.plc4x.plugins.codegenerator.types.fields.DiscriminatorField" -->
<#-- @ftlvariable name="enumField" type="org.apache.plc4x.plugins.codegenerator.types.fields.EnumField" -->
<#-- @ftlvariable name="implicitField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ImplicitField" -->
<#-- @ftlvariable name="manualArrayField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ManualArrayField" -->
<#-- @ftlvariable name="manualField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ManualField" -->
<#-- @ftlvariable name="optionalField" type="org.apache.plc4x.plugins.codegenerator.types.fields.OptionalField" -->
<#-- @ftlvariable name="paddingField" type="org.apache.plc4x.plugins.codegenerator.types.fields.PaddingField" -->
<#-- @ftlvariable name="reservedField" type="org.apache.plc4x.plugins.codegenerator.types.fields.ReservedField" -->
<#-- @ftlvariable name="simpleField" type="org.apache.plc4x.plugins.codegenerator.types.fields.SimpleField" -->
<#-- @ftlvariable name="switchField" type="org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField" -->
<#-- @ftlvariable name="virtualField" type="org.apache.plc4x.plugins.codegenerator.types.fields.VirtualField" -->
<#-- @ftlvariable name="simpleTypeReference" type="org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference" -->
<#-- @ftlvariable name="complexTypeReference" type="org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference" -->
${helper.fileName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.go
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package ${helper.packageName(outputFlavor)}

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"

type ${type.name} struct {
<#-- Prpoerty fields are fields that require a property in the pojo -->
<#if type.propertyFields?has_content>
<#list type.propertyFields as field>
    ${field.name} <#if field.loopType??>[]</#if>${helper.getLanguageTypeNameForField(field)}
</#list>
</#if>
<#--
    This is GO's version of inheritance
    (Makes all properties of the following type available in this tpye)
-->
<#if type.parentType??>    ${type.parentType.name}</#if>
}

func (m ${type.name}) LengthInBits() uint16 {
    var lengthInBits uint16 = <#if type.parentType??>m.${type.parentType.name}.LengthInBits()<#else>0</#if>
    <#list type.fields as field>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field>

    // Array field
    if len(m.${arrayField.name}) > 0 {
                <#if helper.isSimpleTypeReference(arrayField.type)>
                    <#assign simpleTypeReference = arrayField.type>
        lengthInBits += ${simpleTypeReference.sizeInBits} * uint16(len(m.${arrayField.name}))
                <#else>
        for _, element := range m.${arrayField.name} {
            lengthInBits += element.LengthInBits()
        }
                </#if>
    }
                <#break>
            <#case "checksum">
                <#assign checksumField = field>
                <#assign simpleTypeReference = checksumField.type>

    // Checksum Field (checksum)
    lengthInBits += ${simpleTypeReference.sizeInBits}
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    lengthInBits += ${simpleTypeReference.sizeInBits}
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name})
    lengthInBits += ${simpleTypeReference.sizeInBits}
                <#break>
            <#case "enum">
                <#assign enumField = field>

    // Enum Field (${enumField.name})
    lengthInBits += ${helper.getEnumBaseTypeReference(enumField.type).sizeInBits}
                <#break>
            <#case "implicit">
                <#assign implicitField = field>
                <#assign simpleTypeReference = implicitField.type>

    // Implicit Field (${implicitField.name})
    lengthInBits += ${simpleTypeReference.sizeInBits}
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field>

    // Manual Array Field (${manualArrayField.name})
    lengthInBits += ${helper.toParseExpression(manualArrayField, manualArrayField.lengthExpression, type.parserArguments)} * 8
                <#break>
            <#case "manual">
                <#assign manualField = field>

    // Manual Field (${manualField.name})
    lengthInBits += ${helper.toParseExpression(manualField, manualField.lengthExpression, type.parserArguments)} * 8
                <#break>
            <#case "optional">
                <#assign optionalField = field>

    // Optional Field (${optionalField.name})
    if m.${optionalField.name} != nil {
                <#if helper.isSimpleTypeReference(optionalField.type)>
                    <#assign simpleTypeReference = optionalField.type>
        lengthInBits += ${simpleTypeReference.sizeInBits}
                <#else>
        lengthInBits += m.${optionalField.name}.LengthInBits()
                </#if>
    }
                <#break>
            <#case "padding">
                <#assign paddingField = field>
                <#assign simpleTypeReference = paddingField.type>

    // Padding Field (padding)
            <#-- We're replacing the "lastItem" with 'false' here as the item itself can't know if it is the last -->
    _timesPadding:uint8 = uint8(${helper.toParseExpression(paddingField, paddingField.paddingCondition, type.parserArguments)?replace("lastItem", "false")})
    for _timesPadding > 0; _timesPadding-- {
        lengthInBits += ${simpleTypeReference.sizeInBits}
    }
                <#break>
            <#case "reserved">
                <#assign reservedField = field>
                <#assign simpleTypeReference = reservedField.type>

    // Reserved Field (reserved)
    lengthInBits += ${simpleTypeReference.sizeInBits}
                <#break>
            <#case "simple">
                <#assign simpleField = field>

    // Simple field (${simpleField.name})
                <#if helper.isSimpleTypeReference(simpleField.type)>
                    <#assign simpleTypeReference = simpleField.type>
    lengthInBits += ${simpleTypeReference.sizeInBits}
                <#else>
    lengthInBits += m.${simpleField.name}.LengthInBits()
                </#if>
                <#break>
            <#case "switch">
                <#assign switchField = field>

    // Length of sub-type elements will be added by sub-type...
                <#break>
            <#case "virtual">
                <#assign virtualField = field>

    // A virtual field doesn't have any in- or output.
                <#break>
        </#switch>
    </#list>

    return lengthInBits
}

func (m ${type.name}) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func (m ${type.name}) Parse(io spi.ReadBuffer) {
    // TODO: Implement ...
}

func (m ${type.name}) Serialize(io spi.WriteBuffer) {
    // TODO: Implement ...
}
</#outputformat>