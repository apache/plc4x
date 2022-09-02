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
import org.apache.plc4x.java.api.messages.PlcDiscoveryItem;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcBrowseItem implements PlcBrowseItem, Serializable {

    private final String address;

    private final String name;

    private final PlcValueType dataType;

    private final boolean readable;
    private final boolean writable;
    private final boolean subscribable;

    private final List<PlcBrowseItem> children;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcBrowseItem(@JsonProperty("address") String address,
                                @JsonProperty("name") String name,
                                @JsonProperty("dataType") PlcValueType dataType,
                                @JsonProperty("readable") boolean readable,
                                @JsonProperty("writable") boolean writable,
                                @JsonProperty("subscribable") boolean subscribable,
                                @JsonProperty("children") List<PlcBrowseItem> children) {
        this.address = address;
        this.name = name;
        this.dataType = dataType;
        this.readable = readable;
        this.writable = writable;
        this.subscribable = subscribable;
        this.children = children;
    }

    @Override
    public String getAddress() {
        return address;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return dataType;
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
    public List<PlcBrowseItem> getChildren() {
        return children;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());
        writeBuffer.writeString("address", address.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), address);
        writeBuffer.writeString("name", name.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), name);
        // TODO: Find out how to serialize an enum.
        //writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);
        if(children != null && !children.isEmpty()) {
            writeBuffer.pushContext("children");
            for (PlcBrowseItem child : children) {
                writeBuffer.pushContext("child");
                ((DefaultPlcBrowseItem) child).serialize(writeBuffer);
                writeBuffer.popContext("child");
            }
            writeBuffer.popContext("children");
        }
        writeBuffer.popContext(getClass().getSimpleName());
    }

}
