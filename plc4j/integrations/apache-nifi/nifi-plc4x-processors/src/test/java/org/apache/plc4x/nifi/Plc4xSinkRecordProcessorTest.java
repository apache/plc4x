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

import java.io.ByteArrayOutputStream;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.file.DataFileWriter;
import org.apache.avro.generic.GenericData;
import org.apache.avro.generic.GenericDatumWriter;
import org.apache.avro.generic.GenericRecord;
import org.apache.avro.io.DatumWriter;
import org.apache.nifi.avro.AvroReader;
import org.apache.nifi.reporting.InitializationException;
import org.apache.nifi.util.TestRunner;
import org.apache.nifi.util.TestRunners;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class Plc4xSinkRecordProcessorTest {
	
    private TestRunner testRunner;
    private static int NUMBER_OF_CALLS = 5;

	private final AvroReader readerService = new AvroReader();

	static Map<String, String> addressMap = Map.ofEntries(
            Map.entry("BOOL", "STATE/v1:BOOL"),
            // Map.entry("BYTE", "STATE/v2:BYTE(2)"),
            Map.entry("WORD", "STATE/v3:WORD"),
            Map.entry("SINT", "STATE/v4:SINT"),
            Map.entry("USINT", "STATE/v5:USINT"),
            Map.entry("INT", "STATE/v6:INT"),
            Map.entry("UINT", "STATE/v7:UINT"),
            Map.entry("DINT", "STATE/v8:DINT"),
            Map.entry("UDINT", "STATE/v9:UDINT"),
            Map.entry("DWORD", "STATE/v10:DWORD"),
            Map.entry("LINT", "STATE/v11:LINT"),
            Map.entry("ULINT", "STATE/v12:ULINT"),
            Map.entry("LWORD", "STATE/v13:LWORD"),
            Map.entry("REAL", "STATE/v14:REAL"),
            Map.entry("LREAL", "STATE/v15:LREAL"),
            Map.entry("CHAR", "STATE/v16:CHAR"),
            Map.entry("WCHAR", "STATE/v17:WCHAR"),
            Map.entry("STRING", "STATE/v18:STRING"));

    static Map<String, Object> originalMap = Map.ofEntries(
            Map.entry("BOOL", true),
            // Map.entry("BYTE", new short[] {1,2}),
            Map.entry("WORD", 4),
            Map.entry("SINT", -5),
            Map.entry("USINT", "6"),
            Map.entry("INT", 2000),
            Map.entry("UINT", "3000"),
            Map.entry("DINT", "4000"),
            Map.entry("UDINT", "5000"),
            Map.entry("DWORD", 0L),
            Map.entry("LINT", 6000L),
            Map.entry("ULINT", "7000"),
            Map.entry("LWORD", 0L),
            Map.entry("REAL", 1.23456F),
            Map.entry("LREAL", 2.34567),
            Map.entry("CHAR", "c"),
            Map.entry("WCHAR", "d"),
            Map.entry("STRING", "this is a string"));

	static Schema schema = SchemaBuilder.builder()
        .record("tests").fields()
            .nullableBoolean("BOOL", true)
            // .nullableBytes("BYTE", new byte[] {1,2})
            .nullableInt("WORD", 4)
            .nullableInt("SINT", -5)
            .nullableString("USINT", "6")
            .nullableInt("INT", 2000)
            .nullableString("UINT", "3000")
            .nullableString("DINT", "4000")
            .nullableString("UDINT", "5000")
            .nullableLong("DWORD", 0L)
            .nullableLong("LINT", 6000L)
            .nullableString("ULINT", "7000")
            .nullableLong("LWORD", 0L)
            .nullableFloat("REAL", 1.23456F)
            .nullableDouble("LREAL", 2.34567)
            .nullableString("CHAR", "c")
            .nullableString("WCHAR", "d")
            .nullableString("STRING", "this is a string")
        .endRecord();

	static byte[] data;
	static {
        GenericRecord record = new GenericData.Record(schema);
        for (Map.Entry<String, Object> e : originalMap.entrySet()){
            //TODO: complete this part. Byte needs to be a ByteBuffer instance
            // if (e.getKey() == "BYTE"){
            //     record.put(e.getKey(), e.getValue());
            //     continue;
            // }
            record.put(e.getKey(), e.getValue());
        }
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        DatumWriter<GenericRecord> writer = new GenericDatumWriter<GenericRecord>(schema);
        DataFileWriter<GenericRecord> fileWriter = new DataFileWriter<GenericRecord>(writer);

        try {
            fileWriter.create(schema, out);
            fileWriter.append(record);
            fileWriter.close();
        } catch (Exception e) {
            e.printStackTrace();
        }

        data = out.toByteArray();
    }

    
    @BeforeEach
    public void init() throws InitializationException {
    	testRunner = TestRunners.newTestRunner(Plc4xSinkRecordProcessor.class);
    	testRunner.setIncomingConnection(false);
    	testRunner.setValidateExpressionUsage(false);
    	testRunner.setProperty(Plc4xSinkRecordProcessor.PLC_CONNECTION_STRING, "simulated://127.0.0.1");
        testRunner.setProperty(Plc4xSinkRecordProcessor.PLC_WRITE_FUTURE_TIMEOUT_MILISECONDS, "1000");

		addressMap.forEach((k,v) -> testRunner.setProperty(k, v));

    	testRunner.addConnection(Plc4xSinkRecordProcessor.REL_SUCCESS);
    	testRunner.addConnection(Plc4xSinkRecordProcessor.REL_FAILURE);

		for (int i = 0; i<NUMBER_OF_CALLS; i++)
			testRunner.enqueue(data);
    }
    
    
    @Test
    public void testAvroRecordReaderProcessor() throws InitializationException {
    	testRunner.addControllerService("reader", readerService);
    	testRunner.enableControllerService(readerService);
    	testRunner.setProperty(Plc4xSinkRecordProcessor.PLC_RECORD_READER_FACTORY.getName(), "reader");
    	testRunner.run(NUMBER_OF_CALLS,true, true);
    	//validations
    	testRunner.assertTransferCount(Plc4xSinkRecordProcessor.REL_FAILURE, 0);
    	testRunner.assertTransferCount(Plc4xSinkRecordProcessor.REL_SUCCESS, NUMBER_OF_CALLS);
    }
}
