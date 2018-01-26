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
import org.apache.plc4x.java.ads.model.generic.types.Length;
import org.apache.plc4x.java.ads.model.util.ByteReadable;
import org.apache.plc4x.java.ads.model.util.ByteValue;

/**
 * AMS/TCP Header	6 bytes	contains the length of the data packet.
 */
public class AMSTCPHeader implements ByteReadable {

    private final Reserved reserved;

    private final Length length;

    public AMSTCPHeader(Length length) {
        this.reserved = Reserved.CONSTANT;
        this.length = length;
    }

    public static AMSTCPHeader of(int length) {
        return new AMSTCPHeader(Length.of(length));
    }

    @Override
    public ByteBuf getByteBuf() {
        return AMSTCPPaket.buildByteBuff(reserved, length);
    }

    /**
     * Size: 2 bytes
     * These bytes must be set to 0.
     */
    private static class Reserved extends ByteValue {

        private static final Reserved CONSTANT = new Reserved();

        private Reserved() {
            super((byte) 0x00, (byte) 0x00);
            assertLength(2);
        }
    }

}
