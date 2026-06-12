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
package org.apache.plc4x.java.transport.can.socketcan;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Tests for {@link SharedCanManager}.
 * <p>
 * Note: Most tests that require actual CAN socket operations cannot run on non-Linux
 * platforms. These tests focus on the reference counting and lifecycle logic that
 * can be tested without native CAN sockets.
 */
class SharedCanManagerTest {

    @Test
    void managerCanBeInstantiated() {
        SharedCanManager manager = new SharedCanManager();
        assertNotNull(manager);
    }

    @Test
    void acquireHandleOnNonLinuxThrowsRuntimeException() {
        String osName = System.getProperty("os.name", "").toLowerCase();
        if (osName.contains("linux")) {
            // Skip on Linux - would actually try to open a CAN socket
            return;
        }

        SharedCanManager manager = new SharedCanManager();

        // On non-Linux, JavaCAN cannot open a CAN socket
        assertThrows(RuntimeException.class, () ->
                manager.acquireHandle("can0"));
    }

    @Test
    void releaseNonExistentHandleIsNoOp() {
        SharedCanManager manager = new SharedCanManager();

        // Create a mock handle to test release of non-tracked handle
        // This should not throw - just log a warning
        SharedCanManager.SharedCanHandle handle =
                new SharedCanManager.SharedCanHandle(null, "nonexistent");

        // Should not throw
        assertDoesNotThrow(() -> manager.releaseHandle(handle));
    }

    @Test
    void sharedCanHandleRefCountStartsAtOne() {
        SharedCanManager.SharedCanHandle handle =
                new SharedCanManager.SharedCanHandle(null, "test0");

        assertEquals(1, handle.getRefCount());
    }

    @Test
    void sharedCanHandleIncrementRefCount() {
        SharedCanManager.SharedCanHandle handle =
                new SharedCanManager.SharedCanHandle(null, "test0");

        handle.incrementRefCount();
        assertEquals(2, handle.getRefCount());

        handle.incrementRefCount();
        assertEquals(3, handle.getRefCount());
    }

    @Test
    void sharedCanHandleDecrementRefCount() {
        SharedCanManager.SharedCanHandle handle =
                new SharedCanManager.SharedCanHandle(null, "test0");

        handle.incrementRefCount(); // refCount = 2

        int newCount = handle.decrementRefCount();
        assertEquals(1, newCount);
        assertEquals(1, handle.getRefCount());
    }

    @Test
    void sharedCanHandleGetInterfaceName() {
        SharedCanManager.SharedCanHandle handle =
                new SharedCanManager.SharedCanHandle(null, "vcan0");

        assertEquals("vcan0", handle.getInterfaceName());
    }

    @Test
    void sharedCanHandleGetChannel() {
        SharedCanManager.SharedCanHandle handle =
                new SharedCanManager.SharedCanHandle(null, "test0");

        assertNull(handle.getChannel());
    }
}
