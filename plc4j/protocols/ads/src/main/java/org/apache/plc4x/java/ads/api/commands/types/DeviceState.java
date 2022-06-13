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
import org.apache.plc4x.java.ads.api.util.UnsignedShortLEByteValue;

public class DeviceState extends UnsignedShortLEByteValue {

    public static final int NUM_BYTES = UnsignedShortLEByteValue.UNSIGNED_SHORT_LE_NUM_BYTES;

    private DeviceState(byte... values) {
        super(values);
    }

    private DeviceState(int value) {
        super(value);
    }

    private DeviceState(String value) {
        super(value);
    }

    private DeviceState(ByteBuf byteBuf) {
        super(byteBuf);
    }

    public static DeviceState of(byte... values) {
        return new DeviceState(values);
    }

    public static DeviceState of(int value) {
        return new DeviceState(value);
    }

    public static DeviceState of(ByteBuf byteBuf) {
        return new DeviceState(byteBuf);
    }

    public static DeviceState of(String value) {
        return new DeviceState(value);
    }
}
