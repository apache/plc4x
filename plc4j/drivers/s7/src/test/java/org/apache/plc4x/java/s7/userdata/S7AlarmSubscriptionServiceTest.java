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
import org.apache.plc4x.java.s7.readwrite.AlarmMessageObjectPushType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessagePushType;
import org.apache.plc4x.java.s7.readwrite.AssociatedValueType;
import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.DateAndTime;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmS;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest;
import org.apache.plc4x.java.s7.readwrite.State;
import org.apache.plc4x.java.s7.readwrite.SyntaxIdType;
import org.junit.jupiter.api.Test;

import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class S7AlarmSubscriptionServiceTest {

    @Test
    void buildSubscribeRequest_armsAlarmSPathForS7300() {
        S7Message msg = S7AlarmSubscriptionService.buildSubscribeRequest(42, ControllerType.S7_300);
        assertInstanceOf(S7MessageUserData.class, msg);
        assertEquals(42, msg.getTpduReference());

        S7ParameterUserDataItemCPUFunctions cpu = (S7ParameterUserDataItemCPUFunctions)
            ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
        assertEquals(0x04, cpu.getCpuFunctionGroup());
        assertEquals(0x04, cpu.getCpuFunctionType());
        assertEquals(0x02, cpu.getCpuSubfunction());

        S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest req =
            (S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        // ALARM_S subscription: subscription byte = 0x80 (ALM event class), alarm-state appended (dataLength=12).
        assertEquals(S7AlarmSubscriptionService.SUBSCRIPTION_ALM, req.getSubscription());
        assertEquals(12, req.getDataLength());
        assertEquals(S7AlarmSubscriptionService.MAGIC_KEY, req.getMagicKey());
        assertEquals(org.apache.plc4x.java.s7.readwrite.AlarmStateType.ALARM_S_INITIATE,
            req.getAlarmtype());
    }

    @Test
    void buildSubscribeRequest_armsAlarm8PathForS7400() {
        S7Message msg = S7AlarmSubscriptionService.buildSubscribeRequest(1, ControllerType.S7_400);
        S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest req =
            (S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        assertEquals(org.apache.plc4x.java.s7.readwrite.AlarmStateType.ALARM_INITIATE,
            req.getAlarmtype());
    }

    @Test
    void buildUnsubscribeRequest_usesAbortVariant() {
        S7Message msg = S7AlarmSubscriptionService.buildUnsubscribeRequest(7, ControllerType.S7_300);
        S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest req =
            (S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        assertEquals(S7AlarmSubscriptionService.SUBSCRIPTION_ALM, req.getSubscription());
        assertEquals(12, req.getDataLength());
        assertEquals(org.apache.plc4x.java.s7.readwrite.AlarmStateType.ALARM_S_ABORT,
            req.getAlarmtype());
    }

    @Test
    void parseSubscribeResponse_acceptsResultZero() {
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 4,
                /*result=*/(short) 0, (short) 0, org.apache.plc4x.java.s7.readwrite.AlarmStateType.ALARM_S_INITIATE,
                (short) 0, (short) 0),
            /*errorCode=*/null);
        assertTrue(S7AlarmSubscriptionService.parseSubscribeResponse(response));
    }

    @Test
    void parseSubscribeResponse_acceptsResultTwoAsAlreadyInitialized() {
        // S7-300 returns result=0x02 ("already in this state") when ALARM_S is auto-initialized
        // at CPU start. Treated as success; pushes still flow.
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 4,
                /*result=*/(short) 2, (short) 0, null, (short) 0, (short) 0),
            null);
        assertTrue(S7AlarmSubscriptionService.parseSubscribeResponse(response));
    }

    @Test
    void parseSubscribeResponse_rejectsOtherNonZeroResult() {
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 4,
                /*result=*/(short) 1, (short) 0, org.apache.plc4x.java.s7.readwrite.AlarmStateType.ALARM_S_INITIATE,
                (short) 0, (short) 0),
            null);
        assertFalse(S7AlarmSubscriptionService.parseSubscribeResponse(response));
    }

    @Test
    void parseSubscribeResponse_rejectsParameterErrorCode() {
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 4,
                (short) 0, (short) 0, org.apache.plc4x.java.s7.readwrite.AlarmStateType.ALARM_S_INITIATE,
                (short) 0, (short) 0),
            /*errorCode=*/0xD401);
        assertFalse(S7AlarmSubscriptionService.parseSubscribeResponse(response));
    }

    @Test
    void parseAlarmIndication_extractsEventIdAndStateMasksFromAlarmS() {
        AlarmMessageObjectPushType obj = new AlarmMessageObjectPushType(
            /*lengthSpec=*/(short) 0,
            SyntaxIdType.ALARM_INDSET,
            /*numberOfValues=*/(short) 0,
            /*eventId=*/0xABCD1234L,
            /*eventState=*/state(0x05),     // SIG_1 + SIG_3
            /*localState=*/state(0x00),
            /*ackStateGoing=*/state(0xFF),
            /*ackStateComing=*/state(0x80),  // SIG_8
            /*associatedValues=*/Collections.<AssociatedValueType>emptyList());
        AlarmMessagePushType push = new AlarmMessagePushType(
            sampleTimestamp(),
            /*functionId=*/(short) 0x12,
            /*numberOfObjects=*/(short) 1,
            List.of(obj));
        S7Message message = wrapInUserData(
            new S7PayloadAlarmS(DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0, push),
            null);

        PlcValue value = S7AlarmSubscriptionService.parseAlarmIndication(message);
        assertNotNull(value);
        assertEquals(0xABCD1234L, value.getStruct().get("eventId").getLong());
        assertEquals(0x05, value.getStruct().get("eventState").getInt());
        assertEquals(0x00, value.getStruct().get("localState").getInt());
        assertEquals(0xFF, value.getStruct().get("ackStateGoing").getInt());
        assertEquals(0x80, value.getStruct().get("ackStateComing").getInt());
        assertEquals(0x12, value.getStruct().get("functionId").getInt());
    }

    @Test
    void buildSubscribeRequest_serializesCleanly() throws Exception {
        S7Message msg = S7AlarmSubscriptionService.buildSubscribeRequest(99, ControllerType.S7_300);
        // Build a buffer big enough to hold the message and try to serialize.
        org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased wb =
            new org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased(
                new byte[msg.getLengthInBytes()],
                org.apache.plc4x.java.spi.buffers.api.WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                org.apache.plc4x.java.spi.buffers.api.WithOption.WithSignedIntegerEncoding("twos-complement"),
                org.apache.plc4x.java.spi.buffers.api.WithOption.WithFloatEncoding("IEEE754"),
                org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption.WithByteOrder("BIG_ENDIAN"));
        msg.serialize(wb);
    }

    @Test
    void buildAlarmQueryRequest_setsAlarmQueryDiscriminators() {
        S7Message msg = S7AlarmSubscriptionService.buildAlarmQueryRequest(
            5, org.apache.plc4x.java.s7.readwrite.QueryType.ALARM_S);
        assertInstanceOf(S7MessageUserData.class, msg);
        assertEquals(5, msg.getTpduReference());

        S7ParameterUserDataItemCPUFunctions cpu = (S7ParameterUserDataItemCPUFunctions)
            ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
        assertEquals(0x04, cpu.getCpuFunctionGroup());
        assertEquals(0x04, cpu.getCpuFunctionType());
        assertEquals(0x13, cpu.getCpuSubfunction());

        org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionAlarmQueryRequest req =
            (org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionAlarmQueryRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        assertEquals(SyntaxIdType.ALARM_QUERYREQSET, req.getSyntaxId());
        assertEquals(org.apache.plc4x.java.s7.readwrite.QueryType.ALARM_S, req.getQueryType());
        assertEquals(org.apache.plc4x.java.s7.readwrite.AlarmType.ALARM_S, req.getAlarmType());
    }

    @Test
    void parseAlarmQueryResponse_returnsRawBytes() {
        byte[] raw = new byte[] { 0x01, 0x02, 0x03, 0x04 };
        S7Message response = wrapInUserData(
            new org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionAlarmQueryResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, raw.length, raw),
            null);
        byte[] payload = S7AlarmSubscriptionService.parseAlarmQueryResponse(response);
        assertNotNull(payload);
        assertArrayEquals(raw, payload);
    }

    @Test
    void parseAlarmQueryResponse_returnsNullOnParameterError() {
        S7Message response = wrapInUserData(
            new org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionAlarmQueryResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0, new byte[0]),
            /*errorCode=*/0xD403);
        assertNull(S7AlarmSubscriptionService.parseAlarmQueryResponse(response));
    }

    @Test
    void parseAlarmIndication_returnsNullForNonAlarmMessage() {
        S7Message msg = new S7MessageUserData(1,
            new S7ParameterUserData(Collections.emptyList()),
            new S7PayloadUserData(Collections.emptyList()));
        assertNull(S7AlarmSubscriptionService.parseAlarmIndication(msg));
    }

    private static State state(int mask) {
        return new State(
            (mask & 0x80) != 0, (mask & 0x40) != 0, (mask & 0x20) != 0, (mask & 0x10) != 0,
            (mask & 0x08) != 0, (mask & 0x04) != 0, (mask & 0x02) != 0, (mask & 0x01) != 0);
    }

    private static DateAndTime sampleTimestamp() {
        return new DateAndTime(
            (short) 0x24, (short) 0x12, (short) 0x31, (short) 0x12, (short) 0x34, (short) 0x56,
            (short) 123, (byte) 1);
    }

    private static S7Message wrapInUserData(S7PayloadUserDataItem payloadItem, Integer errorCode) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x02,
            (short) 0x00, null, null, errorCode);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.singletonList(payloadItem)));
    }
}
