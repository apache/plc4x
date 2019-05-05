package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.List;

public class IfStatement extends Statement {

    private List<Expression> condition;
    private List<Block> blocks;

    public IfStatement(Expression condition, Block body, Block orElse) {
        this(Arrays.asList(condition), Arrays.asList(body, orElse));
    }

    public IfStatement(List<Expression> condition, List<Block> blocks) {
        assert condition.size() == blocks.size() || condition.size() == (blocks.size() -1);
        this.condition = condition;
        this.blocks = blocks;
    }

    public List<Expression> getConditions() {
        return condition;
    }

    public List<Block> getBlocks() {
        return blocks;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generate(this);
    }

}
