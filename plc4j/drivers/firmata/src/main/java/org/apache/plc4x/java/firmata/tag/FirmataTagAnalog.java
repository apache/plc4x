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
package org.apache.plc4x.java.firmata.tag;

import org.apache.plc4x.java.api.exceptions.PlcInvalidTagException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.spi.drivers.model.ArrayNotationParser;
import org.apache.plc4x.java.spi.drivers.model.DefaultArrayInfo;

import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class FirmataTagAnalog extends FirmataTag {

    public static final Pattern ADDRESS_PATTERN = Pattern.compile("analog:" + FirmataTag.ADDRESS_PATTERN);

    public FirmataTagAnalog(int address, Integer quantity) {
        super(address, quantity);
    }

    public FirmataTagAnalog(int address, Integer quantity, boolean explicitRange) {
        super(address, quantity, explicitRange);
        
    }

    @Override
    public String getAddressString() {
        // The selection terminates the address; Firmata carries no type suffix.
        return "analog:" + getAddress() + ArrayNotationParser.render(getArrayInfo());
    }

    @Override
    public PlcValueType getPlcValueType() {
        return PlcValueType.INT;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        // A range is an array even when it spans one pin; the count cannot say which was written.
        if (isExplicitRange()) {
            return Collections.singletonList(new DefaultArrayInfo(0, getNumberOfElements() - 1, 0, true));
        }
        return Collections.emptyList();    }

    public static FirmataTagAnalog of(String addressString) throws PlcInvalidTagException {
        Matcher matcher = ADDRESS_PATTERN.matcher(addressString);
        if (!matcher.matches()) {
            throw new PlcInvalidTagException(addressString, ADDRESS_PATTERN);
        }
        int address = Integer.parseInt(matcher.group("address"));

        int[] selection = selectionOf(matcher, addressString);
        address += selection[0];
        Integer quantity = selection[1];
        return new FirmataTagAnalog(address, quantity, selection[2] == 1);
    }

}
