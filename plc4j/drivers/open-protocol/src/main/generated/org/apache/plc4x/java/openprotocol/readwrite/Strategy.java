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
package org.apache.plc4x.java.openprotocol.readwrite;

import java.util.HashMap;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum Strategy {
  TorqueControl((int) 1),
  TorqueControlAngleMonitoring((int) 2),
  TorqueControlAngleControlAnd((int) 3),
  AngleControlTorqueMonitoring((int) 4),
  DsControl((int) 5),
  DsControlTorqueMonitoring((int) 6),
  ReverseAngle((int) 7),
  ReverseTorque((int) 8),
  ClickWrench((int) 9),
  RotateSpindleForward((int) 10),
  TorqueControlAngleControlOr((int) 11),
  RotateSpindleReverse((int) 12),
  HomePositionForward((int) 13),
  EpMonitoring((int) 14),
  Yield((int) 15),
  EpFixed((int) 16),
  EpControl((int) 17),
  EpAngleShutoff((int) 18),
  YieldTorqueControlOr((int) 19),
  SnugGradient((int) 20),
  ResidualTorqueTime((int) 21),
  ResidualTorqueAngle((int) 22),
  BreakawayPeak((int) 23),
  LooseAndTightening((int) 24),
  HomePositionReverse((int) 25),
  PvtCompWithSnug((int) 26),
  NoStrategy((int) 99);
  private static final Map<Integer, Strategy> map;

  static {
    map = new HashMap<>();
    for (Strategy value : Strategy.values()) {
      map.put((int) value.getValue(), value);
    }
  }

  private final int value;

  Strategy(int value) {
    this.value = value;
  }

  public int getValue() {
    return value;
  }

  public static Strategy enumForValue(int value) {
    return map.get(value);
  }

  public static Boolean isDefined(int value) {
    return map.containsKey(value);
  }
}