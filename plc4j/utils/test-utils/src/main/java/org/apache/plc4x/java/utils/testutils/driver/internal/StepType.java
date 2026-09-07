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
package org.apache.plc4x.java.utils.testutils.driver.internal;

/**
 * Enum defining the types of test steps that can be executed in a driver test suite.
 */
public enum StepType {
    /** Outgoing PLC message step */
    OUTGOING_PLC_MESSAGE,
    /** Outgoing PLC bytes step */
    OUTGOING_PLC_BYTES,
    /** Incoming PLC message step */
    INCOMING_PLC_MESSAGE,
    /** Incoming PLC bytes step */
    INCOMING_PLC_BYTES,
    /** API request step */
    API_REQUEST,
    /** API response step */
    API_RESPONSE,
    /** API evnet step */
    API_EVENT,
    /** Delay step */
    DELAY,
    /** Terminate step */
    TERMINATE
}
