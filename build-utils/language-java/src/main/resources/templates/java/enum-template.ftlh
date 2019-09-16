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
${helper.packageName(protocolName, languageName, outputFlavor)?replace(".", "/")}/types/${typeName}.java
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

import org.apache.plc4x.java.utils.Message;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public enum ${typeName} {

<#list type.enumValues as enumValue>
    ${enumValue.name}((${helper.getLanguageTypeName(type.type, true)}) ${enumValue.value}<#if type.constantNames?has_content>, <#list type.constantNames as constantName>(${helper.getLanguageTypeName(type.getConstantType(constantName), true)}) ${helper.escapeValue(type.getConstantType(constantName), enumValue.getConstant(constantName))?no_esc}<#sep>, </#sep></#list></#if>)<#sep>,
</#sep></#list>;

    private static final Logger logger = LoggerFactory.getLogger(${typeName}.class);

    private static final Map<${helper.getLanguageTypeName(type.type, false)}, ${typeName}> map;
    static {
        map = new HashMap<>();
        for (${typeName} value : ${typeName}.values()) {
            map.put((${helper.getLanguageTypeName(type.type, true)}) value.getValue(), value);
        }
    }

    private ${helper.getLanguageTypeName(type.type, true)} value;
    <#if type.constantNames?has_content>
        <#list type.constantNames as constantName>
    private ${helper.getLanguageTypeName(type.getConstantType(constantName), true)} ${constantName};
        </#list>
    </#if>

    ${typeName}(${helper.getLanguageTypeName(type.type, true)} value<#if type.constantNames?has_content>, <#list type.constantNames as constantName>${helper.getLanguageTypeName(type.getConstantType(constantName), true)} ${constantName}<#sep>, </#sep></#list></#if>) {
        this.value = value;
        <#if type.constantNames?has_content>
            <#list type.constantNames as constantName>
        this.${constantName} = ${constantName};
            </#list>
        </#if>
    }

    public ${helper.getLanguageTypeName(type.type, true)} getValue() {
        return value;
    }

    <#if type.constantNames?has_content>
        <#list type.constantNames as constantName>
    public ${helper.getLanguageTypeName(type.getConstantType(constantName), true)} get${constantName?cap_first}() {
        return ${constantName};
    }

        </#list>
    </#if>
    public static ${typeName} valueOf(${helper.getLanguageTypeName(type.type, true)} value) {
        if (!map.containsKey(value)) {
            logger.error("No ${typeName} for value {}", value);
        }
        return map.get(value);
    }

}