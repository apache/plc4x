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

public abstract class AMSTCPPaket implements ByteReadable {
    private final AMSTCPHeader amstcpHeader;

    private final AMSHeader amsHeader;

    public AMSTCPPaket(AMSTCPHeader amstcpHeader, AMSHeader amsHeader) {
        this.amstcpHeader = amstcpHeader;
        this.amsHeader = amsHeader;
    }

    public AMSTCPHeader getAmstcpHeader() {
        return amstcpHeader;
    }

    public AMSHeader getAmsHeader() {
        return amsHeader;
    }

    public abstract ADSData getAdsData();

    @Override
    public byte[] getBytes() {
        return getByteBuf().array();
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(amstcpHeader, amsHeader, getAdsData());
    }

    protected ADSData buildADSData(ByteReadable... byteReadables) {
        return () -> buildByteBuff(byteReadables).array();
    }

    public static ByteBuf buildByteBuff(ByteReadable... byteReadables) {
        ByteBuf buffer = Unpooled.buffer();
        for (ByteReadable byteReadable : byteReadables) {
            buffer.writeBytes(byteReadable.getByteBuf());
        }
        return buffer;
    }
}
