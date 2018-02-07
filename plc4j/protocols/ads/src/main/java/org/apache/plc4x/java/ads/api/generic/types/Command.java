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
import java.util.Arrays;

public enum Command implements ByteReadable {
    Invalid(0x00),
    ADS_Read_Device_Info(0x01),
    ADS_Read(0x02),
    ADS_Write(0x03),
    ADS_Read_State(0x04),
    ADS_Write_Control(0x05),
    ADS_Add_Device_Notification(0x06),
    ADS_Delete_Device_Notification(0x07),
    ADS_Device_Notification(0x08),
    ADS_Read_Write(0x09),
    /**
     * Other commands are not defined or are used internally. Therefore the Command Id  is only allowed to contain the above enumerated values!
     */
    UNKNOWN();

    public static final int NUM_BYTES = 2;

    final byte[] value;

    final int intValue;

    Command() {
        // Only used for unkown enum
        value = new byte[0];
        intValue = 0;
    }

    Command(int value) {
        ByteValue.checkUnsignedBounds(value, NUM_BYTES);
        this.intValue = value;
        this.value = ByteBuffer.allocate(NUM_BYTES)
            // LE
            .put((byte) (value & 0xff))
            .put((byte) (value >> 8 & 0xff))
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

    public static Command of(byte... bytes) {
        // TODO: improve by using a map
        for (Command command : values()) {
            if (Arrays.equals(bytes, command.value)) {
                return command;
            }
        }
        return UNKNOWN;
    }

    public static Command of(int intValue) {
        // TODO: improve by using a map
        for (Command command : values()) {
            if (command.intValue == intValue) {
                return command;
            }
        }
        return UNKNOWN;
    }

    public static Command of(ByteBuf byteBuf) {
        return of(byteBuf.readUnsignedShortLE());
    }
}
