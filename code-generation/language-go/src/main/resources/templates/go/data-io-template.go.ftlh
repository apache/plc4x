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
<#-- @ftlvariable name="tracer" type="org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer" -->
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
/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

<#macro importSectionWithContentBelow>
	<#local sectionContent><#nested></#local>
import (
	<#if helper.getRequiredImportsForDataIo()?has_content>
		<#list helper.getRequiredImportsForDataIo() as import>
	${import}
		</#list>
	</#if>
)

	${sectionContent}
</#macro>

<#macro emitImport import>${helper.emitDataIoRequiredImport(import)}</#macro>
<#macro emitImportWithAlias alias import>${helper.emitDataIoRequiredImport(alias, import)}</#macro>

<@importSectionWithContentBelow><@emitImport import="github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils" /><@emitImport import="github.com/apache/plc4x/plc4go/internal/plc4go/spi/values" /><@emitImportWithAlias alias="api" import="github.com/apache/plc4x/plc4go/pkg/plc4go/values" />
// Code generated by code-generation. DO NOT EDIT.

func ${type.name}Parse(readBuffer utils.ReadBuffer<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} <#if helper.isComplexTypeReference(parserArgument.type) && !helper.isEnumTypeReference(parserArgument.type)>I</#if>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) (api.PlcValue, error) {
	readBuffer.PullContext("${type.name}")
	switch {
	<#list type.switchField.cases as case>
        <#if case.discriminatorValues?has_content>case <#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#elseif helper.isComplexTypeReference(type.parserArguments[discriminatorValue?index].type)><#if helper.isEnumTypeReference(type.parserArguments[discriminatorValue?index].type)>${helper.getLanguageTypeNameForTypeReference(type.parserArguments[discriminatorValue?index].type)}_${discriminatorValue}<#else>${discriminatorValue}</#if><#else>${discriminatorValue}</#if><#sep> && </#list> <#else>default</#if>: // ${case.name}
		<#assign valueDefined=false>
		<#if case.name == "Struct">
			_map := map[string]api.PlcValue{}
			<#assign valueDefined=true>
		</#if>
		<#list case.fields as field>
			<#switch field.typeName>
				<#case "array">
					<#assign arrayField = field>

			// Array Field (${field.name})
			var ${field.name} []api.PlcValue
			for i := 0; i < int(${helper.toParseExpression(null, field.loopExpression, type.parserArguments)}); i++ {
				_item, _itemErr := <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.name, field.type)}<#else>Complex type array in data-io parsing currently not implemented</#if>
				if _itemErr != nil {
					return nil, errors.Wrap(_itemErr, "Error parsing 'value' field")<@emitImport import="github.com/pkg/errors" />
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
				return nil, errors.Wrap(_${field.name}Err, "Error parsing '${field.name}' field")<@emitImport import="github.com/pkg/errors" />
			}
					<#if field.name == "value">
						<#assign valueDefined=true>
					</#if>
					<#break>
				<#case "reserved">
					<#assign reservedField = field>

			// Reserved Field (Just skip the bytes)
			if _, _err := ${helper.getReadBufferReadMethodCall("reserved", field.type)}; _err != nil {
				return nil, errors.Wrap(_err, "Error parsing reserved field")<@emitImport import="github.com/pkg/errors" />
			}
					<#break>
				<#case "simple">
					<#assign simpleField = field>

			// Simple Field (${field.name})
			<#if case.name == "Struct" || ((case.name == "DATE_AND_TIME") && ((field.name == "year") || (field.name == "month") || (field.name == "day") || (field.name == "hour") || (field.name == "minutes") || (field.name == "seconds") || (field.name == "secondsSinceEpoch"))) || ((case.name == "DATE") && ((field.name == "year") || (field.name == "month") || (field.name == "day"))) || ((case.name == "TIME_OF_DAY") && ((field.name == "hour") || (field.name == "minutes") || (field.name == "seconds"))) || field.name == "value">${field.name}<#else>_</#if>, _${field.name}Err := <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.name, field.type)}<#else>${field.type.name}Parse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>
			if _${field.name}Err != nil {
				return nil, errors.Wrap(_${field.name}Err, "Error parsing '${field.name}' field")<@emitImport import="github.com/pkg/errors" />
			}
					<#if case.name == "Struct">
			_map["${case.name}"] = ${helper.getPlcValueTypeForTypeReference(field.type)}(${field.name})
					</#if>
					<#assign valueDefined=true>
					<#break>
			</#switch>
		</#list>
		<#if valueDefined>
			<#switch case.name>
				<#case "TIME_OF_DAY">
					<#if helper.hasFieldsWithNames(case.fields, "hour", "minutes", "seconds")>
			value := time.Date(0,0,0, int(hour), int(minutes), int(seconds), 0, nil)<@emitImport import="time" />
					</#if>
			readBuffer.CloseContext("${type.name}")
			return values.NewPlcTIME_OF_DAY(value), nil
					<#break>
				<#case "DATE">
					<#if helper.hasFieldsWithNames(case.fields, "year", "month", "day")>
			value := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, nil)<@emitImport import="time" />
					</#if>
			readBuffer.CloseContext("${type.name}")
			return values.NewPlcDATE(value), nil
					<#break>
				<#case "DATE_AND_TIME">
					<#if helper.hasFieldsWithNames(case.fields, "year", "month", "day", "hour", "minutes", "seconds")>
			value := time.Date(int(year), time.Month(month), int(day), int(hour), int(minutes), int(seconds), 0, nil)<@emitImport import="time" />
                    <#elseif helper.hasFieldsWithNames(case.fields, "secondsSinceEpoch")>
            value := time.Unix(int64(secondsSinceEpoch), 0)<@emitImport import="time" />
					</#if>
			readBuffer.CloseContext("${type.name}")
			return values.NewPlcDATE_AND_TIME(value), nil
					<#break>
			readBuffer.CloseContext("${type.name}")
				<#case "Struct">
			readBuffer.CloseContext("${type.name}")
			return values.NewPlcStruct(_map), nil
					<#break>
				<#case "List">
			readBuffer.CloseContext("${type.name}")
			return values.NewPlcList(value), nil
					<#break>
				<#default>
			readBuffer.CloseContext("${type.name}")
			return values.NewPlc${case.name}(value), nil
			</#switch>
		</#if>
	</#list>
	}
    // TODO: add more info which type it is actually
	return nil, errors.New("unsupported type")<@emitImport import="github.com/pkg/errors" />
}

func ${type.name}Serialize(writeBuffer utils.WriteBuffer, value api.PlcValue<#if type.parserArguments?has_content>, <#list type.parserArguments as parserArgument>${parserArgument.name} <#if helper.isComplexTypeReference(parserArgument.type) && !helper.isEnumTypeReference(parserArgument.type)>I</#if>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#sep>, </#sep></#list></#if>) error {
	writeBuffer.PushContext("${type.name}")
	switch {
	<#assign defaultSet=false>
	<#list type.switchField.cases as case>
        <#if case.discriminatorValues?has_content>case <#list case.discriminatorValues as discriminatorValue>${helper.toParseExpression(null, type.switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == <#if helper.discriminatorValueNeedsStringEqualityCheck(type.switchField.discriminatorExpressions[discriminatorValue?index])>"${discriminatorValue}"<#elseif helper.isComplexTypeReference(type.parserArguments[discriminatorValue?index].type)><#if helper.isEnumTypeReference(type.parserArguments[discriminatorValue?index].type)>${helper.getLanguageTypeNameForTypeReference(type.parserArguments[discriminatorValue?index].type)}_${discriminatorValue}<#else>${discriminatorValue}</#if><#else>${discriminatorValue}</#if><#sep> && </#list> <#else>default<#assign defaultSet=true></#if>: // ${case.name}
		<#list case.fields as field>
			<#switch field.typeName>
				<#case "array">
					<#assign arrayField = field>

			// Array Field (${field.name})
			for i := uint32(0); i < uint32(${helper.toSerializationExpression(null, field.loopExpression, type.parserArguments)}); i++ {
					<#if case.name = "Struct">
				${arrayField.name} := value.GetValue("${arrayField.name}")
					</#if>
				_itemErr := <#if helper.isSimpleTypeReference(arrayField.type)>${helper.getWriteBufferWriteMethodCall("", arrayField.type, arrayField.name + ".GetIndex(i).Get" + helper.getLanguageTypeNameForTypeReference(arrayField.type)?cap_first + "()", field)}<#else>Complex type array in data-io serialization currently not implemented</#if>
				if _itemErr != nil {
					return errors.Wrap(_itemErr, "Error serializing 'value' field")<@emitImport import="github.com/pkg/errors" />
				}
			}
					<#break>
				<#case "manual">
					<#assign manualField = field>

			// Manual Field (${field.name})
			_${field.name}Err := ${helper.toSerializationExpression(field, manualField.serializeExpression, type.parserArguments)}
			if _${field.name}Err != nil {
				return errors.Wrap(_${field.name}Err, "Error serializing '${field.name}' field")<@emitImport import="github.com/pkg/errors" />
			}
					<#break>
				<#case "reserved">
					<#assign reservedField = field>

			// Reserved Field (Just skip the bytes)
			if _err := ${helper.getWriteBufferWriteMethodCall("reserved", field.type, helper.getReservedValue(field), field)}; _err != nil {
				return errors.Wrap(_err, "Error serializing reserved field")<@emitImport import="github.com/pkg/errors" />
			}
					<#break>
				<#case "simple">
					<#assign simpleField = field>

			// Simple Field (${field.name})
			if _err := <#if helper.isSimpleTypeReference(field.type)>${helper.getWriteBufferWriteMethodCall(field.name, field.type, "value.Get" + helper.getLanguageTypeNameForTypeReference(field.type)?cap_first + "()", field)}<#else>${field.type.name}Serialize(io, <#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getLanguageTypeNameForTypeReference(helper.getArgumentType(field.type, parserArgument?index), true)}) (${helper.toParseExpression(field, parserArgument, type.parserArguments)})<#sep>, </#sep></#list></#if>)</#if>; _err != nil {
				return errors.Wrap(_err, "Error serializing '${field.name}' field")<@emitImport import="github.com/pkg/errors" />
			}
					<#break>
			</#switch>
		</#list>
	</#list>
	<#if !defaultSet>
		default:
            // TODO: add more info which type it is actually
			return errors.New("unsupported type")<@emitImport import="github.com/pkg/errors" />
	</#if>
	}
	writeBuffer.PopContext("${type.name}")
	return nil
}
</@importSectionWithContentBelow>

</#outputformat>
