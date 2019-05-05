package org.apache.plc4x.codegen.ast;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

/**
 * Builds a Block.
 */
public class BlockBuilder {

    private final List<Node> statements = new ArrayList<>();

    public BlockBuilder() {
    }

    public BlockBuilder add(Node statements) {
        this.statements.add(statements);
        return this;
    }

    public BlockBuilder add(Collection<Node> statements) {
        this.statements.addAll(statements);
        return this;
    }

    public Block toBlock() {
        return new Block(this.statements);
    }
}
