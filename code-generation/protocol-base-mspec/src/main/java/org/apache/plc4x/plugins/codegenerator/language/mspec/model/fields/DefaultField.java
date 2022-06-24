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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.fields;

import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.Map;
import java.util.Objects;
import java.util.Optional;

public abstract class DefaultField {

    protected TypeDefinition owner;
    protected final Map<String, Term> attributes;

    protected DefaultField(Map<String, Term> attributes) {
        this.attributes = Objects.requireNonNull(attributes);
    }

    public void setOwner(TypeDefinition owner) {
        this.owner = owner;
    }

    public TypeDefinition getOwner() {
        return owner;
    }

    public Optional<Term> getAttribute(String attributeName) {
        if(attributes.containsKey(attributeName)) {
            return Optional.of(attributes.get(attributeName));
        }
        return Optional.empty();
    }

    @Override
    public String toString() {
        return "DefaultField{" +
            "attributes=" + attributes +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultField that = (DefaultField) o;
        return Objects.equals(attributes, that.attributes);
    }

    @Override
    public int hashCode() {
        return Objects.hash(attributes);
    }
}
