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
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

class ValueTransformerRegistryTest {

    private ValueTransformerRegistry registry;

    @BeforeEach
    void setUp() {
        // Create a fresh registry instance for each test to ensure isolation
        registry = ValueTransformerRegistry.createDefault();
    }

    @Test
    void testCreateDefaultRegistersSimpleTransformer() {
        ValueTransformer defaultTransformer = registry.getDefault();

        assertNotNull(defaultTransformer);
        assertEquals("simple", defaultTransformer.getName());
    }

    @Test
    void testEmptyRegistryCanBeCreated() {
        ValueTransformerRegistry emptyRegistry = new ValueTransformerRegistry();

        // Empty registry should have no transformers
        assertTrue(emptyRegistry.getRegisteredNames().isEmpty());
    }

    @Test
    void testCanRegisterToEmptyRegistry() {
        ValueTransformerRegistry emptyRegistry = new ValueTransformerRegistry();
        MockValueTransformer mockTransformer = new MockValueTransformer("test");

        emptyRegistry.register(mockTransformer);

        assertEquals(1, emptyRegistry.getRegisteredNames().size());
        assertSame(mockTransformer, emptyRegistry.get("test"));
    }

    @Test
    void testRegisterNewTransformer() {
        MockValueTransformer mockTransformer = new MockValueTransformer("test");

        registry.register(mockTransformer);

        ValueTransformer retrieved = registry.get("test");
        assertNotNull(retrieved);
        assertEquals("test", retrieved.getName());
        assertSame(mockTransformer, retrieved);
    }

    @Test
    void testRegisterNullTransformerThrowsException() {
        assertThrows(IllegalArgumentException.class, () -> {
            registry.register(null);
        });
    }

    @Test
    void testRegisterTransformerWithNullNameThrowsException() {
        MockValueTransformer mockTransformer = new MockValueTransformer(null);

        assertThrows(IllegalArgumentException.class, () -> {
            registry.register(mockTransformer);
        });
    }

    @Test
    void testRegisterTransformerWithEmptyNameThrowsException() {
        MockValueTransformer mockTransformer = new MockValueTransformer("");

        assertThrows(IllegalArgumentException.class, () -> {
            registry.register(mockTransformer);
        });
    }

    @Test
    void testRegisterDuplicateTransformerKeepsExisting() {
        MockValueTransformer first = new MockValueTransformer("duplicate");
        MockValueTransformer second = new MockValueTransformer("duplicate");

        registry.register(first);
        registry.register(second);

        ValueTransformer retrieved = registry.get("duplicate");
        assertSame(first, retrieved, "Should keep the first registered transformer");
    }

    @Test
    void testGetNonExistentTransformerReturnsNull() {
        ValueTransformer result = registry.get("nonexistent");

        assertNull(result);
    }

    @Test
    void testSetDefault() {
        MockValueTransformer mockTransformer = new MockValueTransformer("custom");
        registry.register(mockTransformer);

        registry.setDefault("custom");

        ValueTransformer defaultTransformer = registry.getDefault();
        assertEquals("custom", defaultTransformer.getName());
        assertSame(mockTransformer, defaultTransformer);
    }

    @Test
    void testSetDefaultToNonExistentTransformerThrowsException() {
        assertThrows(IllegalArgumentException.class, () -> {
            registry.setDefault("nonexistent");
        });
    }

    @Test
    void testGetDefaultWhenDefaultNotRegisteredThrowsException() {
        // Create an empty registry and try to get default without registering anything
        ValueTransformerRegistry emptyRegistry = new ValueTransformerRegistry();

        assertThrows(IllegalStateException.class, () -> {
            emptyRegistry.getDefault();
        });
    }

    @Test
    void testUnregister() {
        MockValueTransformer mockTransformer = new MockValueTransformer("toremove");
        registry.register(mockTransformer);

        boolean removed = registry.unregister("toremove");

        assertTrue(removed);
        assertNull(registry.get("toremove"));
    }

    @Test
    void testUnregisterNonExistentReturnsFalse() {
        boolean removed = registry.unregister("nonexistent");

        assertFalse(removed);
    }

    @Test
    void testUnregisterDefaultTransformerThrowsException() {
        assertThrows(IllegalStateException.class, () -> {
            registry.unregister("simple");
        });
    }

    @Test
    void testGetRegisteredNames() {
        MockValueTransformer mock1 = new MockValueTransformer("test1");
        MockValueTransformer mock2 = new MockValueTransformer("test2");

        registry.register(mock1);
        registry.register(mock2);

        java.util.Set<String> names = registry.getRegisteredNames();

        assertTrue(names.contains("simple"));
        assertTrue(names.contains("test1"));
        assertTrue(names.contains("test2"));
    }

    @Test
    void testGetRegisteredNamesIsUnmodifiable() {
        java.util.Set<String> names = registry.getRegisteredNames();

        assertThrows(UnsupportedOperationException.class, () -> {
            names.add("new");
        });
    }

    /**
     * Mock ValueTransformer for testing
     */
    private static class MockValueTransformer implements ValueTransformer {
        private final String name;

        public MockValueTransformer(String name) {
            this.name = name;
        }

        @Override
        public PlcValue transform(String expression, Map<String, PlcValue> context) throws TransformException {
            return null; // Not used in registry tests
        }

        @Override
        public String getName() {
            return name;
        }
    }
}
