package org.apache.plc4x.codegen.version2;

public class ParameterExpression extends Expression {

    private final String name;

    public ParameterExpression(TypeNode type, String name) {
        super(type);
        this.name = name;
    }

    public String getName() {
        return name;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }

    @Override public void write(Generator generator) {
        generator.generate(this);
    }
}
