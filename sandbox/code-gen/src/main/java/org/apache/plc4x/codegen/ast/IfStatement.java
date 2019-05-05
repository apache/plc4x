package org.apache.plc4x.codegen.ast;

public class IfStatement extends Statement {

    private Expression condition;
    private Block body;
    private Block orElse;

    public IfStatement(Expression condition, Block body, Block orElse) {
        // Check condition returns Binary
        this.condition = condition;
        this.body = body;
        this.orElse = orElse;
    }

    public Expression getCondition() {
        return condition;
    }

    public Block getBody() {
        return body;
    }

    public Block getOrElse() {
        return orElse;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generate(this);
    }

}
