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
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;

import java.nio.charset.StandardCharsets;
import java.util.Map;

public class DefaultListPlcBrowseItem extends DefaultPlcBrowseItem {

    public DefaultListPlcBrowseItem(PlcTag tag,
                                    String name,
                                    boolean readable,
                                    boolean writable,
                                    boolean subscribable,
                                    Map<String, PlcBrowseItem> children,
                                    Map<String, PlcValue> options) {
        super(tag, name, readable, writable, subscribable, children, options);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());
        writeBuffer.writeString("address",
            getTag().getAddressString().getBytes(StandardCharsets.UTF_8).length * 8,
            getTag().getAddressString(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
        writeBuffer.writeString("name",
            getName().getBytes(StandardCharsets.UTF_8).length * 8,
            getName(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
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
                writeBuffer.writeString("name",
                    optionEntry.getKey().getBytes(StandardCharsets.UTF_8).length * 8,
                    optionEntry.getKey(), WithOption.WithEncoding(StandardCharsets.UTF_8.name()));
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
