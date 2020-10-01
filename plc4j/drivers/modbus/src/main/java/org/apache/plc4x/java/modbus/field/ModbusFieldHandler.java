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
package org.apache.plc4x.java.modbus.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.value.*;
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;
import org.apache.plc4x.java.modbus.readwrite.*;

import java.math.BigInteger;
import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;
import java.util.Arrays;
import java.util.stream.Collectors;

public class ModbusFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (ModbusFieldDiscreteInput.matches(fieldQuery)) {
            return ModbusFieldDiscreteInput.of(fieldQuery);
        } else if (ModbusFieldHoldingRegister.matches(fieldQuery)) {
            return ModbusFieldHoldingRegister.of(fieldQuery);
        } else if (ModbusFieldInputRegister.matches(fieldQuery)) {
            return ModbusFieldInputRegister.of(fieldQuery);
        } else if (ModbusFieldCoil.matches(fieldQuery)) {
            return ModbusFieldCoil.of(fieldQuery);
        } else if (ModbusExtendedRegister.matches(fieldQuery)) {
            return ModbusExtendedRegister.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeBoolean(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeLong(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeBigInteger(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeFloat(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeBigDecimal(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        return internalEncode(field, values);
    }

    private PlcValue internalEncode(PlcField field, Object[] values) {
        ModbusField modbusField = (ModbusField) field;
        try {
            switch (modbusField.getDataType()) {
                //Implement Custom PlcValue types here
                default:
                    return PlcValues.of(values, Class.forName(PlcValues.class.getPackage().getName() + ".Plc" + modbusField.getDataType()));
            }
        } catch (ClassNotFoundException e) {
            throw new PlcRuntimeException("Invalid encoder for type " + modbusField.getDataType() + e);
        }
    }
}
