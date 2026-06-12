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
import org.apache.plc4x.java.api.messages.*;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

/**
 * Handler for validating API responses.
 * Retrieves the pending response future from TestContext and validates the result.
 */
public class ApiResponseHandler {

    private static final Logger LOGGER = LoggerFactory.getLogger(ApiResponseHandler.class);
    private static final long DEFAULT_TIMEOUT_MS = 5000;

    private final Element responseXml;

    public ApiResponseHandler(Element responseXml) {
        this.responseXml = responseXml;
    }

    /**
     * Executes the API response validation.
     * Retrieves the pending future from context, waits for it to complete, and validates the response.
     *
     * @param context the test context containing the pending response future
     * @throws DriverTestsuiteException if validation fails or no pending response exists
     */
    public void executeApiResponse(TestContext context) {
        try {
            // Retrieve the pending response future
            CompletableFuture<?> future = context.getPendingResponse();
            if (future == null) {
                throw new DriverTestsuiteException("No pending response found in context");
            }

            // Wait for the response
            Object response = future.get(DEFAULT_TIMEOUT_MS, TimeUnit.MILLISECONDS);

            // Validate based on the response type
            switch (response) {
                case PlcReadResponse readResponse -> validateReadResponse(readResponse, context);
                case PlcWriteResponse writeResponse -> validateWriteResponse(writeResponse);
                case PlcBrowseResponse browseResponse -> validateBrowseResponse(browseResponse, context);
                case PlcSubscriptionResponse subscriptionResponse -> validateSubscriptionResponse(subscriptionResponse);
                default -> LOGGER.warn("Unknown response type: {}", response.getClass().getName());
            }

            LOGGER.debug("API response validated successfully");

        } catch (DriverTestsuiteException e) {
            throw e;
        } catch (Exception e) {
            throw new DriverTestsuiteException("Failed to validate API response", e);
        }
    }

    private void validateReadResponse(PlcReadResponse response, TestContext context) {
        // Get the expected values element from responseXml
        Element plcReadResponseElement = responseXml;
        if (!plcReadResponseElement.getName().equals("PlcReadResponse")) {
            throw new DriverTestsuiteException("Expecting PlcReadResponse element in responseXml");
        }

        Element valuesElement = plcReadResponseElement.element("values");
        if (valuesElement == null) {
            // If no values element and auto-migrate is enabled, we'll create one during migration
            if (!context.isAutoMigrate()) {
                throw new DriverTestsuiteException("No values element found in PlcReadResponse");
            }
        }

        // Generate the actual XML from the read response for comparison and potential migration
        Element actualReadResponseElement = serializeReadResponseToXml(response);

        // Compare only the <values> sections since the expected XML may contain a <request> element
        // that we don't serialize. The values section is what matters for validation.
        Element expectedValuesElement = plcReadResponseElement.element("values");
        Element actualValuesElement = actualReadResponseElement.element("values");
        if (expectedValuesElement == null) {
            expectedValuesElement = plcReadResponseElement;
        }
        if (actualValuesElement == null) {
            actualValuesElement = actualReadResponseElement;
        }

        // Use the shared XML comparison module
        XmlComparator comparator = new DomXmlComparator();
        XmlComparisonResult result = comparator.compare(expectedValuesElement, actualValuesElement);

        if (result.hasDifferences()) {
            // Check if auto-migrate is enabled
            if (context.isAutoMigrate()) {
                LOGGER.warn("Read response validation failed with {} differences. Auto-migrating...",
                    result.getDifferenceCount());
                XmlAutoMigrator.migrateDom(plcReadResponseElement, actualReadResponseElement,
                    context.getTestsuiteUri());
                LOGGER.info("Auto-migration completed successfully");

                // Log the response for visibility
                StringBuilder sb = new StringBuilder("API read response (migrated):\n");
                for (String tagName : response.getTagNames()) {
                    PlcResponseCode actualResponseCode = response.getResponseCode(tagName);
                    PlcValue actualValue = response.getPlcValue(tagName);
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
        for (String tagName : response.getTagNames()) {
            PlcResponseCode actualResponseCode = response.getResponseCode(tagName);
            PlcValue actualValue = response.getPlcValue(tagName);
            if (actualResponseCode != PlcResponseCode.OK) {
                sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append("\n");
            } else {
                sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append(": ").append(actualValue).append("\n");
            }
        }
        LOGGER.info(sb.toString());
    }

    /**
     * Serializes a PlcReadResponse to XML format matching the expected structure.
     */
    private Element serializeReadResponseToXml(PlcReadResponse response) {
        Element readResponseElement = DocumentHelper.createElement("PlcReadResponse");
        Element valuesElement = readResponseElement.addElement("values").addAttribute("isList", "true");

        for (String tagName : response.getTagNames()) {
            Element tagElement = valuesElement.addElement(sanitizeXmlElementName(tagName));
            Element responseItem = tagElement.addElement("PlcResponseItem");

            // Add response code with attributes matching the existing format
            PlcResponseCode responseCode = response.getResponseCode(tagName);
            Element codeElement = responseItem.addElement("code");
            codeElement.addAttribute("dataType", "string");
            codeElement.addAttribute("bitLength", "16");
            codeElement.addAttribute("encoding", "UTF-8");
            codeElement.setText(responseCode.name());

            // Add value if response is OK
            if (responseCode == PlcResponseCode.OK) {
                PlcValue plcValue = response.getPlcValue(tagName);
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

    private void validateWriteResponse(PlcWriteResponse response) {
        // Get the expected response codes element from responseXml
        Element plcWriteResponseElement = responseXml;
        if (!plcWriteResponseElement.getName().equals("PlcWriteResponse")) {
            throw new DriverTestsuiteException("Expecting PlcWriteResponse element in responseXml");
        }

        Element responseCodesElement = plcWriteResponseElement.element("responseCodes");
        if (responseCodesElement == null) {
            throw new DriverTestsuiteException("No responseCodes element found in PlcWriteResponse");
        }

        List<String> errors = new ArrayList<>();
        StringBuilder sb = new StringBuilder("API write response:\n");

        for (String tagName : response.getTagNames()) {
            PlcResponseCode actualResponseCode = response.getResponseCode(tagName);

            // Find the expected response code for this tag (sanitize name for XML lookup since
            // brackets in tag names like MAIN.g_arrBool[1] become underscores in XML elements)
            Element tagElement = responseCodesElement.element(sanitizeXmlElementName(tagName));
            if (tagElement == null) {
                errors.add(String.format("Tag '%s': No expected response code found in responseXml", tagName));
                continue;
            }

            Element responseCodeElement = tagElement.element("ResponseCode");
            if (responseCodeElement == null) {
                errors.add(String.format("Tag '%s': No ResponseCode element found", tagName));
                continue;
            }

            // Get expected code from stringRepresentation attribute or text content
            String expectedCode = responseCodeElement.attributeValue("stringRepresentation");
            if (expectedCode == null) {
                expectedCode = responseCodeElement.getTextTrim();
            }

            if (!actualResponseCode.name().equals(expectedCode)) {
                errors.add(String.format("Tag '%s': Response code mismatch - Expected: %s, Actual: %s",
                    tagName, expectedCode, actualResponseCode.name()));
            }

            sb.append("  ").append(tagName).append(": ").append(actualResponseCode).append("\n");
        }

        LOGGER.info(sb.toString());

        if (!errors.isEmpty()) {
            StringBuilder errorMsg = new StringBuilder("Write response validation failed:\n");
            for (String error : errors) {
                errorMsg.append("  - ").append(error).append("\n");
            }
            throw new AssertionError(errorMsg.toString());
        }
    }

    private void validateBrowseResponse(PlcBrowseResponse response, TestContext context) {
        // Get the expected browse response element from responseXml
        Element expectedBrowseResponseElement = responseXml;
        if (!expectedBrowseResponseElement.getName().equals("PlcBrowseResponse")) {
            throw new DriverTestsuiteException("Expecting PlcBrowseResponse element in responseXml");
        }

        // Generate the actual XML from the browse response
        Element actualBrowseResponseElement = serializeBrowseResponseToXml(response);

        // Use the shared XML comparison module
        XmlComparator comparator = new DomXmlComparator();
        XmlComparisonResult result = comparator.compare(expectedBrowseResponseElement, actualBrowseResponseElement);

        if (result.hasDifferences()) {
            // Check if auto-migrate is enabled
            if (context.isAutoMigrate()) {
                LOGGER.warn("Browse response validation failed with {} differences. Auto-migrating...",
                    result.getDifferenceCount());
                XmlAutoMigrator.migrateDom(expectedBrowseResponseElement, actualBrowseResponseElement,
                    context.getTestsuiteUri());
                LOGGER.info("Auto-migration completed successfully");
                return;
            }

            // Log differences using the shared reporter
            XmlDifferenceReporter.logDiffReport(result, LOGGER);

            // Build error message for exception
            StringBuilder errorMsg = new StringBuilder("Browse response validation failed:\n");
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
        int totalItems = response.getQueryNames().stream()
            .mapToInt(q -> response.getValues(q).size())
            .sum();

        // Output a simplified version to the console.
        if (LOGGER.isInfoEnabled()) {
            StringBuilder sb = new StringBuilder("API read response:\n");
            for (String queryName : response.getQueryNames()) {
                PlcResponseCode responseCode = response.getResponseCode(queryName);
                if (responseCode != PlcResponseCode.OK) {
                    sb.append("  ").append(queryName).append(": ").append(responseCode).append("\n");
                } else {
                    sb.append("  ").append(queryName).append(": ").append(responseCode).append(":\n");
                    List<PlcBrowseItem> values = response.getValues(queryName);
                    for (PlcBrowseItem browseItem : values) {
                        outputBrowseItem(browseItem, sb, 0);
                    }
                }
            }
            LOGGER.info(sb.toString());
        }
    }

    private void outputBrowseItem(PlcBrowseItem browseItem, StringBuilder sb, int level) {
        sb.append("  ").append("  ".repeat(level)).append(browseItem.getName()).append(": ").append(browseItem.getTag().getPlcValueType().name()).append("\n");
        if(browseItem.getChildren() != null) {
            for(Map.Entry<String, PlcBrowseItem> child : browseItem.getChildren().entrySet()) {
                outputBrowseItem(child.getValue(), sb, level + 1);
            }
        }
    }

    /**
     * Serializes a PlcBrowseResponse to XML format matching the expected structure.
     */
    private Element serializeBrowseResponseToXml(PlcBrowseResponse response) {
        Element browseResponseElement = DocumentHelper.createElement("PlcBrowseResponse");
        Element valuesElement = browseResponseElement.addElement("values");

        // Track visited items to prevent infinite recursion (using identity-based comparison)
        Set<Object> visited = Collections.newSetFromMap(new IdentityHashMap<>());

        for (String queryName : response.getQueryNames()) {
            Element queryElement = valuesElement.addElement(queryName);
            Element responseItem = queryElement.addElement("PlcResponseItem");

            // Add response code
            PlcResponseCode responseCode = response.getResponseCode(queryName);
            responseItem.addElement("code").setText(responseCode.name());

            // Add values if response is OK
            if (responseCode == PlcResponseCode.OK) {
                List<?> browseResults = response.getValues(queryName);
                if (!browseResults.isEmpty()) {
                    Element valuesListElement = responseItem.addElement("values");
                    for (Object browseResult : browseResults) {
                        serializeBrowseItem(valuesListElement, browseResult, visited);
                    }
                }
            }
        }

        return browseResponseElement;
    }

    /**
     * Serializes a browse result item to XML.
     * Handles PlcBrowseItem with all its properties.
     * Uses a visited set to prevent infinite recursion from circular references.
     */
    private void serializeBrowseItem(Element parentElement, Object browseResult, Set<Object> visited) {
        if (!(browseResult instanceof PlcBrowseItem browseItem)) {
            // Fallback for non-PlcBrowseItem objects
            Element itemElement = parentElement.addElement("item");
            itemElement.setText(browseResult.toString());
            return;
        }

        // Check for circular reference
        if (visited.contains(browseResult)) {
            Element itemElement = parentElement.addElement("PlcBrowseItem");
            itemElement.addElement("circularReference").setText("true");
            if (browseItem.getName() != null) {
                itemElement.addElement("name").setText(browseItem.getName());
            }
            return;
        }

        // Mark as visited
        visited.add(browseResult);

        try {
            Element itemElement = parentElement.addElement("PlcBrowseItem");

            // Serialize tag address
            if (browseItem.getTag() != null) {
                itemElement.addElement("address").setText(browseItem.getTag().getAddressString());
            }

            // Serialize name
            if (browseItem.getName() != null) {
                itemElement.addElement("name").setText(browseItem.getName());
            }

            // Serialize access flags
            itemElement.addElement("readable").setText(String.valueOf(browseItem.isReadable()));
            itemElement.addElement("writable").setText(String.valueOf(browseItem.isWritable()));
            itemElement.addElement("subscribable").setText(String.valueOf(browseItem.isSubscribable()));
            itemElement.addElement("publishable").setText(String.valueOf(browseItem.isPublishable()));

            // Serialize array information
            List<ArrayInfo> arrayInfo = browseItem.getArrayInformation();
            if (arrayInfo != null && !arrayInfo.isEmpty()) {
                Element arrayInfoElement = itemElement.addElement("arrayInformation");
                for (ArrayInfo info : arrayInfo) {
                    Element infoElement = arrayInfoElement.addElement("ArrayInfo");
                    infoElement.addElement("lowerBound").setText(String.valueOf(info.getLowerBound()));
                    infoElement.addElement("upperBound").setText(String.valueOf(info.getUpperBound()));
                }
            }

            // Serialize children
            Map<String, PlcBrowseItem> children = browseItem.getChildren();
            if (children != null && !children.isEmpty()) {
                Element childrenElement = itemElement.addElement("children");
                for (Map.Entry<String, PlcBrowseItem> entry : children.entrySet()) {
                    Element childWrapper = childrenElement.addElement("child");
                    childWrapper.addElement("key").setText(entry.getKey());
                    serializeBrowseItem(childWrapper, entry.getValue(), visited);
                }
            }

            // Serialize options
            Map<String, PlcValue> options = browseItem.getOptions();
            if (options != null && !options.isEmpty()) {
                Element optionsElement = itemElement.addElement("options");
                for (Map.Entry<String, PlcValue> entry : options.entrySet()) {
                    Element optionElement = optionsElement.addElement("option");
                    optionElement.addElement("name").setText(entry.getKey());
                    optionElement.addElement("value").setText(formatPlcValue(entry.getValue()));
                }
            }
        } finally {
            // Remove from visited after processing to allow the same item in different branches
            visited.remove(browseResult);
        }
    }

    private void validateSubscriptionResponse(PlcSubscriptionResponse response) {
        // TODO: Instead of hard-coding the expected response codes, we should use the responseXml to extract them.
        for (String tagName : response.getTagNames()) {
            PlcResponseCode responseCode = response.getResponseCode(tagName);
            if (responseCode != PlcResponseCode.OK) {
                LOGGER.warn("Subscription for tag {} returned response code: {}", tagName, responseCode);
            } else {
                LOGGER.debug("Subscription for tag {} created successfully", tagName);
            }
        }
    }

    /**
     * Sanitizes a tag name for use as an XML element name.
     * Replaces characters that are invalid in XML element names with underscores.
     * Examples:
     * - "MAIN.g_arrBool[1]" → "MAIN.g_arrBool_1_"
     * - "Program:TestProgram.g_b1" → "Program_TestProgram.g_b1"
     * - "%DB4:0.0:BOOL" → "_DB4_0.0_BOOL"
     * - "%DB4:140:STRING(200)" → "_DB4_140_STRING_200_"
     */
    private String sanitizeXmlElementName(String tagName) {
        return tagName
            .replace('%', '_')
            .replace('[', '_').replace(']', '_')
            .replace('(', '_').replace(')', '_')
            .replace(':', '_');
    }

    /**
     * Formats a PlcValue to a string for comparison.
     */
    private String formatPlcValue(PlcValue value) {
        if (value == null) {
            return "null";
        }

        if (value.isList()) {
            List<String> values = new ArrayList<>();
            for (int i = 0; i < value.getLength(); i++) {
                values.add(value.getIndex(i).toString());
            }
            return values.toString();
        } else {
            return value.toString();
        }
    }

}
