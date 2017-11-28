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
package org.apache.plc4x.java.isotp.netty.model.types;

import java.util.HashMap;
import java.util.Map;

public enum TpduCode {
    DATA((byte) 0xF0),
    DATA_ACKNOWLEDGEMENT((byte) 0x60),
    CONNECTION_REQUEST((byte) 0xE0),
    CONNECTION_CONFIRM((byte) 0xD0),
    DISCONNECT_REQUEST((byte) 0x80),
    DISCONNECT_CONFIRM((byte) 0xC0),
    EXPEDITED_DATA((byte) 0x10),
    EXPEDITED_DATA_ACKNOWLEDGEMENT((byte) 0x20),
    REJECT((byte) 0x50),
    TPDU_ERROR((byte) 0x70);

    private byte code;

    TpduCode(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, TpduCode> map = null;

    public static TpduCode valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (TpduCode tpduCode : TpduCode.values()) {
                map.put(tpduCode.code, tpduCode);
            }
        }
        return map.get(code);
    }

}
