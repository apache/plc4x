/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at
 
   http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.api.commands.types;

import org.apache.commons.codec.binary.Hex;
import org.junit.Ignore;
import org.junit.Test;

import static org.apache.plc4x.java.ads.util.Junit5Backport.assertThrows;
import static org.junit.Assert.assertEquals;


public class ReadLengthTest {

    private final byte NULL_BYTE = 0x0;

    @Test
    public void ofBytes() {
        assertEquals(0l, ReadLength.of(NULL_BYTE, NULL_BYTE, NULL_BYTE, NULL_BYTE).getAsLong());
        assertThrows(IllegalArgumentException.class, () -> ReadLength.of(NULL_BYTE, NULL_BYTE, NULL_BYTE, NULL_BYTE, NULL_BYTE));
    }

    @Test
    public void ofLong() {
        assertByte(ReadLength.of(1), "0x01000000");
        assertByte(ReadLength.of(65535), "0xffff0000");
        assertThrows(IllegalArgumentException.class, () -> ReadLength.of(-1));
        assertThrows(IllegalArgumentException.class, () -> ReadLength.of(4294967296L));
    }

    @Test
    public void ofString() {
        assertByte(ReadLength.of("1"), "0x01000000");
    }

    @Test
    @Ignore("Must be fixed after refactoring of 'toString' method.")
    public void testToString() {
        assertEquals(ReadLength.of("1").getAsLong(), 1l);
    }

    private void assertByte(ReadLength actual, String expected) {
        assertEquals(expected, "0x" + Hex.encodeHexString(actual.getBytes()));
    }
}