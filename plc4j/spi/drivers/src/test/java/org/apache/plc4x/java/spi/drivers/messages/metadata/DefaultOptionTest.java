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
package org.apache.plc4x.java.spi.drivers.messages.metadata;

import org.apache.plc4x.java.api.types.OptionType;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * {@link DefaultOption} gained a {@code secret} field and a seventh constructor parameter, with
 * the original six-argument constructor kept as a delegate. The delegation is the part worth
 * pinning: a wrong argument order there would silently mislabel every option built by the old
 * constructor, and every existing caller uses it.
 */
class DefaultOptionTest {

    private static DefaultOption sixArg() {
        return new DefaultOption("timeout", OptionType.LONG, "request timeout", false, 5000L, "1.0");
    }

    private static DefaultOption sevenArg(boolean secret) {
        return new DefaultOption("timeout", OptionType.LONG, "request timeout", false, 5000L, "1.0", secret);
    }

    @Test
    @DisplayName("the six-argument constructor produces a non-secret option")
    void sixArgConstructorDefaultsToNotSecret() {
        assertFalse(sixArg().isSecret());
    }

    @Test
    @DisplayName("the seven-argument constructor carries the flag through")
    void sevenArgConstructorCarriesTheFlag() {
        assertTrue(sevenArg(true).isSecret());
        assertFalse(sevenArg(false).isSecret());
    }

    @Test
    @DisplayName("the six-argument constructor delegates every other value unchanged")
    void sixArgConstructorDelegatesTheRest() {
        // Guards the delegating call itself: same inputs through both constructors must give the
        // same option, so a transposed argument in the delegate cannot pass unnoticed.
        DefaultOption viaSix = sixArg();
        DefaultOption viaSeven = sevenArg(false);

        assertEquals(viaSeven.getKey(), viaSix.getKey());
        assertEquals(viaSeven.getType(), viaSix.getType());
        assertEquals(viaSeven.getDescription(), viaSix.getDescription());
        assertEquals(viaSeven.isRequired(), viaSix.isRequired());
        assertEquals(viaSeven.getDefaultValue(), viaSix.getDefaultValue());
        assertEquals(viaSeven.getSince(), viaSix.getSince());
        assertEquals(viaSeven.isSecret(), viaSix.isSecret());
    }

    @Test
    @DisplayName("all values arrive where they belong")
    void everyValueLandsInItsOwnAccessor() {
        DefaultOption option = new DefaultOption(
            "password", OptionType.STRING, "the password", true, "hunter2", "0.9", true);

        assertEquals("password", option.getKey());
        assertEquals(OptionType.STRING, option.getType());
        assertEquals("the password", option.getDescription());
        assertTrue(option.isRequired());
        assertEquals(Optional.of("hunter2"), option.getDefaultValue());
        assertEquals(Optional.of("0.9"), option.getSince());
        assertTrue(option.isSecret());
    }

    @Test
    @DisplayName("a secret option may still have no default and no since")
    void nullDefaultAndSinceStayEmpty() {
        DefaultOption option = new DefaultOption(
            "psk-key", OptionType.STRING, "pre-shared key", true, null, null, true);

        assertEquals(Optional.empty(), option.getDefaultValue());
        assertEquals(Optional.empty(), option.getSince());
        assertTrue(option.isSecret());
    }
}
