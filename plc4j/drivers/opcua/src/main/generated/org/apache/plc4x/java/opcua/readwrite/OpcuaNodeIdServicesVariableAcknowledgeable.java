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

public enum OpcuaNodeIdServicesVariableAcknowledgeable {
  AcknowledgeableConditionType_EnabledState((int) 9073L),
  AcknowledgeableConditionType_EnabledState_Id((int) 9074L),
  AcknowledgeableConditionType_EnabledState_Name((int) 9075L),
  AcknowledgeableConditionType_EnabledState_Number((int) 9076L),
  AcknowledgeableConditionType_EnabledState_EffectiveDisplayName((int) 9077L),
  AcknowledgeableConditionType_EnabledState_TransitionTime((int) 9078L),
  AcknowledgeableConditionType_EnabledState_EffectiveTransitionTime((int) 9079L),
  AcknowledgeableConditionType_EnabledState_TrueState((int) 9080L),
  AcknowledgeableConditionType_EnabledState_FalseState((int) 9081L),
  AcknowledgeableConditionType_AckedState((int) 9093L),
  AcknowledgeableConditionType_AckedState_Id((int) 9094L),
  AcknowledgeableConditionType_AckedState_Name((int) 9095L),
  AcknowledgeableConditionType_AckedState_Number((int) 9096L),
  AcknowledgeableConditionType_AckedState_EffectiveDisplayName((int) 9097L),
  AcknowledgeableConditionType_AckedState_TransitionTime((int) 9098L),
  AcknowledgeableConditionType_AckedState_EffectiveTransitionTime((int) 9099L),
  AcknowledgeableConditionType_AckedState_TrueState((int) 9100L),
  AcknowledgeableConditionType_AckedState_FalseState((int) 9101L),
  AcknowledgeableConditionType_ConfirmedState((int) 9102L),
  AcknowledgeableConditionType_ConfirmedState_Id((int) 9103L),
  AcknowledgeableConditionType_ConfirmedState_Name((int) 9104L),
  AcknowledgeableConditionType_ConfirmedState_Number((int) 9105L),
  AcknowledgeableConditionType_ConfirmedState_EffectiveDisplayName((int) 9106L),
  AcknowledgeableConditionType_ConfirmedState_TransitionTime((int) 9107L),
  AcknowledgeableConditionType_ConfirmedState_EffectiveTransitionTime((int) 9108L),
  AcknowledgeableConditionType_ConfirmedState_TrueState((int) 9109L),
  AcknowledgeableConditionType_ConfirmedState_FalseState((int) 9110L),
  AcknowledgeableConditionType_Acknowledge_InputArguments((int) 9112L),
  AcknowledgeableConditionType_Confirm_InputArguments((int) 9114L);
  private static final Map<Integer, OpcuaNodeIdServicesVariableAcknowledgeable> map;

  static {
    map = new HashMap<>();
    for (OpcuaNodeIdServicesVariableAcknowledgeable value :
        OpcuaNodeIdServicesVariableAcknowledgeable.values()) {
      map.put((int) value.getValue(), value);
    }
  }

  private final int value;

  OpcuaNodeIdServicesVariableAcknowledgeable(int value) {
    this.value = value;
  }

  public int getValue() {
    return value;
  }

  public static OpcuaNodeIdServicesVariableAcknowledgeable enumForValue(int value) {
    return map.get(value);
  }

  public static Boolean isDefined(int value) {
    return map.containsKey(value);
  }
}