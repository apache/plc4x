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
import org.apache.plc4x.java.api.model.PlcField;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public abstract class ModbusField implements PlcField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>\\d+)(\\[(?<quantity>\\d+)])?");
    protected static final int protocolAddressOffset = 1;

    private final int address;

    private final int quantity;

    public static ModbusField of(String addressString) throws PlcInvalidFieldException {
        if(ModbusFieldCoil.matches(addressString)) {
            return ModbusFieldCoil.of(addressString);
        }
        if(ModbusFieldDiscreteInput.matches(addressString)) {
            return ModbusFieldDiscreteInput.of(addressString);
        }
        if(ModbusFieldHoldingRegister.matches(addressString)) {
            return ModbusFieldHoldingRegister.of(addressString);
        }
        if(ModbusFieldInputRegister.matches(addressString)) {
            return ModbusFieldInputRegister.of(addressString);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + addressString);
    }

    protected ModbusField(int address, Integer quantity) {
        this.address = address;
        if ((this.address + protocolAddressOffset) <= 0) {
            throw new IllegalArgumentException("address must be greater then zero. Was " + (this.address + protocolAddressOffset));
        }
        this.quantity = quantity != null ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater then zero. Was " + this.quantity);
        }
    }

    public int getAddress() {
        return address;
    }

    public int getQuantity() {
        return quantity;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof ModbusField)) {
            return false;
        }
        ModbusField that = (ModbusField) o;
        return address == that.address;
    }

    @Override
    public int hashCode() {
        return Objects.hash(address);
    }

    @Override
    public String toString() {
        return "ModbusField{" +
            "address=" + address +
            "quantity=" + quantity +
            '}';
    }

}
