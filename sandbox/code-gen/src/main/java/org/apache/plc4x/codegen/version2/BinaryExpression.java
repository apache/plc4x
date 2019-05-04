package org.apache.plc4x.codegen.version2;

public class BinaryExpression extends Expression {

    private final Node left;
    private final Node right;
    private final Operation op;

    protected BinaryExpression(TypeNode type, Node left, Node right, Operation op) {
        super(type);
        this.left = left;
        this.right = right;
        this.op = op;
    }

    public Node getLeft() {
        return left;
    }

    public Node getRight() {
        return right;
    }

    public Operation getOp() {
        return op;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generate(this);
    }

    public enum Operation {
        EQ,
        NEQ,
        GT,
        LT,
        PLUS;
    }
}
