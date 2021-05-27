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
package org.apache.plc4x.java.mock.field;

import org.apache.plc4x.java.api.exceptions.PlcInvalidFieldException;
import org.apache.plc4x.java.api.model.PlcField;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class MockField implements PlcField {

    private final String address;
    private MockPlcValue plcValue;
    private MockType type;

    private static final Pattern PATTERN =
        Pattern.compile("%(?<name>[a-zA-Z_.0-9]+(?:\\[[0-9]*])?):?(?<type>[A-Z]*)");

    public static MockField of(String addressString) throws PlcInvalidFieldException {
        Matcher matcher = PATTERN.matcher(addressString);
        if (matcher.matches()) {
            String addr = matcher.group("name");
            MockType type = MockType.valueOf(matcher.group("type"));
            return new MockField(addr, type);
        }
        return null;
    }

    public MockField(String address) {
        this.address = address;
        this.plcValue = null;
    }

    public MockField(String address, MockType type) {
        this.address = address;
        this.type = type;
    }

    public MockField(String address, MockPlcValue plcValue) {
        this.address = address;
        this.plcValue = plcValue;
    }

    public String getAddress() {
        return address;
    }

    public MockPlcValue getPlcValue() {
        return plcValue;
    }

    @Override
    public String getPlcDataType() {
        return type.toString();
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


    @Override
    public Class<?> getDefaultJavaType() {
        switch (type) {
            case BOOL:
                return Boolean.class;
            case INT:
                return Integer.class;
            case REAL:
                return Double.class;
            default:
                return null;
        }
    }
}
