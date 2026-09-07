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

package org.apache.plc4x.java.spi.transports.api.exceptions;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TransportExceptionTest {

    @Test
    void testConstructorWithMessage() {
        String message = "Transport error occurred";
        TransportException exception = new TransportException(message);

        assertEquals(message, exception.getMessage());
        assertNull(exception.getCause());
    }

    @Test
    void testConstructorWithMessageAndCause() {
        String message = "Transport error occurred";
        Throwable cause = new RuntimeException("Root cause");
        TransportException exception = new TransportException(message, cause);

        assertEquals(message, exception.getMessage());
        assertEquals(cause, exception.getCause());
    }

    @Test
    void testExceptionCanBeThrown() {
        assertThrows(TransportException.class, () -> {
            throw new TransportException("Test exception");
        });
    }

    @Test
    void testExceptionCanBeCaught() {
        try {
            throw new TransportException("Test message");
        } catch (TransportException e) {
            assertEquals("Test message", e.getMessage());
        }
    }
}
