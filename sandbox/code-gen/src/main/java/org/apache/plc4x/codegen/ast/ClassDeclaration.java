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

import java.util.List;

public class ClassDeclaration implements Node {

    private final String namespace;
    private final String className;

    private final List<FieldDeclaration> fields;
    private final List<ConstructorDeclaration> constructors;
    private final List<MethodDefinition> methods;
    private final List<ClassDeclaration> innerClasses;

    public ClassDeclaration(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods, List<ClassDeclaration> innerClasses) {
        this.namespace = namespace;
        this.className = className;
        this.fields = fields;
        this.constructors = constructors;
        this.methods = methods;
        this.innerClasses = innerClasses;
    }

    public String getNamespace() {
        return namespace;
    }

    public String getClassName() {
        return className;
    }

    public List<FieldDeclaration> getFields() {
        return fields;
    }

    public List<ConstructorDeclaration> getConstructors() {
        return constructors;
    }

    public List<MethodDefinition> getMethods() {
        return methods;
    }

    public List<ClassDeclaration> getInnerClasses() {
        return innerClasses;
    }

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override
    public void write(Generator writer) {
        writer.generateClass(this.namespace, this.className, this.fields, this.constructors, this.methods, innerClasses, true);
    }
}
