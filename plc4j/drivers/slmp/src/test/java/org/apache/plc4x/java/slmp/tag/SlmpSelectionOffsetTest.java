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
package org.apache.plc4x.java.slmp.tag;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

/**
 * A selection offset counts elements; an SLMP device number counts 16-bit words. The read length
 * already scales by {@code getWordsPerElement()}, so an unscaled offset does not shorten the read -
 * it moves it to the wrong devices, silently.
 */
class SlmpSelectionOffsetTest {

    @Test
    @DisplayName("a one-word type advances one device per element")
    void oneWordPerElement() {
        assertEquals(104, SlmpTag.of("D100[4]:INT").getDeviceNumber());
    }

    @Test
    @DisplayName("a two-word type advances two devices per element")
    void twoWordsPerElement() {
        // The fifth DINT begins eight words along, at D108 - not D104.
        assertEquals(108, SlmpTag.of("D100[4]:DINT").getDeviceNumber());
    }

    @Test
    @DisplayName("a REAL advances by its two words as well")
    void realsAdvanceByTwoWords() {
        assertEquals(104, SlmpTag.of("D100[2]:REAL").getDeviceNumber());
    }

    @Test
    @DisplayName("a declared base is measured in elements too")
    void aDeclaredBaseIsInElements() {
        // [4..7;4] starts at the declared base, so it shifts nothing regardless of the width.
        assertEquals(100, SlmpTag.of("D100[4..7;4]:DINT").getDeviceNumber());
    }
}
