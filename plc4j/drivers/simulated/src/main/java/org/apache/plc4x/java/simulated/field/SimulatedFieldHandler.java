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
import org.apache.plc4x.java.spi.connection.DefaultPlcFieldHandler;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.util.Arrays;
import java.util.LinkedList;
import java.util.List;

public class SimulatedFieldHandler extends DefaultPlcFieldHandler {

    @Override
    public PlcField createField(String fieldQuery) {
        if (SimulatedField.matches(fieldQuery)) {
            return SimulatedField.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

    @Override
    public PlcValue encodeTime(PlcField field, Object[] values) {
        switch (field.getPlcDataType().toUpperCase()) {
            case "LOCALTIME":
            case "TIME":
                if(values.length == 1) {
                    return new PlcTime((LocalTime) values[0]);
                } else {
                    List<PlcTime> plcValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcValues.add(new PlcTime((LocalTime) values[i]));
                    }
                    return new PlcList(plcValues);
                }
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + field.getPlcDataType());
        }
    }

    @Override
    public PlcValue encodeDate(PlcField field, Object[] values) {
        switch (field.getPlcDataType().toUpperCase()) {
            case "LOCALDATE":
            case "DATE":
                if(values.length == 1) {
                    return new PlcDate((LocalDate) values[0]);
                } else {
                    List<PlcDate> plcValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcValues.add(new PlcDate((LocalDate) values[i]));
                    }
                    return new PlcList(plcValues);
                }
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + field.getPlcDataType());
        }
    }

    @Override
    public PlcValue encodeDateTime(PlcField field, Object[] values) {
        switch (field.getPlcDataType().toUpperCase()) {
            case "LOCALDATETIME":
            case "DATETIME":
                if(values.length == 1) {
                    return new PlcDateTime((LocalDateTime) values[0]);
                } else {
                    List<PlcDateTime> plcValues = new LinkedList<>();
                    for (int i = 0; i < values.length; i++) {
                        plcValues.add(new PlcDateTime((LocalDateTime) values[i]));
                    }
                    return new PlcList(plcValues);
                }
            default:
                throw new PlcRuntimeException("Invalid encoder for type " + field.getPlcDataType());
        }
    }

}
