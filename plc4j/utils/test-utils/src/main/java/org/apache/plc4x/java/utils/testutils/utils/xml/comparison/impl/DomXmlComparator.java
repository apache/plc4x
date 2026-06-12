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
package org.apache.plc4x.java.utils.testutils.utils.xml.comparison.impl;

import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparisonResult;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlDifference;
import org.dom4j.DocumentException;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;

import java.util.ArrayList;
import java.util.List;

/**
 * DOM-based XML comparator that traverses the tree directly.
 * Memory-efficient for very large XML documents as it doesn't build huge strings.
 * Recommended for browse responses and other large structures.
 */
public class DomXmlComparator implements XmlComparator {

    @Override
    public XmlComparisonResult compare(Element expected, Element actual) {
        List<XmlDifference> differences = new ArrayList<>();
        compareElements(expected, actual, "", differences);

        if (differences.isEmpty()) {
            return XmlComparisonResult.identical(expected.asXML());
        }

        return XmlComparisonResult.different(differences, expected.asXML(), actual.asXML());
    }

    @Override
    public XmlComparisonResult compare(String expected, String actual) {
        try {
            Element expectedElement = DocumentHelper.parseText(expected).getRootElement();
            Element actualElement = DocumentHelper.parseText(actual).getRootElement();
            return compare(expectedElement, actualElement);
        } catch (DocumentException e) {
            throw new RuntimeException("Failed to parse XML strings", e);
        }
    }

    /**
     * Recursively compares two DOM4J elements and collects differences.
     *
     * @param expected the expected element
     * @param actual the actual element
     * @param path the current path in the tree (for error reporting)
     * @param differences list to collect found differences
     */
    private void compareElements(Element expected, Element actual, String path, List<XmlDifference> differences) {
        String currentPath = path.isEmpty() ? expected.getName() : path + "/" + expected.getName();

        // Check if actual element exists
        if (actual == null) {
            differences.add(new XmlDifference(
                currentPath,
                XmlDifference.DifferenceType.MISSING_ELEMENT,
                expected.getName(),
                "null",
                "Expected element but got null"
            ));
            return;
        }

        // Check element names match
        if (!expected.getName().equals(actual.getName())) {
            differences.add(new XmlDifference(
                currentPath,
                XmlDifference.DifferenceType.ELEMENT_NAME,
                expected.getName(),
                actual.getName(),
                "Element name mismatch"
            ));
            return;
        }

        // Check text content
        String expectedText = expected.getTextTrim();
        String actualText = actual.getTextTrim();
        if (!expectedText.isEmpty() || !actualText.isEmpty()) {
            if (!expectedText.equals(actualText)) {
                differences.add(new XmlDifference(
                    currentPath,
                    XmlDifference.DifferenceType.TEXT_CONTENT,
                    expectedText,
                    actualText,
                    "Text content mismatch"
                ));
            }
        }

        // Check attributes
        List<?> expectedAttributes = expected.attributes();
        List<?> actualAttributes = actual.attributes();
        if (expectedAttributes.size() != actualAttributes.size()) {
            differences.add(new XmlDifference(
                currentPath,
                XmlDifference.DifferenceType.ATTRIBUTE_COUNT,
                String.valueOf(expectedAttributes.size()),
                String.valueOf(actualAttributes.size()),
                "Attribute count mismatch"
            ));
        }

        // Compare child elements
        List<Element> expectedChildren = expected.elements();
        List<Element> actualChildren = actual.elements();

        if (expectedChildren.size() != actualChildren.size()) {
            differences.add(new XmlDifference(
                currentPath,
                XmlDifference.DifferenceType.CHILD_COUNT,
                String.valueOf(expectedChildren.size()),
                String.valueOf(actualChildren.size()),
                "Child element count mismatch"
            ));
            // Continue to compare what we can
        }

        // Compare children one by one
        int minChildren = Math.min(expectedChildren.size(), actualChildren.size());
        for (int i = 0; i < minChildren; i++) {
            compareElements(expectedChildren.get(i), actualChildren.get(i), currentPath, differences);
        }
    }
}
