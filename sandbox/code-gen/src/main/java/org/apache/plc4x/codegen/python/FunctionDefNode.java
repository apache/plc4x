package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("FunctionDef")
@JsonIgnoreProperties(ignoreUnknown = true)
public class FunctionDefNode extends LineEntryNode {

    @JsonProperty("args")
    private ArgumentsNode args;

    @JsonProperty("body")
    private List<Node> body = new ArrayList<>();

    @JsonProperty("name")
    private String name;

    @JsonProperty("returns")
    private Node returns;

    public ArgumentsNode getArgs() {
        return args;
    }

    public void setArgs(ArgumentsNode args) {
        this.args = args;
    }

    public List<Node> getBody() {
        return body;
    }

    public void setBody(List<Node> body) {
        this.body = body;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public Node getReturns() {
        return returns;
    }

    public void setReturns(Node returns) {
        this.returns = returns;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
