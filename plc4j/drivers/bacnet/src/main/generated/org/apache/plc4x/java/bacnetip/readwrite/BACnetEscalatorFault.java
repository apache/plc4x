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
package org.apache.plc4x.java.bacnetip.readwrite;

import java.util.HashMap;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum BACnetEscalatorFault {
  CONTROLLER_FAULT((int) 0),
  DRIVE_AND_MOTOR_FAULT((int) 1),
  MECHANICAL_COMPONENT_FAULT((int) 2),
  OVERSPEED_FAULT((int) 3),
  POWER_SUPPLY_FAULT((int) 4),
  SAFETY_DEVICE_FAULT((int) 5),
  CONTROLLER_SUPPLY_FAULT((int) 6),
  DRIVE_TEMPERATURE_EXCEEDED((int) 7),
  COMB_PLATE_FAULT((int) 8),
  VENDOR_PROPRIETARY_VALUE((int) 0XFFFF);
  private static final Map<Integer, BACnetEscalatorFault> map;

  static {
    map = new HashMap<>();
    for (BACnetEscalatorFault value : BACnetEscalatorFault.values()) {
      map.put((int) value.getValue(), value);
    }
  }

  private final int value;

  BACnetEscalatorFault(int value) {
    this.value = value;
  }

  public int getValue() {
    return value;
  }

  public static BACnetEscalatorFault enumForValue(int value) {
    return map.get(value);
  }

  public static Boolean isDefined(int value) {
    return map.containsKey(value);
  }
}
