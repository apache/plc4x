package org.apache.plc4x.codegen.version2;

import java.util.Collections;
import java.util.List;
import java.util.Set;

public class MethodDefinition implements Node {

    private final Set<Modifier> modifiers;
    private final String name;
    private final TypeNode resultType;
    private final List<ParameterExpression> parameters;
    private final Block body;

    public MethodDefinition(Set<Modifier> modifiers, String name, TypeNode resultType, List<ParameterExpression> parameters, Block body) {
        this.modifiers = modifiers;
        this.name = name;
        this.resultType = resultType;
        this.parameters = parameters;
        this.body = body;
    }

    public MethodDefinition(String name, TypeNode resultType, List<ParameterExpression> parameters, Block body) {
        this(Collections.emptySet(), name, resultType, parameters, body);
    }

    public Set<Modifier> getModifiers() {
        return modifiers;
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

    public enum Modifier {
        STATIC,
        PRIVATE
    }

}
