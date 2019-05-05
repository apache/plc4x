package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("Assign")
public class AssignNode extends LineEntryNode {

    @JsonProperty("targets")
    private List<Node> targets = new ArrayList<>();

    @JsonProperty("value")
    private Node value;

    public List<Node> getTargets() {
        return targets;
    }

    public void setTargets(List<Node> targets) {
        this.targets = targets;
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
