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

import org.apache.commons.lang3.concurrent.BasicThreadFactory;
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
import java.util.logging.Level;
import java.util.logging.Logger;

public class S7ProtocolEventLogic implements PlcSubscriber {
    private final org.slf4j.Logger logger = LoggerFactory.getLogger(S7ProtocolEventLogic.class);
    private static final int DEFAULT_DELAY = 100;

    private final Map<EventType, Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>>> mapIndex = new HashMap<>();

    private final ObjectProcessor runProcessor;
    private final EventDispatcher runDispatcher;


    private final Thread processor;
    private final Thread dispatcher;

    public S7ProtocolEventLogic(BlockingQueue<S7Event> eventQueue) {
        BlockingQueue<S7Event> dispatchQueue = new ArrayBlockingQueue<>(1024);
        runProcessor = new ObjectProcessor(eventQueue, dispatchQueue);
        runDispatcher = new EventDispatcher(dispatchQueue);
        processor = new BasicThreadFactory.Builder()
            .namingPattern("plc4x-evt-processor-thread-%d")
            .daemon(true)
            .priority(Thread.MAX_PRIORITY)
            .build().newThread(runProcessor);
        dispatcher = new BasicThreadFactory.Builder()
            .namingPattern("plc4x-evt-dispatcher-thread-%d")
            .daemon(true)
            .priority(Thread.MAX_PRIORITY)
            .build().newThread(runDispatcher);
    }

    public void start() {
        processor.start();
        dispatcher.start();
    }

    public void stop() {
        runProcessor.doShutdown();
        runDispatcher.doShutdown();
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
        Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers;
        S7PlcSubscriptionHandle handle = (S7PlcSubscriptionHandle) handles.toArray()[0];
        if (!mapIndex.containsKey(handle.getEventType())) {
            mapConsumers = new HashMap<>();
            mapIndex.put(handle.getEventType(), mapConsumers);
        }
        mapConsumers = mapIndex.get(handle.getEventType());
        //TODO: Check the implementation of "DefaultPlcConsumerRegistration". List<> vs Collection<>
        DefaultPlcConsumerRegistration registration = new DefaultPlcConsumerRegistration(this,
            consumer, handles.toArray(new PlcSubscriptionHandle[0]));
        mapConsumers.put(registration, consumer);
        return registration;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        S7PlcSubscriptionHandle handle = (S7PlcSubscriptionHandle) registration.getSubscriptionHandles().get(0);
        Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(handle.getEventType());
        mapConsumers.remove(registration);
    }

    
    //TODO: Replace with disruptor
    private class ObjectProcessor implements Runnable {

        private final BlockingQueue<S7Event> eventQueue;
        private final BlockingQueue<S7Event> dispatchQueue;
        private boolean shutdown = false;
        public ObjectProcessor(BlockingQueue<S7Event> eventQueue, BlockingQueue<S7Event> dispatchQueue) {
            this.eventQueue = eventQueue;
            this.dispatchQueue = dispatchQueue;
        }

        @Override
        public void run() {
            while (!shutdown) {
                try {
                    S7Event s7Event = eventQueue.poll(DEFAULT_DELAY, TimeUnit.MILLISECONDS);
                    if ((s7Event != null) && (dispatchQueue.remainingCapacity() > 1)) {
                        if (s7Event instanceof S7ModeEvent) {
                            dispatchQueue.add(s7Event);
                        } else if (s7Event instanceof S7UserEvent) {
                            dispatchQueue.add(s7Event);                            
                        } else if (s7Event instanceof S7SysEvent) {                            
                            dispatchQueue.add(s7Event);                            
                        } else if (s7Event instanceof S7CyclicEvent) {
                            dispatchQueue.add(s7Event);
                        } else {
                            dispatchQueue.add(s7Event);
                        }
                    }
                } catch (InterruptedException ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
            logger.info("ObjectProcessor Bye!");
        }

        public void doShutdown() {
            shutdown = true;
        }
    }

    //TODO: Replace with disruptor    
    private class EventDispatcher implements Runnable {
        private final BlockingQueue<S7Event> dispatchQueue;
        private boolean shutdown = false;
        private S7Event cycDelayedObject = null;

        public EventDispatcher(BlockingQueue<S7Event> dispatchQueue) {
            this.dispatchQueue = dispatchQueue;
        }

        @Override
        public void run() {
            try {
                Thread.sleep(500);
            } catch (InterruptedException ex) {
                logger.warn(ex.toString());
            }
            while (!shutdown) {
                try {
                    S7Event s7Event = dispatchQueue.poll(DEFAULT_DELAY, TimeUnit.MILLISECONDS);
                    if ((s7Event == null) && (cycDelayedObject != null)) {
                        s7Event = cycDelayedObject;
                        cycDelayedObject = null;
                    }
                    if (s7Event != null) {
                        if (s7Event instanceof S7ModeEvent) {
                            S7ModeEvent modeEvent = (S7ModeEvent) s7Event;
                            if (mapIndex.containsKey(EventType.MODE)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.MODE);
                                mapConsumers.forEach((x, y) -> y.accept(modeEvent));
                            }
                        } else if (s7Event instanceof S7UserEvent) {
                            S7UserEvent userEvent = (S7UserEvent) s7Event;
                            if (mapIndex.containsKey(EventType.USR)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.USR);
                                mapConsumers.forEach((x, y) -> y.accept(userEvent));
                            }
                        } else if (s7Event instanceof S7SysEvent) {
                            S7SysEvent sysEvent = (S7SysEvent) s7Event;
                            if (mapIndex.containsKey(EventType.SYS)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.SYS);
                                mapConsumers.forEach((x, y) -> y.accept(sysEvent));
                            }
                        } else if (s7Event instanceof S7AlarmEvent) {
                            S7AlarmEvent alarmEvent = (S7AlarmEvent) s7Event;
                            if (mapIndex.containsKey(EventType.ALM)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.ALM);
                                mapConsumers.forEach((x, y) -> y.accept(alarmEvent));
                            }
                        } else if (s7Event instanceof S7CyclicEvent) {
                            S7CyclicEvent cyclicEvent = (S7CyclicEvent) s7Event;
                            if (mapIndex.containsKey(EventType.CYC)) {
                                Map<PlcConsumerRegistration, Consumer<PlcSubscriptionEvent>> mapConsumers = mapIndex.get(EventType.CYC);

                                if (cycDelayedObject != null) {
                                    mapConsumers.forEach((x, y) -> y.accept(cycDelayedObject));
                                    cycDelayedObject = null;
                                }
                                if (mapConsumers.isEmpty()) cycDelayedObject = s7Event;
                                
                                mapConsumers.forEach((x, y) -> {
                                    S7PlcSubscriptionHandle sh = (S7PlcSubscriptionHandle) x.getSubscriptionHandles().get(0);
                                    Short id = Short.parseShort(sh.getEventId());
                                    if (cyclicEvent.getMap().get("JOBID") == id) {
                                        y.accept(cyclicEvent);
                                    }
                                });
                            } else {
                                cycDelayedObject = s7Event;
                            }
                        }
                    }
                } catch (Exception ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
            logger.info("EventDispatcher Bye!");
        }

        public void doShutdown() {
            shutdown = true;
        }

    }

}
