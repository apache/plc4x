package org.apache.plc4x.nifi.util;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.SchemaBuilder.FieldAssembler;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcBYTE;
import org.apache.plc4x.java.spi.values.PlcBigDecimal;
import org.apache.plc4x.java.spi.values.PlcBigInteger;
import org.apache.plc4x.java.spi.values.PlcBitString;
import org.apache.plc4x.java.spi.values.PlcCHAR;
import org.apache.plc4x.java.spi.values.PlcDATE;
import org.apache.plc4x.java.spi.values.PlcDATE_AND_TIME;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcDWORD;
import org.apache.plc4x.java.spi.values.PlcIECValue;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcLREAL;
import org.apache.plc4x.java.spi.values.PlcLTIME;
import org.apache.plc4x.java.spi.values.PlcLWORD;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcNull;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcSINT;
import org.apache.plc4x.java.spi.values.PlcSTRING;
import org.apache.plc4x.java.spi.values.PlcSimpleValue;
import org.apache.plc4x.java.spi.values.PlcStruct;
import org.apache.plc4x.java.spi.values.PlcTIME;
import org.apache.plc4x.java.spi.values.PlcTIME_OF_DAY;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.spi.values.PlcULINT;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.apache.plc4x.java.spi.values.PlcValueAdapter;
import org.apache.plc4x.java.spi.values.PlcValues;
import org.apache.plc4x.java.spi.values.PlcWCHAR;
import org.apache.plc4x.java.spi.values.PlcWORD;
//TODO review remaining datatypes
import org.apache.plc4x.java.api.value.*;



public class Plc4xCommon {

	/**
	 * This method could be used to infer output AVRO schema directly from the PlcReadResponse object. 
	 * And used directly from the RecordPlc4xWriter.writePlcReadResponse() method.
	 * However, to make sure output schema does not change, it is built from the processor configuration (variable memory addresses).
	 */
	
	//TODO this variable name could be configurable in the future
	public static final String PLC4X_RECORD_TIMESTAMP_FIELD_NAME = "ts";
	
	//TODO please review this method so that it may be better way to map values 
	public static Schema createSchema(Map<String, ? extends PlcValue> responseDataStructure){
		//plc and record datatype map
		final FieldAssembler<Schema> builder = SchemaBuilder.record("PlcReadResponse").namespace("any.data").fields();	
		String fieldName = null;
		
		for (Map.Entry<String, ? extends PlcValue> entry : responseDataStructure.entrySet()) {
			 
			fieldName = entry.getKey();
			
			//TODO here many PLC values are mapped to string, just because i am not sure about which avro datatype should be used...
			if (entry.getValue() instanceof PlcBigDecimal) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().floatType().endUnion().noDefault(); 				
			}else if (entry.getValue() instanceof PlcBigInteger) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().longType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcBitString) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcBOOL) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().booleanType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcBYTE) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcCHAR) {	
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDATE_AND_TIME) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();	
			}else if (entry.getValue() instanceof PlcDATE) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcDWORD) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();				
			}else if (entry.getValue() instanceof PlcLINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcList) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcLREAL) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcLTIME) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcLWORD) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcNull) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcREAL) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().doubleType().endUnion().noDefault();		
			}else if (entry.getValue() instanceof PlcSINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().intType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcSTRING) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcStruct) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcTIME_OF_DAY) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcTIME) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcUDINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcUINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcULINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcUSINT) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcWCHAR) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();
			}else if (entry.getValue() instanceof PlcWORD) {
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();				
			}
			//TODO these should be included too?
			//}else if (entry.getValue() instanceof PlcValueAdapter) {
			//}else if (entry.getValue() instanceof PlcValues) {
			//}else if (entry.getValue() instanceof PlcIECValue<T>) {
			//}else if (entry.getValue() instanceof PlcSimpleValue<T>) {
			
			
			else { //TODO try forcing any other datatype to string...
				builder.name(fieldName).type().unionOf().nullBuilder().endNull().and().stringType().endUnion().noDefault();	
			}
			
		}
		
		//add timestamp field to schema
		builder.name(PLC4X_RECORD_TIMESTAMP_FIELD_NAME).type().longType().noDefault();
		
		
		return builder.endRecord();

	}
	

}
