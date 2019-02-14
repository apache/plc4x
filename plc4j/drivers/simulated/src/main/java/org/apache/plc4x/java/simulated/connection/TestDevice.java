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
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.apache.plc4x.java.base.model.InternalPlcSubscriptionHandle;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.*;
import java.util.function.Consumer;

/**
 * Test device storing its state in memory.
 * Values are stored in a HashMap.
 */
public class TestDevice {

    private final Random random = new Random();

    private final String name;

    private final Map<TestField, BaseDefaultFieldItem> state = new HashMap<>();

    private final Map<PlcSubscriptionHandle, ScheduledFuture<?>> cyclicSubscriptions = new HashMap<>();

    private final Map<PlcSubscriptionHandle, Future<?>> eventSubscriptions = new HashMap<>();

    private final IdentityHashMap<PlcSubscriptionHandle, Pair<TestField, Consumer<BaseDefaultFieldItem>>> changeOfStateSubscriptions = new IdentityHashMap<>();

    private final ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(1);

    private final ExecutorService pool = Executors.newCachedThreadPool();

    public TestDevice(String name) {
        this.name = name;
    }

    public Optional<BaseDefaultFieldItem> get(TestField field) {
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

    public void set(TestField field, BaseDefaultFieldItem value) {
        Objects.requireNonNull(field);
        switch (field.getType()) {
            case STATE:
                changeOfStateSubscriptions.values().stream()
                    .filter(pair -> pair.getKey().equals(field))
                    .map(Pair::getValue)
                    .forEach(baseDefaultFieldItemConsumer -> baseDefaultFieldItemConsumer.accept(value));
                state.put(field, value);
                return;
            case STDOUT:
                System.out.printf("TEST PLC STDOUT [%s]: %s%n", field.getName(), Objects.toString(value.getValues()[0]));
                return;
            case RANDOM:
                System.out.printf("TEST PLC RANDOM [%s]: %s%n", field.getName(), Objects.toString(value.getValues()[0]));
                return;
        }
        throw new IllegalArgumentException("Unsupported field type: " + field.getType().name());
    }

    @SuppressWarnings("unchecked")
    private BaseDefaultFieldItem randomValue(Class<?> type) {
        Object result = null;

        if (type.equals(Byte.class)) {
            result = (byte) random.nextInt(1 << 8);
        }

        if (type.equals(Short.class)) {
            result = (short) random.nextInt(1 << 16);
        }

        if (type.equals(Integer.class)) {
            result = random.nextInt();
        }

        if (type.equals(Long.class)) {
            result = random.nextLong();
        }

        if (type.equals(Float.class)) {
            result = random.nextFloat();
        }

        if (type.equals(Double.class)) {
            result = random.nextDouble();
        }

        if (type.equals(Boolean.class)) {
            result = random.nextBoolean();
        }

        if (type.equals(String.class)) {
            int length = random.nextInt(100);
            StringBuilder sb = new StringBuilder(length);
            for (int i = 0; i < length; i++) {
                char c = (char) ('a' + random.nextInt(26));
                sb.append(c);
            }
            result = sb.toString();
        }

        if (type.equals(byte[].class)) {
            int length = random.nextInt(100);
            byte[] bytes = new byte[length];
            random.nextBytes(bytes);
            result = bytes;
        }

        return new TestFieldItem(new Object[]{result});
    }

    @Override
    public String toString() {
        return name;
    }

    public void addCyclicSubscription(Consumer<BaseDefaultFieldItem> consumer, PlcSubscriptionHandle handle, TestField plcField, Duration duration) {
        ScheduledFuture<?> scheduledFuture = scheduler.scheduleAtFixedRate(() -> {
            BaseDefaultFieldItem baseDefaultFieldItem = state.get(plcField);
            if (baseDefaultFieldItem == null) {
                return;
            }
            consumer.accept(baseDefaultFieldItem);
        }, duration.toMillis(), duration.toMillis(), TimeUnit.MILLISECONDS);
        cyclicSubscriptions.put(handle, scheduledFuture);
    }

    public void addChangeOfStateSubscription(Consumer<BaseDefaultFieldItem> consumer, PlcSubscriptionHandle handle, TestField plcField) {
        changeOfStateSubscriptions.put(handle, Pair.of(plcField, consumer));
    }

    public void addEventSubscription(Consumer<BaseDefaultFieldItem> consumer, PlcSubscriptionHandle handle, TestField plcField) {
        Future<?> submit = pool.submit(() -> {
            while (!Thread.currentThread().isInterrupted()) {
                BaseDefaultFieldItem baseDefaultFieldItem = state.get(plcField);
                if (baseDefaultFieldItem == null) {
                    continue;
                }
                consumer.accept(baseDefaultFieldItem);
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
