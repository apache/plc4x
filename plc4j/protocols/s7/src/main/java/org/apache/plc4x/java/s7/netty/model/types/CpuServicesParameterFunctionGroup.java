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

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public enum CpuServicesParameterFunctionGroup {

    CPU_FUNCTIONS((byte) 0x04);

    private static final Logger logger = LoggerFactory.getLogger(CpuServicesParameterFunctionGroup.class);

    private static final Map<Byte, CpuServicesParameterFunctionGroup> map;
    static {
        map = new HashMap<>();
        for (CpuServicesParameterFunctionGroup cpuServicesParameterFunctionGroup : CpuServicesParameterFunctionGroup.values()) {
            map.put(cpuServicesParameterFunctionGroup.code, cpuServicesParameterFunctionGroup);
        }
    }

    private final byte code;

    CpuServicesParameterFunctionGroup(byte code) {
        this.code = code;
    }

    public byte getCode() {
        return code;
    }

    public static CpuServicesParameterFunctionGroup valueOf(byte code) {
        if(!map.containsKey(code)) {
            logger.error("CpuServicesParameterFunctionGroup for code {} not found", code);
        }
        return map.get(code);
    }

}
