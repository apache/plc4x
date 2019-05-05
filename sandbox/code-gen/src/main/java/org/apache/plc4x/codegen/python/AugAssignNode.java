package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("AugAssign")
public class AugAssignNode extends LineEntryNode {

    @JsonProperty("op")
    private Node op;

    @JsonProperty("target")
    private Node target;

    @JsonProperty("value")
    private Node value;

    public Node getOp() {
        return op;
    }

    public void setOp(Node op) {
        this.op = op;
    }

    public Node getTarget() {
        return target;
    }

    public void setTarget(Node target) {
        this.target = target;
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
