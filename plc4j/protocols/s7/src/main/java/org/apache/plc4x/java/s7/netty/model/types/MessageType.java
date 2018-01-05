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
public enum MessageType {
    JOB((byte) 0x01),
    ACK((byte) 0x02),
    ACK_DATA((byte) 0x03),
    USER_DATA((byte) 0x07); /* Renamed from "Userdata" */

    private static Map<Byte, MessageType> map = null;
    
    private final byte code;

    MessageType(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    public static MessageType valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (MessageType messageType : MessageType.values()) {
                map.put(messageType.code, messageType);
            }
        }
        return map.get(code);
    }

}
