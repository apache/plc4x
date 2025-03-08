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
package org.apache.plc4x.java.s7.readwrite.optimizer;

import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.s7.readwrite.*;
import org.apache.plc4x.java.s7.readwrite.context.S7DriverContext;
import org.apache.plc4x.java.s7.readwrite.tag.*;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;
import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.spi.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.spi.messages.utils.*;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
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
        Map<String, Map<PlcTag, String>> sortedTags = new HashMap<>();
        for (String tagName : readRequest.getTagNames()) {
            PlcTag tag = readRequest.getTag(tagName);
            if(tag instanceof S7SzlTag) {
                if(!sortedTags.containsKey("SZL")) {
                    sortedTags.put("SZL", createTagMap());
                }
                sortedTags.get("SZL").put(tag, tagName);
            } else if(tag instanceof S7ClkTag) {
                if(!sortedTags.containsKey("CLK")) {
                    sortedTags.put("CLK", createTagMap());
                }
                sortedTags.get("CLK").put(tag, tagName);
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
                if(!sortedTags.containsKey(areaName)) {
                    sortedTags.put(areaName, createTagMap());
                }
                sortedTags.get(areaName).put(tag, tagName);
            } else {
                System.out.println("Ignored");
            }
        }

        // Build groups of byte-arrays to fetch
        // Strategies:
        // - If two consecutive tags are more than S7_ADDRESS_ANY_SIZE bytes apart,
        //   reading them in one block is less efficient than reading them separately.
        //   Other than that, extend the byte array to include the next tag.
        LinkedHashMap<String, PlcTagItem<PlcTag>> optimizedTagMap = new LinkedHashMap<>();
        for (Map<PlcTag, String> tagList : sortedTags.values()) {
            MemoryArea currentMemoryArea = null;
            int currentDataBlockNumber = -1;
            int currentChunkStartByteOffset = -1;
            int currentChunkEndByteOffset = -1;
            Map<PlcTag, String> currentChunkTags = createTagMap();
            for (PlcTag plcTag : tagList.keySet()) {
                // We don't do anything for Szl and Clk tags.
                if(plcTag instanceof S7SzlTag) {
                    optimizedTagMap.put(tagList.get(plcTag), new DefaultPlcTagItem<>(plcTag));
                } else if (plcTag instanceof S7ClkTag) {
                    optimizedTagMap.put(tagList.get(plcTag), new DefaultPlcTagItem<>(plcTag));
                }
                // Var-length strings, are a performance nightmare. Trying to optimize reading them is probably not
                // worth the effort. For now, we simply handle them as un-chunked tags.
                else if(plcTag instanceof S7StringVarLengthTag) {
                    optimizedTagMap.put(tagList.get(plcTag), new DefaultPlcTagItem<>(plcTag));
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
                        curTagSize = (2 + (stringFixedLengthTag.getStringLength() * bytesPerChar)) * s7Tag.getNumberOfElements();
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
                        optimizedTagMap.put("__chunk__" + optimizedTagMap.size(),
                            new DefaultPlcTagItem<>(
                                new S7TagChunk(TransportSize.BYTE, currentMemoryArea, currentDataBlockNumber,
                                    currentChunkStartByteOffset, (byte) 0,
                                    currentChunkEndByteOffset - currentChunkStartByteOffset,
                                    currentChunkTags)));

                        // Start a new one.
                        currentChunkStartByteOffset = s7Tag.getByteOffset();
                        currentChunkEndByteOffset = s7Tag.getByteOffset() + curTagSize;
                        currentChunkTags = createTagMap();
                    }
                    // Otherwise extend the array size to include this tag.
                    else {
                        currentChunkEndByteOffset = s7Tag.getByteOffset() + curTagSize;
                    }

                    // Add the tag to the list of tags for the current chunk.
                    currentChunkTags.put(s7Tag, tagList.get(s7Tag));
                }
            }
            // Add the new tag.
            // Save the current chunk.
            optimizedTagMap.put("__chunk__" + optimizedTagMap.size(),
                new DefaultPlcTagItem<>(
                    new S7TagChunk(TransportSize.BYTE, currentMemoryArea, currentDataBlockNumber,
                        currentChunkStartByteOffset, (byte) 0,
                        currentChunkEndByteOffset - currentChunkStartByteOffset,
                        currentChunkTags)));
        }

        return super.processReadRequest(
            new DefaultPlcReadRequest(((DefaultPlcReadRequest) readRequest).getReader(), optimizedTagMap),
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
        PlcReadResponse readResponse = super.processReadResponses(new DefaultPlcReadRequest(((DefaultPlcReadRequest) readRequest).getReader(), tags), readResponses, driverContext);

        // If a Tag is a normal tag, just copy it over. However, if it's a S7TagChunk, process it.
        Map<String, PlcResponseItem<PlcValue>> values = new HashMap<>();
        for (String tagName : readResponse.getTagNames()) {
            PlcResponseCode responseCode = readResponse.getResponseCode(tagName);
            PlcValue plcValue = readResponse.getPlcValue(tagName);
            // We need to get the request from the response as only this contains the
            // chunk tags and not the original one passed in via readRequest parameter.
            PlcTag tag = readResponse.getRequest().getTag(tagName);

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

    protected Map<PlcTag, String> createTagMap() {
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
                // TODO: Pass the type of plc into the parse function ...
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

        public S7TagChunk(TransportSize dataType, MemoryArea memoryArea, int blockNumber, int byteOffset, byte bitOffset, int numElements, Map<PlcTag, String> chunkTags) {
            super(dataType, memoryArea, blockNumber, byteOffset, bitOffset, numElements);
            this.chunkTags = chunkTags;
        }

        public Map<PlcTag, String> getChunkTags() {
            return chunkTags;
        }
    }

}
