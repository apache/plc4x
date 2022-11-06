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
package org.apache.plc4x.java.mock.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.ArrayInfo;
import org.apache.plc4x.java.api.model.PlcField;
import org.apache.plc4x.java.api.types.PlcValueType;
import org.apache.plc4x.java.api.value.PlcValue;
import org.apache.plc4x.java.spi.model.DefaultArrayInfo;

import java.util.Collections;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class MockField implements PlcField {

    private final String address;
    private final PlcValueType type;
    private static final Pattern PATTERN =
        Pattern.compile("%(?<address>[a-zA-Z_.0-9]+(?:\\[[0-9]*])?):?(?<type>[A-Z]*)");

    public static MockField of(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = PATTERN.matcher(addressString);
        if (matcher.matches()) {
            String address = matcher.group("address");
            PlcValueType type = PlcValueType.valueOf(matcher.group("type"));
            return new MockField(address, type);
        }
        return null;
    }

    public MockField(String address) {
        this.address = address;
        this.type = null;
    }

    public MockField(String address, PlcValueType type) {
        this.address = address;
        this.type = type;
    }

    public MockField(String address, PlcValue plcValue) {
        this.address = address;
        this.type = plcValue.getPlcValueType();
    }

    @Override
    public String getAddressString() {
        return address;
    }

    @Override
    public PlcValueType getPlcValueType() {
        return type;
    }

    @Override
    public List<ArrayInfo> getArrayInfo() {
        return Collections.emptyList();
    }

    @Override
    public String toString() {
        return "mock field: " + address;
    }

    @Override
    public boolean equals(Object o) {
        return o instanceof MockField
            && ((MockField) o).address.equals(this.address);
    }

    @Override
    public int hashCode() {
        return address.hashCode();
    }

}
