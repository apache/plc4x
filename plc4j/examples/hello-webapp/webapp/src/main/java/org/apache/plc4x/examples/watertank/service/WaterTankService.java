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
package org.apache.plc4x.examples.watertank.service;

import flex.messaging.MessageBroker;
import flex.messaging.MessageDestination;
import flex.messaging.config.NetworkSettings;
import flex.messaging.messages.AsyncMessage;
import flex.messaging.services.MessageService;
import flex.messaging.util.UUIDUtils;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.messages.PlcSubscriptionRequest;
import org.apache.plc4x.java.api.messages.PlcSubscriptionResponse;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.flex.remoting.RemotingDestination;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;
import javax.annotation.PreDestroy;
import java.security.SecureRandom;
import java.time.Instant;
import java.util.Calendar;
import java.util.HashMap;
import java.util.Map;
import java.util.Random;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

@Service("waterTankService")
@RemotingDestination
public class WaterTankService {

    public static final String WATER_LEVEL = "waterLevel";
    @Autowired
    private MessageBroker messageBroker;

    private PlcConnection connection;

    @Value("${plc4x.watertank.connection-string}")
    private String connectionString;

    @Value("${plc4x.watertank.address-string.waterlevel}")
    private String addressStringWaterLevel;

    @Value("${plc4x.watertank.simulation}")
    private boolean simulation;

    private final Random random = new SecureRandom();

    @PostConstruct
    protected void connectAndSubscribe() {
        // Create the BlazeDS destinations the clients can subscribe to.
        createDestinations();

        // Connect to the remote and subscribe to the values.
        if (simulation) {
            connectToFakeDevice();
        } else {
            connectToDevice();
        }
    }

    @PreDestroy
    protected void disconnect() {
        if ((connection != null) && connection.isConnected()) {
            try {
                connection.close();
            } catch (Exception e) {
                throw new PlcRuntimeException("Error stopping");
            }
        }
    }

    protected void connectToDevice() {
        try {
            // Connect to the device
            connection = new PlcDriverManager().getConnection(connectionString);

            // Check if subscriptions are supported by this connection.
            if (!connection.getMetadata().canSubscribe()) {
                throw new PlcRuntimeException("This driver doesn't support subscribing");
            }

            // Prepare a subscription request.
            final PlcSubscriptionRequest subscriptionRequest =
                connection.subscriptionRequestBuilder().addChangeOfStateField(
                    WATER_LEVEL, addressStringWaterLevel).build();

            // Execute the request.
            PlcSubscriptionResponse syncResponse = subscriptionRequest.execute().get();

            // Attach handlers for the incoming data.
            for (String subscriptionName : syncResponse.getFieldNames()) {
                final PlcSubscriptionHandle subscriptionHandle = syncResponse.getSubscriptionHandle(subscriptionName);
                subscriptionHandle.register(new WaterLevelHandler());
            }
        } catch (PlcConnectionException e) {
            throw new PlcRuntimeException("Error connecting to device", e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new PlcRuntimeException("Error subscribing for data", e);
        } catch (ExecutionException e) {
            throw new PlcRuntimeException("Error subscribing for data", e);
        }
    }

    protected void connectToFakeDevice() {
        WaterLevelHandler handler = new WaterLevelHandler();
        Thread thread = new Thread(() -> {
            try {
                while (true) {
                    TimeUnit.MILLISECONDS.sleep(100);

                    short value = (short) random.nextInt(1024);
                    Map<String, ResponseItem<PlcValue>> values = new HashMap<>();
                    values.put(WATER_LEVEL, new ResponseItem<>(PlcResponseCode.OK, new PlcINT(value)));
                    DefaultPlcSubscriptionEvent event = new DefaultPlcSubscriptionEvent(Instant.now(), values);

                    handler.accept(event);
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });
        thread.start();
    }

    protected void createDestinations() {
        // Create a new dynamic-destination for this session.
        final MessageService service = (MessageService) messageBroker.getService("message-service");

        // Create our new destination which we want to post to.
        final MessageDestination destination = (MessageDestination) service.createDestination(WATER_LEVEL);

        // Attach the channels to it.
        destination.setChannels(messageBroker.getChannelIds());

        // Make sure subscriptions to this channel are timed out after 5 minutes.
        NetworkSettings destinationNetworkSettings = new NetworkSettings();
        destinationNetworkSettings.setSubscriptionTimeoutMinutes(5);
        destination.setNetworkSettings(destinationNetworkSettings);

        // Start the destination.
        destination.start();
    }

    private class WaterLevelHandler implements Consumer<PlcSubscriptionEvent> {
        @Override
        public void accept(PlcSubscriptionEvent plcSubscriptionEvent) {
            if (!messageBroker.isStarted()) {
                return;
            }
            // Get the data from the event.
            final Short waterLevel = plcSubscriptionEvent.getShort(WATER_LEVEL);
            if (waterLevel == null) {
                return;
            }

            // Create a new message that will be broadcasted to the clients.
            final AsyncMessage message = new AsyncMessage();
            message.setMessageId(UUIDUtils.createUUID(false));
            message.setTimestamp(Calendar.getInstance().getTimeInMillis());
            message.setDestination(WATER_LEVEL);
            message.setBody(waterLevel);

            // Actually send the message.
            messageBroker.routeMessageToService(message, null);
        }
    }

}
