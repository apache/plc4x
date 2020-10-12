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
    "plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
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
    ${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if>${helper.getLanguageTypeNameForField(field)}
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
    Serialize(io spi.WriteBuffer)
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
    initialize(<#list type.getAllPropertyFields() as field>${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) spi.Message
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
<#if type.parentType??>func (m ${type.name}) initialize(<#list type.parentType.getAllPropertyFields() as field>${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) spi.Message {
    <#list type.parentType.getAllPropertyFields() as field>
    m.${field.name} = ${field.name}
    </#list>
    return m
}
</#if>

func New${type.name}(<#list type.getPropertyFields() as field>${field.name} <#if field.typeName == "optional">*</#if><#if field.loopType??>[]</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) <#if type.parentType??>${type.parentType.name}Initializer<#else>spi.Message</#if> {
    return &${type.name}{<#list type.getPropertyFields() as field>${field.name}: ${field.name}<#sep>, </#sep></#list>}
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

func ${type.name}Parse(io spi.ReadBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) (<#if type.parentType??>${type.parentType.name}Initializer<#else>spi.Message</#if>, error) {
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
    var ${arrayField.name} []${helper.getLanguageTypeNameForField(arrayField)}
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                <#if helper.isCountArrayField(field)>
    // Count array
    {
        ${arrayField.name} := make([]${helper.getLanguageTypeNameForField(field)}, ${helper.toIntegerParseExpression(16, arrayField.loopExpression, type.parserArguments)})
        for curItem := uint16(0); curItem < uint16(${helper.toIntegerParseExpression(16, arrayField.loopExpression, type.parserArguments)}); curItem++ {
<#if (!helper.isSimpleTypeReference(arrayField.type)) && helper.requiresVariable(arrayField, "lastItem")>            lastItem := curItem == uint16(${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)} - 1)</#if>
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
            ${arrayField.name} = append(${arrayField.name}, ${helper.getReadBufferReadMethodCall(simpleTypeReference)})
                    <#else>
            _message, _err := ${arrayField.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(arrayField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
            if _err != nil {
                return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
            }
            var _item ${helper.getLanguageTypeNameForField(field)}
            _item, _ok := _message.(${arrayField.type.name})
            if !_ok {
                return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ${arrayField.type.name}")
            }
            ${arrayField.name} = append(${arrayField.name}, _item)
                    </#if>
        }
    }
                <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                <#elseif helper.isLengthArrayField(field)>
                <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
    // Length array
    _${arrayField.name}Length := ${helper.toIntegerParseExpression(16, arrayField.loopExpression, type.parserArguments)}
    _${arrayField.name}EndPos := io.GetPos() + uint16(_${arrayField.name}Length)
    for ;io.GetPos() < _${arrayField.name}EndPos; {
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
        ${arrayField.name} = append(${arrayField.name}, ${helper.getReadBufferReadMethodCall(simpleTypeReference)})
                    <#else>
        _message, _err := ${arrayField.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(arrayField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {
            return nil, errors.New("Error parsing '${arrayField.name}' field " + _err.Error())
        }
        var _item ${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(${arrayField.type.name})
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
    for ;!bool(${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}); {
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
        ${arrayField.name} = append(${arrayField.name}, ${helper.getReadBufferReadMethodCall(simpleTypeReference)})
                    <#else>
        _message, _err := ${arrayField.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(arrayField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {
        }
        var _item ${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(${arrayField.type.name})
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
        byte[] checksumRawData = io.getBytes(startPos, io.GetPos())
        ${helper.getLanguageTypeNameForField(field)} checksumRef = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(checksumField, checksumField.checksumExpression, type.parserArguments)})
        if checksum != checksumRef {
            return nil, errors.New("Checksum verification failed. Expected " + (checksumRef & 0xFFFF) + " but got " + (checksum & 0xFFFF)")
        }
    }
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    var ${constField.name} ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
    if ${constField.name} != ${type.name}_${constField.name?upper_case} {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(${type.name}_${constField.name?upper_case})) + " but got " + strconv.Itoa(int(${constField.name})))
    }
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
    var ${helper.getVariableName(field)} ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
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
    var ${helper.getVariableName(field)} ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
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
    var ${manualField.name} ${helper.getLanguageTypeNameForField(field)} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualField, manualField.parseExpression, type.parserArguments)})
            <#break>
            <#case "optional">
                <#assign optionalField = field>

    // Optional Field (${optionalField.name}) (Can be skipped, if a given expression evaluates to false)
                <#if optionalField.conditionExpression.contains("curPos")>
    curPos = io.GetPos() - startPos
                </#if>
    var ${optionalField.name} *${helper.getLanguageTypeNameForField(field)} = nil
    if ${helper.toBooleanParseExpression(optionalField.conditionExpression, type.parserArguments)} {
            <#if helper.isSimpleTypeReference(optionalField.type)>
        _val := ${helper.getReadBufferReadMethodCall(optionalField.type)}
        ${optionalField.name} = &_val
            <#else>
        _message, _err := ${helper.getLanguageTypeNameForField(field)}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument>${helper.toTypedParseExpression(helper.getArgumentType(optionalField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
        if _err != nil {
            return nil, errors.New("Error parsing '${optionalField.name}' field " + _err.Error())
        }
        var _item ${helper.getLanguageTypeNameForField(field)}
        _item, _ok := _message.(${helper.getLanguageTypeNameForField(field)})
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ${helper.getLanguageTypeNameForField(field)}")
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
            ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        }
    }
                <#break>
            <#case "reserved">
                <#assign reservedField = field>
                <#assign simpleTypeReference = reservedField.type>

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        var reserved ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
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
    var ${simpleField.name} ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
                <#else>
                    <#assign complexTypeReference = simpleField.type>
    _${simpleField.name}Message, _err := ${complexTypeReference.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument><#if helper.isOptionalField(helper.getFieldForNameFromCurrentOrParent(parserArgument.name))>*</#if>${helper.toTypedParseExpression(helper.getArgumentType(simpleField.type, parserArgument?index), parserArgument, type.parserArguments)}<#sep>, </#sep></#list></#if>)
    if _err != nil {
        return nil, errors.New("Error parsing simple field '${simpleField.name}'. " + _err.Error())
    }
    var ${simpleField.name} ${helper.getLanguageTypeNameForField(field)}
    ${simpleField.name}, _${simpleField.name}Ok := _${simpleField.name}Message.(${helper.getLanguageTypeNameForField(field)})
    if !_${simpleField.name}Ok {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_${simpleField.name}Message).Name() + " to ${helper.getLanguageTypeNameForField(field)}")
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
    var ${virtualField.name} ${helper.getLanguageTypeNameForField(field)} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(virtualField, virtualField.valueExpression, type.parserArguments)})
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

func (m ${type.name}) Serialize(io spi.WriteBuffer<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${parserArgument.name} ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) {
    <#assign arraySizeInBytesHelpers=helper.requiresHelperFunctions("ARRAY_SIZE_IN_BYTES")>
    <#if arraySizeInBytesHelpers?has_content>
        <#list arraySizeInBytesHelpers?keys as key>
            <#assign typeName=arraySizeInBytesHelpers[key]>
    ${key}ArraySizeInBytes := func(items []${typeName}) uint32 {
        var sizeInBytes uint32 = 0
        for _, v := range items {
            sizeInBytes += uint32(v.LengthInBytes())
        }
        return sizeInBytes
    }
        </#list>
    </#if>
    <#if helper.isDiscriminatedParentTypeDefinition()>i${type.name} := CastI${type.name}(m)</#if>
    <#list type.fields as field>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field>
                <#assign simpleTypeReference = arrayField.type>

    // Array Field (${arrayField.name})
    if m.${arrayField.name} != nil {
                <#if helper.isComplexTypeReference(arrayField.type) && (helper.needsVariable(arrayField, "curItem", true) || helper.needsVariable(arrayField, "lastItem", true))>
        itemCount := uint16(len(m.${arrayField.name}))
        var curItem uint16 = 0
                </#if>
        for _, _element := range m.${arrayField.name} {
                <#if helper.isSimpleTypeReference(arrayField.type)>
                    <#assign simpleTypeReference = arrayField.type>
            ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "_element")}
                <#else>
                    <#assign complexTypeReference = arrayField.type>
                    <#if helper.needsVariable(arrayField, "lastItem", true)>
            var lastItem bool = curItem == (itemCount - 1)
                    </#if>
            _element.Serialize(io<#if helper.getSerializerTerms(field.params)?has_content>, <#list helper.getSerializerTerms(field.params) as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>)
                </#if>
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
        ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(_checksum)")}
    }
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, constField.referenceValue)}
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
    ${discriminatorField.name} := ${helper.getLanguageTypeNameForField(field)}(${type.name}${discriminatorField.name?cap_first}(i${type.name}))
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + discriminatorField.name + ")")}
                <#break>
            <#case "enum">
                <#assign enumField = field>

    // Enum field (${enumField.name})
    ${enumField.name} := ${helper.getLanguageTypeNameForField(field)}(m.${enumField.name})
    ${enumField.name}.Serialize(io)
                <#break>
            <#case "implicit">
                <#assign implicitField = field>
                <#assign simpleTypeReference = implicitField.type>

    // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    ${implicitField.name} := ${helper.getLanguageTypeNameForField(field)}(${helper.toSerializationExpression(implicitField, implicitField.serializeExpression, type.parserArguments)})
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + implicitField.name + ")")}
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field>

    // Manual Array Field (${manualArrayField.name})
    if m.${manualArrayField.name} != nil {
        for(${helper.getLanguageTypeNameForField(field)} element : m.${manualArrayField.name}) {
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
    var ${optionalField.name} *${helper.getLanguageTypeNameForField(field)} = nil
    if m.${optionalField.name} != nil {
        ${optionalField.name} = m.${optionalField.name}
                <#if helper.isSimpleTypeReference(optionalField.type)>
                    <#assign simpleTypeReference = optionalField.type>
        ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "*(" + optionalField.name + ")")}
                <#else>
                    <#assign complexTypeReference = optionalField.type>
        ${optionalField.name}.Serialize(io)
                </#if>
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
            ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(_paddingValue)")}
        }
    }
                <#break>
            <#case "reserved">
                <#assign reservedField = field>
                <#assign simpleTypeReference = reservedField.type>

    // Reserved Field (reserved)
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, helper.getReservedValue(field))}
                <#break>
            <#case "simple">
                <#assign simpleField = field>

    // Simple Field (${simpleField.name})
    ${simpleField.name} := ${helper.getLanguageTypeNameForField(field)}(m.${simpleField.name})
                <#if helper.isSimpleTypeReference(simpleField.type)>
                    <#assign simpleTypeReference = simpleField.type>
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + simpleField.name + ")")}
                <#else>
                    <#assign complexTypeReference = simpleField.type>
    ${simpleField.name}.Serialize(io)
                </#if>
                <#break>
            <#case "switch">
                <#assign switchField = field>

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    i${type.name}.Serialize(io)
                <#break>
            <#case "virtual">
                <#break>
        </#switch>
    </#list>
}
</#outputformat>