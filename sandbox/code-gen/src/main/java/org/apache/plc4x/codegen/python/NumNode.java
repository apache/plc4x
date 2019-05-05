package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("Num")
public class NumNode extends LineEntryNode {

    @JsonProperty("n")
    private double n;

    public double getN() {
        return n;
    }

    public void setN(double n) {
        this.n = n;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
