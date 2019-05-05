package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.List;

public class Block extends Statement {

    private final List<Node> statements;

    public Block(List<Node> statements) {
        this.statements = statements;
    }

    public Block(Node... statements) {
        this(Arrays.asList(statements));
    }

    public List<Node> getStatements() {
        return statements;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.writeBlock(this);
    }
}
