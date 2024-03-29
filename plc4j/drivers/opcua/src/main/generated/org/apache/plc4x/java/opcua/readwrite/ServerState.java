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
package org.apache.plc4x.java.opcua.readwrite;

import java.util.HashMap;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum ServerState {
  serverStateRunning((long) 0L),
  serverStateFailed((long) 1L),
  serverStateNoConfiguration((long) 2L),
  serverStateSuspended((long) 3L),
  serverStateShutdown((long) 4L),
  serverStateTest((long) 5L),
  serverStateCommunicationFault((long) 6L),
  serverStateUnknown((long) 7L);
  private static final Map<Long, ServerState> map;

  static {
    map = new HashMap<>();
    for (ServerState value : ServerState.values()) {
      map.put((long) value.getValue(), value);
    }
  }

  private final long value;

  ServerState(long value) {
    this.value = value;
  }

  public long getValue() {
    return value;
  }

  public static ServerState enumForValue(long value) {
    return map.get(value);
  }

  public static Boolean isDefined(long value) {
    return map.containsKey(value);
  }
}
