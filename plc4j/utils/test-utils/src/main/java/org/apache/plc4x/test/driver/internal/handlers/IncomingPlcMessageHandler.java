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
package org.apache.plc4x.test.driver.internal.handlers;

import io.netty.buffer.Unpooled;
import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.DriverTestsuiteConfiguration;
import org.apache.plc4x.test.migration.MessageResolver;
import org.apache.plc4x.test.migration.MessageValidatorAndMigrator;
import org.dom4j.Element;

import java.util.List;

public class IncomingPlcMessageHandler {

    private final DriverTestsuiteConfiguration driverTestsuiteConfiguration;

    private final Element payload;

    private final List<String> parserArguments;

    public IncomingPlcMessageHandler(DriverTestsuiteConfiguration driverTestsuiteConfiguration, Element payload, List<String> parserArguments) {
        this.driverTestsuiteConfiguration = driverTestsuiteConfiguration;
        this.payload = payload;
        this.parserArguments = parserArguments;
    }

    public void executeIncomingPlcMessage(Plc4xEmbeddedChannel embeddedChannel, ByteOrder byteOrder) {
        // Get a byte representation of the incoming message.
        final byte[] data = getBytesFromXml(payload, byteOrder);
        // Send the bytes to the channel.
        embeddedChannel.writeInbound(Unpooled.wrappedBuffer(data));
    }

    @SuppressWarnings({"rawtypes", "unchecked"})
    public byte[] getBytesFromXml(Element referenceXml, ByteOrder byteOrder) throws DriverTestsuiteException {
        final WriteBufferByteBased writeBuffer = new WriteBufferByteBased(1024, byteOrder);
        MessageInput messageInput = MessageResolver.getMessageInput(driverTestsuiteConfiguration.getOptions(), referenceXml.getName());
        // Get Message and Validate
        Message message = MessageValidatorAndMigrator.validateInboundMessageAndGet(messageInput, referenceXml, parserArguments);

        // Get Bytes
        try {
            message.serialize(writeBuffer);
            final byte[] data = new byte[message.getLengthInBytes()];
            System.arraycopy(writeBuffer.getData(), 0, data, 0, writeBuffer.getPos());
            return data;
        } catch (SerializationException e) {
            throw new DriverTestsuiteException("Error serializing message", e);
        }
    }

}
