package org.apache.plc4x.codegen.version2;

import java.util.List;

/**
 * Method Declaration
 */
public class Method implements Node {

    private final TypeNode type;    // Type where this thing is defined on
    private final String name;
    private final TypeNode returnType;
    private final List<TypeNode> parameterTypes;
    private final List<TypeNode> expressionTypes;

    public Method(TypeNode type, String name, TypeNode returnType, List<TypeNode> parameterTypes, List<TypeNode> expressionTypes) {
        this.type = type;
        this.name = name;
        this.returnType = returnType;
        this.parameterTypes = parameterTypes;
        this.expressionTypes = expressionTypes;
    }

    public TypeNode getType() {
        return type;
    }

    public String getName() {
        return name;
    }

    public TypeNode getReturnType() {
        return returnType;
    }

    public List<TypeNode> getParameterTypes() {
        return parameterTypes;
    }

    public List<TypeNode> getExpressionTypes() {
        return expressionTypes;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {

    }
}
