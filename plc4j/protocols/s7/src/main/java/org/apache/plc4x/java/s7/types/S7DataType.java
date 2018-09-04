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
package org.apache.plc4x.java.s7.types;

import org.apache.plc4x.java.s7.netty.model.types.DataTransportSize;

import java.util.*;

public enum S7DataType {

    /**
     * TODO: For the types with code 0x00 we need to put some additional effort in reverse engineering the codes for these types.
     */
    // -----------------------------------------
    // Single bit
    // -----------------------------------------
    BOOL(0x01, "X", 1, null, DataTransportSize.BIT, S7ControllerType.S7_ANY),

    // -----------------------------------------
    // Bit strings
    // -----------------------------------------
    BYTE(0x02, "B", 1, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    WORD(0x04, "W", 2, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    DWORD(0x06, "D", 4, WORD, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    LWORD(0x00, null, 8, null, null, S7ControllerType.S7_1200, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Integers
    // -----------------------------------------
    INT(0x05, "W", 2, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    // Double Precision Int
    DINT(0x07, "D", 4, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    // Unsigned Small Int
    USINT(0x00, "B", 1, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // (Signed) Small Int
    SINT(0x00, "B", 1, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Unsigned Int
    UINT(0x00, "W", 2, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Unsigned Double Precision Int
    UDINT(0x00, "D", 4, INT, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    LINT(0x00, null, 8, INT, null, S7ControllerType.S7_1500),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    ULINT(0x00, null, 16, INT, null, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Reals
    // -----------------------------------------
    REAL(0x08, "D", 4, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    LREAL(0x00, "X", 8, REAL, null, S7ControllerType.S7_1200, S7ControllerType.S7_1500),

    // -----------------------------------------
    // Durations
    // -----------------------------------------
    // IEC time
    TIME(0x0B, "", 4, null, null, S7ControllerType.S7_ANY),

    // -----------------------------------------
    // Date
    // -----------------------------------------
    // IEC date (yyyy-m-d)
    DATE(0x00, "", 2, null, null, S7ControllerType.S7_ANY),

    // -----------------------------------------
    // Time of day
    // -----------------------------------------
    // Time (hh:mm:ss.S)
    TIME_OF_DAY(0x0A, "", 4, null, null, S7ControllerType.S7_ANY),

    // -----------------------------------------
    // Date and time of day
    // -----------------------------------------
    DATE_AND_TIME(0x0F, "", 8, null, null, S7ControllerType.S7_ANY),

    // -----------------------------------------
    // ASCII Strings
    // -----------------------------------------
    // Single-byte character
    CHAR(0x03, "B", 1, null, DataTransportSize.BYTE_WORD_DWORD, S7ControllerType.S7_ANY),
    // Double-byte character
    WCHAR(0x00, "", 2, null, null, S7ControllerType.S7_ANY),
    // Variable-length single-byte character string
    STRING(0x00, "", -1, null, null, S7ControllerType.S7_ANY),
    // Variable-length double-byte character string
    WSTRING(0x00, "", -1, null, null, S7ControllerType.S7_ANY);

    /* TO BE CONTINUED */

    private final byte typeCode;
    private final String sizeCode;
    private int sizeInBytes;
    private final Set<S7ControllerType> supportedControllerTypes;
    private final S7DataType baseType;
    private final DataTransportSize dataTransportSize;

    S7DataType(int typeCode, String sizeCode, int sizeInBytes, S7DataType baseType, DataTransportSize dataTransportSize,
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

    S7DataType getBaseType() {
        // If this is a base-type itself, the baseType is null, in all
        // other cases it is set.
        if (baseType == null) {
            return this;
        } else {
            return baseType;
        }
    }

    S7DataType getSubType(String sizeCode) {
        // Try to find a sub-type with this base type for which the size code matches.
        for (S7DataType value : values()) {
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

    private final static Map<Byte, S7DataType> map;

    static {
        map = new HashMap<>();
        for (S7DataType dataType : S7DataType.values()) {
            map.put(dataType.typeCode, dataType);
        }
    }

    public static S7DataType valueOf(byte code) {
        return map.get(code);
    }

}