/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.messages.PlcWriteResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcList;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.messages.utils.FieldValueItem;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.HashMap;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.TreeMap;
import java.util.concurrent.CompletableFuture;
import java.util.function.BiFunction;
import java.util.function.Supplier;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcWriteRequest implements InternalPlcWriteRequest, InternalPlcFieldRequest {

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

    @Override
    @JsonIgnore
    public List<PlcValue> getPlcValues() {
        return fields.values().stream().map(FieldValueItem::getValue).collect(Collectors.toCollection(LinkedList::new));
    }

    @Override
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

    @Override
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
        if(value instanceof PlcList) {
            PlcList list = (PlcList) value;
            return list.getLength();
        }
        return 1;
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcWriter writer;
        private final PlcFieldHandler fieldHandler;
        private final PlcValueHandler valueHandler;
        private final Map<String, BuilderItem<Object>> fields;

        public Builder(PlcWriter writer, PlcFieldHandler fieldHandler, PlcValueHandler valueHandler) {
            this.writer = writer;
            this.fieldHandler = fieldHandler;
            this.valueHandler = valueHandler;
            fields = new TreeMap<>();
        }

        @Override
        public <T> Builder addItem(String name, String fieldQuery, Object... values) {
            return addItem(name, fieldQuery, values, valueHandler.of(values));
        }

        @Override
        public <T> Builder addItem(String name, PlcField fieldQuery, Object... values) {
            return addItem(name, fieldQuery, values, valueHandler.of(values));
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, FieldValueItem> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, builderItem) -> {
                // Compile the query string.
                PlcField parsedField = builderItem.fieldQuery.get();
                // Encode the payload.
                // TODO: Depending on the field type, handle the PlcValue creation differently.
                PlcValue value = builderItem.encoder.apply(parsedField, builderItem.values);
                parsedFields.put(name, new FieldValueItem(parsedField, value));
            });
            return new DefaultPlcWriteRequest(writer, parsedFields);
        }

        private Builder addItem(String name, String fieldQuery, Object[] values, BiFunction<PlcField, Object[], PlcValue> encoder) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, new BuilderItem<>(() -> fieldHandler.createField(fieldQuery), values, encoder));
            return this;
        }

        private Builder addItem(String name, PlcField field, Object[] values, BiFunction<PlcField, Object[], PlcValue> encoder) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, new BuilderItem<>(() -> field, values, encoder));
            return this;
        }

        private static class BuilderItem<T> {
            private final Supplier<PlcField> fieldQuery;
            private final T[] values;
            private final BiFunction<PlcField, T[], PlcValue> encoder;

            private BuilderItem(Supplier<PlcField> fieldQuery, T[] values, BiFunction<PlcField, T[], PlcValue> encoder) {
                this.fieldQuery = fieldQuery;
                this.values = values;
                this.encoder = encoder;
            }
        }

    }

}
