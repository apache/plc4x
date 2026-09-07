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
import org.apache.plc4x.java.s7.readwrite.CycServiceItemType;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesPush;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeEmptyResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesSubscribeResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesUnsubscribeRequest;
import org.apache.plc4x.java.s7.readwrite.TimeBase;
import org.apache.plc4x.java.s7.tag.S7Tag;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * S7Comm UserData "Cyclic services" — has the PLC push variable values at a fixed
 * interval rather than us polling for them.
 *
 * <p>Wire shape:
 * <ul>
 *   <li>Subscribe request: {@code cpuFunctionGroup=0x02, cpuFunctionType=0x04, cpuSubfunction=0x01},
 *       carrying an itemsCount + {@link TimeBase} + timeFactor + a list of {@link CycServiceItemType}
 *       (one entry per subscribed variable). Interval = {@code timeBase × timeFactor}.</li>
 *   <li>Subscribe response: same group/subfn with type=0x08; the response's parameter
 *       {@code sequenceNumber} is the {@code jobId} the PLC has assigned to this subscription.
 *       Subsequent pushes and the unsubscribe request key off this id.</li>
 *   <li>Push: {@code group=0x02, type=0x00, subfn=0x01} with the same item ordering as
 *       the original subscribe request — caller decodes each {@link AssociatedValueType}
 *       payload using the corresponding tag's data type.</li>
 *   <li>Unsubscribe: {@code group=0x02, subfn=0x04} with {@code function=0x05} and the
 *       {@code jobId} from the original subscribe response.</li>
 * </ul>
 */
public final class S7CyclicSubscriptionService {

    public static final int CPU_FUNCTION_GROUP_CYCLIC = 0x02;
    public static final int CPU_FUNCTION_TYPE_PUSH    = 0x00;
    public static final int CYCLIC_PUSH_SUBFUNCTION   = 0x01;

    /** "Function" byte the PLC expects in an unsubscribe request — 0x05 = unsubscribe single job. */
    private static final short UNSUBSCRIBE_SINGLE_JOB = 0x05;

    /** Variable-spec syntax id for "any address". 0xb0 (DB-block read) is unsupported here. */
    private static final short SYNTAX_ID_ANY = 0x10;

    private S7CyclicSubscriptionService() {
    }

    /** Choice of {@link TimeBase} + timeFactor that approximates the requested interval. */
    public record IntervalSpec(TimeBase base, short factor) {
        /** Effective interval in milliseconds. */
        public long toMillis() {
            long baseMs = switch (base) {
                case B01SEC -> 100L;
                case B1SEC -> 1000L;
                case B10SEC -> 10_000L;
            };
            return baseMs * (factor & 0xFF);
        }
    }

    /**
     * Pick the {@link TimeBase} + factor pair for a requested interval. Prefers the
     * <em>coarsest</em> base whose factor still fits in a byte — empirically required for
     * S7-300 firmware (≥ V3.2) which rejects {@code B01SEC} (100 ms base) requests with
     * {@code errCode=0xD804}, only honouring {@code B1SEC} and {@code B10SEC}. Sub-second
     * cadences still fall back to B01SEC because that's the only way to express them on
     * the wire — those will fail on the strict firmware variants and the caller should be
     * prepared to handle the rejection.
     *
     * <p>Range: 100 ms … 2550 s. Null/negative durations default to 1 s.
     */
    public static IntervalSpec pickInterval(Duration requested) {
        long ms = requested == null ? 1_000L : Math.max(100L, requested.toMillis());
        if (ms >= 10_000L && ms <= 2_550_000L && ms % 10_000L == 0L) {
            return new IntervalSpec(TimeBase.B10SEC, (short) (ms / 10_000L));
        }
        if (ms >= 1_000L && ms <= 255_000L && ms % 1_000L == 0L) {
            return new IntervalSpec(TimeBase.B1SEC, (short) (ms / 1_000L));
        }
        // Fall through to B1SEC for "near-second" durations rather than the strict-rejected
        // B01SEC base whenever rounding to seconds doesn't change the cadence by more than
        // 10%. This is the path that lets a Duration.ofSeconds(1) request work without the
        // caller knowing the firmware quirk.
        if (ms >= 900L && ms <= 255_000L) {
            return new IntervalSpec(TimeBase.B1SEC, (short) Math.max(1, Math.round(ms / 1_000.0)));
        }
        if (ms >= 9_000L && ms <= 2_550_000L) {
            return new IntervalSpec(TimeBase.B10SEC, (short) Math.max(1, Math.round(ms / 10_000.0)));
        }
        // Genuine sub-second request — B01SEC is the only base that can express it. Some
        // firmwares reject this with 0xD804; callers should treat that as "not supported".
        int factor = Math.max(1, (int) Math.round(ms / 100.0));
        return new IntervalSpec(TimeBase.B01SEC, (short) Math.min(255, factor));
    }

    /**
     * Build a "subscribe to cyclic pushes" request for the given list of tags. All tags share
     * the same interval; group on the caller side if you need different cadences (one wire
     * subscription per unique interval).
     */
    public static S7Message buildSubscribeRequest(int tpduRef, List<S7Tag> tags, IntervalSpec interval) {
        List<CycServiceItemType> items = new ArrayList<>(tags.size());
        for (S7Tag tag : tags) {
            items.add(toCycServiceItem(tag));
        }
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,    // method: request
            (byte)  0x04,    // cpuFunctionType: request
            (byte)  CPU_FUNCTION_GROUP_CYCLIC,
            (short) 0x01,    // cpuSubfunction: subscribe
            (short) 0x00,    // sequenceNumber (PLC assigns the real jobId in the response)
            null, null, null);
        // dataLength: itemsCount(2) + timeBase(1) + timeFactor(1) + sum of each item's bytes.
        // Each CycServiceItemAnyType is 12 bytes on the wire (functionId(1)+byteLength(1)+
        // syntaxId(1)+transportSize(1)+length(2)+dbNumber(2)+memoryArea(1)+address(3)).
        int payloadLen = 4 + (items.size() * 12);
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCyclicServicesSubscribeRequest(
            DataTransportErrorCode.OK,
            DataTransportSize.OCTET_STRING,
            payloadLen,
            items.size(),
            interval.base(),
            interval.factor(),
            items);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /** Build a request that cancels a single cyclic subscription identified by jobId. */
    public static S7Message buildUnsubscribeRequest(int tpduRef, short jobId) {
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,    // method: request
            (byte)  0x04,    // cpuFunctionType: request
            (byte)  CPU_FUNCTION_GROUP_CYCLIC,
            (short) 0x04,    // cpuSubfunction: unsubscribe
            (short) 0x00,    // sequenceNumber
            null, null, null);
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCyclicServicesUnsubscribeRequest(
            DataTransportErrorCode.OK,
            DataTransportSize.OCTET_STRING,
            /*dataLength=*/2,
            UNSUBSCRIBE_SINGLE_JOB,
            jobId);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /**
     * Read the {@code sequenceNumber} field from a UserData message's CPUFunctions parameter.
     * For a cyclic subscribe response this is the PLC-assigned {@code jobId}; for incoming
     * pushes it identifies which subscription the values belong to. Returns {@code null} if
     * the message doesn't carry a CPUFunctions parameter.
     */
    public static Short getJobId(S7Message message) {
        if (!(message instanceof S7MessageUserData userData)) {
            return null;
        }
        if (!(userData.getParameter() instanceof S7ParameterUserData param)) {
            return null;
        }
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu) {
                return cpu.getSequenceNumber();
            }
        }
        return null;
    }

    /**
     * Parse a cyclic-subscribe response. Returns the {@code jobId} (sequenceNumber from the
     * parameter header) if accepted, or {@code null} if the PLC reported a non-zero error
     * code in the parameter or the payload isn't a SubscribeResponse.
     */
    public static Short parseSubscribeResponse(S7Message response) {
        if (!(response instanceof S7MessageUserData userData)) {
            return null;
        }
        if (hasParameterErrorCode(response)) {
            return null;
        }
        S7PayloadUserData payload = (S7PayloadUserData) userData.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            // Two success shapes the PLC may return:
            //   - SubscribeResponse with an items list (current values for each subscribed tag)
            //   - SubscribeEmptyResponse (dataLength=0) when the PLC doesn't bundle initial
            //     values into the ack — pushes start arriving on the cycle anyway.
            if (item instanceof S7PayloadUserDataItemCyclicServicesSubscribeResponse
                || item instanceof S7PayloadUserDataItemCyclicServicesSubscribeEmptyResponse) {
                return getJobId(response);
            }
        }
        return null;
    }

    /**
     * Extract the per-item raw byte payloads from a cyclic push. Each entry corresponds, in
     * order, to one tag from the original subscribe request — the caller decodes each byte
     * array using that tag's data type. Returns {@code null} if the message isn't a recognised
     * cyclic push.
     */
    public static List<byte[]> parsePushItems(S7Message message) {
        if (!(message instanceof S7MessageUserData userData)) {
            return null;
        }
        if (!(userData.getPayload() instanceof S7PayloadUserData payload)) {
            return null;
        }
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (item instanceof S7PayloadUserDataItemCyclicServicesPush push) {
                List<AssociatedValueType> values = push.getItems() != null
                    ? push.getItems() : Collections.emptyList();
                List<byte[]> bytes = new ArrayList<>(values.size());
                for (AssociatedValueType v : values) {
                    bytes.add(toByteArray(v.getData()));
                }
                return bytes;
            }
        }
        return null;
    }

    private static CycServiceItemType toCycServiceItem(S7Tag tag) {
        // address is a 24-bit field encoding bit-addressed memory: byteOffset shifted left
        // by 3 plus the bit number. Same encoding as a regular S7AddressAny request.
        int address = (tag.getByteOffset() << 3) | (tag.getBitOffset() & 0x07);
        return new CycServiceItemAnyType(
            /*byteLength=*/(short) 0x0a,        // 10 bytes after byteLength field
            /*syntaxId=*/SYNTAX_ID_ANY,
            tag.getDataType(),
            tag.getNumberOfElements(),
            tag.getBlockNumber(),
            tag.getMemoryArea(),
            address);
    }

    private static byte[] toByteArray(List<Short> data) {
        if (data == null) return new byte[0];
        byte[] out = new byte[data.size()];
        for (int i = 0; i < out.length; i++) {
            out[i] = (byte) (data.get(i) & 0xFF);
        }
        return out;
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
