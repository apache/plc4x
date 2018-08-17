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
package org.apache.plc4x.java.base.messages;

import org.apache.plc4x.java.api.messages.PlcRequest;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.api.messages.ProtocolMessage;
import org.hamcrest.core.IsEqual;
import org.junit.Before;
import org.junit.Test;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

import static org.junit.Assert.*;
import static org.mockito.Mockito.mock;

@SuppressWarnings("unchecked")
public class PlcRequestContainerTest {

    private PlcRequestContainer<PlcRequest, PlcResponse> SUT;

    private PlcResponse plcResponse;
    private PlcRequest plcRequest;
    private CompletableFuture<PlcResponse> future;

    @Before
    public void setUp() {
        plcResponse = mock(PlcResponse.class);
        plcRequest = mock(PlcRequest.class);
        future = CompletableFuture.completedFuture(plcResponse);
        SUT = new PlcRequestContainer<>(plcRequest, future);
    }

    @Test
    public void getRequest() {
        PlcRequest request = SUT.getRequest();
        assertNotNull(request);
    }

    @Test
    public void getResponseFuture() throws Exception {
        CompletableFuture<PlcResponse> responseFuture = SUT.getResponseFuture();
        assertNotNull(responseFuture);
        assertNotNull(responseFuture.get(1, TimeUnit.NANOSECONDS));
    }

    @Test
    public void getParent() {
        ProtocolMessage parent = SUT.getParent();
        assertNull(parent);
    }

    @Test
    public void equalsAndHashCode() {
        PlcRequestContainer<PlcRequest, PlcResponse> other = new PlcRequestContainer<>(plcRequest, future);
        assertThat(SUT.hashCode(), IsEqual.equalTo(other.hashCode()));
        assertThat(SUT.equals(other), IsEqual.equalTo(true));
        assertThat(SUT.equals(new Object()), IsEqual.equalTo(false));
        assertEquals(SUT, SUT);
    }

    @Test
    public void testToString() {
        SUT.toString();
    }

}