package org.apache.plc4x.codegen.ast;

public class ReturnStatement extends Statement {

    private final Expression value;

    public ReturnStatement(Expression value) {
        this.value = value;
    }

    public Expression getValue() {
        return value;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateReturn(this.getValue());
    }
}
