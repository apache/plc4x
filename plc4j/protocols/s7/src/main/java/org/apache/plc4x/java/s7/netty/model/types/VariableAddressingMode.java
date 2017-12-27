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
 * Renamed from "SyntaxId".
 */
public enum VariableAddressingMode {
    S7ANY((byte) 0x10),
    PBC_R_ID((byte) 0x13),
    ALARM_LOCKFREE((byte) 0x15),
    ALARM_IND((byte) 0x16),
    ALARM_ACK((byte) 0x19),
    ALARM_QUERYREQ((byte) 0x1a),
    NOTIFY_IND((byte) 0x1c),
    NIC((byte) 0x82),
    DRIVEESANY((byte) 0xa2),
    DBREAD((byte) 0xb0),
    SYM1200((byte) 0xb2); /* Renamed from "1200SYM" */

    private final byte code;

    VariableAddressingMode(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, VariableAddressingMode> map = null;

    public static VariableAddressingMode valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (VariableAddressingMode variableAddressingMode : VariableAddressingMode.values()) {
                map.put(variableAddressingMode.code, variableAddressingMode);
            }
        }
        return map.get(code);
    }

}
