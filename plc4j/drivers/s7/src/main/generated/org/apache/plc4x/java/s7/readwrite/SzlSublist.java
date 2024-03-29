/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.s7.readwrite;

import java.util.HashMap;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum SzlSublist {
  NONE((short) 0x00),
  MODULE_IDENTIFICATION((short) 0x11),
  CPU_FEATURES((short) 0x12),
  USER_MEMORY_AREA((short) 0x13),
  SYSTEM_AREAS((short) 0x14),
  BLOCK_TYPES((short) 0x15),
  STATUS_MODULE_LEDS((short) 0x19),
  COMPONENT_IDENTIFICATION((short) 0x1C),
  INTERRUPT_STATUS((short) 0x22),
  ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS((short) 0x25),
  COMMUNICATION_STATUS_DATA((short) 0x32),
  H_CPU_GROUP_INFORMATION((short) 0x71),
  STATUS_SINGLE_MODULE_LED((short) 0x74),
  SWITCHED_DP_SLAVES_H_SYSTEM((short) 0x75),
  DP_MASTER_SYSTEM_INFORMATION((short) 0x90),
  MODULE_STATUS_INFORMATION((short) 0x91),
  RACK_OR_STATION_STATUS_INFORMATION((short) 0x92),
  RACK_OR_STATION_STATUS_INFORMATION_2((short) 0x94),
  ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION((short) 0x95),
  MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP((short) 0x96),
  TOOL_CHANGER_INFORMATION_PROFINET((short) 0x9C),
  DIAGNOSTIC_BUFFER((short) 0xA0),
  MODULE_DIAGNOSTIC_INFORMATION_DR0((short) 0xB1),
  MODULE_DIAGNOSTIC_INFORMATION_DR1_GI((short) 0xB2),
  MODULE_DIAGNOSTIC_INFORMATION_DR1_LA((short) 0xB3),
  DIAGNOSTIC_DATA_DP_SLAVE((short) 0xB4);
  private static final Map<Short, SzlSublist> map;

  static {
    map = new HashMap<>();
    for (SzlSublist value : SzlSublist.values()) {
      map.put((short) value.getValue(), value);
    }
  }

  private final short value;

  SzlSublist(short value) {
    this.value = value;
  }

  public short getValue() {
    return value;
  }

  public static SzlSublist enumForValue(short value) {
    return map.get(value);
  }

  public static Boolean isDefined(short value) {
    return map.containsKey(value);
  }
}
