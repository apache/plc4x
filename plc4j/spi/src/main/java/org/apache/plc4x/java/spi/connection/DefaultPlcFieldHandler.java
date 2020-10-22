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

package org.apache.plc4x.java.spi.connection;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;

import java.math.BigInteger;
import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;

/**
 * Base Implementation of {@link PlcFieldHandler} which throws a {@link PlcRuntimeException} for all
 * encodeXXX methods.
 */
public abstract class DefaultPlcFieldHandler implements PlcFieldHandler {

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {

        List<PlcBOOL> booleanValues = new LinkedList<>();
        for (Object value : values) {
            if (value instanceof Boolean) {
                booleanValues.add(new PlcBOOL((Boolean) value));
            } else if (value instanceof Byte) {
                Byte byteValue = (Byte) value;
                BitSet bitSet = BitSet.valueOf(new byte[]{byteValue});
                for (int i = 0; i < 8; i++) {
                    booleanValues.add(new PlcBOOL(bitSet.get(i)));
                }
            } else if (value instanceof Short) {
                Short shortValue = (Short) value;
                BitSet bitSet = BitSet.valueOf(new long[]{shortValue});
                for (int i = 0; i < 16; i++) {
                    booleanValues.add(new PlcBOOL(bitSet.get(i)));
                }
            } else if (value instanceof Integer) {
                Integer integerValue = (Integer) value;
                BitSet bitSet = BitSet.valueOf(new long[]{integerValue});
                for (int i = 0; i < 32; i++) {
                    booleanValues.add(new PlcBOOL(bitSet.get(i)));
                }
            } else if (value instanceof Long) {
                long longValue = (Long) value;
                BitSet bitSet = BitSet.valueOf(new long[]{longValue});
                for (int i = 0; i < 64; i++) {
                    booleanValues.add(new PlcBOOL(bitSet.get(i)));
                }
            } else {
                throw new IllegalArgumentException(
                    "Value of type " + value.getClass().getName() +
                        " is not assignable to " + field + " fields.");
            }
        }
        if(booleanValues.size() == 1) {
            return booleanValues.get(0);
        } else {
            return new PlcList(booleanValues);
        }
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((Byte) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((Byte) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((Byte) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((Byte) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((Byte) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Byte) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((Byte) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((Byte) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((Byte) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Byte) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((Byte) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((Byte) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((Byte) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((Byte) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((Byte) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Byte) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((Byte) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((Byte) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((Byte) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((Byte) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((Byte) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Byte) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((Byte) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((Byte) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((Byte) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((Byte) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((Byte) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Byte) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((Byte) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Byte) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((Byte) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((Byte) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((Byte) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((Byte) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                stringBytes = new byte[values.length];

                for (int i = 0; i < values.length; i++) {
                    Byte byteValue = (Byte) values[i];
                    stringBytes[(2*i)+0] = (byte) (byteValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_8));
            case "WSTRING":
            case "STRING16":
                stringBytes = new byte[values.length];

                for (int i = 0; i < values.length; i++) {
                    Byte byteValue = (Byte) values[i];
                    stringBytes[(2*i)+0] = (byte) (byteValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_8));
            default:
                if(values.length == 1) {
                    return new PlcSINT((Byte) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Byte) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
        }
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((Short) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((Short) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((Short) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((Short) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((Short) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Short) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((Short) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((Short) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((Short) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Short) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((Short) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((Short) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((Short) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((Short) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((Short) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Short) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((Short) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((Short) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((Short) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((Short) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((Short) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Short) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((Short) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((Short) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((Short) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((Short) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((Short) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Short) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((Short) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Short) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((Short) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((Short) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((Short) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((Short) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                stringBytes = new byte[2 * values.length];

                for (int i = 0; i < values.length; i++) {
                    Short shortValue = (Short) values[i];
                    stringBytes[(2*i)+0] = (byte) (shortValue >> 8);
                    stringBytes[(2*i)+1] = (byte) (shortValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_8));
            case "WSTRING":
            case "STRING16":
                stringBytes = new byte[2 * values.length];

                for (int i = 0; i < values.length; i++) {
                    Short shortValue = (Short) values[i];
                    stringBytes[(2*i)+0] = (byte) (shortValue >> 8);
                    stringBytes[(2*i)+1] = (byte) (shortValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_16));
            default:
                if(values.length == 1) {
                    return new PlcINT((Short) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Short) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
        }
    }


    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((Integer) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((Integer) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((Integer) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((Integer) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((Integer) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Integer) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((Integer) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((Integer) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((Integer) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Integer) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((Integer) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((Integer) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((Integer) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((Integer) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((Integer) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Integer) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((Integer) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((Integer) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((Integer) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((Integer) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((Integer) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Integer) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((Integer) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((Integer) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((Integer) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((Integer) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((Integer) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Integer) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((Integer) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Integer) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((Integer) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((Integer) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((Integer) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((Integer) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                stringBytes = new byte[4 * values.length];

                for (int i = 0; i < values.length; i++) {
                    Integer integerValue = (Integer) values[i];
                    stringBytes[(4*i)+0] = (byte) ((integerValue >> 24) & 0xFF);
                    stringBytes[(4*i)+1] = (byte) ((integerValue >> 16) & 0xFF);
                    stringBytes[(4*i)+2] = (byte) ((integerValue >> 8) & 0xFF);
                    stringBytes[(4*i)+3] = (byte) (integerValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_8));
            case "WSTRING":
            case "STRING16":
                    stringBytes = new byte[4 * values.length];

                    for (int i = 0; i < values.length; i++) {
                        Integer integerValue = (Integer) values[i];
                        stringBytes[(4*i)+0] = (byte) ((integerValue >> 24) & 0xFF);
                        stringBytes[(4*i)+1] = (byte) ((integerValue >> 16) & 0xFF);
                        stringBytes[(4*i)+2] = (byte) ((integerValue >> 8) & 0xFF);
                        stringBytes[(4*i)+3] = (byte) (integerValue & 0xFF);
                    }
                    return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_16));
            default:
                if(values.length == 1) {
                    return new PlcDINT((Integer) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Integer) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
        }
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((Long) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((Long) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((Long) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((Long) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((Long) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Long) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((Long) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((Long) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((Long) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Long) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((Long) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((Long) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((Long) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((Long) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((Long) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Long) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((Long) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((Long) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((Long) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((Long) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((Long) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Long) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((Long) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((Long) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((Long) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((Long) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((Long) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Long) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((Long) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Long) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((Long) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((Long) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((Long) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((Long) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                stringBytes = new byte[8 * values.length];

                for (int i = 0; i < values.length; i++) {
                    Long longValue = (Long) values[i];
                    stringBytes[(4*i)+0] = (byte) ((longValue >> 56) & 0xFF);
                    stringBytes[(4*i)+1] = (byte) ((longValue >> 48) & 0xFF);
                    stringBytes[(4*i)+2] = (byte) ((longValue >> 40) & 0xFF);
                    stringBytes[(4*i)+3] = (byte) ((longValue >> 32) & 0xFF);
                    stringBytes[(4*i)+4] = (byte) ((longValue >> 24) & 0xFF);
                    stringBytes[(4*i)+5] = (byte) ((longValue >> 16) & 0xFF);
                    stringBytes[(4*i)+6] = (byte) ((longValue >> 8) & 0xFF);
                    stringBytes[(4*i)+7] = (byte) (longValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_8));
            case "WSTRING":
            case "STRING16":
                stringBytes = new byte[8 * values.length];

                for (int i = 0; i < values.length; i++) {
                    Long longValue = (Long) values[i];
                    stringBytes[(4*i)+0] = (byte) ((longValue >> 56) & 0xFF);
                    stringBytes[(4*i)+1] = (byte) ((longValue >> 48) & 0xFF);
                    stringBytes[(4*i)+2] = (byte) ((longValue >> 40) & 0xFF);
                    stringBytes[(4*i)+3] = (byte) ((longValue >> 32) & 0xFF);
                    stringBytes[(4*i)+4] = (byte) ((longValue >> 24) & 0xFF);
                    stringBytes[(4*i)+5] = (byte) ((longValue >> 16) & 0xFF);
                    stringBytes[(4*i)+6] = (byte) ((longValue >> 8) & 0xFF);
                    stringBytes[(4*i)+7] = (byte) (longValue & 0xFF);
                }
                return new PlcSTRING(new String(stringBytes, StandardCharsets.UTF_16));
            default:
                if(values.length == 1) {
                    return new PlcLINT((Long) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Long) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
        }
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((BigInteger) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((BigInteger) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((BigInteger) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((BigInteger) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((BigInteger) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((BigInteger) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((BigInteger) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((BigInteger) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((BigInteger) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((BigInteger) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((BigInteger) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((BigInteger) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((BigInteger) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((BigInteger) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((BigInteger) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((BigInteger) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((BigInteger) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((BigInteger) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((BigInteger) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((BigInteger) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((BigInteger) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((BigInteger) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((BigInteger) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((BigInteger) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((BigInteger) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((BigInteger) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((BigInteger) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                throw new PlcRuntimeException("Invalid STRING encoder for type " + values.getClass().getName());
            case "WSTRING":
            case "STRING16":
                throw new PlcRuntimeException("Invalid WSTRING encoder for type " + values.getClass().getName());
            default:
                if(values.length == 1) {
                    return new PlcBigInteger((BigInteger) values[0]);
                } else {
                    List<PlcBigInteger> plcBigIntegerValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBigIntegerValues.add(new PlcBigInteger((BigInteger) values[i]));
                    }
                    return new PlcList(plcBigIntegerValues);
                }
        }
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((Float) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((Float) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((Float) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((Float) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((Float) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Float) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((Float) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((Float) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((Float) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Float) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((Float) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((Float) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((Float) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((Float) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((Float) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Float) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((Float) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((Float) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((Float) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((Float) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((Float) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Float) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((Float) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((Float) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((Float) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((Float) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((Float) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Float) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((Float) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Float) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((Float) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((Float) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((Float) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((Float) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                throw new PlcRuntimeException("Invalid STRING encoder for type " + values.getClass().getName());
            case "WSTRING":
            case "STRING16":
                throw new PlcRuntimeException("Invalid WSTRING encoder for type " + values.getClass().getName());
            default:
                if(values.length == 1) {
                    return new PlcREAL((Float) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Float) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
        }
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((Double) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((Double) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((Double) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((Double) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((Double) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((Double) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((Double) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((Double) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((Double) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((Double) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((Double) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((Double) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((Double) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((Double) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((Double) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((Double) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((Double) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((Double) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((Double) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((Double) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((Double) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((Double) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((Double) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((Double) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((Double) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((Double) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((Double) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((Double) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((Double) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Double) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((Double) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((Double) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((Double) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((Double) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                throw new PlcRuntimeException("Invalid STRING encoder for type " + values.getClass().getName());
            case "WSTRING":
            case "STRING16":
                throw new PlcRuntimeException("Invalid WSTRING encoder for type " + values.getClass().getName());
            default:
                if(values.length == 1) {
                    return new PlcLREAL((Double) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((Double) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
        }
    }

    @Override
    public PlcValue encodeBigDecimal(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((BigDecimal) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((BigDecimal) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((BigDecimal) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((BigDecimal) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((BigDecimal) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((BigDecimal) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((BigDecimal) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((BigDecimal) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((BigDecimal) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((BigDecimal) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((BigDecimal) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((BigDecimal) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((BigDecimal) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((BigDecimal) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((BigDecimal) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((BigDecimal) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((BigDecimal) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((BigDecimal) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((BigDecimal) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((BigDecimal) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((BigDecimal) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((BigDecimal) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((BigDecimal) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((BigDecimal) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((BigDecimal) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((BigDecimal) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((BigDecimal) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                throw new PlcRuntimeException("Invalid STRING encoder for type " + values.getClass().getName());
            case "WSTRING":
            case "STRING16":
                throw new PlcRuntimeException("Invalid WSTRING encoder for type " + values.getClass().getName());
            default:
                if(values.length == 1) {
                    return new PlcBigDecimal((BigDecimal) values[0]);
                } else {
                    List<PlcBigDecimal> plcBigDecimalValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBigDecimalValues.add(new PlcBigDecimal((BigDecimal) values[i]));
                    }
                    return new PlcList(plcBigDecimalValues);
                }
        }
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        byte[] stringBytes;
        switch (field.getPlcDataType().toUpperCase()) {
            case "BOOL":
            case "BIT":
                if(values.length == 1) {
                    return new PlcBOOL((String) values[0]);
                } else {
                    List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBOOLValues.add(new PlcBOOL((String) values[i]));
                    }
                    return new PlcList(plcBOOLValues);
                }
            case "BYTE":
            case "BITARR8":
                if(values.length == 1) {
                    return new PlcBYTE((String) values[0]);
                } else {
                    List<PlcBYTE> plcBYTEValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcBYTEValues.add(new PlcBYTE((String) values[i]));
                    }
                    return new PlcList(plcBYTEValues);
                }
            case "SINT":
            case "INT8":
                if(values.length == 1) {
                    return new PlcSINT((String) values[0]);
                } else {
                    List<PlcSINT> plcSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcSINTValues.add(new PlcSINT((String) values[i]));
                    }
                    return new PlcList(plcSINTValues);
                }
            case "USINT":
            case "UINT8":
            case "BIT8":
                if(values.length == 1) {
                    return new PlcUSINT((String) values[0]);
                } else {
                    List<PlcUSINT> plcUSINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUSINTValues.add(new PlcUSINT((String) values[i]));
                    }
                    return new PlcList(plcUSINTValues);
                }
            case "INT":
            case "INT16":
                if(values.length == 1) {
                    return new PlcINT((String) values[0]);
                } else {
                    List<PlcINT> plcINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcINTValues.add(new PlcINT((String) values[i]));
                    }
                    return new PlcList(plcINTValues);
                }
            case "UINT":
            case "UINT16":
                if(values.length == 1) {
                    return new PlcUINT((String) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUINTValues.add(new PlcUINT((String) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
            case "WORD":
            case "BITARR16":
                if(values.length == 1) {
                    return new PlcWORD((String) values[0]);
                } else {
                    List<PlcWORD> plcWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWORDValues.add(new PlcWORD((String) values[i]));
                    }
                    return new PlcList(plcWORDValues);
                }
            case "DINT":
            case "INT32":
                if(values.length == 1) {
                    return new PlcDINT((String) values[0]);
                } else {
                    List<PlcDINT> plcDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDINTValues.add(new PlcDINT((String) values[i]));
                    }
                    return new PlcList(plcDINTValues);
                }
            case "UDINT":
            case "UINT32":
                if(values.length == 1) {
                    return new PlcUDINT((String) values[0]);
                } else {
                    List<PlcUDINT> plcUDINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcUDINTValues.add(new PlcUDINT((String) values[i]));
                    }
                    return new PlcList(plcUDINTValues);
                }
            case "DWORD":
            case "BITARR32":
                if(values.length == 1) {
                    return new PlcDWORD((String) values[0]);
                } else {
                    List<PlcDWORD> plcDWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcDWORDValues.add(new PlcDWORD((String) values[i]));
                    }
                    return new PlcList(plcDWORDValues);
                }
            case "LINT":
            case "INT64":
                if(values.length == 1) {
                    return new PlcLINT((String) values[0]);
                } else {
                    List<PlcLINT> plcLINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLINTValues.add(new PlcLINT((String) values[i]));
                    }
                    return new PlcList(plcLINTValues);
                }
            case "ULINT":
            case "UINT64":
                if(values.length == 1) {
                    return new PlcULINT((String) values[0]);
                } else {
                    List<PlcULINT> plcULINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcULINTValues.add(new PlcULINT((String) values[i]));
                    }
                    return new PlcList(plcULINTValues);
                }
            case "LWORD":
            case "BITARR64":
                if(values.length == 1) {
                    return new PlcLWORD((String) values[0]);
                } else {
                    List<PlcLWORD> plcLWORDValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLWORDValues.add(new PlcLWORD((String) values[i]));
                    }
                    return new PlcList(plcLWORDValues);
                }
            case "REAL":
            case "FLOAT":
                if(values.length == 1) {
                    return new PlcREAL((String) values[0]);
                } else {
                    List<PlcREAL> plcREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcREALValues.add(new PlcREAL((String) values[i]));
                    }
                    return new PlcList(plcREALValues);
                }
            case "LREAL":
            case "DOUBLE":
                if(values.length == 1) {
                    return new PlcLREAL((String) values[0]);
                } else {
                    List<PlcLREAL> plcLREALValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcLREALValues.add(new PlcLREAL((String) values[i]));
                    }
                    return new PlcList(plcLREALValues);
                }
            case "CHAR":
                if(values.length == 1) {
                    return new PlcCHAR((String) values[0]);
                } else {
                    List<PlcCHAR> plcCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcCHARValues.add(new PlcCHAR((String) values[i]));
                    }
                    return new PlcList(plcCHARValues);
                }
            case "WCHAR":
                if(values.length == 1) {
                    return new PlcWCHAR((String) values[0]);
                } else {
                    List<PlcWCHAR> plcWCHARValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcWCHARValues.add(new PlcWCHAR((String) values[i]));
                    }
                    return new PlcList(plcWCHARValues);
                }
            case "STRING":
                if(values.length == 1) {
                    return new PlcSTRING((String) values[0]);
                } else {
                    List<PlcSTRING> plcStringValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcStringValues.add(new PlcSTRING((String) values[i]));
                    }
                    return new PlcList(plcStringValues);
                }
            case "WSTRING":
            case "STRING16":
                if(values.length == 1) {
                    return new PlcSTRING((String) values[0]);
                } else {
                    List<PlcSTRING> plcStringValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcStringValues.add(new PlcSTRING((String) values[i]));
                    }
                    return new PlcList(plcStringValues);
                }
            default:
                if(values.length == 1) {
                    return new PlcSTRING((String) values[0]);
                } else {
                    List<PlcSTRING> plcStringValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcStringValues.add(new PlcSTRING((String) values[i]));
                    }
                    return new PlcList(plcStringValues);
                }
        }
    }

    @Override
    public PlcValue encodeTime(PlcField field, Object[] values) {
        throw new PlcRuntimeException("Invalid encoder for type " + field);
    }

    @Override
    public PlcValue encodeDate(PlcField field, Object[] values) {
        throw new PlcRuntimeException("Invalid encoder for type " + field);
    }

    @Override
    public PlcValue encodeDateTime(PlcField field, Object[] values) {
        throw new PlcRuntimeException("Invalid encoder for type " + field);
    }

}
