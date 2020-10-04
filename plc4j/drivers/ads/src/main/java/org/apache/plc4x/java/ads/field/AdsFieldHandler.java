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
package org.apache.plc4x.java.ads.field;

import org.apache.plc4x.java.ads.readwrite.types.AdsDataType;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;
import java.util.stream.Collectors;

public class AdsFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
        if (DirectAdsField.matches(fieldQuery)) {
            return DirectAdsField.of(fieldQuery);
        } else if (SymbolicAdsField.matches(fieldQuery)) {
            return SymbolicAdsField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        // All of these types are declared as Bit or Bit-String types.
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeBoolean(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeInteger(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeInteger(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeInteger(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeInteger(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeInteger(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeFloatingPoint(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeFloatingPoint(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        String[] stringValues = (String[]) values;
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BIT:
            case BIT8: {
                return internalEncode(field, values, "BOOL");
            }
            case BYTE:
            case BITARR8:
                return internalEncode(field, values, "BYTE");

            case WORD:
            case BITARR16:
                return internalEncode(field, values, "WORD");

            case DWORD:
            case BITARR32:
                return internalEncode(field, values, "DWORD");

            case SINT:
            case INT8: {
                return internalEncode(field, values, "SINT");
            }

            case USINT:
            case UINT8:
            case INT:
            case INT16: {
                return internalEncode(field, values, "INT");
            }

            case UINT:
            case UINT16:
            case DINT:
            case INT32: {
                return internalEncode(field, values, "DINT");
            }

            case UDINT:
            case UINT32:
            case LINT:
            case INT64: {
                return internalEncode(field, values, "LINT");
            }

            case ULINT:
            case UINT64: {
                return internalEncode(field, values, "ULINT");
            }

            case REAL:
            case FLOAT: {
                return internalEncode(field, values, "REAL");
            }

            case LREAL:
            case DOUBLE: {
                return internalEncode(field, values, "LREAL");
            }

            case STRING: {
                if (values.length == 1) {
                    return new PlcString(stringValues[0]);
                } else {
                    return new PlcList(Arrays.stream(stringValues).map(s -> new PlcString(s)).collect(Collectors.toList()));
                }
            }
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalEncodeString(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeTime(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BIT:
            case BIT8:
            case BITARR8:
            case BITARR16:
            case BITARR32:
            case INT8:
            case INT16:
            case INT32:
            case INT64:
            case UINT8:
            case UINT16:
            case UINT32:
            case UINT64:
            case FLOAT:
            case DOUBLE:
            case BOOL:
            case BYTE:
            case WORD:
            case DWORD:
            case SINT:
            case USINT:
            case INT:
            case UINT:
            case DINT:
            case UDINT:
            case LINT:
            case ULINT:
            case REAL:
            case LREAL:
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalTimeTemporal(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeDate(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BIT:
            case BIT8:
            case BITARR8:
            case BITARR16:
            case BITARR32:
            case INT8:
            case INT16:
            case INT32:
            case INT64:
            case UINT8:
            case UINT16:
            case UINT32:
            case UINT64:
            case FLOAT:
            case DOUBLE:
            case BOOL:
            case BYTE:
            case WORD:
            case DWORD:
            case SINT:
            case USINT:
            case INT:
            case UINT:
            case DINT:
            case UDINT:
            case LINT:
            case ULINT:
            case REAL:
            case LREAL:
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalDateTemporal(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeDateTime(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BIT:
            case BIT8:
            case BITARR8:
            case BITARR16:
            case BITARR32:
            case INT8:
            case INT16:
            case INT32:
            case INT64:
            case UINT8:
            case UINT16:
            case UINT32:
            case UINT64:
            case FLOAT:
            case DOUBLE:
            case BOOL:
            case BYTE:
            case WORD:
            case DWORD:
            case SINT:
            case USINT:
            case INT:
            case UINT:
            case DINT:
            case UDINT:
            case LINT:
            case ULINT:
            case REAL:
            case LREAL:
            case STRING:
            /*case TIME:
            case TIME_OF_DAY:
            case DATE:
            case DATE_AND_TIME:
            case ARRAY:
            case POINTER:
            case ENUM:
            case STRUCT:
            case ALIAS:
            case SUB_RANGE_DATA_TYPE:
                return internalDateTimeTemporal(field, values);
            case UNKNOWN:*/
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    private PlcValue internalEncode(PlcField field, Object[] values, String datatype) {
        AdsField adsField = (AdsField) field;
        try {
            switch (adsField.getAdsDataType().name()) {
                //Implement Custom PlcValue types here
                default:
                    return PlcValues.of(values, Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + datatype));
            }
        } catch (ClassNotFoundException e) {
            throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name() + e);
        }
    }

    private PlcValue internalEncodeString(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
//        Number maxLength = adsField.getAdsDataType().getUpperBound();
        //boolean encoding16Bit;
        switch (adsField.getAdsDataType()) {
            case STRING:
                //encoding16Bit = false;
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign string values to " + adsField.getAdsDataType().name() + " fields.");
        }
        List<String> stringValues = new LinkedList<>();
        for (Object value : values) {
            if (value instanceof String) {
                String stringValue = (String) value;
/*                if (stringValue.length() > maxLength.intValue()) {
                    throw new IllegalArgumentException(
                        "String length " + stringValue.length() + " exceeds allowed maximum for type "
                            + adsField.getAdsDataType().name() + " (max " + maxLength + ")");
                }*/
                stringValues.add(stringValue);
            }
            // All other types just translate to max one String character.
            else if (value instanceof Byte) {
                Byte byteValue = (Byte) value;
                byte[] stringBytes = new byte[]{byteValue};
                /*if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {*/
                stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                /*}*/
            } else if (value instanceof Short) {
                Short shortValue = (Short) value;
                byte[] stringBytes = new byte[2];
                stringBytes[0] = (byte) (shortValue >> 8);
                stringBytes[1] = (byte) (shortValue & 0xFF);
                /*if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {*/
                stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                /*}*/
            } else if (value instanceof Integer) {
                Integer integerValue = (Integer) value;
                byte[] stringBytes = new byte[4];
                stringBytes[0] = (byte) ((integerValue >> 24) & 0xFF);
                stringBytes[1] = (byte) ((integerValue >> 16) & 0xFF);
                stringBytes[2] = (byte) ((integerValue >> 8) & 0xFF);
                stringBytes[3] = (byte) (integerValue & 0xFF);
                /*if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {*/
                stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                /*}*/
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
                /*if (encoding16Bit) {
                    stringValues.add(new String(stringBytes, StandardCharsets.UTF_16));
                } else {*/
                stringValues.add(new String(stringBytes, StandardCharsets.UTF_8));
                /*}*/
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + value.getClass().getName() +
                        " is not assignable to " + adsField.getAdsDataType().name() + " fields.");
            }
        }
        if(stringValues.size() == 1) {
            return new PlcString(stringValues.get(0));
        } else {
            return new PlcList(stringValues);
        }
    }

    private PlcValue internalTimeTemporal(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
/*            case TIME:
            case DATE:
            case DATE_AND_TIME:
                break;*/
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + adsField.getAdsDataType().name() + " fields.");
        }
        // TODO: support other types
        /*List<LocalTime> localTimeValues = Arrays.stream(values)
            .filter(LocalTime.class::isInstance)
            .map(LocalTime.class::cast)
            .collect(Collectors.toList());
        if(localTimeValues.size() == 1) {
            return new PlcTime(localTimeValues.get(0));
        } else {
            return new PlcList(localTimeValues);
        }*/
    }

    private PlcValue internalDateTemporal(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
/*            case TIME:
            case DATE:
            case DATE_AND_TIME:
                break;*/
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + adsField.getAdsDataType().name() + " fields.");
        }
        // TODO: support other types
        /*List<LocalDate> localDateValues = Arrays.stream(values)
            .filter(LocalDate.class::isInstance)
            .map(LocalDate.class::cast)
            .collect(Collectors.toList());
        if(localDateValues.size() == 1) {
            return new PlcDate(localDateValues.get(0));
        } else {
            return new PlcList(localDateValues);
        }*/
    }

    private PlcValue internalDateTimeTemporal(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        Class<? extends PlcValue> fieldType;
        switch (adsField.getAdsDataType()) {
/*            case TIME:
                fieldType = PlcTime.class;
                break;
            case DATE:
                fieldType = PlcDate.class;
                break;
            case DATE_AND_TIME:
                fieldType = PlcDateTime.class;
                break;*/
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + adsField.getAdsDataType().name() + " fields.");
        }
        /*if(values.length == 1) {
            // TODO: add type conversion
            if (fieldType == PlcTime.class) {
                return new PlcTime((LocalTime) values[0]);
            } else if (fieldType == PlcDate.class) {
                return new PlcDate((LocalDate) values[0]);
            } else {
                return new PlcDateTime((LocalDateTime) values[0]);
            }
        } else {
            return new PlcList(Arrays.asList(values));
        }*/
    }
}
