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
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.cs.CsLanguageTemplateHelper" -->
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
${helper.fileName(protocolName, languageName, outputFlavor)?replace(".", "/")}/model/${type.name}.cs
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

namespace org.apache.plc4net.drivers.${protocolName?replace("-", "")}.${outputFlavor?replace("-", "")}.model
{

    public<#if helper.isDiscriminatedParentTypeDefinition()> abstract</#if> class ${type.name}<#if type.parentType??> : ${type.parentType.name}</#if>
    {
    <#--
        If this is a discriminated child type, we need to generate methods for accessing it's discriminator
        values, as if they were normal java properties.
    -->
    <#if helper.isDiscriminatedChildTypeDefinition()>
        <#assign discriminatedChildType = type>

        <#-- @ftlvariable name="discriminatedChildType" type="org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition" -->
        // Accessors for discriminator values.
        <#list helper.getDiscriminatorValues(discriminatedChildType) as discriminatorName, discriminatorValue>
        <#-- If the discriminator name matches that of another field, suppress the methods generation -->
            <#if !helper.isNonDiscriminatorField(discriminatorName)>
        public override ${helper.getLanguageTypeNameForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])} Get${discriminatorName?cap_first}() {
            return <#if discriminatorValue??>${discriminatorValue}<#else>${helper.getNullValueForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])}</#if>;
        }
            </#if>
        </#list>
    </#if>
    <#--
        If this is a discriminated parent type, we need to generate the abstract methods for accessing it's
        discriminator values instead.
    -->
    <#if helper.isDiscriminatedParentTypeDefinition()>
        <#assign discriminatedParentType = type>

    <#-- @ftlvariable name="discriminatedParentType" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
        // Abstract accessors for discriminator values.
        <#list helper.getDiscriminatorTypes() as discriminatorName, discriminatorType>
        <#-- If the discriminator name matches that of another field, suppress the methods generation -->
            <#if !helper.isNonDiscriminatorField(discriminatorName)>
        public abstract ${helper.getLanguageTypeNameForTypeReference(discriminatorType)} Get${discriminatorName?cap_first}();
            </#if>
        </#list>
    </#if>
    <#-- If the current type contains "const" fields, generate some java constants for holing their values -->
    <#if type.constFields?has_content>

        // Constant values.
        <#list type.constFields as field>
        public const ${helper.getLanguageTypeNameForField(field)} ${field.name?upper_case} = ${field.referenceValue};
        </#list>
    </#if>
    <#-- Prpoerty fields are fields that require a property in the pojo -->
    <#if type.propertyFields?has_content>

        // Properties.
        <#list type.propertyFields as field>
        public ${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name?cap_first} { get; }
        </#list>
    </#if>

    <#-- getAllPropertyFields() returns not only the property fields of this type but also of it's parents -->
        public ${type.name}(<#list type.getAllPropertyFields() as field>${helper.getLanguageTypeNameForField(field)}<#if field.loopType??>[]</#if> ${field.name}<#sep>, </#sep></#list>)
    <#if type.getParentPropertyFields()?has_content>
            : base(<#list type.getParentPropertyFields() as field>${field.name}<#sep>, </#sep></#list>)
    </#if>
        {
    <#list type.propertyFields as field>
            ${field.name?cap_first} = ${field.name};
    </#list>
        }

    }

}
</#outputformat>