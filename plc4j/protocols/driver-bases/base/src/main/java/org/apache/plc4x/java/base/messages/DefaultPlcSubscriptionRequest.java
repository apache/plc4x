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
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.UncheckedPlcInvalidFieldException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.messages.items.FieldItem;

import java.time.Duration;
import java.util.*;
import java.util.function.BiFunction;

public class DefaultPlcSubscriptionRequest implements InternalPlcSubscriptionRequest {

    @Override
    public int getNumberOfFields() {
        return 0;
    }

    @Override
    public LinkedHashSet<String> getFieldNames() {
        return null;
    }

    @Override
    public PlcField getField(String name) {
        return null;
    }

    @Override
    public LinkedList<PlcField> getFields() {
        return null;
    }

    @Override
    public PlcSubscriptionType getPlcSubscriptionType() {
        return null;
    }

    public static class Builder implements PlcSubscriptionRequest.Builder {

        private final PlcFieldHandler fieldHandler;
        private final Map<String, BuilderItem<Object>> fields;

        public Builder(PlcFieldHandler fieldHandler) {
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
        public PlcSubscriptionRequest build() throws PlcInvalidFieldException {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> parsedFields = new LinkedHashMap<>();
            try {
                fields.forEach((name, builderItem) -> {
                    // Compile the query string.
                    PlcField parsedField = null;
                    try {
                        parsedField = fieldHandler.createField(builderItem.fieldQuery);
                    } catch (PlcInvalidFieldException e) {
                        throw new UncheckedPlcInvalidFieldException(e);
                    }
                    // Encode the payload.
                    // TODO: Depending on the field type, handle the FieldItem creation differently.
                    FieldItem fieldItem = builderItem.encoder.apply(parsedField, null);
                    parsedFields.put(name, new ImmutablePair<>(parsedField, fieldItem));
                });
            } catch (UncheckedPlcInvalidFieldException e) {
                throw e.getWrappedException();
            }
            return new DefaultPlcSubscriptionRequest();
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
