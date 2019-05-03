package org.apache.plc4x.codegen;

import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("NotEq")
public class NotEqNode extends Node {
    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
