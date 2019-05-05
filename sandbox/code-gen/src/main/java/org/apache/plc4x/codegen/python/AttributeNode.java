package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("Attribute")
public class AttributeNode extends LineEntryNode {

    @JsonProperty("attr")
    private String attr;

    @JsonProperty("ctx")
    private ContextNode ctx;

    @JsonProperty("value")
    private Node value;

    public String getAttr() {
        return attr;
    }

    public void setAttr(String attr) {
        this.attr = attr;
    }

    public ContextNode getCtx() {
        return ctx;
    }

    public void setCtx(ContextNode ctx) {
        this.ctx = ctx;
    }

    public Node getValue() {
        return value;
    }

    public void setValue(Node value) {
        this.value = value;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
