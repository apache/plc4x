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
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.BitSet;
import java.util.LinkedList;
import java.util.List;

public class ModbusFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) throws PlcInvalidFieldException {
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
        ModbusField modbusField = (ModbusField) field;
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
                        " is not assignable to " + modbusField + " fields.");
            }
        }
        if(booleanValues.size() == 1) {
            return new PlcBOOL(booleanValues.get(0));
        } else {
            return new PlcList(booleanValues);
        }
    }

    @Override
    public PlcValue encodeByte(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeInteger(PlcField field, Object[] values) {
        ModbusField modbusField = (ModbusField) field;
        switch (modbusField.getDataType()) {
            case "INT":
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
              if(values.length == 1) {
                  return new PlcUINT((Integer) values[0]);
              } else {
                  List<PlcUINT> plcUINTValues = new LinkedList<>();
                  for (int i = 0; i < values.length; i++) {
                    plcUINTValues.add(new PlcUINT((Integer) values[i]));
                  }
                  return new PlcList(plcUINTValues);
              }
            case "REAL":
              if(values.length == 1) {
                  return new PlcREAL((Integer) values[0]);
              } else {
                  List<PlcREAL> plcREALValues = new LinkedList<>();
                  for (int i = 0; i < values.length; i++) {
                    plcREALValues.add(new PlcREAL((Integer) values[i]));
                  }
                  return new PlcList(plcREALValues);
              }
            case "BOOL":
              if(values.length == 1) {
                  return new PlcBOOL((Integer) values[0]);
              } else {
                  List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                  for (int i = 0; i < values.length; i++) {
                    plcBOOLValues.add(new PlcBOOL((Integer) values[i]));
                  }
                  return new PlcList(plcBOOLValues);
              }
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + modbusField.getDataType());
        }
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
      ModbusField modbusField = (ModbusField) field;
      switch (modbusField.getDataType()) {
          case "REAL":
            if(values.length == 1) {
                return new PlcREAL((String) values[0]);
            } else {
                List<PlcREAL> plcREALValues = new LinkedList<>();
                for (int i = 0; i < values.length; i++) {
                  plcREALValues.add(new PlcREAL((String) values[i]));
                }
                return new PlcList(plcREALValues);
            }
          default:
              throw new PlcRuntimeException("Invalid encoder for type " + modbusField.getDataType());
      }
    }

    @Override
    public PlcValue encodeDouble(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeBigDecimal(PlcField field, Object[] values) {
        return encodeShort(field, values);
    }

    @Override
    public PlcValue encodeString(PlcField field, Object[] values) {
        ModbusField modbusField = (ModbusField) field;
        switch (modbusField.getDataType()) {
            case "INT":
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
              if(values.length == 1) {
                  return new PlcUINT((String) values[0]);
              } else {
                  List<PlcUINT> plcUINTValues = new LinkedList<>();
                  for (int i = 0; i < values.length; i++) {
                    plcUINTValues.add(new PlcUINT((String) values[i]));
                  }
                  return new PlcList(plcUINTValues);
              }
            case "REAL":
              if(values.length == 1) {
                  return new PlcREAL((String) values[0]);
              } else {
                  List<PlcREAL> plcREALValues = new LinkedList<>();
                  for (int i = 0; i < values.length; i++) {
                    plcREALValues.add(new PlcREAL((String) values[i]));
                  }
                  return new PlcList(plcREALValues);
              }
            case "BOOL":
              if(values.length == 1) {
                  return new PlcBOOL((String) values[0]);
              } else {
                  List<PlcBOOL> plcBOOLValues = new LinkedList<>();
                  for (int i = 0; i < values.length; i++) {
                    plcBOOLValues.add(new PlcBOOL((String) values[i]));
                  }
                  return new PlcList(plcBOOLValues);
              }
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + modbusField.getDataType());
        }
    }

    @Override
    public PlcValue encodeShort(PlcField field, Object[] values) {
          ModbusField modbusField = (ModbusField) field;
          switch (modbusField.getDataType()) {
              case "INT":
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
                if(values.length == 1) {
                    return new PlcUINT((Short) values[0]);
                } else {
                    List<PlcUINT> plcUINTValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                      plcUINTValues.add(new PlcUINT((Short) values[i]));
                    }
                    return new PlcList(plcUINTValues);
                }
              default:
                  throw new PlcRuntimeException("Invalid encoder for type " + modbusField.getDataType());
          }
    }

}
