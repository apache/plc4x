package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("Name")
public class NameNode extends LineEntryNode {

    @JsonProperty("ctx")
    private ContextNode ctx;

    @JsonProperty("id")
    private String id;

    public ContextNode getCtx() {
        return ctx;
    }

    public void setCtx(ContextNode ctx) {
        this.ctx = ctx;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
