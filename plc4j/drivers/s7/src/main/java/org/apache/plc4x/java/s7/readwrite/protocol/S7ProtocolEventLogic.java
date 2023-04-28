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
package org.apache.plc4x.java.s7.readwrite.protocol;

import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.events.*;
import org.apache.plc4x.java.s7.readwrite.EventType;
import org.apache.plc4x.java.s7.readwrite.S7ParameterModeTransition;
import org.apache.plc4x.java.s7.readwrite.S7PayloadDiagnosticMessage;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.slf4j.LoggerFactory;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

public class S7ProtocolEventLogic implements PlcSubscriber {

    private final org.slf4j.Logger logger = LoggerFactory.getLogger(S7ProtocolEventLogic.class);

    private final BlockingQueue<?> eventQueue;
    private final BlockingQueue<PlcSubscriptionEvent> dispachQueue = new ArrayBlockingQueue<>(1024);

    private final Map<EventType, Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>>> mapIndex = new HashMap<>();
    private final Map<EventType, PlcSubscriptionHandle> eventTypeHandles = new HashMap<>();

    private final Runnable runProcessor;
    private final Runnable runDispacher;

    private final Thread processor;
    private final Thread dispacher;

    public S7ProtocolEventLogic(BlockingQueue<?> eventQueue) {
        this.eventQueue = eventQueue;
        runProcessor = new ObjectProcessor(eventQueue, dispachQueue);
        runDispacher = new EventDispatcher(dispachQueue);
        processor = new Thread(runProcessor);
        dispacher = new Thread(runDispacher);
    }

    public void start() {
        processor.start();
        dispacher.start();
    }

    public void stop() {
        ((ObjectProcessor) runProcessor).doShutdown();
        ((EventDispatcher) runDispacher).doShutdown();
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        throw new UnsupportedOperationException("Not supported yet."); //To change body of generated methods, choose Tools | Templates.
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = null;
        S7PlcSubscriptionHandle handle = (S7PlcSubscriptionHandle) handles.toArray()[0];
        if (!mapIndex.containsKey(handle.getEventType())) {
            mapConsumers = new HashMap<>();
            mapIndex.put(handle.getEventType(), mapConsumers);
        }
        mapConsumers = mapIndex.get(handle.getEventType());
        //TODO: Check the implementation of "DefaultPlcConsumerRegistration". List<> vs Collection<>
        DefaultPlcConsumerRegistration registro = new DefaultPlcConsumerRegistration(this,
            consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        mapConsumers.put(registro, consumer);
        return registro;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        S7PlcSubscriptionHandle handle = (S7PlcSubscriptionHandle) registration.getSubscriptionHandles().get(0);
        Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(handle.getEventType());
        mapConsumers.remove(registration);
    }

    private static class ObjectProcessor implements Runnable {

        private static final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(ObjectProcessor.class);

        private final BlockingQueue<?> eventQueue;
        private final BlockingQueue<PlcSubscriptionEvent> dispatchQueue;
        private boolean shutdown = false;
        private final int delay = 100;

        public ObjectProcessor(BlockingQueue<?> eventQueue, BlockingQueue<PlcSubscriptionEvent> dispatchQueue) {
            this.eventQueue = eventQueue;
            this.dispatchQueue = dispatchQueue;
        }

        @Override
        public void run() {
            while (!shutdown) {
                try {
                    Object obj = eventQueue.poll(delay, TimeUnit.MILLISECONDS);
                    if (obj == null) {
                        continue;
                    }
                    if (obj instanceof S7ParameterModeTransition) {
                        S7ModeEvent modeEvent = new S7ModeEvent((S7ParameterModeTransition) obj);
                        dispatchQueue.add(modeEvent);
                    } else if (obj instanceof S7PayloadDiagnosticMessage) {
                        S7PayloadDiagnosticMessage msg = (S7PayloadDiagnosticMessage) obj;
                        if ((msg.getEventId() >= 0x0A000) & (msg.getEventId() <= 0x0BFFF)) {
                            S7UserEvent userEvent = new S7UserEvent(msg);
                            dispatchQueue.add(userEvent);
                        } else {
                            S7SysEvent sysEvent = new S7SysEvent(msg);
                            dispatchQueue.add(sysEvent);
                        }
                    } else if (obj instanceof S7CyclicEvent) {
                        dispatchQueue.add((S7CyclicEvent) obj);
                    } else {
                        S7AlarmEvent alarmEvent = new S7AlarmEvent(obj);
                        dispatchQueue.add(alarmEvent);
                    }
                } catch (InterruptedException e) {
                    LOGGER.error("oh no", e);
                    Thread.currentThread().interrupt();
                }
            }
            LOGGER.trace("ObjectProcessor Bye!");
        }

        public void doShutdown() {
            shutdown = true;
        }
    }

    private class EventDispatcher implements Runnable {

        private final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(EventDispatcher.class);

        private final BlockingQueue<PlcSubscriptionEvent> dispachQueue;
        private boolean shutdown = false;
        private final int delay = 100;
        private PlcSubscriptionEvent cycDelayedObject = null;

        public EventDispatcher(BlockingQueue<PlcSubscriptionEvent> dispachQueue) {
            this.dispachQueue = dispachQueue;
        }

        @Override
        public void run() {
            try {
                Thread.sleep(500);
            } catch (InterruptedException e) {
                logger.warn("Error sleeping", e);
                Thread.currentThread().interrupt();
            }
            while (!shutdown) {
                try {
                    PlcSubscriptionEvent obj = dispachQueue.poll(delay, TimeUnit.MILLISECONDS);
                    if (obj != null) {
                        if (obj instanceof S7ModeEvent) {
                            if (mapIndex.containsKey(EventType.MODE)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.MODE);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7UserEvent) {
                            if (mapIndex.containsKey(EventType.USR)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.USR);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7SysEvent) {
                            if (mapIndex.containsKey(EventType.SYS)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.SYS);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7AlarmEvent) {
                            if (mapIndex.containsKey(EventType.ALM)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.ALM);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7CyclicEvent) {
                            if (mapIndex.containsKey(EventType.CYC)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.CYC);
                                if (cycDelayedObject != null) {
                                    mapConsumers.forEach((x, y) -> y.accept(cycDelayedObject));
                                    cycDelayedObject = null;
                                }
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            } else {
                                cycDelayedObject = obj;
                            }
                        }
                    }
                } catch (Exception e) {
                    LOGGER.error("oh no", e);
                    Thread.currentThread().interrupt();
                }
            }
            LOGGER.trace("EventDispatcher Bye!");
        }

        public void doShutdown() {
            shutdown = true;
        }

    }


}
