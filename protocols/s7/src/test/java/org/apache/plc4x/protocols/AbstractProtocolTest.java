/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */

package org.apache.plc4x.protocols;

import org.apache.daffodil.tdml.DFDLTestSuite;
import org.apache.daffodil.tdml.Runner;
import org.apache.daffodil.util.Misc;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;
import scala.collection.Iterator;

import java.util.LinkedList;
import java.util.List;

public abstract class AbstractProtocolTest {

    private final String testsuiteName;

    public AbstractProtocolTest(String testsuiteName) {
        this.testsuiteName = testsuiteName;
    }

    @TestFactory
    public List<DynamicTest> getTestsuiteTests() {
        DFDLTestSuite testSuite = new DFDLTestSuite(Misc.getRequiredResource(testsuiteName), true, true, false,
            Runner.defaultRoundTripDefaultDefault(),
            Runner.defaultValidationDefaultDefault(),
            Runner.defaultImplementationsDefaultDefault(),
            Runner.defaultShouldDoErrorComparisonOnCrossTests(),
            Runner.defaultShouldDoWarningComparisonOnCrossTests());
        List<DynamicTest> dynamicTests = new LinkedList<>();
        Iterator<String> iterator = testSuite.testCaseMap().keySet().iterator();
        while (iterator.hasNext()) {
            String testcaseName = iterator.next();
            String testcaseLabel = testSuite.suiteName() + ": " + testcaseName;
            DynamicTest test = DynamicTest.dynamicTest(testcaseLabel, () ->
                testSuite.runOneTest(testcaseName, scala.Option.apply(null), false)
            );
            dynamicTests.add(test);
        }
        return dynamicTests;
    }

}
