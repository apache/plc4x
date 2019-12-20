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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public enum CpuServicesParameterSubFunctionGroup {

    READ_SSL((byte) 0x01),
    MESSAGE_SERVICE((byte) 0x02),
    DIAG_MESSAGE((byte) 0x03),
    ALARM_S((byte) 0x04),
    ALARM8((byte) 0x05),
    NOTIFY((byte) 0x06),
    ALARM8_LOCK((byte) 0x07),
    ALARM8_UNLOCK((byte) 0x08),
    SCAN((byte) 0x04),
    ALARM_ACK((byte) 0x0b),
    ALARM_ACK_IND((byte) 0x0c),    
    ALARM8_LOCK_IND((byte) 0x0d),
    ALARM8_UNLOCK_IND((byte) 0x0e),
    ALARM_SQ_IND((byte) 0x11),  
    ALARM_S_IND((byte) 0x12),    
    ALARM_QUERY((byte) 0x13),    
    NOTIFY8((byte) 0x16); //TODO

    private static final Logger logger = LoggerFactory.getLogger(CpuServicesParameterSubFunctionGroup.class);

    private static final Map<Byte, CpuServicesParameterSubFunctionGroup> map;
    static {
        map = new HashMap<>();
        for (CpuServicesParameterSubFunctionGroup cpuServicesParameterSubFunctionGroup : CpuServicesParameterSubFunctionGroup.values()) {
            map.put(cpuServicesParameterSubFunctionGroup.code, cpuServicesParameterSubFunctionGroup);
        }
    }

    private final byte code;

    CpuServicesParameterSubFunctionGroup(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    public static CpuServicesParameterSubFunctionGroup valueOf(byte code) {
        if(!map.containsKey(code)) {
            logger.error("CpuServicesParameterSubFunctionGroup for code {} not found", code);
        }
        return map.get(code);
    }

}
