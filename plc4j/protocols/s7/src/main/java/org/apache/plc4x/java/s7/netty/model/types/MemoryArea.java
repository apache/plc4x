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
package org.apache.plc4x.java.s7.netty.model.types;

import java.util.HashMap;
import java.util.Map;

/**
 * (Values determined by evaluating generated ".pcapng" files)
 */
public enum MemoryArea {
    COUNTERS((byte) 0x1C), /* Renamed from "S7 Counters" */
    TIMERS((byte) 0x1D), /* Renamed from "S7 Timers" */
    DIRECT_PERIPHERAL_ACCESS((byte) 0x80),
    INPUTS((byte) 0x81),
    OUTPUTS((byte) 0x82),
    FLAGS((byte) 0x83),
    DATA_BLOCKS((byte) 0x84),
    INSTANCE_DATA_BLOCKS((byte) 0x85),
    LOCAL_DATA((byte) 0x86),

    S7_200_SYSTEM_INFO((byte) 0x03), /* Renamed from "System info of 200 family" */
    S7_200_FLAGS((byte) 0x05), /* Renamed from "System flags of 200 family" */
    S7_200_INPUTS((byte) 0x06), /* Renamed from "System inputs of 200 family" */
    S7_200_OUTPUTS((byte) 0x07), /* Renamed from "System outputs of 200 family" */
    S7_200_IEC_COUNTERS((byte) 0x1E), /* Renamed from "IEC counters (200 family)" */
    S7_200_IEC_TIMERS((byte) 0x1F); /* Renamed from "IEC timers (200 family)" */

    private byte code;

    MemoryArea(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, MemoryArea> map = null;

    public static MemoryArea valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (MemoryArea memoryArea : MemoryArea.values()) {
                map.put(memoryArea.code, memoryArea);
            }
        }
        return map.get(code);
    }

}
