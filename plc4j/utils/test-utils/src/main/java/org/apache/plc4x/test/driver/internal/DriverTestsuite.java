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
package org.apache.plc4x.test.driver.internal;

import static org.apache.plc4x.test.xml.XmlHelper.*;

import org.apache.plc4x.java.spi.connection.GeneratedDriverBase;
import org.apache.plc4x.java.spi.generation.ByteOrder;
import org.apache.plc4x.test.dom4j.LocationAwareDocumentFactory;
import org.apache.plc4x.test.dom4j.LocationAwareElement;
import org.apache.plc4x.test.dom4j.LocationAwareSAXReader;
import org.apache.plc4x.test.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.test.driver.internal.utils.Delay;
import org.apache.plc4x.test.driver.internal.utils.Synchronizer;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.Element;
import org.dom4j.QName;
import org.dom4j.io.SAXReader;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.InputStream;
import java.net.URI;
import java.util.*;
import java.util.stream.Collectors;

public class DriverTestsuite {

    private static final Logger LOGGER = LoggerFactory.getLogger(DriverTestsuite.class);

    private final DriverTestsuiteConfiguration driverTestsuiteConfiguration;

    private List<TestStep> setupSteps;
    private List<TestStep> teardownSteps;
    private List<Testcase> testcases;

    private DriverTestsuite(DriverTestsuiteConfiguration driverTestsuiteConfiguration) {
        this.driverTestsuiteConfiguration = driverTestsuiteConfiguration;
    }

    private void setSetupSteps(List<TestStep> setupSteps) {
        this.setupSteps = setupSteps;
    }

    private void setTeardownSteps(List<TestStep> teardownSteps) {
        this.teardownSteps = teardownSteps;
    }

    private void setTestcases(List<Testcase> testcases) {
        this.testcases = testcases;
    }

    public static DriverTestsuite parseTestsuite(URI suiteUri, InputStream testsuiteDocumentXml, boolean autoMigrate ) throws DriverTestsuiteException {
        try {
            SAXReader reader = new LocationAwareSAXReader();
            reader.setDocumentFactory(new LocationAwareDocumentFactory());
            Document document = reader.read(testsuiteDocumentXml);
            Element testsuiteXml = document.getRootElement();

            // Force the driver to not wait for the connection before returning the connection.
            System.setProperty(GeneratedDriverBase.PROPERTY_PLC4X_FORCE_AWAIT_SETUP_COMPLETE, "false");

            // Force the driver to not wait for the disconnection before returning closing the channel.
            System.setProperty(GeneratedDriverBase.PROPERTY_PLC4X_FORCE_AWAIT_DISCONNECT_COMPLETE, "false");

            // TODO: replace with signal
            Delay.delay(200);

            // Shared instance for synchronizing
            Synchronizer synchronizer = new Synchronizer();
            DriverTestsuiteBuilder driverTestsuiteBuilder = new DriverTestsuiteBuilder(suiteUri, testsuiteXml, synchronizer, autoMigrate);
            return driverTestsuiteBuilder.build();
        } catch (DocumentException e) {
            throw new DriverTestsuiteException("Error parsing testsuite xml", e);
        }
    }

    public DriverTestsuiteConfiguration getDriverTestsuiteConfiguration() {
        return driverTestsuiteConfiguration;
    }

    public List<Testcase> getTestcases() {
        return testcases;
    }

    public String getName() {
        return driverTestsuiteConfiguration.getTestSuiteName();
    }

    public List<TestStep> getSetupSteps() {
        return setupSteps;
    }

    public List<TestStep> getTeardownSteps() {
        return teardownSteps;
    }


    private static class DriverTestsuiteBuilder {

        private final URI suiteUri;

        private final Synchronizer synchronizer;

        private final Element testsuiteXml;

        private final String testsuiteName;

        private final String protocolName;

        private final String outputFlavor;

        private final String driverName;

        private final Element optionsElement;

        private final Element driverParametersElement;

        private final ByteOrder byteOrder;

        private final boolean autoMigrate;

        public DriverTestsuiteBuilder(URI suiteUri, Element testsuiteXml, Synchronizer synchronizer, boolean autoMigrate) {
            this.suiteUri = suiteUri;
            this.synchronizer = synchronizer;
            this.testsuiteXml = testsuiteXml;
            this.testsuiteName = extractText(testsuiteXml, "name");
            this.protocolName = extractText(testsuiteXml, "protocolName");
            this.outputFlavor = extractText(testsuiteXml, "outputFlavor");
            this.driverName = extractText(testsuiteXml, "driver-name");
            this.optionsElement = testsuiteXml.element("options");
            this.driverParametersElement = testsuiteXml.element(new QName("driver-parameters"));
            this.byteOrder = ByteOrder.valueOf(testsuiteXml.attributeValue("byteOrder"));
            this.autoMigrate = autoMigrate;
        }

        private DriverTestsuite build() {
            DriverTestsuiteConfiguration driverTestsuiteConfiguration = new DriverTestsuiteConfiguration(
                suiteUri,
                testsuiteName,
                protocolName,
                outputFlavor,
                driverName,
                parseParameters(optionsElement),
                parseParameters(driverParametersElement),
                autoMigrate,
                byteOrder
            );
            DriverTestsuite driverTestsuite = new DriverTestsuite(driverTestsuiteConfiguration);
            driverTestsuite.setTestcases(parseTestCases(driverTestsuite));
            driverTestsuite.setSetupSteps(parseSetupSteps(driverTestsuiteConfiguration));
            driverTestsuite.setTeardownSteps(parseTearDownSteps(driverTestsuiteConfiguration));
            return driverTestsuite;
        }

        private List<Testcase> parseTestCases(DriverTestsuite driverTestsuite) {
            List<Testcase> testcases = testsuiteXml.elements(new QName("testcase")).stream()
                .map(testcaseXml -> {
                    Testcase testcase = new Testcase(
                        driverTestsuite, extractText(testcaseXml, "name"),
                        parseDescription(testcaseXml),
                        parseSteps(testcaseXml, driverTestsuite.driverTestsuiteConfiguration),
                        synchronizer
                    );
                    if (testcaseXml instanceof LocationAwareElement) {
                        testcase.setLocation(((LocationAwareElement) testcaseXml).getLocation());
                    }
                    return testcase;
                })
                .collect(Collectors.toList());
            LOGGER.info("Found {} testcases.", testcases.size());
            return testcases;
        }

        private List<TestStep> parseSteps(Element testcaseXml, DriverTestsuiteConfiguration driverTestsuiteConfiguration) {
            return testcaseXml.element(new QName("steps")).elements().stream()
                .map(subElement -> TestStep.parseTestStep(subElement, synchronizer, driverTestsuiteConfiguration))
                .collect(Collectors.toList());
        }

        private String parseDescription(Element testcaseXml) {
            Element descriptionElement = testcaseXml.element(new QName("description"));
            return (descriptionElement != null) ? descriptionElement.getTextTrim() : null;
        }

        private List<TestStep> parseTearDownSteps(DriverTestsuiteConfiguration driverTestsuiteConfiguration) throws DriverTestsuiteException {
            return parseSteps("teardown", driverTestsuiteConfiguration);
        }

        private List<TestStep> parseSetupSteps(DriverTestsuiteConfiguration driverTestsuiteConfiguration) throws DriverTestsuiteException {
            return parseSteps("setup", driverTestsuiteConfiguration);
        }

        private List<TestStep> parseSteps(String type, DriverTestsuiteConfiguration driverTestsuiteConfiguration) throws DriverTestsuiteException {
            Element element = testsuiteXml.element(new QName(type));
            if (element == null) {
                return Collections.emptyList();
            }
            return element.elements().stream()
                .map(subElement -> TestStep.parseTestStep(subElement, synchronizer, driverTestsuiteConfiguration))
                .collect(Collectors.toCollection(LinkedList::new));
        }

    }

}
