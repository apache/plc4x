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

public enum HeaderErrorClass {
    NO_ERROR((byte) 0x00),
    APPLICATION_RELATIONSHIP_ERROR((byte) 0x81),
    OBJECT_DEFINITION_ERROR((byte) 0x82),
    NO_RESOURCES_AVAILABLE_ERROR((byte) 0x83),
    ERROR_ON_SERVICE_PROCESSING((byte) 0x84),
    ERROR_ON_SUPPLIES((byte) 0x85),
    ACCESS_ERROR((byte) 0x87);

    private static final Map<Byte, HeaderErrorClass> map;
    static {
        map = new HashMap<>();
        for (HeaderErrorClass headerErrorClass : HeaderErrorClass.values()) {
            map.put(headerErrorClass.code, headerErrorClass);
        }
    }

    private final byte code;

    HeaderErrorClass(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    public static HeaderErrorClass valueOf(byte code) {
        return map.get(code);
    }

}
