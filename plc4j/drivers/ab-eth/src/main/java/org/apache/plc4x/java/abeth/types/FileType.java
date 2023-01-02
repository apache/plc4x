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
package org.apache.plc4x.java.abeth.types;

import org.apache.plc4x.java.api.types.PlcValueType;

import java.util.HashMap;
import java.util.Map;

public enum FileType {

    STATUS((short) 0x84, PlcValueType.RAW_BYTE_ARRAY),
    BIT((short) 0x85, PlcValueType.BOOL),
    TIMER((short) 0x86, PlcValueType.TIME),
    COUNTER((short) 0x87, PlcValueType.RAW_BYTE_ARRAY),
    CONTROL((short) 0x88, PlcValueType.RAW_BYTE_ARRAY),
    INTEGER((short) 0x89, PlcValueType.INT),
    FLOAT((short) 0x8A, PlcValueType.REAL),
    OUTPUT((short) 0x8B, PlcValueType.RAW_BYTE_ARRAY),
    INPUT((short) 0x8C, PlcValueType.RAW_BYTE_ARRAY),
    STRING((short) 0x8D, PlcValueType.STRING),
    ASCII((short) 0x8E, PlcValueType.STRING),
    BCD((short) 0x8F, PlcValueType.RAW_BYTE_ARRAY),
    WORD((short) 0x89, PlcValueType.WORD),  // custom 2-byte Integer
    DWORD((short) 0x89, PlcValueType.DWORD), // custom 4-byte Integer
    SINGLEBIT((short) 0x89, PlcValueType.WORD); // custom single bit from Integer file, no official AB type

    private final short typeCode;
    private final PlcValueType plcValueType;

    FileType(short typeCode, PlcValueType plcValueType) {
        this.typeCode = typeCode;
        this.plcValueType = plcValueType;
    }

    public short getTypeCode() {
        return typeCode;
    }

    public PlcValueType getPlcValueType() {
        return plcValueType;
    }

    private static final Map<Short, FileType> map;
    static {
        map = new HashMap<>();
        for (FileType dataType : FileType.values()) {
            map.put(dataType.typeCode, dataType);
        }
    }

    public static FileType valueOf(short code) {
        return map.get(code);
    }

}
