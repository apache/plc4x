/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.plugins.codegenerator.parser;

import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;
import org.antlr.v4.runtime.tree.ParseTree;
import org.antlr.v4.runtime.tree.ParseTreeWalker;
import org.apache.plc4x.codegen.ast.BinaryExpression;
import org.apache.plc4x.codegen.ast.BlockBuilder;
import org.apache.plc4x.codegen.ast.ClassDeclaration;
import org.apache.plc4x.codegen.ast.CodeWriter;
import org.apache.plc4x.codegen.ast.Expression;
import org.apache.plc4x.codegen.ast.Expressions;
import org.apache.plc4x.codegen.ast.Generator;
import org.apache.plc4x.codegen.ast.JavaGenerator;
import org.apache.plc4x.codegen.ast.Method;
import org.apache.plc4x.codegen.ast.ParameterExpression;
import org.apache.plc4x.codegen.ast.Primitive;
import org.apache.plc4x.codegen.ast.PythonGenerator;
import org.apache.plc4x.codegen.util.PojoFactory;
import org.apache.plc4x.codegenerator.parser.imaginary.ImaginaryLexer;
import org.apache.plc4x.codegenerator.parser.imaginary.ImaginaryParser;
import org.apache.plc4x.plugins.codegenerator.model.ComplexType;
import org.apache.plc4x.plugins.codegenerator.model.DiscriminatedComplexType;
import org.apache.plc4x.plugins.codegenerator.model.Type;
import org.apache.plc4x.plugins.codegenerator.model.fields.ArrayField;
import org.apache.plc4x.plugins.codegenerator.model.fields.ConstField;
import org.apache.plc4x.plugins.codegenerator.model.fields.DiscriminatorField;
import org.apache.plc4x.plugins.codegenerator.model.fields.Field;
import org.apache.plc4x.plugins.codegenerator.model.fields.ImplicitField;
import org.apache.plc4x.plugins.codegenerator.model.fields.OptionalField;
import org.apache.plc4x.plugins.codegenerator.model.fields.SimpleField;
import org.apache.plc4x.plugins.codegenerator.model.fields.SwitchField;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class MessageFormatParser {

    /**
     * Turns parsed field into POJO Field.
     * @param field
     * @return
     */
    private static PojoFactory.Field fieldTranslator(Field field) {
        if (field instanceof SimpleField) {
            final String name = ((SimpleField) field).getName();
            final Type type = ((SimpleField) field).getType();
            return new PojoFactory.Field(Expressions.typeOf(type.getName()), name);
        } else if (field instanceof ConstField) {
            final Type type = ((ConstField) field).getType();
            final Object value = ((ConstField) field).getReferenceValue();
            return new PojoFactory.Field(Expressions.typeOf(type.getName()), "const_" + value);
        } else if (field instanceof ImplicitField) {
            final Type type = ((ImplicitField) field).getType();
            final String expr = ((ImplicitField) field).getSerializationExpression();
            return new PojoFactory.Field(Expressions.typeOf(type.getName()), "implicit_" + expr);
        } else if (field instanceof OptionalField) {
            final Type type = ((OptionalField) field).getType();
            final String name = ((OptionalField) field).getName();
            final String expr = ((OptionalField) field).getConditionExpression();
            return new PojoFactory.Field(Expressions.typeOf(type.getName()), name + "?" + expr);
        }
        return new PojoFactory.Field(Primitive.VOID, "UNKNOWN_FIELD");
    }

    public List<Type> parse(InputStream source) {
        try {
            ImaginaryLexer lexer = new ImaginaryLexer(CharStreams.fromStream(source));
            CommonTokenStream tokens = new CommonTokenStream(lexer);
            ImaginaryParser parser = new ImaginaryParser(tokens);
            ParseTree tree = parser.file();
            ParseTreeWalker walker = new ParseTreeWalker();
            MessageFormatListener listener = new MessageFormatListener();
            walker.walk(listener, tree);
            Map<String, ComplexType> complexTypes = listener.getComplexTypes();
            // System.out.println(complexTypes);
            // Do a Pojo for everyone
            final PojoFactory factory = new PojoFactory();
            final List<ClassBuilder> classes = new ArrayList<>();
            for (Map.Entry<String, ComplexType> entry : complexTypes.entrySet()) {
                final CodeWriter writer = new CodeWriter(4);
                final Generator generator = new PythonGenerator(writer);
                System.out.println("// ----------------------");
                System.out.println("//    " + entry.getKey() + "   ");
                System.out.println("// ----------------------");

                // Old Approach
                final ComplexType type = entry.getValue();
                final List<PojoFactory.Field> fields = type.getFields().stream()
                    .map(MessageFormatParser::fieldTranslator)
                    .collect(Collectors.toList());
                final PojoFactory.PojoDescription pojoDesc = new PojoFactory.PojoDescription(entry.getKey(), fields);
                final ClassDeclaration classDeclaration = factory.create(pojoDesc);
                // classDeclaration.write(generator);
                // Output POJO
                // System.out.println(writer.getCode());

                // New Approach
                final ClassBuilder classBuilder = new ClassBuilder();
                classes.add(classBuilder);

                classBuilder.withName(entry.getKey());
                int constCount = 1;
                int implicitCount = 1;
                int tmpCount = 1;
                BlockBuilder writerBuilder = new BlockBuilder();
                BlockBuilder readerBuilder = new BlockBuilder();
                final ParameterExpression _buffer = Expressions.parameter("buffer", Expressions.typeOf("Buffer"));
                final Method _writeUint8 = Expressions.method(Expressions.typeOf("Buffer"), "writeUint8", Primitive.VOID, Collections.singletonList(Primitive.INTEGER), Collections.emptyList());
                final Method _readUint8 = Expressions.method(Expressions.typeOf("Buffer"), "readUint8", Primitive.INTEGER, Collections.emptyList(), Collections.emptyList());
                for (Field field : type.getFields()) {
                    if (field instanceof ArrayField) {
                        // do nothing
                        System.out.println("Skipping Array field...");
                    } else if (field instanceof ConstField) {
                        classBuilder.withConstant("CONST_" + constCount++, ((ConstField) field).getType().getName(), ((ConstField) field).getReferenceValue(), new ClassBuilder.Documentation("Constant expression"));

                        String constString = ((ConstField) field).getReferenceValue().toString();
                        constString = constString.replace("'", "");
                        final Expression _const = Expressions.constant(Integer.decode(constString));

                        // Writer
                        writerBuilder.add(Expressions.comment("Constant"));
                        writerBuilder.add(Expressions.call(
                            _buffer,
                            _writeUint8,
                            _const
                        ));

                        // Reader
                        readerBuilder.add(Expressions.comment("This should be a check on the expected value"));
                        String tmpVar = "tmp" + tmpCount++;
                        readerBuilder.add(Expressions.declaration(
                            tmpVar,
                            Expressions.call(_buffer, _readUint8)
                        ));
                        readerBuilder.add(Expressions.ifElse(
                            Expressions.binaryExpression(Primitive.BOOLEAN, Expressions.parameter(tmpVar, Primitive.VOID), _const, BinaryExpression.Operation.NEQ),
                            Expressions.block(
                                Expressions.comment("We should throw an exception here")
                            )
                        ));

                    } else if (field instanceof SimpleField) {
                        classBuilder.withField(((SimpleField) field).getName(), ((SimpleField) field).getType().getName(), new ClassBuilder.Documentation("Simple Field..."));

                        // Writer

                        final ParameterExpression parameter = Expressions.parameter(
                            ((SimpleField) field).getName(),
                            Expressions.typeOf(((SimpleField) field).getType().getName())
                        );
                        if (((SimpleField) field).getType().getName().startsWith("uint")) {
                            writerBuilder.add(Expressions.call(
                                _buffer,
                                _writeUint8,
                                parameter
                            ));
                        } else {
                            writerBuilder.add(Expressions.call(
                                parameter,
                                "write",
                                Primitive.VOID,
                                _buffer
                            ));
                        }

                        // Reader

                    } else if (field instanceof ImplicitField) {
                        classBuilder.withField("implicit" + implicitCount, ((ImplicitField) field).getType().getName(), ((ImplicitField) field).getSerializationExpression(), new ClassBuilder.Documentation("Implicitly defined field"));
                    } else if (field instanceof DiscriminatorField) {
                        classBuilder.withField(((DiscriminatorField) field).getName(), ((DiscriminatorField) field).getType().getName(), new ClassBuilder.Documentation("Discriminator Field..."));
                    } else if (field instanceof SwitchField) {
                        // Now do all "subtypes"
                        for (DiscriminatedComplexType aCase : ((SwitchField) field).getCases()) {
                            final ClassBuilder childBuilder = new ClassBuilder();
                            classes.add(childBuilder);
                            childBuilder.withName(aCase.getName());
                            childBuilder.withParent(entry.getKey());

                        }
                    } else if (field instanceof OptionalField) {
                        classBuilder.withField(((OptionalField) field).getName(), ((OptionalField) field).getType().getName(),
                            new ClassBuilder.Documentation("Optional field\n" + ((OptionalField) field).getConditionExpression()));
                    } else {
                        throw new RuntimeException("Fields of class " + field.getClass().getSimpleName() + " are not implemented yet!");
                    }
                }

                // System.out.println(classBuilder.toString());

                // Print something here
                writerBuilder.toBlock().write(generator);
                readerBuilder.toBlock().write(generator);
                System.out.println(writer.getCode());
            }

            // Now print all Classes
            for (ClassBuilder aClass : classes) {
                System.out.println(aClass.toString());
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return null;
    }

}
