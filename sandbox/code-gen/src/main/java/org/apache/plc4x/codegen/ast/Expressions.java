package org.apache.plc4x.codegen.ast;

/**
 * General Factory method to use.
 */
public class Expressions {

    /**
     * Assign a value to a target (field or parameter)
     * @param target Where the value is assigned
     * @param value What to assign
     * @return Assignment Expression
     */
    public static Expression assignment(Expression target, Node value) {
        return new AssignementExpression(target, value);
    }

    /**
     * Base for all Binary Expression, i.e., Expressions which take
     * two inputs and return one Output.
     * Examples are Comparators, Math, ... .
     * @param type Type of the returned expression
     * @param left
     * @param right
     * @param op
     * @return
     */
    public static Expression binaryExpression(TypeNode type, Node left, Node right, Operation op) {
        return new BinaryExpression(type, left, right, op);
    }

    /**
     * A Block of code.
     * @param statements
     * @return
     */
    public static Statement block(List<Node> statements) {
        return new Block(statements);
    }

    /**
     * A block of code.
     * @param statements
     * @return
     */
    public static Statement block(Node... statements) {
        return new Block(statements);
    }

    /**
     * Regular (dynamic) call.
     * @param instance
     * @param method
     * @param arguments
     */
    public static Expression call(Node instance, Method method, Arguments... arguments) {
        return new CallExpression(method, target, arguments);
    }

    /**
     * Static call (call to a static method)
     * 
     * TODO check if Method is static
     * 
     * @param method
     * @param arguments
     */
    public static Expression staticCall(Method method, Arguments... arguments) {
        return new CallExpression(method, null, arguments);
    }

    /**
     * Conditional Statement of the form
     * <code>
     * if (cond1) {
     *  ...
     * } else if (cond2) {
     *  ...
     * } else {
     *  ...
     * }
     * </code>
     */
    public static Statement conditionalStatement(List<Expression> condition, List<Block> blocks) {
        return new IfStatement(condition, blocks);
    }
}
