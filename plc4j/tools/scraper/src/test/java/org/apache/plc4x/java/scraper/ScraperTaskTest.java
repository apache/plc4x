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
package org.apache.plc4x.java.scraper;

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Nested;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.mockito.junit.jupiter.MockitoSettings;
import org.mockito.quality.Strictness;

import java.util.Collections;
import java.util.concurrent.*;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
@MockitoSettings(strictness = Strictness.LENIENT)
public class ScraperTaskTest implements WithAssertions {

    @Mock
    MockDevice mockDevice;

    @Test
    public void scrape() throws PlcConnectionException {
        PlcDriverManager driverManager = new PlcDriverManager();
        MockConnection connection = (MockConnection) driverManager.getConnection("mock:scraper");
        connection.setDevice(mockDevice);
        when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcSTRING("hallo")));

        ScraperTask scraperTask = new ScraperTaskImpl(driverManager, "job1", "m1", "mock:scraper", Collections.singletonMap("a", "b"),
            1_000, ForkJoinPool.commonPool(), (j,a,m) -> {});

        scraperTask.run();
    }

    @Nested
    class Exceptions {

        @Test
        public void badResponseCode_shouldHandleException() throws PlcConnectionException {
            // Given
            PlcDriverManager driverManager = new PlcDriverManager();
            MockConnection connection = (MockConnection) driverManager.getConnection("mock:scraper");
            connection.setDevice(mockDevice);
            when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.NOT_FOUND, new PlcSTRING("hallo")));

            ScraperTask scraperTask = new ScraperTaskImpl(driverManager, "job1", "m1",
                "mock:scraper", Collections.singletonMap("a", "b"), 1_000, ForkJoinPool.commonPool(), (j,a,m) -> {});

            // When
            scraperTask.run();
        }

        @Mock
        PlcDriverManager driverManager;

        @Test
        public void handleConnectionException() throws PlcConnectionException {
            // Given
            when(driverManager.getConnection(anyString())).thenThrow(new PlcConnectionException("stfu"));

            ScraperTask scraperTask = new ScraperTaskImpl(driverManager, "job1", "m1", "mock:scraper", Collections.singletonMap("a", "b"),
                1_000, ForkJoinPool.commonPool(), (j,a,m) -> {});

            ScraperTask spy = spy(scraperTask);
            spy.run();

            verify(spy).handleException(any());
        }

        @Test
        void runByScheduler_handledGracefully() throws PlcConnectionException {
            when(driverManager.getConnection(anyString())).thenThrow(new PlcConnectionException("stfu"));
            ScheduledExecutorService pool = Executors.newScheduledThreadPool(1);
            ScraperTask scraperTask = new ScraperTaskImpl(driverManager, "job1", "m1", "mock:scraper", Collections.singletonMap("a", "b"),
                1_000, ForkJoinPool.commonPool(), (j,a,m) -> {});

            Future<?> future = pool.scheduleAtFixedRate(scraperTask, 0, 10, TimeUnit.MILLISECONDS);

            assertThat(future).isNotDone();
        }

    }
}