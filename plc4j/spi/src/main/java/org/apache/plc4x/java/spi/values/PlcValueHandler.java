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
package org.apache.plc4x.java.spi.values;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.value.PlcValue;

import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.List;

public class PlcValueHandler implements org.apache.plc4x.java.api.value.PlcValueHandler {

    public PlcValue newPlcValue(Object value) {
        return of(new Object[]{value});
    }

    public PlcValue newPlcValue(Object[] values) {
        return of(values);
    }

    public PlcValue newPlcValue(PlcTag tag, Object value) {
        return of(tag, new Object[]{value});
    }

    public PlcValue newPlcValue(PlcTag tag, Object[] values) {
        return of(tag, values);
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
            return PlcRawByteArray.of(value);
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
            return PlcLINT.of(value);
        }
        if (value instanceof Float) {
            return PlcREAL.of(value);
        }
        if (value instanceof Double) {
            return PlcLREAL.of(value);
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


    public static PlcValue of(PlcTag tag, Object value) {
        return of(tag, new Object[]{value});
    }


    public static PlcValue of(PlcTag tag, Object[] values) {
        if (values.length == 1) {
            Object value = values[0];
            if(tag.getPlcValueType() == null) {
                // TODO: This is a hacky shortcut ..
                if(value instanceof PlcValue) {
                    return (PlcValue) value;
                }
                return new PlcNull();
            }
            switch (tag.getPlcValueType()) {
                case BOOL:
                    return PlcBOOL.of(value);
                case BYTE:
                    if(value instanceof Short) {
                        return new PlcBYTE((short) value);
                    } else if(value instanceof Integer) {
                        return new PlcBYTE(((Integer) value).shortValue());
                    } else if(value instanceof Long) {
                        return new PlcBYTE(((Long) value).shortValue());
                    } else if(value instanceof BigInteger) {
                        return new PlcBYTE(((BigInteger) value).shortValue());
                    }
                    throw new PlcRuntimeException("BYTE requires short");
                case SINT:
                    return PlcSINT.of(value);
                case USINT:
                    return PlcUSINT.of(value);
                case INT:
                    return PlcINT.of(value);
                case UINT:
                    return PlcUINT.of(value);
                case WORD:
                    if(value instanceof Short) {
                        return new PlcWORD((int) value);
                    } else if(value instanceof Integer) {
                        return new PlcWORD((int) value);
                    } else if(value instanceof Long) {
                        return new PlcWORD(((Long) value).intValue());
                    } else if(value instanceof BigInteger) {
                        return new PlcWORD(((BigInteger) value).intValue());
                    }
                    throw new PlcRuntimeException("WORD requires int");
                case DINT:
                    return PlcDINT.of(value);
                case UDINT:
                    return PlcUDINT.of(value);
                case DWORD:
                    if(value instanceof Short) {
                        return new PlcDWORD((long) value);
                    } else if(value instanceof Integer) {
                        return new PlcDWORD((long) value);
                    } else if(value instanceof Long) {
                        return new PlcDWORD((long) value);
                    } else if(value instanceof BigInteger) {
                        return new PlcDWORD(((BigInteger) value).longValue());
                    }
                    throw new PlcRuntimeException("DWORD requires long");
                case LINT:
                    return PlcLINT.of(value);
                case ULINT:
                    return PlcULINT.of(value);
                case LWORD:
                    if(value instanceof Short) {
                        return new PlcLWORD(BigInteger.valueOf((long) value));
                    } else if(value instanceof Integer) {
                        return new PlcLWORD(BigInteger.valueOf((long) value));
                    } else if(value instanceof Long) {
                        return new PlcLWORD(BigInteger.valueOf((long) value));
                    } else if(value instanceof BigInteger) {
                        return new PlcLWORD((BigInteger) value);
                    }
                    throw new PlcRuntimeException("LWORD requires BigInteger");
                case REAL:
                    return PlcREAL.of(value);
                case LREAL:
                    return PlcLREAL.of(value);
                case CHAR:
                    return PlcCHAR.of(value);
                case WCHAR:
                    return PlcWCHAR.of(value);
                case STRING:
                    return PlcSTRING.of(value);
                case WSTRING:
                    return PlcWSTRING.of(value);
                case TIME:
                    return PlcTIME.of(value);
                case DATE:
                    return PlcDATE.of(value);
                case TIME_OF_DAY:
                    return PlcTIME_OF_DAY.of(value);
                case DATE_AND_TIME:
                    return PlcDATE_AND_TIME.of(value);
                default:
                    return customDataType(new Object[]{value});
            }
        } else {
            PlcList list = new PlcList();
            for (Object value : values) {
                list.add(of(tag, new Object[]{value}));
            }
            return list;
        }
    }

    public static PlcValue customDataType(Object[] values) {
        return of(values);
    }
}
