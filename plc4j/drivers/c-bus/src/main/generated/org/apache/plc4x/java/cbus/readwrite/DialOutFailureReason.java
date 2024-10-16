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
package org.apache.plc4x.java.cbus.readwrite;

import java.util.HashMap;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum DialOutFailureReason {
  NO_DIAL_TONE((short) 0x01),
  NO_ANSWER((short) 0x02),
  NO_VALID_ACKNOWLEDGEMENT_OF_PROMPTS((short) 0x03),
  NUMBER_WAS_UNOBTAINABLE_DOES_NOT_EXIST((short) 0x04),
  NUMBER_WAS_BUSY((short) 0x05),
  INTERNAL_FAILURE((short) 0x06);
  private static final Map<Short, DialOutFailureReason> map;

  static {
    map = new HashMap<>();
    for (DialOutFailureReason value : DialOutFailureReason.values()) {
      map.put((short) value.getValue(), value);
    }
  }

  private final short value;

  DialOutFailureReason(short value) {
    this.value = value;
  }

  public short getValue() {
    return value;
  }

  public static DialOutFailureReason enumForValue(short value) {
    return map.get(value);
  }

  public static Boolean isDefined(short value) {
    return map.containsKey(value);
  }
}
