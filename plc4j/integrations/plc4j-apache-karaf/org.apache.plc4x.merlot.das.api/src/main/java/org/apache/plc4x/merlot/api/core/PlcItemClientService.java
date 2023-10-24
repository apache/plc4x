/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.plc4x.merlot.api.core;

import com.lmax.disruptor.RingBuffer;
import com.lmax.disruptor.TimeoutException;
import com.lmax.disruptor.dsl.Disruptor;
import com.lmax.disruptor.util.DaemonThreadFactory;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.merlot.api.PlcItemClient;
import org.osgi.framework.BundleContext;
import org.slf4j.LoggerFactory;


public class PlcItemClientService {
    private static final org.slf4j.Logger LOGGER = LoggerFactory.getLogger(PlcItemClientService .class);
    private final int BUFFER_SIZE = 1024;
    
    private final BundleContext bc;
    private ConcurrentHashMap<UUID, List<PlcItemClient>> item_clients;
    private Disruptor<PlcItemEvent> disruptor;
    private RingBuffer<PlcItemEvent> ringbuffer;    

    public PlcItemClientService(BundleContext bc) {
        this.bc = bc;
    }
    
    public void init(){
        // 
        item_clients = new ConcurrentHashMap<>();
        
        disruptor = new Disruptor<>(PlcItemEvent::new, BUFFER_SIZE, DaemonThreadFactory.INSTANCE); 
        
        disruptor.handleEventsWith((event, sequence, endOfBatch) -> {
            if (item_clients.containsKey(event.uid)) {
                item_clients.get(event.uid).forEach(c -> c.execute(event.plcvalue));  
                System.out.println("PlcValue: " + event.plcvalue);
            }
        });
        disruptor.start();
        ringbuffer = disruptor.getRingBuffer();
    }
    
    public void destroy() {
        try {
            disruptor.shutdown(1, TimeUnit.SECONDS);
        } catch (TimeoutException ex) {
            LOGGER.info(ex.getMessage());
        }
        item_clients.clear();
    }
    
   //TODO: Check if the Item really exists
    public boolean subcription(PlcItemClient client, UUID item_uid){
        if (item_clients.containsKey(item_uid)) {
            final List<PlcItemClient> clients = item_clients.get(item_uid);
            if (!clients.contains(client)) {
                clients.add(client);
                return true;
            }
        } else {
            final List<PlcItemClient> clients = new ArrayList<PlcItemClient>();
            clients.add(client);
            item_clients.put(item_uid, clients);
            return true;
        }
        return false;
    }
    
    public void putItemEvent(UUID item_uid, PlcValue  plcvalue){
        if (item_clients.containsKey(item_uid)) {        
            ringbuffer.publishEvent((event, sequence, buffer) -> {
                    event.setUid(item_uid);
                    event.setPlcValue(plcvalue);
                });
        }
    }
        
    public class PlcItemEvent  {
        private UUID uid;

        private int event;

        private PlcValue  plcvalue;

        public UUID getUid() {
            return uid;
        }

        public void setUid(UUID uid) {
            this.uid = uid;
        }

        public int getEvent() {
            return event;
        }

        public void setEvent(int event) {
            this.event = event;
        }

        public PlcValue getPlcValue() {
            return plcvalue;
        }

        public void setPlcValue(PlcValue plcvalue) {
            this.plcvalue = plcvalue;
        }                
            
    }
    
}
