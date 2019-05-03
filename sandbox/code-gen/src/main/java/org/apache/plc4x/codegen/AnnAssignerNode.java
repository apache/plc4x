package org.apache.plc4x.codegen;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("AnnAssign")
public class AnnAssignerNode extends LineEntryNode {

    @JsonProperty("annotation")
    private Node annotation;

    @JsonProperty("simple")
    private int simple;

    @JsonProperty("target")
    private Node target;

    @JsonProperty("value")
    private Node value;

    public Node getAnnotation() {
        return annotation;
    }

    public void setAnnotation(Node annotation) {
        this.annotation = annotation;
    }

    public int getSimple() {
        return simple;
    }

    public void setSimple(int simple) {
        this.simple = simple;
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
