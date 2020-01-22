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
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/io/${typeName}IO.java
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

public class ${typeName}IO {

    private static final Logger LOGGER = LoggerFactory.getLogger(${typeName}IO.class);

    public static PlcValue staticParse(ReadBuffer io<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeName(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        <#list type.switchField.cases as case><#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toSwitchExpression(type.switchField.discriminatorNames[discriminatorValue?index])}, ${discriminatorValue})<#sep> && </#sep></#list>) </#if>{ // ${case.name}
            <#list case.fields as field>
                <#assign skipReturn=false>
                <#switch field.typeName>
                    <#case "array">

            // Array field (${field.name})
            <#-- Only update curPos if the length expression uses it -->
                        <#if field.loopExpression.contains("curPos")>
            curPos = io.getPos() - startPos;
                        </#if>
            <#-- If this is a count array, we can directly initialize an array with the given size -->
                        <#if helper.isCountArray(field)>
            // Count array
            if(${helper.toParseExpression(field, field.loopExpression, type.parserArguments)?no_esc} > Integer.MAX_VALUE) {
                throw new ParseException("Array count of " + (${helper.toParseExpression(field, field.loopExpression, type.parserArguments)?no_esc}) + " exceeds the maximum allowed count of " + Integer.MAX_VALUE);
            }
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[] ${field.name};
            {
                int itemCount = (int) ${helper.toParseExpression(field, field.loopExpression, type.parserArguments)?no_esc};
                ${field.name} = new ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[itemCount];
                for(int curItem = 0; curItem < itemCount; curItem++) {
                    ${field.name}[curItem] = <#if helper.isSimpleType(field.type)>${helper.getReadBufferReadMethodCall(field.type)?no_esc}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>;
                }
            }
            <#-- In all other cases do we have to work with a list, that is later converted to an array -->
                        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
                            <#if helper.isLengthArray(field)>
            // Length array
            int _${field.name}Length = ${helper.toParseExpression(field, field.loopExpression, type.parserArguments)?no_esc};
            List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
            int ${field.name}EndPos = io.getPos() + _${field.name}Length;
            while(io.getPos() < ${field.name}EndPos) {
                _${field.name}List.add(<#if helper.isSimpleType(field.type)>${helper.getReadBufferReadMethodCall(field.type)?no_esc}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>);
                <#-- After parsing, update the current position, but only if it's needed -->
                                <#if field.loopExpression.contains("curPos")>
                curPos = io.getPos() - startPos;
                                </#if>
            }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
                            <#elseif helper.isTerminatedArray(field)>
            // Terminated array
            List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
            while(!((boolean) (${helper.toParseExpression(field, field.loopExpression, type.parserArguments)?no_esc}))) {
                _${field.name}List.add(<#if helper.isSimpleType(field.type)>${helper.getReadBufferReadMethodCall(field.type)?no_esc}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>);

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
                            <#if helper.isSimpleType(field.type)>
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
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
            if(${field.name} != ${typeName}.${field.name?upper_case}) {
                throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
            }
                    <#break>
                    <#case "enum">

            // Enum field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = ${helper.getNonPrimitiveLanguageTypeNameForField(field)}.valueOf(${helper.getReadBufferReadMethodCall(helper.getEnumBaseType(field.type))?no_esc});
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
                            <#-- Disable this for now -->
                            <#case "String">
            //return new PlcStruct(${field.name});
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
                ${helper.getLanguageTypeNameForField(field)} reserved = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
                if(reserved != ${helper.getReservedValue(field)}) {
                    LOGGER.info("Expected constant value " + ${field.referenceValue} + " but got " + reserved + " for reserved field.");
                }
            }
                    <#break>
                    <#case "simple">

            // Simple Field (${field.name})
            ${helper.getNonPrimitiveLanguageTypeNameForField(field)} ${field.name} = <#if helper.isSimpleType(field.type)>${helper.getReadBufferReadMethodCall(field.type)?no_esc}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
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
            return new PlcList(Arrays.asList(value));
                    <#break>
            <#-- Disable this for now -->
                    <#case "String">
            //return new PlcStruct(_map);
                    <#break>
                    <#default>
            return new Plc${case.name}(value);
                </#switch>
            </#if>
        }<#sep> else </#sep></#list>
        return null;
    }

<#if outputFlavor != "passive">
    public static WriteBuffer staticSerialize(PlcValue _value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeName(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        <#list type.switchField.cases as case><#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toSwitchExpression(type.switchField.discriminatorNames[discriminatorValue?index])}, ${discriminatorValue})<#sep> && </#sep></#list>) </#if>{ // ${case.name}
            WriteBuffer io = new WriteBuffer(${helper.getSizeInBits(case)});
            <#list case.fields as field>
                <#switch field.typeName>
                    <#case "array">
                    <#break>
                    <#case "const">
            // Const Field (${field.name})
            ${helper.getWriteBufferReadMethodCall(field.type, field.referenceValue)?no_esc};
                    <#break>
                    <#case "enum">
            // Enum field (${field.name})
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${field.name?cap_first}();
            ${helper.getWriteBufferReadMethodCall(helper.getEnumBaseType(field.type), "(" + field.name + ".getValue())")?no_esc};
                    <#break>
                    <#case "manual">
            // Manual Field (${field.name})
            ${helper.toSerializationExpression(field, field.serializeExpression, type.parserArguments)};
                    <#break>
                    <#case "reserved">
            // Reserved Field
            ${helper.getWriteBufferReadMethodCall(field.type, helper.getReservedValue(field))?no_esc};
                    <#break>
                    <#case "simple">
            // Simple Field (${field.name})
                        <#if case.name == "Struct">
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.getStruct().get("${field.name}").get${helper.getLanguageTypeNameForField(field)?cap_first}();
                        <#else>
                            <#if field.name == "value">
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) _value.get${case.name}();
                            <#else>
                                <#-- Just for now -->
            ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getNullValueForType(field.type)};
                            </#if>
                        </#if>
                        <#if helper.isSimpleType(field.type)>
            ${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")?no_esc};
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