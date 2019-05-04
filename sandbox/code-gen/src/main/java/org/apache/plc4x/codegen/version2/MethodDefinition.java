package org.apache.plc4x.codegen.version2;

import java.util.List;

public class MethodDefinition implements Node {

    private final String name;
    private final TypeNode resultType;
    private final List<ParameterExpression> parameters;
    private final Block body;

    public MethodDefinition(String name, TypeNode resultType, List<ParameterExpression> parameters, Block body) {
        this.name = name;
        this.resultType = resultType;
        this.parameters = parameters;
        this.body = body;
    }

    public String getName() {
        return name;
    }

    public TypeNode getResultType() {
        return resultType;
    }

    public List<ParameterExpression> getParameters() {
        return parameters;
    }

    public Block getBody() {
        return body;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generate(this);
    }
}
