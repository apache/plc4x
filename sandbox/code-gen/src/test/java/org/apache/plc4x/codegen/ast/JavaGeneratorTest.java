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

public class JavaGeneratorTest {

    private CodeWriter writer;
    private JavaGenerator generator;

    @BeforeEach
    public void setUp() {
        writer = new CodeWriter(4);
        generator = new JavaGenerator(writer);
    }

    @Test
    public void writeDeclaration() {
        final DeclarationStatement stmt = new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "a"), new ConstantExpression(5.0));
        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("Double a = 5.0", code);
    }

    @Test
    public void ifStatement() {
        final Statement stmt = new IfStatement(
            new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "a"),
                new ConstantExpression(10.0),
                BinaryExpression.Operation.EQ),
            new Block(new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "b"), new ConstantExpression(5.0))),
            new Block(new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "b"), new ConstantExpression(3.0)))
        );

        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("if (a == 10.0) {\n" +
            "    Double b = 5.0;\n" +
            "} else {\n" +
            "    Double b = 3.0;\n" +
            "}\n", code);
    }

    @Test
    public void IfIf() {
        final Statement stmt = new Block(Arrays.asList(
            new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "c"), null),
            new IfStatement(
                new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "a"),
                    new ConstantExpression(10.0),
                    BinaryExpression.Operation.EQ),
                new Block(
                    new IfStatement(
                        new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "b"),
                            new ConstantExpression(10.0),
                            BinaryExpression.Operation.EQ),
                        new Block(
                            new AssignementExpression(new ParameterExpression(Primitive.DOUBLE, "c"), new ConstantExpression(5.0)),
                            new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "d"), new ConstantExpression(100.0))
                        ),
                        null
                    )
                ),
                null
            )));

        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("    Double c;\n" +
            "    if (a == 10.0) {\n" +
            "        if (b == 10.0) {\n" +
            "            c = 5.0;\n" +
            "            Double d = 100.0;\n" +
            "        }\n" +
            "    }\n", code);
    }

    @Test
    public void callStaticMethod() {
        final TypeDefinition myClazz = new TypeDefinition("MyClazz");
        Expression expr = new CallExpression(new Method(myClazz, "toString", Primitive.VOID, Collections.singletonList(Primitive.DOUBLE), Collections.EMPTY_LIST), null, new ConstantExpression(5.0));

        expr.write(generator);

        final String code = writer.getCode();

        assertEquals("MyClazz.toString(5.0)", code);
    }

    @Test
    public void callMethod() {
        final TypeDefinition myClazz = new TypeDefinition("MyClazz");
        Expression expr = new CallExpression(new Method(myClazz, "toString", Primitive.VOID, Collections.singletonList(Primitive.DOUBLE), Collections.EMPTY_LIST), new ParameterExpression(myClazz, "a"), new ConstantExpression(5.0));

        expr.write(generator);

        final String code = writer.getCode();

        assertEquals("a.toString(5.0)", code);
    }

    @Test
    public void complexCallAssignment() {
        final TypeDefinition myClazz = new TypeDefinition("MyClazz");
        final ParameterExpression instance = new ParameterExpression(myClazz, "instance");
        final Method getNumberMethod = new Method(myClazz, "getNumber", Primitive.DOUBLE, Collections.emptyList(), Collections.emptyList());
        final Method staticMethod = new Method(myClazz, "staticMethod", Primitive.DOUBLE, Collections.emptyList(), Collections.emptyList());
        Statement stmt = new Block(Arrays.asList(
            new DeclarationStatement(instance, null),
            new AssignementExpression(instance, new NewExpression(myClazz)),
            new AssignementExpression(instance, new NewExpression(myClazz, new NewExpression(myClazz))),
            new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "a"), new CallExpression(getNumberMethod, instance)),
            new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "b"), new CallExpression(staticMethod, null))
        ));

        stmt.write(generator);
        final String code = writer.getCode();
        assertEquals("    MyClazz instance;\n" +
            "    instance = new MyClazz();\n" +
            "    instance = new MyClazz(new MyClazz());\n" +
            "    Double a = instance.getNumber();\n" +
            "    Double b = MyClazz.staticMethod();\n", code);
    }

    @Test
    public void writeMethodDefinition() {
        final ParameterExpression a = new ParameterExpression(Primitive.DOUBLE, "a");
        final ParameterExpression b = new ParameterExpression(Primitive.DOUBLE, "b");
        final MethodDefinition decl = new MethodDefinition("add", Primitive.DOUBLE,
            Arrays.asList(
                a,
                b
            ),
            new Block(
                new ReturnStatement(
                    new BinaryExpression(a.getType(), a, b, BinaryExpression.Operation.PLUS)
                )
            )
        );

        decl.write(generator);
        final String code = writer.getCode();
        assertEquals("public Double add(Double a, Double b) {\n" +
            "    return a + b;\n" +
            "}\n", code);
    }

    @Test
    public void defineClass() {
        final FieldDeclaration current = new FieldDeclaration(Primitive.DOUBLE, "current");
        final FieldReference currentRef = new FieldReference(Primitive.DOUBLE, "current");

        final ParameterExpression a = new ParameterExpression(Primitive.DOUBLE, "a");
        final ParameterExpression b = new ParameterExpression(Primitive.DOUBLE, "b");
        final MethodDefinition decl = new MethodDefinition("add", Primitive.DOUBLE,
            Arrays.asList(
                a,
                b
            ),
            new Block(
                new ReturnStatement(
                    new BinaryExpression(a.getType(), a, b, BinaryExpression.Operation.PLUS)
                )
            )
        );
        final MethodDefinition inc = new MethodDefinition("inc", Primitive.VOID,
            Collections.EMPTY_LIST,
            new Block(
                new AssignementExpression(
                    currentRef,
                    new BinaryExpression(currentRef.getType(), currentRef, new ConstantExpression(1.0), BinaryExpression.Operation.PLUS)
                )
            )
        );

        final ClassDeclaration clazz = new ClassDeclaration("org.apache.plc4x", "MyClazz", Arrays.asList(current), Collections.emptyList(), Arrays.asList(inc, decl), null);

        clazz.write(generator);
        final String code = writer.getCode();
        assertEquals("public class MyClazz {\n" +
            "    \n" +
            "    public Double current;\n" +
            "    \n" +
            "    public void inc() {\n" +
            "        this.current = this.current + 1.0;\n" +
            "    }\n" +
            "    \n" +
            "    public Double add(Double a, Double b) {\n" +
            "        return a + b;\n" +
            "    }\n" +
            "    \n" +
            "}\n", code);
    }

    @Test
    public void defineClassWithConstructor() {
        final FieldDeclaration current = new FieldDeclaration(Primitive.DOUBLE, "current");

        final FieldReference currentRef = new FieldReference(Primitive.DOUBLE, "current");

        final ParameterExpression value = new ParameterExpression(Primitive.DOUBLE, "value");
        final ClassDeclaration clazz = new ClassDeclaration("org.apache.plc4x",
            "MyClazz",
            Arrays.asList(current),
            Arrays.asList(
                new ConstructorDeclaration(
                    Collections.singletonList(value),
                    new Block(new AssignementExpression(currentRef, value))
                )
            ),
            Collections.emptyList(), null);

        clazz.write(generator);
        final String code = writer.getCode();
        assertEquals("public class MyClazz {\n" +
            "    \n" +
            "    public Double current;\n" +
            "    \n" +
            "    public MyClazz(Double value) {\n" +
            "        this.current = value;\n" +
            "    }\n" +
            "    \n" +
            "}\n", code);
    }

    @Test
    public void defineClassWithInnerClass() {
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
        assertEquals("public class MyClazz {\n" +
            "    \n" +
            "    public Double current;\n" +
            "    \n" +
            "    public MyClazz(Double value) {\n" +
            "        this.current = value;\n" +
            "    }\n" +
            "    \n" +
            "    public static class MyInnerClazz {\n" +
            "        \n" +
            "    }\n" +
            "}\n", code);
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
        assertEquals("if (a == 10.0) {\n" +
            "} else if (a == 5.0) {\n" +
            "} else {\n" +
            "}\n", code);
    }
}