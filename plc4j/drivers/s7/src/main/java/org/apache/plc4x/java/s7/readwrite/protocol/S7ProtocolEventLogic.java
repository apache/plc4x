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
import org.apache.plc4x.java.s7.readwrite.types.EventType;
import org.apache.plc4x.java.spi.messages.PlcSubscriber;
import org.apache.plc4x.java.spi.model.DefaultPlcConsumerRegistration;

/**
 *
 * @author cgarcia
 */
public class S7ProtocolEventLogic implements PlcSubscriber {

    
    private final BlockingQueue eventqueue;
    
    private Map<PlcSubscriptionHandle, Map<PlcConsumerRegistration, Consumer>> mapIndex = new HashMap(); 

    
    private Thread processor;
    private Thread dispacher;

    public S7ProtocolEventLogic(BlockingQueue eventqueue) {
        this.eventqueue = eventqueue;
    }
    
    public void start() {
        processor = new Thread(new ObjectProcessor(eventqueue));
        dispacher = new Thread();
        
        processor.run();
        dispacher.run();
        
    }
    
    public void stop(){
        processor.interrupt();
        dispacher.interrupt();
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
        if (!mapIndex.containsKey((PlcSubscriptionHandle) handles.toArray()[0])) {
            mapConsumers = new HashMap();
            mapIndex.put((PlcSubscriptionHandle) handles.toArray()[0], mapConsumers);
        }
        mapConsumers = mapIndex.get((PlcSubscriptionHandle) handles.toArray()[0]);
        DefaultPlcConsumerRegistration registro = new DefaultPlcConsumerRegistration(this,
                consumer, (PlcSubscriptionHandle[]) handles.toArray());
        mapConsumers.put(registro, consumer);
        return registro;
        
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        Map<PlcConsumerRegistration, Consumer> mapConsumers = mapIndex.get(registration.getSubscriptionHandles().get(0)); 
        mapConsumers.remove(registration);
    }
    
    private class ObjectProcessor implements Runnable {

        private final BlockingQueue eventqueue;
        private int delay = 1000;
        
        public ObjectProcessor(BlockingQueue eventqueue) {
            this.eventqueue = eventqueue;
        }
        
        
        @Override
        public void run() {
            while(true){
                try {
                    System.out.println("ObjectProcessor  Paso por aqui...");
                    Object obj = eventqueue.poll(delay, TimeUnit.MILLISECONDS);                    
                } catch (InterruptedException ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
        }        
    }    
    
    private class EventDispacher implements Runnable {

        @Override
        public void run() {
            while(true){
                try {
                    Thread.sleep(5000);
                } catch (InterruptedException ex) {
                    Logger.getLogger(S7ProtocolEventLogic.class.getName()).log(Level.SEVERE, null, ex);
                }
            }
        }        
    }
    
    
}
