package org.apache.plc4x.codegen.version2;

import org.junit.Test;

import java.sql.SQLOutput;

import static org.junit.Assert.*;

public class PojoFactoryTest {

    @Test
    public void createPojoJava() {
        System.out.println("Java:");
        System.out.println("----------");
        final ClassDefinition pojo = newPojo();

        final CodeWriter writer = new CodeWriter(4);
        final JavaGenerator generator = new JavaGenerator(writer);

        pojo.write(generator);

        System.out.println(writer.getCode());
    }

    @Test
    public void createPojoPython() {
        System.out.println("Python:");
        System.out.println("----------");
        final ClassDefinition pojo = newPojo();

        final CodeWriter writer = new CodeWriter(4);
        final PythonGenerator generator = new PythonGenerator(writer);

        pojo.write(generator);

        System.out.println(writer.getCode());
    }

    private ClassDefinition newPojo() {
        final PojoFactory factory = new PojoFactory();
        final PojoFactory.PojoDescription description = new PojoFactory.PojoDescription("MyPojo",
            new PojoFactory.Field(Primitive.DOUBLE, "field1"),
            new PojoFactory.Field(Primitive.DOUBLE, "field2"),
            new PojoFactory.Field(Primitive.DOUBLE, "field3")
        );
        return factory.create(description);
    }
}