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
import java.util.stream.Collectors;

/**
 * This layer can be used to split a {@link org.apache.plc4x.java.api.messages.PlcRequest} which addresses multiple {@link PlcField}s into multiple subsequent {@link org.apache.plc4x.java.api.messages.PlcRequest}s.
 */
public class SingleItemToSingleRequestProtocol extends ChannelDuplexHandler {

    public static final Logger LOGGER = LoggerFactory.getLogger(SingleItemToSingleRequestProtocol.class);

    private PendingWriteQueue queue;

    // Map to track send subcontainers
    private ConcurrentMap<Integer, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>>> sentButUnacknowledgedSubContainer;

    // Map to map tdpu to original parent container
    private ConcurrentMap<Integer, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>>> correlationToParentContainer;

    // Map to track tdpus per container
    private ConcurrentMap<PlcRequestContainer<?, ?>, Set<Integer>> containerCorrelationIdMap;

    // Map to track a list of responses per parent container
    private ConcurrentMap<PlcRequestContainer<?, ?>, List<InternalPlcResponse<?>>> responsesToBeDelivered;

    private AtomicInteger correlationId;

    public SingleItemToSingleRequestProtocol() {
        this(true);
    }

    public SingleItemToSingleRequestProtocol(boolean betterImplementationPossible) {
        if (betterImplementationPossible) {
            String callStack = Arrays.stream(Thread.currentThread().getStackTrace())
                .map(StackTraceElement::toString)
                .collect(Collectors.joining("\n"));
            LOGGER.warn("Unoptimized Usage of {} detected at:\n{}", this.getClass(), callStack);
            LOGGER.info("Consider implementing item splitting native to the protocol.");
        }
    }

    @Override
    public void channelRegistered(ChannelHandlerContext ctx) throws Exception {
        this.queue = new PendingWriteQueue(ctx);
        this.sentButUnacknowledgedSubContainer = new ConcurrentHashMap<>();
        this.correlationToParentContainer = new ConcurrentHashMap<>();
        this.containerCorrelationIdMap = new ConcurrentHashMap<>();
        this.responsesToBeDelivered = new ConcurrentHashMap<>();
        this.correlationId = new AtomicInteger();
        super.channelRegistered(ctx);
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        this.queue.removeAndWriteAll();
        this.sentButUnacknowledgedSubContainer.clear();
        this.correlationToParentContainer.clear();
        this.containerCorrelationIdMap.clear();
        this.responsesToBeDelivered.clear();
        this.correlationId.set(0);
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

    protected void tryFinish(Integer correlationId, InternalPlcResponse msg, CompletableFuture<InternalPlcResponse<?>> originalResponseFuture) {
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>> subPlcRequestContainer = sentButUnacknowledgedSubContainer.remove(correlationId);
        LOGGER.info("{} got acknowledged", subPlcRequestContainer);
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>> originalPlcRequestContainer = correlationToParentContainer.remove(correlationId);
        if (originalPlcRequestContainer == null) {
            LOGGER.warn("Unrelated package received {}", msg);
            return;
        }
        List<InternalPlcResponse<?>> correlatedResponseItems = responsesToBeDelivered.computeIfAbsent(originalPlcRequestContainer, ignore -> new LinkedList<>());
        correlatedResponseItems.add(msg);
        Set<Integer> integers = containerCorrelationIdMap.get(originalPlcRequestContainer);
        integers.remove(correlationId);
        if (integers.isEmpty()) {
            InternalPlcResponse<?> plcResponse;
            if (originalPlcRequestContainer.getRequest() instanceof InternalPlcReadRequest) {
                InternalPlcReadRequest internalPlcReadRequest = (InternalPlcReadRequest) originalPlcRequestContainer.getRequest();
                HashMap<String, Pair<PlcResponseCode, FieldItem>> fields = new HashMap<>();

                correlatedResponseItems.stream()
                    .map(InternalPlcReadResponse.class::cast)
                    .map(InternalPlcReadResponse::getValues)
                    .forEach(stringPairMap -> stringPairMap.forEach(fields::put));

                plcResponse = new DefaultPlcReadResponse(internalPlcReadRequest, fields);
            } else if (originalPlcRequestContainer.getRequest() instanceof InternalPlcWriteRequest) {
                InternalPlcWriteRequest internalPlcWriteRequest = (InternalPlcWriteRequest) originalPlcRequestContainer.getRequest();
                HashMap<String, PlcResponseCode> values = new HashMap<>();

                correlatedResponseItems.stream()
                    .map(InternalPlcWriteResponse.class::cast)
                    .map(InternalPlcWriteResponse::getValues)
                    .forEach(stringPairMap -> stringPairMap.forEach(values::put));

                plcResponse = new DefaultPlcWriteResponse(internalPlcWriteRequest, values);
            } else {
                throw new PlcRuntimeException("Unknown type detected " + originalPlcRequestContainer.getRequest());
            }
            responsesToBeDelivered.remove(originalPlcRequestContainer);
            containerCorrelationIdMap.remove(originalPlcRequestContainer);
            originalResponseFuture.complete(plcResponse);
        }
    }

    protected void errored(int correlationId, Throwable throwable, CompletableFuture<InternalPlcResponse<?>> originalResponseFuture) {
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse<?>> plcRequestContainer = sentButUnacknowledgedSubContainer.remove(correlationId);
        // TODO: cleanup missing maps as the complete response gets canceled now.
        LOGGER.warn("PlcRequestContainer {} and correlationId {} failed ", plcRequestContainer, correlationId, throwable);
        originalResponseFuture.completeExceptionally(throwable);
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
                    internalPlcReadRequest.getNamedFields().forEach(field -> {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                        Integer tdpu = correlationId.getAndIncrement();
                        CompletableFuture<InternalPlcResponse> correlatedCompletableFuture = new CompletableFuture<>();
                        // Important: don't chain to above as we want the above to be completed not the result of when complete
                        correlatedCompletableFuture
                            .thenApply(InternalPlcResponse.class::cast)
                            .whenComplete((internalPlcResponse, throwable) -> {
                                if (throwable != null) {
                                    errored(tdpu, throwable, in.getResponseFuture());
                                } else {
                                    tryFinish(tdpu, internalPlcResponse, in.getResponseFuture());
                                }
                            });
                        PlcRequestContainer<CorrelatedPlcReadRequest, InternalPlcResponse> correlatedPlcRequestContainer = new PlcRequestContainer<>(CorrelatedPlcReadRequest.of(field, tdpu), correlatedCompletableFuture);
                        correlationToParentContainer.put(tdpu, in);
                        queue.add(correlatedPlcRequestContainer, subPromise);
                        if (!tdpus.add(tdpu)) {
                            throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                        }
                        promiseCombiner.add((Future) subPromise);
                    });
                }
                if (internalPlcFieldRequest instanceof InternalPlcWriteRequest) {
                    InternalPlcWriteRequest internalPlcWriteRequest = (InternalPlcWriteRequest) internalPlcFieldRequest;
                    internalPlcWriteRequest.getNamedFieldTriples().forEach(fieldItemTriple -> {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                        Integer tdpu = correlationId.getAndIncrement();
                        CompletableFuture<InternalPlcResponse> correlatedCompletableFuture = new CompletableFuture<>()
                            .thenApply(InternalPlcResponse.class::cast)
                            .whenComplete((internalPlcResponse, throwable) -> {
                                if (throwable != null) {
                                    errored(tdpu, throwable, in.getResponseFuture());
                                } else {
                                    tryFinish(tdpu, internalPlcResponse, in.getResponseFuture());
                                }
                            });
                        PlcRequestContainer<CorrelatedPlcWriteRequest, InternalPlcResponse> correlatedPlcRequestContainer = new PlcRequestContainer<>(CorrelatedPlcWriteRequest.of(fieldItemTriple, tdpu), correlatedCompletableFuture);
                        correlationToParentContainer.put(tdpu, in);
                        queue.add(correlatedPlcRequestContainer, subPromise);
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

                // Add it to the list of sentButUnacknowledgedSubContainer.
                sentButUnacknowledgedSubContainer.put(correlatedPlcRequest.getTdpu(), currentItem);

                LOGGER.debug("container with id {} sent: ", correlatedPlcRequest.getTdpu(), currentItem);
            }
        }
        ctx.flush();
    }

    protected interface CorrelatedPlcRequest extends InternalPlcRequest {

        int getTdpu();
    }

    protected static class CorrelatedPlcReadRequest extends DefaultPlcReadRequest implements CorrelatedPlcRequest {

        protected final int tdpu;

        protected CorrelatedPlcReadRequest(LinkedHashMap<String, PlcField> fields, int tdpu) {
            super(fields);
            this.tdpu = tdpu;
        }

        protected static CorrelatedPlcReadRequest of(Pair<String, PlcField> stringPlcFieldPair, int tdpu) {
            LinkedHashMap<String, PlcField> fields = new LinkedHashMap<>();
            fields.put(stringPlcFieldPair.getKey(), stringPlcFieldPair.getValue());
            return new CorrelatedPlcReadRequest(fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }

    protected static class CorrelatedPlcWriteRequest extends DefaultPlcWriteRequest implements CorrelatedPlcRequest {

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

    // TODO: maybe export to jmx
    public Map<String, Integer> getStatistics() {
        HashMap<String, Integer> statistics = new HashMap<>();
        statistics.put("queue", queue.size());
        statistics.put("sentButUnacknowledgedSubContainer", sentButUnacknowledgedSubContainer.size());
        statistics.put("correlationToParentContainer", correlationToParentContainer.size());
        statistics.put("containerCorrelationIdMap", containerCorrelationIdMap.size());
        statistics.put("responsesToBeDelivered", responsesToBeDelivered.size());
        statistics.put("currentCorrelationId", correlationId.get());
        return statistics;
    }
}
