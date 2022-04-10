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
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.api.value.PlcValueHandler;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;

public class IEC61131ValueHandler implements PlcValueHandler {


    public PlcValue newPlcValue(Object value) {
        return of(new Object[]{value});
    }

    public PlcValue newPlcValue(Object[] values) {
        return of(values);
    }

    public PlcValue newPlcValue(PlcField field, Object value) {
        return of(field, new Object[]{value});
    }

    public PlcValue newPlcValue(PlcField field, Object[] values) {
        return of(field, values);
    }

    public static PlcValue of(Object value) {
        return of(new Object[]{value});
    }

    public static PlcValue of(List<?> value) {
        return of(value.toArray());
    }

    public static PlcValue of(Object[] values) {
        if (values.length != 1) {
            PlcList list = new PlcList();
            for (Object value : values) {
                list.add(of(new Object[]{value}));
            }
            return list;
        }
        Object value = values[0];
        if (value instanceof Boolean) {
            return PlcBOOL.of(value);
        }
        if (value instanceof Byte) {
            return PlcSINT.of(value);
        }
        if (value instanceof byte[]) {
            return PlcByteArray.of(value);
        }
        if (value instanceof Short) {
            return PlcINT.of(value);
        }
        if (value instanceof Integer) {
            return PlcDINT.of(value);
        }
        if (value instanceof Long) {
            return PlcLINT.of(value);
        }
        if (value instanceof BigInteger) {
            return new PlcBigInteger((BigInteger) value);
        }
        if (value instanceof Float) {
            return PlcREAL.of(value);
        }
        if (value instanceof Double) {
            return PlcLREAL.of(value);
        }
        if (value instanceof BigDecimal) {
            return new PlcBigDecimal((BigDecimal) value);
        }
        if (value instanceof Duration) {
            return new PlcTIME((Duration) value);
        }
        if (value instanceof LocalTime) {
            return new PlcTIME_OF_DAY((LocalTime) value);
        }
        if (value instanceof LocalDate) {
            return new PlcDATE((LocalDate) value);
        }
        if (value instanceof LocalDateTime) {
            return new PlcDATE_AND_TIME((LocalDateTime) value);
        }
        if (value instanceof String) {
            return new PlcSTRING((String) value);
        }
        if (value instanceof PlcValue) {
            return (PlcValue) value;
        }
        throw new PlcUnsupportedDataTypeException("Data Type " + value.getClass()
            + " Is not supported");
    }


    public static PlcValue of(PlcField field, Object value) {
        return of(field, new Object[]{value});
    }


    public static PlcValue of(PlcField field, Object[] values) {
        if (values.length == 1) {
            Object value = values[0];
            switch (field.getPlcDataType().toUpperCase()) {
                case "BOOL":
                case "BIT":
                    return PlcBOOL.of(value);
                case "BYTE":
                case "BITARR8":
                    if(value instanceof Short) {
                        return new PlcBitString((short) value);
                    } else if(value instanceof Integer) {
                        return new PlcBitString(((Integer) value).shortValue());
                    } else if(value instanceof Long) {
                        return new PlcBitString(((Long) value).shortValue());
                    } else if(value instanceof BigInteger) {
                        return new PlcBitString(((BigInteger) value).shortValue());
                    } else if(value instanceof boolean[]) {
                        return new PlcBitString((boolean[]) value);
                    }
                    throw new PlcRuntimeException("BYTE requires short or boolean[8]");
                case "SINT":
                case "INT8":
                    return PlcSINT.of(value);
                case "USINT":
                case "UINT8":
                case "BIT8":
                    return PlcUSINT.of(value);
                case "INT":
                case "INT16":
                    return PlcINT.of(value);
                case "UINT":
                case "UINT16":
                    return PlcUINT.of(value);
                case "WORD":
                case "BITARR16":
                    if(value instanceof Short) {
                        return new PlcBitString((int) value);
                    } else if(value instanceof Integer) {
                        return new PlcBitString((int) value);
                    } else if(value instanceof Long) {
                        return new PlcBitString(((Long) value).intValue());
                    } else if(value instanceof BigInteger) {
                        return new PlcBitString(((BigInteger) value).intValue());
                    } else if(value instanceof boolean[]) {
                        return new PlcBitString((boolean[]) value);
                    }
                    throw new PlcRuntimeException("WORD requires int or boolean[16]");
                case "DINT":
                case "INT32":
                    return PlcDINT.of(value);
                case "UDINT":
                case "UINT32":
                    return PlcUDINT.of(value);
                case "DWORD":
                case "BITARR32":
                    if(value instanceof Short) {
                        return new PlcBitString((long) value);
                    } else if(value instanceof Integer) {
                        return new PlcBitString((long) value);
                    } else if(value instanceof Long) {
                        return new PlcBitString((long) value);
                    } else if(value instanceof BigInteger) {
                        return new PlcBitString(((BigInteger) value).longValue());
                    } else if(value instanceof boolean[]) {
                        return new PlcBitString((boolean[]) value);
                    }
                    throw new PlcRuntimeException("DWORD requires long or boolean[32]");
                case "LINT":
                case "INT64":
                    return PlcLINT.of(value);
                case "ULINT":
                case "UINT64":
                    return PlcULINT.of(value);
                case "LWORD":
                case "BITARR64":
                    if(value instanceof Short) {
                        return new PlcBitString(BigInteger.valueOf((long) value));
                    } else if(value instanceof Integer) {
                        return new PlcBitString(BigInteger.valueOf((long) value));
                    } else if(value instanceof Long) {
                        return new PlcBitString(BigInteger.valueOf((long) value));
                    } else if(value instanceof BigInteger) {
                        return new PlcBitString((BigInteger) value);
                    } else if(value instanceof boolean[]) {
                        return new PlcBitString((boolean[]) value);
                    }
                    throw new PlcRuntimeException("LWORD requires BigInteger or boolean[64]");
                case "REAL":
                case "FLOAT":
                    return PlcREAL.of(value);
                case "LREAL":
                case "DOUBLE":
                    return PlcLREAL.of(value);
                case "CHAR":
                    return PlcCHAR.of(value);
                case "WCHAR":
                    return PlcWCHAR.of(value);
                case "STRING":
                    return PlcSTRING.of(value);
                case "WSTRING":
                case "STRING16":
                    return PlcSTRING.of(value);
                case "TIME":
                    return PlcTIME.of(value);
                case "DATE":
                    return PlcDATE.of(value);
                case "TIME_OF_DAY":
                    return PlcTIME_OF_DAY.of(value);
                case "DATE_AND_TIME":
                    return PlcDATE_AND_TIME.of(value);
                default:
                    return customDataType(new Object[]{value});
            }
        } else {
            PlcList list = new PlcList();
            for (Object value : values) {
                list.add(of(field, new Object[]{value}));
            }
            return list;
        }
    }

    public static PlcValue customDataType(Object[] values) {
        return of(values);
    }
}
