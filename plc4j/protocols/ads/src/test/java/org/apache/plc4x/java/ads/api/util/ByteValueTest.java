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
package org.apache.plc4x.java.ads.api.util;

import io.netty.buffer.ByteBuf;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.math.BigInteger;

import static java.util.Arrays.copyOfRange;
import static org.hamcrest.CoreMatchers.is;
import static org.junit.Assert.assertThat;

public class ByteValueTest {

    private ByteValue byteValue;
    private long upperBound = (long) Math.pow(2, (8 * 4));

    @Before
    public void setUp() throws Exception {
        byteValue = new ByteValue((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4);
    }

    @After
    public void tearDown() throws Exception {
        byteValue = null;
    }

    @Test
    public void assertCorrectLength() {
        byteValue.assertLength(4); // no exception
    }

    @Test(expected = IllegalArgumentException.class)
    public void assertCorrectLengthException() {
        byteValue.assertLength(3);
    }

    @Test
    public void checkUnsignedBoundsLong() {
        ByteValue.checkUnsignedBounds(0, 4);
        ByteValue.checkUnsignedBounds(upperBound - 1, 4);
    }

    @Test
    public void checkUnsignedBoundsLongHex() {
        // Hex representation to visualize valid bounds in bytes
        ByteValue.checkUnsignedBounds(0x0_00_00, 2);
        ByteValue.checkUnsignedBounds(0x0_FF_FF, 2);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsLongNegative() {
        ByteValue.checkUnsignedBounds(-1, 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsLongTooBig() {
        ByteValue.checkUnsignedBounds(upperBound, 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsLongTooBigHex() {
        ByteValue.checkUnsignedBounds(0x1_00_00, 2);
    }

    @Test
    public void checkUnsignedBoundsBig() {
        ByteValue.checkUnsignedBounds(new BigInteger("0"), 4);
        ByteValue.checkUnsignedBounds(new BigInteger(Long.toString(upperBound - 1)), 4);
    }

    @Test
    public void checkUnsignedBoundsBigHex() {
        ByteValue.checkUnsignedBounds(BigInteger.valueOf(0x0_00_00), 2);
        ByteValue.checkUnsignedBounds(BigInteger.valueOf(0x0_FF_FF), 2);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsBigNegative() {
        ByteValue.checkUnsignedBounds(new BigInteger("-1"), 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsBigTooBig() {
        ByteValue.checkUnsignedBounds(new BigInteger(Long.toString(upperBound)).add(BigInteger.ONE), 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsBigTooBigHex() {
        ByteValue.checkUnsignedBounds(BigInteger.valueOf(0x1_00_00), 2);
    }

    @Test
    public void getBytes() {
        byte[] correct = {(byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4};
        assertThat(byteValue.getBytes(), is(correct));
    }

    @Test
    public void getByteBuf() {
        byte[] correct = {(byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4};
        ByteBuf data = byteValue.getByteBuf();

        assertThat(data.readableBytes(), is(4));
        assertThat(data.readByte(), is((byte) 0x1));
        assertThat(data.readByte(), is((byte) 0x2));
        assertThat(data.readByte(), is((byte) 0x3));
        assertThat(data.readByte(), is((byte) 0x4));
        assertThat(copyOfRange(data.array(), 0, 4), is(correct));
    }

    @Test
    public void equals() {
        ByteValue a = new ByteValue(((byte) 0x1));
        ByteValue b = new ByteValue(((byte) 0x1));
        ByteValue c = new ByteValue(((byte) 0x2));
        byte array[] = {(byte) 0x1};

        assertThat(a.equals(a), is(true));
        assertThat(a.equals(b), is(true));
        assertThat(a.equals(c), is(false));
        assertThat(a.equals(1), is(false));
        assertThat(a.equals((byte) 1), is(false));
        assertThat(a.equals(array), is(false));
    }

}