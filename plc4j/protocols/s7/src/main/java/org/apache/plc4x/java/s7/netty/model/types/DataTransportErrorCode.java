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

public enum DataTransportErrorCode {
    RESERVED((byte) 0x00),
    ACCESS_DENIED((byte) 0x03),
    INVALID_ADDRESS((byte) 0x05),
    NOT_FOUND((byte) 0x0A),
    OK((byte) 0xFF);

    private byte code;

    DataTransportErrorCode(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, DataTransportErrorCode> map = null;

    public static DataTransportErrorCode valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (DataTransportErrorCode dataTransportErrorCode : DataTransportErrorCode.values()) {
                map.put(dataTransportErrorCode.code, dataTransportErrorCode);
            }
        }
        if(!map.containsKey(code)) {
            System.out.println("No DataTransportErrorCode for code " + code);
        }
        return map.get(code);
    }

}
