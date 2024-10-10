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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.exceptions.PlcTagNotFoundException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcTagRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagErrorItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithRenderAsList;

public class DefaultPlcReadRequest implements PlcReadRequest, PlcTagRequest, Serializable {

    private final PlcReader reader;
    // This is intentionally a linked hash map in order to keep the order of how elements were added.
    private final LinkedHashMap<String, PlcTagItem<PlcTag>> tags;

    public DefaultPlcReadRequest(PlcReader reader,
                                 LinkedHashMap<String, PlcTagItem<PlcTag>> tags) {
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

    public PlcTagItem<PlcTag> getTagItem(String tagName) {
        return tags.get(tagName);
    }

    @Override
    public PlcTag getTag(String tagName) {
        return tags.get(tagName).getTag();
    }

    @Override
    public PlcResponseCode getTagResponseCode(String tagName) {
        return tags.get(tagName).getResponseCode();
    }

    @Override
    public List<PlcTag> getTags() {
        return tags.values().stream().filter(plcTagItem -> plcTagItem instanceof DefaultPlcTagItem).map(PlcTagItem::getTag).collect(Collectors.toList());
    }

    public PlcReader getReader() {
        return reader;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcReadRequest");

        writeBuffer.pushContext("PlcTagRequest");
        writeBuffer.pushContext("tags", WithRenderAsList(true));
        for (Map.Entry<String, PlcTagItem<PlcTag>> tagEntry : tags.entrySet()) {
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcTagItem<PlcTag> tagItem = tagEntry.getValue();
            if (!(tagItem instanceof Serializable)) {
                throw new RuntimeException("Error serializing. PlcTagItem doesn't implement Serializable");
            }
            ((Serializable) tagItem).serialize(writeBuffer);
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("tags");
        writeBuffer.popContext("PlcTagRequest");

        writeBuffer.popContext("PlcReadRequest");
    }

    public static class Builder implements PlcReadRequest.Builder {

        private final PlcReader reader;
        private final PlcTagHandler tagHandler;
        private final Map<String, Supplier<PlcTagItem<PlcTag>>> tagItems;

        public Builder(PlcReader reader, PlcTagHandler tagHandler) {
            this.reader = Objects.requireNonNull(reader);
            this.tagHandler = Objects.requireNonNull(tagHandler);
            tagItems = new LinkedHashMap<>();
        }

        @Override
        public PlcReadRequest.Builder addTagAddress(String name, String tagAddress) {
            if (tagItems.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tagItems.put(name, () -> {
                try {
                    PlcTag tag = tagHandler.parseTag(tagAddress);
                    return new DefaultPlcTagItem<>(tag);
                } catch (PlcTagNotFoundException e) {
                    return new DefaultPlcTagErrorItem<>(PlcResponseCode.NOT_FOUND);
                } catch (Exception e) {
                    return new DefaultPlcTagErrorItem<>(PlcResponseCode.INVALID_ADDRESS);
                }
            });
            return this;
        }

        @Override
        public PlcReadRequest.Builder addTag(String name, PlcTag tag) {
            if (tagItems.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tagItems.put(name, () -> new DefaultPlcTagItem<>(tag));
            return this;
        }

        @Override
        public PlcReadRequest build() {
            LinkedHashMap<String, PlcTagItem<PlcTag>> parsedTags = new LinkedHashMap<>();
            tagItems.forEach((name, tagItemSupplier) -> {
                PlcTagItem<PlcTag> plcTagItem = tagItemSupplier.get();
                parsedTags.put(name, plcTagItem);
            });
            return new DefaultPlcReadRequest(reader, parsedTags);
        }
    }

}
