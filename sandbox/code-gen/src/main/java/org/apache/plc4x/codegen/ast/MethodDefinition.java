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

import java.util.Collections;
import java.util.List;
import java.util.Set;

public class MethodDefinition implements Node {

    private final Set<Modifier> modifiers;
    private final String name;
    private final TypeDefinition resultType;
    private final List<ParameterExpression> parameters;
    private final Block body;

    public MethodDefinition(Set<Modifier> modifiers, String name, TypeDefinition resultType, List<ParameterExpression> parameters, Block body) {
        this.modifiers = modifiers;
        this.name = name;
        this.resultType = resultType;
        this.parameters = parameters;
        this.body = body;
    }

    public MethodDefinition(String name, TypeDefinition resultType, List<ParameterExpression> parameters, Block body) {
        this(Collections.<Modifier>emptySet(), name, resultType, parameters, body);
    }

    public Set<Modifier> getModifiers() {
        return modifiers;
    }

    public String getName() {
        return name;
    }

    public TypeDefinition getResultType() {
        return resultType;
    }

    public List<ParameterExpression> getParameters() {
        return parameters;
    }

    public Block getBody() {
        return body;
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override
    public void write(Generator writer) {
        writer.generate(this);
    }

}
