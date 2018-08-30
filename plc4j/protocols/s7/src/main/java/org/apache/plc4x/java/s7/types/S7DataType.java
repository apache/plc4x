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

import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

public enum S7DataType {

    /**
     * TODO: For the types with code 0x00 we need to put some additional effort in reverse engineering the codes for these types.
     */
    BOOL(0x01, "X", null, S7ControllerType.S7_ANY),

    BYTE(0x02, "B", null, S7ControllerType.S7_ANY),

    CHAR(0x03, "B", null, S7ControllerType.S7_ANY),

    WORD(0x04, "W", null, S7ControllerType.S7_ANY),
    DWORD(0x06, "D", WORD, S7ControllerType.S7_ANY),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    LWORD(0x00, null, null, S7ControllerType.S7_1200, S7ControllerType.S7_1500),

    INT(0x05, "W", null, S7ControllerType.S7_ANY),
    // Double Precision Int
    DINT(0x07, "D", INT, S7ControllerType.S7_ANY),
    // Unsigned Small Int
    USINT(0x00, "B", INT, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // (Signed) Small Int
    SINT(0x00, "B", INT, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Unsigned Int
    UINT(0x00, "W", INT, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Unsigned Double Precision Int
    UDINT(0x00, "D", INT, S7ControllerType.S7_1200, S7ControllerType.S7_1500),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    UDLINT(0x00, null, INT, S7ControllerType.S7_1500),
    // Only got a basic TIA license (S7-1500 needed to find this out)
    LINT(0x00, null, INT, S7ControllerType.S7_1500),

    REAL(0x08, "D", null, S7ControllerType.S7_ANY),
    // Ok ... this is strange ...
    LREAL(0x00, "X", REAL, S7ControllerType.S7_1200, S7ControllerType.S7_1500);

    /* TO BE CONTINUED */

    private byte typeCode;
    private String sizeCode;
    private Set<S7ControllerType> supportedControllerTypes;
    private S7DataType baseType;

    S7DataType(int typeCode, String sizeCode, S7DataType baseType, S7ControllerType... supportedControllerTypes) {
        this.typeCode = (byte) typeCode;
        this.sizeCode = sizeCode;
        this.supportedControllerTypes = new HashSet<>(Arrays.asList(supportedControllerTypes));
        this.baseType = baseType;
    }

    byte getTypeCode() {
        return typeCode;
    }

    public String getSizeCode() {
        return sizeCode;
    }

    boolean isBaseType() {
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

    boolean isControllerTypeSupported(S7ControllerType controllerType) {
        return supportedControllerTypes.contains(controllerType);
    }

}