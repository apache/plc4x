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
package org.apache.plc4x.java.ads.api.generic.types;

import org.apache.commons.codec.binary.Hex;
import org.junit.Test;

import static org.apache.plc4x.java.ads.util.Junit5Backport.assertThrows;
import static org.junit.Assert.assertEquals;


public class AMSPortTest {

    byte NULL_BYTE = 0x0;

    @Test
    public void ofBytes() {
        assertEquals("0", AMSPort.of(NULL_BYTE, NULL_BYTE).toString());
        assertThrows(IllegalArgumentException.class, () -> AMSPort.of(NULL_BYTE, NULL_BYTE, NULL_BYTE));
    }

    @Test
    public void ofInt() {
        assertByte(AMSPort.of(1), "0x0100");
        assertByte(AMSPort.of(65535), "0xffff");
        assertThrows(IllegalArgumentException.class, () -> AMSPort.of(-1));
        assertThrows(IllegalArgumentException.class, () -> AMSPort.of(65536));
    }

    @Test
    public void ofString() {
        assertByte(AMSPort.of("1"), "0x0100");
    }

    @Test
    public void testToString() {
        assertEquals(AMSPort.of("1").toString(), "1");
    }

    void assertByte(AMSPort actual, String expected) {
        assertEquals(expected, "0x" + Hex.encodeHexString(actual.getBytes()));
    }
}