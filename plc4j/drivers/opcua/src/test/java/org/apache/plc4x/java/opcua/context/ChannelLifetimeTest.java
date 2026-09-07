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

import org.apache.plc4x.java.opcua.config.OpcuaConfiguration;
import org.apache.plc4x.java.spi.config.ConfigurationFactory;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import static org.apache.plc4x.java.opcua.context.SecureChannel.effectiveChannelLifetime;
import static org.apache.plc4x.java.opcua.context.SecureChannel.keepAliveInterval;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Reconciliation of the channel lifetime the server returns against the one that was requested.
 * <p>
 * The renewal schedule is derived from this value and runs on an executor shared by every OPC UA
 * connection in the JVM, so the lifetime decides how much of a shared resource one peer can
 * command. The reconciliation is what keeps that decision on this side of the wire.
 */
class ChannelLifetimeTest {

    private static final long REQUESTED = 3_600_000L;

    /** The shipped default of {@code min-channel-lifetime-ms}. */
    private static final long MIN = 5_000L;

    /**
     * A server may revise the requested lifetime downwards. Anything else - zero, negative, or a
     * value larger than we offered - is not a revision, so the requested lifetime stands.
     */
    @ParameterizedTest
    @ValueSource(longs = {0L, -1L, Long.MIN_VALUE, REQUESTED + 1, Long.MAX_VALUE})
    @DisplayName("a lifetime outside the offer is not adopted")
    void lifetimeOutsideTheOfferFallsBackToTheRequestedValue(long serverValue) {
        assertEquals(REQUESTED, effectiveChannelLifetime(serverValue, REQUESTED, MIN));
    }

    /**
     * The concrete values from the scenario this guards against: a lifetime of 1 ms would schedule
     * a renewal - a full secure-channel exchange - every millisecond, for as long as the connection
     * lives.
     */
    @ParameterizedTest
    @ValueSource(longs = {1L, 2L, 100L, 999L, MIN - 1})
    @DisplayName("an implausibly short lifetime is raised to the floor")
    void shortLifetimeIsRaisedToTheFloor(long serverValue) {
        assertEquals(MIN, effectiveChannelLifetime(serverValue, REQUESTED, MIN));
    }

    @Test
    @DisplayName("a plausible revision downwards is honoured")
    void plausibleRevisionIsHonoured() {
        assertEquals(600_000L, effectiveChannelLifetime(600_000L, REQUESTED, MIN));
        assertEquals(REQUESTED, effectiveChannelLifetime(REQUESTED, REQUESTED, MIN));
        assertEquals(MIN, effectiveChannelLifetime(MIN, REQUESTED, MIN));
    }

    /**
     * The floor must not override the operator. Someone who deliberately configures a short
     * {@code channel-lifetime-ms} gets it - the reconciliation only ever declines to go below what was
     * asked for, never above it.
     */
    @Test
    @DisplayName("a deliberately short configured lifetime is not overridden")
    void configuredLifetimeShorterThanTheFloorIsHonoured() {
        long requested = 1_000L;

        assertEquals(requested, effectiveChannelLifetime(requested, requested, MIN));
        assertEquals(requested, effectiveChannelLifetime(0L, requested, MIN));
        assertEquals(requested, effectiveChannelLifetime(500L, requested, MIN),
            "the floor is bounded by the requested lifetime, so it cannot exceed it");
    }

    /**
     * The renewal period is what actually reaches the scheduler, and the scheduler rejects a
     * non-positive period from inside a completion stage - a failure that is easy to lose.
     */
    @ParameterizedTest
    @ValueSource(longs = {0L, 1L, 2L, 1_000L, 3_600_000L})
    @DisplayName("the derived renewal period is always schedulable")
    void derivedPeriodIsAlwaysSchedulable(long lifetime) {
        long period = keepAliveInterval(lifetime);

        assertTrue(period > 0, "period " + period + " would be rejected by the scheduler");

        ScheduledExecutorService executor = Executors.newSingleThreadScheduledExecutor();
        try {
            assertDoesNotThrow(() -> executor.scheduleAtFixedRate(
                () -> {
                }, period, period, TimeUnit.MILLISECONDS).cancel(true));
        } finally {
            executor.shutdownNow();
        }
    }

    /**
     * End to end: whatever the peer sends, the interval that reaches the shared executor leaves
     * real time between renewals.
     */
    @ParameterizedTest
    @ValueSource(longs = {Long.MIN_VALUE, -1L, 0L, 1L, 100L, REQUESTED, Long.MAX_VALUE})
    @DisplayName("no server-supplied value produces a renewal treadmill")
    void noServerValueProducesATreadmill(long serverValue) {
        long period = keepAliveInterval(effectiveChannelLifetime(serverValue, REQUESTED, MIN));

        assertTrue(period >= (long) Math.ceil(MIN * 0.75f),
            "server value " + serverValue + " produced a renewal period of " + period + " ms");
    }

    /**
     * The point of making the minimum configurable: an operator who has a server that genuinely
     * needs faster renewal can lower it and get exactly what the server asked for. The warning the
     * driver logs names this parameter precisely so that this is discoverable.
     */
    @Test
    @DisplayName("lowering the configured minimum honours the server")
    void loweringTheMinimumHonoursTheServer() {
        assertEquals(1_000L, effectiveChannelLifetime(1_000L, REQUESTED, 1_000L),
            "with the minimum lowered to the server's value, that value is used unchanged");
        assertEquals(100L, effectiveChannelLifetime(100L, REQUESTED, 0L),
            "a minimum of zero disables the floor entirely - the operator's choice");
    }

    /**
     * A default that failed to resolve would arrive as 0, which disables the floor completely and
     * silently. That failure looks identical to working correctly until a server exploits it, so
     * the default is pinned here rather than assumed.
     */
    @Test
    @DisplayName("min-channel-lifetime-ms defaults to 5000 ms")
    void minimumChannelLifetimeDefaultResolves() throws Exception {
        OpcuaConfiguration defaults =
            new ConfigurationFactory().createConfiguration(OpcuaConfiguration.class, "");

        assertEquals(MIN, defaults.getMinChannelLifetime(),
            "an unconfigured connection must get a real floor, never 0");
    }

    @Test
    @DisplayName("min-channel-lifetime-ms can be lowered from the connection string")
    void minimumChannelLifetimeCanBeLowered() throws Exception {
        OpcuaConfiguration lowered = new ConfigurationFactory()
            .createConfiguration(OpcuaConfiguration.class, "min-channel-lifetime-ms=250");

        assertEquals(250L, lowered.getMinChannelLifetime());
    }
}
