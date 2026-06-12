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

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
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
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.s7.tag.S7Tag;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcBrowseItem;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

/**
 * S7Comm UserData "Browse" service. Exposes:
 * <ul>
 *   <li>Static memory areas (M / I / Q) — always available, even on devices without
 *       UserData services.</li>
 *   <li>Data block enumeration via the Block Functions UserData service
 *       ({@code cpuFunctionGroup=0x05}, subfunction {@code 0x02} "list of blocks of
 *       type"). Field-level browse is out of scope this round.</li>
 * </ul>
 */
public final class S7BrowseService {

    /**
     * Two-ASCII-char block-type code packed into a uint16 — what the PLC's "list of
     * blocks of type" request expects in its data section.
     */
    public enum BlockType {
        OB (0x3038, "OB"),
        DB (0x3041, "DB"),
        SDB(0x3042, "SDB"),
        FC (0x3043, "FC"),
        SFC(0x3044, "SFC"),
        FB (0x3045, "FB"),
        SFB(0x3046, "SFB");

        public final int wireValue;
        public final String shortName;

        BlockType(int wireValue, String shortName) {
            this.wireValue = wireValue;
            this.shortName = shortName;
        }
    }

    private S7BrowseService() {
    }

    /**
     * Diagnostic helper: dumps the response's CPU-function header + raw block-list bytes.
     * Drop once browse is reliable on the controller families we care about.
     */
    public static void dumpListBlocksResponse(S7Message response, org.slf4j.Logger logger) {
        if (!(response instanceof S7MessageUserData)) {
            logger.info("Block-list response not a UserData message: {}",
                response == null ? "null" : response.getClass().getSimpleName());
            return;
        }
        S7ParameterUserData param = (S7ParameterUserData) response.getParameter();
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu) {
                logger.info("Block-list CPU params: method=0x{}, type=0x{}, group=0x{}, subfn=0x{}, errCode={}",
                    Integer.toHexString(cpu.getMethod() & 0xFF),
                    Integer.toHexString(cpu.getCpuFunctionType() & 0xFF),
                    Integer.toHexString(cpu.getCpuFunctionGroup() & 0xFF),
                    Integer.toHexString(cpu.getCpuSubfunction() & 0xFFFF),
                    cpu.getErrorCode());
            }
        }
        S7PayloadUserData payload = (S7PayloadUserData) response.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            byte[] bytes = item instanceof S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse resp
                ? resp.getItems()
                : null;
            logger.info("Block-list payload item: type={}, bytes={}",
                item.getClass().getSimpleName(),
                bytes == null ? "<not a ListBlocksOfType response>" : (bytes.length + " bytes"));
            if (bytes != null) {
                StringBuilder sb = new StringBuilder();
                for (byte b : bytes) sb.append(String.format("%02X ", b & 0xFF));
                logger.info("Block-list raw items: {}", sb);
            }
        }
    }

    /**
     * Build the static list of top-level memory areas that exist on every S7 family. These
     * give callers somewhere to start browsing without depending on UserData services —
     * useful for LOGO and similar low-end devices where the SZL probe failed.
     */
    public static List<PlcBrowseItem> staticAreas(boolean subscribable) {
        List<PlcBrowseItem> items = new ArrayList<>(3);
        items.add(area("M", MemoryArea.FLAGS_MARKERS, subscribable));
        items.add(area("I", MemoryArea.INPUTS,        subscribable));
        items.add(area("Q", MemoryArea.OUTPUTS,       subscribable));
        return items;
    }

    private static PlcBrowseItem area(String shortName, MemoryArea area, boolean subscribable) {
        // Placeholder tag for the area root. Caller composes the real offset / data type.
        S7Tag tag = new S7Tag(TransportSize.BYTE, area, 0, 0, (byte) 0, 1);
        Set<PlcSubscriptionType> subscriptionTypes =
            subscribable ? Set.of(PlcSubscriptionType.CYCLIC) : Collections.emptySet();
        return new DefaultPlcBrowseItem(tag, "%" + shortName, true, true, subscriptionTypes, false,
            Collections.emptyList(), Collections.emptyMap(), Collections.emptyMap());
    }

    /**
     * Construct the {@code S7MessageUserData} for a "list of blocks of type" request.
     * Caller assigns the tpduReference and dispatches it through the connection.
     */
    public static S7Message buildListBlocksOfTypeRequest(int tpduRef, BlockType type) {
        S7ParameterUserDataItem parameter = new S7ParameterUserDataItemCPUFunctions(
            (short) 0x11,    // method: request
            (byte)  0x04,    // cpuFunctionType: request
            (byte)  0x03,    // cpuFunctionGroup: Block functions
            (short) 0x02,    // cpuSubfunction: list of blocks of type
            (short) 0x00,    // sequenceNumber
            null, null, null);
        S7PayloadUserDataItem payload = new S7PayloadUserDataItemCpuFunctionListBlocksOfTypeRequest(
            DataTransportErrorCode.OK,
            DataTransportSize.OCTET_STRING,
            2,                // dataLength: 2 bytes for the packed block-type code
            type.wireValue);
        return new S7MessageUserData(tpduRef,
            new S7ParameterUserData(Collections.singletonList(parameter)),
            new S7PayloadUserData(Collections.singletonList(payload)));
    }

    /**
     * Parse the response of a "list of blocks of type" request into a list of block numbers.
     * Each entry on the wire is 4 bytes: {@code blockNumber:uint16 BE} + {@code flags:uint8}
     * + {@code language:uint8}. We only return the block numbers — flags/language can be
     * fetched separately via the GetBlockInfo subfunction if needed later.
     *
     * <p>Returns an empty list when the PLC returned a non-zero error code (some devices
     * legitimately have no blocks of a given type, and we treat that as "empty list" not
     * "browse failed").
     */
    public static List<Integer> parseListBlocksOfTypeResponse(S7Message response) {
        if (!(response instanceof S7MessageUserData)) {
            return Collections.emptyList();
        }
        if (hasParameterErrorCode(response)) {
            return Collections.emptyList();
        }
        S7PayloadUserData payload = (S7PayloadUserData) response.getPayload();
        for (S7PayloadUserDataItem item : payload.getItems()) {
            if (!(item instanceof S7PayloadUserDataItemCpuFunctionListBlocksOfTypeResponse resp)) {
                continue;
            }
            byte[] bytes = resp.getItems();
            if (bytes == null || bytes.length < 4) {
                return Collections.emptyList();
            }
            int count = bytes.length / 4;
            List<Integer> blockNumbers = new ArrayList<>(count);
            for (int i = 0; i < count; i++) {
                int off = i * 4;
                int blockNumber = ((bytes[off] & 0xFF) << 8) | (bytes[off + 1] & 0xFF);
                blockNumbers.add(blockNumber);
            }
            return blockNumbers;
        }
        return Collections.emptyList();
    }

    /**
     * Convert a list of DB block numbers into browse items addressed as {@code DB<n>}. The
     * tag itself uses {@code TransportSize.BYTE} as a placeholder — callers compose the
     * real offset / data type when they request a tag inside the DB.
     */
    public static List<PlcBrowseItem> dataBlocksFromNumbers(List<Integer> blockNumbers, boolean subscribable) {
        List<PlcBrowseItem> items = new ArrayList<>(blockNumbers.size());
        for (int blockNumber : blockNumbers) {
            items.add(dataBlockItem(blockNumber, subscribable));
        }
        return items;
    }

    private static PlcBrowseItem dataBlockItem(int blockNumber, boolean subscribable) {
        S7Tag tag = new S7Tag(TransportSize.BYTE, MemoryArea.DATA_BLOCKS, blockNumber, 0, (byte) 0, 1);
        Map<String, PlcBrowseItem> noChildren = Collections.emptyMap();
        Set<PlcSubscriptionType> subscriptionTypes =
            subscribable ? Set.of(PlcSubscriptionType.CYCLIC) : Collections.emptySet();
        return new DefaultPlcBrowseItem(tag, "DB" + blockNumber, true, true, subscriptionTypes, false,
            Collections.emptyList(), noChildren, new HashMap<>());
    }

    private static boolean hasParameterErrorCode(S7Message response) {
        S7ParameterUserData param = (S7ParameterUserData) response.getParameter();
        for (S7ParameterUserDataItem pItem : param.getItems()) {
            if (pItem instanceof S7ParameterUserDataItemCPUFunctions cpu
                && cpu.getErrorCode() != null && cpu.getErrorCode() != 0) {
                return true;
            }
        }
        return false;
    }
}
