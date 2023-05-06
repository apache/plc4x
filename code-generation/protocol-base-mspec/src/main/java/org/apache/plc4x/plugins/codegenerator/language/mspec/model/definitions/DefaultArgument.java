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
import org.apache.plc4x.plugins.codegenerator.types.references.TypeReference;

import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.CompletionStage;

public class DefaultArgument implements Argument {

    private final String name;

    private TypeReference type;

    protected final CompletableFuture<TypeReference> typeReferenceCompletionStage = new CompletableFuture<>();

    public DefaultArgument(String name) {
        this.name = Objects.requireNonNull(name);
    }

    public DefaultArgument(String name, TypeReference type) {
        this.name = Objects.requireNonNull(name);
        this.type = type;
    }

    public String getName() {
        return name;
    }

    public TypeReference getType() {
        if (type == null) {
            throw new IllegalStateException("type not set");
        }
        return type;
    }

    public void setType(TypeReference typeReference) {
        this.type = typeReference;
        typeReferenceCompletionStage.complete(typeReference);
    }

    public CompletionStage<TypeReference> getTypeReferenceCompletionStage() {
        return typeReferenceCompletionStage;
    }

    @Override
    public String toString() {
        return "DefaultArgument{" +
            "type=" + type +
            ", name='" + name + '\'' +
            '}';
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DefaultArgument that = (DefaultArgument) o;
        return Objects.equals(type, that.type) && Objects.equals(name, that.name);
    }

    @Override
    public int hashCode() {
        return Objects.hash(type, name);
    }
}
