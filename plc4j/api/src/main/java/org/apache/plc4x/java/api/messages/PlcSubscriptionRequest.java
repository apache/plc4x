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
package org.apache.plc4x.java.api.messages;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public interface PlcSubscriptionRequest extends PlcFieldRequest {

    @Override
    CompletableFuture<? extends PlcSubscriptionResponse> execute();

    interface Builder extends PlcRequestBuilder {

        @Override
        PlcSubscriptionRequest build();

        /**
         * Adds a new field to the to be constructed request which should be polled cyclically.
         *
         * @param name            alias of the field.
         * @param fieldQuery      field query string for accessing the field.
         * @param pollingInterval interval, in which the field should be polled.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addCyclicField(String name, String fieldQuery, Duration pollingInterval);

        /**
         * Adds a new field to the to be constructed request which should be updated as soon as
         * a value changes in the PLC.
         *
         * @param name       alias of the field.
         * @param fieldQuery field query string for accessing the field.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addChangeOfStateField(String name, String fieldQuery);

        /**
         * Adds a new subscription to the to be constructed request which should be updated
         * as soon as an event occurs.
         * <p>
         * REMARK: We will have to see if this signature is correct as soon as we start using this type of subscription.
         *
         * @param name       alias of the field.
         * @param fieldQuery field query string for accessing the field.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addEventField(String name, String fieldQuery);
    }

}
