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

package org.apache.plc4x.test.parserserializer;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.apache.commons.lang3.StringUtils;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.test.XmlTestsuiteRunner;
import org.apache.plc4x.test.dom4j.LocationAwareDocumentFactory;
import org.apache.plc4x.test.dom4j.LocationAwareElement;
import org.apache.plc4x.test.dom4j.LocationAwareSAXReader;
import org.apache.plc4x.test.mapper.TestSuiteMappingModule;
import org.apache.plc4x.test.parserserializer.exceptions.ParserSerializerTestsuiteException;
import org.apache.plc4x.test.parserserializer.model.ParserSerializerTestsuite;
import org.apache.plc4x.test.parserserializer.model.Testcase;
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
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.Diff;

import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;
import java.lang.reflect.Modifier;
import java.net.URISyntaxException;
import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;

@SuppressWarnings({"unchecked", "rawtypes"})
public class ParserSerializerTestsuiteRunner extends XmlTestsuiteRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(ParserSerializerTestsuiteRunner.class);

    /**
     * set to true during testcase migration
     */
    private final boolean failOnUnMigrated;

    /**
     * if set to true if will automigrate and on the next run test should be green
     */
    private final boolean autoMigrate;

    private final Set<String> ignoredTestCases = new HashSet<>();

    public ParserSerializerTestsuiteRunner(String testsuiteDocument, String... ignoredTestCases) {
        this(testsuiteDocument, false, ignoredTestCases);
    }

    public ParserSerializerTestsuiteRunner(String testsuiteDocument, boolean failOnUnMigrated, String... ignoredTestCases) {
        this(testsuiteDocument, false, false, ignoredTestCases);
    }

    public ParserSerializerTestsuiteRunner(String testsuiteDocument, boolean failOnUnMigrated, boolean autoMigrate, String... ignoredTestCases) {
        super(testsuiteDocument);
        this.failOnUnMigrated = failOnUnMigrated;
        this.autoMigrate = autoMigrate;
        Collections.addAll(this.ignoredTestCases, ignoredTestCases);
    }

    @TestFactory
    public Iterable<DynamicTest> getTestsuiteTests() throws ParserSerializerTestsuiteException, URISyntaxException {
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
            boolean littleEndian = !"true".equals(testsuiteXml.attributeValue("bigEndian"));
            Element testsuiteName = testsuiteXml.element(new QName("name"));
            Element protocolName = testsuiteXml.element(new QName("protocolName"));
            Element outputFlavor = testsuiteXml.element(new QName("outputFlavor"));
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
                Testcase testcase = new Testcase(testsuiteName.getStringValue(), protocolName.getStringValue(), outputFlavor.getStringValue(), name, description, raw, rootType, parserArguments, xmlElement);
                if (testcaseXml instanceof LocationAwareElement) {
                    // pass source location to test
                    testcase.setLocation(((LocationAwareElement) testcaseXml).getLocation());
                }
                testcases.add(testcase);
            }
            LOGGER.info(String.format("Found %d testcases.", testcases.size()));
            return new ParserSerializerTestsuite(testsuiteName.getTextTrim(), testcases, littleEndian);
        } catch (DocumentException e) {
            throw new ParserSerializerTestsuiteException("Error parsing testsuite xml", e);
        } catch (DecoderException e) {
            throw new ParserSerializerTestsuiteException("Error parsing testcase raw data", e);
        }
    }

    private void run(ParserSerializerTestsuite testSuite, Testcase testcase) throws ParserSerializerTestsuiteException {
        ReadBufferByteBased readBuffer = new ReadBufferByteBased(testcase.getRaw(), testSuite.isLittleEndian());
        String referenceXml = testcase.getXml().elements().get(0).asXML();

        MessageIO messageIO = getMessageIOForTestcase(testcase);
        try {
            Object msg = messageIO.parse(readBuffer, testcase.getParserArguments().toArray());
            {
                try {
                    // First try to use the native xml writer
                    WriteBufferXmlBased writerBufferXmlBased = new WriteBufferXmlBased();
                    messageIO.serialize(writerBufferXmlBased, msg);
                    String xmlString = writerBufferXmlBased.getXmlString();
                    Diff diff = DiffBuilder.compare(referenceXml).withTest(xmlString).ignoreComments().ignoreWhitespace().build();
                    if (diff.hasDifferences()) {
                        String border = StringUtils.repeat("=", 100);
                        String centeredDiffDetectedMessage = StringUtils.center(" Diff detected ", 100, "=");
                        String centeredTestCaseName = StringUtils.center(testcase.getName(), 100, "=");
                        System.err.printf(
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
                            centeredTestCaseName);

                        throw new MigrationException(xmlString);
                    }
                } catch (RuntimeException e) {
                    if (!(e instanceof MigrationException)) {
                        System.err.println("Error in serializer");
                        System.err.println(e.getMessage());
                        e.printStackTrace();
                    }
                    ObjectMapper mapper = new XmlMapper().enableDefaultTyping().registerModule(new TestSuiteMappingModule());
                    String xmlStringFallback = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(msg);
                    Diff diff2 = DiffBuilder.compare(referenceXml).withTest(xmlStringFallback).ignoreComments().ignoreWhitespace().build();
                    if (diff2.hasDifferences()) {
                        System.out.println(xmlStringFallback);
                        throw new ParserSerializerTestsuiteException("Differences were found after parsing.\n" + diff2);
                    } else {
                        System.out.println("No diff detected with old");
                    }
                    if (autoMigrate && e instanceof MigrationException) {
                        Path path = Paths.get(suiteUri);
                        System.out.printf("Migrating %s now", path);
                        Charset charset = StandardCharsets.UTF_8;

                        String content;
                        try {
                            content = new String(Files.readAllBytes(path), charset);
                        } catch (IOException ioException) {
                            throw new RuntimeException(ioException);
                        }
                        // We need to indent the search string properly
                        String indent = StringUtils.repeat(' ', 6);
                        String searchString = StringUtils.replace(xmlStringFallback, "\n", "\n" + indent);
                        searchString = StringUtils.trim(searchString);
                        String newXml = ((MigrationException) e).newXml;
                        newXml = StringUtils.replace(newXml, "\n", "\n" + indent);
                        // Remove last wrong indent
                        newXml = newXml.substring(0, newXml.length() - 7);
                        content = StringUtils.replaceOnce(content, searchString, newXml);
                        try {
                            Files.write(path, content.getBytes(charset));
                        } catch (IOException ioException) {
                            throw new RuntimeException(ioException);
                        }
                        System.out.printf("Done migrating %s", path);
                    }
                    if (failOnUnMigrated) {
                        throw new RuntimeException("fail on un-migrated set to true. Please migrate testcase", e);
                    }
                }
            }

            WriteBufferByteBased writeBuffer = new WriteBufferByteBased(((Message) msg).getLengthInBytes(), testSuite.isLittleEndian());
            messageIO.serialize(writeBuffer, msg);
            byte[] data = writeBuffer.getData();
            if (testcase.getRaw().length != data.length) {
                LOGGER.info("Expected a byte array with a length of " + testcase.getRaw().length +
                    " but got one with " + data.length);
            }
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
        } catch (ParseException e) {
            throw new ParserSerializerTestsuiteException("Unable to parse message", e);
        } catch (JsonProcessingException e) {
            throw new ParserSerializerTestsuiteException("Unable to serialize parsed message as XML string", e);
        }
    }

    private MessageIO getMessageIOForTestcase(Testcase testcase) throws ParserSerializerTestsuiteException {
        String ioClassName, ioRootClassName;
        try {
            String classPackage = String.format("org.apache.plc4x.java.%s.%s", testcase.getProtocolName(), StringUtils.replace(testcase.getOutputFlavor(), "-", ""));
            try {
                Package.getPackage(classPackage);
            } catch (RuntimeException e) {
                System.err.println("Error resolving package for " + classPackage);
                switch (testcase.getTestSuiteName()) {
                    case "Firmata":
                        classPackage = "org.apache.plc4x.java.firmata.readwrite";
                        break;
                    case "Allen-Bradley ETH":
                        classPackage = "org.apache.plc4x.java.abeth.readwrite";
                        break;
                    case "Beckhoff ADS/AMS Discovery":
                        classPackage = "org.apache.plc4x.java.ads.discovery.readwrite";
                        break;
                    case "Beckhoff ADS/AMS":
                        classPackage = "org.apache.plc4x.java.ads.readwrite";
                        break;
                    case "Tests of CANopen frames payload":
                        classPackage = "org.apache.plc4x.java.canopen.readwrite";
                        break;
                    case "Tests of CANopen frames from Wireshark sample PCAP files":
                        classPackage = "org.apache.plc4x.java.canopen.readwrite";
                        break;
                    case "EIP":
                        classPackage = "org.apache.plc4x.java.eip.readwrite";
                        break;
                    case "Modbus":
                        classPackage = "org.apache.plc4x.java.modbus.readwrite";
                        break;
                    case "S7":
                        classPackage = "org.apache.plc4x.java.s7.readwrite";
                        break;
                    case "KNXNet/IP":
                        classPackage = "org.apache.plc4x.java.knxnetip.readwrite";
                        break;
                    default:
                        throw new RuntimeException(String.format("fallback to old. No packageName for '%s' configured.\nAdd the required package to the switch clause above", testcase.getTestSuiteName()));
                }
            }
            String fullQualifiedClassName = classPackage + "." + testcase.getXml().elements().get(0).getName();
            ioRootClassName = fullQualifiedClassName.substring(0, fullQualifiedClassName.lastIndexOf('.') + 1) + testcase.getRootType();
            ioClassName = fullQualifiedClassName.substring(0, fullQualifiedClassName.lastIndexOf('.') + 1) + "io." +
                testcase.getRootType() + "IO";
            try {
                Class.forName(ioRootClassName);
                Class.forName(ioClassName);
            } catch (ClassNotFoundException e) {
                throw new RuntimeException("fallback to old", e);
            }
        } catch (RuntimeException e) {
            System.err.println("Error in serializer");
            System.err.println(e.getMessage());
            e.printStackTrace();
            String fullQualifiedClassName = testcase.getXml().elements().get(0).attributeValue(new QName("className"));
            ioRootClassName = fullQualifiedClassName.substring(0, fullQualifiedClassName.lastIndexOf('.') + 1) + testcase.getRootType();
            ioClassName = fullQualifiedClassName.substring(0, fullQualifiedClassName.lastIndexOf('.') + 1) + "io." +
                testcase.getRootType() + "IO";
        }
        try {
            Class<?> ioRootClass = Class.forName(ioRootClassName);
            Class<?> ioClass = Class.forName(ioClassName);
            Method parseMethodTmp = null;
            Method serializeMethodTmp = null;
            final List<Class<?>> parameterTypes = new LinkedList<>();
            for (Method method : ioClass.getMethods()) {
                if (method.getName().equals("staticParse") && Modifier.isStatic(method.getModifiers()) &&
                    (method.getReturnType() == ioRootClass)) {
                    parseMethodTmp = method;

                    // Get a list of additional parameter types for the parser.
                    for (int i = 1; i < method.getParameterCount(); i++) {
                        Class<?> parameterType = parseMethodTmp.getParameterTypes()[i];
                        parameterTypes.add(parameterType);
                    }
                }
                if (method.getName().equals("staticSerialize") && Modifier.isStatic(method.getModifiers()) &&
                    (method.getParameterTypes()[1] == ioRootClass)) {
                    serializeMethodTmp = method;
                }
            }
            if ((parseMethodTmp == null) || (serializeMethodTmp == null)) {
                throw new ParserSerializerTestsuiteException(
                    "Unable to instantiate IO component. Missing static parse or serialize methods.");
            }
            final Method parseMethod = parseMethodTmp;
            final Method serializeMethod = serializeMethodTmp;
            return new MessageIO() {
                @Override
                public Object parse(ReadBuffer io, Object... args) throws ParseException {
                    try {
                        Object[] argValues = new Object[args.length + 1];
                        argValues[0] = io;
                        for (int i = 1; i <= args.length; i++) {
                            String parameterValue = (String) args[i - 1];
                            Class<?> parameterType = parameterTypes.get(i - 1);
                            if (parameterType == Boolean.class) {
                                argValues[i] = Boolean.parseBoolean(parameterValue);
                            } else if (parameterType == Byte.class) {
                                argValues[i] = Byte.parseByte(parameterValue);
                            } else if (parameterType == Short.class) {
                                argValues[i] = Short.parseShort(parameterValue);
                            } else if (parameterType == Integer.class) {
                                argValues[i] = Integer.parseInt(parameterValue);
                            } else if (parameterType == Long.class) {
                                argValues[i] = Long.parseLong(parameterValue);
                            } else if (parameterType == Float.class) {
                                argValues[i] = Float.parseFloat(parameterValue);
                            } else if (parameterType == Double.class) {
                                argValues[i] = Double.parseDouble(parameterValue);
                            } else if (parameterType == String.class) {
                                argValues[i] = parameterValue;
                            } else if (Enum.class.isAssignableFrom(parameterType)) {
                                argValues[i] = Enum.valueOf((Class<? extends Enum>) parameterType, parameterValue);
                            } else {
                                throw new ParseException("Currently unsupported parameter type");
                            }
                        }

                        return parseMethod.invoke(null, argValues);
                    } catch (IllegalAccessException | InvocationTargetException e) {
                        throw new ParseException("error parsing", e);
                    }
                }

                @Override
                public void serialize(WriteBuffer io, Object value, Object... args) throws ParseException {
                    try {
                        serializeMethod.invoke(null, io, value);
                    } catch (IllegalAccessException | InvocationTargetException e) {
                        throw new ParseException("error serializing", e);
                    }
                }
            };
        } catch (ClassNotFoundException e) {
            throw new ParserSerializerTestsuiteException("Unable to instantiate IO component", e);
        }
    }

    private static class MigrationException extends RuntimeException {
        final String newXml;

        public MigrationException(String newXml) {
            this.newXml = newXml;
        }
    }

}
