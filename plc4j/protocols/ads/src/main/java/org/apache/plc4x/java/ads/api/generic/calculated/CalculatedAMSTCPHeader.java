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
package org.apache.plc4x.java.ads.api.generic.calculated;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.generic.AMSTCPHeader;
import org.apache.plc4x.java.ads.api.generic.types.Length;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;

import static org.apache.plc4x.java.ads.api.util.ByteReadableUtils.buildByteBuff;

/**
 * AMS/TCP Header	6 bytes	contains the length of the data packet.
 * This Header is caluclated. Can be used when sending.
 */
public class CalculatedAMSTCPHeader extends AMSTCPHeader {

    protected final LengthSupplier[] lengthSupplier;

    protected CalculatedAMSTCPHeader(LengthSupplier... lengthSupplier) {
        super(null);
        this.lengthSupplier = lengthSupplier;
    }

    public static CalculatedAMSTCPHeader of(LengthSupplier... lengthSupplier) {
        return new CalculatedAMSTCPHeader(lengthSupplier);
    }

    @Override
    public ByteBuf getByteBuf() {
        return buildByteBuff(reserved, Length.of(getCalculatedLength()));
    }

    @Override
    public long getCalculatedLength() {
        long aggregateLength = 0;
        for (LengthSupplier supplier : lengthSupplier) {
            aggregateLength += supplier.getCalculatedLength();
        }
        return aggregateLength;
    }

    @Override
    public String toString() {
        return "CalculatedAMSTCPHeader{" +
            "reserved=" + reserved +
            ", length=" + getCalculatedLength() +
            '}';
    }
}
