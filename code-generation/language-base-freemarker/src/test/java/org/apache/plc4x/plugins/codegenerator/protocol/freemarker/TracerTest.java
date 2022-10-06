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
package org.apache.plc4x.plugins.codegenerator.protocol.freemarker;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.HashMap;

import static org.junit.jupiter.api.Assertions.*;

class TracerTest {

    @BeforeEach
    void setUp() {
        Tracer.ENABLED = false;
    }

    @Test
    void start() {
        Tracer tracer = Tracer.start("random");
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals("/*random*/", tracer.toString());
    }

    @Test
    void dive() {
        Tracer tracer = Tracer.start("random");
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals("/*random*/", tracer.toString());
        tracer = tracer.dive("something other");
        assertEquals("/*random/something other*/", tracer.toString());
        tracer = tracer.dive("something third");
        assertEquals("/*random/something other/something third*/", tracer.toString());
        String functionExpression = "someFunction(1+2)";
        String someExpression = tracer + functionExpression;
        Tracer subTracer = Tracer.start("subtracing").dive("even deeper");
        String tracedExpression = subTracer + someExpression;
        assertEquals("/*subtracing/even deeper*//*random/something other/something third*/" + functionExpression, tracedExpression);
    }

    @Test
    void removeTraces() {
        Tracer tracer = Tracer.start("random");
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals("/*random*/", tracer.toString());
        String unTracedExpression = "someExpression(123,123)";
        String expression = tracer + unTracedExpression;
        assertEquals(unTracedExpression, tracer.removeTraces(expression));
        Tracer subTracer = Tracer.start("subtracing").dive("even deeper");
        String tracedExpression = subTracer + expression;
        assertEquals(unTracedExpression, tracer.removeTraces(tracedExpression));
    }

    @Test
    void extractTraces() {
        Tracer tracer = Tracer.start("random");
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        String trace = "/*random*/";
        assertEquals(trace, tracer.toString());
        String unTracedExpression = "someExpression(123,123)";
        String expression = tracer + unTracedExpression;
        assertEquals(unTracedExpression, tracer.removeTraces(expression));
        assertEquals(trace, tracer.extractTraces(expression));
        Tracer subTracer = Tracer.start("subtracing").dive("even deeper");
        String tracedExpression = subTracer + expression;
        assertEquals("/*subtracing/even deeper*//*random*/", tracer.extractTraces(tracedExpression));
    }

    @Test
    void separator() {
        String separator = "everythingIsPossible";
        Tracer tracer = new Tracer("random") {
            @Override
            protected String separator() {
                return separator;
            }
        };
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals("/*random*/", tracer.toString());
        tracer = tracer.dive("something other");
        assertEquals("/*random" + separator + "something other*/", tracer.toString());
    }

    @Test
    void prefix() {
        String prefix = "everythingIsPossible";
        Tracer tracer = new Tracer("random") {
            @Override
            protected String prefix() {
                return prefix;
            }
        };
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals(prefix + "random*/", tracer.toString());
        tracer = tracer.dive("something other");
        assertEquals(prefix + "random/something other*/", tracer.toString());
    }

    @Test
    void suffix() {
        String suffix = "everythingIsPossible";
        Tracer tracer = new Tracer("random") {
            @Override
            protected String suffix() {
                return suffix;
            }
        };
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals("/*random" + suffix, tracer.toString());
        tracer = tracer.dive("something other");
        assertEquals("/*random/something other" + suffix, tracer.toString());
    }

    @Test
    void isEnabled() {
        Tracer.ENABLED = false;
        Tracer tracer = Tracer.start("doesn't matter");
        assertFalse(tracer.isEnabled());
        Tracer.ENABLED = true;
        assertTrue(tracer.isEnabled());
        Tracer.ENABLED = false;
        assertFalse(tracer.isEnabled());
    }

    @Test
    void testToString() {
        Tracer tracer = Tracer.start("random");
        assertNotNull(tracer);
        assertEquals("", tracer.toString());
        Tracer.ENABLED = true;
        assertEquals("/*random*/", tracer.toString());
    }
}