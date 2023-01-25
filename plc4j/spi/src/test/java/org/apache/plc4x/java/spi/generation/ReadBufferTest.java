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
package org.apache.plc4x.java.spi.generation;

import org.apache.plc4x.java.spi.codegen.WithOption;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ReadBufferTest {

    /**
     * Test which makes sure that PLC4X-256 is not happening.
     */
    @Test
    void readString() throws ParseException {
        String value = "abcdef";
        final ReadBuffer buffer = new ReadBufferByteBased(value.getBytes(StandardCharsets.UTF_8));
        String answer = buffer.readString("", value.length() * 8, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        assertEquals(value, answer);
    }
}