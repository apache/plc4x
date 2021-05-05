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

package org.apache.plc4x.test.migration;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.mapper.TestSuiteMappingModule;
import org.dom4j.Element;
import org.dom4j.QName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.Diff;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.net.URI;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

// TODO: once migrated from old to new reuse as auto-updater for out-of date testsuites (actual will be used as new expected).
public class MessageValidatorAndMigrator {

    private final static Logger LOGGER = LoggerFactory.getLogger(MessageValidatorAndMigrator.class);

    /**
     * Validates a outbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param testCaseName    name of the testcase
     * @param protocolName    name of the protocol
     * @param outputFlavor    flavor of the output (e.g read-write)
     * @param referenceXml    the xml we expect the outbound message to be
     * @param parserArguments the parser arguments to create an instance of the message
     * @param data            the bytes of the message
     * @param bigEndian       if BOM is big endian
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @throws DriverTestsuiteException if something goes wrong
     */
    @SuppressWarnings({"rawtypes", "unchecked"})
    public static void validateOutboundMessageAndMigrate(String testCaseName, String protocolName, String outputFlavor, Element referenceXml, List<String> parserArguments, byte[] data, boolean bigEndian, boolean autoMigrate, URI siteURI) throws DriverTestsuiteException {
        @Deprecated final String className = referenceXml.attributeValue(new QName("className"));
        MessageIO messageIO = MessageResolver.getMessageIO(protocolName, outputFlavor, referenceXml.getName(), className);
        validateOutboundMessageAndMigrate(testCaseName, messageIO, referenceXml, parserArguments, data, bigEndian, autoMigrate, siteURI);
    }

    /**
     * Validates a outbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param testCaseName    name of the testcase
     * @param messageIO       the pre-constructed MessageIO
     * @param referenceXml    the xml we expect the outbound message to be
     * @param parserArguments the parser arguments to create an instance of the message
     * @param data            the bytes of the message
     * @param bigEndian       if BOM is big endian
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @throws DriverTestsuiteException if something goes wrong
     */
    @SuppressWarnings({"rawtypes", "unchecked"})
    public static void validateOutboundMessageAndMigrate(String testCaseName, MessageIO messageIO, Element referenceXml, List<String> parserArguments, byte[] data, boolean bigEndian, boolean autoMigrate, URI siteURI) throws DriverTestsuiteException {
        final ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, !bigEndian);

        try {
            final Object parsedOutput = messageIO.parse(readBuffer, parserArguments.toArray());
            final String referenceXmlString = referenceXml.asXML();
            try {
                // First try to use the native xml writer
                WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
                messageIO.serialize(writeBufferXmlBased, parsedOutput);
                String xmlString = writeBufferXmlBased.getXmlString();
                final Diff diff = DiffBuilder.compare(referenceXmlString)
                    .withTest(xmlString).checkForSimilar().ignoreComments().ignoreWhitespace()
                    .build();
                if (diff.hasDifferences()) {
                    String border = StringUtils.repeat("=", 100);
                    String centeredDiffDetectedMessage = StringUtils.center(" Diff detected ", 100, "=");
                    String centeredTestCaseName = StringUtils.center(testCaseName, 100, "=");
                    LOGGER.warn(String.format(
                        "\n" +
                            // Border
                            "%1$s\n" +
                            // Testcase name
                            "%5$s\n" +
                            // diff detected message
                            "%2$s\n" +
                            // Border
                            "%1$s\n" +
                            // xml
                            "%3$s\n" +
                            // Border
                            "%1$s\n%1$s\n" +
                            // Text
                            "Differences were found after parsing (Use the above xml in the testsuite to disable this warning).\n" +
                            // Diff
                            "%4$s\n" +
                            // Double Border
                            "%1$s\n%1$s\n" +
                            // Text
                            "Falling back to old jackson based xml mapper\n",
                        border,
                        centeredDiffDetectedMessage,
                        xmlString,
                        diff,
                        centeredTestCaseName));
                    throw new MigrationException(xmlString);
                }
            } catch (RuntimeException | ParseException e) {
                if (!(e instanceof MigrationException)) {
                    LOGGER.error("Error in serializer", e);
                }
                LOGGER.warn("Un-migrated test", e);
                final ObjectMapper mapper = new XmlMapper().enableDefaultTyping().registerModule(new TestSuiteMappingModule());
                String xmlStringFallback = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(parsedOutput);
                final Diff diff2 = DiffBuilder.compare(referenceXmlString)
                    .withTest(xmlStringFallback).checkForSimilar().ignoreComments().ignoreWhitespace()
                    .build();
                if (diff2.hasDifferences()) {
                    String border = StringUtils.repeat("=", 100);
                    String centeredDiffDetectedMessage = StringUtils.center(" Diff detected ", 100, "=");
                    // TODO: insert testcase name
                    String centeredTestCaseName = StringUtils.center("TODO: insert testcase name here", 100, "=");
                    LOGGER.warn(String.format(
                        "\n" +
                            // Border
                            "%1$s\n" +
                            // Testcase name
                            "%5$s\n" +
                            // diff detected message
                            "%2$s\n" +
                            // Border
                            "%1$s\n" +
                            // xml
                            "%3$s\n" +
                            // Border
                            "%1$s\n%1$s\n" +
                            // Text
                            "Differences were found after parsing (Use the above xml in the testsuite to disable this warning).\n" +
                            // Diff
                            "%4$s\n" +
                            // Double Border
                            "%1$s\n%1$s\n" +
                            // Text
                            "Falling back to old jackson based xml mapper\n",
                        border,
                        centeredDiffDetectedMessage,
                        referenceXmlString,
                        diff2,
                        centeredTestCaseName));
                    throw new DriverTestsuiteException("Differences were found after parsing.\n" + diff2);
                } else {
                    LOGGER.info("No diff detected with old");
                }
                if (autoMigrate && e instanceof MigrationException) {
                    Path path = Paths.get(siteURI);
                    LOGGER.info("Migrating {} now", path);
                    Charset charset = StandardCharsets.UTF_8;

                    String content;
                    try {
                        content = new String(Files.readAllBytes(path), charset);
                    } catch (IOException ioException) {
                        throw new RuntimeException(ioException);
                    }
                    String indent = TestCasePatcher.determineIndent(content, xmlStringFallback);
                    String searchString = TestCasePatcher.indent(xmlStringFallback, indent);
                    String newXml = ((MigrationException) e).newXml;
                    newXml = TestCasePatcher.indent(newXml, indent);
                    content = StringUtils.replaceOnce(content, searchString, newXml);
                    try {
                        Files.write(path, content.getBytes(charset));
                    } catch (IOException ioException) {
                        throw new RuntimeException(ioException);
                    }
                    LOGGER.info("Done migrating {}", path);
                }
            }
        } catch (ParseException | JsonProcessingException e) {
            throw new DriverTestsuiteException("Error parsing message", e);
        }
    }

    /**
     * Validates a inbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param protocolName    name of the protocol
     * @param outputFlavor    flavor of the output (e.g read-write)
     * @param referenceXml    the xml we expect the outbound messag
     * @param className       @param className    deprecated fallback classname attribute
     * @param parserArguments the parser arguments to create an instance of the message
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @return the message if all went well
     */
    @SuppressWarnings("rawtypes")
    public static Message validateInboundMessageMigrateAndGet(String protocolName, String outputFlavor, Element referenceXml, @Deprecated String className, List<String> parserArguments, boolean autoMigrate, URI siteURI) {
        MessageIO messageIO = MessageResolver.getMessageIO(protocolName, outputFlavor, referenceXml.getName(), className);
        return validateInboundMessageMigrateAndGet(messageIO, referenceXml, className, parserArguments, autoMigrate, siteURI);
    }

    /**
     * Validates a inbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param messageIO       the pre-constructed MessageIO
     * @param referenceXml    the xml we expect the outbound messag
     * @param className       @param className    deprecated fallback classname attribute
     * @param parserArguments the parser arguments to create an instance of the message
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @return the message if all went well
     */
    @SuppressWarnings({"rawtypes", "unchecked"})
    public static Message validateInboundMessageMigrateAndGet(MessageIO messageIO, Element referenceXml, @Deprecated String className, List<String> parserArguments, boolean autoMigrate, URI siteURI) {

        final String referenceXmlString = referenceXml.asXML();
        try {
            return (Message) messageIO.parse(new ReadBufferXmlBased(new ByteArrayInputStream(referenceXmlString.getBytes(StandardCharsets.UTF_8))), parserArguments.toArray(new String[0]));
        } catch (RuntimeException | ParseException e) {
            LOGGER.warn("Unmigrated test", e);
            final ObjectMapper mapper = new XmlMapper().enableDefaultTyping().registerModule(new TestSuiteMappingModule());
            Message message;
            try {
                message = mapper.readValue(referenceXmlString, getMessageType(className));
            } catch (JsonProcessingException innerE) {
                throw new DriverTestsuiteException("Error parsing message", innerE);
            }
            try {
                WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
                messageIO.serialize(writeBufferXmlBased, message);
                if (autoMigrate) {
                    Path path = Paths.get(siteURI);
                    LOGGER.info("Migrating {} now", path);
                    Charset charset = StandardCharsets.UTF_8;

                    String content;
                    try {
                        content = new String(Files.readAllBytes(path), charset);
                    } catch (IOException ioException) {
                        throw new RuntimeException(ioException);
                    }
                    String indent = TestCasePatcher.determineIndent(content, referenceXmlString);
                    String searchString = indent + referenceXmlString;
                    String newXml = writeBufferXmlBased.getXmlString();
                    newXml = TestCasePatcher.indent(newXml, indent);
                    content = StringUtils.replaceOnce(content, searchString, newXml);
                    try {
                        Files.write(path, content.getBytes(charset));
                    } catch (IOException ioException) {
                        throw new RuntimeException(ioException);
                    }
                    LOGGER.info("Done migrating {}", path);
                }
            } catch (ParseException parseException) {
                LOGGER.warn("could not migrate", e);
            }
            return message;
        }
    }

    @SuppressWarnings("unchecked")
    private static Class<? extends Message> getMessageType(String messageClassName) throws DriverTestsuiteException {
        try {
            final Class<?> messageClass = Class.forName(messageClassName);
            if (Message.class.isAssignableFrom(messageClass)) {
                return (Class<? extends Message>) messageClass;
            }
            throw new DriverTestsuiteException("IO class must implement Message interface");
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException("Error loading message class", e);
        }
    }
}
