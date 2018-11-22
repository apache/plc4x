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
package org.apache.plc4x.java.mock;

import org.apache.commons.lang3.tuple.Pair;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.authentication.PlcAuthentication;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedOperationException;
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.base.messages.DefaultPlcReadRequest;
import org.apache.plc4x.java.base.messages.DefaultPlcReadResponse;
import org.apache.plc4x.java.base.messages.PlcReader;
import org.apache.plc4x.java.base.messages.items.BaseDefaultFieldItem;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Collectors;

public class PlcMockConnection implements PlcConnection, PlcReader {

    private static final Logger logger = LoggerFactory.getLogger(PlcMockConnection.class);

    private final String name;
    private final PlcAuthentication authentication;

    private boolean isConnected = false;
    private MockDevice device;

    PlcMockConnection(String name, PlcAuthentication authentication) {
        this.name = name;
        this.authentication = authentication;
    }

    public MockDevice getDevice() {
        return device;
    }

    public void setDevice(MockDevice device) {
        logger.info("Set Mock Devie on Mock Connection " + this + " with device " + device);
        this.device = device;
    }

    @Override
    public void connect() {
        // do nothing
    }

    @Override
    public boolean isConnected() {
        return true;
    }

    @Override
    public void close() {
        logger.info("Closing MockConnection with device " + device);
    }

    @Override
    public PlcConnectionMetadata getMetadata() {
        return new PlcConnectionMetadata() {
            @Override
            public boolean canRead() {
                return true;
            }

            @Override
            public boolean canWrite() {
                return false;
            }

            @Override
            public boolean canSubscribe() {
                return false;
            }
        };
    }

    @Override
    public PlcReadRequest.Builder readRequestBuilder() {
        return new DefaultPlcReadRequest.Builder(this, new MockFieldHandler());
    }

    @Override
    public CompletableFuture<PlcReadResponse> read(PlcReadRequest readRequest) {
        return CompletableFuture.supplyAsync(new Supplier<PlcReadResponse>() {

            @Override
            public PlcReadResponse get() {
                logger.debug("Sending read request to MockDevice");
                Map<String, Pair<PlcResponseCode, BaseDefaultFieldItem>> response = readRequest.getFieldNames().stream()
                    .collect(Collectors.toMap(
                        Function.identity(),
                        name -> device.read(((MockField) readRequest.getField(name)).getFieldQuery())
                        )
                    );
                return new DefaultPlcReadResponse((DefaultPlcReadRequest)readRequest, response);
            }
        });
    }

    @Override
    public PlcWriteRequest.Builder writeRequestBuilder() {
        throw new PlcUnsupportedOperationException("Write not supported by Mock Driver");
    }

    @Override
    public PlcSubscriptionRequest.Builder subscriptionRequestBuilder() {
        throw new PlcUnsupportedOperationException("Subscription not supported by Mock Driver");
    }

    @Override
    public PlcUnsubscriptionRequest.Builder unsubscriptionRequestBuilder() {
        throw new PlcUnsupportedOperationException("Subscription not supported by Mock Driver");
    }

    public PlcAuthentication getAuthentication() {
        return authentication;
    }
}
