package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

import java.util.ArrayList;
import java.util.List;

@JsonTypeName("Compare")
public class CompareNode extends LineEntryNode {

    @JsonProperty("comparators")
    private List<Node> comparators = new ArrayList<>();

    @JsonProperty("left")
    private Node left;

    @JsonProperty("ops")
    private List<Node> ops;

    public List<Node> getComparators() {
        return comparators;
    }

    public void setComparators(List<Node> comparators) {
        this.comparators = comparators;
    }

    public Node getLeft() {
        return left;
    }

    public void setLeft(Node left) {
        this.left = left;
    }

    public List<Node> getOps() {
        return ops;
    }

    public void setOps(List<Node> ops) {
        this.ops = ops;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
