package org.apache.plc4x.codegen.version2;

public class AssignementExpression extends Expression {

    private final Expression target;
    private final Node value;

    protected AssignementExpression(Expression target, Node value) {
        super(target.getType());
        this.target = target;
        this.value = value;
    }

    public Node getTarget() {
        return target;
    }

    public Node getValue() {
        return value;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generate(this);
    }
}
