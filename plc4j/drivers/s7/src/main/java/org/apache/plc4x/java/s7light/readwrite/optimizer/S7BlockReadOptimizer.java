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
package org.apache.plc4x.java.s7light.readwrite.optimizer;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.tag.*;
import org.apache.plc4x.java.s7light.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcRawByteArray;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.stream.IntStream;

/**
 * In contrast to the S7Optimizer, the S7BlockReadOptimizer only looks at the bytes that require being read
 * and focuses on reading raw blocks of bytes instead of doing individual reads for each item.
 */
public class S7BlockReadOptimizer extends S7Optimizer {

    private static final Logger logger = LoggerFactory.getLogger(S7BlockReadOptimizer.class);

    public static final int EMPTY_READ_REQUEST_SIZE = new S7MessageRequest(0, new S7ParameterReadVarRequest(
        Collections.emptyList()), null).getLengthInBytes();
    public static final int EMPTY_READ_RESPONSE_SIZE = new S7MessageResponseData(0, new S7ParameterReadVarResponse(
        (short) 0), new S7PayloadReadVarResponse(Collections.emptyList()), (short) 0, (short) 0).getLengthInBytes();
    public static final int S7_ADDRESS_ANY_SIZE = 2 +
        new S7AddressAny(TransportSize.INT, 1, 1, MemoryArea.DATA_BLOCKS, 1, (byte) 0).getLengthInBytes();

    @Override
    protected List<PlcReadRequest> processReadRequest(PlcReadRequest readRequest, DriverContext driverContext) {
        S7DriverContext s7DriverContext = (S7DriverContext) driverContext;
        if(!s7DriverContext.isEnableBlockReadOptimizer()) {
            return super.processReadRequest(readRequest, s7DriverContext);
        }

        // Sort the tags by area
        // (We can only read multiple tags in one byte array, if they are located in the same area)
        Map<String, Map<PlcTag, String>> sortedTagsPerArea = new HashMap<>();
        for (String tagName : readRequest.getTagNames()) {
            PlcTag tag = readRequest.getTag(tagName);
            if(tag instanceof S7SzlTag) {
                if(!sortedTagsPerArea.containsKey("SZL")) {
                    sortedTagsPerArea.put("SZL", createSortedTagMap());
                }
                sortedTagsPerArea.get("SZL").put(tag, tagName);
            } else if(tag instanceof S7ClkTag) {
                if(!sortedTagsPerArea.containsKey("CLK")) {
                    sortedTagsPerArea.put("CLK", createSortedTagMap());
                }
                sortedTagsPerArea.get("CLK").put(tag, tagName);
            } else if(tag instanceof S7Tag) {
                S7Tag s7Tag = (S7Tag) tag;
                MemoryArea memoryArea = s7Tag.getMemoryArea();
                // When reading DATA_BLOCKS we need to also use the block number.
                String areaName = memoryArea.getShortName();
                if(memoryArea ==  MemoryArea.DATA_BLOCKS) {
                    areaName += s7Tag.getBlockNumber();
                } else if(memoryArea == MemoryArea.INSTANCE_DATA_BLOCKS) {
                    areaName += s7Tag.getBlockNumber();
                }
                if(!sortedTagsPerArea.containsKey(areaName)) {
                    sortedTagsPerArea.put(areaName, createSortedTagMap());
                }
                sortedTagsPerArea.get(areaName).put(tag, tagName);
            } else {
                System.out.println("Ignored");
            }
        }

        // Build groups of byte-arrays to fetch
        // Strategies:
        // - If two consecutive tags are more than S7_ADDRESS_ANY_SIZE bytes apart,
        //   reading them in one block is less efficient than reading them separately.
        //   Other than that, extend the byte array to include the next tag.
        Map<TagNameSize, PlcTagItem<PlcTag>> optimizedTagMap = new TreeMap<>();
        for (Map<PlcTag, String> tagList : sortedTagsPerArea.values()) {
            MemoryArea currentMemoryArea = null;
            int currentDataBlockNumber = -1;
            int currentChunkStartByteOffset = -1;
            int currentChunkEndByteOffset = -1;
            Map<PlcTag, String> currentChunkTags = createSortedTagMap();
            for (PlcTag plcTag : tagList.keySet()) {
                // We don't do anything for Szl and Clk tags.
                if(plcTag instanceof S7SzlTag) {
                    // TODO: Implement the size
                    optimizedTagMap.put(new TagNameSize(tagList.get(plcTag), 0), new DefaultPlcTagItem<>(plcTag));
                } else if (plcTag instanceof S7ClkTag) {
                    // TODO: Implement the size
                    optimizedTagMap.put(new TagNameSize(tagList.get(plcTag), 0), new DefaultPlcTagItem<>(plcTag));
                }
                // Var-length strings, are a performance nightmare. Trying to optimize reading them is probably not
                // worth the effort. For now, we simply handle them as un-chunked tags.
                else if(plcTag instanceof S7StringVarLengthTag) {
                    // A var-length string tag simply reads 2 or 4 bytes.
                    optimizedTagMap.put(new TagNameSize(tagList.get(plcTag),
                        ((S7StringVarLengthTag) plcTag).getDataType() == TransportSize.STRING ? 2 : 4),
                        new DefaultPlcTagItem<>(plcTag));
                }

                // Only regular tags are optimized.
                else if (plcTag instanceof S7Tag) {
                    S7Tag s7Tag = (S7Tag) plcTag;

                    int curTagSize = s7Tag.getDataType().getSizeInBytes() * s7Tag.getNumberOfElements();
                    // In case of fixed length strings, a string starts with two bytes: max length,
                    // actual length and then the string bytes after that.
                    if(s7Tag instanceof S7StringFixedLengthTag) {
                        S7StringFixedLengthTag stringFixedLengthTag = (S7StringFixedLengthTag) s7Tag;
                        int bytesPerChar = (stringFixedLengthTag.getDataType() == TransportSize.WSTRING) ? 2 : 1;
                        curTagSize = ((2 * bytesPerChar) + (stringFixedLengthTag.getStringLength() * bytesPerChar)) * s7Tag.getNumberOfElements();
                    }

                    // If this is the first tag, use that as starting point.
                    if(currentMemoryArea == null) {
                        currentMemoryArea = s7Tag.getMemoryArea();
                        currentDataBlockNumber = s7Tag.getBlockNumber();
                        currentChunkStartByteOffset = s7Tag.getByteOffset();
                        currentChunkEndByteOffset = s7Tag.getByteOffset() + curTagSize;
                    }
                    // If the next tag would be more bytes away than a s7 address item requires, it's cheaper to
                    // split up into multiple items.
                    else if(currentChunkEndByteOffset + S7_ADDRESS_ANY_SIZE < s7Tag.getByteOffset()) {
                        // Save the current chunk.
                        optimizedTagMap.put(new TagNameSize("__chunk__" + optimizedTagMap.size(), currentChunkEndByteOffset - currentChunkStartByteOffset),
                            new DefaultPlcTagItem<>(
                                new S7TagChunk(TransportSize.BYTE, currentMemoryArea, currentDataBlockNumber,
                                    currentChunkStartByteOffset, (byte) 0,
                                    currentChunkEndByteOffset - currentChunkStartByteOffset,
                                    currentChunkTags, 0, 1, currentChunkEndByteOffset - currentChunkStartByteOffset)));

                        // Start a new one.
                        currentChunkStartByteOffset = s7Tag.getByteOffset();
                        currentChunkEndByteOffset = s7Tag.getByteOffset() + curTagSize;
                        currentChunkTags = createSortedTagMap();
                    }
                    // Otherwise extend the array size to include this tag.
                    else {
                        currentChunkEndByteOffset = s7Tag.getByteOffset() + curTagSize;
                    }

                    // Add the tag to the list of tags for the current chunk.
                    currentChunkTags.put(s7Tag, tagList.get(s7Tag));
                }
            }

            // Finish the last chunk.
            optimizedTagMap.put(new TagNameSize("__chunk__" + optimizedTagMap.size(), currentChunkEndByteOffset - currentChunkStartByteOffset),
                new DefaultPlcTagItem<>(
                    new S7TagChunk(TransportSize.BYTE, currentMemoryArea, currentDataBlockNumber,
                        currentChunkStartByteOffset, (byte) 0,
                        currentChunkEndByteOffset - currentChunkStartByteOffset,
                        currentChunkTags, 0, 1, currentChunkEndByteOffset - currentChunkStartByteOffset)));
        }

        // Go through all chunks. If there are ones larger than the max PDU size, split them up into
        // multiple tags, that utilize the packets to the maximum.
        final int maxRequestSize = ((S7DriverContext) driverContext).getPduSize() - (EMPTY_READ_RESPONSE_SIZE + 4);
        Map<TagNameSize, PlcTagItem<PlcTag>> optimizedTagMap2 = new TreeMap<>();
        for (TagNameSize tagNameSize : optimizedTagMap.keySet()) {
            PlcTagItem<PlcTag> curTagItem = optimizedTagMap.get(tagNameSize);
            S7Tag curTag = (S7Tag) curTagItem.getTag();
            if(tagNameSize.getTagSize() > maxRequestSize) {
                String curTagNameBase = tagNameSize.getTagName();
                int curTagFragmentNumber = 0;
                int curTagSize = tagNameSize.getTagSize();
                int curTagOffset = curTag.getByteOffset();
                int curPartIndex = 0;
                int totalPartCount = (curTagSize / maxRequestSize) + 1;
                while(curTagSize > maxRequestSize) {
                    optimizedTagMap2.put(new TagNameSize(curTagNameBase + "." + curTagFragmentNumber, maxRequestSize),
                        new DefaultPlcTagItem<>(
                            new S7TagChunk(curTag.getDataType(), curTag.getMemoryArea(), curTag.getBlockNumber(), curTagOffset, (byte) 0, maxRequestSize,
                                (curTag instanceof S7TagChunk) ? ((S7TagChunk) curTag).getChunkTags() : Collections.singletonMap(curTag, tagNameSize.getTagName()),
                                curPartIndex, totalPartCount, curTagSize)));

                    curTagOffset += maxRequestSize;
                    curTagSize -= maxRequestSize;
                    curTagFragmentNumber++;
                    curPartIndex++;
                }
                optimizedTagMap2.put(new TagNameSize(curTagNameBase + "." + curTagFragmentNumber, curTagSize),
                    new DefaultPlcTagItem<>(
                        new S7TagChunk(curTag.getDataType(), curTag.getMemoryArea(), curTag.getBlockNumber(), curTagOffset, (byte) 0, curTagSize,
                            (curTag instanceof S7TagChunk) ? ((S7TagChunk) curTag).getChunkTags() : Collections.singletonMap(curTag, tagNameSize.getTagName()),
                            curPartIndex, totalPartCount, curTagSize)));
            }
            // Just copy over tags that fit into a request.
            else {
                optimizedTagMap2.put(tagNameSize, optimizedTagMap.get(tagNameSize));
            }
        }

        // Using the First Fit Decreasing (FFD) bin-packing algorithm try to find the ideal
        // packing for utilizing request sizes.
        // 1. Assign a size to each tag
        // 2. Sort the tags by size (biggest first)
        // 3. Repeat this, till all tags are consumed
        //      1. Take the first packet of the list
        //      2. If the tag itself exceeds the max request size, keep on splitting it into chunks until
        //         the rest would fit into a request. Then proceed with the rest as if it was a normal tag
        //      2. Go through the existing list of requests and check if the current tag would fit
        //          1. If it fits, add it to the request
        //          2. If it doesn't fit go to the next request and check
        //          3. If you reach the end, and it didn't fit any of the previous requests, add a new one
        LinkedHashMap<String, PlcTagItem<PlcTag>> executableTagMap = new LinkedHashMap<>();
        for (TagNameSize tagNameSize : optimizedTagMap2.keySet()) {
            // TODO: Implement the algorithm above.
            executableTagMap.put(tagNameSize.getTagName(), optimizedTagMap2.get(tagNameSize));
        }

        return super.processReadRequest(
            new DefaultPlcReadRequest(((DefaultPlcReadRequest) readRequest).getReader(), executableTagMap),
            driverContext);
    }

    protected PlcReadResponse processReadResponses(PlcReadRequest readRequest, Map<PlcReadRequest, SubResponse<PlcReadResponse>> readResponses, DriverContext driverContext) {
        S7DriverContext s7DriverContext = (S7DriverContext) driverContext;
        if(!s7DriverContext.isEnableBlockReadOptimizer()) {
            return super.processReadResponses(readRequest, readResponses, driverContext);
        }

        // Create the merged read request so the layer down can correctly process it.
        LinkedHashMap<String, PlcTagItem<PlcTag>> tags = new LinkedHashMap<>();
        for (PlcReadRequest plcReadRequest : readResponses.keySet()) {
            for (String tagName : plcReadRequest.getTagNames()) {
                PlcTag tag = plcReadRequest.getTag(tagName);
                tags.put(tagName, new DefaultPlcTagItem<>(tag));
            }
        }

        // Have the upstream optimizer handle its thing.
        PlcReadResponse rawReadResponse = super.processReadResponses(new DefaultPlcReadRequest(((DefaultPlcReadRequest) readRequest).getReader(), tags), readResponses, driverContext);

        // Merge together split-up chunks.
        LinkedHashMap<String, PlcTagItem<PlcTag>> mergedTagItems = new LinkedHashMap<>();
        Map<String, PlcResponseItem<PlcValue>> mergedValues = new LinkedHashMap<>();
        for (String tagName : rawReadResponse.getTagNames()) {
            PlcTag tag = rawReadResponse.getRequest().getTag(tagName);

            // This indicates it's a chunk that's been split up.
            if((tag instanceof S7TagChunk) && (((S7TagChunk) tag).getTotalPartNumber() > 1)){
                S7TagChunk tagChunk = (S7TagChunk) tag;
                // Only handle the first part
                if(tagChunk.getCurPartIndex() == 0) {
                    String tagBaseName = tagName.substring(0, tagName.indexOf("."));

                    // Allocate a chunk of memory able to store the data of all parts.
                    byte[] chunkData = new byte[tagChunk.getTotalPartSize()];

                    int firstPartOffset = tagChunk.getByteOffset();
                    // Copy the payload of each chunk into the total chunk data.
                    for(int curPartIndex = 0; curPartIndex < tagChunk.getTotalPartNumber(); curPartIndex++) {
                        String curPartTagName = tagBaseName + "." + curPartIndex;
                        S7TagChunk curPartTag = (S7TagChunk) rawReadResponse.getRequest().getTag(curPartTagName);

                        PlcResponseCode responseCode = rawReadResponse.getResponseCode(curPartTagName);
                        // If one of the parts is not OK, all in this chunk will be not OK too.
                        if(responseCode != PlcResponseCode.OK) {
                            mergedValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, new PlcNull()));
                            break;
                        }
                        PlcValue plcValue = rawReadResponse.getPlcValue(curPartTagName);
                        byte[] chunkPartData = plcValue.getRaw();
                        int partDataOffset = curPartTag.getByteOffset() - firstPartOffset;
                        System.arraycopy(chunkPartData, 0, chunkData, partDataOffset, chunkPartData.length);
                    }

                    // Store the merged tag.
                    mergedTagItems.put(tagBaseName, new DefaultPlcTagItem<>(
                        new S7TagChunk(tagChunk.getDataType(), tagChunk.getMemoryArea(), tagChunk.getBlockNumber(),
                            firstPartOffset, (byte) 0, tagChunk.getTotalPartSize(), tagChunk.getChunkTags(),
                            0, 1, tagChunk.getTotalPartSize())));
                    mergedValues.put(tagBaseName, new DefaultPlcResponseItem<>(PlcResponseCode.OK, new PlcRawByteArray(chunkData)));
                }
            }
            // All others are just un-split chunks
            else {
                PlcResponseCode responseCode = rawReadResponse.getResponseCode(tagName);
                PlcValue plcValue = rawReadResponse.getPlcValue(tagName);
                mergedTagItems.put(tagName, new DefaultPlcTagItem<>(tag));
                mergedValues.put(tagName, new DefaultPlcResponseItem<>(responseCode, plcValue));
            }
        }
        PlcReadResponse mergedReadResponse = new DefaultPlcReadResponse(new DefaultPlcReadRequest(((DefaultPlcReadRequest)rawReadResponse.getRequest()).getReader(), mergedTagItems), mergedValues);

        // If a Tag is a normal tag, just copy it over. However, if it's a S7TagChunk, process it.
        Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
        for (String tagName : mergedReadResponse.getTagNames()) {
            PlcResponseCode responseCode = mergedReadResponse.getResponseCode(tagName);
            PlcValue plcValue = mergedReadResponse.getPlcValue(tagName);
            // We need to get the request from the response as only this contains the
            // chunk tags and not the original one passed in via readRequest parameter.
            PlcTag tag = mergedReadResponse.getRequest().getTag(tagName);

            // If it's not a tag-chunk, then we just add it to the response.
            if(!(tag instanceof S7TagChunk)) {
                values.put(tagName, new DefaultPlcResponseItem<>(responseCode, plcValue));
                continue;
            }

            // Otherwise it's a tag-chunk and we need to decode the payload.
            // So we decode all the tags that are chunked-together in this.
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(plcValue.getRaw());
            S7TagChunk s7TagChunk = (S7TagChunk) tag;
            int chunkByteOffset = s7TagChunk.getByteOffset();
            for (PlcTag plcTag : s7TagChunk.getChunkTags().keySet()) {
                S7Tag s7Tag = (S7Tag) plcTag;
                String curTagName = s7TagChunk.getChunkTags().get(plcTag);
                int curTagStartPosition = s7Tag.getByteOffset() - chunkByteOffset;
                int curTagDataSize = s7Tag.getDataType().getSizeInBytes() * s7Tag.getNumberOfElements();
                if(s7Tag instanceof S7StringFixedLengthTag) {
                    S7StringFixedLengthTag s7StringFixedLengthTag = (S7StringFixedLengthTag) s7Tag;
                    if(s7Tag.getDataType() == TransportSize.WSTRING) {
                        //curTagStartPosition += 4;
                        curTagDataSize = s7StringFixedLengthTag.getStringLength() * 2;
                    } else {
                        //curTagStartPosition += 2;
                        curTagDataSize = s7StringFixedLengthTag.getStringLength();
                    }
                }
                byte[] curTagData = readBuffer.getBytes(curTagStartPosition, curTagStartPosition + curTagDataSize);
                PlcValue tagValue = parsePlcValue(s7Tag, curTagData, s7DriverContext);
                values.put(curTagName, new DefaultPlcResponseItem<>(responseCode, tagValue));
            }
        }

        return new DefaultPlcReadResponse(readRequest, values);
    }

    protected Map<PlcTag, String> createSortedTagMap() {
        return new TreeMap<>((tag1, tag2) -> {
            if (tag1 instanceof S7SzlTag) {
                S7SzlTag s7Tag1 = (S7SzlTag) tag1;
                S7SzlTag s7Tag2 = (S7SzlTag) tag2;
                if (s7Tag1.getSzlId() == s7Tag2.getSzlId()) {
                    return s7Tag1.getIndex() - s7Tag2.getIndex();
                }
                return s7Tag1.getSzlId() - s7Tag2.getSzlId();
            } else if (tag1 instanceof S7ClkTag) {
                // Technically CLK tags should be identical as there's
                // only one address for reading the PLC clock information.
                return 0;
            } else if (tag1 instanceof S7Tag) {
                S7Tag s7Tag1 = (S7Tag) tag1;
                S7Tag s7Tag2 = (S7Tag) tag2;
                if (s7Tag1.getByteOffset() == s7Tag2.getByteOffset()) {
                    return s7Tag1.getBitOffset() - s7Tag2.getBitOffset();
                }
                return s7Tag1.getByteOffset() - s7Tag2.getByteOffset();
            }
            return 0;
        });
    }

    private PlcValue parsePlcValue(S7Tag tag, byte[] data, S7DriverContext s7DriverContext) {
        ReadBuffer readBuffer = new ReadBufferByteBased(data);
        try {
            int stringLength = (tag instanceof S7StringFixedLengthTag) ? ((S7StringFixedLengthTag) tag).getStringLength() : 254;
            if (tag.getNumberOfElements() == 1) {
                return DataItem.staticParse(readBuffer, tag.getDataType().getDataProtocolId(),
                    s7DriverContext.getControllerType(), stringLength);
            } else {
                // In case of reading an array of bytes, make use of our simpler PlcRawByteArray as the user is
                // probably expecting to process the read raw data.
                if(tag.getDataType() == TransportSize.BYTE) {
                    return new PlcRawByteArray(data);
                } else {
                    // Fetch all
                    final PlcValue[] resultItems = IntStream.range(0, tag.getNumberOfElements()).mapToObj(i -> {
                        try {
                            return DataItem.staticParse(readBuffer, tag.getDataType().getDataProtocolId(),
                                s7DriverContext.getControllerType(), stringLength);
                        } catch (ParseException e) {
                            logger.warn("Error parsing tag item of type: '{}' (at position {}})", tag.getDataType().name(), i, e);
                        }
                        return null;
                    }).toArray(PlcValue[]::new);
                    return DefaultPlcValueHandler.of(tag, resultItems);
                }
            }
        } catch (ParseException e) {
            logger.warn("Error parsing tag item of type: '{}'", tag.getDataType().name(), e);
        }
        return null;
    }

    public static class S7TagChunk extends S7Tag {

        private final Map<PlcTag, String> chunkTags;
        private final int curPartIndex;
        private final int totalPartNumber;
        private final int totalPartSize;

        public S7TagChunk(TransportSize dataType, MemoryArea memoryArea, int blockNumber, int byteOffset, byte bitOffset, int numElements, Map<PlcTag, String> chunkTags, int curPartIndex, int totalPartNumber, int totalPartSize) {
            super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
            this.chunkTags = chunkTags;
            this.curPartIndex = curPartIndex;
            this.totalPartNumber = totalPartNumber;
            this.totalPartSize = totalPartSize;
        }

        public Map<PlcTag, String> getChunkTags() {
            return chunkTags;
        }

        public int getCurPartIndex() {
            return curPartIndex;
        }

        public int getTotalPartNumber() {
            return totalPartNumber;
        }

        public int getTotalPartSize() {
            return totalPartSize;
        }
    }

    public static class TagNameSize implements Comparable<TagNameSize> {
        private final String tagName;
        private final int tagSize;

        public TagNameSize(String tagName, int tagSize) {
            this.tagName = tagName;
            this.tagSize = tagSize;
        }

        public String getTagName() {
            return tagName;
        }

        public int getTagSize() {
            return tagSize;
        }

        @Override
        public int compareTo(TagNameSize o) {
            int sizeComparison = Integer.compare(tagSize, o.tagSize);
            if(sizeComparison != 0) {
                return sizeComparison;
            }
            return tagName.compareTo(o.tagName);
        }

    }

}
