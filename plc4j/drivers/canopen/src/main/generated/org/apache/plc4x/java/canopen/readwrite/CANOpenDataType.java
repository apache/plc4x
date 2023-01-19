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
package org.apache.plc4x.java.canopen.readwrite;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum CANOpenDataType {
  BOOLEAN((long) 0L, (String) "BIT", (short) 1),
  UNSIGNED8((long) 1L, (String) "USINT", (short) 8),
  UNSIGNED16((long) 2L, (String) "UINT", (short) 16),
  UNSIGNED24((long) 3L, (String) "RAW_BYTE_ARRAY", (short) 24),
  UNSIGNED32((long) 4L, (String) "UDINT", (short) 32),
  UNSIGNED40((long) 5L, (String) "RAW_BYTE_ARRAY", (short) 40),
  UNSIGNED48((long) 6L, (String) "RAW_BYTE_ARRAY", (short) 48),
  UNSIGNED56((long) 7L, (String) "RAW_BYTE_ARRAY", (short) 56),
  UNSIGNED64((long) 8L, (String) "ULINT", (short) 64),
  INTEGER8((long) 9L, (String) "SINT", (short) 8),
  INTEGER16((long) 10L, (String) "INT", (short) 16),
  INTEGER24((long) 11L, (String) "RAW_BYTE_ARRAY", (short) 24),
  INTEGER32((long) 12L, (String) "DINT", (short) 32),
  INTEGER40((long) 13L, (String) "RAW_BYTE_ARRAY", (short) 40),
  INTEGER48((long) 14L, (String) "RAW_BYTE_ARRAY", (short) 48),
  INTEGER56((long) 15L, (String) "RAW_BYTE_ARRAY", (short) 56),
  INTEGER64((long) 16L, (String) "LINT", (short) 64),
  REAL32((long) 17L, (String) "REAL", (short) 32),
  REAL64((long) 18L, (String) "LREAL", (short) 64),
  RECORD((long) 19L, (String) "BYTE", (short) 8),
  OCTET_STRING((long) 20L, (String) "BYTE", (short) 8),
  VISIBLE_STRING((long) 21L, (String) "CHAR", (short) 8),
  UNICODE_STRING((long) 22L, (String) "WCHAR", (short) 16),
  TIME_OF_DAY((long) 23L, (String) "TIME_OF_DAY", (short) 48),
  TIME_DIFFERENCE((long) 24L, (String) "TIME", (short) 48);
  private static final Map<Long, CANOpenDataType> map;

  static {
    map = new HashMap<>();
    for (CANOpenDataType value : CANOpenDataType.values()) {
      map.put((long) value.getValue(), value);
    }
  }

  private long value;
  private String plcValueName;
  private short numBits;

  CANOpenDataType(long value, String plcValueName, short numBits) {
    this.value = value;
    this.plcValueName = plcValueName;
    this.numBits = numBits;
  }

  public long getValue() {
    return value;
  }

  public String getPlcValueName() {
    return plcValueName;
  }

  public static CANOpenDataType firstEnumForFieldPlcValueName(String fieldValue) {
    for (CANOpenDataType _val : CANOpenDataType.values()) {
      if (_val.getPlcValueName().equals(fieldValue)) {
        return _val;
      }
    }
    return null;
  }

  public static List<CANOpenDataType> enumsForFieldPlcValueName(String fieldValue) {
    List<CANOpenDataType> _values = new ArrayList();
    for (CANOpenDataType _val : CANOpenDataType.values()) {
      if (_val.getPlcValueName().equals(fieldValue)) {
        _values.add(_val);
      }
    }
    return _values;
  }

  public short getNumBits() {
    return numBits;
  }

  public static CANOpenDataType firstEnumForFieldNumBits(short fieldValue) {
    for (CANOpenDataType _val : CANOpenDataType.values()) {
      if (_val.getNumBits() == fieldValue) {
        return _val;
      }
    }
    return null;
  }

  public static List<CANOpenDataType> enumsForFieldNumBits(short fieldValue) {
    List<CANOpenDataType> _values = new ArrayList();
    for (CANOpenDataType _val : CANOpenDataType.values()) {
      if (_val.getNumBits() == fieldValue) {
        _values.add(_val);
      }
    }
    return _values;
  }

  public static CANOpenDataType enumForValue(long value) {
    return map.get(value);
  }

  public static Boolean isDefined(long value) {
    return map.containsKey(value);
  }
}
