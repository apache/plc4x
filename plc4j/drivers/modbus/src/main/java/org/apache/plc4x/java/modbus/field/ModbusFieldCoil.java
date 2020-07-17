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

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class ModbusFieldCoil extends ModbusField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("coil:" + ModbusField.ADDRESS_PATTERN);
    public static final Pattern ADDRESS_SHORTER_PATTERN = Pattern.compile("0" + ModbusField.ADDRESS_PATTERN);
    public static final Pattern ADDRESS_SHORT_PATTERN = Pattern.compile("0x" + ModbusField.ADDRESS_PATTERN);

    public ModbusFieldCoil(int address, Integer quantity) {
        super(address, quantity);
    }

    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches() ||
            ADDRESS_SHORTER_PATTERN.matcher(addressString).matches() ||
            ADDRESS_SHORT_PATTERN.matcher(addressString).matches();
    }

    public static Matcher getMatcher(String addressString) throws PlcInvalidFieldException {
        Matcher matcher;
        if (ADDRESS_PATTERN.matcher(addressString).matches()) {
          matcher = ADDRESS_PATTERN.matcher(addressString);
        } else if (ADDRESS_SHORT_PATTERN.matcher(addressString).matches()) {
          matcher = ADDRESS_SHORT_PATTERN.matcher(addressString);
        } else if (ADDRESS_SHORTER_PATTERN.matcher(addressString).matches()) {
          matcher = ADDRESS_SHORTER_PATTERN.matcher(addressString);
        } else {
          throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
        }
        return matcher;
    }

    public static ModbusFieldCoil of(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = getMatcher(addressString);
        matcher.find();
        int address = Integer.parseInt(matcher.group("address")) - protocolAddressOffset;

        String quantityString = matcher.group("quantity");
        Integer quantity = quantityString != null ? Integer.valueOf(quantityString) : null;
        return new ModbusFieldCoil(address, quantity);
    }

}
