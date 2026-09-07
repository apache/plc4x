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
package org.apache.plc4x.java.utils.testutils.driver.internal.utils;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.apache.plc4x.java.transport.test.TestTransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Utility class for working with transport instances in tests.
 */
public class ChannelUtil {

    private static final Logger LOGGER = LoggerFactory.getLogger(ChannelUtil.class);

    private ChannelUtil() {
        // Utility class
    }

    /**
     * Reads all outbound bytes from a test transport.
     *
     * @param transportInstance the transport instance
     * @return the bytes that were written
     */
    public static byte[] getOutboundBytes(TransportInstance<?> transportInstance) {
        if (!(transportInstance instanceof TestTransportInstance)) {
            throw new DriverTestsuiteException("Transport must be a TestTransportInstance for testing");
        }

        TestTransportInstance testTransport = (TestTransportInstance) transportInstance;
        return testTransport.getAllWrittenData();
    }

    /**
     * Waits for and reads the expected number of outbound bytes from a test transport.
     * This method blocks until the expected bytes are available or timeout occurs.
     *
     * @param transportInstance the transport instance
     * @param expectedBytes     the number of bytes to wait for
     * @param timeoutMs         timeout in milliseconds
     * @return the bytes that were written
     */
    public static byte[] waitForOutboundBytes(TransportInstance<?> transportInstance, int expectedBytes, long timeoutMs) {
        if (!(transportInstance instanceof TestTransportInstance)) {
            throw new DriverTestsuiteException("Transport must be a TestTransportInstance for testing");
        }

        TestTransportInstance testTransport = (TestTransportInstance) transportInstance;
        try {
            return testTransport.waitForWrittenData(expectedBytes, timeoutMs);
        } catch (TransportException e) {
            throw new DriverTestsuiteException("Failed waiting for outbound bytes: " + e.getMessage(), e);
        }
    }

    /**
     * Writes bytes to the inbound side of a test transport (simulating incoming data).
     *
     * @param transportInstance the transport instance
     * @param bytes             the bytes to inject
     */
    public static void writeInboundBytes(TransportInstance<?> transportInstance, byte[] bytes) {
        if (!(transportInstance instanceof TestTransportInstance)) {
            throw new DriverTestsuiteException("Transport must be a TestTransportInstance for testing");
        }

        TestTransportInstance testTransport = (TestTransportInstance) transportInstance;
        int written = testTransport.injectTestData(bytes);
        if (written < bytes.length) {
            LOGGER.warn("Only wrote {} of {} bytes to transport", written, bytes.length);
        }
    }

    /**
     * Gets the number of bytes that have been written to the transport.
     *
     * @param transportInstance the transport instance
     * @return number of written bytes
     */
    public static int getNumBytesWritten(TransportInstance<?> transportInstance) {
        if (!(transportInstance instanceof TestTransportInstance)) {
            throw new DriverTestsuiteException("Transport must be a TestTransportInstance for testing");
        }

        TestTransportInstance testTransport = (TestTransportInstance) transportInstance;
        return testTransport.getNumBytesWritten();
    }
}
