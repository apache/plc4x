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
package org.apache.plc4x.java.spi.messages;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcTagRequest;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithRenderAsList;

public class DefaultPlcReadRequest implements PlcReadRequest, PlcTagRequest, Serializable {

    private final PlcReader reader;
    // This is intentionally a linked hash map in order to keep the order of how elements were added.
    private final LinkedHashMap<String, PlcTag> tags;

    public DefaultPlcReadRequest(PlcReader reader,
                                 LinkedHashMap<String, PlcTag> tags) {
        this.reader = reader;
        this.tags = tags;
    }

    @Override
    public CompletableFuture<PlcReadResponse> execute() {
        return reader.read(this);
    }

    @Override
    public int getNumberOfTags() {
        return tags.size();
    }

    @Override
    public LinkedHashSet<String> getTagNames() {
        // TODO: Check if this already is a LinkedHashSet.
        return new LinkedHashSet<>(tags.keySet());
    }

    @Override
    public PlcTag getTag(String name) {
        return tags.get(name);
    }

    @Override
    public List<PlcTag> getTags() {
        return new LinkedList<>(tags.values());
    }

    public PlcReader getReader() {
        return reader;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcReadRequest");

        writeBuffer.pushContext("PlcTagRequest");
        writeBuffer.pushContext("tags", WithRenderAsList(true));
        for (Map.Entry<String, PlcTag> tagEntry : tags.entrySet()) {
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcTag tag = tagEntry.getValue();
            if(!(tag instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Tag doesn't implement Serializable");
            }
            ((Serializable) tag).serialize(writeBuffer);
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("tags");
        writeBuffer.popContext("PlcTagRequest");

        writeBuffer.popContext("PlcReadRequest");
    }

    public static class Builder implements PlcReadRequest.Builder {

        private final PlcReader reader;
        private final PlcTagHandler tagHandler;
        private final Map<String, Supplier<PlcTag>> tags;

        public Builder(PlcReader reader, PlcTagHandler tagHandler) {
            this.reader = reader;
            this.tagHandler = tagHandler;
            tags = new LinkedHashMap<>();
        }

        @Override
        public PlcReadRequest.Builder addTagAddress(String name, String tagAddress) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, () -> tagHandler.parseTag(tagAddress));
            return this;
        }

        @Override
        public PlcReadRequest.Builder addTag(String name, PlcTag tag) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, () -> tag);
            return this;
        }

        @Override
        public PlcReadRequest build() {
            LinkedHashMap<String, PlcTag> parsedTags = new LinkedHashMap<>();
            tags.forEach((name, tagSupplier) -> parsedTags.put(name, tagSupplier.get()));
            return new DefaultPlcReadRequest(reader, parsedTags);
        }

        @Override
        public Optional<PlcTag> parseTagAddressSafe(String tagAddress){
            PlcTag plcTag;
            try{
                plcTag = tagHandler.parseTag(tagAddress);
            }catch (PlcInvalidTagException e){
                return Optional.empty();
            }
            return Optional.ofNullable(plcTag);
        }
    }

}
