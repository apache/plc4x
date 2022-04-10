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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.references;

import org.apache.plc4x.plugins.codegenerator.types.definitions.EnumTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.references.EnumTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.references.SimpleTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

public class DefaultEnumTypeReference implements EnumTypeReference {

    protected final String name;

    protected final List<Term> params;

    protected transient EnumTypeDefinition typeDefinition;

    public DefaultEnumTypeReference(String name, List<Term> params) {
        this.name = Objects.requireNonNull(name);
        this.params = params;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public Optional<List<Term>> getParams() {
        return Optional.ofNullable(params);
    }

    @Override
    public Optional<SimpleTypeReference> getBaseTypeReference() {
        return getEnumTypeDefinition().getType();
    }

    @Override
    public TypeDefinition getTypeDefinition() {
        return getEnumTypeDefinition();
    }

    @Override
    public void setTypeDefinition(TypeDefinition typeDefinition) {
        Objects.requireNonNull(typeDefinition);
        if(!(typeDefinition instanceof EnumTypeDefinition)) {
            throw new IllegalArgumentException("DefaultEnumTypeReferences only accept instances of EnumTypeDefinitions");
        }
        this.typeDefinition = ((EnumTypeDefinition) typeDefinition);
    }

    @Override
    public EnumTypeDefinition getEnumTypeDefinition() {
        if (typeDefinition == null) {
            throw new IllegalStateException("Should not happen as this should be initialized. No type for " + name + " set!!!");
        }
        return typeDefinition;
    }

    public void setEnumTypeDefinition(EnumTypeDefinition typeDefinition) {
        Objects.requireNonNull(typeDefinition);
        this.typeDefinition = typeDefinition;
    }

    @Override
    public String toString() {
        return "DefaultEnumTypeReference{" +
                "name='" + name + '\'' +
                ", params=" + params +
                '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultEnumTypeReference that = (DefaultEnumTypeReference) o;
        return Objects.equals(name, that.name) && Objects.equals(params, that.params);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, params);
    }
}
