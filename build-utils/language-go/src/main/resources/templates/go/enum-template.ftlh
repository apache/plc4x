<#--
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
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition" -->
${helper.fileName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.go
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
package ${helper.packageName(outputFlavor)}

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"

type ${type.name} <#if type.type?has_content>${helper.getLanguageTypeNameForTypeReference(type.type)}<#else>string</#if>

const(
<#list type.enumValues as enumValue>
    ${enumValue.name} ${type.name} = ${enumValue.value}
</#list>
)

<#--
If this enum also defines constant values, output one function for each.
Each of these functions consist of a simple switch statement using the enum as switch-value
returning the constant for a given enum value.
-->
<#if type.constantNames?has_content>
    <#list type.constantNames as constantName>

func (e ${type.name}) Get${constantName?cap_first}() ${helper.getLanguageTypeNameForTypeReference(type.getConstantType(constantName))} {
    switch e  {
        <#list helper.getUniqueEnumValues(type.enumValues) as enumValue>
        case ${helper.escapeValue(type.type, enumValue.value)}: { /* '${enumValue.value}' */
            return ${helper.escapeEnumValue(type.getConstantType(constantName), enumValue.getConstant(constantName))}
        }
        </#list>
        default: {
            return 0
        }
    }
}
    </#list>
</#if>

<#if type.type?has_content>
func (e *${type.name}) Parse(io spi.ReadBuffer) {
    // TODO: Implement ...
}

func (e ${type.name}) Serialize(io spi.WriteBuffer) {
    // TODO: Implement ...
}
</#if>
</#outputformat>