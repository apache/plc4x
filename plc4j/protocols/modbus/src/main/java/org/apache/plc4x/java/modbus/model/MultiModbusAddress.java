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

import java.util.Objects;
import java.util.regex.Pattern;

public abstract class MultiModbusAddress extends ModbusAddress {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile(ModbusAddress.ADDRESS_PATTERN + "/" + "(?<quantity>\\d+)");

    private final int quantity;

    protected MultiModbusAddress(int address, int quantity) {
        super(address);
        this.quantity = quantity;
    }

    public int getQuantity() {
        return quantity;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof MultiModbusAddress)) {
            return false;
        }
        if (!super.equals(o)) {
            return false;
        }
        MultiModbusAddress that = (MultiModbusAddress) o;
        return quantity == that.quantity;
    }

    @Override
    public int hashCode() {

        return Objects.hash(super.hashCode(), quantity);
    }

    @Override
    public String toString() {
        return "MultiModbusAddress{" +
            "quantity=" + quantity +
            "} " + super.toString();
    }
}
