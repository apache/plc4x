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
${helper.setConstants(type, protocolName, outputFlavor)}${helper.getIncludesDirectory()?replace(".", "/")}/${helper.camelCaseToSnakeCase(typeName)}.h
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
#ifndef PLC4C_${helper.getCTypeName(type.name)?upper_case}_H_
#define PLC4C_${helper.getCTypeName(type.name)?upper_case}_H_
#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
#include <stdint.h>
#include <plc4c/utils/list.h>
<#if helper.getComplexTypeReferencesInFields()?has_content>
<#list helper.getComplexTypeReferencesInFields() as complexType>
#include "${helper.camelCaseToSnakeCase(complexType.getName())}.h"
</#list>
</#if>
<#if helper.isDiscriminatedType(type)>
#include "${helper.camelCaseToSnakeCase(type.parentType.name)}.h"
</#if>
        <#--
            Abstract types only occur as parents in the case of discriminated types.
            As C doesn't have any form of polymorph type-system, we need to save the type information in the
            data-structure. So with this enum we're doing exactly this. It will be used in the serializers and
            parsers and stored in discriminated types.
        -->
<#if helper.isAbstractType(type)>

// Enum assigning each sub-type an individual id.
enum plc4c_${helper.getCTypeName(type.name)}_type {
    <#list helper.getDiscriminatedSubTypes(type) as subtype>
  plc4c_${helper.getCTypeName(type.name)}_type_${helper.getCTypeName(subtype.name)} = ${subtype?index}<#sep>,
    </#list>
};
typedef enum plc4c_${helper.getCTypeName(type.name)}_type plc4c_${helper.getCTypeName(type.name)}_type;
</#if>
        <#--
            When using const fields, output the constant reference values
            as defines so we can use them elsewhere
        -->
<#if type.constFields?has_content>

// Constant values.
    <#list type.constFields as field>
const ${helper.getLanguageTypeNameForField(field)} ${helper.getCTypeName(type.name)?upper_case}_${helper.camelCaseToSnakeCase(field.name)?upper_case} = ${field.referenceValue};
    </#list>
</#if>

struct plc4c_${helper.getCTypeName(type.name)} {
<#if helper.isAbstractType(type)>
  plc4c_${helper.getCTypeName(type.name)}_type _type;
</#if>
<#if helper.isDiscriminatedType(type)>
  plc4c_${helper.getCTypeName(type.parentType.name)}_type _type;
</#if>
<#list type.allPropertyFields as field>
  ${helper.getLanguageTypeNameForField(field)}<#if !helper.isSimpleType(field.type)>*</#if> ${helper.camelCaseToSnakeCase(field.name)}${helper.getTypeSizeForField(field)}<#if field.loopType??>${helper.getLoopExpressionSuffix(field)}</#if>;
</#list>
};
typedef struct plc4c_${helper.getCTypeName(type.name)} plc4c_${helper.getCTypeName(type.name)};

plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_parse(plc4c_spi_read_buffer* buf, <#if type.parserArguments?has_content><#list type.parserArguments as parserArgument>${helper.getLanguageTypeName(parserArgument.type)} ${parserArgument.name}<#sep>, </#list>, </#if>plc4c_${helper.getCTypeName(type.name)}** message);

plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_serialize(plc4c_spi_write_buffer* buf, plc4c_${helper.getCTypeName(type.name)}* message);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_${helper.getCTypeName(type.name)?upper_case}_H_
