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
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertTrue;

class ValueTransformerTest {

    @Test
    void testDefaultSupportsMethod() throws TransformException {
        ValueTransformer transformer = new ValueTransformer() {
            @Override
            public PlcValue transform(String expression, Map<String, PlcValue> context) throws TransformException {
                return null;
            }

            @Override
            public String getName() {
                return "test";
            }
        };

        // Test the default supports() method
        assertTrue(transformer.supports("any expression"));
        assertTrue(transformer.supports(""));
        assertTrue(transformer.supports(null));
    }
}
