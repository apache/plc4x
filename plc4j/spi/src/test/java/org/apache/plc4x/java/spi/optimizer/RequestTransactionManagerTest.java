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
package org.apache.plc4x.java.spi.optimizer;

import org.apache.plc4x.java.spi.transaction.RequestTransactionManager;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assertions.assertTrue;

public class RequestTransactionManagerTest {

    private static final Logger logger = LoggerFactory.getLogger(RequestTransactionManagerTest.class);

    @Test
    public void simpleExample() throws ExecutionException, InterruptedException {
        CompletableFuture<Void> sendRequest = new CompletableFuture<>();
        CompletableFuture<Void> receiveResponse = new CompletableFuture<>();
        CompletableFuture<Void> transactionIsFinished = new CompletableFuture<>();

        RequestTransactionManager tm = new RequestTransactionManager();

        // Assert there is no request going on
        assertEquals(0, tm.getNumberOfActiveRequests());
        // Send a Request
        sendRequest(tm, sendRequest, receiveResponse, transactionIsFinished);
        // Assert that there is a request going on
        sendRequest.get();
        assertEquals(1, tm.getNumberOfActiveRequests());
        // Finish the Request with a response
        receiveResponse.complete(null);
        // Wait till async operation completed
        transactionIsFinished.get();
        // Here, there should no longer be a running request
        // Assert that there is no request going on
        assertEquals(0, tm.getNumberOfActiveRequests());
    }

    @Test
    public void exampleWithMultipleRequests() throws ExecutionException, InterruptedException {
        CompletableFuture<Void> sendRequest1 = new CompletableFuture<>();
        CompletableFuture<Void> endRequest1 = new CompletableFuture<>();
        CompletableFuture<Void> requestIsEnded1 = new CompletableFuture<>();
        CompletableFuture<Void> sendRequest2 = new CompletableFuture<>();
        CompletableFuture<Void> endRequest2 = new CompletableFuture<>();
        CompletableFuture<Void> requestIsEnded2 = new CompletableFuture<>();

        RequestTransactionManager tm = new RequestTransactionManager();

        // Assert there is no request going on
        assertEquals(0, tm.getNumberOfActiveRequests());
        // Send Request 1
        sendRequest(tm, sendRequest1, endRequest1, requestIsEnded1);
        // Send Request 2
        sendRequest(tm, sendRequest2, endRequest2, requestIsEnded2);

        // Assert that there is a request going on
        sendRequest1.get();
        assertEquals(1, tm.getNumberOfActiveRequests());
        // Finish the Request with a response
        endRequest1.complete(null);
        // Wait till async operation (and transaction end) completed
        requestIsEnded1.get();
        // Request 2 should now be processed and finish
        sendRequest2.get();
        endRequest2.complete(null);
        requestIsEnded2.get();
        assertEquals(0, tm.getNumberOfActiveRequests());
    }

    @Test
    public void version2() throws ExecutionException, InterruptedException {
        CompletableFuture<Void> sendRequest = new CompletableFuture<>();
        CompletableFuture<Void> receiveResponse = new CompletableFuture<>();
        CompletableFuture<Void> transactionIsFinished = new CompletableFuture<>();

        RequestTransactionManager tm = new RequestTransactionManager();
        RequestTransactionManager.RequestTransaction handle = tm.startRequest();
        handle.submit(() -> {
            // ...
            sendRequest.complete(null);
            // Receive
            receiveResponse.thenAccept((n) -> {
                handle.endRequest();
                transactionIsFinished.complete(null);
            });
        });

//        // Exception case
//        handle.failRequest(new RuntimeException());

        // Assert that there is a request going on
        sendRequest.get();
        assertEquals(1, tm.getNumberOfActiveRequests());
        // Finish the Request with a response
        receiveResponse.complete(null);
        // Wait till async operation completed
        transactionIsFinished.get();
        // Here, there should no longer be a running request
        // Assert that there is no request going on
        assertEquals(0, tm.getNumberOfActiveRequests());
    }

    @Test
    //@Disabled("This test is randomly failing on Jenkins")
    public void abortTransactionFromExternally() throws ExecutionException, InterruptedException {
        CompletableFuture<Void> sendRequest = new CompletableFuture<>();
        CompletableFuture<Void> receiveResponse = new CompletableFuture<>();

        RequestTransactionManager tm = new RequestTransactionManager();
        RequestTransactionManager.RequestTransaction handle = tm.startRequest();
        handle.submit(() -> {
            // ...
            sendRequest.complete(null);
            // Receive
            receiveResponse.whenComplete((n,e) -> {
                // never execute
                fail();
            });
            //Wait that the fail is handled internally surely and then interrupt this block execute
            try {
                receiveResponse.get();
            } catch (Exception e) {
                assertInstanceOf(InterruptedException.class,e);
            }
        });

        // Assert that there is a request going on
        sendRequest.get();

        // Exception case
        handle.failRequest(new RuntimeException());

        // Wait that the fail is handled internally surely
        //Thread.sleep(100);

        // Assert that no requests are active
        assertEquals(0, tm.getNumberOfActiveRequests());

        // Assert that its cancelled
        assertTrue(handle.getCompletionFuture().isCancelled());
        assertFalse(receiveResponse.isDone());
    }

    private void sendRequest(RequestTransactionManager tm, CompletableFuture<Void> sendRequest, CompletableFuture<Void> endRequest, CompletableFuture<Void> requestIsEnded) {
        tm.submit(handle -> {
            handle.submit(() -> {
                // Wait till the Request is sent
                sendRequest.complete(null);
                // Receive
                endRequest.thenAccept((n) -> {
                    handle.endRequest();
                    requestIsEnded.complete(null);
                });
            });
        });
    }

}