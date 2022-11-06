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
package org.apache.plc4x.java.modbus.base.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.spi.connection.PlcFieldHandler;

public class ModbusFieldHandler implements PlcFieldHandler {

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
        } else if (ModbusIdentificationRegister.matches(fieldQuery)) {
            return ModbusIdentificationRegister.of(fieldQuery);
        }
        throw new PlcInvalidFieldException(fieldQuery);
    }

}
