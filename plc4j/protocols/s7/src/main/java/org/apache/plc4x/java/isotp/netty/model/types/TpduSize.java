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

public enum TpduSize {
    SIZE_128((byte) 0x07, 128),
    SIZE_256((byte) 0x08, 256),
    SIZE_512((byte) 0x09, 512),
    SIZE_1024((byte) 0x0a, 1024),
    SIZE_2048((byte) 0x0b, 2048),
    SIZE_4096((byte) 0x0c, 4096),
    SIZE_8192((byte) 0x0d, 8192);

    private static Map<Byte, TpduSize> map = null;
    
    private final byte code;
    private final int value;

    TpduSize(byte code, int value) {
        this.code = code;
        this.value = value;
    }

    public byte getCode() {
        return code;
    }

    public int getValue() {
        return value;
    }

    public static TpduSize valueForGivenSize(int pduSize) {
        assert pduSize > 0;
        for (TpduSize tpduSize : values()) {
            if(tpduSize.getValue() <= pduSize) {
                return tpduSize;
            }
        }
        // If the requested pdu size is greater than 8MB,
        // Simply use that as the given size is simple a
        // requested size, if the remote responds with a
        // lower value the application has to live with
        // this anyway.
        return SIZE_8192;
    }

    public static TpduSize valueOf(byte code) {
        if (map == null) {
            map = new HashMap<>();
            for (TpduSize tpduSize : TpduSize.values()) {
                map.put(tpduSize.code, tpduSize);
            }
        }
        return map.get(code);
    }

}
