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
package org.apache.plc4x.nifi.util;

import java.util.List;
import java.util.Map;

import org.apache.avro.Schema;
import org.apache.avro.SchemaBuilder;
import org.apache.avro.SchemaBuilder.BaseTypeBuilder;
import org.apache.avro.SchemaBuilder.FieldAssembler;
import org.apache.avro.SchemaBuilder.NullDefault;
import org.apache.avro.SchemaBuilder.UnionAccumulator;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.values.PlcBOOL;
import org.apache.plc4x.java.spi.values.PlcBYTE;
import org.apache.plc4x.java.spi.values.PlcCHAR;
import org.apache.plc4x.java.spi.values.PlcDATE;
import org.apache.plc4x.java.spi.values.PlcDATE_AND_TIME;
import org.apache.plc4x.java.spi.values.PlcDINT;
import org.apache.plc4x.java.spi.values.PlcDWORD;
import org.apache.plc4x.java.spi.values.PlcINT;
import org.apache.plc4x.java.spi.values.PlcLINT;
import org.apache.plc4x.java.spi.values.PlcLREAL;
import org.apache.plc4x.java.spi.values.PlcLWORD;
import org.apache.plc4x.java.spi.values.PlcList;
import org.apache.plc4x.java.spi.values.PlcREAL;
import org.apache.plc4x.java.spi.values.PlcSINT;
import org.apache.plc4x.java.spi.values.PlcTIME;
import org.apache.plc4x.java.spi.values.PlcTIME_OF_DAY;
import org.apache.plc4x.java.spi.values.PlcUDINT;
import org.apache.plc4x.java.spi.values.PlcUINT;
import org.apache.plc4x.java.spi.values.PlcULINT;
import org.apache.plc4x.java.spi.values.PlcUSINT;
import org.apache.plc4x.java.spi.values.PlcWCHAR;
import org.apache.plc4x.java.spi.values.PlcWORD;

public class Plc4xCommon {

	/**
	 * This method is used to infer output AVRO schema directly from the PlcReadResponse object. 
	 * It is directly used from the RecordPlc4xWriter.writePlcReadResponse() method.
	 * However, to make sure output schema does not change, it is built from the processor configuration (variable memory addresses).
	 * 
	 * At the moment this method does not handle the following Object Types: PlcValueAdapter, PlcIECValue<T>, PlcSimpleValue<T>
	 * 
	 * @param responseDataStructure: a map that reflects the structure of the answer given by the PLC when making a Read Request.
	 * @return AVRO Schema built from responseDataStructure.
	 */
	public static Schema createSchema(Map<String, ? extends PlcValue> responseDataStructure, String timestampFieldName){
		//plc and record datatype map
		final FieldAssembler<Schema> builder = SchemaBuilder.record("PlcReadResponse").namespace("any.data").fields();	
		String fieldName = null;
		
		for (Map.Entry<String, ? extends PlcValue> entry : responseDataStructure.entrySet()) {
			fieldName = entry.getKey();
			PlcValue value = entry.getValue();
			BaseTypeBuilder<UnionAccumulator<NullDefault<Schema>>> fieldBuilder = 
				builder.name(fieldName).type().unionOf().nullType().and();
			
			if (value instanceof PlcList) {
				if(!value.getList().isEmpty()) {
					fieldBuilder = fieldBuilder.array().items();
					value = value.getList().get(0);
				}
			}

			// PlcTYPEs not in here are casted to avro string type.
			UnionAccumulator<NullDefault<Schema>> buildedField = null;
			if (value instanceof PlcBOOL) {
				buildedField = fieldBuilder.booleanType();
			}else if (value instanceof PlcBYTE) {
				buildedField = fieldBuilder.bytesType();
			}else if (value instanceof PlcINT) {
				buildedField = fieldBuilder.intType();				
			}else if (value instanceof PlcLINT) {
				buildedField = fieldBuilder.longType();
			}else if (value instanceof PlcLREAL) {
				buildedField = fieldBuilder.doubleType();
			}else if (value instanceof PlcREAL) {
				buildedField = fieldBuilder.floatType();		
			}else if (value instanceof PlcSINT) {
				buildedField = fieldBuilder.intType();		
			}else  {// Default to string:
				fieldBuilder.stringType().endUnion().nullDefault();
				continue;// In case of null default continue
			}
			buildedField.endUnion().noDefault();
		}
		
		//add timestamp tag to schema
		builder.name(timestampFieldName).type().longType().noDefault();
		
		
		return builder.endRecord();
	}
	
	
	private static Object normalizeBasicTypes(final Object valueOriginal) {
		if (valueOriginal == null) 
			return null;
			
		if (valueOriginal instanceof PlcValue) {
			PlcValue value = (PlcValue) valueOriginal;
			// 8 bits
			if (value instanceof PlcBOOL && value.isBoolean())
				return value.getBoolean();
			if (value instanceof PlcBYTE && (value.isByte() || value.isShort()))
				return new byte[]{value.getByte()};
			if (value instanceof PlcCHAR && value.isShort())
				return value.getString();
			if ((value instanceof PlcSINT || value instanceof PlcUSINT) && value.isShort())
				return value.getShort();


			// 16 bits
			if (value instanceof PlcWORD && (value.isInteger() || value.isShort()))
				return value.getString();
			if (value instanceof PlcINT && value.isInteger())
				return value.getInteger();
			if (value instanceof PlcUINT && value.isInteger())
				return value.getInteger();
			if ((value instanceof PlcWCHAR || value instanceof PlcDWORD) && value.isInteger())
				return value.getString();

			// 32 bits
			if (value instanceof PlcREAL && value.isFloat())
				return value.getFloat();
			if ((value instanceof PlcDINT || value instanceof PlcUDINT) && value.isInteger())
				return value.getInteger();
			if (value instanceof PlcDWORD && value.isInteger())
				return value.getString();
			
			// 64 bits
			if ((value instanceof PlcLINT || value instanceof PlcULINT) && value.isLong())
				return value.getLong();
			if (value instanceof PlcLREAL && value.isDouble())
				return value.getDouble();
			if (value instanceof PlcLWORD && (value.isLong() || value.isBigInteger()))
				return value.getString();

			// Dates and time
			if (value instanceof PlcDATE && value.isDate())
				return value.getDate();
			if (value instanceof PlcDATE_AND_TIME && value.isDateTime())
				return value.getDateTime();
			if (value instanceof PlcTIME && value.isTime())
				return value.getTime();
			if (value instanceof PlcTIME_OF_DAY && value.isTime())
				return value.getTime();

			// Everything else to string
			return value.getString();
		} 
		return valueOriginal;
	}
	
	public static Object normalizeValue(final Object valueOriginal) {
        if (valueOriginal == null) {
            return null;
        }
        if (valueOriginal instanceof List) {
            return ((List<?>) valueOriginal).toArray();
        } else  if (valueOriginal instanceof PlcValue) {
			PlcValue value = (PlcValue) valueOriginal;

			if (value.isList() && value instanceof PlcList) {
	        	Object[] r = new Object[value.getList().size()];
	        	int i = 0;
	        	for (Object element : value.getList()) {
	        		r[i] =  normalizeBasicTypes(element);
	        		i++;
				}
	        	return r;
	        } 	
			return normalizeBasicTypes(value);
        } else {
        	return valueOriginal;
        }
    }
	
}

