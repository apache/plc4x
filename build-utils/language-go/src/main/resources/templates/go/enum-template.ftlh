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
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/types/${type.name}.java
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

package ${helper.packageName(protocolName, languageName, outputFlavor)}.types;

import org.apache.plc4x.java.spi.generation.Message;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public enum ${type.name} {

<#list type.enumValues as enumValue>
    ${enumValue.name}(<#if type.type?has_content>(${helper.getLanguageTypeNameForTypeReference(type.type, true)}) ${enumValue.value}</#if><#if type.constantNames?has_content><#if type.type?has_content>, </#if><#list type.constantNames as constantName>(${helper.getLanguageTypeNameForTypeReference(type.getConstantType(constantName), true)}) ${helper.escapeValue(type.getConstantType(constantName), enumValue.getConstant(constantName))}<#sep>, </#sep></#list></#if>)<#sep>,
</#sep></#list>;

<#if type.type?has_content>
    private static final Logger logger = LoggerFactory.getLogger(${type.name}.class);

    private static final Map<${helper.getLanguageTypeNameForTypeReference(type.type, false)}, ${type.name}> map;
    static {
        map = new HashMap<>();
        for (${type.name} value : ${type.name}.values()) {
            map.put((${helper.getLanguageTypeNameForTypeReference(type.type, true)}) value.getValue(), value);
        }
    }

    private ${helper.getLanguageTypeNameForTypeReference(type.type, true)} value;
</#if>
<#if type.constantNames?has_content>
    <#list type.constantNames as constantName>
        private ${helper.getLanguageTypeNameForTypeReference(type.getConstantType(constantName), true)} ${constantName};
    </#list>
</#if>

    ${type.name}(<#if type.type?has_content>${helper.getLanguageTypeNameForTypeReference(type.type, true)} value</#if><#if type.constantNames?has_content><#if type.type?has_content>, </#if><#list type.constantNames as constantName>${helper.getLanguageTypeNameForTypeReference(type.getConstantType(constantName), true)} ${constantName}<#sep>, </#sep></#list></#if>) {
<#if type.type?has_content>        this.value = value;</#if>
<#if type.constantNames?has_content>
    <#list type.constantNames as constantName>
        this.${constantName} = ${constantName};
    </#list>
</#if>
    }

<#if type.type?has_content>
    public ${helper.getLanguageTypeNameForTypeReference(type.type, true)} getValue() {
        return value;
    }
</#if>

<#if type.constantNames?has_content>
    <#list type.constantNames as constantName>
    public ${helper.getLanguageTypeNameForTypeReference(type.getConstantType(constantName), true)} get${constantName?cap_first}() {
        return ${constantName};
    }

    </#list>
</#if>
<#if type.type?has_content>
    public static ${type.name} valueOf(${helper.getLanguageTypeNameForTypeReference(type.type, true)} value) {
        if (!map.containsKey(value)) {
            logger.error("No ${type.name} for value {}", value);
        }
        return map.get(value);
    }
</#if>

}
</#outputformat>