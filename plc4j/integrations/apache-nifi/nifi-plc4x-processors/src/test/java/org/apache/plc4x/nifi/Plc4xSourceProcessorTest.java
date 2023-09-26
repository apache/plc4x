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
package org.apache.plc4x.nifi;

import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.apache.plc4x.nifi.address.AddressesAccessUtils;
import org.apache.plc4x.nifi.util.Plc4xCommonTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class Plc4xSourceProcessorTest {

    private TestRunner testRunner;
    private static final int NUMBER_OF_CALLS = 5;

    @BeforeEach
    public void init() {
        testRunner = TestRunners.newTestRunner(Plc4xSourceProcessor.class);
        testRunner.setIncomingConnection(false);
        testRunner.setValidateExpressionUsage(false);

        testRunner.setProperty(Plc4xSourceProcessor.PLC_CONNECTION_STRING, "${literal(\"simulated\")}://127.0.0.1");
        testRunner.setProperty(Plc4xSourceProcessor.PLC_FUTURE_TIMEOUT_MILISECONDS, "1000");

        testRunner.addConnection(Plc4xSourceProcessor.REL_SUCCESS);
        testRunner.addConnection(Plc4xSourceProcessor.REL_FAILURE);
    }

    public void testProcessor() {

        testRunner.run(NUMBER_OF_CALLS);
        testRunner.assertTransferCount(Plc4xSourceProcessor.REL_FAILURE, 0);
        testRunner.assertTransferCount(Plc4xSourceProcessor.REL_SUCCESS, NUMBER_OF_CALLS);
    }

    // Test dynamic properties addressess access strategy
    @Test
    public void testWithAddressProperties() {
        testRunner.setProperty(AddressesAccessUtils.PLC_ADDRESS_ACCESS_STRATEGY, AddressesAccessUtils.ADDRESS_PROPERTY);
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, v));
        testProcessor();
    }

    // Test addressess text property access strategy
    @Test
    public void testWithAddressText() throws JsonProcessingException { 
        testRunner.setProperty(AddressesAccessUtils.PLC_ADDRESS_ACCESS_STRATEGY, AddressesAccessUtils.ADDRESS_TEXT);
        testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY, new ObjectMapper().writeValueAsString(Plc4xCommonTest.getAddressMap()).toString());
        testProcessor();
    }

}
