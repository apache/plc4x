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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public enum DataTransportErrorCode {
    RESERVED((byte) 0x00),
    OK((byte) 0xFF),
    ACCESS_DENIED((byte) 0x03),
    INVALID_ADDRESS((byte) 0x05),
    DATA_TYPE_NOT_SUPPORTED((byte) 0x06),
    NOT_FOUND((byte) 0x0A);

    private static final Logger logger = LoggerFactory.getLogger(DataTransportErrorCode.class);

    private static final Map<Byte, DataTransportErrorCode> map;
    static {
        map = new HashMap<>();
        for (DataTransportErrorCode dataTransportErrorCode : DataTransportErrorCode.values()) {
            map.put(dataTransportErrorCode.code, dataTransportErrorCode);
        }
    }

    private byte code;

    DataTransportErrorCode(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    public static DataTransportErrorCode valueOf(byte code) {
        if (!map.containsKey(code)) {
            logger.error("No DataTransportErrorCode for code {}", code);
        }
        return map.get(code);
    }

}
