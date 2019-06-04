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
${packageName?replace(".", "/")}/io/${typeName}IO.java
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

package ${packageName}.io;

import ${packageName}.*;
import org.apache.commons.jexl3.*;
import org.apache.plc4x.java.utils.EvaluationHelper;
import org.apache.plc4x.java.utils.ReadBuffer;
import org.apache.plc4x.java.utils.ParseException;
import org.apache.plc4x.java.utils.WriteBuffer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;

public class ${typeName}IO  {

    private static final Logger LOGGER = LoggerFactory.getLogger(${typeName}IO.class);

<#assign needJexl=false>
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">
    private static JexlExpression ex${field.name?cap_first};
    <#assign needJexl=true>
    <#break>
    <#case "optional">
    private static JexlExpression ex${field.name?cap_first};
    <#assign needJexl=true>
    <#break>
    </#switch>
</#list>
<#if needJexl>
    static {
        JexlEngine jexl = new JexlBuilder().create();
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">
        ex${field.name?cap_first} = jexl.createExpression("${field.lengthExpression}");
        <#break>
    <#case "optional">
        ex${field.name?cap_first} = jexl.createExpression("${field.conditionExpression}");
        <#break>
</#switch>
</#list>
    }
</#if>

    public static <#if helper.isDiscriminatedType(type)>${typeName}Builder<#else>${typeName}</#if> parse(ReadBuffer io<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForSpecType(parserArgument.type)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) throws ParseException {
        int startPos = io.getPos();
        JexlContext jc = new MapContext();
        <#if type.parserArguments?has_content>
        <#list type.parserArguments as parserArgument>
        jc.set("${parserArgument.name}", ${parserArgument.name});
        </#list></#if>
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array field
        <#if field.lengthExpression?contains("curPos")>
        jc.set("curPos", io.getPos() - startPos);
        </#if>
        int size = Integer.valueOf(ex${field.name?cap_first}.evaluate(jc).toString());
        <#if helper.isCountArray(field)>
        ${helper.getLanguageTypeNameForSpecType(field.type)}[] ${field.name} = new ${helper.getLanguageTypeNameForSpecType(field.type)}[size];
        for(int i = 0; i < size; i++) {
            ${field.name}[i] = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>${parserArgument}<#sep>, </#sep></#list></#if>)</#if>;
        }
        <#else>
        List<${helper.getLanguageTypeNameForSpecType(field.type)}> ${field.name}List = <#if helper.isCountArray(field)>new ArrayList<>(size)<#else>new LinkedList<>()</#if>;
        int endPos = io.getPos() + size;
        while(io.getPos() < endPos) {
            ${field.name}List.add(<#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>${parserArgument}<#sep>, </#sep></#list></#if>)</#if>);
        }
        ${helper.getLanguageTypeNameForSpecType(field.type)}[] ${field.name} = ${field.name}List.toArray(new ${helper.getLanguageTypeNameForSpecType(field.type)}[0]);
        </#if>
        jc.set("${field.name}", ${field.name});
        <#break>
    <#case "const">

        // Const Field
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        if(${field.name} != ${typeName}.${field.name?upper_case}) {
            throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
        }
        jc.set("${field.name}", ${field.name});
        <#break>
    <#case "discriminator">

        // Discriminator Field (Used as input to a switch field)
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        jc.set("${field.name}", ${field.name});
        <#break>
    <#case "implicit">

        // Implicit Field (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        jc.set("${field.name}", ${field.name});
        <#break>
    <#case "optional">

        // Optional Field (Can be skipped, if a given expression evaluates to false)
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = ${helper.getNullValueForType(field.type)};
        if(ex${field.name?cap_first}.evaluate(jc) == Boolean.TRUE) {
            ${field.name} = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io);</#if>;
        }
        jc.set("${field.name}", ${field.name});
        <#break>
    <#case "reserved">

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        {
            ${helper.getLanguageTypeNameForSpecType(field.type)} reserved = io.${helper.getReadMethodName(field.type)}(${field.type.size});
            if(reserved != ${field.referenceValue}) {
                LOGGER.info("Expected constant value " + ${field.referenceValue} + " but got " + reserved + " for reserved field.");
            }
        }
        <#break>
    <#case "simple">

        // Simple field
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>${parserArgument}<#sep>, </#sep></#list></#if>)</#if>;
        jc.set("${field.name}", ${field.name});
        <#break>
    <#case "switch">

        // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
        ${typeName}Builder builder = null;
        <#list field.cases as case>
        if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toReadExpression(field.discriminatorNames[discriminatorValue?index])}, ${discriminatorValue})<#sep> && </#sep></#list>) {
            builder = ${case.name}IO.parse(io<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
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

    public static void serialize(WriteBuffer io, ${typeName} value) throws ParseException {
        JexlEngine jexl = new JexlBuilder().create();
        JexlContext jc = new MapContext();
        jc.set("this", value);
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array field
        if(value.get${field.name?cap_first}() != null) {
            for(${helper.getLanguageTypeNameForSpecType(field.type)} element : value.get${field.name?cap_first}()) {
                <#if helper.isSimpleType(field.type)>
                io.${helper.getWriteBufferReadMethodCall(field.type, "element")};
                <#else>
                ${helper.getLanguageTypeNameForSpecType(field.type)}IO.serialize(io, element);
                </#if>
            }
        }
        <#break>
    <#case "const">

        // Const Field
        io.${helper.getWriteBufferReadMethodCall(field.type, field.referenceValue)};
        <#break>
    <#case "context">

        // Context Field (Artificial value not retrieved by the input, but the context)
        <#break>
    <#case "discriminator">

        // Discriminator Field (Used as input to a switch field)
        io.${helper.getWriteBufferReadMethodCall(field.type, "value.getDiscriminatorValues()[0]")};
        <#break>
    <#case "implicit">

        // Implicit Field (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        JexlExpression ex${field.name?cap_first} = jexl.createExpression("${field.serializationExpression}");
        io.${helper.getWriteBufferReadMethodCall(field.type, "ex" + field.name?cap_first + ".evaluate(jc)")};
        <#break>
    <#case "optional">

        // Optional Field (Can be skipped, if a given expression evaluates to false)
        JexlExpression ex${field.name?cap_first} = jexl.createExpression("${field.conditionExpression}");
        if(ex${field.name?cap_first}.evaluate(jc) == Boolean.TRUE) {
            <#if helper.isSimpleType(field.type)>
            io.${helper.getWriteBufferReadMethodCall(field.type, "value.get" + field.name?cap_first + "()")};
            <#else>
            ${helper.getLanguageTypeNameForSpecType(field.type)}IO.serialize(io, value.get${field.name?cap_first}());
            </#if>
        }
        <#break>
    <#case "reserved">

        // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
        io.${helper.getWriteBufferReadMethodCall(field.type, field.referenceValue)};
        <#break>
    <#case "simple">

        // Simple field
        <#if helper.isSimpleType(field.type)>
        io.${helper.getWriteBufferReadMethodCall(field.type, "value.get" + field.name?cap_first + "()")};
        <#else>
        ${helper.getLanguageTypeNameForSpecType(field.type)}IO.serialize(io, value.get${field.name?cap_first}());
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

<#if type.abstract>
    public static interface ${typeName}Builder {
        ${typeName} build(<#list type.propertyFields as field>${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>);
    }

</#if>
<#if helper.isDiscriminatedType(type)>
    public static class ${typeName}Builder implements ${type.parentType.name}IO.${type.parentType.name}Builder {
        <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
        private final ${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> ${field.name};
        </#list>
        </#if>

        public ${typeName}Builder(<#list type.propertyFields as field>${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            <#list type.propertyFields as field>
            this.${field.name} = ${field.name};
            </#list>
        }

        public ${typeName} build(<#list type.parentType.propertyFields as field>${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            return new ${typeName}(<#list type.getAllPropertyFields() as field>${field.name}<#sep>, </#sep></#list>);
        }
    }

</#if>
}