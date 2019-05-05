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
