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
import org.apache.plc4x.java.api.PlcConnection;
import org.apache.plc4x.java.api.exceptions.PlcConnectionException;
import org.apache.plc4x.java.utils.connectionpool.PooledPlcDriverManager;
import org.apache.plc4x.nifi.Plc4xSourceRecordProcessor;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;


public class Plc4xSourceRecordProcessorTest {
	
    private TestRunner testRunner;
    private static int NUMBER_OF_CALLS = 1;
    
        @BeforeEach
    public void init() throws InitializationException {
    	
    	testRunner = TestRunners.newTestRunner(Plc4xSourceRecordProcessor.class);
    	//forzar que no tenga conexiones de entrada
    	testRunner.setIncomingConnection(false);
    	//fijo las propiedades
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_CONNECTION_STRING, "s7://10.105.143.1:102?remote-rack=0&remote-slot=0&controller-type=S7_300");
    	testRunner.setProperty(Plc4xSourceRecordProcessor.PLC_ADDRESS_STRING, "miboolean=%DB20:DBX6.0:BOOL;miint=%DB20:DBW06:INT");
    	testRunner.setProperty(Plc4xSourceRecordProcessor.FORCE_RECONNECT, "true");
    	//filo los servicios
    	final MockRecordWriter writerService = new MockRecordWriter("header", false);
    	testRunner.addControllerService("writer", writerService);
    	testRunner.enableControllerService(writerService);
    	testRunner.setProperty(Plc4xSourceRecordProcessor.RECORD_WRITER_FACTORY.getName(), "writer");
    	
    	
//    	final JsonRecordSetWriter jsonWriter = new JsonRecordSetWriter();
//    	testRunner.addControllerService(Plc4xSourceRecordProcessor.RECORD_WRITER_FACTORY.getName(), jsonWriter);
//    	testRunner.setProperty(jsonWriter, SchemaAccessUtils.SCHEMA_ACCESS_STRATEGY, SchemaAccessUtils.INHERIT_RECORD_SCHEMA);
//    	testRunner.setProperty(jsonWriter, "Pretty Print JSON", "true");
//    	testRunner.setProperty(jsonWriter, "Schema Write Strategy", "full-schema-attribute");
//    	testRunner.enableControllerService(jsonWriter);
    	
    	
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
