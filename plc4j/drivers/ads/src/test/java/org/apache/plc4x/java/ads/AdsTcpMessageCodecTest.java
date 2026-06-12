/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.ads;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.lang.reflect.Method;

import static org.junit.jupiter.api.Assertions.assertEquals;

class AdsTcpMessageCodecTest {

    @Test
    void minimumHeaderSizeIsAmsTcpHeader() throws Exception {
        AdsTcpMessageCodec codec = new AdsTcpMessageCodec(Mockito.mock(TransportInstance.class), p -> { });
        Method m = AdsTcpMessageCodec.class.getDeclaredMethod("getMinimumHeaderSize");
        m.setAccessible(true);
        assertEquals(6, m.invoke(codec));
    }

    @Test
    void calculateTotalMessageSizeUsesLittleEndianLength() throws Exception {
        AdsTcpMessageCodec codec = new AdsTcpMessageCodec(Mockito.mock(TransportInstance.class), p -> { });
        Method m = AdsTcpMessageCodec.class.getDeclaredMethod("calculateTotalMessageSize", byte[].class, int.class);
        m.setAccessible(true);
        // 6-byte header: 2 reserved + 4-byte little-endian length = 0x00000010 (16)
        byte[] header = new byte[]{0, 0, 0x10, 0, 0, 0};
        assertEquals(6 + 16, m.invoke(codec, header, 1024));

        byte[] header2 = new byte[]{0, 0, (byte) 0xFF, (byte) 0xFF, 0, 0};
        assertEquals(6 + 0xFFFF, m.invoke(codec, header2, 1024));
    }
}
