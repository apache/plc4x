package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("arg")
public class ArgNode extends LineEntryNode {

    @JsonProperty("annotation")
    private Node annotation;

    @JsonProperty("arg")
    private String arg;

    public Node getAnnotation() {
        return annotation;
    }

    public void setAnnotation(Node annotation) {
        this.annotation = annotation;
    }

    public String getArg() {
        return arg;
    }

    public void setArg(String arg) {
        this.arg = arg;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
