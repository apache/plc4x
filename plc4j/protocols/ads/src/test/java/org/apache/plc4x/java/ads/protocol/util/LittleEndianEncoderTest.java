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
package org.apache.plc4x.java.ads.protocol.util;

import org.junit.Test;

import java.util.Calendar;
import java.util.Date;

import static org.apache.plc4x.java.ads.util.Assert.assertByteEquals;

public class LittleEndianEncoderTest {

    @Test
    public void encodeData() throws Exception {
        assertByteEquals(new byte[]{0x01, 0x00, 0x01, 0x00}, LittleEndianEncoder.encodeData(Boolean.class, true, false, true, false));

        assertByteEquals(new byte[]{0x12, 0x03, 0x05, 0x7f}, LittleEndianEncoder.encodeData(Byte.class, (byte) 0x12, (byte) 0x03, (byte) 0x05, (byte) 0x7f));

        assertByteEquals(new byte[]{0x1, 0x00}, LittleEndianEncoder.encodeData(Short.class, (short) 1));
        assertByteEquals(new byte[]{0x0e, 0x00, 0x50, 0x00}, LittleEndianEncoder.encodeData(Short.class, (short) 14, (short) 80));

        assertByteEquals(new byte[]{0x5a, 0x0a, 0x00, 0x00}, LittleEndianEncoder.encodeData(Integer.class, 2650));
        assertByteEquals(new byte[]{0x5a, 0x0a, 0x00, 0x00, 0x50, 0x00, 0x00, 0x00}, LittleEndianEncoder.encodeData(Integer.class, 2650, 80));

        assertByteEquals(new byte[]{(byte) 0xc3, (byte) 0xf5, 0x48, 0x40}, LittleEndianEncoder.encodeData(Float.class, 3.14f));
        assertByteEquals(new byte[]{(byte) 0xc3, (byte) 0xf5, 0x48, 0x40, 0x14, (byte) 0xae, 0x07, 0x40}, LittleEndianEncoder.encodeData(Float.class, 3.14f, 2.12f));

        assertByteEquals(new byte[]{0x1F, (byte) 0x85, (byte) 0xEB, 0x51, (byte) 0xB8, 0x1E, 0x09, 0x40}, LittleEndianEncoder.encodeData(Double.class, 3.14));
        assertByteEquals(new byte[]{0x1F, (byte) 0x85, (byte) 0xEB, 0x51, (byte) 0xB8, 0x1E, 0x09, 0x40, (byte) 0xF6, 0x28, 0x5C, (byte) 0x8F, (byte) 0xC2, (byte) 0xF5, 0x00, 0x40}, LittleEndianEncoder.encodeData(Double.class, 3.14, 2.12));

        Calendar calendar1 = Calendar.getInstance();
        //calendar1.set(2003, Calendar.DECEMBER, 23, 13, 3, 0);
        calendar1.setTime(new Date(1072180980436L));
        assertByteEquals(new byte[]{(byte) 0x40, (byte) 0x79, (byte) 0xFB, (byte) 0xB5, (byte) 0x4C, (byte) 0xC9, (byte) 0xC3, (byte) 0x01}, LittleEndianEncoder.encodeData(Calendar.class, calendar1));

        assertByteEquals(new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x00}, LittleEndianEncoder.encodeData(String.class, "plc4x"));
        assertByteEquals(new byte[]{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}, LittleEndianEncoder.encodeData(String.class, "HelloWorld!"));
        assertByteEquals(new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x00, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x00}, LittleEndianEncoder.encodeData(String.class, "plc4x", "HelloWorld!"));
    }
}