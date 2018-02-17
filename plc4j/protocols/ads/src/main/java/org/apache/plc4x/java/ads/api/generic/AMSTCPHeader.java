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
import org.apache.plc4x.java.ads.api.generic.types.TcpLength;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.ads.api.util.ByteValue;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static java.util.Objects.requireNonNull;
import static org.apache.plc4x.java.ads.api.util.ByteReadableUtils.buildByteBuff;

/**
 * AMS/TCP Header	6 bytes	contains the tcpLength of the data packet.
 */
public class AMSTCPHeader implements ByteReadable {

    private final Reserved reserved;

    private final TcpLength tcpLength;

    private final LengthSupplier[] lengthSuppliers;

    private AMSTCPHeader(TcpLength tcpLength) {
        this.reserved = requireNonNull(Reserved.CONSTANT);
        this.tcpLength = requireNonNull(tcpLength);
        lengthSuppliers = null;

    }

    private AMSTCPHeader(LengthSupplier... lengthSuppliers) {
        this.reserved = requireNonNull(Reserved.CONSTANT);
        this.tcpLength = null;
        this.lengthSuppliers = requireNonNull(lengthSuppliers);
    }

    public static AMSTCPHeader of(TcpLength tcpLength) {
        return new AMSTCPHeader(tcpLength);
    }

    public static AMSTCPHeader of(long length) {
        return new AMSTCPHeader(TcpLength.of(length));
    }

    public static AMSTCPHeader of(LengthSupplier... lengthSuppliers) {
        return new AMSTCPHeader(lengthSuppliers);
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
        if (this == o) return true;
        if (!(o instanceof AMSTCPHeader)) return false;

        AMSTCPHeader that = (AMSTCPHeader) o;

        return getTcpLength().equals(that.getTcpLength());
    }

    @Override
    public int hashCode() {
        int result = reserved != null ? reserved.hashCode() : 0;
        result = 31 * result + (getTcpLength() != null ? getTcpLength().hashCode() : 0);
        return result;
    }

    @Override
    public String toString() {
        return "AMSTCPHeader{" +
            "reserved=" + reserved +
            ", tcpLength=" + getTcpLength() +
            '}';
    }
}
