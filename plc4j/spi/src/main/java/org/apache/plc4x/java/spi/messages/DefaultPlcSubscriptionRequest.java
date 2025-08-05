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
package org.apache.plc4x.java.spi.messages;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcTagItem;
import org.apache.plc4x.java.spi.messages.utils.PlcTagItem;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Supplier;
import java.util.stream.Collectors;

public class DefaultPlcSubscriptionRequest implements PlcSubscriptionRequest, Serializable {

    private final PlcSubscriber subscriber;

    private final LinkedHashMap<String, PlcTagItem<PlcSubscriptionTag>> tags;

    private final Consumer<PlcSubscriptionEvent> consumer;
    private final Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers;

    public DefaultPlcSubscriptionRequest(PlcSubscriber subscriber,
                                         LinkedHashMap<String, PlcTagItem<PlcSubscriptionTag>> tags,
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
    public PlcSubscriptionTag getTag(String tagName) {
        return tags.get(tagName).getTag();
    }

    @Override
    public PlcResponseCode getTagResponseCode(String tagName) {
        return tags.get(tagName).getResponseCode();
    }

    @Override
    public List<PlcSubscriptionTag> getTags() {
        return tags.values().stream().map(PlcTagItem::getTag).collect(Collectors.toCollection(LinkedList::new));
    }

    @Override
    public Consumer<PlcSubscriptionEvent> getConsumer() {
        return consumer;
    }

    @Override
    public Consumer<PlcSubscriptionEvent> getTagConsumer(String tagName) {
        return tagConsumers.get(tagName);
    }

    public Map<String, Consumer<PlcSubscriptionEvent>> getTagConsumers() {
        return tagConsumers;
    }

    public PlcSubscriber getSubscriber() {
        return subscriber;
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcSubscriptionRequest");

        writeBuffer.pushContext("tags");
        for (Map.Entry<String, PlcTagItem<PlcSubscriptionTag>> tagEntry : tags.entrySet()) {
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcTagItem<PlcSubscriptionTag> tag = tagEntry.getValue();
            if (!(tag instanceof Serializable)) {
                throw new RuntimeException("Error serializing. Tag doesn't implement XmlSerializable");
            }
            ((Serializable) tag).serialize(writeBuffer);
            writeBuffer.popContext(tagName);
        }
        writeBuffer.popContext("tags");

        writeBuffer.popContext("PlcSubscriptionRequest");
    }

    public static class Builder implements PlcSubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private final PlcTagHandler tagHandler;
        private final Map<String, BuilderItem> tags;
        private Consumer<PlcSubscriptionEvent> consumer;

        public Builder(PlcSubscriber subscriber, PlcTagHandler tagHandler) {
            this.subscriber = Objects.requireNonNull(subscriber);
            this.tagHandler = Objects.requireNonNull(tagHandler);
            this.tags = new TreeMap<>();
        }

        @Override
        public PlcSubscriptionRequest.Builder setConsumer(Consumer<PlcSubscriptionEvent> consumer) {
            this.consumer = Objects.requireNonNull(consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
            addCyclicTagAddress(name, tagAddress, pollingInterval, null);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.CYCLIC, pollingInterval, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
            addCyclicTag(name, tag, pollingInterval, null);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CYCLIC, pollingInterval, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
            addChangeOfStateTagAddress(name, tagAddress, null);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.CHANGE_OF_STATE, consumer));
            return null;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
            addChangeOfStateTag(name, tag, null);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CHANGE_OF_STATE, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
            addEventTagAddress(name, tagAddress, null);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.EVENT, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
            addEventTag(name, tag, null);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.EVENT, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest build() {
            LinkedHashMap<String, PlcTagItem<PlcSubscriptionTag>> parsedTags = new LinkedHashMap<>();

            Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers = new LinkedHashMap<>();
            tags.forEach((name, builderItem) -> {
                PlcTag parsedTag = builderItem.tag.get();
                parsedTags.put(name, new DefaultPlcTagItem<>(new DefaultPlcSubscriptionTag(builderItem.plcSubscriptionType, parsedTag, builderItem.duration)));
                if(builderItem.consumer != null) {
                    tagConsumers.put(name, builderItem.consumer);
                }
            });

            return new DefaultPlcSubscriptionRequest(subscriber, parsedTags, consumer, tagConsumers);
        }

        private static class BuilderItem {
            private final Supplier<PlcTag> tag;
            private final PlcSubscriptionType plcSubscriptionType;
            private final Duration duration;
            private final Consumer<PlcSubscriptionEvent> consumer;

            private BuilderItem(Supplier<PlcTag> tag, PlcSubscriptionType plcSubscriptionType, Consumer<PlcSubscriptionEvent> consumer) {
                this(tag, plcSubscriptionType, null, consumer);
            }

            private BuilderItem(Supplier<PlcTag> tag, PlcSubscriptionType plcSubscriptionType, Duration duration, Consumer<PlcSubscriptionEvent> consumer) {
                this.tag = tag;
                this.plcSubscriptionType = plcSubscriptionType;
                this.duration = duration;
                this.consumer = consumer;
            }
        }
    }

    @Override
    public String toString() {
        return "DefaultPlcSubscriptionRequest{" +
            "subscriber=" + subscriber +
            ", tags=" + tags +
            '}';
    }
}
