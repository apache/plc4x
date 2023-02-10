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
package org.apache.plc4x.nifi.record;

import java.io.Closeable;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.atomic.AtomicReference;

import org.apache.avro.Schema;
import org.apache.nifi.avro.AvroTypeUtil;
import org.apache.nifi.serialization.record.MapRecord;
import org.apache.nifi.serialization.record.Record;
import org.apache.nifi.serialization.record.RecordField;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.nifi.serialization.record.RecordSet;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.nifi.util.Plc4xCommon;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Plc4xReadResponseRecordSet implements RecordSet, Closeable {
    private static final Logger logger = LoggerFactory.getLogger(Plc4xReadResponseRecordSet.class);
    private final PlcReadResponse readResponse;
    private final Set<String> rsColumnNames;
    private boolean moreRows;

   	private AtomicReference<RecordSchema> recordSchema = new AtomicReference<RecordSchema>(null);

    public Plc4xReadResponseRecordSet(final PlcReadResponse readResponse, RecordSchema recordSchema) throws IOException {
        this.readResponse = readResponse;
        moreRows = true;
        
        logger.debug("Creating record schema from PlcReadResponse");
        Map<String, ? extends PlcValue> responseDataStructure = readResponse.getAsPlcValue().getStruct();
        rsColumnNames = responseDataStructure.keySet();
               
        if (recordSchema == null) {
        	Schema avroSchema = Plc4xCommon.createSchema(responseDataStructure);     	
        	this.recordSchema.set(AvroTypeUtil.createSchema(avroSchema));
        } else {
            this.recordSchema.set(recordSchema);
        }
        logger.debug("Record schema from PlcReadResponse successfuly created.");

    }

    
    @Override
    public RecordSchema getSchema() {
        return this.recordSchema.get();
    }

    // Protected methods for subclasses to access private member variables
    protected PlcReadResponse getReadResponse() {
        return readResponse;
    }

    protected boolean hasMoreRows() {
        return moreRows;
    }

    protected void setMoreRows(boolean moreRows) {
        this.moreRows = moreRows;
    }

    @Override
    public Record next() throws IOException {
        if (moreRows) {
             final Record record = createRecord(readResponse);
             setMoreRows(false);
             return record;
        } else {
             return null;
        }
    }

    @Override
    public void close() {
        //do nothing
    }

    protected Record createRecord(final PlcReadResponse readResponse) throws IOException{
        final Map<String, Object> values = new HashMap<>(getSchema().getFieldCount());

        logger.debug("creating record.");

        for (final RecordField tag : getSchema().getFields()) {
            final String tagName = tag.getFieldName();

            final Object value;
            
            if (rsColumnNames.contains(tagName)) {
            	value = normalizeValue(readResponse.getAsPlcValue().getValue(tagName));
            } else {
                value = null;
            }
            
            logger.trace(String.format("Adding %s tag value to record.", tagName));
            values.put(tagName, value);
        }

        //add timestamp tag to schema
        values.put(Plc4xCommon.PLC4X_RECORD_TIMESTAMP_FIELD_NAME, System.currentTimeMillis());
        logger.debug("added timestamp tag to record.");

        	
        return new MapRecord(getSchema(), values);
    }

    @SuppressWarnings("rawtypes")
    private Object normalizeValue(final PlcValue value) {
        Object r = Plc4xCommon.normalizeValue(value);
        if (r != null) {
            logger.trace("Value data type: " + r.getClass());
        }
        return r;
        
    }


}
