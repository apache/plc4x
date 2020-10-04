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
package org.apache.plc4x.java.s7.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.lang.reflect.InvocationTargetException;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
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
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case CHAR:
            case WCHAR:
            case STRING:
            case WSTRING:
                return internalEncodeString(field, values);
            default:
                return internalEncode(field, values);
        }
    }

    @Override
    public PlcValue encodeTime(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case TIME:
                return internalEncodeTemporal(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public PlcValue encodeDate(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case DATE:
                return internalEncodeTemporal(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    @Override
    public PlcValue encodeDateTime(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case DATE_AND_TIME:
                return internalEncodeTemporal(field, values);
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name());
        }
    }

    private PlcValue internalEncode(PlcField field, Object[] values) {
        S7Field s7Field = (S7Field) field;
        try {
            switch (s7Field.getDataType().name()) {
                //Implement Custom PlcValue types here
                default:
                    return PlcValues.of(values, Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + s7Field.getDataType().name()));
            }
        } catch (ClassNotFoundException e) {
            throw new PlcRuntimeException("Invalid encoder for type " + s7Field.getDataType().name() + e);
        }
    }

    private PlcValue internalEncodeString(PlcField field, Object[] values) {
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
        if(stringValues.size() == 1) {
            return new PlcString(stringValues.get(0));
        } else {
            return new PlcList(stringValues);
        }
    }

    private PlcValue internalEncodeTemporal(PlcField field, Object[] values) {
        if(values.length > 1) {
            return new PlcList(Arrays.asList(values));
        }
        S7Field s7Field = (S7Field) field;
        switch (s7Field.getDataType()) {
            case TIME:
                return new PlcTime((LocalTime) values[0]);
            case DATE:
                return new PlcDate((LocalDate) values[0]);
            case DATE_AND_TIME:
                return new PlcDateTime((LocalDateTime) values[0]);
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + s7Field.getDataType().name() + " fields.");
        }
    }

}
