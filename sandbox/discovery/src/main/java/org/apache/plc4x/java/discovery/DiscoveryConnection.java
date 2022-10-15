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
package org.apache.plc4x.java.discovery;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.PlcDriver;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcConsumerRegistration;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.spi.discovery.ActiveDiscovery;
import org.apache.plc4x.java.spi.discovery.BroadcastDiscovery;
import org.apache.plc4x.java.spi.discovery.SupportsDiscovery;
import org.apache.plc4x.java.spi.discovery.PassiveDiscovery;
import org.apache.plc4x.java.spi.messages.*;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.ServiceLoader;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

public class DiscoveryConnection implements PlcConnection, PlcSubscriber, PlcBrowser {

    private boolean connected = false;
    private Map<String, ActiveDiscovery> activeDiscovery;
    private Map<String, PassiveDiscovery> passiveDiscovery;
    private Map<String, BroadcastDiscovery> broadcastDiscovery;

    @Override
    public void connect() {
        if(connected) {
            return;
        }

        activeDiscovery = new HashMap<>();
        passiveDiscovery = new HashMap<>();
        broadcastDiscovery = new HashMap<>();

        // Access all PlcDriver implementations in this VM.
        ServiceLoader<PlcDriver> plcDriverLoader = ServiceLoader.load(
            PlcDriver.class, Thread.currentThread().getContextClassLoader());

        // Iterate over all of them and check which types of discovery they support.
        for (PlcDriver driver : plcDriverLoader) {
            // Check if this driver generally supports discovery at all.
            if(driver instanceof SupportsDiscovery) {
                if(driver instanceof ActiveDiscovery) {
                    activeDiscovery.put(driver.getProtocolCode(), (ActiveDiscovery) driver);
                }
                if(driver instanceof PassiveDiscovery) {
                    passiveDiscovery.put(driver.getProtocolCode(), (PassiveDiscovery) driver);
                }
                if(driver instanceof BroadcastDiscovery) {
                    broadcastDiscovery.put(driver.getProtocolCode(), (BroadcastDiscovery) driver);
                }
            }
        }

        connected = true;
    }

    @Override
    public boolean isConnected() {
        return connected;
    }

    @Override
    public void close() {
        connected = false;
        activeDiscovery = null;
        passiveDiscovery = null;
        broadcastDiscovery = null;
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return new PlcConnectionMetadata() {
            @Override
            public boolean canRead() {
                return false;
            }

            @Override
            public boolean canWrite() {
                return false;
            }

            @Override
            public boolean canSubscribe() {
                return true;
            }

            @Override
            public boolean canBrowse() {
                return true;
            }
        };
    }

    @Override
    public CompletableFuture<Void> ping() {
        CompletableFuture<Void> future = new CompletableFuture<>();
        future.completeExceptionally(new PlcUnsupportedOperationException("The connection does not support pinging"));
        return future;
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        throw new PlcUnsupportedOperationException("The connection does not support reading");
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        throw new PlcUnsupportedOperationException("The connection does not support writing");
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        return new DefaultPlcSubscriptionRequest.Builder(this, new DiscoveryFieldHandler());
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        return new DefaultPlcUnsubscriptionRequest.Builder(this);
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        return new DefaultPlcBrowseRequest.Builder(this);
    }

    @Override
    public CompletableFuture<PlcSubscriptionResponse> subscribe(PlcSubscriptionRequest subscriptionRequest) {
        // TODO: Implement ...
        return null;
    }

    @Override
    public CompletableFuture<PlcUnsubscriptionResponse> unsubscribe(PlcUnsubscriptionRequest unsubscriptionRequest) {
        // TODO: Implement ...
        return null;
    }

    @Override
    public PlcConsumerRegistration register(Consumer<PlcSubscriptionEvent> consumer, Collection<PlcSubscriptionHandle> handles) {
        // TODO: Implement ...
        return null;
    }

    @Override
    public void unregister(PlcConsumerRegistration registration) {
        // TODO: Implement ...
    }

    @Override
    public CompletableFuture<PlcBrowseResponse> browse(PlcBrowseRequest browseRequest) {
        // TODO: Implement ...
        return null;
    }

}
