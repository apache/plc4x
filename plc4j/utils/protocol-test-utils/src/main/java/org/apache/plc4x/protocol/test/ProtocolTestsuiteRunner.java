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

package org.apache.plc4x.protocol.test;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.plc4x.java.utils.*;
import org.apache.plc4x.protocol.test.exceptions.ProtocolTestsuiteException;
import org.apache.plc4x.protocol.test.model.ProtocolTestsuite;
import org.apache.plc4x.protocol.test.model.Testcase;
import org.dom4j.*;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import org.dom4j.io.SAXReader;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.Diff;

import java.io.InputStream;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.util.*;

public class ProtocolTestsuiteRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(ProtocolTestsuiteRunner.class);

    private final String testsuiteDocument;

    public ProtocolTestsuiteRunner(String testsuiteDocument) {
        this.testsuiteDocument = testsuiteDocument;
    }

    @TestFactory
    public Iterable<DynamicTest> getTestsuiteTests() throws ProtocolTestsuiteException {
        ProtocolTestsuite testSuite = parseTestsuite(ProtocolTestsuiteRunner.class.getResourceAsStream(testsuiteDocument));
        List<DynamicTest> dynamicTests = new LinkedList<>();
        for(Testcase testcase : testSuite.getTestcases()) {
            String testcaseName = testcase.getName();
            String testcaseLabel = testSuite.getName() + ": " + testcaseName;
            DynamicTest test = DynamicTest.dynamicTest(testcaseLabel, () ->
                run(testSuite, testcase)
            );
            dynamicTests.add(test);
        }
        return dynamicTests;
    }

    private ProtocolTestsuite parseTestsuite(InputStream testsuiteDocumentXml) throws ProtocolTestsuiteException {
        try {
            SAXReader reader = new SAXReader();
            Document document = reader.read(testsuiteDocumentXml);
            Element testsuiteXml = document.getRootElement();
            boolean littleEndian = !"true".equals(testsuiteXml.attributeValue("bigEndian"));
            Element testsuiteName = testsuiteXml.element(new QName("name"));
            List<Element> testcasesXml = testsuiteXml.elements(new QName("testcase"));
            List<Testcase> testcases = new ArrayList<>(testcasesXml.size());
            for(Element testcaseXml : testcasesXml) {
                Element nameElement = testcaseXml.element(new QName("name"));
                Element descriptionElement = testcaseXml.element(new QName("description"));
                Element rawElement = testcaseXml.element(new QName("raw"));
                Element rootTypeElement = testcaseXml.element(new QName("root-type"));
                Element xmlElement = testcaseXml.element(new QName("xml"));

                String name = nameElement.getTextTrim();
                String description = (descriptionElement != null) ? descriptionElement.getTextTrim() : null;
                byte[] raw = Hex.decodeHex(rawElement.getTextTrim());
                String rootType = rootTypeElement.getTextTrim();

                testcases.add(new Testcase(name, description, raw, rootType, xmlElement));
            }
            LOGGER.info(String.format("Found %d testcases.", testcases.size()));
            return new ProtocolTestsuite(testsuiteName.getTextTrim(), testcases, littleEndian);
        } catch (DocumentException e) {
            throw new ProtocolTestsuiteException("Error parsing testsuite xml", e);
        } catch (DecoderException e) {
            throw new ProtocolTestsuiteException("Error parsing testcase raw data", e);
        }
    }

    private void run(ProtocolTestsuite testSuite, Testcase testcase) throws ProtocolTestsuiteException {
        ObjectMapper mapper = new XmlMapper().enableDefaultTyping();
        ReadBuffer readBuffer = new ReadBuffer(testcase.getRaw(), testSuite.isLittleEndian());
        String referenceXml = testcase.getXml().elements().get(0).asXML();

        MessageIO messageIO = getMessageIOForTestcase(testcase);
        try {
            Object msg = messageIO.parse(readBuffer);
            String xmlString = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(msg);
            Diff diff = DiffBuilder.compare(referenceXml).withTest(xmlString).ignoreWhitespace().build();
            if(diff.hasDifferences()) {
                System.out.println(xmlString);
                throw new ProtocolTestsuiteException("Differences were found after parsing.\n" + diff.toString());
            }
            WriteBuffer writeBuffer = new WriteBuffer(((Message) msg).getLengthInBytes(), testSuite.isLittleEndian());
            messageIO.serialize(writeBuffer, msg);
            byte[] data = writeBuffer.getData();
            if(testcase.getRaw().length != data.length) {
                LOGGER.info("Expected a byte array with a length of " + testcase.getRaw().length +
                    " but got one with " + data.length);
            }
            if(!Arrays.equals(testcase.getRaw(), data)) {
                int i;
                for(i = 0; i < data.length; i++) {
                    if(data[i] != testcase.getRaw()[i]) {
                        break;
                    }
                }
                throw new ProtocolTestsuiteException("Differences were found after serializing.\nExpected: " +
                    Hex.encodeHexString(testcase.getRaw()) + "\nBut Got:  " + Hex.encodeHexString(data) +
                    "\n          " + String.join("", Collections.nCopies(i, "--")) + "^");
            }
        } catch (ParseException e) {
            throw new ProtocolTestsuiteException("Unable to parse message", e);
        } catch (JsonProcessingException e) {
            throw new ProtocolTestsuiteException("Unable to serialize parsed message as XML string", e);
        }
    }

    private MessageIO getMessageIOForTestcase(Testcase testcase) throws ProtocolTestsuiteException {
        String className = testcase.getXml().elements().get(0).attributeValue(new QName("className"));
        String ioRootClassName = className.substring(0, className.lastIndexOf('.') + 1) + testcase.getRootType();
        String ioClassName = className.substring(0, className.lastIndexOf('.') + 1) + "io." +
            testcase.getRootType() + "IO";
        try {
            Class<?> ioRootClass = Class.forName(ioRootClassName);
            Class<?> ioClass = Class.forName(ioClassName);
            final Method parseMethod = ioClass.getMethod("parse", ReadBuffer.class);
            final Method serializeMethod = ioClass.getMethod("serialize", WriteBuffer.class, ioRootClass);
            return new MessageIO() {
                @Override
                public Object parse(ReadBuffer io) throws ParseException {
                    try {
                        return parseMethod.invoke(null, io);
                    } catch (IllegalAccessException | InvocationTargetException e) {
                        throw new ParseException("error parsing", e);
                    }
                }

                @Override
                public void serialize(WriteBuffer io, Object value) throws ParseException {
                    try {
                        serializeMethod.invoke(null, io, value);
                    } catch (IllegalAccessException | InvocationTargetException e) {
                        throw new ParseException("error serializing", e);
                    }
                }
            };
        } catch (NoSuchMethodException | ClassNotFoundException e) {
            throw new ProtocolTestsuiteException("Unable to instantiate IO component", e);
        }
    }

}
