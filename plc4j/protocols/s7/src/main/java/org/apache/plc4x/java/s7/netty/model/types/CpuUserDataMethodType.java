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
import org.apache.plc4x.java.api.model.PlcField;

/**
 *
 * @author cgarcia
 */
public enum CpuUserDataMethodType implements PlcField {
    
    UNDEF("UNDEF", (byte) 0x00),
    REQUEST("REQUEST", (byte) 0x11),
    RESPONSE("RESPONSE", (byte) 0x12),
    ;
    
    
    private static final Map<Byte, CpuUserDataMethodType> map;
    
    static {
        map = new HashMap<>();
        for (CpuUserDataMethodType subevent : CpuUserDataMethodType.values()) {
            map.put(subevent.code, subevent);
        }
    }    
    
    private final String event;
    private final byte code;
    
    CpuUserDataMethodType(String event, byte code){
        this.event = event;
        this.code = code;
    }
    
    public String getEvent(){
        return event;
    }    
    
    public byte getCode() {
        return code;
    }    
    
    public static CpuUserDataMethodType valueOfEvent(String event) {
        for (CpuUserDataMethodType value : CpuUserDataMethodType.values()) {
            if(value.getEvent().equals(event)) {
                return value;
            }
        }
        return null;
    }

    public static CpuUserDataMethodType valueOf(byte code) {
        return map.get(code);
    }
    
    
}
