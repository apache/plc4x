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
package org.apache.plc4x.java.ads.api.serial.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;

public class CRC extends UnsignedShortLEByteValue {

    public static final int NUM_BYTES = UnsignedShortLEByteValue.UNSIGNED_SHORT_LE_NUM_BYTES;

    private CRC(byte... values) {
        super(values);
    }

    private CRC(int value) {
        super(value);
    }

    private CRC(String length) {
        super(length);
    }

    private CRC(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static CRC of(byte... values) {
        return new CRC(values);
    }

    public static CRC of(int value) {
        return new CRC(value);
    }

    public static CRC of(String length) {
        return new CRC(length);
    }

    public static CRC of(ByteBuf byteBuf) {
        return new CRC(byteBuf);
    }
}
