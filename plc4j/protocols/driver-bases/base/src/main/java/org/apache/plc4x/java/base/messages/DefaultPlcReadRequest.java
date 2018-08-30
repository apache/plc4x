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
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;

import java.util.Collection;
import java.util.Map;
import java.util.TreeMap;

public class DefaultPlcReadRequest implements PlcReadRequest {

    private Map<String, PlcField> fields;

    private DefaultPlcReadRequest(Map<String, PlcField> fields) {
        this.fields = fields;
    }

    @Override
    public Collection<String> getFieldNames() {
        return fields.keySet();
    }

    @Override
    public PlcField getField(String name) {
        return fields.get(name);
    }

    public static class DefaultPlcReadRequestBuilder implements PlcReadRequest.Builder {

        private final PlcFieldHandler fieldHandler;
        private final Map<String, String> fields;

        public DefaultPlcReadRequestBuilder(PlcFieldHandler fieldHandler) {
            this.fieldHandler = fieldHandler;
            fields = new TreeMap<>();
        }

        @Override
        public Builder addItem(String name, String fieldQuery) {
            if(fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, fieldQuery);
            return this;
        }

        @Override
        public PlcReadRequest build() {
            Map<String, PlcField> parsedFields = new TreeMap<>();
            fields.forEach((name, fieldQuery) -> {
                PlcField parsedField = fieldHandler.createField(fieldQuery);
                parsedFields.put(name, parsedField);
            });
            return new DefaultPlcReadRequest(parsedFields);
        }

    }

}
