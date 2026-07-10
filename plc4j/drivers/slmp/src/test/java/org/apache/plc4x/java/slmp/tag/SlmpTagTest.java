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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.slmp.SlmpDataType;
import org.apache.plc4x.java.slmp.readwrite.SlmpDeviceCode;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SlmpTagTest {

    @Test
    void parsesDecimalDeviceWithDefaults() {
        SlmpTag tag = SlmpTag.of("D350");
        assertEquals(SlmpDeviceCode.D, tag.getDeviceCode());
        assertEquals(350, tag.getDeviceNumber());
        assertEquals(SlmpDataType.WORD, tag.getDataType());   // default WORD
        assertEquals(1, tag.getQuantity());
        assertEquals(1, tag.getNumberOfPoints());
    }

    @Test
    void parsesDataTypeAndQuantity() {
        SlmpTag tag = SlmpTag.of("R200:REAL[4]");
        assertEquals(SlmpDeviceCode.R, tag.getDeviceCode());
        assertEquals(200, tag.getDeviceNumber());
        assertEquals(SlmpDataType.REAL, tag.getDataType());
        assertEquals(4, tag.getQuantity());
        assertEquals(8, tag.getNumberOfPoints());             // 4 * 2 words
    }

    @Test
    void parsesHexLinkRegisterBareForm() {
        SlmpTag tag = SlmpTag.of("W1A:WORD[10]");
        assertEquals(SlmpDeviceCode.W, tag.getDeviceCode());
        assertEquals(0x1A, tag.getDeviceNumber());            // W is hex
        assertEquals(10, tag.getQuantity());
    }

    @Test
    void parsesHexLinkRegisterExplicit0xForm() {
        assertEquals(0x1A, SlmpTag.of("W0x1A").getDeviceNumber());
        assertEquals(0x1A, SlmpTag.of("W0X1A").getDeviceNumber());
    }

    @Test
    void parses0xPrefixCombinedWithDatatypeAndQuantity() {
        // the documented example form: 0x prefix + datatype + quantity in one address
        SlmpTag tag = SlmpTag.of("W0x1A:WORD[10]");
        assertEquals(SlmpDeviceCode.W, tag.getDeviceCode());
        assertEquals(0x1A, tag.getDeviceNumber());
        assertEquals(SlmpDataType.WORD, tag.getDataType());
        assertEquals(10, tag.getQuantity());
    }

    @Test
    void canonicalAddressStringRoundTrips() {
        assertEquals("D350:INT[2]", SlmpTag.of("D350:INT[2]").getAddressString());
        assertEquals("W0x1A", SlmpTag.of("W1A").getAddressString());   // W prints 0x hex
    }

    @Test
    void rejectsUnsupportedBitDevice() {
        PlcInvalidTagException ex = assertThrows(PlcInvalidTagException.class, () -> SlmpTag.of("M100"));
        assertTrue(ex.getMessage().contains("M"));
    }

    @Test
    void rejects0xPrefixOnDecimalDevice() {
        assertThrows(PlcInvalidTagException.class, () -> SlmpTag.of("D0x10"));
    }

    @Test
    void rejectsHexDigitsInDecimalDevice() {
        assertThrows(PlcInvalidTagException.class, () -> SlmpTag.of("D1A"));
    }

    @Test
    void rejectsZeroQuantity() {
        assertThrows(PlcInvalidTagException.class, () -> SlmpTag.of("D100:WORD[0]"));
    }

    @Test
    void rejectsOverCeiling() {
        assertThrows(PlcInvalidTagException.class, () -> SlmpTag.of("D0:WORD[961]"));
    }
}
