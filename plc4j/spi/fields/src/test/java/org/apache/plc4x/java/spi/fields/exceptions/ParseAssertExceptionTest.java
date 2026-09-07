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

package org.apache.plc4x.java.spi.fields.exceptions;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ParseAssertExceptionTest {

    @Test
    void testConstructorWithMessage() {
        String message = "Assert failed";
        ParseAssertException exception = new ParseAssertException(message);

        assertEquals(message, exception.getMessage());
        assertNull(exception.getCause());
    }

    @Test
    void testConstructorWithMessageAndCause() {
        String message = "Assert failed";
        Throwable cause = new RuntimeException("Root cause");
        ParseAssertException exception = new ParseAssertException(message, cause);

        assertEquals(message, exception.getMessage());
        assertEquals(cause, exception.getCause());
    }

    @Test
    void testExceptionCanBeThrown() {
        assertThrows(ParseAssertException.class, () -> {
            throw new ParseAssertException("Test exception");
        });
    }

    @Test
    void testExceptionCanBeCaught() {
        try {
            throw new ParseAssertException("Test message");
        } catch (ParseAssertException e) {
            assertEquals("Test message", e.getMessage());
        }
    }
}
