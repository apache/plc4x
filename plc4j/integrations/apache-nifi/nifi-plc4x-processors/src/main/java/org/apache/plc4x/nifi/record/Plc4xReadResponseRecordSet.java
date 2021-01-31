package org.apache.plc4x.nifi.record;

import org.apache.nifi.serialization.SimpleRecordSchema;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.Closeable;
import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.sql.Array;
import java.sql.SQLException;
import java.sql.Types;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import org.apache.nifi.serialization.record.DataType;
import org.apache.nifi.serialization.record.MapRecord;
import org.apache.nifi.serialization.record.Record;
import org.apache.nifi.serialization.record.RecordField;
import org.apache.nifi.serialization.record.RecordFieldType;
import org.apache.nifi.serialization.record.RecordSchema;
import org.apache.nifi.serialization.record.RecordSet;
import org.apache.plc4x.java.api.messages.PlcReadResponse;
import org.apache.plc4x.nifi.util.PLC4X_DATA_TYPE;
import org.apache.plc4x.nifi.util.PLC4X_PROTOCOL;
import org.apache.plc4x.nifi.util.Plc4xCommon;
import org.slf4j.LoggerFactory;


public class Plc4xReadResponseRecordSet implements RecordSet, Closeable {
    private static final Logger logger = LoggerFactory.getLogger(Plc4xReadResponseRecordSet.class);
    private final PlcReadResponse readResponse;
    private final RecordSchema schema;
    private final Set<String> rsColumnNames;
    private boolean moreRows;


    public Plc4xReadResponseRecordSet(final Map<String, String> plcAddressMap, final PlcReadResponse readResponse, final RecordSchema readerSchema, PLC4X_PROTOCOL PROTOCOL) throws IOException {
        this.readResponse = readResponse;
        moreRows = true;
        this.schema = createSchema(plcAddressMap, readerSchema, true, PROTOCOL);
        rsColumnNames = plcAddressMap.keySet();
        
    }

    @Override
    public RecordSchema getSchema() {
        return schema;
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
        final Map<String, Object> values = new HashMap<>(schema.getFieldCount());

        for (final RecordField field : schema.getFields()) {
            final String fieldName = field.getFieldName();

            final Object value;
            if (rsColumnNames.contains(fieldName)) {
                value = normalizeValue(readResponse.getObject(fieldName));
            } else {
                value = null;
            }

            values.put(fieldName, value);
        }

        //TODO add timestamp field to schema
        values.put(Plc4xCommon.PLC4X_RECORD_TIMESTAMP_FIELD_NAME, System.currentTimeMillis());
        	
        return new MapRecord(schema, values);
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

    private static RecordSchema createSchema(final Map<String, String> plcAddressMap, final RecordSchema readerSchema, boolean nullable, PLC4X_PROTOCOL PROTOCOL) {
        final List<RecordField> fields = new ArrayList<>(plcAddressMap.size());
        for (Map.Entry<String, String> entry : plcAddressMap.entrySet()) {
            PLC4X_DATA_TYPE plc4xType = Plc4xCommon.inferTypeFromAddressString(entry.getValue(), PROTOCOL);
            final DataType dataType = getDataType(plc4xType, readerSchema);
            final String fieldName = entry.getKey();
            final RecordField field = new RecordField(fieldName, dataType, nullable);
            fields.add(field);
        }
        
        //TODO add timestamp field to schema
        final RecordField timestampField = new RecordField(Plc4xCommon.PLC4X_RECORD_TIMESTAMP_FIELD_NAME, RecordFieldType.LONG.getDataType(), false);
        fields.add(timestampField);
        
        
        return new SimpleRecordSchema(fields);
    }

    private static DataType getDataType(final PLC4X_DATA_TYPE plc4xType, final RecordSchema readerSchema){
    	switch (plc4xType) {
            case INT:
                return RecordFieldType.INT.getDataType();
            case BOOL:
            case BIT:
            	return RecordFieldType.BOOLEAN.getDataType();
            case DECIMAL:
            	return RecordFieldType.DOUBLE.getDataType();
            case BYTE:
            	return RecordFieldType.BYTE.getDataType();
            case CHAR:
            	return RecordFieldType.CHAR.getDataType();
            case DOUBLE:
            	return RecordFieldType.DOUBLE.getDataType();
            case FLOAT:
            	return RecordFieldType.FLOAT.getDataType();
            case LONG:
            	return RecordFieldType.LONG.getDataType();
            case SMALLINT:
            	return RecordFieldType.SHORT.getDataType();
            case STRING:
            	return RecordFieldType.STRING.getDataType(); 
            case ARRAY:
            	//TODO array type -> return RecordFieldType.ARRAY.getArrayDataType(RecordFieldType.BOOLEAN.getDataType());
            	return RecordFieldType.STRING.getDataType(); 
            case SHORT: 
            	return RecordFieldType.SHORT.getDataType();
            default:
            	return RecordFieldType.STRING.getDataType();
        }
    }
}
