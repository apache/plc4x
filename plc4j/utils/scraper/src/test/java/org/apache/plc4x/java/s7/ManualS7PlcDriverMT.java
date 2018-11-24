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

package org.apache.plc4x.java.s7;

import org.apache.commons.math3.stat.descriptive.DescriptiveStatistics;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.junit.Test;

import java.util.ArrayList;
import java.util.Locale;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

public class ManualS7PlcDriverMT {

    public static final String CONN_STRING = "s7://10.10.64.22/0/1";
    public static final String FIELD_STRING = "%DB225:DBW0:INT";

//    public static final String CONN_STRING = "s7://10.10.64.20/0/1";
//    public static final String FIELD_STRING = "%DB3:DBD32:DINT";

    @Test
    public void simpleLoop() {
        PlcDriverManager plcDriverManager = new PooledPlcDriverManager();

        DescriptiveStatistics statistics = new DescriptiveStatistics();
        for (int i = 1; i <= 1000; i++) {

            long start = System.nanoTime();
            try (PlcConnection connection = plcDriverManager.getConnection(CONN_STRING)) {
                CompletableFuture<? extends PlcReadResponse> future = connection.readRequestBuilder()
                    .addItem("distance", FIELD_STRING)
                    .build()
                    .execute();

                PlcReadResponse response = future.get(10, TimeUnit.SECONDS);

                System.out.println(i + " " + response.getLong("distance"));
            } catch (Exception e) {
                e.printStackTrace();
            }
            long end = System.nanoTime();
            statistics.addValue((double)(end-start));
        }

        printStatistics(statistics);
    }

    @Test
    public void scheduledLoop() {

    }

    private void printStatistics(DescriptiveStatistics statistics) {
        System.out.println("Mean response time: " + TimeUnit.NANOSECONDS.toMillis((long)statistics.getMean()) + " ms");
        System.out.println("Median response time: " + TimeUnit.NANOSECONDS.toMillis((long)statistics.getPercentile(50)) + " ms");
        for (int i = 10; i <= 90; i += 10) {
            System.out.println(String.format(Locale.ENGLISH, "Percentile %3d %%: %5d ms", i, TimeUnit.NANOSECONDS.toMillis((long) statistics.getPercentile(i))));
        }
        for (int i = 91; i <= 100; i++) {
            System.out.println(String.format(Locale.ENGLISH, "Percentile %3d %%: %5d ms", i, TimeUnit.NANOSECONDS.toMillis((long) statistics.getPercentile(i))));
        }
    }
}
