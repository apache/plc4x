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
package org.apache.plc4x.java.modbus.base.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.modbus.readwrite.*;
import org.apache.plc4x.java.spi.codegen.WithOption;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.regex.Pattern;

public abstract class ModbusTag implements PlcTag, Serializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>\\d+)(:(?<datatype>[a-zA-Z_]+))?(\\[(?<quantity>\\d+)])?");
    public static final Pattern FIXED_DIGIT_MODBUS_PATTERN = Pattern.compile("(?<address>\\d{4,5})?(:(?<datatype>[a-zA-Z_]+))?(\\[(?<quantity>\\d+)])?");

    public static final int PROTOCOL_ADDRESS_OFFSET = 1;

    private final int address;

    private final int quantity;

    private final ModbusDataType dataType;

    public static ModbusTag of(String addressString) {
        if (ModbusTagCoil.matches(addressString)) {
            return ModbusTagCoil.of(addressString);
        }
        if (ModbusTagDiscreteInput.matches(addressString)) {
            return ModbusTagDiscreteInput.of(addressString);
        }
        if (ModbusTagHoldingRegister.matches(addressString)) {
            return ModbusTagHoldingRegister.of(addressString);
        }
        if (ModbusTagInputRegister.matches(addressString)) {
            return ModbusTagInputRegister.of(addressString);
        }
        if (ModbusTagExtendedRegister.matches(addressString)) {
            return ModbusTagExtendedRegister.of(addressString);
        }
        throw new PlcInvalidTagException("Unable to parse address: " + addressString);
    }

    @Override
    public String getAddressString() {
        String address = String.format("%s%05d", getAddressStringPrefix(), getLogicalAddress());
        if(getDataType() != null) {
            address += ":" + getDataType().name();
        }
        if(!getArrayInfo().isEmpty()) {
            address += "[" + getArrayInfo().get(0).getUpperBound() + "]";
        }
        return address;
    }

    protected abstract String getAddressStringPrefix();

    /**
     * Instantiate a new ModbusTag
     * @param address The WIRE address that is to be used.
     * @param quantity The number of registers
     * @param dataType The type for the interpretation of the registers.
     */
    protected ModbusTag(int address, Integer quantity, ModbusDataType dataType) {
        this.address = address;
        if (getLogicalAddress() <= 0) {
            throw new IllegalArgumentException("address must be greater than zero. Was " + getLogicalAddress());
        }
        this.quantity = quantity != null ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater than zero. Was " + this.quantity);
        }
        this.dataType = dataType != null ? dataType : ModbusDataType.INT;
    }

    /**
     * Get the technical address that must be used 'on the wire'
     * @return The address that is to be used on the wire (shifted by 1 because of the modbus spec).
     */
    public int getAddress() {
        return address;
    }

    /**
     * Get the logical (configured) address
     * @return The address which was configured and is different from what is used on the wire.
     */
    public abstract int getLogicalAddress();

    public int getNumberOfElements() {
        return quantity;
    }

    public int getLengthBytes() {
        return quantity * dataType.getDataTypeSize();
    }

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
        if (!(o instanceof ModbusTag)) {
            return false;
        }
        ModbusTag that = (ModbusTag) o;
        return address == that.address &&
            quantity == that.quantity &&
            dataType == that.dataType &&
            getClass() == that.getClass(); // MUST be identical
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.getClass(), address, quantity, dataType);
    }

    @Override
    public String toString() {
        return this.getClass().getSimpleName() + " {" +
            "address=" + address +
            ", quantity=" + quantity +
            ", dataType=" + dataType +
            " }";
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws SerializationException {
        writeBuffer.pushContext(getClass().getSimpleName());

        writeBuffer.writeUnsignedInt("address", 16, address);
        writeBuffer.writeUnsignedInt("numberOfElements", 16, getNumberOfElements());
        String dataType = getDataType().name();
        writeBuffer.writeString("dataType",
            dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithEncoding(StandardCharsets.UTF_8.name()));

        writeBuffer.popContext(getClass().getSimpleName());
    }

}
