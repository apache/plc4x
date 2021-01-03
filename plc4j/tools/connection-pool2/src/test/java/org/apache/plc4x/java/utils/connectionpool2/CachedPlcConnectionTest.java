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

import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.mockito.invocation.InvocationOnMock;
import org.mockito.stubbing.Answer;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

class CachedPlcConnectionTest {

    @Test
    void whenReadFutureFails_handleGracefully() throws ExecutionException, InterruptedException, TimeoutException {
        final CachedDriverManager driverManager = Mockito.mock(CachedDriverManager.class);
        final PlcConnection mockConnection = Mockito.mock(PlcConnection.class, Mockito.RETURNS_DEEP_STUBS);

        when(mockConnection.readRequestBuilder().build().execute()).thenAnswer(new Answer<Object>() {
            @Override
            public Object answer(InvocationOnMock invocation) throws Throwable {
                final CompletableFuture<? extends PlcReadResponse> future = new CompletableFuture<>();
                final Thread thread = new Thread(() -> {
                    try {
                        Thread.sleep(1_000);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                    future.completeExceptionally(new RuntimeException("abc"));
                });
                thread.setDaemon(true);
                thread.start();
                return future;
            }
        });

        final CachedPlcConnection connection = new CachedPlcConnection(driverManager, mockConnection);

        try {
            connection.readRequestBuilder()
                .addItem("a", "b")
                .build()
                .execute()
                .get(5, TimeUnit.SECONDS);
        } catch (Exception e) {
            // Do nothing...
        }

        verify(driverManager).handleBrokenConnection();
    }

    @Test
    void whenReadFutureTimesOut_handleGracefully() throws ExecutionException, InterruptedException, TimeoutException {
        final CachedDriverManager driverManager = Mockito.mock(CachedDriverManager.class);
        final PlcConnection mockConnection = Mockito.mock(PlcConnection.class, Mockito.RETURNS_DEEP_STUBS);

        when(mockConnection.readRequestBuilder().build().execute()).thenAnswer(new Answer<Object>() {
            @Override
            public Object answer(InvocationOnMock invocation) throws Throwable {
                final CompletableFuture<? extends PlcReadResponse> future = new CompletableFuture<>();
                // Return a Future that will never end!
                return future;
            }
        });

        final CachedPlcConnection connection = new CachedPlcConnection(driverManager, mockConnection);

        try {
            connection.readRequestBuilder()
                .addItem("a", "b")
                .build()
                .execute()
                .get(10, TimeUnit.SECONDS);
        } catch (Exception e) {
            // Do nothing...
        }

        verify(driverManager).handleBrokenConnection();
    }

}