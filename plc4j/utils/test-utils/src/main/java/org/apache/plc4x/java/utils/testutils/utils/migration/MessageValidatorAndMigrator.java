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
package org.apache.plc4x.java.utils.testutils.utils.migration;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.MessageInput;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.xmlbased.ReadBufferXmlBased;
import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlAutoMigrator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparisonResult;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlDifferenceReporter;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.impl.XmlUnitComparator;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayInputStream;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.Map;

import static org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption.WithByteOrder;

public class MessageValidatorAndMigrator {

    private final static Logger LOGGER = LoggerFactory.getLogger(MessageValidatorAndMigrator.class);

    /*
     * Validates an outbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true.
     * <p>
     * Passed options should contain a single 'package' option or 'protocolName' and 'outputFlavor'.
     * If the package is not specified, then protocol-name and output flavor (e.g., read-write) are
     * used to construct lookup package.
     *
     * @param testCaseName    name of the testcase
     * @param options         map with specific test/lookup options.
     * @param referenceXml the XML output we are expecting in the outbound message
     * @param parserArguments the parser arguments to create an instance of the message
     * @param data            the bytes of the message
     * @param byteOrderName   the byte-oder being used
     * @param autoMigrate     indicates if we want to migrate to a new version
     * @param siteURI         the file which we want to auto migrate
     * @throws DriverTestsuiteException if something goes wrong
    public static void validateOutboundMessageAndMigrate(String testCaseName, Map<String, String> options, Element referenceXml, List<String> parserArguments, byte[] data, String byteOrderName, boolean autoMigrate, URI siteURI) throws DriverTestsuiteException {
    MessageInput<?> messageInput = MessageResolver.getMessageInput(options, referenceXml.getName(), parserArguments);
    validateOutboundMessageAndMigrate(testCaseName, parsedMessage, referenceXml, data, byteOrderName, autoMigrate, siteURI);
    }
     */

    /**
     * Validates an outbound message and migrates it to the expectation if the parameter
     * {@code autoMigrate} is set to true
     *
     * @param testCaseName  name of the testcase
     * @param parsedMessage the parsed Message
     * @param referenceXml  the XML output we are expecting in the outbound message
     * @param data          the bytes of the message
     * @param byteOrderName the byte-order being used
     * @param autoMigrate   indicates if we want to migrate to a new version
     * @param siteURI       the file which we want to auto migrate
     * @return true if migration happened
     * @throws DriverTestsuiteException if something goes wrong
     */
    public static boolean validateOutboundMessageAndMigrate(String testCaseName, Message parsedMessage, Element referenceXml, byte[] data, String byteOrderName, boolean autoMigrate, URI siteURI) throws DriverTestsuiteException {
        final ReadBufferByteBased readBuffer = new ReadBufferByteBased(data, WithByteOrder(byteOrderName));

        try {
            final String referenceXmlString = referenceXml.asXML();
            try {
                // First, try to use the native xml writer
                WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
                parsedMessage.serialize(writeBufferXmlBased);
                String xmlString = writeBufferXmlBased.getXmlString();

                // Use the shared XML comparison module
                XmlComparator comparator = new XmlUnitComparator();
                XmlComparisonResult result = comparator.compare(referenceXmlString, xmlString);

                if (result.hasDifferences()) {
                    String diffReport = XmlDifferenceReporter.generateDiffReport(result);
                    LOGGER.error("XML mismatch:\n{}", diffReport);
                    throw new MigrationException(xmlString, diffReport);
                }
                return false;
            } catch (RuntimeException | BufferException e) {
                if (!(e instanceof MigrationException)) {
                    LOGGER.error("Error in serializer", e);
                }
                if (autoMigrate && e instanceof MigrationException me) {
                    LOGGER.info("Migrating {} now", siteURI);

                    // Use the shared auto-migration module
                    XmlAutoMigrator.migrateString(referenceXmlString, me.newXml, siteURI);

                    LOGGER.info("Done migrating {}", siteURI);
                    return true;
                } else if (e instanceof MigrationException me) {
                    String message;
                    if (me.diffReport != null) {
                        message = "XML mismatch detected:\n" + me.diffReport + "\n\nSet autoMigrate=true to fix";
                    } else {
                        message = "Output doesn't match.\nGot:\n" + me.newXml + "\nSet to auto migrate to fix";
                    }
                    LOGGER.error(message);
                    throw new RuntimeException(message, e);
                } else {
                    String message = "Output doesn't match. Set to auto migrate to fix";
                    LOGGER.error(message);
                    throw new RuntimeException(message, e);
                }
            }
        } catch (RuntimeException e) {
            LOGGER.error("Something went wrong: siteURI='{}'", siteURI, e);
            throw e;
        }
    }

    /**
     * Validates an inbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param options      Options which contain custom 'package' name or keys 'protocolName' (name of the protocol)
     *                     and 'outputFlavor' (flavor of the output e.g., read-write) which are used to construct
     *                     class lookup root package.
     * @param referenceXml the XML output we are expecting in the outbound message
     * @return the message if all went well
     */
    public static Message validateInboundMessageAndGet(Map<String, String> options, Element referenceXml) {
        MessageInput<?> messageIO = MessageResolver.getMessageInput(options, referenceXml.getName(), Collections.emptyList());
        return validateInboundMessageAndGet(messageIO, referenceXml);
    }

    /**
     * Validates an inbound message and migrates it to the expectation if the parameter {@code autoMigrate} is set to true
     *
     * @param messageInput the pre-constructed MessageInput
     * @param referenceXml the XML output we are expecting in the outbound message
     * @return the message if all went well
     */
    @SuppressWarnings({"rawtypes"})
    public static Message validateInboundMessageAndGet(MessageInput messageInput, Element referenceXml) {
        final String referenceXmlString = referenceXml.asXML();
        try {
            return (Message) messageInput.parse(new ReadBufferXmlBased(new ByteArrayInputStream(referenceXmlString.getBytes(StandardCharsets.UTF_8))));
        } catch (RuntimeException | BufferException e) {
            throw new DriverTestsuiteException(String.format("Error parsing message from:\n%s", referenceXmlString), e);
        }
    }

}
