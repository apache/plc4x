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
package org.apache.plc4x.nifi.subscription;

import java.time.Duration;
import java.util.Map;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import org.apache.nifi.logging.ComponentLog;
import org.apache.nifi.processor.exception.ProcessException;
import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcConnectionManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcProtocolException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;

public class Plc4xListenerDispatcher implements Runnable {

    private PlcConnectionManager connectionManager;
    private Plc4xSubscriptionType subscriptionType;
    private Long cyclingPollingInterval;
    private ComponentLog logger;
    private boolean running = false;
    private BlockingQueue<PlcSubscriptionEvent> events;
    private PlcConnection connection;
    private Long timeout;
    private BlockingQueue<PlcSubscriptionEvent> queuedEvents;

    public boolean isRunning() {
        return running;
    }

    public Plc4xListenerDispatcher(Long timeout, Plc4xSubscriptionType subscriptionType, Long cyclingPollingInterval, ComponentLog logger, final BlockingQueue<PlcSubscriptionEvent> events) {
        this.timeout = timeout;
        this.subscriptionType = subscriptionType;
        this.cyclingPollingInterval = cyclingPollingInterval;
        this.logger = logger;
        this.events = events;
        this.queuedEvents = new LinkedBlockingQueue<>();
        this.connectionManager = new DefaultPlcDriverManager();
    }

    /**
     * Opens the dispatcher
     *
     * @param plcConnectionString the connection string for the to the plc
     * @param tags                a map of tag identifier and tag address the
     *                            dispatcher will try subscribing
     * @throws Exception
     * @throws PlcConnectionException
     */
    public void open(String plcConnectionString, Map<String, String> tags) throws PlcConnectionException, Exception {
        connection = connectionManager.getConnection(plcConnectionString);

        if (!connection.getMetadata().canSubscribe()) {
            throw new PlcProtocolException("This connection does not support subscription");
        }

        PlcSubscriptionRequest.Builder builder = connection.subscriptionRequestBuilder();

        for (Map.Entry<String, String> entry : tags.entrySet()) {
            switch (subscriptionType) {
                case CHANGE:
                    builder.addChangeOfStateTagAddress(entry.getKey(), entry.getValue());
                    break;
                case CYCLIC:
                    builder.addCyclicTagAddress(entry.getKey(), entry.getValue(), Duration.ofMillis(cyclingPollingInterval));
                    break;
                case EVENT:
                    builder.addEventTagAddress(entry.getKey(), entry.getValue());
            }
        }
        PlcSubscriptionRequest subscriptionRequest = builder.build();
        PlcSubscriptionResponse subscriptionResponse;
        try {
            subscriptionResponse = subscriptionRequest.execute().get(timeout, TimeUnit.MILLISECONDS);
            
        } catch (InterruptedException e) {
            logger.error("InterruptedException reading the data from PLC", e);
            throw e;
        } catch (TimeoutException e) {
            logger.error("Timeout connection to PLC", e);
            throw e;
        } catch (Exception e) {
            logger.error("Exception reading the data from PLC", e);
            throw (e instanceof ProcessException) ? (ProcessException) e : new ProcessException(e);
        }

        for (PlcSubscriptionHandle handle : subscriptionResponse.getSubscriptionHandles()) {
            handle.register(plcSubscriptionEvent -> {
                queuedEvents.offer(plcSubscriptionEvent);
            });
        }

        running = true;
    }

    /**
     * Closes all listeners and stops all handler threads.
     */
    public void close() {
        running = false;
        try {
            connection.close();
        } catch (Exception e) {
            logger.debug(e.getMessage());
        }
    }

    /**
     * Runs the thread. If no subscription events are received in less than timeout milliseconds the dispatcher is closed.
     */
    @Override
    public void run() {
        while (running) {
            try {
                // If there is a new event before timeout save it, else reopen the connection
                PlcSubscriptionEvent event = queuedEvents.poll(timeout, TimeUnit.MILLISECONDS);
                if (event != null){
                    events.put(event);
                } else {
                    close();
                }
            } catch (InterruptedException e){ 
                close();
            }
        }
    }
}
