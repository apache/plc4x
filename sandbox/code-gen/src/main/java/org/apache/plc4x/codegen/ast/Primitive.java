package org.apache.plc4x.codegen.ast;

public class Primitive extends TypeNode {

    public static final TypeNode DOUBLE = new Primitive("double");
    public static final TypeNode VOID = new Primitive("Void");

    public Primitive(String typeString) {
        super(typeString);
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator generator) {
        generator.generate(this);
    }
}
