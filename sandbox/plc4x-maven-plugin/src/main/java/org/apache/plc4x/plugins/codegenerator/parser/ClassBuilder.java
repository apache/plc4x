package org.apache.plc4x.plugins.codegenerator.parser;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

/**
 * Util to Define Classes based on the {@link org.apache.plc4x.codegen.ast.Expressions} API.
 *
 * @author julian
 */
public class ClassBuilder {

    private String className;
    private final List<ConstantField> constants = new ArrayList<>();
    private final List<VariableField> fields = new ArrayList<>();
    // TODO constructor call??
    private String parent;

    public ClassBuilder() {
    }

    public ClassBuilder withName(String name) {
        className = name;
        return this;
    }

    public ClassBuilder withConstant(String name, String type, Object value) {
        this.constants.add(new ConstantField(name, type, value));
        return this;
    }

    public ClassBuilder withConstant(String name, String type, Object value, Documentation docu) {
        this.constants.add(new ConstantField(name, type, value, docu));
        return this;
    }

    public ClassBuilder withField(String name, String type) {
        this.fields.add(new VariableField(name, type));
        return this;
    }

    public ClassBuilder withField(String name, String type, Documentation docu) {
        this.fields.add(new VariableField(name, type, Optional.empty(), docu));
        return this;
    }

    public ClassBuilder withField(String name, String type, Object initialValue, Documentation docu) {
        this.fields.add(new VariableField(name, type, Optional.of(initialValue), docu));
        return this;
    }

    public ClassBuilder withParent(String parent) {
        this.parent = parent;
        return this;
    }

    @Override public String toString() {
        return "ClassBuilder{" +
            "className='" + className + '\'' +
            ", constants=" + constants +
            ", fields=" + fields +
            ", parent='" + parent + '\'' +
            '}';
    }

    // Regular "java" fields
    public static class VariableField {

        private final String name;
        private final String type;
        private final Optional<Object> initialValue;
        private final Documentation documentation;

        public VariableField(String name, String type) {
            this(name, type, Optional.empty(), null);
        }

        public VariableField(String name, String type, Optional<Object> initialValue) {
            this(name, type, initialValue, null);
        }

        public VariableField(String name, String type, Optional<Object> initialValue, Documentation documentation) {
            this.name = name;
            this.type = type;
            this.initialValue = initialValue;
            this.documentation = documentation;
        }

        @Override public String toString() {
            return "VariableField{" +
                "name='" + name + '\'' +
                ", type='" + type + '\'' +
                ", initialValue=" + initialValue +
                ", documentation=" + documentation +
                '}';
        }
    }

    // Consts
    public static class ConstantField {

        private final String name;
        private final String type;
        private final Object value;
        private final Documentation documentation;

        public ConstantField(String name, String type, Object value) {
            this(name, type, value, null);
        }

        public ConstantField(String name, String type, Object value, Documentation documentation) {
            this.name = name;
            this.type = type;
            this.value = value;
            this.documentation = documentation;
        }

        @Override public String toString() {
            return "ConstantField{" +
                "name='" + name + '\'' +
                ", type='" + type + '\'' +
                ", value=" + value +
                ", documentation=" + documentation +
                '}';
        }
    }

    /**
     * Documentation for a declaration
     */
    public static class Documentation {

        private final String text;

        public Documentation(String text) {
            this.text = text;
        }

        @Override public String toString() {
            return "Documentation{" +
                "text='" + text + '\'' +
                '}';
        }
    }

}
