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

import java.util.*;
import java.util.stream.Collectors;

/**
 * Creates Enum Like things...
 */
public class EnumFactory {

    public ClassDeclaration create(PojoDescription desc, List<EnumEntry> enumEntries) {
        // Create all Fields first
        final List<Node> constructorStatements = new ArrayList<>();
        final List<FieldDeclaration> fields = new ArrayList<>();
        for (Field field1 : desc.fields) {
            FieldDeclaration fieldDeclaration = new FieldDeclaration(field1.getType(), field1.getName(), Modifier.PRIVATE, Modifier.FINAL);
            fields.add(fieldDeclaration);
            constructorStatements.add(Expressions.assignment(
                Expressions.field(field1.getName()),
                Expressions.parameter(field1.getName(), field1.getType())
            ));
        }

        // Create the Getters
        final List<MethodDefinition> getters = desc.fields.stream()
            .map(field -> getGetterDefinition(field.getName(), field.getType()))
            .collect(Collectors.toList());

        // Now add all these getters
        final List<MethodDefinition> methods = new ArrayList<>(getters);

        final ArrayList<FieldDeclaration> finalFields = new ArrayList<>();

        // Now add all these static methods on top
        for (EnumEntry enumEntry : enumEntries) {
            final TypeDefinition clazzType = Expressions.typeOf(desc.getName());
            final String fieldName = enumEntry.getName().toUpperCase();
            finalFields.add(new FieldDeclaration(
                new HashSet<>(Arrays.asList(Modifier.STATIC, Modifier.FINAL)),
                clazzType,
                fieldName,
                Expressions.new_(clazzType, enumEntry.getArguments())
            ));
        }

        finalFields.addAll(fields);

        // Add constructor
        final ConstructorDeclaration constructor = new ConstructorDeclaration(Collections.singleton(Modifier.PRIVATE),
            desc.fields.stream()
                .map(field -> Expressions.parameter(field.getName(), field.getType())).collect(Collectors.toList()),
            Block.build().add(constructorStatements).toBlock());


        return new ClassDeclaration("", desc.getName(), finalFields, Arrays.asList(constructor), methods, null);
    }

    static MethodDefinition getGetterDefinition(String name, TypeDefinition type) {
        String getter = "get" + name.substring(0, 1).toUpperCase() + name.substring(1);
        Block body = Block.build().add(Expressions.return_(Expressions.field(name))).toBlock();
        return new MethodDefinition(getter, type, Collections.emptyList(), body);
    }

    static MethodDefinition getSetterDefinition(String name, TypeDefinition type) {
        String getter = "set" + name.substring(0, 1).toUpperCase() + name.substring(1);
        final ParameterExpression param = Expressions.parameter(name, type);
        Block body = Block.build().add(Expressions.assignment(Expressions.field(name), param)).toBlock();
        return new MethodDefinition(getter, Primitive.VOID, Collections.singletonList(param), body);
    }

    public static class EnumEntry {

        private final String name;
        private final List<Node> arguments;

        public EnumEntry(String name, List<Node> arguments) {
            this.name = name;
            this.arguments = arguments;
        }

        public String getName() {
            return name;
        }

        public List<Node> getArguments() {
            return arguments;
        }
    }

    public static class PojoDescription {

        private final String name;
        private final List<Field> fields;

        public PojoDescription(String name, Field... fields) {
            this.name = name;
            this.fields = Arrays.asList(fields);
        }

        public PojoDescription(String name, List<Field> fields) {
            this.name = name;
            this.fields = fields;
        }

        public String getName() {
            return name;
        }

        public List<Field> getFields() {
            return fields;
        }
    }

    public static class Field {

        private final TypeDefinition type;
        private final String name;

        public Field(TypeDefinition type, String name) {
            this.type = type;
            this.name = name;
        }

        public TypeDefinition getType() {
            return type;
        }

        public String getName() {
            return name;
        }
    }

}
