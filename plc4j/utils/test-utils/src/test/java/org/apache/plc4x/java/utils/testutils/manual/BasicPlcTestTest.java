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
package org.apache.plc4x.java.utils.testutils.manual;

import org.apache.plc4x.java.api.authentication.PlcNullAuthentication;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.api.types.PlcResponseCode;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertSame;

class BasicPlcTestTest {

    @Test
    void testBasicPlcTestCreation() {
        BasicPlcTest test = new BasicPlcTest("test://localhost");
        assertNotNull(test);
    }

    @Test
    void testBasicPlcTestWithAllParameters() {
        BasicPlcTest test = new BasicPlcTest("test://localhost", new PlcNullAuthentication(), true, true, true, true, 100);
        assertNotNull(test);
    }

    @Test
    void testAddTestCaseWithExpectedValue() {
        BasicPlcTest test = new BasicPlcTest("test://localhost");
        PlcSTRING value = new PlcSTRING("test");
        BasicPlcTest result = test.addTestCase("tag1", value);
        assertSame(test, result); // Should return this for chaining
    }

    @Test
    void testAddTestCaseWithResponseCode() {
        BasicPlcTest test = new BasicPlcTest("test://localhost");
        BasicPlcTest result = test.addTestCase("tag1", PlcResponseCode.NOT_FOUND);
        assertSame(test, result);
    }

    @Test
    void testAddTestCaseWithResponseCodeAndValue() {
        BasicPlcTest test = new BasicPlcTest("test://localhost");
        PlcSTRING value = new PlcSTRING("test");
        BasicPlcTest result = test.addTestCase("tag1", PlcResponseCode.OK, value);
        assertSame(test, result);
    }

    @Test
    void testChaining() {
        BasicPlcTest test = new BasicPlcTest("test://localhost");
        PlcSTRING value1 = new PlcSTRING("test1");
        PlcSTRING value2 = new PlcSTRING("test2");

        BasicPlcTest result = test
            .addTestCase("tag1", value1)
            .addTestCase("tag2", PlcResponseCode.NOT_FOUND)
            .addTestCase("tag3", PlcResponseCode.OK, value2);

        assertSame(test, result);
    }
}
