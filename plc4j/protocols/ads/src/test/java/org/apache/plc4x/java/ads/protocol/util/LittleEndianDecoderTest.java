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

import org.apache.plc4x.java.ads.api.commands.types.Length;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.junit.Test;

import java.util.Calendar;
import java.util.Date;

import static java.util.Arrays.asList;
import static java.util.Collections.singletonList;
import static org.apache.plc4x.java.base.util.Junit5Backport.assertThrows;
import static org.junit.Assert.assertEquals;

public class LittleEndianDecoderTest {

    @Test
    public void getLengthFor() {
        assertEquals(LittleEndianDecoder.getLengthFor(Boolean.class, 0), Length.of(1));
        assertEquals(LittleEndianDecoder.getLengthFor(Byte.class, 0), Length.of(1));
        assertEquals(LittleEndianDecoder.getLengthFor(Short.class, 0), Length.of(2));
        assertEquals(LittleEndianDecoder.getLengthFor(Integer.class, 0), Length.of(4));
        assertEquals(LittleEndianDecoder.getLengthFor(Float.class, 0), Length.of(4));
        assertEquals(LittleEndianDecoder.getLengthFor(Double.class, 0), Length.of(8));
        assertEquals(LittleEndianDecoder.getLengthFor(Calendar.class, 0), Length.of(8));
        assertEquals(LittleEndianDecoder.getLengthFor(LittleEndianDecoderTest.class, 666), Length.of(666));
    }

    @Test
    public void decodeData() throws Exception {
        assertEquals(asList(true, false), LittleEndianDecoder.decodeData(Boolean.class, new byte[]{0x1, 0x0}));

        assertEquals(asList((byte) 0x1, (byte) 0x0), LittleEndianDecoder.decodeData(Byte.class, new byte[]{0x1, 0x0}));

        assertEquals(singletonList((short) 1), LittleEndianDecoder.decodeData(Short.class, new byte[]{0x1}));
        assertEquals(singletonList((short) 256), LittleEndianDecoder.decodeData(Short.class, new byte[]{0x0, 0x1}));
        assertEquals(asList((short) 256, (short) 256), LittleEndianDecoder.decodeData(Short.class, new byte[]{0x0, 0x1, 0x0, 0x1}));

        assertEquals(singletonList(1), LittleEndianDecoder.decodeData(Integer.class, new byte[]{0x1}));
        assertEquals(singletonList(16777216), LittleEndianDecoder.decodeData(Integer.class, new byte[]{0x0, 0x0, 0x0, 0x1}));
        assertEquals(asList(16777216, 16777216), LittleEndianDecoder.decodeData(Integer.class, new byte[]{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1}));

        assertEquals(singletonList(1.4E-45f), LittleEndianDecoder.decodeData(Float.class, new byte[]{0x1}));
        assertEquals(singletonList(2.3509887E-38f), LittleEndianDecoder.decodeData(Float.class, new byte[]{0x0, 0x0, 0x0, 0x1}));
        assertEquals(asList(2.3509887E-38f, 2.3509887E-38f), LittleEndianDecoder.decodeData(Float.class, new byte[]{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1}));

        assertEquals(singletonList(4.9E-324), LittleEndianDecoder.decodeData(Double.class, new byte[]{0x1}));
        assertEquals(singletonList(7.2911220195563975E-304), LittleEndianDecoder.decodeData(Double.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}));
        assertEquals(asList(7.2911220195563975E-304, 7.2911220195563975E-304), LittleEndianDecoder.decodeData(Double.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}));

        Calendar calendar1 = Calendar.getInstance();
        calendar1.setTime(new Date(-11644473600000L));
        assertEquals(singletonList(calendar1), LittleEndianDecoder.decodeData(Calendar.class, new byte[]{0x1}));
        Calendar calendar0x0001 = Calendar.getInstance();
        calendar0x0001.setTime(new Date(-4438714196208L));
        assertEquals(singletonList(calendar0x0001), LittleEndianDecoder.decodeData(Calendar.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}));
        assertEquals(asList(calendar0x0001, calendar0x0001), LittleEndianDecoder.decodeData(Calendar.class, new byte[]{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}));

        assertEquals(singletonList("plc4x"), LittleEndianDecoder.decodeData(String.class, new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x0}));
        assertEquals(singletonList("plc4xplc4x"), LittleEndianDecoder.decodeData(String.class, new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x70, 0x6c, 0x63, 0x34, 0x78, 0x0}));
        assertEquals(asList("plc4x", "plc4x"), LittleEndianDecoder.decodeData(String.class, new byte[]{0x70, 0x6c, 0x63, 0x34, 0x78, 0x0, 0x70, 0x6c, 0x63, 0x34, 0x78, 0x0}));

        assertThrows(PlcProtocolException.class, () -> LittleEndianDecoder.decodeData(String.class, new byte[]{0x01}));
        assertThrows(PlcUnsupportedDataTypeException.class, () -> LittleEndianDecoder.decodeData(this.getClass(), new byte[10]));
    }

}