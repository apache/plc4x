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

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class PythonGeneratorTest {

    public Generator generator;
    public CodeWriter writer;

    @BeforeEach
    public void setUp() {
        this.writer = new CodeWriter(4);
        this.generator = new PythonGenerator(writer);
    }

    @Test
    public void generateCode() {
        final FieldDeclaration current = new FieldDeclaration(Primitive.DOUBLE, "current");

        final FieldReference currentRef = new FieldReference(Primitive.DOUBLE, "current");

        final ParameterExpression value = new ParameterExpression(Primitive.DOUBLE, "value");

        // Define Inner Class
        final ClassDeclaration innerClass = new ClassDeclaration("", "MyInnerClazz", Collections.emptyList(), Collections.emptyList(), Collections.emptyList(), null);

        final ClassDeclaration clazz = new ClassDeclaration("org.apache.plc4x",
            "MyClazz",
            Arrays.asList(current),
            Arrays.asList(
                new ConstructorDeclaration(
                    Collections.singletonList(value),
                    new Block(new AssignementExpression(currentRef, value))
                )
            ),
            Collections.emptyList(),
            Collections.singletonList(innerClass));

        clazz.write(generator);
        final String code = writer.getCode();
        assertEquals("class MyClazz:\n" +
            "    \n" +
            "    current: float\n" +
            "    \n" +
            "    def __init__(float value):\n" +
            "        self.current = value\n" +
            "    \n" +
            "    class MyInnerClazz:\n" +
            "        pass\n" +
            "        \n", code);
    }

    @Test
    public void ifMultipleElse() {
        final IfStatement stmt = new IfStatement(
            Arrays.asList(
                new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "a"), new ConstantExpression(10.0), BinaryExpression.Operation.EQ),
                new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "a"), new ConstantExpression(5.0), BinaryExpression.Operation.EQ)
            ),
            Arrays.asList(
                new Block(),
                new Block(),
                new Block()
            ));

        stmt.write(generator);
        final String code = writer.getCode();
        assertEquals("if a == 10.0:\n" +
            "    pass\n" +
            "elif a == 5.0:\n" +
            "    pass\n" +
            "else:\n" +
            "    pass\n", code);
    }
}