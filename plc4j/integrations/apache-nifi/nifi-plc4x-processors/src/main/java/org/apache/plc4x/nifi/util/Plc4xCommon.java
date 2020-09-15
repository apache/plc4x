package org.apache.plc4x.nifi.util;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.SchemaBuilder.FieldAssembler;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.*;



public class Plc4xCommon {
	
	public static final String PLC4X_RECORD_TIMESTAMP_FIELD_NAME = "ts";
	
	
	//TODO: do this using plc4x helper if exists
	public static PLC4X_DATA_TYPE inferTypeFromAddressString(String address, PLC4X_PROTOCOL protocol) {
		PLC4X_DATA_TYPE type;
		switch (protocol) {
		case S7:
			//type = PLC4X_DATA_TYPE.valueOf(address.split(":")[address.split(":").length-1]); //use each protocol's datatype conversion strategy
			type = S7_DATATYPE.getPlc4xDatatypeById(address.split(":")[address.split(":").length-1]);
			break;
		case MODBUS:
			type = MODBUS_DATATYPE.getPlc4xDatatypeById(address.split(":")[0]);
			break;
		default:
			type = null;
			break;
		}
		
		return type;
	}
	
	public static Schema createSchema(Map<String, String> plcAddressMap, PLC4X_PROTOCOL protocol){
		final FieldAssembler<Schema> builder = SchemaBuilder.record("PlcReadResponse").namespace("any.data").fields();	
		for (Map.Entry<String,String> address : plcAddressMap.entrySet()) {
		    String fieldName = address.getKey();
		    PLC4X_DATA_TYPE type = inferTypeFromAddressString(address.getValue(), protocol);
		    switch (type) {
                case CHAR:
                case STRING:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
                    break;
                case BOOL:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().booleanType().endUnion().noDefault();
                    break;
                case INT:
                case SHORT:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();
                    break;
                case LONG:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().longType().endUnion().noDefault();
                    break;
                case FLOAT:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().floatType().endUnion().noDefault();
                    break;
                case DOUBLE:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().doubleType().endUnion().noDefault();
                    break;
                case BYTE:
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().bytesType().endUnion().noDefault();
                    break;
                case ARRAY:
                	//TODO array type -> builder.name(fieldName).type().array().items().booleanType().noDefault();
                    builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
                    break;
                default:
                    throw new IllegalArgumentException("createSchema: Unknown AVRO type " + type +") cannot be converted to Avro type");
            }
        }
        return builder.endRecord();
	}
	
	

	/**
	 * This method could be used to infer output AVRO schema directly from the PlcReadResponse object. 
	 * And used directly from the RecordPlc4xWriter.writePlcReadResponse() method.
	 * However, to make sure output schema does not change, it is built from the processor configuration (variable memory addresses).
	 */
	public static Schema createSchema(Map<String, ? extends PlcValue> responseDataStructure){
		//plc and record datatype map
		final FieldAssembler<Schema> builder = SchemaBuilder.record("PlcReadResponse").namespace("any.data").fields();	
		String fieldName = null;
		
		for (Map.Entry<String, ? extends PlcValue> entry : responseDataStructure.entrySet()) {
			 
			fieldName = entry.getKey();
			
			if (entry.getValue() instanceof PlcBigDecimal) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().floatType().endUnion().noDefault(); //TODO correct map for bigdecimal?
			}else if (entry.getValue() instanceof PlcBoolean) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().booleanType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcBigInteger) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().longType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcByte) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault(); //TODO correct map for bigdecimal?
			}else if (entry.getValue() instanceof PlcDate) {
				//TODO avro spec: {"type": "int","logicalType": "date"}
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDateTime) {
				//TODO avro spec: {"type": "int","logicalType": "date"}
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDouble) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().doubleType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcFloat) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().floatType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcInteger) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcList) {
				//builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
				builder.name(fieldName).type().array().items().booleanType();
			}else if (entry.getValue() instanceof PlcLong) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().longType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcNull) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcShort) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcString) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcStruct) {
				//TODO
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcTime) {
				//TODO
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}
							
		}
		
		return builder.endRecord();

	}
	
	
	public static PLC4X_PROTOCOL getConnectionProtocol(String connectionString) {
		if (connectionString != null && !connectionString.isEmpty()) {
			String protocolId = connectionString.split(":")[0];
			PLC4X_PROTOCOL protocol = PLC4X_PROTOCOL.qualifierValueOf(protocolId);
			return protocol;
		}
		return null;
	}
	
}
