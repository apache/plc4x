package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("Call")
public class CallNode extends LineEntryNode {

    @JsonProperty("args")
    private List<Node> args = new ArrayList<>();

    @JsonProperty("func")
    private Node func;

    @JsonProperty("keywords")
    private List<Node> keywords = new ArrayList<>();

    public List<Node> getArgs() {
        return args;
    }

    public void setArgs(List<Node> args) {
        this.args = args;
    }

    public Node getFunc() {
        return func;
    }

    public void setFunc(Node func) {
        this.func = func;
    }

    public List<Node> getKeywords() {
        return keywords;
    }

    public void setKeywords(List<Node> keywords) {
        this.keywords = keywords;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
