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

package org.apache.plc4x.test.driver.xmlunit;

import org.w3c.dom.NamedNodeMap;
import org.w3c.dom.Node;
import org.xmlunit.diff.Comparison;
import org.xmlunit.diff.ComparisonResult;
import org.xmlunit.diff.DifferenceEvaluator;

/**
 * Evaluator of differences which allows to ignore differences for elements annotated with 'plc4x-skip-comparison' attribute.
 */
public class SkipDifferenceEvaluator implements DifferenceEvaluator {

    @Override
    public ComparisonResult evaluate(Comparison comparison, ComparisonResult comparisonResult) {
        if (comparisonResult != ComparisonResult.EQUAL) {
            Node target = comparison.getControlDetails().getTarget();

            // root element
            if (target == null || target.getParentNode() == null) {
                return comparisonResult;
            }

            // verify parent element - help with text nodes
            NamedNodeMap attributes = target.getParentNode().getAttributes();
            if (attributes != null) {
                Node attribute = attributes.getNamedItem(SkipAttributeFilter.IGNORE_ATTRIBUTE_NAME);
                if (attribute != null) {
                    String content = attribute.getTextContent();
                    return Boolean.parseBoolean(content.trim()) ? ComparisonResult.EQUAL : comparisonResult;
                }
            }
        }

        return comparisonResult;
    }
}
