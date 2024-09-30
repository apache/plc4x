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
package org.apache.plc4x.java.can.generic.protocol;

import java.util.Map.Entry;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcSubscriptionType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.can.adapter.Plc4xCANProtocolBase;
import org.apache.plc4x.java.can.generic.tag.GenericCANTag;
import org.apache.plc4x.java.can.generic.tag.GenericCANTagHandler;
import org.apache.plc4x.java.can.generic.transport.GenericFrame;
import org.apache.plc4x.java.genericcan.readwrite.DataItem;
import org.apache.plc4x.java.spi.ConversationContext;
import org.apache.plc4x.java.spi.connection.PlcTagHandler;
import org.apache.plc4x.java.spi.context.DriverContext;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.java.spi.messages.*;
import org.apache.plc4x.java.spi.messages.utils.DefaultPlcResponseItem;
import org.apache.plc4x.java.spi.messages.utils.PlcResponseItem;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;
import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.function.Consumer;

public class GenericCANProtocolLogic extends Plc4xCANProtocolBase<GenericFrame> implements PlcSubscriber {

    private final Logger logger = LoggerFactory.getLogger(GenericCANProtocolLogic.class);

    private RequestTransactionManager tm;
    private final Map<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> consumers = new ConcurrentHashMap<>();

    @Override
    public void setDriverContext(DriverContext driverContext) {
        super.setDriverContext(driverContext);
        this.tm = new RequestTransactionManager(1);
    }

    @Override
    public void close(ConversationContext<GenericFrame> context) {
        tm.shutdown();
    }

    @Override
    public void setConversationContext(ConversationContext<GenericFrame> conversationContext) {
        super.setConversationContext(conversationContext);
    }

    @Override
    public PlcTagHandler getTagHandler() {
        return new GenericCANTagHandler();
    }

    @Override
    public void onConnect(ConversationContext<GenericFrame> context) {
        context.fireConnected();
    }

    @Override
    public void onDisconnect(ConversationContext<GenericFrame> context) {
        context.fireDisconnected();
    }

    @Override
    public void decode(ConversationContext<GenericFrame> context, GenericFrame msg) throws Exception {
        for (Map.Entry<DefaultPlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> entry : consumers.entrySet()) {
            DefaultPlcConsumerRegistration registration = entry.getKey();
            Consumer<PlcSubscriptionEvent> consumer = entry.getValue();
            for (PlcSubscriptionHandle handle : registration.getSubscriptionHandles()) {
                GenericCANSubscriptionHandle subscription = (GenericCANSubscriptionHandle) handle;
                Map<String, PlcResponseItem<PlcValue>> tags = new LinkedHashMap<>();
                ReadBuffer buffer = new ReadBufferByteBased(msg.getData(), ByteOrder.LITTLE_ENDIAN);
                buffer.pullContext("readTags");
                if (subscription.matches(msg.getNodeId())) {
                    byte[] data = msg.getData();
                    ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, ByteOrder.LITTLE_ENDIAN);
                    for (Entry<String, GenericCANTag> tag : subscription.getTags().entrySet()) {
                        try {
                            PlcValue value = read(readBuffer, tag.getValue(), data.length);
                            if (value == null) {
                                tags.put(tag.getKey(), new DefaultPlcResponseItem<>(PlcResponseCode.INTERNAL_ERROR, null));
                            } else {
                                tags.put(tag.getKey(), new DefaultPlcResponseItem<>(PlcResponseCode.OK, value));
                            }
                        } catch (ParseException e) {
                            tags.put(tag.getKey(), new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_DATA, null));
                        }
                    }
                    consumer.accept(new DefaultPlcSubscriptionEvent(Instant.now(), tags));
                }
                buffer.closeContext("readTags");
            }
        }
    }

    private PlcValue read(ReadBuffer buffer, GenericCANTag tag, int length) throws ParseException {
        try {
            buffer.pullContext("read-" + tag);
            return DataItem.staticParse(buffer, tag.getDataType(), length);
        } finally {
            buffer.closeContext("read-" + tag);
        }
    }

    private void write(WriteBuffer buffer, GenericCANTag tag, PlcValue value) throws SerializationException {
        try {
            buffer.pushContext("write-" + tag);
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(DataItem.getLengthInBytes(value, tag.getDataType(), value.getLength()), ByteOrder.LITTLE_ENDIAN);
            DataItem.staticSerialize(writeBuffer, value, tag.getDataType(), value.getLength(), ByteOrder.LITTLE_ENDIAN);
            buffer.writeByteArray(writeBuffer.getBytes());
        } finally {
            buffer.popContext("write-" + tag);
        }
    }

    @Override
    public CompletableFuture<PlcWriteResponse> write(PlcWriteRequest writeRequest) {
        final RequestTransactionManager.RequestTransaction transaction = tm.startRequest();

        CompletableFuture<PlcWriteResponse> response = new CompletableFuture<>();
        transaction.submit(() -> {
            Map<Integer, WriteBufferByteBased> messages = new LinkedHashMap<>();
            Map<Integer, Map<String, PlcResponseCode>> responses = new HashMap<>();

            for (String tagName : writeRequest.getTagNames()) {
                PlcTag plcTag = writeRequest.getTag(tagName);
                if (!(plcTag instanceof GenericCANTag)) {
                    responses.computeIfAbsent(-1, (node) -> new HashMap<>())
                        .put(tagName, PlcResponseCode.UNSUPPORTED);
                    continue;
                }
                GenericCANTag canTag = (GenericCANTag) plcTag;
                WriteBuffer buffer = messages.computeIfAbsent(canTag.getNodeId(), (node) -> new WriteBufferByteBased(8, ByteOrder.LITTLE_ENDIAN));

                Map<String, PlcResponseCode> statusMap = responses.computeIfAbsent(canTag.getNodeId(), (node) -> new HashMap<>());

                PlcValue value = writeRequest.getPlcValue(tagName);
                try {
                    write(buffer, canTag, value);
                    statusMap.put(tagName, PlcResponseCode.OK);
                } catch (SerializationException e) {
                    statusMap.put(tagName, PlcResponseCode.INVALID_DATA);
                }
            }

            Map<String, PlcResponseCode> codes = new HashMap<>();
            for (Map.Entry<Integer, WriteBufferByteBased> message : messages.entrySet()) {
                boolean discarded = false;
                for (Map.Entry<String, PlcResponseCode> entry : responses.get(message.getKey()).entrySet()) {
                    codes.put(entry.getKey(), entry.getValue());
                    if (!discarded && entry.getValue() != PlcResponseCode.OK) {
                        logger.info("Discarding writing of frame with tag {}. Node {} will not be communicated.", entry.getKey(), message.getKey());
                        discarded = true;
                    }
                }
                if (!discarded) {
                    byte[] data = message.getValue().getBytes();
                    logger.debug("Writing message with id {} and {} bytes of data", message.getKey(), data.length);
                    conversationContext.sendToWire(new GenericFrame(message.getKey(), data));
                }
            }

            response.complete(new DefaultPlcWriteResponse(writeRequest, codes));
            transaction.endRequest();
        });

        return response;
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest request) {
        DefaultPlcSubscriptionRequest rq = (DefaultPlcSubscriptionRequest) request;

        Map<String, PlcResponseItem<PlcSubscriptionHandle>> answers = new LinkedHashMap<>();
        DefaultPlcSubscriptionResponse response = new DefaultPlcSubscriptionResponse(rq, answers);

        Map<Integer, GenericCANSubscriptionHandle> handles = new LinkedHashMap<>();
        for (String key : rq.getTagNames()) {
            DefaultPlcSubscriptionTag subscription = (DefaultPlcSubscriptionTag) rq.getTag(key);
            if (subscription.getPlcSubscriptionType() != PlcSubscriptionType.EVENT) {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.UNSUPPORTED, null));
            } else if (subscription.getTag() instanceof GenericCANTag) {
                GenericCANTag canTag = (GenericCANTag) subscription.getTag();
                GenericCANSubscriptionHandle subscriptionHandle = handles.computeIfAbsent(canTag.getNodeId(),
                    node -> new GenericCANSubscriptionHandle(this, node)
                );
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.OK, subscriptionHandle));
                subscriptionHandle.add(key, canTag);
            } else {
                answers.put(key, new DefaultPlcResponseItem<>(PlcResponseCode.INVALID_ADDRESS, null));
            }
        }

        return CompletableFuture.completedFuture(response);
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        final DefaultPlcConsumerRegistration consumerRegistration = new DefaultPlcConsumerRegistration(this, consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        consumers.put(consumerRegistration, consumer);
        return consumerRegistration;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        consumers.remove(registration);
    }
}
