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
package org.apache.plc4x.java.s7.readwrite.protocol;

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

import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcUnsubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.s7.events.S7AlarmEvent;
import org.apache.plc4x.java.s7.events.S7CyclicEvent;
import org.apache.plc4x.java.s7.events.S7ModeEvent;
import org.apache.plc4x.java.s7.events.S7SysEvent;
import org.apache.plc4x.java.s7.events.S7UserEvent;
import org.apache.plc4x.java.s7.readwrite.Alarm8MessageQueryType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageObjectQueryType;
import org.apache.plc4x.java.s7.readwrite.AlarmMessageQueryType;
import org.apache.plc4x.java.s7.readwrite.S7Message;
import org.apache.plc4x.java.s7.readwrite.S7MessageUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterModeTransition;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserData;
import org.apache.plc4x.java.s7.readwrite.S7ParameterUserDataItemCPUFunctions;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarm8Query;
import org.apache.plc4x.java.s7.readwrite.S7PayloadAlarmQuery;
import org.apache.plc4x.java.s7.readwrite.S7PayloadDiagnosticMessage;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserData;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItem;
import org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCyclicServicesPush;
import org.apache.plc4x.java.s7.readwrite.types.EventType;
import org.apache.plc4x.java.s7.readwrite.utils.S7PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.generation.MessageIO;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;
import org.slf4j.LoggerFactory;

public class S7ProtocolEventLogic implements PlcSubscriber {
    private final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(S7ProtocolEventLogic.class);

    private final BlockingQueue eventqueue;
    private final BlockingQueue dispachqueue = new ArrayBlockingQueue<>(1024);

    private Map<EventType, Map<PlcConsumerRegistration, Consumer>> mapIndex = new HashMap();
    private Map<EventType, PlcSubscriptionHandle> eventtypehandles = new HashMap();
    ;

    private final Runnable runnProcessor;
    private final Runnable runnDispacher;


    private Thread processor;
    private Thread dispacher;

    public S7ProtocolEventLogic(BlockingQueue eventqueue) {
        this.eventqueue = eventqueue;
        runnProcessor = new ObjectProcessor(eventqueue, dispachqueue);
        runnDispacher = new EventDispacher(dispachqueue);
        processor = new Thread(runnProcessor);
        dispacher = new Thread(runnDispacher);
    }

    public void start() {
        processor.start();
        dispacher.start();
    }

    public void stop() {
        ((ObjectProcessor) runnProcessor).doShutdown();
        ((EventDispacher) runnDispacher).doShutdown();
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
        Map<PlcConsumerRegistration, Consumer> mapConsumers = null;
        S7PlcSubscriptionHandle handle = (S7PlcSubscriptionHandle) handles.toArray()[0];
        if (!mapIndex.containsKey(handle.getEventType())) {
            mapConsumers = new HashMap();
            mapIndex.put(handle.getEventType(), mapConsumers);
        }
        mapConsumers = mapIndex.get(handle.getEventType());
        //TODO: Check the implementation of "DefaultPlcConsumerRegistration". List<> vs Collection<>
        DefaultPlcConsumerRegistration registro = new DefaultPlcConsumerRegistration(this,
            consumer, handles.toArray(new PlcSubscriptionHandle[handles.size()]));
        mapConsumers.put(registro, consumer);
        return registro;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        S7PlcSubscriptionHandle handle = (S7PlcSubscriptionHandle) registration.getSubscriptionHandles().get(0);
        Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(handle.getEventType());
        mapConsumers.remove(registration);
    }

    private class ObjectProcessor implements Runnable {

        private final BlockingQueue eventqueue;
        private final BlockingQueue dispathqueue;
        private boolean shutdown = false;
        private int delay = 5000;

        public ObjectProcessor(BlockingQueue eventqueue, BlockingQueue dispathqueue) {
            this.eventqueue = eventqueue;
            this.dispathqueue = dispathqueue;
        }

        @Override
        public void run() {
            while (!shutdown) {
                try {
                    Object obj = eventqueue.poll(delay, TimeUnit.MILLISECONDS);
                    if (obj == null) {
                        continue;
                    }
                    if (obj instanceof S7ParameterModeTransition) {
                        S7ModeEvent modeEvent = new S7ModeEvent((S7ParameterModeTransition) obj);
                        dispathqueue.add(modeEvent);
                    } else if (obj instanceof S7PayloadDiagnosticMessage) {
                        S7PayloadDiagnosticMessage msg = (S7PayloadDiagnosticMessage) obj;
                        if ((msg.getEventId() >= 0x0A000) & (msg.getEventId() <= 0x0BFFF)) {
                            S7UserEvent userevent = new S7UserEvent(msg);
                            dispathqueue.add(userevent);
                        } else {
                            S7SysEvent sysevent = new S7SysEvent(msg);
                            dispathqueue.add(sysevent);
                        }
                    } else if (obj instanceof S7PayloadAlarmQuery) {
                        AlarmMessageQueryType msgs = ((S7PayloadAlarmQuery) obj).getAlarmMessage();
                        for (AlarmMessageObjectQueryType msg : msgs.getMessageObjects()) {
                            S7AlarmEvent alarmevent = new S7AlarmEvent(msg);
                            dispathqueue.add(alarmevent);
                        }
                    } else if (obj instanceof S7PayloadAlarm8Query) {
                        Alarm8MessageQueryType msgs = ((S7PayloadAlarm8Query) obj).getAlarmMessage();
                        for (AlarmMessageObjectQueryType msg : msgs.getMessageObjects()) {
                            S7AlarmEvent alarmevent = new S7AlarmEvent(msg);
                            dispathqueue.add(alarmevent);
                        }
                    } else if ((obj instanceof S7CyclicEvent)) {
                        dispathqueue.add(obj);
                    } else {
                        S7AlarmEvent alarmevent = new S7AlarmEvent(obj);
                        dispathqueue.add(alarmevent);
                    }
                } catch (InterruptedException ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
            LOGGER.info("ObjectProcessor Bye!");
        }

        public void doShutdown() {
            shutdown = true;
        }
    }

    private class EventDispacher implements Runnable {
        private final BlockingQueue dispachqueue;
        private boolean shutdown = false;
        private int delay = 5000;
        private Object cycDelayedObject = null;

        public EventDispacher(BlockingQueue dispachqueue) {
            this.dispachqueue = dispachqueue;
        }

        @Override
        public void run() {
            while (!shutdown) {
                try {
                    Object obj = dispachqueue.poll(delay, TimeUnit.MILLISECONDS);
                    if (obj != null) {
                        if (obj instanceof S7ModeEvent) {
                            if (mapIndex.containsKey(EventType.MODE)) {
                                Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(EventType.MODE);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7UserEvent) {
                            if (mapIndex.containsKey(EventType.USR)) {
                                Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(EventType.USR);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7SysEvent) {
                            if (mapIndex.containsKey(EventType.SYS)) {
                                Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(EventType.SYS);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7AlarmEvent) {
                            if (mapIndex.containsKey(EventType.ALM)) {
                                Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(EventType.ALM);
                                mapConsumers.forEach((x, y) -> y.accept(obj));
                            }
                        } else if (obj instanceof S7CyclicEvent) {
                            if (mapIndex.containsKey(EventType.CYC)) {
                                Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(EventType.CYC);
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
                } catch (Exception ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
            LOGGER.info("EventDispacher Bye!");
        }

        public void doShutdown() {
            shutdown = true;
        }

    }


}
