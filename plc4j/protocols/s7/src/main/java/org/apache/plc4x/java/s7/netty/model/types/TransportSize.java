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
    BIT((byte) 0x01),
    BYTE((byte) 0x02),
    CHAR((byte) 0x03),
    WORD((byte) 0x04),
    INT((byte) 0x05),
    DWORD((byte) 0x06),
    DINT((byte) 0x07),
    REAL((byte) 0x08),
    TOD((byte) 0x0A),
    TIME((byte) 0x0B),
    S5TIME((byte) 0x0C),
    DATE_AND_TIME((byte) 0x0F),
    COUNTER((byte) 0x1C),
    TIMER((byte) 0x1D),
    IEC_TIMER((byte) 0x1E),
    IEC_COUNTER((byte) 0x1F),
    HS_COUNTER((byte) 0x20);

    private final byte code;

    TransportSize(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, TransportSize> map = null;

    public static TransportSize valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (TransportSize transportSize : TransportSize.values()) {
                map.put(transportSize.code, transportSize);
            }
        }
        return map.get(code);
    }

}
