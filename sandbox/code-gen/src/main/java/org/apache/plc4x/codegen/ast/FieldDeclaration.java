package org.apache.plc4x.codegen.ast;

public class FieldDeclaration implements Node {

    private final TypeNode type;
    private final String name;

    public FieldDeclaration(TypeNode type, String name) {
        this.type = type;
        this.name = name;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateFieldDeclaration(type, name);
    }
}
