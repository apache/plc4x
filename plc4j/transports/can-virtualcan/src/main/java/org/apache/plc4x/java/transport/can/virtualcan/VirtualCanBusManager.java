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
package org.apache.plc4x.java.transport.can.virtualcan;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArraySet;

/**
 * Static manager that coordinates virtual CAN buses.
 * <p>
 * Each bus is identified by name and holds a set of connected
 * {@link VirtualCanTransportInstance} objects. When one instance writes a frame,
 * the manager broadcasts it to every other instance on the same bus.
 * <p>
 * This class is thread-safe. Bus sets use {@link CopyOnWriteArraySet} so that
 * broadcasts can iterate without holding a lock while instances connect or disconnect.
 */
public final class VirtualCanBusManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(VirtualCanBusManager.class);

    /** Map of bus name to the set of instances currently connected to that bus. */
    private static final ConcurrentHashMap<String, Set<VirtualCanTransportInstance>> BUSES =
            new ConcurrentHashMap<>();

    // Utility class — no instances
    private VirtualCanBusManager() {
    }

    /**
     * Connects a transport instance to the named virtual bus.
     * Creates the bus if it does not already exist.
     *
     * @param busName  the name of the virtual bus
     * @param instance the instance to add
     */
    public static void connect(String busName, VirtualCanTransportInstance instance) {
        BUSES.computeIfAbsent(busName, k -> new CopyOnWriteArraySet<>()).add(instance);
        LOGGER.debug("Instance connected to virtual CAN bus '{}' (bus size: {})",
                busName, BUSES.get(busName).size());
    }

    /**
     * Disconnects a transport instance from the named virtual bus.
     * Removes the bus entirely when the last instance disconnects.
     *
     * @param busName  the name of the virtual bus
     * @param instance the instance to remove
     */
    public static void disconnect(String busName, VirtualCanTransportInstance instance) {
        BUSES.computeIfPresent(busName, (k, instances) -> {
            instances.remove(instance);
            LOGGER.debug("Instance disconnected from virtual CAN bus '{}' (bus size: {})",
                    busName, instances.size());
            // Remove the bus when the last participant leaves
            return instances.isEmpty() ? null : instances;
        });
    }

    /**
     * Broadcasts frame bytes to every instance on the named bus except the sender.
     * <p>
     * Delivery is synchronous: each receiver's {@code onFrameReceived} is called
     * on the sender's thread. This keeps the virtual transport simple and
     * deterministic for testing.
     *
     * @param busName    the name of the virtual bus
     * @param sender     the instance that originated the frame (excluded from delivery)
     * @param frameBytes the raw frame bytes to deliver
     */
    public static void broadcast(String busName, VirtualCanTransportInstance sender, byte[] frameBytes) {
        Set<VirtualCanTransportInstance> instances = BUSES.get(busName);
        if (instances == null) {
            // No bus or no other participants — nothing to deliver
            return;
        }
        for (VirtualCanTransportInstance instance : instances) {
            if (instance != sender) {
                try {
                    instance.onFrameReceived(frameBytes);
                } catch (Exception e) {
                    LOGGER.warn("Error delivering frame to instance on bus '{}': {}",
                            busName, e.getMessage(), e);
                }
            }
        }
    }

    /**
     * Removes all buses and their instances.
     * This is intended for test cleanup only and is package-private.
     */
    static void reset() {
        BUSES.clear();
        LOGGER.debug("VirtualCanBusManager reset — all buses cleared");
    }
}
