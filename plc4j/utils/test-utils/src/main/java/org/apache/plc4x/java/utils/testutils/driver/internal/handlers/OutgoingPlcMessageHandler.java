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
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption;
import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.apache.plc4x.java.spi.transports.api.TransportInstance;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.ChannelUtil;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.Delay;
import org.apache.commons.lang3.ClassUtils;
import org.dom4j.Element;
import org.dom4j.QName;
import org.dom4j.io.OutputFormat;
import org.dom4j.io.XMLWriter;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.*;

import java.io.IOException;
import java.io.StringWriter;
import java.lang.reflect.*;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;

/**
 * Handler for validating outgoing PLC messages.
 * NOTE: Full message parsing/validation not yet implemented.
 * Use outgoing-plc-bytes test steps instead for byte-level testing.
 */
public class OutgoingPlcMessageHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(OutgoingPlcMessageHandler.class);

    private final String basePackage;
    private final Element referenceXml;

    public OutgoingPlcMessageHandler(String basePackage, Element referenceXml) {
        this.basePackage = basePackage;
        this.referenceXml = referenceXml;
    }

    /**
     * Executes the outgoing PLC message validation.
     *
     * @param transportInstance the transport instance
     * @param byteOrder         the byte order name
     */
    public void executeOutgoingPlcMessage(TransportInstance<?> transportInstance, String byteOrder) {
        // Wait until the driver has actually written something to the transport. A fixed
        // short delay is unreliable — depending on scheduling the driver may not have
        // reached the messageCodec.send() call yet.
        long deadline = System.currentTimeMillis() + 5000L;
        byte[] actualBytes = ChannelUtil.getOutboundBytes(transportInstance);
        while (actualBytes.length == 0 && System.currentTimeMillis() < deadline) {
            Delay.shortDelay();
            actualBytes = ChannelUtil.getOutboundBytes(transportInstance);
        }
        // Give a brief grace period so we collect the full message rather than just the
        // first few bytes that happened to flush before we read.
        Delay.shortDelay();
        byte[] more = ChannelUtil.getOutboundBytes(transportInstance);
        if (more.length > 0) {
            byte[] combined = new byte[actualBytes.length + more.length];
            System.arraycopy(actualBytes, 0, combined, 0, actualBytes.length);
            System.arraycopy(more, 0, combined, actualBytes.length, more.length);
            actualBytes = combined;
        }

        // Get the name of the root message type
        Optional<Element> messageElementOptional = referenceXml.elements().stream().filter(
            e -> !e.getQName().equals(new QName("parser-arguments"))).findFirst();
        if (messageElementOptional.isEmpty()) {
            throw new RuntimeException("No message element found in reference XML");
        }
        Element messageElement = messageElementOptional.get();
        String rootMessageTypeName = messageElement.getName();
        String className = basePackage + "." + rootMessageTypeName;

        // Parse the message.
        // In general, we take the first element, that's not "parser-argument", add that to the
        // package name, resolve the staticParse method, extract its parameters, parse and convert
        // the content in the "parser-arguments" element and invoke ith with that information.
        // The result should be a Message object representing the parsed message.
        Class<?> messageTypeClass;
        try {
            messageTypeClass = Thread.currentThread().getContextClassLoader().loadClass(className);

            // Get the method named "staticParse"
            Method staticParseMethod = Arrays.stream(messageTypeClass.getDeclaredMethods())
                .filter(method -> method.getName().equals("staticParse") &&
                    Modifier.isStatic(method.getModifiers()) &&
                    method.getParameterTypes()[0] == ReadBuffer.class)
                .findFirst()
                .orElseThrow(() -> new RuntimeException("No staticParse method found in class " + className));

            // The default integer/float encodings match what every driver codec configures;
            // generated parsers (e.g. S7) fail without them on fields that don't pass
            // explicit per-field options.
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(actualBytes,
                WithOption.WithUnsignedIntegerEncoding("unsigned-binary"),
                WithOption.WithSignedIntegerEncoding("twos-complement"),
                WithOption.WithFloatEncoding("IEEE754"),
                WithByteBasedOption.WithByteOrder(byteOrder));

            Element parserArgumentsElement = referenceXml.element(new QName("parser-arguments"));

            // Get the parameter types for the staticParseMethod
            int numArgs = staticParseMethod.getParameters().length;
            Object[] argValues = new Object[numArgs];
            for (int i = 0; i < numArgs; i++) {
                Parameter parameter = staticParseMethod.getParameters()[i];
                if (i == 0) {
                    if (!parameter.getName().equals("readBuffer")) {
                        throw new RuntimeException("Invalid parameter name for first argument. Expected 'readBuffer' but got: " + parameter.getName());
                    }
                    if (!parameter.getType().equals(ReadBuffer.class)) {
                        throw new RuntimeException("Invalid parameter type for first argument. Expected 'ReadBuffer' but got: " + parameter.getType());
                    }
                    argValues[i] = readBuffer;
                } else {
                    String parameterName = parameter.getName();
                    Class<?> parameterType = parameter.getType();
                    if ((parserArgumentsElement == null) || (parserArgumentsElement.element(new QName(parameterName)) == null)) {
                        throw new RuntimeException("No parser-arguments element or parameterName element found for parameter " + parameterName);
                    } else {
                        String parameterStringValue = parserArgumentsElement.element(new QName(parameterName)).getTextTrim();
                        argValues[i] = parseDynamic(parameterType, parameterStringValue);
                    }
                }
            }

            // Actually parse the message
            Object parsed = staticParseMethod.invoke(null, argValues);
            LOGGER.debug("Parsed message: {}", parsed);

            if (!(parsed instanceof Message message)) {
                throw new RuntimeException("Parsed message is not an instance of Message");
            }

            // Write the parsed message as XML, so we can compare it to the reference XML.
            try {
                WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
                writeBufferXmlBased.writeMessage(message);
                String xmlString = writeBufferXmlBased.getXmlString();
                LOGGER.debug("Serialized message: {}", xmlString);

                // Convert reference XML to string for comparison
                String referenceXmlString = elementToString(messageElement);

                // Compare the XMLs
                compareXml(referenceXmlString, xmlString);
            } catch (BufferException e) {
                throw new RuntimeException("Failed to serialize message", e);
            }
        } catch (ClassNotFoundException e) {
            throw new RuntimeException("Failed to load class " + className, e);
        } catch (InvocationTargetException e) {
            throw new RuntimeException("Failed to invoke staticParse method in class " + className, e.getCause() != null ? e.getCause() : e);
        } catch (IllegalAccessException e) {
            throw new RuntimeException("Failed to access staticParse method in class " + className, e);
        }
    }

    public static Object parseDynamic(Class<?> type, String value) {
        try {
            // Try valueOf(String)
            try {
                Method m = type.getMethod("valueOf", String.class);
                if (Modifier.isStatic(m.getModifiers())) {
                    return m.invoke(null, value);
                }
            } catch (NoSuchMethodException ignored) {}

            // If the type is a primitive, get the corresponding object-type.
            if (type.isPrimitive()) {
                type = ClassUtils.primitiveToWrapper(type);
            }

            // Try parse*(String)
            for (Method m : type.getMethods()) {
                if (Modifier.isStatic(m.getModifiers())
                    && m.getParameterCount() == 1
                    && m.getParameterTypes()[0] == String.class
                    && m.getName().startsWith("parse")) {
                    return m.invoke(null, value);
                }
            }

            // String constructor
            try {
                Constructor<?> c = type.getConstructor(String.class);
                return c.newInstance(value);
            } catch (NoSuchMethodException ignored) {}

            throw new IllegalArgumentException(
                "Don't know how to parse type: " + type.getName()
            );

        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    /**
     * Converts a dom4j Element to a formatted XML string.
     */
    private String elementToString(Element element) {
        try {
            StringWriter writer = new StringWriter();
            OutputFormat format = OutputFormat.createPrettyPrint();
            format.setIndentSize(2);
            format.setNewlines(true);
            XMLWriter xmlWriter = new XMLWriter(writer, format);
            xmlWriter.write(element);
            xmlWriter.close();
            return writer.toString();
        } catch (IOException e) {
            throw new RuntimeException("Failed to convert element to string", e);
        }
    }

    /**
     * Compares two XML strings and throws an exception with a detailed diff if they differ.
     * Ignores whitespace and element ordering differences.
     */
    private void compareXml(String expected, String actual) {
        Diff diff = DiffBuilder.compare(expected)
            .withTest(actual)
            .ignoreWhitespace()
            .ignoreComments()
            .checkForSimilar()
            .withNodeMatcher(new DefaultNodeMatcher(ElementSelectors.byNameAndAllAttributes))
            .build();

        if (diff.hasDifferences()) {
            String diffReport = generateSideBySideDiff(expected, actual, diff);
            LOGGER.error("XML mismatch:\n{}", diffReport);
            throw new AssertionError("XML mismatch:\n" + diffReport);
        }

        LOGGER.info("XML comparison successful - outgoing message matches expected");
    }

    /**
     * Generates a side-by-side diff report for XML differences.
     */
    private String generateSideBySideDiff(String expected, String actual, Diff diff) {
        StringBuilder sb = new StringBuilder();

        // Header
        sb.append("\n");
        sb.append("=".repeat(185)).append("\n");
        sb.append("XML COMPARISON FAILED\n");
        sb.append("=".repeat(185)).append("\n\n");

        // List all differences
        sb.append("DIFFERENCES FOUND:\n");
        sb.append("-".repeat(50)).append("\n");

        int diffNum = 1;
        for (Difference difference : diff.getDifferences()) {
            Comparison comparison = difference.getComparison();
            sb.append(String.format("%d. %s%n", diffNum++, formatDifference(comparison)));
        }
        sb.append("\n");

        // Side-by-side comparison
        sb.append("SIDE-BY-SIDE COMPARISON:\n");
        sb.append("-".repeat(185)).append("\n");

        // Filter out empty lines and trim each line
        List<String> expectedLines = Arrays.stream(expected.split("\n"))
            .map(String::trim)
            .filter(line -> !line.isEmpty())
            .toList();
        List<String> actualLines = Arrays.stream(actual.split("\n"))
            .map(String::trim)
            .filter(line -> !line.isEmpty())
            .toList();

        int maxLines = Math.max(expectedLines.size(), actualLines.size());
        int colWidth = 90;

        sb.append(String.format("%-" + colWidth + "s | %s%n", "EXPECTED", "ACTUAL"));
        sb.append("-".repeat(colWidth)).append(" | ").append("-".repeat(colWidth)).append("\n");

        for (int i = 0; i < maxLines; i++) {
            String expLineFull = i < expectedLines.size() ? expectedLines.get(i) : "";
            String actLineFull = i < actualLines.size() ? actualLines.get(i) : "";
            String expLine = truncate(expLineFull, colWidth);
            String actLine = truncate(actLineFull, colWidth);

            // Mark differing lines - compare full lines, not truncated ones
            String marker = expLineFull.equals(actLineFull) ? " " : "*";
            sb.append(String.format("%-" + colWidth + "s %s %s%n", expLine, marker, actLine));
        }

        sb.append("\n");
        sb.append("=".repeat(185)).append("\n");

        // Full XML outputs for reference
        sb.append("\nFULL EXPECTED XML:\n");
        sb.append("-".repeat(50)).append("\n");
        sb.append(expected).append("\n");

        sb.append("\nFULL ACTUAL XML:\n");
        sb.append("-".repeat(50)).append("\n");
        sb.append(actual).append("\n");

        return sb.toString();
    }

    /**
     * Formats a single difference for display.
     */
    private String formatDifference(Comparison comparison) {
        ComparisonType type = comparison.getType();
        String expectedValue = comparison.getControlDetails().getValue() != null
            ? comparison.getControlDetails().getValue().toString() : "null";
        String actualValue = comparison.getTestDetails().getValue() != null
            ? comparison.getTestDetails().getValue().toString() : "null";
        String xpath = comparison.getControlDetails().getXPath() != null
            ? comparison.getControlDetails().getXPath() : "N/A";

        return String.format("Type: %s%n   XPath: %s%n   Expected: %s%n   Actual: %s",
            type, xpath, expectedValue, actualValue);
    }

    /**
     * Truncates a string to the specified length, adding "..." if truncated.
     */
    private String truncate(String str, int maxLength) {
        if (str.length() <= maxLength) {
            return str;
        }
        return str.substring(0, maxLength - 3) + "...";
    }

}
