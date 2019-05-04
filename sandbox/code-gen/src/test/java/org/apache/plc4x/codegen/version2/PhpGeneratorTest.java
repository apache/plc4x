package org.apache.plc4x.codegen.version2;

import org.junit.Test;

import static org.junit.Assert.assertEquals;

public class PhpGeneratorTest {

    @Test
    public void createSomeCode() {
        final CodeWriter writer = new CodeWriter(4);

        final DeclarationStatement stmt = new DeclarationStatement(new ParameterExpression(Primitive.DOUBLE, "a"), new ConstantNode(5.0));

        final Generator generator = new PhpGenerator(writer);

        stmt.write(generator);

        final String code = writer.getCode();

        assertEquals("$a = 5.0\n", code);
    }
}