/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.nifi;

import org.apache.nifi.avro.AvroRecordSetWriter;
import org.apache.nifi.reporting.InitializationException;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.apache.plc4x.nifi.address.AddressesAccessUtils;
import org.apache.plc4x.nifi.util.Plc4xCommonTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class Plc4xSourceRecordProcessorTest {
	
    private TestRunner testRunner;
    private static int NUMBER_OF_CALLS = 5;

    private final AvroRecordSetWriter writerService = new  AvroRecordSetWriter();
    
    @BeforeEach
    public void init() throws InitializationException {
    	testRunner = TestRunners.newTestRunner(Plc4xSourceRecordProcessor.class);
    	testRunner.setIncomingConnection(false);
    	testRunner.setValidateExpressionUsage(false);

    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_FUTURE_TIMEOUT_MILISECONDS, "100");
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");
		testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_SCHEMA_CACHE_SIZE, "1");

    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_SUCCESS);
    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_FAILURE);

		testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
		testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_RECORD_WRITER_FACTORY.getName(), "writer");
    }

    public void testAvroRecordWriterProcessor() throws InitializationException {  	
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validations
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);

		Plc4xCommonTest.assertAvroContent(testRunner.getFlowFilesForRelationship(Plc4xSourceProcessor.REL_SUCCESS), false, true);
    }

	// Test dynamic properties addressess access strategy
	@Test
    public void testWithAddressProperties() throws InitializationException {
        testRunner.setProperty(AddressesAccessUtils.PLC_ADDRESS_ACCESS_STRATEGY, AddressesAccessUtils.ADDRESS_PROPERTY);
        Plc4xCommonTest.getAddressMap().forEach((k,v) -> testRunner.setProperty(k, v));
        testAvroRecordWriterProcessor();
    }

	// Test addressess text property access strategy
    @Test
    public void testWithAddressText() throws InitializationException, JsonProcessingException { 
        testRunner.setProperty(AddressesAccessUtils.PLC_ADDRESS_ACCESS_STRATEGY, AddressesAccessUtils.ADDRESS_TEXT);
        testRunner.setProperty(AddressesAccessUtils.ADDRESS_TEXT_PROPERTY, new ObjectMapper().writeValueAsString(Plc4xCommonTest.getAddressMap()).toString());
        testAvroRecordWriterProcessor();
    }
}
