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
package org.apache.plc4x.java.s7.userdata;

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlResponse;
import org.apache.plc4x.java.s7.readwrite.SzlId;
import org.apache.plc4x.java.s7.readwrite.SzlModuleTypeClass;
import org.apache.plc4x.java.s7.readwrite.SzlSublist;

import java.nio.charset.StandardCharsets;
import java.util.Collections;

/**
 * S7Comm UserData "CPU functions / Read SZL" service.
 *
 * <p>SZL = Systemzustandsliste (System Status List). The probe issued at connect time
 * doubles as a capability check: if the device responds, we know UserData services
 * (alarms, cyclic, browse) are usable. LOGO and similar low-end CPUs typically don't
 * implement SZL, so a failed probe is the signal to stay in plain Read/Write var mode.
 */
public final class S7SzlService {

    /** Module identification SZL ID — works on S7-300/400. */
    public static final SzlId MODULE_IDENTIFICATION =
        new SzlId(SzlModuleTypeClass.CPU, (byte) 0x00, SzlSublist.MODULE_IDENTIFICATION);

    /**
     * Component identification SZL ID — the standard for S7-1200/1500. Combined with index
     * {@code 0x0001} the response carries the CPU's order number (MLFB) in ASCII.
     */
    public static final SzlId COMPONENT_IDENTIFICATION =
        new SzlId(SzlModuleTypeClass.CPU, (byte) 0x00, SzlSublist.COMPONENT_IDENTIFICATION);

    private S7SzlService() {
    }

    /**
     * Build the {@code S7MessageUserData} for an SZL read request. Caller assigns the
     * tpduReference and dispatches it through the connection.
     */
    public static S7Message buildRequest(int tpduRef, SzlId szlId) {
        return buildRequest(tpduRef, szlId, 0x0000);
    }

    public static S7Message buildRequest(int tpduRef, SzlId szlId, int szlIndex) {
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,    // method
            (byte)  0x04,    // cpuFunctionType: request
            (byte)  0x04,    // cpuFunctionGroup: CPU
            (short) 0x01,    // cpuSubfunction: ReadSZL
            (short) 0x00,    // sequenceNumber
            null, null, null);
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCpuFunctionReadSzlRequest(
            DataTransportErrorCode.OK,
            DataTransportSize.OCTET_STRING,
            4,                // dataLength: SzlId(2) + SzlIndex(2). The Java SPI1 driver had
                              // 0x0C here; that's incorrect — see plc4go for the right value.
                              // PLCs reply with garbage payloads when the declared length is
                              // longer than what's actually on the wire.
            szlId,
            szlIndex);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /** Probe result: identifying information and derived capability flags. */
    public record ProbeResult(String articleNumber, ControllerType controllerType) {
    }

    /**
     * Diagnostic helper: log the response's UserData CPU-function parameter (so we can tell
     * whether it really is a ReadSZL response or an error reply) plus the raw payload bytes.
     * <p>
     * Left in in in case we need it for diagnosing issues in the future.
     */
    public static void dumpRawSzlPayload(S7Message response, org.slf4j.Logger logger) {
        if (!(response instanceof S7MessageUserData)) {
            logger.info("SZL probe got non-UserData response: {}",
                response == null ? "null" : response.getClass().getSimpleName());
            return;
        }
        S7ParameterUserData param = (S7ParameterUserData) response.getParameter();
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu) {
                logger.info("SZL probe CPU params: method=0x{}, type=0x{}, group=0x{}, subfn=0x{}, errCode={}",
                    Integer.toHexString(cpu.getMethod()  & 0xFF),
                    Integer.toHexString(cpu.getCpuFunctionType()  & 0xFF),
                    Integer.toHexString(cpu.getCpuFunctionGroup() & 0xFF),
                    Integer.toHexString(cpu.getCpuSubfunction() & 0xFFFF),
                    cpu.getErrorCode());
            }
        }
        S7PayloadUserData payload = (S7PayloadUserData) response.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            byte[] bytes = item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse szl
                ? szl.getItems()
                : null;
            logger.info("SZL probe payload item: type={}, bytes={}",
                item.getClass().getSimpleName(),
                bytes == null ? "<not a ReadSZL response>" : (bytes.length + " bytes"));
            if (bytes != null) {
                StringBuilder sb = new StringBuilder();
                for (byte b : bytes) sb.append(String.format("%02X ", b & 0xFF));
                logger.info("SZL raw items: {}", sb);
            }
        }
    }

    /**
     * Parse an SZL response (MODULE_IDENTIFICATION or COMPONENT_IDENTIFICATION) and pull the
     * article number out of it.
     *
     * <p>The data section of either response is a header (8 bytes — SZL-ID, SZL-INDEX,
     * item-length, item-count) followed by N items of {@code item-length} bytes. For
     * MODULE_IDENTIFICATION the first item starts with a 2-byte index then 20 bytes of MLFB
     * in ASCII. For COMPONENT_IDENTIFICATION (used on S7-1200/1500) item-index 1 carries the
     * CPU order number as 32 ASCII bytes after the 2-byte index. We scan the items array
     * looking for the first chunk that decodes as printable ASCII starting with {@code "6ES7"}
     * — robust across families without needing a separate parser per SZL ID.
     */
    public static ProbeResult parseProbeResponse(S7Message response) {
        if (!(response instanceof S7MessageUserData)) {
            throw new IllegalArgumentException("Expected S7MessageUserData, got "
                + response.getClass().getSimpleName());
        }
        // The PLC reports unsupported SZL-ID/index combinations via the parameter's errorCode.
        // Treat any non-zero error as "this SZL didn't work" and let the caller try the next one.
        S7ParameterUserData param = (S7ParameterUserData) response.getParameter();
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu
                && cpu.getErrorCode() != null && cpu.getErrorCode() != 0) {
                throw new IllegalStateException("SZL request rejected by PLC, errorCode=0x"
                    + Integer.toHexString(cpu.getErrorCode()));
            }
        }
        S7PayloadUserData payload = (S7PayloadUserData) response.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (!(item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse szl)) {
                continue;
            }
            byte[] bytes = szl.getItems();
            if (bytes == null || bytes.length < 4) {
                throw new IllegalStateException("SZL response too short: "
                    + (bytes == null ? -1 : bytes.length) + " bytes");
            }
            // Best-effort identification: try the order-number prefix first, then the
            // model-name pattern (e.g. "CPU 315-2 PN/DP" on S7-300). If neither hits we
            // still return an empty-string article — the response was structurally valid,
            // which is enough to say UserData services work on this device.
            String article = findArticleNumber(bytes);
            String descriptor = article != null ? article : findCpuModelName(bytes);
            return new ProbeResult(
                descriptor != null ? descriptor : "",
                decodeControllerType(descriptor));
        }
        throw new IllegalStateException("No SZL response item in payload");
    }

    /**
     * Look for a {@code "CPU "} prefix followed by a model number (e.g. {@code "CPU 315-2 PN/DP"},
     * {@code "CPU 1212C"}). S7-300/400 expose their identification this way via
     * {@code COMPONENT_IDENTIFICATION} item 7; S7-1200/1500 sometimes do too. Returns
     * {@code null} if nothing recognisable is found.
     */
    static String findCpuModelName(byte[] bytes) {
        for (int start = 0; start + 6 <= bytes.length; start++) {
            if (bytes[start] != 'C' || bytes[start + 1] != 'P' || bytes[start + 2] != 'U'
                || bytes[start + 3] != ' ') {
                continue;
            }
            int end = start;
            int max = Math.min(start + 32, bytes.length);
            while (end < max && bytes[end] >= 0x20 && bytes[end] < 0x7F) {
                end++;
            }
            String candidate = new String(bytes, start, end - start, StandardCharsets.US_ASCII).trim();
            if (candidate.length() > 4) {
                return candidate;
            }
        }
        return null;
    }

    /**
     * Scan {@code bytes} for the first ASCII window that starts with a known Siemens
     * order-number prefix (currently {@code 6ES7}) and is at least 12 chars long.
     * Returns {@code null} if no such window exists.
     */
    static String findArticleNumber(byte[] bytes) {
        for (int start = 0; start + 12 <= bytes.length; start++) {
            if (bytes[start] != '6') continue;
            // Read up to 20 bytes from this position, stopping at the first non-printable byte.
            int end = start;
            int max = Math.min(start + 32, bytes.length);
            while (end < max && bytes[end] >= 0x20 && bytes[end] < 0x7F) {
                end++;
            }
            if (end - start < 12) continue;
            String candidate = new String(bytes, start, end - start, StandardCharsets.US_ASCII).trim();
            if (candidate.startsWith("6ES7") || candidate.startsWith("6GK")) {
                return candidate;
            }
        }
        return null;
    }

    /**
     * Map a Siemens MLFB article number to a known {@link ControllerType}. Accepts both the
     * space-separated form {@code "6ES7 212-…"} (S7-300/400 SZL output) and the compact form
     * {@code "6ES7212-…"} (S7-1200/1500 SZL-COMPONENT-IDENTIFICATION output). Returns
     * {@code ANY} for unrecognised hardware.
     */
    static ControllerType decodeControllerType(String descriptor) {
        if (descriptor == null) {
            return ControllerType.ANY;
        }
        if (descriptor.startsWith("6ES7")) {
            // Order-number form, e.g. "6ES7 212-1AE40-0XB0". Family digit is the first
            // non-space character after the "6ES7" prefix.
            int i = 4;
            while (i < descriptor.length() && descriptor.charAt(i) == ' ') i++;
            if (i >= descriptor.length()) {
                return ControllerType.ANY;
            }
            return switch (descriptor.charAt(i)) {
                case '2' -> ControllerType.S7_1200;
                case '5' -> ControllerType.S7_1500;
                case '3' -> ControllerType.S7_300;
                case '4' -> ControllerType.S7_400;
                default -> ControllerType.ANY;
            };
        }
        if (descriptor.startsWith("CPU ")) {
            // Model-name form, e.g. "CPU 315-2 PN/DP" (S7-300), "CPU 1212C" (S7-1200).
            // The first digit after "CPU " selects the family; for 4-digit numbers the
            // second digit disambiguates 1200 vs 1500.
            int i = 4;
            while (i < descriptor.length() && descriptor.charAt(i) == ' ') i++;
            if (i >= descriptor.length() || !Character.isDigit(descriptor.charAt(i))) {
                return ControllerType.ANY;
            }
            int digits = 0;
            int n = 0;
            while (i + digits < descriptor.length() && Character.isDigit(descriptor.charAt(i + digits))) {
                n = n * 10 + (descriptor.charAt(i + digits) - '0');
                digits++;
                if (digits >= 4) break;
            }
            if (digits == 3) {
                return switch (descriptor.charAt(i)) {
                    case '3' -> ControllerType.S7_300;
                    case '4' -> ControllerType.S7_400;
                    default -> ControllerType.ANY;
                };
            }
            if (digits == 4) {
                if (n >= 1200 && n <= 1299) return ControllerType.S7_1200;
                if (n >= 1500 && n <= 1599) return ControllerType.S7_1500;
            }
        }
        return ControllerType.ANY;
    }
}
