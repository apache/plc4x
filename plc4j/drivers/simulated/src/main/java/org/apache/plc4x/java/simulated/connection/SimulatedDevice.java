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
package org.apache.plc4x.java.simulated.connection;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValues;
import org.apache.plc4x.java.api.value.IEC61131ValueHandler;
import org.apache.plc4x.java.api.value.PlcValueHandler;
import org.apache.plc4x.java.simulated.field.SimulatedField;
import org.apache.plc4x.java.spi.model.InternalPlcSubscriptionHandle;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.*;
import java.util.function.Consumer;

/**
 * Test device storing its state in memory.
 * Values are stored in a HashMap.
 */
public class SimulatedDevice {

    private final Random random = new Random();

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
        Objects.requireNonNull(field);
        switch (field.getType()) {
            case STATE:
                return Optional.ofNullable(state.get(field));
            case RANDOM:
                return Optional.of(randomValue(field.getDataType()));
            case STDOUT:
                return Optional.empty();
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    public void set(SimulatedField field, PlcValue value) {
        Objects.requireNonNull(field);
        switch (field.getType()) {
            case STATE:
                changeOfStateSubscriptions.values().stream()
                    .filter(pair -> pair.getKey().equals(field))
                    .map(Pair::getValue)
                    .forEach(baseDefaultPlcValueConsumer -> baseDefaultPlcValueConsumer.accept(value));
                state.put(field, value);
                return;
            case STDOUT:
                System.out.printf("TEST PLC STDOUT [%s]: %s%n", field.getName(), value.getString());
                return;
            case RANDOM:
                System.out.printf("TEST PLC RANDOM [%s]: %s%n", field.getName(), value.getString());
                return;
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    @SuppressWarnings("unchecked")
    private PlcValue randomValue(Class<?> type) {
        Object result = null;

        if (type.equals(Byte.class)) {
            return IEC61131ValueHandler.newPlcValue((byte) random.nextInt(1 << 8));
        }

        if (type.equals(Short.class)) {
            return IEC61131ValueHandler.newPlcValue((short) random.nextInt(1 << 16));
        }

        if (type.equals(Integer.class)) {
            return IEC61131ValueHandler.newPlcValue(random.nextInt());
        }

        if (type.equals(Long.class)) {
            return IEC61131ValueHandler.newPlcValue(random.nextLong());
        }

        if (type.equals(Float.class)) {
            return IEC61131ValueHandler.newPlcValue(random.nextFloat());
        }

        if (type.equals(Double.class)) {
            return IEC61131ValueHandler.newPlcValue(random.nextDouble());
        }

        if (type.equals(Boolean.class)) {
            return IEC61131ValueHandler.newPlcValue(random.nextBoolean());
        }

        if (type.equals(String.class)) {
            int length = random.nextInt(100);
            StringBuilder sb = new StringBuilder(length);
            for (int i = 0; i < length; i++) {
                char c = (char) ('a' + random.nextInt(26));
                sb.append(c);
            }
            return IEC61131ValueHandler.newPlcValue(sb.toString());
        }

        return null;
    }

    @Override
    public String toString() {
        return name;
    }

    public void addCyclicSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, SimulatedField plcField, Duration duration) {
        ScheduledFuture<?> scheduledFuture = scheduler.scheduleAtFixedRate(() -> {
            PlcValue baseDefaultPlcValue = state.get(plcField);
            if (baseDefaultPlcValue == null) {
                return;
            }
            consumer.accept(baseDefaultPlcValue);
        }, duration.toMillis(), duration.toMillis(), TimeUnit.MILLISECONDS);
        cyclicSubscriptions.put(handle, scheduledFuture);
    }

    public void addChangeOfStateSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, SimulatedField plcField) {
        changeOfStateSubscriptions.put(handle, Pair.of(plcField, consumer));
    }

    public void addEventSubscription(Consumer<PlcValue> consumer, PlcSubscriptionHandle handle, SimulatedField plcField) {
        Future<?> submit = pool.submit(() -> {
            while (!Thread.currentThread().isInterrupted()) {
                PlcValue baseDefaultPlcValue = state.get(plcField);
                if (baseDefaultPlcValue == null) {
                    continue;
                }
                consumer.accept(baseDefaultPlcValue);
                try {
                    TimeUnit.SECONDS.sleep((long) (Math.random() * 10));
                } catch (InterruptedException ignore) {
                    Thread.currentThread().interrupt();
                    return;
                }
            }
        });

        eventSubscriptions.put(handle, submit);
    }

    public void removeHandles(Collection<? extends InternalPlcSubscriptionHandle> internalPlcSubscriptionHandles) {
        internalPlcSubscriptionHandles.forEach(handle -> {
            ScheduledFuture<?> remove = cyclicSubscriptions.remove(handle);
            if (remove == null) {
                return;
            }
            remove.cancel(true);
        });
        internalPlcSubscriptionHandles.forEach(handle -> {
            Future<?> remove = eventSubscriptions.remove(handle);
            if (remove == null) {
                return;
            }
            remove.cancel(true);
        });
        internalPlcSubscriptionHandles.forEach(changeOfStateSubscriptions::remove);
    }
}
