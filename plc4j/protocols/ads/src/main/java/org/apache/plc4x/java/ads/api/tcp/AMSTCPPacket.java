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
package org.apache.plc4x.java.ads.api.tcp;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.ADSData;
import org.apache.plc4x.java.ads.api.generic.AMSHeader;
import org.apache.plc4x.java.ads.api.generic.AMSPacket;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

import static java.util.Objects.requireNonNull;

public class AMSTCPPacket implements ByteReadable {
    private final AMSTCPHeader amsTcpHeader;

    private final AMSPacket amsPacket;

    private AMSTCPPacket(AMSTCPHeader amsTcpHeader, AMSPacket amsPacket) {
        this.amsTcpHeader = requireNonNull(amsTcpHeader);
        this.amsPacket = amsPacket;
    }

    private AMSTCPPacket(AMSPacket amsPacket) {
        this.amsPacket = amsPacket;
        // It is important that we wrap the ads data call as this will initialized in the constructor
        // so this value will be null if we call adsData now.
        this.amsTcpHeader = AMSTCPHeader.of(requireNonNull(amsPacket.getAmsHeader()), () -> getAdsData().getCalculatedLength());
    }

    public AMSTCPHeader getAmsTcpHeader() {
        return amsTcpHeader;
    }

    public AMSPacket getAmsPacket() {
        return amsPacket;
    }

    public AMSHeader getAmsHeader() {
        return amsPacket.getAmsHeader();
    }

    public ADSData getAdsData() {
        return amsPacket.getAdsData();
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(amsTcpHeader, amsPacket);
    }

    public static AMSTCPPacket of(AMSTCPHeader amsTcpHeader, AMSPacket amsPacket) {
        return new AMSTCPPacket(amsTcpHeader, amsPacket);
    }

    public static AMSTCPPacket of(AMSPacket amsPacket) {
        return new AMSTCPPacket(amsPacket);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof AMSTCPPacket)) return false;

        AMSTCPPacket that = (AMSTCPPacket) o;

        if (!amsTcpHeader.equals(that.amsTcpHeader)) return false;
        return amsPacket.equals(that.amsPacket);
    }

    @Override
    public int hashCode() {
        int result = amsTcpHeader.hashCode();
        result = 31 * result + amsPacket.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AMSTCPPacket{" +
            "amsTcpHeader=" + amsTcpHeader +
            ", amsPacket=" + amsPacket +
            '}';
    }
}
