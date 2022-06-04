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
<#-- @ftlvariable name="helper" type="org.apache.plc4x.language.rust.RustLanguageTemplateHelper" -->
<#-- @ftlvariable name="tracer" type="org.apache.plc4x.plugins.codegenerator.protocol.freemarker.Tracer" -->
<#-- @ftlvariable name="type" type="org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition" -->
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.rs
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
// package ${helper.packageName(protocolName, languageName, outputFlavor)};
use crate::plc4x_enum;
use std::io::{Error, ErrorKind, Read, Write};
use crate::{Message, NoOption, ReadBuffer, WriteBuffer};

<#list helper.generateImports(type) as import>
<#if type.name != import>
<#--  File -> Struct  -->
use crate::${import}::${import};
</#if>
</#list>

<#if type.parserArguments.isEmpty()>
plc4x_enum!
[enum <#if helper.isExactBitLength(type.type)>${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}<#else>${helper.getExactBitLength(type.type)} => ${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}</#if> : ${type.name}
    <#list type.enumValues as enumValue>
        [${enumValue.value} => ${enumValue.name}]<#sep>, </#sep>
    </#list>
];
<#else>
pub struct ${type.name}Arguments {
    <#list type.parserArguments.orElseThrow() as arg>
        ${helper.sink(arg)}${arg.name}: ${helper.getNullsafeLanguageTypeNameForTypeReference(arg, type.enumValues)}<#sep>, </#sep>
    </#list>
}

plc4x_enum!
[enum <#if helper.isExactBitLength(type.type)>${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}<#else>${helper.getExactBitLength(type.type)} => ${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}</#if> : ${type.name}
<#list type.enumValues as enumValue>
    [${enumValue.value} => ${enumValue.name}]<#sep>, </#sep>
</#list>
];

impl ${type.name} {

    fn get_arguments(self) -> ${type.name}Arguments {
        match self {
            <#list type.enumValues as enumValue>
                ${enumValue.name} => {
                    ${type.name}Arguments {
                        <#list type.parserArguments.orElseThrow() as constant>
                            ${constant.name}: ${helper.formatLiteral(constant, type.enumValues, enumValue.getConstant(constant.name).orElseThrow())}<#sep>,</#sep>
                        </#list>
                    }
                }<#sep>, </#sep>
            </#list>
        }
<#--        match value {-->
<#--            <#list type.enumValues as val>-->
<#--                (<#list type.parserArguments.orElseThrow() as constant>${val.getConstant(constant.name).orElseThrow()} <#sep>,</#sep></#list>) => {}-->
<#--            </#list>-->
<#--            ${helper.sink(type)}-->
<#--            <#list type.parserArguments as k, v>${k}<#sep>,</#sep></#list>-->
    }
}
</#if>
</#outputformat>
