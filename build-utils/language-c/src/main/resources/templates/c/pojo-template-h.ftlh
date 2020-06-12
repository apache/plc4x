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
<#if !helper.isDiscriminatedChildTypeDefinition(type)>${helper.getIncludesDirectory()?replace(".", "/")}/${helper.camelCaseToSnakeCase(typeName)}.h
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
<#--
    Add any import statements for partent-types, complex types used in properties or parser arguments.
-->
<#if helper.getComplexTypeReferences()?has_content>
    <#list helper.getComplexTypeReferences() as typeReference>
#include "${helper.camelCaseToSnakeCase(typeReference.name)}.h"
    </#list>
</#if>
<#--
    Abstract types only occur as parents in the case of discriminated types.
    As C doesn't have any form of polymorph type-system, we need to save the type information in the
    data-structure. So with this enum we're doing exactly this. It will be used in the serializers and
    parsers and stored in discriminated types.
-->
<#if helper.isDiscriminatedParentTypeDefinition(type)>

// Structure used to contain the discriminator values for discriminated types using this as a parent
struct plc4c_${helper.getCTypeName(type.name)}_discriminator {
<#if helper.getDiscriminatorTypes()?has_content>
    <#list helper.getDiscriminatorTypes()?keys as discriminatorName>
  ${helper.getLanguageTypeNameForTypeReference(helper.getDiscriminatorTypes()[discriminatorName])} ${discriminatorName};
    </#list>
</#if>
};
typedef struct plc4c_${helper.getCTypeName(type.name)}_discriminator plc4c_${helper.getCTypeName(type.name)}_discriminator;

// Enum assigning each sub-type an individual id.
enum plc4c_${helper.getCTypeName(type.name)}_type {
    <#list helper.getSubTypeDefinitions() as subtype>
  plc4c_${helper.getCTypeName(type.name)}_type_${helper.getCTypeName(subtype.name)} = ${subtype?index}<#sep>,
    </#list>
};
typedef enum plc4c_${helper.getCTypeName(type.name)}_type plc4c_${helper.getCTypeName(type.name)}_type;

// Function to get the discriminator values for a given type.
plc4c_${helper.getCTypeName(type.name)}_discriminator plc4c_${helper.getCTypeName(type.name)}_get_discriminator(plc4c_${helper.getCTypeName(type.name)}_type type);
</#if>
<#--
    When using const fields, output the constant reference values
    as global const values so we can use them elsewhere.
-->
<#if type.constFields?has_content>

// Constant values.
    <#list type.constFields as field>
const ${helper.getLanguageTypeNameForField(field)} ${helper.getCTypeName(type.name)?upper_case}_${helper.camelCaseToSnakeCase(field.name)?upper_case} = ${field.referenceValue};
    </#list>
</#if>

<#--
    Create the general data-structure for this type
-->
struct plc4c_${helper.getCTypeName(type.name)} {
<#if helper.isDiscriminatedParentTypeDefinition(type)>
  /* This is an abstract type so this property saves the type of this typed union */
  plc4c_${helper.getCTypeName(type.name)}_type _type;
</#if>
  /* Properties */
<#list helper.getPropertyAndSwitchFields() as field>
    <#if helper.isSwitchField(field)>
  union {
        <#list field.cases as case>
    struct { /* ${case.name} */
            <#list case.propertyFields as caseField>
      ${helper.getLanguageTypeNameForField(caseField)}<#if !helper.isSimpleTypeReference(caseField.type)>*</#if> ${helper.camelCaseToSnakeCase(case.name)}_${helper.camelCaseToSnakeCase(caseField.name)}${helper.getTypeSizeForField(caseField)}<#if field.loopType??>${helper.getLoopExpressionSuffix(caseField)}</#if>;
            </#list>
    };
        </#list>
  };
    <#elseif helper.isEnumField(field)>
  ${helper.getLanguageTypeNameForField(field)} ${helper.camelCaseToSnakeCase(field.name)};
    <#else>
  ${helper.getLanguageTypeNameForField(field)}<#if !helper.isSimpleTypeReference(field.type)>*</#if> ${helper.camelCaseToSnakeCase(field.name)}${helper.getTypeSizeForField(field)}<#if field.loopType??>${helper.getLoopExpressionSuffix(field)}</#if>;
    </#if>
</#list>
};
typedef struct plc4c_${helper.getCTypeName(type.name)} plc4c_${helper.getCTypeName(type.name)};

<#--
    Define the parse-method for elements of this tpye
-->
plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_parse(plc4c_spi_read_buffer* buf, <#if type.parserArguments?has_content><#list type.parserArguments as parserArgument>${helper.getLanguageTypeNameForTypeReference(parserArgument.type)}<#if !helper.isSimpleTypeReference(parserArgument.type)>*</#if> ${parserArgument.name}<#sep>, </#list>, </#if>plc4c_${helper.getCTypeName(type.name)}** message);

<#--
    Define the serialize-method for elements of this tpye
-->
plc4c_return_code plc4c_${helper.getCTypeName(type.name)}_serialize(plc4c_spi_write_buffer* buf, plc4c_${helper.getCTypeName(type.name)}* message);

#ifdef __cplusplus
}
#endif
#endif  // PLC4C_${helper.getCTypeName(type.name)?upper_case}_H_
</#if>