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

package org.apache.plc4x.java.spi.transports.api;

import org.apache.plc4x.java.spi.transports.api.config.TransportConfiguration;
import org.apache.plc4x.java.spi.transports.api.exceptions.TransportException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class TransportInstanceTest {

    @Test
    void testDefaultGetDriverConfig() {
        // Create a minimal implementation to test the default method
        TransportInstance<TransportConfiguration> transportInstance = new TransportInstance<TransportConfiguration>() {
            @Override
            public TransportConfiguration getConfiguration() {
                return null;
            }

            @Override
            public boolean isOpen() {
                return false;
            }

            @Override
            public int getNumBytesAvailable() throws TransportException {
                return 0;
            }

            @Override
            public byte[] peekReadableBytes(int numBytes) throws TransportException {
                return new byte[0];
            }

            @Override
            public byte[] read(int numBytes) throws TransportException {
                return new byte[0];
            }

            @Override
            public void write(byte[] bytes) throws TransportException {
            }

            @Override
            public void close() throws TransportException {
            }
        };

        // The default implementation returns an empty string
        assertEquals("", transportInstance.getDriverConfig());
    }

}
