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
package org.apache.plc4x.java.firmata.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.firmata.readwrite.PinMode;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertNotEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FirmataTagTest {

    @Test
    void factoryRoutesDigitalAndAnalogPrefixes() {
        assertInstanceOf(FirmataTagDigital.class, FirmataTag.of("digital:5"));
        assertInstanceOf(FirmataTagAnalog.class, FirmataTag.of("analog:2"));
    }

    @Test
    void factoryRejectsUnparseableTagStrings() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTag.of("uart:0"));
    }

    @Test
    void quantityMustBePositive() {
        assertThrows(IllegalArgumentException.class,
            () -> new FirmataTagDigital(0, 0, null));
    }

    @Test
    void digitalDefaults() {
        FirmataTagDigital tag = FirmataTagDigital.of("digital:2");
        assertEquals(2, tag.getAddress());
        assertEquals(1, tag.getNumberOfElements());
        assertEquals(PlcValueType.BOOL, tag.getPlcValueType());
        assertEquals("digital:2", tag.getAddressString());
        // Single-element tags don't emit ArrayInfo.
        assertTrue(tag.getArrayInfo().isEmpty());
        // No explicit mode means INPUT (set by the connection on subscribe).
        assertNull(tag.getPinMode());
    }

    @Test
    void digitalWithRangeAndPullup() {
        FirmataTagDigital tag = FirmataTagDigital.of("digital:8[0..3]:PULLUP");
        assertEquals(8, tag.getAddress());
        assertEquals(4, tag.getNumberOfElements());
        assertEquals(PinMode.PinModePullup, tag.getPinMode());
        assertEquals("digital:8[0..3]", tag.getAddressString());
        assertFalse(tag.getArrayInfo().isEmpty());
        // BitSet should mark pins 8..11.
        for (int pin = 8; pin < 12; pin++) {
            assertTrue(tag.getBitSet().get(pin), "pin " + pin + " should be set");
        }
        assertFalse(tag.getBitSet().get(7));
        assertFalse(tag.getBitSet().get(12));
    }

    @Test
    void digitalRejectsMalformedAddress() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTagDigital.of("digital:foo"));
    }

    @Test
    void analogDefaults() {
        FirmataTagAnalog tag = FirmataTagAnalog.of("analog:3");
        assertEquals(3, tag.getAddress());
        assertEquals(1, tag.getNumberOfElements());
        assertEquals(PlcValueType.INT, tag.getPlcValueType());
        assertEquals("analog:3", tag.getAddressString());
        assertTrue(tag.getArrayInfo().isEmpty());
    }

    @Test
    void analogWithRange() {
        FirmataTagAnalog tag = FirmataTagAnalog.of("analog:0[0..2]");
        assertEquals(0, tag.getAddress());
        assertEquals(3, tag.getNumberOfElements());
        assertEquals("analog:0[0..2]", tag.getAddressString());
        assertFalse(tag.getArrayInfo().isEmpty());
    }

    @Test
    void analogRejectsMalformedAddress() {
        assertThrows(PlcInvalidTagException.class, () -> FirmataTagAnalog.of("analog:x"));
    }

    @Test
    void equalityAndHashAreAddressBased() {
        FirmataTagDigital a = FirmataTagDigital.of("digital:5");
        FirmataTagDigital b = FirmataTagDigital.of("digital:5[0..2]");
        FirmataTagDigital c = FirmataTagDigital.of("digital:6");
        // FirmataTag.equals compares by address only, so a and b match.
        assertEquals(a, b);
        assertEquals(a.hashCode(), b.hashCode());
        assertNotEquals(a, c);
        assertNotEquals(a, "not a tag");
        // toString contains the address — sanity check it's there.
        assertTrue(a.toString().contains("address=5"));
    }

}
