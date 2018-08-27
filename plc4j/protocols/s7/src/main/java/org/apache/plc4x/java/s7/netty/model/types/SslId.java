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

/**
 * An SSL ID is a short value build up the following way:
 * Bits 15-12: Module class
 * Bits 11-8:  Number of the partial list extract
 * Bits 7-0:   Number of the partial list
 *
 * Device Class:
 * - 0000: CPU
 * - 0100: IM
 * - 1000: FM
 * - 1100: CP
 */
public enum SslId {

    // S7 300-400 ID's
    MODULE_IDENTIFICATION((short) 0x0011),
    CPU_CHARACTERISTICS((short) 0x0012),
    USER_MEMORY_AREAS((short) 0x0013),
    SYSTEM_AREAS((short) 0x0014),
    BLOCK_TYPES((short) 0x0015),
    INTERRUPT_STATUS((short) 0x0022),
    ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS((short) 0x0025),
    COMPUTATION_STATUS_DATA((short) 0x0032),
    H_CPU_GROUP_INFORMATION((short) 0x0071),
    STATUS_OF_THE_MODULE_LEDS((short) 0x0074),
    SWITCHED_DP_SLAVES_IN_THE_H_SYSTEM((short) 0x0075),
    MODULE_STATUS_INFORMATION((short) 0x0091),
    RACK_STATION_STATUS_INFORMATION_1((short) 0x0092),
    RACK_STATION_STATUS_INFORMATION_2((short) 0x0094),
    EXTENDED_DP_MASTER_SYSTEM_PROFINET_IO_SYSTEM_INFORMATION((short) 0x0095),
    MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFINET_DP((short) 0x0096),
    MODULE_STATUS_INFORMATION_PROFINET_IO((short) 0x009C),
    DIAGNOSTIC_BUFFER_OF_THE_CPU((short) 0x00A0),
    MODULE_DIAGNOSTIC_INFORMATION_DATA_RECORD_0((short) 0x00B1),
    MODULE_DIAGNOSTIC_INFORMATION_DATA_RECORD_1_GEOGRAPHICAL_ADDRESS((short) 0x00B2),
    MODULE_DIAGNOSTIC_INFORMATION_DATA_RECORD_1_LOGICAL_ADDRESS((short) 0x00B3),
    DIAGNOSTIC_DATA_OF_A_DP_SLAVE((short) 0x00B4),

    IDENTIFICATION_OF_ALL_COMPONENTS((short) 0x001c),
    INFORMATION_ABOUT_COMMUNICATION_UNIT((short) 0x0131),
    CURRENT_MODE_TRANSITION((short) 0x0424);

    private static final Logger logger = LoggerFactory.getLogger(SslId.class);

    private final short code;

    SslId(short code) {
        this.code = code;
    }

    public short getCode() {
        return code;
    }

    private final static Map<Short, SslId> map;

    static {
        map = new HashMap<>();
        for (SslId sslId : SslId.values()) {
            map.put(sslId.code, sslId);
        }
    }

    public static SslId valueOf(short code) {
        if(!map.containsKey(code)) {
            logger.error("SslId for code {} not found", code);
        }
        return map.get(code);
    }


}
