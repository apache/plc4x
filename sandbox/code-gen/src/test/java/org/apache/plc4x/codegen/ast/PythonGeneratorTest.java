package org.apache.plc4x.codegen.ast;

import org.junit.Before;
import org.junit.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.Assert.*;

public class PythonGeneratorTest {

    public Generator generator;
    public CodeWriter writer;

    @Before
    public void setUp() throws Exception {
        this.writer = new CodeWriter(4);
        this.generator = new PythonGenerator(writer);
    }

    @Test
    public void generateCode() {
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
        assertEquals("class MyClazz:\n" +
            "    \n" +
            "    self.current: double\n" +
            "    \n" +
            "    def __init__(double value):\n" +
            "        self.current = value\n" +
            "    \n" +
            "    class MyInnerClazz:\n" +
            "        pass\n" +
            "        \n", code);
    }
}