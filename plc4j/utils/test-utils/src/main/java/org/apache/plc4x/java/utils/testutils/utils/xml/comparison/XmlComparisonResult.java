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

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * Result of an XML comparison operation.
 * Contains information about differences found and provides methods
 * for reporting and auto-migration.
 */
public class XmlComparisonResult {

    private final boolean identical;
    private final List<XmlDifference> differences;
    private final String expectedXml;
    private final String actualXml;

    private XmlComparisonResult(boolean identical, List<XmlDifference> differences,
                               String expectedXml, String actualXml) {
        this.identical = identical;
        this.differences = differences != null ? new ArrayList<>(differences) : Collections.emptyList();
        this.expectedXml = expectedXml;
        this.actualXml = actualXml;
    }

    /**
     * Creates a result indicating the XMLs are identical.
     */
    public static XmlComparisonResult identical(String xml) {
        return new XmlComparisonResult(true, Collections.emptyList(), xml, xml);
    }

    /**
     * Creates a result indicating differences were found.
     */
    public static XmlComparisonResult different(List<XmlDifference> differences,
                                               String expectedXml, String actualXml) {
        return new XmlComparisonResult(false, differences, expectedXml, actualXml);
    }

    /**
     * Returns true if the XMLs are identical.
     */
    public boolean isIdentical() {
        return identical;
    }

    /**
     * Returns true if differences were found.
     */
    public boolean hasDifferences() {
        return !identical;
    }

    /**
     * Returns the list of differences found.
     */
    public List<XmlDifference> getDifferences() {
        return Collections.unmodifiableList(differences);
    }

    /**
     * Returns the number of differences found.
     */
    public int getDifferenceCount() {
        return differences.size();
    }

    /**
     * Returns the expected XML.
     */
    public String getExpectedXml() {
        return expectedXml;
    }

    /**
     * Returns the actual XML.
     */
    public String getActualXml() {
        return actualXml;
    }

    @Override
    public String toString() {
        if (identical) {
            return "XMLs are identical";
        }
        return String.format("Found %d differences between XMLs", differences.size());
    }
}
