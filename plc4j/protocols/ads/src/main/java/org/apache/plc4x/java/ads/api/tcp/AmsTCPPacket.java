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
import org.apache.plc4x.java.ads.api.tcp.types.UserData;
import org.apache.plc4x.java.ads.api.util.ByteReadable;

import static java.util.Objects.requireNonNull;

public class AmsTCPPacket implements ByteReadable {
    /**
     * The ams - tcp to be sent.
     */
    private final AmsTcpHeader amsTcpHeader;

    /**
     * The AMS packet to be sent.
     */
    private final UserData userData;

    private AmsTCPPacket(AmsTcpHeader amsTcpHeader, UserData userData) {
        this.amsTcpHeader = requireNonNull(amsTcpHeader);
        this.userData = requireNonNull(userData);
    }

    private AmsTCPPacket(UserData userData) {
        this.userData = requireNonNull(userData);
        // It is important that we wrap the ads data call as this will initialized in the constructor
        // so this value will be null if we call adsData now.
        this.amsTcpHeader = AmsTcpHeader.of(requireNonNull(userData));
    }

    public AmsTcpHeader getAmsTcpHeader() {
        return amsTcpHeader;
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(amsTcpHeader, userData);
    }

    public static AmsTCPPacket of(UserData userData) {
        return new AmsTCPPacket(userData);
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof AmsTCPPacket)) {
            return false;
        }

        AmsTCPPacket that = (AmsTCPPacket) o;

        if (!amsTcpHeader.equals(that.amsTcpHeader)) {
            return false;
        }
        return userData.equals(that.userData);
    }

    @Override
    public int hashCode() {
        int result = amsTcpHeader.hashCode();
        result = 31 * result + userData.hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AmsTCPPacket{" +
            "amsTcpHeader=" + amsTcpHeader +
            ", amsPacket=" + userData +
            '}';
    }
}
