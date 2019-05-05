package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("If")
public class IfNode extends LineEntryNode {

    @JsonProperty("body")
    private List<Node> body = new ArrayList<>();

    @JsonProperty("orelse")
    private List<Node> orelse = new ArrayList<>();

    @JsonProperty("test")
    private Node test;

    public List<Node> getBody() {
        return body;
    }

    public void setBody(List<Node> body) {
        this.body = body;
    }

    public List<Node> getOrelse() {
        return orelse;
    }

    public void setOrelse(List<Node> orelse) {
        this.orelse = orelse;
    }

    public Node getTest() {
        return test;
    }

    public void setTest(Node test) {
        this.test = test;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
