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
import org.apache.plc4x.java.utils.IoBuffer;
import org.apache.plc4x.java.utils.ParseException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ${typeName}IO  {

    private static final Logger LOGGER = LoggerFactory.getLogger(${typeName}IO.class);

    public static <#if helper.isDiscriminatedType(type)>${typeName}Builder<#else>${typeName}</#if> parse(IoBuffer io) throws ParseException {
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">
        array
        <#break>
    <#case "const">

        // Const Field
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        if(${field.name} != ${typeName}.${field.name?upper_case}) {
            throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
        }
        <#break>
    <#case "discriminator">

        // Discriminator Field (Used as input to a switch field)
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        <#break>
    <#case "implicit">

        // Implicit Field (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
        ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name} = io.${helper.getReadMethodName(field.type)}(${field.type.size});
        <#break>
    <#case "optional">

        // Optional Field (Can be skipped, if a given expression evaluates to false)
        if(${field.conditionExpression}) {
           simple
        }
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
        simple
        <#break>
    <#case "switch">

        // Switch field (Depending on the discriminator values, passes the instantiation to a sub-type)
        ${typeName}Builder builder = null;
        <#list field.cases as case>
        if(<#list case.discriminatorValues as discriminatorValue>(${field.discriminatorNames[discriminatorValue?index]} == ${discriminatorValue})<#sep> && </#sep></#list>) {
            builder = ${case.name}IO.parse(io);
        }
        </#list>
        else {
            throw new ParseException("Unsupported case for discriminated type");
        }
        <#break>
</#switch>
</#list>
        <#if helper.isDiscriminatedType(type)>
        return new ${typeName}Builder(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        <#elseif type.abstract>
        return builder.build(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        <#else>
        return new ${typeName}(<#list type.propertyFields as field>${field.name}<#sep>, </#sep></#list>);
        </#if>
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