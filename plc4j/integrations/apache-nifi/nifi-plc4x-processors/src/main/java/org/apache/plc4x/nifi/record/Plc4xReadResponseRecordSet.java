package org.apache.plc4x.nifi.record;

import java.io.Closeable;
import java.io.IOException;
import java.util.HashMap;
import java.util.List;
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

    // TODO: review this AtomicReference?
	// TODO: this could be enhanced checking if record schema should be updated (via a cache boolean, checking property values is a nifi expression language, etc)
  	private AtomicReference<RecordSchema> recordSchema;

    public Plc4xReadResponseRecordSet(final PlcReadResponse readResponse) throws IOException {
        this.readResponse = readResponse;
        moreRows = true;
        
        logger.debug("Creating record schema from PlcReadResponse");
        Map<String, ? extends PlcValue> responseDataStructure = readResponse.getAsPlcValue().getStruct();
        rsColumnNames = responseDataStructure.keySet();
        
        if (recordSchema == null) {
        	Schema avroSchema = Plc4xCommon.createSchema(responseDataStructure); //TODO review this method as it is the 'mapping' from PlcValues to avro datatypes        	
        	recordSchema = new AtomicReference<RecordSchema>();
        	recordSchema.set(AvroTypeUtil.createSchema(avroSchema));
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

        for (final RecordField field : getSchema().getFields()) {
            final String fieldName = field.getFieldName();

            final Object value;
            
            //TODO
            if (rsColumnNames.contains(fieldName)) {
            	value = normalizeValue(readResponse.getObject(fieldName));
            } else {
                value = null;
            }
            //TODO we are asuming that record schema is always inferred from request, not writen by the user, so maybe previous lines could be changed by the following one
           // value = normalizeValue(readResponse.getObject(fieldName));
            
            logger.debug(String.format("Adding %s field value to record.", fieldName));
            values.put(fieldName, value);
        }

        //add timestamp field to schema
        values.put(Plc4xCommon.PLC4X_RECORD_TIMESTAMP_FIELD_NAME, System.currentTimeMillis());
        logger.debug("added timestamp field to record.");

        	
        return new MapRecord(getSchema(), values);
    }

    @SuppressWarnings("rawtypes")
    private Object normalizeValue(final Object value) {
        if (value == null) {
            return null;
        }
        if (value instanceof List) {
            return ((List) value).toArray();
        }
        return value;
    }


}
