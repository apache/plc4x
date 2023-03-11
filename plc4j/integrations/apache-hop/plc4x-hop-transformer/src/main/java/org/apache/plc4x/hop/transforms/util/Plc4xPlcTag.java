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
package org.apache.plc4x.hop.transforms.util;

import com.fasterxml.jackson.annotation.JsonIgnore;
import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.PlcTag;
import org.apache.plc4x.java.spi.generation.SerializationException;
import org.apache.plc4x.java.spi.generation.WriteBuffer;
import org.apache.plc4x.java.spi.utils.Serializable;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/*
* Code taken from the Modbus driver implementation.
* This allows me to make the correct interpretation of the "Hop" values 
* towards the writing fields.
*
* TODO: The HOP interpreter generates an error when processing 
*       the "Integer" fields.
*/

public  class Plc4xPlcTag implements PlcTag, Serializable {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>[\\%a-zA-Z_\\.0-9]+)(:(?<datatype>[a-zA-Z_]+))?(\\[(?<quantity>\\d+)])?");

    protected static final int PROTOCOL_ADDRESS_OFFSET = 1;

    private final String address;

    private final int quantity;

    private final Plc4xDataType dataType;

    public static Plc4xPlcTag of(String addressString) {
        if (Plc4xPlcTag.matches(addressString)) {
            Matcher matcher = ADDRESS_PATTERN.matcher(addressString); 

            matcher.matches();
            
            String address = matcher.group("address");

            String quantityString = matcher.group("quantity");
            int quantity = quantityString != null ? Integer.parseInt(quantityString) : 1;

            Plc4xDataType dataType = (matcher.group("datatype") != null) ? Plc4xDataType.valueOf(matcher.group("datatype")) : Plc4xDataType.BOOL;

            
            return new Plc4xPlcTag(address, quantity, dataType);
        }

        throw new PlcInvalidTagException("Unable to parse address: " + addressString);
    }

    protected Plc4xPlcTag(String address, Integer quantity, Plc4xDataType dataType) {
        this.address = address;

        this.quantity = quantity != null ? quantity : 1;

        this.dataType = dataType != null ? dataType : Plc4xDataType.INT;
    }
    
    public static boolean matches(String addressString) {
        return ADDRESS_PATTERN.matcher(addressString).matches();
    }    

    public String getAddress() {
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

    public Plc4xDataType getDataType() {
        return dataType;
    }

    public String getPlcDataType() {
        return dataType.name();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof Plc4xPlcTag)) {
            return false;
        }
        Plc4xPlcTag that = (Plc4xPlcTag) o;
        return address == that.address;
    }

    @Override
    public int hashCode() {
        return Objects.hash(address);
    }

    @Override
    public String toString() {
        return "Plc4xPlcField{" +
            " address=" + address +
            " datatype=" + dataType +
            " quantity=" + quantity +
            " }";
    }

    @Override
    public void serialize(WriteBuffer wb) throws SerializationException {
        throw new UnsupportedOperationException("Not supported yet.");
    }

    @Override
    public String getAddressString() {
        throw new UnsupportedOperationException("Not supported yet.");
    }

}
