/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/
package org.apache.plc4x.java.base.messages;

import org.apache.commons.lang3.tuple.ImmutablePair;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.*;
import java.util.function.BiFunction;
import java.util.stream.Collectors;

public class DefaultPlcWriteRequest implements InternalPlcWriteRequest {

    private final LinkedHashMap<String, Pair<PlcField, FieldItem>> fields;

    private DefaultPlcWriteRequest(LinkedHashMap<String, Pair<PlcField, FieldItem>> fields) {
        this.fields = fields;
    }

    @Override
    public int getNumberOfFields() {
        return fields.size();
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        // TODO: Check if this already is a LinkedHashSet.
        return new LinkedHashSet<>(fields.keySet());
    }

    @Override
    public PlcField getField(String name) {
        return fields.get(name).getKey();
    }

    @Override
    public LinkedList<PlcField> getFields() {
        return fields.values().stream().map(Pair::getKey).collect(Collectors.toCollection(LinkedList::new));
    }

    public FieldItem getFieldItem(String name) {
        return fields.get(name).getValue();
    }

    @Override
    public LinkedList<FieldItem> getFieldItems() {
        return fields.values().stream().map(Pair::getValue).collect(Collectors.toCollection(LinkedList::new));
    }

    @Override
    public int getNumberOfValues(String name) {
        return fields.get(name).getValue().getNumValues();
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcFieldHandler fieldHandler;
        private final Map<String, BuilderItem<Object>> fields;
        private final Map<Class<?>, BiFunction<PlcField, Object[], FieldItem>> handlerMap;

        public Builder(PlcFieldHandler fieldHandler) {
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
            BiFunction<PlcField, Object[], FieldItem> plcFieldFieldItemBiFunction = handlerMap.get(checkedClazz);
            if (plcFieldFieldItemBiFunction == null) {
                throw new IllegalArgumentException("no field handler for " + checkedClazz + " found");
            }
            return addItem(name, fieldQuery, values, plcFieldFieldItemBiFunction);
        }

        @Override
        public PlcWriteRequest build() {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, builderItem) -> {
                // Compile the query string.
                PlcField parsedField = fieldHandler.createField(builderItem.fieldQuery);
                // Encode the payload.
                // TODO: Depending on the field type, handle the FieldItem creation differently.
                FieldItem fieldItem = builderItem.encoder.apply(parsedField, builderItem.values);
                parsedFields.put(name, new ImmutablePair<>(parsedField, fieldItem));
            });
            return new DefaultPlcWriteRequest(parsedFields);
        }

        private Builder addItem(String name, String fieldQuery, Object[] values, BiFunction<PlcField, Object[], FieldItem> encoder) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, new BuilderItem<>(fieldQuery, values, encoder));
            return this;
        }

        private static class BuilderItem<T> {
            private final String fieldQuery;
            private final T[] values;
            private final BiFunction<PlcField, T[], FieldItem> encoder;

            private BuilderItem(String fieldQuery, T[] values, BiFunction<PlcField, T[], FieldItem> encoder) {
                this.fieldQuery = fieldQuery;
                this.values = values;
                this.encoder = encoder;
            }
        }

    }

}
