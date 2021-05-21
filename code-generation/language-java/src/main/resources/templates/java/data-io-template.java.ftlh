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

import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.generation.EvaluationHelper;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.WriteBufferByteBased;
import ${helper.packageName(protocolName, languageName, outputFlavor)}.*;
import ${helper.packageName(protocolName, languageName, outputFlavor)}.types.*;
import org.apache.plc4x.java.spi.values.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.math.BigInteger;
import java.time.*;
import java.util.*;
import java.util.function.Supplier;

// Code generated by code-generation. DO NOT EDIT.

public class ${type.name}IO {

    private static final Logger LOGGER = LoggerFactory.getLogger(${type.name}IO.class);
    public static PlcValue staticParse(ReadBuffer readBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        <#assign defaultCaseOutput=false>
        <#list type.switchField.cases as case>
<#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)}, <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#elseif helper.isComplexTypeReference(type.parserArguments[discriminatorValue?index].type)><#if helper.isEnumTypeReference(type.parserArguments[discriminatorValue?index].type)>${helper.getLanguageTypeNameForTypeReference(type.parserArguments[discriminatorValue?index].type, false)}.${discriminatorValue}<#else>${discriminatorValue}</#if><#else>${discriminatorValue}</#if>)<#sep> && </#sep></#list>) <#else><#assign defaultCaseOutput=true></#if>{ // ${case.name}
            <#assign valueDefined=false>
            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">

            // Array field (${field.name})
            <#-- Only update curPos if the length expression uses it -->
                        <#if field.loopExpression.contains("curPos")>
            curPos = readBuffer.getPos() - startPos;
                        </#if>
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                        <#if helper.isCountArrayField(field)>
            // Count array
            if(${helper.toParseExpression(field, field.loopExpression, type.parserArguments)} > Integer.MAX_VALUE) {
                throw new ParseException("Array count of " + (${helper.toParseExpression(field, field.loopExpression, type.parserArguments)}) + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
            }
            List<PlcValue> ${field.name};
            {
                int itemCount = (int) ${helper.toParseExpression(field, field.loopExpression, type.parserArguments)};
                ${field.name} = new LinkedList<>();
                for(int curItem = 0; curItem < itemCount; curItem++) {
                    ${field.name}.add(new ${helper.getPlcValueTypeForTypeReference(field.type)}((${helper.getNonPrimitiveLanguageTypeNameForField(field)}) <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type, "", field)})<#else>${field.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);
                }
            }
            <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                            <#if helper.isLengthArrayField(field)>
            // Length array
            int _${field.name}Length = ${helper.toParseExpression(field, field.loopExpression, type.parserArguments)};
            int ${field.name}EndPos = readBuffer.getPos() + _${field.name}Length;
            List<PlcValue> value = new LinkedList<>();
            while(readBuffer.getPos() < ${field.name}EndPos) {
                value.add(
                <#if helper.isSimpleTypeReference(field.type)>
                    new ${helper.getPlcValueTypeForTypeReference(field.type)}(${helper.getReadBufferReadMethodCall(field.type, "", field)})
                <#else>${field.type.name}IO.staticParse(readBuffer
                    <#if field.params?has_content>,
                        <#list field.params as parserArgument>
                            (${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})
                            <#sep>, </#sep>
                        </#list>
                    </#if>
                    )
                </#if>
                );
            }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                            <#elseif helper.isTerminatedArrayField(field)>
            // Terminated array
            List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> ${field.name} = new LinkedList<>();
            while(!((boolean) (${helper.toParseExpression(field, field.loopExpression, type.parserArguments)}))) {
                ${field.name}.add(<#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type, "", field)}<#else>${field.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);

                <#-- After parsing, update the current position, but only if it's needed -->
                                <#if field.loopExpression.contains("curPos")>
                curPos = readBuffer.getPos() - startPos;
                                </#if>
            }
                            </#if>
                        </#if>
                        <#if field.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "const">

            // Const Field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type, "", field)};
            if(${field.name} != ${type.name}.${field.name?upper_case}) {
                throw new ParseException("Expected constant value " + ${type.name}.${field.name?upper_case} + " but got " + ${field.name});
            }
                        <#if field.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "enum">

            // Enum field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getNonPrimitiveLanguageTypeNameForField(field)}.enumForValue(${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(field.type), "", field)});
                        <#if field.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "manual">

            // Manual Field (${field.name})
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field, field.parseExpression, type.parserArguments)});
                        <#if field.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                    <#case "reserved">

            // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            {
                ${helper.getLanguageTypeNameForField(field)} reserved = ${helper.getReadBufferReadMethodCall(field.type, "", field)};
                if(reserved != ${helper.getReservedValue(field)}) {
                    LOGGER.info("Expected constant value " + ${field.referenceValue} + " but got " + reserved + " for reserved field.");
                }
            }
                    <#break>
                    <#case "simple">

                        <#if helper.isEnumField(field)>
            // Enum field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getNonPrimitiveLanguageTypeNameForField(field)}.enumForValue(${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(field.type), "", field)});
                        <#else>
            // Simple Field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type, "", field)}<#else>${field.type.name}IO.staticParse(readBuffer<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
                        </#if>
                        <#if case.name == "Struct" ||
                            ((case.name == "DATE_AND_TIME") && ((field.name == "year") || (field.name == "month") || (field.name == "day") || (field.name == "hour") || (field.name == "minutes") || (field.name == "seconds"))) ||
                            ((case.name == "DATE_AND_TIME") && (field.name == "secondsSinceEpoch")) ||
                            ((case.name == "DATE") && ((field.name == "year") || (field.name == "month") || (field.name == "day"))) ||
                            ((case.name == "TIME_OF_DAY") && ((field.name == "hour") || (field.name == "minutes") || (field.name == "seconds"))) ||
                            field.name == "value">
                            <#assign valueDefined=true>
                        </#if>
                    <#break>
                </#switch>
            </#list>
            <#if case.name == "Struct">

                <#-- In this case we need to wrap each field in a PlcValue that matches it's natural type -->
            Map<String, PlcValue> _map = new HashMap<>();
                <#list case.fields as field>
                    <#if helper.isArrayField(field)>
            _map.put("${field.name}", new PlcList(${field.name}));
                    <#elseif helper.isPropertyField(field)>
                        <#switch helper.getLanguageTypeNameForTypeReference(field.type)>
                            <#case "Boolean">
            _map.put("${field.name}", new PlcBOOL(${field.name}));
                                <#break>
                            <#case "Byte">
            _map.put("${field.name}", new PlcSINT(${field.name}));
                                <#break>
                            <#case "Short">
            _map.put("${field.name}", new PlcINT(${field.name}));
                                <#break>
                            <#case "Integer">
            _map.put("${field.name}", new PlcDINT(${field.name}));
                                <#break>
                            <#case "Long">
            _map.put("${field.name}", new PlcLINT(${field.name}));
                                <#break>
                            <#case "BigInteger">
            _map.put("${field.name}", new PlcBigInteger(${field.name}));
                                <#break>
                            <#case "Float">
            _map.put("${field.name}", new PlcREAL(${field.name}));
                                <#break>
                            <#case "Double">
            _map.put("${field.name}", new PlcLREAL(${field.name}));
                                <#break>
                            <#case "BigDecimal">
            _map.put("${field.name}", new PlcBigDecimal(${field.name}));
                                <#break>
                            <#case "String">
            _map.put("${field.name}", new PlcSTRING(${field.name}));
                                <#break>
                            <#case "LocalTime">
            _map.put("${field.name}", new PlcTIME_OF_DAY(${field.name}));
                                <#break>
                            <#case "LocalDate">
            _map.put("${field.name}", new PlcDATE(${field.name}));
                                <#break>
                            <#case "LocalDateTime">
            _map.put("${field.name}", new PlcDATE_AND_TIME(${field.name}));
                                <#break>
                        </#switch>
                    </#if>
                </#list>
                <#assign valueDefined=true>
            </#if>

            <#if valueDefined>
                <#switch case.name>
                    <#case "TIME">
            return new PlcTIME(value);
                    <#break>
                    <#case "DATE">
                    <#if helper.hasFieldsWithNames(case.fields, "year", "month", "day")>
            LocalDate value = LocalDate.of(year.intValue(), (month == 0) ? 1 : month.intValue(), (day == 0) ? 1 : day.intValue());
                    </#if>
            return new PlcDATE(value);
                    <#break>
                    <#case "TIME_OF_DAY">
                    <#if helper.hasFieldsWithNames(case.fields, "hour", "minutes", "seconds", "nanos")>
            LocalTime value = LocalTime.of(hour.intValue(), minutes.intValue(), seconds.intValue(), nanos.intValue());
                    <#elseif helper.hasFieldsWithNames(case.fields, "hour", "minutes", "seconds")>
            LocalTime value = LocalTime.of(hour.intValue(), minutes.intValue(), seconds.intValue());
                    </#if>
            return new PlcTIME_OF_DAY(value);
                    <#break>
                    <#case "DATE_AND_TIME">
                    <#if helper.hasFieldsWithNames(case.fields, "year", "month", "day", "hour", "minutes", "seconds", "nanos")>
            LocalDateTime value = LocalDateTime.of(year.intValue(), (month == 0) ? 1 : month.intValue(), (day == 0) ? 1 : day.intValue(), hour.intValue(), minutes.intValue(), seconds.intValue(), nanos.intValue());
                    <#elseif helper.hasFieldsWithNames(case.fields, "year", "month", "day", "hour", "minutes", "seconds")>
            LocalDateTime value = LocalDateTime.of(year.intValue(), (month == 0) ? 1 : month.intValue(), (day == 0) ? 1 : day.intValue(), hour.intValue(), minutes.intValue(), seconds.intValue());
                    <#elseif helper.hasFieldsWithNames(case.fields, "secondsSinceEpoch")>
            LocalDateTime value = LocalDateTime.ofEpochSecond(secondsSinceEpoch, 0, ZoneOffset.UTC);
                    </#if>
            return new PlcDATE_AND_TIME(value);
                    <#break>
                    <#case "Struct">
            return new PlcStruct(_map);
                    <#break>
                    <#case "List">
            return new PlcList(value);
                    <#break>
                    <#default>
            return new Plc${case.name}(value);
                </#switch>
            </#if>
        }<#sep> else </#sep></#list>
        <#if !defaultCaseOutput>
        return null;
        </#if>
    }

<#if outputFlavor != "passive">
    public static WriteBufferByteBased staticSerialize(PlcValue _value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        return staticSerialize(_value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>, false);
    }

    public static WriteBufferByteBased staticSerialize(PlcValue _value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>, boolean littleEndian) throws ParseException {
        <#assign defaultCaseOutput=false>
        <#list type.switchField.cases as case><#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)},<#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#elseif helper.isComplexTypeReference(type.parserArguments[discriminatorValue?index].type)><#if helper.isEnumTypeReference(type.parserArguments[discriminatorValue?index].type)>${helper.getLanguageTypeNameForTypeReference(type.parserArguments[discriminatorValue?index].type, false)}.${discriminatorValue}<#else>${discriminatorValue}</#if><#else>${discriminatorValue}</#if>)<#sep> && </#sep></#list>) <#else><#assign defaultCaseOutput=true></#if>{ // ${case.name}
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased((int) Math.ceil(((float) ${helper.getSizeInBits(case, type.parserArguments)}) / 8.0f), littleEndian);

            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">
            PlcList values = (PlcList) _value;

                        <#if case.name == "Struct">
            for (PlcValue val : ((List<PlcValue>) values.getStruct().get("${field.name}").getList())) {
                ${helper.getLanguageTypeNameForField(field)} value = (${helper.getLanguageTypeNameForField(field)}) val.get${helper.getLanguageTypeNameForField(field)?cap_first}();
                ${helper.getWriteBufferWriteMethodCall(field.type, "value", field)};
            }
                        <#else>
            for (PlcValue val : ((List<PlcValue>) values.getList())) {
                ${helper.getLanguageTypeNameForField(field)} value = (${helper.getLanguageTypeNameForField(field)}) val.get${helper.getLanguageTypeNameForField(field)?cap_first}();
                ${helper.getWriteBufferWriteMethodCall(field.type, "(" + field.name + ")", field)};
            }
                        </#if>

                    <#if case.name == "BOOL">
            while (writeBuffer.getPos() < writeBuffer.getData().length) {
                writeBuffer.writeBit(false);
            }
                    </#if>
                    <#break>
                    <#case "const">
            // Const Field (${field.name})
            ${helper.getWriteBufferWriteMethodCall(field.type, field.referenceValue, field)};
                    <#break>
                    <#case "enum">
            // Enum field (${field.name})
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${field.name?cap_first}();
            ${helper.getWriteBufferWriteMethodCall(helper.getEnumBaseTypeReference(field.type), "(" + field.name + ".getValue())", field)};
                    <#break>
                    <#case "manual">
            // Manual Field (${field.name})
            ${helper.toSerializationExpression(field, field.serializeExpression, type.parserArguments)};
                    <#break>
                    <#case "reserved">
            // Reserved Field
            ${helper.getWriteBufferWriteMethodCall(field.type, helper.getReservedValue(field), field)};
                    <#break>
                    <#case "simple">
            // Simple Field (${field.name})
                        <#if case.name == "Struct">
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.getStruct().get("${field.name}").get${helper.getLanguageTypeNameForField(field)?cap_first}();
                        <#else>
                            <#if field.name == "value">
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${helper.getLanguageTypeNameForField(field)?cap_first}();
                            <#else>
                                <#-- Just for now -->
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getNullValueForTypeReference(field.type)};
                            </#if>
                        </#if>
                        <#if helper.isSimpleTypeReference(field.type)>
            ${helper.getWriteBufferWriteMethodCall(field.type, "(" + field.name + ")", field)};
                        <#else>
            ${field.type.name}IO.staticSerialize(writeBuffer, ${field.name});
                        </#if>
                    <#break>
                </#switch>
            </#list>
            return writeBuffer;
        }<#sep> else </#sep></#list>
        <#if !defaultCaseOutput>
        return null;
        </#if>
    }
</#if>

}
</#outputformat>
