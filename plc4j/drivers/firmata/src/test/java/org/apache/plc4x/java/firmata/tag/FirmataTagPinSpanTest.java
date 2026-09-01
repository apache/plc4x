/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.firmata.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

/**
 * A digital tag turns its element count into a set bit per pin, so the count decides both how
 * long the loop runs and how much the BitSet holds. The wire names a pin in eight bits, which
 * is the only span worth accepting.
 */
public class FirmataTagPinSpanTest {

    @Test
    void aCountBeyondThePinSpaceIsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTag.of("digital:0[100000000]"));
    }

    @Test
    void aCountTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTag.of("digital:0[99999999999]"));
    }

    @Test
    void anAddressTooWideToBeANumberIsAnInvalidTagNotANumberFormatError() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTag.of("analog:99999999999"));
    }

    @Test
    void aSpanRunningPastTheLastPinIsRejected() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTag.of("digital:250[0..99]"));
    }

    @Test
    void theWholePinSpaceIsStillAllowed() {
        FirmataTagDigital tag = FirmataTagDigital.of("digital:0[0..255]");
        assertEquals(256, tag.getNumberOfElements());
        assertEquals(256, tag.getBitSet().cardinality());
    }
}
