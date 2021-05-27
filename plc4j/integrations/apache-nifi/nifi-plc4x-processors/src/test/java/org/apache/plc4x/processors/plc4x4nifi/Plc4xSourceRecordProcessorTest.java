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
package org.apache.plc4x.processors.plc4x4nifi;
import org.apache.nifi.reporting.InitializationException;
import org.apache.nifi.serialization.record.MockRecordWriter;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.apache.plc4x.nifi.Plc4xSourceRecordProcessor;
import org.junit.jupiter.api.BeforeEach;
//import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Test;


public class Plc4xSourceRecordProcessorTest {
	
    private TestRunner testRunner;
    private static int NUMBER_OF_CALLS = 500;
    
        @BeforeEach
    public void init() throws InitializationException {
    	        	
    	testRunner = TestRunners.newTestRunner(Plc4xSourceRecordProcessor.class);
    	//forzar que no tenga conexiones de entrada
    	testRunner.setIncomingConnection(false);
    	//fijo las propiedades
    	testRunner.setValidateExpressionUsage(false);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_CONNECTION_STRING, "s7://10.105.143.7:102?remote-rack=0&remote-slot=1&controller-type=S7_1200");
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_ADDRESS_STRING, "var1=%DB1:DBX0.0:BOOL;var2=%DB1:DBX0.1:BOOL;var3=%DB1:DBB01:BYTE;var4=%DB1:DBW02:WORD;var5=%DB1:DBW04:INT");
    	
    	//filo los servicios
    	final MockRecordWriter writerService = new MockRecordWriter("header", false);
    	testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.RECORD_WRITER_FACTORY.getName(), "writer");
    	
    	//fijo las relaciones
    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_SUCCESS);
    	testRunner.addConnection(Plc4xSourceRecordProcessor.REL_FAILURE);
    }

    @Test
    public void testProcessor() {
    	
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validaciones
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSourceRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);
    }

}
