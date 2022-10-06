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
package org.apache.plc4x.java.s7.readwrite.utils;

import org.apache.plc4x.java.s7.readwrite.DeviceGroup;

public class S7TsapIdEncoder {

    private S7TsapIdEncoder() {
        // Prevent this from being instantiated.
    }

    public static short encodeS7TsapId(DeviceGroup deviceGroup, int rack, int slot) {
        short firstByte = (short) (deviceGroup.getValue() << 8);
        short secondByte = (short) ((rack << 4) | (slot & 0x0F));
        return (short) (firstByte | secondByte);
    }

    public static DeviceGroup decodeDeviceGroup(short tsapId) {
        byte deviceGroupCode = (byte) ((tsapId >> 8) & (0xFF));
        return DeviceGroup.enumForValue(deviceGroupCode);
    }

    public static int decodeRack(short tsapId) {
        return (tsapId >> 4) & 0xF;
    }

    public static int decodeSlot(short tsapId) {
        return tsapId & 0xF;
    }

}
