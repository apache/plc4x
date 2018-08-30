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

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcClientDatatype;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Collection;
import java.util.Map;
import java.util.TreeMap;

public class DefaultPlcWriteRequest implements PlcWriteRequest {

    private final Map<String, Item> fields;

    private DefaultPlcWriteRequest(Map<String, Item> fields) {
        this.fields = fields;
    }

    @Override
    public Collection<String> getFieldNames() {
        return fields.keySet();
    }

    @Override
    public PlcField getField(String name) {
        return fields.get(name).field;
    }

    @Override
    public byte[][] getValues(String name) {
        return fields.get(name).values;
    }

    public static class Builder implements PlcWriteRequest.Builder {

        private final PlcFieldHandler fieldHandler;
        private final Map<String, BuilderItem> fields;

        public Builder(PlcFieldHandler fieldHandler) {
            this.fieldHandler = fieldHandler;
            fields = new TreeMap<>();
        }

        @Override
        public PlcWriteRequest.Builder addItem(String name, String fieldQuery, byte[]... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.RAW);
        }

        @Override
        public PlcWriteRequest.Builder addItem(String name, String fieldQuery, Object... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.OBJECT);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Boolean... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.BOOLEAN);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Byte... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.BYTE);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Short... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.SHORT);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Integer... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.INTEGER);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Long... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.LONG);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Float... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.FLOAT);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, Double... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.DOUBLE);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, String... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.STRING);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, LocalTime... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.TIME);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, LocalDate... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.DATE);
        }

        @Override
        public Builder addItem(String name, String fieldQuery, LocalDateTime... values) {
            return addItem(name, fieldQuery, values, PlcClientDatatype.DATE_TIME);
        }

        @Override
        public PlcWriteRequest build() {
            Map<String, DefaultPlcWriteRequest.Item> parsedFields = new TreeMap<>();
            fields.forEach((name, builderItem) -> {
                // Compile the query string.
                PlcField parsedField = fieldHandler.createField(builderItem.fieldQuery);
                // Encode the payload.
                byte[][] rawData = fieldHandler.encode(parsedField, builderItem.clientDatatype, builderItem.values);
                parsedFields.put(name, new DefaultPlcWriteRequest.Item(parsedField, rawData));
            });
            return new DefaultPlcWriteRequest(parsedFields);
        }

        private Builder addItem(String name, String fieldQuery, Object[] values, PlcClientDatatype clientDatatype) {
            if(fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, new BuilderItem(fieldQuery, clientDatatype, values));
            return this;
        }

        private static class BuilderItem {

            private final String fieldQuery;
            private final PlcClientDatatype clientDatatype;
            private final Object[] values;

            private BuilderItem(String fieldQuery, PlcClientDatatype clientDatatype, Object[] values) {
                this.fieldQuery = fieldQuery;
                this.clientDatatype = clientDatatype;
                this.values = values;
            }

        }

    }

    private static class Item {

        private final PlcField field;
        private final byte[][] values;

        private Item(PlcField field, byte[][] values) {
            this.field = field;
            this.values = values;
        }

    }

}
