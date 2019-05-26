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

import org.apache.plc4x.java.s7.types.S7ControllerType;

import java.util.*;

public enum TransportSize {

    /**
     * TODO: For the types with code 0x00 we need to put some additional effort in reverse engineering the codes for these types.
     */
    // -----------------------------------------
    // Single bit
    // -----------------------------------------
    BOOL(0x01, "X", 1, null, DataTransportSize.BIT, S7ControllerType.ANY),

    // -----------------------------------------
    // Bit strings
    // -----------------------------------------
    BYTE(0x02, "B", 1, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    WORD(0x04, "W", 2, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    DWORD(0x06, "D", 4, WORD, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    // TODO: Find the code
    LWORD(0x00, "X", 8, null, null, S7ControllerType.S7_1200, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Integers
    // -----------------------------------------
    // Signed Int
    INT(0x05, "W", 2, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    // Unsigned Int
    UINT(0x05, "W", 2, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // (Signed) Small Int
    SINT(0x02, "B", 1, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Unsigned Small Int
    USINT(0x02, "B", 1, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Double Precision Int
    DINT(0x07, "D", 4, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    // Unsigned Double Precision Int
    UDINT(0x07, "D", 4, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    // TODO: Find the code
    LINT(0x00, "X", 8, INT, null, S7ControllerType.S7_1500),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    // TODO: Find the code
    ULINT(0x00, "X", 16, INT, null, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Reals
    // -----------------------------------------
    REAL(0x08, "D", 4, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    // TODO: Find the code
    LREAL(0x00, "X", 8, REAL, null, S7ControllerType.S7_1200, S7ControllerType.S7_1200, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Durations
    // -----------------------------------------
    // IEC time
    TIME(0x0B, "X", 4, null, null, S7ControllerType.ANY),
    // TODO: Find the code
    LTIME(0x00, "X", 8, TIME, null, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Date
    // -----------------------------------------
    // IEC date (yyyy-m-d)
    // TODO: Find the code
    DATE(0x00, "X", 2, null, null, S7ControllerType.ANY),

    // -----------------------------------------
    // Time of day
    // -----------------------------------------
    // Time (hh:mm:ss.S)
    TIME_OF_DAY(0x0A, "X", 4, null, null, S7ControllerType.ANY),

    // -----------------------------------------
    // Date and time of day
    // -----------------------------------------
    DATE_AND_TIME(0x0F, "X", 8, null, null, S7ControllerType.ANY),

    // -----------------------------------------
    // ASCII Strings
    // -----------------------------------------
    // Single-byte character
    CHAR(0x03, "B", 1, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    // Double-byte character
    // TODO: Find the code (Perhaps 0x13)
    WCHAR(0x13, "X", 2, null, null, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Variable-length single-byte character string
    STRING(0x03, "X", 1, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.ANY),
    // Variable-length double-byte character string
    // TODO: Find the code (Perhaps 0x13)
    WSTRING(0x00, "X", 1, null, null, S7ControllerType.S7_1200, S7ControllerType.S7_1500);

    /* TO BE CONTINUED */

    // Codes and their types:
    // 0x1C: Counter
    // 0x1D: Timer
    // 0x1E: IEC Timer
    // 0x1F: IEC Counter
    // 0x20: HS Counter
    //

    private static final Map<Byte, TransportSize> map;
    static {
        map = new HashMap<>();
        for (TransportSize dataType : TransportSize.values()) {
            map.put(dataType.typeCode, dataType);
        }
    }

    private final byte typeCode;
    private final String sizeCode;
    private final int sizeInBytes;
    private final Set<S7ControllerType> supportedControllerTypes;
    private final TransportSize baseType;
    private final DataTransportSize dataTransportSize;

    TransportSize(int typeCode, String sizeCode, int sizeInBytes, TransportSize baseType, DataTransportSize dataTransportSize,
                  S7ControllerType... supportedControllerTypes) {
        this.typeCode = (byte) typeCode;
        this.sizeCode = sizeCode;
        this.sizeInBytes = sizeInBytes;
        this.supportedControllerTypes = new HashSet<>(Arrays.asList(supportedControllerTypes));
        this.baseType = baseType;
        this.dataTransportSize = dataTransportSize;
    }

    public byte getTypeCode() {
        return typeCode;
    }

    public String getSizeCode() {
        return sizeCode;
    }

    public int getSizeInBytes() {
        return sizeInBytes;
    }

    public boolean isBaseType() {
        return baseType == null;
    }

    TransportSize getBaseType() {
        // If this is a base-type itself, the baseType is null, in all
        // other cases it is set.
        if (baseType == null) {
            return this;
        } else {
            return baseType;
        }
    }

    TransportSize getSubType(String sizeCode) {
        // Try to find a sub-type with this base type for which the size code matches.
        for (TransportSize value : values()) {
            if ((value.baseType == this) && (value.sizeCode != null) && (value.sizeCode.equals(sizeCode))) {
                return value;
            }
        }
        return null;
    }

    public DataTransportSize getDataTransportSize() {
        return dataTransportSize;
    }

    boolean isControllerTypeSupported(S7ControllerType controllerType) {
        return supportedControllerTypes.contains(controllerType);
    }

    public static TransportSize valueOf(byte code) {
        return map.get(code);
    }

}