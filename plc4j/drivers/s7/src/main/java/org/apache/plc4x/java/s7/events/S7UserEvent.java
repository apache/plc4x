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
package org.apache.plc4x.java.s7.events;

import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import org.apache.plc4x.java.s7.readwrite.S7PayloadDiagnosticMessage;

public class S7UserEvent extends S7SysEvent {

    S7UserEvent(Instant instant, Map<String, Object> map) {
        super(instant, map);
    }

    public static S7UserEvent of(S7PayloadDiagnosticMessage payload) {
        S7SysEvent event = S7SysEvent.of(payload);
        Map<String, Object> map = new HashMap<>(event.getMap());
        map.put(Fields.TYPE.name(), "USER");
        return new S7UserEvent(event.getTimestamp(), map);
    }
}
