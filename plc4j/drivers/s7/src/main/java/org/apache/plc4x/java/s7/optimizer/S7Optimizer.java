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
package org.apache.plc4x.java.s7.optimizer;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.s7.readwrite.MemoryArea;
import org.apache.plc4x.java.s7.readwrite.S7AddressAny;
import org.apache.plc4x.java.s7.readwrite.S7MessageRequest;
import org.apache.plc4x.java.s7.readwrite.S7MessageResponseData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7ParameterReadVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarRequest;
import org.apache.plc4x.java.s7.readwrite.S7ParameterWriteVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadReadVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarResponse;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItem;
import org.apache.plc4x.java.s7.readwrite.S7VarRequestParameterItemAddress;
import org.apache.plc4x.java.s7.readwrite.TransportSize;
import org.apache.plc4x.java.s7.context.S7DriverContext;
import org.apache.plc4x.java.s7.tag.S7StringFixedLengthTag;
import org.apache.plc4x.java.s7.tag.S7StringVarLengthTag;
import org.apache.plc4x.java.s7.tag.S7Tag;

import java.util.ArrayList;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Splits {@link PlcReadRequest}/{@link PlcWriteRequest} instances into PDU-sized pieces.
 *
 * <p>The base optimizer treats each user tag as one S7 protocol item. If a single tag's
 * payload exceeds the PDU size (e.g. a large array of bytes), it is split across multiple
 * chunks; the connection re-assembles the per-fragment payloads into the final user value.
 */
public class S7Optimizer {

    public static final int EMPTY_READ_REQUEST_SIZE = new S7MessageRequest(0,
        new S7ParameterReadVarRequest(Collections.emptyList()), null).getLengthInBytes();
    public static final int EMPTY_READ_RESPONSE_SIZE = new S7MessageResponseData(0,
        new S7ParameterReadVarResponse((short) 0),
        new S7PayloadReadVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();
    public static final int EMPTY_WRITE_REQUEST_SIZE = new S7MessageRequest(0,
        new S7ParameterWriteVarRequest(Collections.emptyList()),
        new org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarRequest(Collections.emptyList())).getLengthInBytes();
    public static final int EMPTY_WRITE_RESPONSE_SIZE = new S7MessageResponseData(0,
        new S7ParameterWriteVarResponse((short) 0),
        new S7PayloadWriteVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();
    public static final int S7_ADDRESS_ANY_SIZE = 2 +
        new S7AddressAny(TransportSize.INT, 1, 1, MemoryArea.DATA_BLOCKS, 1, (byte) 0).getLengthInBytes();

    public List<S7ReadChunk> splitReadRequest(PlcReadRequest request, S7DriverContext context) {
        return splitReadFromMap(toLinkedMap(request), context);
    }

    /**
     * Split entry-point that subclasses use after pre-processing the request.
     */
    protected List<S7ReadChunk> splitReadFromMap(LinkedHashMap<String, PlcTag> userTags, S7DriverContext context) {
        List<S7ReadChunk> chunks = new ArrayList<>();
        List<S7ReadChunk.Slot> currentSlots = new ArrayList<>();
        int curRequestSize = EMPTY_READ_REQUEST_SIZE;
        int curResponseSize = EMPTY_READ_RESPONSE_SIZE;

        for (Map.Entry<String, PlcTag> entry : userTags.entrySet()) {
            String tagName = entry.getKey();
            PlcTag plcTag = entry.getValue();
            if (!(plcTag instanceof S7Tag s7Tag)) {
                throw new PlcRuntimeException("Unsupported tag type for tag " + tagName + ": " +
                    (plcTag == null ? "null" : plcTag.getClass().getName()));
            }

            int requestItemSize = S7_ADDRESS_ANY_SIZE;
            int responseItemElementSize = elementSizeBytes(s7Tag);
            int responseItemSize = roundUpEven(4 + s7Tag.getNumberOfElements() * responseItemElementSize);

            // Whole tag fits; either append to current chunk or roll over.
            if (responseItemSize <= context.getPduSize() - EMPTY_READ_RESPONSE_SIZE) {
                if (curRequestSize + requestItemSize > context.getPduSize() ||
                    curResponseSize + responseItemSize > context.getPduSize()) {
                    if (!currentSlots.isEmpty()) {
                        chunks.add(new S7ReadChunk(currentSlots));
                        currentSlots = new ArrayList<>();
                        curRequestSize = EMPTY_READ_REQUEST_SIZE;
                        curResponseSize = EMPTY_READ_RESPONSE_SIZE;
                    }
                }
                S7ReadChunk.Slot slot = new S7ReadChunk.Slot(
                    new S7VarRequestParameterItemAddress(toS7Address(s7Tag)),
                    s7Tag,
                    Collections.singletonList(new S7ReadChunk.Binding(tagName, s7Tag, 0, 0, false)));
                currentSlots.add(slot);
                curRequestSize += requestItemSize;
                curResponseSize += responseItemSize;
            } else {
                // Tag too big for one PDU: byte-split into fragments.
                if (!currentSlots.isEmpty()) {
                    chunks.add(new S7ReadChunk(currentSlots));
                    currentSlots = new ArrayList<>();
                    curRequestSize = EMPTY_READ_REQUEST_SIZE;
                    curResponseSize = EMPTY_READ_RESPONSE_SIZE;
                }
                int maxPayload = context.getPduSize() - (EMPTY_READ_RESPONSE_SIZE + 4);
                int elementSize = s7Tag.getDataType().getSizeInBytes();
                if (elementSize <= 0) {
                    throw new PlcRuntimeException("Cannot split tag '" + tagName + "': element size unknown");
                }
                int elemsPerRequest = Math.max(1, maxPayload / elementSize);
                int totalElems = s7Tag.getNumberOfElements();
                int curElem = 0;
                int curByteOffset = s7Tag.getByteOffset();
                while (curElem < totalElems) {
                    int n = Math.min(elemsPerRequest, totalElems - curElem);
                    S7Tag fragment = new S7Tag(s7Tag.getDataType(), s7Tag.getMemoryArea(),
                        s7Tag.getBlockNumber(), curByteOffset, (byte) 0, n);
                    S7ReadChunk.Slot slot = new S7ReadChunk.Slot(
                        new S7VarRequestParameterItemAddress(toS7Address(fragment)),
                        fragment,
                        Collections.singletonList(
                            new S7ReadChunk.Binding(tagName, s7Tag, 0, curElem, true)));
                    chunks.add(new S7ReadChunk(Collections.singletonList(slot)));
                    curElem += n;
                    curByteOffset += n * elementSize;
                }
            }
        }
        if (!currentSlots.isEmpty()) {
            chunks.add(new S7ReadChunk(currentSlots));
        }
        return chunks;
    }

    public List<S7WriteChunk> splitWriteRequest(PlcWriteRequest request, S7DriverContext context) {
        List<S7WriteChunk> chunks = new ArrayList<>();
        List<String> currentNames = new ArrayList<>();
        int curRequestSize = EMPTY_WRITE_REQUEST_SIZE;
        int curResponseSize = EMPTY_WRITE_RESPONSE_SIZE;

        for (String tagName : request.getTagNames()) {
            PlcTag plcTag = request.getTag(tagName);
            if (!(plcTag instanceof S7Tag s7Tag)) {
                throw new PlcRuntimeException("Unsupported tag type for tag " + tagName);
            }
            int payloadBytes = computeWritePayloadBytes(s7Tag, request);
            int requestItemSize = roundUpEven(S7_ADDRESS_ANY_SIZE + 4 + payloadBytes);
            int responseItemSize = 1;

            if (curRequestSize + requestItemSize > context.getPduSize() ||
                curResponseSize + responseItemSize > context.getPduSize()) {
                if (currentNames.isEmpty()) {
                    throw new PlcRuntimeException(
                        "Tag '" + tagName + "' exceeds PDU size; write splitting not supported");
                }
                chunks.add(new S7WriteChunk(currentNames));
                currentNames = new ArrayList<>();
                curRequestSize = EMPTY_WRITE_REQUEST_SIZE;
                curResponseSize = EMPTY_WRITE_RESPONSE_SIZE;
            }
            currentNames.add(tagName);
            curRequestSize += requestItemSize;
            curResponseSize += responseItemSize;
        }
        if (!currentNames.isEmpty()) {
            chunks.add(new S7WriteChunk(currentNames));
        }
        return chunks;
    }

    // ------------------------------------------------------------------------
    // Helpers
    // ------------------------------------------------------------------------

    protected static int elementSizeBytes(S7Tag s7Tag) {
        if (s7Tag instanceof S7StringVarLengthTag varLen) {
            return varLen.getDataType() == TransportSize.WSTRING ? 4 : 2;
        }
        if (s7Tag instanceof S7StringFixedLengthTag fixed) {
            int chars = fixed.getStringLength() + 2;
            return fixed.getDataType() == TransportSize.WSTRING ? chars * 2 : chars;
        }
        return s7Tag.getDataType().getSizeInBytes();
    }

    private static int computeWritePayloadBytes(S7Tag s7Tag, PlcWriteRequest request) {
        if (s7Tag.getDataType() == TransportSize.BOOL) {
            return (s7Tag.getNumberOfElements() + 7) / 8;
        }
        if (s7Tag instanceof S7StringFixedLengthTag fixed) {
            int bytesPerChar = fixed.getDataType() == TransportSize.WSTRING ? 2 : 1;
            return s7Tag.getNumberOfElements() * (fixed.getStringLength() + 2) * bytesPerChar;
        }
        if (s7Tag instanceof S7StringVarLengthTag varLen) {
            int len = request.getPlcValue(getTagName(request, s7Tag)).getString().length();
            int bytesPerChar = varLen.getDataType() == TransportSize.WSTRING ? 2 : 1;
            return s7Tag.getNumberOfElements() * (len + 2) * bytesPerChar;
        }
        return s7Tag.getNumberOfElements() * s7Tag.getDataType().getSizeInBytes();
    }

    private static String getTagName(PlcWriteRequest req, S7Tag tag) {
        for (String n : req.getTagNames()) {
            if (req.getTag(n) == tag) return n;
        }
        throw new PlcRuntimeException("Tag not found in request");
    }

    protected static int roundUpEven(int n) {
        return (n & 1) == 0 ? n : n + 1;
    }

    protected static org.apache.plc4x.java.s7.readwrite.S7Address toS7Address(S7Tag tag) {
        TransportSize transportSize = tag.getDataType();
        int numElements = tag.getNumberOfElements();
        if (transportSize == TransportSize.STRING) {
            transportSize = TransportSize.CHAR;
            int len = (tag instanceof S7StringFixedLengthTag f) ? f.getStringLength() : 254;
            numElements = numElements * (len + 2);
        } else if (transportSize == TransportSize.WSTRING) {
            transportSize = TransportSize.CHAR;
            int len = (tag instanceof S7StringFixedLengthTag f) ? f.getStringLength() : 254;
            numElements = numElements * (len + 2) * 2;
        } else if (transportSize == TransportSize.BOOL && numElements > 1) {
            numElements = (numElements + 7) / 8;
            transportSize = TransportSize.BYTE;
        }
        if (transportSize.getCode() == 0x00) {
            numElements = numElements * transportSize.getSizeInBytes();
            transportSize = TransportSize.BYTE;
        }
        return new S7AddressAny(transportSize, numElements, tag.getBlockNumber(),
            tag.getMemoryArea(), tag.getByteOffset(), tag.getBitOffset());
    }

    protected static LinkedHashMap<String, PlcTag> toLinkedMap(PlcReadRequest request) {
        LinkedHashMap<String, PlcTag> out = new LinkedHashMap<>();
        for (String n : request.getTagNames()) {
            out.put(n, request.getTag(n));
        }
        return out;
    }
}
