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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.AlarmType;
import org.apache.plc4x.java.s7.readwrite.QueryType;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionAlarmQueryRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionAlarmQueryResponse;
import org.apache.plc4x.java.s7.readwrite.SyntaxIdType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageAckPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageAckObjectPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageObjectPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessagePushType;
import org.apache.plc4x.java.s7.readwrite.AlarmStateType;
import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarm8;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmAckInd;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmS;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmSC;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmSQ;
import org.apache.plc4x.java.s7.readwrite.S7PayloadNotify;
import org.apache.plc4x.java.s7.readwrite.S7PayloadNotify8;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse;
import org.apache.plc4x.java.s7.readwrite.State;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcUDINT;

import java.time.Instant;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * S7Comm UserData "Alarm subscription" service. Handles:
 * <ul>
 *   <li>The MsgSubscription request that asks the PLC to start pushing alarm indications
 *       ({@code cpuFunctionGroup=0x04}, subfunction {@code 0x02}, subscription byte {@code 0x80}).</li>
 *   <li>Recognising and parsing unsolicited alarm indications the PLC subsequently pushes
 *       ({@code cpuFunctionType=0x00}, subfunctions 0x05/0x06/0x0B/0x11/0x12/0x13/0x16).</li>
 * </ul>
 *
 * <p>Indication subfunctions:
 * <ul>
 *   <li>{@code 0x05} ALARM_8 — acyclic alarm with 8 signals</li>
 *   <li>{@code 0x06} NOTIFY  — block-call notify</li>
 *   <li>{@code 0x0B} ALARM_ACK_IND — operator-acknowledged alarm</li>
 *   <li>{@code 0x11} ALARM_SQ — alarm with acknowledgement (S7-300/400 family)</li>
 *   <li>{@code 0x12} ALARM_S  — alarm without acknowledgement</li>
 *   <li>{@code 0x13} ALARM_SC — alarm clear (state-going)</li>
 *   <li>{@code 0x16} NOTIFY_8 — notify with 8 signals</li>
 * </ul>
 */
public final class S7AlarmSubscriptionService {

    /** Subscription byte: bit 7 (0x80) = subscribe to alarm class messages. */
    public static final short SUBSCRIPTION_ALM = 0x80;

    /** Subscription byte 0x00 cancels all message-class subscriptions on this connection. */
    public static final short SUBSCRIPTION_CANCEL = 0x00;

    /**
     * Magic key the PLC echoes back in subscription confirmations. SPI1 used "HmiRtm  "
     * (HMI runtime, padded to 8 ASCII bytes); we keep the same value for compatibility
     * with PLCs that match the key pattern in some firmware paths.
     */
    public static final String MAGIC_KEY = "HmiRtm  ";

    public static final int CPU_FUNCTION_GROUP_CPU = 0x04;
    public static final int CPU_FUNCTION_TYPE_PUSH = 0x00;

    /** All push subfunctions that carry an alarm indication payload. */
    public static final Set<Integer> ALARM_INDICATION_SUBFUNCTIONS =
        Set.of(0x05, 0x06, 0x0B, 0x11, 0x12, 0x13, 0x16);

    private S7AlarmSubscriptionService() {
    }

    /**
     * Build a "subscribe to ALARM_S/SQ push messages" UserData request. The S7Comm UserData
     * MsgSubscription service has two distinct paths:
     * <ul>
     *   <li><b>Event-class subscription</b> (subscription byte ≥ 0x01, e.g. 0x80 for
     *       MODE/SYS/USR/ALM event-class) — {@code dataLength=10}, no appended alarm-state.
     *       This stream carries mode-change and system/user events; it does <em>not</em>
     *       deliver classic ALARM_S/SQ indications produced by SFC17/18.</li>
     *   <li><b>ALARM_S/SQ subscription</b> (subscription byte = 0x00 plus appended
     *       {@code AlarmStateType}) — {@code dataLength=12}. This is the path that arms
     *       the legacy alarm-message stream that SFC17 (ALARM_SQ) and SFC18 (ALARM_S) feed
     *       into on S7-300/400.</li>
     * </ul>
     *
     * <p>We use the second path — that's what user code triggering SFC17/18 expects to flow
     * through. The {@link AlarmStateType} hint depends on the controller family:
     * <ul>
     *   <li>S7-400 → {@code ALARM_INITIATE} (block-of-8 ALARM_8 family).</li>
     *   <li>Anything else (300 / 1200 / 1500) → {@code ALARM_S_INITIATE}.</li>
     * </ul>
     */
    public static S7Message buildSubscribeRequest(int tpduRef, ControllerType controllerType) {
        return buildMsgSubscriptionMessage(tpduRef, /*dataLength=*/12, SUBSCRIPTION_ALM,
            alarmStateForController(controllerType), /*reserve=*/(short) 0);
    }

    /** Build a request that cancels the ALARM_S/SQ subscription on this connection. */
    public static S7Message buildUnsubscribeRequest(int tpduRef, ControllerType controllerType) {
        // Use the abort variant of the same alarm-state type we subscribed with.
        AlarmStateType abortType = controllerType == ControllerType.S7_400
            ? AlarmStateType.ALARM_ABORT : AlarmStateType.ALARM_S_ABORT;
        return buildMsgSubscriptionMessage(tpduRef, /*dataLength=*/12, SUBSCRIPTION_ALM,
            abortType, /*reserve=*/(short) 0);
    }

    private static AlarmStateType alarmStateForController(ControllerType controllerType) {
        return controllerType == ControllerType.S7_400
            ? AlarmStateType.ALARM_INITIATE : AlarmStateType.ALARM_S_INITIATE;
    }

    private static S7Message buildMsgSubscriptionMessage(int tpduRef, int dataLength,
                                                         short subscription,
                                                         AlarmStateType alarmType, Short reserve) {
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,    // method: request
            (byte)  0x04,    // cpuFunctionType: request
            (byte)  0x04,    // cpuFunctionGroup: CPU functions
            (short) 0x02,    // cpuSubfunction: MsgSubscription
            (short) 0x00,    // sequenceNumber
            null, null, null);
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest(
            org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode.OK,
            org.apache.plc4x.java.s7.readwrite.DataTransportSize.OCTET_STRING,
            dataLength,
            subscription,
            MAGIC_KEY,
            alarmType,
            reserve);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /**
     * Parse a MsgSubscription response. Returns {@code true} iff the PLC accepted the
     * subscription (or reports an idempotent "already-initialized" state — see below).
     *
     * <p>Result-byte semantics observed against real S7-300 firmware (CPU 315-2 PN/DP):
     * <ul>
     *   <li>{@code 0x00} — subscription accepted, alarm-state initialized.</li>
     *   <li>{@code 0x02} — alarm subsystem is already in the requested state (benign on
     *       S7-300, which auto-initializes ALARM_S at CPU start). HMI clients also treat
     *       this as success — pushes still flow once user code calls SFC17/18.</li>
     * </ul>
     * Any other non-zero value is treated as a real rejection.
     */
    public static boolean parseSubscribeResponse(S7Message response) {
        if (!(response instanceof S7MessageUserData)) {
            return false;
        }
        if (hasParameterErrorCode(response)) {
            return false;
        }
        S7PayloadUserData payload = (S7PayloadUserData) response.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (item instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse alarm) {
                int result = alarm.getResult() & 0xFF;
                return result == 0x00 || result == 0x02;
            }
            if (item instanceof S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse) {
                // Plain "ok" response with no body — treat as success.
                return true;
            }
        }
        return false;
    }

    /**
     * Convert an unsolicited alarm-indication push into a {@link PlcStruct} payload.
     * Returns {@code null} if the message isn't a recognised alarm indication. The struct
     * is keyed by name and contains the PLC-side timestamp, the event id, and the four
     * State bytes (eventState / localState / ackStateGoing / ackStateComing) decoded into
     * an 8-bit signal mask each.
     *
     * <p>NOTIFY / NOTIFY_8 / ALARM_ACK_IND share the {@link AlarmMessagePushType} layout, so
     * we route them through the same parser — the calling code can disambiguate via the
     * subfunction recovered from the parameter header if needed.
     */
    /**
     * Build a one-shot AlarmQuery request that asks the PLC for currently-buffered alarms
     * of the given type. The response carries a raw byte payload (a list of
     * {@code AlarmMessageObjectQueryType} records) that callers can inspect via
     * {@link #parseAlarmQueryResponse(S7Message)}.
     *
     * @param queryType {@link QueryType#ALARM_S} for ALARM_S buffered alarms (S7-300/1500),
     *                  {@link QueryType#ALARM_8} for ALARM_8 (S7-400).
     */
    public static S7Message buildAlarmQueryRequest(int tpduRef, QueryType queryType) {
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,    // method: request
            (byte)  0x04,    // cpuFunctionType: request
            (byte)  0x04,    // cpuFunctionGroup: CPU functions
            (short) 0x13,    // cpuSubfunction: AlarmQuery
            (short) 0x00,    // sequenceNumber
            null, null, null);
        AlarmType alarmType = (queryType == QueryType.ALARM_8) ? AlarmType.ALARM_8 : AlarmType.ALARM_S;
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCpuFunctionAlarmQueryRequest(
            org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode.OK,
            org.apache.plc4x.java.s7.readwrite.DataTransportSize.OCTET_STRING,
            // dataLength: 4 const bytes + 5 variable bytes (syntaxId, reserved, queryType,
            //             reserved, alarmType) = 9. SPI1 used 12 here and it worked on real
            //             S7-300; the extra 3 bytes are absorbed by the framing layer.
            12,
            SyntaxIdType.ALARM_QUERYREQSET,
            queryType,
            alarmType);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /**
     * Extract the raw byte payload of an AlarmQuery response. Returns {@code null} if the
     * message isn't an AlarmQueryResponse or the PLC reported a non-zero error code.
     */
    public static byte[] parseAlarmQueryResponse(S7Message message) {
        if (!(message instanceof S7MessageUserData userData)) {
            return null;
        }
        if (hasParameterErrorCode(message)) {
            return null;
        }
        if (!(userData.getPayload() instanceof S7PayloadUserData payload)) {
            return null;
        }
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (item instanceof S7PayloadUserDataItemCpuFunctionAlarmQueryResponse resp) {
                byte[] raw = resp.getItems();
                return raw == null ? new byte[0] : raw;
            }
        }
        return null;
    }

    public static PlcValue parseAlarmIndication(S7Message message) {
        if (!(message instanceof S7MessageUserData userData)) {
            return null;
        }
        if (!(userData.getPayload() instanceof S7PayloadUserData payload)) {
            return null;
        }
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (item instanceof S7PayloadAlarmS  s)        return fromPush(s.getAlarmMessage());
            if (item instanceof S7PayloadAlarmSQ sq)       return fromPush(sq.getAlarmMessage());
            if (item instanceof S7PayloadAlarmSC sc)       return fromPush(sc.getAlarmMessage());
            if (item instanceof S7PayloadAlarm8  a8)       return fromPush(a8.getAlarmMessage());
            if (item instanceof S7PayloadNotify  n)        return fromPush(n.getAlarmMessage());
            if (item instanceof S7PayloadNotify8 n8)       return fromPush(n8.getAlarmMessage());
            // ALARM_ACK_IND uses the AlarmMessageAckPushType (no event/local state, only ack states).
            if (item instanceof S7PayloadAlarmAckInd ack)  return fromAckPush(ack.getAlarmMessage());
        }
        return null;
    }

    private static PlcValue fromPush(AlarmMessagePushType push) {
        if (push == null) {
            return null;
        }
        Map<String, PlcValue> children = new HashMap<>();
        addHeader(children, push.getTimeStamp(), push.getFunctionId(), push.getNumberOfObjects());
        List<AlarmMessageObjectPushType> objects = push.getMessageObjects() != null
            ? push.getMessageObjects() : Collections.emptyList();
        if (!objects.isEmpty()) {
            // Surface the first object directly — the common case is exactly one object per
            // indication. Multi-object indications can be inspected by callers via the raw
            // numberOfObjects field; future revisions can expose them as a list when needed.
            AlarmMessageObjectPushType first = objects.get(0);
            children.put("eventId",        new PlcUDINT(first.getEventId()));
            children.put("eventState",     new PlcUDINT(stateToMask(first.getEventState())));
            children.put("localState",     new PlcUDINT(stateToMask(first.getLocalState())));
            children.put("ackStateGoing",  new PlcUDINT(stateToMask(first.getAckStateGoing())));
            children.put("ackStateComing", new PlcUDINT(stateToMask(first.getAckStateComing())));
        }
        children.put("receivedAt", new PlcLINT(Instant.now().toEpochMilli()));
        return new PlcStruct(children);
    }

    private static PlcValue fromAckPush(AlarmMessageAckPushType push) {
        if (push == null) {
            return null;
        }
        Map<String, PlcValue> children = new HashMap<>();
        addHeader(children, push.getTimeStamp(), push.getFunctionId(), push.getNumberOfObjects());
        List<AlarmMessageAckObjectPushType> objects = push.getMessageObjects() != null
            ? push.getMessageObjects() : Collections.emptyList();
        if (!objects.isEmpty()) {
            AlarmMessageAckObjectPushType first = objects.get(0);
            children.put("eventId",        new PlcUDINT(first.getEventId()));
            children.put("ackStateGoing",  new PlcUDINT(stateToMask(first.getAckStateGoing())));
            children.put("ackStateComing", new PlcUDINT(stateToMask(first.getAckStateComing())));
        }
        children.put("receivedAt", new PlcLINT(Instant.now().toEpochMilli()));
        return new PlcStruct(children);
    }

    private static void addHeader(Map<String, PlcValue> out,
                                  org.apache.plc4x.java.s7.readwrite.DateAndTime ts,
                                  short functionId, short numberOfObjects) {
        out.put("plcTimestamp",    new PlcSTRING(formatDateAndTime(ts)));
        out.put("functionId",      new PlcUDINT(functionId & 0xFF));
        out.put("numberOfObjects", new PlcUDINT(numberOfObjects & 0xFF));
    }

    private static int stateToMask(State state) {
        if (state == null) return 0;
        int m = 0;
        if (state.getSIG_1()) m |= 0x01;
        if (state.getSIG_2()) m |= 0x02;
        if (state.getSIG_3()) m |= 0x04;
        if (state.getSIG_4()) m |= 0x08;
        if (state.getSIG_5()) m |= 0x10;
        if (state.getSIG_6()) m |= 0x20;
        if (state.getSIG_7()) m |= 0x40;
        if (state.getSIG_8()) m |= 0x80;
        return m;
    }

    private static String formatDateAndTime(org.apache.plc4x.java.s7.readwrite.DateAndTime dt) {
        if (dt == null) return "";
        // BCD year is 2-digit; S7 convention: < 90 → 2000s, ≥ 90 → 1900s.
        int year = dt.getYear() & 0xFF;
        year += (year < 90) ? 2000 : 1900;
        return String.format("%04d-%02d-%02dT%02d:%02d:%02d.%03d",
            year, dt.getMonth() & 0xFF, dt.getDay() & 0xFF,
            dt.getHour() & 0xFF, dt.getMinutes() & 0xFF, dt.getSeconds() & 0xFF,
            dt.getMsec() & 0xFFF);
    }

    private static boolean hasParameterErrorCode(S7Message response) {
        if (!(response.getParameter() instanceof S7ParameterUserData param)) {
            return false;
        }
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu
                && cpu.getErrorCode() != null && cpu.getErrorCode() != 0) {
                return true;
            }
        }
        return false;
    }
}
