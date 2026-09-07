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

package org.apache.plc4x.java.spi.config.exceptions;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ConfigurationExceptionTest {

    @Test
    void testDefaultConstructor() {
        ConfigurationException exception = new ConfigurationException();

        assertNull(exception.getMessage());
        assertNull(exception.getCause());
    }

    @Test
    void testMessageConstructor() {
        String message = "Test configuration error";
        ConfigurationException exception = new ConfigurationException(message);

        assertEquals(message, exception.getMessage());
        assertNull(exception.getCause());
    }

    @Test
    void testMessageAndCauseConstructor() {
        String message = "Test configuration error";
        Throwable cause = new IllegalArgumentException("Invalid parameter");
        ConfigurationException exception = new ConfigurationException(message, cause);

        assertEquals(message, exception.getMessage());
        assertEquals(cause, exception.getCause());
    }

    @Test
    void testExceptionInheritance() {
        ConfigurationException exception = new ConfigurationException("Test");

        assertInstanceOf(RuntimeException.class, exception);
        assertInstanceOf(Exception.class, exception);
        assertInstanceOf(Throwable.class, exception);
    }

    @Test
    void testExceptionThrowing() {
        String errorMessage = "Configuration validation failed";

        assertThrows(ConfigurationException.class, () -> {
            throw new ConfigurationException(errorMessage);
        });

        ConfigurationException exception = assertThrows(ConfigurationException.class, () -> {
            throw new ConfigurationException(errorMessage);
        });

        assertEquals(errorMessage, exception.getMessage());
    }

    @Test
    void testExceptionWithNullMessage() {
        ConfigurationException exception = new ConfigurationException(null);

        assertNull(exception.getMessage());
        assertNull(exception.getCause());
    }

    @Test
    void testExceptionWithNullCause() {
        String message = "Test message";
        ConfigurationException exception = new ConfigurationException(message, null);

        assertEquals(message, exception.getMessage());
        assertNull(exception.getCause());
    }

    @Test
    void testExceptionChaining() {
        IllegalArgumentException rootCause = new IllegalArgumentException("Root cause");
        RuntimeException intermediateCause = new RuntimeException("Intermediate", rootCause);
        ConfigurationException topLevelException = new ConfigurationException("Top level", intermediateCause);

        assertEquals("Top level", topLevelException.getMessage());
        assertEquals(intermediateCause, topLevelException.getCause());
        assertEquals(rootCause, topLevelException.getCause().getCause());
    }

    @Test
    void testStackTrace() {
        ConfigurationException exception = new ConfigurationException("Test exception");

        assertNotNull(exception.getStackTrace());
        assertTrue(exception.getStackTrace().length > 0);

        // Test that the stack trace contains this test method
        boolean foundTestMethod = false;
        for (StackTraceElement element : exception.getStackTrace()) {
            if (element.getMethodName().equals("testStackTrace")) {
                foundTestMethod = true;
                break;
            }
        }
        assertTrue(foundTestMethod, "Stack trace should contain the test method");
    }
}