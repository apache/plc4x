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
${helper.fileName(protocolName, languageName, outputFlavor)?replace(".", "/")}/model/${type.name}.go
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
package model

import (
<#if helper.getRequiredImports()?has_content>
    <#list helper.getRequiredImports() as import>
    ${import}
    </#list>
</#if>
)
<#if type.constFields?has_content>

// Constant values.
    <#list type.constFields as field>
const ${type.name}_${field.name?upper_case} ${helper.getLanguageTypeNameForField(field)} = ${field.referenceValue}
    </#list>
</#if>

// The data-structure of this message
type ${type.name} struct {
<#-- Prpoerty fields are fields that require a property in the pojo -->
<#if type.propertyFields?has_content>
<#list type.propertyFields as field>
    ${field.name?cap_first} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if><#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForField(field)}
</#list>
</#if>
<#--
    This is GO's version of inheritance
    (Makes all properties of the following type available in this tpye)
-->
<#if type.parentType??>    ${type.parentType.name}</#if>
}

// The corresponding interface
type I${type.name} interface {
<#if helper.isDiscriminatedChildTypeDefinition()>
    I${type.parentType.name}
<#else>
    spi.Message
</#if>
<#if helper.isDiscriminatedParentTypeDefinition()>
    <#list helper.getDiscriminatorTypes() as discriminatorName, discriminatorType>
    <#-- If the discriminator name matches that of another field, suppress the methods generation -->
        <#if !helper.isNonDiscriminatorField(discriminatorName)>
    ${discriminatorName?cap_first}() ${helper.getLanguageTypeNameForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])}
        </#if>
    </#list>
</#if>
    Serialize(io utils.WriteBuffer<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${parserArgument.name} ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) error
}

<#if helper.isDiscriminatedChildTypeDefinition()>
    <#assign discriminatedChildType = type>
    <#-- @ftlvariable name="discriminatedChildType" type="org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition" -->
// Accessors for discriminator values.
    <#list helper.getDiscriminatorValues(discriminatedChildType) as discriminatorName, discriminatorValue>
        <#-- If the discriminator name matches that of another field, suppress the methods generation -->
        <#if !helper.isNonDiscriminatorField(discriminatorName)>
func (m ${type.name}) ${discriminatorName?cap_first}() ${helper.getLanguageTypeNameForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])} {
    return <#if discriminatorValue??>${discriminatorValue}<#else>${helper.getNullValueForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])}</#if>
}

        </#if>
    </#list>
</#if>
<#if helper.isDiscriminatedParentTypeDefinition()>
type ${type.name}Initializer interface {
    initialize(<#list type.getAllPropertyFields() as field>${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if><#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) spi.Message
}

    <#list helper.getDiscriminatorTypes() as discriminatorName, discriminatorType>
    <#-- If the discriminator name matches that of another field, suppress the methods generation -->
        <#if !helper.isNonDiscriminatorField(discriminatorName)>
func ${type.name}${discriminatorName?cap_first}(m I${type.name}) ${helper.getLanguageTypeNameForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])} {
    return m.${discriminatorName?cap_first}()
}

        </#if>
    </#list>
<#else>
<#if type.parentType??>func (m ${type.name}) initialize(<#list type.parentType.getAllPropertyFields() as field>${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if><#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) spi.Message {
    <#list type.parentType.getAllPropertyFields() as field>
    m.${field.name?cap_first} = ${field.name}
    </#list>
    return m
}
</#if>

func New${type.name}(<#list type.getPropertyFields() as field>${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if><#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) <#if type.parentType??>${type.parentType.name}Initializer<#else>spi.Message</#if> {
    return &${type.name}{<#list type.getPropertyFields() as field>${field.name?cap_first}: ${field.name}<#sep>, </#sep></#list>}
}
</#if>

func CastI${type.name}(structType interface{}) I${type.name} {
    castFunc := func(typ interface{}) I${type.name} {
        if i${type.name}, ok := typ.(I${type.name}); ok {
            return i${type.name}
        }
        return nil
    }
    return castFunc(structType)
}

func Cast${type.name}(structType interface{}) ${type.name} {
    castFunc := func(typ interface{}) ${type.name} {
        if s${type.name}, ok := typ.(${type.name}); ok {
            return s${type.name}
        }
        if s${type.name}, ok := typ.(*${type.name}); ok {
            return *s${type.name}
        }
        return ${type.name}{}
    }
    return castFunc(structType)
}

func (m ${type.name}) LengthInBits() uint16 {
    var lengthInBits uint16 = <#if type.parentType??>m.${type.parentType.name}.LengthInBits()<#else>0</#if>
    <#list type.fields as field>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field>

    // Array field
    if len(m.${arrayField.name?cap_first}) > 0 {
                <#if helper.isSimpleTypeReference(arrayField.type)>
                    <#assign simpleTypeReference = arrayField.type>
        lengthInBits += ${simpleTypeReference.sizeInBits} * uint16(len(m.${arrayField.name?cap_first}))
                <#else>
        for _, element := range m.${arrayField.name?cap_first} {
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
    if m.${optionalField.name?cap_first} != nil {
                <#if helper.isSimpleTypeReference(optionalField.type)>
                    <#assign simpleTypeReference = optionalField.type>
        lengthInBits += ${simpleTypeReference.sizeInBits}
                <#else>
        lengthInBits += (*m.${optionalField.name?cap_first}).LengthInBits()
                </#if>
    }
                <#break>
            <#case "padding">
                <#assign paddingField = field>
                <#assign simpleTypeReference = paddingField.type>

    // Padding Field (padding)
            <#-- We're replacing the "lastItem" with 'false' here as the item itself can't know if it is the last -->
    _timesPadding := uint8(${helper.toSerializationExpression(paddingField, paddingField.paddingCondition, type.parserArguments)?replace("lastItem", "false")})
    for ;_timesPadding > 0; _timesPadding-- {
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
    lengthInBits += m.${simpleField.name?cap_first}.LengthInBits()
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

func ${type.name}Parse(io *utils.ReadBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} <#if helper.isComplexTypeReference(parserArgument.type)>I</#if>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) (<#if type.parentType??>${type.parentType.name}Initializer<#else>spi.Message</#if>, error) {
    <#if helper.requiresStartPosAndCurPos()>
    var startPos = io.GetPos()
    var curPos uint16
    </#if>
    <#list type.fields as field>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field>

    // Array field (${arrayField.name})
            <#-- Only update curPos if the length expression uses it -->
                <#if arrayField.loopExpression.contains("curPos")>
    curPos = io.GetPos() - startPos
                </#if>
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                <#if helper.isCountArrayField(field)>
    // Count array
    ${arrayField.name} := make([]<#if helper.isComplexTypeReference(arrayField.type)>I</#if>${helper.getLanguageTypeNameForField(field)}, ${helper.toIntegerParseExpression(16, arrayField.loopExpression, type.parserArguments)})
    for curItem := uint16(0); curItem < uint16(${helper.toIntegerParseExpression(16, arrayField.loopExpression, type.parserArguments)}); curItem++ {
<#if (!helper.isSimpleTypeReference(arrayField.type)) && helper.requiresVariable(arrayField, "lastItem")>            lastItem := curItem == uint16(${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)} - 1)</#if>
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
        _item, _err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        if _err != nil {
            return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
        }
        ${arrayField.name}[curItem] = _item
                    <#else>
        _message, _err := ${arrayField.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(arrayField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {
            return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
        }
        var _item I${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(I${arrayField.type.name})
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ${arrayField.type.name}")
        }
        ${arrayField.name}[curItem] = _item
                    </#if>
    }
                <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                <#elseif helper.isLengthArrayField(field)>
                <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
    // Length array
    ${arrayField.name} := make([]<#if helper.isComplexTypeReference(arrayField.type)>I</#if>${helper.getLanguageTypeNameForField(field)}, 0)
    _${arrayField.name}Length := ${helper.toIntegerParseExpression(16, arrayField.loopExpression, type.parserArguments)}
    _${arrayField.name}EndPos := io.GetPos() + uint16(_${arrayField.name}Length)
    for ;io.GetPos() < _${arrayField.name}EndPos; {
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
        _item, _err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        if _err != nil {
            return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
        }
        ${arrayField.name} = append(${arrayField.name}, _item)
                    <#else>
        _message, _err := ${arrayField.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(arrayField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {
            return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
        }
        var _item I${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(I${arrayField.type.name})
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ${arrayField.type.name}")
        }
        ${arrayField.name} = append(${arrayField.name}, _item)
                    </#if>
                    <#-- After parsing, update the current position, but only if it's needed -->
                    <#if arrayField.loopExpression.contains("curPos")>
        curPos = io.GetPos() - startPos
                    </#if>
    }
                    <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                <#elseif helper.isTerminatedArrayField(field)>
    // Terminated array
    ${arrayField.name} := make([]<#if helper.isComplexTypeReference(arrayField.type)>I</#if>${helper.getLanguageTypeNameForField(field)}, 0)
    for ;!bool(${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}); {
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
        _item, _err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        if _err != nil {
            return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
        }
        ${arrayField.name} = append(${arrayField.name}, _item)
                    <#else>
        _message, _err := ${arrayField.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(arrayField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {

        }
        var _item I${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(I${arrayField.type.name})
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ${arrayField.type.name}")
        }
        ${arrayField.name} = append(${arrayField.name}, _item)
                    </#if>

                    <#-- After parsing, update the current position, but only if it's needed -->
                    <#if arrayField.loopExpression.contains("curPos")>
        curPos = io.GetPos() - startPos
                    </#if>
    }
                </#if>
                <#break>
            <#case "checksum">
                <#assign checksumField = field>
                <#assign simpleTypeReference = checksumField.type>

    // Checksum Field (checksum)
    {
        checksum = ${helper.getNullValueForTypeReference(checksumField.type)}
        // Create an array of all the bytes read in this message element so far.
        checksumRawData := io.getBytes(startPos, io.GetPos())
        checksumRef, _checksumRefErr := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        if _checksumRefErr != nil {
            return nil, errors.New("Error parsing 'checksum' field " + _checksumRefErr.Error())
        }
        checksum, _checksumErr = ${helper.toParseExpression(checksumField, checksumField.checksumExpression, type.parserArguments)}
        if _checksumErr != nil {
            return nil, errors.New("Error parsing 'checksum' field " + _checksumErr.Error())
        }
        if checksum != checksumRef {
            return nil, errors.New("Checksum verification failed. Expected " + (checksumRef & 0xFFFF) + " but got " + (checksum & 0xFFFF)")
        }
    }
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    ${constField.name}, _${constField.name}Err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
    if _${constField.name}Err != nil {
        return nil, errors.New("Error parsing '${constField.name}' field " + _${constField.name}Err.Error())
    }
    if ${constField.name} != ${type.name}_${constField.name?upper_case} {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(${type.name}_${constField.name?upper_case})) + " but got " + strconv.Itoa(int(${constField.name})))
    }
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
    ${helper.getVariableName(field)}, _${discriminatorField.name}Err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
    if _${discriminatorField.name}Err != nil {
        return nil, errors.New("Error parsing '${discriminatorField.name}' field " + _${discriminatorField.name}Err.Error())
    }
                <#break>
            <#case "enum">
                <#assign enumField = field>

    // Enum field (${enumField.name})
    ${enumField.name}, _${enumField.name}Err := ${helper.getLanguageTypeNameForField(field)}Parse(io)
    if _${enumField.name}Err != nil {
        return nil, errors.New("Error parsing '${enumField.name}' field " + _${enumField.name}Err.Error())
    }
                <#break>
            <#case "implicit">
                <#assign implicitField = field>
                <#assign simpleTypeReference = implicitField.type>

    // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    ${helper.getVariableName(field)}, _${implicitField.name}Err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
    if _${implicitField.name}Err != nil {
        return nil, errors.New("Error parsing '${implicitField.name}' field " + _${implicitField.name}Err.Error())
    }
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field>

    // Manual Array Field (${manualArrayField.name})
            <#-- Only update curPos if the length expression uses it -->
                <#if manualArrayField.loopExpression.contains("curPos")>
    curPos = io.GetPos() - startPos
                </#if>
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                <#if helper.isCountArrayField(field)>
    // Count array
    _${manualArrayField.name}Count := ${helper.toParseExpression(manualArrayField, manualArrayField.loopExpression, type.parserArguments)}
    ${helper.getLanguageTypeNameForField(field)}[] ${manualArrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[_${manualArrayField.name}Count]
    for i := 0; i < _${manualArrayField.name}Count; i++ {
        ${manualArrayField.name}[i] = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)})
    }
                <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                <#else>
                <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                    <#if helper.isLengthArrayField(field)>
    // Length array
    _${manualArrayField.name}Length := ${helper.toIntegerParseExpression(16, manualArrayField.loopExpression, type.parserArguments)}
    List<${helper.getLanguageTypeNameForField(manualArrayField)}> _${manualArrayField.name}List = new LinkedList<>()
    ${manualArrayField.name}EndPos := io.GetPos() + _${manualArrayField.name}Length
    for ;io.GetPos() < ${manualArrayField.name}EndPos; {
        _${manualArrayField.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)}))
                    <#-- After parsing, update the current position, but only if it's needed -->
                        <#if manualArrayField.loopExpression.contains("curPos")>
        curPos = io.GetPos() - startPos
                        </#if>
    }
                    <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                    <#elseif helper.isTerminatedArrayField(field)>
    // Terminated array
    List<${helper.getLanguageTypeNameForField(manualArrayField)}> _${manualArrayField.name}List = new LinkedList<>()
    for ;!((boolean) (${helper.toParseExpression(manualArrayField, manualArrayField.loopExpression, type.parserArguments)})); {
        _${manualArrayField.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)}))

                    <#-- After parsing, update the current position, but only if it's needed -->
                        <#if manualArrayField.loopExpression.contains("curPos")>
        curPos = io.GetPos() - startPos
                        </#if>
    }
                    </#if>
                <#--
                    Convert the list into an array. However if the array is of a primitive
                    type we have to iterate over it's elements and explicitly cast them.
                    Otherwise a simple toArray call is fine.
                -->
                    <#if helper.isSimpleTypeReference(field.type)>
    ${helper.getLanguageTypeNameForField(field)}[] ${manualArrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[_${manualArrayField.name}List.size()]
    for i := 0; i < _${manualArrayField.name}List.size(); i++ {
        ${manualArrayField.name}[i] = (${helper.getLanguageTypeNameForField(field)}) _${manualArrayField.name}List.get(i)
    }
                    <#else>
    ${helper.getLanguageTypeNameForField(field)}[] ${manualArrayField.name} = _${manualArrayField.name}List.toArray(new ${helper.getLanguageTypeNameForField(manualArrayField)}[0])
                    </#if>
                </#if>
                <#break>
            <#case "manual">
                <#assign manualField = field>

    // Manual Field (${manualField.name})
    ${manualField.name}, _${manualField.name}Err = ${helper.toParseExpression(manualField, manualField.parseExpression, type.parserArguments)}
    if _${implicitField.name}Err != nil {
        return nil, errors.New("Error parsing '${manualField.name}' field " + _${manualField.name}Err.Error())
    }
            <#break>
            <#case "optional">
                <#assign optionalField = field>

    // Optional Field (${optionalField.name}) (Can be skipped, if a given expression evaluates to false)
                <#if optionalField.conditionExpression.contains("curPos")>
    curPos = io.GetPos() - startPos
                </#if>
    var ${optionalField.name} *<#if helper.isComplexTypeReference(optionalField.type)>I</#if>${helper.getLanguageTypeNameForField(field)} = nil
    if ${helper.toBooleanParseExpression(optionalField.conditionExpression, type.parserArguments)} {
            <#if helper.isSimpleTypeReference(optionalField.type)>
        _val, _err := ${helper.getReadBufferReadMethodCall(optionalField.type)}
        if _err != nil {
            return nil, errors.New("Error parsing '${optionalField.name}' field " + _err.Error())
        }

        ${optionalField.name} = &_val
            <#else>
        _message, _err := ${helper.getLanguageTypeNameForField(field)}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument>${helper.toTypedParseExpression(helper.getArgumentType(optionalField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {
            return nil, errors.New("Error parsing '${optionalField.name}' field " + _err.Error())
        }
        var _item I${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(I${helper.getLanguageTypeNameForField(field)})
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to I${helper.getLanguageTypeNameForField(field)}")
        }
        ${optionalField.name} = &_item
            </#if>
    }
                <#break>
            <#case "padding">
                <#assign paddingField = field>
                <#assign simpleTypeReference = paddingField.type>

    // Padding Field (padding)
    {
        _timesPadding := (${helper.toParseExpression(paddingField, paddingField.paddingCondition, type.parserArguments)})
        for ;(io.HasMore(${helper.getNumBits(simpleTypeReference)})) && (_timesPadding > 0);_timesPadding-- {
            // Just read the padding data and ignore it
            _, _err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
            if _err != nil {
                return nil, errors.New("Error parsing 'padding' field " + _err.Error())
            }
        }
    }
                <#break>
            <#case "reserved">
                <#assign reservedField = field>
                <#assign simpleTypeReference = reservedField.type>

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != ${helper.getReservedValue(reservedField)} {
            log.WithFields(log.Fields{
                "expected value": ${helper.getReservedValue(reservedField)},
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }
                <#break>
            <#case "simple">
                <#assign simpleField = field>

    // Simple Field (${simpleField.name})
                <#if helper.isSimpleTypeReference(simpleField.type)>
                    <#assign simpleTypeReference = simpleField.type>
    ${simpleField.name}, _${simpleField.name}Err := ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
    if _${simpleField.name}Err != nil {
        return nil, errors.New("Error parsing '${simpleField.name}' field " + _${simpleField.name}Err.Error())
    }
                <#else>
                    <#assign complexTypeReference = simpleField.type>
    _${simpleField.name}Message, _err := ${complexTypeReference.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(simpleField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
    if _err != nil {
        return nil, errors.New("Error parsing simple field '${simpleField.name}'. " + _err.Error())
    }
    var ${simpleField.name} I${helper.getLanguageTypeNameForField(field)}
    ${simpleField.name}, _${simpleField.name}Ok := _${simpleField.name}Message.(I${helper.getLanguageTypeNameForField(field)})
    if !_${simpleField.name}Ok {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_${simpleField.name}Message).Name() + " to I${helper.getLanguageTypeNameForField(field)}")
    }
                </#if>
                <#break>

            <#case "switch">
                <#assign switchField = field>

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer ${type.name}Initializer
    var typeSwitchError error
    switch {
                <#list switchField.cases as case>
    <#if case.discriminatorValues?has_content>case <#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == ${discriminatorValue}<#sep> && </#sep></#list></#if>:
        initializer, typeSwitchError = ${case.name}Parse(io<#if case.parserArguments?has_content>, <#list case.parserArguments as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${parserArgument.name}<#sep>, </#sep></#list></#if>)
                </#list>
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }
                <#break>

            <#case "virtual">
                <#assign virtualField = field>

    // Virtual field (Just declare a local variable so we can access it in the parser)
    ${virtualField.name}, _${virtualField.name}Err := ${helper.toParseExpression(virtualField, virtualField.valueExpression, type.parserArguments)}
    if _${virtualField.name}Err != nil {
        return nil, errors.New("Error parsing '${virtualField.name}' field " + _${virtualField.name}Err.Error())
    }

                <#break>
        </#switch>
    </#list>

    // Create the instance
    <#if helper.isDiscriminatedParentTypeDefinition()>
    return initializer.initialize(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>), nil
    <#else>
    return New${type.name}(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>), nil
    </#if>
}

<#if helper.isDiscriminatedParentTypeDefinition()>
func ${type.name}Serialize(io utils.WriteBuffer, m ${type.name}, i I${type.name}, childSerialize func() error<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${parserArgument.name} ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) error {
<#else>
func (m ${type.name}) Serialize(io utils.WriteBuffer<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${parserArgument.name} ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) error {
</#if>
    <#assign arraySizeInBytesHelpers=helper.requiresHelperFunctions("ARRAY_SIZE_IN_BYTES")>
    <#if arraySizeInBytesHelpers?has_content>
        <#list arraySizeInBytesHelpers?keys as key>
            <#assign typeName=arraySizeInBytesHelpers[key]>
    ${key}ArraySizeInBytes := func(items []I${typeName}) uint32 {
        var sizeInBytes uint32 = 0
        for _, v := range items {
            sizeInBytes += uint32(v.LengthInBytes())
        }
        return sizeInBytes
    }
        </#list>
    </#if>
    <#if helper.isDiscriminatedChildTypeDefinition()>
    ser := func() error {
    </#if>
    <#list type.fields as field>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field>
                <#assign simpleTypeReference = arrayField.type>

    // Array Field (${arrayField.name})
    if m.${arrayField.name?cap_first} != nil {
                <#if helper.isComplexTypeReference(arrayField.type) && (helper.needsVariable(arrayField, "curItem", true) || helper.needsVariable(arrayField, "lastItem", true))>
        itemCount := uint16(len(m.${arrayField.name?cap_first}))
        var curItem uint16 = 0
                </#if>
        for _, _element := range m.${arrayField.name?cap_first} {
                <#if helper.isSimpleTypeReference(arrayField.type)>
                    <#assign simpleTypeReference = arrayField.type>
            _elementErr := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "_element")}
                <#else>
                    <#assign complexTypeReference = arrayField.type>
                    <#if helper.needsVariable(arrayField, "lastItem", true)>
            var lastItem bool = curItem == (itemCount - 1)
                    </#if>
            _elementErr := _element.Serialize(io<#if helper.getSerializerTerms(field.params)?has_content>, <#list helper.getSerializerTerms(field.params) as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>)
                </#if>
            if _elementErr != nil {
                return errors.New("Error serializing '${arrayField.name}' field " + _elementErr.Error())
            }
                <#if helper.needsVariable(arrayField, "curItem", true) || helper.needsVariable(arrayField, "lastItem", true)>
            curItem++
                </#if>
        }
    }
                <#break>
            <#case "checksum">
                <#assign checksumField = field>
                <#assign simpleTypeReference = checksumField.type>

    // Checksum Field (checksum) (Calculated)
    {
        _checksum := ${helper.getLanguageTypeNameForField(field)}(${helper.getNullValueForTypeReference(checksumField.type)})
        // Create an array of all the bytes written in this message element so far.
        checksumRawData []uint8 = io.getBytes(startPos, io.GetPos())
        _checksum = ${helper.getLanguageTypeNameForField(field)}(${helper.toSerializationExpression(checksumField, checksumField.checksumExpression, type.parserArguments)})
        _checksumErr := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(_checksum)")}
        if _elementErr != nil {
            return errors.New("Error serializing 'checksum' field " + _elementErr.Error())
        }
    }
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    _${constField.name}Err := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, constField.referenceValue)}
    if _${constField.name}Err != nil {
        return errors.New("Error serializing '${constField.name}' field " + _${constField.name}Err.Error())
    }
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
    ${discriminatorField.name} := ${helper.getLanguageTypeNameForField(field)}(i.${discriminatorField.name?cap_first}())
    _${discriminatorField.name}Err := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + discriminatorField.name + ")")}
    if _${discriminatorField.name}Err != nil {
        return errors.New("Error serializing '${discriminatorField.name}' field " + _${discriminatorField.name}Err.Error())
    }
                <#break>
            <#case "enum">
                <#assign enumField = field>

    // Enum field (${enumField.name})
    ${enumField.name} := Cast${helper.getLanguageTypeNameForField(field)}(m.${enumField.name?cap_first})
    _${enumField.name}Err := ${enumField.name}.Serialize(io)
    if _${enumField.name}Err != nil {
        return errors.New("Error serializing '${enumField.name}' field " + _${enumField.name}Err.Error())
    }
                <#break>
            <#case "implicit">
                <#assign implicitField = field>
                <#assign simpleTypeReference = implicitField.type>

    // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    ${implicitField.name} := ${helper.getLanguageTypeNameForField(field)}(${helper.toSerializationExpression(implicitField, implicitField.serializeExpression, type.parserArguments)})
    _${implicitField.name}Err := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + implicitField.name + ")")}
    if _${implicitField.name}Err != nil {
        return errors.New("Error serializing '${implicitField.name}' field " + _${implicitField.name}Err.Error())
    }
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field>

    // Manual Array Field (${manualArrayField.name})
    if m.${manualArrayField.name?cap_first} != nil {
        for(${helper.getLanguageTypeNameForField(field)} element : m.${manualArrayField.name?cap_first}) {
            ${helper.toSerializationExpression(manualArrayField, manualArrayField.serializeExpression, type.parserArguments)}
        }
    }
                <#break>
            <#case "manual">
                <#assign manualField = field>

    // Manual Field (${manualField.name})
    ${helper.toSerializationExpression(manualField, manualField.serializeExpression, type.parserArguments)}
                <#break>
            <#case "optional">
                <#assign optionalField = field>

    // Optional Field (${optionalField.name}) (Can be skipped, if the value is null)
    var ${optionalField.name} *<#if helper.isComplexTypeReference(optionalField.type)>I</#if>${helper.getLanguageTypeNameForField(field)} = nil
    if m.${optionalField.name?cap_first} != nil {
                <#if helper.isSimpleTypeReference(optionalField.type)>
                    <#assign simpleTypeReference = optionalField.type>
        ${optionalField.name} = m.${optionalField.name?cap_first}
        _${optionalField.name}Err := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "*(" + optionalField.name + ")")}
                <#else>
                    <#assign complexTypeReference = optionalField.type>
        ${optionalField.name} = m.${optionalField.name?cap_first}
        _${optionalField.name}Err := CastI${helper.getLanguageTypeNameForField(field)}(*${optionalField.name}).Serialize(io)
                </#if>
        if _${optionalField.name}Err != nil {
            return errors.New("Error serializing '${optionalField.name}' field " + _${optionalField.name}Err.Error())
        }
    }
                <#break>
            <#case "padding">
                <#assign paddingField = field>
                <#assign simpleTypeReference = paddingField.type>

    // Padding Field (padding)
    {
        _timesPadding := uint8(${helper.toSerializationExpression(paddingField, paddingField.paddingCondition, type.parserArguments)})
        for ;_timesPadding > 0; _timesPadding-- {
            _paddingValue := ${helper.getLanguageTypeNameForField(field)}(${helper.toSerializationExpression(paddingField, paddingField.paddingValue, type.parserArguments)})
            _paddingErr := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(_paddingValue)")}
            if _paddingErr != nil {
                return errors.New("Error serializing 'padding' field " + _paddingErr.Error())
            }
        }
    }
                <#break>
            <#case "reserved">
                <#assign reservedField = field>
                <#assign simpleTypeReference = reservedField.type>

    // Reserved Field (reserved)
    {
        _err := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, helper.getReservedValue(field))}
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }
                <#break>
            <#case "simple">
                <#assign simpleField = field>

    // Simple Field (${simpleField.name})
                <#if helper.isSimpleTypeReference(simpleField.type)>
                    <#assign simpleTypeReference = simpleField.type>
    ${simpleField.name} := ${helper.getLanguageTypeNameForField(field)}(m.${simpleField.name?cap_first})
    _${simpleField.name}Err := ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + simpleField.name + ")")}
                <#else>
                    <#assign complexTypeReference = simpleField.type>
    ${simpleField.name} := CastI${helper.getLanguageTypeNameForField(field)}(m.${simpleField.name?cap_first})
    _${simpleField.name}Err := ${simpleField.name}.Serialize(io)
                </#if>
    if _${simpleField.name}Err != nil {
        return errors.New("Error serializing '${simpleField.name}' field " + _${simpleField.name}Err.Error())
    }
                <#break>
            <#case "switch">
                <#assign switchField = field>

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }
                <#break>
            <#case "virtual">
                <#break>
        </#switch>
    </#list>

    <#if helper.isDiscriminatedChildTypeDefinition()>
        return nil
    }
    return ${type.parentType.name}Serialize(io, m.${type.parentType.name}, CastI${type.parentType.name}(m), ser)
    <#else>
    return nil
    </#if>
}

func (m *${type.name}) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
    <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
            case "${field.name}":
            <#-- int8 arrays are serialized in Base64 encoding -->
            <#if helper.isArrayField(field) && (helper.getLanguageTypeNameForField(field) == 'int8')>
                var _encoded string
                if err := d.DecodeElement(&_encoded, &tok); err != nil {
                    return err
                }
                _decoded := make([]byte, base64.StdEncoding.DecodedLen(len(_encoded)))
                _len, err := base64.StdEncoding.Decode(_decoded, []byte(_encoded))
                if err != nil {
                    return err
                }
                m.${field.name?cap_first} = utils.ByteToInt8(_decoded[0:_len])
            <#elseif helper.isComplexTypeReference(field.type) && helper.isDiscriminatedParentTypeDefinition(helper.getTypeDefinitionForTypeReference(field.type))>
                <#if helper.isArrayField(field)>
                var _values []<#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForTypeReference(field.type)}
                </#if>
                switch tok.Attr[0].Value {
                <#list helper.getSubTypeDefinitions(helper.getTypeDefinitionForTypeReference(field.type)) as subtype>
                    case "org.apache.plc4x.java.${protocolName}.${outputFlavor?replace("-","")}.${subtype.name}":
                        var dt *${subtype.name}
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                    <#if helper.isArrayField(field)>
                        _values = append(_values, dt)
                    <#else>
                        <#if helper.isOptionalField(field)>*</#if>m.${field.name?cap_first} = dt
                    </#if>
                </#list>
                    }
                <#if helper.isArrayField(field)>
                    m.${field.name?cap_first} = _values
                </#if>
            <#elseif helper.isArrayField(field)>
                var data []<#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForField(field)}
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.${field.name?cap_first} = data
            <#elseif helper.isOptionalField(field)>
                var data *<#if helper.isComplexTypeReference(field.type)>I</#if>${helper.getLanguageTypeNameForField(field)}
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.${field.name?cap_first} = data
            <#else>
                var data <#if helper.isComplexTypeReference(field.type)>*</#if>${helper.getLanguageTypeNameForField(field)}
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                <#if helper.isComplexTypeReference(field.type) && !helper.isEnumField(field)>
                m.${field.name?cap_first} = CastI${helper.getLanguageTypeNameForField(field)}(data)
                <#else>
                m.${field.name?cap_first} = data
                </#if>
            </#if>
        </#list>
    </#if>
            }
        }
    }
}

func (m ${type.name}) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.${protocolName}.${outputFlavor?replace("-","")}.${type.name}"},
        }}); err != nil {
        return err
    }
    <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
            <#-- int8 arrays are serialized in Base64 encoding -->
            <#if helper.isArrayField(field) && (helper.getLanguageTypeNameForField(field) == 'int8')>
    _encoded${field.name?cap_first} := make([]byte, base64.StdEncoding.EncodedLen(len(m.${field.name?cap_first})))
    base64.StdEncoding.Encode(_encoded${field.name?cap_first}, utils.Int8ToByte(m.${field.name?cap_first}))
    if err := e.EncodeElement(_encoded${field.name?cap_first}, xml.StartElement{Name: xml.Name{Local: "${field.name}"}}); err != nil {
        return err
    }
            <#-- normal arrays have one element with the same name wrapping them -->
            <#elseif helper.isArrayField(field)>
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "${field.name}"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.${field.name?cap_first}, xml.StartElement{Name: xml.Name{Local: "${field.name}"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "${field.name}"}}); err != nil {
        return err
    }
            <#else>
    if err := e.EncodeElement(m.${field.name?cap_first}, xml.StartElement{Name: xml.Name{Local: "${field.name}"}}); err != nil {
        return err
    }
            </#if>
        </#list>
    </#if>
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

</#outputformat>