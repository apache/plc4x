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
package org.apache.plc4x.sandbox.java.dynamic.s7.types;

import java.util.HashMap;
import java.util.Map;

/**
 * (Values determined by evaluating generated ".pcap" files)
 */
public enum MemoryArea {
    COUNTERS("C", (byte) 0x1C), /* Renamed from "S7 Counters" */ // TODO: Double check shortName
    TIMERS("T", (byte) 0x1D), /* Renamed from "S7 Timers" */ // TODO: Double check shortName
    DIRECT_PERIPHERAL_ACCESS("D", (byte) 0x80), // TODO: Double check shortName
    INPUTS("I", (byte) 0x81),
    OUTPUTS("Q", (byte) 0x82),
    FLAGS("F", (byte) 0x83), // TODO: Double check shortName
    DATA_BLOCKS("DB", (byte) 0x84),
    INSTANCE_DATA_BLOCKS("DBI", (byte) 0x85), // TODO: Double check shortName
    LOCAL_DATA("LD", (byte) 0x86); // TODO: Double check shortName

    // TODO: I think we should remove these ...
    /*S7_200_SYSTEM_INFO(null, (byte) 0x03), * Renamed from "System info of 200 family" *
    S7_200_FLAGS(null, (byte) 0x05), * Renamed from "System flags of 200 family" *
    S7_200_INPUTS(null, (byte) 0x06), * Renamed from "System inputs of 200 family" *
    S7_200_OUTPUTS(null, (byte) 0x07), * Renamed from "System outputs of 200 family" *
    S7_200_IEC_COUNTERS(null, (byte) 0x1E), * Renamed from "IEC counters (200 family)" *
    S7_200_IEC_TIMERS(null, (byte) 0x1F); * Renamed from "IEC timers (200 family)" */

    private static final Map<Byte, MemoryArea> map;
    static {
        map = new HashMap<>();
        for (MemoryArea memoryArea : MemoryArea.values()) {
            map.put(memoryArea.code, memoryArea);
        }
    }

    private final String shortName;
    private final byte code;

    MemoryArea(String shortName, byte code) {
        this.shortName = shortName;
        this.code = code;
    }

    public String getShortName() {
        return shortName;
    }

    public byte getCode() {
        return code;
    }

    public static MemoryArea valueOfShortName(String shortName) {
        for (MemoryArea value : MemoryArea.values()) {
            if(value.getShortName().equals(shortName)) {
                return value;
            }
        }
        return null;
    }

    public static MemoryArea valueOf(byte code) {
        return map.get(code);
    }

}
