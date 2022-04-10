/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.abeth.types;

import java.util.HashMap;
import java.util.Map;

public enum FileType {

    HURZ((short) 0x64),
    STATUS((short) 0x84),
    BIT((short) 0x85),
    TIMER((short) 0x86),
    COUNTER((short) 0x87),
    CONTROL((short) 0x88),
    INTEGER((short) 0x89),
    FLOAT((short) 0x8A),
    OUTPUT((short) 0x8B),
    INPUT((short) 0x8C),
    STRING((short) 0x8D),
    ASCII((short) 0x8E),
    BCD((short) 0x8F),
    WORD((short) 0x89),  // custom 2-byte Integer
    DWORD((short) 0x89), // custom 4-byte Integer
    SINGLEBIT((short) 0x89); // custom single bit from Integer file, no official AB type

    private final short typeCode;

    FileType(short typeCode) {
        this.typeCode = typeCode;
    }

    public short getTypeCode() {
        return typeCode;
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
