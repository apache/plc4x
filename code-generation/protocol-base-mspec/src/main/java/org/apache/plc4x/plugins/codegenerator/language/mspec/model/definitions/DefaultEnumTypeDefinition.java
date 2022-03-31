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
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class DefaultEnumTypeDefinition extends DefaultTypeDefinition implements EnumTypeDefinition {


    private static final Logger LOGGER = LoggerFactory.getLogger(DefaultEnumTypeDefinition.class);


    private final SimpleTypeReference type;
    private final List<EnumValue> enumValues;
    private final Map<String, TypeReference> constants;

    public DefaultEnumTypeDefinition(String name, SimpleTypeReference type, Map<String, Term> attributes, List<EnumValue> enumValues,
                                     List<Argument> parserArgument) {
        super(name, attributes, parserArgument);
        this.type = Objects.requireNonNull(type);
        this.enumValues = Objects.requireNonNull(enumValues);
        this.constants = new HashMap<>();
        if (parserArgument != null) {
            for (Argument argument : parserArgument) {
                ((DefaultArgument) argument).getTypeReferenceCompletionStage().whenComplete((typeReference, throwable) -> {
                    if (throwable != null) {
                        // TODO: proper error collection in type context error bucket
                        LOGGER.debug("Error setting type for {}", argument, throwable);
                        return;
                    }
                    this.constants.put(argument.getName(), argument.getType());
                });
            }
        }
    }

    public Optional<SimpleTypeReference> getType() {
        return Optional.ofNullable(type);
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

    @Override
    public String toString() {
        return "DefaultEnumTypeDefinition{" +
            "type=" + type +
            ", enumValues=" + enumValues +
            ", constants=" + constants +
            "} " + super.toString();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        if (!super.equals(o)) return false;
        DefaultEnumTypeDefinition that = (DefaultEnumTypeDefinition) o;
        return Objects.equals(type, that.type) && Objects.equals(enumValues, that.enumValues) && Objects.equals(constants, that.constants);
    }

    @Override
    public int hashCode() {
        return Objects.hash(super.hashCode(), type, enumValues, constants);
    }
}
