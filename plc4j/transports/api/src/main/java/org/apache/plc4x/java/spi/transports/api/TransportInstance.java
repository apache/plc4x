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

public interface TransportInstance<T extends TransportConfiguration> {

    /**
     * Returns the configuration for this transport instance.
     *
     * @return the configuration for this transport instance.
     */
    T getConfiguration();

    /**
     * Returns the driver-config: the part of the connection URL after the
     * transport-config (host/port) that the transport consumes and before the
     * parameter section. Unlike the transport-config, this portion is meaningful
     * only to the driver. For example, in
     * {@code opcua:tcp://localhost:4840/milo?discovery=false} (or the equivalent
     * {@code opcua:tls://...}) the transport consumes {@code localhost:4840} and
     * this returns {@code /milo}. Transports whose address carries no such trailing
     * part return an empty string.
     *
     * @return the driver-config, never {@code null}
     */
    default String getDriverConfig() {
        return "";
    }

    /**
     * Executes a check if the underlying transport is still open.
     *
     * @return true if the transport is still open, false otherwise.
     */
    boolean isOpen();

    /**
     * Returns the number of bytes available for reading.
     *
     * @return number of bytes available
     * @throws TransportException something went wrong
     */
    int getNumBytesAvailable() throws TransportException;

    /**
     * Returns the given number of bytes from the buffer without moving the read position.
     * Some protocols only require a header to know the size of an incoming message, so some drivers will just peek
     * the size of the header and try to extract the message size from that. If the required number of bytes is
     * available, they will consume these bytes via the read method.
     *
     * @param numBytes number of bytes to peek
     * @return byte array containing the peeked bytes
     * @throws TransportException something went wrong
     */
    byte[] peekReadableBytes(int numBytes) throws TransportException;

    /**
     * Return the given number of bytes and actively mark them as read.
     *
     * @param numBytes number of bytes to read
     * @return byte array containing the read bytes
     * @throws TransportException something went wrong
     */
    byte[] read(int numBytes) throws TransportException;

    /**
     * Simply writes the given bytes to the transport.
     *
     * @param bytes byte array to write
     * @throws TransportException something went wrong
     */
    void write(byte[] bytes) throws TransportException;

    /**
     * Closes the transport.
     *
     * @throws TransportException something went wrong
     */
    void close() throws TransportException;

}
