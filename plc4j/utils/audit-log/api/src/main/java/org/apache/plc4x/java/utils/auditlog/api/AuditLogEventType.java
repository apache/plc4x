/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package org.apache.plc4x.java.utils.auditlog.api;

/**
 * Enumeration of different types of events that can be logged to the audit log.
 */
public enum AuditLogEventType {

    /**
     * General configuration information
     */
    CONFIG,

    /**
     * General system information
     */
    SYSTEM,

    /**
     * Starts a new test suite.
     */
    CONNECT,

    /**
     * Outgoing raw bytes.
     */
    OUTGOING_BYTES,

    /**
     * Outgoing PLC messages.
     */
    OUTGOING_MESSAGE,

    /**
     * Incoming raw bytes.
     */
    INCOMING_BYTES,

    /**
     * Incoming PLC messages.
     */
    INCOMING_MESSAGE,

    /**
     * Incoming API requests.
     */
    API_REQUEST,

    /**
     * Outgoing API responses.
     */
    API_RESPONSE,

    /**
     * Ends a test suite.
     */
    CLOSE,

    /**
     * Log errors
     */
    ERROR
}
