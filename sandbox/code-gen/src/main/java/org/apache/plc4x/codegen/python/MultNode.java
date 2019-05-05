package org.apache.plc4x.codegen.python;

import com.fasterxml.jackson.annotation.JsonTypeName;

@JsonTypeName("Mult")
public class MultNode extends LineEntryNode {
    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }
}
