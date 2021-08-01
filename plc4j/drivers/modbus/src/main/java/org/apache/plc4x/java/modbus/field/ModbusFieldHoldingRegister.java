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
package org.apache.plc4x.java.modbus.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.modbus.readwrite.types.ModbusDataType;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ModbusFieldHoldingRegister extends ModbusField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("holding-register:" + ModbusField.ADDRESS_PATTERN);
    public static final Pattern ADDRESS_SHORTER_PATTERN = Pattern.compile("4" + ModbusField.FIXED_DIGIT_MODBUS_PATTERN);
    public static final Pattern ADDRESS_SHORT_PATTERN = Pattern.compile("4x" + ModbusField.FIXED_DIGIT_MODBUS_PATTERN);

    protected static final int REGISTER_MAXADDRESS = 65535;

    protected ModbusFieldHoldingRegister(int address, Integer quantity, ModbusDataType dataType) {
        super(address, quantity, dataType);
    }

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches() ||
            ADDRESS_SHORTER_PATTERN.matcher(addressString).matches() ||
            ADDRESS_SHORT_PATTERN.matcher(addressString).matches();
    }

    public static Matcher getMatcher(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (matcher.matches()) {
          return matcher;
        }
        matcher = ADDRESS_SHORT_PATTERN.matcher(addressString);
        if (matcher.matches()) {
          return matcher;
        }
        matcher = ADDRESS_SHORTER_PATTERN.matcher(addressString);
        if (matcher.matches()) {
          return matcher;
        }
        throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
    }

    public static ModbusFieldHoldingRegister of(String addressString) {
        Matcher matcher = getMatcher(addressString);
        int address = Integer.parseInt(matcher.group("address")) - PROTOCOL_ADDRESS_OFFSET;
        if (address > REGISTER_MAXADDRESS) {
            throw new IllegalArgumentException("Address must be less than or equal to " + REGISTER_MAXADDRESS + ". Was " + (address + PROTOCOL_ADDRESS_OFFSET));
        }

        String quantityString = matcher.group("quantity");
        int quantity = quantityString != null ? Integer.parseInt(quantityString) : 1;
        if ((address + quantity) > REGISTER_MAXADDRESS) {
            throw new IllegalArgumentException("Last requested address is out of range, should be between " + PROTOCOL_ADDRESS_OFFSET + " and " + REGISTER_MAXADDRESS + ". Was " + (address + PROTOCOL_ADDRESS_OFFSET + (quantity - 1)));
        }

        ModbusDataType dataType = (matcher.group("datatype") != null) ? ModbusDataType.valueOf(matcher.group("datatype")) : ModbusDataType.INT;

        return new ModbusFieldHoldingRegister(address, quantity, dataType);
    }

}
