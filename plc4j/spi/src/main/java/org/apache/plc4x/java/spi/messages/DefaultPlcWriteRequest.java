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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.messages.utils.TagValueItem;

import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcWriteRequest implements PlcWriteRequest, Serializable {

    private final PlcWriter writer;

    private final LinkedHashMap<String, TagValueItem> tags;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcWriteRequest(@JsonProperty("writer") PlcWriter writer,
                                  @JsonProperty("tags") LinkedHashMap<String, TagValueItem> tags) {
        this.writer = writer;
        this.tags = tags;
    }

    @Override
    @JsonIgnore
    public CompletableFuture<PlcWriteResponse> execute() {
        return writer.write(this);
    }

    @Override
    @JsonIgnore
    public int getNumberOfTags() {
        return tags.size();
    }

    @Override
    @JsonIgnore
    public LinkedHashSet<String> getTagNames() {
        // TODO: Check if this already is a LinkedHashSet.
        return new LinkedHashSet<>(tags.keySet());
    }

    @Override
    @JsonIgnore
    public PlcTag getTag(String name) {
        return tags.get(name).getTag();
    }

    @Override
    @JsonIgnore
    public List<PlcTag> getTags() {
        return tags.values().stream().map(TagValueItem::getTag).collect(Collectors.toCollection(LinkedList::new));
    }

    @Override
    @JsonIgnore
    public PlcValue getPlcValue(String name) {
        return tags.get(name).getValue();
    }

    @JsonIgnore
    public List<PlcValue> getPlcValues() {
        return tags.values().stream().map(TagValueItem::getValue).collect(Collectors.toCollection(LinkedList::new));
    }

    public PlcWriter getWriter() {
        return writer;
    }

    @Override
    @JsonIgnore
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

        writeBuffer.pushContext("tags");
        for (Map.Entry<String, TagValueItem> tagEntry : tags.entrySet()) {
            TagValueItem tagValueItem = tagEntry.getValue();
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcTag tag = tagValueItem.getTag();
            if (!(tag instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Tag doesn't implement XmlSerializable");
            }
            ((Serializable) tag).serialize(writeBuffer);
            final PlcValue value = tagValueItem.getValue();
            if (value instanceof PlcList) {
                PlcList list = (PlcList) value;
                for (PlcValue plcValue : list.getList()) {
                    String plcValueString = plcValue.getString();
                    writeBuffer.writeString("value",
                        plcValueString.getBytes(StandardCharsets.UTF_8).length * 8,
                        plcValueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
                }
            } else {
                String plcValueString = value.getString();
                writeBuffer.writeString("value",
                    plcValueString.getBytes(StandardCharsets.UTF_8).length * 8,
                    plcValueString, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
            }
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("tags");

        writeBuffer.popContext("PlcWriteRequest");
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcWriter writer;
        private final PlcTagHandler tagHandler;
        private final PlcValueHandler valueHandler;
        private final Map<String, Pair<Supplier<PlcTag>, Object[]>> tags;

        public Builder(PlcWriter writer, PlcTagHandler tagHandler, PlcValueHandler valueHandler) {
            this.writer = writer;
            this.tagHandler = tagHandler;
            this.valueHandler = valueHandler;
            tags = new TreeMap<>();
        }

        @Override
        public Builder addTagAddress(String name, String tagAddress, Object... values) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, Pair.of(() -> tagHandler.parseTag(tagAddress), values));
            return this;
        }

        @Override
        public Builder addTag(String name, PlcTag tag, Object... values) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, Pair.of(() -> tag, values));
            return this;
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, TagValueItem> parsedTags = new LinkedHashMap<>();
            tags.forEach((name, tagValues) -> {
                // Compile the query string.
                PlcTag tag = tagValues.getLeft().get();
                Object[] value = tagValues.getRight();
                PlcValue plcValue = valueHandler.newPlcValue(tag, value);
                if (tag.getPlcValueType() == PlcValueType.NULL && value!=null) {
                    tag.setPlcValueType(plcValue.getPlcValueType());
                }
                parsedTags.put(name, new TagValueItem(tag, plcValue));
            });
            return new DefaultPlcWriteRequest(writer, parsedTags);
        }
    }

}
