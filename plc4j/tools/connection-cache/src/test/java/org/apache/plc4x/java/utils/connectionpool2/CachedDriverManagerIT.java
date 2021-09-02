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
package org.apache.plc4x.java.utils.connectionpool2;

import org.apache.plc4x.java.api.exceptions.PlcException;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Multi Threading Test
 *
 * @author julian
 * Created by julian on 06.04.20
 */
class CachedDriverManagerIT {

    private static final Logger LOGGER = LoggerFactory.getLogger(CachedDriverManagerIT.class);

    @Test
    void connectWithMultpleThreads() throws InterruptedException, PlcException {
        ExecutorService executorService = Executors.newFixedThreadPool(4);

        PlcConnectionFactory mock = Mockito.mock(PlcConnectionFactory.class);
        MockConnection plcMockConnection = mock(MockConnection.class);
        when(mock.create()).thenReturn(plcMockConnection);

        CachedDriverManager driverManager = new CachedDriverManager("", mock, 100_000);

        AtomicInteger errorCounter = new AtomicInteger(0);
        AtomicInteger successCounter = new AtomicInteger(0);

        for (int i = 1; i <= 100; i++) {
            executorService.submit(() -> {
                try {
                    driverManager.getConnection("").close();
                    successCounter.incrementAndGet();
                } catch (Exception e) {
                    LOGGER.warn("error", e);
                    errorCounter.incrementAndGet();
                }
            });
        }

        executorService.shutdown();

        executorService.awaitTermination(50, TimeUnit.SECONDS);

        assertEquals(100, successCounter.get());
        assertEquals(0, errorCounter.get());
    }
}