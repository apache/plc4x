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
${packageName?replace(".", "/")}/${typeName}.java
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

package ${packageName};

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.utils.SizeAware;

<#if type.abstract>@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, include = JsonTypeInfo.As.PROPERTY, property = "className")</#if>
public<#if type.abstract> abstract</#if> class ${typeName}<#if type.parentType??> extends ${type.parentType.name}</#if> implements SizeAware {

<#if helper.isDiscriminatedType(type)>

    // Discriminator values used by the parser to determine the type to be used. All values have to apply.
    public static final Object[] DISCRIMINATOR_VALUES = new Object[] {
    <#list type.discriminatorValues as discriminatorValue>
        ${discriminatorValue}<#sep>, </#sep>
    </#list>
    };
</#if>
<#if type.constFields?has_content>

    // Constant values.
<#list type.constFields as field>
    public static final ${helper.getLanguageTypeNameForField(field)} ${field.name?upper_case} = ${field.referenceValue};
</#list>
</#if>
<#if type.propertyFields?has_content>

    // Properties.
<#list type.propertyFields as field>
    private final ${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> ${field.name};
</#list>
</#if>

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public ${typeName}(<#list type.getAllPropertyFields() as field>@JsonProperty("${field.name}") ${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
<#if type.getParentPropertyFields()?has_content>
        super(<#list type.getParentPropertyFields() as field>${field.name}<#sep>, </#sep></#list>);
</#if>
<#list type.propertyFields as field>
        this.${field.name} = ${field.name};
</#list>
    }
<#if type.abstract>

    public abstract Object[] getDiscriminatorValues();
</#if>
<#if helper.isDiscriminatedType(type)>

    @JsonIgnore
    public Object[] getDiscriminatorValues() {
        return DISCRIMINATOR_VALUES;
    }
</#if>

<#list type.propertyFields as field>
    public ${helper.getLanguageTypeNameForField(field)}<#if field.lengthType??>[]</#if> get${field.name?cap_first}() {
        return ${field.name};
    }

</#list>
    @Override
    @JsonIgnore
    public int getLengthInBytes() {
        int lengthInBits = <#if type.parentType??>super.getLengthInBytes() * 8<#else>0</#if>;
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array field
        <#if helper.isSimpleType(field.type)>
        lengthInBits += ${field.type.size} * ${field.name}.length;
        <#else>
        for(SizeAware element : ${field.name}) {
            lengthInBits += element.getLengthInBytes() * 8;
        }
        </#if>
        <#break>
    <#case "const">

        // Const Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "discriminator">

        // Discriminator Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "implicit">

        // Implicit Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "optional">

        // Optional Field (${field.name})
        if(${field.name} != null) {
            lengthInBits += ${field.type.size};
        }
        <#break>
    <#case "reserved">

        // Reserved Field
        lengthInBits += ${field.type.size};
        <#break>
    <#case "simple">

        // Simple field (${field.name})
        <#if helper.isSimpleType(field.type)>
        lengthInBits += ${field.type.size};
        <#else>
        lengthInBits += ${field.name}.getLengthInBytes() * 8;
        </#if>
        <#break>
    <#case "switch">

        // Length of sub-type elements will be added by sub-type...
        <#break>
</#switch>
</#list>

        return lengthInBits / 8;
    }

}