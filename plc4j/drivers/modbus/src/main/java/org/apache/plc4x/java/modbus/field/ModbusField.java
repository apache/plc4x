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

import com.fasterxml.jackson.annotation.JsonIgnore;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.modbus.readwrite.types.*;
import org.apache.plc4x.java.spi.generation.ParseException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.nio.charset.StandardCharsets;
import java.util.Objects;
import java.util.regex.Pattern;

public abstract class ModbusField implements PlcField, Serializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>\\d+)(:(?<datatype>[a-zA-Z_]+))?(\\[(?<quantity>\\d+)])?");
    public static final Pattern FIXED_DIGIT_MODBUS_PATTERN = Pattern.compile("(?<address>\\d{4,5})?(:(?<datatype>[a-zA-Z_]+))?(\\[(?<quantity>\\d+)])?");

    protected static final int PROTOCOL_ADDRESS_OFFSET = 1;

    private final int address;

    private final int quantity;

    private final ModbusDataType dataType;

    public static ModbusField of(String addressString) {
        if (ModbusFieldCoil.matches(addressString)) {
            return ModbusFieldCoil.of(addressString);
        }
        if (ModbusFieldDiscreteInput.matches(addressString)) {
            return ModbusFieldDiscreteInput.of(addressString);
        }
        if (ModbusFieldHoldingRegister.matches(addressString)) {
            return ModbusFieldHoldingRegister.of(addressString);
        }
        if (ModbusFieldInputRegister.matches(addressString)) {
            return ModbusFieldInputRegister.of(addressString);
        }
        if (ModbusExtendedRegister.matches(addressString)) {
            return ModbusExtendedRegister.of(addressString);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + addressString);
    }

    protected ModbusField(int address, Integer quantity, ModbusDataType dataType) {
        this.address = address;
        if ((this.address + PROTOCOL_ADDRESS_OFFSET) <= 0) {
            throw new IllegalArgumentException("address must be greater than zero. Was " + (this.address + PROTOCOL_ADDRESS_OFFSET));
        }
        this.quantity = quantity != null ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater than zero. Was " + this.quantity);
        }
        this.dataType = dataType != null ? dataType : ModbusDataType.INT;
    }

    public int getAddress() {
        return address;
    }

    public int getNumberOfElements() {
        return quantity;
    }

    public int getLengthBytes() {
        return quantity * dataType.getDataTypeSize();
    }

    @JsonIgnore
    public int getLengthWords() {
        return (int) ((quantity * (float) dataType.getDataTypeSize()) / 2.0f);
    }

    public ModbusDataType getDataType() {
        return dataType;
    }

    @Override
    public String getPlcDataType() {
        return dataType.name();
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
            "datatype=" + dataType +
            "quantity=" + quantity +
            '}';
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws ParseException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeUnsignedInt("address", 16, address);
        writeBuffer.writeUnsignedInt("numberOfElements", 16, getNumberOfElements());
        String dataType = getPlcDataType();
        writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
