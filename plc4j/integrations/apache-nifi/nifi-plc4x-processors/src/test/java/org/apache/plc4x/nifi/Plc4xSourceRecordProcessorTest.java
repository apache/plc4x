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

import java.util.Map;

import org.apache.nifi.avro.AvroRecordSetWriter;
import org.apache.nifi.reporting.InitializationException;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.apache.plc4x.nifi.util.Plc4xCommonTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class Plc4xSourceRecordProcessorTest {
	
    private TestRunner testRunner;
    private static int NUMBER_OF_CALLS = 5;

    private final AvroRecordSetWriter writerService = new  AvroRecordSetWriter();
    
    @BeforeEach
    public void init() throws InitializationException {
    	testRunner = TestRunners.newTestRunner(Plc4xSourceRecordProcessor.class);
    	testRunner.setIncomingConnection(false);
    	testRunner.setValidateExpressionUsage(false);

    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_READ_FUTURE_TIMEOUT_MILISECONDS, "100");
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");

    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_SUCCESS);
    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_FAILURE);

		testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
		testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_RECORD_WRITER_FACTORY.getName(), "writer");

		for (Map.Entry<String,String> address :Plc4xCommonTest.addressMap.entrySet()) {
			// TODO: Random generation not working with this types
			if (address.getValue().startsWith("RANDOM/")) {
				if (address.getValue().endsWith("BYTE") ||
					address.getValue().endsWith("CHAR") ||
					address.getValue().endsWith("STRING"))
					continue;
			}
			testRunner.setProperty(address.getKey(), address.getValue());
		}
    }

    @Test
    public void testAvroRecordWriterProcessor() throws InitializationException {  	
    	
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validations
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);

		Plc4xCommonTest.assertAvroContent(testRunner.getFlowFilesForRelationship(Plc4xSourceProcessor.REL_SUCCESS), false, true);
    }
}
