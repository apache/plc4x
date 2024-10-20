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

import org.apache.plc4x.java.api.messages.PlcBrowseItem;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;

public class DefaultPlcBrowseItem implements PlcBrowseItem, Serializable {

    private final PlcTag tag;

    private final String name;

    private final boolean readable;
    private final boolean writable;
    private final boolean subscribable;
    private final boolean publishable;

    private final List<ArrayInfo> arrayInformation;
    private final Map<String, PlcBrowseItem> children;
    private final Map<String, PlcValue> options;

    public DefaultPlcBrowseItem(PlcTag tag,
                                String name,
                                boolean readable,
                                boolean writable,
                                boolean subscribable,
                                boolean publishable,
                                List<ArrayInfo> arrayInformation,
                                Map<String, PlcBrowseItem> children,
                                Map<String, PlcValue> options) {
        this.tag = tag;
        this.name = name;
        this.readable = readable;
        this.writable = writable;
        this.subscribable = subscribable;
        this.publishable = publishable;
        this.arrayInformation = arrayInformation;
        this.children = children;
        this.options = options;
    }

    protected DefaultPlcBrowseItem(PlcBrowseItem original, Map<String, PlcBrowseItem> children) {
        this.tag = original.getTag();
        this.name = original.getName();
        this.readable = original.isReadable();
        this.writable = original.isWritable();
        this.subscribable = original.isSubscribable();
        this.publishable = original.isPublishable();
        this.arrayInformation = original.getArrayInformation();
        this.options = original.getOptions();
        this.children = children;
    }

    @Override
    public PlcTag getTag() {
        return tag;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public boolean isReadable() {
        return readable;
    }

    @Override
    public boolean isWritable() {
        return writable;
    }

    @Override
    public boolean isSubscribable() {
        return subscribable;
    }

    @Override
    public boolean isPublishable() {
        return publishable;
    }

    @Override
    public boolean isArray() {
        return false;
    }

    @Override
    public List<ArrayInfo> getArrayInformation() {
        return arrayInformation;
    }

    @Override
    public Map<String, PlcBrowseItem> getChildren() {
        return children;
    }

    @Override
    public Map<String, PlcValue> getOptions() {
        return options;
    }

    /**
     * For simple non-array elements we usually don't have to do anything here
     * @return a one-element list the unchanged item
     */
    /*@Override
    public PlcBrowseItem resolveArrayItems() {
        return this;
    }*/

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcBrowseItem");
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
        writeBuffer.popContext("PlcBrowseItem");
    }

}
