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
import org.apache.plc4x.java.api.messages.PlcBrowseItemArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Map;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultListPlcBrowseItem extends DefaultPlcBrowseItem {

    private final List<PlcBrowseItemArrayInfo> arrayInfo;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultListPlcBrowseItem(@JsonProperty("tag") PlcTag tag,
                                    @JsonProperty("name") String name,
                                    @JsonProperty("dataType") PlcValueType dataType,
                                    @JsonProperty("arrayInfo") List<PlcBrowseItemArrayInfo> arrayInfo,
                                    @JsonProperty("readable") boolean readable,
                                    @JsonProperty("writable") boolean writable,
                                    @JsonProperty("subscribable") boolean subscribable,
                                    @JsonProperty("children") Map<String, PlcBrowseItem> children,
                                    @JsonProperty("options") Map<String, PlcValue> options) {
        super(tag, name, dataType, readable, writable, subscribable, children, options);
        this.arrayInfo = arrayInfo;
    }

    @Override
    public List<PlcBrowseItemArrayInfo> getArrayInfo() {
        return arrayInfo;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());
        writeBuffer.writeString("address", getTag().getAddressString().getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), getTag().getAddressString());
        writeBuffer.writeString("name", getName().getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), getName());
        // TODO: Find out how to serialize an enum.
        //writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);
        if(getChildren() != null && !getChildren().isEmpty()) {
            writeBuffer.pushContext("children");
            for (PlcBrowseItem child : getChildren().values()) {
                writeBuffer.pushContext("child");
                ((DefaultListPlcBrowseItem) child).serialize(writeBuffer);
                writeBuffer.popContext("child");
            }
            writeBuffer.popContext("children");
        }
        if(getOptions() != null && !getOptions().isEmpty()) {
            writeBuffer.pushContext("options");
            for (Map.Entry<String, PlcValue> optionEntry : getOptions().entrySet()) {
                writeBuffer.pushContext("option");
                writeBuffer.writeString("name", optionEntry.getKey().getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), optionEntry.getKey());
                // TODO: Find out how to serialize a PlcValue
                //writeBuffer.writeString("value", optionEntry.getValue().getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), optionEntry.getValue());
                ((DefaultListPlcBrowseItem) optionEntry).serialize(writeBuffer);
                writeBuffer.popContext("option");
            }
            writeBuffer.popContext("options");
        }
        writeBuffer.popContext(getClass().getSimpleName());
    }

}
