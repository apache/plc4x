package org.apache.plc4x.nifi.util;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.SchemaBuilder.FieldAssembler;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.*;



public class Plc4xCommon {
	
	//TODO: do this using plc4x helper if exists
	public static PLC4X_DATA_TYPE inferTypeFromAddressString(String address, PLC4X_PROTOCOL protocol) {
		PLC4X_DATA_TYPE type;
		switch (protocol) {
		case S7:
			type = PLC4X_DATA_TYPE.valueOf(address.split(":")[address.split(":").length-1]);
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
                default:
                    throw new IllegalArgumentException("createSchema: Unknown AVRO type " + type +") cannot be converted to Avro type");
            }
        }
        return builder.endRecord();
	}
	
	

	public static Schema createSchema(Map<String, ? extends PlcValue> responseDataStructure){
		//plc and record datatype map
		final FieldAssembler<Schema> builder = SchemaBuilder.record("PlcReadResponse").namespace("any.data").fields();	
		String debugDatatype = "";
		String fieldName = null;
		
		for (Map.Entry<String, ? extends PlcValue> entry : responseDataStructure.entrySet()) {
			 
			fieldName = entry.getKey();
			
			if (entry.getValue() instanceof PlcBigDecimal) {
				debugDatatype = "big decimal";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().floatType().endUnion().noDefault(); //TODO correct map for bigdecimal?
			}else if (entry.getValue() instanceof PlcBoolean) {
				debugDatatype = "bool";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().booleanType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcBigInteger) {
				debugDatatype = "big integer";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().longType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcByte) {
				debugDatatype = "byte";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault(); //TODO correct map for bigdecimal?
			}else if (entry.getValue() instanceof PlcDate) {
				debugDatatype = "date";
				//TODO avro spec: {"type": "int","logicalType": "date"}
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDateTime) {
				debugDatatype = "datetime";
				//TODO avro spec: {"type": "int","logicalType": "date"}
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDouble) {
				debugDatatype = "double";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().doubleType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcFloat) {
				debugDatatype = "float";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().floatType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcInteger) {
				debugDatatype = "integer";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcList) {
				debugDatatype = "list";
				//builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
				builder.name(fieldName).type().array().items().booleanType();
			}else if (entry.getValue() instanceof PlcLong) {
				debugDatatype = "Long";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().longType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcNull) {
				debugDatatype = "null";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcShort) {
				debugDatatype = "short";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcString) {
				debugDatatype = "string";
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcStruct) {
				debugDatatype = "struct";
				//TODO
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcTime) {
				debugDatatype = "time";
				//TODO
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}
				
				System.out.println("DEBUG:: "+fieldName+" instanceof "+debugDatatype);
			
		}
		
		return builder.endRecord();

	}
	
}
