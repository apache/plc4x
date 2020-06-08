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

import static org.apache.plc4x.java.spi.generation.StaticHelper.*;

import ${helper.packageName(protocolName, languageName, outputFlavor)}.io.*;
import ${helper.packageName(protocolName, languageName, outputFlavor)}.types.*;

import com.fasterxml.jackson.annotation.*;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.java.spi.generation.MessageIO;

import java.time.*;
import java.util.*;
import java.math.BigInteger;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public<#if type.abstract> abstract</#if> class ${typeName}<#if type.parentType??> extends ${type.parentType.name}</#if> implements Message {

<#if helper.isDiscriminatedType(type)>

    // Discriminator values used by the parser to determine the type to be used. All values have to apply.
    public static final Object[] DISCRIMINATOR_VALUES = new Object[] {
    <#list type.discriminatorValues as discriminatorValue>
        <#-- There are rare occasions where the discriminator is defined by a parser argument, in this case we currently can't detect the type -->
        (${helper.getDiscriminatorConstantType(type, discriminatorValue?index)}) ${discriminatorValue}<#sep>, </#sep>
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
    <#if helper.isAbstractField(field)>
    public ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> abstract get${field.name?cap_first}();
    <#else>
    public ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> get${field.name?cap_first}() {
        return ${field.name};
    }
    </#if>

</#list>
<#list type.virtualFields as field>
    public ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> get${field.name?cap_first}() {
        return (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field, field.valueExpression, type.parserArguments)});
    }

</#list>
    @Override
    @JsonIgnore
    public int getLengthInBytes() {
        return getLengthInBits() / 8;
    }

    @Override
    @JsonIgnore
    public int getLengthInBits() {
        int lengthInBits = <#if type.parentType??>super.getLengthInBits()<#else>0</#if>;
<#list type.fields as field>
<#switch field.typeName>
    <#case "array">

        // Array field
        <#if helper.isSimpleType(field.type)>
        lengthInBits += ${field.type.sizeInBits} * ${field.name}.length;
        <#else>
        if(${field.name} != null) {
            for(Message element : ${field.name}) {
                lengthInBits += element.getLengthInBits();
            }
        }
        </#if>
        <#break>
    <#case "checksum">

        // Checksum Field (${field.name})
        lengthInBits += ${field.type.sizeInBits};
        <#break>
    <#case "const">

        // Const Field (${field.name})
        lengthInBits += ${field.type.sizeInBits};
        <#break>
    <#case "discriminator">

        // Discriminator Field (${field.name})
        lengthInBits += ${field.type.sizeInBits};
        <#break>
    <#case "enum">

        // Enum Field (${field.name})
        lengthInBits += ${helper.getEnumBaseType(field.type).sizeInBits};
        <#break>
    <#case "implicit">

        // Implicit Field (${field.name})
        lengthInBits += ${field.type.sizeInBits};
        <#break>
    <#case "manualArray">

        // Manual Array Field (${field.name})
        lengthInBits += ${helper.toParseExpression(field, field.lengthExpression, null)} * 8;
        <#break>
    <#case "manual">

        // Manual Field (${field.name})
        lengthInBits += ${helper.toParseExpression(field, field.lengthExpression, null)} * 8;
        <#break>
    <#case "optional">

        // Optional Field (${field.name})
        if(${field.name} != null) {
        <#if helper.isSimpleType(field.type)>
            lengthInBits += ${field.type.sizeInBits};
        <#else>
            lengthInBits += ${field.name}.getLengthInBits();
        </#if>
        }
        <#break>
    <#case "padding">

        // Padding Field (${field.name})
        <#-- We're replacing the "lastItem" with 'false' here as the item itself can't know if it is the last -->
        if((boolean) (${helper.toParseExpression(field, field.paddingCondition, type.parserArguments)?replace("lastItem", "false")?no_esc})) {
            lengthInBits += ${field.type.sizeInBits};
        }
        <#break>
    <#case "reserved">

        // Reserved Field
        lengthInBits += ${field.type.sizeInBits};
        <#break>
    <#case "simple">

        // Simple field (${field.name})
        <#if helper.isSimpleType(field.type)>
        lengthInBits += ${field.type.sizeInBits};
        <#else>
        lengthInBits += ${field.name}.getLengthInBits();
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

        return lengthInBits;
    }

    public PlcValue toPlcValue() {
<#if type.propertyFields?has_content>
        Map<String, PlcValue> fieldValues = new TreeMap<>();
<#list type.propertyFields as field>
<#switch field.typeName>
<#case "array">
<#case "manualArray">

        {
            List<PlcValue> plcValues = new ArrayList<>(${field.name}.length);
            for(${helper.getLanguageTypeNameForField(field)} item : ${field.name}) {
                <#if helper.isSimpleType(field.type)>
                plcValues.add(PlcValues.of(item));
                <#else>
                plcValues.add(item.toPlcValue());
                </#if>
            }
            fieldValues.put("${field.name}", new PlcList(plcValues));
        }
<#break>
<#case "enum">

        fieldValues.put("${field.name}", new PlcString(${field.name}.name()));
<#break>
<#case "manual">
<#case "simple">

    <#if helper.isSimpleType(field.type)>
        fieldValues.put("${field.name}", PlcValues.of(${field.name}));
    <#else>
        fieldValues.put("${field.name}", ${field.name}.toPlcValue());
    </#if>
<#break>
<#case "optional">

        if(${field.name} != null) {
            <#if helper.isSimpleType(field.type)>
            fieldValues.put("${field.name}", PlcValues.of(${field.name}));
            <#else>
            fieldValues.put("${field.name}", ${field.name}.toPlcValue());
            </#if>
        } else {
            fieldValues.put("${field.name}", null);
        }
<#break>
<#case "switch">

        PlcStruct superStruct = (PlcStruct) super.toPlcValue();
        fieldValues.putAll(superStruct.getStruct());
<#break>
<#case "virtual">
    <#if helper.isSimpleType(field.type)>

        fieldValues.put("${field.name}", PlcValues.of(${field.name}));
    </#if>
<#break>
</#switch>
</#list>

        return new PlcStruct(fieldValues);
<#else>
        return null;
</#if>
    }

    @Override
    @JsonIgnore
    public MessageIO<<#if type.parentType??>${type.parentType.name}<#else>${typeName}</#if>, <#if type.parentType??>${type.parentType.name}<#else>${typeName}</#if>> getMessageIO() {
        return new <#if type.parentType??>${type.parentType.name}<#else>${typeName}</#if>IO();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof ${typeName})) {
            return false;
        }
        ${typeName} that = (${typeName}) o;
        return
            <#if type.propertyFields?has_content>
            <#list type.propertyFields as field>
            (get${field.name?cap_first}() == that.get${field.name?cap_first}()) &&
            </#list>
            </#if>
            <#if type.parentType??>
            super.equals(that) &&
            </#if>
            true;
    }

    @Override
    public int hashCode() {
        return Objects.hash(
            <#if type.parentType??>
            super.hashCode()<#if type.propertyFields?has_content>,</#if>
            </#if>
            <#if type.propertyFields?has_content>
            <#list type.propertyFields as field>
            get${field.name?cap_first}()<#sep>,</#sep>
            </#list>
            </#if>
        );
    }

    @Override
    public String toString() {
        return toString(ToStringStyle.SHORT_PREFIX_STYLE);
    }

    public String toString(ToStringStyle style) {
        return new ToStringBuilder(this, style)
            <#if type.parentType??>
            .appendSuper(super.toString(style))
            </#if>
            <#if type.propertyFields?has_content>
            <#list type.propertyFields as field>
            .append("${field.name}", get${field.name?cap_first}())
            </#list>
            </#if>
            <#list type.virtualFields as field>
            .append("${field.name}", get${field.name?cap_first}())
            </#list>
            .toString();
    }

}