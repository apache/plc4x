package org.apache.plc4x.nifi.util;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.SchemaBuilder.FieldAssembler;



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
}
