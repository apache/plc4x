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
import org.apache.plc4x.java.api.messages.PlcTagRequest;
import org.apache.plc4x.java.api.types.PlcResponseCode;
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
            int responseItemSize = roundUpEven(4 + payloadBytes(s7Tag));

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
                // Tag too big for one PDU: fetch the bytes it covers in fragments.
                if (!currentSlots.isEmpty()) {
                    chunks.add(new S7ReadChunk(currentSlots));
                    currentSlots = new ArrayList<>();
                    curRequestSize = EMPTY_READ_REQUEST_SIZE;
                    curResponseSize = EMPTY_READ_RESPONSE_SIZE;
                }
                splitIntoByteFragments(tagName, s7Tag, context, chunks);
            }
        }
        if (!currentSlots.isEmpty()) {
            chunks.add(new S7ReadChunk(currentSlots));
        }
        return chunks;
    }

    /**
     * A tag too big for one PDU is fetched as plain byte fragments of the payload a read of it
     * would return, one fragment per chunk. Each binding carries where its fragment goes in that
     * payload; the connection decodes the tag once every fragment is in. Splitting the payload
     * rather than the elements works for every type, including BOOL arrays, whose elements are
     * bits, and strings, whose elements are bigger than their transport size says.
     */
    private static void splitIntoByteFragments(String tagName, S7Tag s7Tag, S7DriverContext context,
                                               List<S7ReadChunk> chunks) {
        int totalBytes = payloadBytes(s7Tag);
        // Four bytes of item header, and the payload is padded to an even length.
        int bytesPerRequest = Math.max(1, context.getPduSize() - EMPTY_READ_RESPONSE_SIZE - 5);
        for (int curByte = 0; curByte < totalBytes; curByte += bytesPerRequest) {
            int n = Math.min(bytesPerRequest, totalBytes - curByte);
            S7Tag fragment = new S7Tag(TransportSize.BYTE, s7Tag.getMemoryArea(),
                s7Tag.getBlockNumber(), s7Tag.getByteOffset() + curByte, (byte) 0, n);
            S7ReadChunk.Slot slot = new S7ReadChunk.Slot(
                new S7VarRequestParameterItemAddress(toS7Address(fragment)),
                fragment,
                Collections.singletonList(new S7ReadChunk.Binding(tagName, s7Tag, 0, curByte, true)));
            chunks.add(new S7ReadChunk(Collections.singletonList(slot)));
        }
    }

    public List<S7WriteChunk> splitWriteRequest(PlcWriteRequest request, S7DriverContext context) {
        List<S7WriteChunk> chunks = new ArrayList<>();
        List<S7WriteChunk.Item> currentItems = new ArrayList<>();
        int curRequestSize = EMPTY_WRITE_REQUEST_SIZE;
        int curResponseSize = EMPTY_WRITE_RESPONSE_SIZE;

        for (String tagName : request.getTagNames()) {
            if (isRejected(request, tagName)) {
                continue;
            }
            PlcTag plcTag = request.getTag(tagName);
            if (!(plcTag instanceof S7Tag s7Tag)) {
                throw new PlcRuntimeException("Unsupported tag type for tag " + tagName);
            }
            for (S7WriteChunk.Item item : writeItemsFor(tagName, s7Tag)) {
                int payloadBytes = computeWritePayloadBytes(item.tag(), request, tagName);
                int requestItemSize = roundUpEven(S7_ADDRESS_ANY_SIZE + 4 + payloadBytes);
                int responseItemSize = 1;

                if (curRequestSize + requestItemSize > context.getPduSize() ||
                    curResponseSize + responseItemSize > context.getPduSize()) {
                    if (currentItems.isEmpty()) {
                        throw new PlcRuntimeException(
                            "Tag '" + tagName + "' exceeds PDU size; write splitting not supported");
                    }
                    chunks.add(new S7WriteChunk(currentItems));
                    currentItems = new ArrayList<>();
                    curRequestSize = EMPTY_WRITE_REQUEST_SIZE;
                    curResponseSize = EMPTY_WRITE_RESPONSE_SIZE;
                }
                currentItems.add(item);
                curRequestSize += requestItemSize;
                curResponseSize += responseItemSize;
            }
        }
        if (!currentItems.isEmpty()) {
            chunks.add(new S7WriteChunk(currentItems));
        }
        return chunks;
    }

    /**
     * The protocol items a tag is written with. A write of whole bytes overwrites every bit in
     * them, so a BOOL array is only written packed for the bytes it covers completely; each bit it
     * shares a byte with something else goes out as an item of its own, leaving the rest of that
     * byte alone. Any other tag is a single item.
     */
    public static List<S7WriteChunk.Item> writeItemsFor(String tagName, S7Tag tag) {
        if (!isPackedBoolArray(tag)) {
            return List.of(new S7WriteChunk.Item(tagName, tag, 0));
        }
        List<S7WriteChunk.Item> items = new ArrayList<>();
        int numElements = tag.getNumberOfElements();
        int firstBit = tag.getByteOffset() * 8 + tag.getBitOffset();
        int element = 0;
        while (element < numElements) {
            int bit = firstBit + element;
            int wholeBytes = (bit % 8 == 0) ? (numElements - element) / 8 : 0;
            int count = (wholeBytes > 0) ? wholeBytes * 8 : 1;
            S7Tag piece = new S7Tag(TransportSize.BOOL, tag.getMemoryArea(), tag.getBlockNumber(),
                bit / 8, (byte) (bit % 8), count);
            items.add(new S7WriteChunk.Item(tagName, piece, element));
            element += count;
        }
        return items;
    }

    // ------------------------------------------------------------------------
    // Helpers
    // ------------------------------------------------------------------------

    /**
     * How many bytes of data a read of the tag returns: what {@link #toS7Address} asks for. A BOOL
     * array returns the bytes its bits are packed into, a string its length bytes and characters,
     * a type without an own transport size its bytes.
     */
    public static int payloadBytes(S7Tag tag) {
        S7AddressAny address = (S7AddressAny) toS7Address(tag);
        return address.getNumberOfElements() * address.getTransportSize().getSizeInBytes();
    }

    /**
     * Whether the tag is a BOOL array, which S7 transfers as the bytes its bits are packed into
     * rather than element by element.
     */
    public static boolean isPackedBoolArray(S7Tag s7Tag) {
        return s7Tag.getDataType() == TransportSize.BOOL && s7Tag.getNumberOfElements() > 1;
    }

    /**
     * How many bytes the bits of a BOOL tag span, counting from bit 0 of its start byte: a tag
     * that starts at a bit offset other than 0 can spill into one more byte than its element
     * count alone needs.
     */
    public static int packedBoolBytes(S7Tag s7Tag) {
        return (s7Tag.getBitOffset() + s7Tag.getNumberOfElements() + 7) / 8;
    }

    private static int computeWritePayloadBytes(S7Tag s7Tag, PlcWriteRequest request, String tagName) {
        if (s7Tag.getDataType() == TransportSize.BOOL) {
            return (s7Tag.getNumberOfElements() + 7) / 8;
        }
        if (s7Tag instanceof S7StringFixedLengthTag fixed) {
            int bytesPerChar = fixed.getDataType() == TransportSize.WSTRING ? 2 : 1;
            return s7Tag.getNumberOfElements() * (fixed.getStringLength() + 2) * bytesPerChar;
        }
        if (s7Tag instanceof S7StringVarLengthTag varLen) {
            int len = request.getPlcValue(tagName).getString().length();
            int bytesPerChar = varLen.getDataType() == TransportSize.WSTRING ? 2 : 1;
            return s7Tag.getNumberOfElements() * (len + 2) * bytesPerChar;
        }
        return s7Tag.getNumberOfElements() * s7Tag.getDataType().getSizeInBytes();
    }

    protected static int roundUpEven(int n) {
        return (n & 1) == 0 ? n : n + 1;
    }

    /**
     * The address a read of the tag asks for. Strings, BOOL arrays and types without an own
     * transport size are requested as the characters or bytes they occupy; cyclic subscriptions
     * address their tags the same way.
     */
    public static org.apache.plc4x.java.s7.readwrite.S7Address toS7Address(S7Tag tag) {
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
        } else if (isPackedBoolArray(tag)) {
            // The bytes the bits are packed into, from bit 0 of the start byte: a byte read has no
            // bit address, so the connection skips the bits before the tag's own bit offset.
            return new S7AddressAny(TransportSize.BYTE, packedBoolBytes(tag), tag.getBlockNumber(),
                tag.getMemoryArea(), tag.getByteOffset(), (byte) 0);
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
            if (isRejected(request, n)) {
                // Address the builder couldn't parse: it has no tag to chunk, and the
                // connection reports its code directly.
                continue;
            }
            out.put(n, request.getTag(n));
        }
        return out;
    }

    /**
     * Whether the request builder already rejected this tag (unparseable address, bad value),
     * in which case it carries an error code and a {@code null} tag.
     */
    protected static boolean isRejected(PlcTagRequest request, String tagName) {
        return request.getTagResponseCode(tagName) != PlcResponseCode.OK;
    }
}
