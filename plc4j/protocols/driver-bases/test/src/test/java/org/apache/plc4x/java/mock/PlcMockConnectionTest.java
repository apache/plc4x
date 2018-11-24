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
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.metadata.PlcConnectionMetadata;
import org.apache.plc4x.java.api.model.PlcSubscriptionHandle;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collections;
import java.util.concurrent.CompletableFuture;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class PlcMockConnectionTest implements WithAssertions {

    PlcMockConnection SUT = new PlcMockConnection(null);

    @Mock
    MockDevice mockDevice;

    @BeforeEach
    void setUp() {
        SUT.setDevice(mockDevice);
    }

    @Test
    void getDevice() {
        assertThat(SUT.getDevice()).isNotNull();
    }

    @Test
    void setDevice() {
        SUT.setDevice(null);
    }

    @Test
    void connect() {
        SUT.connect();
    }

    @Test
    void isConnected() {
        SUT.setDevice(mockDevice);
        assertThat(SUT.isConnected()).isTrue();
        SUT.setDevice(null);
        assertThat(SUT.isConnected()).isFalse();
    }

    @Test
    void close() {
        SUT.close();
    }

    @Test
    void getMetadata() {
        PlcConnectionMetadata metadata = SUT.getMetadata();
        assertThat(metadata).isNotNull();
        assertThat(metadata.canRead()).isTrue();
        assertThat(metadata.canWrite()).isTrue();
        assertThat(metadata.canSubscribe()).isTrue();
    }

    @Test
    void readRequestBuilder() {
        assertThat(SUT.readRequestBuilder()).isNotNull();
    }

    @Test
    void read() {
        PlcReadRequest plcReadRequest = SUT.readRequestBuilder().build();
        assertThat(SUT.read(plcReadRequest)).isNotNull();
    }

    @Test
    void write() throws Exception {
        when(mockDevice.write(any(), any())).thenReturn(PlcResponseCode.OK);

        PlcWriteRequest plcWriteRequest = SUT.writeRequestBuilder()
            .addItem("asd", "asd", "asd")
            .build();
        CompletableFuture<PlcWriteResponse> write = SUT.write(plcWriteRequest);
        assertThat(write).isNotNull();
        PlcWriteResponse plcWriteResponse = write.get();
        assertThat(plcWriteResponse).isNotNull();
    }

    @Test
    void subscribe() throws Exception {
        when(mockDevice.subscribe(any())).thenReturn(Pair.of(PlcResponseCode.OK, mock(PlcSubscriptionHandle.class)));

        PlcSubscriptionRequest plcSubscriptionRequest = SUT.subscriptionRequestBuilder()
            .addChangeOfStateField("asd", "asd")
            .build();
        CompletableFuture<PlcSubscriptionResponse> subscribe = SUT.subscribe(plcSubscriptionRequest);
        assertThat(subscribe).isNotNull();
        PlcSubscriptionResponse plcSubscriptionResponse = subscribe.get();
        assertThat(plcSubscriptionResponse).isNotNull();
    }

    @Test
    void unsubscribe() {
        PlcUnsubscriptionRequest unsubscriptionRequest = SUT.unsubscriptionRequestBuilder().build();
        assertThat(SUT.unsubscribe(unsubscriptionRequest)).isNotNull();
    }

    @Test
    void register() {
        SUT.register(plcSubscriptionEvent -> {
        }, Collections.emptyList());
    }

    @Test
    void unregister() {
        SUT.unregister(null);
    }

    @Test
    void writeRequestBuilder() {
        assertThat(SUT.writeRequestBuilder()).isNotNull();
    }

    @Test
    void subscriptionRequestBuilder() {
        assertThat(SUT.subscriptionRequestBuilder()).isNotNull();
    }

    @Test
    void unsubscriptionRequestBuilder() {
        assertThat(SUT.unsubscriptionRequestBuilder()).isNotNull();
    }

    @Test
    void getAuthentication() {
        assertThat(SUT.getAuthentication()).isNull();
    }
}