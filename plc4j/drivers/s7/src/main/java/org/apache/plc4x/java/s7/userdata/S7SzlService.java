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
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

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

    /**
     * Protection level and mode-selector status — SZL-ID {@code 0x0232}, index {@code 0x0004}.
     * The ID is composed of type class {@code CPU}, partial-list extension {@code 0x02} and
     * partial list {@code 0x32}. Documented for S7-300/400; S7-1200/1500 frequently reject it,
     * so callers must treat a failure here as "not reported" rather than an error.
     */
    public static final SzlId PROTECTION_STATUS =
        new SzlId(SzlModuleTypeClass.CPU, (byte) 0x02, SzlSublist.COMMUNICATION_STATUS_DATA);

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

    /**
     * Build the follow-up request for a chained SZL response.
     *
     * <p>When a list doesn't fit in one PDU the CPU answers with {@code lastDataUnit} set and
     * keeps the sequence open. Until it is drained, every new SZL request on that connection
     * is refused with {@code 0xD406} — which is why this has to be sent rather than ignored.
     *
     * <p>The follow-up is the same ReadSZL request with the sequence number the PLC assigned
     * echoed back; verified against a CPU 315-2 PN/DP, which answers it with the next chunk.
     */
    public static S7Message buildContinuationRequest(int tpduRef, SzlId szlId, int szlIndex,
                                                     int sequenceNumber) {
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,               // method
            (byte)  0x04,               // cpuFunctionType: request
            (byte)  0x04,               // cpuFunctionGroup: CPU
            (short) 0x01,               // cpuSubfunction: ReadSZL
            (short) sequenceNumber,     // the sequence the PLC opened
            null, null, null);
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCpuFunctionReadSzlRequest(
            DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING, 4, szlId, szlIndex);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /**
     * Whether the PLC still has data pending for this SZL read, i.e. whether a continuation
     * request is needed before anything else can be asked of this connection.
     */
    public static boolean moreDataFollows(S7Message response) {
        Short lastDataUnit = cpuFunctionsOf(response, S7ParameterUserDataItemCPUFunctions::getLastDataUnit);
        return lastDataUnit != null && lastDataUnit != 0;
    }

    /** The sequence number the PLC assigned to a chained read, to be echoed back. */
    public static int sequenceNumberOf(S7Message response) {
        Short sequenceNumber = cpuFunctionsOf(response,
            S7ParameterUserDataItemCPUFunctions::getSequenceNumber);
        return sequenceNumber == null ? 0 : sequenceNumber;
    }

    private static Short cpuFunctionsOf(S7Message response,
                                        java.util.function.Function<S7ParameterUserDataItemCPUFunctions, Short> field) {
        if (!(response instanceof S7MessageUserData)
            || !(response.getParameter() instanceof S7ParameterUserData param)) {
            return null;
        }
        for (S7ParameterUserDataItem item : param.getItems()) {
            if (item instanceof S7ParameterUserDataItemCPUFunctions cpu) {
                return field.apply(cpu);
            }
        }
        return null;
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
    public static ProbeResult parseProbeResponse(byte[] bytes) {
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

    // ------------------------------------------------------------------------
    // Full device identification
    // ------------------------------------------------------------------------

    /**
     * Everything we can learn about a device from the identification SZLs. Every field is
     * nullable: devices answer different subsets depending on family and firmware, and a
     * field stays {@code null} when the owning SZL was not answered or did not carry it.
     */
    public record S7DeviceIdentification(String orderCode,
                                         String hardwareVersion,
                                         String firmwareVersion,
                                         String moduleName,
                                         String moduleTypeName,
                                         String systemName,
                                         String plantDesignation,
                                         String serialNumber,
                                         Integer protectionLevel,
                                         String keySwitchPosition,
                                         Boolean passwordProtected) {

        public static final S7DeviceIdentification EMPTY = new S7DeviceIdentification(
            null, null, null, null, null, null, null, null, null, null, null);

        /**
         * Combine two partial identifications. Fields set on {@code this} win; {@code null}
         * fields are filled in from {@code other}. Used to fold the three SZL responses
         * (module, component, protection) into a single result.
         */
        public S7DeviceIdentification mergedWith(S7DeviceIdentification other) {
            if (other == null) {
                return this;
            }
            return new S7DeviceIdentification(
                orderCode          != null ? orderCode          : other.orderCode,
                hardwareVersion    != null ? hardwareVersion    : other.hardwareVersion,
                firmwareVersion    != null ? firmwareVersion    : other.firmwareVersion,
                moduleName         != null ? moduleName         : other.moduleName,
                moduleTypeName     != null ? moduleTypeName     : other.moduleTypeName,
                systemName         != null ? systemName         : other.systemName,
                plantDesignation   != null ? plantDesignation   : other.plantDesignation,
                serialNumber       != null ? serialNumber       : other.serialNumber,
                protectionLevel    != null ? protectionLevel    : other.protectionLevel,
                keySwitchPosition  != null ? keySwitchPosition  : other.keySwitchPosition,
                passwordProtected  != null ? passwordProtected  : other.passwordProtected);
        }
    }

    /** An SZL response split into its header and fixed-length records. */
    record SzlList(int szlId, int szlIndex, int recordLength, List<byte[]> records) {
    }

    /**
     * Split an SZL response into records. The data section is an 8-byte header
     * ({@code szl-id}, {@code szl-index}, {@code record-length}, {@code record-count})
     * followed by {@code record-count} records of {@code record-length} bytes each.
     *
     * <p>Records the PLC announced but did not actually send are skipped rather than causing
     * a failure — some devices cut the response short at the PDU boundary.
     */
    static SzlList parseRecords(byte[] data) {
        if (data.length < 8) {
            throw new IllegalStateException("SZL response too short for a data header: "
                + data.length + " bytes");
        }
        int szlId = readWord(data, 0);
        int szlIndex = readWord(data, 2);
        int recordLength = readWord(data, 4);
        int recordCount = readWord(data, 6);
        if (recordLength <= 0) {
            throw new IllegalStateException("SZL response declares record length " + recordLength);
        }
        List<byte[]> records = new ArrayList<>();
        for (int i = 0; i < recordCount; i++) {
            int offset = 8 + (i * recordLength);
            if (offset + recordLength > data.length) {
                break;
            }
            records.add(Arrays.copyOfRange(data, offset, offset + recordLength));
        }
        return new SzlList(szlId, szlIndex, recordLength, records);
    }

    /**
     * Parse a module-identification response (SZL {@code 0x0011}, read with index
     * {@code 0x0000} so the PLC returns every record).
     *
     * <p>Each record is {@code index(2) + MLFB(20) + BGTyp(2) + Ausbg(2) + Ausbe(2)}. Record
     * index 1 identifies the module itself, 6 the basic hardware and 7 the basic firmware.
     * The version is the record's trailing three bytes — the reading Snap7 and nmap's
     * {@code s7-info} both use, and the one that matches what TIA Portal displays.
     */
    public static S7DeviceIdentification parseModuleIdentification(byte[] szlData) {
        String orderCode = null;
        String hardwareVersion = null;
        String firmwareVersion = null;
        for (byte[] record : parseRecords(szlData).records()) {
            switch (readWord(record, 0)) {
                case 0x0001 -> orderCode = emptyToNull(readAscii(record, 2, 20));
                case 0x0006 -> hardwareVersion = readVersion(record);
                case 0x0007 -> firmwareVersion = readVersion(record);
                default -> { /* other records identify submodules — not part of the inventory */ }
            }
        }
        return new S7DeviceIdentification(orderCode, hardwareVersion, firmwareVersion,
            null, null, null, null, null, null, null, null);
    }

    /**
     * Parse a component-identification response (SZL {@code 0x001C}, read with index
     * {@code 0x0000}). Each record is {@code index(2)} followed by an ASCII value; the record
     * indices we care about are 1 (system/station name), 2 (module name), 3 (plant
     * designation), 5 (serial number) and 7 (module type name, e.g. {@code "CPU 1511-1 PN"}).
     */
    public static S7DeviceIdentification parseComponentIdentification(byte[] szlData) {
        String systemName = null;
        String moduleName = null;
        String plantDesignation = null;
        String serialNumber = null;
        String moduleTypeName = null;
        for (byte[] record : parseRecords(szlData).records()) {
            String value = emptyToNull(readAscii(record, 2, record.length - 2));
            switch (readWord(record, 0)) {
                case 0x0001 -> systemName = value;
                case 0x0002 -> moduleName = value;
                case 0x0003 -> plantDesignation = value;
                case 0x0005 -> serialNumber = value;
                case 0x0007 -> moduleTypeName = value;
                default -> { /* copyright, OEM id, location designation — not inventoried */ }
            }
        }
        return new S7DeviceIdentification(null, null, null, moduleName, moduleTypeName,
            systemName, plantDesignation, serialNumber, null, null, null);
    }

    /**
     * Parse a protection-status response (SZL {@code 0x0232} index {@code 0x0004}). The record
     * is {@code index, sch_schal, sch_par, sch_rel, bart_sch, anl_sch} as 16-bit words.
     *
     * <p>{@code sch_rel} is the protection level actually in force (1 = read and write,
     * 2 = read-only without password, 3 = no access without password) and {@code sch_par} is
     * non-zero exactly when a protection level was configured, i.e. when a password is set.
     */
    public static S7DeviceIdentification parseProtectionStatus(byte[] szlData) {
        List<byte[]> records = parseRecords(szlData).records();
        if (records.isEmpty()) {
            throw new IllegalStateException("Protection-status SZL carried no records");
        }
        byte[] record = records.get(0);
        if (record.length < 10) {
            throw new IllegalStateException("Protection-status record too short: "
                + record.length + " bytes");
        }
        int parameterisedLevel = readWord(record, 4);   // sch_par
        int effectiveLevel = readWord(record, 6);       // sch_rel
        int keySwitch = readWord(record, 8);            // bart_sch
        return new S7DeviceIdentification(null, null, null, null, null, null, null, null,
            effectiveLevel, decodeKeySwitchPosition(keySwitch), parameterisedLevel != 0);
    }

    /**
     * Pull the raw SZL data section out of a response, after checking that the PLC didn't
     * reject the request. For a chained response this is one chunk; callers append the
     * chunks and hand the result to the parsers.
     */
    public static byte[] szlDataOf(S7Message response) {
        if (!(response instanceof S7MessageUserData)) {
            throw new IllegalArgumentException("Expected S7MessageUserData, got "
                + (response == null ? "null" : response.getClass().getSimpleName()));
        }
        checkErrorCode(response);
        S7PayloadUserData payload = (S7PayloadUserData) response.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (item instanceof S7PayloadUserDataItemCpuFunctionReadSzlResponse szl
                && szl.getItems() != null) {
                return szl.getItems();
            }
        }
        throw new IllegalStateException("No SZL response item in payload");
    }

    /**
     * The PLC reports unsupported SZL-ID/index combinations via the parameter's errorCode.
     * Treat any non-zero error as "this SZL didn't work" and let the caller try the next one.
     */
    private static void checkErrorCode(S7Message response) {
        S7ParameterUserData param = (S7ParameterUserData) response.getParameter();
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu
                && cpu.getErrorCode() != null && cpu.getErrorCode() != 0) {
                throw new IllegalStateException("SZL request rejected by PLC, errorCode=0x"
                    + Integer.toHexString(cpu.getErrorCode()));
            }
        }
    }

    private static String decodeKeySwitchPosition(int bartSch) {
        return switch (bartSch) {
            case 1 -> "RUN";
            case 2 -> "RUN-P";
            case 3 -> "STOP";
            case 4 -> "MRES";
            default -> "unknown";
        };
    }

    /** Version as the record's trailing three bytes, e.g. {@code "V 2.6.0"}. */
    private static String readVersion(byte[] record) {
        if (record.length < 3) {
            return null;
        }
        return String.format("V %d.%d.%d",
            record[record.length - 3] & 0xFF,
            record[record.length - 2] & 0xFF,
            record[record.length - 1] & 0xFF);
    }

    private static String readAscii(byte[] record, int offset, int length) {
        if (offset >= record.length || length <= 0) {
            return "";
        }
        int available = Math.min(length, record.length - offset);
        // Devices pad with either spaces or NULs; stop at the first NUL so the trailing
        // padding doesn't end up in the string.
        int end = offset;
        while (end < offset + available && record[end] != 0x00) {
            end++;
        }
        return new String(record, offset, end - offset, StandardCharsets.US_ASCII).trim();
    }

    private static String emptyToNull(String value) {
        return value == null || value.isEmpty() ? null : value;
    }

    private static int readWord(byte[] data, int offset) {
        if (offset + 1 >= data.length) {
            return 0;
        }
        return ((data[offset] & 0xFF) << 8) | (data[offset + 1] & 0xFF);
    }
}
