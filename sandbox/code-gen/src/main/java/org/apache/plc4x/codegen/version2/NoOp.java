package org.apache.plc4x.codegen.version2;

/**
 * Do Nothing command;
 */
public class NoOp extends Expression {

    protected NoOp() {
        super(Primitive.VOID);
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generateNoOp();
    }
}
