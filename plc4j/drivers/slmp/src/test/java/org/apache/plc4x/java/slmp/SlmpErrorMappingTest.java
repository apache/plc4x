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
package org.apache.plc4x.java.slmp;

import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.slmp.tag.SlmpTag;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SlmpErrorMappingTest {

    private static byte[] hex(String s) {
        byte[] out = new byte[s.length() / 2];
        for (int i = 0; i < out.length; i++) {
            out[i] = (byte) Integer.parseInt(s.substring(i * 2, i * 2 + 2), 16);
        }
        return out;
    }

    @Test
    void normalCompletionDecodesValue() {
        SlmpTag tag = SlmpTag.of("D350:WORD[2]");
        var item = SlmpResponseMapper.mapTag(tag, 0x0000, hex("ab560f17"));
        assertEquals(PlcResponseCode.OK, item.getResponseCode());
        assertEquals(0x56AB, item.getValue().getList().get(0).getInt());
    }

    @Test
    void nonZeroEndCodeMapsToRemoteError() {
        SlmpTag tag = SlmpTag.of("D350:WORD");
        var item = SlmpResponseMapper.mapTag(tag, 0xC059, new byte[0]);
        assertEquals(PlcResponseCode.REMOTE_ERROR, item.getResponseCode());
        assertNull(item.getValue());
    }

    @Test
    void shortResponseMapsToInvalidData() {
        SlmpTag tag = SlmpTag.of("D350:WORD[2]");      // wants 2 words
        var item = SlmpResponseMapper.mapTag(tag, 0x0000, hex("ab56")); // only 1
        assertEquals(PlcResponseCode.INVALID_DATA, item.getResponseCode());
    }

    @Test
    void mapWriteTagOkOnZeroEndCode() {
        assertEquals(PlcResponseCode.OK, SlmpResponseMapper.mapWriteTag(SlmpTag.of("D350"), 0x0000));
    }

    @Test
    void mapWriteTagRemoteErrorOnNonZeroEndCode() {
        assertEquals(PlcResponseCode.REMOTE_ERROR, SlmpResponseMapper.mapWriteTag(SlmpTag.of("D350"), 0xC059));
    }
}
