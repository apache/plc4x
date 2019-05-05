package org.apache.plc4x.codegen.ast;

import org.junit.Before;
import org.junit.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.Assert.assertEquals;

public class JavaGeneratorTest {

    private CodeWriter writer;
    private JavaGenerator generator;

    @Before
    public void setUp() throws Exception {
        writer = new CodeWriter(4);
        generator = new JavaGenerator(writer);
    }

    @Test
    public void writeDeclaration() {
        final DeclarationStatement stmt = new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "a"), new ConstantNode(5.0));
        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("double a = 5.0", code);
    }

    @Test
    public void ifStatement() {
        final Statement stmt = new IfStatement(
            new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "a"),
                new ConstantNode(10.0),
                BinaryExpression.Operation.EQ),
            new Block(new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "b"), new ConstantNode(5.0))),
            new Block(new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "b"), new ConstantNode(3.0)))
        );

        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("if (a == 10.0) {\n" +
            "    double b = 5.0;\n" +
            "} else {\n" +
            "    double b = 3.0;\n" +
            "}\n", code);
    }

    @Test
    public void IfIf() {
        final Statement stmt = new Block(Arrays.asList(
            new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "c"), null),
            new IfStatement(
                new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "a"),
                    new ConstantNode(10.0),
                    BinaryExpression.Operation.EQ),
                new Block(
                    new IfStatement(
                        new BinaryExpression(Primitive.DOUBLE, new ParameterExpression(Primitive.DOUBLE, "b"),
                            new ConstantNode(10.0),
                            BinaryExpression.Operation.EQ),
                        new Block(
                            new AssignementExpression(new ParameterExpression(Primitive.DOUBLE, "c"), new ConstantNode(5.0)),
                            new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "d"), new ConstantNode(100.0))
                        ),
                        null
                    )
                ),
                null
            )));

        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("    double c;\n" +
            "    if (a == 10.0) {\n" +
            "        if (b == 10.0) {\n" +
            "            c = 5.0;\n" +
            "            double d = 100.0;\n" +
            "        }\n" +
            "    }\n", code);
    }

    @Test
    public void callStaticMethod() {
        final TypeNode myClazz = new Primitive("MyClazz");
        Expression expr = new CallExpression(new Method(myClazz, "toString", Primitive.VOID, Collections.singletonList(Primitive.DOUBLE), Collections.EMPTY_LIST), null, new ConstantNode(5.0));

        expr.write(generator);

        final String code = writer.getCode();

        assertEquals("MyClazz.toString(5.0)", code);
    }

    @Test
    public void callMethod() {
        final TypeNode myClazz = new Primitive("MyClazz");
        Expression expr = new CallExpression(new Method(myClazz, "toString", Primitive.VOID, Collections.singletonList(Primitive.DOUBLE), Collections.EMPTY_LIST), new ParameterExpression(myClazz, "a"), new ConstantNode(5.0));

        expr.write(generator);

        final String code = writer.getCode();

        assertEquals("a.toString(5.0)", code);
    }

    @Test
    public void complexCallAssignment() {
        final TypeNode myClazz = new Primitive("MyClazz");
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
            "    double a = instance.getNumber();\n" +
            "    double b = MyClazz.staticMethod();\n", code);
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
        assertEquals("public double add(double a, double b) {\n" +
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
                    new BinaryExpression(currentRef.getType(), currentRef, new ConstantNode(1.0), BinaryExpression.Operation.PLUS)
                )
            )
        );

        final ClassDefinition clazz = new ClassDefinition("org.apache.plc4x", "MyClazz", Arrays.asList(current), Collections.emptyList(), Arrays.asList(inc, decl), null);

        clazz.write(generator);
        final String code = writer.getCode();
        assertEquals("public class MyClazz {\n" +
            "    \n" +
            "    public double current;\n" +
            "    \n" +
            "    public void inc() {\n" +
            "        this.current = this.current + 1.0;\n" +
            "    }\n" +
            "    \n" +
            "    public double add(double a, double b) {\n" +
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
        final ClassDefinition clazz = new ClassDefinition("org.apache.plc4x",
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
            "    public double current;\n" +
            "    \n" +
            "    public MyClazz(double value) {\n" +
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
        final ClassDefinition innerClass = new ClassDefinition("", "MyInnerClazz", Collections.emptyList(), Collections.emptyList(), Collections.emptyList(), null);

        final ClassDefinition clazz = new ClassDefinition("org.apache.plc4x",
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
            "    public double current;\n" +
            "    \n" +
            "    public MyClazz(double value) {\n" +
            "        this.current = value;\n" +
            "    }\n" +
            "    \n" +
            "    public static class MyInnerClazz {\n" +
            "        \n" +
            "    }\n" +
            "}\n", code);
    }
}