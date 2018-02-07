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
package org.apache.plc4x.java.ads.api.generic;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.commands.ADSCommandType;
import org.apache.plc4x.java.ads.api.generic.calculated.CalculatedAMSHeader;
import org.apache.plc4x.java.ads.api.generic.calculated.CalculatedAMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

import static org.apache.plc4x.java.ads.api.util.ByteReadableUtils.buildByteBuff;

public abstract class AMSTCPPaket implements ByteReadable {
    private final AMSTCPHeader amstcpHeader;

    private final AMSHeader amsHeader;

    public AMSTCPPaket(AMSTCPHeader amstcpHeader, AMSHeader amsHeader) {
        this.amstcpHeader = amstcpHeader;
        this.amsHeader = amsHeader;
    }

    public AMSTCPPaket(AMSHeader amsHeader) {
        // It is important that we wrap the ads data call as this will initialized in the constructor
        // so this value will be null if we call adsData now.
        this.amstcpHeader = CalculatedAMSTCPHeader.of(amsHeader, () -> getAdsData().getCalculatedLength());
        this.amsHeader = amsHeader;
    }

    public AMSTCPPaket(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, Invoke invokeId) {
        if (!getClass().isAnnotationPresent(ADSCommandType.class)) {
            throw new IllegalArgumentException(ADSCommandType.class + " need to be present.");
        }
        this.amsHeader = CalculatedAMSHeader.of(
            targetAmsNetId,
            targetAmsPort,
            sourceAmsNetId,
            sourceAmsPort,
            getClass().getAnnotation(ADSCommandType.class).value(),
            State.DEFAULT,
            () -> DataLength.of(getAdsData().getCalculatedLength()),
            invokeId);
        this.amstcpHeader = CalculatedAMSTCPHeader.of(amsHeader, () -> getAdsData().getCalculatedLength());
    }

    public AMSTCPHeader getAmstcpHeader() {
        return amstcpHeader;
    }

    public AMSHeader getAmsHeader() {
        return amsHeader;
    }

    public abstract ADSData getAdsData();

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(amstcpHeader, amsHeader, getAdsData());
    }

    protected ADSData buildADSData(ByteReadable... byteReadables) {
        return () -> buildByteBuff(byteReadables);
    }

}
