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
package org.apache.plc4x.java.ads.api.tcp;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.tcp.types.TcpLength;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;

/**
 * AMS/TCP Header	6 bytes	contains the tcpLength of the data packet.
 */
public class AmsTcpHeader implements ByteReadable {

    /**
     * These bytes must be set to 0.
     */
    private final Reserved reserved;

    /**
     * This array contains the length of the data packet. It consists of the AMS-Header and the enclosed ADS data. The unit is bytes.
     * is null if length is supplied by {@link #lengthSuppliers}.
     */
    private final TcpLength tcpLength;

    /**
     * This array contains the length of the data packet. It consists of the AMS-Header and the enclosed ADS data. The unit is bytes.
     * is null if length is supplied by {@link #tcpLength}.
     */
    private final transient LengthSupplier[] lengthSuppliers;

    private AmsTcpHeader(TcpLength tcpLength) {
        this.reserved = Reserved.CONSTANT;
        this.tcpLength = requireNonNull(tcpLength);
        this.lengthSuppliers = null;
    }

    private AmsTcpHeader(LengthSupplier... lengthSuppliers) {
        this.reserved = Reserved.CONSTANT;
        this.tcpLength = null;
        this.lengthSuppliers = requireNonNull(lengthSuppliers);
    }

    public static AmsTcpHeader of(TcpLength tcpLength) {
        return new AmsTcpHeader(tcpLength);
    }

    public static AmsTcpHeader of(long length) {
        return new AmsTcpHeader(TcpLength.of(length));
    }

    public static AmsTcpHeader of(LengthSupplier... lengthSuppliers) {
        return new AmsTcpHeader(lengthSuppliers);
    }

    public TcpLength getTcpLength() {
        return lengthSuppliers == null ? tcpLength : TcpLength.of(getCalculatedLength());
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(reserved, getTcpLength());
    }

    /**
     * Size: 2 bytes
     * These bytes must be set to 0.
     */
    public static class Reserved extends ByteValue {

        public static final int NUM_BYTES = 2;

        private static final Reserved CONSTANT = new Reserved();

        private Reserved() {
            super((byte) 0x00, (byte) 0x00);
            assertLength(NUM_BYTES);
        }
    }

    @Override
    public long getCalculatedLength() {
        if (lengthSuppliers == null) {
            return tcpLength.getAsLong();
        } else {
            long aggregateLength = 0;
            for (LengthSupplier supplier : lengthSuppliers) {
                aggregateLength += supplier.getCalculatedLength();
            }
            return aggregateLength;
        }
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof AmsTcpHeader))
            return false;

        AmsTcpHeader that = (AmsTcpHeader) o;

        return getTcpLength().equals(that.getTcpLength());
    }

    @Override
    public int hashCode() {
        int result = reserved.hashCode();
        result = 31 * result + getTcpLength().hashCode();
        return result;
    }

    @Override
    public String toString() {
        return "AmsTcpHeader{" +
            "reserved=" + reserved +
            ", tcpLength=" + getTcpLength() +
            '}';
    }
}
