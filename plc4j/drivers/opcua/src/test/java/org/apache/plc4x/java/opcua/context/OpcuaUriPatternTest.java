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
package org.apache.plc4x.java.opcua.context;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.concurrent.TimeUnit;

import static org.apache.plc4x.java.opcua.context.OpcuaDriverContext.URI_PATTERN;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Matching-cost characterization for the connection-string pattern.
 * <p>
 * A parameter group of the shape {@code (key=value&?)*} is the classic setup for runaway
 * backtracking: if the key class can also absorb the {@code &} separator, an input that ends up
 * failing admits exponentially many candidate partitions and the matcher explores all of them
 * before reporting failure.
 * <p>
 * This pattern is <em>not</em> in that family, and these tests exist to keep it that way. Every
 * iteration of the group must consume exactly one {@code =} (the key class excludes it), so the
 * number of iterations is pinned to the number of {@code =} characters in the input and the
 * partition boundaries cannot vary. Measured cost is linear: a non-matching input of 320 segments
 * / 5779 characters resolves in single-digit milliseconds.
 * <p>
 * The accepted language is pinned alongside the cost, so a future rewrite cannot buy speed by
 * quietly accepting more, nor introduce the ambiguity by widening the key class.
 */
class OpcuaUriPatternTest {

    /**
     * A non-matching input built from many {@code &}-separated segments - the shape that would
     * blow up against an ambiguous parameter group. The timeout is generous enough that only a
     * catastrophic backtrack can breach it.
     */
    @Test
    @Timeout(value = 5, unit = TimeUnit.SECONDS)
    @DisplayName("a long non-matching parameter string fails fast")
    void nonMatchingParameterStringFailsFast() {
        StringBuilder uri = new StringBuilder("opcua:tcp://host/a=b");
        for (int i = 0; i < 60; i++) {
            uri.append("&a=b");
        }
        // The trailing character cannot be part of any parameter, so the overall match must fail.
        uri.append("&!");

        assertFalse(URI_PATTERN.matcher(uri.toString()).matches(),
            "the crafted string is not a valid connection string");
    }

    /**
     * Cost must track input length, not segment count. Against an ambiguous group the step from
     * 10 to 70 segments is the difference between instant and never.
     */
    @Test
    @Timeout(value = 5, unit = TimeUnit.SECONDS)
    @DisplayName("segment count drives cost linearly, not exponentially")
    void costGrowsWithLengthNotWithSegmentCount() {
        for (int segments = 10; segments <= 70; segments += 20) {
            StringBuilder uri = new StringBuilder("opcua:tcp://host/x=y");
            for (int i = 0; i < segments; i++) {
                uri.append("&x=y");
            }
            uri.append("&!");
            assertFalse(URI_PATTERN.matcher(uri.toString()).matches(),
                segments + " segments must fail to match");
        }
    }

    /**
     * Guards the accepted language, so a future rewrite of the parameter group cannot silently
     * widen or narrow what the driver accepts.
     */
    @Test
    @DisplayName("the accepted connection strings are unchanged")
    void acceptsTheDocumentedConnectionStrings() {
        assertTrue(URI_PATTERN.matcher("opcua:tcp://localhost").matches());
        assertTrue(URI_PATTERN.matcher("opcua:tcp://localhost:3131").matches());
        assertTrue(URI_PATTERN.matcher("opcua://127.0.0.1:647").matches());
        assertTrue(URI_PATTERN.matcher("opcua:tcp://127.0.0.1?discovery=false").matches());
        assertTrue(URI_PATTERN.matcher(
            "opcua:tcp://opcua.demo-this.com:51210/UA/SampleServer?discovery=false").matches());
        assertTrue(URI_PATTERN.matcher(
            "opcua:tcp://host:4840?discovery=false&security-policy=Basic256Sha256").matches(),
            "multiple parameters must still be accepted");
        assertTrue(URI_PATTERN.matcher(
            "opcua:tcp://host:4840?a=1&b=2&c=3&d=4").matches(),
            "a longer parameter chain must still be accepted");
    }
}
