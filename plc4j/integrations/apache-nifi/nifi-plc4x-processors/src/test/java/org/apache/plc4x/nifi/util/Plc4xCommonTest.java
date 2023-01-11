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

import java.io.IOException;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.function.Consumer;

import org.apache.avro.file.DataFileReader;
import org.apache.avro.file.SeekableByteArrayInput;
import org.apache.avro.generic.GenericDatumReader;
import org.apache.avro.generic.GenericRecord;
import org.apache.avro.io.DatumReader;
import org.apache.avro.util.Utf8;
import org.apache.nifi.util.MockFlowFile;

public class Plc4xCommonTest {
    public static Map<String, Object> originalMap = new HashMap<>();
    public static Map<String, String> addressMap = new HashMap<>();
    public static Map<String, Class> typeMap = new HashMap<>();

    static {
        // originalMap values are in the type needed to check type mapping between PlcType and Avro
        originalMap.put("BOOL", true);
        originalMap.put("BYTE", "\u0001");
        originalMap.put("WORD", "4");
        originalMap.put("SINT", Short.valueOf((short)-5));
        originalMap.put("USINT", "6");
        originalMap.put("INT", 2000);
        originalMap.put("UINT", "3000");
        originalMap.put("DINT", "4000");
        originalMap.put("UDINT", "5000");
        originalMap.put("DWORD", "0");
        originalMap.put("LINT", 6000L);
        originalMap.put("ULINT", "7000");
        originalMap.put("LWORD", "ab");
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

    public static void assertAvroContent(List<MockFlowFile> flowfiles, boolean checkValue, boolean checkType) {
        flowfiles.forEach(new Consumer<MockFlowFile>() {
            @Override
            public void accept(MockFlowFile t) {
                DatumReader<GenericRecord> dr = new GenericDatumReader<>();
                try (DataFileReader<GenericRecord> dfr = new DataFileReader<GenericRecord>(new SeekableByteArrayInput(t.toByteArray()), dr)) {
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
                    dfr.close();
                
                } catch (IOException e) {
                    e.printStackTrace();
                }
            }
        });
    }
}
