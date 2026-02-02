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

package org.apache.plc4x.java.transports.api;

import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

class DefaultTransportManagerTest {

    @Test
    void testDefaultConstructor() {
        DefaultTransportManager manager = new DefaultTransportManager();
        assertNotNull(manager);
        assertNotNull(manager.classLoader);
    }

    @Test
    void testConstructorWithClassLoader() {
        ClassLoader classLoader = Thread.currentThread().getContextClassLoader();
        DefaultTransportManager manager = new DefaultTransportManager(classLoader);

        assertNotNull(manager);
        assertEquals(classLoader, manager.classLoader);
    }

    @Test
    void testGetTransport_returnsEmptyWhenNoTransportsAvailable() {
        DefaultTransportManager manager = new DefaultTransportManager();

        Optional<Transport> result = manager.getTransport("tcp");
        assertFalse(result.isPresent());
    }

    @Test
    void testGetTransport_returnsEmptyForUnknownCode() {
        DefaultTransportManager manager = new DefaultTransportManager();

        Optional<Transport> result = manager.getTransport("unknown-transport-code");
        assertFalse(result.isPresent());
    }

    @Test
    void testGetTransport_returnsEmptyForNullCode() {
        DefaultTransportManager manager = new DefaultTransportManager();

        Optional<Transport> result = manager.getTransport(null);
        assertFalse(result.isPresent());
    }

    @Test
    void testConstructorWithCustomClassLoader() {
        // Create a custom classloader that won't find any transports
        ClassLoader emptyClassLoader = new ClassLoader(null) {
            @Override
            public Class<?> loadClass(String name) throws ClassNotFoundException {
                throw new ClassNotFoundException(name);
            }
        };

        DefaultTransportManager manager = new DefaultTransportManager(emptyClassLoader);
        assertNotNull(manager);
        assertEquals(emptyClassLoader, manager.classLoader);
    }
}
