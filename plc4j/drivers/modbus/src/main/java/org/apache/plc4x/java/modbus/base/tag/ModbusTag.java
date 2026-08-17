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
import org.apache.plc4x.java.modbus.types.ModbusByteOrder;
import org.apache.plc4x.java.spi.buffers.api.Serializable;
import org.apache.plc4x.java.spi.buffers.api.WithOption;
import org.apache.plc4x.java.spi.buffers.api.WriteBuffer;
import org.apache.plc4x.java.spi.buffers.api.exceptions.BufferException;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;

import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.regex.Pattern;

public abstract class ModbusTag implements PlcTag, Serializable {

    // STRING and WSTRING carry the length of one string in parentheses, the same way the S7 driver
    // spells it: "holding-register:1:STRING(20)[3]" is three 20-character strings. The quantity in
    // brackets keeps meaning "how many values", as it does for every other data type.
    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>\\d+)(:(?<datatype>[a-zA-Z_]+)(\\((?<stringLength>\\d+)\\))?)?(\\[(?<quantity>\\d+)])?");
    public static final Pattern FIXED_DIGIT_MODBUS_PATTERN = Pattern.compile("(?<address>\\d{4,5})?(:(?<datatype>[a-zA-Z_]+)(\\((?<stringLength>\\d+)\\))?)?(\\[(?<quantity>\\d+)])?");

    public static final int PROTOCOL_ADDRESS_OFFSET = 1;

    private final int address;

    private final int quantity;

    /** The declared length of a single string; 1 for every other data type. */
    private final int stringLength;

    private final ModbusDataType dataType;
    private final Short unitId;
    private final ModbusByteOrder byteOrder;

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
            address += "[" + (getArrayInfo().get(0).getUpperBound() + 1) + "]";
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
        this(address, quantity, dataType, new HashMap<>());
    }

    protected ModbusTag(int address, Integer quantity, ModbusDataType dataType, Map<String, String> config) {
        this(address, quantity, null, dataType, config);
    }

    protected ModbusTag(int address, Integer quantity, Integer stringLength, ModbusDataType dataType,
                        Map<String, String> config) {
        this.address = address;
        if (getLogicalAddress() <= 0) {
            throw new IllegalArgumentException("address must be greater than zero. Was " + getLogicalAddress());
        }
        this.quantity = quantity != null ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater than zero. Was " + this.quantity);
        }
        this.dataType = dataType != null ? dataType : ModbusDataType.INT;
        // A string's length is part of its address, because nothing on the wire announces it. For
        // everything else the notion doesn't apply, and a length of 1 keeps the size arithmetic
        // below unchanged.
        if ((this.dataType == ModbusDataType.STRING) || (this.dataType == ModbusDataType.WSTRING)) {
            if (stringLength == null) {
                throw new PlcInvalidTagException(this.dataType.name() + " requires the length of one string, "
                    + "for example '" + this.dataType.name() + "(20)'");
            }
            if (stringLength <= 0) {
                throw new IllegalArgumentException("string length must be greater than zero. Was " + stringLength);
            }
            this.stringLength = stringLength;
        } else {
            if (stringLength != null) {
                throw new PlcInvalidTagException("A length in parentheses is only supported for STRING and "
                    + "WSTRING, not for " + this.dataType.name());
            }
            this.stringLength = 1;
        }
        this.unitId = Optional.ofNullable(config.get("unit-id"))
            .map(Short::parseShort)
            .orElse(null);
        if(config.containsKey("byte-order")) {
            this.byteOrder = ModbusByteOrder.valueOf(config.get("byte-order"));
        } else {
            this.byteOrder = null;
        }
    }

    /**
     * Get the technical address that must be used 'on the wire'
     * @return The address that is to be used on the wire (shifted by 1 because of the modbus spec).
     */
    public int getAddress() {
        return address;
    }

    public Short getUnitId() {
        return unitId;
    }

    public ModbusByteOrder getByteOrder() {
        return byteOrder;
    }

    /**
     * Get the logical (configured) address
     * @return The address which was configured and is different from what is used on the wire.
     */
    public abstract int getLogicalAddress();

    public int getNumberOfElements() {
        return quantity;
    }

    /** The declared length of a single string, or 1 for a data type that isn't a string. */
    public int getStringLength() {
        return stringLength;
    }

    public int getLengthBytes() {
        return quantity * stringLength * dataType.getDataTypeSize();
    }

    public int getLengthWords() {
        return (int) ((quantity * stringLength * (float) dataType.getDataTypeSize()) / 2.0f);
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
            return Collections.singletonList(new DefaultArrayInfo(0, quantity - 1));
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
            unitId == that.unitId &&
            getClass() == that.getClass(); // MUST be identical
    }

    @Override
    public int hashCode() {
        return Objects.hash(this.getClass(), address, quantity, dataType, unitId);
    }

    @Override
    public String toString() {
        return this.getClass().getSimpleName() + " {" +
            "address=" + address +
            ", quantity=" + quantity +
            ", dataType=" + dataType +
            ", unitId=" + unitId +
            " }";
    }

    @Override
    public void serialize(WriteBuffer writeBuffer) throws BufferException {
        writeBuffer.pushContext(WithOption.WithName(getClass().getSimpleName()));

        writeBuffer.writeUnsignedInt(16, address, WithOption.WithName("address"));
        writeBuffer.writeUnsignedInt(16, getNumberOfElements(), WithOption.WithName("numberOfElements"));
        String dataType = getDataType().name();
        writeBuffer.writeString(
            dataType.getBytes(StandardCharsets.UTF_8).length * 8,
            dataType, WithOption.WithName("dataType"), WithOption.WithEncoding("UTF8"));

        if (unitId != null) {
            writeBuffer.writeUnsignedInt(8, unitId, WithOption.WithName("unitId"));
        }
        writeBuffer.popContext(WithOption.WithName(getClass().getSimpleName()));
    }

}
