package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("Add")
public class AddNode extends Node {

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }

}
