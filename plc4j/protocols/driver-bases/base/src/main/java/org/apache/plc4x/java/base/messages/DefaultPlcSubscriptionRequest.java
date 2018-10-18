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
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;

public class DefaultPlcSubscriptionRequest implements InternalPlcSubscriptionRequest, InternalPlcFieldRequest {

    private final PlcSubscriber subscriber;

    private LinkedHashMap<String, SubscriptionPlcField> fields;

    public DefaultPlcSubscriptionRequest(PlcSubscriber subscriber, LinkedHashMap<String, SubscriptionPlcField> fields) {
        this.subscriber = subscriber;
        this.fields = fields;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> execute() {
        return subscriber.subscribe(this);
    }

    @Override
    public int getNumberOfFields() {
        return fields.size();
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        return new LinkedHashSet<>(fields.keySet());
    }

    @Override
    public PlcField getField(String name) {
        SubscriptionPlcField subscriptionPlcField = fields.get(name);
        if (subscriptionPlcField == null) {
            return null;
        }
        return subscriptionPlcField.getPlcField();
    }

    @Override
    public LinkedList<PlcField> getFields() {
        return fields.values().stream().map(SubscriptionPlcField::getPlcField).collect(Collectors.toCollection(LinkedList::new));
    }

    @Override
    public LinkedList<SubscriptionPlcField> getSubscriptionFields() {
        return new LinkedList<>(fields.values());
    }

    @Override
    public LinkedList<Pair<String, PlcField>> getNamedFields() {
        return fields.entrySet()
            .stream()
            .map(stringPlcFieldEntry -> Pair.of(stringPlcFieldEntry.getKey(), (PlcField) stringPlcFieldEntry.getValue()))
            .collect(Collectors.toCollection(LinkedList::new));
    }

    public static class Builder implements PlcSubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private final PlcFieldHandler fieldHandler;
        private final Map<String, BuilderItem<Object>> fields;

        public Builder(PlcSubscriber subscriber, PlcFieldHandler fieldHandler) {
            this.subscriber = subscriber;
            this.fieldHandler = fieldHandler;
            fields = new TreeMap<>();
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicField(String name, String fieldQuery, Duration pollingInterval) {
            fields.put(name, new BuilderItem<>(fieldQuery, PlcSubscriptionType.CYCLIC, pollingInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateField(String name, String fieldQuery) {
            fields.put(name, new BuilderItem<>(fieldQuery, PlcSubscriptionType.CHANGE_OF_STATE));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventField(String name, String fieldQuery) {
            fields.put(name, new BuilderItem<>(fieldQuery, PlcSubscriptionType.EVENT));
            return this;
        }

        @Override
        public PlcSubscriptionRequest build() {
            LinkedHashMap<String, SubscriptionPlcField> parsedFields = new LinkedHashMap<>();

            fields.forEach((name, builderItem) -> {
                PlcField parsedField = fieldHandler.createField(builderItem.fieldQuery);
                parsedFields.put(name, new SubscriptionPlcField(builderItem.plcSubscriptionType, parsedField, builderItem.duration));
            });
            return new DefaultPlcSubscriptionRequest(subscriber, parsedFields);
        }

        private static class BuilderItem<T> {
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
