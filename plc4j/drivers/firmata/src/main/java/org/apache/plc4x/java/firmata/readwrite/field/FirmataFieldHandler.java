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
package org.apache.plc4x.java.firmata.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.PlcBoolean;
import org.apache.plc4x.java.api.value.PlcList;
import org.apache.plc4x.java.api.value.PlcShort;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.util.ArrayList;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;

public class FirmataFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (FirmataFieldDigital.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return FirmataFieldDigital.of(fieldQuery);
        } else if (FirmataFieldAnalog.ADDRESS_PATTERN.matcher(fieldQuery).matches()) {
            return FirmataFieldAnalog.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        FirmataField firmataField = (FirmataField) field;
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
                        " is not assignable to " + firmataField + " fields.");
            }
        }
        if(booleanValues.size() == 1) {
            return new PlcBoolean(booleanValues.get(0));
        } else {
            return new PlcList(booleanValues);
        }
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        if(values.length == 1) {
            Number numberValue = (Number) values[0];
            // Intentionally checking the next larger type.
            if((numberValue.intValue() < Short.MIN_VALUE) || (numberValue.intValue() > Short.MAX_VALUE)) {
                throw new PlcInvalidFieldException("Value of " + numberValue.toString() + " exceeds the boundaries of a short value.");
            }
            return new PlcShort(numberValue.shortValue());
        } else {
            List<PlcShort> shorts = new ArrayList<>(values.length);
            for (Object value : values) {
                Number numberValue = (Number) value;
                // Intentionally checking the next larger type.
                if((numberValue.intValue() < Short.MIN_VALUE) || (numberValue.intValue() > Short.MAX_VALUE)) {
                    throw new PlcInvalidFieldException("Value of " + numberValue.toString() + " exceeds the boundaries of a short value.");
                }
                shorts.add(new PlcShort(((Number) value).shortValue()));
            }
            return new PlcList(shorts);
        }
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeBigDecimal(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

}
