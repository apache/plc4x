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
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.exceptions.PlcTimeoutException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.*;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;
import org.apache.plc4x.java.base.model.SubscriptionPlcField;
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
// TODO: add split config so we can override special requests that are allready splitted downstream
public class SingleItemToSingleRequestProtocol extends ChannelDuplexHandler {

    public static final Logger LOGGER = LoggerFactory.getLogger(SingleItemToSingleRequestProtocol.class);

    private final Timer timer;

    private final PlcReader reader;

    private final PlcWriter writer;

    private final PlcSubscriber subscriber;

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

    private SplitConfig splitConfig;

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, PlcSubscriber subscriber, Timer timer) {
        this(reader, writer, subscriber, timer, new SplitConfig());
    }

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, PlcSubscriber subscriber, Timer timer, SplitConfig splitConfig) {
        this(reader, writer, subscriber, timer, splitConfig, true);
    }

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, PlcSubscriber subscriber, Timer timer, SplitConfig splitConfig, boolean betterImplementationPossible) {
        this(reader, writer, subscriber, timer, TimeUnit.SECONDS.toMillis(30), splitConfig, betterImplementationPossible);
    }

    public SingleItemToSingleRequestProtocol(PlcReader reader, PlcWriter writer, PlcSubscriber subscriber, Timer timer, long defaultReceiveTimeout, SplitConfig splitConfig, boolean betterImplementationPossible) {
        this.reader = reader;
        this.writer = writer;
        this.subscriber = subscriber;
        this.timer = timer;
        this.defaultReceiveTimeout = defaultReceiveTimeout;
        this.splitConfig = splitConfig;
        if (this.splitConfig == null) {
            this.splitConfig = new SplitConfig();
        }
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
                HashMap<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> fields = new HashMap<>();

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
            } else if (originalPlcRequestContainer.getRequest() instanceof InternalPlcSubscriptionRequest) {
                InternalPlcSubscriptionRequest internalPlcSubscriptionRequest = (InternalPlcSubscriptionRequest) originalPlcRequestContainer.getRequest();
                HashMap<String, Pair<PlcResponseCode, PlcSubscriptionHandle>> fields = new HashMap<>();

                correlatedResponseItems.stream()
                    .map(InternalPlcSubscriptionResponse.class::cast)
                    .map(InternalPlcSubscriptionResponse::getValues)
                    .forEach(stringPairMap -> stringPairMap.forEach(fields::put));

                plcResponse = new DefaultPlcSubscriptionResponse(internalPlcSubscriptionRequest, fields);
            } else if (originalPlcRequestContainer.getRequest() instanceof InternalPlcUnsubscriptionRequest) {
                InternalPlcUnsubscriptionRequest internalPlcUnsubscriptionRequest = (InternalPlcUnsubscriptionRequest) originalPlcRequestContainer.getRequest();
                plcResponse = new DefaultPlcUnsubscriptionResponse(internalPlcUnsubscriptionRequest);
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
            if (request instanceof InternalPlcFieldRequest && (splitConfig.splitRead || splitConfig.splitWrite || splitConfig.splitSubscription)) {
                InternalPlcFieldRequest internalPlcFieldRequest = (InternalPlcFieldRequest) request;

                if (internalPlcFieldRequest instanceof InternalPlcReadRequest && splitConfig.splitRead) {
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
                } else if (internalPlcFieldRequest instanceof InternalPlcWriteRequest && splitConfig.splitWrite) {
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
                } else if (internalPlcFieldRequest instanceof InternalPlcSubscriptionRequest && splitConfig.splitSubscription) {
                    InternalPlcSubscriptionRequest internalPlcSubscriptionRequest = (InternalPlcSubscriptionRequest) internalPlcFieldRequest;
                    internalPlcSubscriptionRequest.getNamedSubscriptionFields().forEach(field -> {
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
                        PlcRequestContainer<CorrelatedPlcSubscriptionRequest, InternalPlcResponse> correlatedPlcRequestContainer = new PlcRequestContainer<>(CorrelatedPlcSubscriptionRequest.of(subscriber, field, tdpu), correlatedCompletableFuture);
                        correlationToParentContainer.put(tdpu, in);
                        queue.add(correlatedPlcRequestContainer, subPromise);
                        if (!tdpus.add(tdpu)) {
                            throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                        }
                        promiseCombiner.add((Future) subPromise);
                    });
                } else {
                    throw new PlcProtocolException("Unmapped request type " + request.getClass());
                }
            } else if (request instanceof InternalPlcUnsubscriptionRequest && splitConfig.splitUnsubscription) {
                InternalPlcUnsubscriptionRequest internalPlcUnsubscriptionRequest = (InternalPlcUnsubscriptionRequest) request;
                internalPlcUnsubscriptionRequest.getInternalPlcSubscriptionHandles().forEach(handle -> {
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
                    PlcRequestContainer<CorrelatedPlcUnsubscriptionRequest, InternalPlcResponse> correlatedPlcRequestContainer = new PlcRequestContainer<>(CorrelatedPlcUnsubscriptionRequest.of(subscriber, handle, tdpu), correlatedCompletableFuture);
                    correlationToParentContainer.put(tdpu, in);
                    queue.add(correlatedPlcRequestContainer, subPromise);
                    if (!tdpus.add(tdpu)) {
                        throw new IllegalStateException("AtomicInteger should not create duplicated ids: " + tdpu);
                    }
                    promiseCombiner.add((Future) subPromise);
                });
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

        public CorrelatedPlcWriteRequest(PlcWriter writer, LinkedHashMap<String, Pair<PlcField, BaseDefaultFieldItem>> fields, int tdpu) {
            super(writer, fields);
            this.tdpu = tdpu;
        }

        public static CorrelatedPlcWriteRequest of(PlcWriter writer, Triple<String, PlcField, BaseDefaultFieldItem> fieldItemTriple, int tdpu) {
            LinkedHashMap<String, Pair<PlcField, BaseDefaultFieldItem>> fields = new LinkedHashMap<>();
            fields.put(fieldItemTriple.getLeft(), Pair.of(fieldItemTriple.getMiddle(), fieldItemTriple.getRight()));
            return new CorrelatedPlcWriteRequest(writer, fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }

    protected static class CorrelatedPlcSubscriptionRequest extends DefaultPlcSubscriptionRequest implements CorrelatedPlcRequest {

        protected final int tdpu;

        protected CorrelatedPlcSubscriptionRequest(PlcSubscriber subscriber, LinkedHashMap<String, SubscriptionPlcField> fields, int tdpu) {
            super(subscriber, fields);
            this.tdpu = tdpu;
        }

        protected static CorrelatedPlcSubscriptionRequest of(PlcSubscriber subscriber, Pair<String, SubscriptionPlcField> stringPlcFieldPair, int tdpu) {
            LinkedHashMap<String, SubscriptionPlcField> fields = new LinkedHashMap<>();
            fields.put(stringPlcFieldPair.getKey(), stringPlcFieldPair.getValue());
            return new CorrelatedPlcSubscriptionRequest(subscriber, fields, tdpu);
        }

        @Override
        public int getTdpu() {
            return tdpu;
        }
    }

    protected static class CorrelatedPlcUnsubscriptionRequest extends DefaultPlcUnsubscriptionRequest implements CorrelatedPlcRequest {

        protected final int tdpu;

        protected CorrelatedPlcUnsubscriptionRequest(PlcSubscriber subscriber, LinkedList<InternalPlcSubscriptionHandle> subscriptionHandles, int tdpu) {
            super(subscriber, subscriptionHandles);
            this.tdpu = tdpu;
        }

        protected static CorrelatedPlcUnsubscriptionRequest of(PlcSubscriber subscriber, InternalPlcSubscriptionHandle subscriptionHandle, int tdpu) {
            LinkedList<InternalPlcSubscriptionHandle> list = new LinkedList<>();
            list.add(subscriptionHandle);
            return new CorrelatedPlcUnsubscriptionRequest(subscriber, list, tdpu);
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

    public static class SplitConfig {
        private final boolean splitRead;
        private final boolean splitWrite;
        private final boolean splitSubscription;
        private final boolean splitUnsubscription;

        public SplitConfig() {
            splitRead = true;
            splitWrite = true;
            splitSubscription = true;
            splitUnsubscription = true;
        }

        private SplitConfig(boolean splitRead, boolean splitWrite, boolean splitSubscription, boolean splitUnsubscription) {
            this.splitRead = splitRead;
            this.splitWrite = splitWrite;
            this.splitSubscription = splitSubscription;
            this.splitUnsubscription = splitUnsubscription;
        }

        public static SplitConfigBuilder builder() {
            return new SplitConfigBuilder();
        }

        public static class SplitConfigBuilder {
            private boolean splitRead = true;
            private boolean splitWrite = true;
            private boolean splitSubscription = true;
            private boolean splitUnsubscription = true;

            public SplitConfigBuilder splitRead() {
                splitRead = true;
                return this;
            }

            public SplitConfigBuilder dontSplitRead() {
                splitRead = false;
                return this;
            }

            public SplitConfigBuilder splitWrite() {
                splitWrite = true;
                return this;
            }

            public SplitConfigBuilder dontSplitWrite() {
                splitWrite = false;
                return this;
            }

            public SplitConfigBuilder splitSubscribe() {
                splitSubscription = true;
                return this;
            }

            public SplitConfigBuilder dontSplitSubscribe() {
                splitSubscription = false;
                return this;
            }

            public SplitConfigBuilder splitUnsubscribe() {
                splitUnsubscription = true;
                return this;
            }

            public SplitConfigBuilder dontSplitUnsubscribe() {
                splitUnsubscription = false;
                return this;
            }

            public SplitConfig build() {
                return new SplitConfig(splitRead, splitWrite, splitSubscription, splitUnsubscription);
            }
        }
    }
}
