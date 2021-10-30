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
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.references.DefaultComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.DefaultVariableLiteral;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;
import java.util.stream.Collectors;

public abstract class DefaultTypeDefinition {

    protected final String name;
    private final Map<String, Term> attributes;
    protected final List<Argument> parserArguments;
    protected TypeDefinition parentType;

    public DefaultTypeDefinition(String name, Map<String, Term> attributes, List<Argument> parserArguments) {
        this.name = Objects.requireNonNull(name);
        this.attributes = attributes;
        this.parserArguments = parserArguments;
        this.parentType = null;
    }

    public String getName() {
        return name;
    }

    public Optional<Term> getAttribute(String attributeName) {
        if (attributes.containsKey(attributeName)) {
            return Optional.of(attributes.get(attributeName));
        }
        return Optional.empty();
    }

    public Optional<List<Argument>> getParserArguments() {
        return Optional.ofNullable(parserArguments);
    }

    public TypeDefinition getParentType() {
        return parentType;
    }

    public void setParentType(TypeDefinition parentType) {
        this.parentType = parentType;
    }

}
