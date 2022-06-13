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
package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;

public class FieldDeclaration implements Node {

    private final Set<Modifier> modifiers;
    private final TypeDefinition type;
    private final String name;

    private final Expression initializer;

    public FieldDeclaration(Set<Modifier> modifiers, TypeDefinition type, String name, Expression initializer) {
        this.modifiers = modifiers;
        this.type = type;
        this.name = name;
        this.initializer = initializer;
    }

    public FieldDeclaration(TypeDefinition type, String name) {
        this(Collections.emptySet(), type, name, null);
    }

    public FieldDeclaration(TypeDefinition type, String name, Modifier... modifiers) {
        this(new HashSet<>(Arrays.asList(modifiers)), type, name, null);
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override
    public void write(Generator writer) {
        writer.generateFieldDeclaration(modifiers, type, name, initializer);
    }
}
