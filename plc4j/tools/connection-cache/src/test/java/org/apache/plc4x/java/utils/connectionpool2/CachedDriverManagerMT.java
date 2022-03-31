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

import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.ExecutionException;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

@Disabled("These should be run manually")
class CachedDriverManagerMT {

    private static final Logger LOGGER = LoggerFactory.getLogger(CachedDriverManagerMT.class);

    public static final String PLC_IP = "s7://192.168.167.210/1/1";
//    public static final String PLC_IP = "s7://127.0.0.1/1/1";

    @Test
    void queryPlc() throws InterruptedException {
        ScheduledExecutorService pool = Executors.newScheduledThreadPool(4);

        CachedDriverManager driverManager = new CachedDriverManager(PLC_IP, () -> {
            try {
                PlcConnection connection = new PlcDriverManager().getConnection(PLC_IP);
                // Kill it every second
                pool.schedule(() -> {
                    try {
                        System.out.println("Close...");
                        connection.close();
                    } catch (Exception e) {
                        LOGGER.warn("error", e);
                    }
                }, 1, TimeUnit.SECONDS);
                return connection;
            } catch (PlcConnectionException e) {
                throw new RuntimeException("", e);
            }
        });

        for (int i = 1; i <= 100_000; i++) {
            pool.submit(() -> {
                try (PlcConnection conn = driverManager.getConnection(PLC_IP)) {
                    PlcReadResponse response = conn.readRequestBuilder().addItem("asdf", "%DB444:14.0:BOOL").build().execute().get(500, TimeUnit.MILLISECONDS);
                } catch (Exception e) {
                    // Intentionally do nothing...
                }
            });
            Thread.sleep(1);
        }

        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.DAYS);
    }

    @Test
    void queryPlcWithPool() throws InterruptedException {
        ScheduledExecutorService pool = Executors.newScheduledThreadPool(4);

        PooledDriverManager driverManager = new PooledDriverManager();

        for (int i = 1; i <= 100_000; i++) {
            pool.submit(() -> {
                try (PlcConnection conn = driverManager.getConnection(PLC_IP)) {
                    PlcReadResponse response = conn.readRequestBuilder().addItem("asdf", "%DB444:14.0:BOOL").build().execute().get(500, TimeUnit.MILLISECONDS);
                } catch (Exception e) {
                    // Intentionally do nothing...
                }
            });
            Thread.sleep(1);
        }

        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.DAYS);
    }

    @Test
    void borrowAndDontReturn() throws InterruptedException {
        ScheduledExecutorService pool = Executors.newScheduledThreadPool(4);

        PooledDriverManager driverManager = new PooledDriverManager();

        for (int i = 1; i <= 10_000; i++) {
            pool.submit(() -> {
                try {
                    PlcConnection conn = driverManager.getConnection(PLC_IP);
                    System.out.println("Successfully got a Connection...");
                    PlcReadResponse response = conn.readRequestBuilder().addItem("asdf", "%DB444:14.0:BOOL").build().execute().get(500, TimeUnit.MILLISECONDS);
                    System.out.println("Response: " + response.getBoolean("asdf"));
                } catch (PlcConnectionException | InterruptedException | ExecutionException | TimeoutException e) {
                    // Intentionally do Nothing...
                }
            });
            Thread.sleep(1);
        }

        pool.shutdown();
        pool.awaitTermination(1, TimeUnit.DAYS);
    }
}