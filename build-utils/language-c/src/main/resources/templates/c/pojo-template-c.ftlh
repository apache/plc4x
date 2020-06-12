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
<#if !helper.isDiscriminatedChildTypeDefinition(type)>${helper.getSourceDirectory()?replace(".", "/")}/${helper.camelCaseToSnakeCase(typeName)}.c
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
plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_parse(plc4c_spi_read_buffer* buf, <#if type.parserArguments?has_content><#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#if !helper.isSimpleTypeReference(parserArgument.type)>*</#if> ${parserArgument.name}<#sep>, </#list>, </#if>plc4c_${helper.getCTypeName(type.name)}** message) {
  uint16_t startPos = plc4c_spi_read_get_pos(buf);
  uint16_t curPos;

  // Pointer to the parsed data structure.
  plc4c_${helper.getCTypeName(type.name)}* msg = malloc(sizeof(plc4c_${helper.getCTypeName(type.name)}));
<#if type.getPropertyFields()?has_content>

</#if>
<#list type.fields as field>
    <#switch field.typeName>
        <#case "array">
  // Array field (${field.name})
            <#-- Only update curPos if the length expression uses it -->
            <#if field.loopExpression.contains("curPos")>
  curPos = plc4c_spi_read_get_pos(buf) - startPos;
            </#if>
            <#-- In all other cases do we have to work with a list, that is later converted to an array -->
  plc4c_list ${field.name};
  {
            <#if helper.isCountArrayField(field)>
    // Count array
    uint8_t itemCount = ${helper.toParseExpression(field.loopExpression, type.parserArguments)?no_esc};
    for(int curItem = 0; curItem < itemCount; curItem++) {
      <#if !helper.isSimpleTypeReference(field.type)>boolean lastItem = curItem == (itemCount - 1);</#if>
      <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
                <#if helper.isSimpleTypeReference(field.type)>
      ${helper.getLanguageTypeNameForTypeReference(field.type)} value = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
                <#else>
      <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
      ${helper.getLanguageTypeNameForField(field)}* value = NULL;
      plc4c_${helper.getCTypeName(field.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(parserTerm, type.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &${field.name});
                </#if>
      plc4c_utils_list_insert_head_value(&${field.name}, &value);
    }
            <#-- For a length array, we read data till the read position of the buffer reaches a given position -->
            <#elseif helper.isLengthArrayField(field)>
    // Length array
    uint8_t _${field.name}Length = ${helper.toParseExpression(field.loopExpression, type.parserArguments)?no_esc};
    uint8_t ${field.name}EndPos = plc4c_spi_read_get_pos(buf) + _${field.name}Length;
    while(plc4c_spi_read_get_pos(buf) < ${field.name}EndPos) {
      plc4c_utils_list_insert_head_value(&${field.name}, <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)?no_esc}<#else>plc4c_${helper.getCTypeName(field.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(parserTerm, type.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &${field.name})</#if>);
                <#-- After parsing, update the current position, but only if it's needed -->
                <#if field.loopExpression.contains("curPos")>
      curPos = plc4c_spi_read_get_pos(buf) - startPos;
                </#if>
    }
            <#-- A terminated array keeps on reading data as long as the termination expression evaluates to false -->
            <#elseif helper.isTerminatedArrayField(field)>
    // Terminated array
    while(!((boolean) (${helper.toParseExpression(field.loopExpression, type.parserArguments)?no_esc}))) {
      plc4c_utils_list_insert_head_value(&${field.name}, <#if helper.isSimpleTypeReference(field.type)>${helper.getReadBufferReadMethodCall(field.type)?no_esc}<#else>${field.type.name}IO.staticParse(io<#if field.params?has_content>, <#list field.params as parserArgument>(${helper.getArgumentType(field.type, parserArgument?index)}) (${helper.toParseExpression(parserArgument, type.parserArguments)?no_esc})<#sep>, </#sep></#list></#if>)</#if>);
              <#-- After parsing, update the current position, but only if it's needed -->
              <#if field.loopExpression.contains("curPos")>
      curPos = plc4c_spi_read_get_pos(buf) - startPos;
              </#if>
    }
            </#if>
  }

            <#break>
        <#case "checksum">
  // Checksum Field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getNullValueForType(field.type)};
  {
    // Create an array of all the bytes read in this message element so far.
    byte[] checksumRawData = plc4c_spi_read_get_bytes(buf, startPos, plc4c_spi_read_get_pos(buf));
    ${helper.getLanguageTypeNameForField(field)} _${field.name}Ref = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
    ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field.checksumExpression, type.parserArguments)});
    if(${field.name} != _${field.name}Ref) {
      return PARSE_ERROR;
      // throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X",_${field.name}Ref & 0xFFFF, ${field.name} & 0xFFFF));
    }
  }

            <#break>
        <#case "const">
  // Const Field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
  if(${field.name} != ${helper.getCTypeName(type.name)?upper_case}_${helper.camelCaseToSnakeCase(field.name)?upper_case}) {
    return PARSE_ERROR;
    // throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
  }

            <#break>
        <#case "enum">
  // Enum field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(helper.getEnumBaseTypeReference(field.type))?no_esc};
  msg->${helper.camelCaseToSnakeCase(field.name)} = ${field.name};

            <#break>
        <#case "discriminator">
  // Discriminator Field (${field.name}) (Used as input to a switch field)
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};

            <#break>
        <#case "implicit">
  // Implicit Field (${field.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};

            <#break>
        <#case "manualArray">
             <#-- TODO: Implement -->

            <#break>
        <#case "manual">
  // Manual Field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field.parseExpression, type.parserArguments)});
  msg->${helper.camelCaseToSnakeCase(field.name)} = ${field.name};

        <#break>
        <#case "optional">
  // Optional Field (${field.name}) (Can be skipped, if a given expression evaluates to false)
            <#if field.conditionExpression.contains("curPos")>
  curPos = plc4c_spi_read_get_pos(buf) - startPos;
            </#if>
  ${helper.getLanguageTypeNameForField(field)}<#if !helper.isSimpleTypeReference(field.type)>*</#if> ${field.name} = <#if !helper.isSimpleTypeReference(field.type)>NULL<#else>0</#if>;
  if(${helper.toParseExpression(field.conditionExpression, type.parserArguments)?no_esc}) {
        <#if helper.isSimpleTypeReference(field.type)>
    ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
        <#else>
    plc4c_${helper.getCTypeName(field.type.name)}* ${field.name} = NULL;
    plc4c_${helper.getCTypeName(field.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(parserTerm, type.parserArguments)}<#sep>, </#sep></#list></#if>, &${field.name});
        </#if>
    msg->${helper.camelCaseToSnakeCase(field.name)} = ${field.name};
  }

            <#break>
        <#case "padding">
  // Padding Field (${field.name})
  bool _${field.name}NeedsPadding = (bool) ((plc4c_spi_read_has_more(buf, ${helper.getNumBits(field.type)})) && (${helper.toParseExpression(field.paddingCondition, type.parserArguments)?no_esc}));
  if(_${field.name}NeedsPadding) {
    // Just read the padding data and ignore it
    ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
  }

            <#break>
        <#case "reserved">
  // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
  {
    ${helper.getLanguageTypeNameForField(field)} reserved = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
    if(reserved != ${helper.getReservedValue(field)}) {
      printf("Expected constant value '%d' but got '%d' for reserved field.", ${field.referenceValue}, reserved);
    }
  }

            <#break>
        <#case "simple">
  // Simple Field (${field.name})
            <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
            <#if helper.isSimpleTypeReference(field.type)>
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
            <#else>
            <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
  ${helper.getLanguageTypeNameForField(field)}* ${field.name} = NULL;
  plc4c_${helper.getCTypeName(field.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(parserTerm, type.parserArguments)}<#sep>, </#sep></#list></#if>, (void*) &${field.name});
            </#if>
  msg->${helper.camelCaseToSnakeCase(field.name)} = ${field.name};

            <#break>
        <#case "switch">
  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
            <#list field.cases as case>
  <#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue><#if case.discriminatorValues?size &gt; 1>(</#if>${helper.toVariableParseExpression(field.discriminatorExpressions[discriminatorValue?index], type.parserArguments)?no_esc}<#--${helper.toVariableParseExpressionHurz(field)}--> == ${discriminatorValue}<#if case.discriminatorValues?size &gt; 1>)</#if><#sep> && </#sep></#list>) </#if>{ /* ${case.name} */
                <#list case.propertyFields as caseField>
    ${helper.getLanguageTypeNameForField(caseField)}<#if !helper.isSimpleTypeReference(caseField.type)>*</#if> ${caseField.name}<#if field.loopType??>${helper.getLoopExpressionSuffix(caseField)}</#if><#if helper.getLanguageTypeNameForField(caseField) == "plc4c_list"><#elseif !helper.isSimpleTypeReference(caseField.type)> = NULL<#else> = -1</#if>;
    msg->${helper.camelCaseToSnakeCase(case.name)}_${helper.camelCaseToSnakeCase(caseField.name)} = ${caseField.name};
                    <#sep >

                </#list>
  }<#sep> else </#sep>
            </#list>

            <#break>
        <#case "virtual">
  // Virtual field (Just declare a local variable so we can access it in the parser)
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field.valueExpression, type.parserArguments)});

            <#break>
    </#switch>
</#list>

  return OK;
}

plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_serialize(plc4c_spi_write_buffer* buf, plc4c_${helper.getCTypeName(type.name)}* message) {
  return OK;
}
</#if>