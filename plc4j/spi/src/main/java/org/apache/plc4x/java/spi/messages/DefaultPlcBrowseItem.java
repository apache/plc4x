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
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeInfo;
import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcBrowseItem implements PlcBrowseItem, Serializable {

    private final PlcTag tag;

    private final String name;

    private final boolean readable;
    private final boolean writable;
    private final boolean subscribable;

    private final Map<String, PlcBrowseItem> children;

    private final Map<String, PlcValue> options;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcBrowseItem(@JsonProperty("tag") PlcTag tag,
                                @JsonProperty("name") String name,
                                @JsonProperty("readable") boolean readable,
                                @JsonProperty("writable") boolean writable,
                                @JsonProperty("subscribable") boolean subscribable,
                                @JsonProperty("children") Map<String, PlcBrowseItem> children,
                                @JsonProperty("options") Map<String, PlcValue> options) {
        this.tag = tag;
        this.name = name;
        this.readable = readable;
        this.writable = writable;
        this.subscribable = subscribable;
        this.children = children;
        this.options = options;
    }

    @Override
    public PlcTag getTag() {
        return tag;
    }

    @Override
    public String getName() {
        return name;
    }

    public boolean isReadable() {
        return readable;
    }

    public boolean isWritable() {
        return writable;
    }

    public boolean isSubscribable() {
        return subscribable;
    }

    @Override
    public Map<String, PlcBrowseItem> getChildren() {
        return children;
    }

    @Override
    public Map<String, PlcValue> getOptions() {
        return options;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());
        writeBuffer.writeString("address",
            tag.getAddressString().getBytes(StandardCharsets.UTF_8).length * 8,
            tag.getAddressString(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        writeBuffer.writeString("name",
            name.getBytes(StandardCharsets.UTF_8).length * 8,
            name, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        // TODO: Find out how to serialize an enum.
        //writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);
        if(children != null && !children.isEmpty()) {
            writeBuffer.pushContext("children");
            for (PlcBrowseItem child : children.values()) {
                writeBuffer.pushContext("child");
                ((DefaultPlcBrowseItem) child).serialize(writeBuffer);
                writeBuffer.popContext("child");
            }
            writeBuffer.popContext("children");
        }
        if(options != null && !options.isEmpty()) {
            writeBuffer.pushContext("options");
            for (Map.Entry<String, PlcValue> optionEntry : options.entrySet()) {
                writeBuffer.pushContext("option");
                writeBuffer.writeString("name",
                    optionEntry.getKey().getBytes(StandardCharsets.UTF_8).length * 8,
                    optionEntry.getKey(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
                // TODO: Find out how to serialize a PlcValue
                //writeBuffer.writeString("value", optionEntry.getValue().getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), optionEntry.getValue());
                ((DefaultPlcBrowseItem) optionEntry).serialize(writeBuffer);
                writeBuffer.popContext("option");
            }
            writeBuffer.popContext("options");
        }
        writeBuffer.popContext(getClass().getSimpleName());
    }

}
