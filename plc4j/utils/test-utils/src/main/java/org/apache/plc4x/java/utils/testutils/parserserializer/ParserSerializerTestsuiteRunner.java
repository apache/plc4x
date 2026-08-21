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
package org.apache.plc4x.java.utils.testutils.parserserializer;

import org.apache.plc4x.java.spi.buffers.api.Message;
import org.apache.plc4x.java.spi.buffers.api.MessageInput;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.ascii.AsciiBoxWriter;
import org.apache.plc4x.java.spi.buffers.asciiboxbased.utils.hex.Hex;
import org.apache.plc4x.java.spi.buffers.bytebased.ReadBufferByteBased;
import org.apache.plc4x.java.spi.buffers.bytebased.WriteBufferByteBased;
import org.apache.plc4x.java.spi.buffers.xmlbased.WriteBufferXmlBased;
import org.apache.plc4x.java.utils.testutils.XmlTestsuiteLoader;
import org.apache.plc4x.java.utils.testutils.parserserializer.exceptions.ParserSerializerTestsuiteException;
import org.apache.plc4x.java.utils.testutils.parserserializer.model.ParserSerializerTestsuite;
import org.apache.plc4x.java.utils.testutils.parserserializer.model.Testcase;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareDocumentFactory;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareElement;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareSAXReader;
import org.apache.plc4x.java.utils.testutils.utils.hex.HexDiff;
import org.apache.plc4x.java.utils.testutils.utils.migration.MessageResolver;
import org.apache.plc4x.java.utils.testutils.utils.migration.MessageValidatorAndMigrator;
import org.apache.plc4x.java.utils.testutils.utils.xml.XmlHelper;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.QName;
import org.dom4j.io.SAXReader;
import org.junit.jupiter.api.Assumptions;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

import static org.apache.plc4x.java.spi.buffers.bytebased.WithByteBasedOption.WithByteOrder;
import static org.apache.plc4x.java.spi.utils.StaticHelper.DECODE_HEX;
import static org.apache.plc4x.java.spi.utils.StaticHelper.ENCODE_HEX;

public class ParserSerializerTestsuiteRunner extends XmlTestsuiteLoader {

    private static final Logger LOGGER = LoggerFactory.getLogger(ParserSerializerTestsuiteRunner.class);

    /**
     * if set to true if will auto-migrate and on the next run test should be green
     */
    private final boolean autoMigrate;

    private final Set<String> ignoredTestCases = new HashSet<>();

    public ParserSerializerTestsuiteRunner(String testsuiteDocument, String... ignoredTestCases) {
        this(testsuiteDocument, false, ignoredTestCases);
    }

    public ParserSerializerTestsuiteRunner(String testsuiteDocument, boolean autoMigrate, String... ignoredTestCases) {
        super(testsuiteDocument);
        this.autoMigrate = autoMigrate;
        Collections.addAll(this.ignoredTestCases, ignoredTestCases);
    }

    @TestFactory
    public Iterable<DynamicTest> getTestsuiteTests() throws ParserSerializerTestsuiteException {
        ParserSerializerTestsuite testSuite = parseTestsuite();
        List<DynamicTest> dynamicTests = new LinkedList<>();
        for (Testcase testcase : testSuite.testcases()) {
            String testcaseName = testcase.getName();
            String testcaseLabel = testSuite.name() + ": " + testcaseName;
            DynamicTest test = DynamicTest.dynamicTest(testcaseLabel, getSourceUri(testcase), () -> {
                    Assumptions.assumeFalse(() -> ignoredTestCases.contains(testcaseName), "Testcase " + testcaseName + " ignored");
                    run(testSuite, testcase);
                }
            );
            dynamicTests.add(test);
        }
        return dynamicTests;
    }

    /**
     * Parses every prefix of every reference message and requires each one to fail, if it fails, with
     * a {@link BufferException} rather than with an unchecked exception.
     * <p>
     * The remote end decides where a message ends: a datagram can arrive short and a stream can be
     * closed mid-message, so every prefix of a well formed message is something a driver can be
     * handed. Truncation is worth exercising in particular because a field read that runs out of
     * bytes is reported as an absent optional rather than as a failure, so anything deriving a value
     * from that field sees it unset in a case the surrounding condition says cannot happen. An
     * unchecked exception from there escapes the receive path's error handling entirely.
     * <p>
     * Whether a prefix parses is left open - a truncated message may legitimately parse into a
     * smaller one. Only the manner of failure is asserted.
     */
    @TestFactory
    public Iterable<DynamicTest> getTruncationTests() throws ParserSerializerTestsuiteException {
        ParserSerializerTestsuite testSuite = parseTestsuite();
        List<DynamicTest> dynamicTests = new LinkedList<>();
        for (Testcase testcase : testSuite.testcases()) {
            String testcaseName = testcase.getName();
            String testcaseLabel = testSuite.name() + ": " + testcaseName + " truncated";
            dynamicTests.add(DynamicTest.dynamicTest(testcaseLabel, getSourceUri(testcase), () -> {
                    Assumptions.assumeFalse(() -> ignoredTestCases.contains(testcaseName), "Testcase " + testcaseName + " ignored");
                    runTruncations(testSuite, testcase);
                }
            ));
        }
        return dynamicTests;
    }

    private void runTruncations(ParserSerializerTestsuite testSuite, Testcase testcase) throws ParserSerializerTestsuiteException {
        byte[] testcaseRaw = testcase.getRaw();
        MessageInput<?> messageInput = MessageResolver.getMessageIOStaticLinked(
            testSuite.options(),
            testcase.getRootType(),
            testcase.getParserArguments()
        );

        // Up to and including the whole message, so the untruncated case is covered too.
        for (int length = 0; length <= testcaseRaw.length; length++) {
            byte[] truncated = Arrays.copyOf(testcaseRaw, length);
            ReadBufferByteBased readBuffer = new ReadBufferByteBased(truncated, WithByteOrder(testSuite.byteOrderName()));
            try {
                messageInput.parse(readBuffer);
            } catch (BufferException e) {
                // A message we cannot parse is allowed to be rejected.
            } catch (Exception | StackOverflowError e) {
                throw new ParserSerializerTestsuiteException(String.format(
                    "Parsing %s truncated to %d of %d bytes (%s) failed with %s",
                    testcase.getName(), length, testcaseRaw.length, ENCODE_HEX(truncated), e), e);
            }
        }
    }

    private ParserSerializerTestsuite parseTestsuite() throws ParserSerializerTestsuiteException {
        try {
            SAXReader reader = new LocationAwareSAXReader();
            reader.setDocumentFactory(new LocationAwareDocumentFactory());
            Document document = reader.read(testsuiteDocumentXml);
            Element testsuiteXml = document.getRootElement();
            String byteOrderName = testsuiteXml.attributeValue("byteOrder", "BIG_ENDIAN");
            String testsuiteName = testsuiteXml.element(new QName("name")).getStringValue();
            String protocolName = testsuiteXml.element(new QName("protocol-name")).getStringValue();
            String outputFlavor = testsuiteXml.element(new QName("output-flavor")).getStringValue();

            Element optionsElement = testsuiteXml.element(new QName("options"));
            Map<String, String> options = new HashMap<>(XmlHelper.parseParameters(optionsElement));
            options.put("protocol-name", protocolName);
            options.put("output-flavor", outputFlavor);

            List<Element> testcasesXml = testsuiteXml.elements(new QName("testcase"));
            List<Testcase> testcases = new ArrayList<>(testcasesXml.size());
            for (Element testcaseXml : testcasesXml) {
                Element nameElement = testcaseXml.element(new QName("name"));
                Element descriptionElement = testcaseXml.element(new QName("description"));
                Element rawElement = testcaseXml.element(new QName("raw"));
                Element rootTypeElement = testcaseXml.element(new QName("root-type"));
                Element parserArgumentsElement = testcaseXml.element(new QName("parser-arguments"));
                Element xmlElement = testcaseXml.element(new QName("xml"));

                String name = nameElement.getTextTrim();
                String description = (descriptionElement != null) ? descriptionElement.getTextTrim() : null;
                byte[] raw = DECODE_HEX(rawElement.getTextTrim());
                String rootType = rootTypeElement.getTextTrim();

                // Parse additional parser arguments.
                List<String> parserArguments = new LinkedList<>();
                if (parserArgumentsElement != null) {
                    for (Element element : parserArgumentsElement.elements()) {
                        parserArguments.add(element.getTextTrim());
                    }
                }
                Testcase testcase = new Testcase(testsuiteName, protocolName, outputFlavor, name, description, raw, rootType, parserArguments, xmlElement);
                if (testcaseXml instanceof LocationAwareElement) {
                    // pass source location to test
                    testcase.setLocation(((LocationAwareElement) testcaseXml).getLocation());
                }
                testcases.add(testcase);
            }
            LOGGER.info(String.format("Found %d testcases.", testcases.size()));
            return new ParserSerializerTestsuite(testsuiteName, testcases, byteOrderName, options);
        } catch (DocumentException e) {
            throw new ParserSerializerTestsuiteException("Error parsing testsuite xml", e);
        }
    }

    private void run(ParserSerializerTestsuite testSuite, Testcase testcase) throws ParserSerializerTestsuiteException {
        LOGGER.info("Running testcase {}", testcase);
        byte[] testcaseRaw = testcase.getRaw();
        if (LOGGER.isDebugEnabled())
            LOGGER.debug("\n{}", AsciiBoxWriter.DEFAULT.boxString("Testcase raw input", Hex.dump(testcaseRaw), 0));
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(testcaseRaw, WithByteOrder(testSuite.byteOrderName()));

        try {
            MessageInput<?> messageInput = MessageResolver.getMessageIOStaticLinked(
                testSuite.options(),
                testcase.getRootType(),
                testcase.getParserArguments()
            );

            ////////////////////////////////////////////////////////////////////////////////////////////////////////////
            // Parse the raw bytes into a message
            ////////////////////////////////////////////////////////////////////////////////////////////////////////////

            LOGGER.trace("Parsing message");
            Message parsedMessage = (Message) messageInput.parse(readBuffer);
            LOGGER.debug("Parsed message {}", parsedMessage);

            ////////////////////////////////////////////////////////////////////////////////////////////////////////////
            // Compare the parsed message with the reference XML
            ////////////////////////////////////////////////////////////////////////////////////////////////////////////

            LOGGER.trace("Validating and migrating");
            // In this case no reference-xml has been provided
            // (This is usually during development)
            if (testcase.getXml().elements().isEmpty()) {
                WriteBufferXmlBased writeBufferXmlBased = new WriteBufferXmlBased();
                parsedMessage.serialize(writeBufferXmlBased);
                String xmlString = writeBufferXmlBased.getXmlString();
                throw new ParserSerializerTestsuiteException("Missing reference xml element. Parsed: \n" + xmlString);
            }
            // If more than one root element is provided, the testcase is corrupt.
            else if (testcase.getXml().elements().size() > 1) {
                throw new ParserSerializerTestsuiteException("Too many element roots in testcase");
            }
            boolean migrated = MessageValidatorAndMigrator.validateOutboundMessageAndMigrate(
                testcase.getName(),
                parsedMessage,
                testcase.getXml().elements().getFirst(),
                testcaseRaw,
                testSuite.byteOrderName(),
                autoMigrate,
                suiteUri
            );
            if (migrated) {
                LOGGER.warn("Migrated testcase {}", testcase);
            }
            LOGGER.info("Parsing passed for testcase {}", testcase);

            ////////////////////////////////////////////////////////////////////////////////////////////////////////////
            // Serialize the parsed message to a byte array
            ////////////////////////////////////////////////////////////////////////////////////////////////////////////

            LOGGER.trace("Writing message back again");
            byte[] buffer = new byte[parsedMessage.getLengthInBytes()];
            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(buffer, WithByteOrder(testSuite.byteOrderName()));
            parsedMessage.serialize(writeBuffer);
            LOGGER.info("Serializing passed for testcase {}", testcase);
            byte[] data = writeBuffer.getBytes();
            if (testcaseRaw.length != data.length) {
                LOGGER.info("Expected a byte array with a length of {} but got one with {}", testcaseRaw.length, data.length);
            }

            ////////////////////////////////////////////////////////////////////////////////////////////////////////////
            // Compare the serialized bytes to the initial raw array
            ////////////////////////////////////////////////////////////////////////////////////////////////////////////

            if (!Arrays.equals(testcaseRaw, data)) {
                String expectedHex = ENCODE_HEX(testcaseRaw);
                String actualHex = ENCODE_HEX(data);
                String diff = HexDiff.hexDiff(expectedHex, actualHex);
                String message = String.format(
                    "Serialized bytes don't match.\nExpected: %s\nActual:   %s\nDiff:\n%s",
                    expectedHex, actualHex, diff
                );
                LOGGER.error(message);
                throw new ParserSerializerTestsuiteException(message);
            }
        } catch (BufferException e) {
            throw new ParserSerializerTestsuiteException("Unable to parse message", e);
        }
    }

}
