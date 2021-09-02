/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

import static java.util.Objects.requireNonNull;

public class Result extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    private Result(byte... values) {
        super(values);
    }

    private Result(long value) {
        super(value);
    }

    private Result(String value) {
        super(value);
    }

    private Result(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static Result of(byte... values) {
        return new Result(values);
    }

    public static Result of(long value) {
        return new Result(value);
    }

    public static Result of(String value) {
        return new Result(value);
    }

    public static Result of(ByteBuf byteBuf) {
        return new Result(byteBuf);
    }

    public static Result of(AdsReturnCode adsReturnCode) {
        return of(requireNonNull(adsReturnCode).getHex());
    }

    public AdsReturnCode toAdsReturnCode() {
        return AdsReturnCode.of(getAsLong());
    }

    @Override
    public String toString() {
        return "Result{" + toAdsReturnCode() + "}";
    }
}
