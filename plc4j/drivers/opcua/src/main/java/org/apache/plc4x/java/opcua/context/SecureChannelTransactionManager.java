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

import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Supplier;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class SecureChannelTransactionManager {

    public static final int DEFAULT_MAX_REQUEST_ID = 0xFFFFFFFF;

    private final AtomicInteger transactionIdentifierGenerator = new AtomicInteger(1);
    private final AtomicInteger sequenceIdGenerator = new AtomicInteger(1);
    private final AtomicInteger requestHandleGenerator = new AtomicInteger(1);

    /**
     * Returns the next transaction identifier.
     *
     * @return the next sequential transaction identifier
     */
    public int getTransactionIdentifier() {
        // transaction identifier must begin with 1, otherwise .NET standard server fails!
        int transactionId = transactionIdentifierGenerator.getAndIncrement();
        if (transactionId == DEFAULT_MAX_REQUEST_ID) {
            transactionIdentifierGenerator.set(0);
        }
        return transactionId;
    }

    /**
     * Returns the next sequence identifier.
     *
     * @return the next sequential identifier
     */
    private int getSequenceIdentifier() {
        int sequenceId = sequenceIdGenerator.getAndIncrement();
        if (sequenceId == DEFAULT_MAX_REQUEST_ID) {
            sequenceIdGenerator.set(0);
        }
        return sequenceId;
    }

    /**
     * Creates sequence supplier for temporary use by message sender.
     *
     * @return Sequence supplier.
     */
    public Supplier<Integer> getSequenceSupplier() {
        return this::getSequenceIdentifier;
    }

    /**
     * Returns the next request handle
     *
     * @return the next sequential request handle
     */
    public int getRequestHandle() {
        int requestHandle = requestHandleGenerator.getAndIncrement();
        if (requestHandleGenerator.get() == SecureChannelTransactionManager.DEFAULT_MAX_REQUEST_ID) {
            requestHandleGenerator.set(0);
        }
        return requestHandle;
    }

}
