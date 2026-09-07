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
package org.apache.plc4x.java.openprotocol;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.openprotocol.config.OpenProtocolConfiguration;
import org.apache.plc4x.java.openprotocol.readwrite.OpenProtocolMessage;
import org.apache.plc4x.java.openprotocol.tag.OpenProtocolTagHandler;
import org.apache.plc4x.java.spi.drivers.ConnectionBase;
import org.apache.plc4x.java.spi.drivers.exceptions.MessageCodecException;
import org.apache.plc4x.java.spi.drivers.tags.PlcTagHandler;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.spi.values.DefaultPlcValueHandler;
import org.apache.plc4x.java.spi.values.PlcValueHandler;
import org.apache.plc4x.java.utils.auditlog.api.AuditLog;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Open-Protocol connection skeleton.
 *
 * <p>The legacy {@code OpenProtocolProtocolLogic} was already a no-op stub —
 * onConnect/decode/read/write/subscribe all returned without doing anything
 * useful. This port preserves that shape: the connection successfully opens
 * the transport, sets up a codec, then sits there waiting for someone to
 * implement the actual MID-by-MID exchange. The wire-format parser/serializer
 * (covered by {@code OpenProtocolParserSerializerTest}) is the only piece
 * that's actually verified today.</p>
 */
public class OpenProtocolConnection extends ConnectionBase<OpenProtocolConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(OpenProtocolConnection.class);

    /**
     * Default revision used when parsing incoming messages — the lowest one
     * defined by the spec. The legacy driver hard-coded {@code 1} in
     * {@code OpenProtocolMessage::staticParse}.
     */
    private static final int DEFAULT_REVISION = 1;

    private OpenProtocolMessageCodec messageCodec;
    private volatile boolean connected = false;

    public OpenProtocolConnection(OpenProtocolConfiguration configuration,
                                  TransportInstance<?> transportInstance,
                                  AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        return new OpenProtocolTagHandler();
    }

    @Override
    protected PlcValueHandler getValueHandler() {
        return new DefaultPlcValueHandler();
    }

    @Override
    public boolean isConnected() {
        return connected && messageCodec != null && messageCodec.isOpen();
    }

    @Override
    protected void onConnect() throws PlcConnectionException {
        messageCodec = new OpenProtocolMessageCodec(
            transportInstance, this::handleIncomingMessage, DEFAULT_REVISION);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming Open-Protocol data", e);
            }
        });
        connected = true;
    }

    @Override
    public void close() throws Exception {
        connected = false;
        stopReceiving();
        if (messageCodec != null) {
            messageCodec.close();
        }
        super.close();
    }

    private void handleIncomingMessage(OpenProtocolMessage message) {
        // The legacy stub had an empty decode(); incoming messages have no
        // consumer yet. Logged at trace so a future implementer can see what
        // showed up on the wire.
        LOGGER.trace("Received Open-Protocol message {}", message);
    }

}
