/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.apache.plc4x.java.ads.model;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;
import java.util.stream.Collectors;

// TODO: implement me acording to ads. currently copy pasta from S7
// Use endian decoders.
// TODO: replace all ifs with switches
public class AdsPlcFieldHandler extends DefaultPlcFieldHandler {

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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
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
            case TIME:
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
            case UNKNOWN:
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
            case TIME:
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
            case UNKNOWN:
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
            case TIME:
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
            case UNKNOWN:
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
            case TIME:
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
            case UNKNOWN:
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + adsField.getAdsDataType().name());
        }
    }

    private PlcValue internalEncodeBoolean(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case BOOL:
            case BYTE:
            case WORD:
            case DWORD:
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign boolean values to " + adsField.getAdsDataType().name() + " fields.");
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
                        " is not assignable to " + adsField.getAdsDataType().name() + " fields.");
            }
        }
        if(booleanValues.size() == 1) {
            return new PlcBoolean(booleanValues.get(0));
        } else {
            return new PlcList(booleanValues);
        }
    }

    private PlcValue internalEncodeInteger(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        Class<? extends PlcValue> fieldType;
        switch (adsField.getAdsDataType()) {
            case BYTE:
                fieldType = PlcList.class;
                break;
            case WORD:
                fieldType = PlcList.class;
                break;
            case DWORD:
                fieldType = PlcList.class;
                break;
            case SINT:
                fieldType = PlcInteger.class;
                break;
            case USINT:
                fieldType = PlcLong.class;
                break;
            case INT:
                fieldType = PlcInteger.class;
                break;
            case UINT:
                fieldType = PlcInteger.class;
                break;
            case DINT:
                fieldType = PlcInteger.class;
                break;
            case UDINT:
                fieldType = PlcLong.class;
                break;
            case LINT:
                fieldType = PlcLong.class;
                break;
            case ULINT:
                fieldType = PlcBigInteger.class;
                break;
            case INT32:
                fieldType = PlcInteger.class;
                break;
            case INT64:
                fieldType = PlcLong.class;
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign integer values to " + adsField.getAdsDataType().name() + " fields.");
        }
        if (fieldType == PlcLong.class) {
            Long[] longValues = new Long[values.length];
            for (int i = 0; i < values.length; i++) {
                if (!((values[i] instanceof Byte) || (values[i] instanceof Short) ||
                    (values[i] instanceof Integer) || (values[i] instanceof BigInteger) || (values[i] instanceof Long))) {
                    throw new IllegalArgumentException(
                        "Value of type " + values[i].getClass().getName() +
                            " is not assignable to " + adsField.getAdsDataType().name() + " fields.");
                }

                Number value = (Number) values[i];
                if (!adsField.getAdsDataType().withinBounds(value.doubleValue())) {
                    throw new IllegalArgumentException("Value " + values[i] + " ist not within bounds of " + adsField.getAdsDataType());
                }
                longValues[i] = value.longValue();
            }
            if(longValues.length == 1) {
                return new PlcLong(longValues[0]);
            } else {
                return new PlcList(Arrays.asList(longValues));
            }
        } else {
            BigInteger[] bigIntegerValues = new BigInteger[values.length];
            for (int i = 0; i < values.length; i++) {
                BigDecimal value;
                if (values[i] instanceof BigInteger) {
                    value = new BigDecimal((BigInteger) values[i]);
                } else if (((values[i] instanceof Byte) || (values[i] instanceof Short) ||
                    (values[i] instanceof Integer) || (values[i] instanceof Long))) {
                    value = new BigDecimal(((Number) values[i]).longValue());
                } else {
                    throw new IllegalArgumentException(
                        "Value of type " + values[i].getClass().getName() +
                            " is not assignable to " + adsField.getAdsDataType().name() + " fields.");
                }
                if (!adsField.getAdsDataType().withinBounds(value.doubleValue())) {
                    throw new IllegalArgumentException("Value " + values[i] + " ist not within bounds of " + adsField.getAdsDataType());
                }
                bigIntegerValues[i] = value.toBigInteger();
            }
            if(bigIntegerValues.length == 1) {
                return new PlcBigInteger(bigIntegerValues[0]);
            } else {
                return new PlcList(Arrays.asList(bigIntegerValues));
            }
        }
    }

    private PlcValue internalEncodeFloatingPoint(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        Class<? extends PlcValue> fieldType;
        AdsDataType adsDataType = adsField.getAdsDataType();
        switch (adsDataType) {
            case REAL:
                fieldType = PlcFloat.class;
                break;
            case LREAL:
                fieldType = PlcDouble.class;
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign floating point values to " + adsDataType.name() + " fields.");
        }
        if (fieldType == PlcDouble.class) {
            Double[] floatingPointValues = new Double[values.length];
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Float) {
                    floatingPointValues[i] = ((Float) values[i]).doubleValue();
                } else if (values[i] instanceof Double) {
                    floatingPointValues[i] = (Double) values[i];
                } else {
                    throw new IllegalArgumentException(
                        "Value of type " + values[i].getClass().getName() +
                            " is not assignable to " + adsDataType.name() + " fields.");
                }

                Number value = (Number) values[i];
                if (!adsDataType.withinBounds(value.doubleValue())) {
                    throw new IllegalArgumentException("Value " + values[i] + " ist not within bounds of " + adsDataType);
                }
            }
            if(floatingPointValues.length == 1) {
                return new PlcDouble(floatingPointValues[0]);
            } else {
                return new PlcList(Arrays.asList(floatingPointValues));
            }
        } else {
            Float[] floatingPointValues = new Float[values.length];
            for (int i = 0; i < values.length; i++) {
                if (values[i] instanceof Double) {
                    Double aDouble = (Double) values[i];
                    if (!adsDataType.withinBounds(aDouble)) {
                        throw new IllegalArgumentException(
                            "Value of " + aDouble + " exceeds allowed minimum for type "
                                + adsDataType.name() + " (min " + adsDataType.getLowerBound() + "/max +" + adsDataType.getUpperBound() + ")");
                    }
                    floatingPointValues[i] = aDouble.floatValue();
                } else if (values[i] instanceof Float) {
                    floatingPointValues[i] = (Float) values[i];
                } else {
                    throw new IllegalArgumentException(
                        "Value of type " + values[i].getClass().getName() +
                            " is not assignable to " + adsDataType.name() + " fields.");
                }

                Number value = (Number) values[i];
                if (!adsDataType.withinBounds(value.doubleValue())) {
                    throw new IllegalArgumentException("Value " + values[i] + " ist not within bounds of " + adsDataType);
                }
            }
            if(floatingPointValues.length == 1) {
                return new PlcFloat(floatingPointValues[0]);
            } else {
                return new PlcList(Arrays.asList(floatingPointValues));
            }
        }
    }

    private PlcValue internalEncodeString(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        Number maxLength = adsField.getAdsDataType().getUpperBound();
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
                if (stringValue.length() > maxLength.intValue()) {
                    throw new IllegalArgumentException(
                        "String length " + stringValue.length() + " exceeds allowed maximum for type "
                            + adsField.getAdsDataType().name() + " (max " + maxLength + ")");
                }
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
            return new PlcList(Arrays.asList(stringValues.toArray(new String[0])));
        }
    }

    private PlcValue internalTimeTemporal(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case TIME:
            case DATE:
            case DATE_AND_TIME:
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + adsField.getAdsDataType().name() + " fields.");
        }
        // TODO: support other types
        List<LocalTime> localTimeValues = Arrays.stream(values)
            .filter(LocalTime.class::isInstance)
            .map(LocalTime.class::cast)
            .collect(Collectors.toList());
        if(localTimeValues.size() == 1) {
            return new PlcTime(localTimeValues.get(0));
        } else {
            return new PlcList(Arrays.asList(localTimeValues.toArray(new LocalTime[0])));
        }
    }

    private PlcValue internalDateTemporal(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        switch (adsField.getAdsDataType()) {
            case TIME:
            case DATE:
            case DATE_AND_TIME:
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + adsField.getAdsDataType().name() + " fields.");
        }
        // TODO: support other types
        List<LocalDate> localDateValues = Arrays.stream(values)
            .filter(LocalDate.class::isInstance)
            .map(LocalDate.class::cast)
            .collect(Collectors.toList());
        if(localDateValues.size() == 1) {
            return new PlcDate(localDateValues.get(0));
        } else {
            return new PlcList(Arrays.asList(localDateValues.toArray(new LocalDate[0])));
        }
    }

    private PlcValue internalDateTimeTemporal(PlcField field, Object[] values) {
        AdsField adsField = (AdsField) field;
        Class<? extends PlcValue> fieldType;
        switch (adsField.getAdsDataType()) {
            case TIME:
                fieldType = PlcTime.class;
                break;
            case DATE:
                fieldType = PlcDate.class;
                break;
            case DATE_AND_TIME:
                fieldType = PlcDateTime.class;
                break;
            default:
                throw new IllegalArgumentException(
                    "Cannot assign temporal values to " + adsField.getAdsDataType().name() + " fields.");
        }
        // TODO: add type conversion
        if (fieldType == PlcDateTime.class) {
            final LocalDateTime[] localDateTimes = Arrays.stream(values)
                .filter(LocalDateTime.class::isInstance)
                .map(LocalDateTime.class::cast).toArray(LocalDateTime[]::new);
            if(localDateTimes.length == 1) {
                return new PlcDateTime(localDateTimes[0]);
            } else {
                return new PlcList(Arrays.asList(localDateTimes));
            }
        } else if (fieldType == PlcDate.class) {
            final LocalDate[] localDates = Arrays.stream(values)
                .filter(LocalDate.class::isInstance)
                .map(LocalDate.class::cast).toArray(LocalDate[]::new);
            if(localDates.length == 1) {
                return new PlcDate(localDates[0]);
            } else {
                return new PlcList(Arrays.asList(localDates));
            }
        } else {
            final LocalTime[] localTimes = Arrays.stream(values)
                .filter(LocalTime.class::isInstance)
                .map(LocalTime.class::cast).toArray(LocalTime[]::new);
            if(localTimes.length == 1) {
                return new PlcTime(localTimes[0]);
            } else {
                return new PlcList(Arrays.asList(localTimes));
            }
        }
    }

}
