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
package org.apache.plc4x.java.ads.api.generic.types;

import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import org.apache.plc4x.java.ads.api.util.ByteReadable;
import org.apache.plc4x.java.ads.api.util.ByteValue;

import java.nio.ByteBuffer;

public enum Command implements ByteReadable {
    Invalid(0x0000),
    ADS_Read_Device_Info(0x0001),
    ADS_Read(0x0002),
    ADS_Write(0x0003),
    ADS_Read_State(0x0004),
    ADS_Write_Control(0x0005),
    ADS_Add_Device_Notification(0x0006),
    ADS_Delete_Device_Notification(0x0007),
    ADS_Device_Notification(0x0008),
    ADS_Read_Write(0x0009),
    /**
     * Other commands are not defined or are used internally. Therefore the Command Id  is only allowed to contain the above enumerated values!
     */
    UNKNOWN();

    public static final int NUM_BYTES = 4;

    final byte[] value;

    Command() {
        // Only used for unkown enum
        value = new byte[0];
    }

    Command(long value) {
        ByteValue.checkUnsignedBounds(value, NUM_BYTES);
        this.value = ByteBuffer.allocate(NUM_BYTES)
            .put((byte) (value >> 24 & 0xff))
            .put((byte) (value >> 16 & 0xff))
            .put((byte) (value >> 8 & 0xff))
            .put((byte) (value & 0xff))
            .array();
    }

    @Override
    public byte[] getBytes() {
        if (this == UNKNOWN) {
            throw new IllegalStateException("Unknown enum can't be serialized");
        }
        return value;
    }

    public ByteBuf getByteBuf() {
        if (this == UNKNOWN) {
            throw new IllegalStateException("Unknown enum can't be serialized");
        }
        return Unpooled.buffer().writeBytes(value);
    }
}
