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

public enum PlcResponseCode {
  OK((short) 0x01),
  NOT_FOUND((short) 0x02),
  ACCESS_DENIED((short) 0x03),
  INVALID_ADDRESS((short) 0x04),
  INVALID_DATATYPE((short) 0x06),
  INVALID_DATA((short) 0x07),
  INTERNAL_ERROR((short) 0x08),
  REMOTE_BUSY((short) 0x09),
  REMOTE_ERROR((short) 0x0A),
  UNSUPPORTED((short) 0x0B),
  RESPONSE_PENDING((short) 0x0C);
  private static final Map<Short, PlcResponseCode> map;

  static {
    map = new HashMap<>();
    for (PlcResponseCode value : PlcResponseCode.values()) {
      map.put((short) value.getValue(), value);
    }
  }

  private final short value;

  PlcResponseCode(short value) {
    this.value = value;
  }

  public short getValue() {
    return value;
  }

  public static PlcResponseCode enumForValue(short value) {
    return map.get(value);
  }

  public static Boolean isDefined(short value) {
    return map.containsKey(value);
  }
}
