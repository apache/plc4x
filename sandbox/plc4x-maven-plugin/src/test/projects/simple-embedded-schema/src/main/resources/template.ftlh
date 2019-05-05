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

public class ${typeName} {

<#list segments as segment>
<#if helper.isFieldSegment(segment)>
    private final ${segment.type.name.name} ${segment.name};
</#if>
</#list>

    public ${typeName}(<#list segments as segment><#if helper.isFieldSegment(segment)>${segment.type.name.name} ${segment.name}</#if><#sep>, </#sep></#list>) {
<#list segments as segment>
<#if helper.isFieldSegment(segment)>
        this.${segment.name} = ${segment.name};
</#if>
</#list>
    }

<#list segments as segment>
<#if helper.isFieldSegment(segment)>
    public ${segment.type.name.name} get${segment.name?cap_first}() {
        return ${segment.name};
    }

</#if>
</#list>
}