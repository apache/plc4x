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

package org.apache.plc4x.java.simulated.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;

public class SimulatedFieldHandler implements PlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (SimulatedField.matches(fieldQuery)) {
            return SimulatedField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Boolean.class) {
            if(values.length == 1) {
                return new PlcBoolean((Boolean) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Byte.class) {
            if(values.length == 1) {
                return new PlcInteger((Byte) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Short.class) {
            if(values.length == 1) {
                return new PlcInteger((Short) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Integer.class) {
            if(values.length == 1) {
                return new PlcInteger((Integer) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == BigInteger.class) {
            if(values.length == 1) {
                return new PlcBigInteger((BigInteger) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Long.class) {
            if(values.length == 1) {
                return new PlcLong((Long) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Float.class) {
            if(values.length == 1) {
                return new PlcFloat((Float) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeBigDecimal(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == BigDecimal.class) {
            if(values.length == 1) {
                return new PlcBigDecimal((BigDecimal) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == Double.class) {
            if(values.length == 1) {
                return new PlcDouble((Double) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == String.class) {
            if(values.length == 1) {
                return new PlcString((String) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeTime(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == LocalTime.class) {
            if(values.length == 1) {
                return new PlcTime((LocalTime) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeDate(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == LocalDate.class) {
            if(values.length == 1) {
                return new PlcDate((LocalDate) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

    @Override
    public PlcValue encodeDateTime(PlcField field, Object[] values) {
        SimulatedField testField = (SimulatedField) field;
        if (testField.getDataType() == LocalDateTime.class) {
            if(values.length == 1) {
                return new PlcDateTime((LocalDateTime) values[0]);
            } else {
                return new PlcList(Arrays.asList(values));
            }
        }
        throw new PlcRuntimeException("Invalid encoder for type " + testField.getDataType().getName());
    }

}
