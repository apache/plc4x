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
package org.apache.plc4x.java.utils.testutils.driver.internal.handlers;

import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.ChannelUtil;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Handler for injecting incoming PLC bytes into the channel.
 */
public class IncomingPlcBytesHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(IncomingPlcBytesHandler.class);

    private final Element payload;

    public IncomingPlcBytesHandler(Element payload) {
        this.payload = payload;
    }

    /**
     * Executes the incoming PLC bytes injection.
     *
     * @param transportInstance the transport instance
     * @param byteOrder         the byte order name
     * @throws DriverTestsuiteException if execution fails
     */
    public void executeIncomingPlcBytes(TransportInstance<?> transportInstance, String byteOrder) {
        String hexString = payload.getTextTrim();
        // Remove whitespace
        hexString = hexString.replaceAll("\\s+", "");

        try {
            byte[] bytes = Hex.decodeHex(hexString);
            ChannelUtil.writeInboundBytes(transportInstance, bytes);
            LOGGER.debug("Injected {} bytes into transport: {}", bytes.length, hexString);
        } catch (DecoderException e) {
            throw new DriverTestsuiteException("Failed to decode hex string: " + hexString, e);
        }
    }
}
