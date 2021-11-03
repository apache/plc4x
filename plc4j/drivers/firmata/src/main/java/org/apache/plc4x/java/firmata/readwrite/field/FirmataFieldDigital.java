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
package org.apache.plc4x.java.firmata.readwrite.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.firmata.readwrite.PinMode;

import java.util.BitSet;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class FirmataFieldDigital extends FirmataField {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("digital:" + FirmataField.ADDRESS_PATTERN +
        "(:(?<mode>PULLUP))?");

    protected final BitSet bitSet;
    protected final PinMode pinMode;

    public FirmataFieldDigital(int address, Integer quantity, PinMode pinMode) {
        super(address, quantity);
        // Translate the address into a bit-set.
        bitSet = new BitSet();
        for(int i = getAddress(); i < getAddress() + getNumberOfElements(); i++) {
            bitSet.set(i, true);
        }
        this.pinMode = pinMode;
    }

    public BitSet getBitSet() {
        return bitSet;
    }

    public PinMode getPinMode() {
        return pinMode;
    }

    @Override
    public String getPlcDataType() {
        return "BOOL";
    }

    public static FirmataFieldDigital of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidFieldException(addressString, ADDRESS_PATTERN);
        }
        int address = Integer.parseInt(matcher.group("address"));

        String quantityString = matcher.group("quantity");
        Integer quantity = quantityString != null ? Integer.valueOf(quantityString) : null;

        PinMode pinMode = ("PULLUP".equals(matcher.group("mode"))) ? PinMode.PinModePullup : null;

        return new FirmataFieldDigital(address, quantity, pinMode);
    }

}
