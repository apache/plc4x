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

package org.apache.plc4x.java.api.types;

public enum ConnectionStateChangeType {
    // Connection lifecycle
    CONNECTED,              // Connection established successfully
    DISCONNECTED,           // Graceful disconnection (close() called)
    CONNECTION_LOST,        // Unexpected disconnection (network error, timeout, etc.)

    // Tag/metadata changes
    TAGS_CHANGED,           // Available tags changed, re-browse needed

    // PLC mode changes
    MODE_RUN,               // PLC in RUN mode (executing program)
    MODE_STOP,              // PLC in STOP mode (program halted)
    MODE_CONFIG             // PLC in CONFIG/PROGRAM mode (being programmed)
}
