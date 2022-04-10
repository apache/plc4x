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
import org.apache.plc4x.java.ads.api.util.LengthSupplier;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

public class WriteLength extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    private WriteLength(byte... values) {
        super(values);
    }

    private WriteLength(long value) {
        super(value);
    }

    private WriteLength(String length) {
        super(length);
    }

    private WriteLength(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static WriteLength of(byte... values) {
        return new WriteLength(values);
    }

    public static WriteLength of(long value) {
        return new WriteLength(value);
    }

    public static WriteLength of(LengthSupplier lengthSupplier) {
        return new WriteLength(lengthSupplier.getCalculatedLength());
    }

    public static WriteLength of(String length) {
        return new WriteLength(length);
    }

    public static WriteLength of(ByteBuf byteBuf) {
        return new WriteLength(byteBuf);
    }
}
