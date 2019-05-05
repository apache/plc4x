package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("arguments")
@JsonIgnoreProperties(ignoreUnknown = true)
public class ArgumentsNode extends Node {

    @JsonProperty("args")
    private List<Node> args = new ArrayList<>();

    public List<Node> getArgs() {
        return args;
    }

    public void setArgs(List<Node> args) {
        this.args = args;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
