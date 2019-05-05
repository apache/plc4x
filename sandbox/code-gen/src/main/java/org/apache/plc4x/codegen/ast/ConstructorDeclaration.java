package org.apache.plc4x.codegen.ast;

import java.util.List;
import java.util.Set;

public class ConstructorDeclaration extends MethodDefinition {

    public ConstructorDeclaration(Set<Modifier> modifiers, List<ParameterExpression> parameters, Block body) {
        super(modifiers, null, null, parameters, body);
    }

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
