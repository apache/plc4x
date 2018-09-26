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
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.PromiseCombiner;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.commons.lang3.tuple.Triple;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;
import java.util.concurrent.atomic.AtomicInteger;

// TODO: write test
public class SingleItemToSingleRequestProtocol extends ChannelDuplexHandler {

    public static final Logger LOGGER = LoggerFactory.getLogger(SingleItemToSingleRequestProtocol.class);

    private PendingWriteQueue queue;

    private ConcurrentMap<Integer, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>>> sentButUnacknowledgedRequestItems;

    private ConcurrentMap<PlcRequestContainer<?, ?>, Set<Integer>> containerCorrelationIdMap;

    private ConcurrentMap<PlcRequestContainer<?, ?>, List<InternalPlcResponse<?>>> responsesToBeDevliered;

    private AtomicInteger correlationId;

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

    private void tryFinish(int correlationId, InternalPlcResponse msg) {
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>> plcRequestContainer = sentButUnacknowledgedRequestItems.remove(correlationId);
        if (plcRequestContainer == null) {
            throw new PlcRuntimeException("Unrelated package received " + msg);
        }
        List<InternalPlcResponse<?>> correlatedResponseItems = responsesToBeDevliered.computeIfAbsent(plcRequestContainer, ignore -> new LinkedList<>());
        correlatedResponseItems.add(msg);
        Set<Integer> integers = containerCorrelationIdMap.get(plcRequestContainer);
        integers.remove(correlationId);
        if (integers.isEmpty()) {
            InternalPlcResponse<?> plcResponse;
            if (plcRequestContainer.getRequest() instanceof InternalPlcReadRequest) {
                InternalPlcReadRequest internalPlcReadRequest = (InternalPlcReadRequest) plcRequestContainer.getRequest();
                HashMap<String, Pair<PlcResponseCode, FieldItem>> fields = new HashMap<>();

                correlatedResponseItems.stream()
                    .map(InternalPlcReadResponse.class::cast)
                    .map(InternalPlcReadResponse::getValues)
                    .forEach(stringPairMap -> stringPairMap.forEach(fields::put));

                plcResponse = new DefaultPlcReadResponse(internalPlcReadRequest, fields);
            } else if (plcRequestContainer.getRequest() instanceof InternalPlcWriteRequest) {
                InternalPlcWriteRequest internalPlcWriteRequest = (InternalPlcWriteRequest) plcRequestContainer.getRequest();
                HashMap<String, PlcResponseCode> values = new HashMap<>();

                correlatedResponseItems.stream()
                    .map(InternalPlcWriteResponse.class::cast)
                    .map(InternalPlcWriteResponse::getValues)
                    .forEach(stringPairMap -> stringPairMap.forEach(values::put));

                plcResponse = new DefaultPlcWriteResponse(internalPlcWriteRequest, values);
            } else {
                throw new PlcRuntimeException("Unknown type detected " + plcRequestContainer.getRequest());
            }
            plcRequestContainer.getResponseFuture().complete(plcResponse);
            responsesToBeDevliered.remove(plcRequestContainer);
        }
    }

    private void errored(int correlationId, Throwable throwable) {

    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) throws Exception {
        if (msg instanceof PlcRequestContainer) {
            @SuppressWarnings("unchecked")
            PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>> in = (PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>>) msg;
            Set<Integer> tdpus = containerCorrelationIdMap.computeIfAbsent(in, plcRequestContainer -> new HashSet<>());

            // Create a promise that has to be called multiple times.
            PromiseCombiner promiseCombiner = new PromiseCombiner();
            InternalPlcRequest request = in.getRequest();
            if (request instanceof InternalPlcFieldRequest) {
                InternalPlcFieldRequest internalPlcFieldRequest = (InternalPlcFieldRequest) request;

                if (internalPlcFieldRequest instanceof InternalPlcReadRequest) {
                    InternalPlcReadRequest internalPlcReadRequest = (InternalPlcReadRequest) internalPlcFieldRequest;
                    // TODO: repackage
                    internalPlcReadRequest.getNamedFields().forEach(field -> {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                        int tdpu = correlationId.getAndIncrement();
                        CompletableFuture<InternalPlcResponse> correlatedCompletableFuture = new CompletableFuture<>()
                            .thenApply(InternalPlcResponse.class::cast)
                            .whenComplete((internalPlcResponse, throwable) -> {
                                if (throwable != null) {
                                    errored(tdpu, throwable);
                                } else {
                                    tryFinish(tdpu, internalPlcResponse);
                                }
                            });
                        queue.add(new PlcRequestContainer<>(CorrelatedPlcReadRequest.of(field, tdpu), correlatedCompletableFuture), subPromise);
                        if (!tdpus.add(tdpu)) {
                            throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                        }
                        promiseCombiner.add((Future) subPromise);
                    });
                }
                if (internalPlcFieldRequest instanceof InternalPlcWriteRequest) {
                    InternalPlcWriteRequest internalPlcWriteRequest = (InternalPlcWriteRequest) internalPlcFieldRequest;
                    // TODO: repackage
                    internalPlcWriteRequest.getNamedFieldTriples().forEach(fieldItemTriple -> {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                        int tdpu = correlationId.getAndIncrement();
                        CompletableFuture<InternalPlcResponse> correlatedCompletableFuture = new CompletableFuture<>()
                            .thenApply(InternalPlcResponse.class::cast)
                            .whenComplete((internalPlcResponse, throwable) -> {
                                if (throwable != null) {
                                    errored(tdpu, throwable);
                                } else {
                                    tryFinish(tdpu, internalPlcResponse);
                                }
                            });
                        queue.add(new PlcRequestContainer<>(CorrelatedPlcWriteRequest.of(fieldItemTriple, tdpu), correlatedCompletableFuture), subPromise);
                        if (!tdpus.add(tdpu)) {
                            throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                        }
                        promiseCombiner.add((Future) subPromise);
                    });
                }
            } else {
                ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());
                queue.add(msg, subPromise);
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
            PlcRequestContainer currentItem = (PlcRequestContainer) queue.current();
            InternalPlcRequest request = currentItem.getRequest();

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

            if (request instanceof CorrelatedPlcRequest) {
                CorrelatedPlcRequest correlatedPlcRequest = (CorrelatedPlcRequest) request;

                // Add it to the list of sentButUnacknowledgedRequestItems.
                sentButUnacknowledgedRequestItems.put(correlatedPlcRequest.getTdpu(), currentItem);

                LOGGER.debug("Item Message with id {} sent", correlatedPlcRequest.getTdpu());
            }
        }
        ctx.flush();
    }

    interface CorrelatedPlcRequest extends InternalPlcRequest {

        int getTdpu();
    }

    private static class CorrelatedPlcReadRequest extends DefaultPlcReadRequest implements CorrelatedPlcRequest {

        private final int tdpu;

        public CorrelatedPlcReadRequest(LinkedHashMap<String, PlcField> fields, int tdpu) {
            super(fields);
            this.tdpu = tdpu;
        }

        public static CorrelatedPlcReadRequest of(Pair<String, PlcField> stringPlcFieldPair, int tdpu) {
            LinkedHashMap<String, PlcField> fields = new LinkedHashMap<>();
            fields.put(stringPlcFieldPair.getKey(), stringPlcFieldPair.getValue());
            return new CorrelatedPlcReadRequest(fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }

    private static class CorrelatedPlcWriteRequest extends DefaultPlcWriteRequest implements CorrelatedPlcRequest {

        private final int tdpu;

        public CorrelatedPlcWriteRequest(LinkedHashMap<String, Pair<PlcField, FieldItem>> fields, int tdpu) {
            super(fields);
            this.tdpu = tdpu;
        }

        public static CorrelatedPlcWriteRequest of(Triple<String, PlcField, FieldItem> fieldItemTriple, int tdpu) {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> fields = new LinkedHashMap<>();
            fields.put(fieldItemTriple.getLeft(), Pair.of(fieldItemTriple.getMiddle(), fieldItemTriple.getRight()));
            return new CorrelatedPlcWriteRequest(fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }
}
