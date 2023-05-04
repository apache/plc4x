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
package org.apache.plc4x.java.opcua.context;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class SecureChannelTransactionManager {

    private static final Logger LOGGER = LoggerFactory.getLogger(SecureChannel.class);
    public static final int DEFAULT_MAX_REQUEST_ID = 0xFFFFFFFF;
    private final AtomicInteger transactionIdentifierGenerator = new AtomicInteger(0);
    private final AtomicInteger requestIdentifierGenerator = new AtomicInteger(0);
    private final AtomicInteger activeTransactionId = new AtomicInteger(0);
    private final Map<Integer, Transaction> queue = new HashMap<>();

    public synchronized void submit(Consumer<Integer> onSend, Integer transactionId) {
        LOGGER.info("Active transaction Number {}", activeTransactionId.get());
        if (activeTransactionId.get() == transactionId) {
            onSend.accept(transactionId);
            int newTransactionId = getActiveTransactionIdentifier();
            if (!queue.isEmpty()) {
                Transaction t = queue.remove(newTransactionId);
                if (t == null) {
                    LOGGER.info("Length of Queue is {}", queue.size());
                    LOGGER.info("Transaction ID is {}", newTransactionId);
                    LOGGER.info("Map  is {}", queue);
                    throw new PlcRuntimeException("Transaction Id not found in queued messages {}");
                }
                submit(t.getConsumer(), t.getTransactionId());
            }
        } else {
            LOGGER.info("Storing out of order transaction {}", transactionId);
            queue.put(transactionId, new Transaction(onSend, transactionId));
        }
    }

    /**
     * Returns the next transaction identifier.
     *
     * @return the next sequential transaction identifier
     */
    public int getTransactionIdentifier() {
        int transactionId = transactionIdentifierGenerator.getAndIncrement();
        if(transactionIdentifierGenerator.get() == DEFAULT_MAX_REQUEST_ID) {
            transactionIdentifierGenerator.set(1);
        }
        return transactionId;
    }

    /**
     * Returns the next transaction identifier.
     *
     * @return the next sequential transaction identifier
     */
    private int getActiveTransactionIdentifier() {
        int transactionId = activeTransactionId.incrementAndGet();
        if(activeTransactionId.get() == DEFAULT_MAX_REQUEST_ID) {
            activeTransactionId.set(1);
        }
        return transactionId;
    }

    public static class Transaction {

        private final Integer transactionId;
        private final Consumer<Integer> consumer;

        public Transaction(Consumer<Integer> consumer, Integer transactionId) {
            this.consumer = consumer;
            this.transactionId = transactionId;
        }

        public Integer getTransactionId() {
            return transactionId;
        }

        public Consumer<Integer> getConsumer() {
            return consumer;
        }
    }

}
