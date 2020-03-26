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
package org.apache.plc4x.test.driver;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import io.netty.buffer.ByteBuf;
import io.netty.buffer.Unpooled;
import io.netty.channel.Channel;
import io.netty.channel.embedded.Plc4xEmbeddedChannel;
import org.apache.plc4x.java.PlcDriverManager;
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.messages.PlcReadRequest;
import org.apache.plc4x.java.api.messages.PlcResponse;
import org.apache.plc4x.java.spi.connection.ChannelExposingConnection;
import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.generation.*;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.model.DriverTestsuite;
import org.apache.plc4x.test.driver.model.StepType;
import org.apache.plc4x.test.driver.model.TestStep;
import org.apache.plc4x.test.driver.model.Testcase;
import org.apache.plc4x.test.driver.model.api.TestField;
import org.apache.plc4x.test.driver.model.api.TestReadRequest;
import org.apache.plc4x.test.driver.model.api.TestRequest;
import org.apache.plc4x.test.driver.model.api.TestWriteRequest;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.QName;
import org.dom4j.io.SAXReader;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.Diff;

import java.io.IOException;
import java.io.InputStream;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

public class DriverTestsuiteRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(DriverTestsuiteRunner.class);

    private final String testsuiteDocument;

    private CompletableFuture<? extends PlcResponse> responseFuture;

    public DriverTestsuiteRunner(String testsuiteDocument) {
        this.testsuiteDocument = testsuiteDocument;
    }

    @TestFactory
    public Iterable<DynamicTest> getTestsuiteTests() throws DriverTestsuiteException {
        DriverTestsuite testSuite = parseTestsuite(DriverTestsuiteRunner.class.getResourceAsStream(testsuiteDocument));
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

    private DriverTestsuite parseTestsuite(InputStream testsuiteDocumentXml) throws DriverTestsuiteException {
        try {
            SAXReader reader = new SAXReader();
            Document document = reader.read(testsuiteDocumentXml);
            Element testsuiteXml = document.getRootElement();
            boolean bigEndian = !"false".equals(testsuiteXml.attributeValue("bigEndian"));
            String testsuiteName = testsuiteXml.element(new QName("name")).getTextTrim();
            String driverName = testsuiteXml.element(new QName("driver-name")).getTextTrim();
            Element driverParametersElement = testsuiteXml.element(new QName("driver-parameters"));
            Map<String, String> driverParameters = null;
            if(driverParametersElement != null) {
                driverParameters = new HashMap<>();
                for (Element parameter : driverParametersElement.elements(new QName("parameter"))) {
                    String parameterName = parameter.element(new QName("name")).getTextTrim();
                    String parameterValue = parameter.element(new QName("value")).getTextTrim();
                    driverParameters.put(parameterName, parameterValue);
                }
            }

            List<TestStep> setupSteps = new LinkedList<>();
            if(testsuiteXml.element(new QName("setup")) != null) {
                Element setupElement = testsuiteXml.element(new QName("setup"));
                for (Element element : setupElement.elements()) {
                    setupSteps.add(parseTestStep(element));
                }
            }

            List<TestStep> teardownSteps = new LinkedList<>();
            if(testsuiteXml.element(new QName("teardown")) != null) {
                Element teardownElement = testsuiteXml.element(new QName("teardown"));
                for (Element element : teardownElement.elements()) {
                    setupSteps.add(parseTestStep(element));
                }
            }

            List<Element> testcasesXml = testsuiteXml.elements(new QName("testcase"));
            List<Testcase> testcases = new ArrayList<>(testcasesXml.size());
            for(Element testcaseXml : testcasesXml) {
                Element nameElement = testcaseXml.element(new QName("name"));
                Element descriptionElement = testcaseXml.element(new QName("description"));
                Element stepsElement = testcaseXml.element(new QName("steps"));

                String name = nameElement.getTextTrim();
                String description = (descriptionElement != null) ? descriptionElement.getTextTrim() : null;

                List<TestStep> steps = new LinkedList<>();
                for (Element element : stepsElement.elements()) {
                    steps.add(parseTestStep(element));
                }
                testcases.add(new Testcase(name, description, steps));
            }
            LOGGER.info(String.format("Found %d testcases.", testcases.size()));

            // Force the driver to not wait for the connection before returning the connection.
            System.setProperty(GeneratedDriverBase.PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE, "false");

            PlcConnection connection = getConnection(driverName, driverParameters);

            TimeUnit.MILLISECONDS.sleep(200);

            return new DriverTestsuite(testsuiteName, connection, setupSteps, teardownSteps, testcases, bigEndian);
        } catch (DocumentException e) {
            throw new DriverTestsuiteException("Error parsing testsuite xml", e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new DriverTestsuiteException("Interruption setting up testsuite xml", e);
        }
    }

    private PlcConnection getConnection(String driverName, Map<String, String> driverParameters)
        throws DriverTestsuiteException {
        try {
            StringBuilder sb = new StringBuilder();
            if(driverParameters != null) {
                for (Map.Entry<String, String> parameter : driverParameters.entrySet()) {
                    sb.append("&").append(parameter.getKey()).append("=").append(parameter.getValue());
                }
            }
            if(sb.length() > 0) {
                sb.replace(0, 1, "?");
            }
            return new PlcDriverManager().getConnection(driverName + ":test://hurz" + sb.toString());
        } catch (PlcConnectionException e) {
            throw new DriverTestsuiteException("Error loading driver", e);
        }
    }

    private void run(DriverTestsuite testsuite, Testcase testcase) throws DriverTestsuiteException {
        final PlcConnection plcConnection = testsuite.getConnection();
        final Plc4xEmbeddedChannel embeddedChannel = getEmbeddedChannel(testsuite);
        final boolean bigEndian = testsuite.isBigEndian();
        if(!testsuite.getSetupSteps().isEmpty()) {
            LOGGER.info("Running setup steps");
            for (TestStep setupStep : testsuite.getSetupSteps()) {
                executeStep(setupStep, plcConnection, embeddedChannel, bigEndian);
            }
            LOGGER.info("Finished setup steps");
        }
        LOGGER.info("Running test steps");
        for (TestStep step : testcase.getSteps()) {
            executeStep(step, plcConnection, embeddedChannel, bigEndian);
        }
        LOGGER.info("Finished test steps");
        if(!testsuite.getTeardownSteps().isEmpty()) {
            LOGGER.info("Running teardown steps");
            for (TestStep teardownStep : testsuite.getTeardownSteps()) {
                executeStep(teardownStep, plcConnection, embeddedChannel, bigEndian);
            }
            LOGGER.info("Finished teardown steps");
        }
    }

    private void executeStep(TestStep testStep, PlcConnection plcConnection, Plc4xEmbeddedChannel embeddedChannel, boolean bigEndian) throws DriverTestsuiteException {
        LOGGER.info(String.format("  - Running step: '%s' - %s", testStep.getName(), testStep.getType()));
        final ObjectMapper mapper = new XmlMapper().enableDefaultTyping();
        final Element payload = testStep.getPayload();
        try {
            switch (testStep.getType()) {
                case OUTGOING_PLC_BYTES: {
                    // As we're in the asynchronous world, give the driver some time to respond.
                    shortDelay();
                    // Prepare a ByteBuf that contains the data which would have been sent to the PLC.
                    final byte[] data = getOutboundBytes(embeddedChannel);
                    // Validate the data actually matches the expected message.
                    validateBytes(payload, data, bigEndian);
                    break;
                }
                case OUTGOING_PLC_MESSAGE: {
                    // As we're in the asynchronous world, give the driver some time to respond.
                    shortDelay();
                    // Prepare a ByteBuf that contains the data which would have been sent to the PLC.
                    final byte[] data = getOutboundBytes(embeddedChannel);
                    // Validate the data actually matches the expected message.
                    validateMessage(payload, data, bigEndian);
                    break;
                }
                case INCOMING_PLC_BYTES:
                case INCOMING_PLC_MESSAGE: {
                    // Get a byte representation of the incoming message.
                    final byte[] data = getBytesFromXml(payload, bigEndian);
                    // Send the bytes to the channel.
                    embeddedChannel.writeInbound(Unpooled.wrappedBuffer(data));
                    break;
                }
                case API_REQUEST: {
                    final String referenceXml = payload.asXML();
                    final TestRequest request = mapper.readValue(referenceXml, TestRequest.class);
                    if(request instanceof TestReadRequest) {
                        final TestReadRequest readRequest = (TestReadRequest) request;
                        final PlcReadRequest.Builder builder = plcConnection.readRequestBuilder();
                        for (TestField testField : readRequest.getFields()) {
                            builder.addItem(testField.getName(), testField.getAddress());
                        }
                        final PlcReadRequest plc4xRequest = builder.build();
                        // Currently we can only process one response at at time, throw an error if more
                        // are submitted.
                        if(responseFuture != null) {
                            throw new DriverTestsuiteException("Previous response not handled.");
                        }
                        // Save the response for being used later on.
                        responseFuture = plc4xRequest.execute();
                    } else if(request instanceof TestWriteRequest) {
                        // TODO: Implement ...
                    }
                    break;
                }
                case API_RESPONSE: {
                    if(responseFuture == null) {
                        throw new DriverTestsuiteException("No response expected.");
                    }
                    try {
                        final PlcResponse plcResponse = responseFuture.get(1000, TimeUnit.MILLISECONDS);
                        final String serializedResponse = mapper.writeValueAsString(plcResponse);
                        // TODO: Implement ...
                    } catch(Exception e) {
                        throw new DriverTestsuiteException("Got no response within 1000ms.");
                    }
                    break;
                }
                case DELAY: {
                    delay(1000);
                    break;
                }
                case TERMINATE: {
                    embeddedChannel.close();
                }
            }
        } catch (IOException e) {
            throw new DriverTestsuiteException("Error processing the xml", e);
        }
        LOGGER.info("    Done");
    }

    private TestStep parseTestStep(Element curElement) {
        final String elementName = curElement.getName();
        final StepType stepType = StepType.valueOf(elementName.toUpperCase().replace("-", "_"));
        final String stepName = curElement.attributeValue(new QName("name"));
        Element definition = null;
        if(curElement.hasMixedContent()) {
            definition = curElement.elementIterator().next();
        }
        return new TestStep(stepType, stepName, definition);
    }

    private Plc4xEmbeddedChannel getEmbeddedChannel(DriverTestsuite testSuite) {
        if(testSuite.getConnection() instanceof ChannelExposingConnection) {
            ChannelExposingConnection connection = (ChannelExposingConnection) testSuite.getConnection();
            Channel channel = connection.getChannel();
            if(channel instanceof Plc4xEmbeddedChannel) {
                return (Plc4xEmbeddedChannel) channel;
            }
            throw new PlcRuntimeException("Expecting EmbeddedChannel");
        }
        throw new PlcRuntimeException("Expecting ChannelExposingConnection");
    }

    @SuppressWarnings("unchecked")
    private Class<? extends Message> getMessageType(String messageClassName) throws DriverTestsuiteException {
        try {
            final Class<?> messageClass = Class.forName(messageClassName);
            if(Message.class.isAssignableFrom(messageClass)) {
                return (Class<? extends Message>) messageClass;
            }
            throw new DriverTestsuiteException("IO class must implement Message interface");
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException("Error loading message class", e);
        }
    }

    @SuppressWarnings("unchecked")
    private Class<? extends MessageIO<?,?>> getMessageIOType(String messageClassName) throws DriverTestsuiteException {
        String ioClassName = messageClassName.substring(0, messageClassName.lastIndexOf('.')) + ".io." +
            messageClassName.substring(messageClassName.lastIndexOf('.') + 1) + "IO";
        try {
            final Class<?> ioClass = Class.forName(ioClassName);
            if(MessageIO.class.isAssignableFrom(ioClass)) {
                return (Class<? extends MessageIO<?,?>>) ioClass;
            }
            throw new DriverTestsuiteException("IO class muss implement MessageIO interface");
        } catch (ClassNotFoundException e) {
            throw new DriverTestsuiteException("Error loading io class", e);
        }
    }

    private byte[] getOutboundBytes(Plc4xEmbeddedChannel embeddedChannel) throws DriverTestsuiteException {
        ByteBuf byteBuf = null;
        for(int i = 0; i < 10; i++) {
            byteBuf = embeddedChannel.readOutbound();
            if(byteBuf != null) {
                break;
            }
            delay(10);
        }
        if(byteBuf == null) {
            throw new DriverTestsuiteException("No outbound message available within 100ms");
        }
        final byte[] data = new byte[byteBuf.readableBytes()];
        byteBuf.readBytes(data);
        return data;
    }

    private byte[] getBytesFromXml(Element referenceXml, boolean bigEndian) throws DriverTestsuiteException {
        final String className = referenceXml.attributeValue(new QName("className"));
        final WriteBuffer writeBuffer = new WriteBuffer(1024, !bigEndian);
        try {
            final MessageIO messageIO = getMessageIOType(className).newInstance();
            final String referenceXmlString = referenceXml.asXML();
            final ObjectMapper mapper = new XmlMapper().enableDefaultTyping();
            final Message message = mapper.readValue(referenceXmlString, getMessageType(className));
            try {
                messageIO.serialize(writeBuffer, message);
                final byte[] data = new byte[message.getLengthInBytes()];
                System.arraycopy(writeBuffer.getData(), 0, data, 0, writeBuffer.getPos());
                return data;
            } catch (ParseException e) {
                throw new DriverTestsuiteException("Error serializing message", e);
            }
        } catch (IllegalAccessException | JsonProcessingException | InstantiationException e) {
            throw new DriverTestsuiteException("Error parsing message", e);
        }
    }

    private void validateBytes(Element referenceXml, byte[] data, boolean bigEndian) throws DriverTestsuiteException {
        // TODO: Implement this ...
    }

    private void validateMessage(Element referenceXml, byte[] data, boolean bigEndian) throws DriverTestsuiteException {
        final ObjectMapper mapper = new XmlMapper().enableDefaultTyping();
        final ReadBuffer readBuffer = new ReadBuffer(data, !bigEndian);
        try {
            final String className = referenceXml.attributeValue(new QName("className"));
            final MessageIO<?,?> messageIO = getMessageIOType(className).newInstance();
            final Object parsedOutput = messageIO.parse(readBuffer);
            final String xmlString = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(parsedOutput);
            final String referenceXmlString = referenceXml.asXML();
            final Diff diff = DiffBuilder.compare(referenceXmlString).withTest(xmlString).ignoreComments().ignoreWhitespace().build();
            if (diff.hasDifferences()) {
                LOGGER.warn(xmlString);
                throw new DriverTestsuiteException("Differences were found after parsing.\n" + diff.toString());
            }
        } catch (ParseException | IllegalAccessException | JsonProcessingException | InstantiationException e) {
            throw new DriverTestsuiteException("Error parsing message", e);
        }
    }

    private void shortDelay() throws DriverTestsuiteException {
        delay(23);
    }

    private void delay(int milliseconds) throws DriverTestsuiteException {
        try {
            TimeUnit.MILLISECONDS.sleep(milliseconds);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new DriverTestsuiteException("Interrupted during delay.");
        }
    }

}
