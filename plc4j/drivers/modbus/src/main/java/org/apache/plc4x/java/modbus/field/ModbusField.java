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

import com.fasterxml.jackson.annotation.JsonIgnore;
import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.modbus.readwrite.types.*;
import org.apache.plc4x.java.spi.utils.XmlSerializable;
import org.w3c.dom.Document;
import org.w3c.dom.Element;

import java.util.Objects;
import java.util.regex.Pattern;

public abstract class ModbusField implements PlcField, XmlSerializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>\\d+)(:(?<datatype>[a-zA-Z_]+))?(\\((?<stringLength>\\d{1,3})\\))?(\\[(?<quantity>\\d+)])?");
    public static final Pattern FIXED_DIGIT_MODBUS_PATTERN = Pattern.compile("(?<address>\\d{4,5})?(:(?<datatype>[a-zA-Z_]+))?(\\((?<stringLength>\\d{1,3})\\))?(\\[(?<quantity>\\d+)])?");

    protected static final int PROTOCOL_ADDRESS_OFFSET = 1;

    private final int address;

    private final int quantity;

    private final int stringLength;

    private final ModbusDataType dataType;

    public static ModbusField of(String addressString) {
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
        if(ModbusExtendedRegister.matches(addressString)) {
            return ModbusExtendedRegister.of(addressString);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + addressString);
    }

    protected ModbusField(int address, Integer quantity, Integer stringLength, ModbusDataType dataType) {
        this.address = address;
        if ((this.address + PROTOCOL_ADDRESS_OFFSET) <= 0) {
            throw new IllegalArgumentException("address must be greater than zero. Was " + (this.address + PROTOCOL_ADDRESS_OFFSET));
        }

        this.stringLength = stringLength != null ? stringLength : 1;
        if (this.stringLength <= 0) {
            throw new IllegalArgumentException("stringLength must be greater than zero. Was " + this.stringLength);
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

    public int getStringLength() {
        return stringLength;
    }

    public int getNumberOfElements() {
        return quantity;
    }

    public int getLengthBytes() {
        return quantity * getStringLength() * dataType.getDataTypeSize();
    }

    @JsonIgnore
    public int getLengthWords() {
        return Math.round(getLengthBytes() / 2.0f);       
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
    public void xmlSerialize(Element parent) {
        Document doc = parent.getOwnerDocument();
        Element messageElement = doc.createElement(getClass().getSimpleName());
        parent.appendChild(messageElement);
        Element addressElement = doc.createElement("address");
        addressElement.appendChild(doc.createTextNode(Integer.toString(getAddress())));
        messageElement.appendChild(addressElement);

        Element stringLengthElement = doc.createElement("stringLength");
        stringLengthElement.appendChild(doc.createTextNode(Integer.toString(getStringLength())));
        messageElement.appendChild(stringLengthElement);

        Element quantityElement = doc.createElement("numberOfElements");
        quantityElement.appendChild(doc.createTextNode(Integer.toString(getNumberOfElements())));
        messageElement.appendChild(quantityElement);

        Element datatypeElement = doc.createElement("dataType");
        datatypeElement.appendChild(doc.createTextNode(getPlcDataType()));
        messageElement.appendChild(datatypeElement);
    }

}
