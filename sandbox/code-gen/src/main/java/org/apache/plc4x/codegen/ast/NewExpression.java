package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.List;

public class NewExpression extends Expression {

    private List<Node> arguments;

    public NewExpression(TypeNode myClazz, Node... arguments) {
        super(myClazz);
        this.arguments = Arrays.asList(arguments);
    }

    public List<Node> getArguments() {
        return arguments;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        writer.generate(this);
    }
}
