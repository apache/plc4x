package org.apache.plc4x.codegen.version2;

public class ConstantNode extends Expression {

    private Object value;

    public ConstantNode(Object value) {
        super(TypeUtil.infer(value));
        this.value = value;
    }

    public Object getValue() {
        return value;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }

    @Override public void write(Generator generator) {
        generator.generate(this);
    }

}
