/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.codegen.ast;

import java.util.List;

/**
 * Method Declaration
 * 
 * TODO is this duplicate to something???
 */
public class Method implements Node {

    private final TypeNode type;    // Type where this thing is defined on
    private final String name;
    private final TypeNode returnType;
    private final List<TypeNode> parameterTypes;
    private final List<TypeNode> expressionTypes;

    public Method(TypeNode type, String name, TypeNode returnType, List<TypeNode> parameterTypes, List<TypeNode> expressionTypes) {
        this.type = type;
        this.name = name;
        this.returnType = returnType;
        this.parameterTypes = parameterTypes;
        this.expressionTypes = expressionTypes;
    }

    public TypeNode getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public TypeNode getReturnType() {
        return returnType;
    }

    public List<TypeNode> getParameterTypes() {
        return parameterTypes;
    }

    public List<TypeNode> getExpressionTypes() {
        return expressionTypes;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {

    }
}
