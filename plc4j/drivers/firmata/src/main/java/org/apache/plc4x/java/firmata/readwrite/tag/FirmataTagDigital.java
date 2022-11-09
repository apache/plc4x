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
package org.apache.plc4x.java.firmata.readwrite.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.firmata.readwrite.PinMode;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;

import java.util.BitSet;
import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class FirmataTagDigital extends FirmataTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("digital:" + FirmataTag.ADDRESS_PATTERN +
        "(:(?<mode>PULLUP))?");

    protected final BitSet bitSet;
    protected final PinMode pinMode;

    public FirmataTagDigital(int address, Integer quantity, PinMode pinMode) {
        super(address, quantity);
        // Translate the address into a bit-set.
        bitSet = new BitSet();
        for(int i = getAddress(); i < getAddress() + getNumberOfElements(); i++) {
            bitSet.set(i, true);
        }
        this.pinMode = pinMode;
    }

    @Override
    public String getAddressString() {
        String address = "digital:" + getAddress();
        if(getNumberOfElements() != 1) {
            address += "[" + getNumberOfElements() + "]";
        }
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.BOOL;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        if(getNumberOfElements() != 1) {
            return Collections.singletonList(new DefaultArrayInfo(0, getNumberOfElements()));
        }
        return Collections.emptyList();
    }

    public BitSet getBitSet() {
        return bitSet;
    }

    public PinMode getPinMode() {
        return pinMode;
    }

    public static FirmataTagDigital of(String addressString) {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(addressString, ADDRESS_PATTERN);
        }
        int address = Integer.parseInt(matcher.group("address"));

        String quantityString = matcher.group("quantity");
        Integer quantity = quantityString != null ? Integer.valueOf(quantityString) : null;

        PinMode pinMode = ("PULLUP".equals(matcher.group("mode"))) ? PinMode.PinModePullup : null;

        return new FirmataTagDigital(address, quantity, pinMode);
    }

}
