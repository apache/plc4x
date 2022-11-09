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
package org.apache.plc4x.java.simulated.connection;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.apache.plc4x.java.simulated.readwrite.DataItem;
import org.apache.plc4x.java.simulated.readwrite.SimulatedDataTypeSizes;
import org.apache.plc4x.java.spi.generation.*;

import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionTag;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.security.SecureRandom;
import java.time.Duration;
import java.util.*;
import java.util.concurrent.*;
import java.util.function.Consumer;

/**
 * Test device storing its state in memory.
 * Values are stored in a HashMap.
 */
public class SimulatedDevice {

    private static final Logger LOGGER = LoggerFactory.getLogger(SimulatedDevice.class);

    private final Random random = new SecureRandom();

    private final String name;

    private final Map<SimulatedTag, PlcValue> state = new HashMap<>();

    private final Map<PlcSubscriptionHandle, ScheduledFuture<?>> cyclicSubscriptions = new HashMap<>();

    private final Map<PlcSubscriptionHandle, Future<?>> eventSubscriptions = new HashMap<>();

    private final IdentityHashMap<PlcSubscriptionHandle, Pair<SimulatedTag, Consumer<PlcValue>>> changeOfStateSubscriptions = new IdentityHashMap<>();

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);

    private final ExecutorService pool = Executors.newCachedThreadPool();

    public SimulatedDevice(String name) {
        this.name = name;
    }

    public Optional<PlcValue> get(SimulatedTag tag) {
        LOGGER.debug("getting tag {}", tag);
        Objects.requireNonNull(tag);
        switch (tag.getType()) {
            case STATE:
                return Optional.ofNullable(state.get(tag));
            case RANDOM:
                return Optional.ofNullable(randomValue(tag));
            case STDOUT:
                return Optional.empty();
        }
        throw new IllegalArgumentException("Unsupported tag type: " + tag.getType().name());
    }

    public void set(SimulatedTag tag, PlcValue value) {
        LOGGER.debug("setting tag {} to {}", tag, value);
        Objects.requireNonNull(tag);
        switch (tag.getType()) {
            case STATE:
                changeOfStateSubscriptions.values().stream()
                    .filter(pair -> pair.getKey().equals(tag))
                    .map(Pair::getValue)
                    .peek(plcValueConsumer -> LOGGER.debug("{} is getting notified with {}", plcValueConsumer, value))
                    .forEach(baseDefaultPlcValueConsumer -> baseDefaultPlcValueConsumer.accept(value));
                state.put(tag, value);
                return;
            case STDOUT:
                LOGGER.info("TEST PLC STDOUT [{}]: {}", tag.getName(), value.toString());
                return;
            case RANDOM:
                switch (tag.getPlcValueType()) {
                    case STRING:
                    case WSTRING:
                        break;
                    default:
                        try {
                            final int lengthInBits = DataItem.getLengthInBits(value, tag.getPlcValueType().name(), tag.getArrayInfo().get(0).GetSize());
                            final WriteBufferByteBased writeBuffer = new WriteBufferByteBased((int) Math.ceil(((float) lengthInBits) / 8.0f));
                            DataItem.staticSerialize(writeBuffer, value, tag.getPlcValueType().name(), tag.getArrayInfo().get(0).GetSize(), ByteOrder.BIG_ENDIAN);
                        } catch (SerializationException e) {
                            LOGGER.info("Write failed");
                        }
                }
                LOGGER.info("TEST PLC RANDOM [{}]: {}", tag.getName(), value);
                return;
        }
        throw new IllegalArgumentException("Unsupported tag type: " + tag.getType().name());
    }

    private PlcValue randomValue(SimulatedTag tag) {

        short tagDataTypeSize = SimulatedDataTypeSizes.valueOf(tag.getPlcValueType().name()).getDataTypeSize();

        byte[] b = new byte[tagDataTypeSize * tag.getArrayInfo().get(0).GetSize()];
        random.nextBytes(b);

        ReadBuffer io = new ReadBufferByteBased(b);
        try {
            return DataItem.staticParse(io, tag.getPlcValueType().name(), tag.getArrayInfo().get(0).GetSize());
        } catch (ParseException e) {
            return null;
        }
    }

    @Override
    public String toString() {
        return name;
    }

    public void addCyclicSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, PlcSubscriptionTag subscriptionTag, Duration duration) {
        LOGGER.debug("Adding cyclic subscription: {}, {}, {}, {}", consumer, handle, subscriptionTag, duration);
        assert subscriptionTag instanceof DefaultPlcSubscriptionTag;
        ScheduledFuture<?> scheduledFuture = scheduler.scheduleAtFixedRate(() -> {
            PlcTag innerPlcTag = ((DefaultPlcSubscriptionTag) subscriptionTag).getTag();
            assert innerPlcTag instanceof SimulatedTag;
            PlcValue baseDefaultPlcValue = state.get(innerPlcTag);
            if (baseDefaultPlcValue == null) {
                return;
            }
            consumer.accept(baseDefaultPlcValue);
        }, duration.toMillis(), duration.toMillis(), TimeUnit.MILLISECONDS);
        cyclicSubscriptions.put(handle, scheduledFuture);
    }

    public void addChangeOfStateSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, PlcSubscriptionTag subscriptionTag) {
        LOGGER.debug("Adding change of state subscription: {}, {}, {}", consumer, handle, subscriptionTag);
        changeOfStateSubscriptions.put(handle, Pair.of((SimulatedTag) ((DefaultPlcSubscriptionTag) subscriptionTag).getTag(), consumer));
    }

    public void addEventSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, PlcSubscriptionTag subscriptionTag) {
        LOGGER.debug("Adding event subscription: {}, {}, {}", consumer, handle, subscriptionTag);
        assert subscriptionTag instanceof DefaultPlcSubscriptionTag;
        Future<?> submit = pool.submit(() -> {
            LOGGER.debug("WORKER: starting for {}, {}, {}", consumer, handle, subscriptionTag);
            while (!Thread.currentThread().isInterrupted()) {
                LOGGER.debug("WORKER: running for {}, {}, {}", consumer, handle, subscriptionTag);
                PlcTag innerPlcTag = ((DefaultPlcSubscriptionTag) subscriptionTag).getTag();
                assert innerPlcTag instanceof SimulatedTag;
                PlcValue baseDefaultPlcValue = state.get(innerPlcTag);
                if (baseDefaultPlcValue == null) {
                    LOGGER.debug("WORKER: no value for {}, {}, {}", consumer, handle, subscriptionTag);
                    continue;
                }
                LOGGER.debug("WORKER: accepting {} for {}, {}, {}", baseDefaultPlcValue, consumer, handle, subscriptionTag);
                consumer.accept(baseDefaultPlcValue);
                try {
                    long sleepTime = Math.min(random.nextInt((int) TimeUnit.SECONDS.toNanos(5)), TimeUnit.MILLISECONDS.toNanos(500));
                    LOGGER.debug("WORKER: sleeping {} milliseconds for {}, {}, {}", TimeUnit.NANOSECONDS.toMillis(sleepTime), consumer, handle, subscriptionTag);
                    TimeUnit.NANOSECONDS.sleep(sleepTime);
                } catch (InterruptedException ignore) {
                    Thread.currentThread().interrupt();
                    LOGGER.debug("WORKER: got interrupted for {}, {}, {}", consumer, handle, subscriptionTag);
                    return;
                }
            }
        });

        eventSubscriptions.put(handle, submit);
    }

    public void removeHandles(Collection<? extends PlcSubscriptionHandle> internalPlcSubscriptionHandles) {
        LOGGER.debug("remove handles {}", internalPlcSubscriptionHandles);
        internalPlcSubscriptionHandles.forEach(handle -> {
            ScheduledFuture<?> remove = cyclicSubscriptions.remove(handle);
            if (remove == null) {
                LOGGER.debug("nothing to cancel {}", handle);
                return;
            }
            remove.cancel(true);
        });
        internalPlcSubscriptionHandles.forEach(handle -> {
            Future<?> remove = eventSubscriptions.remove(handle);
            if (remove == null) {
                LOGGER.debug("nothing to cancel {}", handle);
                return;
            }
            remove.cancel(true);
        });
        internalPlcSubscriptionHandles.forEach(changeOfStateSubscriptions::remove);
    }
}
