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

import java.math.BigDecimal;
import java.math.BigInteger;

public class WriteBufferXmlBasedTest {

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // writeBit

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testWriteBit() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();

        buffer.writeBit(true, WithOption.WithName("bitField"));

        String expected = "<bitField dataType=\"bit\" bitLength=\"1\">true</bitField>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(1, buffer.getPositionInBits());
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // writeBits

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testWriteBits() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();

        byte[] bytes = "Hurz".getBytes();
        buffer.writeBits(32, bytes, WithOption.WithName("bitsField"));

        String expected = "<bitsField dataType=\"byte\" bitLength=\"32\">0x4875727a</bitsField>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        // Note: Implementation advances by value.length * 8 bits
        Assertions.assertEquals(32, buffer.getPositionInBits());
    }

    @Test
    public void testWriteBitsNotEnoughBytes() {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        byte[] tooShort = new byte[]{0x01};
        Assertions.assertThrows(BufferException.class, () -> buffer.writeBits(32, tooShort, WithOption.WithName("bitsField")));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // writeUnsigned*

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testWriteUnsignedByte() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeUnsignedByte(7, (byte) 127, WithOption.WithName("uByte"));
        String expected = "<uByte dataType=\"uint\" bitLength=\"7\">127</uByte>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(7, buffer.getPositionInBits());
    }

    @Test
    public void testWriteUnsignedShort() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeUnsignedShort(15, (short) 32767, WithOption.WithName("uShort"));
        String expected = "<uShort dataType=\"uint\" bitLength=\"15\">32767</uShort>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(15, buffer.getPositionInBits());
    }

    @Test
    public void testWriteUnsignedInt() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeUnsignedInt(31, 123456789, WithOption.WithName("uInt"));
        String expected = "<uInt dataType=\"uint\" bitLength=\"31\">123456789</uInt>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(31, buffer.getPositionInBits());
    }

    @Test
    public void testWriteUnsignedLong() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeUnsignedLong(63, 8446744073709551615L, WithOption.WithName("uLong"));
        String expected = "<uLong dataType=\"uint\" bitLength=\"63\">8446744073709551615</uLong>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(63, buffer.getPositionInBits());
    }

    @Test
    public void testWriteUnsignedBigInteger() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        BigInteger big = new BigInteger("340282366920938463463374607431768211455");
        buffer.writeUnsignedBigInteger(128, big, WithOption.WithName("uBig"));
        String expected = "<uBig dataType=\"uint\" bitLength=\"128\">340282366920938463463374607431768211455</uBig>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(128, buffer.getPositionInBits());
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // writeSigned*

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testWriteSignedByte() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeSignedByte(8, (byte) -42, WithOption.WithName("sByte"));
        String expected = "<sByte dataType=\"int\" bitLength=\"8\">-42</sByte>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(8, buffer.getPositionInBits());
    }

    @Test
    public void testWriteSignedShort() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeSignedShort(16, (short) -12345, WithOption.WithName("sShort"));
        String expected = "<sShort dataType=\"int\" bitLength=\"16\">-12345</sShort>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(16, buffer.getPositionInBits());
    }

    @Test
    public void testWriteSignedInt() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeSignedInt(32, -123456, WithOption.WithName("sInt"));
        String expected = "<sInt dataType=\"int\" bitLength=\"32\">-123456</sInt>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(32, buffer.getPositionInBits());
    }

    @Test
    public void testWriteSignedLong() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeSignedLong(64, -1234567890123456789L, WithOption.WithName("sLong"));
        String expected = "<sLong dataType=\"int\" bitLength=\"64\">-1234567890123456789</sLong>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(64, buffer.getPositionInBits());
    }

    @Test
    public void testWriteSignedBigInteger() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        BigInteger big = new BigInteger("-123456789012345678901234567890");
        buffer.writeSignedBigInteger(128, big, WithOption.WithName("sBig"));
        String expected = "<sBig dataType=\"int\" bitLength=\"128\">-123456789012345678901234567890</sBig>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(128, buffer.getPositionInBits());
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // writeFloat/Double/BigDecimal

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testWriteFloat() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeFloat(32, 3.14159f, WithOption.WithName("flt"));
        String expected = "<flt dataType=\"float\" bitLength=\"32\">3.14159</flt>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(32, buffer.getPositionInBits());
    }

    @Test
    public void testWriteDouble() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeDouble(64, 2.718281828459045, WithOption.WithName("dbl"));
        String expected = "<dbl dataType=\"float\" bitLength=\"64\">2.718281828459045</dbl>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(64, buffer.getPositionInBits());
    }

    @Test
    public void testWriteBigDecimal() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeBigDecimal(96, new BigDecimal("1234567890.123456"), WithOption.WithName("bd"));
        String expected = "<bd dataType=\"float\" bitLength=\"96\">1234567890.123456</bd>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(96, buffer.getPositionInBits());
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // writeString

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testWriteStringDefaultEncoding() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeString(40, "Hello", WithOption.WithName("str"));
        String expected = "<str dataType=\"string\" bitLength=\"40\" encoding=\"UTF8\">Hello</str>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(40, buffer.getPositionInBits());
    }

    @Test
    public void testWriteStringWithEncodingOption() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        buffer.writeString(24, "Hi!", WithOption.WithName("str"), WithOption.WithEncoding("UTF16"));
        String expected = "<str dataType=\"string\" bitLength=\"24\" encoding=\"UTF16\">Hi!</str>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
        Assertions.assertEquals(24, buffer.getPositionInBits());
    }

    @Test
    public void testWriteStringCleansInvalidCharacters() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        String valueWithInvalid = "Hello\u0000World"; // contains a null byte which is invalid in XML 1.0
        buffer.writeString(80, valueWithInvalid, WithOption.WithName("str"));
        String expected = "<str dataType=\"string\" bitLength=\"80\" encoding=\"UTF8\">HelloWorld</str>\n";
        Assertions.assertEquals(expected, buffer.getXmlString());
    }

    @Test
    public void testMissingNameOption() {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();
        Assertions.assertThrows(BufferException.class, () -> buffer.writeUnsignedInt(32, 1));
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Context handling

    /// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Test
    public void testContextPushPopAndInnerWrite() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();

        buffer.pushContext(WithOption.WithName("Outer"));
        buffer.writeUnsignedInt(15, 42, WithOption.WithName("inner"));
        buffer.popContext(WithOption.WithName("Outer"));

        String expected = String.join("\n",
            "<Outer>",
            "  <inner dataType=\"uint\" bitLength=\"15\">42</inner>",
            "  </Outer>",
            ""
        );
        Assertions.assertEquals(expected, buffer.getXmlString());
    }

    @Test
    public void testContextWithListAndAdditionalStringRepresentation() throws Exception {
        WriteBufferXmlBased buffer = new WriteBufferXmlBased();

        buffer.pushContext(WithOption.WithName("Items"), WithOption.WithRenderAsList(true));
        buffer.writeUnsignedInt(7, 23, WithOption.WithName("item"), WithOption.WithAdditionalStringRepresentation("0x17"));
        buffer.popContext(WithOption.WithName("Items"));

        String expected = String.join("\n",
            "<Items isList=\"true\">",
            "  <item dataType=\"uint\" bitLength=\"7\" stringRepresentation=\"0x17\">23</item>",
            "  </Items>",
            ""
        );
        Assertions.assertEquals(expected, buffer.getXmlString());
    }
}
