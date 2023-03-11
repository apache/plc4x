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
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;
import java.util.function.Supplier;

public class DefaultPlcSubscriptionRequest implements PlcSubscriptionRequest, Serializable {

    private final PlcSubscriber subscriber;

    private final LinkedHashMap<String, PlcSubscriptionTag> tags;

    private final LinkedHashMap<String, List<Consumer<PlcSubscriptionEvent>>> preRegisteredConsumers;

    public DefaultPlcSubscriptionRequest(PlcSubscriber subscriber,
                                         LinkedHashMap<String, PlcSubscriptionTag> tags,
                                         LinkedHashMap<String, List<Consumer<PlcSubscriptionEvent>>> preRegisteredConsumers) {
        this.subscriber = subscriber;
        this.tags = tags;
        this.preRegisteredConsumers = preRegisteredConsumers;
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
    public PlcSubscriptionTag getTag(String name) {
        return tags.get(name);
    }

    @Override
    public List<PlcSubscriptionTag> getTags() {
        return new ArrayList<>(tags.values());
    }

    @Override
    public Map<String, List<Consumer<PlcSubscriptionEvent>>> getPreRegisteredConsumers() {
        return new LinkedHashMap<>(preRegisteredConsumers);
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext("PlcSubscriptionRequest");

        writeBuffer.pushContext("tags");
        for (Map.Entry<String, PlcSubscriptionTag> tagEntry : tags.entrySet()) {
            String tagName = tagEntry.getKey();
            writeBuffer.pushContext(tagName);
            PlcTag tag = tagEntry.getValue();
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
        private final LinkedHashMap<String, List<Consumer<PlcSubscriptionEvent>>> preRegisteredConsumers;

        public Builder(PlcSubscriber subscriber, PlcTagHandler tagHandler) {
            this.subscriber = subscriber;
            this.tagHandler = tagHandler;
            this.tags = new TreeMap<>();
            this.preRegisteredConsumers = new LinkedHashMap<>();
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.CYCLIC, pollingInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CYCLIC, pollingInterval));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.CHANGE_OF_STATE));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CHANGE_OF_STATE));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.EVENT));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException("Duplicate tag definition '" + name + "'");
            }
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.EVENT));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addPreRegisteredConsumer(String name, Consumer<PlcSubscriptionEvent> consumer) {
            preRegisteredConsumers.putIfAbsent(name, new LinkedList<>());
            preRegisteredConsumers.get(name).add(consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest build() {
            LinkedHashMap<String, PlcSubscriptionTag> parsedTags = new LinkedHashMap<>();

            tags.forEach((name, builderItem) -> {
                PlcTag parsedTag = builderItem.tag.get();
                parsedTags.put(name, new DefaultPlcSubscriptionTag(builderItem.plcSubscriptionType, parsedTag, builderItem.duration));
            });
            preRegisteredConsumers.forEach((tagName, ignored) -> {
                if (!tags.containsKey(tagName)) {
                    throw new RuntimeException("tagName " + tagName + "for preRegisteredConsumer not found");
                }
            });
            return new DefaultPlcSubscriptionRequest(subscriber, parsedTags, preRegisteredConsumers);
        }

        private static class BuilderItem {
            private final Supplier<PlcTag> tag;
            private final PlcSubscriptionType plcSubscriptionType;
            private final Duration duration;

            private BuilderItem(Supplier<PlcTag> tag, PlcSubscriptionType plcSubscriptionType) {
                this(tag, plcSubscriptionType, null);
            }

            private BuilderItem(Supplier<PlcTag> tag, PlcSubscriptionType plcSubscriptionType, Duration duration) {
                this.tag = tag;
                this.plcSubscriptionType = plcSubscriptionType;
                this.duration = duration;
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
