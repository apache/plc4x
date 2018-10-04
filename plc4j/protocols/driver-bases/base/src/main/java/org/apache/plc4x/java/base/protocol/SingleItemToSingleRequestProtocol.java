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
import io.netty.util.Timeout;
import io.netty.util.Timer;
import io.netty.util.concurrent.Future;
import io.netty.util.concurrent.PromiseCombiner;
import org.apache.commons.lang3.tuple.Pair;
import org.apache.commons.lang3.tuple.Triple;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.PlcWriter;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcTimeoutException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.FieldItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.stream.Collectors;

/**
 * This layer can be used to split a {@link org.apache.plc4x.java.api.messages.PlcRequest} which addresses multiple {@link PlcField}s into multiple subsequent {@link org.apache.plc4x.java.api.messages.PlcRequest}s.
 */
public class SingleItemToSingleRequestProtocol extends ChannelDuplexHandler {

    public static final Logger LOGGER = LoggerFactory.getLogger(SingleItemToSingleRequestProtocol.class);

    private final Timer timer;

    private final PlcReader reader;
    private final PlcWriter writer;

    // TODO: maybe better get from map
    private long defaultReceiveTimeout;

    private PendingWriteQueue queue;

    private ConcurrentMap<PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>, Timeout> scheduledTimeouts;

    // Map to track send subcontainers
    private ConcurrentMap<Integer, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> sentButUnacknowledgedSubContainer;

    // Map to map tdpu to original parent container
    // TODO: currently this could be supplied via param, only reason to keep would be for statistics.
    private ConcurrentMap<Integer, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>> correlationToParentContainer;

    // Map to track tdpus per container
    // TODO: currently this could be supplied via param, only reason to keep would be for statistics.
    private ConcurrentMap<PlcRequestContainer<?, ?>, Set<Integer>> containerCorrelationIdMap;

    // Map to track a list of responses per parent container
    private ConcurrentMap<PlcRequestContainer<?, ?>, Queue<InternalPlcResponse>> responsesToBeDelivered;

    private AtomicInteger correlationIdGenerator;

    // TODO: maybe put in map per day or per hour
    private AtomicLong deliveredContainers;

    private AtomicLong erroredContainers;

    private AtomicLong deliveredItems;

    private AtomicLong erroredItems;

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, Timer timer) {
        this(reader, writer, timer, true);
    }

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, Timer timer, boolean betterImplementationPossible) {
        this(reader, writer, timer, TimeUnit.SECONDS.toMillis(30), betterImplementationPossible);
    }

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, Timer timer, long defaultReceiveTimeout, boolean betterImplementationPossible) {
        this.reader = reader;
        this.writer = writer;
        this.timer = timer;
        this.defaultReceiveTimeout = defaultReceiveTimeout;
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
        this.scheduledTimeouts = new ConcurrentHashMap<>();
        this.sentButUnacknowledgedSubContainer = new ConcurrentHashMap<>();
        this.correlationToParentContainer = new ConcurrentHashMap<>();
        this.containerCorrelationIdMap = new ConcurrentHashMap<>();
        this.responsesToBeDelivered = new ConcurrentHashMap<>();
        this.correlationIdGenerator = new AtomicInteger();
        this.deliveredItems = new AtomicLong();
        this.erroredItems = new AtomicLong();
        this.deliveredContainers = new AtomicLong();
        this.erroredContainers = new AtomicLong();
        super.channelRegistered(ctx);
    }

    @Override
    public void channelUnregistered(ChannelHandlerContext ctx) throws Exception {
        this.queue.removeAndWriteAll();
        this.scheduledTimeouts.clear();
        this.sentButUnacknowledgedSubContainer.clear();
        this.correlationToParentContainer.clear();
        this.containerCorrelationIdMap.clear();
        this.responsesToBeDelivered.clear();
        this.correlationIdGenerator.set(0);
        this.deliveredItems.set(0);
        this.erroredItems.set(0);
        this.deliveredContainers.set(0);
        this.erroredContainers.set(0);
        super.channelUnregistered(ctx);
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        // Send everything so we get a proper failure for those pending writes
        this.queue.removeAndWriteAll();
        this.timer.stop();
        this.scheduledTimeouts.clear();
        this.sentButUnacknowledgedSubContainer.clear();
        this.correlationToParentContainer.clear();
        this.containerCorrelationIdMap.clear();
        this.responsesToBeDelivered.clear();
        this.correlationIdGenerator.set(0);
        this.deliveredItems.set(0);
        this.erroredItems.set(0);
        this.deliveredContainers.set(0);
        this.erroredContainers.set(0);
        super.channelInactive(ctx);
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Decoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    protected void tryFinish(Integer currentTdpu, InternalPlcResponse msg, CompletableFuture<InternalPlcResponse> originalResponseFuture) {
        deliveredItems.incrementAndGet();
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> subPlcRequestContainer = sentButUnacknowledgedSubContainer.remove(currentTdpu);
        LOGGER.info("{} got acknowledged", subPlcRequestContainer);
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> originalPlcRequestContainer = correlationToParentContainer.remove(currentTdpu);
        if (originalPlcRequestContainer == null) {
            LOGGER.warn("Unrelated package received {}", msg);
            return;
        }
        Queue<InternalPlcResponse> correlatedResponseItems = responsesToBeDelivered.computeIfAbsent(originalPlcRequestContainer, ignore -> new ConcurrentLinkedQueue<>());
        correlatedResponseItems.add(msg);
        Set<Integer> integers = containerCorrelationIdMap.get(originalPlcRequestContainer);
        integers.remove(currentTdpu);
        if (integers.isEmpty()) {
            deliveredContainers.incrementAndGet();
            Timeout timeout = scheduledTimeouts.remove(originalPlcRequestContainer);
            if (timeout != null) {
                timeout.cancel();
            }

            InternalPlcResponse plcResponse;
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
                errored(currentTdpu, new PlcProtocolException("Unknown type detected " + originalPlcRequestContainer.getRequest().getClass()), originalResponseFuture);
                return;
            }
            responsesToBeDelivered.remove(originalPlcRequestContainer);
            containerCorrelationIdMap.remove(originalPlcRequestContainer);
            originalResponseFuture.complete(plcResponse);
        }
    }

    protected void errored(Integer currentTdpu, Throwable throwable, CompletableFuture<InternalPlcResponse> originalResponseFuture) {
        erroredItems.incrementAndGet();
        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> subPlcRequestContainer = sentButUnacknowledgedSubContainer.remove(currentTdpu);
        LOGGER.info("{} got errored", subPlcRequestContainer);


        PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> originalPlcRequestContainer = correlationToParentContainer.remove(currentTdpu);
        if (originalPlcRequestContainer == null) {
            LOGGER.warn("Unrelated error received tdpu:{}", currentTdpu, throwable);
        } else {
            erroredContainers.incrementAndGet();
            Timeout timeout = scheduledTimeouts.remove(originalPlcRequestContainer);
            if (timeout != null) {
                timeout.cancel();
            }
            responsesToBeDelivered.remove(originalPlcRequestContainer);

            Set<Integer> tdpus = containerCorrelationIdMap.remove(originalPlcRequestContainer);
            if (tdpus != null) {
                tdpus.forEach(tdpu -> {
                    // TODO: technically the other items didn't error so do we increment?
                    //erroredItems.incrementAndGet();
                    sentButUnacknowledgedSubContainer.remove(tdpu);
                    correlationToParentContainer.remove(tdpu);
                });
            }

            LOGGER.warn("PlcRequestContainer {} and correlationId {} failed ", correlationToParentContainer, currentTdpu, throwable);
            originalResponseFuture.completeExceptionally(throwable);
        }
    }

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Encoding
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) throws Exception {
        if (msg instanceof PlcRequestContainer) {
            @SuppressWarnings("unchecked")
            PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> in = (PlcRequestContainer<InternalPlcRequest, InternalPlcResponse>) msg;
            Set<Integer> tdpus = containerCorrelationIdMap.computeIfAbsent(in, plcRequestContainer -> ConcurrentHashMap.newKeySet());

            Timeout timeout = timer.newTimeout(timeout_ -> handleTimeout(timeout_, in, tdpus, System.nanoTime()), defaultReceiveTimeout, TimeUnit.MILLISECONDS);
            scheduledTimeouts.put(in, timeout);

            // Create a promise that has to be called multiple times.
            PromiseCombiner promiseCombiner = new PromiseCombiner();
            InternalPlcRequest request = in.getRequest();
            if (request instanceof InternalPlcFieldRequest) {
                InternalPlcFieldRequest internalPlcFieldRequest = (InternalPlcFieldRequest) request;

                if (internalPlcFieldRequest instanceof InternalPlcReadRequest) {
                    InternalPlcReadRequest internalPlcReadRequest = (InternalPlcReadRequest) internalPlcFieldRequest;
                    internalPlcReadRequest.getNamedFields().forEach(field -> {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                        Integer tdpu = correlationIdGenerator.getAndIncrement();
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
                        PlcRequestContainer<CorrelatedPlcReadRequest, InternalPlcResponse> correlatedPlcRequestContainer = new PlcRequestContainer<>(CorrelatedPlcReadRequest.of(reader, field, tdpu), correlatedCompletableFuture);
                        correlationToParentContainer.put(tdpu, in);
                        queue.add(correlatedPlcRequestContainer, subPromise);
                        if (!tdpus.add(tdpu)) {
                            throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                        }
                        promiseCombiner.add((Future) subPromise);
                    });
                } else if (internalPlcFieldRequest instanceof InternalPlcWriteRequest) {
                    InternalPlcWriteRequest internalPlcWriteRequest = (InternalPlcWriteRequest) internalPlcFieldRequest;
                    internalPlcWriteRequest.getNamedFieldTriples().forEach(fieldItemTriple -> {
                        ChannelPromise subPromise = new DefaultChannelPromise(promise.channel());

                        Integer tdpu = correlationIdGenerator.getAndIncrement();
                        CompletableFuture<InternalPlcResponse> correlatedCompletableFuture = new CompletableFuture<>()
                            .thenApply(InternalPlcResponse.class::cast)
                            .whenComplete((internalPlcResponse, throwable) -> {
                                if (throwable != null) {
                                    errored(tdpu, throwable, in.getResponseFuture());
                                } else {
                                    tryFinish(tdpu, internalPlcResponse, in.getResponseFuture());
                                }
                            });
                        PlcRequestContainer<CorrelatedPlcWriteRequest, InternalPlcResponse> correlatedPlcRequestContainer = new PlcRequestContainer<>(CorrelatedPlcWriteRequest.of(writer, fieldItemTriple, tdpu), correlatedCompletableFuture);
                        correlationToParentContainer.put(tdpu, in);
                        queue.add(correlatedPlcRequestContainer, subPromise);
                        if (!tdpus.add(tdpu)) {
                            throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                        }
                        promiseCombiner.add((Future) subPromise);
                    });
                    // TODO: add sub/unsub
                } else {
                    throw new PlcProtocolException("Unmapped request type " + request.getClass());
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

    @SuppressWarnings("unchecked")
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

    private void handleTimeout(Timeout timeout, PlcRequestContainer<InternalPlcRequest, InternalPlcResponse> in, Set<Integer> tdpus, long scheduledAt) {
        if (timeout.isCancelled()) {
            LOGGER.debug("container {} with timeout {} got canceled", in, timeout);
            return;
        }
        LOGGER.warn("container {} timed out:{}", in, timeout);
        erroredContainers.incrementAndGet();
        responsesToBeDelivered.remove(in);
        containerCorrelationIdMap.remove(in);
        tdpus.forEach(tdpu -> {
            erroredItems.incrementAndGet();
            sentButUnacknowledgedSubContainer.remove(tdpu);
            correlationToParentContainer.remove(tdpu);
        });
        in.getResponseFuture().completeExceptionally(new PlcTimeoutException(System.nanoTime() - scheduledAt));
    }

    protected interface CorrelatedPlcRequest extends InternalPlcRequest {

        int getTdpu();
    }

    protected static class CorrelatedPlcReadRequest extends DefaultPlcReadRequest implements CorrelatedPlcRequest {

        protected final int tdpu;

        protected CorrelatedPlcReadRequest(PlcReader reader, LinkedHashMap<String, PlcField> fields, int tdpu) {
            super(reader, fields);
            this.tdpu = tdpu;
        }

        protected static CorrelatedPlcReadRequest of(PlcReader reader, Pair<String, PlcField> stringPlcFieldPair, int tdpu) {
            LinkedHashMap<String, PlcField> fields = new LinkedHashMap<>();
            fields.put(stringPlcFieldPair.getKey(), stringPlcFieldPair.getValue());
            return new CorrelatedPlcReadRequest(reader, fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }

    protected static class CorrelatedPlcWriteRequest extends DefaultPlcWriteRequest implements CorrelatedPlcRequest {

        private final int tdpu;

        public CorrelatedPlcWriteRequest(PlcWriter writer, LinkedHashMap<String, Pair<PlcField, FieldItem>> fields, int tdpu) {
            super(writer, fields);
            this.tdpu = tdpu;
        }

        public static CorrelatedPlcWriteRequest of(PlcWriter writer, Triple<String, PlcField, FieldItem> fieldItemTriple, int tdpu) {
            LinkedHashMap<String, Pair<PlcField, FieldItem>> fields = new LinkedHashMap<>();
            fields.put(fieldItemTriple.getLeft(), Pair.of(fieldItemTriple.getMiddle(), fieldItemTriple.getRight()));
            return new CorrelatedPlcWriteRequest(writer, fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }

    // TODO: maybe export to jmx
    public Map<String, Number> getStatistics() {
        HashMap<String, Number> statistics = new HashMap<>();
        statistics.put("queue", queue.size());
        statistics.put("sentButUnacknowledgedSubContainer", sentButUnacknowledgedSubContainer.size());
        statistics.put("correlationToParentContainer", correlationToParentContainer.size());
        statistics.put("containerCorrelationIdMap", containerCorrelationIdMap.size());
        statistics.put("responsesToBeDelivered", responsesToBeDelivered.size());
        statistics.put("correlationIdGenerator", correlationIdGenerator.get());
        statistics.put("deliveredItems", deliveredItems.get());
        statistics.put("erroredItems", erroredItems.get());
        statistics.put("deliveredContainers", deliveredContainers.get());
        statistics.put("erroredContainers", erroredContainers.get());
        return statistics;
    }
}
