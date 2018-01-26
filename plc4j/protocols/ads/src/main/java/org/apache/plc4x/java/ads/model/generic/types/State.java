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
package org.apache.plc4x.java.ads.model.generic.types;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.ads.model.util.ByteReadable;

import java.nio.ByteBuffer;

/**
 * 2 bytes	see below.
 * <p>
 * State Flags
 * Flag     Description
 * 0x0001	0: Request / 1: Response
 * 0x0004	ADS command
 * The first bit marks, whether it´s a request or response. The third bit must be set to 1, to exchange data with ADS commands. The other bits are not defined or were used for other internal purposes.
 * <p>
 * Therefore the other bits must be set to 0!
 * <p>
 * Flag     Description
 * 0x000x	TCP Protocol
 * 0x004x	UDP Protocol
 * Bit number 7 marks, if it should be transfered with TCP or UDP.
 */
public enum State implements ByteReadable {
    ADS_REQUEST_TCP(0x0004),
    ADS_RESPONSE_TCP(0x0005),
    ADS_REQUEST_UDP(0x0044),
    ADS_RESPONSE_UDP(0x0045),
    UNKNOWN(0xffff);

    public static final int NUM_BYTES = 4;

    final byte[] value;

    State(int value) {
        this.value = ByteBuffer.allocate(NUM_BYTES).putInt(value).array();
    }

    @Override
    public byte[] getBytes() {
        return value;
    }

    public ByteBuf getByteBuf() {
        return Unpooled.buffer().writeBytes(value);
    }
}
