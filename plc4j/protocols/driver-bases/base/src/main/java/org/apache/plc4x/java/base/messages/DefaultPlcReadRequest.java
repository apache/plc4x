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

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

public class DefaultPlcReadRequest implements InternalPlcReadRequest, InternalPlcFieldRequest {

    private final PlcReader reader;
    private LinkedHashMap<String, PlcField> fields;

    protected DefaultPlcReadRequest(PlcReader reader, LinkedHashMap<String, PlcField> fields) {
        this.reader = reader;
        this.fields = fields;
    }

    @Override
    public CompletableFuture<PlcReadResponse> execute() {
        return reader.read(this);
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
        return fields.get(name);
    }

    @Override
    public LinkedList<PlcField> getFields() {
        return new LinkedList<>(fields.values());
    }

    @Override
    public LinkedList<Pair<String, PlcField>> getNamedFields() {
        return fields.entrySet()
            .stream()
            .map(stringPlcFieldEntry -> Pair.of(stringPlcFieldEntry.getKey(), stringPlcFieldEntry.getValue()))
            .collect(Collectors.toCollection(LinkedList::new));
    }

    protected PlcReader getReader() {
        return reader;
    }

    public static class Builder implements PlcReadRequest.Builder {

        private final PlcReader reader;
        private final PlcFieldHandler fieldHandler;
        private final Map<String, String> fields;

        public Builder(PlcReader reader, PlcFieldHandler fieldHandler) {
            this.reader = reader;
            this.fieldHandler = fieldHandler;
            fields = new TreeMap<>();
        }

        @Override
        public PlcReadRequest.Builder addItem(String name, String fieldQuery) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, fieldQuery);
            return this;
        }

        @Override
        public PlcReadRequest build() {
            LinkedHashMap<String, PlcField> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, fieldQuery) -> {
                PlcField parsedField = fieldHandler.createField(fieldQuery);
                parsedFields.put(name, parsedField);
            });
            return new DefaultPlcReadRequest(reader, parsedFields);
        }

    }

}
