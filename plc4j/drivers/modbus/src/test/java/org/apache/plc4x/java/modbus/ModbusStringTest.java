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

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.modbus.base.tag.ModbusTagHoldingRegister;
import org.apache.plc4x.java.modbus.readwrite.DataItem;
import org.apache.plc4x.java.modbus.readwrite.ModbusDataType;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Reading a string from holding registers. The data type existed but the type switch had no case
 * for it, so parsing returned null with no exception - a read that looked like it had succeeded and
 * found nothing (GH-2307).
 */
class ModbusStringTest {

    @Test
    void readsASingleString() throws Exception {
        PlcValue value = parse("Hello Toddy!", ModbusDataType.STRING, 1, 12);

        assertNotNull(value, "a STRING used to parse to null");
        assertEquals("Hello Toddy!", value.getString());
    }

    /**
     * numberOfValues is the number of strings, each of stringLength characters - so the registers
     * are read as three consecutive 4-character strings rather than one 12-character one.
     */
    @Test
    void readsAnArrayOfStrings() throws Exception {
        PlcValue value = parse("abcdefghijkl", ModbusDataType.STRING, 3, 4);

        assertEquals(3, value.getList().size());
        assertEquals("abcd", value.getList().get(0).getString());
        assertEquals("efgh", value.getList().get(1).getString());
        assertEquals("ijkl", value.getList().get(2).getString());
    }

    @Test
    void readsAWideString() throws Exception {
        byte[] utf16 = "Hi".getBytes(StandardCharsets.UTF_16BE);
        PlcValue value = DataItem.staticParse(buffer(utf16), ModbusDataType.WSTRING, 1, true, 2);

        assertEquals("Hi", value.getString());
    }

    /**
     * A string occupies stringLength bytes, two per register, so the tag has to ask for the right
     * number of registers - 12 characters are 6 registers, and three of them are 18.
     */
    @Test
    void sizesTheRequestFromTheDeclaredLength() {
        assertEquals(6, ModbusTagHoldingRegister.of("holding-register:1:STRING(12)").getLengthWords());
        assertEquals(18, ModbusTagHoldingRegister.of("holding-register:1:STRING(12)[3]").getLengthWords());
        // A wide character is two bytes, so the same declared length needs twice the registers.
        assertEquals(12, ModbusTagHoldingRegister.of("holding-register:1:WSTRING(12)").getLengthWords());
    }

    @Test
    void keepsTheDeclaredLengthOnTheTag() {
        assertEquals(20, ModbusTagHoldingRegister.of("holding-register:1:STRING(20)").getStringLength());
        assertEquals(3, ModbusTagHoldingRegister.of("holding-register:1:STRING(20)[3]").getNumberOfElements());
    }

    /**
     * Nothing on the wire says how long a string is, so the address has to. Guessing a default
     * would bring back the failure this fixes, only quieter.
     */
    @Test
    void refusesAStringWithoutALength() {
        PlcInvalidTagException exception = assertThrows(PlcInvalidTagException.class,
            () -> ModbusTagHoldingRegister.of("holding-register:8:STRING"));
        assertEquals(true, exception.getMessage().contains("STRING(20)"), exception.getMessage());
    }

    /**
     * The length only means something for strings; accepting it elsewhere would silently change
     * how many registers are requested.
     */
    @Test
    void refusesALengthOnANonStringType() {
        assertThrows(PlcInvalidTagException.class,
            () -> ModbusTagHoldingRegister.of("holding-register:8:INT(20)"));
    }

    /** Every other data type keeps sizing exactly as before. */
    @Test
    void leavesOtherDataTypesUnchanged() {
        assertEquals(1, ModbusTagHoldingRegister.of("holding-register:1:INT").getStringLength());
        assertEquals(7, ModbusTagHoldingRegister.of("holding-register:1:INT[7]").getLengthWords());
        assertEquals(14, ModbusTagHoldingRegister.of("holding-register:1:DINT[7]").getLengthWords());
    }

    private static PlcValue parse(String text, ModbusDataType dataType, int numberOfValues, int stringLength)
        throws Exception {
        return DataItem.staticParse(buffer(text.getBytes(StandardCharsets.UTF_8)),
            dataType, numberOfValues, true, stringLength);
    }

    private static ReadBufferByteBased buffer(byte[] data) {
        return new ReadBufferByteBased(data,
            WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
            WithOption.WithSignedIntegerEncoding("twos-complement"),
            WithOption.WithFloatEncoding("IEEE754"),
            WithOption.WithStringEncoding("UTF8"),
            WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
    }
}
