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
package org.apache.plc4x.java.ads.api.serial.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;

public class MagicCookie extends UnsignedShortLEByteValue {

    public static final int NUM_BYTES = UnsignedShortLEByteValue.UNSIGNED_SHORT_LE_NUM_BYTES;

    private MagicCookie(byte... values) {
        super(values);
    }

    private MagicCookie(int value) {
        super(value);
    }

    private MagicCookie(String length) {
        super(length);
    }

    private MagicCookie(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static MagicCookie of(byte... values) {
        return new MagicCookie(values);
    }

    public static MagicCookie of(int value) {
        return new MagicCookie(value);
    }

    public static MagicCookie of(String length) {
        return new MagicCookie(length);
    }

    public static MagicCookie of(ByteBuf byteBuf) {
        return new MagicCookie(byteBuf);
    }
}
