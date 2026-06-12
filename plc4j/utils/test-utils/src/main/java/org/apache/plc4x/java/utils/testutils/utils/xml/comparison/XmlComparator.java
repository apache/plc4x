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

import org.dom4j.Element;

/**
 * Interface for XML comparison implementations.
 * Provides different strategies for comparing XML documents.
 */
public interface XmlComparator {

    /**
     * Compares two XML elements and returns the result.
     *
     * @param expected the expected XML element
     * @param actual the actual XML element
     * @return the comparison result
     */
    XmlComparisonResult compare(Element expected, Element actual);

    /**
     * Compares two XML strings and returns the result.
     *
     * @param expected the expected XML string
     * @param actual the actual XML string
     * @return the comparison result
     */
    XmlComparisonResult compare(String expected, String actual);
}
