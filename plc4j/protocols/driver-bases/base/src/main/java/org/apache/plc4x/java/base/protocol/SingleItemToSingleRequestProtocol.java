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
package org.apache.plc4x.java.base.protocol;

import io.netty.channel.*;
import io.netty.handler.codec.MessageToMessageDecoder;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.PromiseCombiner;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.messages.items.RequestItem;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcReadResponse;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteRequest;
import org.apache.plc4x.java.api.messages.specific.TypeSafePlcWriteResponse;
import org.apache.plc4x.java.base.messages.item.CorrelatedRequestItem;
import org.apache.plc4x.java.base.messages.item.CorrelatedResponseItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

public class SingleItemToSingleRequestProtocol extends ChannelDuplexHandler {

    public static final Logger LOGGER = LoggerFactory.getLogger(SingleItemToSingleRequestProtocol.class);

    private PendingWriteQueue queue;

    private ConcurrentMap<Integer, CorrelatedRequestItem<?>> sentButUnacknowledgedRequestItems;

    private ConcurrentMap<PlcRequestContainer<?, ?>, Set<Integer>> containerCorrelationIdMap;

    private ConcurrentMap<PlcRequestContainer<?, ?>, List<CorrelatedResponseItem<?>>> responsesToBeDevliered;

    private AtomicInteger correlationId;

    private final MessageToMessageDecoder<CorrelatedResponseItem> decoder = new MessageToMessageDecoder<CorrelatedResponseItem>() {

        @Override
        protected void decode(ChannelHandlerContext ctx, CorrelatedResponseItem msg, List<Object> out) throws Exception {
            SingleItemToSingleRequestProtocol.this.decode(ctx, msg, out);
        }
    };

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) throws Exception {
        this.queue = new PendingWriteQueue(ctx);
        this.sentButUnacknowledgedRequestItems = new ConcurrentHashMap<>();
        this.containerCorrelationIdMap = new ConcurrentHashMap<>();
        this.correlationId = new AtomicInteger();
        super.channelRegistered(ctx);
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        this.queue.removeAndWriteAll();
        super.channelUnregistered(ctx);
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        // Send everything so we get a proper failure for those pending writes
        this.queue.removeAndWriteAll();
        super.channelInactive(ctx);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        decoder.channelRead(ctx, msg);
        super.read(ctx);
    }

    private void decode(ChannelHandlerContext ctx, CorrelatedResponseItem<?> msg, List<Object> out) throws PlcProtocolException {
        int correlationId = msg.getCorrelationId();
        CorrelatedRequestItem<?> correlatedRequestItem = sentButUnacknowledgedRequestItems.remove(correlationId);
        if (correlatedRequestItem == null) {
            throw new PlcProtocolException("Unrelated package received " + msg);
        }
        PlcRequestContainer<?, PlcResponse<?, ?, ?>> plcRequestContainer = correlatedRequestItem.getPlcRequestContainer();
        List<CorrelatedResponseItem<?>> correlatedResponseItems = responsesToBeDevliered.computeIfAbsent(plcRequestContainer, ignore -> new LinkedList<>());
        correlatedResponseItems.add(msg);
        Set<Integer> integers = containerCorrelationIdMap.get(plcRequestContainer);
        integers.remove(correlationId);
        if (integers.isEmpty()) {
            PlcResponse<?, ?, ?> plcResponse;
            if (plcRequestContainer.getRequest() instanceof TypeSafePlcReadRequest) {
                TypeSafePlcReadRequest typeSafePlcReadRequest = (TypeSafePlcReadRequest) plcRequestContainer.getRequest();
                plcResponse = new TypeSafePlcReadResponse((TypeSafePlcReadRequest<?>) plcRequestContainer.getRequest(), correlatedResponseItems.stream().map(correlatedResponseItem -> correlatedResponseItem.getResponseItem()).collect(Collectors.toList()));
            } else if (plcRequestContainer.getRequest() instanceof TypeSafePlcWriteRequest) {
                plcResponse = new TypeSafePlcWriteResponse((TypeSafePlcWriteRequest<?>) plcRequestContainer.getRequest(), correlatedResponseItems.stream().map(correlatedResponseItem -> correlatedResponseItem.getResponseItem()).collect(Collectors.toList()));
            } else if (plcRequestContainer.getRequest() instanceof PlcReadRequest) {
                plcResponse = new PlcReadResponse((PlcReadRequest) plcRequestContainer.getRequest(), (List) correlatedResponseItems.stream().map(correlatedResponseItem -> correlatedResponseItem.getResponseItem()).collect(Collectors.toList()));
            } else if (plcRequestContainer.getRequest() instanceof PlcWriteRequest) {
                plcResponse = new PlcWriteResponse((PlcWriteRequest) plcRequestContainer.getRequest(), (List) correlatedResponseItems.stream().map(correlatedResponseItem -> correlatedResponseItem.getResponseItem()).collect(Collectors.toList()));
            } else {
                throw new PlcProtocolException("Unknown type detected " + plcRequestContainer.getRequest());
            }
            plcRequestContainer.getResponseFuture().complete(plcResponse);
            responsesToBeDevliered.remove(plcRequestContainer);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) throws Exception {
        if (msg instanceof PlcRequestContainer) {
            PlcRequestContainer<?, PlcResponse<?, ?, ?>> in = (PlcRequestContainer<?, PlcResponse<?, ?, ?>>) msg;
            Set<Integer> tdpus = containerCorrelationIdMap.computeIfAbsent(in, plcRequestContainer -> new HashSet<>());

            // Create a promise that has to be called multiple times.
            PromiseCombiner promiseCombiner = new PromiseCombiner();
            PlcRequest<?> request = in.getRequest();
            for (RequestItem<?> item : request.getRequestItems()) {
                ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                int tdpu = correlationId.getAndIncrement();
                queue.add(new CorrelatedRequestItem<>(tdpu, item, in), subPromise);
                if (!tdpus.add(tdpu)) {
                    throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                }
                promiseCombiner.add((Future) subPromise);
            }
            promiseCombiner.finish(promise);

            // Start sending the queue content.
            trySendingMessages(ctx);
        } else {
            super.write(ctx, msg, promise);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Helpers
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    protected synchronized void trySendingMessages(ChannelHandlerContext ctx) {
        while (queue.size() > 0) {
            // Get the RequestItem that is up next in the queue.
            CorrelatedRequestItem<?> currentItem = (CorrelatedRequestItem) queue.current();

            if (currentItem == null) {
                break;
            }
            // Send the TPDU.
            try {
                ChannelFuture channelFuture = queue.removeAndWrite();
                ctx.flush();
                if (channelFuture == null) {
                    break;
                }
            } catch (Exception e) {
                LOGGER.error("Error sending more queues messages", e);
                ctx.fireExceptionCaught(e);
            }

            // Add it to the list of sentButUnacknowledgedRequestItems.
            sentButUnacknowledgedRequestItems.put(currentItem.getCorrelationId(), currentItem);

            LOGGER.debug("Item Message with id {} sent", currentItem.getCorrelationId());
        }
        ctx.flush();
    }
}
