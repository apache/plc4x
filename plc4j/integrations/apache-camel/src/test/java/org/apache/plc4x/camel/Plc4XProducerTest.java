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
package org.apache.plc4x.camel;

import org.apache.camel.Exchange;
import org.apache.camel.ExchangePattern;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcWriteRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.lang.reflect.Field;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

import static org.mockito.Mockito.*;

public class Plc4XProducerTest {

    private Plc4XProducer SUT;

    private Exchange testExchange;

    @BeforeEach
    public void setUp() throws Exception {
        Plc4XEndpoint endpointMock = mock(Plc4XEndpoint.class, RETURNS_DEEP_STUBS);
        when(endpointMock.getEndpointUri()).thenReturn("plc4x:mock:10.10.10.1/1/1");
        PlcConnection mockConnection = mock(PlcConnection.class, RETURNS_DEEP_STUBS);

        when(mockConnection.getMetadata().canRead()).thenReturn(true);
        when(mockConnection.getMetadata().canWrite()).thenReturn(true);
        when(mockConnection.writeRequestBuilder())
            .thenReturn(mock(PlcWriteRequest.Builder.class, RETURNS_DEEP_STUBS));

        when(endpointMock.getConnection()).thenReturn(mockConnection);
        SUT = new Plc4XProducer(endpointMock);
        testExchange = mock(Exchange.class, RETURNS_DEEP_STUBS);
        Map<String, Map<String,Object>> tags = new HashMap();
        tags.put("test1", Collections.singletonMap("testAddress1",0));
        tags.put("test1", Collections.singletonMap("testAddress2",true));
        tags.put("test1", Collections.singletonMap("testAddress3","TestString"));
        when(testExchange.getIn().getBody())
            .thenReturn(tags);
    }

    @Test
    public void process() throws Exception {
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOnly);
        SUT.process(testExchange);
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOut);
        SUT.process(testExchange);
        when(testExchange.getIn().getBody()).thenReturn(2);

    }

    @Test
    public void process_Async() {
        SUT.process(testExchange, doneSync -> {
        });
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOnly);
        SUT.process(testExchange, doneSync -> {
        });
        when(testExchange.getPattern()).thenReturn(ExchangePattern.InOut);
        SUT.process(testExchange, doneSync -> {
        });
    }

    @Test
    public void doStop() throws Exception {
        SUT.doStop();
    }

    @Test
    public void doStopOpenRequest() throws Exception {
        Field openRequests = SUT.getClass().getDeclaredField("openRequests");
        openRequests.setAccessible(true);
        AtomicInteger atomicInteger = (AtomicInteger) openRequests.get(SUT);
        atomicInteger.incrementAndGet();
        SUT.doStop();
    }

    @Test
    public void doStopBadConnection() throws Exception {
        Field openRequests = SUT.getClass().getDeclaredField("plcConnection");
        openRequests.setAccessible(true);
        PlcConnection plcConnectionMock = mock(PlcConnection.class);
        doThrow(new RuntimeException("oh noes")).when(plcConnectionMock).close();
        openRequests.set(SUT, plcConnectionMock);
        SUT.doStop();
    }

}