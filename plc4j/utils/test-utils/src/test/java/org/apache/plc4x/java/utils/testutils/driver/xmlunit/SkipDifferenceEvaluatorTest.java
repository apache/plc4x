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
package org.apache.plc4x.java.utils.testutils.driver.xmlunit;

import org.junit.jupiter.api.Test;
import org.w3c.dom.NamedNodeMap;
import org.w3c.dom.Node;
import org.xmlunit.diff.Comparison;
import org.xmlunit.diff.Comparison.Detail;
import org.xmlunit.diff.ComparisonResult;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class SkipDifferenceEvaluatorTest {

    @Test
    void testEvaluateEqualComparison() {
        SkipDifferenceEvaluator evaluator = new SkipDifferenceEvaluator();
        Comparison comparison = mock(Comparison.class);

        ComparisonResult result = evaluator.evaluate(comparison, ComparisonResult.EQUAL);
        assertEquals(ComparisonResult.EQUAL, result);
    }

    @Test
    void testEvaluateWithNullTarget() {
        SkipDifferenceEvaluator evaluator = new SkipDifferenceEvaluator();
        Comparison comparison = mock(Comparison.class);
        Detail detail = mock(Detail.class);

        when(comparison.getControlDetails()).thenReturn(detail);
        when(detail.getTarget()).thenReturn(null);

        ComparisonResult result = evaluator.evaluate(comparison, ComparisonResult.DIFFERENT);
        assertEquals(ComparisonResult.DIFFERENT, result);
    }

    @Test
    void testEvaluateWithSkipAttribute() {
        SkipDifferenceEvaluator evaluator = new SkipDifferenceEvaluator();
        Comparison comparison = mock(Comparison.class);
        Detail detail = mock(Detail.class);
        Node target = mock(Node.class);
        Node parent = mock(Node.class);
        NamedNodeMap attributes = mock(NamedNodeMap.class);
        Node skipAttr = mock(Node.class);

        when(comparison.getControlDetails()).thenReturn(detail);
        when(detail.getTarget()).thenReturn(target);
        when(target.getParentNode()).thenReturn(parent);
        when(parent.getAttributes()).thenReturn(attributes);
        when(attributes.getNamedItem(SkipAttributeFilter.IGNORE_ATTRIBUTE_NAME)).thenReturn(skipAttr);
        when(skipAttr.getTextContent()).thenReturn("true");

        ComparisonResult result = evaluator.evaluate(comparison, ComparisonResult.DIFFERENT);
        assertEquals(ComparisonResult.EQUAL, result);
    }
}
