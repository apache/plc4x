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
package org.apache.plc4x.java.firmata.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.Objects;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public abstract class FirmataField implements PlcField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("(?<address>\\d+)(\\[(?<quantity>\\d+)])?");

    private final int address;

    private final int quantity;

    public static FirmataField of(String fieldString) {
        Matcher matcher = FirmataFieldAnalog.ADDRESS_PATTERN.matcher(fieldString);
        if (matcher.matches()) {
            return FirmataFieldAnalog.of(fieldString);
        }
        matcher = FirmataFieldDigital.ADDRESS_PATTERN.matcher(fieldString);
        if (matcher.matches()) {
            return FirmataFieldDigital.of(fieldString);
        }
        throw new PlcInvalidFieldException("Unable to parse address: " + fieldString);
    }

    protected FirmataField(int address, Integer quantity) {
        this.address = address;
        this.quantity = quantity != null ? quantity : 1;
        if (this.quantity <= 0) {
            throw new IllegalArgumentException("quantity must be greater then zero. Was " + this.quantity);
        }
    }

    public int getAddress() {
        return address;
    }

    public int getNumberOfElements() {
        return quantity;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (!(o instanceof FirmataField)) {
            return false;
        }
        FirmataField that = (FirmataField) o;
        return address == that.address;
    }

    @Override
    public int hashCode() {
        return Objects.hash(address);
    }

    @Override
    public String toString() {
        return "FirmataField{" +
            "address=" + address +
            "quantity=" + quantity +
            '}';
    }


}
