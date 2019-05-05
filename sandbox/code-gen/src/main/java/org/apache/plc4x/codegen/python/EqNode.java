package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("Eq")
public class EqNode extends Node {
    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
