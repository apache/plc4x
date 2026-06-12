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
package org.apache.plc4x.java.s7.tag;

import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.s7.readwrite.QueryType;

import java.io.Serializable;
import java.util.Objects;

/**
 * Marker tag for the S7Comm alarm subsystem. Two flavours:
 * <ul>
 *   <li>{@link Kind#PUSH} ({@code "ALM"}) — long-lived push subscription. Alarm indications
 *       are delivered by the PLC asynchronously when SFC17/18 fire in user code.</li>
 *   <li>{@link Kind#QUERY} ({@code "QUERY:ALARM_S"} / {@code "QUERY:ALARM_8"}) — one-shot
 *       synchronous fetch of currently-buffered alarms. Returns once and goes inert.</li>
 * </ul>
 * Both ride S7Comm UserData {@code cpuFunctionGroup=0x04}; PUSH uses MsgSubscription
 * (subfn 0x02), QUERY uses AlarmQuery (subfn 0x13).
 */
public class S7AlarmTag implements PlcTag, Serializable {

    public static final String ADDRESS_ALARM     = "ALM";
    public static final String ADDRESS_QUERY_S   = "QUERY:ALARM_S";
    public static final String ADDRESS_QUERY_8   = "QUERY:ALARM_8";

    public enum Kind { PUSH, QUERY }

    private static final S7AlarmTag PUSH_INSTANCE   = new S7AlarmTag(Kind.PUSH, null);
    private static final S7AlarmTag QUERY_S_INSTANCE = new S7AlarmTag(Kind.QUERY, QueryType.ALARM_S);
    private static final S7AlarmTag QUERY_8_INSTANCE = new S7AlarmTag(Kind.QUERY, QueryType.ALARM_8);

    private final Kind kind;
    private final QueryType queryType;

    private S7AlarmTag(Kind kind, QueryType queryType) {
        this.kind = Objects.requireNonNull(kind, "kind");
        this.queryType = queryType;
    }

    public static boolean matches(String tagString) {
        return ADDRESS_ALARM.equalsIgnoreCase(tagString)
            || ADDRESS_QUERY_S.equalsIgnoreCase(tagString)
            || ADDRESS_QUERY_8.equalsIgnoreCase(tagString);
    }

    public static S7AlarmTag of(String tagString) {
        if (ADDRESS_ALARM.equalsIgnoreCase(tagString)) return PUSH_INSTANCE;
        if (ADDRESS_QUERY_S.equalsIgnoreCase(tagString)) return QUERY_S_INSTANCE;
        if (ADDRESS_QUERY_8.equalsIgnoreCase(tagString)) return QUERY_8_INSTANCE;
        throw new IllegalArgumentException("Not an alarm tag: " + tagString);
    }

    public static S7AlarmTag instance() {
        return PUSH_INSTANCE;
    }

    public Kind getKind() {
        return kind;
    }

    /** Non-null for {@link Kind#QUERY} tags, null for {@link Kind#PUSH}. */
    public QueryType getQueryType() {
        return queryType;
    }

    @Override
    public String getAddressString() {
        return switch (kind) {
            case PUSH -> ADDRESS_ALARM;
            case QUERY -> queryType == QueryType.ALARM_8 ? ADDRESS_QUERY_8 : ADDRESS_QUERY_S;
        };
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.Struct;
    }
}
