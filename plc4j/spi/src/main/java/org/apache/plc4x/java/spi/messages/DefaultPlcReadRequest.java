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

import com.fasterxml.jackson.annotation.*;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcFieldRequest;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.Supplier;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcReadRequest implements PlcReadRequest, PlcFieldRequest, Serializable {

    private final PlcReader reader;
    // This is intentionally a linked hash map in order to keep the order of how elements were added.
    private LinkedHashMap<String, PlcField> fields;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcReadRequest(@JsonProperty("reader") PlcReader reader,
                                 @JsonProperty("fields") LinkedHashMap<String, PlcField> fields) {
        this.reader = reader;
        this.fields = fields;
    }

    @Override
    @JsonIgnore
    public CompletableFuture<PlcReadResponse> execute() {
        return reader.read(this);
    }

    @Override
    @JsonIgnore
    public int getNumberOfFields() {
        return fields.size();
    }

    @Override
    @JsonIgnore
    public LinkedHashSet<String> getFieldNames() {
        // TODO: Check if this already is a LinkedHashSet.
        return new LinkedHashSet<>(fields.keySet());
    }

    @Override
    @JsonIgnore
    public PlcField getField(String name) {
        return fields.get(name);
    }

    @Override
    @JsonIgnore
    public List<PlcField> getFields() {
        return new LinkedList<>(fields.values());
    }

    @JsonIgnore
    public List<Pair<String, PlcField>> getNamedFields() {
        return fields.entrySet()
            .stream()
            .map(stringPlcFieldEntry -> Pair.of(stringPlcFieldEntry.getKey(), stringPlcFieldEntry.getValue()))
            .collect(Collectors.toCollection(LinkedList::new));
    }

    @JsonIgnore
    public PlcReader getReader() {
        return reader;
    }

    @JsonAnySetter
    public void add(String key, PlcField value) {
        fields.put(key, value);
    }

    @JsonAnyGetter
    public Map<String, PlcField> getMap() {
        return fields;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcReadRequest");

        writeBuffer.pushContext("fields");
        for (Map.Entry<String, PlcField> fieldEntry : fields.entrySet()) {
            String fieldName = fieldEntry.getKey();
            writeBuffer.pushContext(fieldName);
            PlcField field = fieldEntry.getValue();
            if(!(field instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Field doesn't implement XmlSerializable");
            }
            ((Serializable) field).serialize(writeBuffer);
            writeBuffer.popContext(fieldName);
        }
        writeBuffer.popContext("fields");

        writeBuffer.popContext("PlcReadRequest");
    }

    public static class Builder implements PlcReadRequest.Builder {

        private final PlcReader reader;
        private final PlcFieldHandler fieldHandler;
        private final Map<String, Supplier<PlcField>> fields;

        public Builder(PlcReader reader, PlcFieldHandler fieldHandler) {
            this.reader = reader;
            this.fieldHandler = fieldHandler;
            fields = new LinkedHashMap<>();
        }

        @Override
        public PlcReadRequest.Builder addItem(String name, String fieldQuery) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, () -> fieldHandler.parseField(fieldQuery));
            return this;
        }

        @Override
        public PlcReadRequest.Builder addItem(String name, PlcField fieldQuery) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, () -> fieldQuery);
            return this;
        }

        @Override
        public PlcReadRequest build() {
            LinkedHashMap<String, PlcField> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, fieldQuery) -> {
                parsedFields.put(name, fieldQuery.get());
            });
            return new DefaultPlcReadRequest(reader, parsedFields);
        }

    }

}
