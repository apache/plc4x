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

public enum SecurityCommandTypeContainer {
  SecurityCommandOff_0Bytes((short) 0x00, (byte) 0, SecurityCommandType.OFF),
  SecurityCommandOff_1Bytes((short) 0x01, (byte) 1, SecurityCommandType.OFF),
  SecurityCommandOff_2Bytes((short) 0x02, (byte) 2, SecurityCommandType.OFF),
  SecurityCommandOff_3Bytes((short) 0x03, (byte) 3, SecurityCommandType.OFF),
  SecurityCommandOff_4Bytes((short) 0x04, (byte) 4, SecurityCommandType.OFF),
  SecurityCommandOff_5Bytes((short) 0x05, (byte) 5, SecurityCommandType.OFF),
  SecurityCommandOff_6Bytes((short) 0x06, (byte) 6, SecurityCommandType.OFF),
  SecurityCommandOff_7Bytes((short) 0x07, (byte) 7, SecurityCommandType.OFF),
  SecurityCommandEvent_0Bytes((short) 0x08, (byte) 0, SecurityCommandType.EVENT),
  SecurityCommandEvent_1Bytes((short) 0x09, (byte) 1, SecurityCommandType.EVENT),
  SecurityCommandEvent_2Bytes((short) 0x0A, (byte) 2, SecurityCommandType.EVENT),
  SecurityCommandEvent_3Bytes((short) 0x0B, (byte) 3, SecurityCommandType.EVENT),
  SecurityCommandEvent_4Bytes((short) 0x0C, (byte) 4, SecurityCommandType.EVENT),
  SecurityCommandEvent_5Bytes((short) 0x0D, (byte) 5, SecurityCommandType.EVENT),
  SecurityCommandEvent_6Bytes((short) 0x0E, (byte) 6, SecurityCommandType.EVENT),
  SecurityCommandEvent_7Bytes((short) 0x0F, (byte) 7, SecurityCommandType.EVENT),
  SecurityCommandOn_0Bytes((short) 0x78, (byte) 0, SecurityCommandType.ON),
  SecurityCommandOn_1Bytes((short) 0x79, (byte) 1, SecurityCommandType.ON),
  SecurityCommandOn_2Bytes((short) 0x7A, (byte) 2, SecurityCommandType.ON),
  SecurityCommandOn_3Bytes((short) 0x7B, (byte) 3, SecurityCommandType.ON),
  SecurityCommandOn_4Bytes((short) 0x7C, (byte) 4, SecurityCommandType.ON),
  SecurityCommandOn_5Bytes((short) 0x7D, (byte) 5, SecurityCommandType.ON),
  SecurityCommandOn_6Bytes((short) 0x7E, (byte) 6, SecurityCommandType.ON),
  SecurityCommandOn_7Bytes((short) 0x7F, (byte) 7, SecurityCommandType.ON),
  SecurityCommandLongOff_0Bytes((short) 0x80, (byte) 8, SecurityCommandType.OFF),
  SecurityCommandLongOff_1Bytes((short) 0x81, (byte) 1, SecurityCommandType.OFF),
  SecurityCommandLongOff_2Bytes((short) 0x82, (byte) 2, SecurityCommandType.OFF),
  SecurityCommandLongOff_3Bytes((short) 0x83, (byte) 3, SecurityCommandType.OFF),
  SecurityCommandLongOff_4Bytes((short) 0x84, (byte) 4, SecurityCommandType.OFF),
  SecurityCommandLongOff_5Bytes((short) 0x85, (byte) 5, SecurityCommandType.OFF),
  SecurityCommandLongOff_6Bytes((short) 0x86, (byte) 6, SecurityCommandType.OFF),
  SecurityCommandLongOff_7Bytes((short) 0x87, (byte) 7, SecurityCommandType.OFF),
  SecurityCommandLongOff_8Bytes((short) 0x88, (byte) 8, SecurityCommandType.OFF),
  SecurityCommandLongOff_9Bytes((short) 0x89, (byte) 9, SecurityCommandType.OFF),
  SecurityCommandLongOff_10Bytes((short) 0x8A, (byte) 10, SecurityCommandType.OFF),
  SecurityCommandLongOff_11Bytes((short) 0x8B, (byte) 11, SecurityCommandType.OFF),
  SecurityCommandLongOff_12Bytes((short) 0x8C, (byte) 12, SecurityCommandType.OFF),
  SecurityCommandLongOff_13Bytes((short) 0x8D, (byte) 13, SecurityCommandType.OFF),
  SecurityCommandLongOff_14Bytes((short) 0x8E, (byte) 14, SecurityCommandType.OFF),
  SecurityCommandLongOff_15Bytes((short) 0x8F, (byte) 15, SecurityCommandType.OFF),
  SecurityCommandLongOff_16Bytes((short) 0x90, (byte) 16, SecurityCommandType.OFF),
  SecurityCommandLongOff_17Bytes((short) 0x91, (byte) 17, SecurityCommandType.OFF),
  SecurityCommandLongOff_18Bytes((short) 0x92, (byte) 18, SecurityCommandType.OFF),
  SecurityCommandLongOff_19Bytes((short) 0x93, (byte) 19, SecurityCommandType.OFF),
  SecurityCommandLongOff_20Bytes((short) 0x94, (byte) 20, SecurityCommandType.OFF),
  SecurityCommandLongOff_21Bytes((short) 0x95, (byte) 21, SecurityCommandType.OFF),
  SecurityCommandLongOff_22Bytes((short) 0x96, (byte) 22, SecurityCommandType.OFF),
  SecurityCommandLongOff_23Bytes((short) 0x97, (byte) 23, SecurityCommandType.OFF),
  SecurityCommandLongOff_24Bytes((short) 0x98, (byte) 24, SecurityCommandType.OFF),
  SecurityCommandLongOff_25Bytes((short) 0x99, (byte) 25, SecurityCommandType.OFF),
  SecurityCommandLongOff_26Bytes((short) 0x9A, (byte) 26, SecurityCommandType.OFF),
  SecurityCommandLongOff_27Bytes((short) 0x9B, (byte) 27, SecurityCommandType.OFF),
  SecurityCommandLongOff_28Bytes((short) 0x9C, (byte) 28, SecurityCommandType.OFF),
  SecurityCommandLongOff_29Bytes((short) 0x9D, (byte) 29, SecurityCommandType.OFF),
  SecurityCommandLongOff_30Bytes((short) 0x9E, (byte) 30, SecurityCommandType.OFF),
  SecurityCommandLongOff_31Bytes((short) 0x9F, (byte) 31, SecurityCommandType.OFF),
  SecurityCommandLongEvent_0Bytes((short) 0xA0, (byte) 0, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_1Bytes((short) 0xA1, (byte) 1, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_2Bytes((short) 0xA2, (byte) 2, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_3Bytes((short) 0xA3, (byte) 3, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_4Bytes((short) 0xA4, (byte) 4, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_5Bytes((short) 0xA5, (byte) 5, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_6Bytes((short) 0xA6, (byte) 6, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_7Bytes((short) 0xA7, (byte) 7, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_8Bytes((short) 0xA8, (byte) 8, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_9Bytes((short) 0xA9, (byte) 9, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_10Bytes((short) 0xAA, (byte) 10, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_11Bytes((short) 0xAB, (byte) 11, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_12Bytes((short) 0xAC, (byte) 12, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_13Bytes((short) 0xAD, (byte) 13, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_14Bytes((short) 0xAE, (byte) 14, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_15Bytes((short) 0xAF, (byte) 15, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_16Bytes((short) 0xB0, (byte) 16, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_17Bytes((short) 0xB1, (byte) 17, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_18Bytes((short) 0xB2, (byte) 18, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_19Bytes((short) 0xB3, (byte) 19, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_20Bytes((short) 0xB4, (byte) 20, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_21Bytes((short) 0xB5, (byte) 21, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_22Bytes((short) 0xB6, (byte) 22, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_23Bytes((short) 0xB7, (byte) 23, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_24Bytes((short) 0xB8, (byte) 24, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_25Bytes((short) 0xB9, (byte) 25, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_26Bytes((short) 0xBA, (byte) 26, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_27Bytes((short) 0xBB, (byte) 27, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_28Bytes((short) 0xBC, (byte) 28, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_29Bytes((short) 0xBD, (byte) 29, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_30Bytes((short) 0xBE, (byte) 30, SecurityCommandType.EVENT),
  SecurityCommandLongEvent_31Bytes((short) 0xBF, (byte) 31, SecurityCommandType.EVENT),
  SecurityCommandLongOn_0Bytes((short) 0xE0, (byte) 0, SecurityCommandType.ON),
  SecurityCommandLongOn_1Bytes((short) 0xE1, (byte) 1, SecurityCommandType.ON),
  SecurityCommandLongOn_2Bytes((short) 0xE2, (byte) 2, SecurityCommandType.ON),
  SecurityCommandLongOn_3Bytes((short) 0xE3, (byte) 3, SecurityCommandType.ON),
  SecurityCommandLongOn_4Bytes((short) 0xE4, (byte) 4, SecurityCommandType.ON),
  SecurityCommandLongOn_5Bytes((short) 0xE5, (byte) 5, SecurityCommandType.ON),
  SecurityCommandLongOn_6Bytes((short) 0xE6, (byte) 6, SecurityCommandType.ON),
  SecurityCommandLongOn_7Bytes((short) 0xE7, (byte) 7, SecurityCommandType.ON),
  SecurityCommandLongOn_8Bytes((short) 0xE8, (byte) 8, SecurityCommandType.ON),
  SecurityCommandLongOn_9Bytes((short) 0xE9, (byte) 9, SecurityCommandType.ON),
  SecurityCommandLongOn_10Bytes((short) 0xEA, (byte) 10, SecurityCommandType.ON),
  SecurityCommandLongOn_11Bytes((short) 0xEB, (byte) 11, SecurityCommandType.ON),
  SecurityCommandLongOn_12Bytes((short) 0xEC, (byte) 12, SecurityCommandType.ON),
  SecurityCommandLongOn_13Bytes((short) 0xED, (byte) 13, SecurityCommandType.ON),
  SecurityCommandLongOn_14Bytes((short) 0xEE, (byte) 14, SecurityCommandType.ON),
  SecurityCommandLongOn_15Bytes((short) 0xEF, (byte) 15, SecurityCommandType.ON),
  SecurityCommandLongOn_16Bytes((short) 0xF0, (byte) 16, SecurityCommandType.ON),
  SecurityCommandLongOn_17Bytes((short) 0xF1, (byte) 17, SecurityCommandType.ON),
  SecurityCommandLongOn_18Bytes((short) 0xF2, (byte) 18, SecurityCommandType.ON),
  SecurityCommandLongOn_19Bytes((short) 0xF3, (byte) 19, SecurityCommandType.ON),
  SecurityCommandLongOn_20Bytes((short) 0xF4, (byte) 20, SecurityCommandType.ON),
  SecurityCommandLongOn_21Bytes((short) 0xF5, (byte) 21, SecurityCommandType.ON),
  SecurityCommandLongOn_22Bytes((short) 0xF6, (byte) 22, SecurityCommandType.ON),
  SecurityCommandLongOn_23Bytes((short) 0xF7, (byte) 23, SecurityCommandType.ON),
  SecurityCommandLongOn_24Bytes((short) 0xF8, (byte) 24, SecurityCommandType.ON),
  SecurityCommandLongOn_25Bytes((short) 0xF9, (byte) 25, SecurityCommandType.ON),
  SecurityCommandLongOn_26Bytes((short) 0xFA, (byte) 26, SecurityCommandType.ON),
  SecurityCommandLongOn_27Bytes((short) 0xFB, (byte) 27, SecurityCommandType.ON),
  SecurityCommandLongOn_28Bytes((short) 0xFC, (byte) 28, SecurityCommandType.ON),
  SecurityCommandLongOn_29Bytes((short) 0xFD, (byte) 29, SecurityCommandType.ON),
  SecurityCommandLongOn_30Bytes((short) 0xFE, (byte) 30, SecurityCommandType.ON),
  SecurityCommandLongOn_31Bytes((short) 0xFF, (byte) 31, SecurityCommandType.ON);
  private static final Map<Short, SecurityCommandTypeContainer> map;

  static {
    map = new HashMap<>();
    for (SecurityCommandTypeContainer value : SecurityCommandTypeContainer.values()) {
      map.put((short) value.getValue(), value);
    }
  }

  private final short value;
  private final byte numBytes;
  private final SecurityCommandType commandType;

  SecurityCommandTypeContainer(short value, byte numBytes, SecurityCommandType commandType) {
    this.value = value;
    this.numBytes = numBytes;
    this.commandType = commandType;
  }

  public short getValue() {
    return value;
  }

  public byte getNumBytes() {
    return numBytes;
  }

  public static SecurityCommandTypeContainer firstEnumForFieldNumBytes(byte fieldValue) {
    for (SecurityCommandTypeContainer _val : SecurityCommandTypeContainer.values()) {
      if (_val.getNumBytes() == fieldValue) {
        return _val;
      }
    }
    return null;
  }

  public static List<SecurityCommandTypeContainer> enumsForFieldNumBytes(byte fieldValue) {
    List<SecurityCommandTypeContainer> _values = new ArrayList<>();
    for (SecurityCommandTypeContainer _val : SecurityCommandTypeContainer.values()) {
      if (_val.getNumBytes() == fieldValue) {
        _values.add(_val);
      }
    }
    return _values;
  }

  public SecurityCommandType getCommandType() {
    return commandType;
  }

  public static SecurityCommandTypeContainer firstEnumForFieldCommandType(
      SecurityCommandType fieldValue) {
    for (SecurityCommandTypeContainer _val : SecurityCommandTypeContainer.values()) {
      if (_val.getCommandType() == fieldValue) {
        return _val;
      }
    }
    return null;
  }

  public static List<SecurityCommandTypeContainer> enumsForFieldCommandType(
      SecurityCommandType fieldValue) {
    List<SecurityCommandTypeContainer> _values = new ArrayList<>();
    for (SecurityCommandTypeContainer _val : SecurityCommandTypeContainer.values()) {
      if (_val.getCommandType() == fieldValue) {
        _values.add(_val);
      }
    }
    return _values;
  }

  public static SecurityCommandTypeContainer enumForValue(short value) {
    return map.get(value);
  }

  public static Boolean isDefined(short value) {
    return map.containsKey(value);
  }
}
