package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("keyword")
public class KeywordNode extends Node {

    @JsonProperty("arg")
    private String arg;

    @JsonProperty("value")
    private Node value;

    public String getArg() {
        return arg;
    }

    public void setArg(String arg) {
        this.arg = arg;
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
