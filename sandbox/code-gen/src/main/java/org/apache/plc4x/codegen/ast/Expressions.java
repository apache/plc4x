/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.codegen.ast;

import java.util.Collections;
import java.util.List;

/**
 * General Factory method to use.
 */
public class Expressions {

    private Expressions() {
        // do not instantiate
    }

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
    public static Expression binaryExpression(TypeNode type, Node left, Node right, BinaryExpression.Operation op) {
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
    public static Expression call(Node instance, Method method, Node... arguments) {
        return new CallExpression(method, instance, arguments);
    }

    /**
     * Static call (call to a static method)
     * 
     * TODO check if Method is static
     * 
     * @param method
     * @param arguments
     */
    public static Expression staticCall(Method method, Node... arguments) {
        return new CallExpression(method, null, arguments);
    }

    /**
     * Simple if-then-else.
     * If no else is needed, orElse can be null.
     */
    public static Statement ifElse(Expression condition, Block then, Block orElse) {
        return new IfStatement(condition, then, orElse);
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

    /**
     * Defines a compile time constant and infers the type based on what java would do.
     */
    public static Expression constant(Object value) {
        return new ConstantExpression(value);
    }

    /**
     * Define a compile time constant and also passes
     * the expected type for usage in the code generation later.
     */
    public static Expression constant(TypeNode type, Object value) {
        return new ConstantExpression(value);
    }

    /**
     * Declares a constant (no field).
     * Variable type is infered from the initializing expression.
     */
    public static Statement declaration(String variable, Expression initializer) {
        return new DeclarationStatement(parameter(variable, initializer.getType()), initializer);
    }

    /**
     * Declares a constant (no field), which is not initialized.
     */
    public static Statement declaration(String variable, TypeNode type) {
        return new DeclarationStatement(parameter(variable, type), null);
    }

    /**
     * Reference to a field in the surrounding class.
     */
    public static Expression field(String name) {
        return new FieldReference(new UnknownType(), name);
    }

    /**
     * Reference to a field on the given target
     */
    public static Expression field(Node target, String name) {
        return new FieldReference(new UnknownType(), name, target);
    }

    /**
     * Adds a line of comment.
     */
    public static Node comment(String comment) {
        return new LineComment(comment);
    }

    /**
     * Simple call to a method which throws no exception and that stuff.
     */
    public static Expression call(Node target, String methodName, Node... arguments) {
        return new CallExpression(
            new Method(UnknownType.INSTANCE, methodName, UnknownType.INSTANCE,
                Collections.<TypeNode>emptyList(), Collections.<ExceptionType>emptyList()),
            target,
            arguments
        );
    }

    /**
     * Reference to a Method, similar than field-reference
     * @return
     */
    public static Method method(TypeNode definingClass, String name, TypeNode returnType, List<TypeNode> parameterTypes, List<ExceptionType> exceptions) {
        return new Method(definingClass, name, returnType, parameterTypes, exceptions);
    }

    /**
     * Declares a variable.
     */
    public static ParameterExpression parameter(String name, TypeNode type) {
        return new ParameterExpression(type, name);
    }
}
