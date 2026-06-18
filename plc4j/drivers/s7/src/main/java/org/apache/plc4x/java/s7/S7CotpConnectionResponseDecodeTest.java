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
package org.apache.plc4x.java.s7;

import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.configuration.S7Configuration;
import org.apache.plc4x.java.s7.optimizer.S7ReadChunk;
import org.apache.plc4x.java.s7.readwrite.DataTransportErrorCode;
import org.apache.plc4x.java.s7.readwrite.DataTransportSize;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.S7MessageResponseData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadReadVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7VarPayloadDataItem;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItemAddress;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcResponseItem;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Method;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.mock;

/**
 * Unit tests for the read-response decoding in {@link S7CotpConnection}.
 *
 * <h2>Issue #2266 — NullPointerException / ArrayIndexOutOfBoundsException when reading
 * an undefined DB variable</h2>
 *
 * <p>The S7 protocol returns a per-item {@code returnCode} for every variable in a
 * multi-item read response.  When a DB variable does not exist, the PLC sets
 * {@code returnCode = NOT_FOUND} (or {@code INVALID_ADDRESS}) and provides either a
 * {@code null} or zero-length {@code data} array.
 *
 * <p><b>Before the fix</b> the code in {@code applyChunkResponse} read
 * {@code byte[] data = item.getData()} without null-checking, then passed it straight
 * into {@code decodeBindingInto}, which called {@code data.length} → NPE, or tried
 * {@code System.arraycopy} with a non-zero offset into a zero-length array →
 * {@code ArrayIndexOutOfBoundsException}.  Both exceptions were swallowed by the
 * surrounding {@code try/catch} and silently reported as {@code INTERNAL_ERROR}.
 *
 * <p><b>After the fix</b>:
 * <ul>
 *   <li>A {@code null} data array is normalised to an empty array before any use.</li>
 *   <li>A non-OK {@code returnCode} always reaches {@code mergeBindingFailure} and is
 *       surfaced as {@code INVALID_ADDRESS} (for NOT_FOUND / INVALID_ADDRESS) without
 *       ever calling {@code decodeBindingInto}.</li>
 *   <li>An OK code with an empty payload is treated as a protocol inconsistency and
 *       reported as {@code INVALID_ADDRESS} with a warning log instead of crashing.</li>
 *   <li>A block-merged binding whose {@code payloadByteOffset} exceeds the actual data
 *       length is reported as {@code INTERNAL_ERROR} instead of throwing AIOOBE.</li>
 * </ul>
 */
class S7CotpConnectionResponseDecodeTest {

    // -----------------------------------------------------------------------
    // Test infrastructure
    // -----------------------------------------------------------------------

    /**
     * Reflective handle for the private {@code applyChunkResponse(S7ReadChunk, S7Message, Map)}
     * method.  We test it directly because it carries all the logic changed by the fix and
     * can be exercised without a real network connection.
     */
    private Method applyChunkResponse;
    private S7CotpConnection connection;

    @BeforeEach
    void setUp() throws Exception {
        S7Configuration cfg = new S7Configuration();
        @SuppressWarnings("unchecked")
        TransportInstance<?> transport = mock(TransportInstance.class);
        AuditLog auditLog = mock(AuditLog.class);

        connection = new S7CotpConnection(cfg, transport, auditLog);

        applyChunkResponse = S7CotpConnection.class.getDeclaredMethod(
            "applyChunkResponse",
            S7ReadChunk.class,
            org.apache.plc4x.java.s7.readwrite.S7Message.class,
            Map.class);
        applyChunkResponse.setAccessible(true);
    }

    // -----------------------------------------------------------------------
    // Helpers
    // -----------------------------------------------------------------------

    /**
     * Build a minimal {@link S7AddressAny} for a single INT at {@code DB<dbNum>.DBW<byteOff>}.
     */
    private static S7AddressAny intDbAddress(int dbNum, int byteOff) {
        return new S7AddressAny(TransportSize.INT, 1, dbNum,
            MemoryArea.DATA_BLOCKS, byteOff, (byte) 0);
    }

    /**
     * Build a one-slot {@link S7ReadChunk} whose binding represents a single
     * non-split, non-block-merged user tag.
     *
     * @param tagAddress S7 address string, e.g. {@code "%DB1.DBW0:INT"}
     */
    private static S7ReadChunk singleSlotChunk(String tagAddress) {
        S7Tag tag = S7Tag.of(tagAddress);
        S7VarRequestParameterItemAddress reqItem =
            new S7VarRequestParameterItemAddress(
                intDbAddress(tag.getBlockNumber(), tag.getByteOffset()));
        S7ReadChunk.Binding binding =
            new S7ReadChunk.Binding(tagAddress, tag, 0, 0, false);
        S7ReadChunk.Slot slot =
            new S7ReadChunk.Slot(reqItem, tag, Collections.singletonList(binding));
        return new S7ReadChunk(Collections.singletonList(slot));
    }

    /**
     * Wrap a variable-length list of payload items in a well-formed
     * {@link S7MessageResponseData} with no PDU-level error.
     */
    private static S7MessageResponseData responseWith(S7VarPayloadDataItem... items) {
        return new S7MessageResponseData(
            0,
            new S7ParameterReadVarResponse((short) items.length),
            new S7PayloadReadVarResponse(List.of(items)),
            (short) 0,   // errorClass = 0 (no PDU-level error)
            (short) 0);  // errorCode  = 0
    }

    /**
     * Invoke {@code applyChunkResponse} and return the result map.
     * Any checked exceptions are re-thrown as unchecked so tests stay readable.
     */
    private Map<String, PlcResponseItem<PlcValue>> invoke(
            S7ReadChunk chunk,
            org.apache.plc4x.java.s7.readwrite.S7Message response) {
        Map<String, PlcResponseItem<PlcValue>> out = new LinkedHashMap<>();
        try {
            applyChunkResponse.invoke(connection, chunk, response, out);
        } catch (java.lang.reflect.InvocationTargetException ite) {
            Throwable cause = ite.getCause();
            throw (cause instanceof RuntimeException re) ? re
                : new RuntimeException("applyChunkResponse threw unexpected checked exception", cause);
        } catch (ReflectiveOperationException e) {
            throw new RuntimeException("Reflection failed", e);
        }
        return out;
    }

    // -----------------------------------------------------------------------
    // Tests — Issue #2266 regression cases
    // -----------------------------------------------------------------------

    /**
     * <b>Primary regression: NOT_FOUND with null data</b>
     *
     * <p>Some PLC firmware variants return {@code returnCode = NOT_FOUND} with a
     * genuinely {@code null} data reference.  Before the fix this caused a
     * {@link NullPointerException} deep in {@code decodeBindingInto}.  After the
     * fix the NPE is impossible and the caller sees {@code INVALID_ADDRESS}.
     */
    @Test
    void notFoundWithNullData_returnsInvalidAddress_noNPE() {
        S7VarPayloadDataItem item = new S7VarPayloadDataItem(
            DataTransportErrorCode.NOT_FOUND,
            DataTransportSize.BYTE_WORD_DWORD,
            null  // <-- the null that triggered the NPE
        );

        S7ReadChunk chunk = singleSlotChunk("%DB1.DBW0:INT");

        // Must not throw — before the fix this threw NullPointerException.
        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, responseWith(item));

        PlcResponseItem<PlcValue> tagResult = result.get("%DB1.DBW0:INT");
        assertNotNull(tagResult, "Result entry must be present for the requested tag");
        assertEquals(
            PlcResponseCode.INVALID_ADDRESS,
            tagResult.getResponseCode(),
            "NOT_FOUND from PLC must surface as INVALID_ADDRESS, not as an NPE or INTERNAL_ERROR");
        assertNull(tagResult.getValue(), "No value expected for a missing DB variable");
    }

    /**
     * <b>Secondary regression: INVALID_ADDRESS with zero-length data</b>
     *
     * <p>When the PLC returns {@code returnCode = INVALID_ADDRESS} and an empty
     * byte array, the old code's {@code data.length} didn't NPE but the block-merged
     * path could hit AIOOBE.  The fix guards both paths.
     */
    @Test
    void invalidAddressWithEmptyData_returnsInvalidAddress_noException() {
        S7VarPayloadDataItem item = new S7VarPayloadDataItem(
            DataTransportErrorCode.INVALID_ADDRESS,
            DataTransportSize.BYTE_WORD_DWORD,
            new byte[0]  // zero-length
        );

        S7ReadChunk chunk = singleSlotChunk("%DB1.DBW4:INT");

        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, responseWith(item));

        PlcResponseItem<PlcValue> tagResult = result.get("%DB1.DBW4:INT");
        assertNotNull(tagResult);
        assertEquals(PlcResponseCode.INVALID_ADDRESS, tagResult.getResponseCode());
        assertNull(tagResult.getValue());
    }

    /**
     * <b>OK return-code but zero-length payload (protocol inconsistency)</b>
     *
     * <p>Certain edge-case PLC firmware bugs produce {@code returnCode = OK} with
     * an empty data array.  Before the fix, attempting {@code data.length} checks
     * in {@code decodeBindingInto} led to AIOOBE; the surrounding catch masked it
     * as {@code INTERNAL_ERROR}.  After the fix the early guard in
     * {@code applyChunkResponse} catches this before entering
     * {@code decodeBindingInto} and maps it to {@code INVALID_ADDRESS}.
     */
    @Test
    void okWithEmptyData_returnsInvalidAddress_noException() {
        S7VarPayloadDataItem item = new S7VarPayloadDataItem(
            DataTransportErrorCode.OK,
            DataTransportSize.BYTE_WORD_DWORD,
            new byte[0]  // protocol inconsistency: OK but no bytes
        );

        S7ReadChunk chunk = singleSlotChunk("%DB2.DBW0:INT");

        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, responseWith(item));

        PlcResponseItem<PlcValue> tagResult = result.get("%DB2.DBW0:INT");
        assertNotNull(tagResult);
        assertEquals(
            PlcResponseCode.INVALID_ADDRESS,
            tagResult.getResponseCode(),
            "OK + empty data must map to INVALID_ADDRESS (protocol inconsistency), not INTERNAL_ERROR");
    }

    /**
     * <b>Mixed response: first tag undefined (NOT_FOUND + null), second tag valid</b>
     *
     * <p>Before the fix, the NPE on the first item also prevented the second
     * (perfectly valid) tag from being decoded.  After the fix each slot is
     * processed independently: the missing tag gets {@code INVALID_ADDRESS} and
     * the valid tag is decoded correctly.
     */
    @Test
    void mixedResponse_undefinedFirstTag_definedSecondTag_processedIndependently() {
        // Slot 1: undefined DB variable → null data with NOT_FOUND
        S7Tag tag1 = S7Tag.of("%DB1.DBW0:INT");
        S7VarRequestParameterItemAddress req1 =
            new S7VarRequestParameterItemAddress(intDbAddress(1, 0));
        S7ReadChunk.Binding b1 = new S7ReadChunk.Binding("%DB1.DBW0:INT", tag1, 0, 0, false);
        S7ReadChunk.Slot slot1 = new S7ReadChunk.Slot(req1, tag1, List.of(b1));

        S7VarPayloadDataItem item1 = new S7VarPayloadDataItem(
            DataTransportErrorCode.NOT_FOUND,
            DataTransportSize.BYTE_WORD_DWORD,
            null);

        // Slot 2: valid DB variable — INT value 42 = 0x002A
        S7Tag tag2 = S7Tag.of("%DB1.DBW2:INT");
        S7VarRequestParameterItemAddress req2 =
            new S7VarRequestParameterItemAddress(intDbAddress(1, 2));
        S7ReadChunk.Binding b2 = new S7ReadChunk.Binding("%DB1.DBW2:INT", tag2, 0, 0, false);
        S7ReadChunk.Slot slot2 = new S7ReadChunk.Slot(req2, tag2, List.of(b2));

        S7VarPayloadDataItem item2 = new S7VarPayloadDataItem(
            DataTransportErrorCode.OK,
            DataTransportSize.BYTE_WORD_DWORD,
            new byte[]{0x00, 0x2A}  // INT 42
        );

        S7ReadChunk chunk = new S7ReadChunk(List.of(slot1, slot2));
        S7MessageResponseData response = responseWith(item1, item2);

        // Must not throw
        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, response);

        // Undefined tag → INVALID_ADDRESS
        PlcResponseItem<PlcValue> res1 = result.get("%DB1.DBW0:INT");
        assertNotNull(res1, "Undefined tag must have a result entry");
        assertEquals(PlcResponseCode.INVALID_ADDRESS, res1.getResponseCode(),
            "Undefined DB variable must be reported as INVALID_ADDRESS");
        assertNull(res1.getValue());

        // Valid tag → OK with decoded value 42
        PlcResponseItem<PlcValue> res2 = result.get("%DB1.DBW2:INT");
        assertNotNull(res2, "Valid tag must have a result entry");
        assertEquals(PlcResponseCode.OK, res2.getResponseCode(),
            "Valid DB variable must decode successfully");
        assertNotNull(res2.getValue(), "Valid tag must have a non-null decoded value");
        assertEquals(42, res2.getValue().getInteger(),
            "INT bytes 0x002A must decode to 42");
    }

    /**
     * <b>Block-merged binding: payload truncated (AIOOBE regression)</b>
     *
     * <p>The {@link org.apache.plc4x.java.s7.optimizer.S7BlockReadOptimizer} can
     * merge adjacent DB variables into a single protocol slot to reduce round trips.
     * In that case a slot may have multiple bindings, each at a different
     * {@code payloadByteOffset}.  If the PLC response is shorter than expected (e.g.
     * because one of the merged variables is undefined), the old code threw
     * {@code ArrayIndexOutOfBoundsException} in
     * {@code System.arraycopy(data, b.payloadByteOffset(), ...)}.  After the fix
     * the short-data guard in {@code decodeBindingInto} catches this and reports
     * {@code INTERNAL_ERROR} instead of crashing.
     */
    @Test
    void blockMergedSlot_payloadTooShortForSecondBinding_returnsInternalError_noException() {
        // Simulate a block-read slot that covers two INTs at byte-offsets 0 and 2
        // (4 bytes total), but the PLC only returns 2 bytes.

        // The "block" tag covers bytes 0-3 as a BYTE array of length 4.
        S7Tag blockTag = new S7Tag(
            TransportSize.BYTE,
            MemoryArea.DATA_BLOCKS,
            1, 0, (byte) 0, 4  // DB1.DBB0, 4 bytes
        );
        S7VarRequestParameterItemAddress reqItem =
            new S7VarRequestParameterItemAddress(
                new S7AddressAny(TransportSize.BYTE, 4, 1, MemoryArea.DATA_BLOCKS, 0, (byte) 0));

        // Two user tags merged into this slot, at offsets 0 and 2 within the block.
        S7Tag userTag1 = S7Tag.of("%DB1.DBW0:INT");
        S7Tag userTag2 = S7Tag.of("%DB1.DBW2:INT");

        S7ReadChunk.Binding bind1 = new S7ReadChunk.Binding("first",  userTag1, 0, 0, false);
        S7ReadChunk.Binding bind2 = new S7ReadChunk.Binding("second", userTag2, 2, 0, false);

        S7ReadChunk.Slot slot = new S7ReadChunk.Slot(reqItem, blockTag, List.of(bind1, bind2));
        S7ReadChunk chunk = new S7ReadChunk(List.of(slot));

        // PLC only returns 2 bytes — offset 2 for "second" is out of range.
        S7VarPayloadDataItem item = new S7VarPayloadDataItem(
            DataTransportErrorCode.OK,
            DataTransportSize.BYTE_WORD_DWORD,
            new byte[]{0x01, 0x00}  // only 2 bytes, not 4
        );

        S7MessageResponseData response = responseWith(item);

        // Must not throw — before the fix this threw ArrayIndexOutOfBoundsException
        // which was swallowed and converted to INTERNAL_ERROR (still a bad user experience).
        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, response);

        // "second" binding at offset 2 is beyond the 2-byte payload.
        PlcResponseItem<PlcValue> res2 = result.get("second");
        assertNotNull(res2, "'second' must have a result entry even if the payload is too short");
        assertEquals(
            PlcResponseCode.INTERNAL_ERROR,
            res2.getResponseCode(),
            "Out-of-range block binding must be INTERNAL_ERROR, not a thrown exception");
    }

    /**
     * <b>Happy-path smoke test</b>
     *
     * <p>Confirms that the fix did not regress the normal success path: a single
     * INT read returning two valid bytes decodes to the correct integer value.
     */
    @Test
    void happyPath_validIntResponse_decodesCorrectly() {
        // INT value 1234 = 0x04D2
        S7VarPayloadDataItem item = new S7VarPayloadDataItem(
            DataTransportErrorCode.OK,
            DataTransportSize.BYTE_WORD_DWORD,
            new byte[]{0x04, (byte) 0xD2}
        );

        S7ReadChunk chunk = singleSlotChunk("%DB1.DBW0:INT");

        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, responseWith(item));

        PlcResponseItem<PlcValue> tagResult = result.get("%DB1.DBW0:INT");
        assertNotNull(tagResult, "Result must be present for the requested tag");
        assertEquals(PlcResponseCode.OK, tagResult.getResponseCode());
        assertNotNull(tagResult.getValue(), "A valid INT response must carry a decoded value");
        assertEquals(1234, tagResult.getValue().getInteger(),
            "Bytes 0x04D2 must decode to integer 1234");
    }

    /**
     * <b>NOT_FOUND with empty (not null) data array</b>
     *
     * <p>A second firmware variant that produces zero-length data alongside a
     * non-OK return code — distinct from the null-data case.
     */
    @Test
    void notFoundWithEmptyData_returnsInvalidAddress_noException() {
        S7VarPayloadDataItem item = new S7VarPayloadDataItem(
            DataTransportErrorCode.NOT_FOUND,
            DataTransportSize.BYTE_WORD_DWORD,
            new byte[0]
        );

        S7ReadChunk chunk = singleSlotChunk("%DB3.DBW0:INT");

        Map<String, PlcResponseItem<PlcValue>> result = invoke(chunk, responseWith(item));

        PlcResponseItem<PlcValue> tagResult = result.get("%DB3.DBW0:INT");
        assertNotNull(tagResult);
        assertEquals(PlcResponseCode.INVALID_ADDRESS, tagResult.getResponseCode());
        assertNull(tagResult.getValue());
    }
}
