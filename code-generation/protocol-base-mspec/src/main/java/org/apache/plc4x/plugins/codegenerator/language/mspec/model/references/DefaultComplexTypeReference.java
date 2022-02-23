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

import org.apache.plc4x.plugins.codegenerator.types.definitions.ComplexTypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.definitions.TypeDefinition;
import org.apache.plc4x.plugins.codegenerator.types.references.ComplexTypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

public class DefaultComplexTypeReference implements ComplexTypeReference {

    protected final String name;

    protected final List<Term> params;

    protected transient ComplexTypeDefinition typeDefinition;

    public DefaultComplexTypeReference(String name, List<Term> params) {
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
    public TypeDefinition getTypeDefinition() {
        return getComplexTypeDefinition();
    }

    @Override
    public void setTypeDefinition(TypeDefinition typeDefinition) {
        Objects.requireNonNull(typeDefinition);
        if (!(typeDefinition instanceof ComplexTypeDefinition)) {
            throw new IllegalArgumentException("DefaultComplexTypeReferences only accept instances of ComplexTypeDefinitions. Actual type: " + typeDefinition.getClass());
        }
        this.typeDefinition = ((ComplexTypeDefinition) typeDefinition);
    }

    @Override
    public ComplexTypeDefinition getComplexTypeDefinition() {
        if (typeDefinition == null) {
            throw new IllegalStateException("Should not happen as this should be initialized. No type for " + name + " set!!!");
        }
        return typeDefinition;
    }

    public void setComplexTypeDefinition(ComplexTypeDefinition typeDefinition) {
        Objects.requireNonNull(typeDefinition);
        this.typeDefinition = typeDefinition;
    }

    @Override
    public String toString() {
        return "DefaultComplexTypeReference{" +
            "name='" + name + '\'' +
            ", params=" + params +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultComplexTypeReference that = (DefaultComplexTypeReference) o;
        return Objects.equals(name, that.name) && Objects.equals(params, that.params);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, params);
    }
}
