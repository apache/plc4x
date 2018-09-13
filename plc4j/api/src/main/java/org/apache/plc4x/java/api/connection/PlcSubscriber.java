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
package org.apache.plc4x.java.api.connection;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;

import java.util.Arrays;
import java.util.Collection;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.function.Consumer;

/**
 * Interface implemented by all PlcConnections that are able to receive notifications from remote resources.
 */
public interface PlcSubscriber {

    /**
     * Subscribes to fields on the PLC.
     *
     * @param subscriptionRequest subscription request containing at least one subscription request item.
     * @return subscription response containing a subscription response item for each subscription request item.
     */
    CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest);

    /**
     * Subscribes to fields on the PLC.
     *
     * @param plcSubscriptionRequestBuilderConsumer consumer which can be used to build requests.
     * @return subscription response containing a subscription response item for each subscription request item.
     */
    default CompletableFuture<PlcSubscriptionResponse> subscribe(Consumer<PlcSubscriptionRequest.Builder> plcSubscriptionRequestBuilderConsumer) {
        PlcSubscriptionRequest.Builder builder = subscriptionRequestBuilder();
        plcSubscriptionRequestBuilderConsumer.accept(builder);
        return subscribe(builder.build());
    }

    /**
     * Unsubscribes from fields on the PLC. For unsubscribing the unsubscription request uses the subscription
     * handle returned as part of the subscription response item.
     *
     * @param unsubscriptionRequest unsubscription request containing at least one unsubscription request item.
     * @return unsubscription response containing a unsubscription response item for each unsubscription request item.
     */
    CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest);

    /**
     * Unsubscribes from fields on the PLC. For unsubscribing the unsubscription request uses the subscription
     * handle returned as part of the subscription response item.
     *
     * @param plcSubscriptionRequestBuilderConsumer consumer which can be used to build requests.
     * @return unsubscription response containing a unsubscription response item for each unsubscription request item.
     */
    default CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(Consumer<PlcUnsubscriptionRequest.Builder> plcSubscriptionRequestBuilderConsumer) {
        PlcUnsubscriptionRequest.Builder builder = unsubscriptionRequestBuilder();
        plcSubscriptionRequestBuilderConsumer.accept(builder);
        return unsubscribe(builder.build());
    }

    /**
     * Convenience method to subscribe a {@link Consumer} to all fields of the subscription.
     *
     * @param subscriptionRequest subscription request
     * @param consumer            consumer for all {@link PlcSubscriptionEvent}s
     * @return TODO: document me
     * @throws ExecutionException   something went wrong.
     * @throws InterruptedException something went wrong.
     */
    default PlcConsumerRegistration register(PlcSubscriptionRequest subscriptionRequest, Consumer<PlcSubscriptionEvent> consumer) throws ExecutionException, InterruptedException {
        PlcSubscriptionResponse plcSubscriptionResponse = subscribe(subscriptionRequest).get();
        // TODO: we need to return the plcSubscriptionResponse here too as we need this to unsubscribe...
        return register(consumer, plcSubscriptionResponse.getSubscriptionHandles().toArray(new PlcSubscriptionHandle[0]));
    }

    /**
     * Convenience method to subscribe a {@link Consumer} to all fields of the subscription.
     *
     * @param subscriptionRequestBuilderConsumer consumer for building subscription request.
     * @param consumer                           consumer for all {@link PlcSubscriptionEvent}s
     * @return TODO: document me
     * @throws ExecutionException   something went wrong.
     * @throws InterruptedException something went wrong.
     */
    default PlcConsumerRegistration register(Consumer<PlcSubscriptionRequest.Builder> subscriptionRequestBuilderConsumer, Consumer<PlcSubscriptionEvent> consumer) throws ExecutionException, InterruptedException {
        PlcSubscriptionRequest.Builder builder = subscriptionRequestBuilder();
        subscriptionRequestBuilderConsumer.accept(builder);
        return register(builder.build(), consumer);
    }

    /**
     * @param consumer
     * @param handles
     * @return TODO: document me
     */
    PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles);

    default PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, PlcSubscriptionHandle... handles) {
        return register(consumer, Arrays.asList(handles));
    }

    /**
     * // TODO: document me.
     */
    void unregister(PlcConsumerRegistration registration);

    PlcSubscriptionRequest.Builder subscriptionRequestBuilder();

    PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder();

}
