package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("BinOp")
public class BinOpNode extends LineEntryNode {

    @JsonProperty("left")
    private Node left;

    @JsonProperty("op")
    private Node op;

    @JsonProperty("right")
    private Node right;

    public Node getLeft() {
        return left;
    }

    public void setLeft(Node left) {
        this.left = left;
    }

    public Node getOp() {
        return op;
    }

    public void setOp(Node op) {
        this.op = op;
    }

    public Node getRight() {
        return right;
    }

    public void setRight(Node right) {
        this.right = right;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
