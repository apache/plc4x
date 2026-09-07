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
package org.apache.plc4x.java.utils.testutils.utils.xml.comparison;

/**
 * Represents a single difference found during XML comparison.
 */
public class XmlDifference {

    private final String path;
    private final DifferenceType type;
    private final String expectedValue;
    private final String actualValue;
    private final String description;

    public XmlDifference(String path, DifferenceType type,
                        String expectedValue, String actualValue, String description) {
        this.path = path;
        this.type = type;
        this.expectedValue = expectedValue;
        this.actualValue = actualValue;
        this.description = description;
    }

    public String getPath() {
        return path;
    }

    public DifferenceType getType() {
        return type;
    }

    public String getExpectedValue() {
        return expectedValue;
    }

    public String getActualValue() {
        return actualValue;
    }

    public String getDescription() {
        return description;
    }

    @Override
    public String toString() {
        return String.format("At %s: %s - Expected: '%s', Actual: '%s'",
            path, type.getDescription(), expectedValue, actualValue);
    }

    /**
     * Types of XML differences that can be detected.
     */
    public enum DifferenceType {
        ELEMENT_NAME("Element name mismatch"),
        TEXT_CONTENT("Text content mismatch"),
        ATTRIBUTE_COUNT("Attribute count mismatch"),
        ATTRIBUTE_VALUE("Attribute value mismatch"),
        CHILD_COUNT("Child element count mismatch"),
        MISSING_ELEMENT("Missing element"),
        EXTRA_ELEMENT("Extra element");

        private final String description;

        DifferenceType(String description) {
            this.description = description;
        }

        public String getDescription() {
            return description;
        }
    }
}
