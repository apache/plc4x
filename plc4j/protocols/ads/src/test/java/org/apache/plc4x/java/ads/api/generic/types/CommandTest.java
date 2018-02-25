package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.junit.Test;

import static org.hamcrest.Matchers.is;
import static org.junit.Assert.*;

public class CommandTest {

    @Test
    public void getBytes() {
        byte[] result = {(byte)0x01, (byte)0x00};
        Command command = Command.ofInt("1");
        assertThat(command.getBytes(), is(result));
    }

    @Test
    public void getByteBuf() {
        ByteBuf result = Unpooled.buffer();
        result.writeByte(0x02);
        result.writeByte(0x00);
        Command command = Command.ofInt("2");
        assertThat(command.getByteBuf(), is(result));
    }

    @Test(expected = IllegalStateException.class)
    public void getBytesUnknown() {
        Command command = Command.UNKNOWN;
        command.getBytes();
    }

    @Test(expected = IllegalStateException.class)
    public void getByteBufUnknown() {
        Command command = Command.UNKNOWN;
        command.getByteBuf();
    }
}