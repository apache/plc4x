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
package org.apache.plc4x.java.api.metadata;

import org.apache.plc4x.java.api.types.OptionType;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * {@link Option#isSecret()} is a default method so that implementations written before it existed
 * keep compiling and keep linking. That promise is what these tests hold: {@link LegacyOption}
 * implements only the six methods the interface had before, and the fact that this file compiles
 * at all is half the assertion.
 */
class OptionTest {

    /**
     * An implementation predating {@code isSecret()} - it overrides nothing beyond the originally
     * abstract methods. Adding a new abstract method to {@link Option} would break this class at
     * compile time, which is exactly the regression the default is there to prevent.
     */
    static class LegacyOption implements Option {
        @Override
        public String getKey() {
            return "timeout";
        }

        @Override
        public OptionType getType() {
            return OptionType.LONG;
        }

        @Override
        public String getDescription() {
            return "request timeout";
        }

        @Override
        public boolean isRequired() {
            return false;
        }

        @Override
        public Optional<Object> getDefaultValue() {
            return Optional.of(5000L);
        }

        @Override
        public Optional<String> getSince() {
            return Optional.of("1.0");
        }
    }

    /** An implementation that opts in, to show the default is overridable and not hard-wired. */
    static class SecretOption extends LegacyOption {
        @Override
        public boolean isSecret() {
            return true;
        }
    }

    @Test
    @DisplayName("an implementation that does not override reports not-secret")
    void defaultsToNotSecret() {
        // The safe default: an option nobody has classified must not be *assumed* sensitive,
        // because the machinery that consumes this (redaction, toString masking) would otherwise
        // hide every option of every driver that has not been migrated.
        assertFalse(new LegacyOption().isSecret());
    }

    @Test
    @DisplayName("an implementation can override the default")
    void overridingReportsSecret() {
        assertTrue(new SecretOption().isSecret());
    }

    @Test
    @DisplayName("the default does not disturb the other accessors")
    void leavesTheRestOfTheContractAlone() {
        Option option = new LegacyOption();
        assertFalse(option.isRequired());
        assertTrue(option.getDefaultValue().isPresent());
        assertTrue(option.getSince().isPresent());
    }
}
