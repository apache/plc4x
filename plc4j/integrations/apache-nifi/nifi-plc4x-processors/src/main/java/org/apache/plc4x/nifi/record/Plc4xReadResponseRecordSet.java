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
import java.time.Instant;
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
import org.apache.plc4x.java.spi.messages.DefaultPlcSubscriptionEvent;
import org.apache.plc4x.java.spi.messages.utils.ResponseItem;
import org.apache.plc4x.nifi.util.Plc4xCommon;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Plc4xReadResponseRecordSet implements RecordSet, Closeable {
    private static final Logger logger = LoggerFactory.getLogger(Plc4xReadResponseRecordSet.class);
    private final PlcReadResponse readResponse;
    private Set<String> rsColumnNames;
    private boolean moreRows;
    private final boolean debugEnabled = logger.isDebugEnabled();
    private final String timestampFieldName; 
    private boolean isSubscription = false;
    private Instant timestamp;

   	private final AtomicReference<RecordSchema> recordSchema = new AtomicReference<>(null);

    public Plc4xReadResponseRecordSet(final PlcReadResponse readResponse, RecordSchema recordSchema, String timestampFieldName) {
        this.timestampFieldName = timestampFieldName;
        this.readResponse = readResponse;
        if (!isSubscription) {
            timestamp = Instant.now();
        }
        moreRows = true;
        
        isSubscription = readResponse.getRequest() == null;

        if (debugEnabled)
            logger.debug("Creating record schema from PlcReadResponse");
        
        Map<String, ? extends PlcValue> responseDataStructure;

        responseDataStructure = !isSubscription? 
            readResponse.getAsPlcValue().getStruct():
            plc4xSubscriptionResponseRecordSet((DefaultPlcSubscriptionEvent) readResponse);
  
        rsColumnNames = responseDataStructure.keySet();
               
        if (recordSchema == null) {
        	Schema avroSchema = Plc4xCommon.createSchema(responseDataStructure, this.timestampFieldName);     	
        	this.recordSchema.set(AvroTypeUtil.createSchema(avroSchema));
        } else {
            this.recordSchema.set(recordSchema);
        }
        if (debugEnabled)
            logger.debug("Record schema from PlcReadResponse successfuly created.");

    }

    public Map<String, PlcValue> plc4xSubscriptionResponseRecordSet(final DefaultPlcSubscriptionEvent subscriptionEvent) {
        moreRows = true;
        
        if (debugEnabled)
            logger.debug("Creating record schema from DefaultPlcSubscriptionEvent");
        
        Map<String, PlcValue> responseDataStructure = new HashMap<>();

        for (Map.Entry<String, ResponseItem<PlcValue>> entry : subscriptionEvent.getValues().entrySet()) {
            responseDataStructure.put(entry.getKey(), entry.getValue().getValue());
        }

        return responseDataStructure;
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
            Record record;
            
            record = createRecord(readResponse);

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

    protected Record createRecord(final PlcReadResponse readResponse) {
        final Map<String, Object> values = new HashMap<>(getSchema().getFieldCount());

        if (debugEnabled)
            logger.debug("creating record.");

        for (final RecordField tag : getSchema().getFields()) {
            final String tagName = tag.getFieldName();

            final Object value;
            
            if (rsColumnNames.contains(tagName)) {
                if (!isSubscription) {
                    value = normalizeValue(readResponse.getAsPlcValue().getValue(tagName));
                } else {
                    value = normalizeValue(readResponse.getPlcValue(tagName));
                }
            	
            } else {
                value = null;
            }
            
            logger.trace("Adding {} tag value to record.", tagName);
            values.put(tagName, value);
        }

        //add timestamp tag to schema
        if (isSubscription) {
            values.put(timestampFieldName, ((DefaultPlcSubscriptionEvent) readResponse).getTimestamp().toEpochMilli());
        } else {
            values.put(timestampFieldName, timestamp.toEpochMilli());
        }
        
        if (debugEnabled)
            logger.debug("added timestamp tag to record.");

        	
        return new MapRecord(getSchema(), values);
    }

    private Object normalizeValue(final PlcValue value) {
        Object r = Plc4xCommon.normalizeValue(value);
        if (r != null) {
            logger.trace("Value data type: {}", r.getClass());
        }
        return r;
        
    }


}
