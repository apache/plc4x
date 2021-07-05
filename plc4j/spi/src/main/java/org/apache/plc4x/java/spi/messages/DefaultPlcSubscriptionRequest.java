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
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.w3c.dom.Element;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "className")
public class DefaultPlcSubscriptionRequest implements PlcSubscriptionRequest, Serializable {

    private final PlcSubscriber subscriber;

    private LinkedHashMap<String, PlcSubscriptionField> fields;

    @JsonCreator(mode = JsonCreator.Mode.PROPERTIES)
    public DefaultPlcSubscriptionRequest(@JsonProperty("subscriber") PlcSubscriber subscriber,
                                         @JsonProperty("fields") LinkedHashMap<String, PlcSubscriptionField> fields) {
        this.subscriber = subscriber;
        this.fields = fields;
    }

    @Override
    @JsonIgnore
    public CompletableFuture<PlcSubscriptionResponse> execute() {
        return subscriber.subscribe(this);
    }

    @Override
    @JsonIgnore
    public int getNumberOfFields() {
        return fields.size();
    }

    @Override
    @JsonIgnore
    public LinkedHashSet<String> getFieldNames() {
        return new LinkedHashSet<>(fields.keySet());
    }

    @Override
    @JsonIgnore
    public PlcSubscriptionField getField(String name) {
        return fields.get(name);
    }

    @Override
    @JsonIgnore
    public List<PlcSubscriptionField> getFields() {
        return new ArrayList<>(fields.values());
    }

    @JsonIgnore
    public List<Pair<String, PlcSubscriptionField>> getNamedFields() {
        return fields.entrySet()
            .stream()
            .map(stringPlcFieldEntry -> Pair.of(stringPlcFieldEntry.getKey(), stringPlcFieldEntry.getValue()))
            .collect(Collectors.toCollection(LinkedList::new));
    }

    public PlcSubscriber getSubscriber() {
        return subscriber;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext("PlcSubscriptionRequest");

        writeBuffer.pushContext("fields");
        for (Map.Entry<String, PlcSubscriptionField> fieldEntry : fields.entrySet()) {
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

        writeBuffer.popContext("PlcSubscriptionRequest");
    }

    public static class Builder implements PlcSubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private final PlcFieldHandler fieldHandler;
        private final Map<String, BuilderItem> fields;

        public Builder(PlcSubscriber subscriber, PlcFieldHandler fieldHandler) {
            this.subscriber = subscriber;
            this.fieldHandler = fieldHandler;
            fields = new TreeMap<>();
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicField(String name, String fieldQuery, Duration pollingInterval) {
            fields.put(name, new BuilderItem(fieldQuery, PlcSubscriptionType.CYCLIC, pollingInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateField(String name, String fieldQuery) {
            fields.put(name, new BuilderItem(fieldQuery, PlcSubscriptionType.CHANGE_OF_STATE));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventField(String name, String fieldQuery) {
            if (fields.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate field definition '" + name + "'");
            }
            fields.put(name, new BuilderItem(fieldQuery, PlcSubscriptionType.EVENT));
            return this;
        }

        @Override
        public PlcSubscriptionRequest build() {
            LinkedHashMap<String, PlcSubscriptionField> parsedFields = new LinkedHashMap<>();

            fields.forEach((name, builderItem) -> {
                PlcField parsedField = fieldHandler.createField(builderItem.fieldQuery);
                parsedFields.put(name, new DefaultPlcSubscriptionField(builderItem.plcSubscriptionType, parsedField, builderItem.duration));
            });
            return new DefaultPlcSubscriptionRequest(subscriber, parsedFields);
        }

        private static class BuilderItem {
            private final String fieldQuery;
            private final PlcSubscriptionType plcSubscriptionType;
            private final Duration duration;

            private BuilderItem(String fieldQuery, PlcSubscriptionType plcSubscriptionType) {
                this(fieldQuery, plcSubscriptionType, null);
            }

            private BuilderItem(String fieldQuery, PlcSubscriptionType plcSubscriptionType, Duration duration) {
                this.fieldQuery = fieldQuery;
                this.plcSubscriptionType = plcSubscriptionType;
                this.duration = duration;
            }

        }

    }

}
