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
package org.apache.plc4x.java.s7.readwrite.utils;

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

import org.apache.plc4x.java.s7.readwrite.TimeBase;
import org.apache.plc4x.java.s7.readwrite.tag.S7SubscriptionTag;
import org.apache.plc4x.java.s7.readwrite.tag.S7Tag;
import org.apache.plc4x.java.s7.readwrite.types.S7SubscriptionType;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionRequest;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;

//public class S7PlcSubscriptionRequest implements PlcSubscriptionRequest, Serializable {
public class S7PlcSubscriptionRequest extends DefaultPlcSubscriptionRequest {

    private static final String CONST_DUPLICATE_TAG = "Duplicate tag definition";
    private static final String CONST_INVALID_TYPE = "Tag is not of type S7SubscriptionTag";
    private static final String CONST_TIME_CANNOT_BE_ZERO = "Subscription time cannot be zero.";


//    private final PlcSubscriber subscriber;
//
//    private final LinkedHashMap<String, PlcTagItem<PlcSubscriptionTag>> tags;
//
//    private final Consumer<PlcSubscriptionEvent> consumer;
//    private final Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers;

    public S7PlcSubscriptionRequest(PlcSubscriber subscriber,
                                    LinkedHashMap<String, PlcTagItem<PlcSubscriptionTag>> tags,
                                    Consumer<PlcSubscriptionEvent> consumer,
                                    Map<String, Consumer<PlcSubscriptionEvent>> tagConsumers) {
        super(subscriber, tags, consumer, tagConsumers);
//        this.subscriber = subscriber;
//        this.tags = tags;
//        this.consumer = consumer;
//        this.tagConsumers = tagConsumers;
    }

//    @Override
//    public CompletableFuture<PlcSubscriptionResponse> execute() {
//        return subscriber.subscribe(this);
//    }
//
//    @Override
//    public int getNumberOfTags() {
//        return tags.size();
//    }
//
//    @Override
//    public LinkedHashSet<String> getTagNames() {
//        return new LinkedHashSet<>(tags.keySet());
//    }
//
//    @Override
//    public PlcSubscriptionTag getTag(String name) {
//        return tags.get(name).getTag();
//    }
//
//    @Override
//    public PlcResponseCode getTagResponseCode(String tagName) {
//        return tags.get(tagName).getResponseCode();
//    }
//
//    @Override
//    public Consumer<PlcSubscriptionEvent> getConsumer() {
//        return consumer;
//    }
//
//    @Override
//    public Consumer<PlcSubscriptionEvent> getTagConsumer(String tagName) {
//        return tagConsumers.get(tagName);
//    }
//
//    public Map<String, Consumer<PlcSubscriptionEvent>> getTagConsumers() {
//        return tagConsumers;
//    }
//
//    @Override
//    public List<PlcSubscriptionTag> getTags() {
//        return tags.values().stream().map(PlcTagItem::getTag).collect(Collectors.toCollection(LinkedList::new));
//    }
//
//    @Override
//    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
//        writeBuffer.pushContext("PlcSubscriptionRequest");
//
//        writeBuffer.pushContext("tags");
//        for (Map.Entry<String, PlcTagItem<PlcSubscriptionTag>> tagEntry : tags.entrySet()) {
//            String tagName = tagEntry.getKey();
//            writeBuffer.pushContext(tagName);
//            PlcTagItem<PlcSubscriptionTag> tagItem = tagEntry.getValue();
//            if (!(tagItem instanceof Serializable)) {
//                throw new PlcRuntimeException("Error serializing. Tag doesn't implement XmlSerializable");
//            }
//            ((Serializable) tagItem).serialize(writeBuffer);
//            writeBuffer.popContext(tagName);
//        }
//        writeBuffer.popContext("tags");
//
//        writeBuffer.popContext("PlcSubscriptionRequest");
//    }

    public static class Builder implements PlcSubscriptionRequest.Builder {

        private final PlcSubscriber subscriber;
        private final PlcTagHandler tagHandler;
        private final Map<String, BuilderItem> tags;
        private Consumer<PlcSubscriptionEvent> consumer;

        public Builder(PlcSubscriber subscriber, PlcTagHandler tagHandler) {
            this.subscriber = subscriber;
            this.tagHandler = tagHandler;
            this.tags = new TreeMap<>();
        }

        @Override
        public PlcSubscriptionRequest.Builder setConsumer(Consumer<PlcSubscriptionEvent> consumer) {
            this.consumer = Objects.requireNonNull(consumer);
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
            return addCyclicTagAddress(name, tagAddress, pollingInterval, null);
        }

        /*
        * This method receives a String that describes an S7Tag and the 
        * interval required for its sampling.
        * The value of the "pollingInterval" parameter is adapted to the 
        * cyclical subscription requirements of an S7-300/S7-400, 
        * for which multiples of the time base given by TimeBase 
        * must be handled. To say:
        *
        * . B01SEC -> 100, 200, 300, 400, 500, 600, 700, 800, 900 msec
        * . B1SEC  ->   1,   2,   3,   4,   5,   6,   7,   8,   9 sec
        * . B10SEC ->  10,  20,  30,  40,  50,  60,  70,  80,  90 sec
        *
        * As you can see there are no intermediate values, for example 513 msec,
        * it will actually be 500 msec, or its nearest rounding.
        * 
        * @param name Name of the subscription Tag.
        * @param tagAddress String representing an S7Tag
        * @param pollingInterval Required sampling rate based on the "TimeBase"  
        * @return PlcSubscriptionRequest.Builder S7SubscriptonTag type constructor        
        * 
        */
        @Override
        public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException(CONST_DUPLICATE_TAG + " '" + name + "'");
            }
            TimeBase tb = getTimeBase(pollingInterval);
            short multiplier = getMultiplier(tb, pollingInterval);
            S7Tag[] s7tags = new S7Tag[]{S7Tag.of(tagAddress)};
            S7SubscriptionTag tag = new S7SubscriptionTag(S7SubscriptionType.CYCLIC_SUBSCRIPTION, s7tags, tb, multiplier);
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CYCLIC, pollingInterval, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
            return addCyclicTag(name, tag, pollingInterval, null);
        }

        /*
        * This method receives an S7Tag built by the user, he is responsible 
        * for the construction of the object, so no additional verification 
        * is included.
        *
        * @param name Name of the subscription Tag.
        * @param tag    Tag of S7SubscriptionTag type.
        * @param pollingInterval Required sampling rate based on the "TimeBase"
        * @return PlcSubscriptionRequest.Builder S7SubscriptonTag type constructor
        */
        @Override
        public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException(CONST_DUPLICATE_TAG + " '" + name + "'");
            }
            if (!(tag instanceof S7SubscriptionTag)){
                throw new PlcRuntimeException(CONST_INVALID_TYPE);
            }                
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CYCLIC, pollingInterval, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
            return addChangeOfStateTagAddress(tagAddress, name);
        }

        /*
        *
        */
        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException(CONST_DUPLICATE_TAG + " '" + name + "'");
            }
            S7Tag[] s7tags = new S7Tag[]{S7Tag.of(tagAddress)};   
            S7SubscriptionTag tag = new S7SubscriptionTag(S7SubscriptionType.CYCLIC_SUBSCRIPTION, s7tags, TimeBase.B01SEC, (short) 1);            
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CHANGE_OF_STATE, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
            return addChangeOfStateTag(name, tag, null);
        }

        /*
        *
        */        
        @Override
        public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException(CONST_DUPLICATE_TAG + " '" + name + "'");
            }
            if (!(tag instanceof S7SubscriptionTag)){
                throw new PlcRuntimeException(CONST_INVALID_TYPE);
            }              
            tags.put(name, new BuilderItem(() -> tag, PlcSubscriptionType.CHANGE_OF_STATE, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
            return addEventTagAddress(name, tagAddress, null);
        }

        /*
        * This method is responsible for the subscription to Events associated 
        * with the PLC as well as the preliminary version of cyclical 
        * subscription of values.
        *
        * The type of function performed by the tag is given by the definition 
        * of the "tagAddress", for example:
        *
        * "ACK:16#12345678"
        *
        * Represents an acknowledgment of an alarm whose ID is 16#12345678.
        * The following functions are defined:
        *
        * . MODE
        * . SYS
        * . USR
        * . ALM
        * . ACK
        * . QUERY
        * . CYC
        * . CANCEL
        * 
        * Go to the driver manual for a complete description.
        * 
        * @param name Name of the subscription Tag.
        * @param tag    Tag of S7SubscriptionTag type.        
        * @return PlcSubscriptionRequest.Builder S7SubscriptonTag type constructor
        */
        @Override
        public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException(CONST_DUPLICATE_TAG + " '" + name + "'");
            }
            PlcTag tag = tagHandler.parseTag(tagAddress);
            if (!(tag instanceof S7SubscriptionTag)){
                throw new PlcRuntimeException(CONST_INVALID_TYPE);
            }              
            tags.put(name, new BuilderItem(() -> tagHandler.parseTag(tagAddress), PlcSubscriptionType.EVENT, consumer));
            return this;
        }

        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
            return addEventTag(name, tag, null);
        }

        /*
        * This method receives an S7Tag built by the user, he is responsible 
        * for the construction of the object, so no additional verification 
        * is included.
        *
        * @param name Name of the subscription Tag.
        * @param tag    Tag of S7SubscriptionTag type.
        * @return PlcSubscriptionRequest.Builder S7SubscriptonTag type constructor        
        */
        @Override
        public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag, Consumer<PlcSubscriptionEvent> consumer) {
            if (tags.containsKey(name)) {
                throw new PlcRuntimeException(CONST_DUPLICATE_TAG + " '" + name + "'");
            }
            if (!(tag instanceof S7SubscriptionTag)){
                throw new PlcRuntimeException(CONST_INVALID_TYPE);
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

            return new S7PlcSubscriptionRequest(subscriber, parsedTags, consumer, tagConsumers);
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
        
        private TimeBase getTimeBase(Duration duration)  {
            if (duration.equals(Duration.ZERO)) {
                throw new PlcRuntimeException(CONST_TIME_CANNOT_BE_ZERO);
            }
            long millis = duration.toMillis();
            if (millis < 1000) {
                return TimeBase.B01SEC;
            }
            if (millis < 10000) {
                return TimeBase.B1SEC;                
            }
            if (millis < 100000) {
                return TimeBase.B10SEC;  
            }
            
            throw new PlcRuntimeException("The maximum subscription time is 90 sec.");             
        }
        
        //TODO: Check multiplier is 1-99 in BCD??
        private short getMultiplier(TimeBase timeBase, Duration duration)  {
            short multiplier = 1;
            if (duration.equals(Duration.ZERO)) {
                throw new PlcRuntimeException(CONST_TIME_CANNOT_BE_ZERO);
            }
            long millis = duration.toMillis();
            switch(timeBase) {
                case B01SEC:
                    if (millis > 100) {
                        multiplier = (short) (millis / 100);
                    }
                    break;
                case B1SEC:
                    multiplier = (short) (millis / 1000);
                    break;
                case B10SEC:
                    multiplier = (short) (millis / 10000);
                    break;
            }           
            return multiplier;            
        }        

    }

    @Override
    public String toString() {
        return "DefaultPlcSubscriptionRequest{" +
//            "subscriber=" + super.getSubscriber().subscriber +
//            ", tags=" + tags +
            '}';
    }
}
