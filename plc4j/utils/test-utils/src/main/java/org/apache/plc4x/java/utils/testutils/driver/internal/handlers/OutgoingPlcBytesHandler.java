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

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.ReadBuffer;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.ChannelUtil;
import org.apache.plc4x.java.utils.testutils.utils.hex.HexDiff;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparisonResult;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlDifferenceReporter;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.impl.XmlUnitComparator;
import org.apache.commons.codec.binary.Hex;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.lang.reflect.Method;
import java.lang.reflect.Modifier;

/**
 * Handler for validating outgoing PLC bytes.
 */
public class OutgoingPlcBytesHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(OutgoingPlcBytesHandler.class);

    private final String basePackage;
    private final Element payload;

    public OutgoingPlcBytesHandler(String basePackage, Element payload) {
        this.basePackage = basePackage;
        this.payload = payload;
    }

    /**
     * Executes the outgoing PLC bytes validation.
     *
     * @param transportInstance the transport instance
     * @param byteOrder         the byte order name
     * @throws DriverTestsuiteException if validation fails
     */
    public void executeOutgoingPlcBytes(TransportInstance<?> transportInstance, String byteOrder) {
        // Calculate an expected number of bytes from the hex string
        String expectedHex = payload.getTextTrim().replaceAll("\\s+", "");
        int expectedBytes = expectedHex.length() / 2;

        // Wait for the driver to send the expected number of bytes
        LOGGER.debug("Waiting for {} bytes from driver", expectedBytes);
        final byte[] data = ChannelUtil.waitForOutboundBytes(transportInstance, expectedBytes, 10000);

        validateBytes(payload, data, byteOrder);
    }

    private void validateBytes(Element referenceXml, byte[] data, String byteOrder)
            throws DriverTestsuiteException {
        String expectedHex = referenceXml.getTextTrim();
        // Remove whitespace and make comparison easier
        expectedHex = expectedHex.replaceAll("\\s+", "");

        String actualHex = Hex.encodeHexString(data);

        if (!expectedHex.equalsIgnoreCase(actualHex)) {
            // First show hex diff
            String hexDiff = HexDiff.hexDiff(expectedHex, actualHex);
            String message = String.format(
                "Outgoing bytes don't match.\nExpected: %s\nActual:   %s\nDiff:\n%s",
                expectedHex, actualHex, hexDiff
            );
            LOGGER.error(message);

            // Try to parse and show XML diff for better understanding
            try {
                byte[] expectedBytes = Hex.decodeHex(expectedHex);
                String expectedXml = tryParseAsXml(expectedBytes, byteOrder);
                String actualXml = tryParseAsXml(data, byteOrder);

                if (expectedXml != null && actualXml != null) {
                    LOGGER.error("");
                    LOGGER.error("XML MESSAGE COMPARISON:");
                    LOGGER.error("================================================================================");

                    // Use XML comparison module for side-by-side diff
                    XmlComparator comparator = new XmlUnitComparator();
                    XmlComparisonResult result = comparator.compare(expectedXml, actualXml);
                    XmlDifferenceReporter.logDiffReport(result, LOGGER);
                }
            } catch (Exception e) {
                // If XML parsing fails, just continue with hex diff only
                LOGGER.debug("Could not parse messages as XML: {}", e.getMessage());
            }

            throw new DriverTestsuiteException(message);
        }

        LOGGER.debug("Outgoing bytes match: {}", expectedHex);
    }

    /**
     * Attempts to parse bytes as a protocol message and return XML representation.
     * Returns null if parsing fails.
     */
    private String tryParseAsXml(byte[] data, String byteOrder) {
        try {
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, WithByteBasedOption.WithByteOrder(byteOrder));

            // Try to find a Message class in the base package with a staticParse method
            // We need to introspect to find message types since we don't have explicit type info
            Message message = tryParseMessage(readBuffer, basePackage);

            if (message != null) {
                WriteBufferXmlBased writeBuffer = new WriteBufferXmlBased();
                message.serialize(writeBuffer);
                return writeBuffer.getXmlString();
            }
        } catch (Exception e) {
            LOGGER.trace("Failed to parse message as XML: {}", e.getMessage());
        }
        return null;
    }

    /**
     * Tries to parse a message by attempting different message types from the base package.
     * Returns null if no suitable parser is found.
     */
    private Message tryParseMessage(ReadBuffer readBuffer, String basePackage) {
        // This is a best-effort approach - we try common message class names
        // The proper solution would require message type information in the test XML
        String[] commonMessageTypes = {
            "Packet",
            "Message",
            "Frame",
            "PDU",
            "Request",
            "Response"
        };

        for (String messageType : commonMessageTypes) {
            try {
                String className = basePackage + "." + messageType;
                Class<?> messageClass = Thread.currentThread().getContextClassLoader().loadClass(className);

                // Find staticParse method
                Method staticParseMethod = null;
                for (Method method : messageClass.getDeclaredMethods()) {
                    if (method.getName().equals("staticParse") &&
                        Modifier.isStatic(method.getModifiers()) &&
                        method.getParameterCount() > 0 &&
                        method.getParameterTypes()[0] == ReadBuffer.class) {
                        staticParseMethod = method;
                        break;
                    }
                }

                if (staticParseMethod != null) {
                    // Try to invoke with just the read buffer (no parser arguments)
                    Object[] args = new Object[staticParseMethod.getParameterCount()];
                    args[0] = readBuffer;
                    // Fill remaining args with nulls or defaults
                    for (int i = 1; i < args.length; i++) {
                        args[i] = getDefaultValue(staticParseMethod.getParameterTypes()[i]);
                    }

                    Object result = staticParseMethod.invoke(null, args);
                    if (result instanceof Message) {
                        return (Message) result;
                    }
                }
            } catch (Exception e) {
                // Continue trying other message types
                LOGGER.trace("Failed to parse as {}: {}", messageType, e.getMessage());
            }
        }

        return null;
    }

    /**
     * Gets a default value for a given parameter type.
     */
    private Object getDefaultValue(Class<?> type) {
        if (type == boolean.class || type == Boolean.class) {
            return false;
        } else if (type == byte.class || type == Byte.class) {
            return (byte) 0;
        } else if (type == short.class || type == Short.class) {
            return (short) 0;
        } else if (type == int.class || type == Integer.class) {
            return 0;
        } else if (type == long.class || type == Long.class) {
            return 0L;
        } else if (type == float.class || type == Float.class) {
            return 0.0f;
        } else if (type == double.class || type == Double.class) {
            return 0.0;
        }
        return null;
    }
}
