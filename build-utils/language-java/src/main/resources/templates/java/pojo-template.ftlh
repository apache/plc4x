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
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${typeName}.java
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

package ${helper.packageName(protocolName, languageName, outputFlavor)};

<#if helper.getEnumTypes(type)?has_content>import ${helper.packageName(protocolName, languageName, outputFlavor)}.types.*;</#if>

import com.fasterxml.jackson.annotation.*;
import org.apache.plc4x.java.utils.Message;

<#if type.abstract>@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, include = JsonTypeInfo.As.PROPERTY, property = "className")</#if>
public<#if type.abstract> abstract</#if> class ${typeName}<#if type.parentType??> extends ${type.parentType.name}</#if> implements Message {

<#if helper.isDiscriminatedType(type)>

    // Discriminator values used by the parser to determine the type to be used. All values have to apply.
    public static final Object[] DISCRIMINATOR_VALUES = new Object[] {
    <#list type.discriminatorValues as discriminatorValue>
        <#-- There are rare occasions where the discriminator is defined by a parser argument, in this case we currently can't detect the type -->
        <#if type.discriminatorField??>(${helper.getLanguageTypeNameForField(type.discriminatorField)}) </#if>${discriminatorValue}<#sep>, </#sep>
    </#list>
    };
</#if>
<#if type.constFields?has_content>

    // Constant values.
<#list type.constFields as field>
    public static final ${helper.getLanguageTypeNameForField(field)} ${field.name?upper_case} = ${field.referenceValue};
</#list>
</#if>
<#-- Prpoerty fields are fields that require a property in the pojo -->
<#if type.propertyFields?has_content>

    // Properties.
<#list type.propertyFields as field>
    private final ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name};
</#list>
</#if>

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    <#-- getAllPropertyFields() returns not only the property fields of this type but also of it's parents -->
    public ${typeName}(<#list type.getAllPropertyFields() as field>@JsonProperty("${field.name}") ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
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
    public ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> get${field.name?cap_first}() {
        return ${field.name};
    }

</#list>
<#list type.virtualFields as field>
    public ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> get${field.name?cap_first}() {
        return (${helper.getLanguageTypeNameForField(field)}) (${helper.toDeserializationExpression(field.valueExpression, type.parserArguments)});
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
        if(${field.name} != null) {
            for(Message element : ${field.name}) {
                lengthInBits += element.getLengthInBytes() * 8;
            }
        }
        </#if>
        <#break>
    <#case "checksum">

        // Checksum Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "const">

        // Const Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "discriminator">

        // Discriminator Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "enum">

        // Enum Field (${field.name})
        lengthInBits += ${helper.getEnumBaseType(field.type).size};
        <#break>
    <#case "implicit">

        // Implicit Field (${field.name})
        lengthInBits += ${field.type.size};
        <#break>
    <#case "manualArray">

        // Manual Array Field (${field.name})
        lengthInBits += ${helper.toDeserializationExpression(field.lengthExpression, null)} * 8;
        <#break>
    <#case "manual">

        // Manual Field (${field.name})
        lengthInBits += ${helper.toDeserializationExpression(field.lengthExpression, null)} * 8;
        <#break>
    <#case "optional">

        // Optional Field (${field.name})
        if(${field.name} != null) {
        <#if helper.isSimpleType(field.type)>
            lengthInBits += ${field.type.size};
        <#else>
            lengthInBits += ${field.name}.getLengthInBytes() * 8;
        </#if>
        }
        <#break>
    <#case "padding">

        // Padding Field (${field.name})
        if((boolean) (${helper.toDeserializationExpression(field.paddingCondition, type.parserArguments)})) {
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
    <#case "virtual">

        // A virtual field doesn't have any in- or output.
        <#break>
</#switch>
</#list>

        return lengthInBits / 8;
    }

}