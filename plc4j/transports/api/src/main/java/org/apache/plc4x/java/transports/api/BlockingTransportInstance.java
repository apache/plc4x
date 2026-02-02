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

package org.apache.plc4x.java.transports.api;

import org.apache.plc4x.java.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.transports.api.exceptions.TransportException;

import java.time.Duration;
import java.util.concurrent.TimeoutException;

/**
 * Extension of TransportInstance that supports blocking reads with timeout.
 * This allows the receive loop to block until data is available, eliminating
 * the need for polling with Thread.sleep().
 *
 * @param <T> the configuration type
 */
public interface BlockingTransportInstance<T extends TransportConfiguration> extends TransportInstance<T> {

    /**
     * Blocks until at least one byte is available or the timeout expires.
     * This method eliminates the need for polling loops with Thread.sleep().
     *
     * @param timeout maximum time to wait for data
     * @throws TransportException if an I/O error occurs
     * @throws TimeoutException if the timeout expires before data is available
     * @throws InterruptedException if the thread is interrupted while waiting
     */
    void waitForData(Duration timeout) throws TransportException, TimeoutException, InterruptedException;

    /**
     * Blocks until the specified number of bytes are available or the timeout expires.
     *
     * @param numBytes minimum number of bytes to wait for
     * @param timeout maximum time to wait
     * @throws TransportException if an I/O error occurs
     * @throws TimeoutException if the timeout expires
     * @throws InterruptedException if the thread is interrupted while waiting
     */
    default void waitForBytes(int numBytes, Duration timeout) throws TransportException, TimeoutException, InterruptedException {
        long deadline = System.nanoTime() + timeout.toNanos();
        while (getNumBytesAvailable() < numBytes) {
            long remaining = deadline - System.nanoTime();
            if (remaining <= 0) {
                throw new TimeoutException("Timeout waiting for " + numBytes + " bytes");
            }
            waitForData(Duration.ofNanos(remaining));
        }
    }

}