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
package org.apache.plc4x.java.ads.api.commands.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

public class ReadLength extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    private ReadLength(byte... values) {
        super(values);
    }

    private ReadLength(long value) {
        super(value);
    }

    private ReadLength(String length) {
        super(length);
    }

    private ReadLength(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static ReadLength of(byte... values) {
        return new ReadLength(values);
    }

    public static ReadLength of(long value) {
        return new ReadLength(value);
    }

    public static ReadLength of(String length) {
        return new ReadLength(length);
    }

    public static ReadLength of(ByteBuf byteBuf) {
        return new ReadLength(byteBuf);
    }
}
