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

import ${helper.packageName(protocolName, languageName, outputFlavor)}.*;
<#if helper.getComplexTypes(type)?has_content>import ${helper.packageName(protocolName, languageName, outputFlavor)}.io.*;</#if>
<#if helper.getEnumTypes(type)?has_content>import ${helper.packageName(protocolName, languageName, outputFlavor)}.types.*;</#if>
import org.apache.plc4x.java.utils.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;

public class ${typeName}IO {

    private static final Logger LOGGER = LoggerFactory.getLogger(${typeName}IO.class);

    public static ${typeName}<#if helper.isDiscriminatedType(type)>Builder</#if> parse(ReadBuffer io<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeName(parserArgument.type, false)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        int startPos = io.getPos();
        int curPos;
<#list type.fields as field>
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
        int _${field.name}Count = ${helper.toDeserializationExpression(field.loopExpression, type.parserArguments)?no_esc};
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getLanguageTypeNameForField(field)}[_${field.name}Count];
        for(int i = 0; i < _${field.name}Count; i++) {
            ${field.name}[i] = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>;
        }
        <#-- In all other cases do we have to work with a list, that is later converted to an array -->
        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
            <#if helper.isLengthArray(field)>
        // Length array
        int _${field.name}Length = ${helper.toDeserializationExpression(field.loopExpression, type.parserArguments)?no_esc};
        List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
        int ${field.name}EndPos = io.getPos() + _${field.name}Length;
        while(io.getPos() < ${field.name}EndPos) {
            _${field.name}List.add(<#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>);
            <#-- After parsing, update the current position, but only if it's needed -->
            <#if field.loopExpression.contains("curPos")>
            curPos = io.getPos() - startPos;
            </#if>
        }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
            <#elseif helper.isTerminatedArray(field)>
        // Terminated array
        List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
        while(!((boolean) (${helper.toDeserializationExpression(field.loopExpression, type.parserArguments)?no_esc}))) {
            _${field.name}List.add(<#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>);

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
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getLanguageTypeNameForField(field)}[_${field.name}List.size()];
        for(int i = 0; i < _${field.name}List.size(); i++) {
            ${field.name}[i] = (${helper.getLanguageTypeNameForField(field)}) _${field.name}List.get(i);
        }
            <#else>
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = _${field.name}List.toArray(new ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[0]);
            </#if>
        </#if>
        <#break>
    <#case "checksum">

        // Checksum Field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        ${helper.getLanguageTypeNameForField(field)} _${field.name}Ref = (${helper.getLanguageTypeNameForField(field)}) (${helper.toDeserializationExpression(field.checksumExpression, type.parserArguments)});
        if(${field.name} != _${field.name}Ref) {
            throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X",_${field.name}Ref & 0xFFFF, ${field.name} & 0xFFFF));
        }
        <#break>
    <#case "const">

        // Const Field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        if(${field.name} != ${typeName}.${field.name?upper_case}) {
            throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
        }
        <#break>
    <#case "enum">

        // Enum field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getLanguageTypeNameForField(field)}.valueOf(io.${helper.getReadMethodName(helper.getEnumBaseType(field.type))}(${helper.getEnumBaseType(field.type).size}));
        <#break>
    <#case "discriminator">

        // Discriminator Field (${field.name}) (Used as input to a switch field)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        <#break>
    <#case "implicit">

        // Implicit Field (${field.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        <#break>
    <#case "manualArray">

        // Manual Array Field (${field.name})
        <#-- Only update curPos if the length expression uses it -->
        <#if field.loopExpression.contains("curPos")>
        curPos = io.getPos() - startPos;
        </#if>
        <#-- If this is a count array, we can directly initialize an array with the given size -->
        <#if helper.isCountArray(field)>
        // Count array
        int _${field.name}Count = ${helper.toDeserializationExpression(field.loopExpression, type.parserArguments)};
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getLanguageTypeNameForField(field)}[_${field.name}Count];
        for(int i = 0; i < _${field.name}Count; i++) {
            ${field.name}[i] = (${helper.getLanguageTypeNameForField(field)}) (${helper.toDeserializationExpression(field.deserializationExpression, type.parserArguments)});
        }
        <#-- In all other cases do we have to work with a list, that is later converted to an array -->
        <#else>
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
            <#if helper.isLengthArray(field)>
        // Length array
        int _${field.name}Length = ${helper.toDeserializationExpression(field.loopExpression, type.parserArguments)};
        List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
        int ${field.name}EndPos = io.getPos() + _${field.name}Length;
        while(io.getPos() < ${field.name}EndPos) {
            _${field.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toDeserializationExpression(field.deserializationExpression, type.parserArguments)}));
            <#-- After parsing, update the current position, but only if it's needed -->
            <#if field.loopExpression.contains("curPos")>
            curPos = io.getPos() - startPos;
            </#if>
        }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
            <#elseif helper.isTerminatedArray(field)>
        // Terminated array
        List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> _${field.name}List = new LinkedList<>();
        while(!((boolean) (${helper.toDeserializationExpression(field.loopExpression, type.parserArguments)}))) {
            _${field.name}List.add((${helper.getLanguageTypeNameForField(field)}) (${helper.toDeserializationExpression(field.deserializationExpression, type.parserArguments)}));

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
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getLanguageTypeNameForField(field)}[_${field.name}List.size()];
        for(int i = 0; i < _${field.name}List.size(); i++) {
            ${field.name}[i] = (${helper.getLanguageTypeNameForField(field)}) _${field.name}List.get(i);
        }
            <#else>
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = _${field.name}List.toArray(new ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[0]);
            </#if>
        </#if>
        <#break>
    <#case "manual">

        // Manual Field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toDeserializationExpression(field.deserializationExpression, type.parserArguments)});
        <#break>
    <#case "optional">

        // Optional Field (${field.name}) (Can be skipped, if a given expression evaluates to false)
        <#if field.conditionExpression.contains("curPos")>
        curPos = io.getPos() - startPos;
        </#if>
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = null;
        if(${helper.toDeserializationExpression(field.conditionExpression, type.parserArguments)?no_esc}) {
            ${field.name} = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
        }
        <#break>
    <#case "padding">

        // Padding Field (${field.name})
        boolean _${field.name}NeedsPadding = (boolean) (${helper.toDeserializationExpression(field.paddingCondition, type.parserArguments)});
        if(_${field.name}NeedsPadding) {
            // Just read the padding data and ignore it
            io.${helper.getReadBufferReadMethodCall(field.type)};
        }
        <#break>
    <#case "reserved">

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        {
            ${helper.getLanguageTypeNameForField(field)} reserved = io.${helper.getReadBufferReadMethodCall(field.type)};
            if(reserved != ${field.referenceValue}) {
                LOGGER.info("Expected constant value " + ${field.referenceValue} + " but got " + reserved + " for reserved field.");
            }
        }
        <#break>
    <#case "simple">

        // Simple Field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
        <#break>
    <#case "switch">

        // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
        ${typeName}Builder builder = null;
        <#list field.cases as case>
        <#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toSwitchExpression(field.discriminatorNames[discriminatorValue?index])}, ${discriminatorValue})<#sep> && </#sep></#list>) </#if>{
            builder = ${case.name}IO.parse(io<#if case.parserArguments?has_content>, <#list case.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
        }<#sep> else </#sep>
        </#list>
        if (builder == null) {
            throw new ParseException("Unsupported case for discriminated type");
        }
        <#break>
</#switch>
</#list>

        // Create the instance
        <#if helper.isDiscriminatedType(type)>
        return new ${typeName}Builder(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        <#elseif type.abstract>
        return builder.build(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        <#else>
        return new ${typeName}(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        </#if>
    }

<#if outputFlavor != "passive">
    public static void serialize(WriteBuffer io, ${typeName} value) throws ParseException {
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array Field (${field.name})
        if(value.get${field.name?cap_first}() != null) {
            for(${helper.getLanguageTypeNameForField(field)} element : value.get${field.name?cap_first}()) {
                <#if helper.isSimpleType(field.type)>
                io.${helper.getWriteBufferReadMethodCall(field.type, "element")};
                <#else>
                ${field.type.name}IO.serialize(io, element);
                </#if>
            }
        }
        <#break>
    <#case "checksum">

        // Checksum Field (${field.name}) (Calculated)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(field.checksumExpression, type.parserArguments)});
        io.${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")?no_esc};
        <#break>
    <#case "const">

        // Const Field (${field.name})
        io.${helper.getWriteBufferReadMethodCall(field.type, field.referenceValue)};
        <#break>
    <#case "discriminator">

        // Discriminator Field (${field.name}) (Used as input to a switch field)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) value.getDiscriminatorValues()[0];
        io.${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")};
        <#break>
    <#case "enum">

        // Enum field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) value.get${field.name?cap_first}();
        io.${helper.getWriteBufferReadMethodCall(helper.getEnumBaseType(field.type), "(" + field.name + ".getValue())")};
        <#break>
    <#case "implicit">

        // Implicit Field (${field.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(field.serializationExpression, type.parserArguments)});
        io.${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")?no_esc};
        <#break>
    <#case "manualArray">

        // Manual Array Field (${field.name})
        if(value.get${field.name?cap_first}() != null) {
            for(${helper.getLanguageTypeNameForField(field)} element : value.get${field.name?cap_first}()) {
                ${helper.toSerializationExpression(field.serializationExpression, type.parserArguments)};
            }
        }
        <#break>
    <#case "manual">

        // Manual Field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(field.serializationExpression, type.parserArguments)});
         io.${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")?no_esc};
       <#break>
    <#case "optional">

        // Optional Field (${field.name}) (Can be skipped, if the value is null)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = null;
        if(value.get${field.name?cap_first}() != null) {
            ${field.name} = (${helper.getLanguageTypeNameForField(field)}) value.get${field.name?cap_first}();
            <#if helper.isSimpleType(field.type)>
            io.${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")};
            <#else>
            ${field.type.name}IO.serialize(io, ${field.name});
            </#if>
        }
        <#break>
    <#case "padding">

        // Padding Field (${field.name})
        boolean _${field.name}NeedsPadding = (boolean) (${helper.toSerializationExpression(field.paddingCondition, type.parserArguments)});
        if(_${field.name}NeedsPadding) {
            ${helper.getLanguageTypeNameForField(field)} _${field.name}PaddingValue = (${helper.getLanguageTypeNameForField(field)}) (${helper.toSerializationExpression(field.paddingValue, type.parserArguments)});
            io.${helper.getWriteBufferReadMethodCall(field.type, "(_" + field.name + "PaddingValue)")?no_esc};
        }
        <#break>
    <#case "reserved">

        // Reserved Field
        io.${helper.getWriteBufferReadMethodCall(field.type, field.referenceValue)};
        <#break>
    <#case "simple">

        // Simple Field (${field.name})
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) value.get${field.name?cap_first}();
        <#if helper.isSimpleType(field.type)>
        io.${helper.getWriteBufferReadMethodCall(field.type, "(" + field.name + ")")};
        <#else>
        ${field.type.name}IO.serialize(io, ${field.name});
        </#if>
        <#break>
    <#case "switch">

        // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
        <#list field.cases as case>
        if(value instanceof ${case.name}) {
            ${case.name}IO.serialize(io, (${case.name}) value);
        }<#sep> else </#sep>
        </#list>
        <#break>
</#switch>
</#list>
    }
</#if>

    private static int COUNT(Object obj) {
        if(obj.getClass().isArray()) {
            Object[] arr = (Object[]) obj;
            return arr.length;
        } else if(obj instanceof Collection) {
            Collection col = (Collection) obj;
            return col.size();
        }
        throw new RuntimeException("Unable to count object of type " + obj.getClass().getName());
    }

    private static <T> T CAST(Object obj, Class<T> clazz) {
        try {
            return clazz.cast(obj);
        } catch(ClassCastException e) {
            throw new RuntimeException("Unable to cast object of type " + obj.getClass().getName() + " to " + clazz.getName());
        }
    }

<#if type.abstract>
    public static interface ${typeName}Builder {
        ${typeName} build(<#list type.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>);
    }

</#if>
<#if helper.isDiscriminatedType(type)>
    public static class ${typeName}Builder implements ${type.parentType.name}IO.${type.parentType.name}Builder {
        <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
        private final ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name};
        </#list>
        </#if>

        public ${typeName}Builder(<#list type.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            <#list type.propertyFields as field>
            this.${field.name} = ${field.name};
            </#list>
        }

        public ${typeName} build(<#list type.parentType.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            return new ${typeName}(<#list type.getAllPropertyFields() as field>${field.name}<#sep>, </#sep></#list>);
        }
    }

</#if>
}