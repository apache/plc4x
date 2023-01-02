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

import org.apache.commons.pool2.impl.GenericKeyedObjectPool;
import org.apache.commons.pool2.impl.GenericKeyedObjectPoolConfig;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.mock.connection.MockConnection;
import org.apache.plc4x.java.mock.connection.MockDevice;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.assertj.core.api.WithAssertions;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.mockito.junit.jupiter.MockitoSettings;
import org.mockito.quality.Strictness;

import java.util.Arrays;
import java.util.Collections;
import java.util.concurrent.ScheduledThreadPoolExecutor;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
@MockitoSettings(strictness = Strictness.LENIENT)
class ScraperTest implements WithAssertions {

    @Mock
    MockDevice mockDevice;

    public static final String CONN_STRING_TIM = "s7://10.10.64.22/0/1";
    public static final String TAG_STRING_TIM = "%DB225:DBW0:INT";

    public static final String CONN_STRING_CH = "s7://10.10.64.20/0/1";
    public static final String TAG_STRING_CH = "%DB3:DBD32:DINT";

    @Test
    @Disabled
    void real_stuff() throws InterruptedException {
        PlcDriverManager driverManager = new PooledPlcDriverManager(pooledPlcConnectionFactory -> {
            GenericKeyedObjectPoolConfig<PlcConnection> config = new GenericKeyedObjectPoolConfig<>();
            config.setJmxEnabled(true);
            config.setMaxWaitMillis(-1);
            config.setMaxTotal(3);
            config.setMinIdlePerKey(0);
            config.setBlockWhenExhausted(true);
            config.setTestOnBorrow(true);
            config.setTestOnReturn(true);
            return new GenericKeyedObjectPool<>(pooledPlcConnectionFactory, config);
        });

        Scraper scraper = new ScraperImpl((j, a, m) -> {}, driverManager, Arrays.asList(
            new ScrapeJobImpl("job1",
                10,
                Collections.singletonMap("tim", CONN_STRING_TIM),
                Collections.singletonMap("distance", TAG_STRING_TIM)
            ),
            new ScrapeJobImpl("job2",
                10,
                Collections.singletonMap("chris", CONN_STRING_CH),
                Collections.singletonMap("counter", TAG_STRING_CH)
            )
        ));

        Thread.sleep(30_000_000);
    }

    @Test
    void scraper_schedulesJob() throws InterruptedException, PlcConnectionException {
        PlcDriverManager driverManager = new PlcDriverManager();
        MockConnection connection = (MockConnection) driverManager.getConnection("mock:m1");
        connection.setDevice(mockDevice);

        when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcDINT(1)));

        ScraperImpl scraper = new ScraperImpl((j, a, m) -> {}, driverManager, Collections.singletonList(
            new ScrapeJobImpl("job1",
                10,
                Collections.singletonMap("m1", "mock:m1"),
                Collections.singletonMap("tag1", "qry1")
            )
        ));

        scraper.start();

        Thread.sleep(1_000);

        // Assert that tasks got done.
        assertThat(scraper.getScheduler()).isInstanceOf(ScheduledThreadPoolExecutor.class);
        assertThat(scraper.getNumberOfActiveTasks())
            .isEqualTo(1);
        assertThat(((ScheduledThreadPoolExecutor) scraper.getScheduler()).getCompletedTaskCount())
            .isGreaterThan(10);
    }

    @Test
    void stop_stopsAllJobs() {
        PlcDriverManager driverManager = new PlcDriverManager();

        Scraper scraper = new ScraperImpl((j, a, m) -> {}, driverManager, Collections.singletonList(
            new ScrapeJobImpl("job1",
                1,
                Collections.singletonMap("m1", "mock:m1"),
                Collections.singletonMap("tag1", "qry1")
            )
        ));

        scraper.start();

        assertThat(scraper.getNumberOfActiveTasks())
            .isEqualTo(1);

        scraper.stop();

        assertThat(scraper.getNumberOfActiveTasks())
            .isZero();
    }

    @Test
    void restart_works() throws PlcConnectionException {
        PlcDriverManager driverManager = new PlcDriverManager();
        MockConnection connection = (MockConnection) driverManager.getConnection("mock:m1");
        connection.setDevice(mockDevice);

        when(mockDevice.read(any())).thenReturn(new ResponseItem<>(PlcResponseCode.OK, new PlcDINT(1)));

        Scraper scraper = new ScraperImpl((j, a, m) -> {}, driverManager, Collections.singletonList(
            new ScrapeJobImpl("job1",
                1,
                Collections.singletonMap("m1", "mock:m1"),
                Collections.singletonMap("tag1", "qry1")
            )
        ));

        scraper.start();

        assertThat(scraper.getNumberOfActiveTasks())
            .isEqualTo(1);

        scraper.stop();

        assertThat(scraper.getNumberOfActiveTasks())
            .isZero();

        scraper.start();

        assertThat(scraper.getNumberOfActiveTasks())
            .isEqualTo(1);
    }
}