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
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.drivers.functions.PlcWriter;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagErrorItem;
import org.apache.plc4x.java.spi.drivers.messages.items.DefaultPlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.messages.items.PlcTagValueItem;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;

public class DefaultPlcWriteRequest implements PlcWriteRequest {

    private final PlcWriter writer;
    private final LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags;

    public DefaultPlcWriteRequest(PlcWriter writer, LinkedHashMap<String, PlcTagValueItem<PlcTag>> tags) {
        this.writer = writer;
        this.tags = tags;
    }

    @Override
    public CompletableFuture<PlcWriteResponse> execute() {
        return writer.write(this);
    }

    @Override
    public int getNumberOfTags() {
        return tags.size();
    }

    @Override
    public LinkedHashSet<String> getTagNames() {
        return new LinkedHashSet<>(tags.keySet());
    }

    public PlcTagValueItem<PlcTag> getTagValueItem(String tagName) {
        return tags.get(tagName);
    }

    @Override
    public PlcResponseCode getTagResponseCode(String tagName) {
        return tags.get(tagName) != null ? tags.get(tagName).getResponseCode() : null;
    }

    @Override
    public PlcTag getTag(String tagName) {
        return tags.get(tagName) != null ? tags.get(tagName).getTag() : null;
    }

    @Override
    public List<PlcTag> getTags() {
        return tags.values().stream().map(PlcTagValueItem::getTag).collect(Collectors.toCollection(LinkedList::new));
    }

    public PlcValue getPlcValue(String name) {
        return tags.get(name).getValue();
    }

    @Override
    public int getNumberOfValues(String name) {
        final PlcValue value = tags.get(name).getValue();
        if (value instanceof PlcList list) {
            return list.getLength();
        }
        return 1;
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcWriter writer;
        private final PlcTagHandler tagHandler;
        private final PlcValueHandler valueHandler;
        private final Map<String, Supplier<PlcTagValueItem<PlcTag>>> tagValues;

        public Builder(PlcWriter writer, PlcTagHandler tagHandler, PlcValueHandler valueHandler) {
            this.writer = Objects.requireNonNull(writer);
            this.tagHandler = tagHandler;
            this.valueHandler = valueHandler;
            tagValues = new TreeMap<>();
        }

        @Override
        public Builder addTagAddress(String name, String tagAddress, Object... values) {
            Objects.requireNonNull(tagHandler, "tagHandler must not be null");
            if (tagValues.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tagValues.put(name, () -> {
                try {
                    PlcTag tag = tagHandler.parseTag(tagAddress);
                    try {
                        PlcValue plcValue = parsePlcValue(tag, values);
                        return new DefaultPlcTagValueItem<>(tag, plcValue);
                    } catch (Exception e) {
                        return new DefaultPlcTagErrorItem<>(PlcResponseCode.INVALID_DATA);
                    }
                } catch (Exception e) {
                    return new DefaultPlcTagErrorItem<>(PlcResponseCode.INVALID_ADDRESS);
                }
            });
            return this;
        }

        @Override
        public Builder addTag(String name, PlcTag tag, Object... values) {
            if (tagValues.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tagValues.put(name, () -> {
                try {
                    PlcValue plcValue = parsePlcValue(tag, values);
                    return new DefaultPlcTagValueItem<>(tag, plcValue);
                } catch (Exception e) {
                    return new DefaultPlcTagErrorItem<>(PlcResponseCode.INVALID_DATA);
                }
            });
            return this;
        }

        private PlcValue parsePlcValue(PlcTag tag, Object[] values) {
            if ((values == null) || (values.length != 1)) {
                throw new PlcRuntimeException("Expecting 1 item, but got " + (values != null ? values.length : "null"));
            }
            if (values[0] instanceof PlcValue) {
                return (PlcValue) values[0];
            }
            Objects.requireNonNull(valueHandler, "valueHandler must not be null");
            return valueHandler.newPlcValue(tag, values);
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, PlcTagValueItem<PlcTag>> parsedTags = new LinkedHashMap<>();
            tagValues.forEach((name, supplier) -> parsedTags.put(name, supplier.get()));
            return new DefaultPlcWriteRequest(writer, parsedTags);
        }
    }

}
