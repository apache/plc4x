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
package org.apache.plc4x.java.profinet.readwrite;

import java.util.HashMap;
import java.util.Map;

// Code generated by code-generation. DO NOT EDIT.

public enum PnDcp_FrameId {
  RESERVED((int) 0x0000),
  PTCP_RTSyncPDUWithFollowUp((int) 0x0020),
  PTCP_RTSyncPDU((int) 0x0080),
  RT_CLASS_3((int) 0x0100),
  RT_CLASS_1((int) 0x8000),
  RT_CLASS_UDP((int) 0xC000),
  Alarm_High((int) 0xFC01),
  Alarm_Low((int) 0xFE01),
  DCP_Hello_ReqPDU((int) 0xFEFC),
  DCP_GetSet_PDU((int) 0xFEFD),
  DCP_Identify_ReqPDU((int) 0xFEFE),
  DCP_Identify_ResPDU((int) 0xFEFF),
  PTCP_AnnouncePDU((int) 0xFF00),
  PTCP_FollowUpPDU((int) 0xFF20),
  PTCP_DelayReqPDU((int) 0xFF40),
  PTCP_DelayResPDUWithFollowUp((int) 0xFF41),
  PTCP_DelayFuResPDUWithFollowUp((int) 0xFF42),
  PTCP_DelayResPDUWithoutFollowUp((int) 0xFF43),
  FragmentationFrameId((int) 0xFF80);
  private static final Map<Integer, PnDcp_FrameId> map;

  static {
    map = new HashMap<>();
    for (PnDcp_FrameId value : PnDcp_FrameId.values()) {
      map.put((int) value.getValue(), value);
    }
  }

  private int value;

  PnDcp_FrameId(int value) {
    this.value = value;
  }

  public int getValue() {
    return value;
  }

  public static PnDcp_FrameId enumForValue(int value) {
    return map.get(value);
  }

  public static Boolean isDefined(int value) {
    return map.containsKey(value);
  }
}
