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
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.c.CLanguageTemplateHelper" -->
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
<#if !helper.isDiscriminatedChildTypeDefinition(type)>${helper.getSourceDirectory()?replace(".", "/")}/${helper.camelCaseToSnakeCase(type.name)}.c
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

#include <stdio.h>
#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/spi/evaluation_helper.h>
#include "${helper.camelCaseToSnakeCase(type.name)}.h"

<#-- Helper function to get the discriminator for a given enum type constant -->
<#if helper.isDiscriminatedParentTypeDefinition(type)>
// Array of discriminator values that match the enum type constants.
// (The order is identical to the enum constants so we can use the
// enum constant to directly access a given types discriminator values)
const plc4c_${helper.getCTypeName(type.name)}_discriminator plc4c_${helper.getCTypeName(type.name)}_discriminators[] = {
    <#if helper.getDiscriminatorValues()?has_content>
        <#list helper.getDiscriminatorValues() as name, typeValues>
  {/* ${helper.getCTypeName(name)} */
   <#list typeValues as key, value>.${key} = <#if value??>${value}<#else>-1</#if><#sep>, </#list>}<#sep >,
</#list>
    </#if>

};

// Function returning the discriminator values for a given type constant.
plc4c_${helper.getCTypeName(type.name)}_discriminator plc4c_${helper.getCTypeName(type.name)}_get_discriminator(plc4c_${helper.getCTypeName(type.name)}_type type) {
  return plc4c_${helper.getCTypeName(type.name)}_discriminators[type];
}

</#if>
// Parse function.
plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_parse(plc4c_spi_read_buffer* buf, <#if type.parserArguments?has_content><#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#if !helper.isSimpleTypeReference(parserArgument.type)>*</#if> ${parserArgument.name}<#sep>, </#list>, </#if>plc4c_${helper.getCTypeName(type.name)}** _message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Allocate enough memory to contain this data structure.
  (*_message) = malloc(sizeof(plc4c_${helper.getCTypeName(type.name)}));
  if(*_message == NULL) {
    return NO_MEMORY;
  }
<#macro fieldName baseType field><#if baseType != type>${helper.camelCaseToSnakeCase(baseType.name)}_</#if>${helper.camelCaseToSnakeCase(field.name)}</#macro>
<#macro fieldParser baseType field indentContent>
    <#switch field.typeName>
        <#case "array">
            <#assign arrayField = field>

<#if indentContent>  </#if>  // Array field (${arrayField.name})
            <#-- Only update curPos if the length expression uses it -->
            <#if arrayField.loopExpression.contains("curPos")>
<#if indentContent>  </#if>  curPos = plc4c_spi_read_get_pos(buf) - startPos;
            </#if>
            <#-- In all other cases do we have to work with a list, that is later converted to an array -->
<#if indentContent>  </#if>  plc4c_list* ${arrayField.name} = malloc(sizeof(plc4c_list));
<#if indentContent>  </#if>  if(${arrayField.name} == NULL) {
<#if indentContent>  </#if>    return NO_MEMORY;
<#if indentContent>  </#if>  }
<#if indentContent>  </#if>  {
            <#if helper.isCountArrayField(field)>
<#if indentContent>  </#if>    // Count array
<#if indentContent>  </#if>    uint8_t itemCount = ${helper.toParseExpression(baseType, arrayField, arrayField.loopExpression, baseType.parserArguments)};
<#if indentContent>  </#if>    for(int curItem = 0; curItem < itemCount; curItem++) {
<#if indentContent>  </#if>      <#if !helper.isSimpleTypeReference(arrayField.type)>bool lastItem = curItem == (itemCount - 1);</#if>
      <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
<#if indentContent>  </#if>                <#if helper.isSimpleTypeReference(arrayField.type)>
<#if indentContent>  </#if>      ${helper.getLanguageTypeNameForTypeReference(arrayField.type)} _value = ${helper.getReadBufferReadMethodCall(arrayField.type)};
<#if indentContent>  </#if>      plc4c_utils_list_insert_head_value(${arrayField.name}, &_value);
                <#else>
      <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
<#if indentContent>  </#if>      ${helper.getLanguageTypeNameForField(field)}* _value = NULL;
<#if indentContent>  </#if>      plc4c_return_code _res = plc4c_${helper.getCTypeName(arrayField.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(baseType, field, parserTerm, baseType.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &_value);
<#if indentContent>  </#if>      if(_res != OK) {
<#if indentContent>  </#if>        return _res;
<#if indentContent>  </#if>      }
<#if indentContent>  </#if>      plc4c_utils_list_insert_head_value(${arrayField.name}, _value);
                </#if>
<#if indentContent>  </#if>    }
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
            <#elseif helper.isLengthArrayField(field)>
<#if indentContent>  </#if>    // Length array
<#if indentContent>  </#if>    uint8_t _${arrayField.name}Length = ${helper.toParseExpression(baseType, arrayField, arrayField.loopExpression, baseType.parserArguments)};
<#if indentContent>  </#if>    uint8_t ${arrayField.name}EndPos = plc4c_spi_read_get_pos(buf) + _${arrayField.name}Length;
<#if indentContent>  </#if>    while(plc4c_spi_read_get_pos(buf) < ${arrayField.name}EndPos) {
                <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
                <#if helper.isSimpleTypeReference(arrayField.type)>
<#if indentContent>  </#if>      ${helper.getLanguageTypeNameForTypeReference(arrayField.type)} _value = ${helper.getReadBufferReadMethodCall(arrayField.type)};
<#if indentContent>  </#if>      plc4c_utils_list_insert_head_value(${arrayField.name}, &_value);
                <#else>
                <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
<#if indentContent>  </#if>      ${helper.getLanguageTypeNameForField(field)}* _value = NULL;
<#if indentContent>  </#if>      plc4c_return_code _res = plc4c_${helper.getCTypeName(arrayField.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(baseType, field, parserTerm, baseType.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &_value);
<#if indentContent>  </#if>      if(_res != OK) {
<#if indentContent>  </#if>        return _res;
<#if indentContent>  </#if>      }
<#if indentContent>  </#if>      plc4c_utils_list_insert_head_value(${arrayField.name}, _value);
                </#if>
                <#-- After parsing, update the current position, but only if it's needed -->
                <#if arrayField.loopExpression.contains("curPos")>
<#if indentContent>  </#if>      curPos = plc4c_spi_read_get_pos(buf) - startPos;
                </#if>
<#if indentContent>  </#if>    }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
            <#elseif helper.isTerminatedArrayField(field)>
<#if indentContent>  </#if>    // Terminated array
<#if indentContent>  </#if>    while(!((boolean) (${helper.toParseExpression(baseType, arrayField, arrayField.loopExpression, baseType.parserArguments)}))) {
            <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
                <#if helper.isSimpleTypeReference(arrayField.type)>
<#if indentContent>  </#if>      ${helper.getLanguageTypeNameForTypeReference(arrayField.type)} _value = ${helper.getReadBufferReadMethodCall(arrayField.type)};
<#if indentContent>  </#if>      plc4c_utils_list_insert_head_value(${arrayField.name}, &_value);
                <#else>
                <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
<#if indentContent>  </#if>      ${helper.getLanguageTypeNameForField(field)}* _value = NULL;
<#if indentContent>  </#if>      plc4c_return_code _res = plc4c_${helper.getCTypeName(arrayField.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(baseType, field, parserTerm, baseType.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &_value);
<#if indentContent>  </#if>      if(_res != OK) {
<#if indentContent>  </#if>        return _res;
<#if indentContent>  </#if>      }
<#if indentContent>  </#if>      plc4c_utils_list_insert_head_value(${arrayField.name}, _value);
                </#if>
              <#-- After parsing, update the current position, but only if it's needed -->
              <#if arrayField.loopExpression.contains("curPos")>
<#if indentContent>  </#if>      curPos = plc4c_spi_read_get_pos(buf) - startPos;
              </#if>
<#if indentContent>  </#if>    }
            </#if>
<#if indentContent>  </#if>  }
<#if indentContent>  </#if>  (*_message)-><@fieldName baseType=baseType field=arrayField/> = ${arrayField.name};
            <#break>
        <#case "checksum">
            <#assign checksumField = field>
            <#assign simpleTypeReference = checksumField.type>

<#if indentContent>  </#if>  // Checksum Field (${checksumField.name})
<#if indentContent>  </#if>  {
<#if indentContent>  </#if>    // Create an array of all the bytes read in this message element so far.
<#if indentContent>  </#if>    byte[] checksumRawData = plc4c_spi_read_get_bytes(buf, startPos, plc4c_spi_read_get_pos(buf));
<#if indentContent>  </#if>    ${helper.getLanguageTypeNameForField(field)} _checksumRef = ${helper.getReadBufferReadMethodCall(checksumField.type)};
<#if indentContent>  </#if>    ${helper.getLanguageTypeNameForField(field)} _checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(baseType, checksumField, checksumField.checksumExpression, baseType.parserArguments)});
<#if indentContent>  </#if>    if(_checksum != _checksumRef) {
<#if indentContent>  </#if>      return PARSE_ERROR;
<#if indentContent>  </#if>      // throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X",_checksumRef & 0xFFFF, _checksum & 0xFFFF));
<#if indentContent>  </#if>    }
<#if indentContent>  </#if>  }
            <#break>
        <#case "const">
            <#assign constField = field>
            <#assign simpleTypeReference = constField.type>

<#if indentContent>  </#if>  // Const Field (${constField.name})
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${constField.name} = ${helper.getReadBufferReadMethodCall(constField.type)};
<#if indentContent>  </#if>  if(${constField.name} != ${helper.getCTypeName(baseType.name)?upper_case}_${helper.camelCaseToSnakeCase(constField.name)?upper_case}) {
<#if indentContent>  </#if>    return PARSE_ERROR;
<#if indentContent>  </#if>    // throw new ParseException("Expected constant value " + ${helper.getCTypeName(baseType.name)?upper_case}_${helper.camelCaseToSnakeCase(constField.name)?upper_case} + " but got " + ${constField.name});
<#if indentContent>  </#if>  }
            <#break>
        <#case "discriminator">
            <#assign discriminatorField = field>
            <#assign simpleTypeReference = discriminatorField.type>

<#if indentContent>  </#if>  // Discriminator Field (${discriminatorField.name}) (Used as input to a switch field)
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${discriminatorField.name} = ${helper.getReadBufferReadMethodCall(discriminatorField.type)};
            <#break>
        <#case "enum">
            <#assign enumField = field>

<#if indentContent>  </#if>  // Enum field (${enumField.name})
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${enumField.name} = ${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(enumField.type))};
<#if indentContent>  </#if>  (*_message)-><@fieldName baseType=baseType field=enumField/> = ${enumField.name};
            <#break>
        <#case "implicit">
            <#assign implicitField = field>
            <#assign simpleTypeReference = implicitField.type>

<#if indentContent>  </#if>  // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${implicitField.name} = ${helper.getReadBufferReadMethodCall(implicitField.type)};
            <#break>
        <#case "manualArray">
            <#assign manualArrayField = field>

  <#-- TODO: Implement -->
            <#break>
        <#case "manual">
            <#assign manualField = field>

<#if indentContent>  </#if>  // Manual Field (${manualField.name})<#-- check if complex -->
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${manualField.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(baseType, manualField, manualField.parseExpression, baseType.parserArguments)});
<#if indentContent>  </#if>  (*_message)-><@fieldName baseType=baseType field=manualField/> = ${manualField.name};
        <#break>
        <#case "optional">
            <#assign optionalField = field>

<#if indentContent>  </#if>  // Optional Field (${optionalField.name}) (Can be skipped, if a given expression evaluates to false)
            <#if optionalField.conditionExpression.contains("curPos")>
<#if indentContent>  </#if>  curPos = plc4c_spi_read_get_pos(buf) - startPos;
            </#if>
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)}* ${optionalField.name} = NULL;
<#if indentContent>  </#if>  if(${helper.toParseExpression(baseType, field, optionalField.conditionExpression, baseType.parserArguments)}) {
<#if indentContent>  </#if>    ${optionalField.name} = malloc(sizeof(${helper.getLanguageTypeNameForField(field)}));
<#if indentContent>  </#if>    if(${optionalField.name} == NULL) {
<#if indentContent>  </#if>      return NO_MEMORY;
<#if indentContent>  </#if>    }
        <#if helper.isSimpleTypeReference(optionalField.type)>
<#if indentContent>  </#if>    *${optionalField.name} = ${helper.getReadBufferReadMethodCall(optionalField.type)};
        <#else>
<#if indentContent>  </#if>    plc4c_return_code _res = plc4c_${helper.getCTypeName(optionalField.type.name)}_parse(buf<#if optionalField.params?has_content>, <#list optionalField.params as parserTerm>${helper.toParseExpression(baseType, optionalField, parserTerm, baseType.parserArguments)}<#sep>, </#sep></#list></#if>, &${optionalField.name});
<#if indentContent>  </#if>    if(_res != OK) {
<#if indentContent>  </#if>      return _res;
<#if indentContent>  </#if>    }
        </#if>
<#if indentContent>  </#if>    (*_message)-><@fieldName baseType=baseType field=optionalField/> = ${optionalField.name};
<#if indentContent>  </#if>  }
            <#break>
        <#case "padding">
            <#assign paddingField = field>
            <#assign simpleTypeReference = paddingField.type>

<#if indentContent>  </#if>  // Padding Field (padding)
<#if indentContent>  </#if>  {
<#if indentContent>  </#if>    bool _needsPadding = (bool) ((plc4c_spi_read_has_more(buf, ${helper.getNumBits(paddingField.type)})) && (${helper.toParseExpression(baseType, paddingField, paddingField.paddingCondition, baseType.parserArguments)}));
<#if indentContent>  </#if>    if(_needsPadding) {
<#if indentContent>  </#if>      // Just read the padding data and ignore it
<#if indentContent>  </#if>      ${helper.getReadBufferReadMethodCall(paddingField.type)};
<#if indentContent>  </#if>    }
<#if indentContent>  </#if>  }
            <#break>
        <#case "reserved">
            <#assign reservedField = field>
            <#assign simpleTypeReference = reservedField.type>

<#if indentContent>  </#if>  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
<#if indentContent>  </#if>  {
<#if indentContent>  </#if>    ${helper.getLanguageTypeNameForField(field)} _reserved = ${helper.getReadBufferReadMethodCall(reservedField.type)};
<#if indentContent>  </#if>    if(_reserved != ${reservedField.referenceValue}) {
<#if indentContent>  </#if>      printf("Expected constant value '%d' but got '%d' for reserved field.", ${reservedField.referenceValue}, _reserved);
<#if indentContent>  </#if>    }
<#if indentContent>  </#if>  }
            <#break>
        <#case "simple">
            <#assign simpleField = field>

<#if indentContent>  </#if>  // Simple Field (${simpleField.name})
            <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
            <#if helper.isSimpleTypeReference(simpleField.type)>
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${simpleField.name} = ${helper.getReadBufferReadMethodCall(simpleField.type)};
            <#else>
            <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)}* ${simpleField.name};
<#if indentContent>  </#if>  plc4c_return_code _res = plc4c_${helper.getCTypeName(simpleField.type.name)}_parse(buf<#if simpleField.params?has_content>, <#list simpleField.params as parserTerm>${helper.toParseExpression(baseType, simpleField, parserTerm, baseType.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &${simpleField.name});
<#if indentContent>  </#if>  if(_res != OK) {
<#if indentContent>  </#if>    return _res;
<#if indentContent>  </#if>  }
            </#if>
<#if indentContent>  </#if>  (*_message)-><@fieldName baseType=baseType field=simpleField/> = ${simpleField.name};
            <#break>
        <#case "switch">
            <#assign switchField = field>

<#if indentContent>  </#if>  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
            <#list switchField.cases as case>
<#if indentContent>  </#if>  <#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue><#if case.discriminatorValues?size &gt; 1>(</#if>${helper.toVariableParseExpression(baseType, switchField, switchField.discriminatorExpressions[discriminatorValue?index], type.parserArguments)} == ${discriminatorValue}<#if case.discriminatorValues?size &gt; 1>)</#if><#sep> && </#sep></#list>) </#if>{ /* ${case.name} */
<#if indentContent>  </#if>    (*_message)->_type = plc4c_${helper.getCTypeName(type.name)}_type_${helper.getCTypeName(case.name)};
                <#list case.fields as caseField>
<#if indentContent>  </#if>                    <@fieldParser baseType=case field=caseField indentContent=true/>
                    <#sep >

                </#list>
<#if indentContent>  </#if>  }<#sep> else </#sep>
            </#list>
            <#break>
        <#case "virtual">
            <#assign virtualField = field>

<#if indentContent>  </#if>  // Virtual field (Just declare a local variable so we can access it in the parser)
<#if indentContent>  </#if>  ${helper.getLanguageTypeNameForField(field)} ${virtualField.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(baseType, virtualField, virtualField.valueExpression, baseType.parserArguments)});
            <#break>
    </#switch>
</#macro>
<#list type.fields as field>
<@fieldParser baseType=type field=field indentContent=false/>
</#list>

  return OK;
}

plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_serialize(plc4c_spi_write_buffer* buf, plc4c_${helper.getCTypeName(type.name)}* _message<#if helper.getSerializerArguments(type.parserArguments)?has_content>, <#list helper.getSerializerArguments(type.parserArguments) as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)} ${parserArgument.name}<#sep>, </#sep></#list></#if>) {
<#macro fieldSerializer baseType field indentContent>
    <#switch field.typeName>
        <#case "array">
            <#assign arrayField = field>

<#if indentContent>    </#if>  // Array field (${arrayField.name})
<#if indentContent>    </#if>  {
<#if indentContent>    </#if>    uint8_t itemCount = plc4c_utils_list_size(_message-><@fieldName baseType=baseType field=arrayField/>);
<#if indentContent>    </#if>    for(int curItem = 0; curItem < itemCount; curItem++) {
            <#-- When parsing simple types, there is nothing that could require the "lastItem" -->
<#if !helper.isSimpleTypeReference(arrayField.type)><#if indentContent>    </#if>      bool lastItem = curItem == (itemCount - 1);</#if>
<#if indentContent>    </#if>      ${helper.getLanguageTypeNameForTypeReference(arrayField.type)}* _value = (${helper.getLanguageTypeNameForTypeReference(arrayField.type)}*) plc4c_utils_list_get_value(_message-><@fieldName baseType=baseType field=arrayField/>, curItem);
            <#if helper.isSimpleTypeReference(arrayField.type)>
<#if indentContent>    </#if>      ${helper.getWriteBufferWriteMethodCall(arrayField.type, "*_value")};
            <#else>
<#if indentContent>    </#if>      plc4c_return_code _res = plc4c_${helper.getCTypeName(arrayField.type.name)}_serialize(buf, (void*) &_value<#if helper.getSerializerTerms(field.params)?has_content>, <#list helper.getSerializerTerms(field.params) as serializerTerm>${helper.toSerializationExpression(baseType, field, serializerTerm, baseType.parserArguments)}<#sep>, </#sep></#list></#if>);
<#if indentContent>    </#if>      if(_res != OK) {
<#if indentContent>    </#if>        return _res;
<#if indentContent>    </#if>      }
            </#if>
<#if indentContent>    </#if>    }
<#if indentContent>    </#if>  }
            <#break>
        <#--case "checksum">
            <#assign checksumField = field>
            <#assign simpleTypeReference = checksumField.type>

<#if indentContent>    </#if>  // Checksum Field (${checksumField.name})
<#if indentContent>    </#if>  {
<#if indentContent>    </#if>    // Create an array of all the bytes read in this message element so far.
<#if indentContent>    </#if>    byte[] checksumRawData = plc4c_spi_read_get_bytes(buf, startPos, plc4c_spi_read_get_pos(buf));
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForField(field)} _checksumRef = ${helper.getReadBufferReadMethodCall(checksumField.type)};
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForField(field)} _checksum = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(baseType, checksumField, checksumField.checksumExpression, baseType.parserArguments)});
<#if indentContent>    </#if>    if(_checksum != _checksumRef) {
<#if indentContent>    </#if>      return PARSE_ERROR;
<#if indentContent>    </#if>      // throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X",_checksumRef & 0xFFFF, _checksum & 0xFFFF));
<#if indentContent>    </#if>    }
<#if indentContent>    </#if>  }
            <#break-->
        <#case "const">
            <#assign constField = field>
            <#assign simpleTypeReference = constField.type>

<#if indentContent>    </#if>  // Const Field (${constField.name})
<#if indentContent>    </#if>  ${helper.getWriteBufferWriteMethodCall(constField.type, helper.getCTypeName(baseType.name)?upper_case + "_" + helper.camelCaseToSnakeCase(constField.name)?upper_case)};
            <#break>
        <#case "discriminator">
            <#assign discriminatorField = field>
            <#assign simpleTypeReference = discriminatorField.type>

<#if indentContent>    </#if>  // Discriminator Field (${discriminatorField.name})
<#if indentContent>    </#if>  ${helper.getWriteBufferWriteMethodCall(discriminatorField.type, "plc4c_" + helper.getCTypeName(baseType.name) + "_get_discriminator(_message->_type)." + discriminatorField.name)};
            <#break>
        <#case "enum">
            <#assign enumField = field>

<#if indentContent>    </#if>  // Enum field (${enumField.name})
<#if indentContent>    </#if>  {
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForTypeReference(helper.getEnumBaseTypeReference(enumField.type))} _value = _message-><@fieldName baseType=baseType field=enumField/>;
<#if indentContent>    </#if>    ${helper.getWriteBufferWriteMethodCall(helper.getEnumBaseTypeReference(enumField.type), "_value")};
<#if indentContent>    </#if>  }
            <#break>
        <#case "implicit">
            <#assign implicitField = field>
            <#assign simpleTypeReference = implicitField.type>

<#if indentContent>    </#if>  // Implicit Field (${implicitField.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
<#if indentContent>    </#if>  ${helper.getWriteBufferWriteMethodCall(implicitField.type, helper.toSerializationExpression(baseType, implicitField, implicitField.serializeExpression, baseType.parserArguments))};
            <#break>
        <#case "manualArray">
            <#assign manualArrayField = field>

            <#-- TODO: Implement -->
            <#break>
        <#case "manual">
            <#assign manualField = field>

<#if indentContent>    </#if>  // Manual Field (${manualField.name})
<#if indentContent>    </#if>  {
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForTypeReference(manualField.type)}* _value = (${helper.getLanguageTypeNameForTypeReference(arrayField.type)}*) plc4c_utils_list_get_value(_message-><@fieldName baseType=baseType field=arrayField/>, curItem);
            <#if helper.isSimpleTypeReference(manualField.type)>
<#if indentContent>    </#if>    ${helper.getWriteBufferWriteMethodCall(manualField.type, "*_value")};
            <#else>
<#if indentContent>    </#if>    plc4c_return_code _res = plc4c_${helper.getCTypeName(manualField.type.name)}_serialize(buf, (void*) &_value);
<#if indentContent>    </#if>    if(_res != OK) {
<#if indentContent>    </#if>      return _res;
<#if indentContent>    </#if>    }
            </#if>
<#if indentContent>    </#if>  }
            <#break>
        <#case "optional">
            <#assign optionalField = field>

<#if indentContent>    </#if>  // Optional Field (${optionalField.name})
<#if indentContent>    </#if>  if(_message-><@fieldName baseType=baseType field=optionalField/> != NULL) {
            <#if helper.isSimpleTypeReference(optionalField.type)>
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForTypeReference(optionalField.type)}* _value = _message-><@fieldName baseType=baseType field=optionalField/>;
<#if indentContent>    </#if>    ${helper.getWriteBufferWriteMethodCall(optionalField.type, "*_value")};
            <#else>
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForTypeReference(optionalField.type)}* _value = _message-><@fieldName baseType=baseType field=optionalField/>;
<#if indentContent>    </#if>    plc4c_return_code _res = plc4c_${helper.getCTypeName(optionalField.type.name)}_serialize(buf, _value);
<#if indentContent>    </#if>    if(_res != OK) {
<#if indentContent>    </#if>      return _res;
<#if indentContent>    </#if>    }
            </#if>
<#if indentContent>    </#if>  }
            <#break>
        <#case "padding">
            <#assign paddingField = field>
            <#assign simpleTypeReference = paddingField.type>

<#if indentContent>    </#if>  // Padding Field (padding)
<#if indentContent>    </#if>  {
    <#-- We're replacing the "lastItem" with 'false' here as the item itself can't know if it is the last -->
<#if indentContent>    </#if>    bool _needsPadding = (bool) (${helper.toSerializationExpression(baseType, paddingField, paddingField.paddingCondition, baseType.parserArguments)});
<#if indentContent>    </#if>    if(_needsPadding) {
<#if indentContent>    </#if>      // Just output the default padding data
<#if indentContent>    </#if>      ${helper.getWriteBufferWriteMethodCall(paddingField.type, helper.toParseExpression(baseType, paddingField, paddingField.paddingValue, baseType.parserArguments))};
<#if indentContent>    </#if>    }
<#if indentContent>    </#if>  }
            <#break>
        <#case "reserved">
            <#assign reservedField = field>
            <#assign simpleTypeReference = reservedField.type>

<#if indentContent>    </#if>  // Reserved Field
<#if indentContent>    </#if>  ${helper.getWriteBufferWriteMethodCall(reservedField.type, reservedField.referenceValue)};
            <#break>
        <#case "simple">
            <#assign simpleField = field>

<#if indentContent>    </#if>  // Simple Field (${simpleField.name})
<#if indentContent>    </#if>  {
            <#if helper.isSimpleTypeReference(simpleField.type)>
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForTypeReference(simpleField.type)} _value = _message-><@fieldName baseType=baseType field=simpleField/>;
<#if indentContent>    </#if>    ${helper.getWriteBufferWriteMethodCall(simpleField.type, "_value")};
            <#else>
<#if indentContent>    </#if>    ${helper.getLanguageTypeNameForTypeReference(simpleField.type)}* _value = _message-><@fieldName baseType=baseType field=simpleField/>;
<#if indentContent>    </#if>    plc4c_return_code _res = plc4c_${helper.getCTypeName(simpleField.type.name)}_serialize(buf, _value);
<#if indentContent>    </#if>    if(_res != OK) {
<#if indentContent>    </#if>      return _res;
<#if indentContent>    </#if>    }
            </#if>
<#if indentContent>    </#if>  }
            <#break>
        <#case "switch">
            <#assign switchField = field>

<#if indentContent>    </#if>  // Switch Field (Depending of the current type, serialize the sub-type elements)
<#if indentContent>    </#if>  switch(_message->_type) {
        <#list switchField.cases as case>
<#if indentContent>    </#if>    case plc4c_${helper.getCTypeName(type.name)}_type_${helper.getCTypeName(case.name)}: {
            <#list case.fields as caseField>
                <@fieldSerializer baseType=case field=caseField indentContent=true/>
            </#list>

<#if indentContent>    </#if>      break;
<#if indentContent>    </#if>    }
        </#list>
<#if indentContent>    </#if>  }
            <#break>
    </#switch>
</#macro>
<#list type.fields as field>
  <@fieldSerializer baseType=type field=field indentContent=false/>
</#list>

  return OK;
}

uint8_t plc4c_${helper.getCTypeName(type.name)}_length_in_bytes(plc4c_${helper.getCTypeName(type.name)}* _message) {
  return plc4c_${helper.getCTypeName(type.name)}_length_in_bits(_message) / 8;
}

uint8_t plc4c_${helper.getCTypeName(type.name)}_length_in_bits(plc4c_${helper.getCTypeName(type.name)}* _message) {
    <#macro fieldSize baseType field indentContent>
        <#switch field.typeName>
            <#case "array">
                <#assign arrayField = field>

<#if indentContent>    </#if>  // Array field
                <#if helper.isSimpleTypeReference(arrayField.type)>
                    <#assign simpleTypeReference = arrayField.type>
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits} * plc4c_utils_list_size(_message-><@fieldName baseType=baseType field=arrayField/>);
                <#else>
<#if indentContent>    </#if>  if(_message-><@fieldName baseType=baseType field=arrayField/> != NULL) {
<#if indentContent>    </#if>    plc4c_list_element* curElement = _message-><@fieldName baseType=baseType field=arrayField/>->head;
<#if indentContent>    </#if>    while (curElement != NULL) {
<#if indentContent>    </#if>      lengthInBits += ${helper.getLengthInBitsFunctionNameForComplexTypedField(arrayField)}((${helper.getLanguageTypeNameForTypeReference(arrayField.type)}*) curElement->value);
<#if indentContent>    </#if>      curElement = curElement->next;
<#if indentContent>    </#if>    }
<#if indentContent>    </#if>  }
                </#if>
                <#break>
            <#case "checksum">
                <#assign checksumField = field>
                <#assign simpleTypeReference = checksumField.type>

<#if indentContent>    </#if>  // Checksum Field (checksum)
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits};
                <#break>
            <#case "const">
                <#assign constField = field>
                <#assign simpleTypeReference = constField.type>

<#if indentContent>    </#if>  // Const Field (${constField.name})
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits};
                <#break>
            <#case "discriminator">
                <#assign discriminatorField = field>
                <#assign simpleTypeReference = discriminatorField.type>

<#if indentContent>    </#if>  // Discriminator Field (${discriminatorField.name})
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits};
                <#break>
            <#case "enum">
                <#assign enumField = field>

<#if indentContent>    </#if>  // Enum Field (${enumField.name})
<#if indentContent>    </#if>  lengthInBits += ${helper.getEnumBaseTypeReference(enumField.type).sizeInBits};
                <#break>
            <#case "implicit">
                <#assign implicitField = field>
                <#assign simpleTypeReference = implicitField.type>

<#if indentContent>    </#if>  // Implicit Field (${implicitField.name})
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits};
                <#break>
            <#case "manualArray">
                <#assign manualArrayField = field>

<#if indentContent>    </#if>  // Manual Array Field (${manualArrayField.name})
<#if indentContent>    </#if>  lengthInBits += ${helper.toParseExpression(type, manualArrayField, manualArrayField.lengthExpression, type.parserArguments)} * 8;
                <#break>
            <#case "manual">
                <#assign manualField = field>

<#if indentContent>    </#if>  // Manual Field (${manualField.name})
<#if indentContent>    </#if>  lengthInBits += ${helper.toParseExpression(type, manualField, manualField.lengthExpression, type.parserArguments)} * 8;
                <#break>
            <#case "optional">
                <#assign optionalField = field>

<#if indentContent>    </#if>  // Optional Field (${optionalField.name})
<#if indentContent>    </#if>  if(_message-><@fieldName baseType=baseType field=optionalField/> != NULL) {
                <#if helper.isSimpleTypeReference(optionalField.type)>
                    <#assign simpleTypeReference = optionalField.type>
<#if indentContent>    </#if>    lengthInBits += ${simpleTypeReference.sizeInBits};
                <#else>
<#if indentContent>    </#if>    lengthInBits += ${helper.getLengthInBitsFunctionNameForComplexTypedField(optionalField)}(_message-><@fieldName baseType=baseType field=optionalField/>);
                </#if>
<#if indentContent>    </#if>  }
                <#break>
            <#case "padding">
                <#assign paddingField = field>
                <#assign simpleTypeReference = paddingField.type>

<#if indentContent>    </#if>  // Padding Field (padding)
            <#-- We're replacing the "lastItem" with 'false' here as the item itself can't know if it is the last -->
<#if indentContent>    </#if>  if((bool) (${helper.toSerializationExpression(type, paddingField, paddingField.paddingCondition, type.parserArguments)?replace("lastItem", "false")})) {
<#if indentContent>    </#if>    lengthInBits += ${simpleTypeReference.sizeInBits};
<#if indentContent>    </#if>  }
                <#break>
            <#case "reserved">
                <#assign reservedField = field>
                <#assign simpleTypeReference = reservedField.type>

<#if indentContent>    </#if>  // Reserved Field (reserved)
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits};
                <#break>
            <#case "simple">
                <#assign simpleField = field>

<#if indentContent>    </#if>  // Simple field (${simpleField.name})
                <#if helper.isSimpleTypeReference(simpleField.type)>
                    <#assign simpleTypeReference = simpleField.type>
<#if indentContent>    </#if>  lengthInBits += ${simpleTypeReference.sizeInBits};
                <#else>
<#if indentContent>    </#if>  lengthInBits += ${helper.getLengthInBitsFunctionNameForComplexTypedField(simpleField)}(_message-><@fieldName baseType=baseType field=simpleField/>);
                </#if>
                <#break>
            <#case "switch">
                <#assign switchField = field>

<#if indentContent>    </#if>  // Depending of the current type, add the length of sub-type elements ...
<#if indentContent>    </#if>  switch(_message->_type) {
            <#list switchField.cases as case>
<#if indentContent>    </#if>    case plc4c_${helper.getCTypeName(type.name)}_type_${helper.getCTypeName(case.name)}: {
                <#list case.fields as caseField>
                    <@fieldSize baseType=case field=caseField indentContent=true/>
                    <#sep >

                </#list>

<#if indentContent>    </#if>      break;
<#if indentContent>    </#if>    }
            </#list>
<#if indentContent>    </#if>  }
                <#break>
            <#case "virtual">
                <#assign virtualField = field>

<#if indentContent>    </#if>  // A virtual field doesn't have any in- or output.
                <#break>
        </#switch>
    </#macro>
  uint8_t lengthInBits = 0;
<#list type.fields as field>
  <@fieldSize baseType=type field=field indentContent=false/>
</#list>

  return lengthInBits;
}

</#if>
</#outputformat>