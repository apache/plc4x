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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.commands.types.AdsReturnCode;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

import static java.util.Objects.requireNonNull;

public class AmsError extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    public static final AmsError NONE = of(0);

    private AmsError(byte... values) {
        super(values);
    }

    private AmsError(long value) {
        super(value);
    }

    private AmsError(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static AmsError of(byte... values) {
        return new AmsError(values);
    }

    public static AmsError of(long errorCode) {
        return new AmsError(errorCode);
    }

    public static AmsError of(AdsReturnCode errorCode) {
        return new AmsError(requireNonNull(errorCode).getHex());
    }

    public static AmsError of(String errorCode) {
        return of(Long.parseLong(errorCode));
    }

    public static AmsError of(ByteBuf byteBuf) {
        return new AmsError(byteBuf);
    }

    public AdsReturnCode toAdsReturnCode() {
        return AdsReturnCode.of(getAsLong());
    }

    @Override
    public String toString() {
        return "AmsError{" + toAdsReturnCode() + "}";
    }
}
