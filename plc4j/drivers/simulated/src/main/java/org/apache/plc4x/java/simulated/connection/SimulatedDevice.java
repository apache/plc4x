/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.model.PlcSubscriptionField;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.simulated.field.SimulatedField;
import org.apache.plc4x.java.simulated.readwrite.io.DataItemIO;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.ReadBuffer;

import org.apache.plc4x.java.spi.generation.ReadBufferByteBased;
import org.apache.plc4x.java.spi.model.DefaultPlcSubscriptionField;

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

    private final Map<SimulatedField, PlcValue> state = new HashMap<>();

    private final Map<PlcSubscriptionHandle, ScheduledFuture<?>> cyclicSubscriptions = new HashMap<>();

    private final Map<PlcSubscriptionHandle, Future<?>> eventSubscriptions = new HashMap<>();

    private final IdentityHashMap<PlcSubscriptionHandle, Pair<SimulatedField, Consumer<PlcValue>>> changeOfStateSubscriptions = new IdentityHashMap<>();

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);

    private final ExecutorService pool = Executors.newCachedThreadPool();

    public SimulatedDevice(String name) {
        this.name = name;
    }

    public Optional<PlcValue> get(SimulatedField field) {
        LOGGER.debug("getting field {}", field);
        Objects.requireNonNull(field);
        switch (field.getType()) {
            case STATE:
                return Optional.ofNullable(state.get(field));
            case RANDOM:
                return Optional.ofNullable(randomValue(field));
            case STDOUT:
                return Optional.empty();
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    public void set(SimulatedField field, PlcValue value) {
        LOGGER.debug("setting field {} to {}", field, value);
        Objects.requireNonNull(field);
        switch (field.getType()) {
            case STATE:
                changeOfStateSubscriptions.values().stream()
                    .filter(pair -> pair.getKey().equals(field))
                    .map(Pair::getValue)
                    .peek(plcValueConsumer -> {
                        LOGGER.debug("{} is getting notified with {}", plcValueConsumer, value);
                    })
                    .forEach(baseDefaultPlcValueConsumer -> baseDefaultPlcValueConsumer.accept(value));
                state.put(field, value);
                return;
            case STDOUT:
                LOGGER.info("TEST PLC STDOUT [{}]: {}", field.getName(), value.getString());
                return;
            case RANDOM:
                switch (field.getPlcDataType()) {
                    case "STRING":
                    case "WSTRING":
                        break;
                    default:
                        try {
                            DataItemIO.staticSerialize(value, field.getPlcDataType(), field.getNumberOfElements(), false);
                        } catch (ParseException e) {
                            LOGGER.info("Write failed");
                        }
                }
                LOGGER.info("TEST PLC RANDOM [{}]: {}", field.getName(), value);
                return;
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    private PlcValue randomValue(SimulatedField field) {
        short fieldDataTypeSize = field.getDataType().getDataTypeSize();

        byte[] b = new byte[fieldDataTypeSize * field.getNumberOfElements()];
        random.nextBytes(b);

        ReadBuffer io = new ReadBufferByteBased(b);
        try {
            return DataItemIO.staticParse(io, field.getPlcDataType(), field.getNumberOfElements());
        } catch (ParseException e) {
            return null;
        }
    }

    @Override
    public String toString() {
        return name;
    }

    public void addCyclicSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, PlcSubscriptionField plcField, Duration duration) {
        LOGGER.debug("Adding cyclic subscription: {}, {}, {}, {}", consumer, handle, plcField, duration);
        assert plcField instanceof DefaultPlcSubscriptionField;
        ScheduledFuture<?> scheduledFuture = scheduler.scheduleAtFixedRate(() -> {
            PlcField innerPlcField = ((DefaultPlcSubscriptionField) plcField).getPlcField();
            assert innerPlcField instanceof SimulatedField;
            PlcValue baseDefaultPlcValue = state.get(innerPlcField);
            if (baseDefaultPlcValue == null) {
                return;
            }
            consumer.accept(baseDefaultPlcValue);
        }, duration.toMillis(), duration.toMillis(), TimeUnit.MILLISECONDS);
        cyclicSubscriptions.put(handle, scheduledFuture);
    }

    public void addChangeOfStateSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, PlcSubscriptionField plcField) {
        LOGGER.debug("Adding change of state subscription: {}, {}, {}", consumer, handle, plcField);
        changeOfStateSubscriptions.put(handle, Pair.of((SimulatedField) ((DefaultPlcSubscriptionField) plcField).getPlcField(), consumer));
    }

    public void addEventSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, PlcSubscriptionField plcField) {
        LOGGER.debug("Adding event subscription: {}, {}, {}", consumer, handle, plcField);
        assert plcField instanceof DefaultPlcSubscriptionField;
        Future<?> submit = pool.submit(() -> {
            LOGGER.debug("WORKER: starting for {}, {}, {}", consumer, handle, plcField);
            while (!Thread.currentThread().isInterrupted()) {
                LOGGER.debug("WORKER: running for {}, {}, {}", consumer, handle, plcField);
                PlcField innerPlcField = ((DefaultPlcSubscriptionField) plcField).getPlcField();
                assert innerPlcField instanceof SimulatedField;
                PlcValue baseDefaultPlcValue = state.get(innerPlcField);
                if (baseDefaultPlcValue == null) {
                    LOGGER.debug("WORKER: no value for {}, {}, {}", consumer, handle, plcField);
                    continue;
                }
                LOGGER.debug("WORKER: accepting {} for {}, {}, {}", baseDefaultPlcValue, consumer, handle, plcField);
                consumer.accept(baseDefaultPlcValue);
                try {
                    long sleepTime = Math.min(random.nextInt((int) TimeUnit.SECONDS.toNanos(5)), TimeUnit.MILLISECONDS.toNanos(500));
                    LOGGER.debug("WORKER: sleeping {} milliseconds for {}, {}, {}", TimeUnit.NANOSECONDS.toMillis(sleepTime), consumer, handle, plcField);
                    TimeUnit.NANOSECONDS.sleep(sleepTime);
                } catch (InterruptedException ignore) {
                    Thread.currentThread().interrupt();
                    LOGGER.debug("WORKER: got interrupted for {}, {}, {}", consumer, handle, plcField);
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
