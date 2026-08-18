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
package org.apache.plc4x.java.modbus.base;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcSINT;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A Modbus register is 16 bits wide, so a value narrower than that is padded when it is on its own
 * and packed when it is not. That used to be expressed as separate cases in the protocol
 * description - one per type per byte order plus a list variant - and now lives here, which is why
 * these rules need pinning down directly rather than only through the drivers.
 */
class ModbusRegisterCodecTest {

    /**
     * Big endian puts the padding in front of the value, so the byte carrying it comes second.
     */
    @Test
    void padsALoneByteAtTheFrontWhenBigEndian() throws Exception {
        assertEquals(0x2A, parse(new byte[]{0x00, 0x2A}, ModbusDataType.USINT, 1, true).getInt());
        assertArrayEquals(new byte[]{0x00, 0x2A}, serialize(new PlcUSINT((short) 0x2A), ModbusDataType.USINT, 1, true));
    }

    /**
     * Little endian puts the value first and the padding after it.
     */
    @Test
    void padsALoneByteAtTheEndWhenLittleEndian() throws Exception {
        assertEquals(0x2A, parse(new byte[]{0x2A, 0x00}, ModbusDataType.USINT, 1, false).getInt());
        assertArrayEquals(new byte[]{0x2A, 0x00}, serialize(new PlcUSINT((short) 0x2A), ModbusDataType.USINT, 1, false));
    }

    /** A signed byte is padded the same way, and keeps its sign. */
    @Test
    void padsALoneSignedByte() throws Exception {
        assertEquals(-1, parse(new byte[]{0x00, (byte) 0xFF}, ModbusDataType.SINT, 1, true).getInt());
        assertEquals(-1, parse(new byte[]{(byte) 0xFF, 0x00}, ModbusDataType.SINT, 1, false).getInt());
    }

    /**
     * A lone BOOL occupies a whole register: fifteen bits of padding and the bit itself, which
     * lands in the last bit big endian and in the eighth little endian.
     */
    @Test
    void padsALoneBoolToAWholeRegister() throws Exception {
        assertTrue(parse(new byte[]{0x00, 0x01}, ModbusDataType.BOOL, 1, true).getBoolean());
        assertEquals(false, parse(new byte[]{0x00, 0x00}, ModbusDataType.BOOL, 1, true).getBoolean());
        assertTrue(parse(new byte[]{0x01, 0x00}, ModbusDataType.BOOL, 1, false).getBoolean());
    }

    /**
     * CHAR is not padded - it never has been, unlike the other single-byte types.
     */
    @Test
    void leavesALoneCharUnpadded() throws Exception {
        assertEquals("A", parse(new byte[]{0x41}, ModbusDataType.CHAR, 1, true).getString());
        assertEquals(1, ModbusRegisterCodec.lengthInBytes(parse(new byte[]{0x41}, ModbusDataType.CHAR, 1, true),
            ModbusDataType.CHAR, 1, 1));
    }

    /** Values wide enough to fill a register are never padded, in either byte order. */
    @Test
    void neverPadsARegisterWideValue() throws Exception {
        assertEquals(0x1234, parse(new byte[]{0x12, 0x34}, ModbusDataType.UINT, 1, true).getInt());
        assertEquals(2, ModbusRegisterCodec.lengthInBytes(
            parse(new byte[]{0x12, 0x34}, ModbusDataType.UINT, 1, true), ModbusDataType.UINT, 1, 1));
    }

    /**
     * Several sub-register values are packed rather than padded - two bytes to a register - which
     * is why the layout of one value cannot simply be repeated.
     */
    @Test
    void packsSeveralBytesWithoutPadding() throws Exception {
        PlcValue value = parse(new byte[]{0x01, 0x02, 0x03, 0x04}, ModbusDataType.USINT, 4, true);

        assertEquals(4, value.getList().size());
        assertEquals(1, value.getList().get(0).getInt());
        assertEquals(4, value.getList().get(3).getInt());
    }

    /** Packed BOOLs are one bit each, not one register each. */
    @Test
    void packsBoolsAsBits() throws Exception {
        PlcValue value = parse(new byte[]{(byte) 0b10100000, 0x00}, ModbusDataType.BOOL, 3, true);

        assertEquals(3, value.getList().size());
        assertTrue(value.getList().get(0).getBoolean());
        assertEquals(false, value.getList().get(1).getBoolean());
        assertTrue(value.getList().get(2).getBoolean());
    }

    /**
     * An odd number of bytes leaves half a register, which has to be padded out when writing or
     * the request would be a byte short of a whole register.
     */
    @Test
    void padsTheLastRegisterForAnOddCount() throws Exception {
        PlcValue three = new PlcList(List.of(new PlcUSINT((short) 1), new PlcUSINT((short) 2), new PlcUSINT((short) 3)));

        assertEquals(4, ModbusRegisterCodec.lengthInBytes(three, ModbusDataType.USINT, 3, 1));
        assertArrayEquals(new byte[]{0x01, 0x02, 0x03, 0x00}, serialize(three, ModbusDataType.USINT, 3, true));
    }

    /** An even count fills its registers exactly and needs no trailing pad. */
    @Test
    void addsNoTrailingPadForAnEvenCount() throws Exception {
        PlcValue two = new PlcList(List.of(new PlcUSINT((short) 1), new PlcUSINT((short) 2)));

        assertEquals(2, ModbusRegisterCodec.lengthInBytes(two, ModbusDataType.USINT, 2, 1));
        assertArrayEquals(new byte[]{0x01, 0x02}, serialize(two, ModbusDataType.USINT, 2, true));
    }

    /** Sixteen packed bools are exactly one register. */
    @Test
    void padsPackedBoolsToTheRegisterBoundary() throws Exception {
        PlcValue threeBools = new PlcList(List.of(new PlcBOOL(true), new PlcBOOL(false), new PlcBOOL(true)));

        assertEquals(2, ModbusRegisterCodec.lengthInBytes(threeBools, ModbusDataType.BOOL, 3, 1));
        assertArrayEquals(new byte[]{(byte) 0b10100000, 0x00}, serialize(threeBools, ModbusDataType.BOOL, 3, true));
    }

    /** What is written for a padded value has to read back as the same value. */
    @Test
    void roundTripsAPaddedValue() throws Exception {
        for (boolean bigEndian : new boolean[]{true, false}) {
            byte[] written = serialize(new PlcSINT((byte) -42), ModbusDataType.SINT, 1, bigEndian);

            assertEquals(2, written.length, "a lone byte occupies a whole register");
            assertEquals(-42, parse(written, ModbusDataType.SINT, 1, bigEndian).getInt());
        }
    }

    /** And so does a packed run. */
    @Test
    void roundTripsAPackedRun() throws Exception {
        PlcValue values = new PlcList(List.of(new PlcSINT((byte) -1), new PlcSINT((byte) 2), new PlcSINT((byte) -3)));
        byte[] written = serialize(values, ModbusDataType.SINT, 3, true);

        PlcValue read = parse(written, ModbusDataType.SINT, 3, true);
        assertEquals(-1, read.getList().get(0).getInt());
        assertEquals(2, read.getList().get(1).getInt());
        assertEquals(-3, read.getList().get(2).getInt());
    }

    private static PlcValue parse(byte[] data, ModbusDataType dataType, int numberOfValues, boolean bigEndian)
        throws Exception {
        return ModbusRegisterCodec.parse(readBuffer(data), dataType, numberOfValues, bigEndian, 1);
    }

    private static byte[] serialize(PlcValue value, ModbusDataType dataType, int numberOfValues, boolean bigEndian)
        throws Exception {
        int size = ModbusRegisterCodec.lengthInBytes(value, dataType, numberOfValues, 1);
        WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[size],
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        ModbusRegisterCodec.serialize(writeBuffer, value, dataType, numberOfValues, bigEndian, 1);
        return writeBuffer.getBytes();
    }

    private static ReadBufferByteBased readBuffer(byte[] data) {
        return new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
    }
}
