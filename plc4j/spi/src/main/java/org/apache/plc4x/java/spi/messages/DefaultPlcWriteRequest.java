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
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.DefaultTagErrorItem;
import org.apache.plc4x.java.spi.messages.utils.DefaultTagValueItem;
import org.apache.plc4x.java.spi.messages.utils.TagValueItem;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcValueHandler;

import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;

import static org.apache.plc4x.java.spi.generation.WithReaderWriterArgs.WithRenderAsList;

public class DefaultPlcWriteRequest implements PlcWriteRequest, Serializable {

    private final PlcWriter writer;

    private final LinkedHashMap<String, TagValueItem> tags;

    public DefaultPlcWriteRequest(PlcWriter writer,
                                  LinkedHashMap<String, TagValueItem> tags) {
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
        // TODO: Check if this already is a LinkedHashSet.
        return new LinkedHashSet<>(tags.keySet());
    }

    @Override
    public PlcResponseCode getTagResponseCode(String tagName) {
        return tags.get(tagName).getResponseCode();
    }

    @Override
    public PlcTag getTag(String name) {
        return tags.get(name).getTag();
    }

    @Override
    public List<PlcTag> getTags() {
        return tags.values().stream().map(TagValueItem::getTag).collect(Collectors.toCollection(LinkedList::new));
    }

    public PlcValue getPlcValue(String name) {
        return tags.get(name).getValue();
    }

    public List<PlcValue> getPlcValues() {
        return tags.values().stream().map(TagValueItem::getValue).collect(Collectors.toCollection(LinkedList::new));
    }

    public PlcWriter getWriter() {
        return writer;
    }

    @Override
    public int getNumberOfValues(String name) {
        final PlcValue value = tags.get(name).getValue();
        if (value instanceof PlcList) {
            PlcList list = (PlcList) value;
            return list.getLength();
        }
        return 1;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcWriteRequest");

        writeBuffer.pushContext("PlcTagRequest");
        writeBuffer.pushContext("tags", WithRenderAsList(true));
        for (Map.Entry<String, TagValueItem> tagEntry : tags.entrySet()) {
            TagValueItem tagValueItem = tagEntry.getValue();
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcTag tag = tagValueItem.getTag();
            if (!(tag instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Tag doesn't implement Serializable");
            }
            ((Serializable) tag).serialize(writeBuffer);
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("tags");
        writeBuffer.popContext("PlcTagRequest");

        writeBuffer.pushContext("values", WithRenderAsList(true));
        for (Map.Entry<String, TagValueItem> tagEntry : tags.entrySet()) {
            TagValueItem tagValueItem = tagEntry.getValue();
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            final PlcValue value = tagValueItem.getValue();
            if (value != null) {
                serializePlcValue(value, writeBuffer);
            }
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("values");

        writeBuffer.popContext("PlcWriteRequest");
    }

    protected void serializePlcValue(PlcValue plcValue, WriteBuffer writeBuffer) throws SerializationException {
        if (plcValue instanceof Serializable) {
            Serializable serializable = (Serializable) plcValue;
            serializable.serialize(writeBuffer);
        } else {
            String plcValueString = plcValue.getString();
            writeBuffer.writeString("value",
                plcValueString.getBytes(StandardCharsets.UTF_8).length * 8,
                plcValueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        }
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcWriter writer;
        private final PlcTagHandler tagHandler;
        private final PlcValueHandler valueHandler;
        private final Map<String, Supplier<TagValueItem>> tagValues;

        public Builder(PlcWriter writer, PlcTagHandler tagHandler, PlcValueHandler valueHandler) {
            this.writer = writer;
            this.tagHandler = tagHandler;
            this.valueHandler = valueHandler;
            tagValues = new TreeMap<>();
        }

        @Override
        public Builder addTagAddress(String name, String tagAddress, Object... values) {
            if (tagValues.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tagValues.put(name, () -> {
                try {
                    PlcTag tag = tagHandler.parseTag(tagAddress);
                    try {
                        PlcValue plcValue = parsePlcValue(tag, values);
                        return new DefaultTagValueItem(tag, plcValue);
                    } catch (Exception e) {
                        return new DefaultTagErrorItem(PlcResponseCode.INVALID_DATA);
                    }
                } catch (Exception e) {
                    return new DefaultTagErrorItem(PlcResponseCode.INVALID_ADDRESS);
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
                    return new DefaultTagValueItem(tag, plcValue);
                } catch (Exception e) {
                    return new DefaultTagErrorItem(PlcResponseCode.INVALID_DATA);
                }
            });
            return this;
        }

        protected PlcValue parsePlcValue(PlcTag tag, Object... values) {
            // If this is more than one element the value itself will definitely be a list.
            if(values.length > 1) {
                List<PlcValue> listItems = new ArrayList<>(values.length);
                for (Object value : values) {
                    if(value instanceof PlcValue) {
                        listItems.add((PlcValue) value);
                    } else {
                        PlcValue plcItemValue = valueHandler.newPlcValue(tag, value);
                        listItems.add(plcItemValue);
                    }
                }
                return new PlcList(listItems);
            }
            // If the values are already PlcValues, just use them.
            else if(values[0] instanceof PlcValue) {
                return (PlcValue) values[0];
            }
            // In all other cases use the value-handler.
            else {
                return valueHandler.newPlcValue(tag, values[0]);
            }
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, TagValueItem> parsedTags = new LinkedHashMap<>();
            tagValues.forEach((name, tagValueItemSupplier) -> {
                TagValueItem tagValueItem = tagValueItemSupplier.get();
                parsedTags.put(name, tagValueItem);
            });
            return new DefaultPlcWriteRequest(writer, parsedTags);
        }
    }

}
