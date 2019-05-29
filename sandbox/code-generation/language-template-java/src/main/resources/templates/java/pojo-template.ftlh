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
<#list type.constFields as constField>
    private static final ${helper.getLanguageTypeNameForSpecType(constField.type)} ${constField.name?upper_case} = ${constField.referenceValue};
</#list>
</#if>
<#if type.simpleFields?has_content>

    // Properties.
<#list type.simpleFields as simpleField>
    private final ${helper.getLanguageTypeNameForSpecType(simpleField.type)} ${simpleField.name};
</#list>
</#if>

    public ${typeName}(<#list type.getAllSimpleFields() as simpleField>${helper.getLanguageTypeNameForSpecType(simpleField.type)} ${simpleField.name}<#sep>, </#sep></#list>) {
<#if type.getParentSimpleFields()?has_content>
        super(<#list type.getParentSimpleFields() as simpleField>${simpleField.name}<#sep>, </#sep></#list>);
</#if>
<#list type.simpleFields as simpleField>
        this.${simpleField.name} = ${simpleField.name};
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

<#list type.simpleFields as simpleField>
    public ${helper.getLanguageTypeNameForSpecType(simpleField.type)} get${simpleField.name?cap_first}() {
        return ${simpleField.name};
    }

</#list>
}