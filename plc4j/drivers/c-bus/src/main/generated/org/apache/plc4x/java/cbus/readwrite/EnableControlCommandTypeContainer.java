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

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum EnableControlCommandTypeContainer {
  EnableControlCommandSetNetworkVariable0_2Bytes(
      (short) 0x02, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable1_2Bytes(
      (short) 0x0A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable2_2Bytes(
      (short) 0x12, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable3_2Bytes(
      (short) 0x1A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable4_2Bytes(
      (short) 0x22, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable5_2Bytes(
      (short) 0x2A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable6_2Bytes(
      (short) 0x32, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable7_2Bytes(
      (short) 0x3A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable8_2Bytes(
      (short) 0x42, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable9_2Bytes(
      (short) 0x4A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable10_2Bytes(
      (short) 0x52, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable11_2Bytes(
      (short) 0x5A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable12_2Bytes(
      (short) 0x62, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable13_2Bytes(
      (short) 0x6A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable14_2Bytes(
      (short) 0x72, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE),
  EnableControlCommandSetNetworkVariable15_2Bytes(
      (short) 0x7A, (short) 2, EnableControlCommandType.SET_NETWORK_VARIABLE);
  private static final Map<Short, EnableControlCommandTypeContainer> map;

  static {
    map = new HashMap<>();
    for (EnableControlCommandTypeContainer value : EnableControlCommandTypeContainer.values()) {
      map.put((short) value.getValue(), value);
    }
  }

  private short value;
  private short numBytes;
  private EnableControlCommandType commandType;

  EnableControlCommandTypeContainer(
      short value, short numBytes, EnableControlCommandType commandType) {
    this.value = value;
    this.numBytes = numBytes;
    this.commandType = commandType;
  }

  public short getValue() {
    return value;
  }

  public short getNumBytes() {
    return numBytes;
  }

  public static EnableControlCommandTypeContainer firstEnumForFieldNumBytes(short fieldValue) {
    for (EnableControlCommandTypeContainer _val : EnableControlCommandTypeContainer.values()) {
      if (_val.getNumBytes() == fieldValue) {
        return _val;
      }
    }
    return null;
  }

  public static List<EnableControlCommandTypeContainer> enumsForFieldNumBytes(short fieldValue) {
    List<EnableControlCommandTypeContainer> _values = new ArrayList();
    for (EnableControlCommandTypeContainer _val : EnableControlCommandTypeContainer.values()) {
      if (_val.getNumBytes() == fieldValue) {
        _values.add(_val);
      }
    }
    return _values;
  }

  public EnableControlCommandType getCommandType() {
    return commandType;
  }

  public static EnableControlCommandTypeContainer firstEnumForFieldCommandType(
      EnableControlCommandType fieldValue) {
    for (EnableControlCommandTypeContainer _val : EnableControlCommandTypeContainer.values()) {
      if (_val.getCommandType() == fieldValue) {
        return _val;
      }
    }
    return null;
  }

  public static List<EnableControlCommandTypeContainer> enumsForFieldCommandType(
      EnableControlCommandType fieldValue) {
    List<EnableControlCommandTypeContainer> _values = new ArrayList();
    for (EnableControlCommandTypeContainer _val : EnableControlCommandTypeContainer.values()) {
      if (_val.getCommandType() == fieldValue) {
        _values.add(_val);
      }
    }
    return _values;
  }

  public static EnableControlCommandTypeContainer enumForValue(short value) {
    return map.get(value);
  }

  public static Boolean isDefined(short value) {
    return map.containsKey(value);
  }
}
