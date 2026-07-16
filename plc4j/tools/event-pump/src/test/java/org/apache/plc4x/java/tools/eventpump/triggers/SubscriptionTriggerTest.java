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

package org.apache.plc4x.java.tools.eventpump.triggers;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for SubscriptionTrigger.
 */
class SubscriptionTriggerTest {

    @Test
    void testConstructor() {
        // Act
        SubscriptionTrigger trigger = new SubscriptionTrigger("testTag", "MAIN.test");

        // Assert
        assertEquals("testTag", trigger.getTagName());
        assertEquals("MAIN.test", trigger.getTagAddress());
        assertFalse(trigger.isRunning());
    }

    @Test
    void testConstructorWithNullTagName() {
        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> new SubscriptionTrigger(null, "MAIN.test")
        );
        assertTrue(exception.getMessage().contains("Tag name cannot be null or empty"));
    }

    @Test
    void testConstructorWithEmptyTagName() {
        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> new SubscriptionTrigger("", "MAIN.test")
        );
        assertTrue(exception.getMessage().contains("Tag name cannot be null or empty"));
    }

    @Test
    void testConstructorWithWhitespaceTagName() {
        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> new SubscriptionTrigger("   ", "MAIN.test")
        );
        assertTrue(exception.getMessage().contains("Tag name cannot be null or empty"));
    }

    @Test
    void testConstructorWithNullTagAddress() {
        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> new SubscriptionTrigger("testTag", null)
        );
        assertTrue(exception.getMessage().contains("Tag address cannot be null or empty"));
    }

    @Test
    void testConstructorWithEmptyTagAddress() {
        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> new SubscriptionTrigger("testTag", "")
        );
        assertTrue(exception.getMessage().contains("Tag address cannot be null or empty"));
    }

    @Test
    void testConstructorWithWhitespaceTagAddress() {
        // Act & Assert
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> new SubscriptionTrigger("testTag", "   ")
        );
        assertTrue(exception.getMessage().contains("Tag address cannot be null or empty"));
    }

    @Test
    void testGetTagName() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("myTag", "MAIN.value");

        // Act & Assert
        assertEquals("myTag", trigger.getTagName());
    }

    @Test
    void testGetTagAddress() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("myTag", "MAIN.value");

        // Act & Assert
        assertEquals("MAIN.value", trigger.getTagAddress());
    }

    @Test
    void testStartThrowsUnsupportedOperationException() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("testTag", "MAIN.test");
        Trigger.TriggerListener listener = (t) -> {};

        // Act & Assert
        UnsupportedOperationException exception = assertThrows(
            UnsupportedOperationException.class,
            () -> trigger.start(listener)
        );
        assertTrue(exception.getMessage().contains("SubscriptionTrigger is not yet fully implemented"));
        assertTrue(exception.getMessage().contains("Use TimerTrigger"));
    }

    @Test
    void testStop() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("testTag", "MAIN.test");

        // Act - should not throw
        assertDoesNotThrow(() -> trigger.stop());
    }

    @Test
    void testIsRunning() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("testTag", "MAIN.test");

        // Act & Assert - placeholder never runs
        assertFalse(trigger.isRunning());
    }

    @Test
    void testGetType() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("myTag", "MAIN.value");

        // Act
        String type = trigger.getType();

        // Assert
        assertNotNull(type);
        assertTrue(type.contains("Subscription"));
        assertTrue(type.contains("PLACEHOLDER"));
        assertTrue(type.contains("myTag"));
        assertTrue(type.contains("MAIN.value"));
    }

    @Test
    void testClose() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("testTag", "MAIN.test");

        // Act - should not throw
        assertDoesNotThrow(() -> trigger.close());
    }

    @Test
    void testMultipleCalls() {
        // Arrange
        SubscriptionTrigger trigger = new SubscriptionTrigger("testTag", "MAIN.test");

        // Act & Assert - multiple calls should be safe
        assertDoesNotThrow(() -> {
            trigger.stop();
            trigger.stop();
            trigger.close();
            trigger.close();
        });

        assertFalse(trigger.isRunning());
    }
}
