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

import org.apache.plc4x.java.s7.readwrite.AssociatedValueType;
import org.apache.plc4x.java.s7.readwrite.CycServiceItemAnyType;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesPush;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesUnsubscribeRequest;
import org.apache.plc4x.java.s7.readwrite.TimeBase;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.junit.jupiter.api.Test;

import java.time.Duration;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

class S7CyclicSubscriptionServiceTest {

    @Test
    void pickInterval_oneSecond_prefersB1SecOverB01Sec() {
        // Critical: S7-300 firmware (≥V3.2) rejects B01SEC=100ms base requests with errCode
        // 0xD804. A 1-second request must produce (B1SEC, 1) on the wire, not (B01SEC, 10),
        // even though both produce the same effective cadence mathematically.
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofSeconds(1));
        assertEquals(TimeBase.B1SEC, s.base());
        assertEquals(1, s.factor());
    }

    @Test
    void pickInterval_subSecond_fallsBackToB01Sec() {
        // 500ms can only be expressed via the 100ms base. May fail on strict firmware,
        // but that's the only way to ask for sub-second cadence.
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofMillis(500));
        assertEquals(TimeBase.B01SEC, s.base());
        assertEquals(5, s.factor());
    }

    @Test
    void pickInterval_60Sec_prefersCoarsestBase() {
        // 60s is a clean multiple of 10s — prefer (B10SEC, 6) over (B1SEC, 60). Coarser base
        // is less load on the PLC's cyclic scheduler and avoids unnecessary precision.
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofSeconds(60));
        assertEquals(TimeBase.B10SEC, s.base());
        assertEquals(6, s.factor());
        assertEquals(60_000L, s.toMillis());
    }

    @Test
    void pickInterval_2Sec_usesB1SecBase() {
        // Not a multiple of 10s — must use B1SEC.
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofSeconds(2));
        assertEquals(TimeBase.B1SEC, s.base());
        assertEquals(2, s.factor());
    }

    @Test
    void pickInterval_over255Sec_usesB10SecBase() {
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofSeconds(600));
        assertEquals(TimeBase.B10SEC, s.base());
        assertEquals(60, s.factor());
    }

    @Test
    void pickInterval_nearSecond_roundsToB1SecRatherThanB01Sec() {
        // 1100ms is close enough to a whole second — prefer (B1SEC, 1) over (B01SEC, 11)
        // because B01SEC may be rejected by firmware. Caller's intent is "around 1 second".
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofMillis(1100));
        assertEquals(TimeBase.B1SEC, s.base());
    }

    @Test
    void pickInterval_below100ms_clampsTo100ms() {
        S7CyclicSubscriptionService.IntervalSpec s =
            S7CyclicSubscriptionService.pickInterval(Duration.ofMillis(10));
        assertEquals(100L, s.toMillis());
    }

    @Test
    void pickInterval_nullDuration_defaultsTo1Sec() {
        S7CyclicSubscriptionService.IntervalSpec s = S7CyclicSubscriptionService.pickInterval(null);
        assertEquals(1_000L, s.toMillis());
        assertEquals(TimeBase.B1SEC, s.base());
    }

    @Test
    void buildSubscribeRequest_setsCorrectDiscriminatorsAndItem() {
        S7Tag tag = new S7Tag(TransportSize.WORD, MemoryArea.FLAGS_MARKERS, 0, 10, (byte) 0, 1);
        S7CyclicSubscriptionService.IntervalSpec interval = new S7CyclicSubscriptionService.IntervalSpec(
            TimeBase.B1SEC, (short) 1);
        S7Message msg = S7CyclicSubscriptionService.buildSubscribeRequest(7, List.of(tag), interval);

        S7ParameterUserDataItemCPUFunctions cpu = (S7ParameterUserDataItemCPUFunctions)
            ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
        assertEquals(0x02, cpu.getCpuFunctionGroup());
        assertEquals(0x04, cpu.getCpuFunctionType());
        assertEquals(0x01, cpu.getCpuSubfunction());

        S7PayloadUserDataItemCyclicServicesSubscribeRequest req =
            (S7PayloadUserDataItemCyclicServicesSubscribeRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        assertEquals(1, req.getItemsCount());
        assertEquals(TimeBase.B1SEC, req.getTimeBase());
        assertEquals(1, req.getTimeFactor());
        assertEquals(1, req.getItem().size());

        CycServiceItemAnyType item = (CycServiceItemAnyType) req.getItem().get(0);
        assertEquals(TransportSize.WORD, item.getTransportSize());
        assertEquals(1, item.getLength());
        assertEquals(0, item.getDbNumber());
        assertEquals(MemoryArea.FLAGS_MARKERS, item.getMemoryArea());
        // address = (byteOffset << 3) | bitOffset = (10 << 3) | 0 = 80
        assertEquals(80, item.getAddress());
    }

    @Test
    void buildUnsubscribeRequest_setsCorrectDiscriminatorsAndJobId() {
        S7Message msg = S7CyclicSubscriptionService.buildUnsubscribeRequest(3, (short) 0x42);
        S7ParameterUserDataItemCPUFunctions cpu = (S7ParameterUserDataItemCPUFunctions)
            ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
        assertEquals(0x02, cpu.getCpuFunctionGroup());
        assertEquals(0x04, cpu.getCpuFunctionType());
        assertEquals(0x04, cpu.getCpuSubfunction());

        S7PayloadUserDataItemCyclicServicesUnsubscribeRequest req =
            (S7PayloadUserDataItemCyclicServicesUnsubscribeRequest)
                ((S7PayloadUserData) msg.getPayload()).getItems().get(0);
        assertEquals(0x05, req.getFunction());    // 0x05 = unsubscribe single job
        assertEquals(0x42, req.getJobId());
    }

    @Test
    void parseSubscribeResponse_returnsSequenceNumberAsJobId() {
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCyclicServicesSubscribeResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0,
                /*itemsCount=*/0, Collections.emptyList()),
            /*errorCode=*/null,
            /*sequenceNumber=*/(short) 0x77);
        Short jobId = S7CyclicSubscriptionService.parseSubscribeResponse(response);
        assertNotNull(jobId);
        assertEquals((short) 0x77, jobId);
    }

    @Test
    void parseSubscribeResponse_returnsNullOnParameterErrorCode() {
        S7Message response = wrapInUserData(
            new S7PayloadUserDataItemCyclicServicesSubscribeResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0,
                0, Collections.emptyList()),
            /*errorCode=*/0xD403, (short) 0x77);
        assertNull(S7CyclicSubscriptionService.parseSubscribeResponse(response));
    }

    @Test
    void parsePushItems_extractsRawByteSequencesPerEntry() {
        AssociatedValueType v1 = new AssociatedValueType(
            DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, 16,
            List.of((short) 0xCA, (short) 0xFE));
        AssociatedValueType v2 = new AssociatedValueType(
            DataTransportErrorCode.OK, DataTransportSize.BYTE_WORD_DWORD, 16,
            List.of((short) 0xBA, (short) 0xBE));
        S7Message msg = wrapInUserData(
            new S7PayloadUserDataItemCyclicServicesPush(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0,
                2, List.of(v1, v2)),
            null, (short) 0x10);
        List<byte[]> items = S7CyclicSubscriptionService.parsePushItems(msg);
        assertNotNull(items);
        assertEquals(2, items.size());
        assertArrayEquals(new byte[] { (byte) 0xCA, (byte) 0xFE }, items.get(0));
        assertArrayEquals(new byte[] { (byte) 0xBA, (byte) 0xBE }, items.get(1));
    }

    @Test
    void getJobId_readsFromCpuFunctionsParameterSequenceNumber() {
        S7Message msg = wrapInUserData(
            new S7PayloadUserDataItemCyclicServicesSubscribeResponse(
                DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 0,
                0, Collections.emptyList()),
            null, (short) 0x99);
        assertEquals((short) 0x99, S7CyclicSubscriptionService.getJobId(msg));
    }

    private static S7Message wrapInUserData(S7PayloadUserDataItem payloadItem, Integer errorCode, short seqNumber) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x02, (short) 0x01,
            seqNumber, null, null, errorCode);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.singletonList(payloadItem)));
    }
}
