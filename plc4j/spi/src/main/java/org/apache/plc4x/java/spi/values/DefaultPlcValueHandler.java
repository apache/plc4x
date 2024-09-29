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

import org.apache.commons.lang3.NotImplementedException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.exceptions.PlcUnsupportedDataTypeException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;

import java.math.BigInteger;
import java.time.Duration;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;

public class DefaultPlcValueHandler implements PlcValueHandler {

    @Override
    public PlcValue newPlcValue(PlcTag tag, Object value) {
        return of(tag, new Object[]{value});
    }

    @Override
    public PlcValue newPlcValue(PlcTag tag, Object[] values) {
        return of(tag, values);
    }

    public static PlcValue of(PlcTag tag, Object value) {
        return of(tag, new Object[]{value});
    }

    public static PlcValue of(PlcTag tag, Object[] values) {
        if(tag.getArrayInfo().isEmpty()) {
            // If this is not an array type, but the input is passed in as List, if it's just one element,
            // treat this single element as "value".
            if (values.length == 0) {
                throw new PlcRuntimeException("No value was passed in as argument");
            }
            if (values.length == 1) {
                // If only one element was passed hin, however this is a Collection and contains more than one element,
                // this is also invalid.
                if(values[0] instanceof Collection) {
                    if(((Collection<?>) values[0]).size() > 1) {
                        throw new PlcRuntimeException("The type is not an array type, but a collection of more than one item was passed in as argument");
                    } else if(((Collection<?>) values[0]).isEmpty()) {
                        throw new PlcRuntimeException("An empty collection was passed in as argument");
                    }
                }
                return ofElement(tag.getPlcValueType(), values[0]);
            }
            //
            else {
                throw new PlcRuntimeException("The type is not an array type, but more than one value was passed in as argument");
            }
        }
        // In all other cases, we're dealing with an array type and this needs to be handled separately.
        else {
            return ofElements(tag.getPlcValueType(), tag.getArrayInfo(), values);
        }
    }

    private static PlcValue ofElements(PlcValueType type, List<ArrayInfo> arrayInfos, Object[] values) {
        ArrayInfo arrayInfo = arrayInfos.get(0);

        if(values.length == 1) {
            if(values[0] instanceof Object[]) {
                values = (Object[]) values[0];
            } else if (values[0] instanceof Collection) {
                values = ((Collection<?>) values[0]).toArray();
            } else if(values[0] instanceof PlcList) {
                values = ((PlcList) values[0]).getList().toArray();
            } else if(values[0] instanceof PlcRawByteArray) {
                PlcRawByteArray plcRawByteArray = (PlcRawByteArray) values[0];
                if(plcRawByteArray.getRaw().length != arrayInfo.getSize()) {
                    throw new PlcRuntimeException(String.format("Expecting %d items, but got %d", arrayInfo.getSize(), plcRawByteArray.getRaw().length));
                }
                return plcRawByteArray;
            }
        }

        List<PlcValue> plcValues = new ArrayList<>(arrayInfo.getSize());
        // In the last layer we'll create a list of PlcValues
        if(arrayInfos.size() == 1) {
            if(values.length != arrayInfo.getSize()) {
                throw new PlcRuntimeException(String.format("Expecting %d items, but got %d", arrayInfo.getSize(), values.length));
            }
            // TODO: Add some size-checks here ...
            for (Object value : values) {
                plcValues.add(ofElement(type, value));
            }
        }
        // In intermediate layers we'll add a list of PlcLists
        else {
            // TODO: Add some size-checks here ...
            for (Object value : values) {
                plcValues.add(ofElements(type, arrayInfos.subList(1, arrayInfos.size()), (Object[]) value));
            }
        }
        return new PlcList(plcValues);
    }

    private static PlcValue ofElement(PlcValueType type, Object value) {
        // This is a temporary hack for drivers that don't have type information in their tags (ADS)
        if(type == null) {
            return of(value);
        }
        switch (type) {
            case BOOL:
                return PlcBOOL.of(value);
            case BYTE:
                return PlcBYTE.of(value);
            case SINT:
                return PlcSINT.of(value);
            case USINT:
                return PlcUSINT.of(value);
            case INT:
                return PlcINT.of(value);
            case UINT:
                return PlcUINT.of(value);
            case WORD:
                return PlcWORD.of(value);
            case DINT:
                return PlcDINT.of(value);
            case UDINT:
                return PlcUDINT.of(value);
            case DWORD:
                return PlcDWORD.of(value);
            case LINT:
                return PlcLINT.of(value);
            case ULINT:
                return PlcULINT.of(value);
            case LWORD:
                return PlcLWORD.of(value);
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
            case LTIME:
                return PlcLTIME.of(value);
            case DATE:
                return PlcDATE.of(value);
            case LDATE:
                return PlcLDATE.of(value);
            case TIME_OF_DAY:
                return PlcTIME_OF_DAY.of(value);
            case LTIME_OF_DAY:
                return PlcLTIME_OF_DAY.of(value);
            case DATE_AND_TIME:
                return PlcDATE_AND_TIME.of(value);
            case DATE_AND_LTIME:
                return PlcDATE_AND_LTIME.of(value);
            case LDATE_AND_TIME:
                return PlcLDATE_AND_TIME.of(value);
            case RAW_BYTE_ARRAY:
                return PlcRawByteArray.of(value);
            case List:
                // TODO: A tag type LIST actually doesn't make any sense ...
                //  if it's an array, the array information is provided anyway. So I think we should most
                //  probably remove this type from the PlcValueType enumeration.
                throw new NotImplementedException("Not implemented yet");
            case Struct:
                // TODO: It is pretty much impossible to interpret a java object as struct.
                //  We probably shouldn't even try to do so. It might be an interesting option, if we defined
                //  annotations that could be added to Java types what allow the conversion to PlcStruct types
                //  but right now here is nothing we can do.
                throw new NotImplementedException("Not implemented yet");
            case NULL:
                return new PlcNull();
            default:
                throw new PlcUnsupportedDataTypeException("Data Type " + value.getClass()
                    + " Is not supported");
        }
    }

    /**
     * This is a legacy helper that should help with protocols, that currently don't have PlcValueType information
     * available at the request-building-time. Such as the ADS driver. We should remove this option as soon as all
     * drivers are fully PlcValueType-aware.
     * @param value value
     * @return PlcValue for the given type
     */
    private static PlcValue of(Object value) {
        /*if (values.length != 1) {
        PlcList list = new PlcList();
            for (Object value : values) {
                list.add(of(new Object[]{value}));
            }
            return list;
        }
        Object value = values[0];*/
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

/*    public static PlcValue customDataType(Object[] values) {
        return of(values);
    }*/
}
