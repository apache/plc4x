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

package org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields;


import org.apache.plc4x.plugins.codegenerator.types.definitions.DiscriminatedComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.fields.SwitchField;

import java.util.LinkedList;
import java.util.List;

public class DefaultSwitchField implements SwitchField {

    private final String[] discriminatorNames;
    private final List<DiscriminatedComplexTypeDefinition> cases;

    public DefaultSwitchField(String[] discriminatorNames) {
        this.discriminatorNames = discriminatorNames;
        cases = new LinkedList<>();
    }

    public String[] getDiscriminatorNames() {
        return discriminatorNames;
    }

    public void addCase(DiscriminatedComplexTypeDefinition caseType) {
        cases.add(caseType);
    }

    public List<DiscriminatedComplexTypeDefinition> getCases() {
        return cases;
    }

    public String[] getParams() {
        return new String[0];
    }

}
