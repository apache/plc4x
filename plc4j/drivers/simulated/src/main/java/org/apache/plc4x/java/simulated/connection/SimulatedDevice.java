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
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.simulated.readwrite.DataItem;
import org.apache.plc4x.java.simulated.readwrite.SimulatedDataTypeSizes;
import org.apache.plc4x.java.simulated.tag.SimulatedTag;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.drivers.messages.DefaultPlcSubscriptionTag;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.security.SecureRandom;
import java.time.Duration;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.IdentityHashMap;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;
import java.util.Random;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.ScheduledFuture;
import java.util.concurrent.TimeUnit;
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
                            int numElements = tag.getArrayInfo().isEmpty() ? 1 : tag.getArrayInfo().get(0).getSize();
                            int lengthInBits = DataItem.getLengthInBits(value, tag.getPlcValueType().name(), numElements);
                            int sizeInBytes = (int) Math.ceil(((float) lengthInBits) / 8.0f);
                            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(new byte[sizeInBytes]);
                            DataItem.staticSerialize(writeBuffer, value, tag.getPlcValueType().name(), numElements);
                        } catch (BufferException e) {
                            LOGGER.info("Write failed", e);
                        }
                }
                LOGGER.info("TEST PLC RANDOM [{}]: {}", tag.getName(), value);
                return;
        }
        throw new IllegalArgumentException("Unsupported tag type: " + tag.getType().name());
    }

    private PlcValue randomValue(SimulatedTag tag) {
        int numElements = tag.getArrayInfo().isEmpty() ? 1 : tag.getArrayInfo().get(0).getSize();
        // Strings are length-prefixed, so feeding random bytes to the parser almost always
        // yields a bogus length and a BufferException (which would surface as NOT_FOUND).
        // Generate a valid random string directly instead - mirrors how the write path
        // special-cases STRING/WSTRING in set().
        switch (tag.getPlcValueType()) {
            case STRING:
            case WSTRING:
                if (numElements == 1) {
                    return new PlcSTRING(randomString());
                }
                List<PlcValue> elements = new ArrayList<>(numElements);
                for (int i = 0; i < numElements; i++) {
                    elements.add(new PlcSTRING(randomString()));
                }
                return new PlcList(elements);
            default:
                break;
        }
        short tagDataTypeSize = SimulatedDataTypeSizes.valueOf(tag.getPlcValueType().name()).getDataTypeSize();
        byte[] b = new byte[tagDataTypeSize * numElements];
        random.nextBytes(b);
        ReadBufferByteBased io = new ReadBufferByteBased(b);
        try {
            return DataItem.staticParse(io, tag.getPlcValueType().name(), numElements);
        } catch (BufferException e) {
            return null;
        }
    }

    private String randomString() {
        final String alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        int length = 1 + random.nextInt(10);
        StringBuilder sb = new StringBuilder(length);
        for (int i = 0; i < length; i++) {
            sb.append(alphabet.charAt(random.nextInt(alphabet.length())));
        }
        return sb.toString();
    }

    @Override
    public String toString() {
        return name;
    }

    public void addCyclicSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle,
                                      PlcSubscriptionTag subscriptionTag, Duration duration) {
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

    public void addChangeOfStateSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle,
                                             PlcSubscriptionTag subscriptionTag) {
        LOGGER.debug("Adding change of state subscription: {}, {}, {}", consumer, handle, subscriptionTag);
        changeOfStateSubscriptions.put(handle,
            Pair.of((SimulatedTag) ((DefaultPlcSubscriptionTag) subscriptionTag).getTag(), consumer));
    }

    public void addEventSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle,
                                     PlcSubscriptionTag subscriptionTag) {
        LOGGER.debug("Adding event subscription: {}, {}, {}", consumer, handle, subscriptionTag);
        assert subscriptionTag instanceof DefaultPlcSubscriptionTag;
        Future<?> submit = pool.submit(() -> {
            while (!Thread.currentThread().isInterrupted()) {
                PlcTag innerPlcTag = ((DefaultPlcSubscriptionTag) subscriptionTag).getTag();
                assert innerPlcTag instanceof SimulatedTag;
                PlcValue baseDefaultPlcValue = state.get(innerPlcTag);
                if (baseDefaultPlcValue != null) {
                    consumer.accept(baseDefaultPlcValue);
                }
                try {
                    long sleepTime = Math.min(random.nextInt((int) TimeUnit.SECONDS.toNanos(5)),
                        TimeUnit.MILLISECONDS.toNanos(500));
                    TimeUnit.NANOSECONDS.sleep(sleepTime);
                } catch (InterruptedException ignore) {
                    Thread.currentThread().interrupt();
                    return;
                }
            }
        });
        eventSubscriptions.put(handle, submit);
    }

    public void removeHandles(Collection<? extends PlcSubscriptionHandle> handles) {
        LOGGER.debug("remove handles {}", handles);
        handles.forEach(handle -> {
            ScheduledFuture<?> remove = cyclicSubscriptions.remove(handle);
            if (remove != null) {
                remove.cancel(true);
            }
        });
        handles.forEach(handle -> {
            Future<?> remove = eventSubscriptions.remove(handle);
            if (remove != null) {
                remove.cancel(true);
            }
        });
        handles.forEach(changeOfStateSubscriptions::remove);
    }

    public void shutdown() {
        scheduler.shutdownNow();
        pool.shutdownNow();
    }

}
