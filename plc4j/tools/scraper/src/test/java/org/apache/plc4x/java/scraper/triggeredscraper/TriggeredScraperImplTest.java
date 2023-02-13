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
package org.apache.plc4x.java.scraper.triggeredscraper;

import org.apache.plc4x.java.DefaultPlcDriverManager;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.scraper.config.ScraperConfiguration;
import org.apache.plc4x.java.scraper.config.ScraperConfigurationClassicImpl;
import org.apache.plc4x.java.scraper.exception.ScraperException;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollector;
import org.apache.plc4x.java.scraper.triggeredscraper.triggerhandler.collector.TriggerCollectorImpl;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import java.io.IOException;
import java.util.Random;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.TimeUnit;

import static org.mockito.Mockito.when;

class TriggeredScraperImplTest {

    private DefaultPlcDriverManager driverManager;
    private MockDevice mockDevice1;
    private MockDevice mockDevice2;

    @BeforeEach
    public void setUp() throws Exception {
        driverManager = new DefaultPlcDriverManager();
        MockConnection mock1Connection = ((MockConnection) driverManager.getConnection("mock:1"));
        MockConnection mock2Connection = ((MockConnection) driverManager.getConnection("mock:2"));

        // Create Mocks
        mockDevice1 = Mockito.mock(MockDevice.class);
        mockDevice2 = Mockito.mock(MockDevice.class);
        // Assign to Connections
        mock1Connection.setDevice(mockDevice1);
        mock2Connection.setDevice(mockDevice2);
    }

    /**
     * Test is added because we assume some strange behavior.
     */
    @Test
    void scrapeMultipleTargets() throws ScraperException, IOException, InterruptedException {
        // Prepare the Mocking
        // Scrate Jobs 1 and 2
        when(mockDevice1.read("%DB810:DBB0:USINT")).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcLINT(1L)));
        when(mockDevice2.read("%DB810:DBB0:USINT")).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcLINT(2L)));
        // Trigger Jobs
        // Trigger var
        Random rand = new Random();
        when(mockDevice1.read(("%M0.3:BOOL"))).thenAnswer(invocationOnMock -> {
            boolean trigger = rand.nextBoolean();
            System.out.println(trigger);
            return new ResponseItem<>(PlcResponseCode.OK, new PlcBOOL(trigger));
        });
        when(mockDevice2.read(("%M0.3:BOOL"))).thenAnswer(invocationOnMock -> {
            boolean trigger = rand.nextBoolean();
            System.out.println("\t\t" + trigger);
            return new ResponseItem<>(PlcResponseCode.OK, new PlcBOOL(trigger));
        });
        // Read var
        when(mockDevice1.read("%DB810:DBW0:INT")).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcLINT(3L)));
        when(mockDevice2.read("%DB810:DBW0:INT")).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcLINT(4L)));

        ScraperConfiguration configuration = ScraperConfiguration.fromFile("src/test/resources/mock-scraper-config.yml", ScraperConfigurationClassicImpl.class);
        TriggerCollector triggerCollector = new TriggerCollectorImpl(driverManager);
        TriggeredScraperImpl scraper = new TriggeredScraperImpl((j, a, m) -> System.out.printf("Results from %s/%s: %s%n", j, a, m), driverManager, configuration.getJobs(), triggerCollector, 1000);

        scraper.start();

        new Timer().schedule(new TimerTask() {
            @Override
            public void run() {
                scraper.stop();
            }
        }, TimeUnit.SECONDS.toMillis(2));
    }
}