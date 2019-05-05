package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("Tuple")
public class TupleNode extends LineEntryNode {

    @JsonProperty("ctx")
    private ContextNode ctx;

    @JsonProperty("elts")
    private List<Node> elts = new ArrayList<>();

    public ContextNode getCtx() {
        return ctx;
    }

    public void setCtx(ContextNode ctx) {
        this.ctx = ctx;
    }

    public List<Node> getElts() {
        return elts;
    }

    public void setElts(List<Node> elts) {
        this.elts = elts;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
