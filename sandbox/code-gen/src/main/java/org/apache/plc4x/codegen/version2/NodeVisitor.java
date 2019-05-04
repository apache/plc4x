package org.apache.plc4x.codegen.version2;

public interface NodeVisitor<T> {

    T visit(ConstantNode constantNode);

    T visit(DeclarationStatement declarationStatement);

    T visit(ParameterExpression parameterExpression);
}
