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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import org.apache.plc4x.java.ads.api.util.LengthSupplier;
import org.apache.plc4x.java.ads.api.util.UnsignedIntLEByteValue;

public class DataLength extends UnsignedIntLEByteValue {

    public static final int NUM_BYTES = UnsignedIntLEByteValue.UNSIGNED_INT_LE_NUM_BYTES;

    private DataLength(byte... values) {
        super(values);
    }

    private DataLength(long value) {
        super(value);
    }

    private DataLength(String length) {
        super(length);
    }

    private DataLength(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static DataLength of(byte... values) {
        return new DataLength(values);
    }

    public static DataLength of(long value) {
        return new DataLength(value);
    }

    public static DataLength of(LengthSupplier lengthSupplier) {
        return new DataLength(lengthSupplier.getCalculatedLength());
    }

    public static DataLength of(String length) {
        return new DataLength(length);
    }

    public static DataLength of(ByteBuf byteBuf) {
        return new DataLength(byteBuf);
    }
}
