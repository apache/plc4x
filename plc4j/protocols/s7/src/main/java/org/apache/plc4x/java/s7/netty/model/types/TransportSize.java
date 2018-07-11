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
public enum TransportSize {
    BIT((byte) 0x01, 1),
    BYTE((byte) 0x02, 1),
    CHAR((byte) 0x03, 1),
    WORD((byte) 0x04, 2),
    INT((byte) 0x05, 2),
    DWORD((byte) 0x06, 4),
    DINT((byte) 0x07, 4),
    REAL((byte) 0x08, 4),
    TOD((byte) 0x0A, 4),
    TIME((byte) 0x0B, 4),
    S5TIME((byte) 0x0C, 2),
    DATE_AND_TIME((byte) 0x0F, 4),
    COUNTER((byte) 0x1C, -1),
    TIMER((byte) 0x1D, -1),
    IEC_TIMER((byte) 0x1E, -1),
    IEC_COUNTER((byte) 0x1F, -1),
    HS_COUNTER((byte) 0x20, -1);

    private final byte code;
    private final int sizeInBytes;

    TransportSize(byte code, int sizeInBytes) {
        this.code = code;
        this.sizeInBytes = sizeInBytes;
    }

    public byte getCode() {
        return code;
    }

    public int getSizeInBytes() {
        return sizeInBytes;
    }

    private final static Map<Byte, TransportSize> map;

    static {
        map = new HashMap<>();
        for (TransportSize transportSize : TransportSize.values()) {
            map.put(transportSize.code, transportSize);
        }
    }

    public static TransportSize valueOf(byte code) {
        return map.get(code);
    }

}
