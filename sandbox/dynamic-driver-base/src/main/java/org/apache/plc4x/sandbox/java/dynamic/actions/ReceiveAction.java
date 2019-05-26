/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.sandbox.java.dynamic.actions;

import org.apache.commons.scxml2.ActionExecutionContext;
import org.apache.commons.scxml2.EventBuilder;
import org.apache.commons.scxml2.TriggerEvent;
import org.apache.daffodil.japi.DataProcessor;
import org.apache.daffodil.japi.ParseResult;
import org.apache.daffodil.japi.infoset.JDOMInfosetOutputter;
import org.apache.daffodil.japi.io.InputSourceDataInputStream;
import org.jdom2.Document;
import org.jdom2.Namespace;
import org.jdom2.Text;
import org.jdom2.filter.Filters;
import org.jdom2.output.XMLOutputter;
import org.jdom2.xpath.XPathExpression;
import org.jdom2.xpath.XPathFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.BufferedInputStream;
import java.io.ByteArrayInputStream;
import java.io.DataInputStream;
import java.net.Socket;
import java.util.List;
import java.util.concurrent.TimeUnit;

public abstract class ReceiveAction extends BaseDaffodilAction {

    private long timeout = 5000;
    private int packetLengthStartPosition;
    private int packetLengthSizeInBytes;
    private int packetLengthOffset = 0;

    @Override
    protected Logger getLogger() {
        return LoggerFactory.getLogger(ReceiveAction.class);
    }

    public String getPacketLengthStartPosition() {
        return Integer.toString(packetLengthStartPosition);
    }

    public void setPacketLengthStartPosition(String packetLengthStartPosition) {
        this.packetLengthStartPosition = Integer.valueOf(packetLengthStartPosition);
    }

    public String getPacketLengthSizeInBytes() {
        return Integer.toString(packetLengthSizeInBytes);
    }

    public void setPacketLengthSizeInBytes(String packetLengthSizeInBytes) {
        this.packetLengthSizeInBytes = Integer.valueOf(packetLengthSizeInBytes);
    }

    public String getPacketLengthOffset() {
        return Integer.toString(packetLengthOffset);
    }

    public void setPacketLengthOffset(String packetLengthOffset) {
        this.packetLengthOffset = Integer.valueOf(packetLengthOffset);
    }

    public String getTimeout() {
        return Long.toString(timeout);
    }

    public void setTimeout(String timeout) {
        this.timeout = Long.valueOf(timeout);
    }

    @Override
    public void execute(ActionExecutionContext ctx) {
        getLogger().info(getStateName() + ": Receiving...");

        try {
            DataProcessor dp = getDaffodilDataProcessor(ctx);
            if(dp == null) {
                fireFailureEvent(ctx, "Couldn't initialize daffodil data processor.");
                return;
            }

            Socket connection = getSocket(ctx);
            DataInputStream inputStream = new DataInputStream(new BufferedInputStream(connection.getInputStream()));

            // Remember when we started to receive.
            long startTime = System.currentTimeMillis();

            // Check if enough bytes are available to at least find out how big the full packet is.
            while(inputStream.available() < packetLengthStartPosition + packetLengthSizeInBytes) {
                waitWithTimeout(ctx, startTime, timeout);
            }

            // Read these length bytes and reset the input stream back to the start.
            inputStream.mark(packetLengthStartPosition + packetLengthSizeInBytes);
            // Jump to the start of the length data.
            inputStream.skip(packetLengthStartPosition);

            // Read the packet length.
            int packetLength;
            switch (packetLengthSizeInBytes) {
                case 1:
                    packetLength = inputStream.readUnsignedByte();
                    break;
                case 2:
                    packetLength = inputStream.readUnsignedShort();
                    break;
                default:
                    fireFailureEvent(ctx, "Unsupported size for packet length: " + packetLengthSizeInBytes);
                    return;
            }
            packetLength += packetLengthOffset;

            // Go back to the beginning of the packet.
            inputStream.reset();

            // Wait till the entire packet is available.
            while(inputStream.available() < packetLength) {
                waitWithTimeout(ctx, startTime, timeout);
            }

            byte[] packet = new byte[packetLength];
            if(inputStream.read(packet) != packetLength) {
                TriggerEvent event = new EventBuilder("failure", TriggerEvent.SIGNAL_EVENT).
                    data("Couldn't read entire packet.").build();
                ctx.getInternalIOProcessor().addEvent(event);
                return;
            }

            // After having enough bytes available, process the current package.
            JDOMInfosetOutputter outputter = new JDOMInfosetOutputter();
            ParseResult byteMessage = dp.parse(
                new InputSourceDataInputStream(new ByteArrayInputStream(packet)), outputter);
            if (byteMessage.isError()) {
                logDiagnosticInformation(byteMessage);
                return;
            }

            // Get the resulting XML document from the parser.
            Document message = outputter.getResult();

            // Do any form of processing.
            processMessage(message, ctx);
        } catch (Exception e) {
            fireFailureEvent(ctx, e.getMessage());
        }

        getLogger().info("Received.");
        fireSuccessEvent(ctx);
    }

    protected abstract void processMessage(Document message, ActionExecutionContext ctx);

    String getRuleText(Document message, String xpathExpression) {
        // Get the namespace definitions from the input document.
        List<Namespace> namespaces = message.getRootElement().getNamespacesInScope();

        XPathFactory xPathFactory = XPathFactory.instance();
        XPathExpression<org.jdom2.Text> xpath = xPathFactory.compile(
            xpathExpression, Filters.textOnly(), null, namespaces);
        List<Text> result = xpath.evaluate(message);
        if((result == null) || result.isEmpty()) {
            getLogger().info("Couldn't find value for xpath expression: " + xpathExpression + " in document.");
            if(getLogger().isInfoEnabled()) {
                getLogger().info(new XMLOutputter().outputString(message));
            }
            return null;
        }
        return result.get(0).getTextNormalize();
    }

    private void waitWithTimeout(ActionExecutionContext ctx, long startTime, long timeout) {
        try {
            TimeUnit.MILLISECONDS.sleep(20);
            if(System.currentTimeMillis() - startTime > timeout) {
                fireFailureEvent(ctx, "Receive timed out.");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

}
