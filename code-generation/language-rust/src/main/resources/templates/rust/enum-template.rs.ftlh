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
src/${type.name}.rs
<#--${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/${type.name}.rs-->
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
<#--use crate::plc4x_enum;-->
use std::io::{Error, ErrorKind, Read, Write};
use plc4rust::{Message, NoOption};
use plc4rust::read_buffer::ReadBuffer;
use plc4rust::write_buffer::WriteBuffer;

<#list helper.generateImports(type) as import>
<#if type.name != import[import?size - 1]>
<#--  File -> Struct  -->
use crate::<#list import as path>${path}<#sep>::</#sep></#list>;
</#if>
</#list>

<#--
[enum <#if helper.isExactBitLength(type.type)>${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}<#else>${helper.getExactBitLength(type.type)} => ${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}</#if> : ${type.name}
    <#list type.enumValues as enumValue>
        [${enumValue.value} => ${enumValue.name}]<#sep>, </#sep>
    </#list>
];
-->
#[derive(Copy, Clone, PartialEq, Debug)]
#[allow(non_camel_case_types)]
pub enum ${type.name} {
<#list type.enumValues as enumValue>
    ${enumValue.name}<#sep>, </#sep>
</#list>
}

impl TryFrom<${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}> for ${type.name} {
    type Error = ();

    fn try_from(value: ${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}) -> Result<Self, Self::Error> {
        match value {
            <#list type.enumValues as enumValue>
            ${enumValue.value} => Ok(${type.name}::${enumValue.name}),
            </#list>
            _ => {
                panic!("Unable to deserialize enum!")
            }
        }
    }
}

impl Into<${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}> for ${type.name} {
    fn into(self) -> ${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)} {
        match self {
        <#list type.enumValues as enumValue>
            ${enumValue.name} => ${enumValue.value}<#sep>,</#sep>
        </#list>
        }
    }
}

impl Message for ${type.name} {
    type M = ${type.name};
    type P = NoOption;

    fn get_length_in_bits(&self) -> u32 {
        todo!()
    }

    fn serialize<T: Write>(&self, writer: &mut WriteBuffer<T>) -> Result<usize, Error> {
        writer.write_${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}((*self).into())
    }

    fn parse<T: Read>(reader: &mut ReadBuffer<T>, parameter: Option<Self::P>) -> Result<Self::M, Error> {
        assert!(parameter.is_none());
        let result = reader.read_${helper.getLanguageTypeNameForTypeReference(type.type.orElseThrow(), true)}()?;
        match ${type.name}::try_from(result) {
            Ok(result) => {
                Ok(result)
            }
            Err(_) => {
                panic!("Cannot parse {}", result);
            }
        }
    }
}

<#if !type.parserArguments.isEmpty()>
pub struct ${type.name}Arguments {
<#list type.parserArguments.orElseThrow() as arg>
    ${arg.name}: ${helper.getNullsafeLanguageTypeNameForTypeReference(arg, type.enumValues)}<#sep>, </#sep>
</#list>
}

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
    }
}
</#if>
</#outputformat>
