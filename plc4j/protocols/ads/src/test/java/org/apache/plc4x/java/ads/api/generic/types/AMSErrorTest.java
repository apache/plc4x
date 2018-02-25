package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.junit.Test;

import static org.hamcrest.Matchers.is;
import static org.junit.Assert.*;

public class AMSErrorTest {

    @Test
    public void errorBytes() {
        // note bytes in reverse order
        AMSError error = AMSError.of((byte)0x01, (byte)0x20, (byte)0x00, (byte)0x00);
        assertThat(error.getAsLong(), is(0x2001L));
    }

    @Test
    public void errorLong() {
        AMSError error = AMSError.of(0xFF02L);
        assertThat(error.getAsLong(), is(0xFF02L));
    }

    @Test
    public void errorLongBig() {
        AMSError error = AMSError.of(0xFFFFFFFFL);
        assertThat(error.getAsLong(), is(0xFFFFFFFFL));
    }
    
    @Test
    public void errorString() {
        AMSError error = AMSError.of("255");
        assertThat(error.getAsLong(), is(0xFFL));
    }

    @Test
    public void errorByteBuf() {
        ByteBuf buffer = Unpooled.buffer();

        // note bytes in reverse order
        buffer.writeByte((byte)0x04);
        buffer.writeByte((byte)0x01);
        buffer.writeByte((byte)0x00);
        buffer.writeByte((byte)0x00);

        AMSError error = AMSError.of(buffer);
        assertThat(error.getAsLong(), is(260L));
    }

    @Test(expected = NumberFormatException.class)
    public void noHex() {
        AMSError error = AMSError.of("0xFF000000");
    }
    
    @Test(expected = IllegalArgumentException.class)
    public void errorLongTooBig() {
        AMSError error = AMSError.of(0x100000000L);
    }
    
    @Test(expected = IllegalArgumentException.class)
    public void errorNegative() {
        AMSError error = AMSError.of(-1);
    }
    
    @Test
    public void equals() {
        AMSError a = AMSError.of((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4);
        AMSError b = AMSError.of((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4);
        AMSError c = AMSError.of((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0xFF);
        byte array[] = {(byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4};

        assertThat(a.equals(a), is(true));
        assertThat(a.equals(b), is(true));
        assertThat(a.equals(c), is(false));
        assertThat(a.equals(1), is(false));
        assertThat(a.equals((byte) 1), is(false));
        assertThat(a.equals(array), is(false));
    }
}