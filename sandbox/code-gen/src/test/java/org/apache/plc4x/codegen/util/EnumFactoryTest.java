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
package org.apache.plc4x.codegen.util;

import org.apache.plc4x.codegen.ast.*;
import org.junit.jupiter.api.Test;

import java.util.Arrays;

public class EnumFactoryTest {

    @Test
    public void createEnums() {
        System.out.println("Java:");
        System.out.println("----------");
        final ClassDeclaration MyEnum = newEnum();

        final CodeWriter writer = new CodeWriter(4);
        final JavaGenerator generator = new JavaGenerator(writer);

        MyEnum.write(generator);

        System.out.println(writer.getCode());
    }

    @Test
    public void createEnumsInPython() {
        System.out.println("Python:");
        System.out.println("----------");
        final ClassDeclaration MyEnum = newEnum();

        final CodeWriter writer = new CodeWriter(4);
        final PythonGenerator generator = new PythonGenerator(writer);

        MyEnum.write(generator);

        System.out.println(writer.getCode());
    }

    private ClassDeclaration newEnum() {
        final EnumFactory factory = new EnumFactory();
        final EnumFactory.PojoDescription description = new EnumFactory.PojoDescription("MyPojo",
            new EnumFactory.Field(Primitive.DOUBLE, "field1"),
            new EnumFactory.Field(Primitive.DOUBLE, "field2"),
            new EnumFactory.Field(Primitive.DOUBLE, "field3")
        );
        return factory.create(description,
            Arrays.asList(
                new EnumFactory.EnumEntry("alternative1", Arrays.asList(
                    new ConstantExpression(1.0),
                    new ConstantExpression(2.0),
                    new ConstantExpression(3.0)
                )),
                new EnumFactory.EnumEntry("alternative2", Arrays.asList(
                    new ConstantExpression(1.0),
                    new ConstantExpression(2.0),
                    new ConstantExpression(3.0)
                )),
                new EnumFactory.EnumEntry("alternative3", Arrays.asList(
                    new ConstantExpression(1.0),
                    new ConstantExpression(2.0),
                    new ConstantExpression(3.0)
                )),
                new EnumFactory.EnumEntry("alternative4", Arrays.asList(
                    new ConstantExpression(1.0),
                    new ConstantExpression(2.0),
                    new ConstantExpression(3.0)
                ))
            )
        );
    }

}