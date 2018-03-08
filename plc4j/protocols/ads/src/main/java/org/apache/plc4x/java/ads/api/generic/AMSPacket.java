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
import org.apache.plc4x.java.ads.api.serial.AMSSerialFrame;
import org.apache.plc4x.java.ads.api.serial.types.FragmentNumber;
import org.apache.plc4x.java.ads.api.tcp.AMSTCPPacket;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

import static java.util.Objects.requireNonNull;

public abstract class AMSPacket implements ByteReadable {

    protected final AMSHeader amsHeader;

    protected AMSPacket(AMSHeader amsHeader) {
        this.amsHeader = amsHeader;
    }

    protected AMSPacket(AMSNetId targetAmsNetId, AMSPort targetAmsPort, AMSNetId sourceAmsNetId, AMSPort sourceAmsPort, State stateId, Invoke invokeId) {
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
    }

    public AMSHeader getAmsHeader() {
        return amsHeader;
    }

    public abstract ADSData getAdsData();

    protected ADSData buildADSData(ByteReadable... byteReadables) {
        return () -> buildByteBuff(byteReadables);
    }

    public AMSTCPPacket toAmstcpPacket() {
        return AMSTCPPacket.of(this);
    }

    public AMSSerialFrame toAmsSerialFrame(byte fragmentNumber) {
        return AMSSerialFrame.of(FragmentNumber.of(fragmentNumber), this);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(amsHeader, getAdsData());
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof AMSPacket)) return false;

        AMSPacket amsPacket = (AMSPacket) o;

        return amsHeader.equals(amsPacket.amsHeader);
    }

    @Override
    public int hashCode() {
        return amsHeader.hashCode();
    }

    @Override
    public String toString() {
        return getClass().getSimpleName() + "{" +
            "amsHeader=" + amsHeader +
            '}';
    }
}
