package org.apache.plc4x.codegen.version2;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.function.Function;
import java.util.stream.Collectors;

public class PojoFactory {

    public ClassDefinition create(PojoDescription desc) {
        // Create all Fields first
        final List<FieldDeclaration> fields = desc.fields.stream()
            .map(field -> new FieldDeclaration(field.getType(), field.getName()))
            .collect(Collectors.toList());

        // Create the Getters
        final List<MethodDefinition> getters = desc.fields.stream()
            .map(new Function<Field, MethodDefinition>() {
                @Override public MethodDefinition apply(Field field) {
                    String getter = "get" + field.getName().substring(0, 1).toUpperCase() + field.getName().substring(1, field.getName().length());
                    Block body = new Block(new ReturnStatement(new FieldReference(field.getType(), field.getName())));
                    return new MethodDefinition(getter, field.getType(), Collections.emptyList(), body);
                }
            })
            .collect(Collectors.toList());


        return new ClassDefinition("", desc.getName(), fields, null, getters, null);
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

        private final TypeNode type;
        private final String name;

        public Field(TypeNode type, String name) {
            this.type = type;
            this.name = name;
        }

        public TypeNode getType() {
            return type;
        }

        public String getName() {
            return name;
        }
    }

}
