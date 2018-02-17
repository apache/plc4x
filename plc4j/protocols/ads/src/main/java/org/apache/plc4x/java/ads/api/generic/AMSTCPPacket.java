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
import org.apache.plc4x.java.ads.api.generic.types.*;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

import static java.util.Objects.requireNonNull;
import static org.apache.plc4x.java.ads.api.util.ByteReadableUtils.buildByteBuff;

public abstract class AMSTCPPacket implements ByteReadable {
    private final AMSTCPHeader amsTcpHeader;

    private final AMSHeader amsHeader;

    protected AMSTCPPacket(AMSTCPHeader amsTcpHeader, AMSHeader amsHeader) {
        this.amsTcpHeader = requireNonNull(amsTcpHeader);
        this.amsHeader = requireNonNull(amsHeader);
    }

    protected AMSTCPPacket(AMSHeader amsHeader) {
        // It is important that we wrap the ads data call as this will initialized in the constructor
        // so this value will be null if we call adsData now.
        this.amsTcpHeader = AMSTCPHeader.of(requireNonNull(amsHeader), () -> getAdsData().getCalculatedLength());
        this.amsHeader = requireNonNull(amsHeader);
    }

    protected AMSTCPPacket(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, State stateId, Invoke invokeId) {
        if (!getClass().isAnnotationPresent(ADSCommandType.class)) {
            throw new IllegalArgumentException(ADSCommandType.class + " need to be present.");
        }
        this.amsHeader = AMSHeader.of(
            requireNonNull(targetAmsNetId),
            requireNonNull(targetAmsPort),
            requireNonNull(sourceAmsNetId),
            requireNonNull(sourceAmsPort),
            requireNonNull(getClass().getAnnotation(ADSCommandType.class).value()),
            requireNonNull(stateId),
            () -> getAdsData().getCalculatedLength(),
            requireNonNull(AMSError.NONE),
            requireNonNull(invokeId));
        this.amsTcpHeader = AMSTCPHeader.of(amsHeader, () -> getAdsData().getCalculatedLength());
    }

    public AMSTCPHeader getAmsTcpHeader() {
        return amsTcpHeader;
    }

    public AMSHeader getAmsHeader() {
        return amsHeader;
    }

    protected abstract ADSData getAdsData();

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(amsTcpHeader, amsHeader, getAdsData());
    }

    protected ADSData buildADSData(ByteReadable... byteReadables) {
        return () -> buildByteBuff(byteReadables);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof AMSTCPPacket)) return false;

        AMSTCPPacket that = (AMSTCPPacket) o;

        if (!amsTcpHeader.equals(that.amsTcpHeader)) return false;
        return amsHeader.equals(that.amsHeader);
    }

    @Override
    public int hashCode() {
        int result = amsTcpHeader.hashCode();
        result = 31 * result + amsHeader.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return getClass().getSimpleName() + "{" +
            "amsTcpHeader=" + amsTcpHeader +
            ", amsHeader=" + amsHeader +
            '}';
    }
}
