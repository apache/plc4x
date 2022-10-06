/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.*;

public abstract class DefaultTypeDefinition {

    protected final String name;
    private final Map<String, Term> attributes;
    protected final List<Argument> parserArguments;

    public DefaultTypeDefinition(String name, Map<String, Term> attributes, List<Argument> parserArguments) {
        this.name = Objects.requireNonNull(name);
        this.attributes = attributes;
        this.parserArguments = parserArguments;
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

    public Optional<List<Argument>> getAllParserArguments() {
        List<Argument> allArguments = new ArrayList<>();
        if (parserArguments != null) {
            allArguments.addAll(parserArguments);
        }
        return Optional.of(allArguments);
    }

    @Override
    public String toString() {
        return "DefaultTypeDefinition{" +
            "name='" + name + '\'' +
            ", attributes=" + attributes +
            ", parserArguments=" + parserArguments +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultTypeDefinition that = (DefaultTypeDefinition) o;
        return name.equals(that.name) && Objects.equals(attributes, that.attributes) && Objects.equals(parserArguments, that.parserArguments);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, attributes, parserArguments);
    }
}
