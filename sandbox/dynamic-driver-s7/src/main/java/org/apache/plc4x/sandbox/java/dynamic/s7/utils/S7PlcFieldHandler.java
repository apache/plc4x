/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing,
 software distributed under the License is distributed on an
 "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 KIND, either express or implied.  See the License for the
 specific language governing permissions and limitations
 under the License.
 */
package org.apache.plc4x.sandbox.java.dynamic.s7.utils;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.DefaultPlcFieldHandler;
import org.apache.plc4x.java.base.messages.items.*;

import java.lang.reflect.InvocationTargetException;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;

public class S7PlcFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (S7Field.matches(fieldQuery)) {
            return S7Field.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public BaseDefaultFieldItem encodeBoolean(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        // All of these types are declared as Bit or Bit-String types.
        switch (s7Field.getDataType()) {
            case BOOL:
            case BYTE:
            case WORD:
            case DWORD:
            case LWORD:
                return internalEncodeBoolean(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeByte(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        // All of these types are declared as Bit or Bit-String types.
        switch (s7Field.getDataType()) {
            case BYTE:
            case SINT:
            case USINT:
            case CHAR:
                return internalEncodeInteger(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeShort(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case WORD:
            case INT:
            case UINT:
                return internalEncodeInteger(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeInteger(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case DWORD:
            case DINT:
            case UDINT:
                return internalEncodeInteger(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeBigInteger(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case DWORD:
            case DINT:
            case UDINT:
                return internalEncodeInteger(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeLong(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case LWORD:
            case LINT:
            case ULINT:
                return internalEncodeInteger(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeFloat(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case REAL:
                return internalEncodeFloatingPoint(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeDouble(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case LREAL:
                return internalEncodeFloatingPoint(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeString(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case CHAR:
            case WCHAR:
            case STRING:
            case WSTRING:
                return internalEncodeString(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeTime(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case TIME:
                return internalEncodeTemporal(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeDate(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case DATE:
                return internalEncodeTemporal(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public BaseDefaultFieldItem encodeDateTime(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case DATE_AND_TIME:
                return internalEncodeTemporal(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    private BaseDefaultFieldItem internalEncodeBoolean(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case BOOL:
            case BYTE:
            case WORD:
            case DWORD:
            case LWORD:
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign boolean values to " + s7Field.getDataType().name() + " fields.");
        }
        List<Boolean> booleanValues = new LinkedList<>();
        for (Object value : values) {
            if (value instanceof Boolean) {
                Boolean booleanValue = (Boolean) value;
                booleanValues.add(booleanValue);
            } else if (value instanceof Byte) {
                Byte byteValue = (Byte) value;
                BitSet bitSet = BitSet.valueOf(new byte[]{byteValue});
                for (int i = 0; i < 8; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else if (value instanceof Short) {
                Short shortValue = (Short) value;
                BitSet bitSet = BitSet.valueOf(new long[]{shortValue});
                for (int i = 0; i < 16; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else if (value instanceof Integer) {
                Integer integerValue = (Integer) value;
                BitSet bitSet = BitSet.valueOf(new long[]{integerValue});
                for (int i = 0; i < 32; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else if (value instanceof Long) {
                long longValue = (Long) value;
                BitSet bitSet = BitSet.valueOf(new long[]{longValue});
                for (int i = 0; i < 64; i++) {
                    booleanValues.add(bitSet.get(i));
                }
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + value.getClass().getName() +
                        " is not assignable to " + s7Field.getDataType().name() + " fields.");
            }
        }
        return new DefaultBooleanFieldItem(booleanValues.toArray(new Boolean[0]));
    }

    private BaseDefaultFieldItem internalEncodeInteger(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;

        // Initialize the constraints.
        BigInteger minValue;
        BigInteger maxValue;
        Class<? extends BaseDefaultFieldItem> fieldType;
        Class<?> valueType;
        Object[] castedValues;
        switch (s7Field.getDataType()) {
            case BYTE:
                minValue = BigInteger.valueOf((long) Byte.MIN_VALUE);
                maxValue = BigInteger.valueOf((long) Byte.MAX_VALUE);
                fieldType = DefaultByteFieldItem.class;
                valueType = Byte[].class;
                castedValues = new Byte[values.length];
                break;
            case WORD:
                minValue = BigInteger.valueOf((long) Short.MIN_VALUE);
                maxValue = BigInteger.valueOf((long) Short.MAX_VALUE);
                fieldType = DefaultShortFieldItem.class;
                valueType = Short[].class;
                castedValues = new Short[values.length];
                break;
            case DWORD:
                minValue = BigInteger.valueOf((long) Integer.MIN_VALUE);
                maxValue = BigInteger.valueOf((long) Integer.MAX_VALUE);
                fieldType = DefaultIntegerFieldItem.class;
                valueType = Integer[].class;
                castedValues = new Integer[values.length];
                break;
            case LWORD:
                minValue = BigInteger.valueOf(Long.MIN_VALUE);
                maxValue = BigInteger.valueOf(Long.MAX_VALUE);
                fieldType = DefaultLongFieldItem.class;
                valueType = Long[].class;
                castedValues = new Long[values.length];
                break;
            case SINT:
                minValue = BigInteger.valueOf((long) Byte.MIN_VALUE);
                maxValue = BigInteger.valueOf((long) Byte.MAX_VALUE);
                fieldType = DefaultByteFieldItem.class;
                valueType = Byte[].class;
                castedValues = new Byte[values.length];
                break;
            case USINT:
                minValue = BigInteger.valueOf((long) 0);
                maxValue = BigInteger.valueOf((long) Byte.MAX_VALUE * 2);
                fieldType = DefaultShortFieldItem.class;
                valueType = Short[].class;
                castedValues = new Short[values.length];
                break;
            case INT:
                minValue = BigInteger.valueOf((long) Short.MIN_VALUE);
                maxValue = BigInteger.valueOf((long) Short.MAX_VALUE);
                fieldType = DefaultShortFieldItem.class;
                valueType = Short[].class;
                castedValues = new Short[values.length];
                break;
            case UINT:
                minValue = BigInteger.valueOf((long) 0);
                maxValue = BigInteger.valueOf(((long) Short.MAX_VALUE) * 2);
                fieldType = DefaultIntegerFieldItem.class;
                valueType = Integer[].class;
                castedValues = new Integer[values.length];
                break;
            case DINT:
                minValue = BigInteger.valueOf((long) Integer.MIN_VALUE);
                maxValue = BigInteger.valueOf((long) Integer.MAX_VALUE);
                fieldType = DefaultIntegerFieldItem.class;
                valueType = Integer[].class;
                castedValues = new Integer[values.length];
                break;
            case UDINT:
                minValue = BigInteger.valueOf((long) 0);
                maxValue = BigInteger.valueOf(((long) Integer.MAX_VALUE) * 2);
                fieldType = DefaultLongFieldItem.class;
                valueType = Long[].class;
                castedValues = new Long[values.length];
                break;
            case LINT:
                minValue = BigInteger.valueOf(Long.MIN_VALUE);
                maxValue = BigInteger.valueOf(Long.MAX_VALUE);
                fieldType = DefaultLongFieldItem.class;
                valueType = Long[].class;
                castedValues = new Long[values.length];
                break;
            case ULINT:
                minValue = BigInteger.valueOf((long) 0);
                maxValue = BigInteger.valueOf(Long.MAX_VALUE).multiply(BigInteger.valueOf((long) 2));
                fieldType = DefaultBigIntegerFieldItem.class;
                valueType = BigInteger[].class;
                castedValues = new BigInteger[values.length];
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign integer values to " + s7Field.getDataType().name() + " fields.");
        }

        // Check the constraints
        for (int i = 0; i < values.length; i++) {
            BigInteger value;
            if (values[i] instanceof BigInteger) {
                value = (BigInteger) values[i];
            } else if ((values[i] instanceof Byte) || (values[i] instanceof Short) ||
                (values[i] instanceof Integer) || (values[i] instanceof Long)) {
                value = BigInteger.valueOf(((Number) values[i]).longValue());
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + values[i].getClass().getName() +
                        " is not assignable to " + s7Field.getDataType().name() + " fields.");
            }
            if (minValue.compareTo(value) > 0) {
                throw new IllegalArgumentException(
                    "Value of " + value.toString() + " exceeds allowed minimum for type "
                        + s7Field.getDataType().name() + " (min " + minValue.toString() + ")");
            }
            if (maxValue.compareTo(value) < 0) {
                throw new IllegalArgumentException(
                    "Value of " + value.toString() + " exceeds allowed maximum for type "
                        + s7Field.getDataType().name() + " (max " + maxValue.toString() + ")");
            }
            if (valueType == Byte[].class) {
                castedValues[i] = value.byteValue();
            } else if (valueType == Short[].class) {
                castedValues[i] = value.shortValue();
            } else if (valueType == Integer[].class) {
                castedValues[i] = value.intValue();
            } else if (valueType == Long[].class) {
                castedValues[i] = value.longValue();
            } else {
                castedValues[i] = value;
            }
        }

        // Create the field item.
        try {
            return fieldType.getDeclaredConstructor(valueType).newInstance(new Object[]{castedValues});
        } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException e) {
            throw new PlcRuntimeException("Error initializing field class " + fieldType.getSimpleName(), e);
        }
    }

    private BaseDefaultFieldItem internalEncodeFloatingPoint(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;

        // Initialize the constraints.
        Double minValue;
        Double maxValue;
        Class<? extends BaseDefaultFieldItem> fieldType;
        Class<?> valueType;
        Object[] castedValues;
        switch (s7Field.getDataType()) {
            case REAL:
                minValue = (double) -Float.MAX_VALUE;
                maxValue = (double) Float.MAX_VALUE;
                fieldType = DefaultFloatFieldItem.class;
                valueType = Float[].class;
                castedValues = new Float[values.length];
                break;
            case LREAL:
                minValue = -Double.MAX_VALUE;
                maxValue = Double.MAX_VALUE;
                fieldType = DefaultDoubleFieldItem.class;
                valueType = Double[].class;
                castedValues = new Double[values.length];
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign floating point values to " + s7Field.getDataType().name() + " fields.");
        }

        // Check the constraints
        for (int i = 0; i < values.length; i++) {
            Double value;
            if (values[i] instanceof Float) {
                value = ((Float) values[i]).doubleValue();
            } else if (values[i] instanceof Double) {
                value = (Double) values[i];
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + values[i].getClass().getName() +
                        " is not assignable to " + s7Field.getDataType().name() + " fields.");
            }
            if (value < minValue) {
                throw new IllegalArgumentException(
                    "Value of " + value + " exceeds allowed minimum for type "
                        + s7Field.getDataType().name() + " (min " + minValue.toString() + ")");
            }
            if (value > maxValue) {
                throw new IllegalArgumentException(
                    "Value of " + value + " exceeds allowed maximum for type "
                        + s7Field.getDataType().name() + " (max " + maxValue.toString() + ")");
            }
            if (valueType == Float[].class) {
                castedValues[i] = value.floatValue();
            } else {
                castedValues[i] = value;
            }
        }

        // Create the field item.
        try {
            return fieldType.getDeclaredConstructor(valueType).newInstance(new Object[]{castedValues});
        } catch (InstantiationException | IllegalAccessException | InvocationTargetException | NoSuchMethodException e) {
            throw new PlcRuntimeException("Error initializing field class " + fieldType.getSimpleName(), e);
        }
    }

    private BaseDefaultFieldItem internalEncodeString(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;

        // Initialize the constraints.
        int maxLength;
        boolean encoding16Bit;
        switch (s7Field.getDataType()) {
            case CHAR:
                maxLength = 1;
                encoding16Bit = false;
                break;
            case WCHAR:
                maxLength = 1;
                encoding16Bit = true;
                break;
            case STRING:
                maxLength = 254;
                encoding16Bit = false;
                break;
            case WSTRING:
                maxLength = 254;
                encoding16Bit = true;
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign string values to " + s7Field.getDataType().name() + " fields.");
        }

        // Check the constraints and create the strings.
        List<String> stringValues = new LinkedList<>();
        for (Object value : values) {
            if (value instanceof String) {
                String stringValue = (String) value;
                if (stringValue.length() > maxLength) {
                    throw new IllegalArgumentException(
                        "String length " + stringValue.length() + " exceeds allowed maximum for type "
                            + s7Field.getDataType().name() + " (max " + maxLength + ")");
                }
                stringValues.add(stringValue);
            }
            // All other types just translate to max one String character.
            else if (value instanceof Byte) {
                Byte byteValue = (Byte) value;
                byte[] stringBytes = new byte[]{byteValue};
                if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                }
            } else if (value instanceof Short) {
                Short shortValue = (Short) value;
                byte[] stringBytes = new byte[2];
                stringBytes[0] = (byte) (shortValue >> 8);
                stringBytes[1] = (byte) (shortValue & 0xFF);
                if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                }
            } else if (value instanceof Integer) {
                Integer integerValue = (Integer) value;
                byte[] stringBytes = new byte[4];
                stringBytes[0] = (byte) ((integerValue >> 24) & 0xFF);
                stringBytes[1] = (byte) ((integerValue >> 16) & 0xFF);
                stringBytes[2] = (byte) ((integerValue >> 8) & 0xFF);
                stringBytes[3] = (byte) (integerValue & 0xFF);
                if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                }
            } else if (value instanceof Long) {
                Long longValue = (Long) value;
                byte[] stringBytes = new byte[8];
                stringBytes[0] = (byte) ((longValue >> 56) & 0xFF);
                stringBytes[1] = (byte) ((longValue >> 48) & 0xFF);
                stringBytes[2] = (byte) ((longValue >> 40) & 0xFF);
                stringBytes[3] = (byte) ((longValue >> 32) & 0xFF);
                stringBytes[4] = (byte) ((longValue >> 24) & 0xFF);
                stringBytes[5] = (byte) ((longValue >> 16) & 0xFF);
                stringBytes[6] = (byte) ((longValue >> 8) & 0xFF);
                stringBytes[7] = (byte) (longValue & 0xFF);
                if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                }
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + value.getClass().getName() +
                        " is not assignable to " + s7Field.getDataType().name() + " fields.");
            }
        }

        // Create the field item.
        return new DefaultStringFieldItem(stringValues.toArray(new String[0]));
    }

    private BaseDefaultFieldItem internalEncodeTemporal(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case TIME:
                // TODO: I think I should implement this some time ...
            case DATE:
                // TODO: I think I should implement this some time ...
            case DATE_AND_TIME:
                return new DefaultLocalDateTimeFieldItem();
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + s7Field.getDataType().name() + " fields.");
        }
    }

}
