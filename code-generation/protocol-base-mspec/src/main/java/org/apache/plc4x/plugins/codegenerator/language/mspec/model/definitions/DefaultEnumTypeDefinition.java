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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.definitions;

import org.apache.plc4x.plugins.codegenerator.types.definitions.Argument;
import org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.enums.EnumValue;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;

public class DefaultEnumTypeDefinition extends DefaultTypeDefinition implements EnumTypeDefinition {

    private final TypeReference type;
    private final List<EnumValue> enumValues;
    private final Map<String, TypeReference> constants;

    public DefaultEnumTypeDefinition(String name, TypeReference type, Map<String, Term> attributes, List<EnumValue> enumValues,
                                     List<Argument> constants) {
        super(name, attributes, constants);
        this.type = Objects.requireNonNull(type);
        this.enumValues = Objects.requireNonNull(enumValues);
        this.constants = new HashMap<>();
        if (constants != null) {
            for (Argument constant : constants) {
                this.constants.put(constant.getName(), constant.getType());
            }
        }
    }

    @Override
    public TypeReference getType() {
        return type;
    }

    @Override
    public List<EnumValue> getEnumValues() {
        return enumValues;
    }

    @Override
    public List<String> getConstantNames() {
        return new ArrayList<>(constants.keySet());
    }

    @Override
    public TypeReference getConstantType(String constantName) {
        return constants.get(constantName);
    }

}
