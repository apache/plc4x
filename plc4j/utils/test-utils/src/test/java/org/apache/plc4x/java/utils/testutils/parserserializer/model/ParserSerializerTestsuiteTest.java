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
package org.apache.plc4x.java.utils.testutils.parserserializer.model;

import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

class ParserSerializerTestsuiteTest {

    @Test
    void testParserSerializerTestsuiteCreation() {
        List<Testcase> testcases = new ArrayList<>();
        Map<String, String> options = new HashMap<>();

        ParserSerializerTestsuite suite = new ParserSerializerTestsuite(
            "TestSuite",
            testcases,
            "BIG_ENDIAN",
            options
        );

        assertNotNull(suite);
        assertEquals("TestSuite", suite.name());
        assertEquals(testcases, suite.testcases());
        assertEquals("BIG_ENDIAN", suite.byteOrderName());
        assertEquals(options, suite.options());
    }

    @Test
    void testParserSerializerTestsuiteEquality() {
        List<Testcase> testcases = new ArrayList<>();
        Map<String, String> options = new HashMap<>();

        ParserSerializerTestsuite suite1 = new ParserSerializerTestsuite("Test", testcases, "BIG_ENDIAN", options);
        ParserSerializerTestsuite suite2 = new ParserSerializerTestsuite("Test", testcases, "BIG_ENDIAN", options);

        assertEquals(suite1, suite2);
        assertEquals(suite1.hashCode(), suite2.hashCode());
    }
}
