package org.apache.plc4x.codegen.version2;

public class FieldReference extends ParameterExpression {

    public FieldReference(TypeNode type, String name) {
        super(type, name);
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateFieldReference(type, getName());
    }
}
