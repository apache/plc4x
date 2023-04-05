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
package org.apache.plc4x.java.utils.cache;

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcQuery;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.model.PlcSubscriptionTag;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;

import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.function.Consumer;

public class LeasedPlcConnection implements PlcConnection {

    private final ConnectionContainer connectionContainer;
    private PlcConnection connection;
    private boolean invalidateConnection;
    private Timer usageTimer;

    public LeasedPlcConnection(ConnectionContainer connectionContainer, PlcConnection connection, Duration maxUseTime) {
        this.connectionContainer = connectionContainer;
        this.connection = connection;
        this.invalidateConnection = false;
        this.usageTimer = new Timer();
        this.usageTimer.schedule(new TimerTask() {
            @Override
            public void run() {
                close();
            }
        }, Date.from(LocalDateTime.now().plusNanos(maxUseTime.toNanos()).atZone(ZoneId.systemDefault()).toInstant()));
    }

    @Override
    public synchronized void close() {
        // Cancel automatically timing out.
        usageTimer.cancel();

        // Make the connection unusable.
        connection = null;

        // Tell the connection container that the connection is free to be reused.
        connectionContainer.returnConnection(this, invalidateConnection);
    }

    @Override
    public void connect() throws PlcConnectionException {
        throw new PlcConnectionException("Error connecting leased connection");
    }

    @Override
    public boolean isConnected() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.isConnected();
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.getMetadata();
    }

    @Override
    public CompletableFuture<Void> ping() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        return connection.ping();
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        final PlcReadRequest.Builder innerBuilder = connection.readRequestBuilder();
        return new PlcReadRequest.Builder() {
            @Override
            public PlcReadRequest build() {
                final PlcReadRequest innerPlcReadRequest = innerBuilder.build();
                return new PlcReadRequest(){
                    @Override
                    public CompletableFuture<? extends PlcReadResponse> execute() {
                        CompletableFuture<? extends PlcReadResponse> future = innerPlcReadRequest.execute();
                        final CompletableFuture<PlcReadResponse> responseFuture = new CompletableFuture<>();
                        future.handle((plcReadResponse, throwable) -> {
                            if (throwable == null) {
                                responseFuture.complete(plcReadResponse);
                            } else {
                                // Mark the connection as invalid.
                                invalidateConnection = true;
                                responseFuture.completeExceptionally(throwable);
                            }
                            return null;
                        });
                        return responseFuture;
                    }

                    @Override
                    public int getNumberOfTags() {
                        return innerPlcReadRequest.getNumberOfTags();
                    }

                    @Override
                    public LinkedHashSet<String> getTagNames() {
                        return innerPlcReadRequest.getTagNames();
                    }

                    @Override
                    public PlcTag getTag(String name) {
                        return innerPlcReadRequest.getTag(name);
                    }

                    @Override
                    public List<PlcTag> getTags() {
                        return innerPlcReadRequest.getTags();
                    }
                };
            }

            @Override
            public PlcReadRequest.Builder addTagAddress(String name, String tagAddress) {
                return innerBuilder.addTagAddress(name, tagAddress);
            }

            @Override
            public PlcReadRequest.Builder addTag(String name, PlcTag tag) {
                return innerBuilder.addTag(name, tag);
            }
        };
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        final PlcWriteRequest.Builder innerBuilder = connection.writeRequestBuilder();
        return new PlcWriteRequest.Builder() {
            @Override
            public PlcWriteRequest build() {
                PlcWriteRequest innerPlcWriteRequest = innerBuilder.build();
                return new PlcWriteRequest() {
                    @Override
                    public CompletableFuture<? extends PlcWriteResponse> execute() {
                        CompletableFuture<? extends PlcWriteResponse> future = innerPlcWriteRequest.execute();
                        final CompletableFuture<PlcWriteResponse> responseFuture = new CompletableFuture<>();
                        future.handle((plcWriteResponse, throwable)->{
                            if (throwable == null) {
                                responseFuture.complete(plcWriteResponse);
                            } else {
                                // Mark the connection as invalid.
                                invalidateConnection = true;
                                responseFuture.completeExceptionally(throwable);
                            }
                            return null;
                        });
                        return responseFuture;
                    }

                    @Override
                    public int getNumberOfValues(String name) {
                        return innerPlcWriteRequest.getNumberOfValues(name);
                    }

                    @Override
                    public PlcValue getPlcValue(String name) {
                        return innerPlcWriteRequest.getPlcValue(name);
                    }

                    @Override
                    public int getNumberOfTags() {
                        return innerPlcWriteRequest.getNumberOfTags();
                    }

                    @Override
                    public LinkedHashSet<String> getTagNames() {
                        return innerPlcWriteRequest.getTagNames();
                    }

                    @Override
                    public PlcTag getTag(String name) {
                        return innerPlcWriteRequest.getTag(name);
                    }

                    @Override
                    public List<PlcTag> getTags() {
                        return innerPlcWriteRequest.getTags();
                    }
                };
            }

            @Override
            public PlcWriteRequest.Builder addTagAddress(String name, String tagAddress, Object... values) {
                return innerBuilder.addTagAddress(name, tagAddress, values);
            }

            @Override
            public PlcWriteRequest.Builder addTag(String name, PlcTag tag, Object... values) {
                return innerBuilder.addTag(name, tag, values);
            }
        };
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        final PlcSubscriptionRequest.Builder innerBuilder = connection.subscriptionRequestBuilder();
        return new PlcSubscriptionRequest.Builder() {
            @Override
            public PlcSubscriptionRequest build() {
                PlcSubscriptionRequest innerPlcSubscriptionRequest = innerBuilder.build();
                return new PlcSubscriptionRequest() {
                    @Override
                    public CompletableFuture<? extends PlcSubscriptionResponse> execute() {
                        CompletableFuture<? extends PlcSubscriptionResponse> future = innerPlcSubscriptionRequest.execute();
                        final CompletableFuture<PlcSubscriptionResponse> responseFuture = new CompletableFuture<>();
                        future.handle((plcSubscriptionResponse, throwable)->{
                            if (throwable == null) {
                                responseFuture.complete(plcSubscriptionResponse);
                            } else {
                                // Mark the connection as invalid.
                                invalidateConnection = true;
                                responseFuture.completeExceptionally(throwable);
                            }
                            return null;
                        });
                        return responseFuture;
                    }

                    @Override
                    public int getNumberOfTags() {
                        return innerPlcSubscriptionRequest.getNumberOfTags();
                    }

                    @Override
                    public LinkedHashSet<String> getTagNames() {
                        return innerPlcSubscriptionRequest.getTagNames();
                    }

                    @Override
                    public PlcSubscriptionTag getTag(String name) {
                        return innerPlcSubscriptionRequest.getTag(name);
                    }

                    @Override
                    public List<PlcSubscriptionTag> getTags() {
                        return innerPlcSubscriptionRequest.getTags();
                    }

                    @Override
                    public Map<String, List<Consumer<PlcSubscriptionEvent>>> getPreRegisteredConsumers() {
                        return innerPlcSubscriptionRequest.getPreRegisteredConsumers();
                    }
                };
            }

            @Override
            public PlcSubscriptionRequest.Builder addCyclicTagAddress(String name, String tagAddress, Duration pollingInterval) {
                return innerBuilder.addCyclicTagAddress(name, tagAddress, pollingInterval);
            }

            @Override
            public PlcSubscriptionRequest.Builder addCyclicTag(String name, PlcTag tag, Duration pollingInterval) {
                return innerBuilder.addCyclicTag(name, tag, pollingInterval);
            }

            @Override
            public PlcSubscriptionRequest.Builder addChangeOfStateTagAddress(String name, String tagAddress) {
                return innerBuilder.addChangeOfStateTagAddress(name, tagAddress);
            }

            @Override
            public PlcSubscriptionRequest.Builder addChangeOfStateTag(String name, PlcTag tag) {
                return innerBuilder.addChangeOfStateTag(name, tag);
            }

            @Override
            public PlcSubscriptionRequest.Builder addEventTagAddress(String name, String tagAddress) {
                return innerBuilder.addEventTagAddress(name, tagAddress);
            }

            @Override
            public PlcSubscriptionRequest.Builder addEventTag(String name, PlcTag tag) {
                return innerBuilder.addEventTag(name, tag);
            }

            @Override
            public PlcSubscriptionRequest.Builder addPreRegisteredConsumer(String name, Consumer<PlcSubscriptionEvent> preRegisteredConsumer) {
                return innerBuilder.addPreRegisteredConsumer(name, preRegisteredConsumer);
            }
        };
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        final PlcUnsubscriptionRequest.Builder innerBuilder = connection.unsubscriptionRequestBuilder();
        return new PlcUnsubscriptionRequest.Builder() {
            @Override
            public PlcUnsubscriptionRequest build() {
                PlcUnsubscriptionRequest innerPlcUnsubscriptionRequest = innerBuilder.build();
                return new PlcUnsubscriptionRequest() {
                    @Override
                    public CompletableFuture<PlcUnsubscriptionResponse> execute() {
                        CompletableFuture<? extends PlcUnsubscriptionResponse> future = innerPlcUnsubscriptionRequest.execute();
                        final CompletableFuture<PlcUnsubscriptionResponse> responseFuture = new CompletableFuture<>();
                        future.handle((plcUnsubscriptionResponse, throwable)->{
                            if (throwable == null) {
                                responseFuture.complete(plcUnsubscriptionResponse);
                            } else {
                                // Mark the connection as invalid.
                                invalidateConnection = true;
                                responseFuture.completeExceptionally(throwable);
                            }
                            return null;
                        });
                        return responseFuture;
                    }

                    @Override
                    public List<PlcSubscriptionHandle> getSubscriptionHandles() {
                        return innerPlcUnsubscriptionRequest.getSubscriptionHandles();
                    }
                };
            }

            @Override
            public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle) {
                return innerBuilder.addHandles(plcSubscriptionHandle);
            }

            @Override
            public PlcUnsubscriptionRequest.Builder addHandles(PlcSubscriptionHandle plcSubscriptionHandle1, PlcSubscriptionHandle... plcSubscriptionHandles) {
                return innerBuilder.addHandles(plcSubscriptionHandle1, plcSubscriptionHandles);
            }

            @Override
            public PlcUnsubscriptionRequest.Builder addHandles(Collection<PlcSubscriptionHandle> plcSubscriptionHandle) {
                return innerBuilder.addHandles(plcSubscriptionHandle);
            }
        };
    }

    @Override
    public PlcBrowseRequest.Builder browseRequestBuilder() {
        if(connection == null) {
            throw new PlcRuntimeException("Error using leased connection after returning it to the cache.");
        }
        final PlcBrowseRequest.Builder innerBuilder = connection.browseRequestBuilder();
        return new PlcBrowseRequest.Builder() {
            @Override
            public PlcBrowseRequest build() {
                PlcBrowseRequest innerPlcBrowseRequest = innerBuilder.build();
                return new PlcBrowseRequest() {
                    @Override
                    public CompletableFuture<? extends PlcBrowseResponse> execute() {
                        CompletableFuture<? extends PlcBrowseResponse> future = innerPlcBrowseRequest.execute();
                        final CompletableFuture<PlcBrowseResponse> responseFuture = new CompletableFuture<>();
                        future.handle((plcBrowseResponse, throwable)->{
                            if (throwable == null) {
                                responseFuture.complete(plcBrowseResponse);
                            } else {
                                // Mark the connection as invalid.
                                invalidateConnection = true;
                                responseFuture.completeExceptionally(throwable);
                            }
                            return null;
                        });
                        return responseFuture;
                    }

                    @Override
                    public CompletableFuture<? extends PlcBrowseResponse> executeWithInterceptor(PlcBrowseRequestInterceptor interceptor) {
                        CompletableFuture<? extends PlcBrowseResponse> future = innerPlcBrowseRequest.executeWithInterceptor(interceptor);
                        final CompletableFuture<PlcBrowseResponse> responseFuture = new CompletableFuture<>();
                        future.handle((plcBrowseResponse, throwable)->{
                            if (throwable == null) {
                                responseFuture.complete(plcBrowseResponse);
                            } else {
                                // Mark the connection as invalid.
                                invalidateConnection = true;
                                responseFuture.completeExceptionally(throwable);
                            }
                            return null;
                        });
                        return responseFuture;
                    }

                    @Override
                    public LinkedHashSet<String> getQueryNames() {
                        return innerPlcBrowseRequest.getQueryNames();
                    }

                    @Override
                    public PlcQuery getQuery(String name) {
                        return innerPlcBrowseRequest.getQuery(name);
                    }
                };
            }

            @Override
            public PlcBrowseRequest.Builder addQuery(String name, String query) {
                return innerBuilder.addQuery(name, query);
            }
        };
    }

}
