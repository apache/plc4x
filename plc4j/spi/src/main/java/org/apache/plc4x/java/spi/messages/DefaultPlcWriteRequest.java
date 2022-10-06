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
import org.apache.commons.lang3.tuple.Triple;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.messages.utils.FieldValueItem;

import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcWriteRequest implements PlcWriteRequest, Serializable {

    private final PlcWriter writer;
    private final LinkedHashMap<String, FieldValueItem> fields;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcWriteRequest(@JsonProperty("writer") PlcWriter writer,
                                  @JsonProperty("fields") LinkedHashMap<String, FieldValueItem> fields) {
        this.writer = writer;
        this.fields = fields;
    }

    @Override
    @JsonIgnore
    public CompletableFuture<PlcWriteResponse> execute() {
        return writer.write(this);
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
        return fields.get(name).getField();
    }

    @Override
    @JsonIgnore
    public List<PlcField> getFields() {
        return fields.values().stream().map(FieldValueItem::getField).collect(Collectors.toCollection(LinkedList::new));
    }

    @JsonIgnore
    public PlcValue getPlcValue(String name) {
        return fields.get(name).getValue();
    }

    @JsonIgnore
    public List<PlcValue> getPlcValues() {
        return fields.values().stream().map(FieldValueItem::getValue).collect(Collectors.toCollection(LinkedList::new));
    }

    @JsonIgnore
    public List<Pair<String, PlcField>> getNamedFields() {
        return fields.entrySet()
            .stream()
            .map(stringPairEntry ->
                Pair.of(
                    stringPairEntry.getKey(),
                    stringPairEntry.getValue().getField()
                )
            ).collect(Collectors.toCollection(LinkedList::new));
    }


    public PlcWriter getWriter() {
        return writer;
    }

    @JsonIgnore
    public List<Triple<String, PlcField, PlcValue>> getNamedFieldTriples() {
        return fields.entrySet()
            .stream()
            .map(stringPairEntry ->
                Triple.of(
                    stringPairEntry.getKey(),
                    stringPairEntry.getValue().getField(),
                    stringPairEntry.getValue().getValue()
                )
            ).collect(Collectors.toCollection(LinkedList::new));
    }

    @Override
    @JsonIgnore
    public int getNumberOfValues(String name) {
        final PlcValue value = fields.get(name).getValue();
        if (value instanceof PlcList) {
            PlcList list = (PlcList) value;
            return list.getLength();
        }
        return 1;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcWriteRequest");

        writeBuffer.pushContext("fields");
        for (Map.Entry<String, FieldValueItem> fieldEntry : fields.entrySet()) {
            FieldValueItem fieldValueItem = fieldEntry.getValue();
            String fieldName = fieldEntry.getKey();
            writeBuffer.pushContext(fieldName);
            PlcField field = fieldValueItem.getField();
            if (!(field instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Field doesn't implement XmlSerializable");
            }
            ((Serializable) field).serialize(writeBuffer);
            final PlcValue value = fieldValueItem.getValue();
            if (value instanceof PlcList) {
                PlcList list = (PlcList) value;
                for (PlcValue plcValue : list.getList()) {
                    String plcValueString = plcValue.getString();
                    writeBuffer.writeString("value", plcValueString.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), plcValueString);
                }
            } else {
                String plcValueString = value.getString();
                writeBuffer.writeString("value", plcValueString.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), plcValueString);
            }
            writeBuffer.popContext(fieldName);
        }
        writeBuffer.popContext("fields");

        writeBuffer.popContext("PlcWriteRequest");
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcWriter writer;
        private final PlcFieldHandler fieldHandler;
        private final PlcValueHandler valueHandler;
        private final Map<String, Pair<String, Object[]>> fields;

        public Builder(PlcWriter writer, PlcFieldHandler fieldHandler, PlcValueHandler valueHandler) {
            this.writer = writer;
            this.fieldHandler = fieldHandler;
            this.valueHandler = valueHandler;
            fields = new TreeMap<>();
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Object... values) {
            fields.put(name, Pair.of(fieldQuery, values));
            return this;
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, FieldValueItem> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, fieldValues) -> {
                // Compile the query string.
                String fieldQuery = fieldValues.getLeft();
                PlcField field = fieldHandler.createField(fieldQuery);
                Object[] value = fieldValues.getRight();
                PlcValue plcValue = valueHandler.newPlcValue(field, value);
                parsedFields.put(name, new FieldValueItem(field, plcValue));
            });
            return new DefaultPlcWriteRequest(writer, parsedFields);
        }
    }

}
