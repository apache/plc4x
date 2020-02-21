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

import com.fasterxml.jackson.dataformat.xml.XmlMapper;
import org.apache.plc4x.java.spi.generation.Message;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.model.DriverTestsuite;
import org.apache.plc4x.test.driver.model.TestStep;
import org.apache.plc4x.test.driver.model.Testcase;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.QName;
import org.dom4j.io.SAXReader;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.xml.stream.XMLStreamException;
import java.io.IOException;
import java.io.InputStream;
import java.util.*;

public class DriverTestsuiteRunner {

    private static final Logger LOGGER = LoggerFactory.getLogger(DriverTestsuiteRunner.class);

    private final String testsuiteDocument;

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
            boolean littleEndian = !"true".equals(testsuiteXml.attributeValue("bigEndian"));
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
                setupSteps.add(new TestStep(setupElement));
            }

            List<TestStep> teardownSteps = new LinkedList<>();
            if(testsuiteXml.element(new QName("teardown")) != null) {
                Element teardownElement = testsuiteXml.element(new QName("teardown"));
                setupSteps.add(new TestStep(teardownElement));
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
                    steps.add(new TestStep(element));
                }
                testcases.add(new Testcase(name, description, steps));
            }
            LOGGER.info(String.format("Found %d testcases.", testcases.size()));
            return new DriverTestsuite(
                testsuiteName, driverName, driverParameters, setupSteps, teardownSteps, testcases);
        } catch (DocumentException e) {
            throw new DriverTestsuiteException("Error parsing testsuite xml", e);
        }
    }

    private void run(DriverTestsuite testSuite, Testcase testcase) throws DriverTestsuiteException {
        XmlMapper xmlMapper = new XmlMapper();
        for (TestStep step : testcase.getSteps()) {
            String referenceXml = step.getPayload().asXML();
            try {
                final Message message = xmlMapper.readValue(referenceXml, Message.class);
                System.out.println(message);
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }

}
