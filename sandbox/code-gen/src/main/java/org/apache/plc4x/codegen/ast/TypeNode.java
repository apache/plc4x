package org.apache.plc4x.codegen.ast;

/**
 * Stub for the Type System.
 */
public class TypeNode implements Node {

    private final String typeString;

    public TypeNode(String typeString) {
        this.typeString = typeString;
    }

    String getTypeString() {
        return this.typeString;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateType(typeString);
    }
}
