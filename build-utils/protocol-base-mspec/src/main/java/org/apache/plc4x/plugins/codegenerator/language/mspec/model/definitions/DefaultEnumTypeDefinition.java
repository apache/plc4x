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

package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;

import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;

import java.util.List;

public class DefaultEnumTypeDefinition extends DefaultTypeDefinition implements EnumTypeDefinition {

    private final TypeDefinition baseType;
    private final List<EnumTypeDefinition.EnumValue> values;

    public DefaultEnumTypeDefinition(String name, Argument[] parserArguments, String[] tags, TypeDefinition baseType, List<EnumValue> values) {
        super(name, parserArguments, tags);
        this.baseType = baseType;
        this.values = values;
    }

    public TypeDefinition getBaseType() {
        return baseType;
    }

    public List<EnumValue> getValues() {
        return values;
    }

}
