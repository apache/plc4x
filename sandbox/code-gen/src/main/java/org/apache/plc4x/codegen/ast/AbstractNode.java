package org.apache.plc4x.codegen.ast;

public abstract class AbstractNode implements Node {

    public final TypeNode type;

    protected AbstractNode(TypeNode type) {
        this.type = type;
    }

    public TypeNode getType() {
        return this.type;
    }
}
