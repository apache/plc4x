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
        private final Map<String, BuilderItem<Object>> fields;
        private final Map<Class<?>, BiFunction<PlcField, Object[], PlcValue>> handlerMap;

        public Builder(PlcWriter writer, PlcFieldHandler fieldHandler) {
            this.writer = writer;
            this.fieldHandler = fieldHandler;
            fields = new TreeMap<>();
            handlerMap = new HashMap<>();
            handlerMap.put(Boolean.class, fieldHandler::encodeBoolean);
            handlerMap.put(Byte.class, fieldHandler::encodeByte);
            handlerMap.put(Short.class, fieldHandler::encodeShort);
            handlerMap.put(Integer.class, fieldHandler::encodeInteger);
            handlerMap.put(BigInteger.class, fieldHandler::encodeBigInteger);
            handlerMap.put(Long.class, fieldHandler::encodeLong);
            handlerMap.put(Float.class, fieldHandler::encodeFloat);
            handlerMap.put(Double.class, fieldHandler::encodeDouble);
            handlerMap.put(BigDecimal.class, fieldHandler::encodeBigDecimal);
            handlerMap.put(String.class, fieldHandler::encodeString);
            handlerMap.put(LocalTime.class, fieldHandler::encodeTime);
            handlerMap.put(LocalDate.class, fieldHandler::encodeDate);
            handlerMap.put(LocalDateTime.class, fieldHandler::encodeDateTime);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Boolean... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeBoolean);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Byte... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeByte);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Short... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeShort);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Integer... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeInteger);
        }

        @Override
        public PlcWriteRequest.Builder addItem(String name, String fieldQuery, BigInteger... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeBigInteger);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Long... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeLong);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Float... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeFloat);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Double... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeDouble);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, BigDecimal... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeBigDecimal);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, String... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeString);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, LocalTime... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeTime);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, LocalDate... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeDate);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, LocalDateTime... values) {
            return addItem(name, fieldQuery, values, fieldHandler::encodeDateTime);
        }

        @Override
        public <T> Builder addItem(String name, String fieldQuery, T... values) {
            Objects.requireNonNull(values);
            Class<?> checkedClazz = null;
            for (T value : values) {
                if (checkedClazz == null) {
                    checkedClazz = value.getClass();
                }
                if (value.getClass() != checkedClazz) {
                    throw new IllegalArgumentException("Invalid class found " + value.getClass() + ". should all be " + checkedClazz);
                }
            }
            BiFunction<PlcField, Object[], PlcValue> plcFieldPlcValueBiFunction = handlerMap.get(checkedClazz);
            if (plcFieldPlcValueBiFunction == null) {
                throw new IllegalArgumentException("no field handler for " + checkedClazz + " found");
            }
            return addItem(name, fieldQuery, values, plcFieldPlcValueBiFunction);
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, FieldValueItem> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, builderItem) -> {
                // Compile the query string.
                PlcField parsedField = fieldHandler.createField(builderItem.fieldQuery);
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
            fields.put(name, new BuilderItem<>(fieldQuery, values, encoder));
            return this;
        }

        private static class BuilderItem<T> {
            private final String fieldQuery;
            private final T[] values;
            private final BiFunction<PlcField, T[], PlcValue> encoder;

            private BuilderItem(String fieldQuery, T[] values, BiFunction<PlcField, T[], PlcValue> encoder) {
                this.fieldQuery = fieldQuery;
                this.values = values;
                this.encoder = encoder;
            }
        }

    }

}
