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
public enum DataTransportSize {
    NULL((byte) 0x00, false),
    BIT((byte) 0x03, true),
    BYTE_WORD_DWORD((byte) 0x04, true),
    INTEGER((byte) 0x05, true),
    DINTEGER((byte) 0x06, false),
    REAL((byte) 0x07, false),
    OCTET_STRING((byte) 0x09, false);

    private static final Map<Byte, DataTransportSize> map;
    static {
        map = new HashMap<>();
        for (DataTransportSize dataTransportSize : DataTransportSize.values()) {
            map.put(dataTransportSize.code, dataTransportSize);
        }
    }

    private final byte code;
    private final boolean sizeInBits;

    DataTransportSize(byte code, boolean sizeInBits) {
        this.code = code;
        this.sizeInBits = sizeInBits;
    }

    public byte getCode() {
        return code;
    }

    public boolean isSizeInBits() {
        return sizeInBits;
    }

    public static DataTransportSize valueOf(byte code) {
        return map.get(code);
    }

}
