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
<#if helper.getComplexTypes(type)?has_content>import ${packageName}.io.*;</#if>
import org.apache.plc4x.java.utils.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;

public class ${typeName}IO implements MessageIO<${typeName}<#if helper.isDiscriminatedType(type)>IO.${typeName}Builder</#if>, ${typeName}> {

    private static final Logger LOGGER = LoggerFactory.getLogger(${typeName}IO.class);

<#if helper.getComplexTypes(type)?has_content>
    // IO Helpers.
<#list helper.getComplexTypes(type) as complexType>
    private final ${complexType.name}IO ${complexType.name?uncap_first}IO;
</#list>

    public ${typeName}IO() {
<#list helper.getComplexTypes(type) as complexType>
        ${complexType.name?uncap_first}IO = new ${complexType.name}IO();
</#list>
   }

</#if>
    public ${typeName}<#if helper.isDiscriminatedType(type)>Builder</#if> parse(ReadBuffer io, Object... args) throws ParseException {
<#if type.parserArguments?has_content>
<#list type.parserArguments as parserArgument>
        ${helper.getLanguageTypeNameForSpecType(parserArgument.type)} ${parserArgument.name} = (${helper.getLanguageTypeNameForSpecType(parserArgument.type)}) args[${parserArgument?index}];
</#list>
</#if>
        int startPos = io.getPos();
        int curPos;
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array field
        <#if field.lengthExpression.contains("curPos")>
        curPos = io.getPos() - startPos;
        </#if>
        int ${field.name}Size = ${helper.toDeserializationExpression(field.lengthExpression)};
        <#if helper.isCountArray(field)>
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getLanguageTypeNameForField(field)}[${field.name}Size];
        for(int i = 0; i < ${field.name}Size; i++) {
            ${field.name}[i] = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name?uncap_first}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>${parserArgument}<#sep>, </#sep></#list></#if>)</#if>;
        }
        <#else>
        List<${helper.getNonPrimitiveLanguageTypeNameForField(field)}> ${field.name}List = <#if helper.isCountArray(field)>new ArrayList<>(size)<#else>new LinkedList<>()</#if>;
        int ${field.name}EndPos = io.getPos() + ${field.name}Size;
        while(io.getPos() < ${field.name}EndPos) {
            <#if field.lengthExpression.contains("curPos")>
            curPos = io.getPos() - startPos;
            </#if>
            ${field.name}List.add(<#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name?uncap_first}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument)})<#sep>, </#sep></#list></#if>)</#if>);
        }
        <#if helper.isSimpleType(field.type)>
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = new ${helper.getLanguageTypeNameForField(field)}[${field.name}List.size()];
        for(int i = 0; i < ${field.name}List.size(); i++) {
            ${field.name}[i] = (${helper.getLanguageTypeNameForField(field)}) ${field.name}List.get(i);
        }
        <#else>
        ${helper.getLanguageTypeNameForField(field)}[] ${field.name} = ${field.name}List.toArray(new ${helper.getNonPrimitiveLanguageTypeNameForField(field)}[0]);
        </#if>
        </#if>
        <#break>
    <#case "const">

        // Const Field
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        if(${field.name} != ${typeName}.${field.name?upper_case}) {
            throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
        }
        <#break>
    <#case "discriminator">

        // Discriminator Field (Used as input to a switch field)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
    <#break>
    <#case "implicit">

        // Implicit Field (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        <#break>
    <#case "optional">

        // Optional Field (Can be skipped, if a given expression evaluates to false)
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getNullValueForType(field.type)};
        if(${helper.toDeserializationExpression(field.conditionExpression)}) {
            ${field.name} = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name?uncap_first}IO.parse(io);</#if>;
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

        // Simple field
        ${helper.getLanguageTypeNameForField(field)} ${field.name} = <#if helper.isSimpleType(field.type)>io.${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name?uncap_first}IO.parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toDeserializationExpression(parserArgument)})<#sep>, </#sep></#list></#if>)</#if>;
        <#break>
    <#case "switch">

        // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
        ${typeName}Builder builder = null;
        <#list field.cases as case>
        if(<#list case.discriminatorValues as discriminatorValue>EvaluationHelper.equals(${helper.toSwitchExpression(field.discriminatorNames[discriminatorValue?index])}, ${discriminatorValue})<#sep> && </#sep></#list>) {
            builder = ${case.name?uncap_first}IO.parse(io<#if case.parserArguments?has_content>, <#list case.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
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

    public void serialize(WriteBuffer io, ${typeName} value, Object... args) throws ParseException {
<#if type.parserArguments?has_content>
<#list type.parserArguments as parserArgument>
        ${helper.getLanguageTypeNameForSpecType(parserArgument.type)} ${parserArgument.name} = (${helper.getLanguageTypeNameForSpecType(parserArgument.type)}) args[${parserArgument?index}];
</#list>
</#if>
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array field
        if(value.get${field.name?cap_first}() != null) {
            for(${helper.getLanguageTypeNameForField(field)} element : value.get${field.name?cap_first}()) {
                <#if helper.isSimpleType(field.type)>
                io.${helper.getWriteBufferReadMethodCall(field.type, "element")};
                <#else>
                ${field.type.name?uncap_first}IO.serialize(io, element);
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
        io.${helper.getWriteBufferReadMethodCall(field.type, "(" + helper.toSerializationExpression(field.serializationExpression) + ")")?no_esc};
        <#break>
    <#case "optional">

        // Optional Field (Can be skipped, if the value is null)
        if(value.get${field.name?cap_first}() != null) {
            <#if helper.isSimpleType(field.type)>
            io.${helper.getWriteBufferReadMethodCall(field.type, "value.get" + field.name?cap_first + "()")};
            <#else>
            ${field.type.name?uncap_first}IO.serialize(io, value.get${field.name?cap_first}());
            </#if>
        }
        <#break>
    <#case "reserved">

        // Reserved Field.
        io.${helper.getWriteBufferReadMethodCall(field.type, field.referenceValue)};
        <#break>
    <#case "simple">

        // Simple field
        <#if helper.isSimpleType(field.type)>
        io.${helper.getWriteBufferReadMethodCall(field.type, "value.get" + field.name?cap_first + "()")};
        <#else>
        ${field.type.name?uncap_first}IO.serialize(io, value.get${field.name?cap_first}());
        </#if>
        <#break>
    <#case "switch">

        // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
        <#list field.cases as case>
        if(value instanceof ${case.name}) {
            ${case.name?uncap_first}IO.serialize(io, (${case.name}) value);
        }<#sep> else </#sep>
        </#list>
        <#break>
</#switch>
</#list>
    }
    public static int CRC(Object... args) {
        return 0;
    }

    private static Map<String, Method> staticMethods;

    /**
     * Notice, this method doesn't support using multiple static methods with the same name
     * but with different argument types. In this case it will use the first used method and
     * probably fail for all others.
     *
     * @param fullyQualifiedMethodName Fully qualified Class name with method name (Example "a.b.c.Foo.bar")
     * @param args arguments passed to the static method.
     * @return
     */
    private static Object STATIC_CALL(String fullyQualifiedMethodName, Object... args) {
        if(staticMethods == null) {
            staticMethods = new HashMap<>();
        }
        // Get the method handle and cache it for later usage.
        if(!staticMethods.containsKey(fullyQualifiedMethodName)) {
            // Get the types of the arguments.
            Class<?>[] argTypes = new Class[args.length];
            for(int i = 0; i < args.length; i++) {
                argTypes[i] = args[i].getClass();
            }
            // Get the method handle.
            String className = fullyQualifiedMethodName.substring(0, fullyQualifiedMethodName.lastIndexOf("."));
            String methodName = fullyQualifiedMethodName.substring(className.length());
            try {
                Class clazz = Class.forName(className);
                Method method = clazz.getMethod(methodName, argTypes);
                staticMethods.put(fullyQualifiedMethodName, method);
            } catch (ClassNotFoundException | NoSuchMethodException e) {
                throw new RuntimeException("Error finding static method", e);
            }
        }
        Method method = staticMethods.get(fullyQualifiedMethodName);
        try {
            // Invoke the static method.
            return method.invoke(null, args);
        } catch (InvocationTargetException | IllegalAccessException e) {
            throw new RuntimeException("Error invoking static method", e);
        }
    }

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
        ${typeName} build(<#list type.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>);
    }

</#if>
<#if helper.isDiscriminatedType(type)>
    public static class ${typeName}Builder implements ${type.parentType.name}IO.${type.parentType.name}Builder {
        <#if type.propertyFields?has_content>
        <#list type.propertyFields as field>
        private final ${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> ${field.name};
        </#list>
        </#if>

        public ${typeName}Builder(<#list type.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            <#list type.propertyFields as field>
            this.${field.name} = ${field.name};
            </#list>
        }

        public ${typeName} build(<#list type.parentType.propertyFields as field>${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
            return new ${typeName}(<#list type.getAllPropertyFields() as field>${field.name}<#sep>, </#sep></#list>);
        }
    }

</#if>
}