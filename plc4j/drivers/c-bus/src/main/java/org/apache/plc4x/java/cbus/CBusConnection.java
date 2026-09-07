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
package org.apache.plc4x.java.cbus;

import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.cbus.configuration.CBusConfiguration;
import org.apache.plc4x.java.cbus.readwrite.CBusCommand;
import org.apache.plc4x.java.cbus.readwrite.CBusOptions;
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
 * C-Bus connection skeleton.
 *
 * <p>The legacy {@code CBusProtocolLogic} was a placeholder — read returned an
 * unresolved future, decode was empty, and no tag handler was wired up. This
 * port preserves that shape so the module compiles and the parser/serializer
 * tests continue to verify the wire format; actual command sequencing remains
 * future work.</p>
 */
public class CBusConnection extends ConnectionBase<CBusConfiguration> {

    private static final Logger LOGGER = LoggerFactory.getLogger(CBusConnection.class);

    /** All-false options matches what {@code CBusDriver} hard-coded for {@code staticParse}. */
    private static final CBusOptions DEFAULT_OPTIONS =
        new CBusOptions(false, false, false, false, false, false, false, false, false);

    private CBusMessageCodec messageCodec;
    private volatile boolean connected = false;

    public CBusConnection(CBusConfiguration configuration,
                          TransportInstance<?> transportInstance,
                          AuditLog auditLog) {
        super(configuration, transportInstance, auditLog);
    }

    @Override
    protected PlcTagHandler getTagHandler() {
        // Legacy returned null too; no tag syntax is defined for c-bus yet.
        return null;
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
        messageCodec = new CBusMessageCodec(
            transportInstance, this::handleIncomingMessage, DEFAULT_OPTIONS);
        startReceiving(() -> {
            try {
                messageCodec.processIncomingData();
            } catch (MessageCodecException e) {
                LOGGER.error("Error processing incoming C-Bus data", e);
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

    private void handleIncomingMessage(CBusCommand command) {
        LOGGER.trace("Received C-Bus command {}", command);
    }

}
