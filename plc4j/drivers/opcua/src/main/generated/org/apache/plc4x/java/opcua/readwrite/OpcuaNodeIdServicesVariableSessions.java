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

public enum OpcuaNodeIdServicesVariableSessions {
  SessionsDiagnosticsSummaryType_SessionDiagnosticsArray((int) 2027L),
  SessionsDiagnosticsSummaryType_SessionSecurityDiagnosticsArray((int) 2028L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics((int) 12098L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_SessionId((int) 12099L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_SessionName(
      (int) 12100L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ClientDescription(
      (int) 12101L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ServerUri((int) 12102L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_EndpointUrl(
      (int) 12103L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_LocaleIds((int) 12104L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ActualSessionTimeout(
      (int) 12105L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_MaxResponseMessageSize(
      (int) 12106L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ClientConnectionTime(
      (int) 12107L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ClientLastContactTime(
      (int) 12108L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_CurrentSubscriptionsCount(
      (int) 12109L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_CurrentMonitoredItemsCount(
      (int) 12110L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_CurrentPublishRequestsInQueue(
      (int) 12111L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_TotalRequestCount(
      (int) 12112L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_UnauthorizedRequestCount(
      (int) 12113L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ReadCount((int) 12114L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_HistoryReadCount(
      (int) 12115L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_WriteCount((int) 12116L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_HistoryUpdateCount(
      (int) 12117L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_CallCount((int) 12118L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_CreateMonitoredItemsCount(
      (int) 12119L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ModifyMonitoredItemsCount(
      (int) 12120L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_SetMonitoringModeCount(
      (int) 12121L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_SetTriggeringCount(
      (int) 12122L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_DeleteMonitoredItemsCount(
      (int) 12123L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_CreateSubscriptionCount(
      (int) 12124L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_ModifySubscriptionCount(
      (int) 12125L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_SetPublishingModeCount(
      (int) 12126L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_PublishCount(
      (int) 12127L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_RepublishCount(
      (int) 12128L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_TransferSubscriptionsCount(
      (int) 12129L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_DeleteSubscriptionsCount(
      (int) 12130L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_AddNodesCount(
      (int) 12131L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_AddReferencesCount(
      (int) 12132L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_DeleteNodesCount(
      (int) 12133L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_DeleteReferencesCount(
      (int) 12134L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_BrowseCount(
      (int) 12135L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_BrowseNextCount(
      (int) 12136L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_TranslateBrowsePathsToNodeIdsCount(
      (int) 12137L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_QueryFirstCount(
      (int) 12138L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_QueryNextCount(
      (int) 12139L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_RegisterNodesCount(
      (int) 12140L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionDiagnostics_UnregisterNodesCount(
      (int) 12141L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics((int) 12142L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_SessionId(
      (int) 12143L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_ClientUserIdOfSession(
      (int) 12144L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_ClientUserIdHistory(
      (int) 12145L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_AuthenticationMechanism(
      (int) 12146L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_Encoding(
      (int) 12147L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_TransportProtocol(
      (int) 12148L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_SecurityMode(
      (int) 12149L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_SecurityPolicyUri(
      (int) 12150L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SessionSecurityDiagnostics_ClientCertificate(
      (int) 12151L),
  SessionsDiagnosticsSummaryType_ClientName_Placeholder_SubscriptionDiagnosticsArray((int) 12152L);
  private static final Map<Integer, OpcuaNodeIdServicesVariableSessions> map;

  static {
    map = new HashMap<>();
    for (OpcuaNodeIdServicesVariableSessions value : OpcuaNodeIdServicesVariableSessions.values()) {
      map.put((int) value.getValue(), value);
    }
  }

  private final int value;

  OpcuaNodeIdServicesVariableSessions(int value) {
    this.value = value;
  }

  public int getValue() {
    return value;
  }

  public static OpcuaNodeIdServicesVariableSessions enumForValue(int value) {
    return map.get(value);
  }

  public static Boolean isDefined(int value) {
    return map.containsKey(value);
  }
}
