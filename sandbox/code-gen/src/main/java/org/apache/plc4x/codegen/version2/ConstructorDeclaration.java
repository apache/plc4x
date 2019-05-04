package org.apache.plc4x.codegen.version2;

import java.util.List;

public class ConstructorDeclaration extends MethodDefinition {

    public ConstructorDeclaration(List<ParameterExpression> parameters, Block body) {
        super(null, null, parameters, body);
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        throw new UnsupportedOperationException("This should be called by the Class Implementor!");
    }
}
