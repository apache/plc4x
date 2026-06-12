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
package org.apache.plc4x.java.utils.testutils.driver.internal;

import org.apache.plc4x.java.spi.buffers.bytebased.byteorder.ByteOrderBigEndian;
import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.utils.Synchronizer;
import org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareSAXReader;
import org.apache.plc4x.java.utils.testutils.utils.model.Location;
import org.apache.plc4x.java.utils.testutils.utils.xml.XmlHelper;
import org.dom4j.Document;
import org.dom4j.Element;
import org.dom4j.QName;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.InputStream;
import java.net.URI;
import java.util.*;
import java.util.stream.Collectors;

/**
 * Represents a complete driver test suite parsed from XML.
 */
public class DriverTestsuite {

    private static final Logger LOGGER = LoggerFactory.getLogger(DriverTestsuite.class);

    private final DriverTestsuiteConfiguration configuration;
    private final List<TestStep> setupSteps;
    private final List<Testcase> testcases;
    private final List<TestStep> teardownSteps;
    private final Synchronizer synchronizer;

    private DriverTestsuite(DriverTestsuiteConfiguration configuration,
                           List<TestStep> setupSteps,
                           List<Testcase> testcases,
                           List<TestStep> teardownSteps,
                           Synchronizer synchronizer) {
        this.configuration = configuration;
        this.setupSteps = setupSteps;
        this.testcases = testcases;
        this.teardownSteps = teardownSteps;
        this.synchronizer = synchronizer;
    }

    /**
     * Parses a driver test suite from an XML input stream.
     *
     * @param suiteUri        the URI of the test suite
     * @param testsuiteXml    the XML input stream
     * @param autoMigrate     whether to auto-migrate old test formats
     * @return the parsed driver test suite
     * @throws DriverTestsuiteException if parsing fails
     */
    public static DriverTestsuite parseTestsuite(URI suiteUri, String basePackage, InputStream testsuiteXml, boolean autoMigrate) {
        try {
            // Parse XML with location awareness
            LocationAwareSAXReader reader = new LocationAwareSAXReader();
            Document document = reader.read(testsuiteXml);

            // Build the test suite
            return new DriverTestsuiteBuilder(suiteUri, basePackage, document, autoMigrate).build();
        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to parse test suite", e);
        }
    }

    public DriverTestsuiteConfiguration getConfiguration() {
        return configuration;
    }

    public List<TestStep> getSetupSteps() {
        return setupSteps;
    }

    public List<Testcase> getTestcases() {
        return testcases;
    }

    public List<TestStep> getTeardownSteps() {
        return teardownSteps;
    }

    public Synchronizer getSynchronizer() {
        return synchronizer;
    }

    /**
     * Builder for DriverTestsuite.
     */
    private static class DriverTestsuiteBuilder {

        private final URI suiteUri;
        private final String basePackage;
        private final Document document;
        private final boolean autoMigrate;
        private final Synchronizer synchronizer;

        public DriverTestsuiteBuilder(URI suiteUri, String basePackage, Document document, boolean autoMigrate) {
            this.suiteUri = suiteUri;
            this.basePackage = basePackage;
            this.document = document;
            this.autoMigrate = autoMigrate;
            this.synchronizer = new Synchronizer();

            // Disable driver-level await mechanisms for testing
            System.setProperty("org.apache.plc4x.disableAwaitSetupComplete", "true");
            System.setProperty("org.apache.plc4x.disableAwaitDisconnectComplete", "true");
        }

        public DriverTestsuite build() {
            Element root = document.getRootElement();

            // Parse configuration
            DriverTestsuiteConfiguration config = parseConfiguration(root);

            // Parse setup steps
            List<TestStep> setupSteps = parseSteps(root.element(new QName("setup")), config.getByteOrder());

            // Parse test cases
            List<Testcase> testcases = parseTestcases(root, config);

            // Parse teardown steps
            List<TestStep> teardownSteps = parseSteps(root.element(new QName("teardown")), config.getByteOrder());

            DriverTestsuite suite = new DriverTestsuite(config, setupSteps, testcases, teardownSteps, synchronizer);

            // Set testsuite reference on all testcases
            for (Testcase testcase : testcases) {
                testcase.setTestsuite(suite);
            }

            return suite;
        }

        private DriverTestsuiteConfiguration parseConfiguration(Element root) {
            String testsuiteName = root.attributeValue("name");
            if (testsuiteName == null) {
                testsuiteName = "Unnamed Test Suite";
            }

            String protocolName = XmlHelper.extractText(root, "protocol-name");
            String outputFlavor = extractOptionalText(root, "output-flavor");
            String driverName = XmlHelper.extractText(root, "driver-name");

            // Parse byte order. Prefer the root attribute form (which matches the
            // schema used by the EIP/Modbus testsuites — `<testsuite byteOrder="...">`
            // also used by the ParserSerializerTestsuiteRunner). Fall back to a
            // `<byte-order>` child element for older suites, then to BIG_ENDIAN
            // as the historical default.
            String byteOrder = ByteOrderBigEndian.NAME;
            String byteOrderAttr = root.attributeValue("byteOrder");
            if (byteOrderAttr != null) {
                byteOrder = byteOrderAttr.toUpperCase();
            } else {
                String byteOrderStr = extractOptionalText(root, "byte-order");
                if (byteOrderStr != null) {
                    byteOrder = byteOrderStr.toUpperCase();
                }
            }

            // Parse options
            Map<String, String> options = new HashMap<>();
            Element optionsElement = root.element(new QName("options"));
            if (optionsElement != null) {
                options = XmlHelper.parseParameters(optionsElement);
            }

            // Parse driver parameters
            Map<String, String> driverParameters = new HashMap<>();
            Element driverParametersElement = root.element(new QName("driver-parameters"));
            if (driverParametersElement != null) {
                driverParameters = XmlHelper.parseParameters(driverParametersElement);
            }

            // Parse sequential mode (run all test cases on a single connection)
            boolean sequential = false;
            String sequentialStr = extractOptionalText(root, "sequential");
            if ("true".equalsIgnoreCase(sequentialStr)) {
                sequential = true;
            }

            return new DriverTestsuiteConfiguration(
                suiteUri, testsuiteName, protocolName, outputFlavor, driverName,
                options, driverParameters, autoMigrate, sequential, byteOrder
            );
        }

        private List<Testcase> parseTestcases(Element root, DriverTestsuiteConfiguration config) {
            List<Testcase> testcases = new ArrayList<>();

            for (Element testcaseElement : root.elements(new QName("testcase"))) {
                testcases.add(parseTestcase(testcaseElement, config));
            }

            return testcases;
        }

        private Testcase parseTestcase(Element testcaseElement, DriverTestsuiteConfiguration config) {
            // Try to get name as child element first, then as attribute for backwards compatibility
            String name = extractOptionalText(testcaseElement, "name");
            if (name == null) {
                name = testcaseElement.attributeValue("name");
            }
            if (name == null) {
                name = "Unnamed Test Case";
            }

            String description = extractOptionalText(testcaseElement, "description");

            // Parse steps
            Element stepsElement = testcaseElement.element(new QName("steps"));
            List<TestStep> steps = parseSteps(stepsElement, config.getByteOrder());

            // Extract location
            Location location = null;
            if (testcaseElement instanceof org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareElement) {
                org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareElement locAware =
                    (org.apache.plc4x.java.utils.testutils.utils.dom4j.LocationAwareElement) testcaseElement;
                location = locAware.getLocation();
            }

            return new Testcase(name, description, steps, synchronizer, location);
        }

        private List<TestStep> parseSteps(Element stepsElement, String byteOrder) {
            if (stepsElement == null) {
                return Collections.emptyList();
            }

            return stepsElement.elements().stream()
                .map(element -> TestStep.parseTestStep(basePackage, (Element) element, byteOrder))
                .collect(Collectors.toList());
        }

        private String extractOptionalText(Element element, String name) {
            Element child = element.element(new QName(name));
            return child != null ? child.getTextTrim() : null;
        }
    }
}
