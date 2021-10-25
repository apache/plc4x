package org.apache.plc4x.nifi.util;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.HashMap;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.nifi.avro.AvroTypeUtil;
import org.apache.nifi.serialization.record.MapRecord;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcBYTE;
import org.apache.plc4x.java.spi.values.PlcBigDecimal;
import org.apache.plc4x.java.spi.values.PlcBigInteger;
import org.junit.jupiter.api.Test;

public class Plc4xSchemaTest {
	
    @Test
    public void testSchemaTypesAndValues() {
    	
    	Map<String, PlcValue> map = new HashMap<String, PlcValue>();
    	Map<String, Object> mapValues = new HashMap<String, Object>();
    	
    	Double plcBigDecimalOriginal = 10.0;
    	PlcValue plcBigDecimal = new PlcBigDecimal(new BigDecimal(plcBigDecimalOriginal));
    	map.put("plcBigDecimal", plcBigDecimal);
    	mapValues.put("plcBigDecimal", Plc4xCommon.normalizeValue(plcBigDecimal));
    	
    	Integer plcBigIntegerOriginal = 10;
    	PlcValue plcBigInteger = new PlcBigInteger(BigInteger.valueOf(plcBigIntegerOriginal));
    	map.put("plcBigInteger", plcBigInteger);
    	mapValues.put("plcBigInteger", Plc4xCommon.normalizeValue(plcBigInteger));
    	
    	
    	Boolean plcBooleanOriginal = false;
    	PlcValue plcBoolean = new PlcBOOL(plcBooleanOriginal);
    	map.put("plcBoolean", plcBoolean);
    	mapValues.put("plcBoolean", Plc4xCommon.normalizeValue(plcBoolean));
    	
    	byte plcByteOriginal = 1;
    	PlcValue plcByte = new PlcBYTE(plcByteOriginal);
    	map.put("plcByte", plcByte);
    	mapValues.put("plcByte", Plc4xCommon.normalizeValue(plcByte));
    	
    	/*
    	PlcBigDecimal
		PlcBigInteger
		PlcBitString
		PlcBOOL
		PlcBYTE
		PlcCHAR
		PlcDATE_AND_TIME
		PlcDATE
		PlcDINT
		PlcDWORD
		PlcINT
		PlcLINT
		PlcList
		PlcLREAL
		PlcLTIME
		PlcLWORD
		PlcNull
		PlcREAL
		PlcSINT
		PlcSTRING
		PlcStruct
		PlcTIME_OF_DAY
		PlcTIME
		PlcUDINT
		PlcUINT
		PlcULINT
		PlcUSINT
		PlcWCHAR
		PlcWORD
		*/
    	
    	Schema schema = Plc4xCommon.createSchema(map);
    	MapRecord mr = new MapRecord(AvroTypeUtil.createSchema(schema), mapValues);
    	assertEquals(plcBigDecimalOriginal, mr.getAsDouble("plcBigDecimal"));
    	assertEquals(plcBigIntegerOriginal, mr.getAsInt("plcBigInteger"));
    	assertEquals(plcBooleanOriginal, mr.getAsBoolean("plcBoolean"));
    	assertEquals(plcByteOriginal, mr.getAsInt("plcByte").byteValue());
    }
    
    public static void main(String[] args) {
    	Plc4xSchemaTest a = new Plc4xSchemaTest();
    	a.testSchemaTypesAndValues();
	}
}
