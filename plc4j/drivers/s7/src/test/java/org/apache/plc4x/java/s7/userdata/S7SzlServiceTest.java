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

import org.apache.plc4x.java.s7.readwrite.ControllerType;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageRequest;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlResponse;
import org.apache.plc4x.java.s7.readwrite.SzlId;
import org.apache.plc4x.java.s7.readwrite.SzlModuleTypeClass;
import org.apache.plc4x.java.s7.readwrite.SzlSublist;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.Collections;
import java.util.concurrent.CompletionException;
import java.util.concurrent.TimeoutException;

import static org.junit.jupiter.api.Assertions.*;

class S7SzlServiceTest {

    @Test
    void buildRequest_producesUserDataMessage() {
        S7Message msg = S7SzlService.buildRequest(42, S7SzlService.MODULE_IDENTIFICATION);
        assertInstanceOf(S7MessageUserData.class, msg);
        assertEquals(42, msg.getTpduReference());
        S7ParameterUserData param = (S7ParameterUserData) msg.getParameter();
        S7ParameterUserDataItemCPUFunctions item =
            (S7ParameterUserDataItemCPUFunctions) param.getItems().get(0);
        // CPU function group=0x04 (CPU services), type=0x04 (request), subfunction=0x01 (ReadSZL).
        assertEquals(0x04, item.getCpuFunctionGroup());
        assertEquals(0x04, item.getCpuFunctionType());
        assertEquals(0x01, item.getCpuSubfunction());
    }

    @Test
    void parseProbeResponse_decodesS71200Article() {
        byte[] response = mockSzlResponse("6ES7 212-1AE40-0XB0      ");
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(response);
        assertEquals(ControllerType.S7_1200, result.controllerType());
        assertTrue(result.articleNumber().startsWith("6ES7 2"), result.articleNumber());
    }

    @Test
    void parseProbeResponse_decodesS71500Article() {
        assertEquals(ControllerType.S7_1500,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 511-1AK00-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_decodesS7300Article() {
        assertEquals(ControllerType.S7_300,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 314-6CH04-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_decodesS7400Article() {
        assertEquals(ControllerType.S7_400,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 416-2XK02-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_rejectsResponseWithoutSiemensArticleNumber() {
        // Response payload decodes as ASCII but doesn't contain a "6ES7"/"6GK" order number
        // or a "CPU NNN" model name. The response is still structurally valid — that's enough
        // to mark UserData services as supported. ControllerType falls back to ANY.
        byte[] resp = mockSzlResponse("VENDORX 999            ");
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(resp);
        assertEquals(ControllerType.ANY, result.controllerType());
        assertEquals("", result.articleNumber());
    }

    @Test
    void parseProbeResponse_decodesCpuModelNameForS7300() {
        byte[] resp = mockSzlResponse("CPU 315-2 PN/DP                                  ");
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(resp);
        assertEquals(ControllerType.S7_300, result.controllerType());
        assertTrue(result.articleNumber().startsWith("CPU 315"), result.articleNumber());
    }

    @Test
    void parseProbeResponse_decodesCpuModelNameForS71200() {
        byte[] resp = mockSzlResponse("CPU 1212C                                        ");
        assertEquals(ControllerType.S7_1200,
            S7SzlService.parseProbeResponse(resp).controllerType());
    }

    @Test
    void parseProbeResponse_decodesCpuModelNameForS71500() {
        byte[] resp = mockSzlResponse("CPU 1511-1 PN                                    ");
        assertEquals(ControllerType.S7_1500,
            S7SzlService.parseProbeResponse(resp).controllerType());
    }

    @Test
    void parseProbeResponse_unknownModelDigitIsAny() {
        assertEquals(ControllerType.ANY,
            S7SzlService.parseProbeResponse(mockSzlResponse("6ES7 999-0AA00-0AB0")).controllerType());
    }

    @Test
    void parseProbeResponse_acceptsZeroFilledItemAsAny() {
        // Empty/zeroed item is structurally valid — the probe succeeds with an empty
        // article and ANY controllerType. The connection treats that as "UserData works".
        S7SzlService.ProbeResult result = S7SzlService.parseProbeResponse(new byte[16]);
        assertEquals(ControllerType.ANY, result.controllerType());
        assertEquals("", result.articleNumber());
    }

    @Test
    void szlDataOf_rejectsMissingSzlItem() {
        S7Message resp = new S7MessageUserData(7,
            new S7ParameterUserData(Collections.emptyList()),
            new S7PayloadUserData(Collections.emptyList()));
        assertThrows(IllegalStateException.class, () -> S7SzlService.szlDataOf(resp));
    }

    @Test
    void szlDataOf_rejectsNonUserDataMessage() {
        S7Message reqMessage = new S7MessageRequest(1, null, null);
        assertThrows(IllegalArgumentException.class,
            () -> S7SzlService.szlDataOf(reqMessage));
    }

    @Test
    void szlDataOf_returnsTheRawDataSection() {
        byte[] data = {0x00, 0x11, 0x00, 0x00, 0x00, 0x1C, 0x00, 0x01};
        assertArrayEquals(data, S7SzlService.szlDataOf(szlMessage(data, 0, 0)));
    }

    // ------------------------------------------------------------------------
    // Full device identification (SZL 0x0011 / 0x001C / 0x0232)
    // ------------------------------------------------------------------------

    @Test
    void protectionStatusSzlIdEncodesTo0232() {
        SzlId id = S7SzlService.PROTECTION_STATUS;
        assertEquals(SzlModuleTypeClass.CPU, id.getTypeClass());
        assertEquals(0x02, id.getSublistExtract());
        assertEquals(SzlSublist.COMMUNICATION_STATUS_DATA, id.getSublistList());
    }

    @Test
    void parseModuleIdentification_extractsOrderCode() {
        byte[] resp = szlResponse(0x0011, 0x0000, 28,
            moduleRecord(1, "6ES7 511-1AK02-0AB0", 5, 1, 2),
            moduleRecord(6, "6ES7 511-1AK02-0AB0", 1, 0, 0),
            moduleRecord(7, "6ES7 511-1AK02-0AB0", 2, 6, 0));
        assertEquals("6ES7 511-1AK02-0AB0",
            S7SzlService.parseModuleIdentification(resp).orderCode());
    }

    @Test
    void parseModuleIdentification_extractsFirmwareVersionFromRecordSeven() {
        byte[] resp = szlResponse(0x0011, 0x0000, 28,
            moduleRecord(1, "6ES7 511-1AK02-0AB0", 5, 1, 2),
            moduleRecord(7, "6ES7 511-1AK02-0AB0", 2, 6, 0));
        assertEquals("V 2.6.0",
            S7SzlService.parseModuleIdentification(resp).firmwareVersion());
    }

    @Test
    void parseModuleIdentification_extractsHardwareVersionFromRecordSix() {
        byte[] resp = szlResponse(0x0011, 0x0000, 28,
            moduleRecord(6, "6ES7 511-1AK02-0AB0", 1, 0, 0));
        assertEquals("V 1.0.0",
            S7SzlService.parseModuleIdentification(resp).hardwareVersion());
    }

    @Test
    void parseModuleIdentification_leavesFirmwareNullWhenRecordAbsent() {
        byte[] resp = szlResponse(0x0011, 0x0000, 28,
            moduleRecord(1, "6ES7 314-6CH04-0AB0", 3, 3, 0));
        assertNull(S7SzlService.parseModuleIdentification(resp).firmwareVersion());
    }

    @Test
    void parseComponentIdentification_extractsSystemAndModuleName() {
        byte[] resp = szlResponse(0x001C, 0x0000, 34,
            componentRecord(1, "PLANT-LINE-3"),
            componentRecord(2, "PLC_1"));
        S7SzlService.S7DeviceIdentification id = S7SzlService.parseComponentIdentification(resp);
        assertEquals("PLANT-LINE-3", id.systemName());
        assertEquals("PLC_1", id.moduleName());
    }

    @Test
    void parseComponentIdentification_extractsSerialNumberAndPlantDesignation() {
        byte[] resp = szlResponse(0x001C, 0x0000, 34,
            componentRecord(3, "Wall-Mounted Demo Rig"),
            componentRecord(5, "S C-J2UD57132018"));
        S7SzlService.S7DeviceIdentification id = S7SzlService.parseComponentIdentification(resp);
        assertEquals("Wall-Mounted Demo Rig", id.plantDesignation());
        assertEquals("S C-J2UD57132018", id.serialNumber());
    }

    @Test
    void parseComponentIdentification_extractsModuleTypeName() {
        byte[] resp = szlResponse(0x001C, 0x0000, 34,
            componentRecord(7, "CPU 1511-1 PN"));
        assertEquals("CPU 1511-1 PN",
            S7SzlService.parseComponentIdentification(resp).moduleTypeName());
    }

    @Test
    void parseProtectionStatus_reportsConfiguredProtectionLevel() {
        // sch_schal=0 (no selector), sch_par=2 (password parameterised), sch_rel=2 (level in
        // force), bart_sch=2 (RUN-P), anl_sch=0.
        byte[] resp = szlResponse(0x0232, 0x0004, 12, protectionRecord(0, 2, 2, 2, 0));
        S7SzlService.S7DeviceIdentification id = S7SzlService.parseProtectionStatus(resp);
        assertEquals(Integer.valueOf(2), id.protectionLevel());
        assertEquals(Boolean.TRUE, id.passwordProtected());
        assertEquals("RUN-P", id.keySwitchPosition());
    }

    @Test
    void parseProtectionStatus_reportsCpuWithoutPassword() {
        // sch_par=0 means no protection level was parameterised, i.e. no password is set.
        byte[] resp = szlResponse(0x0232, 0x0004, 12, protectionRecord(0, 0, 1, 1, 0));
        S7SzlService.S7DeviceIdentification id = S7SzlService.parseProtectionStatus(resp);
        assertEquals(Integer.valueOf(1), id.protectionLevel());
        assertEquals(Boolean.FALSE, id.passwordProtected());
        assertEquals("RUN", id.keySwitchPosition());
    }

    @Test
    void parseProtectionStatus_reportsUnknownKeySwitchPosition() {
        byte[] resp = szlResponse(0x0232, 0x0004, 12, protectionRecord(0, 0, 3, 0, 0));
        assertEquals("unknown",
            S7SzlService.parseProtectionStatus(resp).keySwitchPosition());
    }

    @Test
    void szlDataOf_rejectsResponseWithNonZeroErrorCode() {
        S7Message resp = szlResponseWithErrorCode(0xD401);
        assertThrows(IllegalStateException.class, () -> S7SzlService.szlDataOf(resp));
    }

    @Test
    void parseIdentification_rejectsTruncatedResponse() {
        assertThrows(IllegalStateException.class,
            () -> S7SzlService.parseModuleIdentification(new byte[4]));
    }

    @Test
    void parseIdentification_ignoresRecordsTruncatedByThePlc() {
        // Declared count of 2 records but only one record's worth of bytes on the wire.
        // The parser must use what arrived rather than reading past the end of the array.
        byte[] header = {0x00, 0x11, 0x00, 0x00, 0x00, 28, 0x00, 0x02};
        byte[] record = moduleRecord(7, "6ES7 511-1AK02-0AB0", 2, 6, 0);
        byte[] data = new byte[header.length + record.length];
        System.arraycopy(header, 0, data, 0, header.length);
        System.arraycopy(record, 0, data, header.length, record.length);
        assertEquals("V 2.6.0", S7SzlService.parseModuleIdentification(data).firmwareVersion());
    }

    @Test
    void mergedWith_fillsNullFieldsFromTheOtherIdentification() {
        S7SzlService.S7DeviceIdentification module = S7SzlService.parseModuleIdentification(
            szlResponse(0x0011, 0x0000, 28, moduleRecord(7, "6ES7 511-1AK02-0AB0", 2, 6, 0)));
        S7SzlService.S7DeviceIdentification component = S7SzlService.parseComponentIdentification(
            szlResponse(0x001C, 0x0000, 34, componentRecord(2, "PLC_1")));

        S7SzlService.S7DeviceIdentification merged = module.mergedWith(component);

        assertEquals("V 2.6.0", merged.firmwareVersion());
        assertEquals("PLC_1", merged.moduleName());
    }

    @Test
    void mergedWith_keepsOwnValueWhenBothAreSet() {
        S7SzlService.S7DeviceIdentification first = S7SzlService.parseComponentIdentification(
            szlResponse(0x001C, 0x0000, 34, componentRecord(2, "PLC_1")));
        S7SzlService.S7DeviceIdentification second = S7SzlService.parseComponentIdentification(
            szlResponse(0x001C, 0x0000, 34, componentRecord(2, "PLC_2")));
        assertEquals("PLC_1", first.mergedWith(second).moduleName());
    }

    @Test
    void mergedWith_toleratesNullOther() {
        S7SzlService.S7DeviceIdentification id = S7SzlService.parseComponentIdentification(
            szlResponse(0x001C, 0x0000, 34, componentRecord(2, "PLC_1")));
        assertEquals("PLC_1", id.mergedWith(null).moduleName());
    }

    // ------------------------------------------------------------------------
    // Why a list was not answered
    // ------------------------------------------------------------------------

    @Test
    void szlDataOf_rejectionCarriesTheErrorCode() {
        S7SzlService.SzlRejectedException rejected = assertThrows(
            S7SzlService.SzlRejectedException.class,
            () -> S7SzlService.szlDataOf(szlResponseWithErrorCode(0xD401)));
        assertEquals(0xD401, rejected.getErrorCode());
    }

    @Test
    void describeFailure_callsAnUnknownSzlIdNotImplemented() {
        // 0xD401 is what a CPU answers for an SZL-ID it doesn't have - the S7-1200/1500
        // reply for the protection-status list. That is not the same as being refused.
        assertEquals("not implemented (0xD401)",
            S7SzlService.describeFailure(new S7SzlService.SzlRejectedException(0xD401)));
    }

    @Test
    void describeFailure_callsAnUnknownIndexNotImplemented() {
        assertEquals("not implemented (0xD402)",
            S7SzlService.describeFailure(new S7SzlService.SzlRejectedException(0xD402)));
    }

    @Test
    void describeFailure_reportsAnyOtherRejectionWithItsCode() {
        assertEquals("rejected (0xD406)",
            S7SzlService.describeFailure(new S7SzlService.SzlRejectedException(0xD406)));
    }

    @Test
    void describeFailure_reportsNoAnswerForATimeout() {
        assertEquals("no answer", S7SzlService.describeFailure(new TimeoutException()));
    }

    @Test
    void describeFailure_unwrapsTheCompletionStageWrapper() {
        // Failures arrive at the driver's handle() wrapped in a CompletionException.
        assertEquals("not implemented (0xD401)", S7SzlService.describeFailure(
            new CompletionException(new S7SzlService.SzlRejectedException(0xD401))));
    }

    // ------------------------------------------------------------------------
    // Chained (multi-part) SZL responses
    // ------------------------------------------------------------------------

    @Test
    void moreDataFollows_isTrueWhenTheResponseSetsLastDataUnit() {
        assertTrue(S7SzlService.moreDataFollows(szlMessage(new byte[8], 1, 1)));
    }

    @Test
    void moreDataFollows_isFalseWhenTheResponseIsComplete() {
        assertFalse(S7SzlService.moreDataFollows(szlMessage(new byte[8], 0, 1)));
    }

    @Test
    void moreDataFollows_isFalseWhenTheParameterOmitsTheFlag() {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x01, (short) 0x00, null, null, null);
        S7Message resp = new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.emptyList()));
        assertFalse(S7SzlService.moreDataFollows(resp));
    }

    @Test
    void sequenceNumberOf_returnsTheSequenceThePlcAssigned() {
        assertEquals(7, S7SzlService.sequenceNumberOf(szlMessage(new byte[8], 1, 7)));
    }

    @Test
    void buildContinuationRequest_echoesTheSequenceNumber() {
        S7Message msg = S7SzlService.buildContinuationRequest(
            42, S7SzlService.COMPONENT_IDENTIFICATION, 0x0000, 7);
        S7ParameterUserDataItemCPUFunctions item = (S7ParameterUserDataItemCPUFunctions)
            ((S7ParameterUserData) msg.getParameter()).getItems().get(0);
        assertEquals(7, item.getSequenceNumber());
        // Still a ReadSZL request — only the sequence number distinguishes it from the first.
        assertEquals(0x04, item.getCpuFunctionGroup());
        assertEquals(0x04, item.getCpuFunctionType());
        assertEquals(0x01, item.getCpuSubfunction());
    }

    /**
     * A chained response splits records at the PDU boundary, so the second chunk carries no
     * header of its own and can start in the middle of a record. Parsing the concatenation has
     * to recover the record that straddles the boundary.
     */
    @Test
    void parseComponentIdentification_readsRecordSplitAcrossChunkBoundary() {
        byte[] whole = szlResponse(0x001C, 0x0000, 34,
            componentRecord(1, "PLANT-LINE-3"),
            componentRecord(5, "S C-J2UD57132018"));
        // Cut mid-way through the second record, as the PLC does at the PDU boundary.
        int cut = 8 + 34 + 20;
        byte[] firstChunk = Arrays.copyOfRange(whole, 0, cut);
        byte[] secondChunk = Arrays.copyOfRange(whole, cut, whole.length);

        byte[] assembled = new byte[firstChunk.length + secondChunk.length];
        System.arraycopy(firstChunk, 0, assembled, 0, firstChunk.length);
        System.arraycopy(secondChunk, 0, assembled, firstChunk.length, secondChunk.length);

        assertEquals("S C-J2UD57132018",
            S7SzlService.parseComponentIdentification(assembled).serialNumber());
    }

    /** A ReadSZL response carrying {@code data}, with the given lastDataUnit and sequence. */
    private static S7Message szlMessage(byte[] data, int lastDataUnit, int sequenceNumber) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x01,
            (short) sequenceNumber, (short) 0x00, (short) lastDataUnit, 0);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.singletonList(
                new S7PayloadUserDataItemCpuFunctionReadSzlResponse(
                    DataTransportErrorCode.OK, DataTransportSize.OCTET_STRING,
                    data.length, data))));
    }

    /** SZL 0x0011 record: index(2) + MLFB(20) + BGTyp(2) + Ausbg(2) + Ausbe(2). */
    private static byte[] moduleRecord(int index, String mlfb, int v1, int v2, int v3) {
        byte[] record = new byte[28];
        record[0] = (byte) (index >> 8);
        record[1] = (byte) index;
        byte[] mlfbBytes = mlfb.getBytes(StandardCharsets.US_ASCII);
        System.arraycopy(mlfbBytes, 0, record, 2, Math.min(mlfbBytes.length, 20));
        for (int i = 2 + Math.min(mlfbBytes.length, 20); i < 22; i++) {
            record[i] = ' ';
        }
        // Snap7 and nmap both read the version as the record's trailing three bytes:
        // Ausbg's low byte followed by Ausbe.
        record[25] = (byte) v1;
        record[26] = (byte) v2;
        record[27] = (byte) v3;
        return record;
    }

    /** SZL 0x001C record: index(2) + 32 ASCII bytes. */
    private static byte[] componentRecord(int index, String value) {
        byte[] record = new byte[34];
        record[0] = (byte) (index >> 8);
        record[1] = (byte) index;
        byte[] valueBytes = value.getBytes(StandardCharsets.US_ASCII);
        System.arraycopy(valueBytes, 0, record, 2, Math.min(valueBytes.length, 32));
        for (int i = 2 + Math.min(valueBytes.length, 32); i < 34; i++) {
            record[i] = ' ';
        }
        return record;
    }

    /** SZL 0x0232 index 4 record: index, sch_schal, sch_par, sch_rel, bart_sch, anl_sch. */
    private static byte[] protectionRecord(int schSchal, int schPar, int schRel,
                                           int bartSch, int anlSch) {
        return words(0x0004, schSchal, schPar, schRel, bartSch, anlSch);
    }

    private static byte[] words(int... values) {
        byte[] out = new byte[values.length * 2];
        for (int i = 0; i < values.length; i++) {
            out[i * 2] = (byte) (values[i] >> 8);
            out[i * 2 + 1] = (byte) values[i];
        }
        return out;
    }

    /** Wrap records in the 8-byte SZL data header the PLC prefixes them with. */
    private static byte[] szlResponse(int szlId, int szlIndex, int recordLength,
                                      byte[]... records) {
        byte[] data = new byte[8 + records.length * recordLength];
        byte[] header = words(szlId, szlIndex, recordLength, records.length);
        System.arraycopy(header, 0, data, 0, 8);
        for (int i = 0; i < records.length; i++) {
            System.arraycopy(records[i], 0, data, 8 + i * recordLength,
                Math.min(records[i].length, recordLength));
        }
        return data;
    }

    private static S7Message szlResponseWithErrorCode(int errorCode) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x01,
            (short) 0x00, null, null, errorCode);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Collections.emptyList()));
    }

    /** Build a 22-byte SZL data section: 2 bytes header + 20-byte article-number window. */
    private static byte[] mockSzlResponse(String article) {
        byte[] articleBytes = article.getBytes(StandardCharsets.US_ASCII);
        byte[] szlItem = new byte[22];
        // Bytes 0..1 are the 2-byte index header — opaque to us. Bytes 2..21 hold the MLFB.
        int copyLen = Math.min(articleBytes.length, 20);
        System.arraycopy(articleBytes, 0, szlItem, 2, copyLen);
        // Pad with spaces so the trim() hits real content.
        for (int i = 2 + copyLen; i < 22; i++) {
            szlItem[i] = ' ';
        }
        return szlItem;
    }

    private static S7Message wrapInUserData(S7PayloadUserDataItem item) {
        S7ParameterUserDataItem param = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x12, (byte) 0x08, (byte) 0x04, (short) 0x01,
            (short) 0x00, null, null, null);
        return new S7MessageUserData(1,
            new S7ParameterUserData(Collections.singletonList(param)),
            new S7PayloadUserData(Arrays.asList(item)));
    }
}
