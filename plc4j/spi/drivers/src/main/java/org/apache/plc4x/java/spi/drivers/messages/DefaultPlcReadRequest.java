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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.messages.PlcTagRequest;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.drivers.functions.PlcReader;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagErrorItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class DefaultPlcReadRequest implements PlcReadRequest, PlcTagRequest {

    private static final Logger LOGGER = LoggerFactory.getLogger(DefaultPlcReadRequest.class);

    private final PlcReader reader;
    private final LinkedHashMap<String, PlcTagItem<PlcTag>> tags;

    public DefaultPlcReadRequest(PlcReader reader, LinkedHashMap<String, PlcTagItem<PlcTag>> tags) {
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
        return new LinkedHashSet<>(tags.keySet());
    }

    public PlcTagItem<PlcTag> getTagItem(String tagName) {
        return tags.get(tagName);
    }

    @Override
    public PlcTag getTag(String tagName) {
        return tags.get(tagName) != null ? tags.get(tagName).getTag() : null;
    }

    @Override
    public PlcResponseCode getTagResponseCode(String tagName) {
        return tags.get(tagName) != null ? tags.get(tagName).getResponseCode() : null;
    }

    @Override
    public List<PlcTag> getTags() {
        return tags.values().stream()
            .filter(item -> item instanceof DefaultPlcTagItem<PlcTag>)
            .map(PlcTagItem::getTag)
            .collect(Collectors.toList());
    }

    public PlcReader getReader() {
        return reader;
    }

    public static class Builder implements PlcReadRequest.Builder {

        private final PlcReader reader;
        private final PlcTagHandler tagHandler;
        private final Map<String, Supplier<PlcTagItem<PlcTag>>> tagItems;

        public Builder(PlcReader reader, PlcTagHandler tagHandler) {
            this.reader = Objects.requireNonNull(reader);
            this.tagHandler = tagHandler;
            tagItems = new LinkedHashMap<>();
        }

        @Override
        public PlcReadRequest.Builder addTagAddress(String name, String tagAddress) {
            Objects.requireNonNull(tagHandler, "tagHandler must not be null");
            if (tagItems.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tagItems.put(name, () -> {
                try {
                    PlcTag tag = tagHandler.parseTag(tagAddress);
                    return new DefaultPlcTagItem<>(tag);
                } catch (Exception e) {
                    // The tag still takes part in the request, carrying INVALID_ADDRESS. Without
                    // this log line the user is left with a response code and no way of telling
                    // which address was rejected or why.
                    LOGGER.warn("Invalid address '{}' for tag '{}': {}", tagAddress, name, e.getMessage(), e);
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
            tagItems.forEach((name, tagItemSupplier) -> parsedTags.put(name, tagItemSupplier.get()));
            return new DefaultPlcReadRequest(reader, parsedTags);
        }
    }

}
