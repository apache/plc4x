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

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;

import java.util.Collection;
import java.util.concurrent.CompletableFuture;
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
     * Unsubscribes from fields on the PLC. For unsubscribing the unsubscription request uses the subscription
     * handle returned as part of the subscription response item.
     *
     * @param unsubscriptionRequest unsubscription request containing at least one unsubscription request item.
     * @return unsubscription response containing a unsubscription response item for each unsubscription request item.
     */
    CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest);

    /**
     * @param consumer
     * @param handles
     * @return TODO: document me
     */
    PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles);

    /**
     * // TODO: document me.
     */
    void unregister(PlcConsumerRegistration registration);

}
