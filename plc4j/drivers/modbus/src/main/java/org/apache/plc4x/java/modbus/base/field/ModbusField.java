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

import com.fasterxml.jackson.annotation.JsonIgnore;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
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
    public PlcValueType getPlcValueType() {
        return PlcValueType.valueOf(dataType.name());
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        if(quantity != 1) {
            return Collections.singletonList(new DefaultArrayInfo(0, quantity));
        }
        return Collections.emptyList();
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
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeUnsignedInt("address", 16, address);
        writeBuffer.writeUnsignedInt("numberOfElements", 16, getNumberOfElements());
        String dataType = getDataType().name();
        writeBuffer.writeString("dataType", dataType.getBytes(StandardCharsets.UTF_8).length * 8, StandardCharsets.UTF_8.name(), dataType);

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
