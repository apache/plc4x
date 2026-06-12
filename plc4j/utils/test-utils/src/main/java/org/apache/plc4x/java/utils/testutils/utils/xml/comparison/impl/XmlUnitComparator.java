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

import org.apache.plc4x.java.utils.testutils.driver.xmlunit.SkipAttributeFilter;
import org.apache.plc4x.java.utils.testutils.driver.xmlunit.SkipDifferenceEvaluator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparator;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlComparisonResult;
import org.apache.plc4x.java.utils.testutils.utils.xml.comparison.XmlDifference;
import org.dom4j.Element;
import org.xmlunit.builder.DiffBuilder;
import org.xmlunit.diff.Comparison;
import org.xmlunit.diff.Diff;
import org.xmlunit.diff.Difference;

import java.util.ArrayList;
import java.util.List;

/**
 * XMLUnit-based comparator that provides detailed diff information.
 * Better for smaller documents where detailed comparison is needed.
 * Uses XMLUnit for thorough analysis but may have memory issues with huge documents.
 */
public class XmlUnitComparator implements XmlComparator {

    @Override
    public XmlComparisonResult compare(Element expected, Element actual) {
        return compare(expected.asXML(), actual.asXML());
    }

    @Override
    public XmlComparisonResult compare(String expected, String actual) {
        final Diff diff = DiffBuilder.compare(expected)
            .withAttributeFilter(new SkipAttributeFilter())
            .withDifferenceEvaluator(new SkipDifferenceEvaluator())
            .withTest(actual)
            .checkForSimilar()
            .ignoreComments()
            .ignoreWhitespace()
            .build();

        if (!diff.hasDifferences()) {
            return XmlComparisonResult.identical(expected);
        }

        // Convert XMLUnit differences to our format
        List<XmlDifference> differences = new ArrayList<>();
        for (Difference difference : diff.getDifferences()) {
            differences.add(convertDifference(difference));
        }

        return XmlComparisonResult.different(differences, expected, actual);
    }

    /**
     * Converts an XMLUnit Difference to our XmlDifference format.
     */
    private XmlDifference convertDifference(Difference difference) {
        Comparison comparison = difference.getComparison();

        String expectedValue = comparison.getControlDetails().getValue() != null
            ? comparison.getControlDetails().getValue().toString() : "null";
        String actualValue = comparison.getTestDetails().getValue() != null
            ? comparison.getTestDetails().getValue().toString() : "null";

        String path = comparison.getControlDetails().getXPath();

        // Try to map XMLUnit comparison type to our DifferenceType
        XmlDifference.DifferenceType type = mapComparisonType(comparison.getType().name());

        String description = comparison.getType().getDescription();

        return new XmlDifference(path, type, expectedValue, actualValue, description);
    }

    /**
     * Maps XMLUnit comparison types to our DifferenceType enum.
     */
    private XmlDifference.DifferenceType mapComparisonType(String xmlUnitType) {
        return switch (xmlUnitType) {
            case "ELEMENT_TAG_NAME" -> XmlDifference.DifferenceType.ELEMENT_NAME;
            case "TEXT_VALUE" -> XmlDifference.DifferenceType.TEXT_CONTENT;
            case "ATTR_VALUE" -> XmlDifference.DifferenceType.ATTRIBUTE_VALUE;
            case "CHILD_NODELIST_LENGTH" -> XmlDifference.DifferenceType.CHILD_COUNT;
            default -> XmlDifference.DifferenceType.TEXT_CONTENT; // fallback
        };
    }
}
