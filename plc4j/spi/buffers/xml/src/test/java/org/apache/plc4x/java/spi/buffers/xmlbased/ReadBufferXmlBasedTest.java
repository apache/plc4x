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
package org.apache.plc4x.java.spi.buffers.xmlbased;

import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.io.ByteArrayInputStream;
import java.math.BigDecimal;
import java.math.BigInteger;

public class ReadBufferXmlBasedTest {

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readBit

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadBit() throws Exception {
        String input = "<bitField dataType=\"bit\" bitLength=\"1\">true</bitField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        boolean bit = buffer.readBit(WithOption.WithName("bitField"));
        Assertions.assertTrue(bit);

        input = "<bitField dataType=\"bit\" bitLength=\"1\">false</bitField>";
        buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        bit = buffer.readBit(WithOption.WithName("bitField"));
        Assertions.assertFalse(bit);
    }

    @Test
    public void testReadBitInvalidType() {
        String input = "<bitField dataType=\"uint\" bitLength=\"1\">true</bitField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("bitField")));
    }

    @Test
    public void testReadBitInvalidValue() {
        String input = "<bitField dataType=\"bit\" bitLength=\"1\">42</bitField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("bitField")));
    }

    @Test
    public void testReadBitInvalidBitLength() {
        String input = "<bitField dataType=\"bit\" bitLength=\"42\">true</bitField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("bitField")));
    }

    @Test
    public void testReadBitInvalidName() {
        String input = "<lalaField dataType=\"bit\" bitLength=\"1\">true</lalaField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("bitField")));
    }

    @Test
    public void testReadBitInvalidClosingTag() {
        String input = "<bitField dataType=\"bit\" bitLength=\"1\">true</wrongField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("bitField")));
    }

    @Test
    public void testReadBitInvalidTagStructure() {
        String input = "<bitField dataType=\"bit\" bitLength=\"1\">true<lalala/></bitField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("bitField")));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readBits

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadBits() throws Exception {
        String input = "<bitsField dataType=\"byte\" bitLength=\"32\">0x4875727A</bitsField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        byte[] bits = buffer.readBits(32, WithOption.WithName("bitsField"));
        byte[] expected = "Hurz".getBytes();
        Assertions.assertArrayEquals(expected, bits);
    }

    @Test
    public void testReadBitsInvalidType() {
        String input = "<bitsField dataType=\"uint\" bitLength=\"32\">0x4875727A</bitsField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBits(32, WithOption.WithName("bitsField")));
    }

    @Test
    public void testReadBitsInvalidValue() {
        String input = "<bitsField dataType=\"byte\" bitLength=\"32\">lalala</bitsField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBits(32, WithOption.WithName("bitsField")));
    }

    @Test
    public void testReadBitsInvalidBitLength() {
        String input = "<bitsField dataType=\"byte\" bitLength=\"12\">0x4875727A</bitsField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        // Simple Field (startingAddress)
        Assertions.assertThrows(BufferException.class, () -> buffer.readBits(32, WithOption.WithName("bitsField")));
    }

    @Test
    public void testReadBitsTooShortHexString() {
        String input = "<bitsField dataType=\"byte\" bitLength=\"32\">0x01</bitsField>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readBits(32, WithOption.WithName("bitsField")));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readUnsignedByte

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadUnsignedByte() throws Exception {
        String input = "<uByte dataType=\"uint\" bitLength=\"7\">127</uByte>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        short value = buffer.readUnsignedByte(7, WithOption.WithName("uByte"));
        Assertions.assertEquals(127, value);
    }

    @Test
    public void testReadUnsignedByteInvalids() {
        // Wrong type
        String wrongType = "<uByte dataType=\"int\" bitLength=\"8\">255</uByte>";
        ReadBufferXmlBased buf1 = new ReadBufferXmlBased(new ByteArrayInputStream(wrongType.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buf1.readUnsignedByte(8, WithOption.WithName("uByte")));
        // Wrong value
        String wrongVal = "<uByte dataType=\"uint\" bitLength=\"8\">NaN</uByte>";
        ReadBufferXmlBased buf2 = new ReadBufferXmlBased(new ByteArrayInputStream(wrongVal.getBytes()));
        Assertions.assertThrows(Exception.class, () -> buf2.readUnsignedByte(8, WithOption.WithName("uByte")));
        // Wrong bit length
        String wrongBits = "<uByte dataType=\"uint\" bitLength=\"7\">1</uByte>";
        ReadBufferXmlBased buf3 = new ReadBufferXmlBased(new ByteArrayInputStream(wrongBits.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buf3.readUnsignedByte(8, WithOption.WithName("uByte")));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readUnsignedShort

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadUnsignedShort() throws Exception {
        String input = "<uShort dataType=\"uint\" bitLength=\"15\">32767</uShort>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        int value = buffer.readUnsignedShort(15, WithOption.WithName("uShort"));
        Assertions.assertEquals(32767, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readUnsignedInt

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadUnsignedInt() throws Exception {
        String input = "<uInt dataType=\"uint\" bitLength=\"31\">123456789</uInt>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        long value = buffer.readUnsignedInt(31, WithOption.WithName("uInt"));
        Assertions.assertEquals(123456789L, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readUnsignedLong

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadUnsignedLong() throws Exception {
        String input = "<uLong dataType=\"uint\" bitLength=\"63\">8446744073709551615</uLong>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        long value = buffer.readUnsignedLong(63, WithOption.WithName("uLong"));
        Assertions.assertEquals(8446744073709551615L, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readUnsignedBigInteger

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadUnsignedBigInteger() throws Exception {
        String input = "<uBig dataType=\"uint\" bitLength=\"127\">340282366920938463463374607431768211455</uBig>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        BigInteger value = buffer.readUnsignedBigInteger(127, WithOption.WithName("uBig"));
        Assertions.assertEquals(new BigInteger("340282366920938463463374607431768211455"), value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readSignedByte

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadSignedByte() throws Exception {
        String input = "<sByte dataType=\"int\" bitLength=\"8\">-42</sByte>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        byte value = buffer.readSignedByte(8, WithOption.WithName("sByte"));
        Assertions.assertEquals((byte) -42, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readSignedShort

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadSignedShort() throws Exception {
        String input = "<sShort dataType=\"int\" bitLength=\"16\">-12345</sShort>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        short value = buffer.readSignedShort(16, WithOption.WithName("sShort"));
        Assertions.assertEquals((short) -12345, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readSignedInt

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadSignedInt() throws Exception {
        String input = "<sInt dataType=\"int\" bitLength=\"32\">-123456</sInt>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        int value = buffer.readSignedInt(32, WithOption.WithName("sInt"));
        Assertions.assertEquals(-123456, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readSignedLong

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadSignedLong() throws Exception {
        String input = "<sLong dataType=\"int\" bitLength=\"64\">-1234567890123456789</sLong>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        long value = buffer.readSignedLong(64, WithOption.WithName("sLong"));
        Assertions.assertEquals(-1234567890123456789L, value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readSignedBigInteger

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadSignedBigInteger() throws Exception {
        String input = "<sBig dataType=\"int\" bitLength=\"128\">-123456789012345678901234567890</sBig>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        BigInteger value = buffer.readSignedBigInteger(128, WithOption.WithName("sBig"));
        Assertions.assertEquals(new BigInteger("-123456789012345678901234567890"), value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readFloat

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadFloat() throws Exception {
        String input = "<flt dataType=\"float\" bitLength=\"32\">3.14159</flt>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        float value = buffer.readFloat(32, WithOption.WithName("flt"));
        Assertions.assertEquals(3.14159f, value, 0.00001f);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readDouble

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadDouble() throws Exception {
        String input = "<dbl dataType=\"float\" bitLength=\"64\">2.718281828459045</dbl>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        double value = buffer.readDouble(64, WithOption.WithName("dbl"));
        Assertions.assertEquals(2.718281828459045, value, 0.000000000000001);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readBigDecimal

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadBigDecimal() throws Exception {
        String input = "<bd dataType=\"float\" bitLength=\"96\">1234567890.123456</bd>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        BigDecimal value = buffer.readBigDecimal(96, WithOption.WithName("bd"));
        Assertions.assertEquals(new BigDecimal("1234567890.123456"), value);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // readString

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testReadString() throws Exception {
        String input = "<str dataType=\"string\" bitLength=\"40\">Hello</str>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        String value = buffer.readString(40, WithOption.WithName("str"));
        Assertions.assertEquals("Hello", value);
    }

    @Test
    public void testReadStringInvalidType() {
        String input = "<str dataType=\"uint\" bitLength=\"40\">Hello</str>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readString(40, WithOption.WithName("str")));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // contextPush/Pop

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testContextPushPop() throws Exception {
        String input = "<Outer><inner dataType=\"uint\" bitLength=\"16\">42</inner></Outer>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        buffer.pushContext(WithOption.WithName("Outer"));
        long v = buffer.readUnsignedInt(16, WithOption.WithName("inner"));
        Assertions.assertEquals(42L, v);
        buffer.popContext(WithOption.WithName("Outer"));
    }

    @Test
    public void testContextPopMismatch() throws Exception {
        String input = "<Outer><inner dataType=\"uint\" bitLength=\"16\">1</inner></Outer>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        buffer.pushContext(WithOption.WithName("Outer"));
        buffer.readUnsignedInt(16, WithOption.WithName("inner"));
        Assertions.assertThrows(BufferException.class, () -> buffer.popContext(WithOption.WithName("WrongOuter")));
    }

    @Test
    public void testContextPushMismatch() throws Exception {
        String input = "<Wrong><inner dataType=\"uint\" bitLength=\"16\">1</inner></Wrong>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.pushContext(WithOption.WithName("Right")));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Other ReadBuffer methods

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testComplex() throws Exception {
        String input = "<ModbusPDUReadCoilsRequest><startingAddress dataType=\"uint\" bitLength=\"16\">23</startingAddress><quantity dataType=\"uint\" bitLength=\"16\">42</quantity></ModbusPDUReadCoilsRequest>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));

        buffer.pushContext(WithOption.WithName("ModbusPDUReadCoilsRequest"));

        // Simple Field (startingAddress)
        long startingAddress = buffer.readUnsignedInt(16, WithOption.WithName("startingAddress"));
        Assertions.assertEquals(23, startingAddress);

        // Simple Field (quantity)
        long quantity = buffer.readUnsignedInt(16, WithOption.WithName("quantity"));
        Assertions.assertEquals(42, quantity);

        buffer.popContext();
    }

    @Test
    public void testMissingClosingTagEOF() {
        String input = "<str dataType=\"string\" bitLength=\"8\">A"; // No closing tag
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readString(8, WithOption.WithName("str")));
    }

    @Test
    public void testInvalidXmlStartsWithClosingTag() {
        String input = "</foo>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readBit(WithOption.WithName("foo")));
    }

    @Test
    public void testMissingDataTypeAttribute() {
        String input = "<uInt bitLength=\"32\">123</uInt>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readUnsignedInt(32, WithOption.WithName("uInt")));
    }

    @Test
    public void testMissingBitLengthAttribute() {
        String input = "<uInt dataType=\"uint\">123</uInt>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readUnsignedInt(32, WithOption.WithName("uInt")));
    }

    @Test
    public void testMissingNameOption() {
        String input = "<uInt dataType=\"uint\" bitLength=\"32\">123</uInt>";
        ReadBufferXmlBased buffer = new ReadBufferXmlBased(new ByteArrayInputStream(input.getBytes()));
        Assertions.assertThrows(BufferException.class, () -> buffer.readUnsignedInt(32));
    }

}
