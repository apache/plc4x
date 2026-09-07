<#--
  Licensed to the Apache Software Foundation (ASF) under one
  or more contributor license agreements.  See the NOTICE file
  distributed with this work for additional information
  regarding copyright ownership.  The ASF licenses this file
  to you under the Apache License, Version 2.0 (the
  "License") you may not use this file except in compliance
  with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

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
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.java.JavaLanguageTemplateHelper" -->
<#-- @ftlvariable name="tracer" type="org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition" -->
<#-- Declare the name and type of variables declared locally inside the template -->
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.py
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

from dataclasses import dataclass

@dataclass
class ${type.name}:
    <#-- If the current type contains "const" fields, generate some Python constants for holing their values -->
    <#if type.constFields?has_content>
        <#list type.constFields as field>
    ${helper.camelCaseToSnakeCase(field.name)?upper_case}: ${helper.getLanguageTypeNameForTypeReference(field.type)} = ${helper.toParseExpression(field, field.type, field.referenceValue, parserArguments)}
        </#list>
    </#if>

</#outputformat>