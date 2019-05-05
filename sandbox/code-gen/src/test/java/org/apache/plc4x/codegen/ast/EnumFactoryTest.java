package org.apache.plc4x.codegen.ast;

import org.junit.Test;

import java.util.Arrays;

public class EnumFactoryTest {

    @Test
    public void createEnums() {
        System.out.println("Java:");
        System.out.println("----------");
        final ClassDefinition MyEnum = newEnum();

        final CodeWriter writer = new CodeWriter(4);
        final JavaGenerator generator = new JavaGenerator(writer);

        MyEnum.write(generator);

        System.out.println(writer.getCode());
    }

    @Test
    public void createEnumsInPython() {
        System.out.println("Python:");
        System.out.println("----------");
        final ClassDefinition MyEnum = newEnum();

        final CodeWriter writer = new CodeWriter(4);
        final PythonGenerator generator = new PythonGenerator(writer);

        MyEnum.write(generator);

        System.out.println(writer.getCode());
    }

    private ClassDefinition newEnum() {
        final EnumFactory factory = new EnumFactory();
        final EnumFactory.PojoDescription description = new EnumFactory.PojoDescription("MyPojo",
            new EnumFactory.Field(Primitive.DOUBLE, "field1"),
            new EnumFactory.Field(Primitive.DOUBLE, "field2"),
            new EnumFactory.Field(Primitive.DOUBLE, "field3")
        );
        return factory.create(description,
            Arrays.asList(
                new EnumFactory.EnumEntry("alternative1", Arrays.asList(
                    new ConstantNode(1.0),
                    new ConstantNode(2.0),
                    new ConstantNode(3.0)
                )),
                new EnumFactory.EnumEntry("alternative2", Arrays.asList(
                    new ConstantNode(1.0),
                    new ConstantNode(2.0),
                    new ConstantNode(3.0)
                )),
                new EnumFactory.EnumEntry("alternative3", Arrays.asList(
                    new ConstantNode(1.0),
                    new ConstantNode(2.0),
                    new ConstantNode(3.0)
                )),
                new EnumFactory.EnumEntry("alternative4", Arrays.asList(
                    new ConstantNode(1.0),
                    new ConstantNode(2.0),
                    new ConstantNode(3.0)
                ))
            )
        );
    }

}