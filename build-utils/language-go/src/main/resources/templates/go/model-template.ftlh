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

import (
    "math"
    "plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
)

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

<#if helper.isDiscriminatedParentTypeDefinition()>
type ${type.name}Initializer interface {
    initialize(<#list type.getAllPropertyFields() as field>${field.name}: ${field.name}<#sep>, </#sep></#list>) ${type.name}
}
<#else>
<#if type.parentType??>func (m ${type.name}) initialize(<#list type.parentType.getAllPropertyFields() as field>${field.name}: ${field.name}<#sep>, </#sep></#list>) ${type.parentType.name} {
    return m.${type.parentType.name}
}
</#if>

func New${type.name}(<#list type.getAllPropertyFields() as field>${field.name} <#if field.loopType??>[]</#if>${helper.getLanguageTypeNameForField(field)}<#sep>, </#sep></#list>) <#if type.parentType??>${type.parentType.name}Initializer<#else>${type.name}</#if> {
    return &${type.name}{<#list type.getAllPropertyFields() as field>${field.name}: ${field.name}<#sep>, </#sep></#list>}
}
</#if>

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

func ${type.name}Parse(io spi.ReadBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} ${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) <#if type.parentType??>${type.parentType.name}Initializer<#else>${type.name}</#if> {
    var startPos = io.GetPos()
    var curPos uint16
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
    if ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)} > math.MaxUint8 {
        throw new ParseException("Array count of " + (${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}) + " exceeds the maximum allowed count of " + math.MaxUint8);
    }
    ${helper.getLanguageTypeNameForField(field)}[] ${arrayField.name};
    {
        var itemCount := ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}
        ${arrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[itemCount]
        for curItem := 0; curItem < itemCount; curItem++ {
                <#-- When parsing simple types, there is nothing that could require the "lastItem" -->
            <#if !helper.isSimpleTypeReference(arrayField.type)>boolean lastItem = curItem == (itemCount - 1);</#if>
            ${arrayField.name}[curItem] = <#if helper.isSimpleTypeReference(arrayField.type)><#assign simpleTypeReference = arrayField.type>${helper.getReadBufferReadMethodCall(simpleTypeReference)}<#else>${arrayField.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(arrayField.type, parserArgument?index))}) (${helper.toParseExpression(arrayField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)<#if helper.isDiscriminatedChildTypeDefinition(helper.getTypeDefinitionForTypeReference(arrayField.type))>.build()</#if></#if>
        }
    }
                <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                <#else>
                <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                    <#if helper.isLengthArrayField(field)>
    // Length array
    var _${arrayField.name}Length := ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}
    List<${helper.getLanguageTypeNameForField(arrayField)}> _${arrayField.name}List = new LinkedList<>();
    ${arrayField.name}EndPos := io.GetPos() + _${arrayField.name}Length;
    for ;io.getPos() < ${arrayField.name}EndPos; {
        _${arrayField.name}List.add(<#if helper.isSimpleTypeReference(arrayField.type)><#assign simpleTypeReference = arrayField.type>${helper.getReadBufferReadMethodCall(simpleTypeReference)}<#else>${arrayField.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(arrayField.type, parserArgument?index))}) (${helper.toParseExpression(arrayField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);
                    <#-- After parsing, update the current position, but only if it's needed -->
                        <#if arrayField.loopExpression.contains("curPos")>
        curPos = io.getPos() - startPos;
                        </#if>
    }
                    <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                    <#elseif helper.isTerminatedArrayField(field)>
    // Terminated array
    List<${helper.getLanguageTypeNameForField(arrayField)}> _${arrayField.name}List = new LinkedList<>();
    for ;!((boolean) (${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)})); {
        _${arrayField.name}List.add(<#if helper.isSimpleTypeReference(arrayField.type)><#assign simpleTypeReference = arrayField.type>${helper.getReadBufferReadMethodCall(simpleTypeReference)}<#else>${arrayField.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(arrayField.type, parserArgument?index))}) (${helper.toParseExpression(arrayField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>)

                    <#-- After parsing, update the current position, but only if it's needed -->
                        <#if arrayField.loopExpression.contains("curPos")>
        curPos = io.getPos() - startPos
                        </#if>
    }
                    </#if>
                <#--
                    Convert the list into an array. However if the array is of a primitive
                    type we have to iterate over it's elements and explicitly cast them.
                    Otherwise a simple toArray call is fine.
                -->
                    <#if helper.isSimpleTypeReference(arrayField.type)>
    ${helper.getLanguageTypeNameForField(field)}[] ${arrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[_${arrayField.name}List.size()]
    for i := 0; i < _${arrayField.name}List.size(); i++ {
        ${arrayField.name}[i] = (${helper.getLanguageTypeNameForField(field)}) _${arrayField.name}List.get(i)
    }
                    <#else>
    ${helper.getLanguageTypeNameForField(field)}[] ${arrayField.name} = _${arrayField.name}List.toArray(new ${helper.getLanguageTypeNameForField(arrayField)}[0])
                    </#if>
                </#if>
                <#break>
            <#case "checksum">
                <#assign checksumField = field>
                <#assign simpleTypeReference = checksumField.type>

    // Checksum Field (checksum)
    {
        checksum = ${helper.getNullValueForTypeReference(checksumField.type)}
        // Create an array of all the bytes read in this message element so far.
        byte[] checksumRawData = io.getBytes(startPos, io.getPos())
        ${helper.getLanguageTypeNameForField(field)} checksumRef = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
        checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(checksumField, checksumField.checksumExpression, type.parserArguments)})
        if checksum != checksumRef {
            throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X", checksumRef & 0xFFFF, checksum & 0xFFFF));
        }
    }
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    ${helper.getLanguageTypeNameForField(field)} ${constField.name} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
    if ${constField.name} != ${type.name}.${constField.name?upper_case} {
        throw new ParseException("Expected constant value " + ${type.name}.${constField.name?upper_case} + " but got " + ${constField.name})
    }
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
    var ${discriminatorField.name} ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
                <#break>
            <#case "enum">
                <#assign enumField = field>

    // Enum field (${enumField.name})
    var ${enumField.name} ${helper.getLanguageTypeNameForField(field)} = ${helper.getLanguageTypeNameForField(field)}.valueOf(${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(enumField.type))})
                <#break>
            <#case "implicit">
                <#assign implicitField = field>
                <#assign simpleTypeReference = implicitField.type>

    // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    var ${implicitField.name} ${helper.getLanguageTypeNameForField(field)} = ${helper.getReadBufferReadMethodCall(simpleTypeReference)}
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field>

    // Manual Array Field (${manualArrayField.name})
            <#-- Only update curPos if the length expression uses it -->
                <#if manualArrayField.loopExpression.contains("curPos")>
    curPos = io.getPos() - startPos
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
    _${manualArrayField.name}Length := ${helper.toParseExpression(manualArrayField, manualArrayField.loopExpression, type.parserArguments)}
    List<${helper.getLanguageTypeNameForField(manualArrayField)}> _${manualArrayField.name}List = new LinkedList<>()
    ${manualArrayField.name}EndPos := io.getPos() + _${manualArrayField.name}Length
    for ;io.getPos() < ${manualArrayField.name}EndPos; {
        _${manualArrayField.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)}))
                    <#-- After parsing, update the current position, but only if it's needed -->
                        <#if manualArrayField.loopExpression.contains("curPos")>
        curPos = io.getPos() - startPos
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
        curPos = io.getPos() - startPos
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
    curPos = io.getPos() - startPos
                </#if>
    var ${optionalField.name} ${helper.getLanguageTypeNameForField(field)} = null
    if ${helper.toParseExpression(optionalField, optionalField.conditionExpression, type.parserArguments)} {
        ${optionalField.name} = <#if helper.isSimpleTypeReference(optionalField.type)>${helper.getReadBufferReadMethodCall(optionalField.type)}<#else>${optionalField.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(optionalField.type, parserArgument?index))}) (${helper.toParseExpression(optionalField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>
    }
                <#break>
            <#case "padding">
                <#assign paddingField = field>
                <#assign simpleTypeReference = paddingField.type>

    // Padding Field (padding)
    {
        _timesPadding := (${helper.toParseExpression(paddingField, paddingField.paddingCondition, type.parserArguments)})
        for ;(io.hasMore(${helper.getNumBits(simpleTypeReference)})) && (_timesPadding-- > 0); {
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
            LOGGER.info("Expected constant value " + ${reservedField.referenceValue} + " but got " + reserved + " for reserved field.")
        }
    }
                <#break>
            <#case "simple">
                <#assign simpleField = field>

    // Simple Field (${simpleField.name})
    var ${simpleField.name} ${helper.getLanguageTypeNameForField(field)} = <#if helper.isSimpleTypeReference(simpleField.type)><#assign simpleTypeReference = simpleField.type>${helper.getReadBufferReadMethodCall(simpleTypeReference)}<#else><#assign complexTypeReference = simpleField.type>${complexTypeReference.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(simpleField.type, parserArgument?index))}) (${helper.toParseExpression(simpleField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>
                <#break>
            <#case "switch">
                <#assign switchField = field>

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer ${type.name}Initializer
    switch {
                <#list switchField.cases as case>
    <#if case.discriminatorValues?has_content>case <#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == ${discriminatorValue}<#sep> && </#sep></#list></#if>:
        initializer = ${case.name}Parse(io<#if case.parserArguments?has_content>, <#list case.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>)
                </#list>
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
    return initializer.initialize(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>)
    <#else>
    return New${type.name}(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>)
    </#if>
}

    <#--if outputFlavor != "passive">
func (m ${type.name}) Serialize(io spi.WriteBuffer) {
    int startPos = io.getPos();

        <#list type.fields as field>
            <#switch field.typeName>
                <#case "array">
                    <#assign arrayField = field>
                    <#assign simpleTypeReference = arrayField.type>

    // Array Field (${arrayField.name})
    if(_value.get${arrayField.name?cap_first}() != null) {
        int itemCount = (int) _value.get${arrayField.name?cap_first}().length;
        int curItem = 0;
        for(${helper.getLanguageTypeNameForField(field)} element : _value.get${arrayField.name?cap_first}()) {
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
            ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "element")};
                    <#else>
                        <#assign complexTypeReference = arrayField.type>
            boolean lastItem = curItem == (itemCount - 1);
            ${complexTypeReference.name}IO.staticSerialize(io, element<#if helper.getSerializerTerms(field.params)?has_content>, <#list helper.getSerializerTerms(field.params) as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
                    </#if>
            curItem++;
        }
    }
                    <#break>
                <#case "checksum">
                    <#assign checksumField = field>
                    <#assign simpleTypeReference = checksumField.type>

    // Checksum Field (checksum) (Calculated)
    {
        ${helper.getLanguageTypeNameForField(field)} _checksum = ${helper.getNullValueForTypeReference(checksumField.type)};
        // Create an array of all the bytes written in this message element so far.
        byte[] checksumRawData = io.getBytes(startPos, io.getPos());
        _checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(checksumField, checksumField.checksumExpression, type.parserArguments)});
        ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(_checksum)")};
    }
                    <#break>
                <#case "const">
                    <#assign constField = field>
                    <#assign simpleTypeReference = constField.type>

    // Const Field (${constField.name})
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, constField.referenceValue)};
                    <#break>
                <#case "discriminator">
                    <#assign discriminatorField = field>
                    <#assign simpleTypeReference = discriminatorField.type>

    // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
    ${helper.getLanguageTypeNameForField(field)} ${discriminatorField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${discriminatorField.name?cap_first}();
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + discriminatorField.name + ")")};
                    <#break>
                <#case "enum">
                    <#assign enumField = field>

    // Enum field (${enumField.name})
    ${helper.getLanguageTypeNameForField(field)} ${enumField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${enumField.name?cap_first}();
    ${helper.getWriteBufferWriteMethodCall(helper.getEnumBaseTypeReference(enumField.type), "(" + enumField.name + ".getValue())")};
                    <#break>
                <#case "implicit">
                    <#assign implicitField = field>
                    <#assign simpleTypeReference = implicitField.type>

    // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    ${helper.getLanguageTypeNameForField(field)} ${implicitField.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(implicitField, implicitField.serializeExpression, type.parserArguments)});
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + implicitField.name + ")")};
                    <#break>
                <#case "manualArray">
                    <#assign manualArrayField = field>

    // Manual Array Field (${manualArrayField.name})
    if(_value.get${manualArrayField.name?cap_first}() != null) {
        for(${helper.getLanguageTypeNameForField(field)} element : _value.get${manualArrayField.name?cap_first}()) {
            ${helper.toSerializationExpression(manualArrayField, manualArrayField.serializeExpression, type.parserArguments)};
        }
    }
                    <#break>
                <#case "manual">
                    <#assign manualField = field>

    // Manual Field (${manualField.name})
    ${helper.toSerializationExpression(manualField, manualField.serializeExpression, type.parserArguments)};
                    <#break>
                <#case "optional">
                    <#assign optionalField = field>

    // Optional Field (${optionalField.name}) (Can be skipped, if the value is null)
    ${helper.getLanguageTypeNameForField(field)} ${optionalField.name} = null;
    if(_value.get${optionalField.name?cap_first}() != null) {
        ${optionalField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${optionalField.name?cap_first}();
                    <#if helper.isSimpleTypeReference(optionalField.type)>
                        <#assign simpleTypeReference = optionalField.type>
        ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + optionalField.name + ")")};
                    <#else>
                        <#assign complexTypeReference = optionalField.type>
        ${complexTypeReference.name}IO.staticSerialize(io, ${optionalField.name});
                    </#if>
    }
                    <#break>
                <#case "padding">
                    <#assign paddingField = field>
                    <#assign simpleTypeReference = paddingField.type>

    // Padding Field (padding)
    {
        int _timesPadding = (int) (${helper.toSerializationExpression(paddingField, paddingField.paddingCondition, type.parserArguments)});
        while (_timesPadding-- > 0) {
            ${helper.getLanguageTypeNameForField(field)} _paddingValue = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(paddingField, paddingField.paddingValue, type.parserArguments)});
            ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(_paddingValue)")};
        }
    }
                    <#break>
                <#case "reserved">
                    <#assign reservedField = field>
                    <#assign simpleTypeReference = reservedField.type>

    // Reserved Field (reserved)
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, helper.getReservedValue(field))};
                    <#break>
                <#case "simple">
                    <#assign simpleField = field>

    // Simple Field (${simpleField.name})
    ${helper.getLanguageTypeNameForField(field)} ${simpleField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${simpleField.name?cap_first}();
                    <#if helper.isSimpleTypeReference(simpleField.type)>
                        <#assign simpleTypeReference = simpleField.type>
    ${helper.getWriteBufferWriteMethodCall(simpleTypeReference, "(" + simpleField.name + ")")};
                    <#else>
                        <#assign complexTypeReference = simpleField.type>
    ${complexTypeReference.name}IO.staticSerialize(io, ${simpleField.name});
                    </#if>
                    <#break>
                <#case "switch">
                    <#assign switchField = field>

    // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
                    <#list switchField.cases as case>
    if(_value instanceof ${case.name}) {
        ${case.name}IO.staticSerialize(io, (${case.name}) _value);
    }<#sep> else </#sep>
                    </#list>
                    <#break>
                <#case "virtual">
                    <#break>
            </#switch>
        </#list>
}
    </#if-->
</#outputformat>