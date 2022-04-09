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
package org.apache.plc4x.codegen.ast;

import java.util.List;
import java.util.Set;

public interface Generator {

    /**
     * Do preliminary stuff.
     * @param root .
     * @return .
     */
    Node prepare(Node root);

    void generate(ConstantExpression constantExpression);

    void generateDeclarationWithInitializer(DeclarationStatement declarationStatement);

    void generateDeclaration(DeclarationStatement declarationStatement);

    void generate(ParameterExpression parameterExpression);

    void generatePrimitive(Primitive.DataType primitive);

    void generate(IfStatement ifStatement);

    void writeBlock(Block statements);

    void generate(BinaryExpression binaryExpression);

    void generate(AssignementExpression assignementExpression);

    void generateStaticCall(Method method, List<Node> constantNode);

    void generateCall(Node target, Method method, List<Node> constantNode);

    void generate(NewExpression newExpression);

    void generate(MethodDefinition methodDefinition);

    void generateReturn(Expression value);

    void generateClass(String namespace, String className, List<FieldDeclaration> fields, List<ConstructorDeclaration> constructors, List<MethodDefinition> methods, List<ClassDeclaration> innerClasses, boolean mainClass);

    void generateFieldDeclaration(Set<Modifier> modifiers, TypeDefinition type, String name, Expression initializer);

    void generateFieldReference(TypeDefinition type, String name);

    void generateConstructor(Set<Modifier> modifiers, String className, List<ParameterExpression> parameters, Block body);

    void generateFile(ClassDeclaration mainClass, List<ClassDeclaration> innerClasses);

    void generateType(String typeString);

    void generateComment(String comment);

    void generateNoOp();
}
