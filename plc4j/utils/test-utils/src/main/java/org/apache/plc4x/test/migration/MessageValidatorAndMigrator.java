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

import java.util.Map;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.dom4j.Element;
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

public class MessageValidatorAndMigrator {

    private final static Logger LOGGER = LoggerFactory.getLogger(MessageValidatorAndMigrator.class);

    /**
     * Validates a outbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true.
     *
     * Passed options should contain a single 'package' option or 'protocolName' and 'outputFlavor'.
     * In case if package is not specified then protocol name and output flavor (e.g read-write) are
     * used to construct lookup package.
     *
     * @param testCaseName    name of the testcase
     * @param options         map with specific test/lookup options.
     * @param referenceXml    the xml we expect the outbound message to be
     * @param parserArguments the parser arguments to create an instance of the message
     * @param data            the bytes of the message
     * @param byteOrder       the byte-oder being used
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @throws DriverTestsuiteException if something goes wrong
     */
    @SuppressWarnings({"rawtypes", "unchecked"})
    public static void validateOutboundMessageAndMigrate(String testCaseName, Map<String, String> options, Element referenceXml, List<String> parserArguments, byte[] data, ByteOrder byteOrder, boolean autoMigrate, URI siteURI) throws DriverTestsuiteException {
        MessageIO messageIO = MessageResolver.getMessageIO(options, referenceXml.getName());
        validateOutboundMessageAndMigrate(testCaseName, messageIO, referenceXml, parserArguments, data, byteOrder, autoMigrate, siteURI);
    }

    /**
     * Validates a outbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param testCaseName    name of the testcase
     * @param messageIO       the pre-constructed MessageIO
     * @param referenceXml    the xml we expect the outbound message to be
     * @param parserArguments the parser arguments to create an instance of the message
     * @param data            the bytes of the message
     * @param byteOrder       the byte-order being used
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @throws DriverTestsuiteException if something goes wrong
     */
    @SuppressWarnings({"rawtypes", "unchecked"})
    public static void validateOutboundMessageAndMigrate(String testCaseName, MessageIO messageIO, Element referenceXml, List<String> parserArguments, byte[] data, ByteOrder byteOrder, boolean autoMigrate, URI siteURI) throws DriverTestsuiteException {
        final ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, byteOrder);

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
                            "%1$s\n%1$s\n",
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
                    String indent = TestCasePatcher.determineIndent(content, referenceXmlString);
                    String searchString = TestCasePatcher.indent(referenceXmlString, indent);
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
        } catch (ParseException e) {
            throw new DriverTestsuiteException("Error parsing message", e);
        }
    }

    /**
     * Validates a inbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param options         Options which contain custom 'package' name or keys 'protocolName' (name of the protocol)
     *                        and 'outputFlavor' (flavor of the output e.g read-write) which are used to construct
     *                        class lookup root package.
     * @param referenceXml    the xml we expect the outbound message
     * @param parserArguments the parser arguments to create an instance of the message
     * @return the message if all went well
     */
    @SuppressWarnings("rawtypes")
    public static Message validateInboundMessageAndGet(Map<String, String> options, Element referenceXml, List<String> parserArguments) {
        MessageIO messageIO = MessageResolver.getMessageIO(options, referenceXml.getName());
        return validateInboundMessageAndGet(messageIO, referenceXml, parserArguments);
    }

    /**
     * Validates a inbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param messageIO       the pre-constructed MessageIO
     * @param referenceXml    the xml we expect the outbound messag
     * @param parserArguments the parser arguments to create an instance of the message
     * @return the message if all went well
     */
    @SuppressWarnings({"rawtypes", "unchecked"})
    public static Message validateInboundMessageAndGet(MessageIO messageIO, Element referenceXml, List<String> parserArguments) {
        final String referenceXmlString = referenceXml.asXML();
        try {
            return (Message) messageIO.parse(new ReadBufferXmlBased(new ByteArrayInputStream(referenceXmlString.getBytes(StandardCharsets.UTF_8))), parserArguments.toArray(new String[0]));
        } catch (RuntimeException | ParseException e) {
            throw new DriverTestsuiteException(String.format("Error parsing message from:\n%s", referenceXmlString), e);
        }
    }
}
