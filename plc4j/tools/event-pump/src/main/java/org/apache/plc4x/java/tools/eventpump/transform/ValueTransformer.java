/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.tools.eventpump.transform;

import org.apache.plc4x.java.api.value.PlcValue;

import java.util.Map;

/**
 * Interface for transforming PLC values using expressions.
 * <p>
 * Implementations can provide different expression languages, from simple arithmetic
 * to more sophisticated languages like SpEL.
 * <p>
 * The transformer is invoked with an expression string and a context containing
 * available variables (typically the tag values from a PLC read).
 * <p>
 * Example implementations:
 * - SimpleExpressionEvaluator: Basic arithmetic and comparisons
 * - SpELValueTransformer: Spring Expression Language support
 */
public interface ValueTransformer {

    /**
     * Transform a value using the given expression.
     *
     * @param expression The expression to evaluate (e.g., "value * 1.8 + 32")
     * @param context A map of variable names to their PlcValue objects
     * @return The transformed value
     * @throws TransformException if the expression cannot be evaluated
     */
    PlcValue transform(String expression, Map<String, PlcValue> context) throws TransformException;

    /**
     * Check if this transformer supports the given expression.
     * This allows for fallback to simpler transformers if needed.
     *
     * @param expression The expression to check
     * @return true if this transformer can handle the expression
     */
    default boolean supports(String expression) {
        return true;
    }

    /**
     * Get the name of this transformer implementation.
     * Used for registration and debugging.
     *
     * @return The transformer name (e.g., "simple", "spel")
     */
    String getName();
}
