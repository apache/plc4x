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

package org.apache.plc4x.java.api.value;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.PlcField;

import java.math.BigInteger;
import java.math.BigDecimal;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.nio.charset.StandardCharsets;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;


public class IEC61131ValueHandler implements PlcValueHandler {


    public PlcValue of(Object value) {
        return newPlcValue(new Object[] {value});
    }

    public PlcValue of(Object[] values) {
        return newPlcValue(values);
    }

    public static PlcValue newPlcValue(Object value) {
        return newPlcValue(new Object[] {value});
    }


    public static PlcValue newPlcValue(Object[] values) {
        if (values.length == 1) {
            Object value = values[0];
            if (value instanceof Boolean) {
                return PlcBOOL.of(value);
            } else if (value instanceof Byte) {
                return PlcSINT.of(value);
            } else if (value instanceof Short) {
                return PlcINT.of(value);
            } else if (value instanceof Integer) {
                return PlcDINT.of(value);
            } else if (value instanceof Long) {
                return PlcLINT.of(value);
            } else if (value instanceof BigInteger) {
                return new PlcBigInteger((BigInteger) value);
            } else if (value instanceof Float) {
                return PlcREAL.of(value);
            } else if (value instanceof Double) {
                return PlcLREAL.of(value);
            } else if (value instanceof BigDecimal) {
                return new PlcBigDecimal((BigDecimal) value);
            } else if (value instanceof LocalTime) {
                return new PlcTime((LocalTime) value);
            } else if (value instanceof LocalDate) {
                return new PlcDate((LocalDate) value);
            } else if (value instanceof LocalDateTime) {
                return new PlcDateTime((LocalDateTime) value);
            } else {
                throw new PlcUnsupportedDataTypeException("Data Type " + value.getClass()
                    + "Is not supported");
            }
        } else {
            PlcList list = new PlcList();
            for (Object value : values) {
                list.add(newPlcValue(new Object[] {value}));
            }
            return list;
        }
    }


    public static PlcValue newPlcValue(PlcField field, Object value) {
        return newPlcValue(field, new Object[] {value});
    }


    public static PlcValue newPlcValue(PlcField field, Object[] values) {
        if(values.length == 1) {
            Object value = values[0];
            switch (field.getPlcDataType().toUpperCase()) {
                case "IEC61131_BOOL":
                case "IEC61131_BIT":
                    return PlcBOOL.of(value);
                case "IEC61131_BYTE":
                case "IEC61131_BITARR8":
                    return PlcBYTE.of(value);
                case "IEC61131_SINT":
                case "IEC61131_INT8":
                    return PlcSINT.of(value);
                case "IEC61131_USINT":
                case "IEC61131_UINT8":
                case "IEC61131_BIT8":
                    return PlcUSINT.of(value);
                case "IEC61131_INT":
                case "IEC61131_INT16":
                    return PlcINT.of(value);
                case "IEC61131_UINT":
                case "IEC61131_UINT16":
                    return PlcUINT.of(value);
                case "IEC61131_WORD":
                case "IEC61131_BITARR16":
                    return PlcWORD.of(value);
                case "IEC61131_DINT":
                case "IEC61131_INT32":
                    return PlcDINT.of(value);
                case "IEC61131_UDINT":
                case "IEC61131_UINT32":
                    return PlcUDINT.of(value);
                case "IEC61131_DWORD":
                case "IEC61131_BITARR32":
                    return PlcDWORD.of(value);
                case "IEC61131_LINT":
                case "IEC61131_INT64":
                    return PlcLINT.of(value);
                case "IEC61131_ULINT":
                case "IEC61131_UINT64":
                    return PlcULINT.of(value);
                case "IEC61131_LWORD":
                case "IEC61131_BITARR64":
                    return PlcLWORD.of(value);
                case "IEC61131_REAL":
                case "IEC61131_FLOAT":
                    return PlcREAL.of(value);
                case "IEC61131_LREAL":
                case "IEC61131_DOUBLE":
                    return PlcLREAL.of(value);
                case "IEC61131_CHAR":
                    return PlcCHAR.of(value);
                case "IEC61131_WCHAR":
                    return PlcWCHAR.of(value);
                case "IEC61131_STRING":
                    return PlcSTRING.of(value);
                case "IEC61131_WSTRING":
                case "IEC61131_STRING16":
                    return PlcSTRING.of(value);
                default:
                    return customDataType(field, new Object[] {value});
            }
        } else {
            PlcList list = new PlcList();
            for (Object value : values) {
                list.add(newPlcValue(field, new Object[] {value}));
            }
            return list;
        }
    }

    public static PlcValue customDataType(PlcField field, Object[] values) {
        throw new PlcUnsupportedDataTypeException("Data Type " + field.getPlcDataType()
            + "Is not supported");
    }
}
