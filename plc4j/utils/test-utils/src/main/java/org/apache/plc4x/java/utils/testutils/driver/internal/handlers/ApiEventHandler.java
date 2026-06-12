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
package org.apache.plc4x.java.utils.testutils.driver.internal.handlers;

import org.apache.plc4x.java.utils.testutils.driver.exceptions.DriverTestsuiteException;
import org.apache.plc4x.java.utils.testutils.driver.internal.TestContext;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlAutoMigrator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparisonResult;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlDifferenceReporter;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.impl.DomXmlComparator;
import org.apache.plc4x.java.api.messages.PlcSubscriptionEvent;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.TimeUnit;

/**
 * Handler for validating incoming API events.
 */
public class ApiEventHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(ApiEventHandler.class);

    private final Element eventXml;

    public ApiEventHandler(Element eventXml) {
        this.eventXml = eventXml;
    }

    /**
     * Executes the API event validation.
     *
     * @param context the test context containing the pending response future
     * @throws DriverTestsuiteException if validation fails or no pending response exists
     */
    public void executeApiEvent(TestContext context) {
        try {
            // Wait for the response
            PlcSubscriptionEvent subscriptionEvent = context.getSubscriptionEvent().get(1000, TimeUnit.MILLISECONDS);

            validateSubscriptionEvert(subscriptionEvent, context);

            LOGGER.debug("API event validated successfully");

        } catch (DriverTestsuiteException e) {
            throw e;
        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to validate API event", e);
        }
    }

    private void validateSubscriptionEvert(PlcSubscriptionEvent subscriptionEvent, TestContext context) {
        // Get the expected values element from eventXml
        Element plcEventElement = eventXml;
        if (!plcEventElement.getName().equals("PlcSubscriptionEvent")) {
            throw new DriverTestsuiteException("Expecting PlcSubscriptionEvent element in eventXml");
        }

        Element valuesElement = plcEventElement.element("values");
        if (valuesElement == null) {
            // If no values element and auto-migrate is enabled, we'll create one during migration
            if (!context.isAutoMigrate()) {
                throw new DriverTestsuiteException("No values element found in PlcReadResponse");
            }
        }

        // Generate the actual XML from the event for comparison and potential migration
        Element actualEventElement = serializeEventToXml(subscriptionEvent);

        // Use the shared XML comparison module
        XmlComparator comparator = new DomXmlComparator();
        XmlComparisonResult result = comparator.compare(plcEventElement, actualEventElement);

        if (result.hasDifferences()) {
            // Check if auto-migrate is enabled
            if (context.isAutoMigrate()) {
                LOGGER.warn("Read response validation failed with {} differences. Auto-migrating...",
                    result.getDifferenceCount());
                XmlAutoMigrator.migrateDom(plcEventElement, actualEventElement,
                    context.getTestsuiteUri());
                LOGGER.info("Auto-migration completed successfully");

                // Log the response for visibility
                StringBuilder sb = new StringBuilder("API event (migrated):\n");
                for (String tagName : subscriptionEvent.getTagNames()) {
                    PlcResponseCode actualResponseCode = subscriptionEvent.getResponseCode(tagName);
                    PlcValue actualValue = subscriptionEvent.getPlcValue(tagName);
                    if (actualResponseCode != PlcResponseCode.OK) {
                        sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append("\n");
                    } else {
                        sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append(": ").append(actualValue).append("\n");
                    }
                }
                LOGGER.info(sb.toString());
                return;
            }

            // Log differences using the shared reporter
            XmlDifferenceReporter.logDiffReport(result, LOGGER);

            // Build error message for exception
            StringBuilder errorMsg = new StringBuilder("Read response validation failed:\n");
            int maxDifferences = Math.min(result.getDifferenceCount(), 50);
            for (int i = 0; i < maxDifferences; i++) {
                errorMsg.append("  - ").append(result.getDifferences().get(i)).append("\n");
            }
            if (result.getDifferenceCount() > maxDifferences) {
                errorMsg.append(String.format("  ... and %d more differences\n",
                    result.getDifferenceCount() - maxDifferences));
            }

            throw new AssertionError(errorMsg.toString());
        }

        // Log success
        StringBuilder sb = new StringBuilder("API read response:\n");
        for (String tagName : subscriptionEvent.getTagNames()) {
            PlcResponseCode actualResponseCode = subscriptionEvent.getResponseCode(tagName);
            PlcValue actualValue = subscriptionEvent.getPlcValue(tagName);
            if (actualResponseCode != PlcResponseCode.OK) {
                sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append("\n");
            } else {
                sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append(": ").append(actualValue).append("\n");
            }
        }
        LOGGER.info(sb.toString());
    }

    /**
     * Serializes a PlcSubscriptionEvent to XML format matching the expected structure.
     */
    private Element serializeEventToXml(PlcSubscriptionEvent subscriptionEvent) {
        Element readResponseElement = DocumentHelper.createElement("PlcSubscriptionEvent");
        Element valuesElement = readResponseElement.addElement("values");

        for (String tagName : subscriptionEvent.getTagNames()) {
            Element tagElement = valuesElement.addElement(tagName);
            Element responseItem = tagElement.addElement("PlcResponseItem");

            // Add response code with attributes matching the existing format
            PlcResponseCode responseCode = subscriptionEvent.getResponseCode(tagName);
            Element codeElement = responseItem.addElement("code");
            codeElement.addAttribute("dataType", "string");
            codeElement.addAttribute("bitLength", "16");
            codeElement.addAttribute("encoding", "UTF-8");
            codeElement.setText(responseCode.name());

            // Add value if event is OK
            if (responseCode == PlcResponseCode.OK) {
                PlcValue plcValue = subscriptionEvent.getPlcValue(tagName);
                if (plcValue != null) {
                    Element valueElement = responseItem.addElement("value");
                    serializePlcValue(valueElement, plcValue);
                }
            }
        }

        return readResponseElement;
    }

    /**
     * Serializes a PlcValue to XML.
     * Handles simple values, lists, and structs recursively.
     */
    private void serializePlcValue(Element parentElement, PlcValue value) {
        serializePlcValue(parentElement, value, 0);
    }

    private static final int MAX_SERIALIZATION_DEPTH = 50;

    private void serializePlcValue(Element parentElement, PlcValue value, int depth) {
        if (value == null) {
            parentElement.addElement("PlcNull");
            return;
        }

        // Prevent stack overflow for deeply nested structures
        if (depth > MAX_SERIALIZATION_DEPTH) {
            parentElement.addElement("PlcTruncated").setText("[nested too deep]");
            return;
        }

        PlcValueType plcValueType = value.getPlcValueType();

        // Handle RAW_BYTE_ARRAY first since it reports isList()=true but should be serialized as a hex string
        if (plcValueType == PlcValueType.RAW_BYTE_ARRAY) {
            Element valueElement = parentElement.addElement("PlcRAW_BYTE_ARRAY");
            addTypeAttributes(valueElement, plcValueType);
            valueElement.setText(value.toString());
        } else if (value.isStruct()) {
            Element structElement = parentElement.addElement("PlcStruct");
            for (String key : value.getKeys()) {
                Element fieldElement = structElement.addElement(key);
                serializePlcValue(fieldElement, value.getValue(key), depth + 1);
            }
        } else if (value.isList()) {
            Element listElement = parentElement.addElement("PlcList");
            for (int i = 0; i < value.getLength(); i++) {
                serializePlcValue(listElement, value.getIndex(i), depth + 1);
            }
        } else {
            // Simple value - use the PlcValue type name as element name with attributes
            String typeName = (plcValueType != null) ? "Plc" + plcValueType.name() : "PlcValue";
            Element valueElement = parentElement.addElement(typeName);

            // Add dataType and bitLength attributes based on the PlcValueType
            if (plcValueType != null) {
                addTypeAttributes(valueElement, plcValueType);
            }

            valueElement.setText(value.toString());
        }
    }

    /**
     * Adds dataType and bitLength attributes to an element based on the PlcValueType.
     */
    private void addTypeAttributes(Element element, PlcValueType plcValueType) {
        switch (plcValueType) {
            case BOOL -> {
                element.addAttribute("dataType", "bit");
                element.addAttribute("bitLength", "1");
            }
            case BYTE, USINT -> {
                element.addAttribute("dataType", "uint");
                element.addAttribute("bitLength", "8");
            }
            case SINT -> {
                element.addAttribute("dataType", "int");
                element.addAttribute("bitLength", "8");
            }
            case UINT, WORD -> {
                element.addAttribute("dataType", "uint");
                element.addAttribute("bitLength", "16");
            }
            case INT -> {
                element.addAttribute("dataType", "int");
                element.addAttribute("bitLength", "16");
            }
            case UDINT, DWORD -> {
                element.addAttribute("dataType", "uint");
                element.addAttribute("bitLength", "32");
            }
            case DINT -> {
                element.addAttribute("dataType", "int");
                element.addAttribute("bitLength", "32");
            }
            case ULINT, LWORD -> {
                element.addAttribute("dataType", "uint");
                element.addAttribute("bitLength", "64");
            }
            case LINT -> {
                element.addAttribute("dataType", "int");
                element.addAttribute("bitLength", "64");
            }
            case REAL -> {
                element.addAttribute("dataType", "float");
                element.addAttribute("bitLength", "32");
            }
            case LREAL -> {
                element.addAttribute("dataType", "float");
                element.addAttribute("bitLength", "64");
            }
            case STRING, WSTRING, CHAR, WCHAR -> {
                element.addAttribute("dataType", "string");
                element.addAttribute("encoding", "UTF-8");
            }
            case TIME, LTIME, DATE, LDATE, TIME_OF_DAY, LTIME_OF_DAY, DATE_AND_TIME, LDATE_AND_TIME, DATE_AND_LTIME -> {
                element.addAttribute("dataType", "time");
            }
            case RAW_BYTE_ARRAY -> {
                element.addAttribute("dataType", "byte");
            }
            default -> {
                // For unknown types, just use the type name
                element.addAttribute("dataType", plcValueType.name().toLowerCase());
            }
        }
    }

}
