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
package org.apache.plc4x.nifi;

import org.apache.nifi.avro.AvroRecordSetWriter;
import org.apache.nifi.json.JsonRecordSetWriter;
import org.apache.nifi.reporting.InitializationException;
import org.apache.nifi.serialization.record.MockRecordWriter;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Disabled;
import org.junit.jupiter.api.Test;

public class Plc4xSourceRecordProcessorTest {
	
    private TestRunner testRunner;
    private static int NUMBER_OF_CALLS = 5;
    
    @BeforeEach
    public void init() throws InitializationException {
    	testRunner = TestRunners.newTestRunner(Plc4xSourceRecordProcessor.class);
    	testRunner.setIncomingConnection(false);
    	testRunner.setValidateExpressionUsage(false);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_READ_FUTURE_TIMEOUT_MILISECONDS, "100");
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");
    	testRunner.setProperty("var1", "STATE/foo1:BOOL");
    	testRunner.setProperty("var2", "STATE/foo2:BOOL");
    	testRunner.setProperty("var3", "STATE/foo3:BYTE");
    	testRunner.setProperty("var4", "STATE/foo4:WORD");
    	testRunner.setProperty("var5", "STATE/foo5:INT");
    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_SUCCESS);
    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_FAILURE);
    }

    @Test
    public void testMockRecordWriterProcessor() throws InitializationException {
    	final MockRecordWriter writerService = new MockRecordWriter("header", false);
    	testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_RECORD_WRITER_FACTORY.getName(), "writer");
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validations
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);
    }
    
    @Test
    @Disabled("Disabled for now")
    public void testJsonRecordWriterProcessor() throws InitializationException {
    	final JsonRecordSetWriter writerService = new  JsonRecordSetWriter();
    	testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_RECORD_WRITER_FACTORY.getName(), "writer");
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validations
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);
    }
    
    @Test
    @Disabled("Disabled for now")
    public void testAvroRecordWriterProcessor() throws InitializationException {
    	final AvroRecordSetWriter writerService = new  AvroRecordSetWriter();
    	testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_RECORD_WRITER_FACTORY.getName(), "writer");
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validations
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);
    }

}
