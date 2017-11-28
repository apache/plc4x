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

public enum DisconnectReason {
    // Reasons available for all classes:
    REASON_NOT_SPECIFIED((byte) 0x00),
    CONGESTION_AT_TSAP((byte) 0x01),
    SESSION_ENTITY_NOT_ATACHED_TO_TSAP((byte) 0x02),
    ADDRESS_UNKNOWN((byte) 0x03),

    // Reasons only available for classes 1 to 4:
    NORMAL((byte) 0x80),
    REMOTE_TRANSPORT_ENTITY_CONGESTION((byte) 0x81),
    CONNECTION_NEGOTIATION_FAILED((byte) 0x82),
    DUPLICATE_SOURCE_REFERENCE((byte) 0x83),
    MISMATCHED_REFERENCES((byte) 0x84),
    PROTOCOL_ERROR((byte) 0x85),
    REFERENCE_OVERFLOW((byte) 0x87),
    CONNECTION_REQUEST_REFUSED((byte) 0x88),
    HEADER_OR_PARAMETER_LENGTH_INVALID((byte) 0x8A);

    private byte code;

    DisconnectReason(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    private static Map<Byte, DisconnectReason> map = null;

    public static DisconnectReason valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (DisconnectReason disconnectReason : DisconnectReason.values()) {
                map.put(disconnectReason.code, disconnectReason);
            }
        }
        return map.get(code);
    }

}
