package org.apache.plc4x.codegen.ast;

public class DeclarationStatement extends Statement {

    public final ParameterExpression parameterExpression;
    public final Expression initializer;

    public DeclarationStatement(ParameterExpression parameterExpression, Expression initializer) {
        this.parameterExpression = parameterExpression;
        this.initializer = initializer;
    }

    public ParameterExpression getParameterExpression() {
        return parameterExpression;
    }

    public Expression getInitializer() {
        return initializer;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return visitor.visit(this);
    }

    @Override public void write(Generator generator) {
        if (initializer != null) {
            generator.generateDeclarationWithInitializer(this);
        } else {
            generator.generateDeclaration(this);
        }
    }

}
