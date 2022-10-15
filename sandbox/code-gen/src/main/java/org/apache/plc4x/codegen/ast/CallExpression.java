/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.codegen.ast;

import java.util.Arrays;
import java.util.List;

public class CallExpression extends Expression {

    private final Method method;
    private final Node target;
    private final List<Node> arguments;

    /**
     * Static Method ==&gt; target == null
     * @param method .
     * @param target .
     * @param arguments .
     */
    CallExpression(Method method, Node target, Node... arguments) {
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

    @Override
    public <T> T accept(NodeVisitor<T> visitor) {
        return null;
    }

    @Override
    public void write(Generator writer) {
        if (target == null) {
            writer.generateStaticCall(method, arguments);
        } else {
            writer.generateCall(target, method, arguments);
        }
    }
}
