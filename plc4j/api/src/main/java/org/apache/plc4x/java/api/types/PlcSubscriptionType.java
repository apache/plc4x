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
package org.apache.plc4x.java.api.types;

import java.util.HashMap;
import java.util.Map;

public enum PlcSubscriptionType {
  CYCLIC((short) 0x01),
  CHANGE_OF_STATE((short) 0x02),
  EVENT((short) 0x03);
  private static final Map<Short, PlcSubscriptionType> map;

  static {
    map = new HashMap<>();
    for (PlcSubscriptionType value : PlcSubscriptionType.values()) {
      map.put((short) value.getValue(), value);
    }
  }

  private final short value;

  PlcSubscriptionType(short value) {
    this.value = value;
  }

  public short getValue() {
    return value;
  }

  public static PlcSubscriptionType enumForValue(short value) {
    return map.get(value);
  }

  public static Boolean isDefined(short value) {
    return map.containsKey(value);
  }
}
