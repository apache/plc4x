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
${packageName?replace(".", "/")}/${typeName}.java
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

package ${packageName};

public<#if type.abstract> abstract</#if> class ${typeName}<#if type.parentType??> extends ${type.parentType.name}</#if> {
<#if helper.isDiscriminatedType(type)>

    // Discriminator values used by the parser to determine the type to be used. All values have to apply.
    public static final String[] DISCRIMINATOR_VALUES = new String[] {
    <#list type.discriminatorValues as discriminatorValue>
        "${discriminatorValue}"<#sep>, </#sep>
    </#list>
    };
</#if>
<#if type.constFields?has_content>

    // Constant values.
<#list type.constFields as field>
    public static final ${helper.getLanguageTypeNameForSpecType(field.type)} ${field.name?upper_case} = ${field.referenceValue};
</#list>
</#if>
<#if type.propertyFields?has_content>

    // Properties.
<#list type.propertyFields as field>
    private final ${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> ${field.name};
</#list>
</#if>

    public ${typeName}(<#list type.getAllPropertyFields() as field>${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> ${field.name}<#sep>, </#sep></#list>) {
<#if type.getParentPropertyFields()?has_content>
        super(<#list type.getParentPropertyFields() as field>${field.name}<#sep>, </#sep></#list>);
</#if>
<#list type.propertyFields as field>
        this.${field.name} = ${field.name};
</#list>
    }
<#if type.abstract>

    public abstract String[] getDiscriminatorValues();
</#if>
<#if helper.isDiscriminatedType(type)>

    public String[] getDiscriminatorValues() {
        return DISCRIMINATOR_VALUES;
    }
</#if>

<#list type.propertyFields as field>
    public ${helper.getLanguageTypeNameForSpecType(field.type)}<#if field.lengthType??>[]</#if> get${field.name?cap_first}() {
        return ${field.name};
    }

</#list>
}