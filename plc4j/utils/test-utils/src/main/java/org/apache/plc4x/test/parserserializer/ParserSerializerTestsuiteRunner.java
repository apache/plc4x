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
package org.apache.plc4x.test.parserserializer;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.test.XmlTestsuiteLoader;
import org.apache.plc4x.test.dom4j.LocationAwareDocumentFactory;
import org.apache.plc4x.test.dom4j.LocationAwareElement;
import org.apache.plc4x.test.dom4j.LocationAwareSAXReader;
import org.apache.plc4x.test.migration.MessageResolver;
import org.apache.plc4x.test.migration.MessageValidatorAndMigrator;
import org.apache.plc4x.test.parserserializer.exceptions.ParserSerializerTestsuiteException;
import org.apache.plc4x.test.parserserializer.model.ParserSerializerTestsuite;
import org.apache.plc4x.test.parserserializer.model.Testcase;
import org.apache.plc4x.test.xml.XmlHelper;
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

public class ParserSerializerTestsuiteRunner extends XmlTestsuiteLoader {

    private static final Logger LOGGER = LoggerFactory.getLogger(ParserSerializerTestsuiteRunner.class);

    /**
     * if set to true if will automigrate and on the next run test should be green
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
        for (Testcase testcase : testSuite.getTestcases()) {
            String testcaseName = testcase.getName();
            String testcaseLabel = testSuite.getName() + ": " + testcaseName;
            DynamicTest test = DynamicTest.dynamicTest(testcaseLabel, getSourceUri(testcase), () -> {
                    Assumptions.assumeFalse(() -> ignoredTestCases.contains(testcaseName), "Testcase " + testcaseName + " ignored");
                    run(testSuite, testcase);
                }
            );
            dynamicTests.add(test);
        }
        return dynamicTests;
    }

    private ParserSerializerTestsuite parseTestsuite() throws ParserSerializerTestsuiteException {
        try {
            SAXReader reader = new LocationAwareSAXReader();
            reader.setDocumentFactory(new LocationAwareDocumentFactory());
            Document document = reader.read(testsuiteDocumentXml);
            Element testsuiteXml = document.getRootElement();
            ByteOrder byteOrder = ByteOrder.valueOf(testsuiteXml.attributeValue("byteOrder", "BIG_ENDIAN"));
            String testsuiteName = testsuiteXml.element(new QName("name")).getStringValue();
            String protocolName = testsuiteXml.element(new QName("protocolName")).getStringValue();
            String outputFlavor = testsuiteXml.element(new QName("outputFlavor")).getStringValue();

            Element optionsElement = testsuiteXml.element(new QName("options"));
            Map<String, String> options = new HashMap<>(XmlHelper.parseParameters(optionsElement));
            options.put("protocolName", protocolName);
            options.put("outputFlavor", outputFlavor);

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
                byte[] raw = Hex.decodeHex(rawElement.getTextTrim());
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
            return new ParserSerializerTestsuite(testsuiteName, testcases, byteOrder, options);
        } catch (DocumentException e) {
            throw new ParserSerializerTestsuiteException("Error parsing testsuite xml", e);
        } catch (DecoderException e) {
            throw new ParserSerializerTestsuiteException("Error parsing testcase raw data", e);
        }
    }

    private void run(ParserSerializerTestsuite testSuite, Testcase testcase) throws ParserSerializerTestsuiteException {
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(testcase.getRaw(), testSuite.getByteOrder());

        try {
            MessageInput<Message> messageInput = MessageResolver.getMessageIOStaticLinked(
                testSuite.getOptions(),
                testcase.getXml().elements().get(0).getName()
            );
            Message parsedOutput = messageInput.parse(readBuffer, testcase.getParserArguments().toArray());
            MessageValidatorAndMigrator.validateOutboundMessageAndMigrate(
                testcase.getName(),
                messageInput,
                testcase.getXml().elements().get(0),
                testcase.getParserArguments(),
                testcase.getRaw(),
                testSuite.getByteOrder(),
                autoMigrate,
                suiteUri
            );

            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(parsedOutput.getLengthInBytes(), testSuite.getByteOrder());
            parsedOutput.serialize(writeBuffer);
            byte[] data = writeBuffer.getData();
            if (testcase.getRaw().length != data.length) {
                LOGGER.info("Expected a byte array with a length of " + testcase.getRaw().length +
                    " but got one with " + data.length);
            }
            // TODO: improve output
            if (!Arrays.equals(testcase.getRaw(), data)) {
                int i;
                for (i = 0; i < data.length; i++) {
                    if (data[i] != testcase.getRaw()[i]) {
                        break;
                    }
                }
                throw new ParserSerializerTestsuiteException("Differences were found after serializing.\nExpected: " +
                    Hex.encodeHexString(testcase.getRaw()) + "\nBut Got:  " + Hex.encodeHexString(data) +
                    "\n          " + String.join("", Collections.nCopies(i, "--")) + "^");
            }
        } catch (SerializationException | ParseException e) {
            throw new ParserSerializerTestsuiteException("Unable to parse message", e);
        }
    }

}
