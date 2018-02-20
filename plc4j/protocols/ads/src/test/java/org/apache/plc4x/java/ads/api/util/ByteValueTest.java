package org.apache.plc4x.java.ads.api.util;

import io.netty.buffer.ByteBuf;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.math.BigInteger;

import static org.hamcrest.CoreMatchers.is;
import static org.hamcrest.CoreMatchers.not;
import static org.junit.Assert.*;

public class ByteValueTest {

    private ByteValue byteValue;
    private long upperBound = (long) Math.pow( 2, (8 * 4));

    @Before
    public void setUp() throws Exception {
        byteValue = new ByteValue((byte)0x1, (byte)0x2,(byte) 0x3, (byte)0x4);
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
        ByteValue.checkUnsignedBounds(upperBound-1, 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsLongNegative() {
        ByteValue.checkUnsignedBounds(-1, 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsLongTooBig() {
        ByteValue.checkUnsignedBounds(upperBound, 4);
    }

    @Test
    public void checkUnsignedBoundsBig() {
        ByteValue.checkUnsignedBounds(new BigInteger("0"), 4);
        ByteValue.checkUnsignedBounds(new BigInteger(Long.toString(upperBound-1)), 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsBigNegative() {
        ByteValue.checkUnsignedBounds(new BigInteger("-1"), 4);
    }

    @Test(expected = IllegalArgumentException.class)
    public void checkUnsignedBoundsBigTooBig() {
        ByteValue.checkUnsignedBounds(new BigInteger(Long.toString(upperBound)), 4);
    }

    @Test
    public void getBytes() {
        byte[] correct = {(byte)0x1, (byte)0x2, (byte)0x3, (byte)0x4};
        assertThat(byteValue.getBytes(), is(correct));
    }

    @Test
    public void getByteBuf() {
        byte[] correct = {(byte)0x1, (byte)0x2, (byte)0x3, (byte)0x4};
        ByteBuf data = byteValue.getByteBuf();

        assertThat(data.readableBytes(), is(4));
        assertThat(data.readByte(), is((byte)0x1));
        assertThat(data.readByte(), is((byte)0x2));
        assertThat(data.readByte(), is((byte)0x3));
        assertThat(data.readByte(), is((byte)0x4));
    }

    @Test
    public void equals() {
        ByteValue a = new ByteValue(((byte)0x1));
        ByteValue b = new ByteValue(((byte)0x1));
        ByteValue c = new ByteValue(((byte)0x2));
        byte array[] = {(byte)0x1};

        assertThat(a.equals(a), is(true));
        assertThat(a.equals(b), is(true));
        assertThat(a.equals(c), is(false));
        assertThat(a.equals(1), is(false));
        assertThat(a.equals((byte)1), is(false));
        assertThat(a.equals(array), is(false));
    }
    
}