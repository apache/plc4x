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

package org.apache.plc4x.java.simulated.connection;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.base.connection.PlcFieldHandler;
import org.apache.plc4x.java.base.messages.items.*;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;

public class TestFieldHandler implements PlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
        if (TestField.matches(fieldQuery)) {
            return TestField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public BaseDefaultFieldItem encodeBoolean(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Boolean.class) {
            return new DefaultBooleanFieldItem((Boolean[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeByte(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Byte.class) {
            return new DefaultLongFieldItem(Arrays.stream(values).map(x -> new Long((Byte) x)).toArray(Long[]::new));
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeShort(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Short.class) {
            return new DefaultLongFieldItem(Arrays.stream(values).map(x -> new Long((Short) x)).toArray(Long[]::new));
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeInteger(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Integer.class) {
            return new DefaultLongFieldItem(Arrays.stream(values).map(x -> new Long((Integer) x)).toArray(Long[]::new));
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeBigInteger(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == BigInteger.class) {
            return new DefaultLongFieldItem(Arrays.stream(values).map(x -> ((BigInteger) x).longValue()).toArray(Long[]::new));
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeLong(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Long.class) {
            return new DefaultLongFieldItem((Long[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeFloat(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Float.class) {
            return new DefaultFloatFieldItem((Float[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeBigDecimal(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == BigDecimal.class) {
            return new DefaultBigDecimalFieldItem((BigDecimal[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeDouble(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == Double.class) {
            return new DefaultDoubleFieldItem((Double[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeString(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == String.class) {
            return new DefaultStringFieldItem((String[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeTime(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == LocalTime.class) {
            return new DefaultLocalTimeFieldItem((LocalTime[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeDate(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == LocalDate.class) {
            return new DefaultLocalDateFieldItem((LocalDate[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeDateTime(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == LocalDateTime.class) {
            return new DefaultLocalDateTimeFieldItem((LocalDateTime[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public BaseDefaultFieldItem encodeByteArray(PlcField field, Object[] values) {
        TestField testField = (TestField) field;
        if (testField.getDataType() == byte[].class) {
            return new DefaultByteArrayFieldItem(Arrays.stream(values).map(Byte.class::cast).toArray(Byte[]::new));
        }
        if (testField.getDataType() == Byte[].class) {
            return new DefaultByteArrayFieldItem((Byte[]) values);
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

}
