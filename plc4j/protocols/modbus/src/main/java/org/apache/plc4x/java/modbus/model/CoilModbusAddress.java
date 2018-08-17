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
package org.apache.plc4x.java.modbus.model;

import org.apache.plc4x.java.api.exceptions.PlcRuntimeException;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class CoilModbusAddress extends ModbusAddress {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("coil:" + ModbusAddress.ADDRESS_PATTERN);

    protected CoilModbusAddress(int address) {
        super(address);
    }

    public static CoilModbusAddress of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcRuntimeException(addressString + " doesn't match " + ADDRESS_PATTERN);
        }
        int address = Integer.parseInt(matcher.group("address"));
        return new CoilModbusAddress(address);
    }
}
