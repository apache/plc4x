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
public enum ParameterType {
    CPU_SERVICES((byte) 0x00),
    READ_VAR((byte) 0x04),
    WRITE_VAR((byte) 0x05),
    REQUEST_DOWNLOAD((byte) 0x1A),
    DOWNLOAD_BLOCK((byte) 0x1B),
    DOWNLOAD_ENDED((byte) 0x1C),
    START_UPLOAD((byte) 0x1D),
    UPLOAD((byte) 0x1E),
    END_UPLOAD((byte) 0x1F),
    PI_SERVICE((byte) 0x28),
    PLC_STOP((byte) 0x29),
    SETUP_COMMUNICATION((byte) 0xF0);

    private final byte code;

    ParameterType(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, ParameterType> map = null;

    public static ParameterType valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (ParameterType parameterType : ParameterType.values()) {
                map.put(parameterType.code, parameterType);
            }
        }
        if(!map.containsKey(code)) {
            System.out.println("ParameterType for code " + code + " not found");
        }
        return map.get(code);
    }

}
