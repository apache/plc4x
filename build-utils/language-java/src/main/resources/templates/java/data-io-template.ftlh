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
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import ${helper.packageName(protocolName, languageName, outputFlavor)}.*;
import ${helper.packageName(protocolName, languageName, outputFlavor)}.types.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.math.BigInteger;
import java.time.*;
import java.util.*;
import java.util.function.Supplier;

public class ${type.name}IO {

    private static final Logger LOGGER = LoggerFactory.getLogger(${type.name}IO.class);

    public static PlcValue staticParse(ReadBuffer io<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        <#list type.switchField.cases as case><#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)}, ${discriminatorValue})<#sep> && </#sep></#list>) </#if>{ // ${case.name}
            <#assign skipReturn=false>
            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">

            // Array field (${field.name})
            <#-- Only update curPos if the length expression uses it -->
                        <#if field.loopExpression.contains("curPos")>
            curPos = io.getPos() - startPos;
                        </#if>
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                        <#if helper.isCountArrayField(field)>
            // Count array
            if(${helper.toParseExpression(field, field.loopExpression, type.parserArguments)} > Integer.MAX_VALUE) {
                throw new ParseException("Array count of " + (${helper.toParseExpression(field, field.loopExpression, type.parserArguments)}) + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
            }
            List<Plc${case.name}> ${field.name};
            {
                int itemCount = (int) ${helper.toParseExpression(field, field.loopExpression, type.parserArguments)};
                ${field.name} = new LinkedList<>();
                for(int curItem = 0; curItem < itemCount; curItem++) {
                    ${field.name}.add(new Plc${case.name}((${helper.getNonPrimitiveLanguageTypeNameForField(field)}) <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)})<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);
                }
            }
            <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                            <#if helper.isLengthArrayField(field)>
            // Length array
            int _${field.name}Length = ${helper.toParseExpression(field, field.loopExpression, type.parserArguments)};
            List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
            int ${field.name}EndPos = io.getPos() + _${field.name}Length;
            while(io.getPos() < ${field.name}EndPos) {
                _${field.name}List.add(<#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);
                <#-- After parsing, update the current position, but only if it's needed -->
                                <#if field.loopExpression.contains("curPos")>
                curPos = io.getPos() - startPos;
                                </#if>
            }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                            <#elseif helper.isTerminatedArrayField(field)>
            // Terminated array
            List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
            while(!((boolean) (${helper.toParseExpression(field, field.loopExpression, type.parserArguments)}))) {
                _${field.name}List.add(<#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>);

                <#-- After parsing, update the current position, but only if it's needed -->
                                <#if field.loopExpression.contains("curPos")>
                curPos = io.getPos() - startPos;
                                </#if>
            }
                            </#if>
            <#--
                Convert the list into an array. However if the array is of a primitive
                type we have to iterate over it's elements and explicitly cast them.
                Otherwise a simple toArray call is fine.
            -->
                            <#if helper.isSimpleTypeReference(field.type)>
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[_${field.name}List.size()];
            for(int i = 0; i < _${field.name}List.size(); i++) {
                ${field.name}[i] = (${helper.getNonPrimitiveLanguageTypeNameForField(field)}) _${field.name}List.get(i);
            }
                            <#else>
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[] ${field.name} = _${field.name}List.toArray(new ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[0]);
                            </#if>
                        </#if>
                    <#break>
                    <#case "const">

            // Const Field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)};
            if(${field.name} != ${type.name}.${field.name?upper_case}) {
                throw new ParseException("Expected constant value " + ${type.name}.${field.name?upper_case} + " but got " + ${field.name});
            }
                    <#break>
                    <#case "enum">

            // Enum field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getNonPrimitiveLanguageTypeNameForField(field)}.valueOf(${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(field.type))});
                    <#break>
                    <#case "manual">

            // Manual Field (${field.name})
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field, field.parseExpression, type.parserArguments)});
                        <#-- If a manual field is detected, the value returned by this is instantly treated as the value -->
                        <#switch case.name>
                            <#case "Time">
            return new PlcTime(${field.name});
                            <#break>
                            <#case "Date">
            return new PlcDate(${field.name});
                            <#break>
                            <#case "DateTime">
            return new PlcDateTime(${field.name});
                            <#break>
                            <#case "Struct">
            return new PlcStruct(${field.name});
                            <#break>
                            <#case "String">
            return new PlcString(${field.name});
                            <#break>
                            <#default>
            return new Plc${case.name}(${field.name});
                        </#switch>
                        <#-- As we returned early, skip outputting a return at the end -->
                        <#assign skipReturn=true>
                    <#break>
                    <#case "reserved">

            // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
            {
                ${helper.getLanguageTypeNameForField(field)} reserved = ${helper.getReadBufferReadMethodCall(field.type)};
                if(reserved != ${helper.getReservedValue(field)}) {
                    LOGGER.info("Expected constant value " + ${field.referenceValue} + " but got " + reserved + " for reserved field.");
                }
            }
                    <#break>
                    <#case "simple">

            // Simple Field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
                    <#break>
                </#switch>
            </#list>
            <#if case.name == "Struct">
            Map<String, PlcValue> _map = new HashMap<>();
                <#list case.fields as field>
            _map.put("${field.name}", PlcValues.of(${field.name}));
                </#list>
            </#if>
            <#if !skipReturn>
                <#switch case.name>
                    <#case "Time">
            return new PlcTime(LocalTime.of((int) hours, (int) minutes, (int) seconds));
                    <#break>
                    <#case "Date">
            return new PlcDate(LocalDate.of((int) year, (int) month, (int) day));
                    <#break>
                    <#case "DateTime">
            return new PlcDateTime(LocalDateTime.of((int) year, (int) month, (int) day, (int) hours, (int) minutes, (int) seconds));
                    <#break>
                    <#case "Struct">
            return new PlcStruct(_map);
                    <#break>
                    <#case "List">
            return new PlcList(value);
                    <#break>
            <#-- Disable this for now as Strings will only be parsed as manual fields -->
                    <#case "String">
            //return new PlcString(_map);
                    <#break>
                    <#default>
                    <#if helper.isCountArrayField(case.fields[0])>
            return new PlcList(value);
                    <#else>
            return new Plc${case.name}(value);
                    </#if>
                </#switch>
            </#if>
        }<#sep> else </#sep></#list>
        return null;
    }

<#if outputFlavor != "passive">
    public static WriteBuffer staticSerialize(PlcValue _value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        return staticSerialize(_value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>, false);
    }

    public static WriteBuffer staticSerialize(PlcValue _value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>, boolean littleEndian) throws ParseException {
        <#list type.switchField.cases as case><#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)}, ${discriminatorValue})<#sep> && </#sep></#list>) </#if>{ // ${case.name}
            WriteBuffer io = new WriteBuffer((int) Math.ceil(((float) ${helper.getSizeInBits(case, type.parserArguments)}) / 8.0f), littleEndian);

            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">
            PlcList values = (PlcList) _value;

            for (Plc${case.name} val : ((List<Plc${case.name}>) values.getList())) {
                ${helper.getLanguageTypeNameForField(field)} value = (${helper.getLanguageTypeNameForField(field)}) val.get${helper.getLanguageTypeNameForField(field)?cap_first}();
                ${helper.getWriteBufferWriteMethodCall(field.type, "(" + field.name + ")")};
            }
                    <#break>
                    <#case "const">
            // Const Field (${field.name})
            ${helper.getWriteBufferWriteMethodCall(field.type, field.referenceValue)};
                    <#break>
                    <#case "enum">
            // Enum field (${field.name})
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${field.name?cap_first}();
            ${helper.getWriteBufferWriteMethodCall(helper.getEnumBaseTypeReference(field.type), "(" + field.name + ".getValue())")};
                    <#break>
                    <#case "manual">
            // Manual Field (${field.name})
            ${helper.toSerializationExpression(field, field.serializeExpression, type.parserArguments)};
                    <#break>
                    <#case "reserved">
            // Reserved Field
            ${helper.getWriteBufferWriteMethodCall(field.type, helper.getReservedValue(field))};
                    <#break>
                    <#case "simple">
            // Simple Field (${field.name})
                        <#if case.name == "Struct">
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.getStruct().get("${field.name}").get${helper.getLanguageTypeNameForField(field)?cap_first}();
                        <#else>
                            <#if field.name == "value">
            Plc${case.name} plcVal = (Plc${case.name}) _value;
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) plcVal.get${case.name}();
                            <#else>
                                <#-- Just for now -->
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getNullValueForTypeReference(field.type)};
                            </#if>
                        </#if>
                        <#if helper.isSimpleTypeReference(field.type)>
            ${helper.getWriteBufferWriteMethodCall(field.type, "(" + field.name + ")")};
                        <#else>
            ${field.type.name}IO.staticSerialize(io, ${field.name});
                        </#if>
                    <#break>
                </#switch>
            </#list>
            return io;
        }<#sep> else </#sep></#list>
        return null;
    }
</#if>

}
</#outputformat>
