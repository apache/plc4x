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
package org.apache.plc4x.plugins.codegenerator.language.mspec.model.terms;

import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;
import org.apache.plc4x.plugins.codegenerator.types.terms.Term;
import org.apache.plc4x.plugins.codegenerator.types.terms.VariableLiteral;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

public class DefaultVariableLiteral implements VariableLiteral {

    private final String name;
    private TypeReference typeReference;
    private final List<Term> args;
    private final Integer index;
    private final VariableLiteral child;

    public DefaultVariableLiteral(String name, List<Term> args, Integer index, VariableLiteral child) {
        this.name = Objects.requireNonNull(name);
        this.args = args;
        this.index = index;
        this.child = child;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public TypeReference getTypeReference() {
        if (typeReference == null) {
            throw new IllegalStateException("No type reference set for " + name);
        }
        return typeReference;
    }


    public void setTypeReference(TypeReference typeReference) {
        Objects.requireNonNull(typeReference);
        this.typeReference = typeReference;
    }

    @Override
    public Optional<List<Term>> getArgs() {
        return Optional.ofNullable(args);
    }

    @Override
    public Optional<Integer> getIndex() {
        return Optional.ofNullable(index);
    }

    @Override
    public Optional<VariableLiteral> getChild() {
        return Optional.ofNullable(child);
    }

    @Override
    public String stringRepresentation() {
        return name + getChild().map(Term::stringRepresentation).orElse("");
    }

    @Override
    public String toString() {
        return "DefaultVariableLiteral{" +
            "name='" + name + '\'' +
            ", typeReference='" + typeReference + '\'' +
            ", args=" + args +
            ", index=" + index +
            ", child=" + child +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultVariableLiteral that = (DefaultVariableLiteral) o;
        return index == that.index && name.equals(that.name) && typeReference.equals(that.typeReference) && Objects.equals(args, that.args) && Objects.equals(child, that.child);
    }

    @Override
    public int hashCode() {
        return Objects.hash(name, typeReference, args, index, child);
    }
}
