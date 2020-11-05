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
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.go.GoLanguageTemplateHelper" -->
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
${helper.fileName(protocolName, languageName, outputFlavor)?replace(".", "/")}/model/${type.name}.go
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
package model

import (
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/model/values"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    api "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/values"
)

func ${type.name}Parse(io *utils.ReadBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} <#if helper.isComplexTypeReference(parserArgument.type)>I</#if>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) (api.PlcValue, error) {
    switch {
    <#list type.switchField.cases as case>
        case <#if case.discriminatorValues?has_content><#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#else>${discriminatorValue}</#if><#sep> && </#sep></#list></#if>: // ${case.name}
        <#assign valueDefined=false>
        <#if case.name == "Struct">
            _map := map[string]interface{}{}
            <#assign valueDefined=true>
        </#if>
        <#list case.fields as field>
            <#switch field.typeName>
                <#case "array">
                    <#assign arrayField = field>

            // Array Field (${field.name})
            var ${field.name} []api.PlcValue
            for i := 0; i < int(${helper.toParseExpression(null, field.loopExpression, type.parserArguments)}); i++ {
                _item, _itemErr := <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>Complex type array in data-io parsing currently not implemented</#if>
                if _itemErr != nil {
                    return nil, errors.New("Error parsing 'value' field " + _itemErr.Error())
                }
                ${field.name} = append(${field.name}, ${helper.getPlcValueTypeForTypeReference(field.type)}(_item))
            }
                    <#if field.name == "value">
                        <#assign valueDefined=true>
                    </#if>
                    <#break>
                <#case "manual">
                    <#assign arrayField = field>

            // Manual Field (${field.name})
            ${field.name}, _${field.name}Err := ${helper.toParseExpression(field, field.parseExpression, type.parserArguments)}
            if _${field.name}Err != nil {
                return nil, errors.New("Error parsing '${field.name}' field " + _${field.name}Err.Error())
            }
                    <#if field.name == "value">
                        <#assign valueDefined=true>
                    </#if>
                    <#break>
                <#case "reserved">
                    <#assign reservedField = field>

            // Reserved Field (Just skip the bytes)
            if _, _err := ${helper.getReadBufferReadMethodCall(field.type)}; _err != nil {
                return nil, errors.New("Error parsing reserved field " + _err.Error())
            }
                    <#break>
                <#case "simple">
                    <#assign simpleField = field>

            // Simple Field (${field.name})
            <#if case.name == "Struct" || ((case.name == "DATE_AND_TIME") && ((field.name == "year") || (field.name == "month") || (field.name == "day") || (field.name == "hour") || (field.name == "minutes") || (field.name == "seconds"))) || ((case.name == "DATE") && ((field.name == "year") || (field.name == "month") || (field.name == "day"))) || ((case.name == "TIME_OF_DAY") && ((field.name == "hour") || (field.name == "minutes") || (field.name == "seconds"))) || field.name == "value">${field.name}<#else>_</#if>, _${field.name}Err := <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)}<#else>${field.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>
            if _${field.name}Err != nil {
                return nil, errors.New("Error parsing '${field.name}' field " + _${field.name}Err.Error())
            }
                    <#if case.name == "Struct">
            _map["${case.name}"] = ${field.name}
                    </#if>
                    <#assign valueDefined=true>
                    <#break>
            </#switch>
        </#list>
        <#if valueDefined>
            <#switch case.name>
                <#case "TIME_OF_DAY">
            value := time.Date(0,0,0, int(hour), int(minutes), int(seconds), 0, nil)
            return values.NewPlcTIME_OF_DAY(value), nil
                    <#break>
                <#case "DATE">
            value := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, nil)
            return values.NewPlcDATE(value), nil
                    <#break>
                <#case "DATE_AND_TIME">
            value := time.Date(int(year), time.Month(month), int(day), int(hour), int(minutes), int(seconds), 0, nil)
            return values.NewPlcDATE_AND_TIME(value), nil
                    <#break>
                <#case "Struct">
            return values.NewPlcStruct(_map), nil
                    <#break>
                <#case "List">
            return values.NewPlcList(value), nil
                    <#break>
                <#default>
            return values.NewPlc${case.name}(value), nil
            </#switch>
        </#if>
    </#list>
    }
    return nil, errors.New("unsupported type")
}

func ${type.name}Serialize(io *utils.WriteBuffer, value api.PlcValue<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} <#if helper.isComplexTypeReference(parserArgument.type)>I</#if>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) error {
    switch {
    <#list type.switchField.cases as case>
        case <#if case.discriminatorValues?has_content><#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#else>${discriminatorValue}</#if><#sep> && </#sep></#list></#if>: // ${case.name}
        <#list case.fields as field>
            <#switch field.typeName>
                <#case "array">
                    <#assign arrayField = field>

            // Array Field (${field.name})
            for i := uint32(0); i < uint32(${helper.toSerializationExpression(null, field.loopExpression, type.parserArguments)}); i++ {
                _itemErr := <#if helper.isSimpleTypeReference(arrayField.type)>${helper.getWriteBufferWriteMethodCall(arrayField.type, arrayField.name + ".GetIndex(i).Get" + helper.getLanguageTypeNameForTypeReference(arrayField.type)?cap_first + "()")}<#else>Complex type array in data-io serialization currently not implemented</#if>
                if _itemErr != nil {
                    return errors.New("Error serializing 'value' field " + _itemErr.Error())
                }
            }
                    <#break>
                <#case "manual">
                    <#assign manualField = field>

            // Manual Field (${field.name})
            _${field.name}Err := ${helper.toSerializationExpression(field, manualField.serializeExpression, type.parserArguments)}
            if _${field.name}Err != nil {
                return errors.New("Error serializing '${field.name}' field " + _${field.name}Err.Error())
            }
                    <#break>
                <#case "reserved">
                    <#assign reservedField = field>

            // Reserved Field (Just skip the bytes)
            if _err := ${helper.getWriteBufferWriteMethodCall(field.type, helper.getReservedValue(field))}; _err != nil {
                return errors.New("Error serializing reserved field " + _err.Error())
            }
                    <#break>
                <#case "simple">
                    <#assign simpleField = field>

            // Simple Field (${field.name})
            if _err := <#if helper.isSimpleTypeReference(field.type)>${helper.getWriteBufferWriteMethodCall(field.type, "value.Get" + helper.getLanguageTypeNameForTypeReference(field.type)?cap_first + "()")}<#else>${field.type.name}Serialize(io, <#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>; _err != nil {
                return errors.New("Error serializing '${field.name}' field " + _err.Error())
            }
                    <#break>
            </#switch>
        </#list>
    </#list>
        default:

            return errors.New("unsupported type")
    }
    return nil
}

</#outputformat>
