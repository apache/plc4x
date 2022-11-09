/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.api.messages;

import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcTag;

import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

public interface PlcSubscriptionRequest extends PlcSubscriptionTagRequest {

    @Override
    CompletableFuture<? extends PlcSubscriptionResponse> execute();

    interface Builder extends PlcRequestBuilder {

        @Override
        PlcSubscriptionRequest build();

        /**
         * Adds a new tag to the to be constructed request which should be polled cyclically.
         *
         * @param name            alias of the tag.
         * @param tagAddress      tag address string for accessing the tag.
         * @param pollingInterval interval, in which the tag should be polled.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval);

        /**
         * Adds a new tag to the to be constructed request which should be polled cyclically.
         *
         * @param name            alias of the tag.
         * @param tag           tag instance for accessing the tag.
         * @param pollingInterval interval, in which the tag should be polled.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval);

        /**
         * Adds a new tag to the to be constructed request which should be updated as soon as
         * a value changes in the PLC.
         *
         * @param name       alias of the tag.
         * @param tagAddress tag address string for accessing the tag.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress);

        /**
         * Adds a new tag to the to be constructed request which should be updated as soon as
         * a value changes in the PLC.
         *
         * @param name alias of the tag.
         * @param tag  tag instance for accessing the tag.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag);

        /**
         * Adds a new subscription to the to be constructed request which should be updated
         * as soon as an event occurs.
         * <p>
         * REMARK: We will have to see if this signature is correct as soon as we start using this type of subscription.
         *
         * @param name       alias of the tag.
         * @param tagAddress tag address string for accessing the tag.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress);

        /**
         * Adds a new subscription to the to be constructed request which should be updated
         * as soon as an event occurs.
         * <p>
         * REMARK: We will have to see if this signature is correct as soon as we start using this type of subscription.
         *
         * @param name alias of the tag.
         * @param tag  tag instance for accessing the tag.
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag);

        /**
         * Convenience method which attaches the {@link Consumer<PlcSubscriptionEvent>} directly to the handles once the
         * requests succeeds.
         * Note: opposed to register on the {@link org.apache.plc4x.java.api.model.PlcSubscriptionHandle} directly you
         * won't retrieve a {@link PlcConsumerRegistration} which is useful to cancel registrations.
         *
         * @param name                  alias of the tag.
         * @param preRegisteredConsumer {@link Consumer<PlcSubscriptionEvent>} to be attached
         * @return builder.
         */
        PlcSubscriptionRequest.Builder addPreRegisteredConsumer(String name, Consumer<PlcSubscriptionEvent> preRegisteredConsumer);

    }

}
