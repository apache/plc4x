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
import org.apache.plc4x.java.spi.PlcSubscriber;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.BiFunction;

// TODO: request broken needs finishing.
public class DefaultPlcSubscriptionRequest implements InternalPlcSubscriptionRequest, InternalPlcFieldRequest {

    private final PlcSubscriber subscriber;

    public DefaultPlcSubscriptionRequest(PlcSubscriber subscriber) {
        this.subscriber = subscriber;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> execute() {
        return subscriber.subscribe(this);
    }

    @Override
    public int getNumberOfFields() {
        throw new IllegalStateException("not available");
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        throw new IllegalStateException("not available");
    }

    @Override
    public PlcField getField(String name) {
        throw new IllegalStateException("not available");
    }

    @Override
    public LinkedList<PlcField> getFields() {
        throw new IllegalStateException("not available");
    }

    @Override
    public PlcSubscriptionType getPlcSubscriptionType() {
        throw new IllegalStateException("not available");
    }

    @Override
    public LinkedList<Pair<String, PlcField>> getNamedFields() {
        throw new IllegalStateException("not available");
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
            return null;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateField(String name, String fieldQuery) {
            return null;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventField(String name, String fieldQuery) {
            return null;
        }

        @Override
        public PlcSubscriptionRequest build() {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> parsedFields = new LinkedHashMap<>();
            fields.forEach((name, builderItem) -> {
                // Compile the query string.
                PlcField parsedField = fieldHandler.createField(builderItem.fieldQuery);
                // Encode the payload.
                // TODO: Depending on the field type, handle the FieldItem creation differently.
                FieldItem fieldItem = builderItem.encoder.apply(parsedField, null);
                parsedFields.put(name, new ImmutablePair<>(parsedField, fieldItem));
            });
            return new DefaultPlcSubscriptionRequest(subscriber);
        }

        private static class BuilderItem<T> {
            private final String fieldQuery;
            private final BiFunction<PlcField, T[], FieldItem> encoder;

            private BuilderItem(String fieldQuery, BiFunction<PlcField, T[], FieldItem> encoder) {
                this.fieldQuery = fieldQuery;
                this.encoder = encoder;
            }
        }

    }

}
