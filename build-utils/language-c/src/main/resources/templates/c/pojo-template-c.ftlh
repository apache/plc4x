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
${helper.setConstants(type, protocolName, outputFlavor)}${helper.getSourceDirectory()?replace(".", "/")}/${helper.camelCaseToSnakeCase(typeName)}.c
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

#include <plc4c/spi/read_buffer.h>
#include <plc4c/spi/write_buffer.h>
#include <plc4c/spi/evaluation_helper.h>

#include "${helper.camelCaseToSnakeCase(typeName)}.h"

plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_parse(plc4c_read_buffer buf, <#if type.parserArguments?has_content><#list type.parserArguments as parserArgument>${helper.getLanguageTypeName(parserArgument.type)} ${parserArgument.name}<#sep>, </#list>, </#if>plc4c_${helper.getCTypeName(type.name)}** message) {
  uint16_t start_pos = plc4c_spi_read_get_pos(buf);
  uint16_t cur_pos;

  plc4c_${helper.getCTypeName(type.name)}* msg = malloc(sizeof(plc4c_${helper.getCTypeName(type.name)}));
<#list type.fields as field>
    <#switch field.typeName>
        <#case "checksum">

  // Checksum Field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getNullValueForType(field.type)};
  {
    // Create an array of all the bytes read in this message element so far.
    byte[] checksumRawData = io.getBytes(startPos, io.getPos());
    ${helper.getLanguageTypeNameForField(field)} _${field.name}Ref = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
    ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field, field.checksumExpression, type.parserArguments)});
    if(${field.name} != _${field.name}Ref) {
      throw new ParseException(String.format("Checksum verification failed. Expected %04X but got %04X",_${field.name}Ref & 0xFFFF, ${field.name} & 0xFFFF));
    }
  }
            <#break>
        <#case "const">

  // Const Field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
  if(${field.name} != ${typeName}.${field.name?upper_case}) {
    throw new ParseException("Expected constant value " + ${typeName}.${field.name?upper_case} + " but got " + ${field.name});
  }
            <#break>
        <#case "enum">

  // Enum field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getLanguageTypeNameForField(field)}.valueOf(${helper.getReadBufferReadMethodCall(helper.getEnumBaseType(field.type))?no_esc});
            <#break>
        <#case "discriminator">

  // Discriminator Field (${field.name}) (Used as input to a switch field)
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
            <#break>
        <#case "implicit">

  // Implicit Field (${field.name}) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
            <#break>
        <#case "manual">

  // Manual Field (${field.name})
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field, field.parseExpression, type.parserArguments)});
  msg.${helper.camelCaseToSnakeCase(field.name)} = ${field.name};
        <#break>
        <#case "optional">

  // Optional Field (${field.name}) (Can be skipped, if a given expression evaluates to false)
            <#if field.conditionExpression.contains("curPos")>
  curPos = io.getPos() - startPos;
            </#if>
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = NULL;
  if(${helper.toParseExpression(field, field.conditionExpression, type.parserArguments)?no_esc}) {
        <#if helper.isSimpleType(field.type)>
    ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
        <#else>
    ${field.name} = plc4c_${helper.getCTypeName(field.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(field, parserTerm, type.parserArguments)}<#sep>, </#sep></#list></#if>);
        </#if>
  }
            <#break>
        <#case "padding">

  // Padding Field (${field.name})
  bool _${field.name}NeedsPadding = (bool) ((plc4c_spi_read_has_more(buf, ${helper.getNumBits(field.type)})) && (${helper.toParseExpression(field, field.paddingCondition, type.parserArguments)?no_esc}));
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
      LOGGER.info("Expected constant value " + ${field.referenceValue} + " but got " + reserved + " for reserved field.");
    }
  }
            <#break>
        <#case "simple">

  // Simple Field (${field.name})
            <#-- Inizialize a local variable with the simple type (Intentionally keeping the java-style names so they can be used in expressions) -->
            <#if helper.isSimpleType(field.type)>
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = ${helper.getReadBufferReadMethodCall(field.type)?no_esc};
            <#else>
            <#-- Inizialize a local variable with the complex type (Intentionally keeping the java-style names so they can be used in expressions) -->
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = plc4c_${helper.getCTypeName(field.type.name)}_parse(buf<#if field.params?has_content>, <#list field.params as parserTerm>${helper.toParseExpression(field, parserTerm, type.parserArguments)}<#sep>, </#sep></#list></#if>);
            </#if>
            <#-- Update the protpery in the struct -->
  msg.${helper.camelCaseToSnakeCase(field.name)} = ${field.name};
            <#break>
        <#case "switch">

  // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
            <#list field.cases as case>
  <#if case.discriminatorValues?has_content>if(<#list case.discriminatorValues as discriminatorValue>plc4c_spi_evaluation_helper_equals(${helper.toSwitchExpression(field.discriminatorNames[discriminatorValue?index])?no_esc}, ${discriminatorValue})<#sep> && </#sep></#list>) </#if>{
    plc4c_${helper.getCTypeName(case.name)}_parse(buf, msg<#if case.parserArguments?has_content>, <#list case.parserArguments as parserArgument>${parserArgument.name}<#sep>, </#sep></#list></#if>);
  }<#sep> else </#sep>
            </#list>
            <#break>
        <#case "virtual">

  // Virtual field (Just declare a local variable so we can access it in the parser)
  ${helper.getLanguageTypeNameForField(field)} ${field.name} = (${helper.getLanguageTypeNameForField(field)}) (${helper.toParseExpression(field, field.valueExpression, type.parserArguments)});
            <#break>
    </#switch>
</#list>

  return OK;
}

plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_serialize(plc4c_write_buffer buf, plc4c_${helper.getCTypeName(type.name)}* message) {
  return OK;
}
