/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.java.ads.model.generic;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.ads.model.util.ByteReadable;
import org.apache.plc4x.java.ads.model.util.ByteValue;

/**
 * AMS/TCP Header	6 bytes	contains the length of the data packet.
 */
public class AMSTCPHeader implements ByteReadable {

    final Reserved reserved;

    final Length length;

    public AMSTCPHeader(Length length) {
        this.reserved = Reserved.CONSTANT;
        this.length = length;
    }

    @Override
    public byte[] getBytes() {
        return getByteBuf().array();
    }

    @Override
    public ByteBuf getByteBuf() {
        return Unpooled.buffer()
            .writeBytes(reserved.getByteBuf())
            .writeBytes(length.getByteBuf());
    }

    /**
     * Size: 2 bytes
     * These bytes must be set to 0.
     */
    static class Reserved extends ByteValue {

        static final Reserved CONSTANT = new Reserved();

        public Reserved() {
            super((byte) 0x00, (byte) 0x00);
            assertLength(2);
        }
    }

    /**
     * Size: 4 bytes
     * This array contains the length of the data packet. It consists of the AMS-Header and the enclosed ADS data. The unit is bytes.
     */
    static class Length extends ByteValue {

        public Length(byte... value) {
            super(value);
            assertLength(4);
        }
    }
}
