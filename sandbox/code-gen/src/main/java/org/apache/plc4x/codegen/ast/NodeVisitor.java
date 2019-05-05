package org.apache.plc4x.codegen.ast;

public interface NodeVisitor<T> {

    T visit(ConstantNode constantNode);

    T visit(DeclarationStatement declarationStatement);

    T visit(ParameterExpression parameterExpression);
}
