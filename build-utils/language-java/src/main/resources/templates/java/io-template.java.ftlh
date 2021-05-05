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
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.java.JavaLanguageTemplateHelper" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
<#-- Declare the name and type of variables declared locally inside the template -->
<#-- @ftlvariable name="typedField" type="org.apache.plc4x.plugins.codegenerator.types.fields.TypedField" -->
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
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/io/${type.name}IO.java
/*
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
*/

package ${helper.packageName(protocolName, languageName, outputFlavor)}.io;

import static org.apache.plc4x.java.spi.generation.StaticHelper.*;

import ${helper.packageName(protocolName, languageName, outputFlavor)}.*;
<#if helper.getComplexTypeReferences()?has_content>import ${helper.packageName(protocolName, languageName, outputFlavor)}.io.*;</#if>
import ${helper.packageName(protocolName, languageName, outputFlavor)}.types.*;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.api.value.PlcValue;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.math.BigInteger;
import java.time.*;
import java.util.*;
import java.util.function.Supplier;

// Code generated by build-utils. DO NOT EDIT.

public class ${type.name}IO implements <#if outputFlavor != "passive">MessageIO<${type.name}, ${type.name}><#else>MessageInput<${type.name}></#if> {

    private static final Logger LOGGER = LoggerFactory.getLogger(${type.name}IO.class);

<#-- The parse and serialize methods here are just proxies for forwardning the requests to static counterparts -->
<#if !helper.isDiscriminatedChildTypeDefinition()>
    @Override
    public ${type.name} parse(ReadBuffer readBuffer, Object... args) throws ParseException {
        <#if type.parserArguments?has_content>
        if((args == null) || (args.length != ${type.parserArguments?size})) {
            throw new PlcRuntimeException("Wrong number of arguments, expected ${type.parserArguments?size}, but got " + args.length);
        }
        <#list type.parserArguments as parserArgument>
        ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name};
        if(args[${parserArgument?index}] instanceof ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}) {
            ${parserArgument.name} = (${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}) args[${parserArgument?index}];
            <#if helper.isSimpleTypeReference(parserArgument.type)>
        } else if (args[${parserArgument?index}] instanceof String) {
            ${parserArgument.name} = ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)}.valueOf((String) args[${parserArgument?index}]);
            </#if>
        } else {
            throw new PlcRuntimeException("Argument ${parserArgument?index} expected to be of type ${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} or a string which is parseable but was " + args[${parserArgument?index}].getClass().getName());
        }
        </#list>
        </#if>
        return ${type.name}IO.staticParse(readBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
    }

    <#if outputFlavor != "passive">
    @Override
    public void serialize(WriteBuffer writeBuffer, ${type.name} value, Object... args) throws ParseException {
        <#if helper.getSerializerArguments(type.parserArguments)?has_content>
        if((args == null) || (args.length != ${type.parserArguments?size})) {
            throw new PlcRuntimeException("Wrong number of arguments, expected ${type.parserArguments?size}, but got " + args.length);
        }
        <#list helper.getSerializerArguments(type.parserArguments) as serializerArgument>
        if(!(args[${serializerArgument?index}] instanceof ${helper.getLanguageTypeNameForTypeReference(serializerArgument.type, false)})) {
            throw new PlcRuntimeException("Argument ${serializerArgument?index} expected to be of type ${helper.getLanguageTypeNameForTypeReference(serializerArgument.type, false)} but was " + args[${serializerArgument?index}].getClass().getName());
        }
        ${helper.getLanguageTypeNameForTypeReference(serializerArgument.type, false)} ${serializerArgument.name} = (${helper.getLanguageTypeNameForTypeReference(serializerArgument.type, false)}) args[${serializerArgument?index}];
        </#list>
        </#if>
        ${type.name}IO.staticSerialize(writeBuffer, value<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
    }

    </#if>
<#else>
    @Override
    public ${type.name} parse(ReadBuffer readBuffer, Object... args) throws ParseException {
        return (${type.name}) new ${type.parentType.name}IO().parse(readBuffer, args);
    }

<#if outputFlavor != "passive">
    @Override
    public void serialize(WriteBuffer writeBuffer, ${type.name} value, Object... args) throws ParseException {
        new ${type.parentType.name}IO().serialize(writeBuffer, value, args);
    }

</#if>
</#if>
<#-- Here come the actual parse and serialize methods that actually do the parsing and serlaizing -->
    public static ${type.name}<#if helper.isDiscriminatedChildTypeDefinition()>Builder</#if> staticParse(ReadBuffer readBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        readBuffer.pullContext("${type.name}");
        int startPos = readBuffer.getPos();
        int curPos;
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">
        <#assign arrayField = field>
        <#if helper.isByteBased(arrayField.type)>
        // Byte Array field (${arrayField.name})
            <#assign numberOfBytesExpression>
                <#compress>
                    <#if helper.isCountArrayField(field)>
                        int numberOfBytes = ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)};
                    <#elseif helper.isLengthArrayField(field)>
                        int numberOfBytes = ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)};
                    <#elseif helper.isTerminatedArrayField(field)>
                    <#-- TODO: we need to find out to implement this-->
                        int numberOfBytes := ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)};
                    <#else>
                    <#-- TODO: we should throw a exception here-->
                        int numberOfBytes := -1
                    </#if>
                </#compress>
            </#assign>
        ${numberOfBytesExpression}
        byte[] ${arrayField.name} = readBuffer.readByteArray("${arrayField.name}", numberOfBytes);
        <#else>
        // Array field (${arrayField.name})
        readBuffer.pullContext("${arrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
            <#-- Only update curPos if the length expression uses it -->
            <#if arrayField.loopExpression.contains("curPos")>
        curPos = readBuffer.getPos() - startPos;
            </#if>
        <#-- If this is a count array, we can directly initialize an array with the given size -->
            <#if helper.isCountArrayField(field)>
        // Count array
        if(${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)} > Integer.MAX_VALUE) {
            throw new ParseException("Array count of " + (${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}) + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
        }
        ${helper.getLanguageTypeNameForField(field)}[] ${arrayField.name};
        {
            int itemCount = Math.max(0, (int) ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)});
            ${arrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[itemCount];
            for(int curItem = 0; curItem < itemCount; curItem++) {
                <#-- When parsing simple types, there is nothing that could require the "lastItem" -->
                <#if !helper.isSimpleTypeReference(arrayField.type)>boolean lastItem = curItem == (itemCount - 1);</#if>
                ${arrayField.name}[curItem] = <#if helper.isSimpleTypeReference(arrayField.type)><#assign simpleTypeReference = arrayField.type>${helper.getReadBufferReadMethodCall("", simpleTypeReference, "", arrayField)}<#else>${arrayField.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(arrayField.type, parserArgument?index), true)}) (${helper.toParseExpression(arrayField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)<#if helper.isDiscriminatedChildTypeDefinition(helper.getTypeDefinitionForTypeReference(arrayField.type))>.build()</#if></#if>;
            }
        }
        <#-- In all other cases do we have to work with a list, that is later converted to an array -->
            <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                <#if helper.isLengthArrayField(field)>
        // Length array
        int _${arrayField.name}Length = ${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)};
        List<${helper.getNonPrimitiveLanguageTypeNameForField(arrayField)}> _${arrayField.name}List = new LinkedList<>();
        int ${arrayField.name}EndPos = readBuffer.getPos() + _${arrayField.name}Length;
        while(readBuffer.getPos() < ${arrayField.name}EndPos) {
            _${arrayField.name}List.add(<#if helper.isSimpleTypeReference(arrayField.type)><#assign simpleTypeReference = arrayField.type>${helper.getReadBufferReadMethodCall("", simpleTypeReference, "", arrayField)}<#else>${arrayField.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(arrayField.type, parserArgument?index), true)}) (${helper.toParseExpression(arrayField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);
            <#-- After parsing, update the current position, but only if it's needed -->
                <#if arrayField.loopExpression.contains("curPos")>
            curPos = readBuffer.getPos() - startPos;
                </#if>
        }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                <#elseif helper.isTerminatedArrayField(field)>
        // Terminated array
        List<${helper.getNonPrimitiveLanguageTypeNameForField(arrayField)}> _${arrayField.name}List = new LinkedList<>();
        while(!((boolean) (${helper.toParseExpression(arrayField, arrayField.loopExpression, type.parserArguments)}))) {
            _${arrayField.name}List.add(<#if helper.isSimpleTypeReference(arrayField.type)><#assign simpleTypeReference = arrayField.type>${helper.getReadBufferReadMethodCall("", simpleTypeReference, "", arrayField)}<#else>${arrayField.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(arrayField.type, parserArgument?index), true)}) (${helper.toParseExpression(arrayField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);

            <#-- After parsing, update the current position, but only if it's needed -->
                <#if arrayField.loopExpression.contains("curPos")>
            curPos = readBuffer.getPos() - startPos;
                </#if>
        }
                </#if>
                <#--
                    Convert the list into an array. However if the array is of a primitive
                    type we have to iterate over it's elements and explicitly cast them.
                    Otherwise a simple toArray call is fine.
                -->
                <#if helper.isSimpleTypeReference(arrayField.type)>
        ${helper.getLanguageTypeNameForField(field)}[] ${arrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[_${arrayField.name}List.size()];
        for(int i = 0; i < _${arrayField.name}List.size(); i++) {
            ${arrayField.name}[i] = (${helper.getLanguageTypeNameForField(field)}) _${arrayField.name}List.get(i);
        }
                <#else>
        ${helper.getLanguageTypeNameForField(field)}[] ${arrayField.name} = _${arrayField.name}List.toArray(new ${helper.getNonPrimitiveLanguageTypeNameForField(arrayField)}[0]);
                </#if>
            </#if>
            readBuffer.closeContext("${arrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
        </#if>
        <#break>
    <#case "checksum">
        <#assign checksumField = field>
        <#assign simpleTypeReference = checksumField.type>

        // Checksum Field (checksum)
        {
            ${helper.getLanguageTypeNameForField(field)} checksum = ${helper.getNullValueForTypeReference(checksumField.type)};
            ${helper.getLanguageTypeNameForField(field)} checksumRef = ${helper.getReadBufferReadMethodCall("checksum", simpleTypeReference, "", checksumField)};
            checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(checksumField, checksumField.checksumExpression, type.parserArguments)});
            if(checksum != checksumRef) {
                throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X", checksumRef & 0xFFFF, checksum & 0xFFFF));
            }
        }
        <#break>
    <#case "const">
        <#assign constField = field>
        <#assign simpleTypeReference = constField.type>

        // Const Field (${constField.name})
        ${helper.getLanguageTypeNameForField(field)} ${constField.name} = ${helper.getReadBufferReadMethodCall(constField.name, simpleTypeReference, "", constField)};
        if(${constField.name} != ${type.name}.${constField.name?upper_case}) {
            throw new ParseException("Expected constant value " + ${type.name}.${constField.name?upper_case} + " but got " + ${constField.name});
        }
        <#break>
    <#case "discriminator">
        <#assign discriminatorField = field>
        <#assign simpleTypeReference = discriminatorField.type>

        // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
        <#if helper.isEnumField(field)>
        ${helper.getLanguageTypeNameForField(field)} ${discriminatorField.name} = ${helper.getLanguageTypeNameForField(discriminatorField)}.enumForValue(${helper.getReadBufferReadMethodCall(discriminatorField.name, helper.getEnumBaseTypeReference(discriminatorField.type), "", discriminatorField)});
        <#else>
        ${helper.getLanguageTypeNameForField(field)} ${discriminatorField.name} = <#if helper.isSimpleTypeReference(discriminatorField.type)>${helper.getReadBufferReadMethodCall(discriminatorField.name, discriminatorField.type, "", discriminatorField)}<#else>${discriminatorField.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(discriminatorField.type, parserArgument?index), true)}) (${helper.toParseExpression(discriminatorField, parserArgument, type)})<#sep>, </#sep></#list></#if>)</#if>;
        </#if>

        <#break>
    <#case "enum">
        <#assign enumField = field>

        readBuffer.pullContext("${enumField.name}");
        // Enum field (${enumField.name})
        <#if enumField.fieldName?has_content>
        ${helper.getLanguageTypeNameForField(field)} ${enumField.name} = ${helper.getLanguageTypeNameForField(field)}.firstEnumForField${enumField.fieldName?cap_first}(${helper.getReadBufferReadMethodCall(helper.getTypeDefinitionForTypeReference(enumField.type).name, helper.getEnumBaseTypeReference(enumField.type), "", enumField)});
        <#else>
        ${helper.getLanguageTypeNameForField(field)} ${enumField.name} = ${helper.getLanguageTypeNameForField(field)}.enumForValue(${helper.getReadBufferReadMethodCall(helper.getTypeDefinitionForTypeReference(enumField.type).name, helper.getEnumBaseTypeReference(enumField.type), "", enumField)});
        </#if>
        readBuffer.closeContext("${enumField.name}");
        <#break>
    <#case "implicit">
        <#assign implicitField = field>
        <#assign simpleTypeReference = implicitField.type>

        // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForField(field)} ${implicitField.name} = ${helper.getReadBufferReadMethodCall(implicitField.name, simpleTypeReference, "", implicitField)};
        <#break>
    <#case "manualArray">
        <#assign manualArrayField = field>

        readBuffer.pullContext("${manualArrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
        // Manual Array Field (${manualArrayField.name})
        <#-- Only update curPos if the length expression uses it -->
        <#if manualArrayField.loopExpression.contains("curPos")>
        curPos = readBuffer.getPos() - startPos;
        </#if>
        <#-- If this is a count array, we can directly initialize an array with the given size -->
        <#if helper.isCountArrayField(field)>
        // Count array
        int _${manualArrayField.name}Count = ${helper.toParseExpression(manualArrayField, manualArrayField.loopExpression, type.parserArguments)};
        ${helper.getLanguageTypeNameForField(field)}[] ${manualArrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[_${manualArrayField.name}Count];
        for(int i = 0; i < _${manualArrayField.name}Count; i++) {
            ${manualArrayField.name}[i] = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)});
        }
        <#-- In all other cases do we have to work with a list, that is later converted to an array -->
        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
            <#if helper.isLengthArrayField(field)>
        // Length array
        int _${manualArrayField.name}Length = ${helper.toParseExpression(manualArrayField, manualArrayField.loopExpression, type.parserArguments)};
        List<${helper.getNonPrimitiveLanguageTypeNameForField(manualArrayField)}> _${manualArrayField.name}List = new LinkedList<>();
        int ${manualArrayField.name}EndPos = readBuffer.getPos() + _${manualArrayField.name}Length;
        while(readBuffer.getPos() < ${manualArrayField.name}EndPos) {
            _${manualArrayField.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)}));
            <#-- After parsing, update the current position, but only if it's needed -->
            <#if manualArrayField.loopExpression.contains("curPos")>
            curPos = readBuffer.getPos() - startPos;
            </#if>
        }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
            <#elseif helper.isTerminatedArrayField(field)>
        // Terminated array
        List<${helper.getNonPrimitiveLanguageTypeNameForField(manualArrayField)}> _${manualArrayField.name}List = new LinkedList<>();
        while(!((boolean) (${helper.toParseExpression(manualArrayField, manualArrayField.loopExpression, type.parserArguments)}))) {
            _${manualArrayField.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualArrayField, manualArrayField.parseExpression, type.parserArguments)}));

            <#-- After parsing, update the current position, but only if it's needed -->
            <#if manualArrayField.loopExpression.contains("curPos")>
            curPos = readBuffer.getPos() - startPos;
            </#if>
        }
            </#if>
            <#--
                Convert the list into an array. However if the array is of a primitive
                type we have to iterate over it's elements and explicitly cast them.
                Otherwise a simple toArray call is fine.
            -->
            <#if helper.isSimpleTypeReference(field.type)>
        ${helper.getLanguageTypeNameForField(field)}[] ${manualArrayField.name} = new ${helper.getLanguageTypeNameForField(field)}[_${manualArrayField.name}List.size()];
        for(int i = 0; i < _${manualArrayField.name}List.size(); i++) {
            ${manualArrayField.name}[i] = (${helper.getLanguageTypeNameForField(field)}) _${manualArrayField.name}List.get(i);
        }
            <#else>
        ${helper.getLanguageTypeNameForField(field)}[] ${manualArrayField.name} = _${manualArrayField.name}List.toArray(new ${helper.getNonPrimitiveLanguageTypeNameForField(manualArrayField)}[0]);
            </#if>
        </#if>
        readBuffer.closeContext("${manualArrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
        <#break>
    <#case "manual">
        <#assign manualField = field>

        // Manual Field (${manualField.name})
        ${helper.getLanguageTypeNameForField(field)} ${manualField.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(manualField, manualField.parseExpression, type.parserArguments)});
        <#break>
    <#case "optional">
        <#assign optionalField = field>

        // Optional Field (${optionalField.name}) (Can be skipped, if a given expression evaluates to false)
        <#if optionalField.conditionExpression.contains("curPos")>
        curPos = readBuffer.getPos() - startPos;
        </#if>
        ${helper.getLanguageTypeNameForField(field)} ${optionalField.name} = null;
        if(${helper.toParseExpression(optionalField, optionalField.conditionExpression, type.parserArguments)}) {
            <#if helper.isEnumField(field)>
            ${optionalField.name} = ${helper.getLanguageTypeNameForField(optionalField)}.enumForValue(${helper.getReadBufferReadMethodCall(optionalField.name, helper.getEnumBaseTypeReference(optionalField.type), "", optionalField)});
            <#else>
            ${optionalField.name} = <#if helper.isSimpleTypeReference(optionalField.type)>${helper.getReadBufferReadMethodCall(optionalField.name, optionalField.type, "", optionalField)}<#else>${optionalField.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(optionalField.type, parserArgument?index), true)}) (${helper.toParseExpression(optionalField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
            </#if>
        }
        <#break>
    <#case "padding">
        <#assign paddingField = field>
        <#assign simpleTypeReference = paddingField.type>

        // Padding Field (padding)
        {
            readBuffer.pullContext("padding", WithReaderWriterArgs.WithRenderAsList(true));
            int _timesPadding = (int) (${helper.toParseExpression(paddingField, paddingField.paddingCondition, type.parserArguments)});
            while ((readBuffer.hasMore(${helper.getNumBits(simpleTypeReference)})) && (_timesPadding-- > 0)) {
                // Just read the padding data and ignore it
                ${helper.getReadBufferReadMethodCall(simpleTypeReference, "", paddingField)};
            }
            readBuffer.closeContext("padding", WithReaderWriterArgs.WithRenderAsList(true));
        }
        <#break>
    <#case "reserved">
        <#assign reservedField = field>
        <#assign simpleTypeReference = reservedField.type>

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        {
            ${helper.getLanguageTypeNameForField(field)} reserved = ${helper.getReadBufferReadMethodCall("reserved", simpleTypeReference, "", reservedField)};
            if(reserved != ${helper.getReservedValue(reservedField)}) {
                LOGGER.info("Expected constant value " + ${reservedField.referenceValue} + " but got " + reserved + " for reserved field.");
            }
        }
        <#break>
    <#case "simple">
        <#assign simpleField = field>

        <#if !helper.isSimpleTypeReference(simpleField.type)>
        readBuffer.pullContext("${simpleField.name}");
        </#if>

        // Simple Field (${simpleField.name})
        <#if helper.isEnumField(field)>
        // enum based simple field with type ${simpleField.type.name}
        ${helper.getLanguageTypeNameForField(simpleField)} ${simpleField.name} = ${helper.getLanguageTypeNameForField(simpleField)}.enumForValue(${helper.getReadBufferReadMethodCall(simpleField.type.name, helper.getEnumBaseTypeReference(simpleField.type), "", simpleField)});
        <#else>
        <#assign simpleFieldLogicalName><#if helper.isSimpleTypeReference(simpleField.type)>${simpleField.name}<#else>${simpleField.typeName}</#if></#assign>
        ${helper.getLanguageTypeNameForField(simpleField)} ${simpleField.name} = <#if helper.isSimpleTypeReference(simpleField.type)><#assign simpleTypeReference = simpleField.type>${helper.getReadBufferReadMethodCall(simpleFieldLogicalName, simpleTypeReference, "", simpleField)}<#else><#assign complexTypeReference = simpleField.type>${complexTypeReference.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(simpleField.type, parserArgument?index), true)}) (${helper.toParseExpression(simpleField, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
        </#if>
        <#if !helper.isSimpleTypeReference(simpleField.type)>
        readBuffer.closeContext("${simpleField.name}");
        </#if>
        <#break>
    <#case "switch">
        <#assign switchField = field>

        // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        ${type.name}Builder builder = null;
        <#list switchField.cases as case>
                <#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toParseExpression(null, switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)}, <#if helper.discriminatorValueNeedsStringEqualityCheck(switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#elseif helper.isComplexTypeReference(helper.getDiscriminatorTypes()[field.discriminatorExpressions[discriminatorValue?index].name])><#if helper.isEnumTypeReference(helper.getDiscriminatorTypes()[field.discriminatorExpressions[discriminatorValue?index].name])>${helper.getDiscriminatorTypes()[field.discriminatorExpressions[discriminatorValue?index].name].name}.${discriminatorValue}<#else>${discriminatorValue}</#if><#else>${discriminatorValue}</#if>)<#sep> && </#sep></#list>) </#if>{
            builder = ${case.name}IO.staticParse(readBuffer<#if case.parserArguments?has_content>, <#list case.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
        }<#sep> else </#sep>
        </#list>
        if (builder == null) {
            throw new ParseException("Unsupported case for discriminated type");
        }
        <#break>

    <#case "virtual">
        <#assign virtualField = field>

        // Virtual field (Just declare a local variable so we can access it in the parser)
        <#if helper.getLanguageTypeNameForField(field) = 'String'>
        ${helper.getLanguageTypeNameForField(field)} ${virtualField.name} = ${helper.getLanguageTypeNameForField(field)}.valueOf(${helper.toParseExpression(virtualField, virtualField.valueExpression, type.parserArguments)});
        <#else>
        ${helper.getLanguageTypeNameForField(field)} ${virtualField.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(virtualField, virtualField.valueExpression, type.parserArguments)});
        </#if>
        <#break>
</#switch>
</#list>

        readBuffer.closeContext("${type.name}");
        // Create the instance
        <#if helper.isDiscriminatedChildTypeDefinition()>
        return new ${type.name}Builder(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        <#elseif helper.isDiscriminatedParentTypeDefinition()>
        return builder.build(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        <#else>
        return new ${type.name}(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        </#if>
    }

<#if outputFlavor != "passive">
    public static void staticSerialize(WriteBuffer writeBuffer, ${type.name} _value<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        int startPos = writeBuffer.getPos();
        writeBuffer.pushContext("${type.name}");
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">
        <#assign arrayField = field>
        <#assign simpleTypeReference = arrayField.type>

        // Array Field (${arrayField.name})
        if(_value.get${arrayField.name?cap_first}() != null) {
            <#if helper.isByteBased(arrayField.type)>
            // Byte Array field (${arrayField.name})
            writeBuffer.writeByteArray("${arrayField.name}", _value.get${arrayField.name?cap_first}());
            <#else>
            writeBuffer.pushContext("${arrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
            int itemCount = (int) _value.get${arrayField.name?cap_first}().length;
            int curItem = 0;
            for(${helper.getLanguageTypeNameForField(field)} element : _value.get${arrayField.name?cap_first}()) {
                    <#if helper.isSimpleTypeReference(arrayField.type)>
                        <#assign simpleTypeReference = arrayField.type>
                ${helper.getWriteBufferWriteMethodCall("", simpleTypeReference, "element", field)};
                    <#else>
                        <#assign complexTypeReference = arrayField.type>
                boolean lastItem = curItem == (itemCount - 1);
                ${complexTypeReference.name}IO.staticSerialize(writeBuffer, element<#if helper.getSerializerTerms(field.params)?has_content>, <#list helper.getSerializerTerms(field.params) as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
                    </#if>
                curItem++;
            }
            writeBuffer.popContext("${arrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
            </#if>
        }
        <#break>
    <#case "checksum">
        <#assign checksumField = field>
        <#assign simpleTypeReference = checksumField.type>

        // Checksum Field (checksum) (Calculated)
        {
            ${helper.getLanguageTypeNameForField(field)} _checksum = ${helper.getNullValueForTypeReference(checksumField.type)};
            _checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(checksumField, checksumField.checksumExpression, type.parserArguments)});
            ${helper.getWriteBufferWriteMethodCall("checksum", simpleTypeReference, "(_checksum)", field)};
        }
        <#break>
    <#case "const">
        <#assign constField = field>
        <#assign simpleTypeReference = constField.type>

        // Const Field (${constField.name})
        <#if helper.getLanguageTypeNameForField(field) = 'float'>${helper.getWriteBufferWriteMethodCall(simpleTypeReference, constField.referenceValue + "f", field)};
        <#else>${helper.getWriteBufferWriteMethodCall(constField.name, simpleTypeReference, constField.referenceValue, field)};
        </#if>
        <#break>
    <#case "discriminator">
        <#assign discriminatorField = field>
        <#assign simpleTypeReference = discriminatorField.type>

        // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
        ${helper.getLanguageTypeNameForField(field)} ${discriminatorField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${discriminatorField.name?cap_first}();
        <#if helper.isSimpleTypeReference(discriminatorField.type)>
            <#assign simpleTypeReference = discriminatorField.type>
            ${helper.getWriteBufferWriteMethodCall(discriminatorField.name, simpleTypeReference, "(" + discriminatorField.name + ")", field)};
        <#else>
            <#assign complexTypeReference = discriminatorField.type>
            <#if helper.isEnumField(field)>
                ${helper.getWriteBufferWriteMethodCall(discriminatorField.name, helper.getEnumBaseTypeReference(discriminatorField.type), "(" + discriminatorField.name + ".getValue())", field)};
            <#else>
                ${complexTypeReference.name}IO.staticSerialize(writeBuffer, ${discriminatorField.name});
            </#if>
        </#if>
        <#break>
    <#case "enum">
        <#assign enumField = field>

        writeBuffer.pushContext("${enumField.name}");
        // Enum field (${enumField.name})
        ${helper.getLanguageTypeNameForField(field)} ${enumField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${enumField.name?cap_first}();
        <#if enumField.fieldName?has_content>
        ${helper.getWriteBufferWriteMethodCall(helper.getTypeDefinitionForTypeReference(enumField.type).name, helper.getEnumFieldSimpleTypeReference(enumField.type, enumField.fieldName), "(" + enumField.name + ".get" + enumField.fieldName?cap_first + "())", field, "WithReaderWriterArgs.WithAdditionalStringRepresentation(${enumField.name}.name())")};
        <#else>
        ${helper.getWriteBufferWriteMethodCall(helper.getTypeDefinitionForTypeReference(enumField.type).name, helper.getEnumBaseTypeReference(enumField.type), "(" + enumField.name + ".getValue())", field, "WithReaderWriterArgs.WithAdditionalStringRepresentation(${enumField.name}.name())")};
        </#if>
        writeBuffer.popContext("${enumField.name}");
        <#break>
    <#case "implicit">
        <#assign implicitField = field>
        <#assign simpleTypeReference = implicitField.type>

        // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForField(field)} ${implicitField.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(implicitField, implicitField.serializeExpression, type.parserArguments)});
        ${helper.getWriteBufferWriteMethodCall(implicitField.name, simpleTypeReference, "(" + implicitField.name + ")", field)};
        <#break>
    <#case "manualArray">
        <#assign manualArrayField = field>

        // Manual Array Field (${manualArrayField.name})
        if(_value.get${manualArrayField.name?cap_first}() != null) {
            writeBuffer.pushContext("${manualArrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
            for(${helper.getLanguageTypeNameForField(field)} element : _value.get${manualArrayField.name?cap_first}()) {
                ${helper.toSerializationExpression(manualArrayField, manualArrayField.serializeExpression, type.parserArguments)};
            }
            writeBuffer.popContext("${manualArrayField.name}", WithReaderWriterArgs.WithRenderAsList(true));
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
            ${helper.getWriteBufferWriteMethodCall(optionalField.name, simpleTypeReference, "(" + optionalField.name + ")", field)};
            <#else>
                <#assign complexTypeReference = optionalField.type>
                <#if helper.isEnumField(field)>
            ${helper.getWriteBufferWriteMethodCall(optionalField.name, helper.getEnumBaseTypeReference(optionalField.type), "(" + optionalField.name + ".getValue())", field)};
                <#else>
            ${complexTypeReference.name}IO.staticSerialize(writeBuffer, ${optionalField.name});
                </#if>
            </#if>
        }
        <#break>
    <#case "padding">
        <#assign paddingField = field>
        <#assign simpleTypeReference = paddingField.type>

        // Padding Field (padding)
        {
            writeBuffer.pushContext("padding", WithReaderWriterArgs.WithRenderAsList(true));
            int _timesPadding = (int) (${helper.toSerializationExpression(paddingField, paddingField.paddingCondition, type.parserArguments)});
            while (_timesPadding-- > 0) {
                ${helper.getLanguageTypeNameForField(field)} _paddingValue = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(paddingField, paddingField.paddingValue, type.parserArguments)});
                ${helper.getWriteBufferWriteMethodCall("", simpleTypeReference, "(_paddingValue)", field)};
            }
            writeBuffer.popContext("padding", WithReaderWriterArgs.WithRenderAsList(true));
        }
        <#break>
    <#case "reserved">
        <#assign reservedField = field>
        <#assign simpleTypeReference = reservedField.type>

        // Reserved Field (reserved)
        ${helper.getWriteBufferWriteMethodCall("reserved", simpleTypeReference, helper.getReservedValue(field), field)};
        <#break>
    <#case "simple">
        <#assign simpleField = field>

        // Simple Field (${simpleField.name})
        ${helper.getLanguageTypeNameForField(field)} ${simpleField.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${simpleField.name?cap_first}();
        <#if helper.isSimpleTypeReference(simpleField.type)>
            <#assign simpleTypeReference = simpleField.type>
            <#if helper.isEnumField(field)>
        // enum based simple field with type ${simpleField.type.name}
        ${helper.getWriteBufferWriteMethodCall(simpleField.type.name, simpleTypeReference, "(" + simpleField.name + ")", simpleField, "WithReaderWriterArgs.WithAdditionalStringRepresentation(${simpleField.name}.name())")};
            <#else>
        ${helper.getWriteBufferWriteMethodCall(simpleField.name, simpleTypeReference, "(" + simpleField.name + ")", simpleField)};
            </#if>
        <#else>
        writeBuffer.pushContext("${simpleField.name}");
            <#assign complexTypeReference = simpleField.type>
            <#if helper.isEnumField(field)>
        // enum field with type ${simpleField.type.name}
        ${helper.getWriteBufferWriteMethodCall(simpleField.type.name, helper.getEnumBaseTypeReference(simpleField.type), "(" + simpleField.name + ".getValue())", simpleField, "WithReaderWriterArgs.WithAdditionalStringRepresentation(${simpleField.name}.name())")};
            <#else>
        ${complexTypeReference.name}IO.staticSerialize(writeBuffer, ${simpleField.name});
            </#if>
        writeBuffer.popContext("${simpleField.name}");
        </#if>
        <#break>
    <#case "switch">
        <#assign switchField = field>

        // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
        <#list switchField.cases as case>
        if(_value instanceof ${case.name}) {
            ${case.name}IO.staticSerialize(writeBuffer, (${case.name}) _value);
        }<#sep> else </#sep>
        </#list>
        <#break>
    <#case "virtual">
        <#break>
</#switch>
</#list>
        writeBuffer.popContext("${type.name}");
    }
</#if>

<#if helper.isDiscriminatedParentTypeDefinition()>
    public static interface ${type.name}Builder {
        ${type.name} build(<#list type.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>);
    }

</#if>
<#if helper.isDiscriminatedChildTypeDefinition()>
    public static class ${type.name}Builder implements ${type.parentType.name}IO.${type.parentType.name}Builder {
        <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
        private final ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name};
        </#list>
        </#if>

        public ${type.name}Builder(<#list type.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            <#list type.propertyFields as field>
            this.${field.name} = ${field.name};
            </#list>
        }

        public ${type.name} build(<#list type.parentType.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            return new ${type.name}(<#list type.getAllPropertyFields() as field>${field.name}<#sep>, </#sep></#list>);
        }
    }

</#if>
}
</#outputformat>