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

using System;
using System.Collections.Generic;
using System.Text;
using org.apache.plc4net.api.value;
using org.apache.plc4net.spi.generation;
using org.apache.plc4net.spi.model.values;

namespace org.apache.plc4net.drivers.${protocolName?replace("-", "")}.${outputFlavor?replace("-", "")}.model
{

    public class ${type.name}
    {

        public IPlcValue Parse(ReadBuffer io<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)} ${parserArgument.name}<#sep>, </#sep></#list></#if>)
        {
    <#list type.switchField.cases as case>
            if (<#if case.discriminatorValues?has_content><#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#else>${discriminatorValue}</#if><#sep> && </#sep></#list></#if>) { // ${case.name}
        <#assign valueDefined=false>
        <#if case.name == "Struct">
                var internalMap = new Dictionary<string, IPlcValue>();
            <#assign valueDefined=true>
        </#if>
        <#list case.fields as field>
            <#switch field.typeName>
                <#case "array">
                    <#assign arrayField = field>

                // Array Field (${field.name});
                var ${arrayField.name} = new List<IPlcValue>();
                for (int i = 0; i < ${helper.toParseExpression(null, field.loopExpression, type.parserArguments)}; i++) {
                    var internalItem = <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>Complex type array in data-io parsing currently not implemented</#if>;
                    ${field.name}.Add(new ${helper.getPlcValueTypeForTypeReference(field.type)}(internalItem));
                }
                    <#if field.name == "value">
                        <#assign valueDefined=true>
                    </#if>
                    <#break>
                <#case "manual">
                    <#assign manualArrayField = field>

                // Manual Field (${field.name})
                var ${field.name} = ${helper.toParseExpression(field, field.parseExpression, type.parserArguments)};
                    <#if field.name == "value">
                        <#assign valueDefined=true>
                    </#if>
                    <#break>
                <#case "reserved">
                    <#assign reservedField = field>

                // Reserved Field (Just skip the bits)
                ${helper.getReadBufferReadMethodCall(field.type)};
                    <#break>
                <#case "simple">
                    <#assign simpleField = field>

                // Simple Field (${field.name})
                var ${field.name} = <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
                    <#if case.name == "Struct">
                internalMap["${case.name}"] = new ${helper.getPlcValueTypeForTypeReference(field.type)}(${field.name});
                    </#if>
                    <#assign valueDefined=true>
                    <#break>
            </#switch>
        </#list>
        <#if valueDefined>

            <#switch case.name>
                <#case "TIME_OF_DAY">
                    <#if helper.hasFieldsWithNames(case.fields, "hour", "minutes", "seconds")>
                var value = new DateTime(0,0,0, hour, minutes, seconds);
                    </#if>
                return new PlcTIME_OF_DAY(value);
                    <#break>
                <#case "DATE">
                    <#if helper.hasFieldsWithNames(case.fields, "year", "month", "day")>
                var value = new DateTime(year, month, day, 0, 0, 0);
                    </#if>
                return new PlcDATE(value);
                    <#break>
                <#case "DATE_AND_TIME">
                    <#if helper.hasFieldsWithNames(case.fields, "year", "month", "day", "hour", "minutes", "seconds")>
                var value = new DateTime(year, month, day, hour, minutes, seconds);
                    </#if>
                return new PlcDATE_AND_TIME(value);
                    <#break>
                <#case "Struct">
                return new PlcStruct(internalMap);
                    <#break>
                <#case "List">
                return new PlcList(value);
                    <#break>
                <#default>
                return new Plc${case.name}(value);
            </#switch>
        </#if>
            } <#sep>else
    </#list>

            return null;
        }

        public void Serialize(WriteBuffer io, IPlcValue value<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)} ${parserArgument.name}<#sep>, </#sep></#list></#if>)
        {
    <#list type.switchField.cases as case>
            if (<#if case.discriminatorValues?has_content><#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#else>${discriminatorValue}</#if><#sep> && </#sep></#list></#if>) { // ${case.name}
        <#list case.fields as field>
            <#switch field.typeName>
                <#case "array">
                    <#assign arrayField = field>

                // Array Field (${field.name})
                for (int i = 0; i < ${helper.toSerializationExpression(null, field.loopExpression, type.parserArguments)}; i++) {
                    <#if helper.isSimpleTypeReference(arrayField.type)>${helper.getWriteBufferWriteMethodCall(arrayField.type, arrayField.name + ".GetIndex(i).Get" + helper.getLanguageTypeNameForTypeReference(arrayField.type)?cap_first + "()")}<#else>Complex type array in data-io serialization currently not implemented</#if>;
                }
                    <#break>
                <#case "manual">
                    <#assign manualField = field>

                // Manual Field (${field.name})
                ${helper.toSerializationExpression(field, manualField.serializeExpression, type.parserArguments)};
                    <#break>
                <#case "reserved">
                    <#assign reservedField = field>

                // Reserved Field (Just skip the bytes)
                ${helper.getWriteBufferWriteMethodCall(field.type, helper.getReservedValue(field))};
                    <#break>
                <#case "simple">
                    <#assign simpleField = field>

                // Simple Field (${field.name})
                <#if helper.isSimpleTypeReference(field.type)>${helper.getWriteBufferWriteMethodCall(field.type, "value.Get" + helper.getLanguageTypeNameForTypeReference(field.type)?cap_first + "()")}<#else>${field.type.name}Serialize(io, <#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>;
                    <#break>
            </#switch>
        </#list>
            } <#sep>else
    </#list>
        }

    }

}

</#outputformat>
