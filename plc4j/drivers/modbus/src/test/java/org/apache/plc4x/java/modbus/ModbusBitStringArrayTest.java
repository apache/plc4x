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
package org.apache.plc4x.java.modbus;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.base.ModbusRegisterCodec;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.values.PlcDWORD;
import org.apache.plc4x.java.spi.values.PlcLWORD;
import org.apache.plc4x.java.spi.values.PlcWORD;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * WORD, DWORD and LWORD had no list variant in the type switch, so a request for several of them
 * matched the single-value case and returned only the first register - the reporter had to fall
 * back to INT to read more than one (GH-2061).
 * <p>
 * A list case in the mspec would not have been enough: it generates an array over the underlying
 * field, so the elements would have come back as PlcUINT rather than PlcWORD. Every element now
 * goes through the single-value parser, so the assertions below check the element type as well as
 * the value.
 */
class ModbusBitStringArrayTest {

    @Test
    void readsSeveralWords() throws Exception {
        PlcValue value = parse(new byte[]{0x00, 0x2A, (byte) 0xFF, (byte) 0xFF, 0x12, 0x34},
            ModbusDataType.WORD, 3);

        assertEquals(3, value.getList().size(), "used to return a single value");
        assertEquals(42, value.getList().get(0).getInt());
        assertEquals(65535, value.getList().get(1).getInt());
        assertEquals(0x1234, value.getList().get(2).getInt());
        for (PlcValue element : value.getList()) {
            assertInstanceOf(PlcWORD.class, element, "a list case would have yielded PlcUINT");
        }
    }

    @Test
    void readsSeveralDoubleWords() throws Exception {
        PlcValue value = parse(new byte[]{0x00, 0x00, 0x00, 0x2A, (byte) 0xFF, (byte) 0xFF, (byte) 0xFF, (byte) 0xFF},
            ModbusDataType.DWORD, 2);

        assertEquals(2, value.getList().size());
        assertEquals(42L, value.getList().get(0).getLong());
        assertEquals(4294967295L, value.getList().get(1).getLong());
        for (PlcValue element : value.getList()) {
            assertInstanceOf(PlcDWORD.class, element, "a list case would have yielded PlcUDINT");
        }
    }

    @Test
    void readsSeveralLongWords() throws Exception {
        PlcValue value = parse(new byte[]{
            0, 0, 0, 0, 0, 0, 0, 42,
            0, 0, 0, 0, 0, 0, 1, 0}, ModbusDataType.LWORD, 2);

        assertEquals(2, value.getList().size());
        assertEquals(42L, value.getList().get(0).getLong());
        assertEquals(256L, value.getList().get(1).getLong());
        for (PlcValue element : value.getList()) {
            assertInstanceOf(PlcLWORD.class, element, "a list case would have yielded PlcULINT");
        }
    }

    /** A single value keeps returning a value rather than a one-element list. */
    @Test
    void stillReadsASingleWordAsAScalar() throws Exception {
        PlcValue value = parse(new byte[]{0x00, 0x2A}, ModbusDataType.WORD, 1);

        assertInstanceOf(PlcWORD.class, value);
        assertTrue(value.isInteger());
        assertEquals(42, value.getInt());
    }

    /**
     * The register arithmetic was always right - only the parsing stopped after the first value -
     * so a request did ask the device for the full range.
     */
    @Test
    void asksTheDeviceForEveryRegister() {
        assertEquals(3, ModbusTagHoldingRegister.of("holding-register:1[0..2]:WORD").getLengthWords());
        assertEquals(4, ModbusTagHoldingRegister.of("holding-register:1[0..1]:DWORD").getLengthWords());
        assertEquals(8, ModbusTagHoldingRegister.of("holding-register:1[0..1]:LWORD").getLengthWords());
    }

    private static PlcValue parse(byte[] data, ModbusDataType dataType, int numberOfValues) throws Exception {
        ReadBufferByteBased buffer = new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        return ModbusRegisterCodec.parse(buffer, dataType, numberOfValues, true, 1);
    }
}
