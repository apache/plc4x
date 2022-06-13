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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.junit.Test;

import static org.hamcrest.Matchers.is;
import static org.junit.Assert.assertThat;

public class AmsErrorTest {

    @Test
    public void errorBytes() {
        // note bytes in reverse order
        AmsError error = AmsError.of((byte)0x01, (byte)0x20, (byte)0x00, (byte)0x00);
        assertThat(error.getAsLong(), is(0x2001L));
    }

    @Test
    public void errorLong() {
        AmsError error = AmsError.of(0xFF02L);
        assertThat(error.getAsLong(), is(0xFF02L));
    }

    @Test
    public void errorLongBig() {
        AmsError error = AmsError.of(0xFFFFFFFFL);
        assertThat(error.getAsLong(), is(0xFFFFFFFFL));
    }
    
    @Test
    public void errorString() {
        AmsError error = AmsError.of("255");
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

        AmsError error = AmsError.of(buffer);
        assertThat(error.getAsLong(), is(260L));
    }

    @Test(expected = NumberFormatException.class)
    public void noHex() {
        AmsError error = AmsError.of("0xFF000000");
    }
    
    @Test(expected = IllegalArgumentException.class)
    public void errorLongTooBig() {
        AmsError error = AmsError.of(0x100000000L);
    }
    
    @Test(expected = IllegalArgumentException.class)
    public void errorNegative() {
        AmsError error = AmsError.of(-1);
    }
    
    @Test
    public void equals() {
        AmsError a = AmsError.of((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4);
        AmsError b = AmsError.of((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4);
        AmsError c = AmsError.of((byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0xFF);
        byte array[] = {(byte) 0x1, (byte) 0x2, (byte) 0x3, (byte) 0x4};

        assertThat(a.equals(a), is(true));
        assertThat(a.equals(b), is(true));
        assertThat(a.equals(c), is(false));
        assertThat(a.equals(1), is(false));
        assertThat(a.equals((byte) 1), is(false));
        assertThat(a.equals(array), is(false));
    }
}