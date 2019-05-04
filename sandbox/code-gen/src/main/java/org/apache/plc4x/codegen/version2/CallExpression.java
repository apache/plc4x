package org.apache.plc4x.codegen.version2;

import java.util.Arrays;
import java.util.List;

public class CallExpression extends Expression {

    private final Method method;
    private final Node target;
    private final List<Node> arguments;

    /**
     * Static Method ==> target == null
     * @param method
     * @param target
     * @param arguments
     */
    public CallExpression(Method method, Node target, Node... arguments) {
        super(method.getReturnType());
        this.method = method;
        this.target = target;
        this.arguments = Arrays.asList(arguments);
    }

    public Method getMethod() {
        return method;
    }

    public Node getTarget() {
        return target;
    }

    @Override public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override public void write(Generator writer) {
        if (target == null) {
            writer.generateStaticCall(method, arguments);
        } else {
            writer.generateCall(target, method, arguments);
        }
    }
}
