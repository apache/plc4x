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
package org.apache.plc4x.java.spi.drivers.messages;

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.spi.drivers.functions.PlcSubscriber;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

public class DefaultPlcSubscriptionRequest implements PlcSubscriptionRequest {

    private final PlcSubscriber subscriber;
    private final LinkedHashMap<String, PlcSubscriptionTag> tags;
    private final Consumer<PlcSubscriptionEvent> consumer;
    private final Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers;

    public DefaultPlcSubscriptionRequest(PlcSubscriber subscriber,
                                         LinkedHashMap<String, PlcSubscriptionTag> tags,
                                         Consumer<PlcSubscriptionEvent> consumer,
                                         Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers) {
        this.subscriber = subscriber;
        this.tags = tags;
        this.consumer = consumer;
        this.tagConsumers = tagConsumers;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> execute() {
        return subscriber.subscribe(this);
    }

    @Override
    public int getNumberOfTags() {
        return tags.size();
    }

    @Override
    public LinkedHashSet<String> getTagNames() {
        return new LinkedHashSet<>(tags.keySet());
    }

    @Override
    public PlcResponseCode getTagResponseCode(String tagName) {
        if (tags.containsKey(tagName)) {
            return PlcResponseCode.OK;
        }
        return PlcResponseCode.NOT_FOUND;
    }

    @Override
    public PlcSubscriptionTag getTag(String name) {
        return tags.get(name);
    }

    @Override
    public List<PlcSubscriptionTag> getTags() {
        return new ArrayList<>(tags.values());
    }

    @Override
    public Consumer<PlcSubscriptionEvent> getConsumer() {
        return consumer;
    }

    @Override
    public Consumer<PlcSubscriptionEvent> getTagConsumer(String name) {
        return tagConsumers.get(name);
    }

    public static class Builder implements PlcSubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private final PlcTagHandler tagHandler;
        private final LinkedHashMap<String, PlcSubscriptionTag> tags = new LinkedHashMap<>();
        private Consumer<PlcSubscriptionEvent> consumer;
        private final Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers = new LinkedHashMap<>();

        public Builder(PlcSubscriber subscriber, PlcTagHandler tagHandler) {
            this.subscriber = subscriber;
            this.tagHandler = tagHandler;
        }

        @Override
        public PlcSubscriptionRequest.Builder setConsumer(Consumer<PlcSubscriptionEvent> consumer) {
            this.consumer = consumer;
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CYCLIC, tag, pollingInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CYCLIC, tag, pollingInterval));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CYCLIC, tag, pollingInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CYCLIC, tag, pollingInterval));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, null));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, null));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Duration minInterval) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, minInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer, Duration minInterval) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, minInterval));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, null));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, null));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Duration minInterval) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, minInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer, Duration minInterval) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.CHANGE_OF_STATE, tag, minInterval));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.EVENT, tag, null));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            PlcTag tag = tagHandler.parseTag(tagAddress);
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.EVENT, tag, null));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.EVENT, tag, null));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            tags.put(name, new DefaultPlcSubscriptionTag(PlcSubscriptionType.EVENT, tag, null));
            tagConsumers.put(name, consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest build() {
            return new DefaultPlcSubscriptionRequest(subscriber, tags, consumer, tagConsumers);
        }
    }

}
