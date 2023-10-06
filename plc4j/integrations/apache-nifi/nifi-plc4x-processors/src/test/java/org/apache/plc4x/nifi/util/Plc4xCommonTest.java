/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.apache.plc4x.nifi.util;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.file.DataFileReader;
import org.apache.avro.file.DataFileWriter;
import org.apache.avro.file.SeekableByteArrayInput;
import org.apache.avro.generic.GenericData;
import org.apache.avro.generic.GenericDatumReader;
import org.apache.avro.generic.GenericDatumWriter;
import org.apache.avro.generic.GenericRecord;
import org.apache.avro.io.DatumReader;
import org.apache.avro.io.DatumWriter;
import org.apache.avro.util.Utf8;
import org.apache.nifi.util.MockFlowFile;

public class Plc4xCommonTest {
    public static final Map<String, Object> originalMap = new HashMap<>();
    public static final Map<String, String> addressMap = new HashMap<>();
    public static final Map<String, Class<?>> typeMap = new HashMap<>();


    // TODO: BOOL, WORD; DWORD and LWORD are commented because random generation is not working with this types 
    // or a because a reverse type mapping between avro and PlcTypes is not implemented
    public static final Schema schema = SchemaBuilder.builder()
        .record("tests").fields()
            .nullableBoolean("BOOL", true)
            // .nullableBytes("BYTE", new byte[] {1,2})
            // .nullableString("WORD", "4")
            .nullableInt("SINT", -5)
            .nullableString("USINT", "6")
            .nullableInt("INT", 2000)
            .nullableString("UINT", "3000")
            .nullableString("DINT", "4000")
            .nullableString("UDINT", "5000")
            // .nullableString("DWORD", "0")
            .nullableLong("LINT", 6000L)
            .nullableString("ULINT", "7000")
            // .nullableString("LWORD", "0")
            .nullableFloat("REAL", 1.23456F)
            .nullableDouble("LREAL", 2.34567)
            .nullableString("CHAR", "c")
            .nullableString("WCHAR", "d")
            .nullableString("STRING", "this is a string")
        .endRecord();

    static {
        // originalMap values are in the type needed to check type mapping between PlcType and Avro
        originalMap.put("BOOL", true);
        originalMap.put("BYTE", "\u0001");
        originalMap.put("WORD", "4");
        originalMap.put("SINT", -5);
        originalMap.put("USINT", "6");
        originalMap.put("INT", 2000);
        originalMap.put("UINT", "3000");
        originalMap.put("DINT", "4000");
        originalMap.put("UDINT", "5000");
        originalMap.put("DWORD", Long.valueOf("0"));
        originalMap.put("LINT", 6000L);
        originalMap.put("ULINT", "7000");
        originalMap.put("LWORD", Long.valueOf("0"));
        originalMap.put("REAL", 1.23456F);
        originalMap.put("LREAL", 2.34567);
        originalMap.put("CHAR", "c");
        originalMap.put("WCHAR", "d");
        originalMap.put("STRING", "this is a string");

        addressMap.put("BOOL", "RANDOM/v1:BOOL");
        addressMap.put("BYTE", "RANDOM/v2:BYTE");
        addressMap.put("WORD", "RANDOM/v3:WORD");
        addressMap.put("SINT", "RANDOM/v4:SINT");
        addressMap.put("USINT", "RANDOM/v5:USINT");
        addressMap.put("INT", "RANDOM/v6:INT");
        addressMap.put("UINT", "RANDOM/v7:UINT");
        addressMap.put("DINT", "RANDOM/v8:DINT");
        addressMap.put("UDINT", "RANDOM/v9:UDINT");
        addressMap.put("DWORD", "RANDOM/v10:DWORD");
        addressMap.put("LINT", "RANDOM/v11:LINT");
        addressMap.put("ULINT", "RANDOM/v12:ULINT");
        addressMap.put("LWORD", "RANDOM/v13:LWORD");
        addressMap.put("REAL", "RANDOM/v14:REAL");
        addressMap.put("LREAL", "RANDOM/v15:LREAL");
        addressMap.put("CHAR", "RANDOM/v16:CHAR");
        addressMap.put("WCHAR", "RANDOM/v17:WCHAR");
        addressMap.put("STRING", "RANDOM/v18:STRING");

        typeMap.put("BOOL", Boolean.class);
        typeMap.put("BYTE", ByteBuffer.class);
        typeMap.put("WORD", Utf8.class);
        typeMap.put("SINT", Integer.class);
        typeMap.put("USINT", Utf8.class);
        typeMap.put("INT", Integer.class);
        typeMap.put("UINT", Utf8.class);
        typeMap.put("DINT", Utf8.class);
        typeMap.put("UDINT", Utf8.class);
        typeMap.put("DWORD", Utf8.class);
        typeMap.put("LINT", Long.class);
        typeMap.put("ULINT", Utf8.class);
        typeMap.put("LWORD", Utf8.class);
        typeMap.put("REAL", Float.class);
        typeMap.put("LREAL", Double.class);
        typeMap.put("CHAR", Utf8.class);
        typeMap.put("WCHAR", Utf8.class);
        typeMap.put("STRING", Utf8.class);
    }

    public static Map<String, String> getAddressMap(){
        Map<String, String> result = new HashMap<>();

        addressMap.forEach((k,v) -> {
			if (v.startsWith("RANDOM/")) {
				if (!v.endsWith("BYTE") &&
					!v.endsWith("CHAR") &&
                    !v.endsWith("WORD") &&
					!v.endsWith("STRING"))
					result.put(k, v);
			} else {
                result.put(k, v);
            }

		});
        return result;
    }

    public static void assertAvroContent(List<MockFlowFile> flowfiles, boolean checkValue, boolean checkType) {
        flowfiles.forEach(t -> {
            DatumReader<GenericRecord> dr = new GenericDatumReader<>();
            try (DataFileReader<GenericRecord> dfr = new DataFileReader<>(new SeekableByteArrayInput(t.toByteArray()), dr)) {
                GenericRecord data = null;
                while (dfr.hasNext()) {
                    data = dfr.next(data);

                    for (String tag : Plc4xCommonTest.addressMap.keySet()) {
                        if (data.hasField(tag)) {
                            // Check value after string conversion
                            if (checkValue)
                                assert data.get(tag).toString().equalsIgnoreCase(Plc4xCommonTest.originalMap.get(tag).toString());

                            // Check type
                            if (checkType)
                                assert data.get(tag).getClass().equals(Plc4xCommonTest.typeMap.get(tag));
                        }
                    }
                }
            } catch (IOException e) {
                e.printStackTrace();
            }
        });
    }

    public static GenericRecord getTestRecord() {
        GenericRecord record = new GenericData.Record(schema);
        record.put("BOOL", true);
        // record.put("BYTE", "\u0001");
        // record.put("WORD", "4");
        record.put("SINT", -5);
        record.put("USINT", "6");
        record.put("INT", 2000);
        record.put("UINT", "3000");
        record.put("DINT", "4000");
        record.put("UDINT", "5000");
        // record.put("DWORD", "0");
        record.put("LINT", 6000L);
        record.put("ULINT", "7000");
        // record.put("LWORD", "0");
        record.put("REAL", 1.23456F);
        record.put("LREAL", 2.34567);
        record.put("CHAR", "c");
        record.put("WCHAR", "d");
        record.put("STRING", "this is a string");
        return record;
    }

    public static byte[] encodeRecord(GenericRecord record){
        ByteArrayOutputStream out = new ByteArrayOutputStream();
        DatumWriter<GenericRecord> writer = new GenericDatumWriter<>(schema);
        DataFileWriter<GenericRecord> fileWriter = new DataFileWriter<>(writer);

        try {
            fileWriter.create(schema, out);
            fileWriter.append(record);
            fileWriter.close();
        } catch (Exception e) {
            e.printStackTrace();
        }

        return out.toByteArray();
    }
}
